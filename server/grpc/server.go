package grpc

import (
	"context"
	"net"

	hrpc "github.com/seb7887/heimdallr/proto/heimdallr"
	"github.com/seb7887/heimdallr/service"
	"google.golang.org/grpc"
)

type GRPCServer interface {
	Serve(ctx context.Context) error
}

type grpcServer struct {
	grpcAddr string
	grpcService service.Service
}

func New(addr string, service service.Service) GRPCServer {
	return &grpcServer{
		grpcAddr: addr,
		grpcService: service,
	}
}

func (s *grpcServer) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	serviceServer := NewHeimdallrGRPCServer(s.grpcService)

	hrpc.RegisterHeimdallrServiceServer(grpcServer, serviceServer)

	if err := grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}
