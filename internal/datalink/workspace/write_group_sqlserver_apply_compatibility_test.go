package workspace

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

type sqlServerApplyCompatibilityFixture struct {
	h       *circuitProfileHarness
	group   *WriteGroup
	digest  string
	initial *recordingplan.SQLServerPartitionEvidence
	current *recordingplan.SQLServerPartitionEvidence
}

func newSQLServerApplyCompatibilityFixture(t *testing.T, mode WriteGroupStorageMode) *sqlServerApplyCompatibilityFixture {
	t.Helper()
	h := newCircuitProfileHarness(t, ":memory:")
	ctx := t.Context()
	connectorID := seedSQLServerConnector(t, h.db, `{"host":"sql-apply.example","database":"recordings","schema":"dbo","username":"u","password":"p"}`)
	draft := h.longDraft(WriteGroupProfileInstantaneous, "sqlserver-apply", h.primary())
	draft.Destination.ConnectorID, draft.Destination.ConnectorRevision = connectorID, "connector-revision-sqlserver"
	draft.Destination.Database, draft.Destination.TableSchema = "recordings", "dbo"
	draft.StorageMode = mode
	savedResult, err := ensureSQLServerGroup(t, h, draft)
	require.NoError(t, err)
	saved := savedResult.Group

	capability := sqlServerTestCapability(connectorID, saved.Destination.ConnectorRevision)
	inspector := &fakeSQLServerReadinessInspector{capability: capability}
	h.service.WithManagedTableInspector(inspector)
	h.service.now = func() time.Time { return time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC) }

	layout, err := h.service.managedSQLServerSchemaLayout(ctx, saved)
	require.NoError(t, err)
	require.NotNil(t, layout.SQLServer)
	options := *layout.SQLServer
	groupProof, groupTarget := sqlServerTestRuntimeCatalog(layout.Columns, layout.Indexes, saved.Destination.TableName, layout.OwnerColumn, 101, options)
	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptTarget := sqlServerTestRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 102, options)
	require.NotNil(t, receiptProof)
	require.NotNil(t, receiptTarget)
	if mode == WriteGroupStorageModeNativeMonthly {
		groupProof.Evidence.NativeMonthly = sqlServerManagedPartitionEvidenceForTest(layout.SQLServer.NativeMonthly)
		groupProof.Evidence.PartitionInspected = true
		groupProof.Evidence.Plain = false
		groupTarget.SQLServer.NativeMonthly = sqlServerTargetPartitionEvidenceForTest(layout.SQLServer.NativeMonthly)
		groupTarget.SQLServer.PartitionInspected = true
		groupTarget.SQLServer.Plain = false
	}
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(*layout, "recordings", groupTarget, receiptTarget)
	require.NoError(t, err)

	var current *recordingplan.SQLServerPartitionEvidence
	if mode == WriteGroupStorageModeNativeMonthly {
		current, err = recordingplan.PlanSQLServerNativeMonthlyExtension(layout.SQLServer.NativeMonthly, time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		scope := StorageScopeOf(saved)
		baseDigest := seedConfirmedStorageScope(t, h.db, scope, *layout, "schema-op-sqlserver")
		_, err = h.db.ExecContext(ctx, `UPDATE write_group_storage_scopes SET confirmed_digest = ?, confirmed_connector_revision = ? WHERE workspace_id = ? AND connector_id = ? AND database_name = ? AND table_schema = ? AND storage_mode = ? AND layout_version = ?`,
			digest, saved.Destination.ConnectorRevision, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion)
		require.NoError(t, err)
		seedSQLServerReadinessMaintenanceState(t, h.db, scope, baseDigest, saved.Destination.ConnectorRevision, current)
		groupProof.Evidence.NativeMonthly = sqlServerManagedPartitionEvidenceForTest(current)
		groupTarget.SQLServer.NativeMonthly = sqlServerTargetPartitionEvidenceForTest(current)
	}
	inspector.group, inspector.receipt = groupProof, receiptProof
	_, err = h.db.ExecContext(ctx, `UPDATE write_groups SET destination_schema_revision = ?, destination_schema_digest = ? WHERE id = ?`,
		"schema-op-sqlserver", digest, saved.ID)
	require.NoError(t, err)
	return &sqlServerApplyCompatibilityFixture{
		h: h, group: saved, digest: digest, initial: layout.SQLServer.NativeMonthly, current: current,
	}
}

func (f *sqlServerApplyCompatibilityFixture) applyMutation(ctx context.Context, t *testing.T) WriteGroupMutation {
	t.Helper()
	current, err := f.h.workspace.GetOrCreate(ctx)
	require.NoError(t, err)
	return WriteGroupMutation{
		WorkspaceID: f.group.WorkspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision: f.group.Revision, ExpectedConnectorRevision: f.group.Destination.ConnectorRevision,
	}
}

func TestSQLServerApplyCompatibilityFreezesReadinessProof(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			f := newSQLServerApplyCompatibilityFixture(t, mode)
			ctx := t.Context()
			readiness, err := f.h.service.Readiness(ctx, f.group.ID)
			require.NoError(t, err)
			require.True(t, readiness.Ready, "readiness issues: %+v", readiness.Issues)
			require.Equal(t, f.digest, readiness.SchemaDigest)
			require.NotNil(t, readiness.runtimeLayout)
			require.NotNil(t, readiness.runtimeLayout.SQLServer)

			applied, err := f.h.service.Apply(ctx, f.group.ID, f.applyMutation(ctx, t))
			require.NoError(t, err)
			require.NotNil(t, applied.Group)
			require.NotEmpty(t, applied.Group.AppliedRevision)
			require.Equal(t, f.digest, applied.Group.Destination.SchemaDigest)

			var payload string
			require.NoError(t, f.h.db.QueryRowContext(ctx, `SELECT payload FROM wg_runtime_versions WHERE group_id = ? AND group_revision = ?`,
				applied.Group.ID, applied.Group.AppliedRevision).Scan(&payload))
			var frozen WriteGroupAppliedSnapshot
			require.NoError(t, json.Unmarshal([]byte(payload), &frozen))
			require.NotNil(t, frozen.RuntimeLayout)
			require.Equal(t, dbtarget.SQLDialectSQLServer, frozen.RuntimeLayout.Dialect)
			require.Equal(t, f.digest, frozen.RuntimeLayout.SchemaDigest)
			require.NotNil(t, frozen.RuntimeLayout.SQLServer)
			require.NotNil(t, frozen.RuntimeLayout.SQLServer.Group)
			require.NotNil(t, frozen.RuntimeLayout.SQLServer.Receipt)
			if mode == WriteGroupStorageModeNativeMonthly {
				require.NotNil(t, frozen.RuntimeLayout.SQLServer.Group.Evidence.NativeMonthly)
				require.Equal(t, f.initial.Partitions, recordingPlanSQLServerLayoutPartition(frozen.RuntimeLayout.SQLServer.Group.Evidence.NativeMonthly).Partitions)
				require.NotEqual(t, f.current.Partitions, recordingPlanSQLServerLayoutPartition(frozen.RuntimeLayout.SQLServer.Group.Evidence.NativeMonthly).Partitions)
			}
		})
	}
}

func TestSQLServerApplyCompatibilityRejectsStaleCASBeforeApply(t *testing.T) {
	for _, stale := range []string{"workspace", "group", "connector"} {
		t.Run(stale, func(t *testing.T) {
			f := newSQLServerApplyCompatibilityFixture(t, WriteGroupStorageModePlain)
			ctx := t.Context()
			mutation := f.applyMutation(ctx, t)
			switch stale {
			case "workspace":
				mutation.ExpectedWorkspaceRevision = "stale-workspace-revision"
			case "group":
				mutation.ExpectedGroupRevision = "stale-group-revision"
			case "connector":
				mutation.ExpectedConnectorRevision = "stale-connector-revision"
			}
			applied, err := f.h.service.Apply(ctx, f.group.ID, mutation)
			require.Error(t, err)
			require.Nil(t, applied)
			saved, err := f.h.service.repo.Get(ctx, f.group.WorkspaceID, f.group.ID)
			require.NoError(t, err)
			require.Empty(t, saved.AppliedRevision, "a stale CAS must not apply the group")
		})
	}
}
