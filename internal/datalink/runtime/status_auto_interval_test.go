package runtime

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/schema"
)

func TestRuntimeStatusReportsConfiguredAndEffectiveIntervalMs(t *testing.T) {
	ctx := context.Background()
	devID := "dev-auto-runtime"
	groupID := "grp-auto-runtime"
	now := time.Now().UTC()

	dev := &schema.Device{
		ID:                      devID,
		Name:                    "Runtime Auto Device",
		Protocol:                schema.ProtocolModbusTCP,
		Status:                  schema.DeviceStatusActive,
		ConnectionConfig:        `{"host":"127.0.0.1","port":1}`,
		AcquisitionAutoInterval: true,
	}
	pt := &schema.Point{
		ID:             "pt-auto-runtime",
		DeviceID:       devID,
		Name:           "D0",
		Address:        "40001",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupID,
		Enabled:        true,
		LastReadAt:     &now,
	}
	grp := &schema.PollingGroup{
		ID:            groupID,
		Name:          "dev-auto-runtime-group",
		OwnerDeviceID: &devID,
		IntervalMs:    1000,
		Enabled:       true,
	}

	svc, err := NewService(Config{
		Writer: &mockWriter{},
		Snapshot: Snapshot{
			Devices:       []*schema.Device{dev},
			Points:        []*schema.Point{pt},
			PollingGroups: []*schema.PollingGroup{grp},
		},
	})
	require.NoError(t, err)
	require.NoError(t, svc.Start(ctx))
	t.Cleanup(func() {
		_ = svc.Stop(context.Background())
	})

	snap, err := svc.RuntimeStatusSnapshot(ctx, devID)
	require.NoError(t, err)
	require.Len(t, snap.Collectors, 1)
	require.Equal(t, 1000, snap.Collectors[0].ConfiguredIntervalMs)
	require.Equal(t, 1000, snap.Collectors[0].EffectiveIntervalMs)

	// Stretch via 1250ms scan -> effective 1500ms.
	svc.scheduler.RecordScanDuration(devID, 1250*time.Millisecond)
	snap, err = svc.RuntimeStatusSnapshot(ctx, devID)
	require.NoError(t, err)
	require.Len(t, snap.Collectors, 1)
	require.Equal(t, 1000, snap.Collectors[0].ConfiguredIntervalMs)
	require.Equal(t, 1500, snap.Collectors[0].EffectiveIntervalMs)

	raw, err := json.Marshal(snap.Collectors[0])
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	require.Equal(t, float64(1000), payload["configured_interval_ms"])
	require.Equal(t, float64(1500), payload["effective_interval_ms"])

	// 5 consecutive on-time scans shrink effective interval back to 1000ms.
	for range 5 {
		svc.scheduler.RecordScanDuration(devID, 700*time.Millisecond)
	}
	snap, err = svc.RuntimeStatusSnapshot(ctx, devID)
	require.NoError(t, err)
	require.Len(t, snap.Collectors, 1)
	require.Equal(t, 1000, snap.Collectors[0].ConfiguredIntervalMs)
	require.Equal(t, 1000, snap.Collectors[0].EffectiveIntervalMs)
}
