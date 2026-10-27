package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

type fakeSQLServerReadinessInspector struct {
	capability *dbtarget.SQLServerRecordingCapability
	group      *dbtarget.SQLServerManagedInspection
	receipt    *dbtarget.SQLServerManagedInspection
}

func (f *fakeSQLServerReadinessInspector) InspectTableAtRevision(context.Context, string, string, string, string) (*dbtarget.TableInspection, error) {
	return &dbtarget.TableInspection{Status: dbtarget.TableInspectionMissing}, nil
}

func (f *fakeSQLServerReadinessInspector) ProbeSQLServerRecordingAtRevision(context.Context, string, string) (*dbtarget.SQLServerRecordingCapability, error) {
	return f.capability, nil
}

func (f *fakeSQLServerReadinessInspector) InspectSQLServerManagedTableAtRevision(_ context.Context, _, _, _, _, _, table string) (*dbtarget.SQLServerManagedInspection, error) {
	if table == dbtarget.EffectReceiptTable {
		return f.receipt, nil
	}
	return f.group, nil
}

func seedSQLServerReadinessMaintenanceState(t *testing.T, db *sql.DB, scope recordingplan.StorageScope, baseDigest, connectorRevision string, current *recordingplan.SQLServerPartitionEvidence) {
	t.Helper()
	ensureSQLServerPartitionTable(t, db)
	encoded, err := json.Marshal(current)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_sqlserver_partition_maintenance (
		workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version, base_layout_digest,
		connector_revision, generation, current_evidence, revision, last_operation_id, status, reason, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, 1, 'op-maint-sql', 'ready', '', ?)`,
		scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion,
		baseDigest, connectorRevision, string(encoded), time.Now().UTC().Format(time.RFC3339Nano))
	require.NoError(t, err)
}

func TestSQLServerReadinessBuildsFrozenManagedProof(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			ctx := t.Context()
			connectorID := seedSQLServerConnector(t, h.db, `{"host":"sql.example","database":"recordings","schema":"dbo","username":"u","password":"p"}`)
			draft := h.longDraft(WriteGroupProfileInstantaneous, "sqlserver-readiness", h.primary())
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
			if mode == WriteGroupStorageModeNativeMonthly {
				require.NotNil(t, layout.SQLServer.NativeMonthly)
			}
			options := *layout.SQLServer
			groupProof, groupTarget := sqlServerTestRuntimeCatalog(layout.Columns, layout.Indexes, saved.Destination.TableName, layout.OwnerColumn, 101, options)
			receiptColumns := []recordingplan.GroupSchemaColumn{
				{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
				{Name: "payload_digest", SQLType: "CHAR(64)"},
				{Name: "committed_at", SQLType: "DATETIME2(6)"},
				{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
			}
			receiptProof, receiptTarget := sqlServerTestRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 102, options)
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
				require.Greater(t, len(current.Partitions), len(layout.SQLServer.NativeMonthly.Partitions))
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

			readiness, err := h.service.Readiness(ctx, saved.ID)
			require.NoError(t, err)
			require.True(t, readiness.ConfigReady, "config issues: %+v", readiness.Issues)
			require.True(t, readiness.SchemaReady, "schema issues: %+v", readiness.Issues)
			require.True(t, readiness.Ready, "readiness issues: %+v", readiness.Issues)
			require.NotNil(t, readiness.runtimeLayout)
			require.Equal(t, dbtarget.SQLDialectSQLServer, readiness.runtimeLayout.Dialect)
			require.Equal(t, digest, readiness.SchemaDigest)
			require.NotNil(t, readiness.runtimeLayout.SQLServer)
			require.Equal(t, saved.Destination.ConnectorID, readiness.runtimeLayout.SQLServer.ConnectorID)
			require.Equal(t, digest, readiness.runtimeLayout.SchemaDigest)
			require.NotNil(t, readiness.runtimeLayout.SQLServer.Group)
			require.NotNil(t, readiness.runtimeLayout.SQLServer.Receipt)
			if mode == WriteGroupStorageModeNativeMonthly {
				require.NotNil(t, readiness.runtimeLayout.SQLServer.Group.Evidence.NativeMonthly)
				require.Equal(t, layout.SQLServer.NativeMonthly.Partitions, recordingPlanSQLServerLayoutPartition(readiness.runtimeLayout.SQLServer.Group.Evidence.NativeMonthly).Partitions)
				require.Equal(t, current.Partitions, readiness.SQLServerCurrent.Partitions, "current evidence stays in the sidecar")
			}
		})
	}
}
