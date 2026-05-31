import { useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2RuntimeContextAPI } from '../../services/studioV2RuntimeContext';

export function useStudioV2RuntimeContextQuery() {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.runtimeContext(),
    queryFn: () => studioV2RuntimeContextAPI.get(),
    retry: false,
    refetchOnWindowFocus: false,
  });
}
