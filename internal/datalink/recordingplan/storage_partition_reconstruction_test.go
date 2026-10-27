package recordingplan

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func nativeStorageRestoreRequest(t *testing.T, f *mysqlPartitionMaintenanceFixture) StorageRestoreRequest {
	t.Helper()
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	return StorageRestoreRequest{
		Scope: f.request.Scope, Dialect: dialectMySQL, ConnectorRevision: lineage.ConnectorRevision,
		TableName: lineage.TableName, OwnerColumn: lineage.OwnerColumn, GroupID: f.request.GroupID,
	}
}

func TestMySQLStorageLineageRefusalNativeMonthlyScope(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	req := nativeStorageRestoreRequest(t, f)
	require.Nil(t, mysqlStorageLineageRefusal(lineage, req), "a confirmed native scope must be accepted")

	wrongScope := req
	wrongScope.Scope.Database, wrongScope.Scope.TableSchema = "other", "other"
	require.NotNil(t, mysqlStorageLineageRefusal(lineage, wrongScope))

	wrongRevision := req
	wrongRevision.ConnectorRevision = "other-revision"
	require.NotNil(t, mysqlStorageLineageRefusal(lineage, wrongRevision))

}

func expandPartitionSidecar(t *testing.T, f *mysqlPartitionMaintenanceFixture) (*MySQLPartitionMaintenanceState, *MySQLPartitionEvidence) {
	t.Helper()
	state, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.NotNil(t, state)
	expanded, err := PlanMySQLNativeMonthlyExtension(state.Current, f.request.At)
	require.NoError(t, err)
	require.NotEqual(t, state.Current, expanded)
	next := *state
	next.Current, next.Revision, next.UpdatedAt = cloneMySQLPartitionEvidence(expanded), state.Revision+1, time.Now().UTC()
	next.Status, next.Reason, next.LastOperationID = mysqlPartitionMaintenanceReady, "", "maintenance-1"
	require.NoError(t, updateMySQLPartitionMaintenanceState(t.Context(), f.db, state.Revision, &next))
	read, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, expanded, read.Current)
	return read, expanded
}

func storageEventCount(t *testing.T, f *mysqlPartitionMaintenanceFixture) int {
	t.Helper()
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(),
		numberedPlaceholders(`SELECT COUNT(*) FROM write_group_storage_events WHERE `+storageScopeWhere),
		scopeArgs(f.request.Scope)...).Scan(&count))
	return count
}

func TestMySQLNativeRepairResetsPartitionSidecar(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	require.Len(t, f.base.Partitions, mysqlPartitionBoundaryCount, "fixture starts with the immutable six-partition layout")
	before, _ := expandPartitionSidecar(t, f)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	layout, layoutDigest, confirmedDigest := lineage.Layout, lineage.LayoutDigest, lineage.ConfirmedDigest
	receipt := cloneTargetTableInspection(f.target.receipt)
	eventsBefore := storageEventCount(t, f)

	changed, err := RecordStorageReconstruction(t.Context(), f.db, f.request.Scope, lineage.Generation, "repair-native", f.request.GroupID, time.Now().UTC())
	require.NoError(t, err)
	require.True(t, changed)

	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	afterState, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, before.Generation+1, afterLineage.Generation)
	require.Equal(t, f.base, afterState.Current, "verified repair resets the expanded sidecar to the immutable initial layout")
	require.Equal(t, afterLineage.Generation, afterState.Generation)
	require.Equal(t, layout, afterLineage.Layout)
	require.Equal(t, layoutDigest, afterLineage.LayoutDigest)
	require.Equal(t, confirmedDigest, afterLineage.ConfirmedDigest)
	require.Equal(t, receipt, f.target.receipt)
	require.Equal(t, eventsBefore+1, storageEventCount(t, f))

	changed, err = RecordStorageReconstruction(t.Context(), f.db, f.request.Scope, lineage.Generation, "repair-native", f.request.GroupID, time.Now().UTC())
	require.NoError(t, err)
	require.False(t, changed, "a repeated completion does not advance the generation again")
	again, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, afterState.Generation, again.Generation)
	require.Equal(t, f.base, again.Current)
	require.Equal(t, eventsBefore+1, storageEventCount(t, f), "a repeated completion does not duplicate events")
}

func TestMySQLNativeRepairSidecarFailureIsAtomic(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	state, _ := expandPartitionSidecar(t, f)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	staleGeneration, staleRevision := lineage.Generation+1, state.Revision+1
	_, err = f.db.ExecContext(t.Context(), numberedPlaceholders(`UPDATE `+mysqlPartitionMaintenanceStateTable+`
		SET generation = $7, revision = $8, updated_at = $9 WHERE `+storageScopeWhere),
		append(scopeArgs(f.request.Scope), staleGeneration, staleRevision, time.Now().UTC().Format(time.RFC3339Nano))...)
	require.NoError(t, err)
	sidecarBefore, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	eventsBefore, reconstructionsBefore, generationBefore := storageEventCount(t, f), lineage.Reconstructions, lineage.Generation

	tx, err := f.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	changed, err := RecordStorageReconstruction(t.Context(), tx, f.request.Scope, lineage.Generation, "repair-atomic", f.request.GroupID, time.Now().UTC())
	require.Error(t, err, "a stale sidecar generation must fail closed")
	require.False(t, changed)
	require.NoError(t, tx.Rollback())

	after, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	afterSidecar, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, generationBefore, after.Generation)
	require.Equal(t, reconstructionsBefore, after.Reconstructions)
	require.Equal(t, eventsBefore, storageEventCount(t, f))
	require.Equal(t, sidecarBefore.Generation, afterSidecar.Generation)
	require.Equal(t, sidecarBefore.Current, afterSidecar.Current)
}

func TestMySQLStorageRestorePolicyRefusalUsesPublicEngine(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	req := nativeStorageRestoreRequest(t, f)
	req.Deny = &StorageDenial{Code: StorageCodeRepairBlocked, Reason: StorageReasonRecoveryDisabled}
	executed := 0
	groupMissing := cloneTargetTableInspection(f.target.group)
	groupMissing.Status = inspectionMissing
	groupMissing.MySQL.Partitions = nil
	receiptExact := cloneTargetTableInspection(f.target.receipt)
	receiptExact.MySQL.Partitions = nil
	target := StorageRestoreTarget{
		Inspect: func(_ context.Context, table string) (TargetTableInspection, error) {
			if table == f.target.layout.TableName {
				return groupMissing, nil
			}
			if table == managedEffectReceiptTable {
				return receiptExact, nil
			}
			return TargetTableInspection{}, errors.New("unexpected table")
		},
		Execute: func(_ context.Context, _ []string) (SchemaExecution, error) {
			executed++
			return SchemaExecution{}, errors.New("execute must not run")
		},
	}
	result, err := f.service.RestoreManagedStorage(t.Context(), req, target)
	require.NoError(t, err)
	require.Equal(t, StorageBlocked, result.Outcome)
	require.Equal(t, StorageCodeRepairBlocked, result.Code)
	require.Equal(t, StorageReasonRecoveryDisabled, result.Reason)
	require.Zero(t, executed)
}

func TestMySQLNativeReconfirmationResetsPartitionSidecar(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	expandPartitionSidecar(t, f)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	confirmation := StorageConfirmation{
		Layout: *lineage.Layout, ConnectorRevision: lineage.ConnectorRevision,
		OperationID: "reconfirm-created", VerifiedDigest: strings.Repeat("a", 64), CreatedTable: true, GroupID: f.request.GroupID,
	}
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC()))
	after, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	state, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, lineage.Generation+1, after.Generation)
	require.Equal(t, f.base, state.Current)
	require.Equal(t, after.Generation, state.Generation)

	expandPartitionSidecar(t, f)
	confirmation.OperationID, confirmation.CreatedTable = "reconfirm-same-table", false
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC()))
	unchangedLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	unchangedSidecar, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, after.Generation, unchangedLineage.Generation)
	require.NotEqual(t, f.base, unchangedSidecar.Current, "ordinary same-table reconfirmation does not reset maintenance state")
}
