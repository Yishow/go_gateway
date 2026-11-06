package runtime

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-gateway/internal/datalink/anomaly"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/schema"
)

// anomalyScanInterval is the detection period (design condition).
const anomalyScanInterval = time.Second

// observeAcquisitionAnomaly hands an accepted acquisition result to the
// detector. It only updates memory and never blocks on storage.
func (s *Service) observeAcquisitionAnomaly(cv collector.CollectedValue, quality schema.QualityFlag) {
	if s.anomalies == nil {
		return
	}
	s.anomalies.Observe(anomaly.Observation{
		DeviceID: cv.DeviceID,
		PointID:  cv.PointID,
		At:       cv.Timestamp,
		OK:       quality == schema.QualityGood && cv.Error == "",
		Value:    cv.Value,
		Reason:   cv.QualityReason,
	})
}

// anomalyLoop scans the runtime's latest point state once per interval. It
// never issues device reads. Storage failures are logged when they change and
// retried on the next tick; they are never reported as recorded events.
func (s *Service) anomalyLoop() {
	ticker := time.NewTicker(anomalyScanInterval)
	defer ticker.Stop()
	loaded := false
	lastErr := ""
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
		}
		err := s.runAnomalyScan(&loaded)
		message := ""
		if err != nil {
			message = err.Error()
		}
		if message != lastErr {
			if err != nil {
				log.Printf("acquisition anomaly detection failed: %v", err)
			} else if lastErr != "" {
				log.Printf("acquisition anomaly detection recovered")
			}
			lastErr = message
		}
	}
}

func (s *Service) runAnomalyScan(loaded *bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*anomalyScanInterval)
	defer cancel()
	if !*loaded {
		// Adopt events left active by a previous run before opening new ones.
		if err := s.anomalies.Load(ctx); err != nil {
			return err
		}
		*loaded = true
	}
	input, err := s.anomalyScanInput(ctx, s.scheduler.CollectionStatus())
	if err != nil {
		return err
	}
	return s.anomalies.Scan(ctx, input)
}

// anomalyScanInput resolves the scheduled points, each device's thresholds and
// breaker state. A failed configuration read aborts the scan instead of
// evaluating with guessed thresholds.
func (s *Service) anomalyScanInput(ctx context.Context, collection map[string]collector.CollectionPointStatus) (anomaly.ScanInput, error) {
	input := anomaly.ScanInput{Monitored: map[string]string{}, Thresholds: map[string]anomaly.Thresholds{}, BreakerOpen: map[string]bool{}}
	for pointID, status := range collection {
		if status.Running && status.Readable && status.DeviceID != "" {
			input.Monitored[pointID] = status.DeviceID
		}
	}
	if len(input.Monitored) == 0 {
		return input, nil
	}
	devices, err := s.listRuntimeStatusDevices(ctx, "")
	if err != nil {
		return anomaly.ScanInput{}, fmt.Errorf("list devices for anomaly thresholds: %w", err)
	}
	groups, err := s.listRuntimeStatusGroups(ctx)
	if err != nil {
		return anomaly.ScanInput{}, fmt.Errorf("list polling groups for anomaly thresholds: %w", err)
	}
	for _, dev := range devices {
		if dev == nil {
			continue
		}
		// Owned-group interval, else scheduler interval; auto adjustment can only
		// lengthen it, so the longer effective interval avoids false stale events.
		_, effectiveMs := s.resolveDeviceIntervalStatus(dev, nil, groups)
		input.Thresholds[dev.ID] = anomaly.Thresholds{
			Stale:  anomaly.StaleThreshold(dev.StaleThresholdMs, effectiveMs),
			Frozen: anomaly.FrozenThreshold(dev.FrozenThresholdMs),
		}
		input.BreakerOpen[dev.ID] = s.runtimeBreakerState(dev.ID) == breakerStateOpen
	}
	return input, nil
}
