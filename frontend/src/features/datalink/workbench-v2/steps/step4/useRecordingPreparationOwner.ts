import { useEffect, useRef, useState } from 'react';
import { inspectBasicRecordingIntent } from '../../state/basicRecording';
import type { RecordingOwnerObservationState } from '../../state/recordingOwnerObservation';
import { recordingTestWriteRecovery } from '../../state/recordingTestWriteRecovery';

/** Inspection never returns the persisted request or token. Known live intent wins storage failures. */
export function recordingStartObservation(ownerKey: string, live: {
  unresolved: boolean; operationId?: string; status?: string; queryFailed?: boolean; busy?: boolean;
}): RecordingOwnerObservationState {
  if (live.unresolved) return { state: 'unresolved', reason: live.queryFailed ? 'query_failed' :
    live.status === 'pending' || live.status === 'running' ? live.status : live.operationId ? 'unknown' : 'lost_ack',
    ...(live.operationId ? { operationId: live.operationId } : {}) };
  const inspection = inspectBasicRecordingIntent(ownerKey);
  const resolved = Boolean(live.operationId && (live.status === 'succeeded' || live.status === 'failed') &&
    (inspection.state === 'empty' || inspection.state === 'retained' && inspection.operationId === live.operationId));
  if (inspection.state === 'retained' && !resolved) return { state: 'unresolved', reason: inspection.operationId ? 'unknown' : 'lost_ack',
    ...(inspection.operationId ? { operationId: inspection.operationId } : {}) };
  if (inspection.state === 'unavailable') return inspection;
  if (live.busy) return { state: 'unavailable', reason: 'loading' };
  if (resolved) return { state: 'observed', result: 'resolved', status: live.status === 'succeeded' ? 'succeeded' : 'failed',
    operationId: live.operationId, reason: 'terminal' };
  return { state: 'observed', result: 'none' };
}

export function useRecordingPreparationOwner(identity: string, workspaceId: string, groupIds: readonly string[]) {
  const epoch = useRef({ identity, value: 0 });
  if (epoch.current.identity !== identity) epoch.current = { identity, value: epoch.current.value + 1 };
  const [online, setOnline] = useState(() => navigator.onLine);
  const [recovery, setRecovery] = useState(() => recordingTestWriteRecovery.read());
  useEffect(() => {
    const refresh = () => { setOnline(navigator.onLine); setRecovery(recordingTestWriteRecovery.read()); };
    const unsubscribe = recordingTestWriteRecovery.subscribe(refresh);
    window.addEventListener('online', refresh); window.addEventListener('offline', refresh); window.addEventListener('storage', refresh);
    return () => { unsubscribe(); window.removeEventListener('online', refresh); window.removeEventListener('offline', refresh); window.removeEventListener('storage', refresh); };
  }, []);
  const references = recovery.references.filter(value => value.workspace_id === workspaceId && groupIds.includes(value.group_id));
  const reference = references[0];
  const testWrite: RecordingOwnerObservationState = reference
    ? { state: 'unresolved', reason: 'cleanup_pending', operationId: reference.operation_id }
    : recovery.unavailable ? { state: 'unavailable', reason: 'unreadable' } : { state: 'observed', result: 'none' };
  return { online, epoch: epoch.current.value, testWrite, reference, references, recoveryUnavailable: recovery.unavailable,
    isCurrent: (value: number) => epoch.current.value === value,
    recoveryAvailable: () => {
      const current = recordingTestWriteRecovery.read();
      return !current.unavailable && !current.references.some(value => value.workspace_id === workspaceId && groupIds.includes(value.group_id));
    } };
}
