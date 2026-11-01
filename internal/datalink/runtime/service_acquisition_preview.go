package runtime

import (
	"context"
	"go-gateway/internal/datalink/collector"
)

// ReadAcquisitionPreview observes the installed scheduler without starting it.
// The returned verifier retains the original native proof, never a new snapshot.
func (s *Service) ReadAcquisitionPreview(ctx context.Context, workspaceID, deviceID, sourceRevision string) (collector.AcquisitionPreview, func() bool, error) {
	if s == nil || s.scheduler == nil {
		return collector.AcquisitionPreview{DeviceID: deviceID, Status: "unavailable", Reason: "installation_unavailable"}, nil, nil
	}
	view, proof, err := s.scheduler.PreviewAcquisition(ctx, workspaceID, deviceID, sourceRevision)
	if err != nil || view.Status != "available" {
		return view, nil, err
	}
	return view, func() bool { return s.scheduler.AcquisitionPreviewCurrent(proof) }, nil
}
