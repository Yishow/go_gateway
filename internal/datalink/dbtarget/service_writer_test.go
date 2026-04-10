package dbtarget

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestMappingService_CreateValidateAndWriter(t *testing.T) {
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
		Key:         "temperature",
		DisplayName: "Temperature",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorStatusReady, connector.Status)

	mappingEntity, err := mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           tagEntity.ID,
		ConnectorID:     connector.ID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)
	require.Equal(t, "sensor_values", mappingEntity.TableName)
	require.Equal(t, "value", mappingEntity.ColumnName)

	validation, err := mappingSvc.Validate(ctx, connector.ID)
	require.NoError(t, err)
	require.True(t, validation.Ready)
	require.Empty(t, validation.Issues)

	ts := time.Date(2026, 3, 16, 12, 30, 0, 0, time.UTC)
	require.NoError(t, writer.WriteTagValue(ctx, tagEntity.ID, 42.5, ts))
	require.NoError(t, writer.WriteTagValue(ctx, tagEntity.ID, 41.0, ts))

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	defer targetDB.Close()

	var value float64
	err = targetDB.QueryRowContext(ctx, `SELECT value FROM sensor_values LIMIT 1`).Scan(&value)
	require.NoError(t, err)
	var rowCount int
	err = targetDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM sensor_values`).Scan(&rowCount)
	require.NoError(t, err)
	require.Equal(t, 41.0, value)
	require.Equal(t, 1, rowCount)
}

func TestMappingService_CreateRejectsUpsertWithoutUniqueTimestamp(t *testing.T) {
	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "target-non-unique.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	defer targetDB.Close()

	_, err = targetDB.Exec(`
		CREATE TABLE sensor_values (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ts DATETIME NOT NULL,
			value REAL NOT NULL
		)
	`)
	require.NoError(t, err)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "flow-rate",
		DisplayName: "Flow Rate",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-non-unique",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           tagEntity.ID,
		ConnectorID:     connector.ID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "primary key 或 single-column unique")
}

func TestWriter_GroupedMappingsFlushOneRowPerBucketOnTimer(t *testing.T) {
	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "target-grouped-timer.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	defer targetDB.Close()
	_, err = targetDB.Exec(`
		CREATE TABLE meter_rows (
			ts DATETIME PRIMARY KEY,
			a1 REAL,
			kw REAL
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
		Name: "grouped-timer",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
		DefaultWriteIntervalSeconds: &interval,
	})
	require.NoError(t, err)
	resetWriteHistory(connector.ID)

	tagA1, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "meter/A1",
		DisplayName: "Meter A1",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	tagKW, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "meter/kw",
		DisplayName: "Meter kW",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	groupKey := "meter"
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           tagA1.ID,
		ConnectorID:     connector.ID,
		TableName:       "meter_rows",
		ColumnName:      "a1",
		GroupKey:        &groupKey,
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           tagKW.ID,
		ConnectorID:     connector.ID,
		TableName:       "meter_rows",
		ColumnName:      "kw",
		GroupKey:        &groupKey,
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)

	firstObservedAt := time.Date(2026, 4, 6, 11, 0, 3, 0, time.UTC)
	secondObservedAt := time.Date(2026, 4, 6, 11, 0, 10, 0, time.UTC)
	bucketStart := time.Date(2026, 4, 6, 11, 0, 0, 0, time.UTC)

	require.NoError(t, writer.WriteTagValue(ctx, tagA1.ID, 42.5, firstObservedAt))
	require.NoError(t, writer.WriteTagValue(ctx, tagKW.ID, 7.25, secondObservedAt))

	var rowCount int
	err = targetDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM meter_rows`).Scan(&rowCount)
	require.NoError(t, err)
	require.Equal(t, 0, rowCount)

	tickCh <- time.Date(2026, 4, 6, 11, 0, 15, 0, time.UTC)
	require.Eventually(t, func() bool {
		queryErr := targetDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM meter_rows`).Scan(&rowCount)
		return queryErr == nil && rowCount == 1
	}, time.Second, 10*time.Millisecond)

	var storedTS time.Time
	var a1Value float64
	var kwValue float64
	err = targetDB.QueryRowContext(ctx, `SELECT ts, a1, kw FROM meter_rows LIMIT 1`).Scan(&storedTS, &a1Value, &kwValue)
	require.NoError(t, err)
	assert.Equal(t, bucketStart, storedTS.UTC())
	assert.Equal(t, 42.5, a1Value)
	assert.Equal(t, 7.25, kwValue)

	records, err := connectorSvc.ListWriteHistory(ctx, connector.ID, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "success", records[0].Status)
	assert.Equal(t, 1, records[0].RowCount)
	assert.Equal(t, bucketStart, records[0].ObservedAt.UTC())
	require.NotNil(t, records[0].GroupKey)
	assert.Equal(t, "meter", *records[0].GroupKey)
	assert.Equal(t, "meter_rows", records[0].TableName)
	assert.Equal(t, 15, records[0].EffectiveIntervalSeconds)
}

func TestWriter_CloseFlushesPendingGroupedBuckets(t *testing.T) {
	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "target-grouped-close.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	defer targetDB.Close()
	_, err = targetDB.Exec(`
		CREATE TABLE meter_rows (
			ts DATETIME PRIMARY KEY,
			a1 REAL
		)
	`)
	require.NoError(t, err)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)
	writer := NewWriterWithConfig(connectorRepo, mappingRepo, WriterConfig{})

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "grouped-close",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)
	resetWriteHistory(connector.ID)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "meter/A1",
		DisplayName: "Meter A1",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	groupKey := "meter"
	_, err = mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           tagEntity.ID,
		ConnectorID:     connector.ID,
		TableName:       "meter_rows",
		ColumnName:      "a1",
		GroupKey:        &groupKey,
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)

	observedAt := time.Date(2026, 4, 6, 12, 5, 4, 0, time.UTC)
	require.NoError(t, writer.WriteTagValue(ctx, tagEntity.ID, 18.5, observedAt))

	var rowCount int
	err = targetDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM meter_rows`).Scan(&rowCount)
	require.NoError(t, err)
	require.Equal(t, 0, rowCount)

	require.NoError(t, writer.Close(ctx))

	err = targetDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM meter_rows`).Scan(&rowCount)
	require.NoError(t, err)
	require.Equal(t, 1, rowCount)

	records, err := connectorSvc.ListWriteHistory(ctx, connector.ID, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "success", records[0].Status)
	require.NotNil(t, records[0].GroupKey)
	assert.Equal(t, "meter", *records[0].GroupKey)
	assert.Equal(t, "meter_rows", records[0].TableName)
}

func TestMappingService_UpdateClearsTimestampWhenSwitchingToInsert(t *testing.T) {
	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := NewMappingService(mappingRepo, connectorRepo, tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "pressure",
		DisplayName: "Pressure",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	mappingEntity, err := mappingSvc.Create(ctx, CreateTargetMappingRequest{
		TagID:           tagEntity.ID,
		ConnectorID:     connector.ID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtr("ts"),
	})
	require.NoError(t, err)
	require.NotNil(t, mappingEntity.TimestampColumn)

	insertMode := schema.DatabaseWriteModeInsert
	updatedMapping, err := mappingSvc.Update(ctx, mappingEntity.ID, UpdateTargetMappingRequest{
		WriteMode: &insertMode,
	})
	require.NoError(t, err)
	require.Equal(t, schema.DatabaseWriteModeInsert, updatedMapping.WriteMode)
	require.Nil(t, updatedMapping.TimestampColumn)

	reloadedMapping, err := mappingSvc.GetByID(ctx, mappingEntity.ID)
	require.NoError(t, err)
	require.Nil(t, reloadedMapping.TimestampColumn)
}

func TestBuildWriteStatement_InsertIgnoresTimestampColumn(t *testing.T) {
	observedAt := time.Date(2026, 3, 16, 12, 45, 0, 0, time.UTC)
	mapping := &schema.DatabaseTargetMapping{
		TableSchema:     "main",
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeInsert,
		TimestampColumn: stringPtr("ts"),
	}

	query, args, err := buildWriteStatement(
		schema.DatabaseConnectorKindSQLite,
		mapping,
		42.5,
		observedAt,
	)
	require.NoError(t, err)
	require.Equal(t, `INSERT INTO "sensor_values" ("value") VALUES (?)`, query)
	require.Len(t, args, 1)
	require.Equal(t, 42.5, args[0])
}

func TestConnectorService_UpdatePreservesPasswordWhenRedactedValueIsSubmitted(t *testing.T) {
	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	connectorRepo := NewSQLConnectorRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn":      targetDSN,
			"password": "secret",
		},
	})
	require.NoError(t, err)

	updatedConnector, err := connectorSvc.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{
			"dsn":      targetDSN,
			"password": "",
		},
	})
	require.NoError(t, err)

	connectionConfig, err := parseConnectionConfig(updatedConnector.ConnectionConfig)
	require.NoError(t, err)
	require.Equal(t, "secret", stringConfigValue(connectionConfig, "password"))
}

func TestConnectorService_UpdateAllowsClearingPassword(t *testing.T) {
	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := createTargetSQLite(t)

	connectorRepo := NewSQLConnectorRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo)

	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name: "sqlite-target",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{
			"dsn":      targetDSN,
			"password": "secret",
		},
	})
	require.NoError(t, err)

	clearPassword := true
	updatedConnector, err := connectorSvc.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{
			"dsn":      targetDSN,
			"password": "",
		},
		ClearPassword: &clearPassword,
	})
	require.NoError(t, err)

	connectionConfig, err := parseConnectionConfig(updatedConnector.ConnectionConfig)
	require.NoError(t, err)
	require.Empty(t, stringConfigValue(connectionConfig, "password"))
	_, hasPassword := connectionConfig["password"]
	require.False(t, hasPassword)
}

func openMigratedTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "main.db")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	migrator := datalinkbase.NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	return db
}

func createTargetSQLite(t *testing.T) string {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "target.db")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE sensor_values (
			ts DATETIME PRIMARY KEY,
			value REAL NOT NULL
		)
	`)
	require.NoError(t, err)

	return dbPath
}

func stringPtr(value string) *string {
	return &value
}
