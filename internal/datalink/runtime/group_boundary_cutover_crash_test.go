package runtime

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

type crashIdentity struct {
	origin groupdelivery.GroupKey
	sample string
}
type crashAccepted struct{ payload, digest, acceptedAt string }

// Recovery does not retain fixture Counters, cached boundary state or current
// definitions. Only the original persisted descriptor supplies measurement meaning.
func crashFrozenBoundary(t *testing.T, store *groupdelivery.Store, key groupdelivery.GroupKey) *GroupBoundary {
	t.Helper()
	raw, err := store.RuntimeVersion(t.Context(), key)
	require.NoError(t, err)
	var frozen workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &frozen))
	require.NotNil(t, frozen.RuntimeLayout)
	interval := time.Duration(frozen.Group.RowPolicy.IntervalSeconds) * time.Second
	cfg := GroupBoundaryConfig{Group: frozen.Group, TagTypes: frozen.RuntimeLayout.TagTypes, CounterBindings: frozen.CounterBindings,
		Dialect: frozen.RuntimeLayout.Dialect, Columns: frozen.RuntimeLayout.Columns, ReceiptTableReady: frozen.RuntimeLayout.ReceiptTableReady,
		FirstBucket: snapshot.BucketStart(frozen.EffectiveAt, interval), StartsAt: frozen.EffectiveAt, Ledger: store.Ledger(key), MaxFutureSkew: 2 * time.Second,
		Clock: func() time.Time { return at(700 * time.Second) }}
	var end string
	err = store.DB().QueryRow(`SELECT ends_at FROM write_group_revision_ends WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&end)
	if !errors.Is(err, sql.ErrNoRows) {
		require.NoError(t, err)
		cfg.EndsAt, err = time.Parse(time.RFC3339Nano, end)
		require.NoError(t, err)
		cfg.Until = snapshot.BucketStart(cfg.EndsAt, interval)
		if cfg.Until.Before(cfg.EndsAt) {
			cfg.Until = cfg.Until.Add(interval)
		}
	}
	b, err := NewGroupBoundaryContext(t.Context(), cfg)
	require.NoError(t, err)
	return b
}

func crashRows(t *testing.T, db *sql.DB, table string) []string {
	t.Helper()
	rows, err := db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	var result []string
	for rows.Next() {
		values, refs := make([]any, len(columns)), make([]any, len(columns))
		for i := range values {
			refs[i] = &values[i]
		}
		require.NoError(t, rows.Scan(refs...))
		raw, err := json.Marshal(values)
		require.NoError(t, err)
		result = append(result, string(raw))
	}
	require.NoError(t, rows.Err())
	return result
}

func crashProgress(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	result := map[string][]string{}
	for _, table := range []string{"wg_delivery_samples", "wg_delivery_sample_handoffs", "wg_delivery_checkpoints", "wg_delivery_counter_baselines", "wg_delivery_buckets", "wg_delivery_outbox", "wg_runtime_versions"} {
		result[table] = crashRows(t, db, table)
	}
	return result
}

// Identity always includes the original workspace/group/revision namespace;
// copied physical row addresses are not a second logical accepted observation.
func crashAssertJournal(t *testing.T, f *runtimeHandoffFixture, originals map[crashIdentity]crashAccepted, finalized bool) {
	t.Helper()
	rows, err := f.store.DB().Query(`SELECT workspace_id,group_id,group_revision,COALESCE(NULLIF(origin_workspace_id,''),workspace_id),COALESCE(NULLIF(origin_group_id,''),group_id),COALESCE(NULLIF(origin_revision,''),group_revision),COALESCE(NULLIF(origin_sample_id,''),sample_id),payload,payload_digest,created_at,consumed FROM wg_delivery_samples ORDER BY seq`)
	require.NoError(t, err)
	defer rows.Close()
	union, consumers := map[crashIdentity]bool{}, map[crashIdentity]groupdelivery.GroupKey{}
	copies := map[groupdelivery.GroupKey]int{}
	for rows.Next() {
		var owner groupdelivery.GroupKey
		var identity crashIdentity
		var accepted crashAccepted
		var consumed int
		require.NoError(t, rows.Scan(&owner.WorkspaceID, &owner.GroupID, &owner.GroupRevision, &identity.origin.WorkspaceID, &identity.origin.GroupID, &identity.origin.GroupRevision, &identity.sample, &accepted.payload, &accepted.digest, &accepted.acceptedAt, &consumed))
		require.True(t, owner == f.sourceKey || owner == f.targetKey, "no unrelated scope may acquire an original")
		original, ok := originals[identity]
		require.True(t, ok, "unknown logical origin: %+v", identity)
		require.Equal(t, original, accepted, "payload, proof, digest and original accepted_at must remain byte-identical")
		union[identity] = true
		copies[owner]++
		if consumed != 0 {
			require.NotContains(t, consumers, identity, "source and target consuming sets must be disjoint")
			consumers[identity] = owner
			want := f.targetKey
			if identity.sample == "before" {
				want = f.sourceKey
			}
			require.Equal(t, want, owner, "only the half-open owner consumes the original observation")
		}
	}
	require.NoError(t, rows.Err())
	require.Len(t, union, len(originals), "accepted logical union never shrinks across crash/restart")
	require.Equal(t, 3, copies[f.sourceKey], "handoff retains every original journal row")
	if finalized {
		require.Len(t, consumers, 3)
		require.Equal(t, 2, copies[f.targetKey])
	}
}

func crashReopen(t *testing.T, f *runtimeHandoffFixture) {
	t.Helper()
	require.NoError(t, f.store.DB().Close())
	f.old, f.emptyTarget = nil, nil
	f.source.config.Ledger, f.target.config.Ledger = nil, nil
	f.source.config.Counters, f.target.config.Counters = nil, nil
	f.store = groupdelivery.NewStore(openDurableDB(t, f.databasePath))
}

func crashAssertFinal(t *testing.T, f *runtimeHandoffFixture) {
	t.Helper()
	end := at(time.Duration(f.target.config.Group.RowPolicy.IntervalSeconds) * time.Second)
	for _, tc := range []struct {
		key             groupdelivery.GroupKey
		last            string
		windowEnd, next time.Time
	}{
		{f.sourceKey, "before", f.boundary, at(15 * time.Second)}, {f.targetKey, "after", end, end},
	} {
		restored, err := f.store.Restore(t.Context(), tc.key)
		require.NoError(t, err)
		require.Empty(t, restored.Samples)
		require.Equal(t, tc.next, restored.NextClose)
		require.Len(t, restored.Baselines, 1)
		checkpoint, err := aggregation.DecodeCounterCheckpoint(restored.Baselines[0].Payload)
		require.NoError(t, err)
		require.Equal(t, tc.last, checkpoint.Reading.SampleID)
		require.Equal(t, tc.windowEnd, checkpoint.WindowEnd)
	}
	require.Len(t, crashRows(t, f.store.DB(), "wg_delivery_sample_handoffs"), 2)
	require.Len(t, crashRows(t, f.store.DB(), "wg_delivery_outbox"), 2)
	reductionReadback(t, f.source, reductionOutbox(t, f.store.DB(), f.sourceKey.GroupID), "100", "end", f.boundary)
	reductionReadback(t, f.target, reductionOutbox(t, f.store.DB(), f.targetKey.GroupID), "200", "start", f.boundary)
}

func TestGroupBoundaryCutoverCrashRecoveryMatrix(t *testing.T) {
	for _, wide := range []bool{false, true} {
		layout := "long"
		if wide {
			layout = "wide"
		}
		for _, phase := range []string{"rollback-before-handoff", "lost-response-after-handoff", "rollback-before-closure", "lost-response-after-closure"} {
			t.Run(layout+"/"+phase, func(t *testing.T) {
				f := newReductionHandoffFixture(t, wide)
				originals := map[crashIdentity]crashAccepted{}
				for _, id := range []string{"before", "exact", "after"} {
					raw := reductionOriginal(t, f.store.DB(), f.sourceKey.GroupID, id)
					originals[crashIdentity{f.sourceKey, id}] = crashAccepted{raw[0], raw[1], raw[2]}
				}
				// Replace the connection with the existing fault-injection driver, before
				// arming it. No old worker or current counter reader participates in recovery.
				require.NoError(t, f.store.DB().Close())
				f.old, f.emptyTarget = nil, nil
				faults := &lostcommit.Faults{Lose: &atomic.Bool{}}
				faultDB := lostcommit.OpenFaults(f.databasePath, faults)
				t.Cleanup(func() { _ = faultDB.Close() })
				require.NoError(t, datalink.NewMigrator().Migrate(faultDB))
				f.store = groupdelivery.NewStore(faultDB)
				targetEnd := at(time.Duration(f.target.config.Group.RowPolicy.IntervalSeconds) * time.Second)
				require.ErrorIs(t, crashFrozenBoundary(t, f.store, f.sourceKey).Tick(t.Context(), at(15*time.Second)), groupdelivery.ErrOwnershipConflict)
				require.ErrorIs(t, crashFrozenBoundary(t, f.store, f.targetKey).Tick(t.Context(), targetEnd), groupdelivery.ErrOwnershipConflict, "pending original acceptance must prevent premature empty target closure")
				closure := strings.Contains(phase, "closure")
				lost := strings.HasPrefix(phase, "lost-response")
				if closure {
					_, err := f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
					require.NoError(t, err)
				}
				var boundary *GroupBoundary
				if closure {
					boundary = crashFrozenBoundary(t, f.store, f.sourceKey)
				}
				before := crashProgress(t, f.store.DB())
				var committed, hits atomic.Int32
				afterCommit := func() { committed.Add(1) }
				faults.AfterCommit.Store(&afterCommit)
				injected := errors.New("injected crash before durable commit")
				reject := func(query string) error {
					statement := "INSERT INTO wg_delivery_sample_handoffs"
					if closure {
						statement = "INSERT INTO wg_delivery_checkpoints"
					}
					if strings.Contains(query, statement) {
						hits.Add(1)
						return injected
					}
					return nil
				}
				if lost {
					faults.Lose.Store(true)
				} else {
					faults.Reject.Store(&reject)
				}
				var err error
				if closure {
					err = boundary.Tick(t.Context(), at(15*time.Second))
				} else {
					_, err = f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
				}
				if lost {
					require.ErrorIs(t, err, lostcommit.ErrResponseLost)
					require.EqualValues(t, 1, committed.Load())
				} else {
					require.ErrorIs(t, err, injected)
					require.EqualValues(t, 1, hits.Load())
					require.Zero(t, committed.Load())
					require.Equal(t, before, crashProgress(t, f.store.DB()), "late transaction failure rolls back every partial write")
				}
				if closure {
					require.Equal(t, boundaryStart, boundary.NextClose(), "a lost response must not acknowledge the old worker's memory")
				}
				boundary = nil // discard even when its transaction actually committed
				crashAssertJournal(t, f, originals, false)
				atCrash := crashProgress(t, f.store.DB())
				if lost && !closure {
					require.Len(t, atCrash["wg_delivery_sample_handoffs"], 1)
					require.Len(t, atCrash["wg_delivery_samples"], 4)
				}
				if lost && closure {
					require.Len(t, atCrash["wg_delivery_outbox"], 1)
					require.Len(t, atCrash["wg_delivery_counter_baselines"], 1)
				}
				crashReopen(t, f)
				require.Equal(t, atCrash, crashProgress(t, f.store.DB()), "physical restart preserves exact committed versus rolled-back state")
				crashAssertJournal(t, f, originals, false)
				_, err = f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
				require.NoError(t, err)
				require.NoError(t, crashFrozenBoundary(t, f.store, f.sourceKey).Tick(t.Context(), at(15*time.Second)))
				require.NoError(t, crashFrozenBoundary(t, f.store, f.targetKey).Tick(t.Context(), targetEnd))
				crashAssertJournal(t, f, originals, true)
				crashAssertFinal(t, f)
				finalized := crashProgress(t, f.store.DB())
				require.Equal(t, atCrash["wg_runtime_versions"], finalized["wg_runtime_versions"], "recovery never rewrites original frozen measurement meaning")
				for _, row := range atCrash["wg_delivery_outbox"] {
					require.Contains(t, finalized["wg_delivery_outbox"], row, "previously committed effects retain every byte and frozen destination")
				}
				for _, row := range atCrash["wg_delivery_sample_handoffs"] {
					require.Contains(t, finalized["wg_delivery_sample_handoffs"], row, "lost response retries retain transfer identity and original proof")
				}
				crashReopen(t, f)
				_, err = f.store.PrepareCutoverHandoffs(t.Context(), f.sourceKey, 0, 128)
				require.NoError(t, err)
				old, target := crashFrozenBoundary(t, f.store, f.sourceKey), crashFrozenBoundary(t, f.store, f.targetKey)
				for _, input := range []struct {
					id    string
					at    time.Time
					value uint64
				}{{"exact", f.boundary, 150}, {"after", at(13 * time.Second), 200}} {
					require.NoError(t, target.AcceptOwnedSample(t.Context(), f.target.envelope(0, input.id, input.at, input.value)))
				}
				require.NoError(t, old.Tick(t.Context(), at(30*time.Second)))
				require.NoError(t, target.Tick(t.Context(), targetEnd))
				require.Equal(t, finalized, crashProgress(t, f.store.DB()), "second reopen, resends and retries add no accepted copies or effects")
				crashAssertJournal(t, f, originals, true)
				crashAssertFinal(t, f)
			})
		}
	}
}
