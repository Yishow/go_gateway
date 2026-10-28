package runtime

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type runtimeHandoffFixture struct {
	source, target       *longFixture
	store                *groupdelivery.Store
	sourceKey, targetKey groupdelivery.GroupKey
	old, emptyTarget     *GroupBoundary
	boundary             time.Time
	databasePath         string
}

// These are complete frozen fixtures established before intake, not public
// Apply wiring. The transfer path must use their original proofs unchanged.
func freezeHandoffRuntime(t *testing.T, store *groupdelivery.Store, f *longFixture, effective time.Time) groupdelivery.GroupKey {
	t.Helper()
	g := f.config.Group
	key := groupdelivery.GroupKey{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.AppliedRevision}
	f.config.Ledger = &frozenFixtureLedger{Ledger: store.Ledger(key), store: store}
	freezeDurableFixture(t, f)
	raw, err := store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	var frozen workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &frozen))
	frozen.EffectiveAt = effective
	frozen.RuntimeLayout.Columns = f.config.Columns
	frozen.RuntimeLayout.Dialect = f.config.Dialect
	frozen.RuntimeLayout.ReceiptTableReady = f.config.ReceiptTableReady
	raw, err = json.Marshal(frozen)
	require.NoError(t, err)
	// Initial fixture construction only: neither owner has accepted any input.
	_, err = store.DB().Exec(`UPDATE wg_runtime_versions SET payload=? WHERE workspace_id=? AND group_id=? AND group_revision=?`, string(raw), key.WorkspaceID, key.GroupID, key.GroupRevision)
	require.NoError(t, err)
	return key
}

func runtimeHandoffVersion(t *testing.T, tx *sql.Tx, f *longFixture, effective time.Time) {
	t.Helper()
	g := f.config.Group
	_, err := tx.Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,status,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,?,?,?,'ready',?,?,?)`, g.ID, g.WorkspaceID, g.Revision, g.AppliedRevision, g.ID, g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.TableName)
	require.NoError(t, err)
	raw, err := json.Marshal(g)
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES(?,?,?,?)`, g.ID, g.AppliedRevision, effective, string(raw))
	require.NoError(t, err)
}

func newRuntimeHandoffFixture(t *testing.T) *runtimeHandoffFixture {
	t.Helper()
	return newRuntimeHandoffFixtureWithLegacy(t, false)
}

func newRuntimeHandoffFixtureWithLegacy(t *testing.T, legacy bool) *runtimeHandoffFixture {
	t.Helper()
	source, target := wideCounterFixture(t), newCounterFixture(t)
	source.members = source.members[:1]
	source.config.Group.Members = source.members
	target.members = append([]workspace.WriteGroupMember(nil), source.members...)
	// Cumulative long uses the existing role-free member shape; its input type
	// remains proven by frozen TagTypes and the unchanged counter facts.
	target.members[0].SourceValueType, target.members[0].ResultValueType = "", ""
	target.config.Group.Members = target.members
	target.config.TagTypes = source.config.TagTypes
	target.config.Counters = source.config.Counters
	source.clock, target.clock = at(11*time.Second), at(11*time.Second)
	databasePath := filepath.Join(t.TempDir(), "handoff.db")
	store := groupdelivery.NewStore(openDurableDB(t, databasePath)).WithClock(func() time.Time { return source.clock })
	boundary := at(12500 * time.Millisecond)
	f := &runtimeHandoffFixture{source: source, target: target, store: store, boundary: boundary, databasePath: databasePath}
	if legacy {
		// Original pre-proof journal format: publish the unchanged original
		// descriptor afterward, as in the existing legacy association fixtures.
		g, m := source.config.Group, source.members[0]
		key := groupdelivery.GroupKey{WorkspaceID: g.WorkspaceID, GroupID: g.ID, GroupRevision: g.AppliedRevision}
		memberJSON, err := json.Marshal([]string{m.DeviceID, m.PointID, m.TagID})
		require.NoError(t, err)
		sample := snapshot.Sample{SampleID: "counter-after", MemberKey: string(memberJSON), ObservedAt: at(13 * time.Second), ReceivedAt: at(13 * time.Second).Add(40 * time.Millisecond), Value: measurement.NewUint64(200), Quality: schema.QualityGood, SourceRevision: m.SourceRevision, MappingRevision: m.MappingRevision}
		require.NoError(t, store.AppendSample(t.Context(), key, boundaryStart, sample))
	}
	f.sourceKey = freezeHandoffRuntime(t, store, source, boundaryStart)
	f.targetKey = freezeHandoffRuntime(t, store, target, boundary)
	target.config.StartsAt = boundary
	f.old, f.emptyTarget = source.build(t), target.build(t)
	if !legacy {
		require.NoError(t, f.old.AcceptSample(t.Context(), source.envelope(0, "counter-after", at(13*time.Second), uint64(200))))
	}
	require.NoError(t, f.old.AcceptSample(t.Context(), source.envelope(0, "counter-before", boundary.Add(-time.Millisecond), uint64(100))))
	tx, err := store.DB().BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	runtimeHandoffVersion(t, tx, source, boundaryStart)
	runtimeHandoffVersion(t, tx, target, boundary)
	m := source.members[0]
	detail, err := json.Marshal(map[string]any{
		"ownership_fence_revision": "runtime-handoff",
		"ended":                    []workspace.CutoverEnd{{GroupID: f.sourceKey.GroupID, Revision: f.sourceKey.GroupRevision}},
		"started":                  []workspace.CutoverStart{{GroupID: f.targetKey.GroupID, Revision: f.targetKey.GroupRevision}},
		"moved":                    []workspace.CutoverMemberRef{{DeviceID: m.DeviceID, PointID: m.PointID, TagID: m.TagID}},
	})
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('runtime-handoff',?,'forward',?,?,'runtime-review',?,'{}')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, boundary.Format(time.RFC3339Nano), string(detail))
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,'runtime-handoff')`, f.sourceKey.WorkspaceID, f.sourceKey.GroupID, f.sourceKey.GroupRevision, boundary.Format(time.RFC3339Nano))
	require.NoError(t, err)
	for direction, owner := range map[string]*longFixture{"ended": source, "started": target} {
		g := owner.config.Group
		_, err = tx.Exec(`INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('runtime-handoff',?,?,?,?,?,?,?,?,?,?,?,?,?)`, g.WorkspaceID, direction, g.ID, g.AppliedRevision, m.DeviceID, m.PointID, m.TagID, string(g.RecordingProfile.Kind), g.Destination.ConnectorID, g.Destination.ConnectorRevision, g.Destination.Database, g.Destination.TableSchema, g.Destination.TableName)
		require.NoError(t, err)
	}
	require.NoError(t, tx.Commit())
	f.old.SetEndsAt(boundary)
	source.config.EndsAt, source.config.Until = boundary, at(15*time.Second)
	return f
}

func TestGroupBoundaryCutoverReplayAuthorizesOriginalCounterProofWithoutRekey(t *testing.T) {
	f := newRuntimeHandoffFixture(t)
	original, err := f.store.Restore(t.Context(), f.sourceKey)
	require.NoError(t, err)
	require.Len(t, original.Samples, 2)
	batch, err := f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	restored, err := f.store.Restore(t.Context(), f.targetKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, snapshot.AcceptanceOrigin{WorkspaceID: f.sourceKey.WorkspaceID, GroupID: f.sourceKey.GroupID, GroupRevision: f.sourceKey.GroupRevision}, restored.Samples[0].AcceptanceOrigin)
	originalWithOrigin := original.Samples[0]
	originalWithOrigin.AcceptanceOrigin = restored.Samples[0].AcceptanceOrigin
	require.Equal(t, originalWithOrigin, restored.Samples[0], "only ephemeral origin differs; proof, original received_at and exact counter value are immutable")
	require.Equal(t, f.sourceKey.GroupID, restored.Samples[0].CounterProof.Member.GroupID)
	// Ordinary append remains forbidden, even after a transfer exists.
	require.Error(t, f.store.AppendSample(t.Context(), f.targetKey, boundaryStart, original.Samples[0]))
	require.Error(t, f.old.Tick(t.Context(), at(700*time.Second)))
	require.Error(t, f.emptyTarget.Tick(t.Context(), at(700*time.Second)))
	f.source.clock, f.target.clock = at(700*time.Second), at(700*time.Second)
	old, target := f.source.build(t), f.target.build(t)
	// The incoming request generates the target proof. Only verified receipt
	// recognition can identify it as the same original accepted observation.
	require.NoError(t, target.AcceptOwnedSample(t.Context(), f.target.envelope(0, "counter-after", at(13*time.Second), uint64(200))))
	require.NoError(t, old.Tick(t.Context(), at(700*time.Second)))
	require.NoError(t, target.Tick(t.Context(), at(700*time.Second)))
	var oldConsumed, targetConsumed int
	require.NoError(t, f.store.DB().QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')='counter-after'`, f.sourceKey.GroupID).Scan(&oldConsumed))
	require.NoError(t, f.store.DB().QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')='counter-after'`, f.targetKey.GroupID).Scan(&targetConsumed))
	require.Zero(t, oldConsumed)
	require.Equal(t, 1, targetConsumed)
	_, err = f.store.Reclaim(t.Context(), 0)
	require.NoError(t, err)
	require.NoError(t, target.Tick(t.Context(), at(1300*time.Second)))
	require.NoError(t, target.AcceptOwnedSample(t.Context(), f.target.envelope(0, "counter-after", at(13*time.Second), uint64(200))))
	var rows int
	require.NoError(t, f.store.DB().QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE json_extract(payload,'$.sample_id')='counter-after'`).Scan(&rows))
	require.Zero(t, rows, "finalized replay cannot reconstruct or reappend GC'd accepted bytes")
}

type pausedHandoffRecognition struct {
	*groupdelivery.Ledger
	reached, resume chan struct{}
	paused          atomic.Bool
}

func (l *pausedHandoffRecognition) RecognizeTransferredSample(ctx context.Context, incoming snapshot.Sample) (groupdelivery.HandoffRecognition, error) {
	result, err := l.Ledger.RecognizeTransferredSample(ctx, incoming)
	if err == nil && !result.Found && l.paused.CompareAndSwap(false, true) {
		close(l.reached)
		<-l.resume
	}
	return result, err
}

func TestGroupBoundaryCutoverReplayTransferBetweenRecognitionAndAppend(t *testing.T) {
	f := newRuntimeHandoffFixture(t)
	paused := &pausedHandoffRecognition{Ledger: f.store.Ledger(f.targetKey), reached: make(chan struct{}), resume: make(chan struct{})}
	f.emptyTarget.ledger = paused
	var release sync.Once
	defer release.Do(func() { close(paused.resume) })
	done := make(chan error, 1)
	input := f.target.envelope(0, "counter-after", at(13*time.Second), uint64(200))
	go func() { done <- f.emptyTarget.AcceptOwnedSample(t.Context(), input) }()
	select {
	case <-paused.reached:
	case err := <-done:
		t.Fatalf("intake missed no-receipt barrier: %v", err)
	}
	_, err := f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
	require.NoError(t, err)
	release.Do(func() { close(paused.resume) })
	require.ErrorIs(t, <-done, groupdelivery.ErrOwnershipConflict, "a no-hit is advisory; ordinary Append must not adopt history committed meanwhile")
	require.Zero(t, f.emptyTarget.assembler.OpenBuckets(), "failed append cannot mutate memory")
	restored, err := f.store.Restore(t.Context(), f.targetKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.NoError(t, f.emptyTarget.AcceptOwnedSample(t.Context(), input), "retry recognizes the durable receipt without a second append")
	require.Error(t, f.emptyTarget.Tick(t.Context(), at(700*time.Second)), "resend must not silently upgrade the stale assembler view")
	fresh := f.target.build(t)
	require.NoError(t, fresh.Tick(t.Context(), at(700*time.Second)))
	var copies, consumed int
	require.NoError(t, f.store.DB().QueryRow(`SELECT COUNT(*),SUM(consumed) FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')='counter-after'`, f.targetKey.GroupID).Scan(&copies, &consumed))
	require.Equal(t, 1, copies)
	require.Equal(t, 1, consumed)
}

func TestGroupBoundaryCutoverReplayLegacyOriginCannotBorrowAssociation(t *testing.T) {
	f := newRuntimeHandoffFixtureWithLegacy(t, true)
	_, err := f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
	require.NoError(t, err)
	restored, err := f.store.Restore(t.Context(), f.targetKey)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.True(t, restored.Samples[0].CounterProof.IsZero())
	require.NoError(t, validateCounterReplay(f.emptyTarget.members, restored, true))
	for _, origin := range []snapshot.AcceptanceOrigin{{}, {WorkspaceID: f.sourceKey.WorkspaceID, GroupID: "wrong-original-owner", GroupRevision: f.sourceKey.GroupRevision}} {
		changed := restored
		changed.Samples = append([]snapshot.Sample(nil), restored.Samples...)
		changed.Samples[0].AcceptanceOrigin = origin
		require.ErrorIs(t, validateCounterReplay(f.emptyTarget.members, changed, true), counterproof.ErrCounterBindingProofConflict, "global legacy association cannot authorize a changed or erased receipt origin")
		_, err := prepareCounterAssociations(t.Context(), f.target.config, changed)
		require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict, "cannot repair a forged transfer origin by associating current target meaning")
	}
	fresh := f.target.build(t)
	require.NoError(t, fresh.Tick(t.Context(), at(700*time.Second)))
	var targetAssociations int
	require.NoError(t, f.store.DB().QueryRow(`SELECT COUNT(*) FROM wg_counter_proof_associations WHERE group_id=?`, f.targetKey.GroupID).Scan(&targetAssociations))
	require.Zero(t, targetAssociations, "legacy transfer keeps its original association, never a target rewrite")
}

func TestGroupBoundaryCutoverReplayReopensDurableTransfer(t *testing.T) {
	f := newRuntimeHandoffFixture(t)
	_, err := f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
	require.NoError(t, err)
	require.NoError(t, f.store.DB().Close())
	reopened := groupdelivery.NewStore(openDurableDB(t, f.databasePath))
	f.target.config.Ledger = &frozenFixtureLedger{Ledger: reopened.Ledger(f.targetKey), store: reopened}
	fresh := f.target.build(t)
	require.NoError(t, fresh.AcceptOwnedSample(t.Context(), f.target.envelope(0, "counter-after", at(13*time.Second), uint64(200))))
	require.NoError(t, fresh.Tick(t.Context(), at(700*time.Second)))
	batch, err := reopened.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	var copies, consumed int
	require.NoError(t, reopened.DB().QueryRow(`SELECT COUNT(*),SUM(consumed) FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')='counter-after'`, f.targetKey.GroupID).Scan(&copies, &consumed))
	require.Equal(t, 1, copies)
	require.Equal(t, 1, consumed)
}
