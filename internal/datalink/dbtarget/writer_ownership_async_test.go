package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_OwnershipCheckDoesNotBlockIngest(t *testing.T) {
	mainDB := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(mainDB), NewSQLTargetMappingRepository(mainDB)
	path := filepath.Join(t.TempDir(), "target.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	now := time.Now()
	require.NoError(t, connectors.Create(t.Context(), &schema.DatabaseConnector{
		ID: "ownership-connector", Name: "ownership", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: fmt.Sprintf(`{"dsn":%q}`, path), Enabled: true,
		DefaultWriteIntervalSeconds: 15, CreatedAt: now, UpdatedAt: now,
	}))
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "ownership-mapping", TagID: "ownership-tag", ConnectorID: "ownership-connector",
		TableName: "snapshot_values", ColumnName: "d0", WriteMode: schema.DatabaseWriteModeInsert,
		Enabled: true, CreatedAt: now, UpdatedAt: now,
	}))
	entered, release := make(chan struct{}, 1), make(chan struct{})
	var once sync.Once
	releaseCheck := func() { once.Do(func() { close(release) }) }
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
		FlushTick: tick,
		SuppressMapping: func(string, string) bool {
			entered <- struct{}{}
			<-release
			return true
		},
	})
	t.Cleanup(func() { releaseCheck(); _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	done := make(chan error, 1)
	go func() { done <- writer.WriteTagValue(t.Context(), "ownership-tag", 12, at) }()
	select {
	case err := <-done:
		if err != nil {
			require.EqualError(t, err, "snapshot accepted")
		}
	case <-time.After(50 * time.Millisecond):
		releaseCheck()
		<-done
		t.Fatal("ingest waited for the ownership check")
	}
	tick <- at.Add(time.Minute)
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not check ownership")
	}
	releaseCheck()
	require.NoError(t, writer.Close(context.Background()))
	var rows int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
	require.Zero(t, rows, "a group-owned output must not be written by the legacy writer")
}
