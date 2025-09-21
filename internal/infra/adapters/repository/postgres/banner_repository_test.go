package postgres_test

import (
	"context"
	"fmt"
	"github.com/SelferRy/ranking_system/internal/infra/adapters/repository/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dsn := "postgres://user:password@localhost:5432/ranking_test?sslmode=disable&search_path=ranking_system"

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	return pool
}

func setupTestTx(t *testing.T, pool *pgxpool.Pool) pgx.Tx {
	tx, err := pool.Begin(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(context.Background())
	})
	return tx
}

func TestBannerRepository(t *testing.T) {
	pool := setupTestDB(t)

	tx := setupTestTx(t, pool)
	repo := postgres.NewBannerRepository(tx)

	ctx := context.Background()

	// insert in banners table
	commandTag, err := tx.Exec(
		ctx,
		`INSERT INTO ranking_system.banners (id, description) VALUES (1, 'test banner') RETURNING *`,
	)
	t.Logf("banners. Inserted %d row(s)", commandTag.RowsAffected())
	require.NoError(t, err)

	// insert in slots table
	commandTag, err = tx.Exec(
		ctx,
		`INSERT INTO ranking_system.slots (id, description) VALUES (1, 'test slot') RETURNING *`,
	)
	t.Logf("slots. Inserted %d row(s)", commandTag.RowsAffected())
	require.NoError(t, err)

	// insert in banner_slot link-table
	commandTag, err = tx.Exec(
		ctx,
		`INSERT INTO ranking_system.banner_slot (banner_id, slot_id) VALUES (1, 1) RETURNING *`,
	)
	t.Logf("banner_slot. Inserted %d row(s)", commandTag.RowsAffected())
	require.NoError(t, err)

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
