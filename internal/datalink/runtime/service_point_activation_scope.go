package runtime

import (
	"context"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// PointRuntimeInstallationRequired distinguishes a valid saved draft from an
// unprovable active configuration. A successful diagnostic is not activation.
// All authoritative projection reads occur outside metadata/scheduler locks.
func (s *Service) PointRuntimeInstallationRequired(ctx context.Context, p *schema.Point) (bool, error) {
	if s == nil || s.workspace == nil || p == nil {
		return false, connector.ErrReadProfile
	}
	projection, err := s.workspace.RuntimeProjection(ctx)
	if err != nil {
		return false, connector.ErrReadProfile
	}
	if _, err = projectionReadProfile(projection, p, true); err != nil {
		return false, err
	}
	for _, d := range projection.Devices {
		if d != nil && d.ID == p.DeviceID {
			switch d.Status {
			case schema.DeviceStatusActive:
				return true, nil
			case schema.DeviceStatusDraft, schema.DeviceStatusDisabled:
				return false, nil
			default:
				return false, connector.ErrReadProfile
			}
		}
	}
	return false, connector.ErrReadProfile
}
