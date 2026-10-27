package workspace

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/recordingplan"
)

func ensureSQLServerPartitionTable(t *testing.T, db *sql.DB) {
	t.Helper()
	ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", "042_sqlserver_partition_maintenance_sqlite.up.sql"))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(ddl))
	require.NoError(t, err)
}

func sqlServerNativeTestLayout(t *testing.T, initial *recordingplan.SQLServerPartitionEvidence) recordingplan.GroupSchemaLayout {
	t.Helper()
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout("gw_long_sqlserver_normalize", "_gw_owner_0123456789abcdef01234567",
		recordingplan.SQLServerGroupSchemaOptions{
			Database: "recordings", Schema: "dbo", Collation: "Latin1_General_100_BIN2", ServerVersion: "16.0.4180.1",
			ServerEdition: "Enterprise Edition", CapabilityRevision: "cap-sqlserver", EngineEdition: 3,
			CompatibilityLevel: 160, NativeMonthly: initial,
		})
	require.NoError(t, err)
	return layout
}

func sqlServerNativeTestState(layout recordingplan.GroupSchemaLayout, native *recordingplan.SQLServerPartitionEvidence) recordingplan.TargetTableInspection {
	columns := make([]string, len(layout.Columns))
	types, nullable, primary := map[string]string{}, map[string]bool{}, map[string]bool{}
	collations, defaults := map[string]string{}, map[string]*string{}
	for i, column := range layout.Columns {
		columns[i] = column.Name
		types[column.Name], nullable[column.Name], primary[column.Name] = column.SQLType, column.Nullable, column.PrimaryKey
		if len(column.SQLType) >= 8 && (column.SQLType[:8] == "NVARCHAR" || column.SQLType[:4] == "CHAR") {
			collations[column.Name] = layout.SQLServer.Collation
		}
	}
	ownerDefault := "gateway"
	defaults[layout.OwnerColumn] = &ownerDefault
	return recordingplan.TargetTableInspection{
		Status: "exists", Columns: columns, ColumnTypes: types, ColumnNullable: nullable, ColumnPrimaryKey: primary,
		Indexes:          []recordingplan.TargetIndexInspection{{Name: layout.Indexes[0].Name, Columns: layout.Indexes[0].Columns}},
		IndexesInspected: true,
		SQLServer: &recordingplan.TargetSQLServerTableEvidence{
			ServerVersion: layout.SQLServer.ServerVersion, ServerEdition: layout.SQLServer.ServerEdition,
			CapabilityRevision: layout.SQLServer.CapabilityRevision, Database: layout.SQLServer.Database, Schema: layout.SQLServer.Schema,
			EngineEdition: layout.SQLServer.EngineEdition, CompatibilityLevel: layout.SQLServer.CompatibilityLevel,
			MetadataVerified: true, TableType: "USER_TABLE", ObjectID: 11, Plain: native == nil,
			PrimaryKey: []string{"bucket_start", "record_id"}, ColumnCollations: collations, ColumnDefaults: defaults,
			ProvenanceJSONChecked: true, NativeMonthly: sqlServerPartitionTarget(native), PartitionInspected: native != nil,
		},
	}
}

func sqlServerPartitionTarget(evidence *recordingplan.SQLServerPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName, SchemeName: evidence.SchemeName,
		Filegroup: evidence.Filegroup, Partitions: append([]recordingplan.SQLServerPartitionBoundary(nil), evidence.Partitions...),
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func sqlServerReceiptTestState(layout recordingplan.GroupSchemaLayout) recordingplan.TargetTableInspection {
	columns := []string{"effect_key", "payload_digest", "committed_at", "_gw_effect_receipts_v1"}
	types := map[string]string{"effect_key": "CHAR(64)", "payload_digest": "CHAR(64)", "committed_at": "DATETIME2(6)", "_gw_effect_receipts_v1": "NVARCHAR(16)"}
	nullable := map[string]bool{"effect_key": false, "payload_digest": false, "committed_at": false, "_gw_effect_receipts_v1": false}
	primary := map[string]bool{"effect_key": true, "payload_digest": false, "committed_at": false, "_gw_effect_receipts_v1": false}
	ownerDefault := "gateway"
	return recordingplan.TargetTableInspection{
		Status: "exists", Columns: columns, ColumnTypes: types, ColumnNullable: nullable, ColumnPrimaryKey: primary,
		IndexesInspected: true,
		SQLServer: &recordingplan.TargetSQLServerTableEvidence{
			ServerVersion: layout.SQLServer.ServerVersion, ServerEdition: layout.SQLServer.ServerEdition,
			CapabilityRevision: layout.SQLServer.CapabilityRevision, Database: layout.SQLServer.Database, Schema: layout.SQLServer.Schema,
			EngineEdition: layout.SQLServer.EngineEdition, CompatibilityLevel: layout.SQLServer.CompatibilityLevel,
			MetadataVerified: true, TableType: "USER_TABLE", ObjectID: 22, Plain: true, PrimaryKey: []string{"effect_key"},
			ColumnCollations: map[string]string{"effect_key": layout.SQLServer.Collation, "payload_digest": layout.SQLServer.Collation, "_gw_effect_receipts_v1": layout.SQLServer.Collation},
			ColumnDefaults:   map[string]*string{"_gw_effect_receipts_v1": &ownerDefault},
		},
	}
}

func TestNormalizeSQLServerPartitionInspectionNormalizesOnlyAuthorizedCurrent(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	db := h.db
	ensureSQLServerPartitionTable(t, db)

	initial, err := recordingplan.NewSQLServerNativeMonthlyPartition(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), "pf_gw_long_sqlserver_normalize", "ps_gw_long_sqlserver_normalize")
	require.NoError(t, err)
	expanded, err := recordingplan.PlanSQLServerNativeMonthlyExtension(initial, time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	layout := sqlServerNativeTestLayout(t, initial)
	scope := recordingplan.StorageScope{
		WorkspaceID: "ws-sql", ConnectorID: "conn-sql", Database: "recordings", TableSchema: "dbo",
		StorageMode: "native_monthly", LayoutVersion: "managed_long_v1",
	}
	group := sqlServerNativeTestState(layout, initial)
	receipt := sqlServerReceiptTestState(layout)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, "recordings", group, receipt)
	require.NoError(t, err)
	baseDigest := seedConfirmedStorageScope(t, db, scope, layout, "op-confirm-sql")
	_, err = db.ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET confirmed_digest=? WHERE workspace_id=? AND connector_id=? AND database_name=? AND table_schema=? AND storage_mode=? AND layout_version=?`,
		digest, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion)
	require.NoError(t, err)
	encoded := recordingplan.CloneSQLServerPartitionEvidence(expanded)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_sqlserver_partition_maintenance (
		workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version, base_layout_digest,
		connector_revision, generation, current_evidence, revision, last_operation_id, status, reason, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, 1, 'op-maint', 'ready', '', ?)
	`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion,
		baseDigest, "connector-revision-mysql", mustJSON(t, encoded), time.Now().UTC().Format(time.RFC3339Nano))
	require.NoError(t, err)

	actual := sqlServerNativeTestState(layout, expanded)
	_, err = recordingplan.VerifyManagedSQLServerCatalog(layout, "recordings", actual, receipt)
	require.Error(t, err, "raw post-maintenance catalog must not verify against the frozen initial layout")
	normalized, err := h.service.NormalizeSQLServerPartitionInspection(t.Context(),
		recordingplan.SchemaPreviewScope{WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID, Database: scope.Database,
			Schema: scope.TableSchema, GroupLayout: &layout}, actual, receipt)
	require.NoError(t, err)
	require.NotNil(t, normalized.SQLServer.NativeMonthly)
	require.Len(t, normalized.SQLServer.NativeMonthly.Partitions, len(initial.Partitions))
	require.Len(t, actual.SQLServer.NativeMonthly.Partitions, len(expanded.Partitions), "actual catalog must remain untouched")

	drifted := sqlServerNativeTestState(layout, expanded)
	drifted.SQLServer.NativeMonthly.Partitions[0].LessThan = "2000-01-01T00:00:00.000000"
	_, err = h.service.NormalizeSQLServerPartitionInspection(t.Context(),
		recordingplan.SchemaPreviewScope{WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID, Database: scope.Database,
			Schema: scope.TableSchema, GroupLayout: &layout}, drifted, receipt)
	require.ErrorIs(t, err, recordingplan.ErrSQLServerPartitionMaintenanceStateMismatch)
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	return string(encoded)
}
