package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	datalinksnapshot "go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMySQLLongRuntimeProofRequiresFrozenPhysicalEvidence(t *testing.T) {
	snap := mysqlRuntimeSnapshot(t)
	require.NoError(t, workspace.ValidateMySQLRuntimeLayout(snap.Group, snap.RuntimeLayout))
	require.False(t, recordingReductionUnsupported(snap.Group, snap.RuntimeLayout.TagTypes))
	profile := snap.Group.RecordingProfile
	members := make([]dbtarget.LongRowMember, 0, len(snap.Group.Members))
	for _, member := range snap.Group.Members {
		tagType := snap.RuntimeLayout.TagTypes[member.TagID]
		stored, ok := workspace.LongMemberStoredType(profile.EffectiveSampleReduction(), member.ValueRole, tagType)
		require.True(t, ok)
		members = append(members, dbtarget.LongRowMember{
			MemberKey: member.DeviceID + "/" + member.PointID + "/" + member.TagID,
			DeviceID:  member.DeviceID, PointID: member.PointID, TagID: member.TagID, Type: stored,
			Reduction: string(profile.EffectiveSampleReduction()), ValueRole: string(member.ValueRole), SourceValueType: member.SourceValueType,
		})
	}
	_, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: snap.RuntimeLayout.Dialect, Columns: snap.RuntimeLayout.Columns, WorkspaceID: snap.Group.WorkspaceID,
		GroupID: snap.Group.ID, GroupRevision: snap.Group.Revision, ProfileKind: string(profile.Kind),
		ProfileInterval: profile.Interval, ProfileUnit: string(profile.Unit), Members: members,
	})
	require.Empty(t, issues)
	require.True(t, validLongRuntimeProof(snap.Group, snap.RuntimeLayout))

	missing := *snap.RuntimeLayout
	missing.MySQL = nil
	require.False(t, validLongRuntimeProof(snap.Group, &missing), "a MySQL long layout without frozen physical evidence cannot recover")

	wrongConnector := *snap.RuntimeLayout
	proof := *snap.RuntimeLayout.MySQL
	proof.ConnectorRevision = "edited-revision"
	wrongConnector.MySQL = &proof
	require.False(t, validLongRuntimeProof(snap.Group, &wrongConnector), "the proof stays bound to the applied connector revision")
}

func TestPipelineBuildAcceptsMySQLOnlyWithCanonicalFrozenLongProof(t *testing.T) {
	fixture := newFixture(t)
	fixture.dest.connector.Kind = schema.DatabaseConnectorKindMySQL
	pipeline := fixture.pipeline()
	snap := mysqlRuntimeSnapshot(t)
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
	frozen, err := pipeline.loadRuntimeSnapshot(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.NotNil(t, frozen)
	require.True(t, validLongRuntimeProof(frozen.Group, frozen.RuntimeLayout))
	require.True(t, sameAppliedSnapshot(snap, frozen))

	built, reason := pipeline.build(t.Context(), snap)
	require.Empty(t, reason)
	require.NotNil(t, built)
}

func TestPipelineDoesNotFallbackToMySQLForLegacyWideOrCustomGroups(t *testing.T) {
	t.Run("legacy wide", func(t *testing.T) {
		fixture := newFixture(t)
		fixture.dest.connector.Kind = schema.DatabaseConnectorKindMySQL
		pipeline := fixture.pipeline()
		legacy := snapshotOf(groupAt("rev-1"), at(0))

		built, reason := pipeline.build(t.Context(), legacy)

		require.Nil(t, built)
		require.Equal(t, reasonDestinationUnsupported, reason)
	})

	t.Run("custom long destination", func(t *testing.T) {
		fixture := newFixture(t)
		fixture.dest.connector.Kind = schema.DatabaseConnectorKindMySQL
		pipeline := fixture.pipeline()
		snap := mysqlRuntimeSnapshot(t)
		snap.Group.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyCustom
		require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))

		built, reason := pipeline.build(t.Context(), snap)

		require.Nil(t, built)
		require.Equal(t, reasonSnapshotUnavailable, reason)
	})
}

func TestMySQLDestinationResolverForwardsOnlyTheMatchingFrozenProof(t *testing.T) {
	snap := mysqlRuntimeSnapshot(t)
	capability := mysqlRuntimeCapability(snap.RuntimeLayout.MySQL)
	resolver, item, opened := newMySQLResolverFixture(t, snap, capability)
	var loaded groupdelivery.GroupKey
	resolver.loadRuntimeSnapshot = func(_ context.Context, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
		loaded = key
		return snap, nil
	}

	target, err := resolver.Resolve(t.Context(), item)

	require.NoError(t, err)
	require.Equal(t, item.Key, loaded, "load the exact frozen group revision named by the outbox row")
	require.Same(t, snap.RuntimeLayout.MySQL, target.MySQL)
	require.Equal(t, schema.DatabaseConnectorKindMySQL, target.Kind)
	require.False(t, capability.SupportsManagedLong, "the frozen proof does not open the full-capability gate")
	require.NoError(t, target.Close())
	require.Equal(t, 1, opened.closeCalls)
}

func TestMySQLDestinationResolverBlocksMissingOrDriftedProof(t *testing.T) {
	for _, test := range []struct {
		name    string
		mutate  func(*workspace.WriteGroupAppliedSnapshot, *dbtarget.MySQLRecordingCapability, *groupdelivery.OutboxItem)
		loadErr error
		missing bool
	}{
		{name: "missing frozen version", missing: true},
		{name: "runtime version read failed", loadErr: errors.New("store unavailable")},
		{name: "outbox revision differs", mutate: func(_ *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.MySQLRecordingCapability, item *groupdelivery.OutboxItem) {
			item.Key.GroupRevision = "older-revision"
		}},
		{name: "frozen table differs", mutate: func(snap *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.MySQLRecordingCapability, _ *groupdelivery.OutboxItem) {
			snap.Group.Destination.TableName += "_other"
		}},
		{name: "actual capability revision drift", mutate: func(_ *workspace.WriteGroupAppliedSnapshot, capability *dbtarget.MySQLRecordingCapability, _ *groupdelivery.OutboxItem) {
			capability.CapabilityRevision = "cap-2"
		}},
		{name: "actual connector revision drift", mutate: func(_ *workspace.WriteGroupAppliedSnapshot, capability *dbtarget.MySQLRecordingCapability, _ *groupdelivery.OutboxItem) {
			capability.ConnectorRevision = "edited-revision"
		}},
		{name: "metadata not verified", mutate: func(_ *workspace.WriteGroupAppliedSnapshot, capability *dbtarget.MySQLRecordingCapability, _ *groupdelivery.OutboxItem) {
			capability.MetadataStatus = "unconfirmed"
		}},
		{name: "non-receipt strategy", mutate: func(_ *workspace.WriteGroupAppliedSnapshot, _ *dbtarget.MySQLRecordingCapability, item *groupdelivery.OutboxItem) {
			item.Destination.DedupeCapability = groupdelivery.DedupeUniqueKey
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			snap := mysqlRuntimeSnapshot(t)
			capability := mysqlRuntimeCapability(snap.RuntimeLayout.MySQL)
			resolver, item, opened := newMySQLResolverFixture(t, snap, capability)
			if test.mutate != nil {
				test.mutate(snap, capability, &item)
			}
			resolver.loadRuntimeSnapshot = func(context.Context, groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
				if test.loadErr != nil {
					return nil, test.loadErr
				}
				if test.missing {
					return nil, nil
				}
				return snap, nil
			}

			_, err := resolver.Resolve(t.Context(), item)

			require.ErrorIs(t, err, groupdelivery.ErrTargetBlocked)
			require.Equal(t, 1, opened.closeCalls, "rejecting a target closes its opened connection")
		})
	}
}

func mysqlRuntimeSnapshot(t *testing.T) *workspace.WriteGroupAppliedSnapshot {
	t.Helper()
	group, frozen := longProof(t, workspace.WriteGroupStorageModePlain, nil)
	group.Destination.Database, group.Destination.TableSchema = "gwtest", "gwtest"
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	group.Destination.SchemaRevision = "confirmed-schema-revision"

	options := recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1",
		Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin",
	}
	managed, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	groupState := mysqlRuntimeCatalogState(managed.Columns, managed.Indexes, owner, "gwtest", options)
	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "VARCHAR(191)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "VARCHAR(64)"},
	}
	receiptState := mysqlRuntimeCatalogState(receiptColumns, nil, "_gw_effect_receipts_v1", "gwtest", options)
	digest, err := recordingplan.VerifyManagedMySQLCatalog(managed, "gwtest", groupState, receiptState)
	require.NoError(t, err)

	groupInspection := mysqlRuntimePhysicalInspection(groupState, table, "gwtest")
	receiptInspection := mysqlRuntimePhysicalInspection(receiptState, dbtarget.EffectReceiptTable, "gwtest")
	proof := &dbtarget.MySQLManagedRowProof{
		ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
		Group: groupInspection, Receipt: receiptInspection,
	}
	roundTripDigest, err := recordingplan.VerifyManagedMySQLCatalog(
		managed, "gwtest", mysqlRuntimeTargetState(groupInspection), mysqlRuntimeTargetState(receiptInspection),
	)
	require.NoError(t, err)
	require.Equal(t, digest, roundTripDigest)
	group.Destination.SchemaDigest = digest
	frozen.Dialect = dbtarget.SQLDialectMySQL
	frozen.Columns = groupInspection.Inspection.Columns
	frozen.SchemaDigest = digest
	frozen.MySQL = proof
	frozen.ReceiptTableReady = true
	require.Equal(t, workspace.WriteGroupStorageModePlain, group.StorageMode)
	require.Equal(t, workspace.WriteGroupStorageStrategyManaged, group.Destination.StorageStrategy)
	require.Equal(t, workspace.WriteGroupLayoutManagedLongV1, group.EffectiveLayoutVersion())
	require.Equal(t, workspace.WriteGroupProfileInstantaneous, group.RecordingProfile.Kind)
	require.Equal(t, group.RowPolicy.IntervalSeconds, group.RecordingProfile.IntervalSeconds())
	require.Equal(t, group.Destination.Database, group.Destination.TableSchema)
	require.NotEmpty(t, group.Destination.ConnectorID)
	require.NotEmpty(t, group.Destination.ConnectorRevision)
	require.Equal(t, group.Destination.ConnectorID, proof.ConnectorID)
	require.Equal(t, group.Destination.ConnectorRevision, proof.ConnectorRevision)
	require.Equal(t, groupdelivery.DedupeReceipt, group.WritePolicy.DedupeCapability)
	require.True(t, frozen.ReceiptTableReady)
	expectedTable, _, identityOK := workspace.LongStorageIdentity(group)
	require.True(t, identityOK)
	require.Equal(t, expectedTable, group.Destination.TableName)
	require.Len(t, group.Members, 1)
	require.Equal(t, dbtarget.SQLDialectMySQL, frozen.Dialect)
	require.Equal(t, group.Destination.SchemaDigest, frozen.SchemaDigest)
	require.Equal(t, group.Destination.SchemaDigest, frozen.SchemaDigest)
	require.Equal(t, dbtarget.EffectReceiptTable, proof.Receipt.Inspection.Table)
	require.Equal(t, group.Destination.Database, proof.Receipt.Inspection.Schema)
	require.Equal(t, proof.Group.Indexes, proof.Group.Inspection.Indexes)
	require.Equal(t, proof.Receipt.Indexes, proof.Receipt.Inspection.Indexes)
	require.Equal(t, "mysql", proof.Group.Evidence.ServerFamily)
	require.Equal(t, "8.4.13", proof.Group.Evidence.ServerVersion)
	require.True(t, longGroup(group))
	require.True(t, workspace.LongRecordingShapeSupported(group, frozen.TagTypes))
	require.Equal(t, group.Destination.TableName, proof.Group.Inspection.Table)
	require.Equal(t, group.Destination.Database, proof.Group.Inspection.Schema)
	require.NoError(t, workspace.ValidateMySQLRuntimeLayout(group, frozen))

	return &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.AppliedRevision,
		EffectiveAt: at(0), Group: group, RuntimeLayout: frozen,
	}
}

func mysqlRuntimeCatalogState(
	columns []recordingplan.GroupSchemaColumn,
	indexes []recordingplan.GroupSchemaIndex,
	owner, database string,
	options recordingplan.MySQLGroupSchemaOptions,
) recordingplan.TargetTableInspection {
	state := recordingplan.TargetTableInspection{
		Status: "exists", ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{},
		ColumnPrimaryKey: map[string]bool{}, IndexesInspected: true,
		MySQL: &recordingplan.TargetMySQLTableEvidence{
			ServerFamily: options.ServerFamily, ServerVersion: options.ServerVersion,
			CapabilityRevision: options.CapabilityRevision, Database: database, Engine: options.Engine,
			Charset: options.Charset, Collation: options.Collation, TableType: "BASE TABLE",
			ColumnCollations: map[string]string{}, ColumnDefaults: map[string]*string{}, JSONColumns: map[string]bool{},
		},
	}
	for _, column := range columns {
		state.Columns = append(state.Columns, column.Name)
		typeName := strings.ToLower(column.SQLType)
		if column.Name == "value_bool" {
			typeName = "tinyint(1)"
		}
		state.ColumnTypes[column.Name] = typeName
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
		if column.PrimaryKey {
			state.MySQL.PrimaryKey = append(state.MySQL.PrimaryKey, column.Name)
		}
		upperType := strings.ToUpper(column.SQLType)
		if strings.HasPrefix(upperType, "VARCHAR(") || strings.HasPrefix(upperType, "CHAR(") || upperType == "LONGTEXT" {
			state.MySQL.ColumnCollations[column.Name] = options.Collation
		}
		if column.Name == "provenance" {
			state.MySQL.JSONColumns[column.Name] = true
		}
	}
	ownerDefault := "gateway"
	state.MySQL.ColumnDefaults[owner] = &ownerDefault
	for _, index := range indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	return state
}

func mysqlRuntimePhysicalInspection(state recordingplan.TargetTableInspection, table, database string) *dbtarget.MySQLManagedInspection {
	columns := make([]dbtarget.ColumnInfo, 0, len(state.Columns))
	for _, name := range state.Columns {
		defaultValue := state.MySQL.ColumnDefaults[name]
		columns = append(columns, dbtarget.ColumnInfo{
			Name: name, DataType: state.ColumnTypes[name], Nullable: state.ColumnNullable[name],
			PrimaryKey: state.ColumnPrimaryKey[name], HasDefault: defaultValue != nil,
		})
	}
	indexes := make([]dbtarget.IndexInfo, 0, len(state.Indexes))
	for _, index := range state.Indexes {
		indexes = append(indexes, dbtarget.IndexInfo{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	evidence := state.MySQL
	return &dbtarget.MySQLManagedInspection{
		Inspection: &dbtarget.TableInspection{
			Status: dbtarget.TableInspectionExists, Schema: database, Table: table, Columns: columns,
			Indexes: indexes, IndexesInspected: true,
		},
		Evidence: &dbtarget.MySQLManagedTableEvidence{
			ServerFamily: evidence.ServerFamily, ServerVersion: evidence.ServerVersion, CapabilityRevision: evidence.CapabilityRevision,
			Database: evidence.Database, Engine: evidence.Engine, Charset: evidence.Charset, Collation: evidence.Collation,
			TableType: evidence.TableType, PrimaryKey: evidence.PrimaryKey,
			ColumnCollations: evidence.ColumnCollations, ColumnDefaults: evidence.ColumnDefaults, JSONColumns: evidence.JSONColumns,
		},
		Indexes: indexes,
	}
}

func mysqlRuntimeTargetState(physical *dbtarget.MySQLManagedInspection) recordingplan.TargetTableInspection {
	inspection, evidence := physical.Inspection, physical.Evidence
	state := recordingplan.TargetTableInspection{
		Status: string(inspection.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{},
		ColumnPrimaryKey: map[string]bool{}, IndexesInspected: inspection.IndexesInspected,
		MySQL: &recordingplan.TargetMySQLTableEvidence{
			ServerFamily: evidence.ServerFamily, ServerVersion: evidence.ServerVersion,
			CapabilityRevision: evidence.CapabilityRevision, Database: evidence.Database, Engine: evidence.Engine,
			Charset: evidence.Charset, Collation: evidence.Collation, TableType: evidence.TableType,
			PrimaryKey: evidence.PrimaryKey, ColumnCollations: evidence.ColumnCollations,
			ColumnDefaults: evidence.ColumnDefaults, JSONColumns: evidence.JSONColumns,
		},
	}
	for _, column := range inspection.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name] = column.DataType
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
	}
	for _, index := range inspection.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	return state
}

func mysqlRuntimeCapability(proof *dbtarget.MySQLManagedRowProof) *dbtarget.MySQLRecordingCapability {
	evidence := proof.Group.Evidence
	return &dbtarget.MySQLRecordingCapability{
		ConnectorID: proof.ConnectorID, ConnectorRevision: proof.ConnectorRevision, Database: evidence.Database,
		ServerFamily: evidence.ServerFamily, ServerVersion: evidence.ServerVersion, CapabilityRevision: evidence.CapabilityRevision,
		Connected: true, MetadataStatus: "verified", TransactionalEngine: true,
	}
}

type mysqlResolverDestinations struct {
	opened     *dbtarget.OpenedDestination
	openErr    error
	closeCalls int
}

func (d *mysqlResolverDestinations) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return &schema.DatabaseConnector{ID: "connector-1", Kind: schema.DatabaseConnectorKindMySQL, IdentityRevision: "crev-1", Enabled: true}, nil
}

func (d *mysqlResolverDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	if d.openErr != nil {
		return nil, d.openErr
	}
	return d.opened, nil
}

func newMySQLResolverFixture(
	t *testing.T,
	snap *workspace.WriteGroupAppliedSnapshot,
	capability *dbtarget.MySQLRecordingCapability,
) (destinationResolver, groupdelivery.OutboxItem, *mysqlResolverDestinations) {
	t.Helper()
	database, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = database.Close() })
	opened := &mysqlResolverDestinations{}
	opened.opened = &dbtarget.OpenedDestination{
		DB: database, Kind: schema.DatabaseConnectorKindMySQL, MySQL: capability,
		Close: func() error {
			opened.closeCalls++
			return database.Close()
		},
	}
	item := mysqlOutboxItem(t, snap)
	return destinationResolver{destinations: opened}, item, opened
}

func mysqlOutboxItem(t *testing.T, snap *workspace.WriteGroupAppliedSnapshot) groupdelivery.OutboxItem {
	t.Helper()
	group := snap.Group
	destination := group.Destination
	scope, err := datalinksnapshot.DestinationScope(
		destination.ConnectorID, destination.ConnectorRevision, destination.Database, destination.TableSchema, destination.TableName,
	)
	require.NoError(t, err)
	return groupdelivery.OutboxItem{
		Key: groupKey(snap),
		Destination: groupdelivery.FrozenDestination{
			Scope: scope, ConnectorID: destination.ConnectorID, ConnectorRevision: destination.ConnectorRevision,
			Database: destination.Database, TableSchema: destination.TableSchema, TableName: destination.TableName,
			DedupeCapability: group.WritePolicy.DedupeCapability, RecordKeyColumn: group.RowPolicy.RecordKeyColumn,
		},
	}
}
