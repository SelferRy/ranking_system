package config

import (
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository"
	"github.com/SelferRy/ranking_system/internal/logger"
	server "github.com/SelferRy/ranking_system/internal/server/grpc"
)

type Config struct {
	Logger   logger.Config     `mapstructure:"logger"`
	Database repository.Config `mapstructure:"database"`
	Server   server.Config     `mapstructure:"server"`
	//App    app.Config    `mapstructure:"app"`
	//Test   string        `yaml:"test"`
	//One    any           `yaml:"one"`
}
