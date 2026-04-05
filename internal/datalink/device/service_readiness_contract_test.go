package device

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CheckReadiness_AllowsPlanningAfterConnectSuccessProbeFailure(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)

	lastSuccess := false
	now := time.Now()
	require.NoError(t, repo.Create(context.Background(), &schema.Device{
		ID:               "dev-probe-failed",
		Name:             "dev-probe-failed",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":2,"probe_address":"49999","probe_function":"03"}`,
		LastTestSuccess:  &lastSuccess,
		LastTestError:    "讀取探測失敗: host=127.0.0.1 port=502 slave=1 function=03 address=49999, read 失敗: timeout",
		CreatedAt:        now,
		UpdatedAt:        now,
	}))

	readiness, err := svc.CheckReadiness(context.Background(), "dev-probe-failed")
	require.NoError(t, err)

	assert.Equal(t, schema.ReadinessStageStatusSuccess, readiness.ConnectStatus)
	assert.Equal(t, schema.ReadinessStageStatusFailed, readiness.ProbeStatus)
	assert.True(t, readiness.PlanningAllowed)
	assert.False(t, readiness.ActivationAllowed)
	assert.False(t, readiness.ApplyAllowed)
	assert.Equal(t, "warning", readiness.Status)
	require.NotEmpty(t, readiness.BlockingReasons)
	assert.Contains(t, readiness.BlockingReasons[0], "讀取探測失敗")
}

func TestService_CheckReadiness_AllowsActivationAfterProbeSuccess(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)

	lastSuccess := true
	now := time.Now()
	require.NoError(t, repo.Create(context.Background(), &schema.Device{
		ID:               "dev-probe-ok",
		Name:             "dev-probe-ok",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":2}`,
		LastTestSuccess:  &lastSuccess,
		CreatedAt:        now,
		UpdatedAt:        now,
	}))

	readiness, err := svc.CheckReadiness(context.Background(), "dev-probe-ok")
	require.NoError(t, err)

	assert.Equal(t, schema.ReadinessStageStatusSuccess, readiness.ConnectStatus)
	assert.Equal(t, schema.ReadinessStageStatusSuccess, readiness.ProbeStatus)
	assert.True(t, readiness.PlanningAllowed)
	assert.True(t, readiness.ActivationAllowed)
	assert.True(t, readiness.ApplyAllowed)
	assert.Equal(t, "ready", readiness.Status)
	assert.Empty(t, readiness.BlockingReasons)
}

func TestService_CheckReadiness_ReportsBlockingReasonsWhenConfigMissing(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)

	now := time.Now()
	require.NoError(t, repo.Create(context.Background(), &schema.Device{
		ID:        "dev-missing-config",
		Name:      "dev-missing-config",
		Protocol:  schema.ProtocolModbusTCP,
		Status:    schema.DeviceStatusDraft,
		CreatedAt: now,
		UpdatedAt: now,
	}))

	readiness, err := svc.CheckReadiness(context.Background(), "dev-missing-config")
	require.NoError(t, err)

	assert.Equal(t, schema.ReadinessStageStatusUnknown, readiness.ConnectStatus)
	assert.Equal(t, schema.ReadinessStageStatusUnknown, readiness.ProbeStatus)
	assert.False(t, readiness.PlanningAllowed)
	assert.False(t, readiness.ActivationAllowed)
	assert.False(t, readiness.ApplyAllowed)
	assert.Equal(t, "error", readiness.Status)
	require.Equal(t, []string{"missing connection configuration"}, readiness.BlockingReasons)
}

func TestService_CheckReadiness_TransitionsAcrossConnectAndProbeOutcomes(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)

	now := time.Now()
	require.NoError(t, repo.Create(context.Background(), &schema.Device{
		ID:               "dev-transition",
		Name:             "dev-transition",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":2,"probe_address":"49999","probe_function":"03"}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}))

	readiness, err := svc.CheckReadiness(context.Background(), "dev-transition")
	require.NoError(t, err)
	assert.Equal(t, schema.ReadinessStageStatusUnknown, readiness.ConnectStatus)
	assert.Equal(t, schema.ReadinessStageStatusUnknown, readiness.ProbeStatus)
	assert.False(t, readiness.PlanningAllowed)
	assert.False(t, readiness.ActivationAllowed)

	record, err := repo.GetByID(context.Background(), "dev-transition")
	require.NoError(t, err)
	connectFailed := false
	record.LastTestSuccess = &connectFailed
	record.LastTestError = "dial tcp 127.0.0.1:502: connect: connection refused"
	require.NoError(t, repo.Update(context.Background(), record))

	readiness, err = svc.CheckReadiness(context.Background(), "dev-transition")
	require.NoError(t, err)
	assert.Equal(t, schema.ReadinessStageStatusFailed, readiness.ConnectStatus)
	assert.Equal(t, schema.ReadinessStageStatusUnknown, readiness.ProbeStatus)
	assert.False(t, readiness.PlanningAllowed)
	assert.False(t, readiness.ActivationAllowed)
	require.NotEmpty(t, readiness.BlockingReasons)
	assert.Contains(t, readiness.BlockingReasons[0], "connection refused")

	record, err = repo.GetByID(context.Background(), "dev-transition")
	require.NoError(t, err)
	probeFailed := false
	record.LastTestSuccess = &probeFailed
	record.LastTestError = "讀取探測失敗: host=127.0.0.1 port=502 slave=1 function=03 address=49999, read 失敗: timeout"
	require.NoError(t, repo.Update(context.Background(), record))

	readiness, err = svc.CheckReadiness(context.Background(), "dev-transition")
	require.NoError(t, err)
	assert.Equal(t, schema.ReadinessStageStatusSuccess, readiness.ConnectStatus)
	assert.Equal(t, schema.ReadinessStageStatusFailed, readiness.ProbeStatus)
	assert.True(t, readiness.PlanningAllowed)
	assert.False(t, readiness.ActivationAllowed)

	record, err = repo.GetByID(context.Background(), "dev-transition")
	require.NoError(t, err)
	probeSucceeded := true
	record.LastTestSuccess = &probeSucceeded
	record.LastTestError = ""
	require.NoError(t, repo.Update(context.Background(), record))

	readiness, err = svc.CheckReadiness(context.Background(), "dev-transition")
	require.NoError(t, err)
	assert.Equal(t, schema.ReadinessStageStatusSuccess, readiness.ConnectStatus)
	assert.Equal(t, schema.ReadinessStageStatusSuccess, readiness.ProbeStatus)
	assert.True(t, readiness.PlanningAllowed)
	assert.True(t, readiness.ActivationAllowed)
	assert.True(t, readiness.ApplyAllowed)
	assert.Empty(t, readiness.BlockingReasons)
}
