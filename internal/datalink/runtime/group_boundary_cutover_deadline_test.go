package runtime

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// Both original descriptors are frozen before acceptance. The supported
// instantaneous 15-second profiles have actual successor window [10,15),
// lateness 1 and deadline 16; counter_raw preserves proof without inventing a
// five-second cumulative profile or backdating a public cutover.
func deadlineFixture(t *testing.T, wide bool) (*runtimeHandoffFixture, func()) {
	t.Helper()
	source, target := wideProfileFixture(t), newLongFixture(t)
	if wide {
		target = wideProfileFixture(t)
		target.config.Group.ID = "deadline-target-wide"
	}
	source.members = source.members[2:3]
	source.config.Group.Members = source.members
	target.members = append([]workspace.WriteGroupMember(nil), source.members...)
	target.config.Group.Members = target.members
	target.config.TagTypes, target.config.Counters = source.config.TagTypes, source.config.Counters
	for _, owner := range []*longFixture{source, target} {
		owner.config.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
		owner.config.Group.RowPolicy.AllowedLatenessSeconds = 1
	}
	source.clock, target.clock = at(9*time.Second), at(9*time.Second)
	path := filepath.Join(t.TempDir(), "deadline.db")
	store := groupdelivery.NewStore(openDurableDB(t, path)).WithClock(func() time.Time { return source.clock })
	f := &runtimeHandoffFixture{source: source, target: target, store: store, boundary: at(10 * time.Second), databasePath: path}
	f.sourceKey = freezeHandoffRuntime(t, store, source, boundaryStart)
	f.targetKey = freezeHandoffRuntime(t, store, target, f.boundary)
	target.config.StartsAt = f.boundary
	f.old, f.emptyTarget = source.build(t), target.build(t)
	return f, func() {
		require.Equal(t, at(9500*time.Millisecond), source.clock, "cutover commits before its future boundary10")
		tx, err := store.DB().BeginTx(t.Context(), nil)
		require.NoError(t, err)
		defer tx.Rollback()
		runtimeHandoffVersion(t, tx, source, boundaryStart)
		runtimeHandoffVersion(t, tx, target, f.boundary)
		m := source.members[0]
		detail, err := json.Marshal(map[string]any{"ownership_fence_revision": "runtime-deadline",
			"ended":   []workspace.CutoverEnd{{GroupID: f.sourceKey.GroupID, Revision: f.sourceKey.GroupRevision}},
			"started": []workspace.CutoverStart{{GroupID: f.targetKey.GroupID, Revision: f.targetKey.GroupRevision}},
			"moved":   []workspace.CutoverMemberRef{{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}}})
		require.NoError(t, err)
		_, err = tx.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('runtime-deadline',?,'forward',?,?,'deadline-review',?,'{}')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, f.boundary.Format(time.RFC3339Nano), string(detail))
		require.NoError(t, err)
		_, err = tx.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,'runtime-deadline')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, f.sourceKey.GroupRevision, f.boundary.Format(time.RFC3339Nano))
		require.NoError(t, err)
		for direction, owner := range map[string]*longFixture{"ended": source, "started": target} {
			g := owner.config.Group
			_, err = tx.Exec(`INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('runtime-deadline',?,?,?,?,?,?,?,?,?,?,?,?,?)`, g.WorkspaceID, direction, g.ID, g.AppliedRevision, m.DeviceID, m.PointID, m.TagID, string(g.RecordingProfile.Kind), g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.Database, g.Destination.TableSchema, g.Destination.TableName)
			require.NoError(t, err)
		}
		require.NoError(t, tx.Commit())
	}
}

func deadlineBoundary(t *testing.T, f *runtimeHandoffFixture, key groupdelivery.GroupKey) *GroupBoundary {
	t.Helper()
	raw, err := f.store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	var frozen workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &frozen))
	cfg := GroupBoundaryConfig{Group: frozen.Group, TagTypes: frozen.RuntimeLayout.TagTypes, CounterBindings: frozen.CounterBindings,
		Columns: frozen.RuntimeLayout.Columns, Dialect: frozen.RuntimeLayout.Dialect, ReceiptTableReady: frozen.RuntimeLayout.ReceiptTableReady,
		FirstBucket: boundaryStart, StartsAt: frozen.EffectiveAt, MaxFutureSkew: 2 * time.Second, Ledger: f.store.Ledger(key), Clock: func() time.Time { return f.source.clock }}
	if key == f.sourceKey {
		cfg.EndsAt, cfg.Until = f.boundary, at(15*time.Second)
	}
	b, err := NewGroupBoundaryContext(t.Context(), cfg)
	require.NoError(t, err)
	return b
}

func deadlineReopen(t *testing.T, f *runtimeHandoffFixture) {
	t.Helper()
	require.NoError(t, f.store.DB().Close())
	f.old, f.emptyTarget = nil, nil
	f.source.config.Ledger, f.target.config.Ledger = nil, nil
	f.source.config.Counters, f.target.config.Counters = nil, nil
	f.store = groupdelivery.NewStore(openDurableDB(t, f.databasePath)).WithClock(func() time.Time { return f.source.clock })
}

func deadlineSQL(t *testing.T, f *runtimeHandoffFixture, sample string) {
	t.Helper()
	frozen := reductionOutbox(t, f.store.DB(), f.targetKey.GroupID)
	row, err := groupdelivery.DecodeRowPayload([]byte(frozen.payload))
	require.NoError(t, err)
	_, err = dbtarget.InsertGroupRow(t.Context(), f.target.db, dbtarget.GroupInsertRequest{Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: frozen.table, Row: row, PayloadDigest: frozen.digest, Strategy: dbtarget.GroupEffectReceipt, CommittedAt: f.source.clock.Format(time.RFC3339Nano)})
	require.NoError(t, err)
	var value, provenance, observed string
	args := []any{&value, &provenance}
	query := `SELECT v_raw,provenance FROM ` + wideFixtureTable
	if f.target.config.Group.LayoutVersion == workspace.WriteGroupLayoutManagedLongV1 {
		query = `SELECT value_decimal,provenance,observed_at FROM ` + longFixtureTable
		args = append(args, &observed)
	}
	require.NoError(t, f.target.db.QueryRow(query).Scan(args...))
	require.Equal(t, "200", value)
	var evidence []struct {
		SampleID   string                                 `json:"sample_id"`
		ObservedAt string                                 `json:"observed_at"`
		Window     struct{ Start, End, Truncated string } `json:"window"`
	}
	if f.target.config.Group.LayoutVersion == workspace.WriteGroupLayoutManagedLongV1 {
		provenance = "[" + provenance + "]"
	}
	require.NoError(t, json.Unmarshal([]byte(provenance), &evidence))
	require.Len(t, evidence, 1)
	require.Equal(t, sample, evidence[0].SampleID)
	if observed == "" {
		observed = evidence[0].ObservedAt
	}
	observedAt, err := time.Parse(time.RFC3339Nano, observed)
	require.NoError(t, err)
	require.Equal(t, at(11*time.Second), observedAt)
	require.Equal(t, at(10*time.Second).Format(longTimeLayout), evidence[0].Window.Start)
	require.Equal(t, at(15*time.Second).Format(longTimeLayout), evidence[0].Window.End)
	require.Equal(t, "start", evidence[0].Window.Truncated)
}

func TestGroupBoundaryCutoverDeadlineDelayedHandoffRemainsEligible(t *testing.T) {
	for _, wide := range []bool{false, true} {
		name := "long"
		if wide {
			name = "wide"
		}
		t.Run(name, func(t *testing.T) {
			f, commit := deadlineFixture(t, wide)
			input := f.source.envelope(0, "s-1", at(11*time.Second), uint64(200))
			input.ReceivedAt = at(9 * time.Second)
			require.NoError(t, f.old.AcceptOwnedSample(t.Context(), input), "future-skew allowance accepts11 at gateway9 before cutover")
			original := reductionOriginal(t, f.store.DB(), f.sourceKey.GroupID, "s-1")
			var observation struct {
				ObservedAt time.Time `json:"observed_at"`
				ReceivedAt time.Time `json:"received_at"`
			}
			require.NoError(t, json.Unmarshal([]byte(original[0]), &observation))
			require.Equal(t, at(11*time.Second), observation.ObservedAt)
			require.Equal(t, at(9*time.Second), observation.ReceivedAt)
			accepted, err := time.Parse(time.RFC3339Nano, original[2])
			require.NoError(t, err)
			require.Equal(t, at(9*time.Second), accepted)
			f.source.clock = at(9500 * time.Millisecond)
			commit()
			f.source.clock = at(16 * time.Second)
			source, target := deadlineBoundary(t, f, f.sourceKey), deadlineBoundary(t, f, f.targetKey)
			before := crashProgress(t, f.store.DB())
			require.ErrorIs(t, source.Tick(t.Context(), f.source.clock), groupdelivery.ErrOwnershipConflict)
			require.ErrorIs(t, target.Tick(t.Context(), f.source.clock), groupdelivery.ErrOwnershipConflict, "deadline16 cannot close an unresolved on-time original")
			require.Equal(t, before, crashProgress(t, f.store.DB()))
			source, target = nil, nil
			f.source.clock = at(20 * time.Second)
			deadlineReopen(t, f)
			target = deadlineBoundary(t, f, f.targetKey)
			require.ErrorIs(t, target.Tick(t.Context(), f.source.clock), groupdelivery.ErrOwnershipConflict)
			batch, err := f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
			require.NoError(t, err)
			require.Equal(t, 1, batch.Accepted)
			require.Equal(t, original, reductionOriginal(t, f.store.DB(), f.sourceKey.GroupID, "s-1"))
			require.Equal(t, original, reductionOriginal(t, f.store.DB(), f.targetKey.GroupID, "s-1"), "handoff20 preserves proof/payload/digest and actual durable accepted_at9")
			require.ErrorIs(t, target.Tick(t.Context(), f.source.clock), groupdelivery.ErrOwnershipConflict, "transfer completion requires a fresh receipt view, not a token update on an empty assembler")
			source, target = deadlineBoundary(t, f, f.sourceKey), deadlineBoundary(t, f, f.targetKey)
			require.Zero(t, target.assembler.Late())
			require.NoError(t, source.Tick(t.Context(), f.source.clock))
			require.NoError(t, target.Tick(t.Context(), f.source.clock))
			var oldConsumed, newConsumed int
			require.NoError(t, f.store.DB().QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=?`, f.sourceKey.GroupID).Scan(&oldConsumed))
			require.NoError(t, f.store.DB().QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=?`, f.targetKey.GroupID).Scan(&newConsumed))
			require.Zero(t, oldConsumed)
			require.Equal(t, 1, newConsumed)
			require.Len(t, crashRows(t, f.store.DB(), "wg_delivery_outbox"), 2)
			oldRow, err := groupdelivery.DecodeRowPayload([]byte(reductionOutbox(t, f.store.DB(), f.sourceKey.GroupID).payload))
			require.NoError(t, err)
			oldCells := map[string]any{}
			for _, cell := range oldRow.Cells {
				oldCells[cell.Column] = cell.Value
			}
			require.Contains(t, oldCells, "v_raw")
			require.Nil(t, oldCells["v_raw"], "existing partial policy may emit a missing old row, never the successor sample")
			deadlineSQL(t, f, "s-1")
			final := crashProgress(t, f.store.DB())
			deadlineReopen(t, f)
			target = deadlineBoundary(t, f, f.targetKey)
			// At20 this is a verified original resend, not a first arrival at20.
			require.NoError(t, target.AcceptOwnedSample(t.Context(), input))
			require.NoError(t, target.Tick(t.Context(), f.source.clock))
			require.Equal(t, final, crashProgress(t, f.store.DB()))
			deadlineSQL(t, f, "s-1")
		})
	}
}

func TestGroupBoundaryCutoverDeadlineFirstArrivalAfterFinalizationIsLate(t *testing.T) {
	for _, wide := range []bool{false, true} {
		name := "long"
		if wide {
			name = "wide"
		}
		t.Run(name, func(t *testing.T) {
			f, commit := deadlineFixture(t, wide)
			f.source.clock = at(9500 * time.Millisecond)
			commit()
			f.source.clock = at(12 * time.Second)
			target := deadlineBoundary(t, f, f.targetKey)
			input := f.target.envelope(0, "finalized", at(11*time.Second), uint64(200))
			input.ReceivedAt = f.source.clock
			require.NoError(t, target.AcceptOwnedSample(t.Context(), input))
			f.source.clock = at(16 * time.Second)
			require.NoError(t, target.Tick(t.Context(), f.source.clock), "actually finalize at the existing deadline before sending a new late input")
			require.Equal(t, at(15*time.Second), target.NextClose())
			deadlineSQL(t, f, "finalized")
			final := crashProgress(t, f.store.DB())
			f.source.clock = at(17 * time.Second)
			late := f.target.envelope(0, "first-at17", at(13*time.Second), uint64(999))
			late.ReceivedAt = f.source.clock
			var refused *GroupSampleError
			require.ErrorAs(t, target.AcceptOwnedSample(t.Context(), late), &refused)
			require.Equal(t, snapshot.OfferLate, refused.Outcome)
			require.Equal(t, snapshot.ReasonLate, refused.Reason)
			require.EqualValues(t, 1, target.assembler.Late())
			require.NoError(t, target.Tick(t.Context(), f.source.clock))
			require.Equal(t, final, crashProgress(t, f.store.DB()), "no ACK, append, checkpoint movement, effect rewrite or reopened finalized window")
			f.source.clock = at(20 * time.Second)
			target = nil
			deadlineReopen(t, f)
			require.NoError(t, deadlineBoundary(t, f, f.targetKey).Tick(t.Context(), f.source.clock))
			require.Equal(t, final, crashProgress(t, f.store.DB()))
			deadlineSQL(t, f, "finalized")
		})
	}
}
