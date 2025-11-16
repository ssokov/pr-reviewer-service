package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	connString := os.Getenv("TEST_DATABASE_URL")
	if connString == "" {
		connString = "postgres://max_superuser:max_superuser@localhost:5432/pr_system?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to database: %v", err)
		return nil
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}
