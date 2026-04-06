import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { sourceRuleTagReviewDecisionAPI } from '../../services/sourceRuleTagReviewDecisions';
import { sourceRuleKeys } from './keys';
import type {
  SourceRuleTagReviewDecision,
  UpsertSourceRuleTagReviewDecisionRequest,
} from '../../types/sourceRuleTagReviewDecisions';

function upsertReviewDecision(
  current: SourceRuleTagReviewDecision[] | undefined,
  next: SourceRuleTagReviewDecision,
) {
  const decisions = current ?? [];
  return [
    ...decisions.filter((decision) => decision.candidate_id !== next.candidate_id),
    next,
  ];
}

export function useSourceRuleTagReviewDecisionsQuery(ruleId?: string | null) {
  return useQuery({
    queryKey: sourceRuleKeys.reviewDecisions(ruleId ?? ''),
    queryFn: () => sourceRuleTagReviewDecisionAPI.list(ruleId!),
    enabled: Boolean(ruleId),
  });
}

export function useUpsertSourceRuleTagReviewDecisionMutation(ruleId?: string | null) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: UpsertSourceRuleTagReviewDecisionRequest) => {
      if (!ruleId) {
        throw new Error('source rule id is required');
      }
      return sourceRuleTagReviewDecisionAPI.upsert(ruleId, request);
    },
    onSuccess: (decision) => {
      if (!ruleId) {
        return;
      }
      queryClient.setQueryData<SourceRuleTagReviewDecision[]>(
        sourceRuleKeys.reviewDecisions(ruleId),
        (current) => upsertReviewDecision(current, decision),
      );
    },
  });
}
