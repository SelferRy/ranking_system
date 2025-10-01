package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"github.com/SelferRy/ranking_system/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagementUseCase_AddBannerToSlot_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockManagementRepository(ctrl)
	uc := NewManagementUseCase(mockRepo)
	uc.logger = logger.NewDefault()

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)

	// expectations
	mockRepo.EXPECT().
		BannerExistsInSlot(ctx, bannerID, slotID).
		Return(false, nil)

	mockRepo.EXPECT().
		AddBannerToSlot(ctx, bannerID, slotID).
		Return(nil)

	// act
	err := uc.AddBannerToSlot(ctx, bannerID, slotID)

	// assert
	require.NoError(t, err)
}

func TestManagementUseCase_AddBannerToSlot_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockManagementRepository(ctrl)
	uc := NewManagementUseCase(mockRepo)
	uc.logger = logger.NewDefault()

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)

	// expectations
	mockRepo.EXPECT().
		BannerExistsInSlot(ctx, bannerID, slotID).
		Return(true, nil)

	// act
	err := uc.AddBannerToSlot(ctx, bannerID, slotID)

	// assert
	require.NoError(t, err)
}

func TestManagementUseCase_AddBannerToSlot_ErrorOnExistsCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockManagementRepository(ctrl)
	uc := NewManagementUseCase(mockRepo)
	uc.logger = logger.NewDefault()

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)

	// expectations
	mockRepo.EXPECT().
		BannerExistsInSlot(ctx, bannerID, slotID).
		Return(false, errors.New("db error"))

	// act
	err := uc.AddBannerToSlot(ctx, bannerID, slotID)

	// assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "check banner existence")
}

func TestManagementUseCase_AddBannerToSlot_ErrorOnInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockManagementRepository(ctrl)
	uc := NewManagementUseCase(mockRepo)
	uc.logger = logger.NewDefault()

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)

	// expectations
	mockRepo.EXPECT().
		BannerExistsInSlot(ctx, bannerID, slotID).
		Return(false, nil)

	mockRepo.EXPECT().
		AddBannerToSlot(ctx, bannerID, slotID).
		Return(errors.New("insert failed"))

	// act
	err := uc.AddBannerToSlot(ctx, bannerID, slotID)

	// assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "add banner to slot")
}

func TestManagementUseCase_RemoveBannerFromSlot_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockManagementRepository(ctrl)
	uc := NewManagementUseCase(mockRepo)
	uc.logger = logger.NewDefault()

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)

	// expectations
	mockRepo.EXPECT().
		RemoveBannerFromSlot(ctx, bannerID, slotID).
		Return(nil)

	// act
	err := uc.RemoveBannerFromSlot(ctx, bannerID, slotID)

	// assert
	require.NoError(t, err)
}

func TestManagementUseCase_RemoveBannerFromSlot_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockManagementRepository(ctrl)
	uc := NewManagementUseCase(mockRepo)
	uc.logger = logger.NewDefault()

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)

	// expectations
	mockRepo.EXPECT().
		RemoveBannerFromSlot(ctx, bannerID, slotID).
		Return(errors.New("delete failed"))

	// act
	err := uc.RemoveBannerFromSlot(ctx, bannerID, slotID)

	// assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "remove banner from slot")
}
