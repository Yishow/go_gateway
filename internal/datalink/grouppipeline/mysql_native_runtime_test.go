package grouppipeline

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// mysqlNativePhysicalInspection mirrors a catalog state into the frozen MySQL
// managed inspection, carrying the table's native partition evidence.
func mysqlNativePhysicalInspection(
	state recordingplan.TargetTableInspection, table, database string,
	partitions *dbtarget.MySQLPartitionEvidence,
) *dbtarget.MySQLManagedInspection {
	physical := mysqlRuntimePhysicalInspection(state, table, database)
	physical.Evidence.Partitions = partitions
	return physical
}

// mysqlNativeTargetState mirrors frozen physical evidence back into the
// recordingplan inspection shape, keeping the partition evidence so the
// verified digest covers it.
func mysqlNativeTargetState(physical *dbtarget.MySQLManagedInspection) recordingplan.TargetTableInspection {
	state := mysqlRuntimeTargetState(physical)
	if physical.Evidence != nil {
		state.MySQL.Partitions = mirrorMySQLPartitionEvidence(physical.Evidence.Partitions)
	}
	return state
}

// mysqlNativeDBPartitions mirrors the confirmed recordingplan native evidence
// into the dbtarget evidence shape the frozen physical inspection carries.
func mysqlNativeDBPartitions(native *recordingplan.MySQLPartitionEvidence) *dbtarget.MySQLPartitionEvidence {
	boundaries := make([]dbtarget.MySQLPartitionBoundary, 0, len(native.Partitions))
	for _, boundary := range native.Partitions {
		boundaries = append(boundaries, dbtarget.MySQLPartitionBoundary{
			Name: boundary.Name, LessThan: boundary.LessThan, Ordinal: boundary.Ordinal,
		})
	}
	return &dbtarget.MySQLPartitionEvidence{Method: native.Method, Expression: native.Expression, Partitions: boundaries}
}

func mysqlEvidenceWithPartitions(
	evidence *dbtarget.MySQLManagedTableEvidence, partitions *dbtarget.MySQLPartitionEvidence,
) *dbtarget.MySQLManagedTableEvidence {
	copied := *evidence
	copied.Partitions = partitions
	return &copied
}

// mysqlNativeMonthlyFixture builds the frozen snapshot an Apply of a
// native_monthly long group leaves behind: canonical RP options with the
// native partition evidence, and physical catalog evidence that mirrors it —
// partitions on the group table, an unpartitioned shared receipt.
func mysqlNativeMonthlyFixture(t *testing.T) (*workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
	t.Helper()
	group, frozen := longProof(t, workspace.WriteGroupStorageModeNativeMonthly, nil)
	group.Destination.Database, group.Destination.TableSchema = "gwnative", "gwnative"
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	group.Destination.SchemaRevision = "confirmed-schema-revision"

	native, err := recordingplan.NewMySQLNativeMonthlyPartition(time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	options := recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1",
		Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin",
		NativeMonthly: native,
	}
	managed, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	groupState := mysqlRuntimeCatalogState(managed.Columns, managed.Indexes, owner, "gwnative", options)
	groupState.MySQL.Partitions = native
	receiptColumns := []recordingplan.GroupSchemaColumn{
		{Name: "effect_key", SQLType: "VARCHAR(191)", PrimaryKey: true},
		{Name: "payload_digest", SQLType: "CHAR(64)"},
		{Name: "committed_at", SQLType: "DATETIME(6)"},
		{Name: "_gw_effect_receipts_v1", SQLType: "VARCHAR(64)"},
	}
	receiptState := mysqlRuntimeCatalogState(receiptColumns, nil, "_gw_effect_receipts_v1", "gwnative", options)
	digest, err := recordingplan.VerifyManagedMySQLCatalog(managed, "gwnative", groupState, receiptState)
	require.NoError(t, err)

	groupInspection := mysqlNativePhysicalInspection(groupState, table, "gwnative", mysqlNativeDBPartitions(native))
	receiptInspection := mysqlNativePhysicalInspection(receiptState, dbtarget.EffectReceiptTable, "gwnative", nil)
	proof := &dbtarget.MySQLManagedRowProof{
		ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
		Group: groupInspection, Receipt: receiptInspection,
	}
	roundTripDigest, err := recordingplan.VerifyManagedMySQLCatalog(
		managed, "gwnative", mysqlNativeTargetState(groupInspection), mysqlNativeTargetState(receiptInspection),
	)
	require.NoError(t, err)
	require.Equal(t, digest, roundTripDigest)
	group.Destination.SchemaDigest = digest
	frozen.Dialect = dbtarget.SQLDialectMySQL
	frozen.Columns = groupInspection.Inspection.Columns
	frozen.SchemaDigest = digest
	frozen.MySQL = proof
	frozen.ReceiptTableReady = true
	require.NoError(t, workspace.ValidateMySQLRuntimeLayout(group, frozen),
		"the canonical validator is the authority and accepts this native proof")
	return group, frozen
}

func TestMySQLNativeMonthlyStorageModeProof(t *testing.T) {
	group, frozen := mysqlNativeMonthlyFixture(t)
	require.True(t, validStorageModeProof(group, frozen))

	t.Run("missing group partition evidence", func(t *testing.T) {
		broken := *frozen
		broken.MySQL = &dbtarget.MySQLManagedRowProof{
			ConnectorID: frozen.MySQL.ConnectorID, ConnectorRevision: frozen.MySQL.ConnectorRevision,
			Group: &dbtarget.MySQLManagedInspection{
				Inspection: frozen.MySQL.Group.Inspection,
				Evidence:   mysqlEvidenceWithPartitions(frozen.MySQL.Group.Evidence, nil),
				Indexes:    frozen.MySQL.Group.Indexes,
			},
			Receipt: frozen.MySQL.Receipt,
		}
		require.False(t, validStorageModeProof(group, &broken))
	})
	t.Run("receipt must stay unpartitioned", func(t *testing.T) {
		native, err := recordingplan.NewMySQLNativeMonthlyPartition(time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		broken := *frozen
		broken.MySQL = &dbtarget.MySQLManagedRowProof{
			ConnectorID: frozen.MySQL.ConnectorID, ConnectorRevision: frozen.MySQL.ConnectorRevision,
			Group: frozen.MySQL.Group,
			Receipt: &dbtarget.MySQLManagedInspection{
				Inspection: frozen.MySQL.Receipt.Inspection,
				Evidence:   mysqlEvidenceWithPartitions(frozen.MySQL.Receipt.Evidence, mysqlNativeDBPartitions(native)),
				Indexes:    frozen.MySQL.Receipt.Indexes,
			},
		}
		require.False(t, validStorageModeProof(group, &broken))
	})
	t.Run("no hypertable claim", func(t *testing.T) {
		broken := *frozen
		broken.Hypertable = &workspace.WriteGroupHypertableEvidence{
			TimeColumn: "bucket_start", ChunkIntervalMicros: 7 * 24 * 3600 * 1000000,
		}
		require.False(t, validStorageModeProof(group, &broken))
	})
	t.Run("mysql dialect only", func(t *testing.T) {
		broken := *frozen
		broken.Dialect = dbtarget.SQLDialectPostgres
		require.False(t, validStorageModeProof(group, &broken))
	})
}

func TestMySQLPlainRuntimeProofRejectsPartitionEvidence(t *testing.T) {
	snap := mysqlRuntimeSnapshot(t)
	group, frozen := snap.Group, snap.RuntimeLayout
	require.Equal(t, workspace.WriteGroupStorageModePlain, group.StorageMode)
	require.True(t, validStorageModeProof(group, frozen))

	native, err := recordingplan.NewMySQLNativeMonthlyPartition(time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	broken := *frozen
	broken.MySQL = &dbtarget.MySQLManagedRowProof{
		ConnectorID: frozen.MySQL.ConnectorID, ConnectorRevision: frozen.MySQL.ConnectorRevision,
		Group: &dbtarget.MySQLManagedInspection{
			Inspection: frozen.MySQL.Group.Inspection,
			Evidence:   mysqlEvidenceWithPartitions(frozen.MySQL.Group.Evidence, mysqlNativeDBPartitions(native)),
			Indexes:    frozen.MySQL.Group.Indexes,
		},
		Receipt: frozen.MySQL.Receipt,
	}
	require.False(t, validStorageModeProof(group, &broken),
		"plain storage never carries native partition evidence")
	require.False(t, validLongRuntimeProof(group, &broken))
}

func TestPlainPostgresStorageProofUnchanged(t *testing.T) {
	group, frozen := longProof(t, workspace.WriteGroupStorageModePlain, nil)
	require.True(t, validStorageModeProof(group, frozen),
		"the partition checks must not disturb non-MySQL plain storage")
}

// TestMySQLNativeMonthlyResolvesOnlyThroughCanonicalValidator pins the
// resolver's acceptance to the canonical validator's verdict: native_monthly
// never resolves on the pipeline's own checks.
func TestMySQLNativeMonthlyResolvesOnlyThroughCanonicalValidator(t *testing.T) {
	group, frozen := mysqlNativeMonthlyFixture(t)
	canonical := workspace.ValidateMySQLRuntimeLayout(group, frozen) == nil
	require.True(t, canonical)
	require.Equal(t, canonical, validLongRuntimeProof(group, frozen))

	snap := &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
		AppliedRevision: group.AppliedRevision, EffectiveAt: at(0), Group: group, RuntimeLayout: frozen,
	}
	item := mysqlOutboxItem(t, snap)
	require.Equal(t, canonical, mysqlFrozenScopeMatches(item, snap))
}
