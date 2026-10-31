package groupdelivery

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
)

func mixedHandoffDescriptor(t *testing.T, key GroupKey, method string, interval int, proof *counterproof.CounterBindingProof) []byte {
	t.Helper()
	var doc map[string]any
	require.NoError(t, json.Unmarshal(handoffProofDescriptor(t, key, proof, false), &doc))
	g := doc["group"].(map[string]any)
	g["recording_profile"] = map[string]any{"kind": "mixed", "result_shape": "member-calculation-v1"}
	g["layout_version"] = "managed_long_v1"
	delete(g, "row_policy")
	m := g["members"].([]any)[0].(map[string]any)
	m["calculation"] = method
	m["source_value_type"] = "int64"
	m["result_value_type"] = "int64"
	m["value_role"] = "gauge"
	m["cadence"] = map[string]any{"interval": interval, "unit": "second", "alignment": "utc_epoch"}
	if method == "mean" {
		m["result_value_type"] = "decimal"
	}
	if method == "delta" {
		m["value_role"] = "counter_raw"
		m["result_value_type"] = "decimal"
	}
	raw, err := json.Marshal(doc)
	require.NoError(t, err)
	return raw
}
func TestMixedHandoffSameIntentTransfersAndRestarts(t *testing.T) {
	for _, method := range []string{"latest", "mean"} {
		t.Run(method, func(t *testing.T) {
			store, db, path := newTestStore(t)
			store = store.WithClock(func() time.Time { return at(21) })
			for _, key := range []GroupKey{testKey, otherKey} {
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, mixedHandoffDescriptor(t, key, method, 15, nil)))
			}
			sample := sampleAt("mixed-handoff", ownershipMember, 23, 7)
			sample.ReceivedAt = at(21)
			require.NoError(t, store.AppendSample(t.Context(), testKey, at(15), sample))
			seedOwnershipCutover(t, db, "f-mixed", testKey, otherKey, at(22))
			var original string
			require.NoError(t, db.QueryRow(`SELECT payload FROM wg_delivery_samples`).Scan(&original))
			batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			require.NoError(t, err)
			require.Equal(t, 1, batch.Accepted)
			var payload, bucket string
			require.NoError(t, db.QueryRow(`SELECT payload,bucket_start FROM wg_delivery_samples WHERE group_id=?`, otherKey.GroupID).Scan(&payload, &bucket))
			require.Equal(t, original, payload)
			require.Equal(t, timestamp(at(15)), bucket)
			require.NoError(t, db.Close())
			reopened := NewStore(openStoreDB(t, path))
			restored, err := reopened.Restore(t.Context(), otherKey)
			require.NoError(t, err)
			require.Len(t, restored.Samples, 1)
			a, err := snapshot.NewAssembler(snapshot.Config{WorkspaceID: otherKey.WorkspaceID, GroupID: otherKey.GroupID, GroupRevision: otherKey.GroupRevision, DestinationScope: "fixture", StartsAt: at(22), Mode: snapshot.ModeMixedLong, Members: []snapshot.Member{{Key: ownershipMember, Interval: 15 * time.Second, Calculation: snapshot.Calculation(method), SourceValueType: schema.DataTypeInt64}}, Storage: snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}})
			require.NoError(t, err)
			require.Equal(t, snapshot.OfferSelected, a.Replay(restored.Samples[0]).Outcome)
			plan := a.Plan(at(30))
			require.NoError(t, plan.Err)
			require.Len(t, plan.Outcomes, 1)
			require.True(t, plan.Outcomes[0].Members[0].Usable())
			again, err := reopened.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			require.NoError(t, err)
			require.Equal(t, 1, again.Accepted)
		})
	}
}
func TestMixedHandoffRejectsChangedMeaning(t *testing.T) {
	for _, change := range []string{
		`json_set(payload,'$.group.members[0].calculation','mean','$.group.members[0].result_value_type','decimal')`,
		`json_set(payload,'$.group.members[0].cadence.interval',60)`,
		`json_set(payload,'$.group.members[0].cadence.alignment','local')`,
		`json_set(payload,'$.group.recording_profile.result_shape','future-v2')`,
		`json_set(payload,'$.group.recording_profile.kind','instantaneous','$.group.row_policy.interval_seconds',15)`,
		`json_set(payload,'$.group.members[0].measurement_id','other')`,
		`json_remove(payload,'$.group.members[0].source_value_type')`,
		`json_remove(payload,'$.group.members[0].result_value_type')`,
		`json_remove(payload,'$.group.members[0].value_role')`,
	} {
		t.Run(change, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			for _, key := range []GroupKey{testKey, otherKey} {
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, mixedHandoffDescriptor(t, key, "latest", 15, nil)))
			}
			_, err := db.Exec(`UPDATE wg_runtime_versions SET payload=`+change+` WHERE group_id=?`, otherKey.GroupID)
			require.NoError(t, err)
			sample := sampleAt("mixed-proof", ownershipMember, 23, 7)
			raw, digest, err := encodeSample(sample)
			require.NoError(t, err)
			_, err = validateHandoffProof(t.Context(), db, handoffProofInput{Origin: testKey, Sequence: 1, Payload: string(raw), Digest: digest, AcceptedAt: timestamp(at(21))}, otherKey)
			require.ErrorIs(t, err, ErrHandoffProofConflict)
		})
	}
}

func TestMixedHandoffDeltaKeepsOriginalSealedProof(t *testing.T) {
	store, db, path := newTestStore(t)
	store = store.WithClock(func() time.Time { return at(21) })
	origin, sample, _ := journalProofFixture(t)
	facts := origin.CounterBindingFacts
	facts.Member.GroupID = otherKey.GroupID
	target, err := counterproof.NewCounterBindingProof(facts)
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, mixedHandoffDescriptor(t, testKey, "delta", 15, &origin)))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, mixedHandoffDescriptor(t, otherKey, "delta", 15, &target)))
	sample.ObservedAt = at(23)
	sample.ReceivedAt = at(21)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(15), sample))
	seedOwnershipCutover(t, db, "f-delta", testKey, otherKey, at(22))
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	require.NoError(t, db.Close())
	reopened := NewStore(openStoreDB(t, path))
	r, err := reopened.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, r.Samples, 1)
	require.Equal(t, sample.CounterProof, r.Samples[0].CounterProof)
	targetRef, err := target.Reference()
	require.NoError(t, err)
	require.True(t, r.AuthorizesHandoff(r.Samples[0], targetRef))
	require.NotEqual(t, targetRef, r.Samples[0].CounterProof)
}

func TestMixedHandoffRejectsAlreadyClosedTargetMember(t *testing.T) {
	store, db, _ := newTestStore(t)
	for _, key := range []GroupKey{testKey, otherKey} {
		require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, mixedHandoffDescriptor(t, key, "latest", 15, nil)))
	}
	sample := sampleAt("late-handoff", ownershipMember, 23, 7)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(15), sample))
	seedOwnershipCutover(t, db, "f-closed", testKey, otherKey, at(22))
	a, err := snapshot.NewAssembler(snapshot.Config{WorkspaceID: otherKey.WorkspaceID, GroupID: otherKey.GroupID, GroupRevision: otherKey.GroupRevision, DestinationScope: "fixture", StartsAt: at(22), Mode: snapshot.ModeMixedLong, Members: []snapshot.Member{{Key: ownershipMember, Interval: 15 * time.Second, Calculation: snapshot.CalculationLatest, SourceValueType: schema.DataTypeInt64}}, Storage: snapshot.StorageCapabilities{NullableValues: true, MemberQuality: true}})
	require.NoError(t, err)
	cp, err := a.MixedCheckpointForPlan(a.Plan(at(30)))
	require.NoError(t, err)
	raw, err := json.Marshal(cp)
	require.NoError(t, err)
	descriptor, err := store.RuntimeVersion(t.Context(), otherKey)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO wg_delivery_mixed_progress(workspace_id,group_id,group_revision,version,descriptor_digest,closed_through,checkpoint,checkpoint_digest,last_closure_digest,updated_at) VALUES(?,?,?,1,?,?,?,?,?,?)`, otherKey.WorkspaceID, otherKey.GroupID, otherKey.GroupRevision, handoffHash(string(descriptor)), timestamp(at(30)), string(raw), handoffHash(string(raw)), "fixture", timestamp(at(31)))
	require.NoError(t, err)
	_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipConflict)
	require.Equal(t, 0, openCount(t, db, "wg_delivery_samples", "group_id='"+otherKey.GroupID+"'"))
}
