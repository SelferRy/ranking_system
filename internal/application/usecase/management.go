package usecase

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/interfaces/repository"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
)

type ManagementUseCase struct {
	log  logger.Logger
	repo repository.ManagementRepository
}

func NewManagementUseCase(logger logger.Logger, repo repository.ManagementRepository) *ManagementUseCase {
	return &ManagementUseCase{log: logger, repo: repo}
}

func (uc *ManagementUseCase) AddBannerToSlot(ctx context.Context, banner entity.Banner, slotID entity.SlotID) error {
	exists, err := uc.repo.BannerExistsInSlot(ctx, banner.ID, slotID)
	if err != nil {
		return fmt.Errorf("check banner existence: %w", err)
	}
	if exists {
		uc.log.Info("banner already exists",
			logger.Int64("banner_id", int64(banner.ID)),
			logger.Int64("slot_id", int64(slotID)),
		)
		return nil
	}

	if err := uc.repo.AddBannerToSlot(ctx, banner, slotID); err != nil {
		return fmt.Errorf("add banner to slot: %w", err)
	}
	uc.log.Info(
		"banner added to slot",
		logger.Int64("bannerID", int64(banner.ID)),
		logger.Int64("slotID", int64(slotID)),
	)
	return nil
}

func (uc *ManagementUseCase) RemoveBannerFromSlot(ctx context.Context, banner entity.Banner, slotID entity.SlotID) error {
	if err := uc.repo.RemoveBannerFromSlot(ctx, banner, slotID); err != nil {
		return fmt.Errorf("remove banner from slot: %w", err)
	}
	uc.log.Info(
		"banner removed from slot",
		logger.Int64("bannerID", int64(banner.ID)),
		logger.Int64("slotID", int64(slotID)),
	)
	return nil
}
