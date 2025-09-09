package server

import (
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"github.com/SelferRy/ranking_system/internal/server/grpc/handler"
)

//	type App struct {
//		logger logger.Logger
//		uc     repository.EventUseCase
//	}
//
//	func New(logger logger.Logger, uc repository.EventUseCase) *App {
//		return &App{
//			logger: logger,
//			uc:     uc,
//		}
//	}
//
//	func GetEventUseCase(conf repository.Config) repository.EventUseCase {
//		panic("Not implemented")
//	}
type Server interface {
	Start() error
	Stop() error
}

type ConcreteServer struct {
	//ctx context.Context
	//conf config.Config
	//logger logger.Logger
}

func New(logger logger.Logger, services ...service.Service) (Server, error) {
	return ConcreteServer{}, nil // fixme
}

func (cs ConcreteServer) Start() error {
	panic("Not implemented")
}

func (cs ConcreteServer) Stop() error {
	panic("Not implemented")
}
