import { cachedHomeReads } from './recordingHomeCachedReads';
import { recordingHomeReadQueue } from './recordingHomeReadQueue';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import { useQueries, useQueryClient } from '@tanstack/react-query';
import { useStudioV2WorkspaceQuery } from '../../../../hooks/datalink/useStudioV2Workspace';
import { useWriteGroupsQuery, writeGroupDeliveryQueryOptions } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import { writeGroupLifecycleQueryOptions } from '../../../../hooks/datalink/useStudioV2WriteGroupLifecycle';
import { validRecordingReadScope } from '../../../../hooks/datalink/studioV2RecordingReadScope';
import { recordingHomeProjection, type RecordingHomeEvidence, type RecordingHomeScope, type RecordingHomeGroupScope } from './recordingHomeProjection';
import { recordingHomeGuidance, type RecordingHomeGuidanceInput } from './recordingHomeGuidance';
import type { GroupCandidates } from '../state/writeGroup/candidates';

interface ReadResult<T> {
  data: T | undefined;
  isError: boolean;
  isFetching: boolean;
  isPending: boolean;
}
/** Error/loading takes precedence over retained query data. Never surface raw errors. */
function evidence<T, S extends RecordingHomeScope>(query: ReadResult<T>, scope: S): RecordingHomeEvidence<T> & { scope: S } {
  if (query.isError) return { state: 'error', scope };
  if (query.isFetching || query.isPending) return { state: 'loading', scope };
  if (query.data === undefined) return { state: 'unknown', scope };
  return { state: 'available', scope, value: query.data as RecordingReadonly<T> };
}
export interface RecordingHomeReadInput {
  /** Database setup revision. The owner changes epoch on context/authority change and leave/return. */
  scope: RecordingHomeScope;
  /** Optional independently observed Share revision, never DatabaseSetupRevision or source freshness. */
  bootstrapRevision?: string;
  /** Existing buildGroupCandidates output with its owner's captured scope. No second hydration. */
  sources?: RecordingHomeEvidence<GroupCandidates>;
  owners?: RecordingHomeGuidanceInput['owners'];
  /** Explicit current-page demand. Omitted defaults to first25; never more than50. */
  detailGroupIds?: readonly string[];
  enabled?: boolean;
}
/** Read-only adapter. No mutation hook, autosave owner, command or persisted state is created. */
export function useRecordingHomeEvidence({ scope, bootstrapRevision, sources, owners, detailGroupIds, enabled = true }: RecordingHomeReadInput) {
  const readable = enabled && validRecordingReadScope(scope);
  const workspace = useStudioV2WorkspaceQuery(readable, { ...scope, bootstrapRevision });
  const groups = useWriteGroupsQuery(readable, scope);
  let workspaceEvidence = evidence(workspace, scope);
  // Share/bootstrap and database setup have independent revisions. Never equate them.
  if (workspaceEvidence.state === 'available' && bootstrapRevision !== undefined && (
    workspace.data?.modbus_share?.workspace_id !== scope.workspaceId ||
    workspace.data?.modbus_share?.workspace_revision !== bootstrapRevision
  )) workspaceEvidence = { state: 'unknown', scope };
  const base = { scope, workspace: readable ? workspaceEvidence : { state: 'unknown' as const, scope },
    sources, groups: readable ? evidence(groups, scope) : { state: 'unknown' as const, scope } };
  const validated = recordingHomeProjection(base);
  const savedGroups = validated.groups.state === 'available' ? validated.groups.value.map(item => item.group) : [];
  const queryClient = useQueryClient();
  const requested = new Set((detailGroupIds ?? savedGroups.slice(0, 25).map(group => group.id)).slice(0, 50));
  const allScopes: RecordingHomeGroupScope[] = savedGroups.map(group => ({ ...scope, groupId: group.id,
    groupRevision: group.revision, appliedRevision: group.applied_revision }));
  const scopes = allScopes.filter(bound => requested.has(bound.groupId));
  const lifecycle = useQueries({ queries: scopes.map(bound => ({
    ...writeGroupLifecycleQueryOptions(bound.workspaceId, bound.groupId, readable, bound),
    queryFn: ({ signal }: { signal: AbortSignal }) => recordingHomeReadQueue.run(() =>
      writeGroupLifecycleQueryOptions(bound.workspaceId, bound.groupId, readable, bound, signal).queryFn(), signal),
  })) });
  const delivery = useQueries({ queries: scopes.map(bound => ({
    ...writeGroupDeliveryQueryOptions(bound.groupId, readable, undefined, bound),
    queryFn: ({ signal }: { signal: AbortSignal }) => recordingHomeReadQueue.run(() =>
      writeGroupDeliveryQueryOptions(bound.groupId, readable, undefined, bound, signal).queryFn(), signal),
  })) });
  const cached = cachedHomeReads(queryClient, allScopes);
  const input = { ...base,
    lifecycle: { ...cached.lifecycle, ...Object.fromEntries(scopes.map((bound, index) => [bound.groupId, evidence(lifecycle[index], bound)])) },
    delivery: { ...cached.delivery, ...Object.fromEntries(scopes.map((bound, index) => [bound.groupId, evidence(delivery[index], bound)])) },
  };
  const projection = recordingHomeProjection(input);
  const guidance = recordingHomeGuidance({ evidence: input, owners });
  const failed = [workspace, groups, ...lifecycle, ...delivery].some(query => query.isError);
  return { projection, guidance,
    error: failed ? { code: 'recording_home_read_failed' as const, action: 'retry' as const } : null,
    retry: async () => {
      if (!readable) return;
      await Promise.all([workspace.refetch(), groups.refetch(), ...lifecycle.map(query => query.refetch()), ...delivery.map(query => query.refetch())]);
    },
  };
}
