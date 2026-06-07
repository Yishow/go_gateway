import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2RulesAPI, type StudioV2UpdateSourceRuleRequest } from '../../services/studioV2Rules';
import type { CreateSourceRuleRequest } from '../../types/datalink';

export function useStudioV2RulesQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.sourceRules(),
    queryFn: () => studioV2RulesAPI.list(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useCreateStudioV2RuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: CreateSourceRuleRequest) => studioV2RulesAPI.create(request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useUpdateStudioV2RuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ ruleId, request }: { ruleId: string; request: StudioV2UpdateSourceRuleRequest }) =>
      studioV2RulesAPI.update(ruleId, request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useDeleteStudioV2RuleMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (ruleId: string) => studioV2RulesAPI.remove(ruleId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}
