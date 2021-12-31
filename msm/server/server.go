/*
Package server is used to start a new MSM server.

It contains two server implementation:
 * gRPC server
 * socket server

*/
package server

import (
	"fmt"

	"github.com/creatachain/augusteum/libs/service"
	"github.com/creatachain/augusteum/msm/types"
)

func NewServer(protoAddr, transport string, app types.Application) (service.Service, error) {
	var s service.Service
	var err error
	switch transport {
	case "socket":
		s = NewSocketServer(protoAddr, app)
	case "grpc":
		s = NewGRPCServer(protoAddr, types.NewGRPCApplication(app))
	default:
		err = fmt.Errorf("unknown server type %s", transport)
	}
	return s, err
}