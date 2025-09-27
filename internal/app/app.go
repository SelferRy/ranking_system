package app

import (
	"context"
	"fmt"
	ucbanner "github.com/SelferRy/ranking_system/internal/application/usecase"
	"github.com/SelferRy/ranking_system/internal/config"
	"github.com/SelferRy/ranking_system/internal/domain/service/bandit"
	internalkafka "github.com/SelferRy/ranking_system/internal/infra/adapters/broker/kafka"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository/postgres"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/server"
	grpcserver "github.com/SelferRy/ranking_system/internal/infra/adapters/server/grpc"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
)

type App struct {
	grpcServer server.Server
	logger     logger.Logger
	dbPool     postgres.Pool
	producer   *internalkafka.Producer
}

func New(ctx context.Context, conf config.Config, logger logger.Logger) (*App, error) {
	// TODO: fill it like in nimoism imp (maybe return api (server) and then listen with new goroutine <-ctx.Done(),etc?
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dbPool, err := postgres.NewPool(ctx, conf.Database)
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %w", err)
	}

	// repo init
	bannerRepo := postgres.NewBannerRepository(dbPool)
	statsRepo := postgres.NewStatsRepository(dbPool)

	// init broker
	producer := internalkafka.NewProducerFromConfig(conf.Broker, logger)

	// init bandit service
	selector := bandit.NewUCB1Service()

	deliveryUC := ucbanner.NewDeliveryUseCase(
		logger,
		bannerRepo,
		statsRepo,
		selector,
		producer,
	)

	grpcServer, err := grpcserver.NewServer(
		grpcserver.Config{
			Host: conf.Server.Host,
			Port: conf.Server.Port,
		},
		logger,
		grpcserver.RegisterServices(deliveryUC),
	)
	if err != nil {
		producer.Close()
		dbPool.Close()
		return nil, fmt.Errorf("grpc server initialization failed: %w", err)
	}

	//return ConcreteServer{ctx, conf, logger}
	return &App{
		grpcServer: grpcServer,
		logger:     logger,
		dbPool:     dbPool,
		producer:   producer,
	}, nil
}

func (a App) Run(ctx context.Context) error {
	a.logger.Info("Starting application...")

	// Start gRPC server in a goroutine
	go func() {
		a.logger.Info("Starting gRPC server", logger.String("address", a.grpcServer.GetAddress()))
		if err := a.grpcServer.Start(); err != nil {
			a.logger.Error("gRPC server failed", logger.Error("grpcServer.Start error: ", err))
		}
	}()

	<-ctx.Done()
	a.logger.Info("Shutdown signal received")

	// Graceful shutdown
	a.grpcServer.Stop(ctx)
	a.producer.Close()
	a.dbPool.Close()

	a.logger.Info("Application stopped gracefully")
	return nil
}
