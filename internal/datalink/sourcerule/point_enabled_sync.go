package sourcerule

import (
	"context"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

// An unchanged activation projection is not a configuration edit. In particular
// restart must not manufacture a new decoding/source revision for every point.
func (s *Service) syncUnchangedRulePoint(ctx context.Context, current *schema.Point, update point.UpdatePointRequest) bool {
	if current == nil || update.Enabled == nil || current.Enabled != *update.Enabled || update.PollingGroupID != nil || update.ReplacePollingGroup {
		return false
	}
	// A no-op cannot silently accept a point changed after the plan read. Reject
	// the optimization and let the existing expected-revision CAS report conflict;
	// this read is never adopted as replacement authority.
	latest, err := s.pointSvc.GetByID(ctx, current.ID)
	if err != nil || latest.RevisionID != current.RevisionID {
		return false
	}
	// Own the runtime installation and its possible compensation using the same
	// observed configuration, without issuing a durable UPDATE or new revision.
	s.observePointWrite(ctx)(current, current)(nil)
	s.syncPoints(ctx, []*schema.Point{current})
	return true
}
