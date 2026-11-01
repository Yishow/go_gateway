package collector

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/wirebudget"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestScheduledCapacityDeadlineEmitsOnlyCurrentSnapshot(t *testing.T) {
	for _, mode := range []string{"current", "replaced", "canceled"} {
		t.Run(mode, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s, scan, _, _ := scheduledExecutionFixture(t, `{"timeout":1}`, "D0")
				scan.snapshot.capturedAt = time.Now().Add(-980 * time.Millisecond)
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				scan.snapshot.ctx = ctx
				s.active = newActiveCapacity(1)
				lease, err := s.active.acquire(t.Context())
				require.NoError(t, err)
				defer lease.release()
				go func() {
					time.Sleep(10 * time.Millisecond)
					if mode == "canceled" {
						cancel()
					}
					if mode == "replaced" {
						point := scan.snapshot.points[0].identityPoint()
						s.AddPointWithObservationLease(&point, scan.snapshot.points[0].readProfile)
					}
				}()
				s.pollScheduledSnapshot(scan.snapshot)
				values := collectedExecutionValues(s)
				if mode != "current" {
					require.Empty(t, values)
					return
				}
				require.Len(t, values, 1, "permit deadline must not silently lose a current scan")
				require.Equal(t, schema.QualityBad, values[0].Quality)
				require.Equal(t, "scan-deadline", values[0].QualityReason)
			})
		})
	}
}

func TestScheduledCapacityReacquiresAfterRequestTimeoutDuringBackoff(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0", "D10")
		require.Len(t, scan.plan.Requests(), 2)
		s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 1, 100*time.Millisecond
		p.timeout = 20 * time.Millisecond
		capacity := newActiveCapacity(1)
		lease, err := capacity.acquire(t.Context())
		require.NoError(t, err)
		ctx := wirebudget.WithWaitCooperator(t.Context(), lease)
		releaseB := make(chan struct{})
		unblockB := sync.OnceFunc(func() { close(releaseB) })
		defer unblockB()
		bDone := make(chan struct{})
		go func() {
			defer close(bDone)
			b, acquireErr := capacity.acquire(t.Context())
			if acquireErr != nil {
				return
			}
			defer b.release()
			<-releaseB
		}()
		synctest.Wait()
		p.hook = func(_ context.Context, req connector.PlannedReadRequest, send int) ([]connector.PlannedPointResult, error) {
			if send == 1 {
				failure := errors.New("retryable first request")
				return executionResults(req, time.Now(), failure), failure
			}
			return executionResults(req, time.Now(), nil), nil
		}
		done := make(chan struct{})
		go func() { defer close(done); defer lease.release(); s.executeScheduledPlan(ctx, conn, scan, "scan", nil) }()
		time.Sleep(25 * time.Millisecond)
		synctest.Wait()
		require.EqualValues(t, 1, p.sends.Load(), "request timeout during yielded backoff must not let next request bypass B's active permit")
		unblockB()
		<-done
		<-bDone
		require.EqualValues(t, 2, p.sends.Load())
		values := collectedExecutionValues(s)
		require.Len(t, values, 2)
		require.Equal(t, schema.QualityBad, values[0].Quality)
		require.Equal(t, schema.QualityGood, values[1].Quality)
	})
}

func TestScheduledCapacityUnsentWaitDoesNotChargeDeviceHealth(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0", "D10")
		s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 1, 100*time.Millisecond
		p.timeout = 20 * time.Millisecond
		capacity := newActiveCapacity(1)
		lease, err := capacity.acquire(t.Context())
		require.NoError(t, err)
		ctx := wirebudget.WithWaitCooperator(t.Context(), lease)
		releaseB := make(chan struct{})
		defer close(releaseB)
		go func() {
			b, acquireErr := capacity.acquire(t.Context())
			if acquireErr != nil {
				return
			}
			defer b.release()
			<-releaseB
		}()
		synctest.Wait()
		p.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
			failure := errors.New("one failed wire attempt")
			return executionResults(req, time.Now(), failure), failure
		}
		done := make(chan struct{})
		go func() {
			defer close(done)
			defer lease.release()
			s.executeScheduledPlan(ctx, conn, scan, "scan", s.deviceBreakers["device"])
		}()
		time.Sleep(45 * time.Millisecond)
		synctest.Wait()
		<-done
		require.EqualValues(t, 1, p.sends.Load())
		values := collectedExecutionValues(s)
		require.Len(t, values, 2)
		for _, value := range values {
			require.Equal(t, schema.QualityBad, value.Quality)
		}
		require.Equal(t, 1, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests, "only first request's actual wire failure belongs to device health; second had no send")
	})
}

func TestScheduledCapacityResumeTimeoutRetainsActualWireFailure(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0", "D10")
		s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 1, 5*time.Millisecond
		p.timeout = 20 * time.Millisecond
		capacity := newActiveCapacity(1)
		lease, err := capacity.acquire(t.Context())
		require.NoError(t, err)
		ctx := wirebudget.WithWaitCooperator(t.Context(), lease)
		releaseB := make(chan struct{})
		defer close(releaseB)
		go func() {
			b, acquireErr := capacity.acquire(t.Context())
			if acquireErr != nil {
				return
			}
			defer b.release()
			<-releaseB
		}()
		synctest.Wait()
		p.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
			failure := errors.New("one failed wire attempt")
			return executionResults(req, time.Now(), failure), failure
		}
		done := make(chan struct{})
		go func() {
			defer close(done)
			defer lease.release()
			s.executeScheduledPlan(ctx, conn, scan, "scan", s.deviceBreakers["device"])
		}()
		time.Sleep(45 * time.Millisecond)
		synctest.Wait()
		<-done
		require.EqualValues(t, 1, p.sends.Load())
		values := collectedExecutionValues(s)
		require.Len(t, values, 2)
		for _, value := range values {
			require.Equal(t, schema.QualityBad, value.Quality)
		}
		require.Equal(t, 1, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests, "only first request's actual wire failure belongs to device health; second had no send")
	})
}
