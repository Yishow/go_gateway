package sourcerule

import (
	"context"

	"go-gateway/internal/datalink/schema"
)

// sourceLinksChangeAcquisition reuses the canonical source-only identity. Tags,
// mappings and link installation IDs cannot become gap authorization inputs.
func sourceLinksChangeAcquisition(before, after []*schema.SourceRuleLink) bool {
	old, oldErr := AcquisitionSourceDigest(AcquisitionSourceRows{Links: before}, false)
	next, nextErr := AcquisitionSourceDigest(AcquisitionSourceRows{Links: after}, false)
	return oldErr != nil || nextErr != nil || old != next
}

func (s *Service) beginSourceLinkAcquisitionMutation(ctx context.Context, ruleID string, before, after []*schema.SourceRuleLink) (context.Context, func(), error) {
	if !sourceLinksChangeAcquisition(before, after) {
		return ctx, func() {}, nil
	}
	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return ctx, func() {}, err
	}
	next, finish := s.beginSourceAcquisitionMutation(ctx, ruleID, rule.DeviceID)
	return next, finish, nil
}
