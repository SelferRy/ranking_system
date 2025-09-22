package postgres_test

import (
	"context"
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

func prepareData(t *testing.T, tx pgx.Tx, ctx context.Context) error {
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

	// insert in groups table
	commandTag, err = tx.Exec(
		ctx,
		`INSERT INTO ranking_system.groups (id, description) VALUES (1, 'test group') RETURNING *`,
	)
	t.Logf("groups. Inserted %d row(s)", commandTag.RowsAffected())
	require.NoError(t, err)
	return err
}
