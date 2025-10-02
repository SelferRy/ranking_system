// Package grpc: service registration
package grpc

import (
	"github.com/SelferRy/ranking_system/api/gen"
	"github.com/SelferRy/ranking_system/internal/application/usecase"
	"google.golang.org/grpc"
)

func RegisterServices(
	ucDelivery *usecase.DeliveryUseCase,
	ucManagement *usecase.ManagementUseCase,
	ucInteraction *usecase.InteractionUseCase,
) func(*grpc.Server) {
	return func(s *grpc.Server) {
		// registration of BannerRotatorService
		handlerBannerRotator := NewBannerRotatorHandler(ucDelivery)
		gen.RegisterBannerRotatorServiceServer(s, handlerBannerRotator)

		// registration of BannerManagementService
		handlerBannerManager := NewBannerManagementHandler(ucManagement)
		gen.RegisterBannerManagementServiceServer(s, handlerBannerManager)

		// registration of BannerInteractionService
		handlerBannerInteraction := NewBannerInteractionHandler(ucInteraction)
		gen.RegisterBannerInteractionServiceServer(s, handlerBannerInteraction)
	}
}
