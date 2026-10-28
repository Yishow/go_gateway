package grouppipeline

import (
	"database/sql"
	"encoding/json"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"
)

// Logical identity includes the original namespace, not the physical copy ID.
// Every copy must retain the exact original payload, digest and acceptance time.
func cutoverCrashJournal(t *testing.T, db *sql.DB) (map[string]string, map[string][]string) {
	t.Helper()
	rows, err := db.Query(`SELECT json_array(
 CASE WHEN origin_workspace_id='' THEN workspace_id ELSE origin_workspace_id END,
 CASE WHEN origin_group_id='' THEN group_id ELSE origin_group_id END,
 CASE WHEN origin_revision='' THEN group_revision ELSE origin_revision END,
 CASE WHEN origin_sample_id='' THEN sample_id ELSE origin_sample_id END),
 json_array(payload,payload_digest,created_at),json_array(workspace_id,group_id,group_revision),consumed
 FROM wg_delivery_samples ORDER BY seq`)
	require.NoError(t, err)
	defer rows.Close()
	originals, consumers := map[string]string{}, map[string][]string{}
	for rows.Next() {
		var identity, original, owner string
		var consumed int
		require.NoError(t, rows.Scan(&identity, &original, &owner, &consumed))
		if previous, ok := originals[identity]; ok {
			require.Equal(t, previous, original, "a physical handoff copy cannot change its original measurement")
		}
		originals[identity] = original
		if consumed != 0 {
			consumers[identity] = append(consumers[identity], owner)
		}
	}
	require.NoError(t, rows.Err())
	return originals, consumers
}

func cutoverCrashEffects(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query(`SELECT json_array(effect_key,record_id,payload,payload_digest,
 destination_scope,connector_id,connector_revision,database_name,table_schema,table_name,dedupe_capability,record_key_column)
 FROM wg_delivery_outbox ORDER BY effect_key`)
	require.NoError(t, err)
	defer rows.Close()
	var effects []string
	for rows.Next() {
		var effect string
		require.NoError(t, rows.Scan(&effect))
		effects = append(effects, effect)
	}
	require.NoError(t, rows.Err())
	return effects
}

func cutoverCrashDescriptors(t *testing.T, f *ownershipPipelineFixture) map[string]string {
	t.Helper()
	result := map[string]string{}
	for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
		canonical, err := json.Marshal(owner)
		require.NoError(t, err)
		result[owner.GroupID+"/canonical"] = string(canonical)
		var applied string
		require.NoError(t, f.local.QueryRow(`SELECT payload FROM write_group_versions WHERE group_id=? AND group_revision=?`, owner.GroupID, owner.AppliedRevision).Scan(&applied))
		result[owner.GroupID+"/applied"] = applied
		frozen, err := f.store.RuntimeVersion(t.Context(), groupKey(owner))
		require.NoError(t, err)
		result[owner.GroupID+"/frozen"] = string(frozen)
	}
	return result
}

func cutoverCrashBaselines(t *testing.T, db *sql.DB) map[string]string {
	t.Helper()
	rows, err := db.Query(`SELECT json_array(workspace_id,group_id,group_revision,member_key),json_array(window_end,payload) FROM wg_delivery_counter_baselines ORDER BY group_id,member_key`)
	require.NoError(t, err)
	defer rows.Close()
	result := map[string]string{}
	for rows.Next() {
		var owner, baseline string
		require.NoError(t, rows.Scan(&owner, &baseline))
		result[owner] = baseline
	}
	require.NoError(t, rows.Err())
	return result
}

func reopenCutoverPipeline(t *testing.T, f *ownershipPipelineFixture) {
	t.Helper()
	deps, config := f.pipe.deps, f.pipe.config
	config.Sender.Recovery = nil // Derived runtime service also holds the old DB; rebuild it.
	oldDB := f.local
	f.pipe, f.store = nil, nil // No worker or assembler from the crashed process is reused.
	require.NoError(t, oldDB.Close())
	require.Error(t, oldDB.Ping(), "the previous physical database handle is closed")
	plain, err := sql.Open("sqlite", f.localPath+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	require.NoError(t, datalink.NewMigrator().Migrate(plain))
	require.NoError(t, plain.Close())
	f.local = lostcommit.OpenFaults(f.localPath+"?_pragma=busy_timeout(5000)", f.localFault)
	f.store = groupdelivery.NewStore(f.local).WithClock(func() time.Time { return f.clock })
	deps.Store = f.store
	f.pipe = New(deps, config)
	require.Empty(t, f.pipe.boundaries)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
}

// The runtime matrix owns all four fault positions. This live integration runs
// the two uncertain-commit boundaries in sequence, with real physical reopen
// after each, then drains immutable output to its original SQLite destinations.
func TestCutoverReconcileCrashRestoresCommittedHandoffAndFrozenOutput(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		for _, reduction := range []string{"latest", "cumulative"} {
			t.Run(string(shape)+"/"+reduction, func(t *testing.T) {
				f, reader := reductionHandoffFixture(t, shape, reduction)
				if reduction == "cumulative" {
					f.clock = at(-2)
					require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "previous", at(-1), uint64(5))))
					require.NoError(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].boundary.Tick(t.Context(), at(0)))
					f.clock = at(7)
				}
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "old-only", at(9), uint64(20))))
				require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "successor", at(11), uint64(100))))
				originals, _ := cutoverCrashJournal(t, f.local)
				beforeEffects := cutoverCrashEffects(t, f.local)
				f.commit()
				f.source.EndsAt = f.boundary
				descriptors := cutoverCrashDescriptors(t, f)
				f.groups.groups = append(f.groups.groups, f.target.Group)
				f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
				f.clock = at(12)
				var armed atomic.Bool
				reject := func(query string) error {
					if strings.Contains(query, "INSERT INTO wg_delivery_sample_handoffs") {
						armed.Store(true)
						f.localFault.Lose.Store(true) // Only after real transfer work, never an earlier advisory read.
					}
					return nil
				}
				f.localFault.Reject.Store(&reject)
				err := f.pipe.Reconcile(t.Context())
				f.localFault.Reject.Store(nil)
				f.localFault.Lose.Store(false)
				require.True(t, armed.Load())
				require.ErrorIs(t, err, lostcommit.ErrResponseLost)
				require.Equal(t, 1, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE state='accepted'`), "handoff really committed despite the lost response")
				require.Equal(t, beforeEffects, cutoverCrashEffects(t, f.local), "handoff is not closure or external commit")
				after, _ := cutoverCrashJournal(t, f.local)
				require.Equal(t, originals, after)
				var receipt string
				require.NoError(t, f.local.QueryRow(`SELECT json_array(transfer_id,record,target_seq) FROM wg_delivery_sample_handoffs`).Scan(&receipt))
				reopenCutoverPipeline(t, f)
				for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
					state, err := f.store.Restore(t.Context(), groupKey(owner))
					require.NoError(t, err)
					view := f.pipe.boundaries[boundaryKey(owner.GroupID, owner.AppliedRevision)].boundary.ReplayState()
					require.Equal(t, state.OwnershipView, view.OwnershipView)
					require.Equal(t, state.HandoffView, view.HandoffView)
				}
				f.clock = at(15)
				if reduction == "cumulative" {
					f.clock = at(70)
				}
				f.localFault.Lose.Store(true)
				for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
					err := f.pipe.boundaries[boundaryKey(owner.GroupID, owner.AppliedRevision)].boundary.Tick(t.Context(), f.clock)
					require.ErrorIs(t, err, lostcommit.ErrResponseLost)
				}
				f.localFault.Lose.Store(false)
				finalEffects := cutoverCrashEffects(t, f.local)
				finalBaselines := cutoverCrashBaselines(t, f.local)
				require.Greater(t, len(finalEffects), len(beforeEffects), "closure really committed despite the lost replies")
				after, consumers := cutoverCrashJournal(t, f.local)
				require.Equal(t, originals, after)
				require.Len(t, consumers, len(originals), "the consuming union equals the original accepted set")
				for identity := range originals {
					var origin []string
					require.NoError(t, json.Unmarshal([]byte(identity), &origin))
					owner := f.source
					if origin[3] == "successor" {
						owner = f.target
					}
					want, err := json.Marshal([]string{owner.WorkspaceID, owner.GroupID, owner.AppliedRevision})
					require.NoError(t, err)
					require.Equal(t, []string{string(want)}, consumers[identity], "exactly the eligible owner consumes each original identity")
				}
				_, err = f.plainTarget().Exec(`CREATE TABLE unapplied_destination(marker TEXT)`)
				require.NoError(t, err)
				// Current drafts can differ; original applied descriptors/history do not.
				f.groups.groups = nil
				for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
					raw, err := json.Marshal(owner.Group)
					require.NoError(t, err)
					var current workspace.WriteGroup
					require.NoError(t, json.Unmarshal(raw, &current))
					current.Revision += "-draft"
					current.Destination.TableName = "unapplied_destination"
					f.groups.groups = append(f.groups.groups, &current)
				}
				require.Equal(t, descriptors, cutoverCrashDescriptors(t, f), "editable drafts share no mutable state with original descriptors")
				reopenCutoverPipeline(t, f)
				for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
					require.NoError(t, f.pipe.boundaries[boundaryKey(owner.GroupID, owner.AppliedRevision)].boundary.Tick(t.Context(), f.clock), "recovery resumes successfully, rather than hiding repeated closure refusals")
				}
				require.Equal(t, finalEffects, cutoverCrashEffects(t, f.local), "recovery never recomputes finalized bytes or retargets frozen destinations")
				require.Equal(t, finalBaselines, cutoverCrashBaselines(t, f.local), "committed baselines survive physical reopen byte-for-byte")
				f.deliver()
				oldValue, oldEvidence := readReductionRow(t, f, f.source, at(0), reduction)
				newValue, newEvidence := readReductionRow(t, f, f.target, at(10), reduction)
				require.Equal(t, "20", oldValue)
				require.Equal(t, "100", newValue)
				require.Equal(t, "old-only", oldEvidence.SampleID)
				require.Equal(t, "successor", newEvidence.SampleID)
				require.Equal(t, "end", oldEvidence.Window.Truncated)
				require.Equal(t, "start", newEvidence.Window.Truncated)
				require.Equal(t, oldEvidence.Window.End, newEvidence.Window.Start)
				require.Equal(t, at(10).Format("2006-01-02T15:04:05.000000Z"), oldEvidence.Window.End)
				if reduction == "cumulative" {
					require.Len(t, finalBaselines, 2)
					require.Equal(t, sql.NullString{String: "15", Valid: true}, oldEvidence.delta)
					require.Equal(t, "baseline", newEvidence.status)
					require.False(t, newEvidence.delta.Valid || newEvidence.deltaFrom.Valid || newEvidence.deltaTo.Valid)
				} else {
					require.Empty(t, finalBaselines)
				}
				require.Zero(t, f.count(f.plainTarget(), `SELECT COUNT(*) FROM unapplied_destination`))
				var retriedReceipt string
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				require.NoError(t, f.local.QueryRow(`SELECT json_array(transfer_id,record,target_seq) FROM wg_delivery_sample_handoffs`).Scan(&retriedReceipt))
				require.Equal(t, receipt, retriedReceipt)
				require.Equal(t, finalEffects, cutoverCrashEffects(t, f.local))
				require.Equal(t, descriptors, cutoverCrashDescriptors(t, f))
				require.Zero(t, reader.calls)
			})
		}
	}
}
