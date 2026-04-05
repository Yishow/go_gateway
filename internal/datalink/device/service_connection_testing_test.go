package device

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_TestConnectionWithResult_ReturnsPlanningHintsAlongsideProbeFailure(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))

	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()

	cfg := map[string]interface{}{
		"host":           "127.0.0.1",
		"port":           server.Port(),
		"slave_id":       1,
		"timeout":        2,
		"probe_address":  "49999",
		"probe_function": "03",
	}
	cfgBytes, err := json.Marshal(cfg)
	require.NoError(t, err)

	repo := NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := NewService(repo, connMgr)

	now := time.Now()
	device := &schema.Device{
		ID:               "dev-probe-failure-hints",
		Name:             "dev-probe-failure-hints",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: string(cfgBytes),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	require.NoError(t, repo.Create(context.Background(), device))

	result, err := svc.TestConnectionWithResult(context.Background(), device.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "讀取探測失敗")
	assert.Equal(t, TestConnectionStageSuccess, result.Connect.Status)
	assert.Equal(t, TestConnectionStageFailed, result.Probe.Status)
	assert.True(t, result.PlanningHints.SourceRulePlanningSupported)
	assert.True(t, result.PlanningHints.ProbeSupported)
	assert.Empty(t, result.PlanningHints.PlanningBlockedReason)
}

func TestBuildPlanningCapabilityHints_ReportsUnsupportedPlanningForMQTT(t *testing.T) {
	hints := buildPlanningCapabilityHints(schema.ProtocolMQTT)

	assert.False(t, hints.SourceRulePlanningSupported)
	assert.False(t, hints.ProbeSupported)
	assert.Equal(t, "protocol does not expose addressed source-rule planning", hints.PlanningBlockedReason)
}
