package repository

import "time"

type Config struct {
	Type        string        `mapstructure:"type,omitempty"`
	Driver      string        `mapstructure:"driver,omitempty"`
	Host        string        `mapstructure:"host,omitempty"`
	Port        string        `mapstructure:"port,omitempty"`
	Database    string        `mapstructure:"database,omitempty"`
	User        string        `mapstructure:"user,omitempty"`
	Password    string        `mapstructure:"password,omitempty"`
	Ssl         string        `mapstructure:"ssl,omitempty"`
	MaxConn     int           `mapstructure:"max_conn,omitempty"`
	MaxIdleConn int           `mapstructure:"max_idle_conn,omitempty"`
	MaxLifetime time.Duration `mapstructure:"max_lifetime,omitempty"`
}
