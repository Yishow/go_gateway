import { useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2WorkspaceAPI } from '../../services/studioV2Workspace';

export function useStudioV2WorkspaceQuery(enabled = true) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.bootstrap(),
    queryFn: () => studioV2WorkspaceAPI.get(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5000,
  });
}
