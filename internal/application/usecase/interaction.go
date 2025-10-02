package usecase

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"time"

	"github.com/SelferRy/ranking_system/internal/domain/interfaces/broker"
	"github.com/SelferRy/ranking_system/internal/domain/interfaces/repository"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
)

type InteractionUseCase struct {
	log            logger.Logger
	managementRepo repository.ManagementRepository
	statsRepo      repository.StatsRepository
	producer       broker.EventProducer
}

func NewInteractionUseCase(
	log logger.Logger,
	managementRepo repository.ManagementRepository,
	statsRepo repository.StatsRepository,
	producer broker.EventProducer,
) *InteractionUseCase {
	return &InteractionUseCase{
		log:            log,
		managementRepo: managementRepo,
		statsRepo:      statsRepo,
		producer:       producer,
	}
}

func (uc *InteractionUseCase) RegisterClick(
	ctx context.Context,
	bannerID entity.BannerID,
	slotID entity.SlotID,
	groupID entity.GroupID,
) error {
	// 1. Check that the banner exists in the slot
	exists, err := uc.managementRepo.BannerExistsInSlot(ctx, bannerID, slotID)
	if err != nil {
		return fmt.Errorf("check banner existence: %w", err)
	}
	if !exists {
		uc.log.Info("banner does not exists in the slot",
			logger.Int64("banner_id", int64(bannerID)),
			logger.Int64("slot_id", int64(slotID)),
		)
		return nil
	}

	// 2. Write to stats
	if err = uc.statsRepo.RecordClick(ctx, bannerID, slotID, groupID); err != nil {
		return fmt.Errorf("record click: %w", err)
	}

	ev := entity.BannerClickRecorded{
		BannerID: bannerID,
		SlotID:   slotID,
		GroupID:  groupID,
		Time:     time.Now(),
	}
	if err = uc.producer.Send(ctx, ev); err != nil {
		return fmt.Errorf("send event: %w", err)
	}

	return nil
}
