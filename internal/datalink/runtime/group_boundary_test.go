package runtime

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

var boundaryStart = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func boundaryAt(seconds int) time.Time {
	return boundaryStart.Add(time.Duration(seconds) * time.Second)
}

type capturedGroupSink struct {
	mu       sync.Mutex
	rows     []GroupRow
	outcomes []snapshot.Outcome
	rowErr   error
}

func (s *capturedGroupSink) AcceptRow(_ context.Context, row GroupRow) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = append(s.rows, row)
	return s.rowErr
}

func (s *capturedGroupSink) ReportOutcome(_ context.Context, outcome snapshot.Outcome) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.outcomes = append(s.outcomes, outcome)
	return nil
}

type boundaryFixture struct {
	config  GroupBoundaryConfig
	sink    *capturedGroupSink
	clock   time.Time
	members []workspace.WriteGroupMember
}

// Two devices both expose 40001 and their tags share a display name, so only
// persisted IDs can tell the members apart.
func newBoundaryFixture(t *testing.T) *boundaryFixture {
	t.Helper()
	members := []workspace.WriteGroupMember{
		{DeviceID: "device-A", PointID: "point-A", TagID: "tag-A", SourceRevision: "src-A", MappingRevision: "map-A", TargetColumn: "temperature", Required: true},
		{DeviceID: "device-B", PointID: "point-B", TagID: "tag-B", SourceRevision: "src-B", MappingRevision: "map-B", TargetColumn: "pressure", Required: true},
		{DeviceID: "device-A", PointID: "point-A2", TagID: "tag-run", SourceRevision: "src-run", MappingRevision: "map-run", TargetColumn: "running", Required: true},
		{DeviceID: "device-B", PointID: "point-B2", TagID: "tag-batch", SourceRevision: "src-batch", MappingRevision: "map-batch", TargetColumn: "batch", Required: true},
		{DeviceID: "device-B", PointID: "point-B3", TagID: "tag-count", SourceRevision: "src-count", MappingRevision: "map-count", TargetColumn: "counter", Required: true},
	}
	sink := &capturedGroupSink{}
	fixture := &boundaryFixture{sink: sink, clock: boundaryAt(5), members: members}
	fixture.config = GroupBoundaryConfig{
		Group: &workspace.WriteGroup{
			ID: "group-G", WorkspaceID: "workspace-A", Revision: "rev-1", AppliedRevision: "rev-1",
			Members: members,
			Destination: workspace.WriteGroupDestination{
				ConnectorID: "connector-1", ConnectorRevision: "crev-1", Database: "db", TableName: "readings",
			},
			RowPolicy: workspace.WriteGroupRowPolicy{IntervalSeconds: 10},
		},
		TagTypes: map[string]schema.DataType{
			"tag-A": schema.DataTypeFloat64, "tag-B": schema.DataTypeInt64, "tag-run": schema.DataTypeBool,
			"tag-batch": schema.DataTypeString, "tag-count": schema.DataTypeUint64,
		},
		Dialect: dbtarget.SQLDialectSQLite,
		Columns: []dbtarget.ColumnInfo{
			{Name: "temperature", DataType: "REAL"}, {Name: "pressure", DataType: "INTEGER"},
			{Name: "running", DataType: "INTEGER"}, {Name: "batch", DataType: "TEXT"},
			{Name: "counter", DataType: "INTEGER"}, {Name: "provenance", DataType: "TEXT", Nullable: true},
		},
		FirstBucket:   boundaryStart,
		MaxFutureSkew: 2 * time.Second,
		Sink:          sink,
		Clock:         func() time.Time { return fixture.clock },
	}
	return fixture
}

func (f *boundaryFixture) envelope(memberIndex int, id string, seconds int, value any) measurement.SampleEnvelope {
	member := f.members[memberIndex]
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: "workspace-A", DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: boundaryAt(seconds), ReceivedAt: boundaryAt(seconds), TimeOrigin: "gateway",
		Value: value, Quality: schema.QualityGood,
	}
}

func (f *boundaryFixture) build(t *testing.T) *GroupBoundary {
	t.Helper()
	boundary, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	return boundary
}

func (f *boundaryFixture) feedAll(t *testing.T, boundary *GroupBoundary) {
	t.Helper()
	values := []any{21.5, int64(101), true, "batch-001", uint64(9007199254740993)}
	for i, value := range values {
		require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(i, "s"+f.members[i].TagID, 4, value)))
	}
}

func cellsOf(row GroupRow) map[string]any {
	cells := map[string]any{}
	for _, cell := range row.Encoded.Cells {
		cells[cell.Column] = cell.Value
	}
	return cells
}

func TestBasicGroupBoundaryClosesOneExactMixedRow(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	f.feedAll(t, boundary)

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(9)))
	require.Empty(t, f.sink.rows, "the bucket is still open")
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Len(t, f.sink.rows, 1)
	require.Empty(t, f.sink.outcomes)

	row := f.sink.rows[0]
	require.Equal(t, snapshot.OutcomeRow, row.Outcome.Kind)
	cells := cellsOf(row)
	require.Equal(t, 21.5, cells["temperature"])
	require.Equal(t, int64(101), cells["pressure"])
	require.Equal(t, int64(1), cells["running"])
	require.Equal(t, "batch-001", cells["batch"])
	require.Equal(t, int64(9007199254740993), cells["counter"], "exact uint64 digits survive the group path")
	wantRecord, err := snapshot.RecordID("workspace-A", "group-G", "rev-1", "", boundaryStart)
	require.NoError(t, err)
	require.Equal(t, wantRecord, row.Encoded.RecordID)
	require.NotEmpty(t, row.Encoded.EffectKey)
}

func TestBasicGroupBoundaryIgnoresSamplesOfOtherGroupsButRefusesForeignWorkspace(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)

	other := f.envelope(0, "other", 4, 1.0)
	other.TagID = "tag-unrelated"
	require.NoError(t, boundary.AcceptSample(t.Context(), other), "a sample for another group is not this boundary's business")

	sameAddressOtherDevice := f.envelope(0, "twin", 4, 1.0)
	sameAddressOtherDevice.DeviceID = "device-B" // tag-A is not on device-B in this group
	require.NoError(t, boundary.AcceptSample(t.Context(), sameAddressOtherDevice), "identity is the persisted IDs, never the address")

	foreign := f.envelope(0, "foreign", 4, 1.0)
	foreign.WorkspaceID = "workspace-B"
	var refused *GroupSampleError
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), foreign), &refused)
	require.Equal(t, "workspace-mismatch", refused.Reason)
	require.Zero(t, f.sink.rowCount())
}

func (s *capturedGroupSink) rowCount() int { return len(s.rows) }

func TestBasicGroupBoundaryMissingMemberReportsSkippedAndWritesNoRow(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	for i, value := range []any{21.5, int64(101), true, "batch-001"} { // counter never arrives
		require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(i, "m"+f.members[i].TagID, 4, value)))
	}

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Empty(t, f.sink.rows, "default skip_row writes no SQL row")
	require.Len(t, f.sink.outcomes, 1)
	require.Equal(t, snapshot.OutcomeSkipped, f.sink.outcomes[0].Kind)
	require.Equal(t, snapshot.ReasonIncompleteRequired, f.sink.outcomes[0].Reason)
}

func TestBasicGroupBoundaryTypeMismatchAndNonFiniteBecomeBadNotZero(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	f.feedAll(t, boundary)
	// A later, wrong-typed reading for pressure and a NaN temperature.
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(1, "pressure-text", 6, "101")))
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "temperature-nan", 7, math.NaN())))

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Empty(t, f.sink.rows)
	require.Len(t, f.sink.outcomes, 1)
	statuses := map[string]snapshot.MemberResult{}
	for _, member := range f.sink.outcomes[0].Members {
		statuses[member.MemberKey] = member
	}
	require.Len(t, statuses, 5)
	badCount := 0
	for _, member := range statuses {
		if member.Status == snapshot.MemberBad {
			badCount++
			require.Equal(t, "type-mismatch", member.Reason)
			require.Equal(t, schema.QualityBad, member.Sample.Quality)
			require.Equal(t, "", string(member.Sample.Value.Type()), "an unconvertible value is not carried as a good value")
		}
	}
	require.Equal(t, 2, badCount)
}

func TestBasicGroupBoundaryFailedReadKeepsCollectorReason(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	f.feedAll(t, boundary)
	f.clock = boundaryAt(8)
	failed := f.envelope(2, "run-failed", 8, nil)
	failed.Quality, failed.QualityReason = schema.QualityBad, "connection-failed"
	require.NoError(t, boundary.AcceptSample(t.Context(), failed))

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Empty(t, f.sink.rows)
	var running snapshot.MemberResult
	for _, member := range f.sink.outcomes[0].Members {
		if member.Status == snapshot.MemberBad {
			running = member
		}
	}
	require.Equal(t, "connection-failed", running.Reason)
}

func TestBasicGroupBoundaryRefusesStaleRevisionAndLateSamples(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)

	stale := f.envelope(0, "stale-scale", 4, 1.0)
	stale.MappingRevision = "map-A-old"
	var refused *GroupSampleError
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), stale), &refused)
	require.Equal(t, snapshot.ReasonRevisionMismatch, refused.Reason)

	f.feedAll(t, boundary)
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Len(t, f.sink.rows, 1)

	f.clock = boundaryAt(12)
	late := f.envelope(0, "late", 9, 99.0)
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), late), &refused)
	require.Equal(t, snapshot.ReasonLate, refused.Reason)
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(12)))
	require.Len(t, f.sink.rows, 1, "a late sample never rewrites or re-emits the closed row")
}

func TestBasicGroupBoundaryDuplicateSampleIsNoOpAndConflictIsRefused(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	first := f.envelope(0, "dup", 4, 21.5)
	require.NoError(t, boundary.AcceptSample(t.Context(), first))
	require.NoError(t, boundary.AcceptSample(t.Context(), first))
	changed := f.envelope(0, "dup", 4, 22.5)
	var refused *GroupSampleError
	require.ErrorAs(t, boundary.AcceptSample(t.Context(), changed), &refused)
	require.Equal(t, snapshot.ReasonIdentityConflict, refused.Reason)
}

func TestBasicGroupBoundaryTimeDrivenSilentBucketNeedsNoSample(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	ticks := make(chan time.Time)
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		boundary.RunTicks(ctx, ticks)
		close(done)
	}()
	ticks <- boundaryAt(25)
	ticks <- boundaryAt(25) // second delivery is processed only after the first completed
	cancel()
	<-done

	f.sink.mu.Lock()
	defer f.sink.mu.Unlock()
	require.Empty(t, f.sink.rows)
	require.Len(t, f.sink.outcomes, 2, "buckets [0,10) and [10,20) each close once")
	for _, outcome := range f.sink.outcomes {
		require.Equal(t, snapshot.OutcomeNoData, outcome.Kind)
	}
}

func TestBasicGroupBoundaryEncodeFailureIsReportedAsBlockedNotDropped(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	f.feedAll(t, boundary)
	// Replace the uint64 reading with MaxUint64, which the INTEGER column cannot hold.
	huge := f.envelope(4, "huge", 6, uint64(math.MaxUint64))
	require.NoError(t, boundary.AcceptSample(t.Context(), huge))

	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Empty(t, f.sink.rows)
	require.Len(t, f.sink.outcomes, 1)
	require.Equal(t, snapshot.OutcomeSkipped, f.sink.outcomes[0].Kind)
	require.Equal(t, "encode-blocked:sql-value-blocked", f.sink.outcomes[0].Reason)
}

func TestBasicGroupBoundarySinkFailureIsReturned(t *testing.T) {
	f := newBoundaryFixture(t)
	f.sink.rowErr = errors.New("sink down")
	boundary := f.build(t)
	f.feedAll(t, boundary)
	require.ErrorContains(t, boundary.Tick(t.Context(), boundaryAt(10)), "sink down")
}

func TestBasicGroupBoundaryActivationBlocksUnsupportedColumnsAndIdentity(t *testing.T) {
	cases := map[string]struct {
		mutate func(*GroupBoundaryConfig)
		code   string
		issue  string
		column string
	}{
		"uint64 on REAL": {func(c *GroupBoundaryConfig) { c.Columns[4].DataType = "REAL" }, "layout-blocked", "unsupported-sql-type", "counter"},
		"missing column": {func(c *GroupBoundaryConfig) { c.Columns = c.Columns[:4] }, "layout-blocked", "column-missing", "counter"},
		"partial without provenance": {func(c *GroupBoundaryConfig) {
			c.Group.RowPolicy.IncompletePolicy = "partial"
			c.Group.Members[4].Required = false
			c.Columns[4].Nullable = true
		}, "layout-blocked", "partial-requires-provenance", ""},
		"entity without identity column": {func(c *GroupBoundaryConfig) { c.Group.Members[0].EntityKey = "plant-1" }, "layout-blocked", "identity-column-missing", ""},
	}
	for name, c := range cases {
		f := newBoundaryFixture(t)
		f.config.Group = copyGroup(f.config.Group)
		f.config.Columns = append([]dbtarget.ColumnInfo(nil), f.config.Columns...)
		c.mutate(&f.config)
		_, err := NewGroupBoundary(f.config)
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked, name)
		require.Equal(t, c.code, blocked.Code, name)
		found := false
		for _, issue := range blocked.Issues {
			if issue.Code == c.issue && issue.Column == c.column {
				found = true
			}
		}
		require.True(t, found, "%s: %+v", name, blocked.Issues)
	}
}

func copyGroup(group *workspace.WriteGroup) *workspace.WriteGroup {
	cloned := *group
	cloned.Members = append([]workspace.WriteGroupMember(nil), group.Members...)
	return &cloned
}

func TestBasicGroupBoundaryActivationBlocksInvalidGroupDefinitions(t *testing.T) {
	cases := map[string]struct {
		mutate func(*GroupBoundaryConfig)
		code   string
	}{
		"not the applied version":     {func(c *GroupBoundaryConfig) { c.Group.Revision = "rev-2" }, "group-not-applied"},
		"never applied":               {func(c *GroupBoundaryConfig) { c.Group.AppliedRevision = "" }, "group-not-applied"},
		"no interval":                 {func(c *GroupBoundaryConfig) { c.Group.RowPolicy.IntervalSeconds = 0 }, "snapshot-config-invalid"},
		"unaligned first bucket":      {func(c *GroupBoundaryConfig) { c.FirstBucket = boundaryAt(3) }, "snapshot-config-invalid"},
		"unsupported tag type":        {func(c *GroupBoundaryConfig) { c.TagTypes["tag-A"] = schema.DataType("decimal") }, "tag-type-unsupported"},
		"unknown tag type":            {func(c *GroupBoundaryConfig) { delete(c.TagTypes, "tag-B") }, "tag-type-unsupported"},
		"no sink":                     {func(c *GroupBoundaryConfig) { c.Sink = nil }, "sink-missing"},
		"no group":                    {func(c *GroupBoundaryConfig) { c.Group = nil }, "group-missing"},
		"member without identity":     {func(c *GroupBoundaryConfig) { c.Group.Members[1].PointID = "" }, "member-identity-incomplete"},
		"optional member in skip_row": {func(c *GroupBoundaryConfig) { c.Group.Members[1].Required = false }, "snapshot-config-invalid"},
	}
	for name, c := range cases {
		f := newBoundaryFixture(t)
		f.config.Group = copyGroup(f.config.Group)
		tagTypes := map[string]schema.DataType{}
		for k, v := range f.config.TagTypes {
			tagTypes[k] = v
		}
		f.config.TagTypes = tagTypes
		c.mutate(&f.config)
		_, err := NewGroupBoundary(f.config)
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked, name)
		require.Equal(t, c.code, blocked.Code, name)
	}
}

func TestProductionGroupOutageRecoveryBoundaryOnlyOwnsItsEffectiveInterval(t *testing.T) {
	f := newBoundaryFixture(t)
	f.config.FirstBucket = boundaryAt(10)
	f.config.Until = boundaryAt(30)
	f.clock = boundaryAt(40)
	boundary := f.build(t)

	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "before", 4, 1.0)), "an earlier revision's sample is ignored, not refused")
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "after", 31, 1.0)), "a later revision's sample is ignored, not refused")
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(40)))
	require.NotZero(t, len(f.sink.outcomes))
	for _, outcome := range f.sink.outcomes {
		require.NotEqual(t, snapshot.OutcomeRow, outcome.Kind)
		require.False(t, outcome.BucketStart.Before(boundaryAt(10)), "nothing is emitted before the effective boundary")
	}
	require.True(t, boundary.Retired(), "after closing through Until the boundary has nothing left to do")
	require.False(t, newBoundaryFixture(t).build(t).Retired(), "a boundary without Until never retires on its own")
}

func TestProductionGroupOutageRecoveryDedupeCapabilityNeedsVerifiedStorage(t *testing.T) {
	cases := map[string]struct {
		mutate func(*GroupBoundaryConfig)
		ok     bool
	}{
		"none": {func(c *GroupBoundaryConfig) {}, true},
		"receipt verified": {func(c *GroupBoundaryConfig) {
			c.Group.WritePolicy.DedupeCapability = "receipt"
			c.ReceiptTableReady = true
		}, true},
		"receipt unverified": {func(c *GroupBoundaryConfig) { c.Group.WritePolicy.DedupeCapability = "receipt" }, false},
		"unique without key": {func(c *GroupBoundaryConfig) { c.Group.WritePolicy.DedupeCapability = "unique_key" }, false},
		"unique key not unique": {func(c *GroupBoundaryConfig) {
			c.Group.WritePolicy.DedupeCapability = "unique_key"
			c.RecordKeyColumn = "record_id"
			c.Columns = append(c.Columns, dbtarget.ColumnInfo{Name: "record_id", DataType: "TEXT"})
		}, false},
		"unique key verified": {func(c *GroupBoundaryConfig) {
			c.Group.WritePolicy.DedupeCapability = "unique_key"
			c.RecordKeyColumn = "record_id"
			c.Columns = append(c.Columns, dbtarget.ColumnInfo{Name: "record_id", DataType: "TEXT", Unique: true})
		}, true},
		"unknown capability": {func(c *GroupBoundaryConfig) { c.Group.WritePolicy.DedupeCapability = "magic" }, false},
	}
	for name, c := range cases {
		f := newBoundaryFixture(t)
		f.config.Group = copyGroup(f.config.Group)
		f.config.Columns = append([]dbtarget.ColumnInfo(nil), f.config.Columns...)
		c.mutate(&f.config)
		_, err := NewGroupBoundary(f.config)
		if c.ok {
			require.NoError(t, err, name)
			continue
		}
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked, name)
		require.Equal(t, "dedupe-unsupported", blocked.Code, name)
	}
}
