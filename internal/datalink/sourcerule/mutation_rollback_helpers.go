package sourcerule

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

// finishRollbackErrors decides whether the original mutation error is still
// safe to return after every compensation call has run. A dirty projection is
// reported when any call cannot prove that its state was restored.
func (s *Service) finishRollbackErrors(ctx context.Context, originalErr error, rollbackErrs ...error) error {
	if s.noteTransactionFailure(ctx, originalErr) {
		return s.dirtyUnknown(ctx)
	}
	for _, rollbackErr := range rollbackErrs {
		if rollbackErr != nil {
			s.noteTransactionFailure(ctx, rollbackErr)
			return s.dirtyUnknown(ctx)
		}
	}
	// A newer point revision is not owned by this materializer or its outer snapshot.
	if errors.Is(originalErr, point.ErrFormatConflict) || errors.Is(originalErr, ErrRuleRevisionConflict) {
		return s.dirtyUnknown(ctx)
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
	var failures []error
	for _, restore := range []func() error{
		func() error { return s.rollbackTagMappingSync(ctx, tagMappings) },
		func() error { return s.rollbackUpdatedPoints(ctx, updatedPoints) },
		func() error { return s.rollbackRuleState(ctx, rule, links) },
		func() error { return s.rollbackCreatedPoints(ctx, createdPointIDs) },
	} {
		if err := restore(); err != nil {
			failures = append(failures, err)
			if s.noteTransactionFailure(ctx, err) {
				break
			}
		}
	}
	return errors.Join(failures...)

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
	if s.noteTransactionFailure(ctx, originalErr) {
		return s.dirtyUnknown(ctx)
	}
	return s.finishRollbackErrors(ctx, originalErr, s.rollbackUpdateState(ctx, tagMappings, updatedPoints, rule, links, createdPointIDs))
}
