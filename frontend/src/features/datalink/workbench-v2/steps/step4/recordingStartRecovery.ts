import type { RecordingStartOperation, RecordingStartRequest } from '../../../../../types/studioV2RecordingStart';
import { inspectBasicRecordingIntent } from '../../state/basicRecording';
import { recordingTestWriteRecovery } from '../../state/recordingTestWriteRecovery';
import { inspectSchemaOperation } from './writeGroup/schemaOperationRecovery';

const identity = (value: unknown): value is string => typeof value === 'string' && value.trim().length > 0;
function uniqueIdentities(values: readonly unknown[]): boolean {
  return values.every(identity) && new Set(values).size === values.length;
}
function sameIdentities(left: readonly unknown[], right: readonly unknown[]): boolean {
  return left.length === right.length && uniqueIdentities(left) && uniqueIdentities(right) && left.every(value => right.includes(value));
}

/** Exact retained intent only; partial progress may legitimately advance group/setup revisions. */
export function isOriginalRequestRetry(operation: RecordingStartOperation | undefined, operationId: string | undefined,
  request: RecordingStartRequest | undefined): boolean {
  if (!operation || !request || !identity(operationId) || operation.operation_id !== operationId ||
    operation.action !== 'recording_start' || !['failed', 'partial'].includes(operation.status) || operation.next_action !== 'retry_same_request' ||
    !identity(request.request_id) || !identity(request.workspace_id) || operation.workspace_id !== request.workspace_id ||
    !Array.isArray(request.device_ids) || request.device_ids.length === 0 || !Array.isArray(operation.device_ids) ||
    !Array.isArray(request.groups) || !Array.isArray(operation.groups) || !Array.isArray(operation.devices)) return false;
  return sameIdentities(operation.device_ids, request.device_ids) &&
    sameIdentities(operation.devices.map(device => device?.device_id), request.device_ids) &&
    sameIdentities(operation.groups.map(group => group?.group_id), request.groups.map(group => group?.group_id));
}

/** Denial only: the native caller must retain and validate its real original request before using this. */
export function recordingStartRecoveryBlocked(ownerKey: string, workspaceId: string, groupIds: readonly string[],
  operationId: string | undefined): boolean {
  if (!identity(ownerKey) || !identity(workspaceId) || !identity(operationId) ||
    ownerKey.split('/')[0] !== workspaceId || !uniqueIdentities(groupIds)) return true;
  const start = inspectBasicRecordingIntent(ownerKey);
  if (start.state === 'unavailable' || start.state === 'retained' && start.operationId !== operationId) return true;
  // Matching O1 is the operation being recovered, never proof that any other owner is clear.
  if (groupIds.some(groupId => inspectSchemaOperation(workspaceId, groupId).state !== 'empty')) return true;
  const recovery = recordingTestWriteRecovery.read();
  return recovery.unavailable || recovery.references.some(reference => reference.workspace_id === workspaceId && groupIds.includes(reference.group_id));
}
