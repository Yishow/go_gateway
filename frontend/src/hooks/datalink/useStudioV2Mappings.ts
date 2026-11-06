import { observeRecordingHomeSourceMutation } from './useRecordingHomeSources';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import {
  studioV2MappingsAPI,
  type StudioV2WorkspaceMappingRequest,
} from '../../services/studioV2Mappings';

export function useStudioV2MappingsQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.mappings(),
    queryFn: () => studioV2MappingsAPI.list(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5000,
  });
}

export function useCreateStudioV2MappingMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: (request: StudioV2WorkspaceMappingRequest) => observeRecordingHomeSourceMutation(queryClient, () => studioV2MappingsAPI.create(request), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.mappings(), refetchType: 'none' });
      // Refresh readiness in the background; a slow GET cannot hold this save.
      void queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useUpdateStudioV2MappingMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: ({ mappingId, request }: { mappingId: string; request: StudioV2WorkspaceMappingRequest }) =>
      observeRecordingHomeSourceMutation(queryClient, () => studioV2MappingsAPI.update(mappingId, request), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.mappings(), refetchType: 'none' });
      // Refresh readiness in the background; a slow GET cannot hold this save.
      void queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useDeleteStudioV2MappingMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: (mappingId: string) => observeRecordingHomeSourceMutation(queryClient, () => studioV2MappingsAPI.remove(mappingId), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.mappings(), refetchType: 'none' });
      // Refresh readiness in the background; a slow GET cannot hold this save.
      void queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}
