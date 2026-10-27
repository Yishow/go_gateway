package grouppipeline

import (
	"context"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func sqlServerRuntimeSnapshot(t *testing.T) *workspace.WriteGroupAppliedSnapshot {
	t.Helper()
	group, frozen := longProof(t, workspace.WriteGroupStorageModePlain, nil)
	group.Destination.Database, group.Destination.TableSchema = "fixture_db", "recording"
	group.Destination.SchemaRevision = "confirmed-schema-revision"
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database: "fixture_db", Schema: "recording", ServerVersion: "16.0.1000.6", ServerEdition: "Developer Edition",
		EngineEdition: 3, CompatibilityLevel: 160, CapabilityRevision: "cap-1", Collation: "Latin1_General_100_BIN2",
	}
	managed, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	groupProof, groupState := sqlServerRuntimeCatalog(managed.Columns, managed.Indexes, table, owner, 101, options)
	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"}, {Name: "committed_at", SQLType: "DATETIME2(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"},
	}
	receiptProof, receiptState := sqlServerRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 102, options)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(managed, options.Database, groupState, receiptState)
	require.NoError(t, err)
	group.Destination.SchemaDigest = digest
	frozen.Dialect, frozen.Columns, frozen.SchemaDigest, frozen.ReceiptTableReady = dbtarget.SQLDialectSQLServer, groupProof.Inspection.Columns, digest, true
	frozen.SQLServer = &dbtarget.SQLServerManagedRowProof{
		ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision, Group: groupProof, Receipt: receiptProof,
	}
	return &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.AppliedRevision,
		EffectiveAt: at(0), Group: group, RuntimeLayout: frozen,
	}
}

func sqlServerRuntimeCatalog(columns []recordingplan.GroupSchemaColumn, indexes []recordingplan.GroupSchemaIndex, table, owner string, objectID int64,
	options recordingplan.SQLServerGroupSchemaOptions,
) (*dbtarget.SQLServerManagedInspection, recordingplan.TargetTableInspection) {
	e := &dbtarget.SQLServerManagedTableEvidence{
		Database: options.Database, Schema: options.Schema, ServerVersion: options.ServerVersion, ServerEdition: options.ServerEdition,
		EngineEdition: options.EngineEdition, CompatibilityLevel: options.CompatibilityLevel, CapabilityRevision: options.CapabilityRevision,
		MetadataVerified: true, TableType: "USER_TABLE", ObjectID: objectID, Plain: true,
		ColumnDefaults: map[string]*string{}, ColumnCollations: map[string]string{}, ProvenanceJSONChecked: table != dbtarget.EffectReceiptTable,
	}
	state := recordingplan.TargetTableInspection{
		Status: "exists", ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}, IndexesInspected: true,
	}
	physical := &dbtarget.SQLServerManagedInspection{
		Inspection: &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Schema: options.Schema, Table: table, IndexesInspected: true},
		Evidence:   e, Indexes: []dbtarget.IndexInfo{},
	}
	ownerDefault := "gateway"
	for _, col := range columns {
		e.ColumnDefaults[col.Name] = nil
		if col.Name == owner {
			e.ColumnDefaults[col.Name] = &ownerDefault
		}
		if strings.HasPrefix(col.SQLType, "NVARCHAR(") || strings.HasPrefix(col.SQLType, "CHAR(") {
			e.ColumnCollations[col.Name] = options.Collation
		}
		if col.PrimaryKey {
			e.PrimaryKey = append(e.PrimaryKey, col.Name)
		}
		physical.Inspection.Columns = append(physical.Inspection.Columns, dbtarget.ColumnInfo{
			Name: col.Name, DataType: col.SQLType, Nullable: col.Nullable, PrimaryKey: col.PrimaryKey, Unique: col.PrimaryKey, HasDefault: col.Name == owner,
		})
		state.Columns = append(state.Columns, col.Name)
		state.ColumnTypes[col.Name], state.ColumnNullable[col.Name], state.ColumnPrimaryKey[col.Name] = col.SQLType, col.Nullable, col.PrimaryKey
	}
	for _, idx := range indexes {
		physical.Indexes = append(physical.Indexes, dbtarget.IndexInfo{Name: idx.Name, Columns: idx.Columns})
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: idx.Name, Columns: idx.Columns})
	}
	physical.Inspection.Indexes = physical.Indexes
	state.SQLServer = &recordingplan.TargetSQLServerTableEvidence{
		Database: e.Database, Schema: e.Schema, ServerVersion: e.ServerVersion, ServerEdition: e.ServerEdition,
		EngineEdition: e.EngineEdition, CompatibilityLevel: e.CompatibilityLevel, CapabilityRevision: e.CapabilityRevision,
		MetadataVerified: e.MetadataVerified, TableType: e.TableType, ObjectID: e.ObjectID, Plain: e.Plain, PrimaryKey: e.PrimaryKey,
		ColumnDefaults: e.ColumnDefaults, ColumnCollations: e.ColumnCollations, ProvenanceJSONChecked: e.ProvenanceJSONChecked,
	}
	return physical, state
}

func sqlServerResolverFixture(t *testing.T, snap *workspace.WriteGroupAppliedSnapshot) (destinationResolver, groupdelivery.OutboxItem, *mysqlResolverDestinations) {
	t.Helper()
	resolver, item, opened := newMySQLResolverFixture(t, snap, nil)
	proof, e := snap.RuntimeLayout.SQLServer, snap.RuntimeLayout.SQLServer.Group.Evidence
	opened.opened.Kind = schema.DatabaseConnectorKindSQLServer
	// A synthetic future released projection exercises the downstream branch;
	// the actual backend projection still always disables SupportsManagedLong.
	opened.opened.SQLServer = &dbtarget.SQLServerRecordingCapability{
		ConnectorID: proof.ConnectorID, ConnectorRevision: proof.ConnectorRevision, Database: e.Database, Schema: e.Schema,
		ServerVersion: e.ServerVersion, ServerEdition: e.ServerEdition, EngineEdition: e.EngineEdition,
		CompatibilityLevel: e.CompatibilityLevel, CapabilityRevision: e.CapabilityRevision,
		Connected: true, MetadataStatus: "verified", SupportsManagedLong: true,
		Permissions: dbtarget.SQLServerRecordingPermissions{Verified: true, SchemaViewDefinition: true, CreateTable: true,
			AlterSchema: true, SelectSchema: true, InsertSchema: true, DeleteSchema: true},
	}
	resolver.loadRuntimeSnapshot = func(context.Context, groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
		return snap, nil
	}
	return resolver, item, opened
}
