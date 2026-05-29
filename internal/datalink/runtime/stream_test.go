package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestService_SubscribeValueEvents_BroadcastsMatchingPoint(t *testing.T) {
	svc := &Service{
		config:       Config{UpdatePointState: false},
		writer:       &mockWriter{},
		mappingIndex: map[string][]mappingBinding{},
		pointMetaIndex: map[string]pointMeta{
			"point-1": {
				DeviceID: "device-1",
				Address:  "40001",
			},
			"point-2": {
				DeviceID: "device-2",
				Address:  "40002",
			},
		},
	}

	stream, unsubscribe := svc.SubscribeValueEvents("device-1", []string{"point-1"})
	defer unsubscribe()

	ts := time.Date(2026, 3, 16, 6, 0, 0, 0, time.UTC)
	svc.handleCollectedValue(context.Background(), collector.CollectedValue{
		PointID:   "point-1",
		DeviceID:  "device-1",
		Value:     37.5,
		Timestamp: ts,
		Quality:   schema.QualityGood,
	})

	select {
	case evt := <-stream:
		require.Equal(t, "device-1", evt.DeviceID)
		require.Equal(t, "point-1", evt.PointID)
		require.Equal(t, "40001", evt.Address)
		require.Equal(t, 37.5, evt.RawValue)
		require.Equal(t, 37.5, evt.TransformedValue)
		require.Equal(t, schema.QualityGood, evt.Quality)
		require.False(t, evt.Stale)
		require.Equal(t, ts, evt.Timestamp)
	case <-time.After(time.Second):
		t.Fatal("expected runtime value event")
	}

	svc.handleCollectedValue(context.Background(), collector.CollectedValue{
		PointID:   "point-2",
		DeviceID:  "device-2",
		Value:     88.8,
		Timestamp: ts,
		Quality:   schema.QualityGood,
	})

	select {
	case evt := <-stream:
		t.Fatalf("unexpected event for non-matching point: %+v", evt)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestService_SubscribeValueEvents_LossTolerantWithSlowSubscriber(t *testing.T) {
	svc := &Service{
		config: Config{UpdatePointState: false},
		writer: &mockWriter{},
	}

	slowStream, unsubscribeSlow := svc.SubscribeValueEvents("device-1", []string{"point-1"})
	defer unsubscribeSlow()

	fastStream, unsubscribeFast := svc.SubscribeValueEvents("device-1", []string{"point-1"})
	defer unsubscribeFast()

	startedAt := time.Now()
	for index := 0; index < 64; index++ {
		svc.broadcastValueEvent(ValueEvent{
			DeviceID:         "device-1",
			PointID:          "point-1",
			Address:          "40001",
			RawValue:         index,
			TransformedValue: index,
			Quality:          schema.QualityGood,
			Stale:            false,
			Timestamp:        time.Date(2026, 3, 16, 6, 0, index, 0, time.UTC),
		})
	}

	if time.Since(startedAt) > 200*time.Millisecond {
		t.Fatalf("expected non-blocking broadcast, took %s", time.Since(startedAt))
	}

	select {
	case evt := <-fastStream:
		require.Equal(t, "device-1", evt.DeviceID)
		require.Equal(t, "point-1", evt.PointID)
	case <-time.After(time.Second):
		t.Fatal("expected fast subscriber to keep receiving events")
	}

	select {
	case <-slowStream:
	case <-time.After(time.Second):
		t.Fatal("expected slow subscriber channel to remain open")
	}
}
