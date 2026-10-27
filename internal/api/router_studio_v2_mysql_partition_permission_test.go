package api

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestMySQLPartitionPermissionLostContinuesIntegration injects the permission
// refusal over real catalog inspection, then proves real SQL/receipt delivery
// into the unchanged MAXVALUE partition. It is not a restricted-role test.
func TestMySQLPartitionPermissionLostContinuesIntegration(t *testing.T) {
	runMySQLPartitionPermissionLostContinuesIntegration(t, mysqlCircuitFamilyMySQL)
}

func runMySQLPartitionPermissionLostContinuesIntegration(t *testing.T, family string) {
	now := time.Now().UTC()
	anchor := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
		Family:          family,
		SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"},
		Inputs:          mysqlCircuitDecimalInputs, StorageMode: workspace.WriteGroupStorageModeNativeMonthly,
		BaseTime: anchor.AddDate(0, 6, 0), OpenFailures: 1_000_000,
	})
	f.accept(0)
	f.closeAfter(10 * time.Minute)
	f.accept(1)
	f.closeAfter(20 * time.Minute)
	requireEventuallyCumulativeRetry(t, f)
	require.NoError(t, f.pipe.Stop(5*time.Second))
	require.Len(t, mysqlPartitionPayloadDigests(t, f), 82)
	assertMySQLPartitionRemoteCounts(t, f, 0)
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.f.db))
	capability, err := targets.ProbeMySQLRecordingAtRevision(t.Context(), f.instant.group.Destination.ConnectorID,
		f.instant.group.Destination.ConnectorRevision)
	require.NoError(t, err)
	attempts := 0
	result, err := recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db)).MaintainMySQLPartitions(t.Context(),
		recordingplan.MySQLPartitionMaintenanceRequest{
			Scope: workspace.StorageScopeOf(f.instant.group), GroupID: f.instant.id,
			ConnectorRevision: f.instant.group.Destination.ConnectorRevision, At: f.base,
		}, recordingplan.MySQLPartitionMaintenanceTarget{
			Inspect: func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
				physical, err := dbtarget.InspectOpenMySQLManagedTable(ctx, f.f.target, capability, capability.Database, table)
				if err != nil {
					return recordingplan.TargetTableInspection{}, err
				}
				return mysqlPartitionTargetInspection(physical), nil
			},
			Reorganize: func(context.Context, *recordingplan.MySQLPartitionEvidence, *recordingplan.MySQLPartitionEvidence) error {
				attempts++
				return recordingplan.ErrTargetPermissionDenied
			},
		})
	require.NoError(t, err)
	require.Equal(t, 1, attempts)
	require.Equal(t, "partition_maintenance_blocked", result.Code)
	require.Equal(t, "permission_lost", result.Reason)
	require.False(t, result.Changed)
	require.Len(t, result.State.Current.Partitions, 6)
	lineage, err := recordingplan.ReadStorageLineage(t.Context(), f.f.db, result.State.Scope)
	require.NoError(t, err)
	require.Empty(t, lineage.RepairState, "maintenance refusal must not become a fatal storage-repair refusal")
	ready, err := f.f.groups.Readiness(t.Context(), f.instant.id)
	require.NoError(t, err)
	require.True(t, ready.SchemaReady, "%+v", ready.Issues)
	require.NotNil(t, ready.PartitionMaintenance)
	require.Equal(t, result.Code, ready.PartitionMaintenance.Code)

	// This fixture-only resolver keeps the verified original physical proof.
	// The production Sender and SQL transaction still perform catalog checks,
	// receipts, real inserts and local completion for every previously accepted row.
	rig := &mysqlRepairRig{targets: targets, store: f.store, frozen: f.instant.snapshot.RuntimeLayout,
		recovery: grouppipeline.NewStorageRecovery(f.store, recordingplan.NewService(recordingplan.NewSQLRepository(f.f.db)), targets)}
	sender := rig.sender()
	for key := range mysqlPartitionPayloadDigests(t, f) {
		require.Eventually(t, func() bool {
			delivered, err := sender.Deliver(t.Context(), key)
			return err == nil && delivered.State == groupdelivery.StateCommitted
		}, 10*time.Second, 20*time.Millisecond)
	}
	assertMySQLPartitionRemoteCounts(t, f, 82)
	var catchAll int
	require.NoError(t, f.f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `"+f.table+"` PARTITION (`p_future`)").Scan(&catchAll))
	require.Equal(t, 82, catchAll, "the verified catch-all remains writable even though maintenance was refused")
	assertMySQLCircuitDecimalJournal(t, f)
	rows := assertMySQLCircuitDecimalOutboxAndReceipts(t, f)
	assertMySQLCircuitDecimalSQLRows(t, f, rows)
}

func mysqlPartitionTargetInspection(physical *dbtarget.MySQLManagedInspection) recordingplan.TargetTableInspection {
	state := recordingplan.TargetTableInspection{Status: string(physical.Inspection.Status), IndexesInspected: physical.Indexes != nil,
		ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}}
	for _, column := range physical.Inspection.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.DataType, column.Nullable, column.PrimaryKey
	}
	for _, index := range physical.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	e := physical.Evidence
	state.MySQL = &recordingplan.TargetMySQLTableEvidence{
		ServerFamily: e.ServerFamily, ServerVersion: e.ServerVersion, CapabilityRevision: e.CapabilityRevision, Database: e.Database,
		Engine: e.Engine, Charset: e.Charset, Collation: e.Collation, TableType: e.TableType, PrimaryKey: e.PrimaryKey,
		ColumnCollations: e.ColumnCollations, ColumnDefaults: e.ColumnDefaults, JSONColumns: e.JSONColumns,
	}
	if e.Partitions != nil {
		state.MySQL.Partitions = &recordingplan.MySQLPartitionEvidence{Method: e.Partitions.Method, Expression: e.Partitions.Expression}
		for _, p := range e.Partitions.Partitions {
			state.MySQL.Partitions.Partitions = append(state.MySQL.Partitions.Partitions, recordingplan.MySQLPartitionBoundary{
				Name: p.Name, LessThan: p.LessThan, Ordinal: p.Ordinal})
		}
	}
	return state
}
