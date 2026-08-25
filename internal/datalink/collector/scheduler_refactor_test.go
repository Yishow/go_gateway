package collector

import (
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSchedulerConfig(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	assert.Equal(t, 3, cfg.DefaultRetryCount)
	assert.Equal(t, 1, cfg.MaxConcurrentPerDevice)
	assert.Equal(t, 1000, cfg.ValueBufferSize)
	assert.NotZero(t, cfg.DefaultRetryDelay)
}

func TestScheduler_AddAndRemoveDevice(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	device := &schema.Device{
		ID:               "dev-1",
		Protocol:         schema.ProtocolModbusTCP,
		ConnectionConfig: `{"host":"127.0.0.1","port":502}`,
	}

	s.AddDevice(device)

	s.mu.RLock()
	_, hasCfg := s.deviceConfigs[device.ID]
	_, hasLock := s.deviceLocks[device.ID]
	_, hasBreaker := s.deviceBreakers[device.ID]
	s.mu.RUnlock()

	assert.True(t, hasCfg)
	assert.True(t, hasLock)
	assert.True(t, hasBreaker)

	s.RemoveDevice(device.ID)

	s.mu.RLock()
	_, hasCfg = s.deviceConfigs[device.ID]
	_, hasBreaker = s.deviceBreakers[device.ID]
	s.mu.RUnlock()

	assert.False(t, hasCfg)
	assert.False(t, hasBreaker)
}

func TestScheduler_AddAndRemovePoint(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	groupID := "group-1"
	point := &schema.Point{
		ID:             "point-1",
		DeviceID:       "dev-1",
		Address:        "40001",
		Function:       "03",
		DataType:       schema.DataTypeInt16,
		PollingGroupID: &groupID,
	}

	s.AddPoint(point)

	s.mu.RLock()
	info, exists := s.pointInfos[point.ID]
	s.mu.RUnlock()

	require.True(t, exists)
	assert.Equal(t, groupID, info.PollingGroupID)
	assert.Equal(t, "dev-1", info.DeviceID)

	s.RemovePoint(point.ID)

	s.mu.RLock()
	_, exists = s.pointInfos[point.ID]
	s.mu.RUnlock()
	assert.False(t, exists)
}

func TestScheduler_ReplacingPollingGroupStopsPreviousTicker(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	group := &schema.PollingGroup{ID: "group-replace", Enabled: true, IntervalMs: 20}
	require.NoError(t, s.Start([]*schema.PollingGroup{group}))
	old := s.groupTickers[group.ID]
	require.NotNil(t, old)

	s.AddPollingGroup(group)
	select {
	case <-old.done:
	case <-time.After(time.Second):
		t.Fatal("replaced polling group ticker did not stop")
	}
	s.Stop()
}

func TestScheduler_DisablingPollingGroupStopsTicker(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	group := &schema.PollingGroup{ID: "group-disable", Enabled: true, IntervalMs: 20}
	require.NoError(t, s.Start([]*schema.PollingGroup{group}))
	active := s.groupTickers[group.ID]
	require.NotNil(t, active)

	s.AddPollingGroup(&schema.PollingGroup{ID: group.ID, Enabled: false, IntervalMs: group.IntervalMs})
	select {
	case <-active.done:
	case <-time.After(time.Second):
		t.Fatal("disabled polling group ticker did not stop")
	}
	s.Stop()
}

func TestScheduler_ResetDeviceBreaker(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	assert.False(t, s.ResetDeviceBreaker("missing"))

	s.AddDevice(&schema.Device{
		ID:               "dev-2",
		Protocol:         schema.ProtocolModbusTCP,
		ConnectionConfig: `{"host":"127.0.0.1"}`,
	})
	assert.True(t, s.ResetDeviceBreaker("dev-2"))
}

func TestScheduler_PollNow_RespectsDeadBreaker(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.BreakerConfig = health.BreakerConfig{
		ErrorThreshold:    0.2,
		WindowSize:        3,
		UnstableThreshold: 0.5,
		CooldownPeriod:    time.Hour,
	}

	s := NewScheduler(cfg, nil)
	s.AddDevice(&schema.Device{
		ID:               "dev-breaker",
		Protocol:         schema.ProtocolModbusTCP,
		ConnectionConfig: `{"host":"127.0.0.1","port":502}`,
	})
	s.AddPoint(&schema.Point{
		ID:       "point-breaker",
		DeviceID: "dev-breaker",
		Address:  "40001",
		Function: "03",
		DataType: schema.DataTypeInt16,
	})

	for i := 0; i < 3; i++ {
		s.deviceBreakers["dev-breaker"].ReportResult(errors.New("timeout"))
	}

	results := s.PollNow([]string{"point-breaker"})
	require.Len(t, results, 1)
	assert.Equal(t, schema.QualityBad, results[0].Quality)
	assert.Contains(t, results[0].Error, "設備熔斷中")
}
