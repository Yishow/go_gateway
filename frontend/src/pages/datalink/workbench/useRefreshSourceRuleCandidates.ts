import { useQueryClient } from '@tanstack/react-query';
import { useCallback } from 'react';
import { sourceRuleKeys } from '../../../hooks/datalink/keys';

export function useRefreshSourceRuleCandidates() {
  const queryClient = useQueryClient();
  return useCallback(
    async () => queryClient.invalidateQueries({ queryKey: sourceRuleKeys.all }),
    [queryClient],
  );
}
