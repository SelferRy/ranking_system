package repository

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

type StatsRepository interface {
	IncrementShow(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID, groupID entity.GroupID) error
	IncrementClick(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID, groupID entity.GroupID) error
	GetBannerStats(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID, groupID entity.GroupID) (*entity.BannerStat, error)
	GetSlotBannersStats(ctx context.Context, slotID entity.SlotID, groupID entity.GroupID) ([]entity.BannerStat, error)
}
