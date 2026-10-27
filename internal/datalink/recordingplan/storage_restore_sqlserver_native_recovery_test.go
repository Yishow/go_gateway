package recordingplan

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativeRepairRefusesUntrustedOrStaleEvidence(t *testing.T) {
	for _, test := range []struct {
		name        string
		prepare     func(*testing.T, *sqlServerRestoreFixture)
		wantOutcome StorageRestoreOutcome
		wantCode    string
		wantReason  string
		wantError   error
	}{
		{
			name: "shared objects",
			prepare: func(_ *testing.T, fixture *sqlServerRestoreFixture) {
				setSQLServerNativeMissingGroup(fixture, fixture.layout.SQLServer.NativeMonthly, true, true, false)
			},
			wantOutcome: StorageBlocked, wantCode: StorageCodeOwnerMismatch, wantReason: StorageReasonTableOccupied,
		},
		{
			name: "unconfirmed partition lookup",
			prepare: func(_ *testing.T, fixture *sqlServerRestoreFixture) {
				setSQLServerNativeMissingGroup(fixture, fixture.layout.SQLServer.NativeMonthly, true, true, true)
				fixture.states[fixture.layout.TableName].SQLServer.PartitionInspected = false
			},
			wantError: ErrTargetInspectionUnconfirmed,
		},
		{
			name: "forbidden table",
			prepare: func(_ *testing.T, fixture *sqlServerRestoreFixture) {
				setSQLServerNativeMissingGroup(fixture, fixture.layout.SQLServer.NativeMonthly, true, true, true)
				forbidden := fixture.states[fixture.layout.TableName]
				forbidden.Status = inspectionForbidden
				fixture.states[fixture.layout.TableName] = forbidden
			},
			wantOutcome: StorageBlocked, wantCode: StorageCodeRepairBlocked, wantReason: StorageReasonPermissionDenied,
		},
		{
			name: "policy denied",
			prepare: func(_ *testing.T, fixture *sqlServerRestoreFixture) {
				fixture.request.Deny = &StorageDenial{Code: StorageCodeRepairBlocked, Reason: StorageReasonPolicyUnconfirmed}
			},
			wantOutcome: StorageBlocked, wantCode: StorageCodeRepairBlocked, wantReason: StorageReasonPolicyUnconfirmed,
		},
		{
			name: "connector revision changed",
			prepare: func(_ *testing.T, fixture *sqlServerRestoreFixture) {
				fixture.request.ConnectorRevision = "rev-changed"
			},
			wantOutcome: StorageBlocked, wantCode: StorageCodeRepairBlocked, wantReason: StorageReasonConnectorRevision,
		},
		{
			name: "partition generation changed",
			prepare: func(t *testing.T, fixture *sqlServerRestoreFixture) {
				current := extendSQLServerNativeRestoreFixture(t, fixture)
				setSQLServerNativeMissingGroup(fixture, current, true, true, true)
				_, err := fixture.db.ExecContext(t.Context(), `UPDATE write_group_sqlserver_partition_maintenance SET generation = 9 WHERE workspace_id = ?`, fixture.request.Scope.WorkspaceID)
				require.NoError(t, err)
			},
			wantOutcome: StorageBlocked, wantCode: StorageCodeRepairBlocked, wantReason: StorageReasonVerificationMismatch,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newSQLServerNativeRestoreFixture(t)
			test.prepare(t, fixture)
			result, err := fixture.service.RestoreManagedStorage(t.Context(), fixture.request,
				sqlServerNativeRestoreTarget(fixture, fixture.layout.SQLServer.NativeMonthly, nil))
			if test.wantError != nil {
				require.ErrorIs(t, err, test.wantError)
				require.Empty(t, fixture.calls)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.wantOutcome, result.Outcome)
			require.Equal(t, test.wantCode, result.Code)
			require.Equal(t, test.wantReason, result.Reason)
			require.Empty(t, fixture.calls, "unsafe or stale evidence must not execute repair DDL")
		})
	}
}

func TestSQLServerNativeRestoreRecoversLostDDLAcknowledgement(t *testing.T) {
	fixture := newSQLServerNativeRestoreFixture(t)
	current := CloneSQLServerPartitionEvidence(fixture.layout.SQLServer.NativeMonthly)
	setSQLServerNativeMissingGroup(fixture, current, false, false, true)
	frozenReceipt, err := json.Marshal(fixture.request.SQLServer.Receipt)
	require.NoError(t, err)
	target := StorageRestoreTarget{
		Inspect: fixture.inspect,
		ExecuteSQLServerRepair: func(_ context.Context, _, operationID string, statements []string) (SchemaExecution, error) {
			executed := sqlServerNativeExecutedStatements(statements, nil)
			fixture.calls = append(fixture.calls, executed...)
			recreated := sqlServerRestoreCatalog(fixture.layout.Columns, fixture.layout.Indexes,
				fixture.layout.TableName, fixture.layout.OwnerColumn, 201, fixture.options)
			recreated.SQLServer.Plain = false
			recreated.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(current, true, true, true, true)
			recreated.SQLServer.PartitionInspected = true
			recreated.SQLServer.RepairOperationID = operationID
			fixture.states[fixture.layout.TableName] = recreated
			return SchemaExecution{Committed: len(executed)}, errors.New("simulated commit acknowledgement loss")
		},
	}
	result, err := fixture.service.RestoreManagedStorage(t.Context(), fixture.request, target)
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome)
	require.Equal(t, SchemaReasonRecovered, result.Operation.Reason)
	require.NotEmpty(t, result.Operation.VerifiedDigest)
	require.Equal(t, 2, mustReadStorageLineage(t, fixture).Generation)
	afterReceipt, err := json.Marshal(fixture.request.SQLServer.Receipt)
	require.NoError(t, err)
	require.Equal(t, frozenReceipt, afterReceipt, "lost DDL acknowledgement cannot recreate or mutate the receipt proof")
	require.Equal(t, 1, strings.Count(strings.Join(fixture.calls, "\n"), "CREATE PARTITION FUNCTION"))
}

func TestSQLServerNativeRestoreResolvesStaleCommittedOperationAfterRestart(t *testing.T) {
	fixture := newSQLServerNativeRestoreFixture(t)
	current := CloneSQLServerPartitionEvidence(fixture.layout.SQLServer.NativeMonthly)
	lineage := mustReadStorageLineage(t, fixture)
	repository, ok := fixture.service.repo.(storageRepairRepository)
	require.True(t, ok)
	op, outcome, err := fixture.service.claimStorageRepair(t.Context(), repository, fixture.request, lineage, fixture.layout)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	var detail storageRepairDetail
	require.NoError(t, json.Unmarshal([]byte(op.Detail), &detail))
	detail.SQLServerStarted = true
	progress, ok := fixture.service.repo.(operationExecutionRepository)
	require.True(t, ok)
	encoded, err := json.Marshal(detail)
	require.NoError(t, err)
	require.NoError(t, progress.SaveSchemaOperationProgress(t.Context(), op.OperationID, op.Owner, string(encoded)))
	recreated := sqlServerRestoreCatalog(fixture.layout.Columns, fixture.layout.Indexes,
		fixture.layout.TableName, fixture.layout.OwnerColumn, 201, fixture.options)
	recreated.SQLServer.Plain = false
	recreated.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(current, true, true, true, true)
	recreated.SQLServer.PartitionInspected = true
	recreated.SQLServer.RepairOperationID = op.OperationID
	fixture.states[fixture.layout.TableName] = recreated
	staleAt := time.Now().UTC().Add(-2 * SchemaOperationLease)
	_, err = fixture.db.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`,
		staleAt, op.OperationID)
	require.NoError(t, err)

	var executeCalls int
	target := StorageRestoreTarget{
		Inspect: fixture.inspect,
		ExecuteSQLServerRepair: func(context.Context, string, string, []string) (SchemaExecution, error) {
			executeCalls++
			return SchemaExecution{}, errors.New("stale recovery must not run DDL")
		},
	}
	resolvedStale, err := fixture.service.ResolveStaleStorageRepair(t.Context(), fixture.request, target)
	require.NoError(t, err)
	require.True(t, resolvedStale, "restart resolves an expired repair from its live operation marker")
	lineageAfter := mustReadStorageLineage(t, fixture)
	resolvedOperation, err := fixture.service.repo.GetSchemaOperation(t.Context(), fixture.request.Scope.WorkspaceID, op.OperationID)
	require.NoError(t, err)
	inspected, err := inspectSQLServerNativeRestoreGroupTables(t.Context(), *lineageAfter.Layout, fixture.request.Scope.Database, fixture.inspect)
	require.NoError(t, err)
	currentDigest := schemaTablesDigest(inspected.tables)
	authorizedDigest, err := fixture.service.authorizedSQLServerDigest(t.Context(), fixture.request.Scope, lineageAfter)
	require.NoError(t, err)
	require.Equal(t, resolvedOperation.VerifiedDigest, currentDigest,
		"lineage=%+v operation=%+v inspection=%+v", lineageAfter, resolvedOperation, inspected)
	require.Equal(t, resolvedOperation.VerifiedDigest, authorizedDigest)
	result, err := fixture.service.RestoreManagedStorage(t.Context(), fixture.request, target)
	require.NoError(t, err)
	require.Equal(t, StorageNotMissing, result.Outcome, "recovered catalog is recognized on the next restore: code=%s reason=%s", result.Code, result.Reason)
	require.Zero(t, executeCalls, "restart recovery resolves only from the committed catalog")
	resolved, err := fixture.service.repo.GetSchemaOperation(t.Context(), fixture.request.Scope.WorkspaceID, op.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationSucceeded, resolved.Status)
	require.Equal(t, SchemaReasonRecovered, resolved.Reason)
	after := mustReadStorageLineage(t, fixture)
	require.Equal(t, 2, after.Generation)
	require.Equal(t, op.OperationID, after.LastReconstructionOperationID)
}

func TestSQLServerNativeRestoreSerializesConcurrentSameScopeClaims(t *testing.T) {
	fixture := newSQLServerNativeRestoreFixture(t)
	current := CloneSQLServerPartitionEvidence(fixture.layout.SQLServer.NativeMonthly)
	setSQLServerNativeMissingGroup(fixture, current, false, false, true)
	var stateMu sync.RWMutex
	started := make(chan struct{}, 1)
	releaseChannel := make(chan struct{})
	release := sync.OnceFunc(func() { close(releaseChannel) })
	defer release()
	executions := make(chan struct{}, 2)
	target := StorageRestoreTarget{
		Inspect: func(ctx context.Context, table string) (TargetTableInspection, error) {
			stateMu.RLock()
			defer stateMu.RUnlock()
			return fixture.inspect(ctx, table)
		},
		ExecuteSQLServerRepair: func(_ context.Context, _, operationID string, statements []string) (SchemaExecution, error) {
			executions <- struct{}{}
			select {
			case started <- struct{}{}:
			default:
			}
			<-releaseChannel
			executed := sqlServerNativeExecutedStatements(statements, nil)
			stateMu.Lock()
			defer stateMu.Unlock()
			fixture.calls = append(fixture.calls, executed...)
			recreated := sqlServerRestoreCatalog(fixture.layout.Columns, fixture.layout.Indexes,
				fixture.layout.TableName, fixture.layout.OwnerColumn, 201, fixture.options)
			recreated.SQLServer.Plain = false
			recreated.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(current, true, true, true, true)
			recreated.SQLServer.PartitionInspected = true
			recreated.SQLServer.RepairOperationID = operationID
			fixture.states[fixture.layout.TableName] = recreated
			return SchemaExecution{Committed: len(executed)}, nil
		},
	}
	type restoreResult struct {
		result StorageRestoreResult
		err    error
	}
	firstDone := make(chan restoreResult, 1)
	go func() {
		result, err := fixture.service.RestoreManagedStorage(t.Context(), fixture.request, target)
		firstDone <- restoreResult{result: result, err: err}
	}()
	select {
	case <-started:
	case <-time.After(10 * time.Second):
		t.Fatal("first same-scope repair did not reach the target")
	}
	secondDone := make(chan restoreResult, 1)
	go func() {
		result, err := fixture.service.RestoreManagedStorage(t.Context(), fixture.request, target)
		secondDone <- restoreResult{result: result, err: err}
	}()
	var second restoreResult
	select {
	case second = <-secondDone:
	case <-time.After(10 * time.Second):
		release()
		t.Fatal("competing same-scope repair did not resolve to the ledger claim")
	}
	require.NoError(t, second.err)
	require.Equal(t, StorageInProgress, second.result.Outcome)
	release()
	first := <-firstDone
	require.NoError(t, first.err)
	require.Equal(t, StorageRestored, first.result.Outcome)
	require.Len(t, executions, 1, "only one same-scope claimant may execute repair DDL")
}

func mustReadStorageLineage(t *testing.T, fixture *sqlServerRestoreFixture) *StorageLineage {
	t.Helper()
	lineage, err := ReadStorageLineage(t.Context(), fixture.db, fixture.request.Scope)
	require.NoError(t, err)
	return lineage
}
