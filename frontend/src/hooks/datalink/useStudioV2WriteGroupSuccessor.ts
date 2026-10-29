import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WriteGroupSuccessorAPI as api } from '../../services/studioV2WriteGroupSuccessor';
import type { WriteGroupSuccessorOperationStatusRequest } from '../../types/studioV2WriteGroupSuccessor';
import { studioV2WorkspaceKeys } from './keys';
import { isRecord } from '../../utils/studioV2WriteGroupJsonPrimitives';
import { studioV2WriteGroupLifecycleKeys } from './studioV2WriteGroupLifecycleKeys';
export type WriteGroupSuccessorVariables<T> = { groupId: string; request: T };

function useSuccessorStage<T, R>(run: (groupId: string, request: T) => Promise<R>, invalidate = false, schema = false) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ groupId, request }: WriteGroupSuccessorVariables<T>) => run(groupId, request),
    retry: false,
    onSettled: async (_data, _error, variables) => {
      if (!invalidate) return;
      const operationId = schema && isRecord(variables.request) && typeof variables.request.schema_operation_id === 'string'
        ? variables.request.schema_operation_id : undefined;
      await Promise.all([
        client.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups() }),
        client.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() }),
        ...(schema ? [client.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseMetadataAll() })] : []),
        ...(operationId ? [
          client.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroupSchemaOperation(operationId) }),
          client.invalidateQueries({ queryKey: studioV2WorkspaceKeys.schemaOperation(operationId) }),
        ] : []),
      ]);
    },
  });
}
/** Read-only reviews never activate; every mutation explicitly disables inherited retries. */
export function useWriteGroupSuccessorActions() {
  const review = useSuccessorStage(api.review);
  const reserve = useSuccessorStage(api.reserve, true);
  const schemaPreview = useSuccessorStage(api.schemaPreview);
  const schemaConfirm = useSuccessorStage(api.schemaConfirm, true, true);
  const activationReview = useSuccessorStage(api.activationReview);
  const activationReadiness = useSuccessorStage(api.activationReadiness);
  const confirmCutover = useSuccessorStage(api.confirmCutover, true);
  const confirmFirstStart = useSuccessorStage(api.confirmFirstStart, true);
  return { review, reserve, schemaPreview, schemaConfirm, activationReview, activationReadiness, confirmCutover, confirmFirstStart };
}
export function useCheckSuccessorOperation() {
  const client = useQueryClient();
  return (groupId: string, request: WriteGroupSuccessorOperationStatusRequest) => client.fetchQuery({
    queryKey: studioV2WriteGroupLifecycleKeys.operation(request.workspace_id, groupId, request.lineage_id, request.operation_id),
    queryFn: () => api.operationStatus(groupId, request), staleTime: 0, retry: false,
  });
}
export function useSuccessorOperationStatusQuery(groupId: string | undefined, request: WriteGroupSuccessorOperationStatusRequest | undefined, enabled = true) {
  return useQuery({
    queryKey: studioV2WriteGroupLifecycleKeys.operation(request?.workspace_id ?? '', groupId ?? '', request?.lineage_id ?? '', request?.operation_id ?? ''),
    queryFn: () => {
      if (!groupId || !request) throw new Error('Recording operation scope is required.');
      return api.operationStatus(groupId, request);
    },
    enabled: enabled && Boolean(groupId && request), retry: false, refetchOnWindowFocus: false,
  });
}
