package runtime

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
)

func (s *Service) handleCollectedValue(ctx context.Context, cv collector.CollectedValue) {
	if cv.Timestamp.IsZero() {
		cv.Timestamp = time.Now()
	}
	quality := cv.Quality
	if quality == "" {
		if cv.Error != "" {
			quality = schema.QualityBad
		} else {
			quality = schema.QualityGood
		}
	}

	if s.config.UpdatePointState && s.pointSvc != nil {
		if err := s.pointSvc.UpdateReadResult(ctx, cv.PointID, cv.Value, cv.Error); err != nil {
			s.pointStateError.Add(1)
		}
	}

	meta := s.lookupPointMeta(cv.PointID)
	transformedValue := cv.Value

	s.mappingMu.RLock()
	bindings := s.mappingIndex[cv.PointID]
	s.mappingMu.RUnlock()
	rawValue := buildRawValue(cv)
	for index, b := range bindings {
		finalValue := cv.Value
		if shouldRunPipeline(b.TransformPipeline) {
			out, err := mapping.ExecutePipeline(cv.Value, b.TransformPipeline)
			if err != nil {
				s.mappingError.Add(1)
				continue
			}
			finalValue = out.CurrentValue
		}
		if index == 0 {
			transformedValue = finalValue
		}

		record := storage.ValueToRecord(b.TagID, finalValue, rawValue, cv.Timestamp, quality, b.TagDataType)
		if err := s.writer.Write(ctx, record); err != nil {
			s.writeError.Add(1)
		} else {
			s.writeSuccess.Add(1)
		}

		if s.target != nil {
			if err := s.target.WriteTagValue(ctx, b.TagID, finalValue, cv.Timestamp); err != nil {
				s.writeError.Add(1)
			}
		}
	}

	deviceID := cv.DeviceID
	if deviceID == "" {
		deviceID = meta.DeviceID
	}
	s.broadcastValueEvent(ValueEvent{
		DeviceID:         deviceID,
		PointID:          cv.PointID,
		Address:          meta.Address,
		RawValue:         cv.Value,
		TransformedValue: transformedValue,
		Quality:          quality,
		Stale:            false,
		Timestamp:        cv.Timestamp,
	})
	if deviceID != "" {
		s.publishDerivedStatus(deviceID)
	}
}

func shouldRunPipeline(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "[]" {
		return false
	}
	return true
}

func buildRawValue(cv collector.CollectedValue) interface{} {
	if len(cv.RawBytes) == 0 {
		return cv.Value
	}
	payload := map[string]interface{}{
		"value":     cv.Value,
		"raw_bytes": cv.RawBytes,
	}
	buf, err := json.Marshal(payload)
	if err != nil {
		return cv.RawBytes
	}
	return string(buf)
}

func (s *Service) publishDerivedStatus(deviceID string) {
	if deviceID == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	status, found, err := s.DeviceStatus(ctx, deviceID)
	if err != nil || !found {
		return
	}

	s.emitStatusIfChanged(status)
}
