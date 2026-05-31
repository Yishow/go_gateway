import { useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2WorkspaceAPI } from '../../services/studioV2Workspace';

export function useStudioV2WorkspaceQuery() {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.bootstrap(),
    queryFn: () => studioV2WorkspaceAPI.get(),
    retry: false,
    refetchOnWindowFocus: false,
  });
}
