package grpc

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer struct {
	server  *grpc.Server
	lis     net.Listener
	logger  logger.Logger
	address string
}

func NewServer(cfg Config, logger logger.Logger, registerServices func(*grpc.Server)) (*GrpcServer, error) {
	// create grPC server
	grpcServer := grpc.NewServer()

	// registered services
	registerServices(grpcServer)

	// create listener
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return &GrpcServer{
		server:  grpcServer,
		lis:     lis,
		logger:  logger,
		address: address,
	}, nil
}

func (s GrpcServer) Start() error {
	s.logger.Info("Starting gRPC server", logger.String("address", s.address))

	if err := s.server.Serve(s.lis); err != nil {
		return fmt.Errorf("gRPC server failed to serve: %w", err)
	}
	return nil
}

func (s GrpcServer) Stop(ctx context.Context) error {
	s.logger.Info("Stopping gRPC server")

	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("Force stopping gRPC server")
		s.server.Stop()
		return ctx.Err()
	}
}

func (s *GrpcServer) GetAddress() string {
	return s.address
}
