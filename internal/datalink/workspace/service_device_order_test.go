package workspace

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_AttachDeviceAndReplaceOrderPersistAcrossRestart(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "studio-v2-device-order.db")
	ctx := context.Background()

	db := openWorkspaceTestDB(t, dbPath)
	svc := NewService(NewSQLRepository(db))

	record, err := svc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	require.Equal(t, WorkspaceStatusReady, record.Status)
	require.Equal(t, []string{"dev-A"}, record.OrderedDeviceIDs)

	record, err = svc.AttachDevice(ctx, "dev-B")
	require.NoError(t, err)
	require.Equal(t, []string{"dev-A", "dev-B"}, record.OrderedDeviceIDs)

	record, err = svc.AttachDevice(ctx, "dev-C")
	require.NoError(t, err)
	require.Equal(t, []string{"dev-A", "dev-B", "dev-C"}, record.OrderedDeviceIDs)

	record, err = svc.ReplaceDeviceOrder(ctx, []string{"dev-C", "dev-A", "dev-B"})
	require.NoError(t, err)
	require.Equal(t, []string{"dev-C", "dev-A", "dev-B"}, record.OrderedDeviceIDs)
	require.NoError(t, db.Close())

	restartedDB := openWorkspaceTestDB(t, dbPath)
	defer restartedDB.Close()

	restartedSvc := NewService(NewSQLRepository(restartedDB))
	restartedRecord, err := restartedSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, WorkspaceStatusReady, restartedRecord.Status)
	require.Equal(t, []string{"dev-C", "dev-A", "dev-B"}, restartedRecord.OrderedDeviceIDs)
}

func TestService_DetachDeviceRemovesOrderAndResetsEmptyStatus(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	_, err := svc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	_, err = svc.AttachDevice(ctx, "dev-B")
	require.NoError(t, err)

	record, err := svc.DetachDevice(ctx, "dev-A")
	require.NoError(t, err)
	require.Equal(t, WorkspaceStatusReady, record.Status)
	require.Equal(t, []string{"dev-B"}, record.OrderedDeviceIDs)

	record, err = svc.DetachDevice(ctx, "dev-B")
	require.NoError(t, err)
	require.Equal(t, WorkspaceStatusEmpty, record.Status)
	require.Empty(t, record.OrderedDeviceIDs)
}

func TestService_ReplaceDeviceOrderRejectsUnknownDeviceIDs(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	_, err := svc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)
	_, err = svc.AttachDevice(ctx, "dev-B")
	require.NoError(t, err)

	_, err = svc.ReplaceDeviceOrder(ctx, []string{"dev-B", "dev-missing"})
	require.Error(t, err)
}
