import type { WriteGroupListResponse } from '../types/studioV2WriteGroup';
import { parseBoundedJson } from './safeJson';
import { isRecord } from './studioV2WriteGroupJsonPrimitives';
import { MAX_WRITE_GROUP_LIST_LENGTH, parseWriteGroupAPIData, parseWriteGroupListData, WriteGroupResponseError } from './studioV2WriteGroupJson';

const MAX_LIST_BYTES = 8 * 1024 * 1024;
const MAX_RESPONSE_DEPTH = 8;
const operation = 'write-group list';
const bytes = (value: string) => new TextEncoder().encode(value).byteLength;

/** Only the LIST groups collection gets a larger bound; all nested JSON limits remain unchanged. */
export function parseWriteGroupListAPIData(value: unknown): WriteGroupListResponse {
  let raw = value;
  if (typeof raw === 'string') {
    if (raw.length > MAX_LIST_BYTES || bytes(raw) > MAX_LIST_BYTES) throw new WriteGroupResponseError(operation);
    try { raw = JSON.parse(raw) as unknown; } catch { throw new WriteGroupResponseError(operation); }
  }
  if (!isRecord(raw) || raw.success !== true || !isRecord(raw.data) || !Array.isArray(raw.data.groups)) {
    return parseWriteGroupAPIData(raw, operation, parseWriteGroupListData, MAX_RESPONSE_DEPTH);
  }
  const groups = raw.data.groups;
  if (groups.length > MAX_WRITE_GROUP_LIST_LENGTH) throw new WriteGroupResponseError(operation);
  // Validate the entire envelope, including unknown fields, without relaxing any other array.
  const shell = parseBoundedJson({ ...raw, data: { ...raw.data, groups: [] } }, MAX_RESPONSE_DEPTH);
  if (!isRecord(shell) || !isRecord(shell.data)) throw new WriteGroupResponseError(operation);
  let size = bytes(JSON.stringify(shell));
  if (size > MAX_LIST_BYTES) throw new WriteGroupResponseError(operation);
  const boundedGroups: unknown[] = [];
  for (const group of groups) {
    if (!isRecord(group)) throw new WriteGroupResponseError(operation);
    // A group starts at depth three in the original envelope.
    const bounded = parseBoundedJson(group, MAX_RESPONSE_DEPTH - 3);
    if (!isRecord(bounded)) throw new WriteGroupResponseError(operation);
    size += bytes(JSON.stringify(bounded)) + (boundedGroups.length ? 1 : 0);
    if (size > MAX_LIST_BYTES) throw new WriteGroupResponseError(operation);
    boundedGroups.push(bounded);
  }
  shell.data.groups = boundedGroups;
  const parsed = parseWriteGroupListData(shell.data);
  if (!parsed) throw new WriteGroupResponseError(operation, shell);
  return parsed;
}
