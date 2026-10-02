import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceWriteGroupsAPI } from '../../services/studioV2WorkspaceWriteGroups';
import type {
  WriteGroupCreateRequest,
  WriteGroupDeleteRequest,
  WriteGroupUpdateRequest,
} from '../../types/studioV2WriteGroup';
import { studioV2WorkspaceKeys } from './keys';

export function useWriteGroupsQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.writeGroups(),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.list(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

/** Read-only local and destination-schema readiness of one saved group. */
export function useWriteGroupReadinessQuery(groupId: string | undefined, enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.writeGroupReadiness(groupId ?? ''),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.readiness(groupId as string),
    enabled: enabled && Boolean(groupId),
    retry: false,
    refetchOnWindowFocus: false,
  });
}

/** Where the group's accepted rows stand; only `sql_committed` is destination evidence. */
export function useWriteGroupDeliveryQuery(groupId: string | undefined, enabled: boolean, refetchIntervalMs?: number) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.writeGroupDelivery(groupId ?? ''),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.delivery(groupId as string),
    enabled: enabled && Boolean(groupId),
    retry: false,
    refetchOnWindowFocus: false,
    refetchInterval: refetchIntervalMs,
  });
}

function useInvalidateWriteGroups() {
  const queryClient = useQueryClient();
  return async () => {
    await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups() });
    await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
  };
}

export function useCreateWriteGroupMutation() {
  const invalidate = useInvalidateWriteGroups();
  return useMutation({
    mutationFn: (request: WriteGroupCreateRequest) => studioV2WorkspaceWriteGroupsAPI.create(request),
    retry: false,
    onSettled: invalidate,
  });
}

export function useUpdateWriteGroupMutation() {
  const invalidate = useInvalidateWriteGroups();
  return useMutation({
    mutationFn: (variables: { id: string; request: WriteGroupUpdateRequest }) =>
      studioV2WorkspaceWriteGroupsAPI.update(variables.id, variables.request),
    retry: false,
    onSettled: invalidate,
  });
}

/** Lifecycle actions all carry expected revisions only and never retry on their own. */
export function useWriteGroupLifecycleMutation(action: 'apply' | 'disable' | 'remove') {
  const invalidate = useInvalidateWriteGroups();
  return useMutation({
    mutationFn: (variables: { id: string; request: WriteGroupDeleteRequest }) =>
      studioV2WorkspaceWriteGroupsAPI[action](variables.id, variables.request),
    retry: false,
    onSettled: invalidate,
  });
}
