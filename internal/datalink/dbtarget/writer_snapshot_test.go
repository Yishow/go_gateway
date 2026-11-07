package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// blockingDestination is a connector whose connection attempts never complete
// until the test releases it.
type blockingDestination struct {
	entered chan struct{}
	release chan struct{}
}

func newBlockingDestination(t *testing.T) *blockingDestination {
	t.Helper()
	d := &blockingDestination{entered: make(chan struct{}, 64), release: make(chan struct{})}
	t.Cleanup(d.Release)
	return d
}

func (d *blockingDestination) Open(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
	d.entered <- struct{}{}
	<-d.release
	return nil, errors.New("connection attempt timed out")
}

func (d *blockingDestination) Release() {
	select {
	case <-d.release:
	default:
		close(d.release)
	}
}

func createSnapshotConnector(t *testing.T, repo *SQLConnectorRepository, kind schema.DatabaseConnectorKind, config string, interval int) *schema.DatabaseConnector {
	t.Helper()
	connector := &schema.DatabaseConnector{
		ID: "connector-" + string(kind), Name: "snapshot-" + string(kind), Kind: kind,
		ConnectionConfig: config, Status: schema.DatabaseConnectorStatusReady, Enabled: true,
		DefaultWriteIntervalSeconds: interval, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(t, repo.Create(t.Context(), connector))
	resetWriteHistory(connector.ID)
	t.Cleanup(func() { resetWriteHistory(connector.ID) })
	return connector
}

func createSnapshotMapping(t *testing.T, repo *SQLTargetMappingRepository, connectorID, tagID, column string) {
	t.Helper()
	require.NoError(t, repo.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "mapping-" + tagID, TagID: tagID, ConnectorID: connectorID,
		TableName: "snapshot_values", ColumnName: column, WriteMode: schema.DatabaseWriteModeInsert,
		Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}))
}

func TestWriter_UnreachableConnectorDoesNotBlockIngest(t *testing.T) {
	mainDB := openMigratedTestDB(t)
	connectors := NewSQLConnectorRepository(mainDB)
	mappings := NewSQLTargetMappingRepository(mainDB)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindPostgres, `{"host":"192.0.2.1","database":"lost"}`, 15)
	for i := range 32 {
		createSnapshotMapping(t, mappings, connector.ID, fmt.Sprintf("tag-%02d", i), fmt.Sprintf("d%d", i))
	}
	destination := newBlockingDestination(t)
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, OpenDestination: destination.Open})
	t.Cleanup(func() {
		destination.Release()
		_ = writer.Close(context.Background())
	})
	require.NoError(t, writer.RefreshMappings(t.Context()))

	base := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	ingestScan := func(round int) {
		t.Helper()
		for i := range 32 {
			done := make(chan error, 1)
			started := time.Now()
			go func() {
				done <- writer.WriteTagValue(t.Context(), fmt.Sprintf("tag-%02d", i), float64(round*100+i), base.Add(time.Duration(round)*time.Second))
			}()
			select {
			case err := <-done:
				require.ErrorIs(t, err, ErrSnapshotAccepted)
				require.Less(t, time.Since(started), 50*time.Millisecond, "ingest of value %d waited on the connector", i)
			case <-time.After(50 * time.Millisecond):
				t.Fatalf("ingest of value %d (round %d) blocked on an unreachable connector", i, round)
			}
		}
	}

	ingestScan(0)
	tick <- base.Add(time.Minute)
	select {
	case <-destination.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("connector worker never attempted the due snapshot")
	}
	// The worker is now stuck inside a connection attempt; ingest must not wait for it.
	ingestScan(1)

	destination.Release()
	require.Eventually(t, func() bool {
		updated, err := connectors.GetByID(t.Context(), connector.ID)
		return err == nil && updated.LastWriteStatus == deliveryOutcomeFailed
	}, 2*time.Second, 10*time.Millisecond, "connector status must report the failed write attempt")
}

func TestWriter_WritesOnlyLatestValuePerInterval(t *testing.T) {
	mainDB := openMigratedTestDB(t)
	connectors := NewSQLConnectorRepository(mainDB)
	mappings := NewSQLTargetMappingRepository(mainDB)
	targetDSN := filepath.Join(t.TempDir(), "snapshot.db")
	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	t.Cleanup(func() { _ = targetDB.Close() })
	_, err = targetDB.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)

	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, targetDSN), 15)
	createSnapshotMapping(t, mappings, connector.ID, "tag-d0", "d0")
	tick := make(chan time.Time, 1)
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.NoError(t, writer.RefreshMappings(t.Context()))
	resetWriteHistory(connector.ID)

	intervalStart := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	for i, value := range []float64{10, 11, 12} {
		require.ErrorIs(t, writer.WriteTagValue(t.Context(), "tag-d0", value, intervalStart.Add(time.Duration(1+4*i)*time.Second)), ErrSnapshotAccepted)
	}
	countRows := func() int {
		var count int
		require.NoError(t, targetDB.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
		return count
	}
	require.Zero(t, countRows(), "ingest must not write synchronously")

	require.NoError(t, writer.flushDue(t.Context(), intervalStart.Add(14*time.Second)))
	require.Zero(t, countRows(), "the interval has not ended yet")

	tick <- intervalStart.Add(15 * time.Second)
	require.Eventually(t, func() bool { return countRows() == 1 }, 2*time.Second, 10*time.Millisecond)
	var written float64
	require.NoError(t, targetDB.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&written))
	require.Equal(t, 12.0, written)

	require.NoError(t, writer.flushDue(t.Context(), intervalStart.Add(30*time.Second)))
	require.Equal(t, 1, countRows(), "an interval without a new value writes no row")

	require.Eventually(t, func() bool { return len(listWriteHistory(connector.ID, 10)) == 1 }, time.Second, 10*time.Millisecond)
	history := listWriteHistory(connector.ID, 10)
	require.Len(t, history, 1)
	require.Equal(t, deliveryOutcomeSuccess, history[0].Status)
	require.Equal(t, 15, history[0].EffectiveIntervalSeconds)
}

func TestWriter_DoesNotWriteTwiceWithinOneInterval(t *testing.T) {
	db := openMigratedTestDB(t)
	connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
	path := filepath.Join(t.TempDir(), "interval.db")
	target, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
	require.NoError(t, err)
	connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), 15)
	createSnapshotMapping(t, mappings, connector.ID, "tag-interval", "d0")
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	at := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "tag-interval", 10, at.Add(time.Second)), ErrSnapshotAccepted)
	require.NoError(t, writer.flushDue(t.Context(), at.Add(15*time.Second)))
	mapping, err := mappings.GetByID(t.Context(), "mapping-tag-interval")
	require.NoError(t, err)
	mapping.Enabled = false
	require.NoError(t, mappings.Update(t.Context(), mapping))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	mapping.Enabled = true
	require.NoError(t, mappings.Update(t.Context(), mapping))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "tag-interval", 12, at.Add(10*time.Second)), ErrSnapshotAccepted)
	require.NoError(t, writer.flushDue(t.Context(), at.Add(16*time.Second)))
	var count int
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, writer.flushDue(t.Context(), at.Add(30*time.Second)))
	require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
	require.Equal(t, 2, count)
}
