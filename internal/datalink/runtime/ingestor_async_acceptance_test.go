package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type asyncAcceptingTarget struct{}

func (asyncAcceptingTarget) WriteTagValue(context.Context, string, any, time.Time) error {
	return nil
}

func (asyncAcceptingTarget) AcceptTagValue(context.Context, string, any, time.Time) (bool, error) {
	return true, nil
}

func TestHandleCollectedValue_AsyncAcceptanceDoesNotClaimDatabaseWrite(t *testing.T) {
	s := targetOutcomeService(&mockWriter{}, asyncAcceptingTarget{})
	s.handleCollectedValue(t.Context(), collector.CollectedValue{
		PointID: "p1", DeviceID: "dev-1", Value: int16(7),
		Timestamp: time.Date(2026, 11, 4, 10, 0, 1, 0, time.UTC), Quality: schema.QualityGood,
	})

	diagnostics := s.databaseDeliveryDiagnostics("")
	require.Len(t, diagnostics, 1)
	require.Equal(t, DatabaseDeliveryStatus("accepted"), diagnostics[0].Status)
	require.Equal(t, []DatabaseDeliveryStage{"collected", "mapped", "buffer_accepted"}, diagnostics[0].Stages)
	require.Nil(t, diagnostics[0].LastSuccessAt)
	flow := s.runtimeFlowDiagnostics("")
	require.Len(t, flow, 1)
	require.NotEqual(t, RuntimeFlowStageDatabaseDelivery, flow[0].LatestSuccessfulStage)
}
