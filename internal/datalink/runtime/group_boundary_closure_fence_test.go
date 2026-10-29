package runtime

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// Establish both original frozen descriptors before any acceptance. The returned
// transaction is the deterministic cutover barrier, never a descriptor rewrite.
func closureFenceFixture(t *testing.T, wideTarget bool) (*runtimeHandoffFixture, func()) {
	t.Helper()
	source, target := wideCounterFixture(t), newCounterFixture(t)
	if wideTarget {
		target = wideCounterFixture(t)
		target.config.Group.ID = "closure-next-wide"
	}
	source.members = source.members[:1]
	source.config.Group.Members = source.members
	target.members = append([]workspace.WriteGroupMember(nil), source.members...)
	if !wideTarget {
		target.members[0].SourceValueType, target.members[0].ResultValueType = "", ""
	}
	target.config.Group.Members = target.members
	target.config.TagTypes, target.config.Counters = source.config.TagTypes, source.config.Counters
	source.clock, target.clock = at(11*time.Second), at(11*time.Second)
	path := filepath.Join(t.TempDir(), "closure-fence.db")
	store := groupdelivery.NewStore(openDurableDB(t, path)).WithClock(func() time.Time { return source.clock })
	f := &runtimeHandoffFixture{source: source, target: target, store: store, boundary: at(12500 * time.Millisecond), databasePath: path}
	f.sourceKey = freezeHandoffRuntime(t, store, source, boundaryStart)
	f.targetKey = freezeHandoffRuntime(t, store, target, f.boundary)
	target.config.StartsAt = f.boundary
	f.old, f.emptyTarget = source.build(t), target.build(t)
	return f, func() {
		tx, err := store.DB().BeginTx(t.Context(), nil)
		require.NoError(t, err)
		defer tx.Rollback()
		runtimeHandoffVersion(t, tx, source, boundaryStart)
		runtimeHandoffVersion(t, tx, target, f.boundary)
		m := source.members[0]
		detail, err := json.Marshal(map[string]any{
			"ownership_fence_revision": "runtime-closure",
			"ended":                    []workspace.CutoverEnd{{GroupID: f.sourceKey.GroupID, Revision: f.sourceKey.GroupRevision}},
			"started":                  []workspace.CutoverStart{{GroupID: f.targetKey.GroupID, Revision: f.targetKey.GroupRevision}},
			"moved":                    []workspace.CutoverMemberRef{{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}},
		})
		require.NoError(t, err)
		_, err = tx.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('runtime-closure',?,'forward',?,?,'closure-review',?,'{}')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, f.boundary.Format(time.RFC3339Nano), string(detail))
		require.NoError(t, err)
		_, err = tx.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,'runtime-closure')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, f.sourceKey.GroupRevision, f.boundary.Format(time.RFC3339Nano))
		require.NoError(t, err)
		for direction, owner := range map[string]*longFixture{"ended": source, "started": target} {
			g := owner.config.Group
			_, err = tx.Exec(`INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('runtime-closure',?,?,?,?,?,?,?,?,?,?,?,?,?)`, g.WorkspaceID, direction, g.ID, g.AppliedRevision, m.DeviceID, m.PointID, m.TagID, string(g.RecordingProfile.Kind), g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.Database, g.Destination.TableSchema, g.Destination.TableName)
			require.NoError(t, err)
		}
		require.NoError(t, tx.Commit())
	}
}

type closureFencePause struct {
	*frozenFixtureLedger
	reached, resume chan struct{}
	once            sync.Once
	beforeRestore   bool
}

func (l *closureFencePause) wait() { l.once.Do(func() { close(l.reached); <-l.resume }) }
func (l *closureFencePause) Restore(ctx context.Context) (groupdelivery.Restored, error) {
	if l.beforeRestore {
		l.wait()
	}
	return l.frozenFixtureLedger.Restore(ctx)
}
func (l *closureFencePause) CommitClosure(ctx context.Context, closure groupdelivery.Closure) error {
	if !l.beforeRestore {
		l.wait()
	}
	return l.frozenFixtureLedger.CommitClosure(ctx, closure)
}

// Include every stored column, not just counts: failed closure must not change
// checkpoint timestamps, consumed flags, payload/proof/digest or baseline bytes.
func closureFenceState(t *testing.T, f *runtimeHandoffFixture) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for _, table := range []string{"wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_counter_baselines"} {
		rows, err := f.store.DB().Query(`SELECT * FROM `+table+` WHERE group_id=? ORDER BY rowid`, f.sourceKey.GroupID)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		state[table] = nil
		for rows.Next() {
			values := make([]any, len(columns))
			refs := make([]any, len(columns))
			for i := range values {
				refs[i] = &values[i]
			}
			require.NoError(t, rows.Scan(refs...))
			raw, err := json.Marshal(values)
			require.NoError(t, err)
			state[table] = append(state[table], string(raw))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	var receipts int
	require.NoError(t, f.store.DB().QueryRow(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs`).Scan(&receipts))
	require.Zero(t, receipts, "general fence must not depend on transfer receipts")
	return state
}

func TestGroupBoundaryClosureFenceCutoverWinsWithoutHandoff(t *testing.T) {
	for _, wide := range []bool{false, true} {
		layout := "long"
		if wide {
			layout = "wide"
		}
		for _, accepted := range []bool{false, true} {
			intake := "empty"
			if accepted {
				intake = "pre-boundary-only"
			}
			t.Run(layout+"/"+intake, func(t *testing.T) {
				f, cutover := closureFenceFixture(t, wide)
				if accepted {
					require.NoError(t, f.old.AcceptSample(t.Context(), f.source.envelope(0, "before", at(11*time.Second), uint64(100))))
				}
				before := closureFenceState(t, f)
				oldReplay := f.old.ReplayState()
				pause := &closureFencePause{frozenFixtureLedger: f.old.ledger.(*frozenFixtureLedger), reached: make(chan struct{}), resume: make(chan struct{})}
				f.old.ledger = pause
				var release sync.Once
				defer release.Do(func() { close(pause.resume) })
				done := make(chan error, 1)
				go func() { done <- f.old.Tick(t.Context(), at(15*time.Second)) }()
				select {
				case <-pause.reached:
				case err := <-done:
					t.Fatalf("closure missed plan barrier: %v", err)
				}
				cutover() // plan already contains unbounded [0,15); durable owner now ends at 12.5
				release.Do(func() { close(pause.resume) })
				require.ErrorIs(t, <-done, groupdelivery.ErrOwnershipConflict)
				require.Equal(t, before, closureFenceState(t, f), "stale plan cannot advance any durable state")
				require.Equal(t, oldReplay, f.old.ReplayState(), "failed closure leaves the in-memory cursor and captured view untouched")
				f.old.SetEndsAt(f.boundary)
				require.ErrorIs(t, f.old.Tick(t.Context(), at(15*time.Second)), groupdelivery.ErrOwnershipConflict, "editing cache bounds cannot upgrade the captured owner view")
				require.Equal(t, before, closureFenceState(t, f))
				f.source.config.EndsAt, f.source.config.Until = f.boundary, at(15*time.Second)
				fresh := f.source.build(t)
				staleReplay, freshReplay := f.old.ReplayState(), fresh.ReplayState()
				require.Equal(t, staleReplay.NextClose, freshReplay.NextClose)
				require.Equal(t, staleReplay.Until, freshReplay.Until)
				require.Equal(t, staleReplay.HandoffView, freshReplay.HandoffView)
				require.NotEqual(t, staleReplay.OwnershipView, freshReplay.OwnershipView, "receipt-free ownership refresh must be visible despite equal cache bounds")
				require.NoError(t, fresh.Tick(t.Context(), at(15*time.Second)), "aligned NextClose > actual end is a legal partial closure")
				restored, err := f.store.Restore(t.Context(), f.sourceKey)
				require.NoError(t, err)
				require.Equal(t, at(15*time.Second), restored.NextClose)
				require.Empty(t, restored.Samples)
				if accepted {
					require.Len(t, restored.Baselines, 1)
					checkpoint, err := aggregation.DecodeCounterCheckpoint(restored.Baselines[0].Payload)
					require.NoError(t, err)
					require.Equal(t, "before", checkpoint.Reading.SampleID)
					require.Equal(t, f.boundary, checkpoint.WindowEnd)
					reductionReadback(t, f.source, reductionOutbox(t, f.store.DB(), f.sourceKey.GroupID), "100", "end", f.boundary)
				} else {
					require.Empty(t, restored.Baselines)
				}
				frozen := closureFenceState(t, f)
				require.NoError(t, f.source.build(t).Tick(t.Context(), at(30*time.Second)))
				require.Equal(t, frozen, closureFenceState(t, f), "restart never reopens a finalized partial window")
			})
		}
	}
}

func TestGroupBoundaryClosureFenceNewViewCannotAuthorizeOldBounds(t *testing.T) {
	for _, wide := range []bool{false, true} {
		layout := "long"
		if wide {
			layout = "wide"
		}
		t.Run(layout, func(t *testing.T) {
			f, cutover := closureFenceFixture(t, wide)
			require.NoError(t, f.old.AcceptSample(t.Context(), f.source.envelope(0, "before", at(11*time.Second), uint64(100))))
			before := closureFenceState(t, f)
			staleConfig := f.source.config // canonical descriptor/bounds were read before cutover
			pause := &closureFencePause{frozenFixtureLedger: staleConfig.Ledger.(*frozenFixtureLedger), reached: make(chan struct{}), resume: make(chan struct{}), beforeRestore: true}
			staleConfig.Ledger = pause
			var release sync.Once
			defer release.Do(func() { close(pause.resume) })
			type result struct {
				boundary *GroupBoundary
				err      error
			}
			done := make(chan result, 1)
			go func() { b, err := NewGroupBoundaryContext(t.Context(), staleConfig); done <- result{b, err} }()
			select {
			case <-pause.reached:
			case got := <-done:
				t.Fatalf("constructor missed Restore barrier: %v", got.err)
			}
			cutover()
			release.Do(func() { close(pause.resume) }) // Restore must see new authority, but config is still old
			got := <-done
			require.ErrorIs(t, got.err, groupdelivery.ErrOwnershipConflict)
			require.Nil(t, got.boundary, "a fresh epoch token must never bless stale assembler bounds")
			require.Equal(t, before, closureFenceState(t, f))
			f.source.config.EndsAt, f.source.config.Until = f.boundary, at(15*time.Second)
			require.NoError(t, f.source.build(t).Tick(t.Context(), at(15*time.Second)))
			reductionReadback(t, f.source, reductionOutbox(t, f.store.DB(), f.sourceKey.GroupID), "100", "end", f.boundary)
		})
	}
}
