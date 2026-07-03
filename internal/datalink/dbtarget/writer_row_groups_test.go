package dbtarget

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestWriter_InsertModeRowGroupsEmitSeparateRowsForSharedColumn(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "target-row-groups.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	defer targetDB.Close()
	_, err = targetDB.Exec(`
		CREATE TABLE row_group_values (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			temperature_c REAL
		)
	`)
	require.NoError(t, err)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)
	tickCh := make(chan time.Time, 2)
	writer := NewWriterWithConfig(connectorRepo, mappingRepo, WriterConfig{FlushTick: tickCh})
	t.Cleanup(func() {
		require.NoError(t, writer.Close(context.Background()))
	})

	interval := 15
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "row-group-insert",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
		DefaultWriteIntervalSeconds: &interval,
	})
	require.NoError(t, err)
	resetWriteHistory(connector.ID)

	tagA, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line-a.temperature",
		DisplayName: "Line A Temperature",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	tagB, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "line-b.temperature",
		DisplayName: "Line B Temperature",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	groupA := "group-temperature:point-a"
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:       tagA.ID,
		ConnectorID: connector.ID,
		TableSchema: "main",
		TableName:   "row_group_values",
		ColumnName:  "temperature_c",
		WriteMode:   schema.DatabaseWriteModeInsert,
		GroupKey:    &groupA,
	})
	require.NoError(t, err)
	groupB := "group-temperature:point-b"
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:       tagB.ID,
		ConnectorID: connector.ID,
		TableSchema: "main",
		TableName:   "row_group_values",
		ColumnName:  "temperature_c",
		WriteMode:   schema.DatabaseWriteModeInsert,
		GroupKey:    &groupB,
	})
	require.NoError(t, err)

	observedAt := time.Date(2026, 6, 9, 10, 30, 2, 0, time.UTC)
	require.NoError(t, writer.WriteTagValue(ctx, tagA.ID, 21.5, observedAt))
	require.NoError(t, writer.WriteTagValue(ctx, tagB.ID, 23.75, observedAt.Add(time.Second)))

	tickCh <- time.Date(2026, 6, 9, 10, 30, 15, 0, time.UTC)
	require.Eventually(t, func() bool {
		var rowCount int
		queryErr := targetDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM row_group_values`).Scan(&rowCount)
		return queryErr == nil && rowCount == 2
	}, time.Second, 10*time.Millisecond)

	rows, err := targetDB.QueryContext(ctx, `SELECT temperature_c FROM row_group_values ORDER BY id`)
	require.NoError(t, err)
	defer rows.Close()

	values := make([]float64, 0, 2)
	for rows.Next() {
		var value float64
		require.NoError(t, rows.Scan(&value))
		values = append(values, value)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []float64{21.5, 23.75}, values)
}
