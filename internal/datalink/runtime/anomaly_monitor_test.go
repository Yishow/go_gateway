package runtime

import (
	"context"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/anomaly"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

// memoryAnomalyRepo is a minimal in-memory anomaly.Repository for runtime tests.
type memoryAnomalyRepo struct {
	mu     sync.Mutex
	events []anomaly.Event
}

func (r *memoryAnomalyRepo) Open(_ context.Context, ev anomaly.Event) (anomaly.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.events {
		if existing.EndedAt == nil && existing.DeviceID == ev.DeviceID && existing.PointID == ev.PointID && existing.Kind == ev.Kind {
			return anomaly.Event{}, anomaly.ErrAlreadyActive
		}
	}
	ev.ID = int64(len(r.events) + 1)
	r.events = append(r.events, ev)
	return ev, nil
}

func (r *memoryAnomalyRepo) End(_ context.Context, id int64, endedAt time.Time, _ *string, _ string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.events {
		if r.events[i].ID == id && r.events[i].EndedAt == nil {
			r.events[i].EndedAt = &endedAt
			return nil
		}
	}
	return anomaly.ErrNotActive
}

func (r *memoryAnomalyRepo) ListActive(context.Context) ([]anomaly.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var active []anomaly.Event
	for _, ev := range r.events {
		if ev.EndedAt == nil {
			active = append(active, ev)
		}
	}
	return active, nil
}

func (r *memoryAnomalyRepo) List(context.Context, anomaly.ListFilter) ([]anomaly.Event, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]anomaly.Event(nil), r.events...), len(r.events), nil
}

func (r *memoryAnomalyRepo) ActiveCounts(context.Context) (map[string]int, error) {
	return map[string]int{}, nil
}

func TestCollectedConnectionFailuresFeedAnomalyDetector(t *testing.T) {
	repo := &memoryAnomalyRepo{}
	at := time.Date(2026, 11, 5, 8, 0, 0, 0, time.UTC)
	now := at
	s := targetOutcomeService(&mockWriter{}, nil)
	s.anomalies = anomaly.NewDetector(repo, func() time.Time { return now })
	input := anomaly.ScanInput{Monitored: map[string]string{"p1": "dev-1"}, Thresholds: map[string]anomaly.Thresholds{"dev-1": {Stale: time.Hour}}}

	for i := range 30 {
		s.handleCollectedValue(t.Context(), collector.CollectedValue{PointID: "p1", DeviceID: "dev-1", Timestamp: at.Add(time.Duration(i) * time.Second), Quality: schema.QualityBad, Error: "dial tcp 10.0.0.9:502: i/o timeout", QualityReason: "connection-failed"})
		now = at.Add(time.Duration(i+1) * time.Second)
		require.NoError(t, s.anomalies.Scan(t.Context(), input))
	}
	recovered := at.Add(30 * time.Second)
	s.handleCollectedValue(t.Context(), collector.CollectedValue{PointID: "p1", DeviceID: "dev-1", Value: int16(1), Timestamp: recovered, Quality: schema.QualityGood})
	now = recovered.Add(time.Second)
	require.NoError(t, s.anomalies.Scan(t.Context(), input))

	events, _, err := repo.List(t.Context(), anomaly.ListFilter{})
	require.NoError(t, err)
	require.Len(t, events, 1, "30 seconds of connection failures record one disconnected event")
	require.Equal(t, anomaly.KindDisconnected, events[0].Kind)
	require.Equal(t, at, events[0].StartedAt)
	require.Equal(t, recovered, *events[0].EndedAt)
	require.Equal(t, anomaly.ReasonConnectionFailed, events[0].Detail, "raw connection errors never reach the anomaly detail")
}

func TestRejectedSamplesDoNotFeedAnomalyDetector(t *testing.T) {
	repo := &memoryAnomalyRepo{}
	s := targetOutcomeService(&mockWriter{}, nil)
	s.scheduler = collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	s.anomalies = anomaly.NewDetector(repo, time.Now)
	s.handleCollectedValue(t.Context(), collector.CollectedValue{PointID: "p1", DeviceID: "dev-1", Timestamp: time.Now(), Quality: schema.QualityBad, QualityReason: "connection-failed"})
	require.NoError(t, s.anomalies.Scan(t.Context(), anomaly.ScanInput{Monitored: map[string]string{"p1": "dev-1"}, Thresholds: map[string]anomaly.Thresholds{"dev-1": {Stale: time.Hour}}}))
	require.Empty(t, repo.events, "a sample from a replaced installation is not evidence of a disconnect")
}

func TestAnomalyScanInputUsesDeviceThresholdsAndScheduledPoints(t *testing.T) {
	ctx := t.Context()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	devRepo := device.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	staleB := 2500
	require.NoError(t, devRepo.Create(ctx, &schema.Device{ID: "dev-a", Name: "A", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: "{}"}))
	require.NoError(t, devRepo.Create(ctx, &schema.Device{ID: "dev-b", Name: "B", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: "{}", StaleThresholdMs: &staleB, FrozenThresholdMs: 600000}))
	owner := "dev-a"
	require.NoError(t, groupRepo.Create(ctx, &schema.PollingGroup{ID: "grp-a", Name: "device:dev-a", OwnerDeviceID: &owner, IntervalMs: 3000, Enabled: true}))
	svc, err := NewService(DefaultConfig(), Dependencies{
		Scheduler:           collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr),
		Writer:              storage.NewMemoryStorage(16),
		DeviceService:       device.NewService(devRepo, connMgr),
		PointService:        point.NewService(point.NewMemoryRepository(), groupRepo),
		MappingService:      mapping.NewService(mapping.NewMemoryRepository()),
		TagService:          tag.NewService(tag.NewMemoryRepository()),
		PollingGroupService: pollinggroup.NewService(groupRepo),
		AnomalyRepository:   &memoryAnomalyRepo{},
	})
	require.NoError(t, err)
	require.NotNil(t, svc.anomalies, "an anomaly repository enables detection")

	input, err := svc.anomalyScanInput(ctx, map[string]collector.CollectionPointStatus{
		"pa":     {DeviceID: "dev-a", Running: true, Readable: true},
		"pb":     {DeviceID: "dev-b", Running: true, Readable: true},
		"p-idle": {DeviceID: "dev-a", Readable: true},
	})
	require.NoError(t, err)
	require.Equal(t, map[string]string{"pa": "dev-a", "pb": "dev-b"}, input.Monitored, "only scheduled points are monitored")
	require.Equal(t, anomaly.Thresholds{Stale: 15 * time.Second}, input.Thresholds["dev-a"], "3000ms owned-group interval x5")
	require.Equal(t, anomaly.Thresholds{Stale: 2500 * time.Millisecond, Frozen: 600 * time.Second}, input.Thresholds["dev-b"])
	require.False(t, input.BreakerOpen["dev-a"])
}

type loadCountingAnomalyRepo struct {
	memoryAnomalyRepo
	loads chan struct{}
}

func (r *loadCountingAnomalyRepo) ListActive(ctx context.Context) ([]anomaly.Event, error) {
	select {
	case r.loads <- struct{}{}:
	default:
	}
	return r.memoryAnomalyRepo.ListActive(ctx)
}

func TestRuntimeStartRunsAnomalyLoopUntilStop(t *testing.T) {
	repo := &loadCountingAnomalyRepo{loads: make(chan struct{}, 1)}
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	groupRepo := pollinggroup.NewMemoryRepository()
	svc, err := NewService(DefaultConfig(), Dependencies{
		Scheduler:           collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr),
		Writer:              storage.NewMemoryStorage(16),
		DeviceService:       device.NewService(device.NewMemoryRepository(), connMgr),
		PointService:        point.NewService(point.NewMemoryRepository(), groupRepo),
		MappingService:      mapping.NewService(mapping.NewMemoryRepository()),
		TagService:          tag.NewService(tag.NewMemoryRepository()),
		PollingGroupService: pollinggroup.NewService(groupRepo),
		AnomalyRepository:   repo,
	})
	require.NoError(t, err)
	require.NoError(t, svc.Start(t.Context()))
	select {
	case <-repo.loads:
	case <-time.After(3 * anomalyScanInterval):
		t.Fatal("anomaly loop did not adopt active events after start")
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, svc.Stop(stopCtx), "the anomaly worker exits with the runtime")
}
