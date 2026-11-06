package main

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func productionSnapshotTarget(t *testing.T, f scheduledProductionFixture) (*sql.DB, string, *schema.DatabaseConnector, *schema.DatabaseTargetMapping) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "target.db")
	destination, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, destination.Close()) })
	_, err = destination.ExecContext(t.Context(), `CREATE TABLE readings (first INTEGER, second INTEGER)`)
	require.NoError(t, err)
	interval := 1
	connector, err := f.services.dbTarget.Create(t.Context(), dbtarget.CreateConnectorRequest{
		Name: "snapshot target", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{"dsn": path}, DefaultWriteIntervalSeconds: &interval,
	})
	require.NoError(t, err)
	bindings, err := f.services.mapping.List(t.Context(), mapping.ListFilter{})
	require.NoError(t, err)
	require.NotEmpty(t, bindings)
	targetMapping, err := f.services.dbMapping.Create(t.Context(), dbtarget.CreateTargetMappingRequest{
		TagID: bindings[0].TagID, ConnectorID: connector.ID, TableSchema: "main", TableName: "readings",
		ColumnName: "first", WriteMode: schema.DatabaseWriteModeInsert,
	})
	require.NoError(t, err)
	return destination, path, connector, targetMapping
}

func snapshotRowCount(db *sql.DB, where string) int {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM readings WHERE ` + where).Scan(&count); err != nil {
		return -1
	}
	return count
}

func TestProductionSnapshotWorker_LoadsSavedMappingAtStartup(t *testing.T) {
	f, _, _ := newAcceptanceScheduledLoad(t, 1, 3, time.Millisecond, true)
	destination, _, _, _ := productionSnapshotTarget(t, f)
	require.NoError(t, f.services.runtime.Stop(context.Background()))
	services, err := wireGatewayServicesContext(t.Context(), f.db, f.manager, nil, f.capture)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(f.db)
	hydrateStartupShare(t.Context(), services, revisions, modbusshare.NewReconciler(services.modbusShare, revisions))
	require.NoError(t, services.runtime.Start(t.Context()))
	require.Eventually(t, func() bool { return snapshotRowCount(destination, "first IS NOT NULL") > 0 }, 8*time.Second, 20*time.Millisecond, "startup must load a saved connector and mapping")
}

func TestProductionSnapshotWorker_Lifecycle(t *testing.T) {
	f, _, _ := newAcceptanceScheduledLoad(t, 1, 3, time.Millisecond)
	destination, _, connector, targetMapping := productionSnapshotTarget(t, f)
	require.NoError(t, f.services.runtime.Start(t.Context()))
	require.Eventually(t, func() bool { return snapshotRowCount(destination, "first IS NOT NULL") > 0 }, 8*time.Second, 20*time.Millisecond, "production worker must pick up the new mapping")

	disabled := false
	connector, err := f.services.dbTarget.Update(t.Context(), connector.ID, dbtarget.UpdateConnectorRequest{Enabled: &disabled, ExpectedIdentityRevision: &connector.IdentityRevision})
	require.NoError(t, err)
	// Allow the production refresh tick and any in-flight write to complete.
	time.Sleep(2 * time.Second)
	before := snapshotRowCount(destination, "1=1")
	time.Sleep(2 * time.Second)
	require.Equal(t, before, snapshotRowCount(destination, "1=1"), "disabled connector must not write")

	enabled := true
	connector, err = f.services.dbTarget.Update(t.Context(), connector.ID, dbtarget.UpdateConnectorRequest{Enabled: &enabled, ExpectedIdentityRevision: &connector.IdentityRevision})
	require.NoError(t, err)
	require.Eventually(t, func() bool { return snapshotRowCount(destination, "1=1") > before }, 8*time.Second, 20*time.Millisecond, "re-enabled connector must resume")

	column := "second"
	_, err = f.services.dbMapping.Update(t.Context(), targetMapping.ID, dbtarget.UpdateTargetMappingRequest{ColumnName: &column})
	require.NoError(t, err)
	require.Eventually(t, func() bool { return snapshotRowCount(destination, "second IS NOT NULL") > 0 }, 8*time.Second, 20*time.Millisecond, "updated mapping must take effect")
	require.NoError(t, f.services.runtime.Stop(context.Background()))
	closedCount := snapshotRowCount(destination, "1=1")
	time.Sleep(1200 * time.Millisecond)
	require.Equal(t, closedCount, snapshotRowCount(destination, "1=1"), "runtime shutdown must stop the worker")
}
