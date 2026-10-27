package recordingplan

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newSQLServerCapabilityProjection(t *testing.T, f *sqlServerPartitionMaintenanceFixture, capRev string) (resultLayout GroupSchemaLayout, resultGroup, resultReceipt TargetTableInspection, resultDigest string) {
	t.Helper()
	options := *f.layout.SQLServer
	options.CapabilityRevision = capRev
	layout, err := ManagedSQLServerLongGroupLayout(f.layout.TableName, f.layout.OwnerColumn, options)
	require.NoError(t, err)

	group := sqlServerRestoreCatalog(layout.Columns, layout.Indexes, layout.TableName, layout.OwnerColumn, 401, options)
	group.SQLServer.NativeMonthly = sqlServerTargetPartitionFromPlan(f.base)
	group.SQLServer.PartitionInspected = true
	group.SQLServer.Plain = false
	receipt := sqlServerRestoreReceiptCatalog(402, options)

	digest, err := VerifyManagedSQLServerCatalog(layout, f.request.Scope.Database, group, receipt)
	require.NoError(t, err)
	return layout, group, receipt, digest
}

func TestSQLServerNativeReconfirmationRebindsSidecarWithoutResettingCurrent(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)

	// Step 1: Run initial maintenance to expand partitions beyond the base
	initialResult, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, initialResult.Changed)

	stateBeforeReconfirm, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.NotNil(t, stateBeforeReconfirm)
	expanded := CloneSQLServerPartitionEvidence(stateBeforeReconfirm.Current)
	require.True(t, sameSQLServerPartitionPrefix(f.base, expanded))
	require.NotEqual(t, f.base, expanded, "maintenance must have expanded partitions beyond initial base")

	// Step 2: Reconfirm with new connector revision and capability projection (CreatedTable: false)
	newLayout, newGroup, newReceipt, verifiedDigest := newSQLServerCapabilityProjection(t, f, "cap-2")
	confirmation := StorageConfirmation{
		Layout:            newLayout,
		ConnectorRevision: "rev-2",
		OperationID:       "op-reconfirm-rev-2",
		VerifiedDigest:    verifiedDigest,
		CreatedTable:      false,
		GroupID:           f.request.GroupID,
	}
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC()))

	// Step 3: Check that sidecar was rebound to new layout digest and connector revision, keeping expanded current boundaries
	afterState, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, afterState.Generation)
	require.Equal(t, "rev-2", afterState.ConnectorRevision)
	require.Equal(t, digestGroupSchemaLayout(newLayout), afterState.BaseLayoutDigest)
	require.Equal(t, expanded, afterState.Current, "same-table reconfirmation must preserve authorized current partition projection")
	require.Equal(t, sqlServerPartitionMaintenanceReady, afterState.Status)

	// Step 4: Verify subsequent MaintainSQLServerPartitions service call succeeds with rev-2
	f.request.ConnectorRevision = "rev-2"
	f.request.At = time.Date(2027, time.January, 20, 0, 0, 0, 0, time.UTC)
	f.target.layout = newLayout
	f.target.group = newGroup
	f.target.receipt = newReceipt
	mResult, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, mResult.Changed)
}

func TestSQLServerNativeRecreateResetsInitialShapeOnce(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)

	// Step 1: Run initial maintenance to expand partitions beyond the base
	initialResult, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, initialResult.Changed)

	stateBeforeReconfirm, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	expanded := CloneSQLServerPartitionEvidence(stateBeforeReconfirm.Current)
	require.NotEqual(t, f.base, expanded)

	// Step 2: Reconfirm with CreatedTable: true (generation advance from 1 to 2)
	newLayout, newGroup, newReceipt, verifiedDigest := newSQLServerCapabilityProjection(t, f, "cap-2")
	confirmation := StorageConfirmation{
		Layout:            newLayout,
		ConnectorRevision: "rev-2",
		OperationID:       "op-recreate-rev-2",
		VerifiedDigest:    verifiedDigest,
		CreatedTable:      true,
		GroupID:           f.request.GroupID,
	}
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC()))

	// Step 3: Check lineage generation advanced to 2 and sidecar was reset to generation 2 with base partition evidence
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, lineage.Generation)

	afterState, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, afterState.Generation)
	require.Equal(t, "rev-2", afterState.ConnectorRevision)
	require.Equal(t, digestGroupSchemaLayout(newLayout), afterState.BaseLayoutDigest)
	require.Equal(t, f.base, afterState.Current, "generation recreation must reset to base partition evidence")
	require.Equal(t, sqlServerPartitionMaintenanceReady, afterState.Status)

	// Repeat of same confirmation operation is idempotent
	repeated := confirmation
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, repeated, time.Now().UTC()))
	repeatState, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, repeatState.Generation)
	require.Equal(t, f.base, repeatState.Current)

	// Step 4: Verify subsequent MaintainSQLServerPartitions service call succeeds in generation 2
	f.request.ConnectorRevision = "rev-2"
	f.request.At = time.Date(2027, time.January, 20, 0, 0, 0, 0, time.UTC)
	f.target.layout = newLayout
	f.target.group = newGroup
	f.target.receipt = newReceipt
	f.target.partitions = CloneSQLServerPartitionEvidence(f.base)
	mResult, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, mResult.Changed)
}

func TestSQLServerNativeReconfirmationRejectsStaleOrForeignSidecar(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)

	// Step 1: Run initial maintenance
	_, err := f.maintain(t)
	require.NoError(t, err)

	beforeLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	// Corrupt sidecar's base_layout_digest
	_, err = f.db.ExecContext(t.Context(), numberedPlaceholders(`UPDATE `+sqlServerPartitionMaintenanceStateTable+`
		SET base_layout_digest = $7 WHERE `+storageScopeWhere), append(scopeArgs(f.request.Scope), strings.Repeat("0", 64))...)
	require.NoError(t, err)

	staleState, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	newLayout, _, _, verifiedDigest := newSQLServerCapabilityProjection(t, f, "cap-2")
	confirmation := StorageConfirmation{
		Layout:            newLayout,
		ConnectorRevision: "rev-2",
		OperationID:       "op-rollback",
		VerifiedDigest:    verifiedDigest,
		CreatedTable:      false,
		GroupID:           f.request.GroupID,
	}

	// Step 2: RecordStorageConfirmation should fail due to state mismatch
	err = RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC())
	require.Error(t, err)
	require.ErrorIs(t, err, ErrSQLServerPartitionMaintenanceStateMismatch)

	// Step 3: Verify atomic rollback - storage scope and sidecar untouched
	afterLineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, beforeLineage.LayoutDigest, afterLineage.LayoutDigest)
	require.Equal(t, beforeLineage.ConnectorRevision, afterLineage.ConnectorRevision)

	afterState, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, staleState.BaseLayoutDigest, afterState.BaseLayoutDigest)
	require.Equal(t, staleState.ConnectorRevision, afterState.ConnectorRevision)
}
