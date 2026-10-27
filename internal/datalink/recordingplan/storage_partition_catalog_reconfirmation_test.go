package recordingplan

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newCapabilityProjection(t *testing.T, f *mysqlPartitionMaintenanceFixture) GroupSchemaLayout {
	t.Helper()
	options := *f.target.layout.MySQL
	options.CapabilityRevision = "cap-2"
	layout, err := ManagedMySQLLongGroupLayout(f.target.layout.TableName, f.target.layout.OwnerColumn, options)
	require.NoError(t, err)
	return layout
}

func actualNativeInspection(f *mysqlPartitionMaintenanceFixture, layout GroupSchemaLayout, partitions *MySQLPartitionEvidence) TargetTableInspection {
	actual := cloneTargetTableInspection(f.target.group)
	actual.MySQL.ServerFamily = layout.MySQL.ServerFamily
	actual.MySQL.ServerVersion = layout.MySQL.ServerVersion
	actual.MySQL.CapabilityRevision = layout.MySQL.CapabilityRevision
	actual.MySQL.Partitions = cloneMySQLPartitionEvidence(partitions)
	return actual
}

func TestExplicitNativeReconfirmationKeepsAuthorizedCurrentPartitions(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	state, expanded := expandPartitionSidecar(t, f)
	newLayout := newCapabilityProjection(t, f)

	actual := actualNativeInspection(f, newLayout, expanded)
	normalized, err := NormalizeMySQLPartitionCatalogForExplicitReconfirmation(t.Context(), f.db, f.request.Scope, newLayout, actual, "rev-2")
	require.NoError(t, err)
	require.Equal(t, f.base, normalized.MySQL.Partitions)
	require.Equal(t, "cap-2", normalized.MySQL.CapabilityRevision)
	require.Equal(t, expanded, actual.MySQL.Partitions, "the explicit normalizer is read-only")
	afterState, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, state, afterState, "preview normalization never mutates the sidecar")

	_, err = NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, newLayout, actual)
	require.ErrorIs(t, err, ErrMySQLPartitionMaintenanceStateMismatch, "runtime/automatic repair stays strict")
}

func TestExplicitNativeReconfirmationRequiresExactAuthorizedCurrentOrInitial(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, expanded := expandPartitionSidecar(t, f)
	newLayout := newCapabilityProjection(t, f)

	actual := actualNativeInspection(f, newLayout, f.base)
	_, err := NormalizeMySQLPartitionCatalogForExplicitReconfirmation(t.Context(), f.db, f.request.Scope, newLayout, actual, "rev-2")
	require.ErrorIs(t, err, ErrIncompatibleExistingTable, "an existing sidecar current state cannot fall back to the initial prefix")

	_, err = f.db.ExecContext(t.Context(), `DELETE FROM `+mysqlPartitionMaintenanceStateTable+` WHERE workspace_id = ?`, f.request.Scope.WorkspaceID)
	require.NoError(t, err)
	actual.MySQL.Partitions = cloneMySQLPartitionEvidence(f.base)
	normalized, err := NormalizeMySQLPartitionCatalogForExplicitReconfirmation(t.Context(), f.db, f.request.Scope, newLayout, actual, "rev-2")
	require.NoError(t, err, "a normal unmaintained native table has only its confirmed initial shape")
	require.Equal(t, f.base, normalized.MySQL.Partitions)
	actual.MySQL.Partitions = cloneMySQLPartitionEvidence(expanded)
	_, err = NormalizeMySQLPartitionCatalogForExplicitReconfirmation(t.Context(), f.db, f.request.Scope, newLayout, actual, "rev-2")
	require.ErrorIs(t, err, ErrIncompatibleExistingTable, "an expanded catalog without a sidecar is not adoptable")
}

func TestExplicitNativeReconfirmationRejectsCapabilityEvidenceMismatch(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, expanded := expandPartitionSidecar(t, f)
	newLayout := newCapabilityProjection(t, f)
	actual := actualNativeInspection(f, newLayout, expanded)
	actual.MySQL.CapabilityRevision = "cap-1"
	_, err := NormalizeMySQLPartitionCatalogForExplicitReconfirmation(t.Context(), f.db, f.request.Scope, newLayout, actual, "rev-2")
	require.ErrorIs(t, err, ErrMySQLPartitionMaintenanceStateMismatch)
}

func TestExplicitNativeConfirmationRebindsSidecarWithoutResettingCurrent(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, expanded := expandPartitionSidecar(t, f)
	newLayout := newCapabilityProjection(t, f)
	before, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	confirmation := StorageConfirmation{
		Layout: newLayout, ConnectorRevision: "rev-2", OperationID: "reconfirm-cap-2",
		VerifiedDigest: strings.Repeat("e", 64), CreatedTable: false, GroupID: f.request.GroupID,
	}
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC()))
	after, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	state, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, before.Generation, after.Generation)
	require.Equal(t, "cap-2", after.Layout.MySQL.CapabilityRevision)
	require.Equal(t, "rev-2", after.ConnectorRevision)
	require.Equal(t, digestGroupSchemaLayout(newLayout), after.LayoutDigest)
	require.Equal(t, digestGroupSchemaLayout(newLayout), state.BaseLayoutDigest)
	require.Equal(t, "rev-2", state.ConnectorRevision)
	require.Equal(t, expanded, state.Current, "same-table confirmation preserves the authorized current partition projection")
	require.Equal(t, mysqlPartitionMaintenanceReady, state.Status)
}

func TestExplicitNativeRecreateResetsInitialShapeOnce(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, expanded := expandPartitionSidecar(t, f)
	newLayout := newCapabilityProjection(t, f)
	before, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	confirmation := StorageConfirmation{
		Layout: newLayout, ConnectorRevision: "rev-2", OperationID: "recreate-cap-2",
		VerifiedDigest: strings.Repeat("f", 64), CreatedTable: true, GroupID: f.request.GroupID,
	}
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC()))
	after, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	state, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, before.Generation+1, after.Generation)
	require.Equal(t, f.base, state.Current)
	require.Equal(t, after.Generation, state.Generation)
	require.Equal(t, digestGroupSchemaLayout(newLayout), state.BaseLayoutDigest)
	require.Equal(t, "rev-2", state.ConnectorRevision)
	require.NotEqual(t, expanded, state.Current)

	repeated := confirmation
	repeated.OperationID = confirmation.OperationID
	require.NoError(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, repeated, time.Now().UTC()))
	again, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	againState, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, after.Generation, again.Generation)
	require.Equal(t, f.base, againState.Current)
}

func TestExplicitNativeReconfirmationSidecarFailureRollsBackAll(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, _ = expandPartitionSidecar(t, f)
	newLayout := newCapabilityProjection(t, f)
	before, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), numberedPlaceholders(`UPDATE `+mysqlPartitionMaintenanceStateTable+`
		SET base_layout_digest = $7 WHERE `+storageScopeWhere), append(scopeArgs(f.request.Scope), strings.Repeat("0", 64))...)
	require.NoError(t, err)
	stale, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	confirmation := StorageConfirmation{
		Layout: newLayout, ConnectorRevision: "rev-2", OperationID: "reconfirm-rollback",
		VerifiedDigest: strings.Repeat("c", 64), CreatedTable: false, GroupID: f.request.GroupID,
	}
	require.Error(t, RecordStorageConfirmation(t.Context(), f.db, f.request.Scope, confirmation, time.Now().UTC()))
	after, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	afterState, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, before.LayoutDigest, after.LayoutDigest)
	require.Equal(t, before.ConfirmedOperationID, after.ConfirmedOperationID)
	require.Equal(t, stale.BaseLayoutDigest, afterState.BaseLayoutDigest)
	require.Equal(t, stale.Current, afterState.Current)
}
