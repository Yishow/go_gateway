package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/diagnostics"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/stretchr/testify/require"
)

type mappingEventSource struct {
	bank     *memory.MemoryBank
	rule     *schema.SourceRule
	mapping  *schema.Mapping
	revision string
}

type mappingEventFixture struct {
	services gatewayServices
	db       *sql.DB
	broker   *diagnostics.Broker
	record   *workspace.Record
	sources  []mappingEventSource
	logPath  string
}

func newMappingEventFixture(t *testing.T, count int, fileLog ...bool) mappingEventFixture {
	t.Helper()
	dbPath := filepath.Join(canonicalTempDir(t), "mapping-events.db")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	options, logPath := diagnostics.Options{}, ""
	if len(fileLog) > 0 && fileLog[0] {
		logPath = filepath.Join(filepath.Dir(dbPath), "logs", "runtime.jsonl")
		options.File, err = diagnostics.OpenFileSink(diagnostics.FileOptions{Path: logPath, DatabasePath: dbPath, DatabaseID: strings.Repeat("a", 64)})
		require.NoError(t, err)
	}
	broker, err := diagnostics.New(options)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, broker.Close(context.Background())) })
	services, err := wireGatewayServicesContext(t.Context(), db, manager, broker)
	require.NoError(t, err)
	f := mappingEventFixture{services: services, db: db, broker: broker, logPath: logPath}
	for i := range count {
		bank := memory.NewMemoryBank(64)
		require.NoError(t, bank.WriteWord(0, 40000))
		server := virtualmodbus.NewServer(bank)
		require.NoError(t, server.Start(0))
		t.Cleanup(func() { require.NoError(t, server.Stop()) })
		config, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1})
		require.NoError(t, err)
		id := fmt.Sprintf("mapping-device-%d", i)
		_, err = db.ExecContext(t.Context(), `INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES(?, 'password_fixture', '', 'modbus_tcp', 'active', ?, '{}')`, id, string(config))
		require.NoError(t, err)
		probe, err := services.device.TestConnectionWithResult(t.Context(), id)
		require.NoError(t, err)
		require.True(t, probe.CanActivate)
		record, err := services.workspace.AttachDevice(t.Context(), id)
		require.NoError(t, err)
		f.record = record
		targetType := schema.DataTypeInt16
		rule, err := services.sourceRule.Create(t.Context(), sourcerule.CreateRuleRequest{ID: fmt.Sprintf("mapping-rule-%d", i), DeviceID: id, StartAddress: "40001", Count: 1, DataType: schema.DataTypeUint16, TargetDataType: &targetType, NamingPrefix: fmt.Sprintf("Mapping%d", i), Enabled: true})
		require.NoError(t, err)
		require.True(t, rule.Enabled)
		view, err := services.sourceRule.GetCandidateView(t.Context(), rule.ID)
		require.NoError(t, err)
		require.Len(t, view.Tags.Candidates, 1)
		candidate, ok := view.Tags.Candidates[0].(*sourcerule.TagCandidateView)
		require.True(t, ok, "native candidate must retain its typed source identity")
		applied, err := services.sourceRule.ApplyTagCandidates(t.Context(), rule.ID, sourcerule.ApplyTagCandidatesRequest{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, RevisionID: rule.RevisionID, CandidateIDs: []string{candidate.ID}})
		require.NoError(t, err)
		require.Len(t, applied.Results, 1)
		mappingRecord, err := services.mapping.GetByID(t.Context(), applied.Results[0].MappingID)
		require.NoError(t, err)
		require.True(t, mappingRecord.Enabled)
		require.Contains(t, mappingRecord.TransformPipeline, `"cast"`)
		require.Contains(t, mappingRecord.TransformPipeline, `"int16"`)
		revision, err := measurement.MappingRevision(*mappingRecord)
		require.NoError(t, err)
		f.sources = append(f.sources, mappingEventSource{bank, rule, mappingRecord, revision})
	}
	f.record.DatabaseSetupRevision = "observed-workspace-revision"
	require.NoError(t, workspace.NewSQLRepository(db).Save(t.Context(), f.record))
	groups, err := services.pollingGroup.List(t.Context())
	require.NoError(t, err)
	for _, group := range groups {
		interval, enabled := 100, true
		_, err := services.pollingGroup.Update(t.Context(), group.ID, pollinggroup.UpdateRequest{IntervalMs: &interval, Enabled: &enabled})
		require.NoError(t, err)
	}
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	return f
}

func (f mappingEventFixture) start(t *testing.T) {
	t.Helper()
	require.NoError(t, f.services.runtime.Start(t.Context()))
}

func (f mappingEventFixture) waitFailure(t *testing.T) {
	t.Helper()
	require.Eventually(t, func() bool { return f.services.runtime.Snapshot().MappingError > 0 }, 3*time.Second, 10*time.Millisecond, "fixture must execute real failing cast, not merely configure a mapping")
}

func (f mappingEventFixture) setValidWord(t *testing.T, i int) {
	t.Helper()
	require.NoError(t, f.sources[i].bank.WriteWord(0, 123))
}

func (f mappingEventFixture) records(t *testing.T) []diagnostics.Event {
	t.Helper()
	return collectorEventRecords(t, f.broker)
}

func findMappingEvent(records []diagnostics.Event, code, mappingID string) *diagnostics.Event {
	for _, event := range records {
		if event.Code == code && event.Fields["mapping_id"] == mappingID {
			matched := event
			return &matched
		}
	}
	return nil
}

func (f *mappingEventFixture) updatePipeline(t *testing.T, i int, pipeline []schema.TransformStep) {
	t.Helper()
	record, err := f.services.mapping.Update(t.Context(), f.sources[i].mapping.ID, mapping.UpdateMappingRequest{TransformPipeline: pipeline, PreviewRawValue: 1})
	require.NoError(t, err)
	revision, err := measurement.MappingRevision(*record)
	require.NoError(t, err)
	f.sources[i].mapping, f.sources[i].revision = record, revision
}

func mappingCastPipeline(target string, extra map[string]any) []schema.TransformStep {
	params := map[string]any{"target_type": target}
	for key, value := range extra {
		params[key] = value
	}
	return []schema.TransformStep{{Type: schema.TransformCast, Order: 0, Params: params}}
}
