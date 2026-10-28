package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/testutil/lostcommit"
)

// Both connections address the actual workspace database. Only the delivery
// connection pauses before write authority; HTTP apply remains able to commit.
// No worker is running, so Reconcile, closure and delivery are explicit stages.
func barrierAPIRig(t *testing.T, wideTarget bool) (*cutoverRig, *groupdelivery.Store, *lostcommit.Faults, chan error) {
	t.Helper()
	r := newCutoverRigWith(t, "sqlite", cutoverOptions{wideTarget: wideTarget})
	require.NoError(t, r.pipe.Stop(time.Second))
	r.anchor = r.wideEffective.Add(time.Minute).Truncate(time.Minute)
	var seq int
	var name, path string
	require.NoError(t, r.f.db.QueryRow(`PRAGMA database_list`).Scan(&seq, &name, &path))
	faults := &lostcommit.Faults{}
	deliveryDB := lostcommit.OpenFaults(path, faults)
	deliveryDB.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, deliveryDB.Close()) })
	now := func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }
	r.f.groups.WithClock(now)
	store := groupdelivery.NewStore(deliveryDB).WithClock(now)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(r.f.db))
	reports := make(chan error, 32)
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{
		Groups: r.f.groups, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Destinations: targets,
		Inspector: dbtarget.NewReadOnlyTableInspector(targets), Store: store, Ledger: ledger,
	}, grouppipeline.Config{Now: now, MaxFutureSkew: time.Minute, OnError: func(err error) { reports <- err }})
	pipe := r.pipe
	t.Cleanup(func() { require.NoError(t, pipe.Stop(time.Second)) })
	r.f.router = NewRouter(&DatalinkServices{
		Workspace: r.f.workspace, WriteGroups: r.f.groups, DBTarget: targets,
		WriteGroupDelivery: pipe, RecordingPlan: ledger,
	})
	barrierAPIReconcile(t, r)
	require.Empty(t, reports)
	return r, store, faults, reports
}

func barrierAPIReconcile(t *testing.T, r *cutoverRig) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	require.NoError(t, r.pipe.Reconcile(ctx), "explicit reconciliation must finish with one delivery connection")
}

func barrierAPIRefusal(t *testing.T, reports <-chan error) {
	t.Helper()
	select {
	case err := <-reports:
		require.ErrorIs(t, err, groupdelivery.ErrOwnershipConflict)
	default:
		t.Fatal("completed stale TickAll did not report an ownership refusal")
	}
	require.Empty(t, reports, "one source worker must report exactly one refusal")
}

type barrierAPIResolver struct{ rig *cutoverRig }

func (r barrierAPIResolver) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	require.Contains(r.rig.t, []string{r.rig.wideTable(), r.rig.longTable()}, item.Destination.TableName)
	require.Equal(r.rig.t, "connector-A", item.Destination.ConnectorID)
	require.Equal(r.rig.t, "connector-1", item.Destination.ConnectorRevision)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.rig.f.db))
	opened, err := targets.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close}, nil
}

func barrierAPISample(r *cutoverRig, id string, observed time.Time, value float64) measurement.SampleEnvelope {
	return measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: r.f.record.ID, DeviceID: r.member.DeviceID, PointID: r.member.PointID,
		TagID: r.member.TagID, SourceRevision: r.member.SourceRevision, MappingRevision: r.member.MappingRevision,
		ObservedAt: observed, ReceivedAt: time.Unix(0, r.clock.Load()).UTC(), TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}
}

// This includes all payload/proof bytes and acceptance timestamps, but not the
// mutable consumed flag. Copies must preserve it under the original namespace.
func barrierAPIOriginals(t *testing.T, r *cutoverRig) map[string]string {
	t.Helper()
	rows, err := r.f.db.Query(`SELECT sample_id,payload,payload_digest,created_at FROM wg_delivery_samples WHERE group_id=? ORDER BY seq`, r.wide["id"])
	require.NoError(t, err)
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var id, payload, digest, created string
		require.NoError(t, rows.Scan(&id, &payload, &digest, &created))
		out[id] = payload + "|" + digest + "|" + created
	}
	require.NoError(t, rows.Err())
	return out
}

func barrierAPIDeliveryState(t *testing.T, r *cutoverRig) map[string][]string {
	t.Helper()
	all := finalizedAPIState(t, r)
	state := map[string][]string{}
	for _, table := range []string{"wg_delivery_samples", "wg_delivery_buckets", "wg_delivery_checkpoints", "wg_delivery_counter_baselines", "wg_delivery_outbox"} {
		state[table] = all[table]
	}
	return state
}

func TestCutoverPublicBarrierDelayedReconcileSQLite(t *testing.T) {
	for _, wide := range []bool{false, true} {
		layout := "long"
		if wide {
			layout = "wide"
		}
		for _, example := range []struct {
			name               string
			boundary, after    time.Duration
			oldValue, newValue float64
		}{
			{"22.5-old6-new7", 22500 * time.Millisecond, 23 * time.Second, 6, 7},
			{"12.5-old100-new200", 12500 * time.Millisecond, 13 * time.Second, 100, 200},
		} {
			t.Run(layout+"/"+example.name, func(t *testing.T) {
				r, store, faults, reports := barrierAPIRig(t, wide)
				boundary := r.at(example.boundary)
				r.advance(boundary.Add(-1500 * time.Millisecond))
				inputs := []measurement.SampleEnvelope{
					barrierAPISample(r, "after", r.at(example.after), example.newValue),
					barrierAPISample(r, "equal", boundary, example.newValue),
					barrierAPISample(r, "before", boundary.Add(-time.Millisecond), example.oldValue),
				}
				for _, input := range inputs { // Reverse observation order, accepted before the future cutover.
					require.NoError(t, r.pipe.AcceptSample(t.Context(), input))
				}
				originals := barrierAPIOriginals(t, r)
				require.Len(t, originals, 3)
				sourceBefore := getGroup(t, r.f, r.wide["id"].(string))
				review := r.review(boundary)
				require.Equal(t, true, review["applicable"])
				request := r.cutoverBody(boundary, map[string]any{"review_digest": review["digest"]})
				entered, release, done := make(chan struct{}), make(chan struct{}), make(chan struct{})
				var once atomic.Bool
				var released sync.Once
				t.Cleanup(func() {
					released.Do(func() { close(release) })
					select {
					case <-done:
					case <-time.After(5 * time.Second):
						t.Error("closure goroutine did not stop before delivery connection cleanup")
					}
					faults.Reject.Store(nil)
				})
				hook := func(query string) error {
					if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0") && once.CompareAndSwap(false, true) {
						close(entered)
						<-release
					}
					return nil
				}
				faults.Reject.Store(&hook)
				// Only earlier due windows are planned here. All clocks remain before
				// boundary; no time reversal or retroactive public cutover is needed.
				go func() { r.pipe.TickAll(t.Context()); close(done) }()
				select {
				case <-entered:
				case <-time.After(10 * time.Second):
					t.Fatal("closure did not reach the real pre-write barrier")
				}
				code, body := r.post(writeGroupsPath+"/"+r.wide["id"].(string)+"/cutover", request)
				require.Equal(t, http.StatusOK, code, "%v", body)
				cutover := body["data"].(map[string]any)
				require.Equal(t, 2, r.localCount(`SELECT COUNT(*) FROM write_group_ownership_fences WHERE fence_revision=?`, cutover["id"]))
				var persistedEnd string
				require.NoError(t, r.f.db.QueryRow(`SELECT ends_at FROM write_group_revision_ends WHERE group_id=?`, r.wide["id"]).Scan(&persistedEnd))
				require.Equal(t, boundary, storedInstant(t, persistedEnd))
				require.Equal(t, sourceBefore, getGroup(t, r.f, r.wide["id"].(string)), "forward cutover preserves source status and historical descriptor")
				beforeRefusal := barrierAPIDeliveryState(t, r)
				released.Do(func() { close(release) })
				select {
				case <-done:
				case <-time.After(5 * time.Second):
					t.Fatal("released closure did not finish after public cutover committed")
				}
				faults.Reject.Store(nil)
				barrierAPIRefusal(t, reports)
				require.Equal(t, beforeRefusal, barrierAPIDeliveryState(t, r))
				require.Equal(t, originals, barrierAPIOriginals(t, r))
				require.Zero(t, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs`), "HTTP apply cannot secretly reconcile the journal")
				require.Zero(t, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=?`, r.long["id"]))
				unready := barrierAPISample(r, "unready-owner", r.at(example.after).Add(time.Millisecond), 999)
				require.ErrorIs(t, r.pipe.AcceptSample(t.Context(), unready), groupdelivery.ErrOwnershipUnavailable)
				require.Equal(t, beforeRefusal, barrierAPIDeliveryState(t, r), "no ready successor means no unpersisted ACK or old-owner append")
				r.advance(r.at(41 * time.Second))
				r.pipe.TickAll(t.Context())
				barrierAPIRefusal(t, reports)
				require.Equal(t, beforeRefusal, barrierAPIDeliveryState(t, r), "even the now-due cross-boundary stale plan cannot commit")
				require.Zero(t, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs`))
				barrierAPIReconcile(t, r)
				require.Equal(t, originals, barrierAPIOriginals(t, r))
				barrierAPIAssertHandoff(t, r, originals, false)
				r.pipe.TickAll(t.Context())
				require.Empty(t, reports)
				barrierAPIAssertHandoff(t, r, originals, true)
				targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
				recovery := grouppipeline.NewStorageRecovery(store, recordingplan.NewService(recordingplan.NewSQLRepository(r.f.db)), targets)
				now := func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }
				sender := groupdelivery.NewSender(store, barrierAPIResolver{r}, groupdelivery.SenderConfig{Owner: "public-barrier", Recovery: recovery, Now: now})
				// Delivery is a separate, sequential stage; the controlled race is
				// public cutover versus the old closure, not concurrent senders.
				dispatch := groupdelivery.NewDispatcher(store, sender, groupdelivery.DispatcherConfig{Now: now, MaxPartitions: 1})
				report, err := dispatch.RunOnce(t.Context())
				require.NoError(t, err)
				require.GreaterOrEqual(t, report.Delivered, 2)
				barrierAPIAssertRows(t, r, wide, boundary, inputs[2].ObservedAt, inputs[0].ObservedAt, example.oldValue, example.newValue)
				digests, receipts := r.outboxDigests(), r.receipts()
				barrierAPIReconcile(t, r)
				r.pipe.TickAll(t.Context())
				report, err = dispatch.RunOnce(t.Context())
				require.NoError(t, err)
				require.Zero(t, report.Delivered)
				require.Empty(t, reports)
				require.Equal(t, digests, r.outboxDigests())
				require.Equal(t, receipts, r.receipts())
				barrierAPIAssertHandoff(t, r, originals, true)
			})
		}
	}
}

func barrierAPIAssertHandoff(t *testing.T, r *cutoverRig, originals map[string]string, consumed bool) {
	t.Helper()
	require.Equal(t, originals, barrierAPIOriginals(t, r))
	require.Equal(t, 2, r.localCount(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE state='accepted' AND source_group_id=? AND target_group_id=?`, r.wide["id"], r.long["id"]))
	rows, err := r.f.db.Query(`SELECT s.sample_id,s.payload,s.payload_digest,s.created_at,s.consumed,t.origin_sample_id,t.payload,t.payload_digest,t.created_at,t.consumed,h.transfer_id,t.handoff_id,t.origin_group_id
		FROM wg_delivery_sample_handoffs h JOIN wg_delivery_samples s ON s.seq=h.source_seq JOIN wg_delivery_samples t ON t.seq=h.target_seq ORDER BY s.sample_id`)
	require.NoError(t, err)
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var sourceID, payload, digest, created, originID, targetPayload, targetDigest, targetCreated, transferID, handoffID, originGroup string
		var sourceConsumed, targetConsumed bool
		require.NoError(t, rows.Scan(&sourceID, &payload, &digest, &created, &sourceConsumed, &originID, &targetPayload, &targetDigest, &targetCreated, &targetConsumed, &transferID, &handoffID, &originGroup))
		require.Equal(t, sourceID, originID)
		require.Equal(t, originals[sourceID], targetPayload+"|"+targetDigest+"|"+targetCreated)
		require.Equal(t, payload+"|"+digest+"|"+created, originals[sourceID])
		require.Equal(t, r.wide["id"], originGroup)
		require.NotEmpty(t, transferID)
		require.Equal(t, transferID, handoffID)
		require.False(t, sourceConsumed, "transferred input must never be consumed by old owner")
		require.Equal(t, consumed, targetConsumed)
		ids = append(ids, sourceID)
	}
	require.NoError(t, rows.Err())
	require.Equal(t, []string{"after", "equal"}, ids, "exactly-boundary input belongs only to successor")
	var oldConsumed bool
	require.NoError(t, r.f.db.QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=? AND sample_id='before'`, r.wide["id"]).Scan(&oldConsumed))
	require.Equal(t, consumed, oldConsumed)
	require.Equal(t, 5, r.localCount(`SELECT COUNT(*) FROM wg_delivery_samples`), "three originals plus two exact copies; none lost or duplicated")
}

func barrierAPIAssertRows(t *testing.T, r *cutoverRig, wide bool, boundary, oldObserved, newObserved time.Time, oldValue, newValue float64) {
	t.Helper()
	oldRows := r.wideRowsInOrder()
	require.Len(t, oldRows, 1)
	require.NotNil(t, oldRows[0].First)
	require.Equal(t, oldValue, *oldRows[0].First)
	var entries []map[string]any
	require.NoError(t, json.Unmarshal([]byte(oldRows[0].Provenance), &entries))
	require.Len(t, entries, 1)
	require.Equal(t, "before", entries[0]["sample_id"])
	require.Equal(t, oldObserved, storedInstant(t, entries[0]["observed_at"].(string)))
	oldWindow := windowOfWide(t, oldRows[0].Provenance)
	require.Equal(t, "end", oldWindow["truncated"])
	require.Equal(t, boundary, storedInstant(t, oldWindow["end"].(string)))
	require.Equal(t, boundary.Truncate(15*time.Second), storedInstant(t, oldWindow["start"].(string)))
	if !wide {
		var found []longRecord
		for _, row := range r.longRecords() {
			if row.Start.Equal(boundary) {
				found = append(found, row)
			}
		}
		require.Len(t, found, 1)
		require.Equal(t, newValue, found[0].Value)
		require.Equal(t, newObserved, found[0].Observed)
		require.Equal(t, boundary.Truncate(20*time.Second).Add(20*time.Second), found[0].End)
		require.Equal(t, "start", windowOf(t, found[0].Provenance)["truncated"])
		return
	}
	member := r.long["members"].([]any)[0].(map[string]any)
	var value float64
	var provenance string
	require.NoError(t, r.dest.QueryRow(`SELECT "`+member["target_column"].(string)+`",provenance FROM `+r.quoted(r.longTable())+` WHERE bucket_start=?`, boundary.Format(time.RFC3339Nano)).Scan(&value, &provenance))
	require.Equal(t, newValue, value)
	require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
	require.Len(t, entries, 1)
	require.Equal(t, newObserved, storedInstant(t, entries[0]["observed_at"].(string)))
	window := windowOfWide(t, provenance)
	require.Equal(t, "start", window["truncated"])
	require.Equal(t, boundary, storedInstant(t, window["start"].(string)))
	require.Equal(t, boundary.Truncate(15*time.Second).Add(15*time.Second), storedInstant(t, window["end"].(string)))
	mean := entries[0]["mean"].(map[string]any)
	require.Equal(t, float64(2), mean["valid_count"], "boundary and after observations both enter the successor mean")
}
