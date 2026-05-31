package handlers

import (
	"context"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

type defaultPointDirectReader struct {
	pointSvc  *point.Service
	deviceSvc *device.Service
}

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

		readResult, readErr := r.deviceSvc.ReadPointNow(ctx, pt)
		collected := collector.CollectedValue{
			PointID:   pt.ID,
			DeviceID:  pt.DeviceID,
			Value:     readResult.Value,
			RawBytes:  readResult.RawBytes,
			Timestamp: readResult.Timestamp,
			Quality:   readResult.Quality,
			Error:     readResult.Error,
		}
		if collected.Timestamp.IsZero() {
			collected.Timestamp = time.Now()
		}
		if readErr != nil {
			collected.Quality = schema.QualityBad
			collected.Error = readErr.Error()
		}

		results = append(results, collected)
	}
	return results
}
