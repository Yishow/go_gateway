import { useMutation, useQuery } from '@tanstack/react-query';
import { studioV2WorkspaceWriteGroupsAPI } from '../../services/studioV2WorkspaceWriteGroups';
import type { WriteGroupTestWriteConfirmation } from '../../types/studioV2WriteGroupTestWrite';
import { studioV2WorkspaceKeys } from './keys';

/** Previews a test write; nothing is written to the target. */
export function useWriteGroupTestWritePreviewMutation() {
  return useMutation({
    mutationFn: (groupId: string) => studioV2WorkspaceWriteGroupsAPI.testWritePreview(groupId),
    retry: false,
  });
}

/**
 * Confirms a previewed test write. It never retries: a repeat after an unknown
 * result must be checked through the operation, not sent again blindly.
 */
export function useWriteGroupTestWriteMutation() {
  return useMutation({
    mutationFn: (request: { groupId: string; confirmation: WriteGroupTestWriteConfirmation }) =>
      studioV2WorkspaceWriteGroupsAPI.testWrite(request.groupId, request.confirmation),
    retry: false,
  });
}

/** Reads one test-write operation on demand, e.g. to resolve an unconfirmed result. */
export function useWriteGroupTestWriteOperationQuery(operationId: string | undefined) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.writeGroupTestWriteOperation(operationId ?? ''),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.testWriteOperation(operationId as string),
    enabled: false,
    retry: false,
    refetchOnWindowFocus: false,
  });
}
