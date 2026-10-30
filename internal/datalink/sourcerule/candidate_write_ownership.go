package sourcerule

import "context"

func (s *Service) noteCandidateWrite(ctx context.Context, revision string) {
	if j := s.pointJournal(ctx); j != nil {
		if j.candidateWrites == nil {
			j.candidateWrites = map[string]bool{}
		}
		j.candidateWrites[revision] = true
	}
}
func (s *Service) deleteOwnedCandidates(ctx context.Context, ruleID, keepRevision string) error {
	j := s.pointJournal(ctx)
	if j == nil {
		return nil
	}
	for revision := range j.candidateWrites {
		if revision != keepRevision {
			if err := s.repo.DeleteCandidateSnapshots(ctx, ruleID, revision); err != nil {
				return err
			}
		}
	}
	return nil
}
