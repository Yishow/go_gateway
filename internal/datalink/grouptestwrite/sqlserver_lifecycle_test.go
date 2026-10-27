package grouptestwrite

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// sqlServerTestWriteGroups is a test-only saved-scope adapter. Positive
// fixtures must explicitly opt into its fake capability; production uses the
// workspace adapter and its live SupportsGroupTestWrite proof.
type sqlServerTestWriteGroups struct {
	*fakeGroups
	layout                         *workspace.WriteGroupRuntimeLayout
	testOnlySupportsGroupTestWrite bool
}

func (g *sqlServerTestWriteGroups) SQLServerTestWriteLayout(_ context.Context, _, _ string) (*workspace.WriteGroupRuntimeLayout, error) {
	if !g.testOnlySupportsGroupTestWrite {
		return nil, workspace.ErrWriteGroupValidation
	}
	return g.layout, nil
}

type sqlServerTestWriteDestinations struct{}

func (sqlServerTestWriteDestinations) GetByID(_ context.Context, id string) (*schema.DatabaseConnector, error) {
	if id != "conn-1" {
		return nil, dbtarget.ErrConnectorNotFound
	}
	return &schema.DatabaseConnector{ID: id, Kind: schema.DatabaseConnectorKindSQLServer, IdentityRevision: "crev-1", Enabled: true}, nil
}

func (sqlServerTestWriteDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	return nil, errors.New("destination must not be opened during preview")
}

// TestSQLServerGroupProbeLifecycle is the task 4.1 behavioral RED gate: the
// public Preview path must accept a verified SQL Server managed-long layout.
func TestSQLServerGroupProbeLifecycle(t *testing.T) {
	group, layout, tagTypes := sqlServerTestWriteFixture(t)
	require.NoError(t, workspace.ValidateSQLServerRuntimeLayout(group, layout), "fixture must be a verified SQL Server proof")
	groups := &sqlServerTestWriteGroups{
		fakeGroups: &fakeGroups{group: group}, layout: layout, testOnlySupportsGroupTestWrite: true,
	}
	service := New(Dependencies{
		Groups: groups, Tags: fakeTags(tagTypes), Destinations: sqlServerTestWriteDestinations{},
		Ledger: recordingplan.NewService(recordingplan.NewMemoryRepository()),
	}, Config{})

	preview, err := service.Preview(t.Context(), Workspace{ID: group.WorkspaceID, Revision: testRevision}, group.ID)
	require.NoError(t, err, "SQL Server managed-long preview should use its verified saved-scope adapter")
	require.Equal(t, "sqlserver", preview.Target.Dialect)
	require.Equal(t, "gateway", preview.OwnerValue, "the operation owner must not replace the canonical table owner")
	require.True(t, strings.HasPrefix(preview.OwnerColumn, "_gw_owner_"))
	require.Len(t, preview.RecordIDs, len(group.Members))
	require.NotContains(t, preview.Cleanup, preview.OwnerColumn+" =")
}

func TestSQLServerNativeMonthlyPreviewAcceptsAuthorizedExtendedAttemptProof(t *testing.T) {
	group, frozen, attempt, tagTypes := sqlServerNativeTestWriteFixture(t)
	require.NoError(t, workspace.ValidateSQLServerRuntimeLayout(group, frozen), "the immutable initial proof remains canonical")
	require.Len(t, frozen.SQLServer.Group.Evidence.NativeMonthly.Partitions, 6)
	require.Greater(t, len(attempt.SQLServer.Group.Evidence.NativeMonthly.Partitions), len(frozen.SQLServer.Group.Evidence.NativeMonthly.Partitions))
	groups := &sqlServerTestWriteGroups{
		fakeGroups: &fakeGroups{group: group}, layout: attempt, testOnlySupportsGroupTestWrite: true,
	}
	service := New(Dependencies{
		Groups: groups, Tags: fakeTags(tagTypes), Destinations: sqlServerTestWriteDestinations{},
		Ledger: recordingplan.NewService(recordingplan.NewMemoryRepository()),
	}, Config{})

	preview, err := service.Preview(t.Context(), Workspace{ID: group.WorkspaceID, Revision: testRevision}, group.ID)
	require.NoError(t, err, "an authorized current partition overlay must not be revalidated as the immutable initial layout")
	require.Equal(t, "sqlserver", preview.Target.Dialect)
	require.Len(t, preview.RecordIDs, len(group.Members))
}

func sqlServerTestWriteFixture(t *testing.T) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout, map[string]schema.DataType) {
	t.Helper()
	group := &workspace.WriteGroup{
		ID: "group-1", WorkspaceID: testWorkspaceID, Revision: "rev-1", AppliedRevision: "rev-1",
		Status: workspace.WriteGroupStatusReady, LayoutVersion: workspace.WriteGroupLayoutManagedLongV1,
		StorageMode: workspace.WriteGroupStorageModePlain,
		RecordingProfile: &workspace.WriteGroupRecordingProfile{
			Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond,
		},
		Members: []workspace.WriteGroupMember{
			{DeviceID: "device-1", PointID: "point-1", TagID: "tag-int", TargetColumn: "value_int"},
			{DeviceID: "device-1", PointID: "point-2", TagID: "tag-text", TargetColumn: "value_text"},
		},
		Destination: workspace.WriteGroupDestination{
			StorageStrategy: workspace.WriteGroupStorageStrategyManaged,
			ConnectorID:     "conn-1", ConnectorRevision: "crev-1", Database: "analytics_db", TableSchema: "dbo",
			SchemaRevision: "schema-rev-1",
		},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 15, RecordKeyColumn: "record_id", ProvenanceColumn: "provenance"},
		WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeReceipt},
	}
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database: group.Destination.Database, Schema: group.Destination.TableSchema,
		Collation: "Latin1_General_100_BIN2", ServerVersion: "16.0.1000", ServerEdition: "Standard Edition",
		CapabilityRevision: "sqlserver-cap-1", EngineEdition: 2, CompatibilityLevel: 160,
	}
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	groupProof, groupState := sqlServerTestWriteCatalog(layout.Columns, layout.Indexes, table, owner, 101, options)
	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptState := sqlServerTestWriteCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 102, options)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, options.Database, groupState, receiptState)
	require.NoError(t, err)
	group.Destination.SchemaDigest = digest
	tagTypes := map[string]schema.DataType{"tag-int": schema.DataTypeInt64, "tag-text": schema.DataTypeString}
	runtimeLayout := &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: groupProof.Inspection.Columns, TagTypes: tagTypes,
		ReceiptTableReady: true, SchemaDigest: digest,
		SQLServer: &dbtarget.SQLServerManagedRowProof{
			ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
			Group: groupProof, Receipt: receiptProof,
		},
	}
	return group, runtimeLayout, tagTypes
}

func sqlServerNativeTestWriteFixture(t *testing.T) (group *workspace.WriteGroup, frozen, attemptLayout *workspace.WriteGroupRuntimeLayout, tagTypes map[string]schema.DataType) {
	t.Helper()
	group, _, tagTypes = sqlServerTestWriteFixture(t)
	group.StorageMode = workspace.WriteGroupStorageModeNativeMonthly
	tableName, ownerColumn, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = tableName
	functionName, schemeName, err := recordingplan.SQLServerNativeMonthlyObjectNames(group.Destination.TableName)
	require.NoError(t, err)
	initial, err := recordingplan.NewSQLServerNativeMonthlyPartition(
		time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC), functionName, schemeName,
	)
	require.NoError(t, err)
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database: group.Destination.Database, Schema: group.Destination.TableSchema,
		Collation: "Latin1_General_100_BIN2", ServerVersion: "16.0.1000", ServerEdition: "Standard Edition",
		CapabilityRevision: "sqlserver-cap-1", EngineEdition: 2, CompatibilityLevel: 160, NativeMonthly: initial,
	}
	schemaLayout, err := recordingplan.ManagedSQLServerLongGroupLayout(group.Destination.TableName, ownerColumn, options)
	require.NoError(t, err)
	groupProof, groupState := sqlServerTestWriteCatalog(schemaLayout.Columns, schemaLayout.Indexes,
		group.Destination.TableName, schemaLayout.OwnerColumn, 101, options)
	groupProof.Evidence.NativeMonthly = sqlServerTestWriteManagedPartitions(initial)
	groupProof.Evidence.PartitionInspected = true
	groupProof.Evidence.Plain = false
	groupState.SQLServer.NativeMonthly = sqlServerTestWriteTargetPartitions(initial)
	groupState.SQLServer.PartitionInspected = true
	groupState.SQLServer.Plain = false

	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptState := sqlServerTestWriteCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 102, options)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(schemaLayout, options.Database, groupState, receiptState)
	require.NoError(t, err)
	group.Destination.SchemaDigest = digest
	frozen = &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: groupProof.Inspection.Columns, TagTypes: tagTypes,
		ReceiptTableReady: true, SchemaDigest: digest,
		SQLServer: &dbtarget.SQLServerManagedRowProof{
			ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
			Group: groupProof, Receipt: receiptProof,
		},
	}
	require.NoError(t, workspace.ValidateSQLServerRuntimeLayout(group, frozen))
	expanded, err := recordingplan.PlanSQLServerNativeMonthlyExtension(initial,
		time.Date(2026, time.December, 20, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	attempt := *frozen
	attempt.SQLServer = dbtarget.CloneSQLServerManagedRowProof(frozen.SQLServer)
	groupProofCopy := *attempt.SQLServer.Group
	groupEvidenceCopy := *groupProofCopy.Evidence
	groupEvidenceCopy.NativeMonthly = sqlServerTestWriteManagedPartitions(expanded)
	groupProofCopy.Evidence = &groupEvidenceCopy
	attempt.SQLServer.Group = &groupProofCopy
	return group, frozen, &attempt, tagTypes
}

func sqlServerTestWriteManagedPartitions(evidence *recordingplan.SQLServerPartitionEvidence) *dbtarget.SQLServerManagedPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]dbtarget.SQLServerPartitionBoundary, len(evidence.Partitions))
	for index, partition := range evidence.Partitions {
		partitions[index] = dbtarget.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &dbtarget.SQLServerManagedPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName, SchemeName: evidence.SchemeName,
		Filegroup: evidence.Filegroup, Partitions: partitions,
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func sqlServerTestWriteTargetPartitions(evidence *recordingplan.SQLServerPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName, SchemeName: evidence.SchemeName,
		Filegroup: evidence.Filegroup, Partitions: slices.Clone(evidence.Partitions),
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func sqlServerTestWriteCatalog(
	columns []recordingplan.GroupSchemaColumn,
	indexes []recordingplan.GroupSchemaIndex,
	table, owner string,
	objectID int64,
	options recordingplan.SQLServerGroupSchemaOptions,
) (*dbtarget.SQLServerManagedInspection, recordingplan.TargetTableInspection) {
	evidence := &dbtarget.SQLServerManagedTableEvidence{
		Database: options.Database, Schema: options.Schema, ServerVersion: options.ServerVersion,
		ServerEdition: options.ServerEdition, CapabilityRevision: options.CapabilityRevision,
		EngineEdition: options.EngineEdition, CompatibilityLevel: options.CompatibilityLevel,
		MetadataVerified: true, TableType: "USER_TABLE", ObjectID: objectID, Plain: true,
		ColumnDefaults: map[string]*string{}, ColumnCollations: map[string]string{},
		ProvenanceJSONChecked: table != dbtarget.EffectReceiptTable,
	}
	inspection := &dbtarget.SQLServerManagedInspection{
		Inspection: &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Schema: options.Schema, Table: table, IndexesInspected: true},
		Evidence:   evidence, Indexes: []dbtarget.IndexInfo{},
	}
	state := recordingplan.TargetTableInspection{
		Status: "exists", ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{},
		ColumnPrimaryKey: map[string]bool{}, IndexesInspected: true,
	}
	ownerDefault := "gateway"
	for _, column := range columns {
		evidence.ColumnDefaults[column.Name] = nil
		if column.Name == owner {
			evidence.ColumnDefaults[column.Name] = &ownerDefault
		}
		if strings.HasPrefix(column.SQLType, "NVARCHAR(") || strings.HasPrefix(column.SQLType, "CHAR(") {
			evidence.ColumnCollations[column.Name] = options.Collation
		}
		if column.PrimaryKey {
			evidence.PrimaryKey = append(evidence.PrimaryKey, column.Name)
		}
		inspection.Inspection.Columns = append(inspection.Inspection.Columns, dbtarget.ColumnInfo{
			Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey,
			Unique: column.PrimaryKey, HasDefault: column.Name == owner,
		})
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.SQLType, column.Nullable, column.PrimaryKey
	}
	for _, index := range indexes {
		inspection.Indexes = append(inspection.Indexes, dbtarget.IndexInfo{Name: index.Name, Columns: index.Columns})
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	inspection.Inspection.Indexes = inspection.Indexes
	state.SQLServer = &recordingplan.TargetSQLServerTableEvidence{
		Database: evidence.Database, Schema: evidence.Schema, ServerVersion: evidence.ServerVersion, ServerEdition: evidence.ServerEdition,
		EngineEdition: evidence.EngineEdition, CompatibilityLevel: evidence.CompatibilityLevel,
		CapabilityRevision: evidence.CapabilityRevision, MetadataVerified: evidence.MetadataVerified,
		TableType: evidence.TableType, ObjectID: evidence.ObjectID, Plain: evidence.Plain, PrimaryKey: evidence.PrimaryKey,
		ColumnDefaults: evidence.ColumnDefaults, ColumnCollations: evidence.ColumnCollations,
		ProvenanceJSONChecked: evidence.ProvenanceJSONChecked,
	}
	return inspection, state
}
