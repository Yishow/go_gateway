package collector

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/protocol/wirebudget"

	"github.com/stretchr/testify/require"
)

// This controlled failure represents ensureConnectionContext returning before
// Take. Source review separately checks the real MC/FATEK/Modbus call paths.
type beforeSendFailureProtocol struct {
	connector.Protocol
	calls  atomic.Int32
	budget *wirebudget.Budget
}

func (*beforeSendFailureProtocol) AcquisitionReadTimeout() time.Duration {
	return 20 * time.Millisecond
}
func (p *beforeSendFailureProtocol) ReadPlanned(ctx context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	p.calls.Add(1)
	p.budget = wirebudget.From(ctx)
	err := errors.New("connection failure before protocol send")
	return executionResults(req, time.Now(), err), err
}

func TestScheduledCapacityRetainsInvokedConnectionFailureBeforeSend(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, scan, _, conn := scheduledExecutionFixture(t, `{}`, "D0", "D10")
		protocol := &beforeSendFailureProtocol{}
		conn.Protocol = protocol
		s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 1, 5*time.Millisecond
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
		done := make(chan struct{})
		go func() {
			defer close(done)
			defer lease.release()
			s.executeScheduledPlan(ctx, conn, scan, "scan", s.deviceBreakers["device"])
		}()
		time.Sleep(45 * time.Millisecond)
		synctest.Wait()
		<-done
		require.EqualValues(t, 1, protocol.calls.Load())
		require.NotNil(t, protocol.budget)
		require.Zero(t, protocol.budget.Used(), "a connection failure is not a protocol send")
		require.Len(t, collectedExecutionValues(s), 2)
		require.Equal(t, 1, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests, "request one invoked a real execution phase; request two only waited for capacity")
	})
}

func TestScheduledCapacityInheritedBudgetDoesNotChargeUninvokedPoint(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0")
		p.timeout = 20 * time.Millisecond
		shared := wirebudget.With(t.Context(), 3, 0)
		require.NoError(t, wirebudget.Take(shared), "represents a previous request's consumed shared allowance")
		capacity := newActiveCapacity(1)
		a, err := capacity.acquire(t.Context())
		require.NoError(t, err)
		a.Suspend()
		b, err := capacity.acquire(t.Context())
		require.NoError(t, err)
		defer b.release()
		ctx := wirebudget.WithWaitCooperator(shared, a)
		s.executeScheduledPlan(ctx, conn, scan, "scan", s.deviceBreakers["device"])
		require.Zero(t, p.sends.Load(), "current point never entered its connector")
		require.Equal(t, 1, wirebudget.From(ctx).Used(), "capacity waiting consumes no additional attempt")
		require.Len(t, collectedExecutionValues(s), 1)
		require.Zero(t, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests, "earlier shared-budget use is not this point's communication failure")
	})
}
