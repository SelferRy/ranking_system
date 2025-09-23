package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
	"github.com/segmentio/kafka-go"
	"time"
)

// broker interface for tests and mocks
type brokerWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Producer - Kafka-backend implementation of infra broker.EventProducer.
type Producer struct {
	writer    brokerWriter
	brokers   []string
	topic     string
	retries   int
	backoff   func(attempt int) time.Duration
	marshaler func(any) ([]byte, error)
	keyFunc   func(entity.DomainEvent) []byte
	logger    logger.Logger
	validator func(entity.DomainEvent) error
}

// Option allows customizing the producer.
type Option func(*Producer)

// NewProducer creates a new Kafka producer (segmentio/kafka-go).
// brokers: list of "host:port"
func NewProducer(brokers []string, topic string, opts ...Option) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}

	p := &Producer{
		writer:    writer,
		brokers:   brokers, //append([]string(nil), brokers...),
		topic:     topic,
		retries:   3,
		backoff:   defaultBackoff,
		marshaler: json.Marshal,
		keyFunc:   defaultKeyFunc,
		logger:    logger.NewDefault(),
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithRetries(n int) Option {
	return func(p *Producer) { p.retries = n }
}

func WithBackoff(fn func(int) time.Duration) Option {
	return func(p *Producer) { p.backoff = fn }
}

func WithMarshaler(m func(any) ([]byte, error)) Option {
	return func(p *Producer) { p.marshaler = m }
}

func WithLogger(log logger.Logger) Option {
	return func(p *Producer) { p.logger = log }
}

func WithKeyFunc(keyFunc func(entity.DomainEvent) []byte) Option {
	return func(p *Producer) { p.keyFunc = keyFunc }
}

func WithValidator(validator func(entity.DomainEvent) error) Option {
	return func(p *Producer) { p.validator = validator }
}

func (p *Producer) Send(ctx context.Context, event entity.DomainEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if p.validator != nil {
		if err := p.validator(event); err != nil {
			return fmt.Errorf("event validation failed: %w", err)
		}
	}

	data, err := p.marshaler(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	key := p.keyFunc(event)
	eventType := event.EventType()

	msg := kafka.Message{
		Key:     key,
		Value:   data,
		Time:    time.Now(),
		Headers: []kafka.Header{{Key: "event_type", Value: []byte(eventType)}},
	}

	var lastErr error
	for attempt := 1; attempt <= p.retries; attempt++ {
		if ctx.Err() != nil {
			return fmt.Errorf("context canceled: %w", ctx.Err())
		}

		lastErr = p.writer.WriteMessages(ctx, msg)
		if lastErr == nil {
			// TODO: обернуть в значения zap.Field + добавить интерфейс в logger
			p.logger.Debug("kafka: message sent",
				logger.StringVal("topic", p.topic),
				logger.StringVal("event_type", eventType),
				logger.StringVal("key", string(key)),
				logger.IntVal("attempt", attempt),
			)
			return nil
		}

		p.logger.Warn("kafka: send failed, will retry",
			logger.StringVal("topic", p.topic),
			logger.IntVal("attempt", attempt),
			logger.ErrorVal("error", lastErr),
		)

		time.Sleep(p.backoff(attempt))
	}

	p.logger.Error("kafka: send failed after retries",
		logger.StringVal("topic", p.topic),
		logger.IntVal("attempts", p.retries),
		logger.ErrorVal("error", lastErr),
	)

	return fmt.Errorf("failed to send message after %d attempts: %w", p.retries, lastErr)
}

// HealthCheck performs a lightweight connectivity check against the cluster.
// It dials the first broker (discovery) and returns any connection error.
// This is suitable for a simple readiness probe. If you need topic-specific checks,
// we can extend this to DialLeader(topic, partition) or use Client.CreateTopics.
func (p *Producer) HealthCheck(ctx context.Context) error {
	if len(p.brokers) == 0 {
		return fmt.Errorf("no brokers configured")
	}

	// Dial the first broker (discovery). kafka.DialContext will fetch metadata.
	conn, err := kafka.DialContext(ctx, "tcp", p.brokers[0])
	if err != nil {
		return fmt.Errorf("kafka dial failed: %w", err)
	}
	_ = conn.Close()
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func defaultKeyFunc(event entity.DomainEvent) []byte {
	return event.Key()
}

func defaultBackoff(attempt int) time.Duration {
	return time.Duration(attempt*attempt) * 100 * time.Millisecond
}

// TODO:
//  +1. Logger values
//   2. Mock tests
//   3. Kafka conteiner UP
//	 4. Integration tests
//   5. Makefile setup for kafka tests
//   6. Скорректировать API под usecase
//	 7. Скорректировать main/root - запуск процесса, инициализации и пр. стартующую инфру
//	 8. CI/CD: Dockerfile + общая сборка в контейнерах
//	 9. Какой-то тест end-to-end
