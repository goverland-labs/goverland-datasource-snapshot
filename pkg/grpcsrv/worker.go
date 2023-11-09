package grpcsrv

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type GrpcServerWorker struct {
	name       string
	grpcServer *grpc.Server
	bind       string
}

func NewGrpcServerWorker(name string, grpcServer *grpc.Server, bind string) *GrpcServerWorker {
	return &GrpcServerWorker{name: name, grpcServer: grpcServer, bind: bind}
}

func (g *GrpcServerWorker) Start() error {
	log.Info().Fields(map[string]interface{}{
		"bind": g.bind,
		"name": g.name,
	}).Msg("start grpc server worker")

	return listenAndServe(g.grpcServer, g.bind)
}

func listenAndServe(s *grpc.Server, bind string) error {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return fmt.Errorf("gRPC listen: %v", err)
	}

	log.Info().Msg("gRPC server started")
	defer log.Info().Msg("gRPC server exited")
	if err := s.Serve(listener); err != nil {
		return fmt.Errorf("gRPC serve: %v", err)
	}

	return nil
}

func (g *GrpcServerWorker) Stop() error {
	g.grpcServer.GracefulStop()
	return nil
}
