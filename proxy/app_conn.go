package proxy

import (
	msmcli "github.com/creatachain/augusteum/msm/client"
	"github.com/creatachain/augusteum/msm/types"
)

//go:generate mockery --case underscore --name AppConnConsensus|AppConnMempool|AppConnQuery|AppConnSnapshot

//----------------------------------------------------------------------------------------
// Enforce which msm msgs can be sent on a connection at the type level

type AppConnConsensus interface {
	SetResponseCallback(msmcli.Callback)
	Error() error

	InitChainSync(types.RequestInitChain) (*types.ResponseInitChain, error)

	BeginBlockSync(types.RequestBeginBlock) (*types.ResponseBeginBlock, error)
	DeliverTxAsync(types.RequestDeliverTx) *msmcli.ReqRes
	EndBlockSync(types.RequestEndBlock) (*types.ResponseEndBlock, error)
	CommitSync() (*types.ResponseCommit, error)
}

type AppConnMempool interface {
	SetResponseCallback(msmcli.Callback)
	Error() error

	CheckTxAsync(types.RequestCheckTx) *msmcli.ReqRes
	CheckTxSync(types.RequestCheckTx) (*types.ResponseCheckTx, error)

	FlushAsync() *msmcli.ReqRes
	FlushSync() error
}

type AppConnQuery interface {
	Error() error

	EchoSync(string) (*types.ResponseEcho, error)
	InfoSync(types.RequestInfo) (*types.ResponseInfo, error)
	QuerySync(types.RequestQuery) (*types.ResponseQuery, error)

	//	SetOptionSync(key string, value string) (res types.Result)
}

type AppConnSnapshot interface {
	Error() error

	ListSnapshotsSync(types.RequestListSnapshots) (*types.ResponseListSnapshots, error)
	OfferSnapshotSync(types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error)
	LoadSnapshotChunkSync(types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error)
	ApplySnapshotChunkSync(types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error)
}

//-----------------------------------------------------------------------------------------
// Implements AppConnConsensus (subset of msmcli.Client)

type appConnConsensus struct {
	appConn msmcli.Client
}

func NewAppConnConsensus(appConn msmcli.Client) AppConnConsensus {
	return &appConnConsensus{
		appConn: appConn,
	}
}

func (app *appConnConsensus) SetResponseCallback(cb msmcli.Callback) {
	app.appConn.SetResponseCallback(cb)
}

func (app *appConnConsensus) Error() error {
	return app.appConn.Error()
}

func (app *appConnConsensus) InitChainSync(req types.RequestInitChain) (*types.ResponseInitChain, error) {
	return app.appConn.InitChainSync(req)
}

func (app *appConnConsensus) BeginBlockSync(req types.RequestBeginBlock) (*types.ResponseBeginBlock, error) {
	return app.appConn.BeginBlockSync(req)
}

func (app *appConnConsensus) DeliverTxAsync(req types.RequestDeliverTx) *msmcli.ReqRes {
	return app.appConn.DeliverTxAsync(req)
}

func (app *appConnConsensus) EndBlockSync(req types.RequestEndBlock) (*types.ResponseEndBlock, error) {
	return app.appConn.EndBlockSync(req)
}

func (app *appConnConsensus) CommitSync() (*types.ResponseCommit, error) {
	return app.appConn.CommitSync()
}

//------------------------------------------------
// Implements AppConnMempool (subset of msmcli.Client)

type appConnMempool struct {
	appConn msmcli.Client
}

func NewAppConnMempool(appConn msmcli.Client) AppConnMempool {
	return &appConnMempool{
		appConn: appConn,
	}
}

func (app *appConnMempool) SetResponseCallback(cb msmcli.Callback) {
	app.appConn.SetResponseCallback(cb)
}

func (app *appConnMempool) Error() error {
	return app.appConn.Error()
}

func (app *appConnMempool) FlushAsync() *msmcli.ReqRes {
	return app.appConn.FlushAsync()
}

func (app *appConnMempool) FlushSync() error {
	return app.appConn.FlushSync()
}

func (app *appConnMempool) CheckTxAsync(req types.RequestCheckTx) *msmcli.ReqRes {
	return app.appConn.CheckTxAsync(req)
}

func (app *appConnMempool) CheckTxSync(req types.RequestCheckTx) (*types.ResponseCheckTx, error) {
	return app.appConn.CheckTxSync(req)
}

//------------------------------------------------
// Implements AppConnQuery (subset of msmcli.Client)

type appConnQuery struct {
	appConn msmcli.Client
}

func NewAppConnQuery(appConn msmcli.Client) AppConnQuery {
	return &appConnQuery{
		appConn: appConn,
	}
}

func (app *appConnQuery) Error() error {
	return app.appConn.Error()
}

func (app *appConnQuery) EchoSync(msg string) (*types.ResponseEcho, error) {
	return app.appConn.EchoSync(msg)
}

func (app *appConnQuery) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	return app.appConn.InfoSync(req)
}

func (app *appConnQuery) QuerySync(reqQuery types.RequestQuery) (*types.ResponseQuery, error) {
	return app.appConn.QuerySync(reqQuery)
}

//------------------------------------------------
// Implements AppConnSnapshot (subset of msmcli.Client)

type appConnSnapshot struct {
	appConn msmcli.Client
}

func NewAppConnSnapshot(appConn msmcli.Client) AppConnSnapshot {
	return &appConnSnapshot{
		appConn: appConn,
	}
}

func (app *appConnSnapshot) Error() error {
	return app.appConn.Error()
}

func (app *appConnSnapshot) ListSnapshotsSync(req types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	return app.appConn.ListSnapshotsSync(req)
}

func (app *appConnSnapshot) OfferSnapshotSync(req types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	return app.appConn.OfferSnapshotSync(req)
}

func (app *appConnSnapshot) LoadSnapshotChunkSync(
	req types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	return app.appConn.LoadSnapshotChunkSync(req)
}

func (app *appConnSnapshot) ApplySnapshotChunkSync(
	req types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	return app.appConn.ApplySnapshotChunkSync(req)
}
