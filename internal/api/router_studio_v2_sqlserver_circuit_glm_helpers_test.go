package api

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// sqlServerGLMExpectedInitialEvidence pins the canonical initial native
// layout: RANGE RIGHT datetime2(6) on bucket_start, p_past + four month
// ranges + p_future MAXVALUE, ordinals 1..6, strictly consecutive months.
func sqlServerGLMExpectedInitialEvidence(evidence *recordingplan.SQLServerPartitionEvidence) error {
	invalid := func(reason string) error {
		return fmt.Errorf("initial native evidence is not canonical: %s", reason)
	}
	if evidence == nil || evidence.Method != "RANGE RIGHT" || evidence.Column != "bucket_start" ||
		len(evidence.Partitions) != 6 {
		return invalid("method, column, or boundary count")
	}
	if evidence.Partitions[0].Name != "p_past" || evidence.Partitions[0].Ordinal != 1 {
		return invalid("past catch-all")
	}
	previous, err := parseGLMBoundary(evidence.Partitions[0].LessThan)
	if err != nil {
		return invalid("past catch-all boundary text")
	}
	if evidence.Partitions[0].Name != "p_past" || previous.Format("200601") == "" {
		return invalid("past catch-all boundary month")
	}
	for index, part := range evidence.Partitions[1:5] {
		boundary, parseErr := parseGLMBoundary(part.LessThan)
		if parseErr != nil || !boundary.Equal(previous.AddDate(0, 1, 0)) {
			return invalid("finite boundaries are not consecutive UTC month starts")
		}
		want := "p_" + previous.Format("200601")
		if part.Name != want || part.Ordinal != index+2 {
			return invalid("finite range name or ordinal")
		}
		previous = boundary
	}
	future := evidence.Partitions[5]
	if future.Name != "p_future" || future.LessThan != "MAXVALUE" || future.Ordinal != 6 {
		return invalid("future catch-all")
	}
	return nil
}

func parseGLMBoundary(text string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000000", text)
}

// sqlServerGLMPlacementQuery builds the partition-placement readback for the
// real fixture: row counts grouped by partition number of the owned function.
// Identifiers are bracket-quoted; no values are interpolated.
func sqlServerGLMPlacementQuery(schemaName, tableName, functionName string) (string, error) {
	quoted := func(value string) (string, error) {
		if strings.TrimSpace(value) == "" || strings.TrimSpace(value) != value {
			return "", fmt.Errorf("placement identifier is not canonical: %q", value)
		}
		return "[" + strings.ReplaceAll(value, "]", "]]") + "]", nil
	}
	schemaQuoted, err := quoted(schemaName)
	if err != nil {
		return "", err
	}
	tableQuoted, err := quoted(tableName)
	if err != nil {
		return "", err
	}
	functionQuoted, err := quoted(functionName)
	if err != nil {
		return "", err
	}
	return "SELECT $partition." + functionQuoted + "([bucket_start]) AS partition_number, COUNT(*) AS rows_in_partition" +
		" FROM " + schemaQuoted + "." + tableQuoted +
		" GROUP BY $partition." + functionQuoted + "([bucket_start]) ORDER BY partition_number", nil
}

// glmDestinationRef is the destination identity the GLM segments bind to.
type glmDestinationRef struct {
	ConnectorID, ConnectorRevision, Database, TableSchema, TableName string
}

func glmDestination(group *workspace.WriteGroup) glmDestinationRef {
	return glmDestinationRef{
		ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
		Database: group.Destination.Database, TableSchema: group.Destination.TableSchema,
		TableName: group.Destination.TableName,
	}
}

func glmSQLServerInspector(t *testing.T, f *sqlServerCircuitPlainFixture) *dbtarget.ManagedTableInspector {
	t.Helper()
	return dbtarget.NewManagedTableInspector(f.targetService, f.db)
}

// inspectGLMManagedTable reads the real SQL Server catalog for one managed
// table at the saved revisions. Missing owned tables surface as a missing
// status inspection, never as an error.
func inspectGLMManagedTable(ctx context.Context, t *testing.T, inspector *dbtarget.ManagedTableInspector, dest glmDestinationRef,
	capabilityRevision, table string) *dbtarget.SQLServerManagedInspection {
	t.Helper()
	managed, err := inspector.InspectSQLServerManagedTableAtRevision(ctx,
		dest.ConnectorID, dest.ConnectorRevision, dest.Database, dest.TableSchema, capabilityRevision, table)
	require.NoError(t, err)
	require.NotNil(t, managed)
	return managed
}

// glmConvertInspection mirrors the production API conversion from a managed
// inspection into the canonical target inspection, including the full native
// partition evidence.
func glmConvertInspection(managed *dbtarget.SQLServerManagedInspection) recordingplan.TargetTableInspection {
	state := recordingplan.TargetTableInspection{
		Status: string(managed.Inspection.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{},
		ColumnPrimaryKey: map[string]bool{}, IndexesInspected: managed.Indexes != nil,
	}
	for _, column := range managed.Inspection.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name] = column.DataType
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
	}
	for _, index := range managed.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	evidence := managed.Evidence
	state.SQLServer = &recordingplan.TargetSQLServerTableEvidence{
		ServerVersion: evidence.ServerVersion, ServerEdition: evidence.ServerEdition, CapabilityRevision: evidence.CapabilityRevision,
		Database: evidence.Database, Schema: evidence.Schema, EngineEdition: evidence.EngineEdition,
		CompatibilityLevel: evidence.CompatibilityLevel, MetadataVerified: evidence.MetadataVerified,
		TableType: evidence.TableType, ObjectID: evidence.ObjectID, Plain: evidence.Plain,
		PrimaryKey: evidence.PrimaryKey, ColumnCollations: evidence.ColumnCollations, ColumnDefaults: evidence.ColumnDefaults,
		ProvenanceJSONChecked: evidence.ProvenanceJSONChecked, RepairOperationID: evidence.RepairOperationID,
		NativeMonthly: glmTargetPartition(evidence.NativeMonthly), PartitionInspected: evidence.PartitionInspected,
	}
	return state
}

func glmTargetPartition(evidence *dbtarget.SQLServerManagedPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]recordingplan.SQLServerPartitionBoundary, len(evidence.Partitions))
	for i, partition := range evidence.Partitions {
		partitions[i] = recordingplan.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName, SchemeName: evidence.SchemeName,
		Filegroup: evidence.Filegroup, Partitions: partitions, FunctionExists: evidence.FunctionExists,
		SchemeExists: evidence.SchemeExists, Exclusive: evidence.Exclusive, IndexesAligned: evidence.IndexesAligned,
	}
}

func glmPartitionFromTarget(evidence *recordingplan.TargetSQLServerPartitionEvidence) *recordingplan.SQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]recordingplan.SQLServerPartitionBoundary, len(evidence.Partitions))
	copy(partitions, evidence.Partitions)
	return &recordingplan.SQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName,
		SchemeName: evidence.SchemeName, Filegroup: evidence.Filegroup, Partitions: partitions,
	}
}

// glmMaintenanceTarget binds the real catalog inspector and Luna's typed
// maintenance executor into the ledger-driven maintenance target.
func glmMaintenanceTarget(t *testing.T, inspector *dbtarget.ManagedTableInspector, dest glmDestinationRef,
	capabilityRevision string, target *sql.DB, testCapability *dbtarget.SQLServerRecordingCapability) recordingplan.SQLServerPartitionMaintenanceTarget {
	t.Helper()
	return recordingplan.SQLServerPartitionMaintenanceTarget{
		Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
			managed, err := inspector.InspectSQLServerManagedTableAtRevision(ctx,
				dest.ConnectorID, dest.ConnectorRevision, dest.Database, dest.TableSchema, capabilityRevision, table)
			if err != nil {
				return recordingplan.TargetTableInspection{}, fmt.Errorf("%w: %w", recordingplan.ErrTargetInspectionUnconfirmed, err)
			}
			return glmConvertInspection(managed), nil
		},
		Split: func(ctx context.Context, before, after *recordingplan.SQLServerPartitionEvidence) error {
			managed, err := inspector.InspectSQLServerManagedTableAtRevision(ctx,
				dest.ConnectorID, dest.ConnectorRevision, dest.Database, dest.TableSchema, capabilityRevision, dest.TableName)
			if err != nil {
				return fmt.Errorf("%w: %w", recordingplan.ErrTargetInspectionUnconfirmed, err)
			}
			afterProof := *after
			partitions := make([]dbtarget.SQLServerPartitionBoundary, len(afterProof.Partitions))
			for i, partition := range afterProof.Partitions {
				partitions[i] = dbtarget.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
			}
			return dbtarget.SplitOpenSQLServerNativeMonthly(ctx, target, testCapability, managed, &dbtarget.SQLServerManagedPartitionEvidence{
				Method: afterProof.Method, Column: afterProof.Column, FunctionName: afterProof.FunctionName,
				SchemeName: afterProof.SchemeName, Filegroup: afterProof.Filegroup, Partitions: partitions,
				FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
			})
		},
	}
}

// glmRestoreTarget binds the restore engine to the real catalog and Luna's
// repair batch executor.
func glmRestoreTarget(t *testing.T, inspector *dbtarget.ManagedTableInspector, dest glmDestinationRef,
	capabilityRevision string, target *sql.DB, testCapability *dbtarget.SQLServerRecordingCapability) recordingplan.StorageRestoreTarget {
	t.Helper()
	inspect := func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		managed, err := inspector.InspectSQLServerManagedTableAtRevision(ctx,
			dest.ConnectorID, dest.ConnectorRevision, dest.Database, dest.TableSchema, capabilityRevision, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		return glmConvertInspection(managed), nil
	}
	execute := func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
		return recordingplan.SchemaExecution{}, fmt.Errorf("generic SQL execution must not run for SQL Server repair")
	}
	return recordingplan.StorageRestoreTarget{
		Inspect: inspect, Execute: execute,
		ExecuteSQLServerRepair: func(ctx context.Context, tableName, operationID string, statements []string) (recordingplan.SchemaExecution, error) {
			execution, err := dbtarget.RunSQLServerManagedRepairBatch(ctx, target, testCapability,
				dbtarget.SQLServerManagedRepairRequest{TableName: tableName, OperationID: operationID, Statements: statements})
			if execution == nil {
				return recordingplan.SchemaExecution{}, err
			}
			return recordingplan.SchemaExecution{Committed: execution.Committed, RolledBack: execution.RolledBack}, err
		},
	}
}

func waitGLMOutboxSettled(t *testing.T, f *sqlServerCircuitPlainFixture, want int) {
	t.Helper()
	require.Eventually(t, func() bool {
		var total, unfinished int
		err := f.db.QueryRowContext(t.Context(), `
			SELECT COUNT(*), SUM(CASE WHEN state <> 'sql_committed' THEN 1 ELSE 0 END)
			FROM wg_delivery_outbox`).Scan(&total, &unfinished)
		return err == nil && total == want && unfinished == 0
	}, 20*time.Second, 25*time.Millisecond, "delivered rows must reach durable SQL receipts")
}

func glmReceiptDigests(t *testing.T, db *sql.DB, schemaName string) map[string]string {
	t.Helper()
	rows, err := db.QueryContext(t.Context(),
		"SELECT effect_key, payload_digest FROM "+sqlServerQualified(schemaName, dbtarget.EffectReceiptTable))
	require.NoError(t, err)
	defer rows.Close()
	digests := map[string]string{}
	for rows.Next() {
		var effectKey, digest string
		require.NoError(t, rows.Scan(&effectKey, &digest))
		digests[effectKey] = digest
	}
	require.NoError(t, rows.Err())
	return digests
}

func queryGLMPlacement(t *testing.T, db *sql.DB, schemaName, tableName, functionName string) map[int]int64 {
	t.Helper()
	query, err := sqlServerGLMPlacementQuery(schemaName, tableName, functionName)
	require.NoError(t, err)
	rows, err := db.QueryContext(t.Context(), query)
	require.NoError(t, err)
	defer rows.Close()
	placement := map[int]int64{}
	for rows.Next() {
		var number int
		var count int64
		require.NoError(t, rows.Scan(&number, &count))
		placement[number] = count
	}
	require.NoError(t, rows.Err())
	return placement
}

func glmNativeScope(group *workspace.WriteGroup) recordingplan.StorageScope {
	return recordingplan.StorageScope{
		WorkspaceID: group.WorkspaceID, ConnectorID: group.Destination.ConnectorID,
		Database: group.Destination.Database, TableSchema: group.Destination.TableSchema,
		StorageMode: string(group.StorageMode), LayoutVersion: string(group.LayoutVersion),
	}
}

func glmMaintenanceRequest(scope recordingplan.StorageScope, group *workspace.WriteGroup, connectorRevision string, at time.Time) recordingplan.SQLServerPartitionMaintenanceRequest {
	return recordingplan.SQLServerPartitionMaintenanceRequest{
		Scope: scope, GroupID: group.ID, ConnectorRevision: connectorRevision, At: at,
	}
}

func glmRestoreRequest(scope recordingplan.StorageScope, group *workspace.WriteGroup, connectorRevision string,
	groupInspection, receiptInspection recordingplan.TargetTableInspection, fallback ...*recordingplan.GroupSchemaLayout,
) recordingplan.StorageRestoreRequest {
	table, owner, ok := workspace.LongStorageIdentity(group)
	if !ok {
		table, owner = group.Destination.TableName, ""
	}
	return recordingplan.StorageRestoreRequest{
		Scope: scope, Dialect: "sqlserver", ConnectorRevision: connectorRevision,
		TableName: table, OwnerColumn: owner, GroupID: group.ID,
		FallbackLayout: func() *recordingplan.GroupSchemaLayout {
			if len(fallback) > 0 {
				return fallback[0]
			}
			return nil
		}(),
		SQLServer: &recordingplan.SQLServerStorageRestoreEvidence{Group: groupInspection, Receipt: receiptInspection},
	}
}

func insertGLMAdapterPlacement(t *testing.T, f *sqlServerCircuitPlainFixture, snap *workspace.WriteGroupAppliedSnapshot,
	group *workspace.WriteGroup, at time.Time) (string, error) {
	t.Helper()
	if snap == nil || snap.RuntimeLayout == nil || snap.RuntimeLayout.SQLServer == nil || group == nil || len(group.Members) == 0 {
		return "", fmt.Errorf("adapter placement requires a frozen SQL Server proof")
	}
	member := group.Members[0]
	layout, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: snap.RuntimeLayout.Columns,
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
		ProfileKind: string(group.RecordingProfile.Kind), ProfileInterval: group.RecordingProfile.Interval,
		ProfileUnit: string(group.RecordingProfile.Unit),
		Members: []dbtarget.LongRowMember{{
			MemberKey: glmMemberKey(member.DeviceID, member.PointID, member.TagID),
			DeviceID:  member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			Type: measurement.ExactDecimal, Reduction: "",
		}},
	})
	if len(issues) > 0 {
		return "", fmt.Errorf("adapter placement layout is not canonical")
	}
	value, err := measurement.NewDecimal("4.0")
	if err != nil {
		return "", err
	}
	recSum := sha256.Sum256([]byte(group.WorkspaceID + "|" + group.ID + "|record|" + at.UTC().Format(time.RFC3339Nano)))
	recordID := hex.EncodeToString(recSum[:])
	effSum := sha256.Sum256([]byte(group.WorkspaceID + "|" + group.ID + "|effect|" + at.UTC().Format(time.RFC3339Nano)))
	effectKey := hex.EncodeToString(effSum[:])
	outcome := snapshot.Outcome{
		Kind: snapshot.OutcomeRow, EntityKey: "adapter-late", RecordID: recordID, EffectKey: effectKey,
		BucketStart: at.UTC(), BucketEnd: at.UTC().Add(10 * time.Minute),
		Members: []snapshot.MemberResult{{
			MemberKey: glmMemberKey(member.DeviceID, member.PointID, member.TagID), Status: snapshot.MemberOK,
			Sample: &snapshot.Sample{SampleID: "adapter-late", MemberKey: glmMemberKey(member.DeviceID, member.PointID, member.TagID),
				ObservedAt: at.UTC(), ReceivedAt: at.UTC(), Quality: schema.QualityGood, Value: value},
		}},
	}
	row, err := layout.EncodeRow(outcome)
	if err != nil {
		return "", err
	}
	digest, err := groupdelivery.RowPayloadDigest(row)
	if err != nil {
		return "", err
	}
	_, err = dbtarget.InsertGroupRow(t.Context(), f.target, dbtarget.GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLServer, SchemaName: group.Destination.TableSchema, TableName: group.Destination.TableName,
		Row: row, PayloadDigest: digest, Strategy: dbtarget.GroupEffectReceipt,
		CommittedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z"), SQLServer: snap.RuntimeLayout.SQLServer,
	})
	if err != nil {
		return "", err
	}
	return effectKey, nil
}

func glmMemberKey(deviceID, pointID, tagID string) string {
	return deviceID + "|" + pointID + "|" + tagID
}
