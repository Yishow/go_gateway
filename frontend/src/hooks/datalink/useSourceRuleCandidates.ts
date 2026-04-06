import { useQuery } from '@tanstack/react-query';
import { sourceRuleKeys } from './keys';
import { sourceRuleCandidateAPI } from '../../services/sourceRuleCandidates';

export function useSourceRuleCandidatesQuery(ruleId?: string | null) {
  return useQuery({
    queryKey: sourceRuleKeys.candidates(ruleId ?? ''),
    queryFn: () => sourceRuleCandidateAPI.get(ruleId!),
    enabled: Boolean(ruleId),
  });
}
