package config

import "github.com/SelferRy/ranking_system/internal/logger"

type Config struct {
	Logger logger.Config `mapstructure:"log"`
	Test   string        `yaml:"test"`
	One    any           `yaml:"one"`
}
