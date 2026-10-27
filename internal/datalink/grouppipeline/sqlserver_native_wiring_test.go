package grouppipeline

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func sqlServerNativeRuntimeSnapshot(t *testing.T) *workspace.WriteGroupAppliedSnapshot {
	t.Helper()
	group := &workspace.WriteGroup{
		ID: "group-native", WorkspaceID: "workspace-native", Revision: "rev-1", AppliedRevision: "rev-1",
		LayoutVersion: workspace.WriteGroupLayoutManagedLongV1, StorageMode: workspace.WriteGroupStorageModeNativeMonthly,
		RecordingProfile: &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond},
		Members:          []workspace.WriteGroupMember{{DeviceID: "device-1", PointID: "point-1", TagID: "tag-1", TargetColumn: "v"}},
		Destination: workspace.WriteGroupDestination{
			StorageStrategy: workspace.WriteGroupStorageStrategyManaged, ConnectorID: "connector-native", ConnectorRevision: "crev-1",
			Database: "fixture_db", TableSchema: "dbo", SchemaRevision: "srev-1",
		},
		RowPolicy:   workspace.WriteGroupRowPolicy{IntervalSeconds: 15, RecordKeyColumn: "record_id"},
		WritePolicy: workspace.WriteGroupWritePolicy{Mode: "append", DedupeCapability: groupdelivery.DedupeReceipt},
	}
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database: group.Destination.Database, Schema: group.Destination.TableSchema, Collation: "Latin1_General_100_BIN2",
		ServerVersion: "16.0.1000.6", ServerEdition: "Developer Edition", EngineEdition: 3,
		CompatibilityLevel: 160, CapabilityRevision: "cap-native",
	}
	native, err := recordingplan.NewSQLServerNativeMonthlyPartition(time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC),
		"pf_"+table, "ps_"+table)
	require.NoError(t, err)
	options.NativeMonthly = native
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	groupProof, groupState := sqlServerRuntimeCatalog(layout.Columns, layout.Indexes, table, owner, 301, options)
	groupProof.Evidence.NativeMonthly = sqlServerManagedPartitionProof(native)
	groupProof.Evidence.PartitionInspected = true
	groupProof.Evidence.Plain = false
	groupState.SQLServer.NativeMonthly = sqlServerTargetPartitionProof(native)
	groupState.SQLServer.PartitionInspected = true
	groupState.SQLServer.Plain = false

	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"}, {Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptState := sqlServerRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 302, options)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, options.Database, groupState, receiptState)
	require.NoError(t, err)
	group.Destination.SchemaDigest = digest
	frozen := &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: groupProof.Inspection.Columns,
		TagTypes: map[string]schema.DataType{"tag-1": schema.DataTypeInt64}, ReceiptTableReady: true, SchemaDigest: digest,
		SQLServer: &dbtarget.SQLServerManagedRowProof{
			ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
			Group: groupProof, Receipt: receiptProof,
		},
	}
	return &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.AppliedRevision,
		EffectiveAt: at(0), Group: group, RuntimeLayout: frozen,
	}
}

func sqlServerManagedPartitionProof(native *recordingplan.SQLServerPartitionEvidence) *dbtarget.SQLServerManagedPartitionEvidence {
	partitions := make([]dbtarget.SQLServerPartitionBoundary, len(native.Partitions))
	for i, partition := range native.Partitions {
		partitions[i] = dbtarget.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &dbtarget.SQLServerManagedPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: partitions, FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func sqlServerTargetPartitionProof(native *recordingplan.SQLServerPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: append([]recordingplan.SQLServerPartitionBoundary(nil), native.Partitions...),
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func TestSQLServerNativeRuntimeProof(t *testing.T) {
	snap := sqlServerNativeRuntimeSnapshot(t)
	require.True(t, validLongRuntimeProof(snap.Group, snap.RuntimeLayout))
	require.True(t, validStorageModeProof(snap.Group, snap.RuntimeLayout))

	encoded, err := json.Marshal(snap)
	require.NoError(t, err)
	var decoded workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.True(t, validLongRuntimeProof(decoded.Group, decoded.RuntimeLayout))

	decoded.RuntimeLayout.SQLServer.Group.Evidence.NativeMonthly.Exclusive = false
	require.False(t, validLongRuntimeProof(decoded.Group, decoded.RuntimeLayout))

	plain := sqlServerRuntimeSnapshot(t)
	require.True(t, validLongRuntimeProof(plain.Group, plain.RuntimeLayout))
	require.True(t, validStorageModeProof(plain.Group, plain.RuntimeLayout))
}
