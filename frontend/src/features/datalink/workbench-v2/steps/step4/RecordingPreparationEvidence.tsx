import { useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '../../../../../hooks/datalink/keys';
import { studioV2WriteGroupLifecycleKeys } from '../../../../../hooks/datalink/studioV2WriteGroupLifecycleKeys';
import { useCallback, useEffect, useLayoutEffect, useMemo, useRef } from 'react';
import { useWriteGroupReadinessQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import { useWriteGroupLifecycleQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroupLifecycle';
import type { RecordingWorkflowEvidence, RecordingWorkflowScope, RecordingLifecycleAvailability, RecordingLifecycleAction } from '../../state/recordingWorkflowView';
import type { WriteGroup, WriteGroupReadiness } from '../../../../../types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '../../../../../types/studioV2WriteGroupLifecycle';
import { code } from '../../../../../utils/studioV2LifecycleJsonPrimitives';

export interface RecordingPreparationEvidenceView {
  readonly scope: RecordingWorkflowScope;
  readonly readiness: RecordingWorkflowEvidence<WriteGroupReadiness>;
  readonly lifecycle: RecordingWorkflowEvidence<RecordingLifecycleAvailability>;
  readonly refresh: () => Promise<void>;
}
export interface RecordingPreparationEvidenceProps {
  scope: RecordingWorkflowScope; group: WriteGroup;
  onEvidence: (evidence: RecordingPreparationEvidenceView | null) => void;
}
interface QueryState<T> { data?: T; isLoading: boolean; isFetching: boolean; isError: boolean }
function queryState(query: QueryState<unknown>, valid: boolean): 'unavailable' | 'loading' | 'error' | undefined {
  if (!valid) return 'unavailable';
  if (query.isError) return 'error';
  if (query.isLoading || query.isFetching) return 'loading';
  if (!query.data) return 'unavailable';
}
function bound(value: WriteGroupReadiness | WriteGroupLifecycle, scope: RecordingWorkflowScope) {
  return value.workspace_id === scope.workspaceId && value.workspace_revision === scope.workspaceRevision &&
    value.group_id === scope.groupId && value.group_revision === scope.groupRevision && value.applied_revision === scope.appliedRevision;
}
function readinessEvidence(query: QueryState<WriteGroupReadiness>, scope: RecordingWorkflowScope, valid: boolean): RecordingWorkflowEvidence<WriteGroupReadiness> {
  const state = queryState(query, valid);
  if (state || !query.data) return { state: state ?? 'unavailable' };
  const value = query.data;
  if (!bound(value, scope) || ![value.config_ready, value.schema_ready, value.ready].every(flag => typeof flag === 'boolean') ||
    !Array.isArray(value.issues) || value.issues.some(issue => !code(issue.code) || !['blocking', 'warning'].includes(issue.severity) ||
      !['Step 1', 'Step 2', 'Step 3', 'Step 4'].includes(issue.step))) return { state: 'unavailable' };
  // The owner needs typed blockers, not raw backend messages, source paths or schema digests.
  return { state: 'available', scope, value: {
    workspace_id: value.workspace_id, workspace_revision: value.workspace_revision, group_id: value.group_id,
    group_revision: value.group_revision, applied_revision: value.applied_revision, config_ready: value.config_ready,
    schema_ready: value.schema_ready, ready: value.ready,
    issues: value.issues.map(issue => ({ code: issue.code, severity: issue.severity, step: issue.step, scope: '', message: '' })),
  } };
}
function lifecycleEvidence(query: QueryState<WriteGroupLifecycle>, scope: RecordingWorkflowScope, valid: boolean,
  connectorId: string, connectorRevision: string, status: WriteGroup['status']): RecordingWorkflowEvidence<RecordingLifecycleAvailability> {
  const state = queryState(query, valid);
  if (state || !query.data) return { state: state ?? 'unavailable' };
  const value = query.data;
  if (!bound(value, scope) || value.connector_id !== connectorId || value.connector_revision !== connectorRevision ||
    value.status !== status || !Array.isArray(value.actions)) return { state: 'unavailable' };
  const actions: Partial<Record<RecordingLifecycleAction, { state: 'available' | 'blocked'; reasonCodes: string[] }>> = {};
  for (const action of value.actions) {
    if (!['edit', 'stop', 'remove', 'recreate', 'transition'].includes(action.action) || actions[action.action] ||
      typeof action.allowed !== 'boolean' || code(action.reason_code, true) === null ||
      (action.allowed ? action.reason_code !== '' : !action.reason_code)) return { state: 'unavailable' };
    actions[action.action] = { state: action.allowed ? 'available' : 'blocked', reasonCodes: action.reason_code ? [action.reason_code] : [] };
  }
  return { state: 'available', scope, value: actions };
}

/** Read-only current evidence reporter. It never creates preview, DDL, start or runtime queries. */
export function RecordingPreparationEvidence({ scope, group, onEvidence }: RecordingPreparationEvidenceProps) {
  const queryClient = useQueryClient();
  const selected = useMemo<RecordingWorkflowScope>(() => ({ workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision,
    deviceId: scope.deviceId, circuitId: scope.circuitId, groupId: scope.groupId, groupRevision: scope.groupRevision, appliedRevision: scope.appliedRevision }),
  [scope.workspaceId, scope.workspaceRevision, scope.deviceId, scope.circuitId, scope.groupId, scope.groupRevision, scope.appliedRevision]);
  const valid = Boolean(selected.workspaceId && selected.workspaceRevision && selected.deviceId && selected.groupId && selected.groupRevision &&
    selected.circuitId !== '' && group.workspace_id === selected.workspaceId && group.id === selected.groupId &&
    group.revision === selected.groupRevision && group.applied_revision === selected.appliedRevision && group.status !== 'deleted' &&
    (group.basic_managed_device_id === undefined || (group.basic_managed_device_id === selected.deviceId &&
      group.members.every(member => member.device_id === group.basic_managed_device_id))) &&
    (group.members.some(member => member.device_id === selected.deviceId && (selected.circuitId === undefined || member.circuit_id === selected.circuitId)) ||
      (group.members.length === 0 && group.basic_managed_device_id === selected.deviceId && selected.circuitId === undefined)));
  const readinessQuery = useWriteGroupReadinessQuery(group.id, valid);
  const lifecycleQuery = useWriteGroupLifecycleQuery(selected.workspaceId, group.id, valid);
  const connectorId = group.destination.connector_id, connectorRevision = group.destination.connector_revision, status = group.status;
  const scopeKey = JSON.stringify([selected, connectorId, connectorRevision, status, valid]);
  const queryKey = JSON.stringify([selected.workspaceId, group.id]);
  const revisionKey = JSON.stringify([selected.workspaceId, selected.workspaceRevision, selected.groupId,
    selected.groupRevision, selected.appliedRevision, connectorId, connectorRevision, status]);
  const fence = useRef({ key: scopeKey, epoch: 0 });
  if (fence.current.key !== scopeKey) fence.current = { key: scopeKey, epoch: fence.current.epoch + 1 };
  const epoch = fence.current.epoch;
  const lease = useMemo(() => ({ active: true, epoch, report: onEvidence }), [epoch, onEvidence]);
  const currentLease = useRef(lease); currentLease.current = lease;
  useLayoutEffect(() => { lease.active = true; return () => { lease.active = false; }; }, [lease]);
  const current = useCallback(() => valid && lease.active && currentLease.current === lease && fence.current.epoch === epoch, [valid, lease, epoch]);
  const readinessRefetch = readinessQuery.refetch, lifecycleRefetch = lifecycleQuery.refetch;
  const refresh = useCallback(async () => {
    if (!current()) return;
    await Promise.all([readinessRefetch({ cancelRefetch: false }), lifecycleRefetch({ cancelRefetch: false })]);
  }, [current, readinessRefetch, lifecycleRefetch]);
  const previous = useRef<{ queryKey: string; groupId: string; revisionKey: string } | undefined>(undefined);
  useEffect(() => {
    const prior = previous.current; previous.current = { queryKey, groupId: group.id, revisionKey };
    if (!prior || prior.revisionKey === revisionKey || !current()) return;
    // Existing query keys omit revisions. Wait out an earlier read, then perform
    // one current revalidation; an A → B → A old continuation cannot issue it.
    const revalidate = async (refetch: typeof readinessRefetch | typeof lifecycleRefetch, key: readonly unknown[]) => {
      // Read the live fetch status, not a render-time flag that may have settled before this effect.
      if (queryClient.getQueryState(key)?.fetchStatus === 'fetching') await refetch({ cancelRefetch: false });
      if (current()) await refetch({ cancelRefetch: false });
    };
    void Promise.all([
      ...(prior.groupId === group.id ? [revalidate(readinessRefetch, studioV2WorkspaceKeys.writeGroupReadiness(group.id))] : []),
      ...(prior.queryKey === queryKey ? [revalidate(lifecycleRefetch, studioV2WriteGroupLifecycleKeys.lifecycle(selected.workspaceId, group.id))] : []),
    ]);
  }, [queryClient, queryKey, revisionKey, current, readinessRefetch, lifecycleRefetch, selected.workspaceId, group.id]);
  const readiness = useMemo(() => readinessEvidence({ data: readinessQuery.data, isLoading: readinessQuery.isLoading, isFetching: readinessQuery.isFetching, isError: readinessQuery.isError }, selected, valid),
    [readinessQuery.data, readinessQuery.isLoading, readinessQuery.isFetching, readinessQuery.isError, selected, valid]);
  const lifecycle = useMemo(() => lifecycleEvidence({ data: lifecycleQuery.data, isLoading: lifecycleQuery.isLoading, isFetching: lifecycleQuery.isFetching, isError: lifecycleQuery.isError }, selected, valid, connectorId, connectorRevision, status),
    [lifecycleQuery.data, lifecycleQuery.isLoading, lifecycleQuery.isFetching, lifecycleQuery.isError, selected, valid, connectorId, connectorRevision, status]);
  useEffect(() => {
    onEvidence({ scope: selected, readiness, lifecycle, refresh });
    return () => onEvidence(null);
  }, [selected, readiness, lifecycle, refresh, onEvidence]);
  return null;
}
