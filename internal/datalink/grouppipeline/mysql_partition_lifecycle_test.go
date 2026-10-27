package grouppipeline

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// mysqlMaintenanceGroups counts List calls so the test can prove the shared
// select loop keeps reconciling while a maintenance flight is blocked.
type mysqlMaintenanceGroups struct {
	*fakeGroups
	listCalls atomic.Int32
}

func (g *mysqlMaintenanceGroups) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	g.listCalls.Add(1)
	return g.fakeGroups.List(ctx)
}

// mysqlMaintenanceDestinations opens nothing: the first call signals entry and
// blocks on the caller's context, the way a hung destination would. Any second
// call is a single-flight violation and fails the pass loudly.
type mysqlMaintenanceDestinations struct {
	*fakeDestinations
	openCalls atomic.Int32
	entered   chan struct{}
	enterOnce sync.Once
}

func (d *mysqlMaintenanceDestinations) OpenDestination(
	ctx context.Context, connectorID, expectedRevision string,
) (*dbtarget.OpenedDestination, error) {
	if d.openCalls.Add(1) != 1 {
		return nil, errors.New("unexpected second OpenDestination during one maintenance flight")
	}
	d.enterOnce.Do(func() { close(d.entered) })
	<-ctx.Done()
	return nil, ctx.Err()
}

func newMySQLMaintenanceFixture(t *testing.T) (*fixture, *mysqlMaintenanceGroups, *mysqlMaintenanceDestinations) {
	t.Helper()
	f := newFixture(t)
	groups := &mysqlMaintenanceGroups{fakeGroups: f.groups}
	destinations := &mysqlMaintenanceDestinations{fakeDestinations: f.dest, entered: make(chan struct{})}
	return f, groups, destinations
}

// mysqlNativeMaintenanceSnapshot freezes the applied native snapshot an Apply
// leaves behind: the canonical native proof saved as the runtime version.
func mysqlNativeMaintenanceSnapshot(t *testing.T, p *Pipeline) *workspace.WriteGroupAppliedSnapshot {
	t.Helper()
	group, frozen := mysqlNativeMonthlyFixture(t)
	snap := &workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
		AppliedRevision: group.AppliedRevision, EffectiveAt: at(0), Group: group, RuntimeLayout: frozen,
	}
	require.NoError(t, p.saveRuntimeSnapshot(t.Context(), snap))
	return snap
}

// TestMySQLPartitionMaintenanceSingleFlightAndShutdown proves the maintenance
// flight is one asynchronous lifetime on the shared loop: exactly one blocked
// OpenDestination, Reconcile still listing while it blocks, a bounded shutdown
// that cancels the open and waits for the flight, and no second open after.
func TestMySQLPartitionMaintenanceSingleFlightAndShutdown(t *testing.T) {
	f, groups, destinations := newMySQLMaintenanceFixture(t)
	group, frozen := mysqlNativeMonthlyFixture(t)
	groups.groups = []*workspace.WriteGroup{liveGroup(group, workspace.WriteGroupStatusReady)}
	groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{
		group.ID: {{
			WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
			AppliedRevision: group.AppliedRevision, EffectiveAt: at(0), Group: group, RuntimeLayout: frozen,
		}},
	}

	pipeline := New(Dependencies{
		Groups: groups, Tags: f.tags, Destinations: destinations, Inspector: f.inspect, Store: f.store,
	}, Config{
		NodeID: "node-1", Owner: "node-1/test", Now: func() time.Time { return f.clock },
		TickInterval: 5 * time.Millisecond, ReconcileInterval: 5 * time.Millisecond,
		PartitionMaintenanceInterval: 5 * time.Millisecond,
	})
	snap := mysqlNativeMaintenanceSnapshot(t, pipeline)
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
		t.Fatal("maintenance never reached OpenDestination")
	}
	require.Equal(t, int32(1), destinations.openCalls.Load(), "single flight: exactly one open")

	listAtEntry := groups.listCalls.Load()
	require.Eventually(t, func() bool { return groups.listCalls.Load() > listAtEntry+3 },
		5*time.Second, 5*time.Millisecond,
		"the shared loop must keep reconciling while the maintenance open is blocked")
	require.Equal(t, int32(1), destinations.openCalls.Load(),
		"overlapping maintenance ticks must not open a second flight")

	stopCtx, cancelStop := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancelStop()
	require.NoError(t, pipeline.StopContext(stopCtx), "shutdown cancels the blocked open and waits for the flight")
	require.Equal(t, int32(1), destinations.openCalls.Load(), "no open after shutdown")
	require.NoError(t, pipeline.WaitStopped(stopCtx))
	require.Equal(t, int32(1), destinations.openCalls.Load())

	var maintenanceRuns int
	require.NoError(t, f.store.DB().QueryRowContext(t.Context(),
		`SELECT COUNT(*) FROM managed_schema_operations WHERE action = ?`, recordingplan.MySQLPartitionMaintenanceAction).Scan(&maintenanceRuns))
	require.Zero(t, maintenanceRuns,
		"a blocked open must not reach the partition ledger; the flight stays at the destination seam")
}
