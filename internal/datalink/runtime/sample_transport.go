package runtime

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// SampleSink is the optional typed acquisition transport seam. Production
// wiring may leave it nil while the legacy writer remains the active path.
type SampleSink interface {
	AcceptSample(context.Context, measurement.SampleEnvelope) error
}

// SampleInterest is an optional extension of SampleSink: a sink that knows
// which sources it consumes lets the runtime skip building typed envelopes for
// everything else, so unrelated tags cannot produce rejections or failures.
type SampleInterest interface {
	WantsSample(deviceID, pointID, tagID string) bool
}

// ErrTypedSampleRejected identifies a sample that cannot be tied to the
// installed runtime identity without fabricating missing metadata.
var ErrTypedSampleRejected = errors.New("typed sample rejected")

func mappingBindingForRecords(
	workspaceID string,
	mappingRecord *schema.Mapping,
	deviceRecord *schema.Device,
	pointRecord *schema.Point,
	tagRecord *schema.Tag,
) (mappingBinding, error) {
	if mappingRecord == nil || deviceRecord == nil || pointRecord == nil || tagRecord == nil {
		return mappingBinding{}, fmt.Errorf("%w: mapping identity is incomplete", ErrTypedSampleRejected)
	}
	if strings.TrimSpace(deviceRecord.ID) == "" || strings.TrimSpace(pointRecord.ID) == "" || strings.TrimSpace(tagRecord.ID) == "" {
		return mappingBinding{}, fmt.Errorf("%w: mapping identity is empty", ErrTypedSampleRejected)
	}
	if pointRecord.DeviceID != deviceRecord.ID || mappingRecord.PointID != pointRecord.ID || mappingRecord.TagID != tagRecord.ID {
		return mappingBinding{}, fmt.Errorf("%w: mapping identity does not match installed records", ErrTypedSampleRejected)
	}
	sourceRevision, err := measurement.SourceRevision(*deviceRecord, *pointRecord, *tagRecord)
	if err != nil {
		return mappingBinding{}, fmt.Errorf("derive source revision: %w", err)
	}
	mappingRevision, err := measurement.MappingRevision(*mappingRecord)
	if err != nil {
		return mappingBinding{}, fmt.Errorf("derive mapping revision: %w", err)
	}
	configFingerprint, err := measurement.AcquisitionConfigFingerprint(*deviceRecord, *pointRecord)
	if err != nil {
		return mappingBinding{}, fmt.Errorf("derive acquisition config fingerprint: %w", err)
	}
	return mappingBinding{
		WorkspaceID:       workspaceID,
		DeviceID:          deviceRecord.ID,
		PointID:           pointRecord.ID,
		TagID:             tagRecord.ID,
		TagDataType:       tagRecord.DataType,
		TransformPipeline: mappingRecord.TransformPipeline,
		SourceRevision:    sourceRevision,
		MappingRevision:   mappingRevision,
		ConfigFingerprint: configFingerprint,
	}, nil
}

func typedSampleFromValue(cv collector.CollectedValue, binding mappingBinding, value, rawValue any) (measurement.SampleEnvelope, error) {
	if strings.TrimSpace(binding.WorkspaceID) == "" || strings.TrimSpace(binding.DeviceID) == "" || strings.TrimSpace(binding.PointID) == "" || strings.TrimSpace(binding.TagID) == "" {
		return measurement.SampleEnvelope{}, fmt.Errorf("%w: runtime binding identity is incomplete", ErrTypedSampleRejected)
	}
	if strings.TrimSpace(binding.SourceRevision) == "" || strings.TrimSpace(binding.MappingRevision) == "" || strings.TrimSpace(binding.ConfigFingerprint) == "" {
		return measurement.SampleEnvelope{}, fmt.Errorf("%w: runtime binding revision is incomplete", ErrTypedSampleRejected)
	}
	if strings.TrimSpace(cv.AcquisitionID) == "" || cv.ObservedAt.IsZero() || cv.ReceivedAt.IsZero() {
		return measurement.SampleEnvelope{}, fmt.Errorf("%w: acquisition timing or identity is incomplete", ErrTypedSampleRejected)
	}
	if cv.DeviceID != binding.DeviceID || cv.PointID != binding.PointID {
		return measurement.SampleEnvelope{}, fmt.Errorf("%w: collected identity does not match runtime binding", ErrTypedSampleRejected)
	}
	if strings.TrimSpace(cv.ConfigFingerprint) == "" || cv.ConfigFingerprint != binding.ConfigFingerprint {
		return measurement.SampleEnvelope{}, fmt.Errorf("%w: acquisition configuration fingerprint mismatch", ErrTypedSampleRejected)
	}
	timeOrigin := measurement.GatewayTimeOrigin(cv.TimeOrigin)
	if timeOrigin != measurement.GatewayTimeOriginSource && timeOrigin != measurement.GatewayTimeOriginGateway {
		return measurement.SampleEnvelope{}, fmt.Errorf("%w: unknown time origin", ErrTypedSampleRejected)
	}
	quality := cv.Quality
	qualityReason := strings.TrimSpace(cv.QualityReason)
	unknownQuality := !isKnownQuality(quality)
	if unknownQuality {
		quality = schema.QualityBad
		qualityReason = "quality-unknown"
	}
	if cv.Error != "" && !unknownQuality {
		quality = schema.QualityBad
		if !isSafeQualityReason(qualityReason, quality) {
			qualityReason = "read-failed"
		}
	}
	if quality == schema.QualityGood {
		qualityReason = ""
	} else if !isSafeQualityReason(qualityReason, quality) {
		qualityReason = fmt.Sprintf("quality-%s", quality)
	}
	sampleID, err := measurement.StableSampleID(
		cv.AcquisitionID, binding.WorkspaceID, binding.DeviceID, binding.PointID,
		binding.TagID, binding.SourceRevision, binding.MappingRevision,
	)
	if err != nil {
		return measurement.SampleEnvelope{}, fmt.Errorf("derive stable sample identity: %w", err)
	}
	return measurement.SampleEnvelope{
		SampleID:          sampleID,
		WorkspaceID:       binding.WorkspaceID,
		DeviceID:          binding.DeviceID,
		PointID:           binding.PointID,
		TagID:             binding.TagID,
		SourceRevision:    binding.SourceRevision,
		MappingRevision:   binding.MappingRevision,
		ConfigFingerprint: binding.ConfigFingerprint,
		AcquisitionID:     cv.AcquisitionID,
		ObservedAt:        cv.ObservedAt.UTC(),
		ReceivedAt:        cv.ReceivedAt.UTC(),
		TimeOrigin:        string(timeOrigin),
		ValueType:         string(binding.TagDataType),
		Value:             value,
		RawValue:          rawValue,
		Quality:           quality,
		QualityReason:     qualityReason,
	}, nil
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

const reasonMappingFailed = "mapping-failed"

func isSafeQualityReason(reason string, quality schema.QualityFlag) bool {
	switch reason {
	case "read-failed", "quality-unknown", "connection-failed", "circuit-open", reasonMappingFailed:
		return true
	}
	return reason == fmt.Sprintf("quality-%s", quality)
}

func (s *Service) refreshMappings(ctx context.Context) error {
	if s.mappingSvc != nil {
		mappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{Enabled: common.Ptr(true), Limit: 100000})
		if err != nil {
			return fmt.Errorf("載入 mappings 失敗: %w", err)
		}

		next := make(map[string][]mappingBinding)
		for _, m := range mappings {
			tg, err := s.tagSvc.GetByID(ctx, m.TagID)
			if err != nil {
				return fmt.Errorf("載入 tag(%s) 失敗: %w", m.TagID, err)
			}
			pt, err := s.pointSvc.GetByID(ctx, m.PointID)
			if err != nil {
				return fmt.Errorf("載入 point(%s) 失敗: %w", m.PointID, err)
			}
			dev, err := s.deviceSvc.GetByID(ctx, pt.DeviceID)
			if err != nil {
				return fmt.Errorf("載入 device(%s) 失敗: %w", pt.DeviceID, err)
			}
			binding, err := mappingBindingForRecords(s.config.WorkspaceID, m, dev, pt, tg)
			if err != nil {
				return fmt.Errorf("建立 mapping(%s) runtime identity 失敗: %w", m.ID, err)
			}
			next[m.PointID] = append(next[m.PointID], binding)
		}

		s.mappingMu.Lock()
		s.mappingIndex = next
		s.mappingMu.Unlock()
		return nil
	}

	deviceByID := make(map[string]*schema.Device, len(s.snapshot.Devices))
	for _, d := range s.snapshot.Devices {
		if d != nil {
			deviceByID[d.ID] = d
		}
	}
	pointByID := make(map[string]*schema.Point, len(s.snapshot.Points))
	for _, p := range s.snapshot.Points {
		if p != nil {
			pointByID[p.ID] = p
		}
	}
	tagByID := make(map[string]*schema.Tag, len(s.snapshot.Tags))
	tagType := make(map[string]schema.DataType, len(s.snapshot.Tags))
	for _, t := range s.snapshot.Tags {
		if t == nil {
			continue
		}
		tagByID[t.ID] = t
		tagType[t.ID] = t.DataType
	}

	next := make(map[string][]mappingBinding)
	workspaceID := s.snapshot.workspaceID(s.config.WorkspaceID)
	for _, m := range s.snapshot.Mappings {
		if m == nil || !m.Enabled {
			continue
		}
		pointRecord := pointByID[m.PointID]
		var deviceRecord *schema.Device
		if pointRecord != nil {
			deviceRecord = deviceByID[pointRecord.DeviceID]
		}
		if binding, err := mappingBindingForRecords(workspaceID, m, deviceRecord, pointRecord, tagByID[m.TagID]); err == nil {
			next[m.PointID] = append(next[m.PointID], binding)
			continue
		}
		// Keep the legacy snapshot path available for old fixtures with no
		// complete typed identity; the optional sink will reject it safely.
		next[m.PointID] = append(next[m.PointID], mappingBinding{
			WorkspaceID:       workspaceID,
			PointID:           m.PointID,
			TagID:             m.TagID,
			TagDataType:       tagType[m.TagID],
			TransformPipeline: m.TransformPipeline,
		})
	}

	s.mappingMu.Lock()
	s.mappingIndex = next
	s.mappingMu.Unlock()
	return nil
}
