package collector

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/schema"
)

func TestSchedulerAutoIntervalStretchesAndShrinksBack(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	devID := "dev-auto-1"
	groupID := "grp-auto-1"

	s.AddDevice(&schema.Device{
		ID:                      devID,
		Name:                    "Auto Interval PLC",
		Protocol:                schema.ProtocolModbusTCP,
		Status:                  schema.DeviceStatusActive,
		AcquisitionAutoInterval: true,
	})
	s.AddPoint(&schema.Point{
		ID:             "pt-1",
		DeviceID:       devID,
		Name:           "D0",
		Address:        "40001",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupID,
		Enabled:        true,
	})
	s.AddPollingGroup(&schema.PollingGroup{
		ID:            groupID,
		Name:          "device-auto-group",
		OwnerDeviceID: &devID,
		IntervalMs:    1000,
		Enabled:       true,
	})

	status, ok := s.DeviceIntervalStatus(devID)
	require.True(t, ok)
	require.Equal(t, 1000, status.ConfiguredIntervalMs)
	require.Equal(t, 1000, status.EffectiveIntervalMs)
	require.True(t, status.AutoInterval)

	var logBuf bytes.Buffer
	origOut := log.Writer()
	log.SetOutput(&logBuf)
	t.Cleanup(func() { log.SetOutput(origOut) })

	// Slow scan of 1250ms -> max(1000, ceil_100ms(1250 * 1.2)) = 1500ms.
	effective := s.RecordScanDuration(devID, 1250*time.Millisecond)
	require.Equal(t, 1500, effective)
	status, ok = s.DeviceIntervalStatus(devID)
	require.True(t, ok)
	require.Equal(t, 1000, status.ConfiguredIntervalMs)
	require.Equal(t, 1500, status.EffectiveIntervalMs)
	require.Contains(t, logBuf.String(), devID)
	require.Contains(t, logBuf.String(), "1500")

	// Next 4 on-time scans keep effective interval at 1500ms.
	for i := range 4 {
		effective = s.RecordScanDuration(devID, 700*time.Millisecond)
		require.Equalf(t, 1500, effective, "on-time scan %d should keep 1500ms before 5th consecutive on-time scan", i+1)
	}

	// 5th consecutive on-time scan shrinks effective interval back to 1000ms.
	effective = s.RecordScanDuration(devID, 700*time.Millisecond)
	require.Equal(t, 1000, effective)
	status, ok = s.DeviceIntervalStatus(devID)
	require.True(t, ok)
	require.Equal(t, 1000, status.ConfiguredIntervalMs)
	require.Equal(t, 1000, status.EffectiveIntervalMs)
}

// TestSchedulerAutoIntervalMultiStepShrink documents the step-back rule:
// After every 5 consecutive on-time scans (elapsed <= configured_interval_ms),
// effective_interval_ms shrinks toward configured_interval_ms by
// max(500ms, ceil_100ms((effective - configured) / 2)) (which is >= max(100ms, (effective-configured)/2)
// rounded to 100ms multiples and returns 1500ms -> 1000ms in one 5-scan window).
func TestSchedulerAutoIntervalMultiStepShrink(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	devID := "dev-auto-step"
	groupID := "grp-auto-step"

	s.AddDevice(&schema.Device{
		ID:                      devID,
		Protocol:                schema.ProtocolModbusTCP,
		Status:                  schema.DeviceStatusActive,
		AcquisitionAutoInterval: true,
	})
	s.AddPollingGroup(&schema.PollingGroup{
		ID:            groupID,
		Name:          "device-step-group",
		OwnerDeviceID: &devID,
		IntervalMs:    1000,
		Enabled:       true,
	})

	// 2500ms * 1.2 = 3000ms.
	require.Equal(t, 3000, s.RecordScanDuration(devID, 2500*time.Millisecond))

	// Step 1 after 5 consecutive on-time scans: 3000 - max(500, ceil100((3000-1000)/2=1000)) = 2000ms.
	for range 4 {
		require.Equal(t, 3000, s.RecordScanDuration(devID, 600*time.Millisecond))
	}
	require.Equal(t, 2000, s.RecordScanDuration(devID, 600*time.Millisecond))

	// Step 2 after 5 more on-time scans: 2000 - max(500, ceil100((2000-1000)/2=500)) = 1500ms.
	for range 4 {
		require.Equal(t, 2000, s.RecordScanDuration(devID, 600*time.Millisecond))
	}
	require.Equal(t, 1500, s.RecordScanDuration(devID, 600*time.Millisecond))

	// Step 3 after 5 more on-time scans: 1500 - max(500, ceil100((1500-1000)/2=300)) = 1000ms.
	for range 4 {
		require.Equal(t, 1500, s.RecordScanDuration(devID, 600*time.Millisecond))
	}
	require.Equal(t, 1000, s.RecordScanDuration(devID, 600*time.Millisecond))
}

func TestSchedulerAutoIntervalDropsExpiredOutlier(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	deviceID := "device-window"
	s.AddDevice(&schema.Device{ID: deviceID, Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, AcquisitionAutoInterval: true})
	s.AddPollingGroup(&schema.PollingGroup{ID: "window", Name: "window", OwnerDeviceID: &deviceID, IntervalMs: 1000, Enabled: true})
	require.Equal(t, 3000, s.RecordScanDuration(deviceID, 2500*time.Millisecond))
	for range 4 {
		require.Equal(t, 3000, s.RecordScanDuration(deviceID, 1250*time.Millisecond))
	}
	require.Equal(t, 1500, s.RecordScanDuration(deviceID, 1250*time.Millisecond))
}

func TestSchedulerManualFixedIntervalKeepsConfiguredAndReportsOverrun(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	devID := "dev-fixed"
	groupID := "grp-fixed"

	s.AddDevice(&schema.Device{
		ID:                      devID,
		Protocol:                schema.ProtocolModbusTCP,
		Status:                  schema.DeviceStatusActive,
		AcquisitionAutoInterval: false,
	})
	s.AddPollingGroup(&schema.PollingGroup{
		ID:            groupID,
		Name:          "device-fixed-group",
		OwnerDeviceID: &devID,
		IntervalMs:    1000,
		Enabled:       true,
	})

	effective := s.RecordScanDuration(devID, 1250*time.Millisecond)
	require.Equal(t, 1000, effective)

	status, ok := s.DeviceIntervalStatus(devID)
	require.True(t, ok)
	require.Equal(t, 1000, status.ConfiguredIntervalMs)
	require.Equal(t, 1000, status.EffectiveIntervalMs)
	require.False(t, status.AutoInterval)
}

func TestSchedulerRemoveDeviceOwnedPollingGroupsStopsOnlyOwnedTickers(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	s.mu.Lock()
	s.running = true // 只建立 ticker，不啟動完整排程迴圈。
	s.mu.Unlock()
	owner := "dev-a"
	s.AddPollingGroup(&schema.PollingGroup{ID: "owned", IntervalMs: 1000, Enabled: true, OwnerDeviceID: &owner})
	s.AddPollingGroup(&schema.PollingGroup{ID: "shared", IntervalMs: 1000, Enabled: true})
	if _, ok := s.PollingGroup("owned"); !ok {
		t.Fatal("precondition: owned ticker should be running")
	}

	s.RemoveDeviceOwnedPollingGroups(owner)

	if _, ok := s.PollingGroup("owned"); ok {
		t.Fatal("expected device-owned polling group ticker to be removed")
	}
	s.mu.RLock()
	_, ownedConfigured := s.configuredGroups["owned"]
	s.mu.RUnlock()
	if ownedConfigured {
		t.Fatal("expected device-owned polling group configuration to be removed")
	}
	if _, ok := s.PollingGroup("shared"); !ok {
		t.Fatal("expected shared polling group ticker to remain")
	}
	s.RemovePollingGroup("shared")
}
