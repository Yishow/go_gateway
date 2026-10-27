package workspace

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func validSQLServerNativeRuntimeFixture(t *testing.T) (*WriteGroup, *WriteGroupRuntimeLayout) {
	t.Helper()
	group := &WriteGroup{
		WorkspaceID: "ws-1", ID: "grp-native", Revision: "rev-1", AppliedRevision: "rev-1",
		LayoutVersion: WriteGroupLayoutManagedLongV1, StorageMode: WriteGroupStorageModeNativeMonthly,
		RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileInstantaneous, Interval: 15, Unit: WriteGroupIntervalUnitSecond},
		Members: []WriteGroupMember{
			{DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", TargetColumn: "val1"},
		},
		Destination: WriteGroupDestination{
			StorageStrategy: WriteGroupStorageStrategyManaged, ConnectorID: "conn-1", ConnectorRevision: "crev-1",
			Database: "analytics_db", TableSchema: "dbo", SchemaRevision: "srev-1",
		},
		RowPolicy:   WriteGroupRowPolicy{IntervalSeconds: 15, RecordKeyColumn: "record_id"},
		WritePolicy: WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeReceipt},
	}
	table, owner, ok := LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database: group.Destination.Database, Schema: group.Destination.TableSchema, Collation: "Latin1_General_100_BIN2",
		ServerVersion: "16.0.1000", ServerEdition: "Standard Edition", CapabilityRevision: "sqlservercap-native",
		EngineEdition: 2, CompatibilityLevel: 160,
	}
	native, err := recordingplan.NewSQLServerNativeMonthlyPartition(time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC),
		"pf_"+table, "ps_"+table)
	require.NoError(t, err)
	options.NativeMonthly = native
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	require.NoError(t, err)

	groupProof, groupState := sqlServerTestRuntimeCatalog(layout.Columns, layout.Indexes, table, owner, 201, options)
	groupProof.Evidence.NativeMonthly = sqlServerManagedPartitionEvidenceForTest(native)
	groupProof.Evidence.PartitionInspected = true
	groupProof.Evidence.Plain = false
	groupState.SQLServer.NativeMonthly = sqlServerTargetPartitionEvidenceForTest(native)
	groupState.SQLServer.PartitionInspected = true
	groupState.SQLServer.Plain = false

	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptState := sqlServerTestRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 202, options)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, options.Database, groupState, receiptState)
	require.NoError(t, err)
	group.Destination.SchemaDigest = digest
	frozen := &WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: groupProof.Inspection.Columns,
		TagTypes: map[string]schema.DataType{"tag-1": schema.DataTypeInt64}, ReceiptTableReady: true, SchemaDigest: digest,
		SQLServer: &dbtarget.SQLServerManagedRowProof{
			ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
			Group: groupProof, Receipt: receiptProof,
		},
	}
	return group, frozen
}

func sqlServerManagedPartitionEvidenceForTest(native *recordingplan.SQLServerPartitionEvidence) *dbtarget.SQLServerManagedPartitionEvidence {
	partitions := make([]dbtarget.SQLServerPartitionBoundary, len(native.Partitions))
	for i, partition := range native.Partitions {
		partitions[i] = dbtarget.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &dbtarget.SQLServerManagedPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: partitions, FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func sqlServerTargetPartitionEvidenceForTest(native *recordingplan.SQLServerPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: append([]recordingplan.SQLServerPartitionBoundary(nil), native.Partitions...),
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func TestSQLServerNativeRuntimeLayoutValidation(t *testing.T) {
	group, frozen := validSQLServerNativeRuntimeFixture(t)
	require.NoError(t, ValidateSQLServerRuntimeLayout(group, frozen))
	encoded, err := json.Marshal(frozen)
	require.NoError(t, err)
	var persisted WriteGroupRuntimeLayout
	require.NoError(t, json.Unmarshal(encoded, &persisted))
	require.NoError(t, ValidateSQLServerRuntimeLayout(group, &persisted))

	for name, mutate := range map[string]func(*WriteGroupRuntimeLayout){
		"missing native proof": func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.NativeMonthly = nil },
		"not exclusive":        func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.NativeMonthly.Exclusive = false },
		"indexes unaligned":    func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.NativeMonthly.IndexesAligned = false },
		"function missing":     func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.NativeMonthly.FunctionExists = false },
		"scheme missing":       func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.NativeMonthly.SchemeExists = false },
		"group plain true":     func(l *WriteGroupRuntimeLayout) { l.SQLServer.Group.Evidence.Plain = true },
		"receipt partition": func(l *WriteGroupRuntimeLayout) {
			l.SQLServer.Receipt.Evidence.NativeMonthly = &dbtarget.SQLServerManagedPartitionEvidence{}
		},
	} {
		t.Run(name, func(t *testing.T) {
			var candidate WriteGroupRuntimeLayout
			require.NoError(t, json.Unmarshal(encoded, &candidate))
			mutate(&candidate)
			require.Error(t, ValidateSQLServerRuntimeLayout(group, &candidate))
		})
	}
}
