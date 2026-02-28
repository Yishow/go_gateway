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

	s.mappingMu.RLock()
	bindings := s.mappingIndex[cv.PointID]
	s.mappingMu.RUnlock()
	if len(bindings) == 0 {
		return
	}

	rawValue := buildRawValue(cv)
	for _, b := range bindings {
		finalValue := cv.Value
		if shouldRunPipeline(b.TransformPipeline) {
			out, err := mapping.ExecutePipeline(cv.Value, b.TransformPipeline)
			if err != nil {
				s.mappingError.Add(1)
				continue
			}
			finalValue = out.CurrentValue
		}

		record := storage.ValueToRecord(b.TagID, finalValue, rawValue, cv.Timestamp, quality, b.TagDataType)
		if err := s.writer.Write(ctx, record); err != nil {
			s.writeError.Add(1)
			continue
		}
		s.writeSuccess.Add(1)
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
