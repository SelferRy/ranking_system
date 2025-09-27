package server

import "context"

// Server интерфейс для запуска/остановки сервера
type Server interface {
	Start() error
	Stop(ctx context.Context) error
	GetAddress() string
}
