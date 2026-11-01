package collector

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestScheduledPartialScalarFailurePreservesIndependentPointResultsAndHealth(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0", "D1", "D2")
	completed := time.Unix(123, 0).UTC()
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, _ int) ([]connector.PlannedPointResult, error) {
		results := executionResults(req, completed, nil)
		results[1].ReadResult = connector.ReadResult{Quality: schema.QualityBad, Error: "local scalar decode failure"}
		return results, nil // The shared wire exchange succeeded, not every decoder.
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "partial", s.deviceBreakers["device"])
	require.EqualValues(t, 1, p.sends.Load())
	values := collectedExecutionValues(s)
	require.Len(t, values, 3)
	good, bad := 0, 0
	ids := map[string]bool{}
	for _, value := range values {
		require.False(t, ids[value.PointID])
		ids[value.PointID] = true
		require.Equal(t, "partial/initial/0", value.RequestID)
		require.Equal(t, completed, value.ObservedAt)
		require.Equal(t, completed, value.ReceivedAt)
		require.Equal(t, scan.bindings[value.PointID].ConfigFingerprint, value.ConfigFingerprint)
		require.Equal(t, scan.bindings[value.PointID].PointRevision, value.PointRevision)
		if value.Quality == schema.QualityBad {
			bad++
			require.Nil(t, value.Value)
			require.Contains(t, value.Error, "local scalar")
		} else {
			good++
			require.Equal(t, schema.QualityGood, value.Quality)
			require.NotNil(t, value.Value)
		}
	}
	require.Equal(t, 2, good)
	require.Equal(t, 1, bad)
	stats := s.deviceBreakers["device"].GetStats().TrackerStats
	require.Equal(t, 3, stats.TotalRequests)
	require.Equal(t, 2, stats.SuccessCount)
	require.Equal(t, 1, stats.FailureCount)
}

func TestScheduledPartialRequestFailureDoesNotContaminateOtherRequests(t *testing.T) {
	s, scan, p, conn := scheduledExecutionFixture(t, `{}`, "D0", "D10", "D20")
	require.Len(t, scan.plan.Requests(), 3, "unknown capability cannot invent a covering interval")
	p.hook = func(_ context.Context, req connector.PlannedReadRequest, send int) ([]connector.PlannedPointResult, error) {
		var err error
		if send == 2 {
			err = errors.New("bad request frame")
		}
		return executionResults(req, time.Unix(int64(100+send), 0), err), err
	}
	s.executeScheduledPlan(t.Context(), conn, scan, "separate", s.deviceBreakers["device"])
	require.EqualValues(t, 3, p.sends.Load(), "one failed request must not trigger scalar fallback")
	values := collectedExecutionValues(s)
	require.Len(t, values, 3)
	requestIDs, completions := map[string]bool{}, map[time.Time]bool{}
	good, bad := 0, 0
	for _, value := range values {
		require.Equal(t, "separate", value.AcquisitionID)
		require.False(t, requestIDs[value.RequestID])
		requestIDs[value.RequestID], completions[value.ReceivedAt] = true, true
		require.Equal(t, value.ReceivedAt, value.ObservedAt)
		if value.Quality == schema.QualityBad {
			bad++
			require.Nil(t, value.Value)
		} else {
			good++
			require.Equal(t, schema.QualityGood, value.Quality)
		}
	}
	require.Len(t, completions, 3, "separate wire completions must not become one atomic scan timestamp")
	require.Equal(t, 2, good)
	require.Equal(t, 1, bad)
	stats := s.deviceBreakers["device"].GetStats().TrackerStats
	require.Equal(t, 3, stats.TotalRequests)
	require.Equal(t, 2, stats.SuccessCount)
	require.Equal(t, 1, stats.FailureCount)
}
