package sourcerule

import (
	"context"
	"go-gateway/internal/datalink/schema"
)

type pointRuntimeActivationScope interface {
	PointRuntimeInstallationRequired(context.Context, *schema.Point) (bool, error)
}

// A known inactive device retains its durable rule/point intent for explicit
// activation. It is not an installation failure and must not manufacture a lease.
func (s *Service) pointRuntimeRequired(ctx context.Context, p *schema.Point) (bool, error) {
	if !p.Enabled {
		return false, nil
	}
	if scope, ok := s.runtimeSync.(pointRuntimeActivationScope); ok {
		return scope.PointRuntimeInstallationRequired(ctx, p)
	}
	return true, nil
}
