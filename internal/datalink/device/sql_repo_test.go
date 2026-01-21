package device

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	migrator := datalink.NewMigrator()
	err = migrator.Migrate(db)
	require.NoError(t, err)

	return db
}

func TestSQLRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)

	configJSON, _ := json.Marshal(map[string]interface{}{
		"ip":   "127.0.0.1",
		"port": 502,
	})

	dev := schema.Device{
		Name:             "Test Device",
		Description:      "A test device",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: string(configJSON),
	}

	err := repo.Create(context.Background(), &dev)
	require.NoError(t, err)
}

func TestSQLRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	ctx := context.Background()

	// 1. Create
	configJSON, _ := json.Marshal(map[string]interface{}{
		"host": "localhost",
	})
	newDev := schema.Device{
		ID:               "test-id-1", // Manually assigning ID for test simplicity if allowed
		Name:             "Test Device 1",
		Description:      "Desc 1",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: string(configJSON),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	// If ID generation is inside Repo, we might need a different approach.
	// But let's assume we can pass ID.
	err := repo.Create(ctx, &newDev)
	require.NoError(t, err)

	// 2. Get
	fetched, err := repo.GetByID(ctx, "test-id-1")
	require.NoError(t, err)
	assert.Equal(t, newDev.Name, fetched.Name)
	assert.Equal(t, newDev.Protocol, fetched.Protocol)

	// 3. Update
	fetched.Name = "Updated Name"
	err = repo.Update(ctx, fetched)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, "test-id-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)

	// 4. List
	list, err := repo.List(ctx, ListFilter{})
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// 5. Delete
	err = repo.Delete(ctx, "test-id-1")
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, "test-id-1")
	assert.Error(t, err)
}
