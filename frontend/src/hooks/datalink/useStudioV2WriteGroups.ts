import { recordingReadKey, validRecordingReadScope, type StudioV2RecordingReadScope } from './studioV2RecordingReadScope';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceWriteGroupsAPI } from '../../services/studioV2WorkspaceWriteGroups';
import type {
  WriteGroupCreateRequest,
  WriteGroupDeleteRequest,
  WriteGroupUpdateRequest,
  WriteGroupSchemaApplyRequest,
  WriteGroupSchemaPreviewRequest,
} from '../../types/studioV2WriteGroup';
import { studioV2WorkspaceKeys } from './keys';

export function useWriteGroupsQuery(enabled: boolean, scope?: StudioV2RecordingReadScope) {
  return useQuery({
    queryKey: recordingReadKey(studioV2WorkspaceKeys.writeGroups(), scope),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.list(),
    enabled: enabled && (!scope || validRecordingReadScope(scope)),
    retry: false,
    refetchOnWindowFocus: false,
  });
}

/** Read-only local and destination-schema readiness of one saved group. */
export function useWriteGroupReadinessQuery(groupId: string | undefined, enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.writeGroupReadiness(groupId ?? ''),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.readiness(groupId as string),
    enabled: enabled && Boolean(groupId),
    retry: false,
    refetchOnWindowFocus: false,
  });
}

/** Where the group's accepted rows stand; only `sql_committed` is destination evidence. */
export function writeGroupDeliveryQueryOptions(groupId: string | undefined, enabled: boolean, refetchIntervalMs?: number, scope?: StudioV2RecordingReadScope, signal?: AbortSignal) {
  return {
    queryKey: recordingReadKey(studioV2WorkspaceKeys.writeGroupDelivery(groupId ?? ''), scope),
    queryFn: () => signal ? studioV2WorkspaceWriteGroupsAPI.delivery(groupId as string, signal) : studioV2WorkspaceWriteGroupsAPI.delivery(groupId as string),
    enabled: enabled && Boolean(groupId) && (!scope || validRecordingReadScope(scope)),
    retry: false, refetchOnWindowFocus: false, refetchInterval: refetchIntervalMs, staleTime: 3000,
  };
}
export function useWriteGroupDeliveryQuery(groupId: string | undefined, enabled: boolean, refetchIntervalMs?: number, scope?: StudioV2RecordingReadScope) {
  return useQuery(writeGroupDeliveryQueryOptions(groupId, enabled, refetchIntervalMs, scope));
}

function useInvalidateWriteGroups() {
  const queryClient = useQueryClient();
  return async () => {
    await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups() });
    await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
  };
}

function useInvalidateManagedSchema() {
  const queryClient = useQueryClient();
  return async () => {
    await Promise.all([
      queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups() }),
      queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() }),
      queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.databaseMetadataAll() }),
    ]);
  };
}

export function useWriteGroupSchemaPreviewMutation() {
  return useMutation({
    mutationFn: (variables: { id: string; request: WriteGroupSchemaPreviewRequest }) =>
      studioV2WorkspaceWriteGroupsAPI.schemaPreview(variables.id, variables.request),
    retry: false,
  });
}

export function useWriteGroupSchemaApplyMutation() {
  const invalidate = useInvalidateManagedSchema();
  return useMutation({
    mutationFn: (variables: { id: string; request: WriteGroupSchemaApplyRequest }) =>
      studioV2WorkspaceWriteGroupsAPI.schemaApply(variables.id, variables.request),
    retry: false,
    onSuccess: invalidate,
  });
}

/** Reads the exact operation once; schema apply callers must not silently retry DDL. */
export function useWriteGroupSchemaOperationQuery(operationId: string | undefined) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.writeGroupSchemaOperation(operationId ?? ''),
    queryFn: () => studioV2WorkspaceWriteGroupsAPI.schemaOperation(operationId as string),
    enabled: false,
    retry: false,
    refetchOnWindowFocus: false,
  });
}

export function useCreateWriteGroupMutation() {
  const invalidate = useInvalidateWriteGroups();
  return useMutation({
    mutationFn: (request: WriteGroupCreateRequest) => studioV2WorkspaceWriteGroupsAPI.create(request),
    retry: false,
    onSettled: invalidate,
  });
}

export function useEnsureCircuitProfileMutation() {
  const invalidate = useInvalidateWriteGroups();
  return useMutation({
    mutationFn: (variables: { request: WriteGroupCreateRequest } | WriteGroupCreateRequest) =>
      studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile('request' in variables ? variables.request : variables),
    retry: false,
    onSettled: invalidate,
  });
}

export function useUpdateWriteGroupMutation() {
  const invalidate = useInvalidateWriteGroups();
  return useMutation({
    mutationFn: (variables: { id: string; request: WriteGroupUpdateRequest }) =>
      studioV2WorkspaceWriteGroupsAPI.update(variables.id, variables.request),
    retry: false,
    onSettled: invalidate,
  });
}

/** Lifecycle actions all carry expected revisions only and never retry on their own. */
export function useWriteGroupLifecycleMutation(action: 'apply' | 'disable' | 'remove') {
  const invalidate = useInvalidateWriteGroups();
  return useMutation({
    mutationFn: (variables: { id: string; request: WriteGroupDeleteRequest }) =>
      studioV2WorkspaceWriteGroupsAPI[action](variables.id, variables.request),
    retry: false,
    onSettled: invalidate,
  });
}
