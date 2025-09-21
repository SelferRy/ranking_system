package postgres_test

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository/postgres"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBannerRepository(t *testing.T) {
	pool := setupTestDB(t)

	tx := setupTestTx(t, pool)
	repo := postgres.NewBannerRepository(tx)

	ctx := context.Background()

	err := prepareData(t, tx, ctx)

	// Run SELECT for check inserted data
	rows, err := tx.Query(ctx, "SELECT id, description FROM ranking_system.banners")
	require.NoError(t, err)
	defer rows.Close()

	// Print results in stdout
	fmt.Println("Данные в таблице banners:")
	for rows.Next() {
		var id int64
		var description string
		err := rows.Scan(&id, &description)
		require.NoError(t, err)
		fmt.Printf("ID: %d, Description: %s\n", id, description)
	}
	require.NoError(t, rows.Err())

	// test repo strategy
	banners, err := repo.RequestBanner(ctx, 1)
	require.NoError(t, err)
	require.Len(t, banners, 1)
	require.Equal(t, "test banner", banners[0].Description)
}
