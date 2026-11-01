import type { RecordingStartRequest } from '../../../../types/studioV2RecordingStart';
/** Explicit extension of the original start owner, not a new store or API payload. */
export interface SavedRecordingStartOwner { kind: 'saved_record'; groupId: string; deviceIds: string[] }
const identity = (value: unknown): value is string => typeof value === 'string' && value.trim().length > 0 && value.length <= 256 && !value.includes('/');
export function sameRecordingDeviceIds(left: readonly string[], right: readonly string[]): boolean {
  return left.length > 0 && left.length <= 64 && left.length === right.length && left.every(identity) && right.every(identity) &&
    new Set(left).size === left.length && new Set(right).size === right.length && left.every(id => right.includes(id));
}
export function savedRecordingOwnerMatches(value: unknown, scopeKey: string, request: RecordingStartRequest): value is SavedRecordingStartOwner {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return false;
  const owner = value as Partial<SavedRecordingStartOwner>;
  return owner.kind === 'saved_record' && identity(owner.groupId) && Array.isArray(owner.deviceIds) &&
    scopeKey === `${request.workspace_id}/saved-record/${owner.groupId}` && Array.isArray(request.device_ids) &&
    sameRecordingDeviceIds(owner.deviceIds, request.device_ids) && Array.isArray(request.groups) && request.groups.length === 1 &&
    request.groups[0]?.group_id === owner.groupId && request.groups[0]?.draft === undefined;
}
