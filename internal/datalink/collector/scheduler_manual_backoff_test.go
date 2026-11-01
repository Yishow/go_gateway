package collector

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestManualPlannedBackoffYieldsCapacityAndKeepsDeviceGate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s, protocol := manualAdmissionFixture(t)
		s.active = newActiveCapacity(1)
		s.config.DefaultRetryCount, s.config.DefaultRetryDelay = 1, 10*time.Millisecond
		otherGroup := "other-group"
		s.AddDevice(&schema.Device{ID: "other", Protocol: schema.ProtocolMC3E, ConnectionConfig: "{}"})
		s.AddPoint(&schema.Point{ID: "other-point", DeviceID: "other", Address: "D0", DataType: schema.DataTypeInt16, PollingGroupID: &otherGroup})
		var firstDevice, otherDevice atomic.Int32
		protocol.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
			request, _, err := req.Snapshot()
			if err != nil {
				return nil, err
			}
			if request.Group.Identity.DeviceID == "device" {
				if firstDevice.Add(1) == 1 {
					failure := errors.New("retryable")
					return executionResults(req, time.Now(), failure), failure
				}
			} else {
				otherDevice.Add(1)
			}
			return executionResults(req, time.Now(), nil), nil
		}
		done := make(chan []CollectedValue, 1)
		go func() { done <- s.PollNowContext(t.Context(), []string{"point"}) }()
		synctest.Wait()
		require.EqualValues(t, 1, firstDevice.Load(), "manual must reach the shared planned request before retry wait")
		require.NoError(t, s.Start([]*schema.PollingGroup{{ID: otherGroup, Enabled: true, IntervalMs: 5}}))
		defer s.Stop()
		blockedCtx, cancel := context.WithCancel(t.Context())
		defer cancel()
		blocked := make(chan []CollectedValue, 1)
		go func() { blocked <- s.PollNowContext(blockedCtx, []string{"point"}) }()
		time.Sleep(10 * time.Millisecond)
		synctest.Wait()
		require.Positive(t, otherDevice.Load(), "manual backoff must yield its active permit to scheduled device")
		require.EqualValues(t, 1, firstDevice.Load(), "same-device caller cannot interleave during retry backoff")
		cancel()
		blockedValues := <-blocked
		require.Len(t, blockedValues, 1)
		require.Equal(t, schema.QualityBad, blockedValues[0].Quality)
		require.NoError(t, s.Stop())
		values := <-done
		require.Len(t, values, 1)
		require.Equal(t, schema.QualityGood, values[0].Quality, values[0].Error)
		require.EqualValues(t, 2, firstDevice.Load(), "the original manual call resumes within its existing request budget")
	})
}
