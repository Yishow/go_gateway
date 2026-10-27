package recordingplan

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestMySQLStorageRestoreStaleRepairAlteredDefaultBlocksAndRetainsReceipt verifies that when
// an expired running repair finds a physical table present whose actual catalog defaults have
// been altered from the confirmed layout, ResolveStaleStorageRepair refuses to mark the attempt
// Succeeded or advance the generation; it sets Partial, blocks the repair state, executes no DDL,
// and leaves the receipt table untouched.
func TestMySQLStorageRestoreStaleRepairAlteredDefaultBlocksAndRetainsReceipt(t *testing.T) {
	f := newMySQLRestoreFixture(t)

	// Simulate physical table present with an altered column default
	group := cloneTargetTableInspection(f.ddl.groupExact)
	alteredDefault := "altered_default_val"
	group.MySQL.ColumnDefaults["value_text"] = &alteredDefault
	f.ddl.states[f.ddl.layout.TableName] = group

	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	repo := f.service.repo.(storageRepairRepository)
	claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)

	// Age the operation lease so that ResolveStaleStorageRepair sees it as expired
	old := time.Now().UTC().Add(-2 * SchemaOperationLease)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
	require.NoError(t, err)

	beforeReceipt, err := json.Marshal(f.ddl.states[managedEffectReceiptTable])
	require.NoError(t, err)

	resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
		Inspect: f.ddl.inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
		},
	})
	require.NoError(t, err)
	require.True(t, resolved)

	// Operation transitions to Partial with VerificationMismatch
	op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationPartial, op.Status)
	require.Equal(t, StorageReasonVerificationMismatch, op.Reason)

	// Generation does not advance; repair state is marked blocked
	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, afterLineage.Generation, "generation must not advance")
	require.False(t, afterLineage.HistoryMayBeLost)
	require.Equal(t, StorageRepairStateBlocked, afterLineage.RepairState)
	require.Equal(t, StorageReasonVerificationMismatch, afterLineage.RepairReason)

	require.Empty(t, f.ddl.calls, "no DDL statements executed")

	afterReceipt, err := json.Marshal(f.ddl.states[managedEffectReceiptTable])
	require.NoError(t, err)
	require.Equal(t, beforeReceipt, afterReceipt, "receipt table must remain unchanged")
}

// TestMySQLStorageRestoreStaleRepairTamperedDetailRejected verifies that if an abandoned
// repair operation's detail has tampered generation, to_generation, or scope,
// ResolveStaleStorageRepair returns ErrSchemaOperationMismatch, executes no DDL, leaves
// lineage generation unchanged, and keeps the operation active in the ledger.
func TestMySQLStorageRestoreStaleRepairTamperedDetailRejected(t *testing.T) {
	for name, mutate := range map[string]func(*storageRepairDetail){
		"tampered generation": func(d *storageRepairDetail) {
			d.Generation = 99
		},
		"tampered to_generation": func(d *storageRepairDetail) {
			d.ToGeneration = 99
		},
		"tampered scope": func(d *storageRepairDetail) {
			d.Scope.WorkspaceID = "foreign-workspace"
		},
	} {
		t.Run(name, func(t *testing.T) {
			f := newMySQLRestoreFixture(t)
			f.ddl.states[f.ddl.layout.TableName] = cloneTargetTableInspection(f.ddl.groupExact)

			lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
			require.NoError(t, err)

			repo := f.service.repo.(storageRepairRepository)
			claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
			require.NoError(t, err)
			require.Equal(t, ClaimAcquired, outcome)

			old := time.Now().UTC().Add(-2 * SchemaOperationLease)
			_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
			require.NoError(t, err)

			var detail storageRepairDetail
			require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
			mutate(&detail)
			mutatedBytes, err := json.Marshal(detail)
			require.NoError(t, err)

			_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(mutatedBytes), claim.OperationID)
			require.NoError(t, err)

			_, err = f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
				Inspect: f.ddl.inspect,
				Execute: func(context.Context, []string) (SchemaExecution, error) {
					return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
				},
			})
			require.Error(t, err)
			require.ErrorIs(t, err, ErrSchemaOperationMismatch)

			require.Empty(t, f.ddl.calls, "no DDL statements executed")

			afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
			require.NoError(t, err)
			require.Equal(t, 1, afterLineage.Generation, "generation must remain unchanged")

			op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
			require.NoError(t, err)
			require.Equal(t, SchemaOperationRunning, op.Status, "operation must remain active")
		})
	}
}

// TestMySQLStorageRestoreHealthyCanonicalTableReturnsNotMissing verifies that when the catalog
// shows the canonical table already healthy and matching confirmed lineage, RestoreManagedStorage
// returns StorageNotMissing with zero statements executed, zero reconstructions, and HistoryMayBeLost false.
func TestMySQLStorageRestoreHealthyCanonicalTableReturnsNotMissing(t *testing.T) {
	f := newMySQLRestoreFixture(t)
	f.ddl.states[f.ddl.layout.TableName] = cloneTargetTableInspection(f.ddl.groupExact)

	result, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageNotMissing, result.Outcome)
	require.Empty(t, f.ddl.calls, "no DDL statements executed")

	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, lineage.Generation)
	require.False(t, lineage.HistoryMayBeLost)
	require.Zero(t, lineage.Reconstructions)
	require.Empty(t, lineage.RepairState)
}
