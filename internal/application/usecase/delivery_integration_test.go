//go:build integration
// +build integration

package usecase_test

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/application/usecase"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/SelferRy/ranking_system/internal/domain/interfaces/broker"
	"github.com/SelferRy/ranking_system/internal/domain/service/bandit"
	internalkafka "github.com/SelferRy/ranking_system/internal/infra/adapters/broker/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository/postgres"
	"github.com/SelferRy/ranking_system/internal/infra/logger"
)

func setupTestDB(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()
	err := godotenv.Load("../../../../.env")
	require.NoError(t, err)

	dsn := os.Getenv("GOOSE_DBSTRING")
	require.NotEmpty(t, dsn, "GOOSE_DBSTRING must be set in .env")

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	require.NoError(t, pool.Ping(ctx))

	return pool
}

// prepareData inserts fixture data directly through pool.Exec
func prepareData(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
	t.Helper()

	queries := []struct {
		sql  string
		args []interface{}
		log  string
	}{
		{
			sql:  `INSERT INTO ranking_system.banners (id, description) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`,
			args: []interface{}{1, "test banner"},
			log:  "banners",
		},
		{
			sql:  `INSERT INTO ranking_system.slots (id, description) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`,
			args: []interface{}{1, "test slot"},
			log:  "slots",
		},
		{
			sql:  `INSERT INTO ranking_system.banner_slot (banner_id, slot_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			args: []interface{}{1, 1},
			log:  "banner_slot",
		},
		{
			sql:  `INSERT INTO ranking_system.groups (id, description) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`,
			args: []interface{}{1, "test group"},
			log:  "groups",
		},
		{
			sql:  `INSERT INTO ranking_system.banner_stats (banner_id, slot_id, group_id, impressions, clicks) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (banner_id,slot_id,group_id) DO NOTHING`,
			args: []interface{}{1, 1, 1, 0, 0},
			log:  "banner_stats",
		},
	}

	for _, q := range queries {
		ct, err := pool.Exec(ctx, q.sql, q.args...)
		require.NoError(t, err)
		t.Logf("%s. Inserted %d row(s)", q.log, ct.RowsAffected())
	}
}

func setupTestBroker(t *testing.T, ctx context.Context, log logger.Logger) broker.EventProducer {
	t.Helper()

	brokerAddr := "localhost:9092"
	topic := "test-integration-topic"

	// Create a topic if it does not already exist
	conn, err := kafka.DialContext(ctx, "tcp", brokerAddr)
	require.NoError(t, err)
	defer conn.Close()

	controller, err := conn.Controller()
	require.NoError(t, err)

	controllerAddr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, err := kafka.DialContext(ctx, "tcp", controllerAddr)
	require.NoError(t, err)
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	_ = controllerConn.CreateTopics(topicConfig)

	producer := internalkafka.NewProducer(
		[]string{brokerAddr},
		topic,
		internalkafka.WithLogger(log),
	)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	require.NoError(t, producer.HealthCheck(ctxTimeout))

	return producer
}

func TestDeliveryUseCase_IntegrationSelectBanner_HappyPath(t *testing.T) {
	ctx := context.Background()
	log := logger.NewDefault()

	pool := setupTestDB(t, ctx)
	defer pool.Close()

	// Loading fixtures
	prepareData(t, pool, ctx)

	bannerRepo := postgres.NewBannerRepository(pool)
	statsRepo := postgres.NewStatsRepository(pool)

	selector := bandit.NewUCB1Service()

	producer := setupTestBroker(t, ctx, log)
	defer producer.Close()

	uc := usecase.NewDeliveryUseCase(
		log,
		bannerRepo,
		statsRepo,
		selector,
		producer,
	)

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)
	groupID := entity.GroupID(1)

	selectedBanner, err := uc.SelectBanner(ctx, slotID, groupID)
	require.NoError(t, err)
	assert.Equal(t, bannerID, selectedBanner.ID)
}
