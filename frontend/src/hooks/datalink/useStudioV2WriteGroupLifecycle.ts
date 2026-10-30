import { recordingReadKey, validRecordingReadScope, type StudioV2RecordingReadScope } from './studioV2RecordingReadScope';
import { useQuery } from '@tanstack/react-query';
import { studioV2WriteGroupLifecycleAPI } from '../../services/studioV2WriteGroupLifecycle';
import { requireLifecycleMatch } from '../../services/studioV2LifecycleTransport';
import { studioV2WriteGroupLifecycleKeys } from './studioV2WriteGroupLifecycleKeys';
export function writeGroupLifecycleQueryOptions(workspaceId: string | undefined, groupId: string | undefined, enabled = true, scope?: StudioV2RecordingReadScope) {
  return {
    queryKey: recordingReadKey(studioV2WriteGroupLifecycleKeys.lifecycle(workspaceId ?? '', groupId ?? ''), scope),
    queryFn: async () => {
      if (!workspaceId || !groupId) throw new Error('Recording scope is required.');
      const result = await studioV2WriteGroupLifecycleAPI.lifecycle(groupId);
      return requireLifecycleMatch(result, result.workspace_id === workspaceId && result.group_id === groupId);
    },
    enabled: enabled && Boolean(workspaceId && groupId) && (!scope || validRecordingReadScope(scope)), retry: false, refetchOnWindowFocus: false,
  };
}

export function useWriteGroupLifecycleQuery(workspaceId: string | undefined, groupId: string | undefined, enabled = true, scope?: StudioV2RecordingReadScope) {
  return useQuery(writeGroupLifecycleQueryOptions(workspaceId, groupId, enabled, scope));
}
