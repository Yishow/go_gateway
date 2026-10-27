package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type counterApplyInspector struct {
	layout       *recordingplan.GroupSchemaLayout
	otherLayouts map[string]*recordingplan.GroupSchemaLayout
	mutate       func()
}

func (i *counterApplyInspector) InspectTableAtRevision(_ context.Context, _, _, schemaName, table string) (*dbtarget.TableInspection, error) {
	if i.mutate != nil {
		mutate := i.mutate
		i.mutate = nil
		mutate()
	}
	if table == dbtarget.EffectReceiptTable {
		return &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Schema: schemaName, Table: table,
			Columns: []dbtarget.ColumnInfo{
				{Name: "effect_key", DataType: "TEXT", PrimaryKey: true},
				{Name: "payload_digest", DataType: "TEXT"}, {Name: "committed_at", DataType: "TEXT"},
				{Name: "_gw_effect_receipts_v1", DataType: "TEXT"},
			}}, nil
	}
	layout := i.layout
	if other := i.otherLayouts[table]; other != nil {
		layout = other
	}
	inspection := managedTableInspection(layout, "")
	inspection.Schema = schemaName
	return inspection, nil
}

func (i *counterApplyInspector) InspectTableIndexesAtRevision(_ context.Context, _, _, _, table string) ([]dbtarget.IndexInfo, error) {
	var indexes []dbtarget.IndexInfo
	layout := i.layout
	if other := i.otherLayouts[table]; other != nil {
		layout = other
	}
	for _, index := range layout.Indexes {
		indexes = append(indexes, dbtarget.IndexInfo{Name: index.Name, Columns: index.Columns})
	}
	return indexes, nil
}

type counterApplyFixture struct {
	db        *sql.DB
	service   *WriteGroupService
	group     *WriteGroup
	mutation  WriteGroupMutation
	inspector *counterApplyInspector
}

func newCounterApplyFixture(t *testing.T, layout WriteGroupLayoutVersion, kind WriteGroupProfileKind, dataType schema.DataType) counterApplyFixture {
	t.Helper()
	return newCounterApplyFixtureWithScale(t, layout, kind, dataType, "")
}

func newCounterApplyFixtureWithScale(t *testing.T, layout WriteGroupLayoutVersion, kind WriteGroupProfileKind, dataType schema.DataType, scale string) counterApplyFixture {
	t.Helper()
	h := newCircuitProfileHarness(t, ":memory:")
	_, err := h.db.ExecContext(t.Context(), `UPDATE database_connectors SET kind = 'postgres', connection_config = '{"database":"persisted-db","schema":"persisted-schema"}' WHERE id = ?`, h.fixture.connectorID)
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `UPDATE points SET data_type = ? WHERE id = ?`, dataType, h.fixture.pointID)
	require.NoError(t, err)
	_, err = h.db.ExecContext(t.Context(), `UPDATE tags SET data_type = ? WHERE id = ?`, dataType, h.fixture.tagID)
	require.NoError(t, err)
	definition := &measurement.MeasurementDefinition{
		ID: "measurement-1", WorkspaceID: h.fixture.workspaceID, DeviceID: h.fixture.deviceID, PointID: h.fixture.pointID, TagID: &h.fixture.tagID,
		DefinitionRevision: "def-1", SourceBindingRevision: "binding-1", SeriesEpoch: "epoch-1", Unit: "kWh",
		SemanticKind: measurement.SemanticKindCounter, NumericEncoding: string(dataType), ExactRecordingScale: scale,
	}
	if dataType != schema.DataTypeFloat32 && dataType != schema.DataTypeFloat64 {
		definition.CounterPolicy = &measurement.CounterPolicy{RolloverModulus: "9007199254740993", MaxRatePerSecond: "0.125"}
	}
	require.NoError(t, measurement.NewSQLRepository(h.db).Create(t.Context(), definition))
	group := makeManagedLongGroup(h.fixture, h.fixture.tagID, definition.ID)
	group.LayoutVersion, group.RecordingProfile.Kind = layout, kind
	if layout == WriteGroupLayoutManagedWideProfileV1 {
		group.StorageMode = ""
		group.Destination.TableName = "counter_wide"
	}
	if kind == WriteGroupProfileInstantaneous {
		group.Members[0].ValueRole = WriteGroupValueRoleCounterRaw
		group.RecordingProfile.Interval, group.RecordingProfile.Unit = 15, WriteGroupIntervalUnitSecond
		group.RowPolicy.IntervalSeconds = 15
	}
	created, err := h.service.Create(t.Context(), h.mutation(t.Context(), group))
	require.NoError(t, err)
	saveManagedReadinessProof(t, h.db, created.Group.ID)
	loaded, err := h.service.Get(t.Context(), created.Group.ID)
	require.NoError(t, err)
	physical, err := managedSchemaLayout(loaded.Group, map[string]schema.DataType{h.fixture.tagID: dataType}, "postgres")
	require.NoError(t, err)
	inspector := &counterApplyInspector{layout: physical}
	h.service.WithManagedTableInspector(inspector).WithClock(func() time.Time { return time.Date(2026, 10, 1, 0, 0, 1, 0, time.UTC) })
	ready, err := h.service.Readiness(t.Context(), created.Group.ID)
	require.NoError(t, err)
	require.True(t, ready.Ready, "fixture must reach the real apply path: %+v", ready.Issues)
	return counterApplyFixture{db: h.db, service: h.service, group: loaded.Group, inspector: inspector,
		mutation: WriteGroupMutation{WorkspaceID: h.fixture.workspaceID, ExpectedWorkspaceRevision: created.WorkspaceRevision,
			ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: "connector-revision-1"}}
}

func (f counterApplyFixture) assertUnapplied(t *testing.T) {
	t.Helper()
	group, err := f.service.Get(t.Context(), f.group.ID)
	require.NoError(t, err)
	require.Empty(t, group.Group.AppliedRevision)
	require.Equal(t, WriteGroupStatusDraft, group.Group.Status)
	current, err := f.service.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, f.mutation.ExpectedWorkspaceRevision, current.DatabaseSetupRevision)
	for _, table := range []string{"wg_runtime_versions", "write_group_versions"} {
		var count int
		require.NoError(t, f.db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+table+" WHERE group_id = ?", f.group.ID).Scan(&count))
		require.Zero(t, count, table)
	}
}

func readCounterApplySnapshot(t *testing.T, query func(context.Context, string, ...any) *sql.Row, groupID string) (WriteGroupAppliedSnapshot, string) {
	t.Helper()
	var payload string
	require.NoError(t, query(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id = ?`, groupID).Scan(&payload))
	var snapshot WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(payload), &snapshot))
	return snapshot, payload
}

type counterApplyRaceRepository struct {
	*SQLRepository
	begins      int
	triggerAt   int
	beforeApply func()
}

func (r *counterApplyRaceRepository) BeginDatabaseSetup(ctx context.Context) (DatabaseSetupTx, error) {
	r.begins++
	// Apply first performs the two readiness snapshots; this is the deterministic
	// boundary after reviewed facts exist but before the apply transaction starts.
	at := r.triggerAt
	if at == 0 {
		at = 3
	}
	if r.begins == at {
		r.beforeApply()
	}
	return r.SQLRepository.BeginDatabaseSetup(ctx)
}
