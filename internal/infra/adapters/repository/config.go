package repository

import "time"

type Config struct {
	Type        string        `json:"type,omitempty"`
	Driver      string        `json:"driver,omitempty"`
	Host        string        `json:"host,omitempty"`
	Port        string        `json:"port,omitempty"`
	Database    string        `json:"database,omitempty"`
	User        string        `json:"user,omitempty"`
	Password    string        `json:"password,omitempty"`
	Ssl         string        `json:"ssl,omitempty"`
	MaxConn     int           `json:"max_conn,omitempty"`
	MaxIdleConn int           `json:"max_idle_conn,omitempty"`
	MaxLifetime time.Duration `json:"max_lifetime,omitempty"`
}
