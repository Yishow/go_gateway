package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestService_RuntimeStatusSnapshot_MapsBreakerStateToOpen(t *testing.T) {
	cfg := collector.DefaultSchedulerConfig()
	cfg.BreakerConfig = health.BreakerConfig{
		ErrorThreshold:    0.2,
		UnstableThreshold: 0.5,
		WindowSize:        3,
		CooldownPeriod:    time.Hour,
	}

	scheduler := collector.NewScheduler(cfg, nil)
	deviceRecord := &schema.Device{
		ID:               "device-1",
		Name:             "Breaker Device",
		Protocol:         schema.ProtocolModbusTCP,
		ConnectionConfig: `{"host":"127.0.0.1","port":1,"slave_id":1,"timeout":1}`,
	}
	groupID := "group-1"
	pointRecord := &schema.Point{
		ID:             "point-1",
		DeviceID:       deviceRecord.ID,
		Name:           "Pressure",
		Address:        "40001",
		Function:       "03",
		DataType:       schema.DataTypeInt16,
		PollingGroupID: &groupID,
		Enabled:        true,
	}

	scheduler.AddDevice(deviceRecord)
	scheduler.AddPoint(pointRecord)
	for attempt := 0; attempt < 3; attempt++ {
		scheduler.PollNow([]string{pointRecord.ID})
	}

	svc := &Service{
		config:    Config{UpdatePointState: false},
		writer:    &mockWriter{},
		scheduler: scheduler,
		snapshot: Snapshot{
			Devices: []*schema.Device{deviceRecord},
			Points:  []*schema.Point{pointRecord},
			PollingGroups: []*schema.PollingGroup{
				{ID: groupID, Name: "runtime", IntervalMs: 1000, Enabled: true},
			},
		},
	}
	svc.running.Store(true)

	snapshot, err := svc.RuntimeStatusSnapshot(context.Background(), deviceRecord.ID)
	require.NoError(t, err)
	require.Len(t, snapshot.Collectors, 1)
	require.Equal(t, "open", snapshot.Collectors[0].BreakerState)
	require.Equal(t, "error", snapshot.Collectors[0].Status)
}

func TestService_StatusEventMatchesAlignedRuntimeSnapshot(t *testing.T) {
	ctx := t.Context()

	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	groupSvc := pollinggroup.NewService(pollinggroup.NewMemoryRepository())

	deviceRecord, err := deviceSvc.Create(ctx, device.CreateDeviceRequest{
		Name:     "Aligned Device",
		Protocol: schema.ProtocolModbusTCP,
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	})
	require.NoError(t, err)
	deviceRecord.Status = schema.DeviceStatusActive
	require.NoError(t, deviceRepo.Update(ctx, deviceRecord))

	groupRecord, err := groupSvc.Create(ctx, pollinggroup.CreateRequest{
		Name:       "aligned-fast",
		IntervalMs: 100,
		Priority:   1,
	})
	require.NoError(t, err)

	pointRecord, err := pointSvc.Create(ctx, point.CreatePointRequest{
		DeviceID:       deviceRecord.ID,
		Name:           "Flow",
		Address:        "40001",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupRecord.ID,
	})
	require.NoError(t, err)

	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	scheduler.AddDevice(deviceRecord)
	scheduler.AddPoint(pointRecord)

	svc := &Service{
		config: Config{UpdatePointState: true},
		writer: &mockWriter{},

		scheduler: scheduler,
		deviceSvc: deviceSvc,
		pointSvc:  pointSvc,
		groupSvc:  groupSvc,

		mappingIndex:      map[string][]mappingBinding{},
		pointMetaIndex:    map[string]pointMeta{pointRecord.ID: {DeviceID: deviceRecord.ID, Address: pointRecord.Address}},
		subscribers:       make(map[int64]valueSubscriber),
		statusSubscribers: make(map[int64]statusSubscriber),
		lastStatuses:      make(map[string]DeviceRuntimeStatus),
		stopCh:            make(chan struct{}),
	}
	svc.running.Store(true)

	svc.wg.Add(1)
	go svc.statusLoop()
	t.Cleanup(func() {
		close(svc.stopCh)
		svc.wg.Wait()
	})

	svc.handleCollectedValue(ctx, collector.CollectedValue{
		PointID:   pointRecord.ID,
		DeviceID:  deviceRecord.ID,
		Value:     12.5,
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
	})

	statusStream, unsubscribe := svc.SubscribeStatusEvents(deviceRecord.ID)
	defer unsubscribe()

	select {
	case initial := <-statusStream:
		require.Equal(t, "running", initial.Status)
		require.Equal(t, 1, initial.PointsTotal)
		require.Equal(t, 1, initial.PointsHealthy)
		require.Equal(t, "closed", initial.BreakerState)
	case <-time.After(time.Second):
		t.Fatal("expected initial runtime status event")
	}

	var staleEvent DeviceStatusEvent
	select {
	case staleEvent = <-statusStream:
	case <-time.After(2 * time.Second):
		t.Fatal("expected stale runtime status event")
	}

	require.Equal(t, "warning", staleEvent.Status)
	require.Equal(t, 1, staleEvent.PointsTotal)
	require.Equal(t, 0, staleEvent.PointsHealthy)
	require.Equal(t, 1, staleEvent.PointsStale)
	require.Equal(t, 0, staleEvent.PointsError)
	require.Equal(t, "closed", staleEvent.BreakerState)

	snapshot, err := svc.RuntimeStatusSnapshot(ctx, deviceRecord.ID)
	require.NoError(t, err)
	require.Len(t, snapshot.Collectors, 1)
	require.Equal(t, snapshot.Collectors[0].Status, staleEvent.Status)
	require.Equal(t, snapshot.Collectors[0].PointsTotal, staleEvent.PointsTotal)
	require.Equal(t, snapshot.Collectors[0].PointsHealthy, staleEvent.PointsHealthy)
	require.Equal(t, snapshot.Collectors[0].PointsStale, staleEvent.PointsStale)
	require.Equal(t, snapshot.Collectors[0].PointsError, staleEvent.PointsError)
	require.Equal(t, snapshot.Collectors[0].BreakerState, staleEvent.BreakerState)
}
