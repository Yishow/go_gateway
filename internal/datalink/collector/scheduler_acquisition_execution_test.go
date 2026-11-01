package collector

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/wirebudget"

	"github.com/stretchr/testify/require"
)

// This controlled adapter supplements, not replaces, cmd/test_ui wire fixtures.
// Its send counter advances only after the production admission/budget boundary.
type scheduledExecutionProtocol struct {
	connector.Protocol
	sends, connects atomic.Int32
	generation      atomic.Uint64
	timeout         time.Duration
	hook            func(context.Context, connector.PlannedReadRequest, int) ([]connector.PlannedPointResult, error)
	requests        []connector.PlannedReadRequest
	budgets         []*wirebudget.Budget
	deadlines       []time.Time
}

func (p *scheduledExecutionProtocol) Connect(context.Context, string) error {
	p.connects.Add(1)
	p.generation.Add(1)
	return nil
}
func (*scheduledExecutionProtocol) Close() error      { return nil }
func (*scheduledExecutionProtocol) IsConnected() bool { return true }
func (p *scheduledExecutionProtocol) AcquisitionTransportGeneration() uint64 {
	return p.generation.Load()
}
func (p *scheduledExecutionProtocol) AcquisitionReadTimeout() time.Duration { return p.timeout }
func (p *scheduledExecutionProtocol) ReadPlanned(ctx context.Context, req connector.PlannedReadRequest) ([]connector.PlannedPointResult, error) {
	if err := wirebudget.Take(ctx); err != nil {
		return executionResults(req, time.Now(), err), err
	}
	count := int(p.sends.Add(1))
	p.requests = append(p.requests, req)
	p.budgets = append(p.budgets, wirebudget.From(ctx))
	deadline, _ := ctx.Deadline()
	p.deadlines = append(p.deadlines, deadline)
	if p.hook != nil {
		return p.hook(ctx, req, count)
	}
	return executionResults(req, time.Unix(int64(count), 0), nil), nil
}

func executionResults(req connector.PlannedReadRequest, at time.Time, err error) []connector.PlannedPointResult {
	r, bindings, snapshotErr := req.Snapshot()
	if snapshotErr != nil {
		return nil
	}
	out := make([]connector.PlannedPointResult, len(bindings))
	for i, binding := range bindings {
		out[i] = connector.PlannedPointResult{Binding: binding, Group: r.Group, AcquisitionID: req.AcquisitionID, RequestID: req.RequestID, PlanFingerprint: req.Plan.Fingerprint(), ReceivedAt: at, ObservedAt: at, ReadResult: connector.ReadResult{Quality: schema.QualityGood, Value: int16(i + 10)}}
		if err != nil {
			out[i].ReadResult = connector.ReadResult{Quality: schema.QualityBad, Error: err.Error()}
		}
	}
	return out
}

func scheduledExecutionFixture(t *testing.T, config string, addresses ...string) (*Scheduler, scheduledPlan, *scheduledExecutionProtocol, *connector.ManagedConnection) {
	t.Helper()
	s, snapshot := scheduledPlanFixture(t, config, addresses...)
	s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 0, 0
	scan, err := buildScheduledPlan(snapshot)
	require.NoError(t, err)
	require.Empty(t, scan.failures)
	p := &scheduledExecutionProtocol{timeout: time.Second}
	conn := &connector.ManagedConnection{DeviceID: snapshot.device.ID, ProtocolType: snapshot.device.Protocol, Config: snapshot.device.Config, Protocol: p}
	return s, scan, p, conn
}

func collectedExecutionValues(s *Scheduler) []CollectedValue {
	var values []CollectedValue
	for {
		select {
		case value := <-s.ValueChannel():
			values = append(values, value)
		default:
			return values
		}
	}
}

func TestScheduledExecutionRetriesShareOneRequestBudgetAndFinalCompletion(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0", "D1")
	s.config.DefaultRetryCount = 2
	final := time.Unix(30, 0).UTC()
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, send int) ([]connector.PlannedPointResult, error) {
		if send == 1 {
			err := errors.New("transient read")
			return executionResults(req, time.Unix(20, 0), err), err
		}
		return executionResults(req, final, nil), nil
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "scan", s.deviceBreakers["device"])
	require.EqualValues(t, 2, p.sends.Load())
	require.Same(t, p.budgets[0], p.budgets[1])
	require.Equal(t, 2, p.budgets[0].Used())
	require.Equal(t, p.requests[0].RequestID, p.requests[1].RequestID)
	require.Equal(t, p.deadlines[0], p.deadlines[1], "retry must not restart the request timeout")
	values := collectedExecutionValues(s)
	require.Len(t, values, 2)
	for _, value := range values {
		require.Equal(t, schema.QualityGood, value.Quality)
		require.Equal(t, final, value.ReceivedAt)
		require.Equal(t, final, value.ObservedAt)
		require.Equal(t, "scan", value.AcquisitionID)
		require.Equal(t, scan.bindings[value.PointID].InstallationLease, value.InstallationLease)
	}
	require.Equal(t, 2, s.deviceBreakers["device"].GetStats().TrackerStats.TotalRequests, "health retains final Point units, not wire attempts")
}

func TestScheduledExecutionZeroTokenFallsBackOnceWithoutPersistentDowngrade(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`, "D0", "D10", "D20")
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, send int) ([]connector.PlannedPointResult, error) {
		r, _ := req.Plan.Request(req.RequestIndex)
		if r.Kind == "random" {
			err := &mcprotocol.MCError{Code: 0xC059}
			return executionResults(req, time.Unix(int64(send), 0), err), err
		}
		return executionResults(req, time.Unix(int64(send), 0), nil), nil
	}
	for _, acquisitionID := range []string{"first", "second"} {
		s.executeScheduledPlan(t.Context(), conn, scan, acquisitionID, nil)
		values := collectedExecutionValues(s)
		require.Len(t, values, 3)
		ids := map[string]bool{}
		for _, value := range values {
			require.Equal(t, schema.QualityGood, value.Quality)
			require.Contains(t, value.RequestID, acquisitionID+"/fallback/")
			require.False(t, ids[value.PointID])
			ids[value.PointID] = true
		}
		require.Empty(t, s.acquisitionDowngrades, "an unmanaged fixture cannot invent a successful-connect owner")
	}
	require.EqualValues(t, 8, p.sends.Load(), "each scan gets one random and exactly three safe requests")
}

func TestScheduledExecutionAdapterResendConsumesSchedulerBudget(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0", "D1")
	s.config.DefaultRetryCount = 1
	var resendErr error
	p.hook = func(ctx context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		// Represent an adapter reconnect resend using the real shared budget.
		resendErr = wirebudget.Take(ctx)
		if resendErr == nil {
			p.sends.Add(1)
		}
		err := errors.New("adapter exhausted its reconnect resend")
		return executionResults(req, time.Unix(200, 0), err), err
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "bounded", nil)
	require.NoError(t, resendErr)
	require.EqualValues(t, 2, p.sends.Load())
	require.Len(t, p.requests, 1, "scheduler must not mint another call after the adapter used both sends")
	require.Equal(t, 2, p.budgets[0].Used())
	values := collectedExecutionValues(s)
	require.Len(t, values, 2)
	for _, value := range values {
		require.Equal(t, schema.QualityBad, value.Quality)
		require.Nil(t, value.Value)
	}
}

func TestScheduledExecutionMalformedRejectionCannotAuthorizeFallback(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`, "D0", "D10", "D20")
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		err := &mcprotocol.MCError{Code: 0xC059}
		results := executionResults(req, time.Unix(10, 0), err)
		results[0].RequestID = "foreign-request"
		return results, err
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "malformed", nil)
	require.EqualValues(t, 1, p.sends.Load())
	require.Empty(t, s.acquisitionDowngrades)
	values := collectedExecutionValues(s)
	require.Len(t, values, 3)
	for _, value := range values {
		require.Equal(t, schema.QualityBad, value.Quality)
		require.Nil(t, value.Value)
		require.Equal(t, "malformed/initial/0", value.RequestID)
	}
}
