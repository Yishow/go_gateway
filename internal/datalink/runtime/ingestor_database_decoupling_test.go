package runtime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// TestHandleCollectedValue_UnreachableDatabaseConnectorDoesNotDelayValueEvents
// drives one 32-point scan through the real legacy database writer while the
// connector's connection attempt never completes.
func TestHandleCollectedValue_UnreachableDatabaseConnectorDoesNotDelayValueEvents(t *testing.T) {
	const points = 32
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "main.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, datalink.NewMigrator().Migrate(db))

	connectors := dbtarget.NewSQLConnectorRepository(db)
	mappings := dbtarget.NewSQLTargetMappingRepository(db)
	connector := &schema.DatabaseConnector{
		ID: "connector-unreachable", Name: "unreachable", Kind: schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: `{"host":"192.0.2.1","database":"lost"}`, Status: schema.DatabaseConnectorStatusReady,
		Enabled: true, DefaultWriteIntervalSeconds: 15, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(t, connectors.Create(t.Context(), connector))

	index := make(map[string][]mappingBinding, points)
	for i := range points {
		tagID := fmt.Sprintf("tag-%02d", i)
		require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
			ID: "mapping-" + tagID, TagID: tagID, ConnectorID: connector.ID, TableName: "scan_values",
			ColumnName: fmt.Sprintf("d%d", i), WriteMode: schema.DatabaseWriteModeInsert, Enabled: true,
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}))
		index[fmt.Sprintf("p%02d", i)] = []mappingBinding{{TagID: tagID, TagDataType: schema.DataTypeInt16}}
	}

	entered, release := make(chan struct{}, points*2), make(chan struct{})
	releaseOnce := func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}
	tick := make(chan time.Time, 1)
	writer := dbtarget.NewWriterWithConfig(connectors, mappings, dbtarget.WriterConfig{
		FlushTick: tick,
		OpenDestination: func(schema.DatabaseConnectorKind, dbtarget.ConnectionConfig) (*datalink.DBManager, error) {
			entered <- struct{}{}
			<-release
			return nil, errors.New("connection attempt timed out")
		},
	})
	t.Cleanup(func() {
		releaseOnce()
		_ = writer.Close(context.Background())
	})
	s := &Service{config: Config{UpdatePointState: false}, writer: &mockWriter{}, target: writer, mappingIndex: index}
	scanAt := time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC)
	ingest := func(pointIndex int, at time.Time) {
		t.Helper()
		done := make(chan struct{})
		started := time.Now()
		go func() {
			defer close(done)
			s.handleCollectedValue(t.Context(), collector.CollectedValue{
				PointID: fmt.Sprintf("p%02d", pointIndex), DeviceID: "dev-1", Value: int16(pointIndex),
				Timestamp: at, Quality: schema.QualityGood,
			})
		}()
		select {
		case <-done:
			require.Less(t, time.Since(started), 50*time.Millisecond)
		case <-time.After(50 * time.Millisecond):
			t.Fatalf("ingest of point %d waited on the unreachable connector", pointIndex)
		}
	}

	// Make the connector worker start a connection attempt that never completes.
	ingest(0, scanAt)
	tick <- scanAt.Add(time.Minute)
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("connector worker never attempted the due snapshot")
	}

	events, unsubscribe := s.SubscribeValueEvents("", nil)
	defer unsubscribe()
	for i := range points {
		ingest(i, scanAt.Add(time.Second))
	}
	require.Len(t, events, points, "every ingested value must broadcast its runtime value event")
	diagnostics := s.databaseDeliveryDiagnostics("")
	require.Len(t, diagnostics, points)
	for _, diagnostic := range diagnostics {
		require.Equal(t, DatabaseDeliveryStatusAccepted, diagnostic.Status)
		require.NotContains(t, diagnostic.Stages, DatabaseDeliveryStageDBWriteSucceeded)
	}

	releaseOnce()
	require.Eventually(t, func() bool {
		updated, getErr := connectors.GetByID(t.Context(), connector.ID)
		return getErr == nil && updated.LastWriteStatus == "failed"
	}, 2*time.Second, 10*time.Millisecond, "connector status must report the failed write attempt")
}
