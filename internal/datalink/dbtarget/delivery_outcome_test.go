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
)

func TestConnectorService_GenerateSchemaRecordsLastSchemaOutcome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "schema-outcome.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	connectorSvc.SetTagReader(tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.outcome.schema",
		DisplayName: "DB Outcome Schema",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "schema-outcome",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "map-schema-outcome",
		TagID:       tagEntity.ID,
		ConnectorID: connector.ID,
		TableSchema: "main",
		TableName:   "measurements",
		ColumnName:  "line_a",
		WriteMode:   schema.DatabaseWriteModeInsert,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}))

	_, err = connectorSvc.GenerateSchema(ctx, connector.ID, SchemaGenerateRequest{DryRun: false})
	require.NoError(t, err)

	updated, err := connectorRepo.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, "success", updated.LastSchemaEnsureStatus)
	require.NotNil(t, updated.LastSchemaEnsureAt)
	require.Empty(t, updated.LastSchemaEnsureError)
}

func TestConnectorService_GenerateSchemaRecordsLastSchemaFailureOutcome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "missing", "schema-outcome.db")

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	connectorSvc.SetTagReader(tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.outcome.schema.failed",
		DisplayName: "DB Outcome Schema Failed",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "schema-outcome-failed",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "map-schema-outcome-failed",
		TagID:       tagEntity.ID,
		ConnectorID: connector.ID,
		TableSchema: "main",
		TableName:   "measurements",
		ColumnName:  "line_a",
		WriteMode:   schema.DatabaseWriteModeInsert,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}))

	_, err = connectorSvc.GenerateSchema(ctx, connector.ID, SchemaGenerateRequest{DryRun: false})
	require.Error(t, err)

	updated, err := connectorRepo.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, "failed", updated.LastSchemaEnsureStatus)
	require.NotNil(t, updated.LastSchemaEnsureAt)
	require.Contains(t, updated.LastSchemaEnsureError, "檢查資料庫目標表結構失敗")
}

func TestWriter_RecordsLastWriteFailureOutcome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "write-outcome.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	writer := NewWriterWithConfig(
		connectorRepo,
		mappingRepo,
		WriterConfig{FlushTick: make(chan time.Time), TagReader: tagSvc},
	)
	t.Cleanup(func() {
		_ = writer.Close(context.Background())
	})

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.outcome.write",
		DisplayName: "DB Outcome Write",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "write-outcome",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "map-write-outcome",
		TagID:       tagEntity.ID,
		ConnectorID: connector.ID,
		TableSchema: "main",
		TableName:   "missing_measurements",
		ColumnName:  "line_a",
		WriteMode:   schema.DatabaseWriteModeInsert,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}))

	observedAt := time.Date(2026, 6, 7, 13, 0, 0, 0, time.UTC)
	err = writer.WriteTagValue(ctx, tagEntity.ID, 42.5, observedAt)
	require.Error(t, err)

	updated, err := connectorRepo.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, "failed", updated.LastWriteStatus)
	require.NotNil(t, updated.LastWriteAt)
	require.Equal(t, observedAt, updated.LastWriteAt.UTC())
	require.Contains(t, updated.LastWriteError, "執行資料庫寫入失敗")
}

func TestWriter_RecordsLastFlushFailureOutcome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "flush-outcome.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	writer := NewWriterWithConfig(
		connectorRepo,
		mappingRepo,
		WriterConfig{FlushTick: make(chan time.Time), TagReader: tagSvc},
	)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.outcome.flush",
		DisplayName: "DB Outcome Flush",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "flush-outcome",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	groupKey := "meter"
	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "map-flush-outcome",
		TagID:       tagEntity.ID,
		ConnectorID: connector.ID,
		TableSchema: "main",
		TableName:   "missing_grouped_measurements",
		ColumnName:  "line_a",
		WriteMode:   schema.DatabaseWriteModeInsert,
		GroupKey:    &groupKey,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}))

	observedAt := time.Date(2026, 6, 7, 14, 0, 0, 0, time.UTC)
	require.NoError(t, writer.WriteTagValue(ctx, tagEntity.ID, 42.5, observedAt))
	require.Error(t, writer.Close(ctx))

	updated, err := connectorRepo.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, "failed", updated.LastFlushStatus)
	require.NotNil(t, updated.LastFlushAt)
	require.Equal(t, observedAt, updated.LastFlushAt.UTC())
	require.Contains(t, updated.LastFlushError, "執行資料庫分組寫入失敗")
}
