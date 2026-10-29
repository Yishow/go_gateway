import { useQuery } from '@tanstack/react-query';
import { studioV2WriteGroupLifecycleAPI } from '../../services/studioV2WriteGroupLifecycle';
import { requireLifecycleMatch } from '../../services/studioV2LifecycleTransport';
import { studioV2WriteGroupLifecycleKeys } from './studioV2WriteGroupLifecycleKeys';
export function useWriteGroupLifecycleQuery(workspaceId: string | undefined, groupId: string | undefined, enabled = true) {
  return useQuery({
    queryKey: studioV2WriteGroupLifecycleKeys.lifecycle(workspaceId ?? '', groupId ?? ''),
    queryFn: async () => {
      if (!workspaceId || !groupId) throw new Error('Recording scope is required.');
      const result = await studioV2WriteGroupLifecycleAPI.lifecycle(groupId);
      return requireLifecycleMatch(result, result.workspace_id === workspaceId && result.group_id === groupId);
    },
    enabled: enabled && Boolean(workspaceId && groupId), retry: false, refetchOnWindowFocus: false,
  });
}
