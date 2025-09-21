package repository

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

//go:generate mockgen -source=stats.go -destination=../../../mocks/mock_stats_repo.go -package=mocks

type StatsRepository interface {
	RecordImpression(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID, groupID entity.GroupID) error
	GetBannerStats(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID, groupID entity.GroupID) (entity.BannerStat, error)
	//IncrementClick(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID, groupID entity.GroupID) error
	//GetSlotBannersStats(ctx context.Context, slotID entity.SlotID, groupID entity.GroupID) ([]entity.BannerStat, error)
}
