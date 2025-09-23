package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"

	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/mocks"
)

func TestProducer_Send_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWriter := mocks.NewMockBrokerWriter(ctrl)

	mockWriter.
		EXPECT().
		WriteMessages(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, msgs ...kafka.Message) error {
			// optional: add checks: len(msgs), etc
			return nil
		})

	p := &Producer{
		writer:    mockWriter,
		brokers:   []string{"127.0.0.1:9092"},
		topic:     "test-topic",
		retries:   1,
		backoff:   func(int) time.Duration { return 0 },
		marshaler: json.Marshal,
		keyFunc:   func(e entity.DomainEvent) []byte { return []byte("k") },
		logger:    logger.NewDefault(),
	}

	ctx := context.Background()
	event := entity.BannerImpressionRecorded{
		BannerID: 1,
		SlotID:   2,
		GroupID:  3,
		Time:     time.Now(),
	}

	if err := p.Send(ctx, event); err != nil {
		t.Fatal("expected nil error, gow %w", err)
	}
}

func TestProducer_Send_RetryFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWriter := mocks.NewMockBrokerWriter(ctrl)

	mockWriter.
		EXPECT().
		WriteMessages(gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("boom")).
		AnyTimes()

	p := &Producer{
		writer:    mockWriter,
		brokers:   []string{"127.0.0.1:9092"},
		topic:     "test-topic",
		retries:   2,
		backoff:   func(int) time.Duration { return 0 },
		marshaler: json.Marshal,
		keyFunc:   func(e entity.DomainEvent) []byte { return []byte("k") },
		logger:    logger.NewDefault(),
	}

	ctx := context.Background()
	event := entity.BannerImpressionRecorded{
		BannerID: 1,
		SlotID:   2,
		GroupID:  3,
		Time:     time.Now(),
	}

	if err := p.Send(ctx, event); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
