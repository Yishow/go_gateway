package runtime

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const (
	wideFixtureTable = "gw_group_wide_profile"
	wideFixtureOwner = "_gw_owner_0123456789abcdef01234567"
)

// wideProfileFixture is an applied managed_wide_profile_v1 mean group over a real
// SQLite table: a numeric gauge, a text state and a confirmed counter_raw.
func wideProfileFixture(t *testing.T) *longFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "wide.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	column := func(name string, nullable bool) recordingplan.GroupSchemaColumn {
		return recordingplan.GroupSchemaColumn{Name: name, SQLType: "TEXT", Nullable: nullable}
	}
	key := column("record_id", false)
	key.PrimaryKey = true
	statements, err := recordingplan.GenerateGroupSchemaDDL("sqlite", recordingplan.GroupSchemaLayout{
		TableName: wideFixtureTable, OwnerColumn: wideFixtureOwner,
		Columns: []recordingplan.GroupSchemaColumn{
			key, column("group_id", false), column("device_id", false), column("bucket_start", false),
			column("provenance", false), column(wideFixtureOwner, false),
			column("v_gauge", true), column("v_state", true), column("v_raw", true),
		},
	})
	require.NoError(t, err)
	for _, statement := range statements {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}

	members := []workspace.WriteGroupMember{
		{
			DeviceID: "device-A", PointID: "point-G", TagID: "tag-G", SourceRevision: "src-G", MappingRevision: "map-G", TargetColumn: "v_gauge",
			ValueRole: workspace.WriteGroupValueRoleGauge, SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactDecimal,
		},
		{
			DeviceID: "device-A", PointID: "point-S", TagID: "tag-S", SourceRevision: "src-S", MappingRevision: "map-S", TargetColumn: "v_state",
			ValueRole: workspace.WriteGroupValueRoleState, SourceValueType: schema.DataTypeString, ResultValueType: measurement.ExactText,
		},
		{
			DeviceID: "device-A", PointID: "point-C", TagID: "tag-C", SourceRevision: "src-C", MappingRevision: "map-C", TargetColumn: "v_raw",
			MeasurementID: ptr("measurement-count"), ValueRole: workspace.WriteGroupValueRoleCounterRaw,
			SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactUint64,
		},
	}
	sink := &capturedGroupSink{}
	fixture := &longFixture{sink: sink, clock: boundaryAt(600), members: members, db: db}
	fixture.config = GroupBoundaryConfig{
		Group: &workspace.WriteGroup{
			ID: "group-W", WorkspaceID: "workspace-A", Revision: "rev-1", AppliedRevision: "rev-1", Members: members,
			LayoutVersion: workspace.WriteGroupLayoutManagedWideProfileV1,
			RecordingProfile: &workspace.WriteGroupRecordingProfile{
				Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond,
				SampleReduction: workspace.WriteGroupSampleReductionMean,
			},
			Destination: workspace.WriteGroupDestination{
				ConnectorID: "connector-1", ConnectorRevision: "crev-1", Database: "db", TableSchema: "main",
				TableName: wideFixtureTable, StorageStrategy: workspace.WriteGroupStorageStrategyManaged,
			},
			RowPolicy: workspace.WriteGroupRowPolicy{
				IntervalSeconds: 15, IncompletePolicy: "partial", RecordKeyColumn: "record_id", GroupIDColumn: "group_id",
				DeviceIDColumn: "device_id", BucketStartColumn: "bucket_start", ProvenanceColumn: "provenance",
			},
			WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: "receipt"},
		},
		TagTypes: map[string]schema.DataType{"tag-G": schema.DataTypeInt64, "tag-S": schema.DataTypeString, "tag-C": schema.DataTypeUint64},
		Counters: map[string]CounterSource{"tag-C": {
			MeasurementID: "measurement-count", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
			SeriesEpoch: "epoch-1", Unit: "kWh", Kind: aggregation.CounterKindMonotonic,
		}},
		Dialect: dbtarget.SQLDialectSQLite, Columns: sqliteColumns(t, db, wideFixtureTable),
		FirstBucket: boundaryStart, MaxFutureSkew: 2 * time.Second, ReceiptTableReady: true,
		Sink: sink, Clock: func() time.Time { return fixture.clock },
	}
	return fixture
}

func TestWideProfileBoundaryWritesOneRowPerWindow(t *testing.T) {
	f := wideProfileFixture(t)
	boundary := f.build(t)
	f.accept(t, boundary,
		f.envelope(0, "g1", at(time.Second), int64(1)), f.envelope(0, "g2", at(2*time.Second), int64(2)),
		f.envelope(1, "s1", at(time.Second), "idle"), f.envelope(1, "s2", at(5*time.Second), "run-42"),
		f.envelope(2, "c1", at(time.Second), uint64(10)), f.envelope(2, "c2", at(5*time.Second), uint64(30)),
	)
	require.NoError(t, boundary.Tick(t.Context(), at(14*time.Second)))
	require.Empty(t, f.sink.rows, "the window is still open")
	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	require.Len(t, f.sink.rows, 1, "one record per window, not one per member")
	cells := cellsOf(f.sink.rows[0])
	require.Equal(t, [3]any{"1.5", "run-42", "30"}, [3]any{cells["v_gauge"], cells["v_state"], cells["v_raw"]},
		"the gauge averages exactly; state and counter_raw keep the newest reading")
	require.Equal(t, boundaryStart.Format(time.RFC3339Nano), cells["bucket_start"])

	f.sink.rows = nil
	require.NoError(t, boundary.Tick(t.Context(), at(30*time.Second)))
	require.Len(t, f.sink.rows, 1, "a window with no sample is still one row")
	silent := cellsOf(f.sink.rows[0])
	require.Equal(t, [3]any{nil, nil, nil}, [3]any{silent["v_gauge"], silent["v_state"], silent["v_raw"]}, "NULL, never zero or an older value")
}

// TestWideProfileBoundaryFailsClosed pins every activation refusal: an unknown
// reduction, a cumulative profile of role-bearing members, a legacy member, an unconfirmed counter, a
// missing frozen fact and a target that cannot hold NULL or provenance never
// become a writer.
func TestWideProfileBoundaryFailsClosed(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate func(*longFixture)
		code   string
	}{
		"an unknown reduction": {func(f *longFixture) {
			f.config.Group.RecordingProfile.SampleReduction = "future"
		}, boundaryMemberMetadata},
		// A cumulative profile is served now (group_boundary_wide_counter_test.go), but
		// only for role-free raw members: gauge, state and counter_raw members stay
		// refused instead of being counted.
		"a cumulative profile over role-bearing members": {func(f *longFixture) {
			f.config.Group.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
		}, boundaryMemberMetadata},
		"an interval that is not the profile's": {func(f *longFixture) { f.config.Group.RowPolicy.IntervalSeconds = 30 }, boundaryProfileInterval},
		"a legacy member under the latest reduction": {func(f *longFixture) {
			f.config.Group.RecordingProfile.SampleReduction = ""
			f.members[0].ValueRole, f.members[0].SourceValueType, f.members[0].ResultValueType = "", "", ""
		}, boundaryProfileUnsupported},
		"a missing result type":        {func(f *longFixture) { f.members[0].ResultValueType = "" }, boundaryMemberMetadata},
		"a source type of another tag": {func(f *longFixture) { f.members[0].SourceValueType = schema.DataTypeInt16 }, boundaryMemberMetadata},
		"a missing tag type":           {func(f *longFixture) { delete(f.config.TagTypes, "tag-G") }, boundaryTagTypeUnsupport},
		"no counter provider output":   {func(f *longFixture) { f.config.Counters = nil }, boundaryCounterSemantics},
		"a counter_raw without a saved measurement": {func(f *longFixture) {
			f.members[2].MeasurementID = nil
		}, boundaryMemberMetadata},
		"a member column that cannot hold NULL": {func(f *longFixture) {
			for i := range f.config.Columns {
				if f.config.Columns[i].Name == "v_gauge" {
					f.config.Columns[i].Nullable = false
				}
			}
		}, boundaryLayoutBlocked},
		"a table without a provenance column": {func(f *longFixture) {
			f.config.Group.RowPolicy.ProvenanceColumn = ""
		}, boundaryLayoutBlocked},
		"a float column for an exact decimal": {func(f *longFixture) {
			for i := range f.config.Columns {
				if f.config.Columns[i].Name == "v_gauge" {
					f.config.Columns[i].DataType = "REAL"
				}
			}
		}, boundaryLayoutBlocked},
	} {
		t.Run(name, func(t *testing.T) {
			f := wideProfileFixture(t)
			tc.mutate(f)
			f.config.Group.Members = f.members
			require.Equal(t, tc.code, blockedCode(t, f.config))
		})
	}
}
