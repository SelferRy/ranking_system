// Package grpc: service registration
package grpc

import (
	"github.com/SelferRy/ranking_system/api/gen"
	"github.com/SelferRy/ranking_system/internal/application/usecase"
	"google.golang.org/grpc"
)

func RegisterServices(uc *usecase.DeliveryUseCase) func(*grpc.Server) {
	return func(s *grpc.Server) {
		// registration of BannerRotatorService
		handler := NewBannerRotatorHandler(uc)
		gen.RegisterBannerRotatorServiceServer(s, handler)
	}
}
