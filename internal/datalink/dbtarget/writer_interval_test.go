package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_EffectiveSnapshotIntervalResolution(t *testing.T) {
	override := 10
	invalid := 0
	for _, tc := range []struct {
		name              string
		mappingInterval   *int
		connectorInterval int
		wantSeconds       int
	}{
		{"mapping override", &override, 30, 10},
		{"connector default", nil, 30, 30},
		{"fallback default", nil, 0, 15},
		{"invalid mapping interval", &invalid, 30, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := openMigratedTestDB(t)
			connectors, mappings := NewSQLConnectorRepository(db), NewSQLTargetMappingRepository(db)
			path := filepath.Join(t.TempDir(), "interval.db")
			target, err := sql.Open("sqlite", path)
			require.NoError(t, err)
			t.Cleanup(func() { _ = target.Close() })
			_, err = target.ExecContext(t.Context(), `CREATE TABLE snapshot_values (d0 REAL)`)
			require.NoError(t, err)
			connector := createSnapshotConnector(t, connectors, schema.DatabaseConnectorKindSQLite, fmt.Sprintf(`{"dsn":%q}`, path), tc.connectorInterval)
			require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
				ID: "interval-mapping", TagID: "interval-tag", ConnectorID: connector.ID,
				TableName: "snapshot_values", ColumnName: "d0", WriteMode: schema.DatabaseWriteModeInsert,
				WriteIntervalSeconds: tc.mappingInterval, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}))
			writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time)})
			t.Cleanup(func() { _ = writer.Close(context.Background()) })
			at := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
			writeErr := writer.WriteTagValue(t.Context(), "interval-tag", 12, at.Add(time.Second))
			if tc.wantSeconds == 0 {
				require.ErrorContains(t, writeErr, "寫入間隔")
				return
			}
			require.ErrorIs(t, writeErr, ErrSnapshotAccepted)
			require.NoError(t, writer.flushDue(t.Context(), at.Add(time.Duration(tc.wantSeconds-1)*time.Second), false))
			var count int
			require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
			require.Zero(t, count)
			require.NoError(t, writer.flushDue(t.Context(), at.Add(time.Duration(tc.wantSeconds)*time.Second), false))
			require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&count))
			require.Equal(t, 1, count)
		})
	}
}
