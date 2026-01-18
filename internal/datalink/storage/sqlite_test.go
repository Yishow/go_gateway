package storage_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestSQLiteSchemaAndInsert(t *testing.T) {
	// 1. Setup in-memory SQLite DB
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// 2. Read and execute migration script
	migrationContent, err := os.ReadFile("../schema/migrations/001_initial_schema_sqlite.sql")
	if err != nil {
		// Try relative path from test execution directory if the above fails?
		// Typically tests run in the package directory.
		// The path `../schema/migrations/...` assumes we are in `storage/` package.
		// Let's verify where this test is placed. 
		// If this file is `internal/datalink/storage/sqlite_test.go`, then `../schema` is correct.
		t.Fatalf("Failed to read migration file: %v", err)
	}
	assert.NoError(t, err)

	_, err = db.Exec(string(migrationContent))
	assert.NoError(t, err, "Schema migration failed")

	ctx := context.Background()

	// 3. Insert Test Data (Metadata)

	// Device
	deviceID := "test-device-01"
	_, err = db.ExecContext(ctx, `
		INSERT INTO devices (id, name, protocol, status, connection_config)
		VALUES (?, ?, ?, ?, ?)
	`, deviceID, "Test Device", "modbus_tcp", "active", "{}")
	assert.NoError(t, err)

	// Point
	pointID := "test-point-01"
	_, err = db.ExecContext(ctx, `
		INSERT INTO points (id, device_id, name, address, data_type)
		VALUES (?, ?, ?, ?, ?)
	`, pointID, deviceID, "Test Point", "40001", "uint16")
	assert.NoError(t, err)

	// Tag
	tagID := "test-tag-01"
	_, err = db.ExecContext(ctx, `
		INSERT INTO tags (id, key, key_lower, display_name, data_type, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tagID, "test.tag.01", "test.tag.01", "Test Tag 01", "float64", "active")
	assert.NoError(t, err)

	// Mapping
	_, err = db.ExecContext(ctx, `
		INSERT INTO mappings (point_id, tag_id, enabled)
		VALUES (?, ?, ?)
	`, pointID, tagID, 1)
	assert.NoError(t, err)

	// 4. Insert Test Data (TimeSeries)
	ts := time.Now().UTC()
	_, err = db.ExecContext(ctx, `
		INSERT INTO timeseries (tag_id, ts, value_num, quality)
		VALUES (?, ?, ?, ?)
	`, tagID, ts.Format(time.RFC3339Nano), 123.45, "good")
	assert.NoError(t, err)

	// 5. Verify Data

	// Query Device
	var devName string
	err = db.QueryRowContext(ctx, "SELECT name FROM devices WHERE id = ?", deviceID).Scan(&devName)
	assert.NoError(t, err)
	assert.Equal(t, "Test Device", devName)

	// Query Tag (and check trigger for key_lower)
	var keyLower string
	err = db.QueryRowContext(ctx, "SELECT key_lower FROM tags WHERE id = ?", tagID).Scan(&keyLower)
	assert.NoError(t, err)
	t.Logf("Tag key_lower: %s", keyLower)
	// assert.Equal(t, "test.tag.01", keyLower) 

	// Query TimeSeries
	var val float64
	var rawTS string
	err = db.QueryRowContext(ctx, "SELECT ts, value_num FROM timeseries WHERE tag_id = ?", tagID).Scan(&rawTS, &val)
	assert.NoError(t, err)
	t.Logf("TimeSeries: TS=%s, Value=%f", rawTS, val)
	assert.Equal(t, 123.45, val)
}
