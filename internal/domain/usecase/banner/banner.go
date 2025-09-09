package banner

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/repository"
	"github.com/SelferRy/ranking_system/internal/domain/service/bandit"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/broker"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"time"
)

type UseCase struct {
	log        logger.Logger
	bannerRepo repository.BannerRepository
	statsRepo  repository.StatsRepository
	bandit     bandit.BannerSelector
	producer   broker.EventProducer
}

func New(
	log logger.Logger,
	bannerRepo repository.BannerRepository,
	statsRepo repository.StatsRepository,
	banditCalc bandit.BannerSelector,
	producer broker.EventProducer,
) *UseCase {
	return &UseCase{
		log:        log,
		bannerRepo: bannerRepo,
		statsRepo:  statsRepo,
		bandit:     banditCalc,
		producer:   producer,
	}
}

func (uc UseCase) AddBannerToSlot(ctx context.Context, slotID, bannerID int64) error {
	exists, err := uc.bannerRepo.ExistsInSlot(ctx, slotID, bannerID)
	if err != nil {
		return fmt.Errorf("check banner existence: %w", err)
	}
	if exists {
		uc.log.Info(
			"banner already exists in the slot",
			logger.Int64("slotID", slotID),
			logger.Int64("bannerID", bannerID),
		)
		return nil
	}
	if err := uc.bannerRepo.AddToSlot(ctx, slotID, bannerID); err != nil {
		return fmt.Errorf("add banner to slot: %w", err)
	}
	uc.log.Info(
		"banner added to slot",
		logger.Int64("slotID", slotID),
		logger.Int64("bannerID", bannerID),
	)
	return nil
}

func (uc UseCase) RemoveBannerFromSlot(ctx context.Context, slotID, bannerID int64) error {
	if err := uc.bannerRepo.RemoveFromSlot(ctx, slotID, bannerID); err != nil {
		return fmt.Errorf("remove banner from slot: %w", err)
	}
	uc.log.Info(
		"banner removed from slot",
		logger.Int64("slotID", slotID),
		logger.Int64("bannerID", bannerID),
	)
	return nil
}

func (uc UseCase) RegisterClick(ctx context.Context, slotID, bannerID, groupID int64) error {
	if err := uc.statsRepo.IncrementClick(ctx, slotID, bannerID, groupID); err != nil {
		return fmt.Errorf("increment banner click: %w", err)
	}
	event := entity.RotationEvent{
		Type:      "click",
		BannerID:  bannerID,
		SlotID:    slotID,
		GroupID:   groupID,
		EventTime: time.Now().UTC(),
	}
	if err := uc.producer.SendEvent(ctx, event); err != nil {
		uc.log.Error("failed to send click event", logger.Error(err))
	}
	return nil
}

func (uc UseCase) SelectBanner() error {
	panic("implement me")
}
