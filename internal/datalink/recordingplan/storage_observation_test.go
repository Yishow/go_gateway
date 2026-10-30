package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepairObservationRequiresNewCommittedLedgerSettlement(t *testing.T) {
	f := newRestoreFixture(t)
	f.drop(f.layout.TableName)
	var got []StorageSettlement
	ctx, finish := CaptureStorageSettlements(t.Context(), StorageObservationConfig{
		RecordLimit: 4, InfoEnabled: true, WarnEnabled: true,
		Observe: func(s StorageSettlement) {
			// Actual repository reentry occurs only after its transaction ended.
			_, err := f.svc.GetSchemaOperation(t.Context(), s.WorkspaceID, s.OperationID)
			require.NoError(t, err)
			got = append(got, s)
		},
	})
	result, err := f.restore(ctx)
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome)
	require.Empty(t, got, "repository must only capture; outer lock owner flushes")
	operation := *result.Operation
	result.Operation.OperationID = "mutated-return-object"
	finish()
	require.Len(t, got, 1, "new durable repair must produce a deferred settlement")
	require.Equal(t, StorageSettlement{Available: true, OwnedExecution: true, OperationID: operation.OperationID, WorkspaceID: f.scope.WorkspaceID,
		Status: SchemaOperationSucceeded, VerifiedDigest: f.digest, UpdatedAt: operation.UpdatedAt, CompletedAt: *operation.CompletedAt}, got[0])
	finish()
	require.Len(t, got, 1, "flush is idempotent")
	require.True(t, f.lineage().HistoryMayBeLost, "recreated schema does not recover historical data")
	require.Equal(t, int32(1), f.ddlRuns.Load())
	ctx, finish = CaptureStorageSettlements(t.Context(), StorageObservationConfig{RecordLimit: 4, InfoEnabled: true, WarnEnabled: true, Observe: func(s StorageSettlement) { got = append(got, s) }})
	again, err := f.restore(ctx)
	require.NoError(t, err)
	require.Equal(t, StorageNotMissing, again.Outcome)
	finish()
	require.Len(t, got, 1, "already restored is not a new operation settlement")
}

func TestRepairObservationRollbackUnacknowledgedAndStaleWinnerCannotSucceed(t *testing.T) {
	for _, mode := range []string{"rollback", "commit-ack", "readback", "stale-winner"} {
		t.Run(mode, func(t *testing.T) {
			f := newSettlementFixture(t)
			op := f.claim(t, "operation-1")
			var got []StorageSettlement
			ctx, finish := CaptureStorageSettlements(t.Context(), StorageObservationConfig{RecordLimit: 2, InfoEnabled: true, WarnEnabled: true, Observe: func(s StorageSettlement) { got = append(got, s) }})
			result := SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: "digest-1"}
			var apply func(context.Context, DBTX) error
			switch mode {
			case "rollback":
				apply = func(context.Context, DBTX) error { return errors.New("lineage rollback") }
			case "commit-ack":
				f.fault.commitAck.Store(true)
			case "readback":
				f.fault.readback.Store(true)
			case "stale-winner":
				_, err := f.repo.finishStorageRepair(t.Context(), op.OperationID, op.Owner, result, nil)
				require.NoError(t, err)
			}
			var settled *SchemaOperation
			var err error
			if mode == "stale-winner" {
				settled, err = f.repo.finishStaleStorageRepair(ctx, op.OperationID, result, time.Now().Add(time.Hour), nil)
				require.NoError(t, err)
				require.Nil(t, settled)
			} else {
				settled, err = f.repo.finishStorageRepair(ctx, op.OperationID, op.Owner, result, apply)
				require.Error(t, err)
				require.Nil(t, settled)
			}
			require.Empty(t, got)
			finish()
			if mode == "commit-ack" || mode == "readback" {
				require.Equal(t, []StorageSettlement{{}}, got, "uncertainty must not become completed")
			} else {
				require.Empty(t, got, "no newly committed row means no settlement")
			}
			stored, err := f.repo.GetSchemaOperation(t.Context(), op.WorkspaceID, op.OperationID)
			require.NoError(t, err)
			if mode == "rollback" {
				require.Equal(t, SchemaOperationRunning, stored.Status)
			} else {
				require.Equal(t, SchemaOperationSucceeded, stored.Status, "observation cannot rewrite native result")
			}
		})
	}
}

func TestRepairObservationRejectsMismatchedReadback(t *testing.T) {
	for _, column := range []string{"operation_id", "status", "verified_digest", "updated_at", "completed_at", "action"} {
		t.Run(column, func(t *testing.T) {
			f := newSettlementFixture(t)
			op := f.claim(t, "operation-1")
			var got []StorageSettlement
			ctx, finish := CaptureStorageSettlements(t.Context(), StorageObservationConfig{RecordLimit: 1, InfoEnabled: true, WarnEnabled: true, Observe: func(s StorageSettlement) { got = append(got, s) }})
			f.fault.afterCommit = func(conn settlementConn) {
				value := "changed"
				if column == "status" {
					value = "failed"
				}
				if column == "updated_at" || column == "completed_at" {
					value = "2000-01-01 00:00:00+00:00"
				}
				_, err := conn.ExecContext(t.Context(), "UPDATE managed_schema_operations SET "+column+" = ? WHERE operation_id = ?", namedValues(value, op.OperationID))
				require.NoError(t, err)
			}
			_, _ = f.repo.finishStorageRepair(ctx, op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: "digest-1"}, nil)
			finish()
			require.Equal(t, []StorageSettlement{{}}, got)
		})
	}
}
func TestRepairObservationBoundsCaptureAndFiltersBeforeAdmission(t *testing.T) {
	for _, mode := range []string{"bounded", "zero-cap", "info-filtered", "warn-filtered", "all-filtered"} {
		t.Run(mode, func(t *testing.T) {
			f := newSettlementFixture(t)
			var got []StorageSettlement
			var dropped uint64
			cfg := StorageObservationConfig{RecordLimit: 1, InfoEnabled: true, WarnEnabled: true,
				Observe: func(s StorageSettlement) { got = append(got, s) }, OnDrop: func(n uint64) { dropped += n }}
			switch mode {
			case "zero-cap":
				cfg.RecordLimit = 0
			case "info-filtered":
				cfg.InfoEnabled = false
			case "warn-filtered":
				cfg.WarnEnabled = false
			case "all-filtered":
				cfg.InfoEnabled, cfg.WarnEnabled = false, false
			}
			ctx, finish := CaptureStorageSettlements(t.Context(), cfg)
			for i, status := range []SchemaOperationStatus{SchemaOperationSucceeded, SchemaOperationSucceeded, SchemaOperationFailed} {
				op := f.claim(t, fmt.Sprintf("operation-%d", i))
				_, err := f.repo.finishStorageRepair(ctx, op.OperationID, op.Owner, SchemaOperationResult{Status: status, VerifiedDigest: "digest-1"}, nil)
				require.NoError(t, err)
			}
			require.Empty(t, got)
			require.Zero(t, dropped, "no callbacks inside native operation")
			finish()
			switch mode {
			case "bounded":
				require.Len(t, got, 1)
				require.Equal(t, uint64(2), dropped)
			case "zero-cap":
				require.Empty(t, got)
				require.Equal(t, uint64(3), dropped)
			case "info-filtered":
				require.Len(t, got, 1)
				require.Equal(t, SchemaOperationFailed, got[0].Status)
				require.Zero(t, dropped)
			case "warn-filtered":
				require.Len(t, got, 1)
				require.Equal(t, SchemaOperationSucceeded, got[0].Status)
				require.Equal(t, uint64(1), dropped)
			case "all-filtered":
				require.Empty(t, got)
				require.Zero(t, dropped)
			}
		})
	}
}

func TestRepairObservationNilAndFilteredCaptureAddNoReads(t *testing.T) {
	for _, mode := range []string{"nil", "filtered", "enabled"} {
		t.Run(mode, func(t *testing.T) {
			f := newSettlementFixture(t)
			op := f.claim(t, "operation-1")
			cfg := StorageObservationConfig{RecordLimit: 1, InfoEnabled: true, WarnEnabled: true}
			if mode != "nil" {
				cfg.Observe = func(StorageSettlement) {}
			}
			if mode == "filtered" {
				cfg.InfoEnabled, cfg.WarnEnabled = false, false
			}
			ctx, finish := CaptureStorageSettlements(t.Context(), cfg)
			before := f.fault.queries.Load()
			_, err := f.repo.finishStorageRepair(ctx, op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: "digest-1"}, nil)
			require.NoError(t, err)
			finish()
			require.Equal(t, int64(1), f.fault.queries.Load()-before, "only the existing native postcommit query")
		})
	}
}

func TestRepairObservationCallbackPanicAndReentryPreserveNativeOutcome(t *testing.T) {
	f := newSettlementFixture(t)
	op := f.claim(t, "operation-1")
	var finish func()
	var callbackCount int
	var dropped uint64
	cfg := StorageObservationConfig{RecordLimit: 1, InfoEnabled: true, WarnEnabled: true,
		Observe: func(s StorageSettlement) {
			callbackCount++
			finish() // Idempotence must not hold a mutex while invoking callbacks.
			stored, err := f.repo.GetSchemaOperation(t.Context(), s.WorkspaceID, s.OperationID)
			require.NoError(t, err)
			require.Equal(t, SchemaOperationSucceeded, stored.Status)
			panic("sink failed")
		}, OnDrop: func(n uint64) { dropped += n; panic("loss sink failed") }}
	ctx, finishCapture := CaptureStorageSettlements(t.Context(), cfg)
	finish = finishCapture
	// The policy copy is immutable, not a live configuration pointer.
	cfg.Observe, cfg.OnDrop = nil, nil
	nested, innerFinish := CaptureStorageSettlements(ctx, StorageObservationConfig{RecordLimit: 1, InfoEnabled: true, Observe: func(StorageSettlement) { t.Error("nested callback must not replace outer owner") }})
	stored, err := f.repo.finishStorageRepair(nested, op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: "digest-1"}, nil)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationSucceeded, stored.Status)
	innerFinish()
	require.Zero(t, callbackCount)
	require.NotPanics(t, finish)
	finish()
	require.Equal(t, 1, callbackCount)
	require.Equal(t, uint64(1), dropped)
}

func TestRepairObservationInvalidMetadataIsIdentityFreeUnavailable(t *testing.T) {
	for _, field := range []string{"workspace", "digest", "unknown-status"} {
		t.Run(field, func(t *testing.T) {
			f := newSettlementFixture(t)
			op := f.claim(t, "operation-1")
			result := SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: "digest-1", Reason: "password_fixture", Detail: `{"raw":"private-path"}`}
			switch field {
			case "workspace":
				_, err := f.db.ExecContext(t.Context(), "UPDATE managed_schema_operations SET workspace_id = ?", "http://private.example")
				require.NoError(t, err)
			case "digest":
				result.VerifiedDigest = "password=fixture"
			case "unknown-status":
				result.Status = "not-native-terminal"
			}
			var got []StorageSettlement
			ctx, finish := CaptureStorageSettlements(t.Context(), StorageObservationConfig{RecordLimit: 1, InfoEnabled: true, WarnEnabled: true, Observe: func(s StorageSettlement) { got = append(got, s) }})
			_, err := f.repo.finishStorageRepair(ctx, op.OperationID, op.Owner, result, nil)
			require.NoError(t, err, "diagnostics must not change existing native result policy")
			finish()
			require.Equal(t, []StorageSettlement{{}}, got)
		})
	}
}

func TestRepairObservationDistinguishesOwnedExecutionFromStaleSettlement(t *testing.T) {
	for _, mode := range []string{"owned", "stale", "stale-warn-only"} {
		t.Run(mode, func(t *testing.T) {
			f := newSettlementFixture(t)
			op := f.claim(t, "operation-1")
			var got []StorageSettlement
			ctx, finish := CaptureStorageSettlements(t.Context(), StorageObservationConfig{RecordLimit: 1,
				InfoEnabled: mode != "stale-warn-only", WarnEnabled: true, Observe: func(s StorageSettlement) { got = append(got, s) }})
			result := SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: "digest-1"}
			var err error
			if mode == "owned" {
				_, err = f.repo.finishStorageRepair(ctx, op.OperationID, op.Owner, result, nil)
			} else {
				_, err = f.repo.finishStaleStorageRepair(ctx, op.OperationID, result, time.Now().Add(time.Hour), nil)
			}
			require.NoError(t, err)
			finish()
			require.Len(t, got, 1, "stale settlement cannot borrow the current Restore input, even with info disabled")
			require.True(t, got[0].Available, "the real settled ledger is known independently of input ownership")
			require.Equal(t, mode == "owned", got[0].OwnedExecution)
		})
	}
}
