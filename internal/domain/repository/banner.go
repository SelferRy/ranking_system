package repository

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

type BannerRepository interface {
	GetForSlot(ctx context.Context, slotID entity.SlotID) ([]entity.Banner, error)
	AddToSlot(ctx context.Context, slotID entity.SlotID, bannerID entity.BannerID) error
	RemoveFromSlot(ctx context.Context, slotID entity.SlotID, bannerID entity.BannerID) error
	ExistsInSlot(ctx context.Context, slotID entity.SlotID, bannerID entity.BannerID) (bool, error)
}
