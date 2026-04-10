package dbtarget

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestConnectorService_CreateGroupingDefaultsWriteInterval(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	connectorRepo := NewSQLConnectorRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "grouping-default-interval",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 15, connector.DefaultWriteIntervalSeconds)

	reloaded, err := connectorRepo.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, 15, reloaded.DefaultWriteIntervalSeconds)
}

func TestSQLTargetMappingRepository_GroupingRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo)

	groupedTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "meter/A1",
		DisplayName: "Meter A1",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	legacyTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "legacy_value",
		DisplayName: "Legacy Value",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "grouping-roundtrip",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	now := time.Now()
	grouped := &schema.DatabaseTargetMapping{
		ID:                   "map-grouped",
		TagID:                groupedTag.ID,
		ConnectorID:          connector.ID,
		TableSchema:          "main",
		TableName:            "meter_rows",
		ColumnName:           "a1",
		WriteMode:            schema.DatabaseWriteModeUpsert,
		TimestampColumn:      stringPtr("ts"),
		GroupKey:             stringPtr("meter"),
		WriteIntervalSeconds: intPtr(15),
		Enabled:              true,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	legacy := &schema.DatabaseTargetMapping{
		ID:              "map-legacy",
		TagID:           legacyTag.ID,
		ConnectorID:     connector.ID,
		TableSchema:     "main",
		TableName:       "legacy_rows",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeInsert,
		TimestampColumn: nil,
		GroupKey:        nil,
		Enabled:         true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	require.NoError(t, mappingRepo.Create(ctx, grouped))
	require.NoError(t, mappingRepo.Create(ctx, legacy))

	reloadedGrouped, err := mappingRepo.GetByID(ctx, grouped.ID)
	require.NoError(t, err)
	require.NotNil(t, reloadedGrouped.GroupKey)
	require.Equal(t, "meter", *reloadedGrouped.GroupKey)
	require.NotNil(t, reloadedGrouped.WriteIntervalSeconds)
	require.Equal(t, 15, *reloadedGrouped.WriteIntervalSeconds)

	reloadedLegacy, err := mappingRepo.GetByID(ctx, legacy.ID)
	require.NoError(t, err)
	require.Nil(t, reloadedLegacy.GroupKey)
	require.Nil(t, reloadedLegacy.WriteIntervalSeconds)

	mappings, err := mappingRepo.List(ctx, TargetMappingListFilter{ConnectorID: &connector.ID})
	require.NoError(t, err)
	require.Len(t, mappings, 2)
}

func intPtr(value int) *int {
	return &value
}
