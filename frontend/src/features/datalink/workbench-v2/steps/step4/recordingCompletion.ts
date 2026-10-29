import type { RecordingStartOperation } from '../../../../../types/studioV2RecordingStart';
import type { BasicRecordingIntentRecord } from '../../state/basicRecording';

export type RecordingCompletionView = { readonly state: 'unavailable'; readonly reason: 'missing' | 'scope_mismatch' }
  | { readonly state: 'available'; readonly operationId: string; readonly status: RecordingStartOperation['status'];
    readonly activatedDeviceIds: readonly string[]; readonly incompleteDeviceIds: readonly string[];
    readonly incompleteGroupIds: readonly string[]; readonly deviceIds: readonly string[]; readonly groupIds: readonly string[];
    readonly handoffDeviceId?: string };

function sameIds(left: readonly string[], right: readonly string[]): boolean {
  return left.length === right.length && left.every(id => typeof id === 'string' && id.trim().length > 0) &&
    right.every(id => typeof id === 'string' && id.trim().length > 0) && new Set(left).size === left.length &&
    new Set(right).size === right.length && left.every(id => right.includes(id));
}

/** Read-only projection of the original owner's exact operation; never authorizes a command. */
export function recordingCompletion(record: BasicRecordingIntentRecord | undefined, operation: RecordingStartOperation | undefined,
  workspaceId: string | undefined, deviceId: string): RecordingCompletionView {
  if (!record?.operation_id || !operation) return { state: 'unavailable', reason: 'missing' };
  const request = record.request, owner = record.scope_key.split('/');
  const deviceIds = request.device_ids, groupIds = request.groups.map(group => group.group_id);
  if (!workspaceId || !deviceId || owner.length !== 3 || owner[0] !== workspaceId || owner[1] !== deviceId ||
    request.workspace_id !== workspaceId || deviceIds.length !== 1 || deviceIds[0] !== deviceId ||
    operation.workspace_id !== workspaceId || operation.operation_id !== record.operation_id ||
    !sameIds(deviceIds, operation.device_ids) || !sameIds(deviceIds, operation.devices.map(device => device.device_id)) ||
    !sameIds(groupIds, operation.groups.map(group => group.group_id))) return { state: 'unavailable', reason: 'scope_mismatch' };
  const activatedDeviceIds = operation.devices.filter(device => device.activated).map(device => device.device_id);
  const incompleteDeviceIds = deviceIds.filter(id => !activatedDeviceIds.includes(id));
  const incompleteGroupIds = operation.groups.filter(group => !group.saved || !group.ready || !group.applied || !group.applied_revision).map(group => group.group_id);
  if (operation.status === 'succeeded' && (operation.stage !== 'complete' || incompleteDeviceIds.length || incompleteGroupIds.length)) {
    return { state: 'unavailable', reason: 'scope_mismatch' };
  }
  const handoffDeviceId = ['succeeded', 'partial'].includes(operation.status) && activatedDeviceIds.length === 1 ? activatedDeviceIds[0] : undefined;
  return { state: 'available', operationId: operation.operation_id, status: operation.status,
    deviceIds: [...deviceIds], groupIds: [...groupIds], activatedDeviceIds, incompleteDeviceIds, incompleteGroupIds,
    ...(handoffDeviceId ? { handoffDeviceId } : {}) };

}
