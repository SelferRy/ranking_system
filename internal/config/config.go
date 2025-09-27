package config

import (
	broker "github.com/SelferRy/ranking_system/internal/infra/adapters/broker/kafka"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository"
	server "github.com/SelferRy/ranking_system/internal/infra/adapters/server/grpc"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
)

type Config struct {
	Logger   logger.Config     `mapstructure:"logger"`
	Database repository.Config `mapstructure:"database"`
	Server   server.Config     `mapstructure:"server"`
	Broker   broker.Config     `mapstructure:"broker"`
}
