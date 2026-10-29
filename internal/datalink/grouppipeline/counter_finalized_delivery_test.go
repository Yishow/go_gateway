package grouppipeline

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func firstCounterEffect(t *testing.T, f *longPipeline) groupdelivery.OutboxItem {
	t.Helper()
	var effect string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT effect_key FROM wg_delivery_outbox WHERE group_id=? AND entity_key=? ORDER BY bucket_start LIMIT 1`, f.group.ID, `["device-1","point-1","tag-kwh-1"]`).Scan(&effect))
	item, err := f.store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	return item
}

func counterFinalizedSender(f *longPipeline) *groupdelivery.Sender {
	return groupdelivery.NewSender(f.store, destinationResolver{destinations: f.dest, loadRuntimeSnapshot: f.pipe.loadFinalizedSnapshot}, f.pipe.senderConfig())
}

func TestCounterFinalizedDeliveryRetainsMixedBacklogAndIndependentReceipts(t *testing.T) {
	for _, reason := range []string{"counter-replay-proof-missing", "counter-replay-proof-conflict", "counter-replay-proof-unsupported"} {
		t.Run(reason, func(t *testing.T) {
			f := newLongPipeline(t)
			healthy := healthyCounterRevision(t, f, "healthy-group", "healthy-revision")
			require.NoError(t, f.reading(0, "finalized-a", at(1), 100))
			f.closeUntil(at(600))
			original := firstCounterEffect(t, f)
			baselineEnd, baselinePayload := f.baseline(0)
			require.NoError(t, f.reading(0, "quarantined-a", at(601), 105))
			damageCounterProof(t, f, reason)
			before := counterDurableState(t, f.store)
			f.groups.groups = append(f.groups.groups, liveGroup(healthy.Group, workspace.WriteGroupStatusReady))
			f.groups.snapshots[healthy.GroupID] = []*workspace.WriteGroupAppliedSnapshot{healthy}
			// Current draft settings cannot redirect the already-finalized effect.
			f.groups.groups[0].Destination.TableName = "edited_current_destination"
			reader := &unavailableExactCounterReader{}
			restart := func() {
				f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }, Sender: groupdelivery.SenderConfig{Backoff: func(int) (string, time.Time) { return groupdelivery.StateRetrying, at(0) }}})
				require.NoError(t, f.pipe.Reconcile(t.Context()))
			}
			restart()
			f.targetFlt.Lose.Store(true)
			result, err := counterFinalizedSender(f).Deliver(t.Context(), original.EffectKey)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateRetrying, result.State)
			require.Equal(t, 1, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
			require.Equal(t, 0, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_receipts`))
			f.targetFlt.Lose.Store(false)
			restart()
			result, err = counterFinalizedSender(f).Deliver(t.Context(), original.EffectKey)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, result.State)
			after, err := f.store.GetOutbox(t.Context(), original.EffectKey)
			require.NoError(t, err)
			require.Equal(t, original.Payload, after.Payload)
			require.Equal(t, original.PayloadDigest, after.PayloadDigest)
			require.Equal(t, original.Destination, after.Destination)
			var raw, unit, epoch string
			require.NoError(t, f.plainTarget().QueryRowContext(t.Context(), `SELECT CAST(value_decimal AS TEXT),unit,series_epoch FROM `+f.table+` WHERE record_id=?`, original.RecordID).Scan(&raw, &unit, &epoch))
			require.Equal(t, []string{"100", "kWh", "epoch-1"}, []string{raw, unit, epoch})
			require.Equal(t, 1, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
			var receipts int
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_receipts WHERE effect_key=? AND payload_digest=?`, original.EffectKey, original.PayloadDigest).Scan(&receipts))
			require.Equal(t, 1, receipts)
			require.NoError(t, f.plainTarget().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM gw_effect_receipts WHERE effect_key=? AND payload_digest=?`, original.EffectKey, original.PayloadDigest).Scan(&receipts))
			require.Equal(t, 1, receipts, "lost ACK retry retains exactly one matching destination receipt")
			current := counterDurableState(t, f.store)
			for _, table := range []string{"wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_counter_baselines", "wg_runtime_versions", "wg_counter_proof_associations"} {
				require.Equal(t, before[table], current[table], table)
			}
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "healthy-1", at(602), uint64(7))))
			f.closeUntil(at(1200))
			f.deliver()
			var healthyRows int
			require.NoError(t, f.plainTarget().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+f.table+` WHERE group_id='healthy-group' AND tag_id='healthy-tag' AND CAST(value_decimal AS TEXT)='7' AND unit='kWh' AND series_epoch='epoch-1'`).Scan(&healthyRows))
			require.Equal(t, 1, healthyRows)
			view, err := f.pipe.Delivery(t.Context(), f.group.ID)
			require.NoError(t, err)
			require.Equal(t, []CounterReplayQuarantine{{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1", PendingSamples: 1, Reason: reason, Action: proofRepairAction}}, view.CounterReplayQuarantines)
			require.Equal(t, 1, view.Stages.Collecting)
			require.Equal(t, 2, view.Stages.SQLCommitted)
			require.Zero(t, view.Stages.Quarantined)
			require.Equal(t, StateBlocked, view.Intake.State)
			require.NotNil(t, view.LastSQLCommittedAt)
			finalEnd, finalBaseline := f.baseline(0)
			require.Equal(t, baselineEnd, finalEnd)
			require.Equal(t, baselinePayload, finalBaseline)
			var nextClose string
			require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT next_close FROM wg_delivery_checkpoints WHERE group_id='group-C' AND group_revision='rev-1'`).Scan(&nextClose))
			require.Equal(t, "2026-01-01T00:10:00.000000000Z", nextClose, "healthy closure and finalized delivery cannot advance quarantined inputs")
			require.Zero(t, reader.calls)
		})
	}
}

func TestCounterFinalizedProofDamageDoesNotBypassStorageRepairGate(t *testing.T) {
	for _, reason := range []string{"counter-replay-proof-missing", "counter-replay-proof-conflict", "counter-replay-proof-unsupported"} {
		t.Run(reason, func(t *testing.T) {
			f := newLongPipeline(t)
			require.NoError(t, f.reading(0, "finalized", at(1), 100))
			f.closeUntil(at(600))
			original := firstCounterEffect(t, f)
			scope := workspace.StorageScopeOf(f.group)
			_, err := f.local.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes (workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,repair_state) VALUES (?,?,?,?,?,?,?,?,?,?)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion, f.table, "owner", "original-layout", recordingplan.StorageRepairStateBlocked)
			require.NoError(t, err)
			damageCounterProof(t, f, reason)
			result, err := counterFinalizedSender(f).Deliver(t.Context(), original.EffectKey)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateBlocked, result.State)
			require.Equal(t, recordingplan.StorageCodeRepairBlocked, result.Code)
			require.Equal(t, 0, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
			require.Equal(t, 0, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_receipts`))
		})
	}
}
