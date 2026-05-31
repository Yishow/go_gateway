package workspace

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func openWorkspaceTestDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)

	migrator := datalink.NewMigrator()
	require.NoError(t, migrator.Migrate(db))

	return db
}

func TestService_GetOrCreatePersistsSingletonAcrossRestart(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "studio-v2-workspace.db")
	ctx := context.Background()

	db := openWorkspaceTestDB(t, dbPath)
	svc := NewService(NewSQLRepository(db))

	first, err := svc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, first.ID)
	require.Equal(t, WorkspaceKindSingle, first.Kind)
	require.Equal(t, WorkspaceStatusEmpty, first.Status)
	require.Empty(t, first.OrderedDeviceIDs)
	require.NoError(t, db.Close())

	restartedDB := openWorkspaceTestDB(t, dbPath)
	defer restartedDB.Close()

	restartedSvc := NewService(NewSQLRepository(restartedDB))
	second, err := restartedSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, first.ID, second.ID)
	require.Equal(t, first.Kind, second.Kind)
	require.Equal(t, first.Status, second.Status)
	require.Equal(t, first.OrderedDeviceIDs, second.OrderedDeviceIDs)
}

func TestService_GetOrCreateDoesNotImportLegacyDevices(t *testing.T) {
	db := openWorkspaceTestDB(t, ":memory:")
	defer db.Close()

	ctx := context.Background()
	deviceRepo := device.NewSQLRepository(db)

	err := deviceRepo.Create(ctx, &schema.Device{
		Name:             "Legacy PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
	})
	require.NoError(t, err)

	svc := NewService(NewSQLRepository(db))
	record, err := svc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, WorkspaceKindSingle, record.Kind)
	require.Equal(t, WorkspaceStatusEmpty, record.Status)
	require.Empty(t, record.OrderedDeviceIDs)
}
