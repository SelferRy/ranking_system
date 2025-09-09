package repository

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

type BannerRepository interface {
	GetForSlot(ctx context.Context, slotID int64) ([]entity.Banner, error)
	AddToSlot(ctx context.Context, slotID, bannerID int64) error
	RemoveFromSlot(ctx context.Context, slotID, bannerID int64) error
	ExistsInSlot(ctx context.Context, slotID, bannerID int64) (bool, error)
}
