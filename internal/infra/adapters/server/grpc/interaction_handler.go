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

type BannerInteractionHandler struct {
	gen.UnimplementedBannerInteractionServiceServer
	interactionUC *usecase.InteractionUseCase
}

func NewBannerInteractionHandler(interactionUC *usecase.InteractionUseCase) *BannerInteractionHandler {
	return &BannerInteractionHandler{
		interactionUC: interactionUC,
	}
}

func (s *BannerInteractionHandler) RegisterClick(
	ctx context.Context,
	req *gen.ClickBannerRequest,
) (*gen.ClickBannerResponse, error) {
	err := s.interactionUC.RegisterClick(
		ctx,
		entity.BannerID(req.BannerId),
		entity.SlotID(req.SlotId),
		entity.GroupID(req.GroupId),
	)
	if err != nil {
		return &gen.ClickBannerResponse{
			Success: false,
			Error:   fmt.Sprintf("click banner: %v", err),
		}, status.Error(codes.Internal, fmt.Sprintf("register click: %v", err))
	}

	return &gen.ClickBannerResponse{
		Success: true,
	}, nil
}
