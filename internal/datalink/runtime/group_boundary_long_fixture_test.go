package runtime

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

const (
	longFixtureTable = "gw_long_0123456789abcdef"
	longFixtureOwner = "_gw_owner_0123456789abcdef01234567"
	longTimeLayout   = "2006-01-02T15:04:05.000000Z"
)

// longFixture is an applied managed_long_v1 instantaneous group over a real
// SQLite long table, so column metadata comes from the real catalog.
type longFixture struct {
	config  GroupBoundaryConfig
	sink    *capturedGroupSink
	clock   time.Time
	members []workspace.WriteGroupMember
	db      *sql.DB
}

func newLongFixture(t *testing.T) *longFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "long.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	layout, err := recordingplan.ManagedLongGroupLayout("sqlite", longFixtureTable, longFixtureOwner)
	require.NoError(t, err)
	statements, err := recordingplan.GenerateGroupSchemaDDL("sqlite", layout)
	require.NoError(t, err)
	for _, statement := range statements {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}

	members := []workspace.WriteGroupMember{
		{DeviceID: "device-A", PointID: "point-A", TagID: "tag-A", SourceRevision: "src-A", MappingRevision: "map-A", TargetColumn: "v_a", Required: true, CircuitID: "circuit-a"},
		{DeviceID: "device-B", PointID: "point-B", TagID: "tag-B", SourceRevision: "src-B", MappingRevision: "map-B", TargetColumn: "v_b", Required: true, CircuitID: "circuit-b"},
		{DeviceID: "device-A", PointID: "point-A2", TagID: "tag-run", SourceRevision: "src-run", MappingRevision: "map-run", TargetColumn: "v_run", Required: true},
		{DeviceID: "device-B", PointID: "point-B2", TagID: "tag-batch", SourceRevision: "src-batch", MappingRevision: "map-batch", TargetColumn: "v_batch", Required: false},
		{DeviceID: "device-B", PointID: "point-B3", TagID: "tag-count", SourceRevision: "src-count", MappingRevision: "map-count", TargetColumn: "v_count", Required: true, MeasurementID: ptr("measurement-count")},
	}
	sink := &capturedGroupSink{}
	fixture := &longFixture{sink: sink, clock: boundaryAt(600), members: members, db: db}
	fixture.config = GroupBoundaryConfig{
		Group: &workspace.WriteGroup{
			ID: "group-L", WorkspaceID: "workspace-A", Revision: "rev-1", AppliedRevision: "rev-1",
			Members:       members,
			LayoutVersion: workspace.WriteGroupLayoutManagedLongV1, StorageMode: workspace.WriteGroupStorageModePlain,
			RecordingProfile: &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond},
			Destination: workspace.WriteGroupDestination{
				ConnectorID: "connector-1", ConnectorRevision: "crev-1", Database: "db", TableSchema: "main",
				TableName: longFixtureTable, StorageStrategy: workspace.WriteGroupStorageStrategyManaged,
			},
			RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 15},
			WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: "receipt"},
		},
		TagTypes: map[string]schema.DataType{
			"tag-A": schema.DataTypeFloat64, "tag-B": schema.DataTypeInt64, "tag-run": schema.DataTypeBool,
			"tag-batch": schema.DataTypeString, "tag-count": schema.DataTypeUint64,
		},
		Dialect: dbtarget.SQLDialectSQLite, Columns: sqliteColumns(t, db, longFixtureTable),
		FirstBucket: boundaryStart, MaxFutureSkew: 2 * time.Second, ReceiptTableReady: true,
		Sink: sink, Clock: func() time.Time { return fixture.clock },
	}
	return fixture
}

func ptr[T any](value T) *T { return &value }

// sqliteColumns reads the real catalog the way production inspection does.
func sqliteColumns(t *testing.T, db *sql.DB, table string) []dbtarget.ColumnInfo {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), fmt.Sprintf(`PRAGMA table_info('%s')`, table))
	require.NoError(t, err)
	defer rows.Close()
	var columns []dbtarget.ColumnInfo
	for rows.Next() {
		var cid, notNull, pk int
		var name, dataType string
		var defaultValue sql.NullString
		require.NoError(t, rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk))
		columns = append(columns, dbtarget.ColumnInfo{Name: name, DataType: dataType, Nullable: notNull == 0, PrimaryKey: pk > 0, Unique: pk > 0})
	}
	require.NoError(t, rows.Err())
	require.NotEmpty(t, columns)
	return columns
}

func (f *longFixture) build(t *testing.T) *GroupBoundary {
	t.Helper()
	freezeDurableFixture(t, f)
	boundary, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	return boundary
}

// at is an offset from the fixture's first window start.
func at(offset time.Duration) time.Time { return boundaryStart.Add(offset) }

func (f *longFixture) envelope(memberIndex int, id string, observed time.Time, value any) measurement.SampleEnvelope {
	member := f.members[memberIndex]
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: "workspace-A", DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: observed, ReceivedAt: observed.Add(40 * time.Millisecond), TimeOrigin: "gateway",
		Value: value, Quality: schema.QualityGood,
	}
}

func (f *longFixture) accept(t *testing.T, boundary *GroupBoundary, envelopes ...measurement.SampleEnvelope) {
	t.Helper()
	for _, envelope := range envelopes {
		require.NoError(t, boundary.AcceptSample(t.Context(), envelope), envelope.SampleID)
	}
}

// longRows returns the closed rows keyed by tag and window start.
type longRowKey struct{ tag, bucket string }

func longRows(t *testing.T, sink *capturedGroupSink) map[longRowKey]map[string]any {
	t.Helper()
	rows := map[longRowKey]map[string]any{}
	for _, row := range sink.rows {
		cells := cellsOf(row)
		key := longRowKey{tag: fmt.Sprint(cells["tag_id"]), bucket: fmt.Sprint(cells["bucket_start"])}
		require.NotContains(t, rows, key, "a member has exactly one record per window")
		rows[key] = cells
	}
	return rows
}

func bucketText(offset time.Duration) string { return at(offset).Format(longTimeLayout) }

// valueSlots lists which typed value columns of a row hold a value.
func valueSlots(cells map[string]any) []string {
	var slots []string
	for _, column := range []string{"value_int", "value_decimal", "value_float", "value_bool", "value_text"} {
		if cells[column] != nil {
			slots = append(slots, column)
		}
	}
	slices.Sort(slots)
	return slots
}

func provenanceText(t *testing.T, cells map[string]any) string {
	t.Helper()
	text, ok := cells["provenance"].(string)
	require.True(t, ok, "provenance is stored as JSON text on SQLite")
	require.True(t, strings.HasPrefix(text, "{"), text)
	return text
}
