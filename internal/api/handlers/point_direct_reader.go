package handlers

import (
	"context"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

type pointReadProfileResolver interface {
	ResolvePointReadProfile(context.Context, *schema.Point) (connector.ReadProfile, error)
}

type defaultPointDirectReader struct {
	profiles  pointReadProfileResolver
	pointSvc  *point.Service
	deviceSvc *device.Service
}

// NewPointDirectReader creates a connector-backed point reader.
//
//nolint:revive // Preserve the existing exported constructor signature.
func NewPointDirectReader(
	pointSvc *point.Service,
	deviceSvc *device.Service,
) *defaultPointDirectReader {
	return &defaultPointDirectReader{
		pointSvc:  pointSvc,
		deviceSvc: deviceSvc,
	}
}

func (r *defaultPointDirectReader) PollDirect(
	ctx context.Context,
	pointIDs []string,
) []collector.CollectedValue {
	if r == nil || r.pointSvc == nil || r.deviceSvc == nil {
		return nil
	}

	results := make([]collector.CollectedValue, 0, len(pointIDs))
	for _, pointID := range pointIDs {
		pt, err := r.pointSvc.GetByID(ctx, pointID)
		if err != nil || pt == nil {
			continue
		}

		evidence, profileErr := r.captureDirectEvidence(ctx, pt)
		var profile connector.ReadProfile
		if r.profiles != nil && profileErr == nil {
			profile, profileErr = r.profiles.ResolvePointReadProfile(ctx, pt)
		}
		evidence.Profile = profile
		var readResult connector.ReadResult
		readErr := profileErr
		if readErr == nil {
			readResult, readErr = r.deviceSvc.ReadPointNow(ctx, pt, profile)
		}
		collected := collector.CollectedValue{
			DirectEvidence: evidence,
			PointID:        pt.ID,
			DeviceID:       pt.DeviceID,
			Value:          readResult.Value,
			RawBytes:       readResult.RawBytes,
			Timestamp:      readResult.Timestamp,
			Quality:        readResult.Quality,
			Error:          readResult.Error,
		}
		if collected.Timestamp.IsZero() {
			collected.Timestamp = time.Now()
		}
		if readErr != nil {
			collected.Quality = schema.QualityBad
			collected.Error = readErr.Error()
		}

		if readErr == nil && !r.AcceptsDirectRead(ctx, collected) {
			collected.Value = nil
			collected.RawBytes = nil
			collected.Quality = schema.QualityBad
			collected.Error = "Read result superseded by configuration change"
		}
		results = append(results, collected)
	}
	return results
}

func (r *defaultPointDirectReader) WithReadProfileResolver(resolver pointReadProfileResolver) *defaultPointDirectReader {
	r.profiles = resolver
	return r
}
