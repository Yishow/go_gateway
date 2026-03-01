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

func TestService_Activate_WithProbeSuccess(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))

	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()

	cfg := map[string]interface{}{
		"host":     "127.0.0.1",
		"port":     server.Port(),
		"slave_id": 1,
		"timeout":  2,
	}
	cfgBytes, err := json.Marshal(cfg)
	require.NoError(t, err)

	repo := NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := NewService(repo, connMgr)

	now := time.Now()
	dev := &schema.Device{
		ID:               "dev-probe-ok",
		Name:             "dev-probe-ok",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: string(cfgBytes),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	require.NoError(t, repo.Create(context.Background(), dev))

	err = svc.Activate(context.Background(), dev.ID)
	require.NoError(t, err)

	updated, err := repo.GetByID(context.Background(), dev.ID)
	require.NoError(t, err)
	assert.Equal(t, schema.DeviceStatusActive, updated.Status)
	require.NotNil(t, updated.LastTestSuccess)
	assert.True(t, *updated.LastTestSuccess)
}

func TestService_Activate_BlockWhenProbeReadFails(t *testing.T) {
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
	dev := &schema.Device{
		ID:               "dev-probe-fail",
		Name:             "dev-probe-fail",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: string(cfgBytes),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	require.NoError(t, repo.Create(context.Background(), dev))

	err = svc.Activate(context.Background(), dev.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "讀取探測失敗")

	updated, getErr := repo.GetByID(context.Background(), dev.ID)
	require.NoError(t, getErr)
	assert.Equal(t, schema.DeviceStatusDraft, updated.Status)
	require.NotNil(t, updated.LastTestSuccess)
	assert.False(t, *updated.LastTestSuccess)
}
