package mock

import (
	"context"

	"github.com/creatachain/augusteum/libs/bytes"
	msm "github.com/creatachain/augusteum/msm/types"
	"github.com/creatachain/augusteum/proxy"
	"github.com/creatachain/augusteum/rpc/client"
	ctypes "github.com/creatachain/augusteum/rpc/core/types"
	"github.com/creatachain/augusteum/types"
)

// MSMApp will send all msm related request to the named app,
// so you can test app behavior from a client without needing
// an entire augusteum node
type MSMApp struct {
	App msm.Application
}

var (
	_ client.MSMClient = MSMApp{}
	_ client.MSMClient = MSMMock{}
	_ client.MSMClient = (*MSMRecorder)(nil)
)

func (a MSMApp) MSMInfo(ctx context.Context) (*ctypes.ResultMSMInfo, error) {
	return &ctypes.ResultMSMInfo{Response: a.App.Info(proxy.RequestInfo)}, nil
}

func (a MSMApp) MSMQuery(ctx context.Context, path string, data bytes.HexBytes) (*ctypes.ResultMSMQuery, error) {
	return a.MSMQueryWithOptions(ctx, path, data, client.DefaultMSMQueryOptions)
}

func (a MSMApp) MSMQueryWithOptions(
	ctx context.Context,
	path string,
	data bytes.HexBytes,
	opts client.MSMQueryOptions) (*ctypes.ResultMSMQuery, error) {
	q := a.App.Query(msm.RequestQuery{
		Data:   data,
		Path:   path,
		Height: opts.Height,
		Prove:  opts.Prove,
	})
	return &ctypes.ResultMSMQuery{Response: q}, nil
}

// NOTE: Caller should call a.App.Commit() separately,
// this function does not actually wait for a commit.
// TODO: Make it wait for a commit and set res.Height appropriately.
func (a MSMApp) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	res := ctypes.ResultBroadcastTxCommit{}
	res.CheckTx = a.App.CheckTx(msm.RequestCheckTx{Tx: tx})
	if res.CheckTx.IsErr() {
		return &res, nil
	}
	res.DeliverTx = a.App.DeliverTx(msm.RequestDeliverTx{Tx: tx})
	res.Height = -1 // TODO
	return &res, nil
}

func (a MSMApp) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	c := a.App.CheckTx(msm.RequestCheckTx{Tx: tx})
	// and this gets written in a background thread...
	if !c.IsErr() {
		go func() { a.App.DeliverTx(msm.RequestDeliverTx{Tx: tx}) }()
	}
	return &ctypes.ResultBroadcastTx{
		Code:      c.Code,
		Data:      c.Data,
		Log:       c.Log,
		Codespace: c.Codespace,
		Hash:      tx.Hash(),
	}, nil
}

func (a MSMApp) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	c := a.App.CheckTx(msm.RequestCheckTx{Tx: tx})
	// and this gets written in a background thread...
	if !c.IsErr() {
		go func() { a.App.DeliverTx(msm.RequestDeliverTx{Tx: tx}) }()
	}
	return &ctypes.ResultBroadcastTx{
		Code:      c.Code,
		Data:      c.Data,
		Log:       c.Log,
		Codespace: c.Codespace,
		Hash:      tx.Hash(),
	}, nil
}

// MSMMock will send all msm related request to the named app,
// so you can test app behavior from a client without needing
// an entire augusteum node
type MSMMock struct {
	Info            Call
	Query           Call
	BroadcastCommit Call
	Broadcast       Call
}

func (m MSMMock) MSMInfo(ctx context.Context) (*ctypes.ResultMSMInfo, error) {
	res, err := m.Info.GetResponse(nil)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultMSMInfo{Response: res.(msm.ResponseInfo)}, nil
}

func (m MSMMock) MSMQuery(ctx context.Context, path string, data bytes.HexBytes) (*ctypes.ResultMSMQuery, error) {
	return m.MSMQueryWithOptions(ctx, path, data, client.DefaultMSMQueryOptions)
}

func (m MSMMock) MSMQueryWithOptions(
	ctx context.Context,
	path string,
	data bytes.HexBytes,
	opts client.MSMQueryOptions) (*ctypes.ResultMSMQuery, error) {
	res, err := m.Query.GetResponse(QueryArgs{path, data, opts.Height, opts.Prove})
	if err != nil {
		return nil, err
	}
	resQuery := res.(msm.ResponseQuery)
	return &ctypes.ResultMSMQuery{Response: resQuery}, nil
}

func (m MSMMock) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	res, err := m.BroadcastCommit.GetResponse(tx)
	if err != nil {
		return nil, err
	}
	return res.(*ctypes.ResultBroadcastTxCommit), nil
}

func (m MSMMock) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := m.Broadcast.GetResponse(tx)
	if err != nil {
		return nil, err
	}
	return res.(*ctypes.ResultBroadcastTx), nil
}

func (m MSMMock) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := m.Broadcast.GetResponse(tx)
	if err != nil {
		return nil, err
	}
	return res.(*ctypes.ResultBroadcastTx), nil
}

// MSMRecorder can wrap another type (MSMApp, MSMMock, or Client)
// and record all MSM related calls.
type MSMRecorder struct {
	Client client.MSMClient
	Calls  []Call
}

func NewMSMRecorder(client client.MSMClient) *MSMRecorder {
	return &MSMRecorder{
		Client: client,
		Calls:  []Call{},
	}
}

type QueryArgs struct {
	Path   string
	Data   bytes.HexBytes
	Height int64
	Prove  bool
}

func (r *MSMRecorder) addCall(call Call) {
	r.Calls = append(r.Calls, call)
}

func (r *MSMRecorder) MSMInfo(ctx context.Context) (*ctypes.ResultMSMInfo, error) {
	res, err := r.Client.MSMInfo(ctx)
	r.addCall(Call{
		Name:     "msm_info",
		Response: res,
		Error:    err,
	})
	return res, err
}

func (r *MSMRecorder) MSMQuery(
	ctx context.Context,
	path string,
	data bytes.HexBytes,
) (*ctypes.ResultMSMQuery, error) {
	return r.MSMQueryWithOptions(ctx, path, data, client.DefaultMSMQueryOptions)
}

func (r *MSMRecorder) MSMQueryWithOptions(
	ctx context.Context,
	path string,
	data bytes.HexBytes,
	opts client.MSMQueryOptions) (*ctypes.ResultMSMQuery, error) {
	res, err := r.Client.MSMQueryWithOptions(ctx, path, data, opts)
	r.addCall(Call{
		Name:     "msm_query",
		Args:     QueryArgs{path, data, opts.Height, opts.Prove},
		Response: res,
		Error:    err,
	})
	return res, err
}

func (r *MSMRecorder) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	res, err := r.Client.BroadcastTxCommit(ctx, tx)
	r.addCall(Call{
		Name:     "broadcast_tx_commit",
		Args:     tx,
		Response: res,
		Error:    err,
	})
	return res, err
}

func (r *MSMRecorder) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := r.Client.BroadcastTxAsync(ctx, tx)
	r.addCall(Call{
		Name:     "broadcast_tx_async",
		Args:     tx,
		Response: res,
		Error:    err,
	})
	return res, err
}

func (r *MSMRecorder) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := r.Client.BroadcastTxSync(ctx, tx)
	r.addCall(Call{
		Name:     "broadcast_tx_sync",
		Args:     tx,
		Response: res,
		Error:    err,
	})
	return res, err
}
