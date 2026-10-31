package grouppipeline

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// Reuse the production Store/fence/version fixture; only its initial descriptors
// change before any accepted input. Both sides retain real mixed profile kinds.
func newMixedCutoverPipelineFixture(t *testing.T) *ownershipPipelineFixture {
	t.Helper()
	f := newHandoffPipelineFixture(t, shapeMeanWide)
	f.boundary = at(23)
	f.target.EffectiveAt = f.boundary
	for i, s := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
		g := s.Group
		g.RecordID = []string{"11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222"}[i]
		g.LayoutVersion = workspace.WriteGroupLayoutManagedWideProfileV1
		g.StorageMode = workspace.WriteGroupStorageModePlain
		g.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileMixed, ResultShape: workspace.WriteGroupResultShapeMemberCalculationV1}
		g.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyManaged
		g.RowPolicy.RecordKeyColumn = "record_id"
		g.RowPolicy.ProvenanceColumn = "provenance"
		g.RowPolicy.IntervalSeconds = 0
		g.RowPolicy.BucketStartColumn = ""
		g.RowPolicy.DueAtColumn = "due_at"
		g.RowPolicy.IncompletePolicy = "partial"
		g.RowPolicy.AllowedLatenessSeconds = 0
		g.Destination.SchemaRevision = "confirmed-operation"
		g.Destination.SchemaDigest = strings.Repeat("b", 64)
		types := map[string]schema.DataType{}
		for j := range g.Members {
			m := &g.Members[j]
			m.Calculation = workspace.WriteGroupMemberCalculationMean
			m.Cadence = &workspace.WriteGroupMemberCadence{Interval: []int64{15, 60}[j], Unit: workspace.WriteGroupIntervalUnitSecond, Alignment: workspace.WriteGroupAlignmentUTCEpoch}
			m.SourceValueType = schema.DataTypeUint64
			m.ResultValueType = measurement.ExactDecimal
			m.ValueRole = workspace.WriteGroupValueRoleGauge
			m.Required = false
			m.EntityKey = ""
			types[m.TagID] = schema.DataTypeUint64
		}
		expected, err := workspace.MixedManagedSchemaLayout(g, types, "sqlite")
		require.NoError(t, err)
		s.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectSQLite, TagTypes: types, SchemaDigest: strings.Repeat("a", 64), ReceiptTableReady: true}
		for _, c := range expected.Columns {
			s.RuntimeLayout.Columns = append(s.RuntimeLayout.Columns, dbtarget.ColumnInfo{Name: c.Name, DataType: c.SQLType, Nullable: c.Nullable, PrimaryKey: c.PrimaryKey})
		}
		raw, err := json.Marshal(s)
		require.NoError(t, err)
		_, err = f.local.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=? WHERE group_id=? AND group_revision=?`, string(raw), s.GroupID, s.AppliedRevision)
		require.NoError(t, err)
	}
	f.groups.groups = []*workspace.WriteGroup{f.source.Group}
	f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{f.source.GroupID: {f.source}}
	f.pipe = New(f.pipe.deps, f.pipe.config)
	f.clock = at(10)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.NotNil(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)])
	return f
}

func TestMixedCutoverPipelineRestartKeepsMemberWindowsAndExactOwner(t *testing.T) {
	for _, mode := range []string{"restart", "live", "cached-target", "disabled-source"} {
		t.Run(mode, func(t *testing.T) { mixedCutoverKeepsWindows(t, mode) })
	}
}
func mixedCutoverKeepsWindows(t *testing.T, mode string) {
	f := newMixedCutoverPipelineFixture(t)
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "first", at(10), uint64(2))))
	slow := *f.source.Group
	slow.Members = slow.Members[1:]
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(&slow, "slow-before", at(10), uint64(6))))
	f.clock = at(15)
	f.pipe.TickAll(t.Context())
	before, err := f.store.Restore(t.Context(), groupKey(f.source))
	require.NoError(t, err)
	require.Equal(t, at(15), before.Mixed.Progress.ClosedThrough)
	f.clock = at(21)
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "exact", at(23), uint64(7))))
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "before", at(22), uint64(4))))
	if mode == "disabled-source" {
		f.groups.groups[0] = liveGroup(f.source.Group, workspace.WriteGroupStatusDisabled)
		require.NoError(t, f.pipe.Reconcile(t.Context()))
		require.True(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].localDrain.Load())
	}
	f.commit()
	f.source.EndsAt = f.boundary
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	if mode == "restart" {
		f.pipe = New(f.pipe.deps, f.pipe.config)
	}
	if mode == "cached-target" {
		m, reason := f.pipe.build(t.Context(), f.target)
		require.Empty(t, reason)
		f.pipe.installManaged(m)
	}
	f.clock = at(24)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	old := f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)]
	next := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)]
	require.NotNil(t, old)
	require.NotNil(t, next)
	if mode == "disabled-source" {
		require.True(t, old.localDrain.Load(), "handoff rebuild must preserve disabled local intake fence")
	}
	require.True(t, old.boundary.NextClose().IsZero())
	require.True(t, next.boundary.NextClose().IsZero())
	restored, err := f.store.Restore(t.Context(), groupKey(f.target))
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, "exact", restored.Samples[0].SampleID)
	f.clock = at(60)
	f.pipe.TickAll(t.Context())
	targetAfter, err := f.store.Restore(t.Context(), groupKey(f.target))
	require.NoError(t, err)
	require.NotNil(t, targetAfter.Mixed)
	require.Equal(t, at(60), targetAfter.Mixed.Progress.ClosedThrough)
	require.Empty(t, targetAfter.Samples)
	r, err := f.store.Restore(t.Context(), groupKey(f.source))
	require.NoError(t, err)
	require.Equal(t, at(23), r.Mixed.EndsAt)
	require.True(t, r.NextClose.IsZero())
	rows, err := f.local.QueryContext(t.Context(), `SELECT due_at,bucket_start,payload FROM wg_delivery_outbox WHERE group_id=?`, f.source.GroupID)
	require.NoError(t, err)
	defer rows.Close()
	slowPartial := false
	for rows.Next() {
		var due, bucket string
		var raw []byte
		require.NoError(t, rows.Scan(&due, &bucket, &raw))
		require.Empty(t, bucket)
		row, err := groupdelivery.DecodeRowPayload(raw)
		require.NoError(t, err)
		if row.DueAt.Equal(at(60)) {
			slowPartial = true
			for _, cell := range row.Cells {
				if cell.Column != "provenance" {
					continue
				}
				var facts []struct {
					Interval int64 `json:"interval_seconds"`
					Window   *struct {
						ActualEnd time.Time `json:"actual_end"`
					} `json:"window"`
				}
				require.NoError(t, json.Unmarshal([]byte(cell.Value.(string)), &facts))
				for _, fact := range facts {
					if fact.Interval == 60 {
						require.NotNil(t, fact.Window)
						require.Equal(t, at(23), fact.Window.ActualEnd)
					}
				}
			}
		}
	}
	require.NoError(t, rows.Err())
	require.True(t, slowPartial)
}

func TestMixedCutoverPipelineRejectsForeignFenceAndChangedIntent(t *testing.T) {
	for _, mode := range []string{"foreign-fence", "changed-cadence", "changed-method"} {
		t.Run(mode, func(t *testing.T) {
			f := newMixedCutoverPipelineFixture(t)
			if mode != "foreign-fence" {
				if mode == "changed-cadence" {
					f.target.Group.Members[0].Cadence.Interval = 30
				} else {
					f.target.Group.Members[0].Calculation = workspace.WriteGroupMemberCalculationLatest
					f.target.Group.Members[0].ResultValueType = measurement.ExactUint64
				}
				raw, err := json.Marshal(f.target)
				require.NoError(t, err)
				_, err = f.local.ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(raw), f.target.GroupID)
				require.NoError(t, err)
			}
			f.clock = at(21)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "exact", at(23), uint64(7))))
			f.commit()
			if mode == "foreign-fence" {
				_, err := f.local.ExecContext(t.Context(), `UPDATE write_group_ownership_fences SET connector_revision='foreign' WHERE direction='started'`)
				require.NoError(t, err)
			}
			f.source.EndsAt = f.boundary
			f.groups.groups = append(f.groups.groups, f.target.Group)
			f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
			f.clock = at(24)
			err := f.pipe.Reconcile(t.Context())
			if mode != "foreign-fence" {
				require.ErrorIs(t, err, groupdelivery.ErrHandoffProofConflict)
			} else {
				require.ErrorIs(t, err, groupdelivery.ErrOwnershipUnavailable)
			}
			var copied int
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=?`, f.target.GroupID).Scan(&copied))
			require.Zero(t, copied, "refused proof must not transfer accepted sample")
		})
	}
}
