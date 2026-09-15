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

	// 2. Clean up before test (Optional: use a test DB)
	// Be careful not to wipe production if DSN points there.
	// We assume DSN points to a test DB.
	_, err = db.ExecContext(ctx, "DROP TABLE IF EXISTS timeseries CASCADE; DROP TABLE IF EXISTS mappings CASCADE; DROP TABLE IF EXISTS points CASCADE; DROP TABLE IF EXISTS polling_groups CASCADE; DROP TABLE IF EXISTS tags CASCADE; DROP TABLE IF EXISTS devices CASCADE; DROP TABLE IF EXISTS system_settings;")
	assert.NoError(t, err)

	// 3. Read and execute migration script
	migrationContent, err := os.ReadFile("../schema/migrations/001_initial_schema.up.sql")
	if err != nil {
		t.Fatalf("Failed to read migration file: %v", err)
	}

	_, err = db.ExecContext(ctx, string(migrationContent))
	assert.NoError(t, err, "Migration failed")

	// 4. Verify Partition Creation Functionality
	// The migration calls Create_timeseries_partition_monthly for current and next month.

	// Check if current month partition exists
	currentMonth := time.Now().Format("2006_01") // YYYY_MM
	partitionName := fmt.Sprintf("timeseries_%s", currentMonth)

	exists, err := tableExists(ctx, db, partitionName)
	assert.NoError(t, err)
	assert.True(t, exists, "Partition %s should exist", partitionName)

	// 5. Test Manual Partition Creation (Future)
	futureDate := time.Now().AddDate(0, 2, 0) // 2 months later
	var createdPartitionName string
	err = db.QueryRowContext(ctx, "SELECT create_timeseries_partition_monthly($1)", futureDate).Scan(&createdPartitionName)
	assert.NoError(t, err)

	expectedName := fmt.Sprintf("timeseries_%s", futureDate.Format("2006_01"))
	assert.Equal(t, expectedName, createdPartitionName)

	exists, err = tableExists(ctx, db, expectedName)
	assert.NoError(t, err)
	assert.True(t, exists, "Future partition %s should exist", expectedName)

	// 6. Insert Data into Partitioned Table
	tagID := "00000000-0000-0000-0000-000000000001" // Valid UUID
	ts := time.Now().UTC()
	_, err = db.ExecContext(ctx, `
		INSERT INTO timeseries (tag_id, ts, value_num, quality)
		VALUES ($1, $2, $3, $4)
	`, tagID, ts, 123.45, "good")
	assert.NoError(t, err, "Insert into partitioned table failed")

	// 7. Verify Data is retrievable
	var val float64
	err = db.QueryRowContext(ctx, "SELECT value_num FROM timeseries WHERE tag_id = $1", tagID).Scan(&val)
	assert.NoError(t, err)
	assert.Equal(t, 123.45, val)
}

func tableExists(ctx context.Context, db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE  table_schema = 'public'
			AND    table_name   = $1
		);
	`
	err := db.QueryRowContext(ctx, query, tableName).Scan(&exists)
	return exists, err
}
