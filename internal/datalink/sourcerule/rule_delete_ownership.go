package sourcerule

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/schema"
)

func (s *Service) deleteOwnedRule(ctx context.Context, expected *schema.SourceRule) error {
	s.noteRuleBaseline(ctx, expected)
	if !s.ownsRuleState(ctx, expected) {
		return ErrRuleRevisionConflict
	}
	journal := s.pointJournal(ctx)
	var err error
	if repo, ok := s.repo.(ruleRevisionDeleter); ok {
		err = repo.DeleteRevision(ctx, expected)
	} else {
		err = s.repo.Delete(ctx, expected.ID)
	}
	if err != nil {
		if !errors.Is(err, ErrRuleRevisionConflict) && journal != nil {
			journal.ruleUncertain = true
		}
		return err
	}
	if journal != nil {
		journal.ruleDeleted = true
	}
	return nil
}
