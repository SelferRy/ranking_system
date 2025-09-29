package postgres

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, conf repository.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	pool, err := pgxpool.New(ctx, connString)
	return pool, err
}
