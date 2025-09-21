package postgres

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/interfaces/repository"
)

type statsRepository struct {
	pool Pool // abstraction over *pgxpool.Pool for tests
}

func NewStatsRepository(pool Pool) repository.StatsRepository {
	return &statsRepository{pool: pool}
}

func (r *statsRepository) RecordImpression(
	ctx context.Context,
	bannerID entity.BannerID,
	slotID entity.SlotID,
	groupID entity.GroupID,
) error {
	const query = `
		INSERT INTO ranking_system.banner_stats (banner_id, slot_id, group_id, impressions, clicks)
		VALUES ($1, $2, $3, 1, 0)
		ON CONFLICT (banner_id, slot_id, group_id)
    		DO UPDATE SET impressions = banner_stats.impressions + 1
	`

	_, err := r.pool.Exec(ctx, query, bannerID, slotID, groupID)
	if err != nil {
		return fmt.Errorf("record impression: %w", err)
	}
	return nil
}

func (r *statsRepository) GetBannerStats(
	ctx context.Context,
	bannerID entity.BannerID,
	slotID entity.SlotID,
	groupID entity.GroupID,
) (entity.BannerStat, error) {
	const query = `
		SELECT banner_id, slot_id, group_id, impressions, clicks 
		FROM ranking_system.banner_stats
		WHERE banner_id = $1
		AND slot_id = $2
		AND group_id = $3
		LIMIT 1
`
	var stat entity.BannerStat
	row := r.pool.QueryRow(ctx, query, bannerID, slotID, groupID)
	if err := row.Scan(&stat.BannerID, &stat.SlotID, &stat.GroupID, &stat.Impressions, &stat.Clicks); err != nil {
		return entity.BannerStat{}, fmt.Errorf("get banner stats: %w", err)
	}
	return stat, nil
}
