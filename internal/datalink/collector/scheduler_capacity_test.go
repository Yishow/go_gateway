package collector

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"
)

func TestSchedulerCapacityRejectsExplicitInvalidButPreservesOmitted(t *testing.T) {
	for _, value := range []int{0, -1} {
		for _, active := range []bool{false, true} {
			config := DefaultSchedulerConfig()
			if active {
				config.MaxActiveScans = &value
			} else {
				config.PendingScanCapacity = &value
			}
			s := NewScheduler(config, nil)
			err := s.Start(nil)
			_ = s.Stop()
			require.Error(t, err)
		}
	}
	s := NewScheduler(SchedulerConfig{ValueBufferSize: 1}, nil)
	require.NoError(t, s.Start(nil))
	require.NoError(t, s.Stop())
}

func TestSchedulerCapacityFreezesCallerPointers(t *testing.T) {
	value := 1
	config := DefaultSchedulerConfig()
	config.MaxActiveScans = &value
	config.PendingScanCapacity = &value
	s := NewScheduler(config, nil)
	value = 0
	require.NoError(t, s.Start(nil))
	require.NoError(t, s.Stop())
}

func TestActiveCapacityBackoffFairResumeAndCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		capacity := newActiveCapacity(1)
		a, err := capacity.acquire(t.Context())
		require.NoError(t, err)
		defer a.release()
		var bEntered atomic.Bool
		bRelease := make(chan struct{})
		releaseB := sync.OnceFunc(func() { close(bRelease) })
		defer releaseB()
		bDone := make(chan struct{})
		go func() {
			defer close(bDone)
			b, e := capacity.acquire(t.Context())
			if e != nil {
				return
			}
			bEntered.Store(true)
			<-bRelease
			b.release()
		}()
		synctest.Wait()
		require.False(t, bEntered.Load(), "capacity one must bound active work")
		a.Suspend()
		synctest.Wait()
		require.True(t, bEntered.Load(), "backoff frees capacity for ready B")
		ctx, cancel := context.WithCancel(t.Context())
		resumed := make(chan error, 1)
		go func() { resumed <- a.Resume(ctx) }()
		synctest.Wait()
		select {
		case <-resumed:
			t.Fatal("resume bypassed B")
		default:
		}
		cancel()
		require.ErrorIs(t, <-resumed, context.Canceled)
		releaseB()
		<-bDone
		c, err := capacity.acquire(t.Context())
		require.NoError(t, err, "canceled resume cannot leak capacity")
		c.release()
	})
}

func TestActiveCapacityResumeCannotOvertakeReadyDevices(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		capacity := newActiveCapacity(1)
		a, err := capacity.acquire(t.Context())
		require.NoError(t, err)
		defer a.release()
		entered := make(chan string, 3)
		releaseB, releaseC := make(chan struct{}), make(chan struct{})
		unblockB := sync.OnceFunc(func() { close(releaseB) })
		unblockC := sync.OnceFunc(func() { close(releaseC) })
		defer unblockB()
		defer unblockC()
		for _, item := range []struct {
			id      string
			release chan struct{}
		}{{"b", releaseB}, {"c", releaseC}} {
			go func() {
				lease, acquireErr := capacity.acquire(t.Context())
				if acquireErr != nil {
					return
				}
				defer lease.release()
				entered <- item.id
				<-item.release
			}()
			synctest.Wait()
		}
		a.Suspend()
		synctest.Wait()
		require.Equal(t, "b", <-entered)
		resumed := make(chan error, 1)
		go func() { resumed <- a.Resume(t.Context()) }()
		synctest.Wait()
		unblockB()
		synctest.Wait()
		require.Equal(t, "c", <-entered, "already-ready C precedes resuming A")
		select {
		case <-resumed:
			t.Fatal("A resumed before C released")
		default:
		}
		unblockC()
		require.NoError(t, <-resumed)
	})
}
