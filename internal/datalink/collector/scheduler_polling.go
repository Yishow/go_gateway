package collector

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

func (s *Scheduler) pollPointWithAcquisition(ctx context.Context, conn *connector.ManagedConnection, pt pointInfo, acquisitionID string) error {
	req := connector.ReadRequest{
		Address:    pt.Address,
		Function:   pt.Function,
		DataType:   pt.DataType,
		DataFormat: pt.DataFormat,
	}

	req.Count = schema.RegisterCountForDataType(pt.DataType)

	var result connector.ReadResult
	var err error

	// 重試邏輯 (指數退避)
	for attempt := 0; attempt <= s.config.DefaultRetryCount; attempt++ {
		if attempt > 0 {
			// 指數退避: delay * 2^attempt
			backoff := s.config.DefaultRetryDelay * time.Duration(1<<uint(attempt))
			if backoff > 30*time.Second {
				backoff = 30 * time.Second // 最大 30 秒
			}
			timer := time.NewTimer(backoff)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
		result, err = conn.Read(ctx, req)
		if err == nil && result.Quality == schema.QualityGood {
			break
		}
	}

	cv := s.collectedValue(managedConfig(conn), pt, result, err, acquisitionID, s.now())
	s.emitValue(cv)
	return err
}

func (s *Scheduler) collectedValue(
	cfg deviceConfig,
	pt pointInfo,
	result connector.ReadResult,
	readErr error,
	acquisitionID string,
	completedAt time.Time,
) CollectedValue {
	return s.collectedValueWithReason(cfg, pt, result, readErr, acquisitionID, completedAt, "")
}

func (s *Scheduler) collectedValueWithReason(
	cfg deviceConfig,
	pt pointInfo,
	result connector.ReadResult,
	readErr error,
	acquisitionID string,
	completedAt time.Time,
	safeReason string,
) CollectedValue {
	quality := result.Quality
	qualityReason := safeReason
	errText := result.Error
	if readErr != nil {
		quality = schema.QualityBad
		qualityReason = "read-failed"
		errText = readErr.Error()
	}
	switch {
	case !isKnownQuality(quality):
		quality = schema.QualityBad
		qualityReason = "quality-unknown"
	case errText != "" && quality == schema.QualityGood:
		quality = schema.QualityBad
		qualityReason = "read-failed"
	case quality != schema.QualityGood && qualityReason == "":
		qualityReason = fmt.Sprintf("quality-%s", quality)
	}

	sourceAt := result.Timestamp
	sourceOrigin := result.TimeOrigin
	if readErr != nil {
		// A failed read cannot establish a trusted source observation.
		sourceAt = time.Time{}
		sourceOrigin = ""
	}
	observedAt, timeOrigin := measurement.ResolveAcquisitionTime(sourceAt, sourceOrigin, completedAt)

	configFingerprint := ""
	if cfg.ID != "" {
		device := schema.Device{
			ID:               cfg.ID,
			Protocol:         cfg.Protocol,
			ConnectionConfig: cfg.Config,
		}
		point := schema.Point{
			ID:             pt.ID,
			DeviceID:       pt.DeviceID,
			Address:        pt.Address,
			Function:       pt.Function,
			DataType:       pt.DataType,
			DataFormat:     pt.DataFormat,
			Mode:           pt.Mode,
			PollingGroupID: pollingGroupPointer(pt.PollingGroupID),
		}
		if fingerprint, err := measurement.AcquisitionConfigFingerprint(device, point); err == nil {
			configFingerprint = fingerprint
		}
	}

	return CollectedValue{
		PointID:           pt.ID,
		DeviceID:          pt.DeviceID,
		Value:             result.Value,
		RawBytes:          result.RawBytes,
		Timestamp:         result.Timestamp,
		AcquisitionID:     acquisitionID,
		ObservedAt:        observedAt,
		ReceivedAt:        completedAt.UTC(),
		TimeOrigin:        string(timeOrigin),
		ConfigFingerprint: configFingerprint,
		Quality:           quality,
		Error:             errText,
		QualityReason:     qualityReason,
	}
}

// managedConfig reports the configuration the live connection actually uses,
// which can differ from the scheduler's latest config when a connection is reused.
func managedConfig(conn *connector.ManagedConnection) deviceConfig {
	if conn == nil {
		return deviceConfig{}
	}
	return deviceConfig{ID: conn.DeviceID, Protocol: conn.ProtocolType, Config: conn.Config}
}

func isKnownQuality(quality schema.QualityFlag) bool {
	switch quality {
	case schema.QualityGood, schema.QualityBad, schema.QualityMissing,
		schema.QualityStale, schema.QualityInvalid, schema.QualityUncertain:
		return true
	default:
		return false
	}
}

func pollingGroupPointer(groupID string) *string {
	if groupID == "" {
		return nil
	}
	return &groupID
}

// emitValue 發送收集到的值
func (s *Scheduler) emitValue(cv CollectedValue) {
	select {
	case s.valueChan <- cv:
	default:
		// 緩衝區已滿，嘗試丟棄舊值後再非阻塞寫入
		select {
		case <-s.valueChan:
		default:
		}
		select {
		case s.valueChan <- cv:
		default:
		}
	}
}
