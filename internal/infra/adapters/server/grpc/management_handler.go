// Package grpc: gRPC handler
package grpc

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/api/gen"
	"github.com/SelferRy/ranking_system/internal/application/usecase"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BannerManagementHandler struct {
	gen.UnimplementedBannerManagementServiceServer
	managementUC *usecase.ManagementUseCase
}

func NewBannerManagementHandler(managementUC *usecase.ManagementUseCase) *BannerManagementHandler {
	return &BannerManagementHandler{
		managementUC: managementUC,
	}
}

func (h *BannerManagementHandler) AddBannerToSlot(
	ctx context.Context,
	req *gen.AddBannerToSlotRequest,
) (*gen.AddBannerToSlotResponse, error) {
	if req.BannerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "banner_id is required")
	}
	if req.SlotId == 0 {
		return nil, status.Error(codes.InvalidArgument, "slot_id is required")
	}

	err := h.managementUC.AddBannerToSlot(
		ctx,
		entity.BannerID(req.BannerId),
		entity.SlotID(req.SlotId),
	)
	if err != nil {
		return &gen.AddBannerToSlotResponse{
			Success: false,
			Error:   fmt.Sprintf("add banner to slot: %v", err),
		}, status.Error(codes.Internal, err.Error())
	}

	return &gen.AddBannerToSlotResponse{
		Success: true,
	}, nil
}

func (h *BannerManagementHandler) RemoveBannerFromSlot(
	ctx context.Context,
	req *gen.RemoveBannerFromSlotRequest,
) (*gen.RemoveBannerFromSlotResponse, error) {
	if req.BannerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "banner_id is required")
	}
	if req.SlotId == 0 {
		return nil, status.Error(codes.InvalidArgument, "slot_id is required")
	}

	err := h.managementUC.RemoveBannerFromSlot(
		ctx,
		entity.BannerID(req.BannerId),
		entity.SlotID(req.SlotId),
	)
	if err != nil {
		return &gen.RemoveBannerFromSlotResponse{
			Success: false,
			Error:   fmt.Sprintf("remove banner from slot: %v", err),
		}, status.Error(codes.Internal, err.Error())
	}

	return &gen.RemoveBannerFromSlotResponse{
		Success: true,
	}, nil
}
