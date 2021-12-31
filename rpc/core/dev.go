package core

import (
	ctypes "github.com/creatachain/augusteum/rpc/core/types"
	rpctypes "github.com/creatachain/augusteum/rpc/jsonrpc/types"
)

// UnsafeFlushMempool removes all transactions from the mempool.
func UnsafeFlushMempool(ctx *rpctypes.Context) (*ctypes.ResultUnsafeFlushMempool, error) {
	env.Mempool.Flush()
	return &ctypes.ResultUnsafeFlushMempool{}, nil
}
