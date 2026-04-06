package dbtarget

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestMappingService_DryRun_ReturnsSchemaMissingForMissingTable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "target-missing-table.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.dryrun.schema.missing",
		DisplayName: "DB DryRun Missing Table",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "dry-run-missing-table",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "map-dryrun-missing-table",
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

	result, err := mappingSvc.DryRun(ctx, connector.ID, MappingDryRunRequest{
		CandidateIDs: []string{"map-dryrun-missing-table"},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Results, 1)
	assert.Equal(t, "blocked", result.Results[0].Status)
	assert.Equal(t, "schema_missing", result.Results[0].Code)
	assert.Equal(t, "map-dryrun-missing-table", result.Results[0].CandidateID)
}

func TestMappingService_DryRun_ReturnsConnectorUnavailableWhenTargetCannotOpen(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "missing", "target-unavailable.db")

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.dryrun.connector.unavailable",
		DisplayName: "DB DryRun Connector Unavailable",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "dry-run-connector-unavailable",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "map-dryrun-connector-unavailable",
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

	result, err := mappingSvc.DryRun(ctx, connector.ID, MappingDryRunRequest{
		CandidateIDs: []string{"map-dryrun-connector-unavailable"},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Results, 1)
	assert.Equal(t, "blocked", result.Results[0].Status)
	assert.Equal(t, "connector_unavailable", result.Results[0].Code)
	assert.NotEmpty(t, result.Results[0].Reason)
}

func TestConnectorService_GenerateSchema_DryRunAndExecute(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "target-generate-schema.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	connectorSvc.SetTagReader(tagSvc)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.schema.generate",
		DisplayName: "DB Schema Generate",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "schema-generate",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "map-schema-generate",
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

	dryRunResult, err := connectorSvc.GenerateSchema(ctx, connector.ID, SchemaGenerateRequest{DryRun: true})
	require.NoError(t, err)
	require.NotNil(t, dryRunResult)
	assert.True(t, dryRunResult.DryRun)
	assert.Greater(t, len(dryRunResult.Statements), 0)
	assert.Equal(t, 0, dryRunResult.Executed)

	executeResult, err := connectorSvc.GenerateSchema(ctx, connector.ID, SchemaGenerateRequest{DryRun: false})
	require.NoError(t, err)
	require.NotNil(t, executeResult)
	assert.False(t, executeResult.DryRun)
	assert.Equal(t, len(executeResult.Statements), executeResult.Executed)
	assert.Greater(t, executeResult.Executed, 0)

	dryRunValidation, err := mappingSvc.DryRun(ctx, connector.ID, MappingDryRunRequest{
		CandidateIDs: []string{"map-schema-generate"},
	})
	require.NoError(t, err)
	require.Len(t, dryRunValidation.Results, 1)
	assert.Equal(t, "ready", dryRunValidation.Results[0].Status)
}

func TestConnectorService_ListWriteHistory_ReturnsLatestWriteRecord(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)
	writer := NewWriter(connectorRepo, mappingRepo)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.write.history",
		DisplayName: "DB Write History",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "write-history",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	resetWriteHistory(connector.ID)

	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           tagEntity.ID,
		ConnectorID:     connector.ID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)

	ts := time.Date(2026, 4, 6, 10, 0, 0, 0, time.UTC)
	require.NoError(t, writer.WriteTagValue(ctx, tagEntity.ID, 42.5, ts))
	require.NoError(t, writer.WriteTagValue(ctx, tagEntity.ID, 41.0, ts.Add(time.Minute)))

	records, err := connectorSvc.ListWriteHistory(ctx, connector.ID, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "success", records[0].Status)
	assert.Equal(t, 1, records[0].RowCount)
	assert.Empty(t, records[0].ErrorSummary)
}
