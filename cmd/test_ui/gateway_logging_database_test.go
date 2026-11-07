package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/diagnostics"

	"github.com/stretchr/testify/require"
)

func TestManagedSlogDatabaseFailuresReachWarnAndFileWithoutSecrets(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	dbPath := filepath.Join(root, "config.db")
	require.NoError(t, os.WriteFile(dbPath, []byte("fixture"), 0o600))
	logPath := filepath.Join(root, "logs", "runtime.jsonl")
	sink, err := diagnostics.OpenFileSink(diagnostics.FileOptions{Path: logPath, DatabasePath: dbPath, DatabaseID: strings.Repeat("a", 64)})
	require.NoError(t, err)
	broker, err := diagnostics.New(diagnostics.Options{Level: "warn", File: sink})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	logger := slog.New(managedSlog{broker: broker})
	first := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	logger.Error("flush database target snapshot", "connector_id", "connector-A", "error", errors.New("password=fixture-secret host=private.example"))
	logger.Warn("database target snapshot is still failing", "connector_id", "connector-A", "failure_count", 11, "first_failure_at", first)
	require.NoError(t, broker.Flush(t.Context()))
	snapshot, err := broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	require.Len(t, snapshot.Records, 2)
	require.Equal(t, "runtime.database_write_failed", snapshot.Records[0].Code)
	require.Equal(t, "runtime.database_write_failing", snapshot.Records[1].Code)
	require.Equal(t, "connector-A", snapshot.Records[1].Fields["connector_id"])
	require.EqualValues(t, 11, snapshot.Records[1].Fields["failure_count"])
	require.Equal(t, first.Format(time.RFC3339Nano), snapshot.Records[1].Fields["first_failure_at"])
	logger.Warn("database target snapshot is still failing", "connector_id", "private.example",
		"failure_count", -1, "first_failure_at", "password=fixture-secret", "failure_reason", "password=fixture-secret")
	require.NoError(t, broker.Flush(t.Context()))
	snapshot, err = broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	require.Empty(t, snapshot.Records[len(snapshot.Records)-1].Fields)
	data, err := os.ReadFile(logPath)
	require.NoError(t, err)
	require.Contains(t, string(data), `"code":"runtime.database_write_failing"`)
	require.NotContains(t, string(data), "fixture-secret")
	require.NotContains(t, string(data), "private.example")
}

func TestManagedSlogDatabaseRecoveryIsVisibleAtInfoLevel(t *testing.T) {
	broker, err := diagnostics.New(diagnostics.Options{Level: "info"})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	logger := slog.New(managedSlog{broker: broker})
	logger.Info("database target snapshot recovered", "connector_id", "connector-A", "failure_count", 3)
	require.NoError(t, broker.Flush(t.Context()))
	snapshot, err := broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	require.Len(t, snapshot.Records, 1)
	require.Equal(t, "runtime.database_write_recovered", snapshot.Records[0].Code)
	require.Equal(t, "connector-A", snapshot.Records[0].Fields["connector_id"])
	require.EqualValues(t, 3, snapshot.Records[0].Fields["failure_count"])
}

func TestLegacyWriterFailureReachesManagedLog(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "config.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalinkbase.NewMigrator().Migrate(db))
	connectors := dbtarget.NewSQLConnectorRepository(db)
	mappings := dbtarget.NewSQLTargetMappingRepository(db)
	base := time.Date(2026, 11, 7, 10, 0, 0, 0, time.UTC)
	require.NoError(t, connectors.Create(t.Context(), &schema.DatabaseConnector{
		ID: "connector-A", Name: "test target", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"dsn":"unused.db"}`, Enabled: true, DefaultWriteIntervalSeconds: 15,
		CreatedAt: base, UpdatedAt: base,
	}))
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "mapping-A", TagID: "tag-A", ConnectorID: "connector-A", TableName: "readings", ColumnName: "value",
		WriteMode: schema.DatabaseWriteModeInsert, Enabled: true, CreatedAt: base, UpdatedAt: base,
	}))
	broker, err := diagnostics.New(diagnostics.Options{Level: "warn"})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	previous := slog.Default()
	slog.SetDefault(slog.New(managedSlog{broker: broker}))
	t.Cleanup(func() { slog.SetDefault(previous) })
	clock := &activationSnapshotClock{now: base}
	tick := make(chan time.Time, 1)
	writer := dbtarget.NewWriterWithConfig(connectors, mappings, dbtarget.WriterConfig{
		FlushTick: tick, Now: clock.Now,
		OpenDestination: func(schema.DatabaseConnectorKind, dbtarget.ConnectionConfig) (*datalinkbase.DBManager, error) {
			return nil, fmt.Errorf("password=fixture-secret host=private.example: %w", dbtarget.ErrDestinationUnreachable)
		},
	})
	t.Cleanup(func() { _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "tag-A", 1, base.Add(time.Second)), dbtarget.ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	var event diagnostics.Event
	require.Eventually(t, func() bool {
		if broker.Flush(t.Context()) != nil {
			return false
		}
		snapshot, snapshotErr := broker.Snapshot(diagnostics.Query{})
		if snapshotErr != nil {
			return false
		}
		for _, record := range snapshot.Records {
			if record.Code == "runtime.database_write_failed" {
				event = record
				return true
			}
		}
		return false
	}, 3*time.Second, 10*time.Millisecond)
	require.Equal(t, "connector-A", event.Fields["connector_id"])
	require.Equal(t, "connection", event.Fields["failure_reason"])
	require.EqualValues(t, 1, event.Fields["failure_count"])
	require.Equal(t, base.Add(15*time.Second).Format(time.RFC3339Nano), event.Fields["first_failure_at"])
	require.NotContains(t, fmt.Sprint(event), "fixture-secret")
	require.NotContains(t, fmt.Sprint(event), "private.example")
}
