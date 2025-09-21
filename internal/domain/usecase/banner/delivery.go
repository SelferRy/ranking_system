package banner

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/interfaces/repository"
	"github.com/SelferRy/ranking_system/internal/domain/service/bandit"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/broker"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"time"
)

type DeliveryUseCase struct {
	log        logger.Logger
	bannerRepo repository.BannerRepository
	statsRepo  repository.StatsRepository
	selector   bandit.BannerSelector
	producer   broker.EventProducer
}

func NewDeliveryUseCase(
	log logger.Logger,
	bannerRepo repository.BannerRepository,
	statsRepo repository.StatsRepository,
	selector bandit.BannerSelector,
	producer broker.EventProducer,
) *DeliveryUseCase {
	return &DeliveryUseCase{
		log:        log,
		bannerRepo: bannerRepo,
		statsRepo:  statsRepo,
		selector:   selector,
		producer:   producer,
	}
}

func (uc *DeliveryUseCase) SelectBanner(
	ctx context.Context,
	slotID entity.SlotID,
	groupID entity.GroupID,
) (entity.Banner, error) {
	// 1. Get all banners for the slot
	banners, err := uc.bannerRepo.RequestBanner(ctx, slotID)
	if err != nil {
		return entity.Banner{}, fmt.Errorf("get banners: %w", err)
	}
	if len(banners) == 0 {
		return entity.Banner{}, fmt.Errorf("no banners for slot %d", slotID)
	}

	// 2. Stats
	stats, err := uc.collectStats(ctx, banners, slotID, groupID)
	if err != nil {
		return entity.Banner{}, fmt.Errorf("collect stats: %w", err)
	}

	// 3. Select
	selectedID, err := uc.selector.SelectBanner(stats)
	if err != nil {
		return entity.Banner{}, fmt.Errorf("select banner: %w", err)
	}

	// 4.Upd stats
	if err := uc.statsRepo.RecordImpression(ctx, selectedID, slotID, groupID); err != nil {
		return entity.Banner{}, fmt.Errorf("record impression: %w", err)
	}

	// 5.Event public
	ev := entity.BannerImpressionRecorded{
		BannerID: selectedID,
		SlotID:   slotID,
		GroupID:  groupID,
		Time:     time.Now(),
	}
	if err := uc.producer.Send(ctx, ev); err != nil {
		return entity.Banner{}, fmt.Errorf("send event: %w", err)
	}

	// 6. Check the banner
	banner, ok := uc.findBanner(banners, selectedID)
	if !ok {
		return entity.Banner{}, fmt.Errorf("selected banner not found")
	}

	return banner, nil
}

func (uc *DeliveryUseCase) collectStats(
	ctx context.Context,
	banners []entity.Banner,
	slotID entity.SlotID,
	groupID entity.GroupID,
) ([]entity.BannerStat, error) {
	stats := make([]entity.BannerStat, 0, len(banners))
	for _, b := range banners {
		stat, err := uc.statsRepo.GetBannerStats(ctx, b.ID, slotID, groupID)
		if err != nil {
			return nil, fmt.Errorf("get stats for banner %d: %w", b.ID, err)
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

func (uc *DeliveryUseCase) findBanner(banners []entity.Banner, id entity.BannerID) (entity.Banner, bool) {
	for _, b := range banners {
		if b.ID == id {
			return b, true
		}
	}
	return entity.Banner{}, false
}
