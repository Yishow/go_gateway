package device

import (
	"errors"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestReadReadinessDerivesDiagnosticsWithoutUpdatingRepository(t *testing.T) {
	failed := false
	repo := &stubDeviceRepository{
		device: &schema.Device{
			ID: "device-read", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusDraft,
			ConnectionConfig: `{"host":"127.0.0.1","port":1,"slave_id":1}`,
			LastTestSuccess:  &failed, LastTestError: "dial tcp private.invalid:502: password=fixture-secret",
			ReadinessStatus: `{"device_id":"device-read","status":"ready","connect_status":"success","probe_status":"success","activation_allowed":true,"availability_status":"unavailable","availability_reason":"operator disabled"}`,
		},
		updateErr: errors.New("observation must not invoke Update"),
	}
	before := *repo.device
	readiness, err := NewService(repo, nil).ReadReadiness(t.Context(), "device-read")
	require.NoError(t, err)
	require.Equal(t, before, *repo.device)
	require.Equal(t, "error", readiness.Status)
	require.Equal(t, schema.ReadinessStageStatusFailed, readiness.ConnectStatus)
	require.Equal(t, schema.ReadinessStageStatusUnknown, readiness.ProbeStatus)
	require.False(t, readiness.PlanningAllowed)
	require.False(t, readiness.ActivationAllowed)
	require.False(t, readiness.ApplyAllowed)
	require.Equal(t, AvailabilityStatusUnavailable, readiness.AvailabilityStatus)
	require.Equal(t, "operator disabled", readiness.AvailabilityReason)
	require.NotEmpty(t, readiness.BlockingReasons)
	for _, reason := range readiness.BlockingReasons {
		require.NotContains(t, reason, "private.invalid")
		require.NotContains(t, reason, "fixture-secret")
	}
	for _, check := range readiness.Checks {
		require.NotContains(t, check.Message, "private.invalid")
		require.NotContains(t, check.Message, "fixture-secret")
	}
}

func TestReadReadinessReportsRepositoryReadFailure(t *testing.T) {
	readiness, err := NewService(&stubDeviceRepository{}, nil).ReadReadiness(t.Context(), "missing")
	require.ErrorContains(t, err, "device not found")
	require.Nil(t, readiness)
}
