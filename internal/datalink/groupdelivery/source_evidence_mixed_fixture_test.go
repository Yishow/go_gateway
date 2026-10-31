package groupdelivery

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"testing"
	"time"
)

func mixedSourceFixture(t *testing.T, skip bool) (*Store, SourceScope, map[string]SourceRef, Closure) {
	return mixedSourceFixtureMode(t, skip, false, false)
}
func mixedSourceFixtureMode(t *testing.T, skip, long, mean bool) (*Store, SourceScope, map[string]SourceRef, Closure) {
	t.Helper()
	store, _, _ := newTestStore(t)
	scope := SourceScope{WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision, ProfileKind: "mixed", ResultShape: dbtarget.MixedResultShapeV1, RecordIdentity: "8b4e4790-dd08-4cbf-8d0f-dfd036a14572", ConnectorID: testDestination.ConnectorID, ConnectorRevision: testDestination.ConnectorRevision, Wide: true, DueAtColumn: "due_at", RecordKeyColumn: "record_id", ProvenanceColumn: "provenance"}
	if long {
		scope.Wide = false
	}
	refs := map[string]SourceRef{}
	members := []snapshot.Member{}
	frozen := []any{}
	sqlMembers := []dbtarget.GroupRowMember{}
	columns := []dbtarget.ColumnInfo{{Name: "record_id", DataType: "TEXT", PrimaryKey: true}, {Name: "due_at", DataType: "TEXT"}, {Name: "provenance", DataType: "TEXT"}}
	for _, tag := range []string{"fast", "required", "slow"} {
		seconds := int64(15)
		if tag == "slow" {
			seconds = 60
		}
		required := tag == "required"
		method := snapshot.CalculationLatest
		resultType := measurement.ExactInt64
		sqlType := "INTEGER"
		if mean && tag == "fast" {
			method = snapshot.CalculationMean
			resultType = measurement.ExactDecimal
			sqlType = "TEXT"
		}
		ref := SourceRef{DeviceID: "d", PointID: "p", TagID: tag, Column: tag, Calculation: string(method), CadenceSeconds: seconds, Alignment: "utc_epoch", Required: required, SourceType: "int64", ResultType: string(resultType)}
		refs[tag] = ref
		key := SourceEntityKey("d", "p", tag)
		members = append(members, snapshot.Member{Key: key, Required: required, Calculation: method, Interval: time.Duration(seconds) * time.Second, SourceValueType: schema.DataTypeInt64})
		sqlMembers = append(sqlMembers, dbtarget.GroupRowMember{MemberKey: key, Column: tag, Type: resultType, Required: required, Calculation: method, Interval: time.Duration(seconds) * time.Second, SourceValueType: schema.DataTypeInt64})
		columns = append(columns, dbtarget.ColumnInfo{Name: tag, DataType: sqlType, Nullable: true})
		frozen = append(frozen, map[string]any{"device_id": "d", "point_id": "p", "tag_id": tag, "target_column": tag, "required": required, "source_revision": "src-1", "mapping_revision": "map-1", "source_value_type": "int64", "result_value_type": resultType, "calculation": method, "cadence": map[string]any{"interval": seconds, "unit": "second", "alignment": "utc_epoch"}})
	}

	layoutName := "managed_wide_profile_v1"
	mode := snapshot.ModeMixedWide
	if long {
		layoutName = "managed_long_v1"
		mode = snapshot.ModeMixedLong
		definition, err := recordingplan.ManagedLongGroupLayout("sqlite", "source_test", "_gw_owner_0123456789abcdef01234567")
		require.NoError(t, err)
		columns = nil
		for _, c := range definition.Columns {
			columns = append(columns, dbtarget.ColumnInfo{Name: c.Name, DataType: c.SQLType, Nullable: c.Nullable, PrimaryKey: c.PrimaryKey})
		}
	}
	raw, err := json.Marshal(map[string]any{"workspace_id": testKey.WorkspaceID, "group_id": testKey.GroupID, "applied_revision": testKey.GroupRevision, "group": map[string]any{"workspace_id": testKey.WorkspaceID, "id": testKey.GroupID, "revision": testKey.GroupRevision, "applied_revision": testKey.GroupRevision, "record_id": scope.RecordIdentity, "layout_version": layoutName, "recording_profile": map[string]any{"kind": "mixed", "result_shape": scope.ResultShape}, "destination": map[string]any{"connector_id": scope.ConnectorID, "connector_revision": scope.ConnectorRevision}, "row_policy": map[string]any{"due_at_column": scope.DueAtColumn, "record_key_column": scope.RecordKeyColumn, "provenance_column": scope.ProvenanceColumn}, "members": frozen}, "runtime_layout": map[string]any{"dialect": "sqlite", "columns": columns}})
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, raw))
	a, err := snapshot.NewAssembler(snapshot.Config{WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision, DestinationScope: testDestination.Scope, StartsAt: t0, Mode: mode, Members: members, IncompletePolicy: snapshot.IncompletePartial, Storage: snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}})
	require.NoError(t, err)
	for _, m := range members {
		if skip && m.Required {
			continue
		}
		samples := []snapshot.Sample{sampleAt(m.Key, m.Key, 10, 9007199254740993)}
		if m.Calculation == snapshot.CalculationMean {
			samples = []snapshot.Sample{sampleAt(m.Key+"-a", m.Key, 10, 0), sampleAt(m.Key+"-b", m.Key, 11, 0), sampleAt(m.Key+"-c", m.Key, 12, 1)}
		}
		for _, sample := range samples {
			require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
			a.Offer(sample, sample.ObservedAt)
		}
	}
	before, err := a.MixedCheckpoint()
	require.NoError(t, err)
	plan := a.Plan(at(15))
	require.NoError(t, plan.Err)
	after, err := a.MixedCheckpointForPlan(plan)
	require.NoError(t, err)
	var encode func(snapshot.Outcome) (dbtarget.EncodedRow, error)
	if long {
		lm := []dbtarget.LongRowMember{}
		for _, m := range sqlMembers {
			r := refs[m.Column]
			lm = append(lm, dbtarget.LongRowMember{MemberKey: m.MemberKey, DeviceID: r.DeviceID, PointID: r.PointID, TagID: r.TagID, Type: m.Type, Required: m.Required, Calculation: m.Calculation, Interval: m.Interval, SourceValueType: m.SourceValueType})
		}
		l, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{ResultShape: scope.ResultShape, WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision, ProfileKind: "mixed", Dialect: dbtarget.SQLDialectSQLite, Columns: columns, Members: lm})
		require.Empty(t, issues)
		encode = l.EncodeRow
	} else {
		l, issues := dbtarget.NewGroupRowLayout(dbtarget.GroupRowSpec{ResultShape: scope.ResultShape, WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision, Dialect: dbtarget.SQLDialectSQLite, Columns: columns, Members: sqlMembers, Partial: true, RecordKeyColumn: scope.RecordKeyColumn, DueAtColumn: scope.DueAtColumn, ProvenanceColumn: scope.ProvenanceColumn})
		require.Empty(t, issues)
		encode = l.EncodeRow
	}
	c := Closure{Destination: testDestination, Mixed: &MixedClosureProgress{Before: before, After: after, DescriptorDigest: handoffHash(string(raw))}}
	for _, o := range plan.Outcomes {
		b := ClosedBucket{Outcome: o}
		if o.Kind == snapshot.OutcomeRow {
			row, e := encode(o)
			require.NoError(t, e)
			b.Row = &row
		}
		c.Buckets = append(c.Buckets, b)
	}
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	return store, scope, refs, c
}
