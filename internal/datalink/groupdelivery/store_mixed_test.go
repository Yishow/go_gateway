package groupdelivery

import (
	"encoding/json"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mixedStorePlan(t *testing.T, store *Store) (*snapshot.Assembler, Closure) {
	t.Helper()
	members := []snapshot.Member{{Key: `["d","p","fast"]`, SourceValueType: schema.DataTypeInt64, Calculation: snapshot.CalculationLatest, Interval: 15 * time.Second}, {Key: `["d","p","slow"]`, SourceValueType: schema.DataTypeInt64, Calculation: snapshot.CalculationLatest, Interval: 60 * time.Second}}
	a, err := snapshot.NewAssembler(snapshot.Config{WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision, DestinationScope: "destination", StartsAt: t0, Mode: snapshot.ModeMixedWide, Members: members, IncompletePolicy: snapshot.IncompletePartial, Storage: snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}})
	require.NoError(t, err)
	payload := []byte(`{"workspace_id":"workspace-A","group_id":"group-G","applied_revision":"rev-1","group":{"workspace_id":"workspace-A","id":"group-G","revision":"rev-1","applied_revision":"rev-1","recording_profile":{"kind":"mixed","result_shape":"member-calculation-v1"},"members":[{"device_id":"d","point_id":"p","tag_id":"fast","calculation":"latest","source_value_type":"int64","cadence":{"interval":15,"unit":"second","alignment":"utc_epoch"}},{"device_id":"d","point_id":"p","tag_id":"slow","calculation":"latest","source_value_type":"int64","cadence":{"interval":60,"unit":"second","alignment":"utc_epoch"}}]}}`)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
	for _, m := range members {
		sample := sampleAt(m.Key, m.Key, 10, 42)
		require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
		a.Offer(sample, at(10))
	}
	before, err := a.MixedCheckpoint()
	require.NoError(t, err)
	plan := a.Plan(at(15))
	require.NoError(t, plan.Err)
	after, err := a.MixedCheckpointForPlan(plan)
	require.NoError(t, err)
	// Required baseline gaps use skipped outcomes; this fixture needs no encoder.
	buckets := make([]ClosedBucket, len(plan.Outcomes))
	for i, o := range plan.Outcomes {
		o.Kind = snapshot.OutcomeSkipped
		o.EffectKey = ""
		buckets[i] = ClosedBucket{Outcome: o}
	}
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	return a, Closure{Destination: testDestination, Buckets: buckets, Mixed: &MixedClosureProgress{Before: before, After: after, DescriptorDigest: restored.DescriptorDigest}}
}
func TestMixedClosureAtomicSkipProgressKeepsSlowWindow(t *testing.T) {
	store, db, path := newTestStore(t)
	_, closure := mixedStorePlan(t, store)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure))
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
	require.Zero(t, countRows(t, db, "wg_delivery_buckets"))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_mixed_progress"))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_mixed_outcomes"))
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closure), "exact replay")
	restored, err := NewStore(openStoreDB(t, path)).Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.NotNil(t, restored.Mixed)
	require.Equal(t, closure.Mixed.After, *restored.Mixed)
	require.True(t, restored.NextClose.IsZero())
	require.Len(t, restored.Samples, 1)
	require.Equal(t, `["d","p","slow"]`, restored.Samples[0].MemberKey)
	changed := closure
	copyProgress := *closure.Mixed
	changed.Mixed = &copyProgress
	changed.Mixed.DescriptorDigest = "foreign"
	require.Error(t, store.CommitClosure(t.Context(), testKey, changed))
	var data string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT checkpoint FROM wg_delivery_mixed_progress`).Scan(&data))
	var cp snapshot.MixedCheckpoint
	require.NoError(t, json.Unmarshal([]byte(data), &cp))
	require.Equal(t, closure.Mixed.After, cp)
}
func TestMixedClosureAtomicFaultRollsBackProgressAndConsumption(t *testing.T) {
	for _, table := range []string{"wg_delivery_mixed_outcomes", "wg_delivery_mixed_progress"} {
		t.Run(table, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			_, c := mixedStorePlan(t, store)
			failOn(t, db, table)
			require.Error(t, store.CommitClosure(t.Context(), testKey, c))
			require.Zero(t, countRows(t, db, "wg_delivery_mixed_progress"))
			require.Zero(t, countRows(t, db, "wg_delivery_mixed_outcomes"))
			r, err := store.Restore(t.Context(), testKey)
			require.NoError(t, err)
			require.Len(t, r.Samples, 2)
		})
	}
}

func TestMixedClosureOutboxKeepsDispatchWithoutFakeBucket(t *testing.T) {
	store, db, _ := newTestStore(t)
	_, c := mixedStorePlan(t, store)
	b := &c.Buckets[0]
	b.Outcome.Kind = snapshot.OutcomeRow
	b.Outcome.EffectKey, _ = snapshot.EffectKey(c.Destination.Scope, b.Outcome.RecordID)
	row := rowBucket("", t0, 9007199254740993).Row
	row.RecordID = b.Outcome.RecordID
	row.EffectKey = b.Outcome.EffectKey
	row.BucketStart = time.Time{}
	row.DueAt = b.Outcome.DueAt
	b.Row = row
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	var bucket string
	var due *string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT bucket_start,due_at FROM wg_delivery_outbox`).Scan(&bucket, &due))
	require.Empty(t, bucket)
	require.NotNil(t, due)
	require.Equal(t, timestamp(at(15)), *due)
	item, err := store.GetOutbox(t.Context(), row.EffectKey)
	require.NoError(t, err)
	require.True(t, item.BucketStart.IsZero())
	decoded, err := DecodeRowPayload(item.Payload)
	require.NoError(t, err)
	require.Equal(t, at(15), decoded.DueAt)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, c))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"))
}
