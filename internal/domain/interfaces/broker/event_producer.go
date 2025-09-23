package broker

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
)

type EventProducer interface {
	Send(ctx context.Context, event entity.DomainEvent) error
}
