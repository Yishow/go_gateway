package pollinggroup

import (
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSQLRepositoryOwnerDeviceRoundTrip(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := t.Context()
	_, err := db.ExecContext(ctx, `INSERT INTO devices (id, name, protocol) VALUES ('dev-a','A','modbus_tcp')`)
	require.NoError(t, err)
	repo := NewSQLRepository(db)
	owner := "dev-a"
	now := time.Now()
	require.NoError(t, repo.Create(ctx, &schema.PollingGroup{ID: "owned", Name: "device:dev-a", IntervalMs: 500, Priority: 50, Enabled: true, OwnerDeviceID: &owner, CreatedAt: now, UpdatedAt: now}))
	require.NoError(t, repo.Create(ctx, &schema.PollingGroup{ID: "shared", Name: "Auto Polling", IntervalMs: 1000, Priority: 50, Enabled: true, CreatedAt: now, UpdatedAt: now}))

	got, err := repo.GetByID(ctx, "owned")
	require.NoError(t, err)
	require.NotNil(t, got.OwnerDeviceID)
	require.Equal(t, "dev-a", *got.OwnerDeviceID)
	shared, err := repo.GetByID(ctx, "shared")
	require.NoError(t, err)
	require.Nil(t, shared.OwnerDeviceID)

	listed, err := repo.List(ctx)
	require.NoError(t, err)
	owners := map[string]*string{}
	for _, g := range listed {
		owners[g.ID] = g.OwnerDeviceID
	}
	require.Equal(t, "dev-a", *owners["owned"])
	require.Nil(t, owners["shared"])

	// Generic updates never move or clear ownership.
	got.OwnerDeviceID = nil
	got.IntervalMs = 700
	require.NoError(t, repo.Update(ctx, got))
	again, err := repo.GetByID(ctx, "owned")
	require.NoError(t, err)
	require.Equal(t, 700, again.IntervalMs)
	require.Equal(t, "dev-a", *again.OwnerDeviceID)

	second := "dev-a"
	err = repo.Create(ctx, &schema.PollingGroup{ID: "owned-2", Name: "dup", IntervalMs: 500, Enabled: true, OwnerDeviceID: &second, CreatedAt: now, UpdatedAt: now})
	require.Error(t, err, "a device owns at most one polling group")
}

func TestSQLRepositoryReadsMigrationBackfilledGroup(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	ctx := t.Context()
	_, err := db.ExecContext(ctx, `INSERT INTO devices (id, name, protocol) VALUES ('dev-a','A','modbus_tcp')`)
	require.NoError(t, err)
	// Startup migration creates device-owned groups for devices without one.
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	groups, err := NewSQLRepository(db).List(ctx)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, "dev-a", *groups[0].OwnerDeviceID)
	require.Equal(t, 1000, groups[0].IntervalMs)
	require.False(t, groups[0].CreatedAt.IsZero())
}
