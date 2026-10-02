package dbtarget

import (
	"os"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestReadOnlyTableInspectorSQLiteReadsExistingTable(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE sensor_values (reactor_temp REAL NOT NULL)`)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindSQLite, `{"dsn":"`+target+`"}`)
	inspector := NewReadOnlyTableInspector(service)

	result, err := inspector.InspectTable(t.Context(), id, "main", "sensor_values")

	requireInspection(t, result, err, TableInspectionExists, "")
	require.Equal(t, []string{"reactor_temp"}, []string{result.Columns[0].Name})
}

func TestReadOnlyTableInspectorSQLiteMissingFileDoesNotCreate(t *testing.T) {
	target := filepath.Join(t.TempDir(), "missing.db")
	service, id := newInspectionService(t, schema.DatabaseConnectorKindSQLite, `{"dsn":"`+target+`"}`)
	inspector := NewReadOnlyTableInspector(service)

	result, err := inspector.InspectTable(t.Context(), id, "main", "sensor_values")

	requireInspection(t, result, err, TableInspectionFailed, "connection_failed")
	_, statErr := os.Stat(target)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestReadOnlyTableInspectorSQLiteUsesDsnBeforePath(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE sensor_values (reactor_temp REAL)`)
	missingPath := filepath.Join(t.TempDir(), "must-not-open.db")
	service, id := newInspectionService(t, schema.DatabaseConnectorKindSQLite,
		`{"dsn":"`+target+`","path":"`+missingPath+`"}`)
	inspector := NewReadOnlyTableInspector(service)

	result, err := inspector.InspectTable(t.Context(), id, "main", "sensor_values")

	requireInspection(t, result, err, TableInspectionExists, "")
	_, statErr := os.Stat(missingPath)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestSQLiteReadOnlyDSNRejectsMemoryScope(t *testing.T) {
	for _, raw := range []string{":memory:", "file::memory:?cache=shared", "memory?cache=shared"} {
		_, err := sqliteReadOnlyDSN(ConnectionConfig{"dsn": raw})
		require.Error(t, err, raw)
	}
}
