package core

import (
	"github.com/creatachain/augusteum/libs/bytes"
	msm "github.com/creatachain/augusteum/msm/types"
	"github.com/creatachain/augusteum/proxy"
	ctypes "github.com/creatachain/augusteum/rpc/core/types"
	rpctypes "github.com/creatachain/augusteum/rpc/jsonrpc/types"
)

// MSMQuery queries the application for some information.
// More: https://docs.augusteum.com/master/rpc/#/MSM/msm_query
func MSMQuery(
	ctx *rpctypes.Context,
	path string,
	data bytes.HexBytes,
	height int64,
	prove bool,
) (*ctypes.ResultMSMQuery, error) {
	resQuery, err := env.ProxyAppQuery.QuerySync(msm.RequestQuery{
		Path:   path,
		Data:   data,
		Height: height,
		Prove:  prove,
	})
	if err != nil {
		return nil, err
	}
	env.Logger.Info("MSMQuery", "path", path, "data", data, "result", resQuery)
	return &ctypes.ResultMSMQuery{Response: *resQuery}, nil
}

// MSMInfo gets some info about the application.
// More: https://docs.augusteum.com/master/rpc/#/MSM/msm_info
func MSMInfo(ctx *rpctypes.Context) (*ctypes.ResultMSMInfo, error) {
	resInfo, err := env.ProxyAppQuery.InfoSync(proxy.RequestInfo)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultMSMInfo{Response: *resInfo}, nil
}
