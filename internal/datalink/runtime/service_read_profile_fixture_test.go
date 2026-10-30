package runtime

import (
	"context"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// Standalone device fixtures expose their real point service readback to the
// same projection interface used in production, including auto-assigned groups.
func attachDeviceProfileFixture(s *Service, d *schema.Device) {
	s.workspace = compensationWorkspaceFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
		points, err := s.pointSvc.ListByDevice(ctx, d.ID)
		if err != nil {
			return nil, err
		}
		return &workspace.RuntimeProjection{WorkspaceID: "fixture", Version: "fixture", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: points}, nil
	})
}
