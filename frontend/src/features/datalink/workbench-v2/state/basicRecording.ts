import type { RecordingStartOperation, RecordingStartRequest } from '../../../../types/studioV2RecordingStart';
import type { WriteGroup, WriteGroupDestinationDraft, WriteGroupDraft } from '../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from './writeGroup/candidates';

export interface BasicManagedDraft {
  scope_key: string;
  device_id: string;
  draft: WriteGroupDraft;
}

export interface BasicRecordingIntentRecord {
  version?: 1;
  scope_key: string;
  request: RecordingStartRequest;
  operation_id?: string;
}

/** Readiness token rotation preserves intent; persisted Share revisions do not. */
export function recordingStartRevisionsMatch(
  request: RecordingStartRequest | undefined,
  settingsRevision?: string,
  workspaceRevision?: string,
): boolean {
  return Boolean(request) && request?.settings_revision === settingsRevision && request?.workspace_revision === workspaceRevision;
}

function operationMatchesGroup(operation: RecordingStartOperation, group: WriteGroup): boolean {
  const progress = operation.groups.find((item) => item.group_id === group.id);
  return Boolean(progress && progress.group_revision === group.revision &&
    (!progress.applied_revision || progress.applied_revision === group.applied_revision));
}

/** Keeps a request reusable when the operation itself produced the current setup revision. */
export function recordingStartOperationOwnsCurrentSetupRevision(
  operation: RecordingStartOperation | undefined,
  workspaceRevision: string,
  group: WriteGroup | undefined,
): boolean {
  return Boolean(operation?.setup_revision && workspaceRevision && group &&
    operation.setup_revision === workspaceRevision && operationMatchesGroup(operation, group));
}

/** Prefer the newest local POST result while an operation GET cache catches up. */
export function latestRecordingStartOperation(
  operationId: string | undefined,
  queried: RecordingStartOperation | undefined,
  submitted: RecordingStartOperation | undefined,
): RecordingStartOperation | undefined {
  const currentQuery = queried?.operation_id === operationId ? queried : undefined;
  const currentSubmission = submitted?.operation_id === operationId ? submitted : undefined;
  if (!currentQuery) return currentSubmission;
  if (!currentSubmission) return currentQuery;
  const queriedAt = Date.parse(currentQuery.updated_at);
  const submittedAt = Date.parse(currentSubmission.updated_at);
  return Number.isFinite(queriedAt) && queriedAt > submittedAt ? currentQuery : currentSubmission;
}

/** Stable local identity used to reuse one basic managed group per device. */
export function basicRecordingScopeKey(workspaceId: string, deviceId: string, role: string): string {
  return [workspaceId, deviceId, role].map((part) => part.trim()).join('/');
}

export function basicRecordingIntentKey(scopeKey: string): string {
  return `wbv2.basic-recording.${encodeURIComponent(scopeKey)}`;
}

function storage(): Storage | undefined {
  if (typeof window === 'undefined') return undefined;
  try {
    return window.sessionStorage;
  } catch {
    return undefined;
  }
}

export function saveBasicRecordingIntent(record: BasicRecordingIntentRecord): void {
  try {
    storage()?.setItem(basicRecordingIntentKey(record.scope_key), JSON.stringify({ ...record, version: 1 }));
  } catch {
    // Offline recovery is best effort; the request still remains in component state.
  }
}

function safeText(value: unknown, maxLength: number): value is string {
  return typeof value === 'string' && value.trim().length > 0 && value.length <= maxLength;
}

function validStoredRequest(value: unknown, workspaceId: string, deviceId: string): value is RecordingStartRequest {
  if (!value || typeof value !== 'object' || Array.isArray(value)) return false;
  const request = value as Partial<RecordingStartRequest>;
  if (!safeText(request.request_id, 128) || request.workspace_id !== workspaceId ||
    !Array.isArray(request.device_ids) || request.device_ids.length !== 1 ||
    request.device_ids[0] !== deviceId || !Array.isArray(request.groups) ||
    request.groups.length > 1) return false;
  if (!safeText(request.expected_workspace_revision, 256) &&
    !(request.expected_workspace_revision === '' && request.groups.length === 0)) return false;
  if (request.groups.length === 1) {
    const group = request.groups[0];
    if (!group || typeof group !== 'object' || Array.isArray(group) ||
      !safeText(group.group_id, 256) || !safeText(group.expected_group_revision, 256) ||
      !safeText(group.expected_connector_revision, 256) ||
      (group.draft !== undefined && (!group.draft || typeof group.draft !== 'object' || Array.isArray(group.draft)))) return false;
  }
  for (const field of ['readiness_token', 'settings_revision', 'workspace_revision'] as const) {
    if (request[field] !== undefined && !safeText(request[field], 256)) return false;
  }
  return true;
}

export function loadBasicRecordingIntent(scopeKey: string): BasicRecordingIntentRecord | undefined {
  try {
    const raw = storage()?.getItem(basicRecordingIntentKey(scopeKey));
    if (!raw) return undefined;
    const value: unknown = JSON.parse(raw);
    if (!value || typeof value !== 'object' || Array.isArray(value)) return undefined;
    const record = value as Partial<BasicRecordingIntentRecord>;
    const scopeParts = scopeKey.split('/');
    if (record.version !== 1 || scopeParts.length !== 3 || !safeText(scopeParts[0], 256) || !safeText(scopeParts[1], 256) ||
      record.scope_key !== scopeKey || !validStoredRequest(record.request, scopeParts[0], scopeParts[1]) ||
      (record.operation_id !== undefined && !safeText(record.operation_id, 256))) return undefined;
    return { version: 1, scope_key: scopeKey, request: record.request, ...(record.operation_id ? { operation_id: record.operation_id } : {}) };
  } catch {
    return undefined;
  }
}

export function createBasicRecordingRequestId(): string {
  const cryptoObject = typeof globalThis.crypto !== 'undefined' ? globalThis.crypto : undefined;
  if (cryptoObject?.randomUUID) return cryptoObject.randomUUID();
  const bytes = new Uint8Array(16);
  if (cryptoObject?.getRandomValues) cryptoObject.getRandomValues(bytes);
  else bytes.forEach((_value, index) => { bytes[index] = Math.floor(Math.random() * 256); });
  bytes[6] = (bytes[6] & 0x0f) | 0x40;
  bytes[8] = (bytes[8] & 0x3f) | 0x80;
  const hex = Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('');
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}

/**
 * Builds the minimum managed request for each device represented by persisted
 * Point-to-Tag candidates. The server assigns the table and member columns.
 */
export function buildBasicManagedDrafts(
  workspaceId: string,
  candidates: GroupCandidate[],
  destination: WriteGroupDestinationDraft,
  role: string,
): BasicManagedDraft[] {
  const byDevice = new Map<string, GroupCandidate[]>();
  for (const candidate of candidates) {
    const deviceCandidates = byDevice.get(candidate.device_id) ?? [];
    deviceCandidates.push(candidate);
    byDevice.set(candidate.device_id, deviceCandidates);
  }

  return Array.from(byDevice, ([deviceId, deviceCandidates]) => {
    const first = deviceCandidates[0];
    return {
      scope_key: basicRecordingScopeKey(workspaceId, deviceId, role),
      device_id: deviceId,
      draft: {
        workspace_id: workspaceId,
        name: first.device_name.trim() || deviceId,
        members: deviceCandidates.map((candidate) => ({
          device_id: candidate.device_id,
          point_id: candidate.point_id,
          tag_id: candidate.tag_id,
          target_column: '',
          required: true,
        })),
        destination: { ...destination, table_name: '', storage_strategy: 'managed' },
        row_policy: {
          interval_seconds: 60,
          allowed_lateness_seconds: 0,
          incomplete_policy: 'skip_row',
        },
        write_policy: { mode: 'append', dedupe_capability: 'receipt' },
      },
    };
  });
}

/** Converts a saved canonical managed group to the optional start-time draft. */
export function basicManagedGroupDraft(
  group: WriteGroup,
  intervalSeconds: number,
  incompletePolicy: string,
): WriteGroupDraft {
  return {
    workspace_id: group.workspace_id,
    name: group.name,
    members: group.members.map((member) => ({
      device_id: member.device_id,
      point_id: member.point_id,
      tag_id: member.tag_id,
      ...(member.entity_key ? { entity_key: member.entity_key } : {}),
      ...(member.measurement_id ? { measurement_id: member.measurement_id } : {}),
      target_column: member.target_column,
      required: member.required,
      ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
    })),
    destination: {
      connector_id: group.destination.connector_id,
      connector_revision: group.destination.connector_revision,
      table_schema: group.destination.table_schema,
      table_name: group.destination.table_name,
      storage_strategy: group.destination.storage_strategy,
    },
    row_policy: {
      ...group.row_policy,
      interval_seconds: intervalSeconds,
      allowed_lateness_seconds: group.row_policy.allowed_lateness_seconds,
      incomplete_policy: incompletePolicy,
    },
    write_policy: { ...group.write_policy },
  };
}

/** Existing custom, entity-keyed and cross-device groups remain advanced. */
export function isBasicManagedWriteGroup(group: WriteGroup, deviceId: string): boolean {
  if (group.basic_managed_device_id !== deviceId || group.destination.storage_strategy !== 'managed' || group.members.length === 0) {
    return false;
  }
  if (group.members.some((member) => member.device_id !== deviceId || Boolean(member.entity_key?.trim()))) {
    return false;
  }
  const rowPolicy = group.row_policy;
  return !rowPolicy.entity_key_column?.trim() &&
    !(rowPolicy.group_key_columns?.length) &&
    !(rowPolicy.unique_key_columns?.length) &&
    !rowPolicy.value_column?.trim() &&
    !rowPolicy.quality_column?.trim();
}
