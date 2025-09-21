package postgres

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/interfaces/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type bannerRepository struct {
	//pool *pgxpool.Pool
	pool DB
}

func NewBannerRepository(pool DB) repository.BannerRepository {
	return &bannerRepository{pool: pool}
}

func (r *bannerRepository) RequestBanner(ctx context.Context, slotID entity.SlotID) ([]entity.Banner, error) {
	const query = `
		SELECT b.id, b.description
		FROM ranking_system.banners b
		JOIN ranking_system.banner_slot bs ON b.id = bs.banner_id
		WHERE bs.slot_id = $1
	`

	rows, err := r.pool.Query(ctx, query, slotID)
	if err != nil {
		return nil, fmt.Errorf("query banners: %w", err)
	}
	defer rows.Close()

	var banners []entity.Banner
	for rows.Next() {
		var banner entity.Banner
		if err := rows.Scan(&banner.ID, &banner.Description); err != nil {
			return nil, fmt.Errorf("scan banner: %w", err)
		}
		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return banners, nil
}
