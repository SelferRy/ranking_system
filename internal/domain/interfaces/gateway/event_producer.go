package gateway

import "github.com/SelferRy/ranking_system/internal/domain/entity"

type EventProducer interface {
	Send(event entity.DomainEvent) error
}
