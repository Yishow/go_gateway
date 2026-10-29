package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/tag"
)

// Keep the public setup/apply routes, then use explicit Reconcile/Tick/Dispatch
// rather than a long worker interval (which does not pause the initial cycle).
func finalizedAPIRig(t *testing.T, options cutoverOptions) (*cutoverRig, *groupdelivery.Store) {
	t.Helper()
	r := newCutoverRigWith(t, "sqlite", options)
	require.NoError(t, r.pipe.Stop(time.Second))
	r.anchor = r.wideEffective.Add(time.Minute).Truncate(time.Minute)
	r.f.groups.WithClock(func() time.Time { return r.wideEffective })
	r.clock.Store(r.wideEffective.Add(time.Second).UnixNano())
	store := groupdelivery.NewStore(r.f.db)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{Groups: r.f.groups, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Destinations: targets, Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: store, Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(r.f.db))}, grouppipeline.Config{Now: func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute, OnError: func(err error) { t.Errorf("manual pipeline: %v", err) }})
	r.reconcile()
	return r, store
}

type finalizedAPIResolver struct{ rig *cutoverRig }

func (r finalizedAPIResolver) Resolve(_ context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	require.Equal(r.rig.t, r.rig.wideTable(), item.Destination.TableName)
	return groupdelivery.Target{DB: r.rig.dest, Kind: schema.DatabaseConnectorKindSQLite}, nil
}

func TestCutoverFinalizedPublicAPI(t *testing.T) {
	for _, mode := range []string{"finalized future", "post reclaim", "authority unavailable"} {
		t.Run(mode, func(t *testing.T) {
			r, store := finalizedAPIRig(t, cutoverOptions{})
			boundary := r.at(22500 * time.Millisecond)
			review := r.review(boundary)
			require.Equal(t, true, review["applicable"])
			request := r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]})
			// A future observation at second 23 closes after review, before apply.
			r.accept(r.at(23*time.Second), 7.0)
			r.advance(r.at(31 * time.Second))
			r.pipe.TickAll(t.Context())
			require.Equal(t, 1, r.outbox()["pending"])
			sender := groupdelivery.NewSender(store, finalizedAPIResolver{r}, groupdelivery.SenderConfig{Owner: "finalized-test"})
			report, err := groupdelivery.NewDispatcher(store, sender, groupdelivery.DispatcherConfig{}).RunOnce(t.Context())
			require.NoError(t, err)
			require.Equal(t, 1, report.Delivered)
			rowsBefore := r.wideRowsInOrder()
			require.Len(t, rowsBefore, 1)
			require.Equal(t, 7.0, *rowsBefore[0].First)
			if mode == "post reclaim" {
				removed, err := store.Reclaim(t.Context(), 0)
				require.NoError(t, err)
				require.Positive(t, removed)
				require.Empty(t, r.outbox())
			}
			if mode == "authority unavailable" {
				_, err := r.f.db.Exec(`UPDATE wg_delivery_checkpoints SET next_close='dsn=private;value=secret'`)
				require.NoError(t, err)
			}
			counts := r.cutoverCounts()
			wide := getGroup(t, r.f, r.wide["id"].(string))
			digests := r.outboxDigests()
			receipts := r.receipts()
			code, body := r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover", request)
			expectedStatus, expectedCode := http.StatusConflict, "recording-cutover-finalized-overlap"
			if mode == "authority unavailable" {
				expectedStatus, expectedCode = http.StatusServiceUnavailable, "recording-ownership-unavailable"
			}
			require.Equal(t, expectedStatus, code, "%v", body)
			apiErr := errorOf(t, body)
			require.Equal(t, expectedCode, apiErr["code"])
			require.Equal(t, "review_cutover", apiErr["action"])
			require.Equal(t, false, apiErr["retryable"])
			require.NotContains(t, apiErr["message"], "private")
			require.NotContains(t, apiErr["message"], "secret")
			require.NotContains(t, apiErr["message"], "SQLITE")
			require.Equal(t, counts, r.cutoverCounts())
			require.Equal(t, wide, getGroup(t, r.f, r.wide["id"].(string)))
			require.Equal(t, digests, r.outboxDigests())
			require.Equal(t, receipts, r.receipts())
			require.Equal(t, rowsBefore, r.wideRowsInOrder())
		})
	}
}

func TestCutoverFinalizedRollbackIncludesRetainedSuccessor(t *testing.T) {
	r, _ := finalizedAPIRig(t, cutoverOptions{second: true})
	wideID := r.wide["id"].(string)
	r.saveRemainingDraft(getGroup(t, r.f, wideID), "tag-second")
	forward := r.cutover(r.at(22500 * time.Millisecond))
	require.Len(t, forward["started"], 2)
	r.reconcile()
	boundary := r.at(82500 * time.Millisecond)
	review := r.rollbackReview(forward["id"].(string), boundary)
	require.Equal(t, true, review["applicable"], "%v", review)
	request := r.rollbackBody(forward["id"].(string), boundary, map[string]any{"review_digest": review["digest"]})
	r.acceptFor(r.second, r.at(83*time.Second), 99.0)
	r.advance(r.at(91 * time.Second))
	r.pipe.TickAll(t.Context())
	// Only the retained 15-second successor closes past 82.5. The moved
	// 20-second target has a cursor of 80 and must not cause this refusal.
	var next string
	require.NoError(t, r.f.db.QueryRow(`SELECT next_close FROM wg_delivery_checkpoints WHERE group_id=?`, r.long["id"]).Scan(&next))
	require.Equal(t, r.at(80*time.Second), storedInstant(t, next))
	before := r.cutoverCounts()
	saved := getGroup(t, r.f, wideID)
	digests := r.outboxDigests()
	code, body := r.post(writeGroupsPath+"/"+wideID+"/cutover", request)
	require.Equal(t, http.StatusConflict, code, "%v", body)
	require.Equal(t, "recording-cutover-finalized-overlap", errorOf(t, body)["code"])
	require.Equal(t, before, r.cutoverCounts())
	require.Equal(t, saved, getGroup(t, r.f, wideID))
	require.Equal(t, digests, r.outboxDigests())
}

// A second SQLite connection commits closure after the configuration
// transaction has read its snapshot but before it takes write authority. WAL
// permits the closure to win; the stale snapshot must refuse its lock upgrade.
func TestCutoverFinalizedStaleTransactionSnapshotRefuses(t *testing.T) {
	r, _ := finalizedAPIRig(t, cutoverOptions{lostCommit: true})
	_, err := r.f.db.Exec(`PRAGMA journal_mode=WAL`)
	require.NoError(t, err)
	var seq int
	var name, path string
	require.NoError(t, r.f.db.QueryRow(`PRAGMA database_list`).Scan(&seq, &name, &path))
	writer, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer writer.Close()
	boundary := r.at(22500 * time.Millisecond)
	review := r.review(boundary)
	request := r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]})
	wide := getGroup(t, r.f, r.wide["id"].(string))
	before := r.cutoverCounts()
	ran := false
	hook := func(query string) error {
		if !strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at=updated_at WHERE 0") {
			return nil
		}
		ran = true
		key := groupdelivery.GroupKey{WorkspaceID: r.f.record.ID, GroupID: wide["id"].(string), GroupRevision: wide["applied_revision"].(string)}
		return groupdelivery.NewStore(writer).CommitClosure(t.Context(), key, groupdelivery.Closure{
			Destination: groupdelivery.FrozenDestination{Scope: "source", ConnectorID: "connector-A", ConnectorRevision: "connector-1", TableName: r.wideTable()}, NextClose: r.at(30 * time.Second),
			Buckets: []groupdelivery.ClosedBucket{{Outcome: snapshot.Outcome{Kind: snapshot.OutcomeNoData, BucketStart: r.at(15 * time.Second), BucketEnd: r.at(30 * time.Second)}}},
		})
	}
	r.faults.Reject.Store(&hook)
	code, body := r.post(writeGroupsPath+"/"+wide["id"].(string)+"/cutover", request)
	r.faults.Reject.Store(nil)
	require.True(t, ran, "closure must interleave inside the actual configuration transaction")
	require.Equal(t, http.StatusServiceUnavailable, code, "%v", body)
	require.Equal(t, "recording-ownership-unavailable", errorOf(t, body)["code"])
	require.Equal(t, before, r.cutoverCounts())
	require.Equal(t, wide, getGroup(t, r.f, wide["id"].(string)))
	var cursor string
	require.NoError(t, r.f.db.QueryRow(`SELECT next_close FROM wg_delivery_checkpoints`).Scan(&cursor))
	require.Equal(t, r.at(30*time.Second), storedInstant(t, cursor), "the winning closure survives the rejected cutover")
	// The same request now reads the committed authority and gives the precise
	// overlap code; no automatic retry or write was performed by the first call.
	code, body = r.post(writeGroupsPath+"/"+wide["id"].(string)+"/cutover", request)
	require.Equal(t, http.StatusConflict, code, "%v", body)
	require.Equal(t, "recording-cutover-finalized-overlap", errorOf(t, body)["code"])
	require.Equal(t, before, r.cutoverCounts())
}

// Snapshot every local authority/effect byte for rollback refusal assertions.
func finalizedAPIState(t *testing.T, r *cutoverRig) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for _, table := range []string{"system_settings", "write_groups", "write_group_versions", "wg_runtime_versions", "write_group_cutovers", "write_group_revision_ends", "write_group_ownership_fences", "wg_delivery_samples", "wg_delivery_buckets", "wg_delivery_outbox", "wg_delivery_checkpoints", "wg_delivery_counter_baselines"} {
		rows, err := r.f.db.Query(`SELECT * FROM ` + table + ` ORDER BY rowid`)
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		for rows.Next() {
			values := make([]any, len(columns))
			ptrs := make([]any, len(columns))
			for i := range values {
				ptrs[i] = &values[i]
			}
			require.NoError(t, rows.Scan(ptrs...))
			raw, err := json.Marshal(values)
			require.NoError(t, err)
			state[table] = append(state[table], string(raw))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return state
}
