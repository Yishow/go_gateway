import { studioV2WorkspaceKeys } from './keys';
/** No current-head/latest fallback: recovery is always keyed by the original root. */
export const studioV2WriteGroupLifecycleKeys = {
  all: () => [...studioV2WorkspaceKeys.writeGroups(), 'lifecycle'] as const,
  lifecycle: (workspaceId: string, groupId: string) => [...studioV2WriteGroupLifecycleKeys.all(), workspaceId, groupId] as const,
  operation: (workspaceId: string, predecessorGroupId: string, lineageId: string, operationId: string) =>
    [...studioV2WriteGroupLifecycleKeys.all(), workspaceId, predecessorGroupId, 'successor', lineageId, operationId] as const,
};
