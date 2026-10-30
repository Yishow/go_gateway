package grouppipeline

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestRepairObservationRequiresNewCommittedLedgerSettlement(t *testing.T) {
	f := newStorageObservationFixture(t)
	f.drop()
	decision, err := f.restore()
	require.NoError(t, err)
	require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
	require.Len(t, f.events, 1)
	o := f.events[0]
	require.True(t, o.Available)
	require.Equal(t, StorageObservationCompleted, o.Outcome)
	require.Equal(t, f.item.Key, o.Key)
	require.Equal(t, f.item.Destination.ConnectorRevision, o.ConnectorRevision)
	require.Equal(t, "observed-wrev", o.WorkspaceRevision)
	op, err := f.ledger.GetSchemaOperation(t.Context(), o.Key.WorkspaceID, o.OperationID)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status)
	require.Equal(t, op.OperationID, f.lineage().LastReconstructionOperationID)
	f.events = nil
	_, err = f.restore()
	require.NoError(t, err)
	require.Empty(t, f.events, "already present/no-op must not replay old settlement")
}
func TestRepairObservationKeepsOriginalFrozenRevisionAfterDraftChange(t *testing.T) {
	f := newStorageObservationFixture(t)
	f.groups.groups[0].Revision = "saved-latest"
	f.groups.groups[0].AppliedRevision = "requested-next"
	f.drop()
	_, err := f.restore()
	require.NoError(t, err)
	require.Len(t, f.events, 1)
	require.True(t, f.events[0].Available)
	require.Equal(t, "rev-1", f.events[0].Key.GroupRevision)
	require.Equal(t, "crev-1", f.events[0].ConnectorRevision)
}
func TestRepairObservationRecreationDoesNotClaimHistoricalDataRecovered(t *testing.T) {
	f := newStorageObservationFixture(t)
	f.drop()
	_, err := f.restore()
	require.NoError(t, err)
	require.True(t, f.lineage().HistoryMayBeLost)
	require.Zero(t, f.count(f.target.DB, `SELECT COUNT(*) FROM "`+f.table+`"`))
	require.Positive(t, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state!='sql_committed'`))
	require.Len(t, f.events, 1)
	require.Equal(t, StorageObservationCompleted, f.events[0].Outcome)
	require.Equal(t, "none", f.events[0].Reason, "completion carries no fabricated delivery/readback result")
}
func TestRepairObservationCallbacksRunAfterScopeLock(t *testing.T) {
	f := newStorageObservationFixture(t)
	entered := false
	f.recovery.observation.Observe = func(o StorageOperationObservation) {
		f.events = append(f.events, o)
		if entered {
			return
		}
		entered = true
		release := f.recovery.lock(f.scope)
		release()
		_, err := f.restore()
		require.NoError(t, err)
		panic("capture-password")
	}
	f.drop()
	_, err := f.restore()
	require.NoError(t, err)
	require.True(t, entered)
	require.Equal(t, 1, f.lineage().Reconstructions)
}
func TestRepairObservationMissingChangedAndForeignProofIsUnavailable(t *testing.T) {
	for _, mode := range []string{"missing-wrev", "changed-wrev", "foreign", "panic", "frozen-change", "ledger-read"} {
		t.Run(mode, func(t *testing.T) {
			f := newStorageObservationFixture(t)
			f.drop()
			f.recovery.observation.ScopeRevision = func(ctx context.Context, key groupdelivery.GroupKey) (string, bool) {
				f.reads++
				switch mode {
				case "missing-wrev":
					return "", true
				case "changed-wrev":
					if f.reads == 2 {
						return "changed", true
					}
				case "foreign":
					return "foreign", false
				case "panic":
					panic("private")
				case "frozen-change":
					if f.reads == 2 {
						raw, err := f.store.RuntimeVersion(ctx, key)
						require.NoError(t, err)
						_, err = f.local.ExecContext(ctx, `UPDATE wg_runtime_versions SET payload=?`, strings.Replace(string(raw), "crev-1", "crev-2", 1))
						require.NoError(t, err)
					}
				case "ledger-read":
					reject := func(q string) error {
						if strings.Contains(q, "FROM managed_schema_operations") {
							return errors.New("private provider")
						}
						return nil
					}
					f.localFault.Reject.Store(&reject)
				}
				return "observed-wrev", true
			}
			_, err := f.restore()
			f.localFault.Reject.Store(nil)
			require.NoError(t, err)
			require.Equal(t, []StorageOperationObservation{{}}, f.events)
		})
	}
}
func TestRepairObservationPreclaimRefusalWithoutLedgerIsUnavailable(t *testing.T) {
	f := newStorageObservationFixture(t)
	disabled := false
	f.snapshot.Group.WritePolicy.RecreateMissing = &disabled
	f.recovery.snapshots.Store(f.item.Key, f.snapshot)
	f.drop()
	decision, err := f.restore()
	require.NoError(t, err)
	require.Equal(t, groupdelivery.RestoreBlock, decision.Action)
	require.Equal(t, []StorageOperationObservation{{}}, f.events)
	require.Zero(t, f.count(f.local, `SELECT COUNT(*) FROM managed_schema_operations WHERE action='storage_repair'`))
}

func TestRepairObservationRestartWithoutOriginalRevisionIsUnavailable(t *testing.T) {
	f := newStorageObservationFixture(t)
	f.recovery.observation = StorageObservationConfig{}
	f.drop()
	reject := func(q string) error {
		if strings.Contains(q, "UPDATE managed_schema_operations") && strings.Contains(q, "SET status") {
			return errors.New("local settlement blocked")
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	_, err := f.restore()
	f.localFault.Reject.Store(nil)
	require.Error(t, err)
	_, err = f.local.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at=? WHERE action=?`, time.Now().UTC().Add(-2*recordingplan.SchemaOperationLease), recordingplan.StorageRepairAction)
	require.NoError(t, err)
	r := NewStorageRecoveryWithObservation(f.store, f.ledger, f.dest, StorageObservationConfig{RecordLimit: 1024, InfoEnabled: true, WarnEnabled: true, ScopeRevision: func(context.Context, groupdelivery.GroupKey) (string, bool) { f.reads++; return "observed-wrev", true }, Observe: func(o StorageOperationObservation) { f.events = append(f.events, o) }})
	require.NoError(t, r.ResolveStale(t.Context()))
	require.Equal(t, []StorageOperationObservation{{}}, f.events)
	require.Zero(t, f.reads, "no original revision proof means no borrowed workspace attribution")
	require.True(t, f.lineage().HistoryMayBeLost)
}

func TestRepairObservationOldStaleInsideFreshRestoreCannotBorrowInputRevision(t *testing.T) {
	f := newStorageObservationFixture(t)
	cfg := f.recovery.observation
	f.recovery.observation = StorageObservationConfig{}
	f.drop()
	reject := func(q string) error {
		if strings.Contains(q, "UPDATE managed_schema_operations") && strings.Contains(q, "SET status") {
			return errors.New("local settlement blocked")
		}
		return nil
	}
	f.localFault.Reject.Store(&reject)
	_, err := f.restore()
	f.localFault.Reject.Store(nil)
	require.Error(t, err)
	var oldID string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT operation_id FROM managed_schema_operations WHERE action=?`, recordingplan.StorageRepairAction).Scan(&oldID))
	_, err = f.local.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at=? WHERE operation_id=?`, time.Now().UTC().Add(-2*recordingplan.SchemaOperationLease), oldID)
	require.NoError(t, err)
	f.drop()
	f.recovery = NewStorageRecoveryWithObservation(f.store, f.ledger, f.dest, cfg)
	decision, err := f.restore()
	require.NoError(t, err)
	require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
	require.Len(t, f.events, 2)
	require.Equal(t, StorageOperationObservation{}, f.events[0])
	require.True(t, f.events[1].Available)
	require.NotEqual(t, oldID, f.events[1].OperationID)
}

func TestRepairObservationRollbackAndUnacknowledgedCannotSucceed(t *testing.T) {
	for _, mode := range []string{"rollback", "unacknowledged"} {
		t.Run(mode, func(t *testing.T) {
			f := newStorageObservationFixture(t)
			f.drop()
			reject := func(q string) error {
				if mode == "rollback" && strings.HasPrefix(strings.TrimSpace(q), "CREATE TABLE") {
					return errors.New("private target failure")
				}
				if mode == "unacknowledged" && strings.Contains(q, "UPDATE managed_schema_operations") && strings.Contains(q, "SET status") {
					return errors.New("private commit failure")
				}
				return nil
			}
			if mode == "rollback" {
				f.targetFlt.Reject.Store(&reject)
			} else {
				f.localFault.Reject.Store(&reject)
			}
			decision, err := f.restore()
			f.targetFlt.Reject.Store(nil)
			f.localFault.Reject.Store(nil)
			if mode == "rollback" {
				require.NoError(t, err)
				require.Equal(t, groupdelivery.RestoreBlock, decision.Action)
			} else {
				require.Error(t, err)
			}
			require.Len(t, f.events, 1)
			require.NotEqual(t, StorageObservationCompleted, f.events[0].Outcome)
			if mode == "rollback" {
				require.True(t, f.events[0].Available)
				require.Equal(t, StorageObservationBlocked, f.events[0].Outcome)
			} else {
				require.Equal(t, StorageOperationObservation{}, f.events[0])
			}
		})
	}
}

func TestRepairObservationNilFilteredAndFullCapturePreserveNativeResult(t *testing.T) {
	for _, mode := range []string{"nil", "filtered", "full", "drop-panic"} {
		t.Run(mode, func(t *testing.T) {
			f := newStorageObservationFixture(t)
			var drops uint64
			f.recovery.observation.OnDrop = func(n uint64) {
				drops += n
				if mode == "drop-panic" {
					panic("private")
				}
			}
			switch mode {
			case "nil":
				f.recovery.observation.Observe = nil
			case "filtered":
				f.recovery.observation.InfoEnabled = false
			default:
				f.recovery.observation.RecordLimit = 0
			}
			f.drop()
			decision, err := f.restore()
			require.NoError(t, err)
			require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
			require.Empty(t, f.events)
			require.Zero(t, f.reads, "nil/filtered/dropped performs no diagnostic scope reads")
			require.Equal(t, 1, f.lineage().Reconstructions)
			if mode == "full" || mode == "drop-panic" {
				require.Equal(t, uint64(1), drops)
			} else {
				require.Zero(t, drops)
			}
		})
	}
}

func TestRepairObservationFinalProofChangesCannotBorrowCompletedResult(t *testing.T) {
	for _, mode := range []string{"operation", "outbox", "lineage"} {
		t.Run(mode, func(t *testing.T) {
			f := newStorageObservationFixture(t)
			f.recovery.observation.ScopeRevision = func(ctx context.Context, _ groupdelivery.GroupKey) (string, bool) {
				f.reads++
				if f.reads == 2 {
					var err error
					switch mode {
					case "operation":
						_, err = f.local.ExecContext(ctx, `UPDATE managed_schema_operations SET verified_digest='wrong' WHERE action=?`, recordingplan.StorageRepairAction)
					case "outbox":
						_, err = f.local.ExecContext(ctx, `UPDATE wg_delivery_outbox SET payload='private-payload' WHERE effect_key=?`, f.item.EffectKey)
					case "lineage":
						_, err = f.local.ExecContext(ctx, `UPDATE write_group_storage_scopes SET confirmed_connector_revision='other'`)
					}
					require.NoError(t, err)
				}
				return "observed-wrev", true
			}
			f.drop()
			decision, err := f.restore()
			require.NoError(t, err)
			require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
			require.Equal(t, []StorageOperationObservation{{}}, f.events)
		})
	}
}

func TestRepairObservationUnknownCatalogCannotBecomeCompleted(t *testing.T) {
	f := newStorageObservationFixture(t)
	f.drop()
	// Fail only post-DDL catalog inspection, not the actual transaction.
	//nolint:unparam // Reject's callback signature requires an error return.
	reject := func(q string) error {
		if strings.HasPrefix(strings.TrimSpace(q), "CREATE TABLE") {
			post := func(q string) error {
				if strings.Contains(q, "sqlite_master") {
					return errors.New("private-catalog")
				}
				return nil
			}
			f.targetFlt.Reject.Store(&post)
		}
		return nil
	}
	f.targetFlt.Reject.Store(&reject)
	decision, err := f.restore()
	f.targetFlt.Reject.Store(nil)
	require.NoError(t, err)
	require.Equal(t, groupdelivery.RestoreBlock, decision.Action)
	require.Equal(t, []StorageOperationObservation{{}}, f.events)
	var status string
	require.NoError(t, f.local.QueryRowContext(t.Context(), `SELECT status FROM managed_schema_operations WHERE action=?`, recordingplan.StorageRepairAction).Scan(&status))
	require.Equal(t, string(recordingplan.SchemaOperationPartial), status)
}
