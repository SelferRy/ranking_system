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

type BannerRotatorService struct {
	gen.UnimplementedBannerRotatorServiceServer
	deliveryUC *usecase.DeliveryUseCase
}

func NewBannerRotatorService(deliveryUC *usecase.DeliveryUseCase) *BannerRotatorService {
	return &BannerRotatorService{
		deliveryUC: deliveryUC,
	}
}

func (s *BannerRotatorService) SelectBanner(
	ctx context.Context,
	req *gen.SelectBannerRequest,
) (*gen.SelectBannerResponse, error) {
	if req.SlotId == 0 {
		return nil, status.Error(codes.InvalidArgument, "slot_id is required")
	}
	if req.GroupId == 0 {
		return nil, status.Error(codes.InvalidArgument, "group_id is required")
	}

	banner, err := s.deliveryUC.SelectBanner(
		ctx,
		entity.SlotID(req.SlotId),
		entity.GroupID(req.GroupId),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("select banner: %v", err))
	}

	return &gen.SelectBannerResponse{
		Banner: &gen.Banner{
			Id:          uint64(banner.ID),
			Description: banner.Description,
		},
	}, nil
}
