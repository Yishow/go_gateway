package grouppipeline

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

// averageShape is one applied recording shape whose closure carries a frozen
// average or counter result: mean on the long table, mean on a wide profile
// table, and a cumulative wide profile (the shape that carries counter baselines
// beside the wide row).
type averageShape string

const (
	shapeMeanLong       averageShape = "mean long"
	shapeMeanWide       averageShape = "mean wide"
	shapeCumulativeWide averageShape = "cumulative wide"
)

var averageShapes = []averageShape{shapeMeanLong, shapeMeanWide, shapeCumulativeWide}

func (s averageShape) wide() bool       { return s != shapeMeanLong }
func (s averageShape) cumulative() bool { return s == shapeCumulativeWide }

// meanWindow is the legal instantaneous interval the mean shapes close on; the
// cumulative shape keeps the 10 minute interval of the existing counter fixtures.
const meanWindow = 15 * time.Second

func (s averageShape) window() time.Duration {
	if s.cumulative() {
		return longWindow
	}
	return meanWindow
}

// windowEnd is the instant window n closes.
func (s averageShape) windowEnd(window int) time.Time {
	return base.Add(time.Duration(window+1) * s.window())
}

// recordsPerWindow is how many outbox effects one closed window leaves: the long
// table keeps one record per source, a wide profile one row for all of them.
func (s averageShape) recordsPerWindow() int {
	if s.wide() {
		return 1
	}
	return 2
}

// requireLegalProfile holds the fixture to the saved-profile contract the real
// save enforces (workspace.WriteGroupRecordingProfile.validate): an instantaneous
// profile counts whole seconds from 1 to 300, a cumulative one whole minutes from
// 1 to 600, a mean is instantaneous only, and the row interval is the profile's
// own. The runtime then freezes a group that Apply could really have saved.
func requireLegalProfile(t *testing.T, group *workspace.WriteGroup) {
	t.Helper()
	profile := group.RecordingProfile
	require.NotNil(t, profile)
	switch profile.Kind {
	case workspace.WriteGroupProfileInstantaneous:
		require.Equal(t, workspace.WriteGroupIntervalUnitSecond, profile.Unit)
		require.GreaterOrEqual(t, profile.Interval, 1)
		require.LessOrEqual(t, profile.Interval, 300, "an instantaneous profile is at most 300 seconds")
	case workspace.WriteGroupProfileCumulative:
		require.Equal(t, workspace.WriteGroupIntervalUnitMinute, profile.Unit)
		require.GreaterOrEqual(t, profile.Interval, 1)
		require.LessOrEqual(t, profile.Interval, 600, "a cumulative profile is at most 600 minutes")
		require.NotEqual(t, workspace.WriteGroupSampleReductionMean, profile.SampleReduction, "a cumulative profile never averages")
	default:
		require.Fail(t, "unsupported profile kind", string(profile.Kind))
	}
	require.Equal(t, profile.IntervalSeconds(), group.RowPolicy.IntervalSeconds, "the profile is the only authority for the row interval")
}

func averageGroup(shape averageShape) *workspace.WriteGroup {
	profile := &workspace.WriteGroupRecordingProfile{
		Kind: workspace.WriteGroupProfileInstantaneous, Interval: int(meanWindow / time.Second), Unit: workspace.WriteGroupIntervalUnitSecond,
		SampleReduction: workspace.WriteGroupSampleReductionMean,
	}
	if shape.cumulative() {
		profile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileCumulative, Interval: 10, Unit: workspace.WriteGroupIntervalUnitMinute}
	}
	member := func(i int, tag string) workspace.WriteGroupMember {
		m := workspace.WriteGroupMember{
			DeviceID: "device-1", PointID: fmt.Sprintf("point-%d", i+1), TagID: tag, TargetColumn: fmt.Sprintf("kwh_%d", i+1),
			SourceRevision: fmt.Sprintf("src-%d", i+1), MappingRevision: fmt.Sprintf("map-%d", i+1),
			SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactDecimal,
		}
		if !shape.wide() {
			m.CircuitID = "circuit-1"
		}
		if shape.cumulative() {
			m.ResultValueType, m.MeasurementID = measurement.ExactUint64, ptr(fmt.Sprintf("meas-%d", i+1))
		} else {
			m.ValueRole = workspace.WriteGroupValueRoleGauge
		}
		return m
	}
	group := &workspace.WriteGroup{
		ID: "group-A", WorkspaceID: "workspace-A", Revision: "rev-1", AppliedRevision: "rev-1", Status: workspace.WriteGroupStatusReady,
		LayoutVersion: workspace.WriteGroupLayoutManagedLongV1, StorageMode: workspace.WriteGroupStorageModePlain, RecordingProfile: profile,
		Members:     []workspace.WriteGroupMember{member(0, longKwhFirst), member(1, longKwhSecond)},
		Destination: workspace.WriteGroupDestination{ConnectorID: "connector-1", ConnectorRevision: "crev-1", TableSchema: "main", StorageStrategy: workspace.WriteGroupStorageStrategyManaged},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: profile.IntervalSeconds(), RecordKeyColumn: "record_id"},
		WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeReceipt},
	}
	if shape.wide() {
		group.LayoutVersion, group.StorageMode = workspace.WriteGroupLayoutManagedWideProfileV1, ""
		group.RowPolicy.GroupIDColumn, group.RowPolicy.BucketStartColumn, group.RowPolicy.ProvenanceColumn = "group_id", "bucket_start", "provenance"
		group.RowPolicy.IncompletePolicy = "partial"
		group.Destination.TableName = "average_wide"
	}
	return group
}

// wideAverageLayout is the managed wide profile table of the group: text key,
// owner and evidence columns and one nullable member column per source (SQLite
// holds the exact decimal and uint64 results as text).
func wideAverageLayout(group *workspace.WriteGroup) recordingplan.GroupSchemaLayout {
	layout := recordingplan.GroupSchemaLayout{TableName: group.Destination.TableName, OwnerColumn: canonicalManagedOwnerColumn(group)}
	add := func(name string, nullable, primary bool) {
		layout.Columns = append(layout.Columns, recordingplan.GroupSchemaColumn{Name: name, SQLType: "TEXT", Nullable: nullable, PrimaryKey: primary})
	}
	add(group.RowPolicy.RecordKeyColumn, false, true)
	add(group.RowPolicy.GroupIDColumn, false, false)
	add(group.RowPolicy.BucketStartColumn, false, false)
	add(group.RowPolicy.ProvenanceColumn, false, false)
	add(layout.OwnerColumn, false, false)
	for _, member := range group.Members {
		add(member.TargetColumn, true, false)
	}
	return layout
}

// newAveragePipeline builds the production pipeline over an applied group of the
// shape, with a real local database whose commits can fail or lose their reply
// and a real SQLite destination whose commits can lose theirs. The group is
// restored from its frozen proof, exactly as a restart does.
func newAveragePipeline(t *testing.T, shape averageShape) *longPipeline {
	t.Helper()
	f := &longPipeline{t: t, localPath: filepath.Join(t.TempDir(), "gateway.db"), clock: at(1)}
	f.localFault = &lostcommit.Faults{Lose: &atomic.Bool{}}
	f.targetFlt = &lostcommit.Faults{Lose: &atomic.Bool{}}
	f.group = averageGroup(shape)
	requireLegalProfile(t, f.group)

	var layout recordingplan.GroupSchemaLayout
	var err error
	if shape.wide() {
		f.table, layout = f.group.Destination.TableName, wideAverageLayout(f.group)
	} else {
		table, owner, ok := workspace.LongStorageIdentity(f.group)
		require.True(t, ok)
		f.table = table
		f.group.Destination.TableName = table
		layout, err = recordingplan.ManagedLongGroupLayout("sqlite", table, owner)
		require.NoError(t, err)
	}
	f.counters = confirmedCounters{sources: map[string]runtime.CounterSource{}}
	if shape.cumulative() {
		for _, member := range f.group.Members {
			f.counters.sources[member.TagID] = runtime.CounterSource{
				MeasurementID: *member.MeasurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
				SeriesEpoch: "epoch-1", Unit: "kWh", Kind: aggregation.CounterKindMonotonic,
			}
		}
	}

	statements, err := recordingplan.GenerateGroupSchemaDDL("sqlite", layout)
	require.NoError(t, err)
	targetPath := filepath.Join(t.TempDir(), "target.db")
	plain, err := sql.Open("sqlite", targetPath+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	for _, statement := range statements {
		_, err := plain.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}
	f.columns = sqliteColumnInfo(t, plain, f.table)
	require.NoError(t, plain.Close())
	target := lostcommit.OpenFaults(targetPath+"?_pragma=busy_timeout(5000)", f.targetFlt)
	t.Cleanup(func() { _ = target.Close() })
	f.dest = &liveDestinations{
		connector: &schema.DatabaseConnector{ID: "connector-1", Kind: schema.DatabaseConnectorKindSQLite, IdentityRevision: "crev-1", Enabled: true},
		target:    target, kind: schema.DatabaseConnectorKindSQLite,
	}

	f.snapshot = &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: f.group.WorkspaceID, GroupID: f.group.ID, GroupRevision: f.group.Revision, AppliedRevision: f.group.AppliedRevision,
		EffectiveAt: at(0), Group: f.group,
	}
	digest := sha256.Sum256([]byte(f.table))
	frozen := *f.snapshot
	frozen.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLite, Columns: f.columns,
		TagTypes:          map[string]schema.DataType{longKwhFirst: schema.DataTypeUint64, longKwhSecond: schema.DataTypeUint64},
		ReceiptTableReady: true, SchemaDigest: hex.EncodeToString(digest[:]),
	}
	// The applied snapshot the group source resolves carries the same frozen proof the restored copy does.
	f.snapshot.RuntimeLayout = frozen.RuntimeLayout
	f.groups = &fakeGroups{
		groups:    []*workspace.WriteGroup{liveGroup(f.group, workspace.WriteGroupStatusReady)},
		snapshots: map[string][]*workspace.WriteGroupAppliedSnapshot{f.group.ID: {f.snapshot}},
	}
	if shape.cumulative() {
		freezeFixtureCounters(t, &frozen, f.counters.sources)
	}
	f.openLocal(&frozen)
	return f
}

// feed accepts the scripted readings of one window through the production
// intake. The mean shapes get two or one readings per source inside their 15
// second window, the cumulative shape one counter reading per source; every ACK
// is a durable journal write.
func (f *longPipeline) feed(shape averageShape, window int) {
	f.t.Helper()
	if shape.cumulative() {
		if window == 0 {
			f.closeFirstWindowReadings()
		} else {
			f.acceptSecondWindow()
		}
		return
	}
	start := window * int(meanWindow/time.Second)
	script := [][]struct {
		id    string
		after int
		value uint64
	}{
		{{"m0-0a", 1, 9007199254740992}, {"m0-0b", 2, 9007199254740993}, {"m1-0a", 1, 10}, {"m1-0b", 3, 21}},
		{{"m0-1a", 1, 100}, {"m0-1b", 5, 200}, {"m1-1a", 2, 7}},
	}
	for _, reading := range script[window] {
		member := 0
		if reading.id[1] == '1' {
			member = 1
		}
		require.NoError(f.t, f.reading(member, reading.id, at(start+reading.after), reading.value), reading.id)
	}
}

// closeFirstWindowReadings accepts window 0's counter readings without closing it.
func (f *longPipeline) closeFirstWindowReadings() {
	f.t.Helper()
	require.NoError(f.t, f.reading(0, "m1-0", at(9*60), firstStart))
	require.NoError(f.t, f.reading(1, "m2-0", at(9*60), secondStart))
}

// averageRow is one stored member result of a window, however the layout keeps it.
type averageRow struct {
	value string
	// sum and count are the mean evidence; status and delta the counter evidence.
	sum, status, delta string
	count              int64
}

type averageEvidence struct {
	Mean *struct {
		Sum        string `json:"sum"`
		ValidCount int64  `json:"valid_count"`
	} `json:"mean"`
	Counter *struct {
		Status string `json:"status"`
		Delta  string `json:"delta"`
	} `json:"counter"`
}

func windowOfBucket(t *testing.T, shape averageShape, text string) int {
	t.Helper()
	instant, err := time.Parse(time.RFC3339Nano, text)
	require.NoError(t, err, text)
	return int(instant.Sub(base) / shape.window())
}

// averageRows reads the destination back into "m<member>@w<window>" rows.
func (f *longPipeline) averageRows(shape averageShape) map[string]averageRow {
	f.t.Helper()
	rows := map[string]averageRow{}
	add := func(member, window int, value any, evidence averageEvidence) {
		row := averageRow{}
		if text, ok := value.(string); ok {
			row.value = text
		}
		if evidence.Mean != nil {
			row.sum, row.count = evidence.Mean.Sum, evidence.Mean.ValidCount
		}
		if evidence.Counter != nil {
			row.status, row.delta = evidence.Counter.Status, evidence.Counter.Delta
		}
		rows[fmt.Sprintf("m%d@w%d", member, window)] = row
	}
	if !shape.wide() {
		query := fmt.Sprintf(`SELECT tag_id, bucket_start, CAST(value_decimal AS TEXT), provenance FROM %s`, f.table)
		result, err := f.plainTarget().QueryContext(f.t.Context(), query)
		require.NoError(f.t, err)
		defer result.Close()
		for result.Next() {
			var tag, bucket, value, provenance string
			require.NoError(f.t, result.Scan(&tag, &bucket, &value, &provenance))
			var evidence averageEvidence
			require.NoError(f.t, json.Unmarshal([]byte(provenance), &evidence), provenance)
			member := 0
			if tag == longKwhSecond {
				member = 1
			}
			add(member, windowOfBucket(f.t, shape, bucket), value, evidence)
		}
		require.NoError(f.t, result.Err())
		return rows
	}
	query := fmt.Sprintf(`SELECT bucket_start, kwh_1, kwh_2, provenance FROM %s`, f.table)
	result, err := f.plainTarget().QueryContext(f.t.Context(), query)
	require.NoError(f.t, err)
	defer result.Close()
	for result.Next() {
		var bucket, provenance string
		var first, second any
		require.NoError(f.t, result.Scan(&bucket, &first, &second, &provenance))
		var evidence []averageEvidence
		require.NoError(f.t, json.Unmarshal([]byte(provenance), &evidence), provenance)
		require.Len(f.t, evidence, 2, "a wide row carries evidence for each member")
		window := windowOfBucket(f.t, shape, bucket)
		add(0, window, first, evidence[0])
		add(1, window, second, evidence[1])
	}
	require.NoError(f.t, result.Err())
	return rows
}

// want is the exact result of every scripted member window.
func (s averageShape) want(windows ...int) map[string]averageRow {
	all := map[string]averageRow{
		"m0@w0": {value: "9007199254740992.5", sum: "18014398509481985", count: 2},
		"m1@w0": {value: "15.5", sum: "31", count: 2},
		"m0@w1": {value: "150", sum: "300", count: 2},
		"m1@w1": {value: "7", sum: "7", count: 1},
	}
	if s.cumulative() {
		all = map[string]averageRow{
			"m0@w0": {value: "9007199254740992", status: "baseline"},
			"m1@w0": {value: "100", status: "baseline"},
			"m0@w1": {value: "9007199254740993", status: "ok", delta: "1"},
			"m1@w1": {value: "130", status: "ok", delta: "30"},
		}
	}
	want := map[string]averageRow{}
	for _, window := range windows {
		for member := range 2 {
			key := fmt.Sprintf("m%d@w%d", member, window)
			want[key] = all[key]
		}
	}
	return want
}
