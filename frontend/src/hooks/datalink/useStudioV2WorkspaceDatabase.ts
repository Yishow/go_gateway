import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import {
  studioV2WorkspaceDatabaseAPI,
  type StudioV2WorkspaceDatabaseConfigRequest,
  type StudioV2WorkspaceDatabaseTargetRequest,
} from '../../services/studioV2WorkspaceDatabase';

export function useStudioV2DatabaseConfigQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.databaseConfig(),
    queryFn: () => studioV2WorkspaceDatabaseAPI.getConfig(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useStudioV2DatabaseTargetsQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.databaseTargets(),
    queryFn: () => studioV2WorkspaceDatabaseAPI.listTargets(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useUpdateStudioV2DatabaseConfigMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: StudioV2WorkspaceDatabaseConfigRequest) => studioV2WorkspaceDatabaseAPI.updateConfig(request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseConfig() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseTargets() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}

export function useUpsertStudioV2DatabaseTargetMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ pointId, request }: { pointId: string; request: StudioV2WorkspaceDatabaseTargetRequest }) =>
      studioV2WorkspaceDatabaseAPI.upsertTarget(pointId, request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseTargets() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
    },
  });
}
