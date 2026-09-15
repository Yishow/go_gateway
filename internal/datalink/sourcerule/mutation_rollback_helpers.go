package sourcerule

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/schema"
)

// finishRollbackErrors decides whether the original mutation error is still
// safe to return after every compensation call has run. A dirty projection is
// reported when any call cannot prove that its state was restored.
func (s *Service) finishRollbackErrors(ctx context.Context, originalErr error, rollbackErrs ...error) error {
	for _, rollbackErr := range rollbackErrs {
		if rollbackErr != nil {
			return s.dirtyUnknown(ctx)
		}
	}
	return originalErr
}

func (s *Service) rollbackUpdateState(
	ctx context.Context,
	tagMappings tagMappingSyncResult,
	updatedPoints []pointUpdatePlan,
	rule *schema.SourceRule,
	links []*schema.SourceRuleLink,
	createdPointIDs []string,
) error {
	return errors.Join(
		s.rollbackTagMappingSync(ctx, tagMappings),
		s.rollbackUpdatedPoints(ctx, updatedPoints),
		s.rollbackRuleState(ctx, rule, links),
		s.rollbackCreatedPoints(ctx, createdPointIDs),
	)
}

func (s *Service) finishUpdateRollback(
	ctx context.Context,
	originalErr error,
	tagMappings tagMappingSyncResult,
	updatedPoints []pointUpdatePlan,
	rule *schema.SourceRule,
	links []*schema.SourceRuleLink,
	createdPointIDs []string,
) error {
	return s.finishRollbackErrors(ctx, originalErr, s.rollbackUpdateState(ctx, tagMappings, updatedPoints, rule, links, createdPointIDs))
}
