package repository

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

type StatsRepository interface {
	IncrementShow(ctx context.Context, bannerID, slotID, groupID int64) error
	IncrementClick(ctx context.Context, bannerID, slotID, groupID int64) error
	GetBannerStats(ctx context.Context, bannerID, slotID, groupID int64) (*entity.BannerStat, error)
	GetSlotBannersStats(ctx context.Context, slotID, groupID int64) ([]entity.BannerStat, error)
}
