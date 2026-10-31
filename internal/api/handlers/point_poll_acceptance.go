package handlers

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/point"
)

type pointInstallationSampleVerifier interface {
	AcceptsInstallationSample(string, uint64, string) bool
}
type pointDirectSampleVerifier interface {
	AcceptsDirectRead(context.Context, collector.CollectedValue) bool
}

func (h *PointHandler) acceptsPollResult(ctx context.Context, cv collector.CollectedValue) bool {
	if cv.DirectEvidence != nil {
		verifier, ok := h.directReader.(pointDirectSampleVerifier)
		return ok && verifier.AcceptsDirectRead(ctx, cv)
	}
	if cv.InstallationLease != 0 {
		verifier, ok := h.manualPoller.(pointInstallationSampleVerifier)
		if !ok || !verifier.AcceptsInstallationSample(cv.PointID, cv.InstallationLease, cv.InstallationScope) {
			return false
		}
		current, err := h.svc.GetByID(ctx, cv.PointID)
		return err == nil && current != nil && current.RevisionID == cv.PointRevision && current.DeviceID == cv.DeviceID && verifier.AcceptsInstallationSample(cv.PointID, cv.InstallationLease, cv.InstallationScope)
	}
	if verifier, ok := h.manualPoller.(interface{ AcceptsLegacyUnscopedSample(string) bool }); ok {
		return verifier.AcceptsLegacyUnscopedSample(cv.PointID)
	}
	return true // Historical non-installed poller API compatibility only.
}
func supersededPollResult(id string) PollResult {
	return PollResult{PointID: id, Quality: pollQualityBad, Error: "Read result superseded by configuration change"}
}

// persistPollResult uses only identity captured before provider/wire I/O.
// It never holds scheduler locks while writing SQL.
func (h *PointHandler) persistPollResult(ctx context.Context, cv collector.CollectedValue) (bool, error) {
	var err error
	if cv.DirectEvidence != nil {
		captured := cv.DirectEvidence.Point
		err = h.svc.UpdateReadResultAtRevision(ctx, cv.PointID, captured.DeviceID, captured.RevisionID, cv.Value, cv.Error)
	} else if cv.InstallationLease != 0 {
		err = h.svc.UpdateReadResultAtRevision(ctx, cv.PointID, cv.DeviceID, cv.PointRevision, cv.Value, cv.Error)
	} else {
		err = h.svc.UpdateReadResult(ctx, cv.PointID, cv.Value, cv.Error)
	}
	if errors.Is(err, point.ErrReadResultSuperseded) {
		return false, nil
	}
	return err == nil, err
}
