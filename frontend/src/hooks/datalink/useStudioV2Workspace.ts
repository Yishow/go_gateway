import { recordingReadKey, validRecordingReadScope, type StudioV2RecordingReadScope } from './studioV2RecordingReadScope';
import { useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { studioV2WorkspaceAPI } from '../../services/studioV2Workspace';

export function useStudioV2WorkspaceQuery(enabled = true, scope?: StudioV2RecordingReadScope) {
  return useQuery({
    queryKey: recordingReadKey(studioV2WorkspaceKeys.bootstrap(), scope),
    queryFn: () => studioV2WorkspaceAPI.get(),
    enabled: enabled && (!scope || validRecordingReadScope(scope)),
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5000,
  });
}
