import { configReadContext, configReadKey, fencedConfigRead } from './recordingConfigRefresh';
import { recordingReadKey, validRecordingReadScope, type StudioV2RecordingReadScope } from './studioV2RecordingReadScope';
import { useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2WorkspaceAPI } from '../../services/studioV2Workspace';

export function useStudioV2WorkspaceQuery(enabled = true, scope?: StudioV2RecordingReadScope) {
  const context = configReadContext(scope?.workspaceId);
  return useQuery({
    queryKey: recordingReadKey(configReadKey(studioV2WorkspaceKeys.bootstrap(), context), scope),
    queryFn: ({ signal, client }) => fencedConfigRead(client, context, signal, async () => {
      const result = await studioV2WorkspaceAPI.get({ signal, baseURL: context.apiOrigin });
      if (!result?.id || context.workspaceId && result.id !== context.workspaceId) throw new Error('Workspace observation unavailable');
      return result;
    }),
    enabled: enabled && (!scope || validRecordingReadScope(scope)),
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5000,
  });
}
