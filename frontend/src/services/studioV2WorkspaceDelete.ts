import type { DeleteBacklogDetails, DeviceDeleteReferences, SuccessorDeleteReferences, WriteGroupDeleteReceipt } from '../types/studioV2Delete';
import type { WriteGroupDeleteRequest } from '../types/studioV2WriteGroup';
import { parseWriteGroupAPIData } from '../utils/studioV2WriteGroupJson';
import { isRecord, readRequiredText } from '../utils/studioV2WriteGroupJsonPrimitives';
import { studioV2DatalinkApi } from './studioV2Workspace';

/** Bounded deletion failure, detached from raw backend and transport diagnostics. */
export class DeleteResponseError extends Error {
  readonly outcome: 'failed' | 'unconfirmed' | 'deleted';
  readonly code?: string;
  readonly action?: string;
  readonly details?: DeleteBacklogDetails | DeviceDeleteReferences | SuccessorDeleteReferences;

  constructor(outcome: 'failed' | 'unconfirmed' | 'deleted', code?: string, action?: string, details?: DeleteBacklogDetails | DeviceDeleteReferences | SuccessorDeleteReferences) {
    super('Deletion outcome requires review');
    this.name = 'DeleteResponseError';
    this.outcome = outcome;
    this.code = code;
    this.action = action;
    this.details = details;
  }
}

function parseReferences(value: unknown): DeviceDeleteReferences | undefined {
  if (!isRecord(value) || !Array.isArray(value.recordings) || value.recordings.length > 20 ||
    typeof value.total !== 'number' || !Number.isSafeInteger(value.total) || value.total < 1 || value.total < value.recordings.length) return undefined;
  const recordings = value.recordings.map(item => {
    if (!isRecord(item)) return null;
    const id = readRequiredText(item, 'id');
    const name = readRequiredText(item, 'name');
    return id && name ? { id, name } : null;
  });
  if (recordings.some(item => !item) || new Set(recordings.map(item => item?.id)).size !== recordings.length) return undefined;
  return { recordings: recordings as DeviceDeleteReferences['recordings'], total: value.total };
}

function parseBacklog(value: unknown): DeleteBacklogDetails | undefined {
  if (!isRecord(value) || !Number.isSafeInteger(value.samples_unconsumed) || typeof value.samples_unconsumed !== 'number' ||
    value.samples_unconsumed < 0 || !isRecord(value.outbox)) return undefined;
  const entries = Object.entries(value.outbox);
  if (entries.length > 64 || entries.some(([state, count]) => state.length < 1 || state.length > 128 ||
    typeof count !== 'number' || !Number.isSafeInteger(count) || count < 1)) return undefined;
  return { samples_unconsumed: value.samples_unconsumed, outbox: Object.fromEntries(entries) as Record<string, number> };
}

function deleteError(error: unknown, device = false): DeleteResponseError {
  if (device && isRecord(error) && isRecord(error.response) && error.response.status === 500 &&
    isRecord(error.response.data) && error.response.data.success === false && isRecord(error.response.data.error)) {
    const code = error.response.data.error.code;
    if (code === 'DEVICE_DELETED_SHARE_SCOPE_REFRESH_FAILED' || code === 'DEVICE_DELETED_CONNECTION_CLOSE_FAILED') {
      return new DeleteResponseError('deleted', code, 'reload');
    }
  }
  if (isRecord(error) && isRecord(error.response) && error.response.status === 409 && isRecord(error.response.data) &&
    error.response.data.success === false && isRecord(error.response.data.error)) {
    const payload = error.response.data.error;
    if (payload.code === 'revision_mismatch') return new DeleteResponseError('failed', 'revision_mismatch', 'reload');
    if (payload.code === 'WRITE_GROUP_DELETE_BACKLOG' && payload.action === 'stop_and_wait') {
      const details = parseBacklog(payload.details);
      if (details) return new DeleteResponseError('failed', payload.code, payload.action, details);
    }
    if ((payload.code === 'WRITE_GROUP_DELETE_ACTIVE' && payload.action === 'stop_first') ||
      (payload.code === 'WRITE_GROUP_DELETE_OPERATION_PENDING' && payload.action === 'open_recording'))
      return new DeleteResponseError('failed', payload.code, payload.action);
    if (payload.code === 'WRITE_GROUP_DELETE_HAS_SUCCESSOR' && payload.action === 'open_recording') {
      const raw = isRecord(payload.details) ? payload.details : undefined;
      const refs = raw && parseReferences({ recordings: raw.referenced_by, total: raw.total });
      return new DeleteResponseError('failed', payload.code, payload.action,
        refs ? { referenced_by: refs.recordings, total: refs.total } : undefined);
    }
    if (payload.code === 'DEVICE_REFERENCED_BY_RECORDINGS' && payload.action === 'remove_recordings_first') {
      const details = parseReferences(payload.details);
      if (details) return new DeleteResponseError('failed', payload.code, payload.action, details);
    }
    const code = readRequiredText(payload, 'code');
    if (code) return new DeleteResponseError('failed', code, 'reload');
  }
  if (isRecord(error) && isRecord(error.response) && typeof error.response.status === 'number' &&
    error.response.status >= 400 && error.response.status <= 599) {
    const data = error.response.data;
    const code = isRecord(data) && data.success === false && isRecord(data.error) ? readRequiredText(data.error, 'code') : null;
    return new DeleteResponseError('failed', code ?? undefined, 'reload');
  }
  return new DeleteResponseError('unconfirmed');
}

/** Device deletion succeeds only when the server confirms it explicitly. */
export async function removeWorkspaceDevice(id: string): Promise<void> {
  const res = await studioV2DatalinkApi.delete<unknown>(`/studio-v2/workspace/devices/${encodeURIComponent(id)}`)
    .catch(error => { throw deleteError(error, true); });
  if (!isRecord(res.data) || res.data.success !== true) throw new DeleteResponseError('unconfirmed');
}

/** Only an exact, well-formed receipt establishes that this group was deleted. */
export async function removeWriteGroup(id: string, request: WriteGroupDeleteRequest): Promise<WriteGroupDeleteReceipt> {
  const res = await studioV2DatalinkApi.delete<unknown>(
    `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`, { data: { ...request } },
  ).catch(error => { throw deleteError(error); });
  return parseWriteGroupAPIData(res.data, 'write-group delete', value => {
    if (!isRecord(value) || value.deleted !== true || value.group_id !== id || !readRequiredText(value, 'group_id') ||
      typeof value.lineage_id !== 'string' || !Number.isSafeInteger(value.deleted_generations) ||
      typeof value.deleted_generations !== 'number' || value.deleted_generations < 1 ||
      (value.lineage_id === '' && value.deleted_generations !== 1)) return null;
    return { deleted: true, group_id: id, lineage_id: value.lineage_id, deleted_generations: value.deleted_generations };
  });
}
