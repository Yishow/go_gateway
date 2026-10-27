package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const sqlServerPartitionBlockedCode = "partition_maintenance_blocked"

// TestSQLServerCircuitPartitionIntegration is the original 3.3 entrypoint:
// real sys catalog/placement, boundary and late rows, concurrent maintenance,
// missing owned object repair and foreign ownership refusal on the disposable
// SQL Server fixture. The shared builder skips NOT_RUN when the manifest
// is missing or the managed-long gate is closed, before any write; otherwise
// this body delivers real rows, reads the real catalog and runs maintenance
// and repair through the production services.
func TestSQLServerCircuitPartitionIntegration(t *testing.T) {
	fixture := SQLServerCircuitFixtureOrSkip(t)
	f := NewSQLServerCircuitPlainFixture(t, fixture, SQLServerCircuitFixtureOptions{
		StorageMode: "native_monthly", TagType: "uint64",
	})
	ctx := t.Context()
	instantID, cumulativeID, _ := f.GroupIDs()
	saved, err := f.groups.Get(ctx, cumulativeID)
	require.NoError(t, err)
	require.Equal(t, "native_monthly", string(saved.Group.StorageMode))
	dest := glmDestination(saved.Group)

	targetService := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	inspector := glmSQLServerInspector(t, f)
	capability, err := targetService.ProbeSQLServerRecordingAtRevision(ctx, dest.ConnectorID, dest.ConnectorRevision)
	require.NoError(t, err)
	require.NoError(t, fixture.ValidateCapability(capability))

	// Catalog/proof: the native scope exposes exclusive aligned PF/PS on
	// PRIMARY with the canonical initial M0..M4 boundaries.
	managed := inspectGLMManagedTable(ctx, t, inspector, dest, capability.CapabilityRevision, dest.TableName)
	require.NotNil(t, managed.Evidence.NativeMonthly, "native scope must carry live partition evidence")
	require.True(t, managed.Evidence.NativeMonthly.Exclusive, "PF/PS must be exclusively owned")
	require.True(t, managed.Evidence.NativeMonthly.IndexesAligned)
	require.Equal(t, "PRIMARY", managed.Evidence.NativeMonthly.Filegroup)
	require.NoError(t, sqlServerGLMExpectedInitialEvidence(glmPartitionFromTarget(glmTargetPartition(managed.Evidence.NativeMonthly))))

	// Adapter placement case: use the production SQL writer with the frozen
	// proof for the late past-catch-all row. Current rows still come from the
	// normal source profile pipeline.
	payload, err := f.store.RuntimeVersion(ctx, groupdelivery.GroupKey{
		WorkspaceID:   saved.Group.WorkspaceID,
		GroupID:       saved.Group.ID,
		GroupRevision: saved.Group.Revision,
	})
	require.NoError(t, err)
	var snapshot workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(payload, &snapshot))
	require.NotNil(t, snapshot.RuntimeLayout)
	require.NotNil(t, snapshot.RuntimeLayout.SQLServer)
	initialNativeBytes, err := json.Marshal(snapshot.RuntimeLayout.SQLServer.Group.Evidence.NativeMonthly)
	require.NoError(t, err)
	beforeReceipts := glmReceiptDigests(t, f.target, dest.TableSchema)
	adapterKey, err := insertGLMAdapterPlacement(t, f, &snapshot, saved.Group, f.base.Add(-30*24*time.Hour+12*time.Hour))
	require.NoError(t, err)
	f.accept(t, instantID, sqlServerCircuitMeasurementID, sqlServerCircuitInput{
		sampleID: "sql-partition-current-1", offset: time.Second, value: 5,
	})
	f.accept(t, instantID, sqlServerCircuitMeasurementID, sqlServerCircuitInput{
		sampleID: "sql-partition-current-2", offset: 16 * time.Second, value: 6,
	})
	f.closeAt(t, 30*time.Second)
	waitGLMOutboxSettled(t, f, 2)

	placement := queryGLMPlacement(t, f.target, dest.TableSchema, dest.TableName, "pf_"+dest.TableName)
	require.Equal(t, map[int]int64{1: 1, 2: 2}, placement, "the adapter late row lands in the past catch-all, current rows in M0")
	receipts := glmReceiptDigests(t, f.target, dest.TableSchema)
	for key, digest := range beforeReceipts {
		require.Equal(t, digest, receipts[key], "pre-existing shared receipt must be preserved")
	}
	require.Contains(t, receipts, adapterKey)
	require.Len(t, receipts, len(beforeReceipts)+3)

	// Concurrent maintenance over the scoped lease/ledger: use real probed capability.
	require.True(t, capability.SupportsManagedLong, "production capability must support managed long")
	require.True(t, capability.SupportsNativePartition, "production capability must support native partition")
	require.True(t, capability.SupportsMissingRepair, "production capability must support missing repair")
	scope := glmNativeScope(saved.Group)
	maintenance := recordingplan.NewService(recordingplan.NewSQLRepository(f.db))
	maintenanceTarget := glmMaintenanceTarget(t, inspector, dest, capability.CapabilityRevision, f.target, capability)

	var wg sync.WaitGroup
	results := make([]recordingplan.SQLServerPartitionMaintenanceResult, 2)
	failures := make([]error, 2)
	for slot := range results {
		wg.Go(func() {
			results[slot], failures[slot] = maintenance.MaintainSQLServerPartitions(ctx, glmMaintenanceRequest(scope, saved.Group, dest.ConnectorRevision, f.base.AddDate(0, 1, 0)), maintenanceTarget)
		})
	}
	wg.Wait()
	successes, blockedCount := 0, 0
	for slot := range results {
		require.NoError(t, failures[slot])
		switch results[slot].Code {
		case "":
			successes++
		case sqlServerPartitionBlockedCode:
			blockedCount++
		default:
			require.FailNow(t, "unexpected maintenance result", results[slot].Code)
		}
	}
	require.Equal(t, 2, successes+blockedCount)
	require.GreaterOrEqual(t, successes, 1, "the scoped lease must let a racing caller extend without corruption")
	require.Equal(t, receipts, glmReceiptDigests(t, f.target, dest.TableSchema), "maintenance never touches receipts")
	require.Equal(t, placement, queryGLMPlacement(t, f.target, dest.TableSchema, dest.TableName, "pf_"+dest.TableName),
		"SPLIT is future-append-only: delivered rows keep their placement")

	// Missing owned group table: restore from the immutable INITIAL snapshot;
	// the authorized CURRENT projection remains the repair target.
	initialGroup := snapshot.RuntimeLayout.SQLServer.Group
	initialReceipt := snapshot.RuntimeLayout.SQLServer.Receipt
	// A scope_busy concurrent maintenance call reads state before checking the active lease,
	// which may yield a stale pre-split snapshot. Read the authoritative committed state from DB.
	state, err := recordingplan.ReadSQLServerPartitionMaintenanceState(ctx, f.db, scope)
	require.NoError(t, err)
	require.NotNil(t, state, "maintenance state must exist in database")
	require.NotNil(t, state.Current, "maintenance must expose its authorized current projection")
	authorizedCurrent := recordingplan.CloneSQLServerPartitionEvidence(state.Current)
	lineage, err := recordingplan.ReadStorageLineage(ctx, f.db, scope)
	require.NoError(t, err)
	require.NotNil(t, lineage)
	require.True(t, f.createdTable)
	frozenGroup := inspectGLMManagedTable(ctx, t, inspector, dest, capability.CapabilityRevision, dest.TableName)
	require.Equal(t, initialGroup.Evidence.ObjectID, frozenGroup.Evidence.ObjectID, "repair must target the test-created owned object")
	require.True(t, sqlServerGroupProbeOwnerDefault(frozenGroup.Evidence.ColumnDefaults), "the actual table owner marker must remain gateway")
	_, err = f.target.ExecContext(ctx, "DROP TABLE "+sqlServerQualified(dest.TableSchema, dest.TableName))
	require.NoError(t, err)
	restore, err := maintenance.RestoreManagedStorage(ctx, glmRestoreRequest(scope, saved.Group, dest.ConnectorRevision,
		glmConvertInspection(initialGroup), glmConvertInspection(initialReceipt), lineage.Layout),
		glmRestoreTarget(t, inspector, dest, capability.CapabilityRevision, f.target, capability))
	require.NoError(t, err)
	require.Equal(t, recordingplan.StorageRestored, restore.Outcome, "owned missing table must restore: %s/%s", restore.Code, restore.Reason)
	restored := glmConvertInspection(inspectGLMManagedTable(ctx, t, inspector, dest, capability.CapabilityRevision, dest.TableName))
	require.Equal(t, string(dbtarget.TableInspectionExists), restored.Status)
	require.Equal(t, authorizedCurrent.Partitions, glmPartitionFromTarget(restored.SQLServer.NativeMonthly).Partitions,
		"repair preserves the authorized current expanded boundaries")
	currentStoredPayload, err := f.store.RuntimeVersion(ctx, groupdelivery.GroupKey{
		WorkspaceID: saved.Group.WorkspaceID, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision,
	})
	require.NoError(t, err)
	var currentStoredSnapshot workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(currentStoredPayload, &currentStoredSnapshot))
	afterNativeBytes, err := json.Marshal(currentStoredSnapshot.RuntimeLayout.SQLServer.Group.Evidence.NativeMonthly)
	require.NoError(t, err)
	require.Equal(t, string(initialNativeBytes), string(afterNativeBytes),
		"the immutable initial proof in store must not be mutated by maintenance or restore")
	require.Len(t, currentStoredSnapshot.RuntimeLayout.SQLServer.Group.Evidence.NativeMonthly.Partitions, 6,
		"the immutable initial proof stays six boundaries")
	require.Equal(t, receipts, glmReceiptDigests(t, f.target, dest.TableSchema), "receipts are preserved verbatim")

	// Original missing-object combinations are executable, using the same
	// production restore path. Each drop is preceded by this fixture's creation
	// and object proof, and only this test-owned native object set is touched.
	for _, missing := range []struct {
		name          string
		partScheme    bool
		partitionFunc bool
	}{
		{name: "table_and_scheme", partScheme: true},
		{name: "all_three", partScheme: true, partitionFunc: true},
	} {
		t.Run(missing.name, func(t *testing.T) {
			require.True(t, f.createdTable)
			current := inspectGLMManagedTable(ctx, t, inspector, dest, capability.CapabilityRevision, dest.TableName)
			require.NotZero(t, current.Evidence.ObjectID)
			require.True(t, sqlServerGroupProbeOwnerDefault(current.Evidence.ColumnDefaults), "the actual table owner marker must remain gateway")
			_, err := f.target.ExecContext(ctx, "IF OBJECT_ID(@p1,N'U')=@p2 DROP TABLE "+sqlServerQualified(dest.TableSchema, dest.TableName),
				dest.TableSchema+"."+dest.TableName, current.Evidence.ObjectID)
			require.NoError(t, err)
			psName := "ps_" + dest.TableName
			pfName := "pf_" + dest.TableName
			if missing.partScheme {
				var psObjectID int64
				require.NoError(t, f.target.QueryRowContext(ctx, "SELECT data_space_id FROM sys.partition_schemes WHERE name = @p1", psName).Scan(&psObjectID))
				require.NotZero(t, psObjectID)
				_, err = f.target.ExecContext(ctx, "IF EXISTS (SELECT 1 FROM sys.partition_schemes WHERE name=@p1 AND data_space_id=@p2) DROP PARTITION SCHEME ["+strings.ReplaceAll(psName, "]", "]]")+"]", psName, psObjectID)
				require.NoError(t, err)
			}
			if missing.partitionFunc {
				var pfObjectID int64
				require.NoError(t, f.target.QueryRowContext(ctx, "SELECT function_id FROM sys.partition_functions WHERE name = @p1", pfName).Scan(&pfObjectID))
				require.NotZero(t, pfObjectID)
				_, err = f.target.ExecContext(ctx, "IF EXISTS (SELECT 1 FROM sys.partition_functions WHERE name=@p1 AND function_id=@p2) DROP PARTITION FUNCTION ["+strings.ReplaceAll(pfName, "]", "]]")+"]", pfName, pfObjectID)
				require.NoError(t, err)
			}
			restore, err := maintenance.RestoreManagedStorage(ctx, glmRestoreRequest(scope, saved.Group, dest.ConnectorRevision,
				glmConvertInspection(initialGroup), glmConvertInspection(initialReceipt), lineage.Layout),
				glmRestoreTarget(t, inspector, dest, capability.CapabilityRevision, f.target, capability))
			require.NoError(t, err)
			require.Equal(t, recordingplan.StorageRestored, restore.Outcome, "owned missing %s must restore: %s/%s", missing.name, restore.Code, restore.Reason)
			restored := glmConvertInspection(inspectGLMManagedTable(ctx, t, inspector, dest, capability.CapabilityRevision, dest.TableName))
			require.Equal(t, string(dbtarget.TableInspectionExists), restored.Status)
		})
	}

	// Foreign ownership refusal: target the same PF object, then restore the
	// original marker in cleanup so the fixture is not left foreign-owned.
	pfName := "pf_" + dest.TableName
	var originalMarker sql.NullString
	_ = f.target.QueryRowContext(ctx, "SELECT CAST(value AS nvarchar(4000)) FROM sys.extended_properties WHERE class=21 AND major_id=(SELECT function_id FROM sys.partition_functions WHERE name=@p1) AND name=N'gateway.recording.partition_owner'", pfName).Scan(&originalMarker)
	var originalPFObjectID int64
	require.NoError(t, f.target.QueryRowContext(ctx, "SELECT function_id FROM sys.partition_functions WHERE name=@p1", pfName).Scan(&originalPFObjectID))
	require.NotZero(t, originalPFObjectID)
	t.Cleanup(func() {
		if originalMarker.Valid {
			_, restoreErr := f.target.ExecContext(context.Background(), `EXEC sys.sp_updateextendedproperty @name=N'gateway.recording.partition_owner', @value=@p1, @level0type=N'PARTITION FUNCTION', @level0name=@p2`, originalMarker.String, pfName)
			require.NoError(t, restoreErr)
		} else {
			_, dropErr := f.target.ExecContext(context.Background(), `EXEC sys.sp_dropextendedproperty @name=N'gateway.recording.partition_owner', @level0type=N'PARTITION FUNCTION', @level0name=@p1`, pfName)
			require.NoError(t, dropErr)
		}
	})
	_, err = f.target.ExecContext(ctx, `EXEC sys.sp_updateextendedproperty @name=N'gateway.recording.partition_owner', @value=N'dbo.foreign', @level0type=N'PARTITION FUNCTION', @level0name=N'pf_`+dest.TableName+`'`)
	require.NoError(t, err)
	var afterPFObjectID int64
	require.NoError(t, f.target.QueryRowContext(ctx, "SELECT function_id FROM sys.partition_functions WHERE name=@p1", pfName).Scan(&afterPFObjectID))
	require.Equal(t, originalPFObjectID, afterPFObjectID, "foreign marker fault must target the same PF object")
	blocked, err := maintenance.MaintainSQLServerPartitions(ctx, glmMaintenanceRequest(scope, saved.Group, dest.ConnectorRevision, f.base.AddDate(0, 2, 0)), maintenanceTarget)
	require.NoError(t, err)
	require.Equal(t, sqlServerPartitionBlockedCode, blocked.Code, "a foreign partition owner must block the scope")
	require.Equal(t, receipts, glmReceiptDigests(t, f.target, dest.TableSchema))

}
