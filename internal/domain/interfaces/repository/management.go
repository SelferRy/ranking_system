package repository

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

// ManagementRepository defines methods for managing banners within slots.
type ManagementRepository interface {
	AddBannerToSlot(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID) error
	BannerExistsInSlot(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID) (bool, error)
	RemoveBannerFromSlot(ctx context.Context, bannerID entity.BannerID, slotID entity.SlotID) error
}
