package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type dbTargetToolingFixture struct {
	router        *gin.Engine
	connectorRepo *dbtarget.SQLConnectorRepository
	connectorSvc  *dbtarget.ConnectorService
	mappingSvc    *dbtarget.MappingService
	writer        *dbtarget.Writer
	mappingRepo   *dbtarget.SQLTargetMappingRepository
	tagSvc        *tag.Service
}

func setupDBTargetToolingFixture(t *testing.T) *dbTargetToolingFixture {
	t.Helper()

	gin.SetMode(gin.TestMode)
	mainDB := openMigratedDBTargetMainDB(t)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := dbtarget.NewSQLConnectorRepository(mainDB)
	mappingRepo := dbtarget.NewSQLTargetMappingRepository(mainDB)
	connectorSvc := dbtarget.NewConnectorService(connectorRepo, mappingRepo)
	connectorSvc.SetTagReader(tagSvc)
	mappingSvc := dbtarget.NewMappingService(mappingRepo, connectorRepo, tagSvc)
	writer := dbtarget.NewWriter(connectorRepo, mappingRepo)

	handler := NewDatabaseTargetHandler(connectorSvc, mappingSvc)
	router := gin.Default()
	router.POST("/datalink/db-targets/connectors/:id/schema/generate", handler.GenerateSchema)
	router.POST("/datalink/db-targets/connectors/:id/mappings/dry-run", handler.DryRunMappings)
	router.GET("/datalink/db-targets/connectors/:id/write-history", handler.ListWriteHistory)

	return &dbTargetToolingFixture{
		router:        router,
		connectorRepo: connectorRepo,
		connectorSvc:  connectorSvc,
		mappingSvc:    mappingSvc,
		writer:        writer,
		mappingRepo:   mappingRepo,
		tagSvc:        tagSvc,
	}
}

func TestDatabaseTargetHandler_DryRunMappings_ReturnsSchemaMissing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fixture := setupDBTargetToolingFixture(t)
	targetDSN := filepath.Join(t.TempDir(), "target-handler-dryrun-missing.db")
	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	connector, err := fixture.connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name: "handler-dryrun-schema-missing",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	tagEntity, err := fixture.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "handler.dryrun.schema.missing",
		DisplayName: "Handler DryRun Missing",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	now := time.Now()
	require.NoError(t, fixture.mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "handler-map-schema-missing",
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

	body := []byte(`{"candidate_ids":["handler-map-schema-missing"]}`)
	req, err := http.NewRequest(
		http.MethodPost,
		"/datalink/db-targets/connectors/"+connector.ID+"/mappings/dry-run",
		bytes.NewBuffer(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, true, payload["success"])
	data := payload["data"].(map[string]any)
	results := data["results"].([]any)
	require.Len(t, results, 1)
	assert.Equal(t, "blocked", results[0].(map[string]any)["status"])
	assert.Equal(t, "schema_missing", results[0].(map[string]any)["code"])
}

func TestDatabaseTargetHandler_GenerateSchema_DryRunAndExecute(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fixture := setupDBTargetToolingFixture(t)
	targetDSN := filepath.Join(t.TempDir(), "target-handler-generate.db")
	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	connector, err := fixture.connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name: "handler-generate-schema",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	tagEntity, err := fixture.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "handler.generate.schema",
		DisplayName: "Handler Generate Schema",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	now := time.Now()
	require.NoError(t, fixture.mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:          "handler-map-generate",
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

	dryRunReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/db-targets/connectors/"+connector.ID+"/schema/generate",
		bytes.NewBufferString(`{"dry_run":true}`),
	)
	require.NoError(t, err)
	dryRunReq.Header.Set("Content-Type", "application/json")
	dryRunResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(dryRunResp, dryRunReq)
	require.Equal(t, http.StatusOK, dryRunResp.Code)

	var dryRunPayload map[string]any
	require.NoError(t, json.Unmarshal(dryRunResp.Body.Bytes(), &dryRunPayload))
	assert.Equal(t, true, dryRunPayload["success"])
	dryRunData := dryRunPayload["data"].(map[string]any)
	assert.Equal(t, true, dryRunData["dry_run"])
	assert.Greater(t, len(dryRunData["statements"].([]any)), 0)

	execReq, err := http.NewRequest(
		http.MethodPost,
		"/datalink/db-targets/connectors/"+connector.ID+"/schema/generate",
		bytes.NewBufferString(`{"dry_run":false}`),
	)
	require.NoError(t, err)
	execReq.Header.Set("Content-Type", "application/json")
	execResp := httptest.NewRecorder()
	fixture.router.ServeHTTP(execResp, execReq)
	require.Equal(t, http.StatusOK, execResp.Code)

	var execPayload map[string]any
	require.NoError(t, json.Unmarshal(execResp.Body.Bytes(), &execPayload))
	assert.Equal(t, true, execPayload["success"])
	execData := execPayload["data"].(map[string]any)
	assert.Equal(t, false, execData["dry_run"])
	assert.Greater(t, int(execData["executed"].(float64)), 0)
}

func TestDatabaseTargetHandler_ListWriteHistory_ReturnsLatestRecords(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fixture := setupDBTargetToolingFixture(t)
	targetDSN := filepath.Join(t.TempDir(), "target-handler-history.db")
	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	_, err = targetDB.Exec(`
		CREATE TABLE sensor_values (
			ts DATETIME PRIMARY KEY,
			value REAL NOT NULL
		)
	`)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	connector, err := fixture.connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name: "handler-write-history",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	tagEntity, err := fixture.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "handler.write.history",
		DisplayName: "Handler Write History",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	_, err = fixture.mappingSvc.Create(ctx, dbtarget.CreateTargetMappingRequest{
		TagID:           tagEntity.ID,
		ConnectorID:     connector.ID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtrForDBTarget("ts"),
	})
	require.NoError(t, err)

	ts := time.Date(2026, 4, 6, 11, 0, 0, 0, time.UTC)
	require.NoError(t, fixture.writer.WriteTagValue(ctx, tagEntity.ID, 42.5, ts))

	req, err := http.NewRequest(
		http.MethodGet,
		"/datalink/db-targets/connectors/"+connector.ID+"/write-history?limit=1",
		nil,
	)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, true, payload["success"])
	data := payload["data"].(map[string]any)
	records := data["records"].([]any)
	require.Len(t, records, 1)
	assert.Equal(t, "success", records[0].(map[string]any)["status"])
	assert.Equal(t, float64(1), records[0].(map[string]any)["row_count"])
}

func TestDatabaseTargetHandler_ListWriteHistory_ReturnsGroupedFlushMetadata(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fixture := setupDBTargetToolingFixture(t)
	targetDSN := filepath.Join(t.TempDir(), "target-handler-grouped-history.db")
	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	_, err = targetDB.Exec(`
		CREATE TABLE meter_rows (
			ts DATETIME PRIMARY KEY,
			a1 REAL,
			kw REAL
		)
	`)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	tickCh := make(chan time.Time, 2)

	connector, err := fixture.connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name: "handler-grouped-history",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{
			"dsn": targetDSN,
		},
	})
	require.NoError(t, err)

	tagA1, err := fixture.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "meter/A1",
		DisplayName: "Meter A1",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	tagKW, err := fixture.tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "meter/kw",
		DisplayName: "Meter kW",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	groupKey := "meter"
	_, err = fixture.mappingSvc.Create(ctx, dbtarget.CreateTargetMappingRequest{
		TagID:           tagA1.ID,
		ConnectorID:     connector.ID,
		TableName:       "meter_rows",
		ColumnName:      "a1",
		GroupKey:        &groupKey,
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtrForDBTarget("ts"),
	})
	require.NoError(t, err)
	_, err = fixture.mappingSvc.Create(ctx, dbtarget.CreateTargetMappingRequest{
		TagID:           tagKW.ID,
		ConnectorID:     connector.ID,
		TableName:       "meter_rows",
		ColumnName:      "kw",
		GroupKey:        &groupKey,
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: stringPtrForDBTarget("ts"),
	})
	require.NoError(t, err)

	writer := dbtarget.NewWriterWithConfig(
		fixture.connectorRepo,
		fixture.mappingRepo,
		dbtarget.WriterConfig{FlushTick: tickCh},
	)
	t.Cleanup(func() {
		require.NoError(t, writer.Close(context.Background()))
	})

	observedAt := time.Date(2026, 4, 6, 11, 0, 5, 0, time.UTC)
	require.NoError(t, writer.WriteTagValue(ctx, tagA1.ID, 42.5, observedAt))
	require.NoError(t, writer.WriteTagValue(ctx, tagKW.ID, 7.25, observedAt.Add(3*time.Second)))
	tickCh <- time.Date(2026, 4, 6, 11, 0, 15, 0, time.UTC)

	require.Eventually(t, func() bool {
		records, listErr := fixture.connectorSvc.ListWriteHistory(ctx, connector.ID, 1)
		return listErr == nil && len(records) == 1
	}, time.Second, 10*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodGet,
		"/datalink/db-targets/connectors/"+connector.ID+"/write-history?limit=1",
		nil,
	)
	require.NoError(t, err)
	resp := httptest.NewRecorder()
	fixture.router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	records := payload["data"].(map[string]any)["records"].([]any)
	require.Len(t, records, 1)
	record := records[0].(map[string]any)
	assert.Equal(t, "success", record["status"])
	assert.Equal(t, float64(1), record["row_count"])
	assert.Equal(t, "meter", record["group_key"])
	assert.Equal(t, "meter_rows", record["table_name"])
	assert.Equal(t, float64(15), record["effective_interval_seconds"])
	assert.Equal(t, "2026-04-06T11:00:00Z", record["observed_at"])
}

func openMigratedDBTargetMainDB(t *testing.T) *sql.DB {
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

func stringPtrForDBTarget(value string) *string {
	return &value
}
