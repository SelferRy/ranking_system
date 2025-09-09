package internalgrpc

import (
	router "github.com/SelferRy/ranking_system/internal/api"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"google.golang.org/grpc"
)

type Config struct {
	Type string `json:"type,omitempty"`
	Host string `json:"host,omitempty"`
	Port string `json:"port,omitempty"`
}

type Server interface {
	Start() error
	Stop() error
}

type Application interface {
}

type GrpcServer struct {
	Address string
	logger  logger.Logger
	server  *grpc.Server
}

func NewServer(logger logger.Logger, app router.Application, conf Config) *Server {
	panic("Not implemented")
	//return GrpcServer{}
}

func (gs GrpcServer) Start() error {
	gs.logger.Info("gRPC server started.")
	return nil // fixme
}

func (gs GrpcServer) Stop() error {
	gs.logger.Info("gRPC server has started to stop.")
	return nil // fixme
}
