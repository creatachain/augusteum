package core

import (
	ctypes "github.com/creatachain/augusteum/rpc/core/types"
	rpctypes "github.com/creatachain/augusteum/rpc/jsonrpc/types"
)

// Health gets node health. Returns empty result (200 OK) on success, no
// response - in case of an error.
// More: https://docs.augusteum.com/master/rpc/#/Info/health
func Health(ctx *rpctypes.Context) (*ctypes.ResultHealth, error) {
	return &ctypes.ResultHealth{}, nil
}
