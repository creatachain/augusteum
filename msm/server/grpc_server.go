package server

import (
	"net"

	"google.golang.org/grpc"

	tmnet "github.com/creatachain/augusteum/libs/net"
	"github.com/creatachain/augusteum/libs/service"
	"github.com/creatachain/augusteum/msm/types"
)

type GRPCServer struct {
	service.BaseService

	proto    string
	addr     string
	listener net.Listener
	server   *grpc.Server

	app types.MSMApplicationServer
}

// NewGRPCServer returns a new gRPC MSM server
func NewGRPCServer(protoAddr string, app types.MSMApplicationServer) service.Service {
	proto, addr := tmnet.ProtocolAndAddress(protoAddr)
	s := &GRPCServer{
		proto:    proto,
		addr:     addr,
		listener: nil,
		app:      app,
	}
	s.BaseService = *service.NewBaseService(nil, "MSMServer", s)
	return s
}

// OnStart starts the gRPC service.
func (s *GRPCServer) OnStart() error {

	ln, err := net.Listen(s.proto, s.addr)
	if err != nil {
		return err
	}

	s.listener = ln
	s.server = grpc.NewServer()
	types.RegisterMSMApplicationServer(s.server, s.app)

	s.Logger.Info("Listening", "proto", s.proto, "addr", s.addr)
	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			s.Logger.Error("Error serving gRPC server", "err", err)
		}
	}()
	return nil
}

// OnStop stops the gRPC server.
func (s *GRPCServer) OnStop() {
	s.server.Stop()
}
