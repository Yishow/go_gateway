package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresPartitionMigration(t *testing.T) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		t.Skip("Skipping Postgres test: POSTGRES_DSN not set")
	}

	// 1. Connect to Postgres
	db, err := sql.Open("pgx", dsn)
	assert.NoError(t, err)
	defer db.Close()
	ctx := t.Context()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping Postgres: %v", err)
	}

	// Isolate every object in one rollback-only schema; never drop public tables.
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	namespace := fmt.Sprintf("gw_partition_%d", time.Now().UnixNano())
	_, err = tx.ExecContext(ctx, `CREATE SCHEMA "`+namespace+`"`)
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `SET LOCAL search_path TO "`+namespace+`", public`)
	require.NoError(t, err)

	// 3. Read and execute migration script
	migrationContent, err := os.ReadFile("../schema/migrations/001_initial_schema.up.sql")
	if err != nil {
		t.Fatalf("Failed to read migration file: %v", err)
	}

	_, err = tx.ExecContext(ctx, string(migrationContent))
	require.NoError(t, err, "Migration failed")

	// 4. Verify Partition Creation Functionality
	// The migration calls Create_timeseries_partition_monthly for current and next month.

	// Check if current month partition exists
	currentMonth := time.Now().Format("2006_01") // YYYY_MM
	partitionName := fmt.Sprintf("timeseries_%s", currentMonth)

	exists, err := tableExists(ctx, tx, partitionName)
	assert.NoError(t, err)
	assert.True(t, exists, "Partition %s should exist", partitionName)

	// 5. Test Manual Partition Creation (Future)
	futureDate := time.Now().AddDate(0, 2, 0) // 2 months later
	var createdPartitionName string
	err = tx.QueryRowContext(ctx, "SELECT create_timeseries_partition_monthly($1::date)", futureDate).Scan(&createdPartitionName)
	assert.NoError(t, err)

	expectedName := fmt.Sprintf("timeseries_%s", futureDate.Format("2006_01"))
	assert.Equal(t, expectedName, createdPartitionName)

	exists, err = tableExists(ctx, tx, expectedName)
	assert.NoError(t, err)
	assert.True(t, exists, "Future partition %s should exist", expectedName)

	// 6. Insert Data into Partitioned Table
	tagID := "00000000-0000-0000-0000-000000000001" // Valid UUID
	ts := time.Now().UTC()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO timeseries (tag_id, ts, value_num, quality)
		VALUES ($1, $2, $3, $4)
	`, tagID, ts, 123.45, "good")
	assert.NoError(t, err, "Insert into partitioned table failed")

	// 7. Verify Data is retrievable
	var val float64
	err = tx.QueryRowContext(ctx, "SELECT value_num FROM timeseries WHERE tag_id = $1", tagID).Scan(&val)
	assert.NoError(t, err)
	assert.Equal(t, 123.45, val)
}

func tableExists(ctx context.Context, tx *sql.Tx, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE  table_schema = current_schema()
			AND    table_name   = $1
		);
	`
	err := tx.QueryRowContext(ctx, query, tableName).Scan(&exists)
	return exists, err
}
