package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

func ensurePartitionMaintenanceTable(t *testing.T, db *sql.DB) {
	t.Helper()
	ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", "041_mysql_partition_maintenance_sqlite.up.sql"))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(ddl))
	require.NoError(t, err)
}

func seedPartitionMaintenanceState(
	t *testing.T, db *sql.DB, scope recordingplan.StorageScope,
	baseDigest, connectorRev string, gen int,
	current *recordingplan.MySQLPartitionEvidence, rev int64,
	opID, status, reason string, updatedAt time.Time,
) {
	t.Helper()
	ensurePartitionMaintenanceTable(t, db)
	encoded, err := json.Marshal(current)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `
		INSERT INTO write_group_mysql_partition_maintenance (
			workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
			base_layout_digest, connector_revision, generation, current_evidence, revision,
			last_operation_id, status, reason, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion,
		baseDigest, connectorRev, gen, string(encoded), rev, opID, status, reason, updatedAt.Format(time.RFC3339Nano))
	require.NoError(t, err)
}

func seedConfirmedStorageScope(t *testing.T, db *sql.DB, scope recordingplan.StorageScope, layout recordingplan.GroupSchemaLayout, opID string) string {
	t.Helper()
	encoded, err := json.Marshal(layout)
	require.NoError(t, err)
	sum := sha256.Sum256(encoded)
	digest := hex.EncodeToString(sum[:])
	_, err = db.ExecContext(t.Context(), `
		INSERT INTO write_group_storage_scopes (
			workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version,
			table_name, owner_column, layout_digest, generation, confirmed_layout,
			confirmed_connector_revision, confirmed_operation_id, confirmed_digest
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?, ?)
		ON CONFLICT (workspace_id, connector_id, database_name, table_schema, storage_mode, layout_version)
		DO UPDATE SET confirmed_layout = excluded.confirmed_layout, layout_digest = excluded.layout_digest,
			confirmed_connector_revision = excluded.confirmed_connector_revision,
			confirmed_operation_id = excluded.confirmed_operation_id, confirmed_digest = excluded.confirmed_digest
	`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion,
		layout.TableName, layout.OwnerColumn, digest, string(encoded), "connector-revision-mysql", opID, digest)
	require.NoError(t, err)
	return digest
}

func initialNativePartitions(t *testing.T) *recordingplan.MySQLPartitionEvidence {
	t.Helper()
	p, err := recordingplan.NewMySQLNativeMonthlyPartition(time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	return p
}

func expandedNativePartitions(t *testing.T, initial *recordingplan.MySQLPartitionEvidence) *recordingplan.MySQLPartitionEvidence {
	t.Helper()
	expanded, err := recordingplan.PlanMySQLNativeMonthlyExtension(initial, time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	return expanded
}

type fakeMySQLManagedRuntimeInspector struct {
	capability *dbtarget.MySQLRecordingCapability
	group      *dbtarget.MySQLManagedInspection
	receipt    *dbtarget.MySQLManagedInspection
	err        error
}

func (f *fakeMySQLManagedRuntimeInspector) InspectTableAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) (*dbtarget.TableInspection, error) {
	return &dbtarget.TableInspection{Status: "exists"}, nil
}

func (f *fakeMySQLManagedRuntimeInspector) ProbeMySQLRecordingAtRevision(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.MySQLRecordingCapability, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.capability, nil
}

func (f *fakeMySQLManagedRuntimeInspector) InspectMySQLManagedTableAtRevision(ctx context.Context, connectorID, expectedRevision, db, capRev, table string) (*dbtarget.MySQLManagedInspection, error) {
	if f.err != nil {
		return nil, f.err
	}
	if table == dbtarget.EffectReceiptTable {
		return f.receipt, nil
	}
	return f.group, nil
}

func TestNormalizeMySQLPartitionInspection_SkipsMissingReceiptsAndPlain(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	ctx := t.Context()

	// Missing table
	missingState := recordingplan.TargetTableInspection{Status: "missing"}
	scope := recordingplan.SchemaPreviewScope{}
	got, err := h.service.NormalizeMySQLPartitionInspection(ctx, scope, missingState)
	require.NoError(t, err)
	require.Equal(t, "missing", got.Status)

	// Plain scope (NativeMonthly is nil)
	const testOwnerCol = "_gw_owner_0123456789abcdef01234567"
	plainLayout, err := recordingplan.ManagedMySQLLongGroupLayout("gw_plain", testOwnerCol, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1",
		Engine: "InnoDB", Collation: "utf8mb4_0900_bin", Charset: "utf8mb4",
	})
	require.NoError(t, err)
	plainScope := recordingplan.SchemaPreviewScope{
		WorkspaceID: "ws-1", ConnectorID: "conn-1", Database: "testdb", Schema: "testdb",
		GroupLayout: &plainLayout,
	}
	existingState := recordingplan.TargetTableInspection{
		Status: "exists", Columns: []string{testOwnerCol}, MySQL: &recordingplan.TargetMySQLTableEvidence{ServerFamily: "mysql"},
	}
	got, err = h.service.NormalizeMySQLPartitionInspection(ctx, plainScope, existingState)
	require.NoError(t, err)
	require.Equal(t, existingState, got)

	// Receipt table (missing owner column)
	nativeLayout, err := recordingplan.ManagedMySQLLongGroupLayout("gw_native", testOwnerCol, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1",
		Engine: "InnoDB", Collation: "utf8mb4_0900_bin", Charset: "utf8mb4",
		NativeMonthly: initialNativePartitions(t),
	})
	require.NoError(t, err)
	nativeScope := recordingplan.SchemaPreviewScope{
		WorkspaceID: "ws-1", ConnectorID: "conn-1", Database: "testdb", Schema: "testdb",
		GroupLayout: &nativeLayout,
	}
	receiptState := recordingplan.TargetTableInspection{
		Status: "exists", Columns: []string{"effect_id", "applied_at"}, MySQL: &recordingplan.TargetMySQLTableEvidence{ServerFamily: "mysql"},
	}
	got, err = h.service.NormalizeMySQLPartitionInspection(ctx, nativeScope, receiptState)
	require.NoError(t, err)
	require.Equal(t, receiptState, got)
}

func TestNormalizeMySQLPartitionInspection_NormalizesExistingNativeTable(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	ctx := t.Context()

	initial := initialNativePartitions(t)
	expanded := expandedNativePartitions(t, initial)

	scope := recordingplan.StorageScope{
		WorkspaceID: "ws-1", ConnectorID: "conn-1", Database: "testdb", TableSchema: "testdb",
		StorageMode: "native_monthly", LayoutVersion: "managed_long_v1",
	}
	tableObj := longStorageScope{
		WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID,
		Database: scope.Database, TableSchema: scope.TableSchema,
		Mode: WriteGroupStorageModeNativeMonthly, Layout: WriteGroupLayoutManagedLongV1,
	}
	tableName := tableObj.tableName()
	ownerColumn := tableObj.ownerColumn()

	layout, err := recordingplan.ManagedMySQLLongGroupLayout(tableName, ownerColumn, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1",
		Engine: "InnoDB", Collation: "utf8mb4_0900_bin", Charset: "utf8mb4",
		NativeMonthly: initial,
	})
	require.NoError(t, err)

	baseDigest := seedConfirmedStorageScope(t, h.db, scope, layout, "op-confirm-1")
	seedPartitionMaintenanceState(t, h.db, scope, baseDigest, "connector-revision-mysql", 1, expanded, 1, "op-maint-1", "ready", "", time.Now())

	previewScope := recordingplan.SchemaPreviewScope{
		WorkspaceID: "ws-1", ConnectorID: "conn-1", Database: "testdb", Schema: "testdb",
		GroupLayout: &layout,
	}

	physical := makeMySQLPhysical(tableName, ownerColumn, layout, expanded)
	actualState := recordingPlanMySQLInspection(physical)

	got, err := h.service.NormalizeMySQLPartitionInspection(ctx, previewScope, actualState)
	require.NoError(t, err)
	require.NotNil(t, got.MySQL)
	require.NotNil(t, got.MySQL.Partitions)
	// Normalized must project back to initial partition count
	require.Equal(t, len(initial.Partitions), len(got.MySQL.Partitions.Partitions))
	// Original actualState must be unmodified
	require.Equal(t, len(expanded.Partitions), len(actualState.MySQL.Partitions.Partitions))
}

func TestMySQLReadiness_PartitionMaintenanceSidecarProjection(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	ensurePartitionMaintenanceTable(t, h.db)

	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "native-maint", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-mysql"
	group.StorageMode = WriteGroupStorageModeNativeMonthly
	saved := mustEnsureMySQLGroup(t, h, group)

	scope := StorageScopeOf(saved)

	// 1. When sidecar row is absent
	result := &WriteGroupReadiness{ConfigReady: true}
	h.service.populateMySQLPartitionMaintenanceReadiness(t.Context(), saved, result)
	require.Nil(t, result.PartitionMaintenance)

	// 2. When sidecar row is "ready"
	initial := initialNativePartitions(t)
	seedPartitionMaintenanceState(t, h.db, scope, strings.Repeat("a", 64), "connector-revision-mysql", 1, initial, 1, "op-1", "ready", "", time.Now())
	result = &WriteGroupReadiness{ConfigReady: true}
	h.service.populateMySQLPartitionMaintenanceReadiness(t.Context(), saved, result)
	require.NotNil(t, result.PartitionMaintenance)
	require.Equal(t, "ready", result.PartitionMaintenance.Status)
	require.Empty(t, result.PartitionMaintenance.Code)
	require.True(t, result.ConfigReady)
	require.Empty(t, result.Issues)

	// 3. When sidecar row is "blocked"
	_, err := h.db.ExecContext(t.Context(), `UPDATE write_group_mysql_partition_maintenance SET status = 'blocked', reason = 'permission_lost'`)
	require.NoError(t, err)
	result = &WriteGroupReadiness{ConfigReady: true}
	h.service.populateMySQLPartitionMaintenanceReadiness(t.Context(), saved, result)
	require.NotNil(t, result.PartitionMaintenance)
	require.Equal(t, "blocked", result.PartitionMaintenance.Status)
	require.Equal(t, "partition_maintenance_blocked", result.PartitionMaintenance.Code)
	require.Equal(t, "permission_lost", result.PartitionMaintenance.Reason)
	// Blocked maintenance must NOT append fatal issues or change ConfigReady
	require.Empty(t, result.Issues)
	require.True(t, result.ConfigReady)
}

func makeMySQLPhysical(tableName, ownerColumn string, layout recordingplan.GroupSchemaLayout, partitions *recordingplan.MySQLPartitionEvidence) *dbtarget.MySQLManagedInspection {
	marker := "gateway"
	evidence := &dbtarget.MySQLManagedTableEvidence{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: layout.MySQL.CapabilityRevision, Database: "partition_test",
		Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
		ColumnCollations: map[string]string{}, ColumnDefaults: map[string]*string{}, JSONColumns: map[string]bool{},
		Partitions: dbtargetMySQLPartitions(partitions),
	}
	tableColumns := make([]dbtarget.ColumnInfo, 0, len(layout.Columns))
	for _, col := range layout.Columns {
		tableColumns = append(tableColumns, dbtarget.ColumnInfo{
			Name: col.Name, DataType: col.SQLType, PrimaryKey: col.PrimaryKey, Nullable: col.Nullable,
		})
		evidence.ColumnDefaults[col.Name] = nil
		if col.PrimaryKey {
			evidence.PrimaryKey = append(evidence.PrimaryKey, col.Name)
		}
		if strings.HasPrefix(col.SQLType, "VARCHAR(") || strings.HasPrefix(col.SQLType, "CHAR(") || col.SQLType == "LONGTEXT" {
			evidence.ColumnCollations[col.Name] = "utf8mb4_0900_bin"
		}
		if col.SQLType == "JSON" {
			evidence.JSONColumns[col.Name] = true
		}
	}
	evidence.ColumnDefaults[ownerColumn] = &marker

	indexes := make([]dbtarget.IndexInfo, 0, len(layout.Indexes))
	for _, idx := range layout.Indexes {
		indexes = append(indexes, dbtarget.IndexInfo{Name: idx.Name, Columns: idx.Columns})
	}

	return &dbtarget.MySQLManagedInspection{
		Inspection: &dbtarget.TableInspection{
			Status: dbtarget.TableInspectionExists, Schema: "partition_test", Table: tableName,
			Columns: tableColumns, IndexesInspected: true, Indexes: indexes,
		},
		Evidence: evidence,
		Indexes:  indexes,
	}
}

func makeMySQLReceiptPhysical() *dbtarget.MySQLManagedInspection {
	marker := "gateway"
	cols := []dbtarget.ColumnInfo{
		{Name: "effect_key", DataType: "VARCHAR(191)", PrimaryKey: true},
		{Name: "payload_digest", DataType: "CHAR(64)"},
		{Name: "committed_at", DataType: "DATETIME(6)"},
		{Name: "_gw_effect_receipts_v1", DataType: "VARCHAR(64)"},
	}
	evidence := &dbtarget.MySQLManagedTableEvidence{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "partition-cap-1", Database: "partition_test",
		Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
		PrimaryKey: []string{"effect_key"},
		ColumnCollations: map[string]string{
			"effect_key":             "utf8mb4_0900_bin",
			"payload_digest":         "utf8mb4_0900_bin",
			"_gw_effect_receipts_v1": "utf8mb4_0900_bin",
		},
		ColumnDefaults: map[string]*string{
			"effect_key":             nil,
			"payload_digest":         nil,
			"committed_at":           nil,
			"_gw_effect_receipts_v1": &marker,
		},
		JSONColumns: map[string]bool{},
		Partitions:  nil,
	}
	return &dbtarget.MySQLManagedInspection{
		Inspection: &dbtarget.TableInspection{
			Status: dbtarget.TableInspectionExists, Schema: "partition_test", Table: dbtarget.EffectReceiptTable,
			Columns: cols, IndexesInspected: true, Indexes: []dbtarget.IndexInfo{},
		},
		Evidence: evidence,
		Indexes:  []dbtarget.IndexInfo{},
	}
}

func TestReadMySQLRuntimeLayout_ClonesEvidenceAndPreservesInitialPartitions(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	ensurePartitionMaintenanceTable(t, h.db)
	h.service.now = func() time.Time { return time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC) }

	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "native-runtime-layout", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-mysql"
	group.StorageMode = WriteGroupStorageModeNativeMonthly
	saved := mustEnsureMySQLGroup(t, h, group)

	capability := &dbtarget.MySQLRecordingCapability{
		ConnectorID: connectorID, ConnectorRevision: saved.Destination.ConnectorRevision, Database: "partition_test",
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "partition-cap-1",
		Connected: true, MetadataStatus: "verified", TransactionalEngine: true,
	}

	initial := initialNativePartitions(t)
	expandedRP := expandedNativePartitions(t, initial)

	layout, err := recordingplan.ManagedMySQLLongGroupLayout(saved.Destination.TableName, "_gw_owner_"+longStorageScopeOf(saved).digest("owner")[:24], recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "partition-cap-1",
		Engine: "InnoDB", Collation: "utf8mb4_0900_bin", Charset: "utf8mb4",
		NativeMonthly: initial,
	})
	require.NoError(t, err)

	scope := StorageScopeOf(saved)
	baseDigest := seedConfirmedStorageScope(t, h.db, scope, layout, "op-confirm-1")
	seedPartitionMaintenanceState(t, h.db, scope, baseDigest, "connector-revision-mysql", 1, expandedRP, 1, "op-maint-1", "ready", "", time.Now())

	groupPhysical := makeMySQLPhysical(saved.Destination.TableName, layout.OwnerColumn, layout, expandedRP)
	receiptPhysical := makeMySQLReceiptPhysical()

	initialGroupInspection := recordingPlanMySQLInspection(groupPhysical)
	initialGroupInspection.MySQL.Partitions = initial
	receiptInspection := recordingPlanMySQLInspection(receiptPhysical)

	initialDigest, err := recordingplan.VerifyManagedMySQLCatalog(layout, saved.Destination.Database, initialGroupInspection, receiptInspection)
	require.NoError(t, err)

	saved.Destination.SchemaRevision = "schema-rev-1"
	saved.Destination.SchemaDigest = initialDigest
	saved.WritePolicy.DedupeCapability = "receipt"

	inspector := &fakeMySQLManagedRuntimeInspector{
		capability: capability,
		group:      groupPhysical,
		receipt:    receiptPhysical,
	}
	h.service.WithManagedTableInspector(inspector)

	snapshot, err := h.service.readinessSnapshot(t.Context(), saved.ID, saved)
	require.NoError(t, err)

	frozen, err := h.service.readMySQLRuntimeLayout(t.Context(), snapshot)
	require.NoError(t, err)
	require.NotNil(t, frozen)

	// 1. Frozen proof carries the INITIAL partitions, NOT the expanded partitions
	frozenPartitions := frozen.MySQL.Group.Evidence.Partitions
	require.NotNil(t, frozenPartitions)
	require.Equal(t, len(initial.Partitions), len(frozenPartitions.Partitions), "frozen proof must carry initial partition count")

	// 2. The live actual evidence object remains unmodified (still carrying expanded partitions)
	require.Equal(t, len(expandedRP.Partitions), len(groupPhysical.Evidence.Partitions.Partitions), "live actual evidence object must be unmodified")
}
