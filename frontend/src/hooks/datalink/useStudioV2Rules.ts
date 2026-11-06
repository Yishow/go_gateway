import { observeRecordingHomeSourceMutation } from './useRecordingHomeSources';
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
    staleTime: 5000,
  });
}

export function useCreateStudioV2RuleMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: (request: CreateSourceRuleRequest) => observeRecordingHomeSourceMutation(queryClient, () => studioV2RulesAPI.create(request), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useUpdateStudioV2RuleMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: ({ ruleId, request }: { ruleId: string; request: StudioV2UpdateSourceRuleRequest }) =>
      observeRecordingHomeSourceMutation(queryClient, () => studioV2RulesAPI.update(ruleId, request), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useDeleteStudioV2RuleMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: (ruleId: string) => observeRecordingHomeSourceMutation(queryClient, () => studioV2RulesAPI.remove(ruleId), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.sourceRules() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}
