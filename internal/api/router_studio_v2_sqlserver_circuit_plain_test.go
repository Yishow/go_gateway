package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const (
	sqlServerCircuitMeasurementID        = "meas-sql-counter"
	sqlServerCircuitIntegerMeasurementID = "meas-sql-counter-integer"
	sqlServerCircuitWindow               = 10 * time.Minute
)

type sqlServerCircuitInput struct {
	sampleID string
	offset   time.Duration
	value    uint64
}

var sqlServerCircuitInputs = [...]sqlServerCircuitInput{
	{sampleID: "sql-baseline", offset: 9 * time.Minute, value: 120000},
	{sampleID: "sql-after-gap", offset: 29 * time.Minute, value: 120004},
	{sampleID: "sql-reset", offset: 59 * time.Minute, value: 7},
}

var sqlServerCircuitIntegerInputs = [...]sqlServerCircuitInput{
	{sampleID: "sql-2^53", offset: 69 * time.Minute, value: 9007199254740992},
	{sampleID: "sql-2^53-plus-1", offset: 79 * time.Minute, value: 9007199254740993},
}

type sqlServerCircuitPlainFixture struct {
	managedGroupFixture
	fixture             *SQLServerCircuitFixture
	target              *sql.DB
	pipe                *grouppipeline.Pipeline
	clock               *atomic.Int64
	base                time.Time
	measurementID       string
	database            string
	instantID           string
	cumulativeID        string
	integerCumulativeID string
	targetService       *dbtarget.ConnectorService
	tags                *tag.Service
	store               *groupdelivery.Store
	pipeConfig          grouppipeline.Config
	createdTable        bool
}

// TestSQLServerCircuitPlainIntegration is the real plain SQL Server entrypoint.
// It stays NOT_RUN unless GW_SQLSERVER_TEST_FIXTURE names a validated
// disposable manifest and the production managed-long operation gate is open.
func TestSQLServerCircuitPlainIntegration(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	f := newSQLServerCircuitPlainFixture(t, fixture)
	runSQLServerCircuitPlainIntegration(t, f)
}

// NewSQLServerCircuitPlainFixture is the reusable plain/native setup entrypoint
// for GLM/root consumers. It keeps the closed managed-long gate: a real probe
// and capability validation happen before schema DDL, and a closed gate skips
// truthfully before writes.
func NewSQLServerCircuitPlainFixture(t *testing.T, fixture *SQLServerCircuitFixture, options ...SQLServerCircuitFixtureOptions) *sqlServerCircuitPlainFixture {
	return newSQLServerCircuitPlainFixture(t, fixture, options...)
}

func newSQLServerCircuitPlainFixture(t *testing.T, fixture *SQLServerCircuitFixture, options ...SQLServerCircuitFixtureOptions) *sqlServerCircuitPlainFixture {
	t.Helper()
	var option SQLServerCircuitFixtureOptions
	if len(options) > 0 {
		option = options[0]
	}
	mode := option.StorageMode
	if mode == "" {
		mode = workspace.WriteGroupStorageModePlain
	}
	tagType := option.TagType
	if tagType == "" {
		tagType = "uint64"
	}
	f := newManagedDestinationFixture(t, "sqlite", tagType)
	ctx := t.Context()
	_, err := f.db.ExecContext(ctx, `UPDATE database_connectors SET kind = 'sqlserver', connection_config = ?, identity_revision = 'connector-1' WHERE id = 'connector-A'`,
		fixture.ConnectionConfigJSON())
	require.NoError(t, err)
	f.record, err = f.workspace.BindDatabaseConnector(ctx, "connector-A")
	require.NoError(t, err)
	f.namespace = fixture.Manifest.Schema
	seedProfileSource(t, f.writeGroupRouterFixture, "integer")
	_, err = f.db.ExecContext(ctx, `UPDATE tags SET data_type='uint64' WHERE id='tag-integer'`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(ctx, `UPDATE points SET data_type='uint64' WHERE id='point-integer'`)
	require.NoError(t, err)

	// Probe and gate-check before any schema DDL. The capability is derived
	// from the saved connector, never from the manifest alone.
	targetService := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	capability, err := targetService.ProbeSQLServerRecordingAtRevision(ctx, "connector-A", "connector-1")
	require.NoError(t, err)
	require.NoError(t, fixture.ValidateCapability(capability))
	checkSQLServerCircuitGate(t, mode, capability)
	opened, err := targetService.OpenDestination(ctx, "connector-A", "connector-1")
	require.NoError(t, err)
	f.target = opened.DB
	if opened.Close != nil {
		t.Cleanup(func() { require.NoError(t, opened.Close()) })
	}

	_, err = f.db.ExecContext(ctx, `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
			series_epoch, name, quantity, unit, semantic_kind, numeric_encoding, exact_recording_scale
		) VALUES (?, ?, 'device-A', 'point-A', 'tag-A', 'device-A', 'def-sql', 'bind-sql', 'epoch-sql', 'SQL Counter', 'energy', 'kWh', 'counter', 'uint64', '0.1')`,
		sqlServerCircuitMeasurementID, f.record.ID)
	require.NoError(t, err)
	definition, err := measurement.NewSQLRepository(f.db).GetByID(ctx, sqlServerCircuitMeasurementID)
	require.NoError(t, err)
	require.Equal(t, "0.1", definition.ExactRecordingScale)
	_, err = f.db.ExecContext(ctx, `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
			series_epoch, name, quantity, unit, semantic_kind, numeric_encoding, exact_recording_scale
		) VALUES (?, ?, 'device-A', 'point-integer', 'tag-integer', 'device-A', 'def-sql-int', 'bind-sql-int', 'epoch-sql-int', 'SQL Integer Counter', 'energy', 'kWh', 'counter', 'uint64', '1')`,
		sqlServerCircuitIntegerMeasurementID, f.record.ID)
	require.NoError(t, err)

	instantID := ensureSQLServerCircuitProfile(t, f, fixture.Manifest.Database, mode, "instantaneous", false, sqlServerCircuitMeasurementID)
	cumulativeID := ensureSQLServerCircuitProfile(t, f, fixture.Manifest.Database, mode, "cumulative", true, sqlServerCircuitMeasurementID)
	integerCumulativeID := ensureSQLServerCircuitProfile(t, f, fixture.Manifest.Database, mode, "cumulative", true, sqlServerCircuitIntegerMeasurementID)
	instantGroup, err := f.groups.Get(ctx, instantID)
	require.NoError(t, err)
	cumulativeGroup, err := f.groups.Get(ctx, cumulativeID)
	require.NoError(t, err)
	integerGroup, err := f.groups.Get(ctx, integerCumulativeID)
	require.NoError(t, err)
	require.Equal(t, instantGroup.Group.Destination.TableName, cumulativeGroup.Group.Destination.TableName)
	require.Equal(t, cumulativeGroup.Group.Destination.TableName, integerGroup.Group.Destination.TableName)

	base := time.Now().UTC().Truncate(sqlServerCircuitWindow).Add(2 * sqlServerCircuitWindow)
	clock := &atomic.Int64{}
	clock.Store(base.Add(-time.Second).UnixNano())
	f.groups.WithClock(func() time.Time { return time.Unix(0, clock.Load()).UTC() })

	instantToken := previewLong(t, f, instantID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, instantID, instantToken).Status)
	registerSQLServerCircuitCreatedTableCleanup(t, f, targetService, instantToken)
	cumulativeToken := previewLong(t, f, cumulativeID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, cumulativeID, cumulativeToken).Status)
	integerToken := previewLong(t, f, integerCumulativeID)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, integerCumulativeID, integerToken).Status)

	for _, id := range []string{instantID, cumulativeID, integerCumulativeID} {
		apply := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
		require.Equal(t, http.StatusOK, apply.Code, apply.Body.String())
	}
	snapshot, err := f.groups.ResolveAppliedAt(ctx, cumulativeID, base.Add(24*time.Hour))
	require.NoError(t, err)
	base = snapshot.EffectiveAt
	clock.Store(base.UnixNano())

	tags := tag.NewService(tag.NewSQLRepository(f.db))
	store := groupdelivery.NewStore(f.db)
	pipeConfig := grouppipeline.Config{
		NodeID: "sqlserver-circuit-plain", Owner: "sqlserver-circuit-plain/test",
		Now:          func() time.Time { return time.Unix(0, clock.Load()).UTC() },
		TickInterval: time.Hour, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour,
		MaxFutureSkew: time.Minute, OnError: func(err error) { t.Logf("PIPE_ON_ERROR: %v", err) },
		Sender: groupdelivery.SenderConfig{Backoff: func(int) (string, time.Time) {
			return groupdelivery.StateRetrying, time.Now().UTC().Add(3 * time.Second)
		}},
	}
	pipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: f.groups, Tags: tags, Destinations: targetService,
		Inspector: dbtarget.NewReadOnlyTableInspector(targetService), Store: store,
		Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db)),
	}, pipeConfig)
	require.NoError(t, pipe.Start(ctx))
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	return &sqlServerCircuitPlainFixture{
		managedGroupFixture: f, fixture: fixture, target: f.target, pipe: pipe,
		clock: clock, base: base, measurementID: sqlServerCircuitMeasurementID,
		database: fixture.Manifest.Database, instantID: instantID, cumulativeID: cumulativeID,
		integerCumulativeID: integerCumulativeID, targetService: targetService, tags: tags,
		store: store, pipeConfig: pipeConfig, createdTable: sqlServerCircuitPreviewCreatesGroup(instantToken),
	}
}

func TestSQLServerCircuitPlainMissingTableIsNotReady(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	f := newSQLServerCircuitPlainFixture(t, fixture)
	require.True(t, f.createdTable, "missing-table fault requires this test to have created the table")
	require.NoError(t, f.pipe.Stop(time.Second))
	group, err := f.groups.Get(t.Context(), f.cumulativeID)
	require.NoError(t, err)
	payload, err := f.store.RuntimeVersion(t.Context(), groupdelivery.GroupKey{
		WorkspaceID:   group.Group.WorkspaceID,
		GroupID:       group.Group.ID,
		GroupRevision: group.Group.Revision,
	})
	require.NoError(t, err)
	var snapshot workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(payload, &snapshot))
	require.NotNil(t, snapshot.RuntimeLayout)
	require.NotNil(t, snapshot.RuntimeLayout.SQLServer)
	capability, err := f.targetService.ProbeSQLServerRecordingAtRevision(t.Context(), group.Group.Destination.ConnectorID, group.Group.Destination.ConnectorRevision)
	require.NoError(t, err)
	current, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, group.Group.Destination.TableName)
	require.NoError(t, err)
	require.Equal(t, snapshot.RuntimeLayout.SQLServer.Group.Evidence.ObjectID, current.Evidence.ObjectID)
	require.Equal(t, snapshot.RuntimeLayout.SQLServer.Group.Evidence.ColumnDefaults, current.Evidence.ColumnDefaults,
		"the actual scope owner proof must still match the test-created table")
	require.Regexp(t, `^gw_long_[0-9a-f]{16}$`, group.Group.Destination.TableName)
	var existingRows int
	require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+sqlServerQualified(group.Group.Destination.TableSchema, group.Group.Destination.TableName)).Scan(&existingRows))
	require.Zero(t, existingRows, "only an empty fixture-created table may be removed for the missing-table case")
	_, err = f.target.ExecContext(t.Context(), "IF OBJECT_ID(@p1,N'U')=@p2 AND NOT EXISTS (SELECT 1 FROM "+
		sqlServerQualified(group.Group.Destination.TableSchema, group.Group.Destination.TableName)+") DROP TABLE "+
		sqlServerQualified(group.Group.Destination.TableSchema, group.Group.Destination.TableName)+
		" ELSE THROW 50000, 'fixture table identity or contents changed', 1", group.Group.Destination.TableSchema+"."+group.Group.Destination.TableName, current.Evidence.ObjectID)
	require.NoError(t, err)
	readiness, err := f.groups.Readiness(t.Context(), f.cumulativeID)
	require.NoError(t, err)
	require.False(t, readiness.SchemaReady, "a missing owned SQL Server table must not be reported ready")
	require.False(t, readiness.Ready)
	require.NotEmpty(t, readiness.Issues)
}

func ensureSQLServerCircuitProfile(t *testing.T, f managedGroupFixture, database string, mode workspace.WriteGroupStorageMode, kind string, cumulative bool, measurementID string) string {
	t.Helper()
	interval, unit := any(15), "second"
	if cumulative {
		interval, unit = 10, "minute"
	}
	circuitID := "sql-main"
	source := ""
	if measurementID == sqlServerCircuitIntegerMeasurementID {
		circuitID = "sql-main-integer"
		source = "integer"
	}
	request := profileRequest(t, f.writeGroupRouterFixture, kind, interval, unit, circuitID, source)
	group := request["group"].(map[string]any)
	group["layout_version"], group["storage_mode"] = "managed_long_v1", string(mode)
	destination := group["destination"].(map[string]any)
	destination["connector_id"], destination["connector_revision"] = "connector-A", "connector-1"
	destination["database"], destination["table_schema"] = database, f.namespace
	member := group["members"].([]map[string]any)[0]
	member["measurement_id"] = measurementID
	if !cumulative {
		member["value_role"] = "counter_raw"
	}
	request["expected_connector_revision"] = "connector-1"
	response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	_, saved := groupSaveData(t, decodeJSONBody(t, response))
	return saved["id"].(string)
}

func (f *sqlServerCircuitPlainFixture) restart(t *testing.T) {
	t.Helper()
	require.NoError(t, f.pipe.Stop(time.Second))
	pipe := grouppipeline.New(grouppipeline.Dependencies{
		Groups: f.groups, Tags: f.tags, Destinations: f.targetService,
		Inspector: dbtarget.NewReadOnlyTableInspector(f.targetService), Store: f.store,
		Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db)),
	}, f.pipeConfig)
	require.NoError(t, pipe.Start(t.Context()))
	t.Cleanup(func() { _ = pipe.Stop(time.Second) })
	f.pipe = pipe
}

func (f *sqlServerCircuitPlainFixture) GroupIDs() (instant, cumulative, integer string) {
	return f.instantID, f.cumulativeID, f.integerCumulativeID
}

func (f *sqlServerCircuitPlainFixture) GroupReadiness(ctx context.Context, groupID string) (*workspace.WriteGroupReadiness, error) {
	return f.groups.Readiness(ctx, groupID)
}

func (f *sqlServerCircuitPlainFixture) definition(t *testing.T, measurementID string) *measurement.MeasurementDefinition {
	t.Helper()
	definition, err := measurement.NewSQLRepository(f.db).GetByID(t.Context(), measurementID)
	require.NoError(t, err)
	return definition
}

func (f *sqlServerCircuitPlainFixture) accept(t *testing.T, groupID, measurementID string, input sqlServerCircuitInput) {
	t.Helper()
	group, err := f.groups.Get(t.Context(), groupID)
	require.NoError(t, err)
	member := group.Group.Members[0]
	definition := f.definition(t, measurementID)
	observed := f.base.Add(input.offset)
	f.clock.Store(observed.Add(time.Second).UnixNano())
	require.NoError(t, f.pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
		SampleID: input.sampleID, WorkspaceID: f.record.ID, DeviceID: member.DeviceID,
		PointID: member.PointID, TagID: member.TagID, SourceRevision: member.SourceRevision,
		MappingRevision: member.MappingRevision, MeasurementID: definition.ID,
		DefinitionRevision: definition.DefinitionRevision, SourceBindingRevision: definition.SourceBindingRevision,
		SeriesEpoch: definition.SeriesEpoch, AcquisitionID: input.sampleID + "-acquisition",
		ValueType: string(measurement.ExactUint64), ObservedAt: observed,
		ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: input.value,
		Quality: schema.QualityGood,
	}))
	var payload, digest string
	if scanErr := f.db.QueryRowContext(t.Context(), `SELECT payload,payload_digest FROM wg_delivery_samples WHERE group_id=? AND sample_id=?`,
		groupID, input.sampleID).Scan(&payload, &digest); scanErr != nil {
		t.Logf("ACCEPT_QUERY_FAILED: group=%s sample=%s err=%v", groupID, input.sampleID, scanErr)
		for _, s := range f.pipe.Status() {
			t.Logf("PIPE_STATUS_ON_FAIL: group=%s rev=%s state=%s reason=%s", s.GroupID, s.Revision, s.State, s.Reason)
		}
		require.NoError(t, scanErr)
	}
	var journal struct {
		Value measurement.ExactValue `json:"value"`
	}
	require.NoError(t, json.Unmarshal([]byte(payload), &journal))
	expected, err := measurement.ScaleExactInteger(measurement.NewUint64(input.value), definition.ExactRecordingScale)
	require.NoError(t, err)
	require.True(t, expected.Equal(journal.Value), "journal intake must preserve the confirmed exact scale")
	require.Len(t, digest, 64)
}

func (f *sqlServerCircuitPlainFixture) closeAt(t *testing.T, offset time.Duration) {
	t.Helper()
	f.clock.Store(f.base.Add(offset + time.Second).UnixNano())
	f.pipe.TickAll(t.Context())
}

func runSQLServerCircuitPlainIntegration(t *testing.T, f *sqlServerCircuitPlainFixture) {
	t.Helper()
	f.accept(t, f.cumulativeID, sqlServerCircuitMeasurementID, sqlServerCircuitInputs[0])
	f.closeAt(t, 10*time.Minute)
	f.restart(t)                 // offline/process restart before the gap continuation
	f.closeAt(t, 20*time.Minute) // gap: no sample in this window
	f.accept(t, f.cumulativeID, sqlServerCircuitMeasurementID, sqlServerCircuitInputs[1])
	f.closeAt(t, 30*time.Minute)
	f.accept(t, f.cumulativeID, sqlServerCircuitMeasurementID, sqlServerCircuitInputs[2])
	f.closeAt(t, 60*time.Minute)
	f.accept(t, f.integerCumulativeID, sqlServerCircuitIntegerMeasurementID, sqlServerCircuitIntegerInputs[0])
	f.closeAt(t, 70*time.Minute)
	f.accept(t, f.integerCumulativeID, sqlServerCircuitIntegerMeasurementID, sqlServerCircuitIntegerInputs[1])
	f.closeAt(t, 80*time.Minute)

	// Bound set to 60s based on empirical 237/20s throughput for 328 closed rows under ARM emulation; does not claim performance pass.
	require.Eventually(t, func() bool {
		var total, unfinished int
		err := f.db.QueryRowContext(t.Context(), `
			SELECT COUNT(*), SUM(CASE WHEN state <> 'sql_committed' THEN 1 ELSE 0 END)
			FROM wg_delivery_outbox WHERE group_id IN (?, ?, ?)`, f.instantID, f.cumulativeID, f.integerCumulativeID).Scan(&total, &unfinished)
		return err == nil && total >= 4 && unfinished == 0
	}, 60*time.Second, 25*time.Millisecond, "all closed SQL Server rows must reach a durable SQL receipt")

	group, err := f.groups.Get(t.Context(), f.cumulativeID)
	require.NoError(t, err)
	schemaName, table := group.Group.Destination.TableSchema, group.Group.Destination.TableName
	rows, err := f.target.QueryContext(t.Context(), "SELECT record_id, value_decimal, delta_value, status FROM "+sqlServerQualified(schemaName, table)+" WHERE workspace_id = @p1 AND group_id IN (@p2, @p3) ORDER BY bucket_start",
		group.Group.WorkspaceID, f.cumulativeID, f.integerCumulativeID)
	require.NoError(t, err)
	defer rows.Close()
	seen := map[string]bool{}
	statuses := map[string]int{}
	deltas := map[string]bool{}
	var recordIDs []string
	for rows.Next() {
		var recordID, status string
		var valueRaw, deltaRaw any
		require.NoError(t, rows.Scan(&recordID, &valueRaw, &deltaRaw, &status))
		recordIDs = append(recordIDs, recordID)
		statuses[status]++
		if valueRaw == nil {
			require.Equal(t, "missing", status)
			require.Nil(t, deltaRaw)
			continue
		}
		value, err := dbtarget.DecodeExactValue(dbtarget.SQLDialectSQLServer, "DECIMAL(38,18)", measurement.ExactDecimal, valueRaw)
		require.NoError(t, err)
		digits, ok := value.CanonicalDigits()
		require.True(t, ok)
		seen[digits] = true
		if status == "baseline" || status == "reset_or_uncertain" {
			require.Nil(t, deltaRaw)
		}
		if deltaRaw != nil {
			delta, err := dbtarget.DecodeExactValue(dbtarget.SQLDialectSQLServer, "DECIMAL(38,18)", measurement.ExactDecimal, deltaRaw)
			require.NoError(t, err)
			deltaDigits, ok := delta.CanonicalDigits()
			require.True(t, ok)
			deltas[deltaDigits] = true
		}
	}
	require.NoError(t, rows.Err())
	for _, expected := range []string{"12000", "12000.4", "0.7", "9007199254740992", "9007199254740993"} {
		require.True(t, seen[expected], "missing exact value %s; seen=%v", expected, seen)
	}
	for _, expected := range []string{"0.4", "1"} {
		require.True(t, deltas[expected], "missing exact delta %s; seen=%v", expected, deltas)
	}
	require.Positive(t, statuses["baseline"])
	require.Positive(t, statuses["missing"])
	require.Positive(t, statuses["gap"])
	require.Positive(t, statuses["reset_or_uncertain"])
	require.NotEmpty(t, recordIDs)

	effectDigests := map[string]string{}
	frozenRows := map[string]dbtarget.EncodedRow{}
	recordIDs = nil
	effectRows, err := f.db.QueryContext(t.Context(), `SELECT effect_key, payload_digest, payload FROM wg_delivery_outbox WHERE group_id IN (?, ?, ?)`, f.instantID, f.cumulativeID, f.integerCumulativeID)
	require.NoError(t, err)
	for effectRows.Next() {
		var effectKey, digest, payload string
		require.NoError(t, effectRows.Scan(&effectKey, &digest, &payload))
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		frozenRows[row.RecordID] = row
		recordIDs = append(recordIDs, row.RecordID)
		effectDigests[effectKey] = digest
	}
	require.NoError(t, effectRows.Err())
	require.NoError(t, effectRows.Close())
	require.NoError(t, f.pipe.Stop(time.Second))
	var receipts int
	require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+sqlServerQualified(schemaName, dbtarget.EffectReceiptTable)).Scan(&receipts))
	require.GreaterOrEqual(t, receipts, len(effectDigests))
	cleanupSQLServerCircuitEffects(t, f.target, schemaName, table, recordIDs, effectDigests, frozenRows)
	var outboxPayload string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_outbox WHERE group_id = ? LIMIT 1`, f.cumulativeID).Scan(&outboxPayload))
	require.True(t, json.Valid([]byte(outboxPayload)))
	require.NotContains(t, outboxPayload, "password")
}

func sqlServerParameters(count int) string {
	parts := make([]string, count)
	for i := range parts {
		parts[i] = fmt.Sprintf("@p%d", i+1)
	}
	return strings.Join(parts, ",")
}

func sqlServerArgs(values []string) []any {
	args := make([]any, len(values))
	for i, value := range values {
		args[i] = value
	}
	return args
}

func sqlServerQualified(schemaName, table string) string {
	quote := func(value string) string { return "[" + strings.ReplaceAll(value, "]", "]]") + "]" }
	return quote(schemaName) + "." + quote(table)
}
