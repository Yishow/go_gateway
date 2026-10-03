import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2RecordingStartAPI } from '../../services/studioV2RecordingStart';
import type { RecordingStartRequest } from '../../types/studioV2RecordingStart';
import type { WriteGroupCreateRequest } from '../../types/studioV2WriteGroup';
import { studioV2WorkspaceKeys } from './keys';

function useInvalidateStart() {
  const queryClient = useQueryClient();
  return () => Promise.all([
    queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups() }),
    queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() }),
  ]);
}

export function useEnsureBasicManagedMutation() {
  const invalidate = useInvalidateStart();
  return useMutation({
    mutationFn: ({ deviceId, request }: { deviceId: string; request: WriteGroupCreateRequest }) =>
      studioV2RecordingStartAPI.ensureBasic(deviceId, request),
    retry: false, onSettled: invalidate,
  });
}

export function useRecordingStartMutation() {
  const invalidate = useInvalidateStart();
  return useMutation({
    mutationFn: (request: RecordingStartRequest) => studioV2RecordingStartAPI.start(request),
    retry: false, onSettled: invalidate,
  });
}

export function useRecordingStartOperationQuery(operationId: string | undefined) {
  return useQuery({
    queryKey: [...studioV2WorkspaceKeys.all, 'recording-start', operationId ?? ''],
    queryFn: () => studioV2RecordingStartAPI.operation(operationId as string),
    enabled: Boolean(operationId), retry: false, refetchOnWindowFocus: false,
    refetchInterval: (query) =>
      query.state.data?.status === 'running' || query.state.data?.status === 'pending' ? 2000 : false,
  });
}
