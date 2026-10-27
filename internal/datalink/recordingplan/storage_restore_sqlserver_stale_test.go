package recordingplan

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLServerStaleRecovery_StartedAdvancesGeneration(t *testing.T) {
	f := newSQLServerRestoreFixture(t)

	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	repo := f.service.repo.(storageRepairRepository)
	claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)

	// Recreated table present with new ObjectID 201 and bound to this repair operation ID
	recreated := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 201, f.options)
	recreated.SQLServer.RepairOperationID = claim.OperationID
	f.states[f.layout.TableName] = recreated

	// Persist SQLServerStarted = true on claim
	var detail storageRepairDetail
	require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
	detail.SQLServerStarted = true
	encodedDetail, err := json.Marshal(detail)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim.OperationID)
	require.NoError(t, err)

	// Age the operation lease
	old := time.Now().UTC().Add(-2 * SchemaOperationLease)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
	require.NoError(t, err)

	resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
		Inspect: f.inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
		},
	})
	require.NoError(t, err)
	require.True(t, resolved)

	// Operation transitions to Succeeded
	op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationSucceeded, op.Status)
	require.Equal(t, SchemaReasonRecovered, op.Reason)

	// Generation advances to 2
	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, afterLineage.Generation)
	require.True(t, afterLineage.HistoryMayBeLost)
	require.Equal(t, claim.OperationID, afterLineage.LastReconstructionOperationID)
}

func TestSQLServerStaleRecovery_NotStartedBlocks(t *testing.T) {
	f := newSQLServerRestoreFixture(t)

	// Table present with new ObjectID, but DDL was never started
	recreated := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 201, f.options)
	f.states[f.layout.TableName] = recreated

	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	repo := f.service.repo.(storageRepairRepository)
	claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)

	// Age lease without setting SQLServerStarted
	old := time.Now().UTC().Add(-2 * SchemaOperationLease)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
	require.NoError(t, err)

	resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
		Inspect: f.inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
		},
	})
	require.NoError(t, err)
	require.True(t, resolved)

	op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationPartial, op.Status)
	require.Equal(t, StorageReasonVerificationMismatch, op.Reason)

	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, afterLineage.Generation, "generation must not advance")
	require.Equal(t, StorageRepairStateBlocked, afterLineage.RepairState)
}

func TestSQLServerStaleRecovery_CurrentAuthorizedClosesWithoutReconstruction(t *testing.T) {
	t.Run("generation 1 confirmed table", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)

		// Table present matching original confirmed catalog
		f.states[f.layout.TableName] = f.group

		lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
		require.NoError(t, err)

		repo := f.service.repo.(storageRepairRepository)
		claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
		require.NoError(t, err)
		require.Equal(t, ClaimAcquired, outcome)

		var detail storageRepairDetail
		require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
		detail.SQLServerStarted = true
		encodedDetail, err := json.Marshal(detail)
		require.NoError(t, err)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim.OperationID)
		require.NoError(t, err)

		old := time.Now().UTC().Add(-2 * SchemaOperationLease)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
		require.NoError(t, err)

		resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
			Inspect: f.inspect,
			Execute: func(context.Context, []string) (SchemaExecution, error) {
				return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
			},
		})
		require.NoError(t, err)
		require.True(t, resolved)

		// Must close as abandoned without advancing generation or blocking
		op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
		require.NoError(t, err)
		require.Equal(t, SchemaOperationFailed, op.Status)
		require.Equal(t, StorageReasonAttemptAbandoned, op.Reason)

		afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
		require.NoError(t, err)
		require.Equal(t, 1, afterLineage.Generation, "generation must not advance")
		require.Empty(t, afterLineage.RepairState, "must not be blocked")
		require.Zero(t, afterLineage.Reconstructions, "zero reconstructions recorded")
	})

	t.Run("generation 2 authorized table", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)

		// First, perform a legitimate restore to reach generation 2
		res, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageRestored, res.Outcome)

		lineage2, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
		require.NoError(t, err)
		require.Equal(t, 2, lineage2.Generation)

		// Now simulate a stale claim for generation 2 while the table is already authorized
		repo := f.service.repo.(storageRepairRepository)
		claim2, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage2, *lineage2.Layout)
		require.NoError(t, err)
		require.Equal(t, ClaimAcquired, outcome)

		var detail storageRepairDetail
		require.NoError(t, json.Unmarshal([]byte(claim2.Detail), &detail))
		detail.SQLServerStarted = true
		encodedDetail, err := json.Marshal(detail)
		require.NoError(t, err)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim2.OperationID)
		require.NoError(t, err)

		old := time.Now().UTC().Add(-2 * SchemaOperationLease)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim2.OperationID)
		require.NoError(t, err)

		resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
			Inspect: f.inspect,
			Execute: func(context.Context, []string) (SchemaExecution, error) {
				return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
			},
		})
		require.NoError(t, err)
		require.True(t, resolved)

		// Must close without reconstruction at gen 2
		op2, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim2.OperationID)
		require.NoError(t, err)
		require.Equal(t, SchemaOperationFailed, op2.Status)
		require.Equal(t, StorageReasonAttemptAbandoned, op2.Reason)

		afterLineage2, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
		require.NoError(t, err)
		require.Equal(t, 2, afterLineage2.Generation, "generation must stay at 2")
		require.Empty(t, afterLineage2.RepairState)
	})
}

func TestSQLServerStaleRecovery_ForeignReplacementBlocks(t *testing.T) {
	f := newSQLServerRestoreFixture(t)

	// Foreign replacement table with no or wrong repair marker
	foreign := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 999, f.options)
	foreign.SQLServer.RepairOperationID = "foreign-operation-id"
	f.states[f.layout.TableName] = foreign

	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	repo := f.service.repo.(storageRepairRepository)
	claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)

	var detail storageRepairDetail
	require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
	detail.SQLServerStarted = true
	encodedDetail, err := json.Marshal(detail)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim.OperationID)
	require.NoError(t, err)

	old := time.Now().UTC().Add(-2 * SchemaOperationLease)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
	require.NoError(t, err)

	resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
		Inspect: f.inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
		},
	})
	require.NoError(t, err)
	require.True(t, resolved)

	op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationPartial, op.Status)
	require.Equal(t, StorageReasonVerificationMismatch, op.Reason)

	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, afterLineage.Generation)
	require.Equal(t, StorageRepairStateBlocked, afterLineage.RepairState)
}

func TestSQLServerStaleRecovery_CorruptOrMismatchedDetail(t *testing.T) {
	t.Run("nil nested sqlserver", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
		require.NoError(t, err)

		repo := f.service.repo.(storageRepairRepository)
		claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
		require.NoError(t, err)
		require.Equal(t, ClaimAcquired, outcome)

		// Nil out detail.SQLServer.Group.SQLServer
		var detail storageRepairDetail
		require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
		detail.SQLServer.Group.SQLServer = nil
		encodedDetail, err := json.Marshal(detail)
		require.NoError(t, err)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim.OperationID)
		require.NoError(t, err)

		old := time.Now().UTC().Add(-2 * SchemaOperationLease)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
		require.NoError(t, err)

		_, err = f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.inspect})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrSchemaOperationMismatch)
	})

	t.Run("tampered original proof digest mismatch", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
		require.NoError(t, err)

		repo := f.service.repo.(storageRepairRepository)
		claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
		require.NoError(t, err)
		require.Equal(t, ClaimAcquired, outcome)

		// Tamper columns in original proof group
		var detail storageRepairDetail
		require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
		detail.SQLServer.Group.Columns = detail.SQLServer.Group.Columns[:len(detail.SQLServer.Group.Columns)-1]
		encodedDetail, err := json.Marshal(detail)
		require.NoError(t, err)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim.OperationID)
		require.NoError(t, err)

		old := time.Now().UTC().Add(-2 * SchemaOperationLease)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
		require.NoError(t, err)

		_, err = f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.inspect})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrSchemaOperationMismatch)
	})

	t.Run("tampered receipt catalog mismatch", func(t *testing.T) {
		f := newSQLServerRestoreFixture(t)
		lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
		require.NoError(t, err)

		repo := f.service.repo.(storageRepairRepository)
		claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
		require.NoError(t, err)
		require.Equal(t, ClaimAcquired, outcome)

		var detail storageRepairDetail
		require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
		detail.SQLServer.Receipt.SQLServer.ObjectID = 99999
		encodedDetail, err := json.Marshal(detail)
		require.NoError(t, err)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim.OperationID)
		require.NoError(t, err)

		old := time.Now().UTC().Add(-2 * SchemaOperationLease)
		_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
		require.NoError(t, err)

		_, err = f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.inspect})
		require.Error(t, err)
		require.ErrorIs(t, err, ErrSchemaOperationMismatch)
	})
}

func TestSQLServerStaleRecovery_StillMissingAbandons(t *testing.T) {
	f := newSQLServerRestoreFixture(t)

	// Table still missing
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	repo := f.service.repo.(storageRepairRepository)
	claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)

	old := time.Now().UTC().Add(-2 * SchemaOperationLease)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
	require.NoError(t, err)

	resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{
		Inspect: f.inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
		},
	})
	require.NoError(t, err)
	require.True(t, resolved)

	op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationFailed, op.Status)
	require.Equal(t, StorageReasonAttemptAbandoned, op.Reason)

	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, afterLineage.Generation)
	require.Empty(t, afterLineage.RepairState, "abandoned repair leaves state unblocked for retry")
}

func TestSQLServerStaleRecovery_GenerationCASMissRollsBack(t *testing.T) {
	f := newSQLServerRestoreFixture(t)

	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	repo := f.service.repo.(storageRepairRepository)
	claim, outcome, err := f.service.claimStorageRepair(t.Context(), repo, f.request, lineage, *lineage.Layout)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)

	recreated := sqlServerRestoreCatalog(f.layout.Columns, f.layout.Indexes, f.layout.TableName, f.layout.OwnerColumn, 201, f.options)
	recreated.SQLServer.RepairOperationID = claim.OperationID
	f.states[f.layout.TableName] = recreated

	var detail storageRepairDetail
	require.NoError(t, json.Unmarshal([]byte(claim.Detail), &detail))
	detail.SQLServerStarted = true
	encodedDetail, err := json.Marshal(detail)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, string(encodedDetail), claim.OperationID)
	require.NoError(t, err)

	old := time.Now().UTC().Add(-2 * SchemaOperationLease)
	_, err = f.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
	require.NoError(t, err)

	_, err = f.db.ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET generation = 99 WHERE workspace_id = ?`, f.request.Scope.WorkspaceID)
	require.NoError(t, err)

	_, err = f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.inspect})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrSchemaOperationMismatch)

	op, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, claim.OperationID)
	require.NoError(t, err)
	require.NotEqual(t, SchemaOperationSucceeded, op.Status)
}
