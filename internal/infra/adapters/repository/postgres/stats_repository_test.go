//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"github.com/SelferRy/ranking_system/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository/postgres"
)

func TestStatsRepository__RecordAndGet__Integration(t *testing.T) {
	pool := setupTestDB(t)

	tx := setupTestTx(t, pool)
	repo := postgres.NewStatsRepository(tx)

	ctx := context.Background()

	bannerID := entity.BannerID(1)
	slotID := entity.SlotID(1)
	groupID := entity.GroupID(1)

	err := prepareData(t, tx, ctx)

	err = repo.RecordImpression(ctx, bannerID, slotID, groupID)
	require.NoError(t, err)

	stat, err := repo.GetBannerStats(ctx, bannerID, slotID, groupID)
	require.NoError(t, err)
	require.Equal(t, int64(1), stat.Impressions)
	require.Equal(t, int64(0), stat.Clicks)
}
