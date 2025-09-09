package broker

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

type EventProducer interface {
	SendEvent(ctx context.Context, event entity.RotationEvent) error
}
