package dbtarget

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

func TestProductionGroupOutageRecoveryOpenDestinationHonorsFrozenRevision(t *testing.T) {
	ctx := t.Context()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)
	connectorRepo := NewSQLConnectorRepository(mainDB)
	service := NewConnectorService(connectorRepo, NewSQLTargetMappingRepository(mainDB))
	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: ConnectionConfig{"dsn": targetDSN},
	})
	require.NoError(t, err)

	opened, err := service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorKindSQLite, opened.Kind)
	require.NoError(t, opened.DB.PingContext(ctx))
	require.NoError(t, opened.Close())

	_, err = service.OpenDestination(ctx, connector.ID, "an-older-revision")
	require.ErrorIs(t, err, ErrDestinationBlocked, "an edited connector never receives rows accepted for the old identity")

	_, err = service.OpenDestination(ctx, "missing-connector", "r")
	require.ErrorIs(t, err, ErrDestinationBlocked)

	connector.Enabled = false
	require.NoError(t, connectorRepo.Update(ctx, connector))
	_, err = service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.ErrorIs(t, err, ErrDestinationBlocked, "a disabled destination holds its backlog")
}

func TestProductionGroupOutageRecoveryOpenDestinationUnreachableIsTransient(t *testing.T) {
	ctx := t.Context()
	mainDB := openMigratedTestDB(t)
	connectorRepo := NewSQLConnectorRepository(mainDB)
	service := NewConnectorService(connectorRepo, NewSQLTargetMappingRepository(mainDB))
	// A SQLite connector pointing at a directory that does not exist: the
	// destination is offline, which is a retryable outage, not a blocked identity.
	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name: "offline", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{"dsn": createTargetSQLite(t)},
	})
	require.NoError(t, err)
	require.NoError(t, removeSQLiteFile(connector))

	_, err = service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrDestinationBlocked)
}

func removeSQLiteFile(connector *schema.DatabaseConnector) error {
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return err
	}
	return os.Remove(stringConfigValue(config, "dsn"))
}

func TestProductionGroupOutageRecoveryLegacyWriterSkipsMappingsOwnedByAGroup(t *testing.T) {
	ctx := t.Context()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)
	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)
	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "temperature", DisplayName: "Temperature", DataType: schema.DataTypeFloat64})
	require.NoError(t, err)
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: ConnectionConfig{"dsn": targetDSN},
	})
	require.NoError(t, err)
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID: tagEntity.ID, ConnectorID: connector.ID, TableName: "sensor_values", ColumnName: "value",
		WriteMode: schema.DatabaseWriteModeUpsert, TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)

	owned := false
	writer := NewWriterWithConfig(connectorRepo, mappingRepo, WriterConfig{
		TagReader: tagSvc,
		SuppressMapping: func(connectorID, tagID string) bool {
			return owned && connectorID == connector.ID && tagID == tagEntity.ID
		},
	})
	ts := time.Date(2026, 3, 16, 12, 30, 0, 0, time.UTC)
	rows := func() int {
		db, err := sql.Open("sqlite", targetDSN)
		require.NoError(t, err)
		defer db.Close()
		var n int
		require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values`).Scan(&n))
		return n
	}

	require.NoError(t, writer.WriteTagValue(ctx, tagEntity.ID, 1.5, ts))
	require.Equal(t, 1, rows(), "before a group owns the output the legacy writer works as before")

	owned = true
	require.ErrorIs(t, writer.WriteTagValue(ctx, tagEntity.ID, 2.5, ts.Add(time.Minute)), ErrOutputOwnedByWriteGroup)
	require.Equal(t, 1, rows(), "once a group owns the output the legacy writer stops, so there is exactly one writer")

	otherTag, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "pressure", DisplayName: "Pressure", DataType: schema.DataTypeFloat64})
	require.NoError(t, err)
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID: otherTag.ID, ConnectorID: connector.ID, TableName: "sensor_values", ColumnName: "value",
		WriteMode: schema.DatabaseWriteModeUpsert, TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)
	require.NoError(t, writer.WriteTagValue(ctx, otherTag.ID, 9.5, ts.Add(2*time.Minute)))
	require.Equal(t, 2, rows(), "tags the group does not own are unaffected")
}

func TestRevisionBoundBacklogLegacyWriterReportsOwnedOutputAsNotApplicable(t *testing.T) {
	ctx := t.Context()
	mainDB := openMigratedTestDB(t)
	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)
	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "temperature", DisplayName: "Temperature", DataType: schema.DataTypeFloat64})
	require.NoError(t, err)
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: ConnectionConfig{"dsn": createTargetSQLite(t)},
	})
	require.NoError(t, err)
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID: tagEntity.ID, ConnectorID: connector.ID, TableName: "sensor_values", ColumnName: "value",
		WriteMode: schema.DatabaseWriteModeUpsert, TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)
	writer := NewWriterWithConfig(connectorRepo, mappingRepo, WriterConfig{
		TagReader:       tagSvc,
		SuppressMapping: func(string, string) bool { return true },
	})

	err = writer.WriteTagValue(ctx, tagEntity.ID, 1.5, time.Now())
	require.ErrorIs(t, err, ErrOutputOwnedByWriteGroup, "a write that never happened must not look like success (nil)")
}

func TestRevisionBoundBacklogCredentialRejectionBlocksOnlyThatDestination(t *testing.T) {
	ctx := t.Context()
	mainDB := openMigratedTestDB(t)
	connectorRepo := NewSQLConnectorRepository(mainDB)
	service := NewConnectorService(connectorRepo, NewSQLTargetMappingRepository(mainDB))
	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: ConnectionConfig{"dsn": createTargetSQLite(t)},
	})
	require.NoError(t, err)

	original := openExternalDBManagerFunc
	t.Cleanup(func() { openExternalDBManagerFunc = original })
	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		return nil, errors.New("pq: password authentication failed for user \"gateway\"")
	}
	_, err = service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.ErrorIs(t, err, ErrDestinationBlocked, "a rejected credential keeps the backlog blocked for the same identity")
	require.NotContains(t, err.Error(), "gateway", "the blocked error never carries the driver text")

	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		return nil, errors.New("dial tcp 10.0.0.5:5432: i/o timeout")
	}
	_, err = service.OpenDestination(ctx, connector.ID, connector.IdentityRevision)
	require.ErrorIs(t, err, ErrDestinationUnreachable)
	require.NotErrorIs(t, err, ErrDestinationBlocked)
}
