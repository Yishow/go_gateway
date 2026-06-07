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
  });
}

export function useCreateStudioV2MappingMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: StudioV2WorkspaceMappingRequest) => studioV2MappingsAPI.create(request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.mappings() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useUpdateStudioV2MappingMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ mappingId, request }: { mappingId: string; request: StudioV2WorkspaceMappingRequest }) =>
      studioV2MappingsAPI.update(mappingId, request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.mappings() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useDeleteStudioV2MappingMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (mappingId: string) => studioV2MappingsAPI.remove(mappingId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.mappings() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}
