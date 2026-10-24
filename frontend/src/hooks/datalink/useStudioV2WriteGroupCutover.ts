import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceWriteGroupsAPI } from '../../services/studioV2WorkspaceWriteGroups';
import type { WriteGroupCutoverRequest } from '../../types/studioV2WriteGroupCutover';
import { studioV2WorkspaceKeys } from './keys';

type CutoverVariables = { sourceId: string; request: WriteGroupCutoverRequest };

/** Under the write-groups prefix, so every group invalidation also refreshes the recorded cutovers. */
const cutoversKey = (sourceId: string) => [...studioV2WorkspaceKeys.writeGroups(), sourceId, 'cutovers'] as const;

/** Read-only difference of a forward cutover (or its rollback) at an exact instant; it changes nothing. */
export function useCutoverReviewMutation() {
  return useMutation({
    mutationFn: (variables: CutoverVariables) => studioV2WorkspaceWriteGroupsAPI.reviewCutover(variables.sourceId, variables.request),
    retry: false,
  });
}

/** The one transactional switch. A lost answer is never retried here: the caller checks the recorded cutovers first. */
export function useApplyCutoverMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (variables: CutoverVariables) => studioV2WorkspaceWriteGroupsAPI.applyCutover(variables.sourceId, variables.request),
    retry: false,
    onSettled: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

/** Reads the recorded cutovers of one source on demand, bypassing the cache, to settle a lost answer. */
export function useCheckCutovers() {
  const queryClient = useQueryClient();
  return (sourceId: string) => queryClient.fetchQuery({
    queryKey: cutoversKey(sourceId), queryFn: () => studioV2WorkspaceWriteGroupsAPI.listCutovers(sourceId), staleTime: 0,
  });
}

/** The cutovers and rollbacks recorded for one source group. */
export function useCutoversQuery(sourceId: string | undefined, enabled: boolean) {
  return useQuery({
    queryKey: cutoversKey(sourceId ?? ''),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.listCutovers(sourceId as string),
    enabled: enabled && Boolean(sourceId),
    retry: false,
    refetchOnWindowFocus: false,
  });
}
