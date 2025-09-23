//go:build integration
// +build integration

package kafka_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"testing"
	"time"

	internalkafka "github.com/SelferRy/ranking_system/internal/infra/adapters/broker/kafka"
	kafka "github.com/segmentio/kafka-go"
)

func TestProducer_Integration_SendAndConsume(t *testing.T) {
	broker := "localhost:9092"
	topic := "test-integration-topic"

	p := internalkafka.NewProducer([]string{broker}, topic)
	defer p.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := p.HealthCheck(ctx)
	if err != nil {
		t.Fatalf("health check: %v", err)
	}

	conn, _ := kafka.DialContext(ctx, "tcp", broker) // err skipped because already check via HealthCheck
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		t.Fatalf("controller: %v", err)
	}

	controllerAddr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, err := kafka.DialContext(ctx, "tcp", controllerAddr)
	if err != nil {
		t.Fatalf("dial controller: %v", err)
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	_ = controllerConn.CreateTopics(topicConfig)

	// produce one event
	event := entity.BannerImpressionRecorded{
		BannerID: 123,
		SlotID:   7,
		GroupID:  11,
		Time:     time.Now(),
	}

	if err := p.Send(ctx, event); err != nil {
		t.Fatalf("producer send failed: %v", err)
	}

	// consume the event
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})
	defer reader.Close()

	// read the msg
	readCtx, readCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer readCancel()

	msg, err := reader.ReadMessage(readCtx)
	if err != nil {
		t.Fatalf("read message failed: %v", err)
	}

	var got entity.BannerImpressionRecorded
	if err := json.Unmarshal(msg.Value, &got); err != nil {
		t.Fatalf("unmarshal message failed: %v", err)
	}

	if got.BannerID != event.BannerID || got.SlotID != event.SlotID {
		t.Fatalf("unexpected payload: got %+v, want %+v", got, event)
	}

	t.Logf("integration test success: got message key=%s, value=%s", string(msg.Key), string(msg.Value))
}
