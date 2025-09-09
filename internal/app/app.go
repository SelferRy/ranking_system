package app

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/config"
	ucbanner "github.com/SelferRy/ranking_system/internal/domain/usecase/banner"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/broker/kafka"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository/sql/pg"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"github.com/SelferRy/ranking_system/internal/server"
	"github.com/SelferRy/ranking_system/internal/server/grpc"
	bannerservice "github.com/SelferRy/ranking_system/internal/server/grpc/handler/banner"
)

type App struct {
	grpcServer server.Server
	logger     logger.Logger
}

func New(ctx context.Context, conf config.Config, logger logger.Logger) (*App, error) {
	// TODO: fill it like in nimoism imp (maybe return api (server) and then listen with new goroutine <-ctx.Done(),etc?
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	db, err := pg.New(logger, conf.Database)
	if err != nil {
		return nil, fmt.Errorf("repo initialization problem: %w", err)
	}

	bannerRepo, err := panic()

	broker, err := kafka.New(logger)
	if err != nil {
		return nil, fmt.Errorf("message broker gateway initialization problem: %w", err)
	}

	useCase, err := ucbanner.New(logger, repo, broker)
	if err != nil {
		return nil, fmt.Errorf("use case initialization problem: %w", err)
	}

	service, err := bannerservice.New(logger, useCase)
	if err != nil {
		return nil, fmt.Errorf("service initialization problem: %w", err)
	}

	grpcServer, err := grpc.New(logger, service)
	if err != nil {
		return nil, fmt.Errorf("server initialization problem: %w", err)
	}

	//return ConcreteServer{ctx, conf, logger}
	return &App{
		grpcServer: grpcServer,
		logger:     logger,
	}, nil
}
