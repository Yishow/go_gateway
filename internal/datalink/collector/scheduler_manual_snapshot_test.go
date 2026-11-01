package collector

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestManualSameInstalledGroupHasScheduledPlanIdentity(t *testing.T) {
	s, _ := scheduledPlanFixture(t, "{}", "D0", "D1")
	require.NoError(t, s.Start([]*schema.PollingGroup{{ID: "source-cadence", Enabled: true, IntervalMs: 60000}}))
	defer s.Stop()
	scheduled, err := buildScheduledPlan(s.captureScheduledGroup("source-cadence")[0])
	require.NoError(t, err)
	snapshots := s.captureManualPoints(t.Context(), []string{"point-1", "point-0"})
	require.Len(t, snapshots, 1)
	manual, err := buildScheduledPlan(snapshots[0])
	require.NoError(t, err)
	require.Equal(t, scheduled.input.Sources, manual.input.Sources)
	require.Equal(t, scheduled.input.ConfigurationRevision, manual.input.ConfigurationRevision, "execution entrypoint is not configuration authority")
	require.Equal(t, scheduled.plan.Fingerprint(), manual.plan.Fingerprint())
}
