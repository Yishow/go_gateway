import type { EvidenceResult, RuleEvidenceView } from '../../../../types/runtimeRuleEvidence';
import { inspectBasicRecordingDevice } from '../state/basicRecording';
import { useEffect, useState } from 'react';
import { useQueries, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2RuntimeContextAPI } from '../../../../services/studioV2RuntimeContext';
import { listRecordingTestWriteRecovery } from '../../../../services/recordingTestWriteRecovery';
import { studioV2WorkspaceKeys } from '../../../../hooks/datalink/keys';
import { recordingTestWriteRecovery } from '../state/recordingTestWriteRecovery';
import { recordingHomeReadQueue } from './recordingHomeReadQueue';
import { readRecordingHomeOwners } from './recordingHomeOwnerReads';
import { recordingHomeRuntimeGroup } from './recordingHomeRuntimeEvidence';
import type { RecordingHomeScope, RecordingHomeGroupScope } from './recordingHomeProjection';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import { writeGroupDeliveryQueryOptions } from '../../../../hooks/datalink/useStudioV2WriteGroups';
import { recordingHomeResultSummary, type RecordingHomeCommitHistory } from './recordingHomeResultSummary';

/** Demand is captured by the route, never inferred from changing query results.
 * One existing complete recovery GET plus <=50 deduplicated device reads, all
 * under the shared six-transport limiter. No polling or operation resolution. */
export function useRecordingHomeFacts(scope: RecordingHomeScope, groups: readonly RecordingReadonly<WriteGroup>[], ids: readonly string[]) {
  const client = useQueryClient(); const [, update] = useState(0);
  const identity = JSON.stringify(scope);
  const [retainedHistory, retainHistory] = useState<{ identity: string; values: Record<string, RecordingHomeCommitHistory> }>({ identity, values: {} });
  useEffect(() => {
    const refresh = () => update(value => value + 1);
    const unsubscribe = recordingTestWriteRecovery.subscribe(refresh);
    window.addEventListener('storage', refresh);
    return () => { unsubscribe(); window.removeEventListener('storage', refresh); };
  }, []);
  useEffect(() => client.getQueryCache().subscribe(event => {
    if (event.type !== 'updated') return;
    const key = event.query.queryKey, marker = key.indexOf('home-rule-evidence');
    if (marker >= 0 && key[marker + 1] === scope.workspaceId && key[marker + 2] === scope.workspaceRevision && key[marker + 3] === scope.epoch) update(value => value + 1);
  }), [client, scope.workspaceId, scope.workspaceRevision, scope.epoch]);
  const canonical = groups.filter(group => group.workspace_id === scope.workspaceId);
  const selected = canonical.filter(group => ids.slice(0, 50).includes(group.id));
  const devices = [...new Set(selected.map(group => group.members[0]?.device_id ?? group.basic_managed_device_id).filter((id): id is string => Boolean(id)))];
  const recovery = useQuery({ queryKey: [...studioV2WorkspaceKeys.all, 'home-owner-recovery', scope.workspaceId, scope.workspaceRevision, scope.epoch],
    queryFn: ({ signal }) => recordingHomeReadQueue.run(() => listRecordingTestWriteRecovery(scope.workspaceId,
      { workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision }), signal), enabled: selected.length > 0,
    retry: false, refetchOnWindowFocus: false });
  const runtime = useQueries({ queries: devices.map(deviceId => ({
    queryKey: [...studioV2WorkspaceKeys.runtimeContext(), 'home-rule-evidence', scope.workspaceId, scope.workspaceRevision, scope.epoch, deviceId],
    queryFn: ({ signal }: { signal: AbortSignal }) => recordingHomeReadQueue.run(() => studioV2RuntimeContextAPI.getRuleEvidence({
      workspace_id: scope.workspaceId, workspace_revision: scope.workspaceRevision, device_id: deviceId,
    }, signal), signal), retry: false, refetchOnWindowFocus: false,
  })) });
  const starts = new Map([...new Set(canonical.flatMap(group => [...group.members.map(member => member.device_id), ...(group.basic_managed_device_id ? [group.basic_managed_device_id] : [])]))]
    .map(device => [device, inspectBasicRecordingDevice(scope.workspaceId, device)]));
  const pendingDevices = [...starts].filter(([, inspection]) => inspection.retained.some(intent => !selected.some(group => intent.groupIds.includes(group.id)))).map(([device]) => device);
  const owners = Object.fromEntries(canonical.map(group => {
    const bound: RecordingHomeGroupScope = { ...scope, groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision };
    return [group.id, readRecordingHomeOwners(bound, group, recovery.isSuccess && !recovery.isFetching
      ? { state: 'available', scope, value: recovery.data } : { state: 'unknown', scope }, starts)];
  }));
  const runtimeGroups = Object.fromEntries(canonical.map(group => {
    const device = group.members[0]?.device_id ?? group.basic_managed_device_id;
    const query = runtime[devices.indexOf(device ?? '')];
    // Read already-observed exact-scope cache without mounting off-page observers.
    // Errors/loading/invalidation supersede retained data; revision validation below
    // prevents a cached device response from blessing a changed group head.
    const cached = client.getQueryState<EvidenceResult<RuleEvidenceView>>([...studioV2WorkspaceKeys.runtimeContext(), 'home-rule-evidence', scope.workspaceId, scope.workspaceRevision, scope.epoch, device]);
    const value = query ? query.isSuccess && !query.isFetching ? query.data : undefined
      : cached?.status === 'success' && cached.fetchStatus === 'idle' && !cached.isInvalidated ? cached.data : undefined;
    return [group.id, recordingHomeRuntimeGroup(scope, group, value ? { state: 'available', scope, value } : undefined)];
  }));
  const observedHistories = Object.fromEntries(canonical.map(group => {
    const bound = { ...scope, groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision };
    const options = writeGroupDeliveryQueryOptions(group.id, false, undefined, bound);
    const previous = client.getQueryData<Awaited<ReturnType<typeof options.queryFn>>>(options.queryKey);
    const result = recordingHomeResultSummary({ group, lifecycle: { state: 'unavailable', reason: 'unknown' }, delivery: previous ? { state: 'available', value: previous } : { state: 'unavailable', reason: 'unknown' } }, { runtime: runtimeGroups[group.id] });
    const history: RecordingHomeCommitHistory | undefined = result?.lastCommittedAt ? { workspaceId: scope.workspaceId, groupId: group.id,
      committedAt: result.lastCommittedAt, ...(result.historicalDestination ? { destination: result.historicalDestination } : {}) } : undefined;
    return [group.id, history];
  }));
  useEffect(() => {
    retainHistory(previous => {
      const values = previous.identity === identity ? { ...previous.values } : {};
      for (const [id, value] of Object.entries(observedHistories)) {
        if (value && (!values[id] || Date.parse(value.committedAt) >= Date.parse(values[id].committedAt))) values[id] = value;
      }
      return previous.identity === identity && JSON.stringify(values) === JSON.stringify(previous.values) ? previous : { identity, values };
    });
  }, [identity, observedHistories]);
  const histories = Object.fromEntries(canonical.map(group => {
    const observed = observedHistories[group.id], retained = retainedHistory.identity === identity ? retainedHistory.values[group.id] : undefined;
    return [group.id, retained && (!observed || Date.parse(retained.committedAt) > Date.parse(observed.committedAt)) ? retained : observed];
  }));
  return { owners, runtimeGroups, histories, pendingDevices, retry: () => selected.length ? Promise.all([recovery.refetch(), ...runtime.map(query => query.refetch())]) : Promise.resolve([]) };
}
