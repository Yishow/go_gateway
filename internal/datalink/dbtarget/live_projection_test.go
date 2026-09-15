package dbtarget

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestMappingService_ValidateIgnoresStaleMissingTagMappings(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)

	liveTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.live.validate",
		DisplayName: "DB Live Validate",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "live-projection-validate",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           liveTag.ID,
		ConnectorID:     connector.ID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)
	createStaleTargetMapping(ctx, t, mappingRepo, connector.ID)

	validation, err := mappingSvc.Validate(ctx, connector.ID)
	require.NoError(t, err)
	require.True(t, validation.Ready)
	for _, issue := range validation.Issues {
		require.NotEqual(t, "tag_missing", issue.Code)
		require.NotEqual(t, "tag-old", issue.TagID)
	}
}

func TestMappingService_DryRunExcludesStaleMissingTagMappingsByDefault(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)

	liveTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.live.dryrun",
		DisplayName: "DB Live DryRun",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "live-projection-dryrun",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	liveMapping, err := mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           liveTag.ID,
		ConnectorID:     connector.ID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)
	createStaleTargetMapping(ctx, t, mappingRepo, connector.ID)

	result, err := mappingSvc.DryRun(ctx, connector.ID, MappingDryRunRequest{})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Results, 1)
	require.Equal(t, liveMapping.ID, result.Results[0].MappingID)
	require.Equal(t, "ready", result.Results[0].Status)
}

func TestWriter_IgnoresStaleMissingTagMappings(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	tickCh := make(chan time.Time)
	writer := NewWriterWithConfig(
		connectorRepo,
		mappingRepo,
		WriterConfig{FlushTick: tickCh, TagReader: tagSvc},
	)
	t.Cleanup(func() {
		require.NoError(t, writer.Close(context.Background()))
	})

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "live-projection-writer",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	createStaleTargetMapping(ctx, t, mappingRepo, connector.ID)

	observedAt := time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC)
	require.NoError(t, writer.WriteTagValue(ctx, "tag-old", 99.5, observedAt))
}

func createStaleTargetMapping(
	ctx context.Context,
	t *testing.T,
	mappingRepo TargetMappingRepository,
	connectorID string,
) {
	t.Helper()

	now := time.Now()
	err := mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:              "stale-map-tag-old",
		TagID:           "tag-old",
		ConnectorID:     connectorID,
		TableSchema:     "main",
		TableName:       "stale_values",
		ColumnName:      "old_value",
		WriteMode:       schema.DatabaseWriteModeInsert,
		TimestampColumn: nil,
		Enabled:         true,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	require.NoError(t, err)
}
