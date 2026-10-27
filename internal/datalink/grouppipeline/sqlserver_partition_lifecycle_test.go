package grouppipeline

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionMaintenanceUsesSharedWorkerAndStops(t *testing.T) {
	f := newFixture(t)
	groups := &mysqlMaintenanceGroups{fakeGroups: f.groups}
	destinations := &mysqlMaintenanceDestinations{fakeDestinations: f.dest, entered: make(chan struct{})}
	pipeline := New(Dependencies{
		Groups: groups, Tags: f.tags, Destinations: destinations, Inspector: f.inspect, Store: f.store,
	}, Config{
		NodeID: "node-1", Owner: "node-1/sqlserver-maintenance", Now: func() time.Time { return f.clock },
		TickInterval: 5 * time.Millisecond, ReconcileInterval: 5 * time.Millisecond,
		PartitionMaintenanceInterval: 5 * time.Millisecond,
	})
	snap := sqlServerNativeRuntimeSnapshot(t)
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
	groups.groups = []*workspace.WriteGroup{liveGroup(snap.Group, workspace.WriteGroupStatusReady)}
	groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{snap.Group.ID: {snap}}

	runCtx, cancelRun := context.WithCancel(t.Context())
	defer cancelRun()
	require.NoError(t, pipeline.Start(runCtx))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), 5*time.Second)
		defer cancel()
		require.NoError(t, pipeline.StopContext(ctx))
	})
	select {
	case <-destinations.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("shared maintenance worker never reached the SQL Server destination")
	}
	require.Equal(t, int32(1), destinations.openCalls.Load())
	listedAtEntry := groups.listCalls.Load()
	require.Eventually(t, func() bool { return groups.listCalls.Load() > listedAtEntry+3 },
		5*time.Second, 5*time.Millisecond, "maintenance must not block ordinary reconciliation")
	require.Equal(t, int32(1), destinations.openCalls.Load(), "one shared maintenance flight")

	stopCtx, cancelStop := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancelStop()
	require.NoError(t, pipeline.StopContext(stopCtx))
	require.NoError(t, pipeline.WaitStopped(stopCtx))
	require.Equal(t, int32(1), destinations.openCalls.Load(), "no destination open after shutdown")
	var operations int
	require.NoError(t, f.store.DB().QueryRowContext(t.Context(),
		`SELECT COUNT(*) FROM managed_schema_operations WHERE action = ?`, recordingplan.SQLServerPartitionMaintenanceAction).Scan(&operations))
	require.Zero(t, operations, "a blocked connection cannot authorize maintenance DDL")
}
