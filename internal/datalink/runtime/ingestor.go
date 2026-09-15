package runtime

import (
	"context"
	"encoding/json"
	"errors"
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
	deviceID := cv.DeviceID
	if deviceID == "" {
		deviceID = meta.DeviceID
	}
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
			if outcomeWriter, ok := s.target.(TargetOutcomeWriter); ok {
				for _, outcome := range outcomeWriter.WriteTagValueOutcomes(ctx, b.TagID, finalValue, cv.Timestamp) {
					s.recordTargetDeliveryOutcome(deviceID, cv.PointID, b.TagID, outcome, cv.Timestamp)
				}
			} else if err := s.target.WriteTagValue(ctx, b.TagID, finalValue, cv.Timestamp); err != nil {
				s.writeError.Add(1)
				var deliveryErr *DeliveryError
				if errors.As(err, &deliveryErr) && deliveryErr.Target == "modbus_share" {
					s.recordModbusShareDeliveryDiagnostic(ModbusShareDeliveryDiagnostic{DeviceID: deviceID, PointID: cv.PointID, TagID: b.TagID, Status: ModbusShareDeliveryStatusFailed, Stage: deliveryErr.Stage, Error: deliveryErr.Error(), ObservedAt: cv.Timestamp})
					continue
				}
				s.recordDatabaseDeliveryDiagnostic(DatabaseDeliveryDiagnostic{
					DeviceID:    deviceID,
					PointID:     cv.PointID,
					TagID:       b.TagID,
					Status:      DatabaseDeliveryStatusFailed,
					Stages:      databaseDeliveryFailureStages(),
					FailedStage: DatabaseDeliveryStageDBWrite,
					Error:       err.Error(),
					ObservedAt:  cv.Timestamp,
				})
			} else {
				s.recordDatabaseDeliveryDiagnostic(DatabaseDeliveryDiagnostic{
					DeviceID:   deviceID,
					PointID:    cv.PointID,
					TagID:      b.TagID,
					Status:     DatabaseDeliveryStatusSucceeded,
					Stages:     databaseDeliverySuccessStages(),
					ObservedAt: cv.Timestamp,
				})
			}
		}
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
		s.publishDerivedStatus(ctx, deviceID)
	}
}

func (s *Service) recordTargetDeliveryOutcome(deviceID, pointID, tagID string, outcome TargetDeliveryOutcome, observedAt time.Time) {
	if outcome.Err != nil {
		s.writeError.Add(1)
	}
	switch outcome.Target {
	case TargetModbusShare:
		diagnostic := ModbusShareDeliveryDiagnostic{DeviceID: deviceID, PointID: pointID, TagID: tagID, Status: ModbusShareDeliveryStatusSucceeded, Stage: "modbus_write", ObservedAt: observedAt}
		if outcome.Err != nil {
			diagnostic.Status = ModbusShareDeliveryStatusFailed
			diagnostic.Error = outcome.Err.Error()
			var deliveryErr *DeliveryError
			if errors.As(outcome.Err, &deliveryErr) {
				diagnostic.Stage = deliveryErr.Stage
			}
		}
		s.recordModbusShareDeliveryDiagnostic(diagnostic)
	case TargetDatabase:
		diagnostic := DatabaseDeliveryDiagnostic{DeviceID: deviceID, PointID: pointID, TagID: tagID, Status: DatabaseDeliveryStatusSucceeded, Stages: databaseDeliverySuccessStages(), ObservedAt: observedAt}
		if outcome.Err != nil {
			diagnostic.Status = DatabaseDeliveryStatusFailed
			diagnostic.Stages = databaseDeliveryFailureStages()
			diagnostic.FailedStage = DatabaseDeliveryStageDBWrite
			diagnostic.Error = outcome.Err.Error()
		}
		s.recordDatabaseDeliveryDiagnostic(diagnostic)
	}
}

func databaseDeliveryFailureStages() []DatabaseDeliveryStage {
	return []DatabaseDeliveryStage{
		DatabaseDeliveryStageCollected,
		DatabaseDeliveryStageMapped,
		DatabaseDeliveryStageDBWriteFailed,
	}
}

func databaseDeliverySuccessStages() []DatabaseDeliveryStage {
	return []DatabaseDeliveryStage{
		DatabaseDeliveryStageCollected,
		DatabaseDeliveryStageMapped,
		DatabaseDeliveryStageDBWriteSucceeded,
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

func (s *Service) publishDerivedStatus(ctx context.Context, deviceID string) {
	if deviceID == "" {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	status, found, err := s.DeviceStatus(ctx, deviceID)
	if err != nil || !found {
		return
	}

	s.emitStatusIfChanged(status)
}
