package usecase

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/service/bandit"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"github.com/SelferRy/ranking_system/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeliveryUseCase__SelectBanner__HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// deps
	log := logger.NewDefault()
	bannerRepo := mocks.NewMockBannerRepository(ctrl)
	statsRepo := mocks.NewMockStatsRepository(ctrl)
	selector := bandit.NewUCB1Service() //mocks.NewMockBannerSelector(ctrl)
	producer := mocks.NewMockEventProducer(ctrl)

	uc := NewDeliveryUseCase(
		log,
		bannerRepo,
		statsRepo,
		selector,
		producer,
	)

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)
	groupID := entity.GroupID(1)

	// expectations
	bannerRepo.EXPECT().
		RequestBanner(ctx, slotID).
		Return([]entity.Banner{{ID: bannerID, Description: "test"}}, nil)

	statsRepo.EXPECT().
		GetBannerStats(ctx, bannerID, slotID, groupID).
		Return(entity.BannerStat{BannerID: bannerID, Impressions: 1}, nil)

	statsRepo.EXPECT().
		RecordImpression(ctx, bannerID, slotID, groupID).
		Return(nil)

	producer.EXPECT().
		Send(ctx, gomock.AssignableToTypeOf(entity.BannerImpressionRecorded{})).
		Return(nil)

	// act
	b, err := uc.SelectBanner(ctx, slotID, groupID)

	// assert
	require.NoError(t, err)
	assert.Equal(t, bannerID, b.ID)
}

func TestDeliveryUseCase__SelectBanner__HappyPathSelectorMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// deps
	log := logger.NewDefault()
	bannerRepo := mocks.NewMockBannerRepository(ctrl)
	statsRepo := mocks.NewMockStatsRepository(ctrl)
	selector := mocks.NewMockBannerSelector(ctrl)
	producer := mocks.NewMockEventProducer(ctrl)

	uc := NewDeliveryUseCase(
		log,
		bannerRepo,
		statsRepo,
		selector,
		producer,
	)

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)
	groupID := entity.GroupID(1)

	// expectations
	bannerRepo.EXPECT().
		RequestBanner(ctx, slotID).
		Return([]entity.Banner{{ID: bannerID, Description: "test"}}, nil)

	statsRepo.EXPECT().
		GetBannerStats(ctx, bannerID, slotID, groupID).
		Return(entity.BannerStat{BannerID: bannerID, Impressions: 1}, nil)

	selector.EXPECT().
		SelectBanner(gomock.Any()).
		Return(bannerID, nil)

	statsRepo.EXPECT().
		RecordImpression(ctx, bannerID, slotID, groupID).
		Return(nil)

	producer.EXPECT().
		Send(ctx, gomock.AssignableToTypeOf(entity.BannerImpressionRecorded{})).
		Return(nil)

	// act
	b, err := uc.SelectBanner(ctx, slotID, groupID)

	// assert
	require.NoError(t, err)
	assert.Equal(t, bannerID, b.ID)
}
