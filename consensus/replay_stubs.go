package consensus

import (
	"github.com/creatachain/augusteum/libs/clist"
	mempl "github.com/creatachain/augusteum/mempool"
	msm "github.com/creatachain/augusteum/msm/types"
	tmstate "github.com/creatachain/augusteum/proto/augusteum/state"
	"github.com/creatachain/augusteum/proxy"
	"github.com/creatachain/augusteum/types"
)

//-----------------------------------------------------------------------------

type emptyMempool struct{}

var _ mempl.Mempool = emptyMempool{}

func (emptyMempool) Lock()     {}
func (emptyMempool) Unlock()   {}
func (emptyMempool) Size() int { return 0 }
func (emptyMempool) CheckTx(_ types.Tx, _ func(*msm.Response), _ mempl.TxInfo) error {
	return nil
}
func (emptyMempool) ReapMaxBytesMaxGas(_, _ int64) types.Txs { return types.Txs{} }
func (emptyMempool) ReapMaxTxs(n int) types.Txs              { return types.Txs{} }
func (emptyMempool) Update(
	_ int64,
	_ types.Txs,
	_ []*msm.ResponseDeliverTx,
	_ mempl.PreCheckFunc,
	_ mempl.PostCheckFunc,
) error {
	return nil
}
func (emptyMempool) Flush()                        {}
func (emptyMempool) FlushAppConn() error           { return nil }
func (emptyMempool) TxsAvailable() <-chan struct{} { return make(chan struct{}) }
func (emptyMempool) EnableTxsAvailable()           {}
func (emptyMempool) TxsBytes() int64               { return 0 }

func (emptyMempool) TxsFront() *clist.CElement    { return nil }
func (emptyMempool) TxsWaitChan() <-chan struct{} { return nil }

func (emptyMempool) InitWAL() error { return nil }
func (emptyMempool) CloseWAL()      {}

//-----------------------------------------------------------------------------
// mockProxyApp uses MSMResponses to give the right results.
//
// Useful because we don't want to call Commit() twice for the same block on
// the real app.

func newMockProxyApp(appHash []byte, msmResponses *tmstate.MSMResponses) proxy.AppConnConsensus {
	clientCreator := proxy.NewLocalClientCreator(&mockProxyApp{
		appHash:      appHash,
		msmResponses: msmResponses,
	})
	cli, _ := clientCreator.NewMSMClient()
	err := cli.Start()
	if err != nil {
		panic(err)
	}
	return proxy.NewAppConnConsensus(cli)
}

type mockProxyApp struct {
	msm.BaseApplication

	appHash      []byte
	txCount      int
	msmResponses *tmstate.MSMResponses
}

func (mock *mockProxyApp) DeliverTx(req msm.RequestDeliverTx) msm.ResponseDeliverTx {
	r := mock.msmResponses.DeliverTxs[mock.txCount]
	mock.txCount++
	if r == nil {
		return msm.ResponseDeliverTx{}
	}
	return *r
}

func (mock *mockProxyApp) EndBlock(req msm.RequestEndBlock) msm.ResponseEndBlock {
	mock.txCount = 0
	return *mock.msmResponses.EndBlock
}

func (mock *mockProxyApp) Commit() msm.ResponseCommit {
	return msm.ResponseCommit{Data: mock.appHash}
}
