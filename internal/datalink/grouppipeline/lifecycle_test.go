package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type shutdownGroups struct {
	GroupSource
	entered chan struct{}
	release chan struct{}
	calls   atomic.Int32
	blockAt int32
}

func (g *shutdownGroups) List(ctx context.Context) (*workspace.WriteGroupListResult, error) {
	if g.calls.Add(1) == g.blockAt {
		close(g.entered)
		<-g.release // Non-cooperative IO must remain represented after timeout.
	}
	return g.GroupSource.List(ctx)
}

func TestPipelineShutdownDeadlineIncludesBlockedLoop(t *testing.T) {
	f := newFixture(t)
	groups := &shutdownGroups{GroupSource: f.groups, entered: make(chan struct{}), release: make(chan struct{}), blockAt: 2}
	release := sync.OnceFunc(func() { close(groups.release) })
	defer release()
	p := f.pipeline()
	p.deps.Groups = groups
	p.config.ReconcileInterval = time.Millisecond
	require.NoError(t, p.Start(t.Context()))
	require.Error(t, p.WaitStopped(t.Context()), "a live pipeline is not stopped")
	select {
	case <-groups.entered:
	case <-time.After(2 * time.Second):
		t.Fatal("reconcile loop did not enter the barrier")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, p.StopContext(ctx), context.DeadlineExceeded)
	require.ErrorIs(t, p.StopContext(ctx), context.DeadlineExceeded)
	require.ErrorIs(t, p.WaitStopped(ctx), context.DeadlineExceeded)
	require.Error(t, p.Start(t.Context()))
	require.Error(t, p.AcceptSample(t.Context(), envelope(f.groups.snapshots["group-G"][0].Group, "rejected", at(1), 1.0)))
	release()
	require.NoError(t, p.WaitStopped(t.Context()))
	require.NoError(t, p.Stop(time.Second))
	require.NoError(t, p.Start(t.Context()))
	require.NoError(t, p.Stop(time.Second))
}

func TestPipelineShutdownDuringStartupRemainsObservable(t *testing.T) {
	f := newFixture(t)
	groups := &shutdownGroups{GroupSource: f.groups, entered: make(chan struct{}), release: make(chan struct{}), blockAt: 1}
	release := sync.OnceFunc(func() { close(groups.release) })
	defer release()
	p := f.pipeline()
	p.deps.Groups = groups
	started := make(chan error, 1)
	go func() { started <- p.Start(t.Context()) }()
	select {
	case <-groups.entered:
	case <-time.After(time.Second):
		t.Fatal("startup did not enter the barrier")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.ErrorIs(t, p.StopContext(ctx), context.Canceled)
	require.Error(t, p.Start(t.Context()))
	require.ErrorIs(t, p.WaitStopped(ctx), context.Canceled)
	release()
	require.Error(t, <-started)
	require.NoError(t, p.WaitStopped(t.Context()))
}

func TestPipelineShutdownNeverStartedAndConcurrentCalls(t *testing.T) {
	p := newFixture(t).pipeline()
	results := make(chan error, 20)
	var wg sync.WaitGroup
	for range cap(results) {
		wg.Go(func() { results <- p.StopContext(t.Context()) })
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}
	require.NoError(t, p.WaitStopped(t.Context()))
}

func TestPipelineStartupCancellationDoesNotStartDeliveryAfterReconcile(t *testing.T) {
	f := newFixture(t)
	groups := &shutdownGroups{GroupSource: f.groups, entered: make(chan struct{}), release: make(chan struct{}), blockAt: 1}
	release := sync.OnceFunc(func() { close(groups.release) })
	defer release()
	p := f.pipeline()
	p.deps.Groups = groups
	startup, cancel := context.WithCancel(t.Context())
	defer cancel()
	started := make(chan error, 1)
	go func() { started <- p.Start(startup) }()
	select {
	case <-groups.entered:
	case <-time.After(time.Second):
		t.Fatal("startup did not enter reconcile")
	}
	cancel()
	release()
	require.ErrorIs(t, <-started, context.Canceled)
	p.lifecycle.Lock()
	worker := p.worker
	p.lifecycle.Unlock()
	require.Nil(t, worker, "canceled startup cannot construct a delivery worker")
	require.NoError(t, p.StopContext(t.Context()))
	require.NoError(t, p.WaitStopped(t.Context()))
}

func TestPipelineStartupLateCancellationRollsBackLaunchedWorker(t *testing.T) {
	f := newFixture(t)
	// Reconcile refuses before reading storage so the held connection blocks
	// precisely the subsequent worker recovery.
	f.groups.listErr = errors.New("unavailable test groups")
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "recovery.db"))
	require.NoError(t, err)
	defer db.Close()
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	db.SetMaxOpenConns(1)
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	defer conn.Close()
	f.store = groupdelivery.NewStore(db)
	p := f.pipeline()
	loopEntered, loopRelease := make(chan struct{}), make(chan struct{})
	release := sync.OnceFunc(func() { close(loopRelease) })
	defer release()
	var nowCalls atomic.Int32
	p.config.Delivery.Now = func() time.Time {
		if nowCalls.Add(1) == 2 {
			close(loopEntered)
			<-loopRelease
		}
		return time.Now()
	}
	startup, cancel := context.WithCancel(t.Context())
	defer cancel()
	started := make(chan error, 1)
	baseline := db.Stats().WaitCount
	go func() { started <- p.Start(startup) }()
	require.Eventually(t, func() bool { return db.Stats().WaitCount > baseline }, time.Second, time.Millisecond)
	// Start is blocked in recovery with lifecycle unlocked. Hold the lock so it
	// cannot publish success before cancellation of its original startup request.
	p.lifecycle.Lock()
	require.NoError(t, conn.Close())
	select {
	case <-loopEntered:
	case <-time.After(time.Second):
		p.lifecycle.Unlock()
		t.Fatal("worker loop did not reach its barrier")
	}
	cancel()
	p.lifecycle.Unlock()
	require.ErrorIs(t, <-started, context.Canceled)
	require.ErrorIs(t, p.WaitStopped(startup), context.Canceled, "rollback remains observable while original worker is live")
	require.Error(t, p.Start(t.Context()))
	release()
	require.NoError(t, p.WaitStopped(t.Context()))
	require.NoError(t, p.Start(t.Context()), "retry is allowed after rollback really finishes")
	require.NoError(t, p.StopContext(t.Context()))
}
