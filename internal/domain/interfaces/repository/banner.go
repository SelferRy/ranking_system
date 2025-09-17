package repository

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

//go:generate mockgen -source=banner.go -destination=../../../mocks/mock_banner2_repo.go -package=mocks

type BannerRepository interface {
	GetForSlot(ctx context.Context, slotID entity.SlotID) ([]entity.Banner, error)
	AddToSlot(ctx context.Context, slotID entity.SlotID, bannerID entity.BannerID) error
	RemoveFromSlot(ctx context.Context, slotID entity.SlotID, bannerID entity.BannerID) error
	ExistsInSlot(ctx context.Context, slotID entity.SlotID, bannerID entity.BannerID) (bool, error)
}
