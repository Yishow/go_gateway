import type {
  WriteGroup,
  WriteGroupDestination,
  WriteGroupLayoutVersion,
  WriteGroupListResponse,
  WriteGroupMember,
  WriteGroupMigration,
  WriteGroupMutationResponse,
  WriteGroupReadiness,
  WriteGroupProfileKind,
  WriteGroupResultValueType,
  WriteGroupRowPolicy,
  WriteGroupStatus,
  WriteGroupStorageStrategy,
  WriteGroupValueRole,
  WriteGroupWritePolicy,
} from '../types/studioV2WriteGroup';
import { DATALINK_DATA_TYPES, type DataType } from '../types/datalink';
import type {
  StudioV2WorkspaceReadinessIssue,
  StudioV2WorkspaceReadinessSeverity,
  StudioV2WorkspaceReadinessStep,
} from '../types/studioV2WorkspaceReadiness';
import {
  boundedString,
  MAX_SAFE_JSON_ARRAY_LENGTH,
  MAX_SAFE_JSON_DEPTH,
  MAX_SAFE_JSON_OBJECT_KEYS,
  normalizeTypedEnvelope,
  parseBoundedJson,
} from './safeJson';
import {
  isRecord,
  type JsonRecord,
  readOptionalBoolean,
  readOptionalInteger,
  readOptionalText,
  readPresentText,
  readRequiredInteger,
  readRequiredText,
} from './studioV2WriteGroupJsonPrimitives';
import { parseRecordingProfile, recordingProfileIntervalSeconds } from './studioV2RecordingProfile';
import { parseWriteGroupExactRecordingSource } from './studioV2ExactRecordingSourceJson';

const WRITE_GROUP_STATUSES = new Set<WriteGroupStatus>([
  'draft', 'ready', 'running', 'disabled', 'deleted',
]);
const WRITE_GROUP_STORAGE_STRATEGIES = new Set<WriteGroupStorageStrategy>(['managed', 'custom']);
const WRITE_GROUP_LAYOUT_VERSIONS = new Set<string>(['managed_wide_v1', 'managed_wide_profile_v1', 'managed_long_v1']);
const WRITE_GROUP_VALUE_ROLES = new Set<WriteGroupValueRole>(['gauge', 'counter_raw', 'state']);
const WRITE_GROUP_SOURCE_VALUE_TYPES = new Set<DataType>(DATALINK_DATA_TYPES);
const WRITE_GROUP_RESULT_VALUE_TYPES = new Set<WriteGroupResultValueType>([
  'bool', 'text', 'int64', 'uint64', 'decimal', 'float64',
]);
const READINESS_SEVERITIES = new Set<StudioV2WorkspaceReadinessSeverity>(['blocking', 'warning']);
const READINESS_STEPS = new Set<StudioV2WorkspaceReadinessStep>(['Step 1', 'Step 2', 'Step 3', 'Step 4']);
const PROTOTYPE_KEYS = new Set(['__proto__', 'constructor', 'prototype']);

type ResponseOutcome = 'failed' | 'unconfirmed';

function parseOptionalStringArray(record: JsonRecord, key: string): string[] | undefined | null {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  if (!Array.isArray(value) || value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const result = value.map((item) => boundedString(item));
  return result.some((item) => !item) ? null : result as string[];
}

function parseOptionalStringMap(record: JsonRecord, key: string): Record<string, string> | undefined | null {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  if (!isRecord(value)) return null;
  const entries = Object.entries(value);
  if (entries.length > MAX_SAFE_JSON_OBJECT_KEYS) return null;
  const parsed: Array<[string, string]> = [];
  for (const [rawKey, rawValue] of entries) {
    const mapKey = boundedString(rawKey);
    const mapValue = boundedString(rawValue);
    if (!mapKey || !mapValue || PROTOTYPE_KEYS.has(rawKey) || PROTOTYPE_KEYS.has(mapKey)) return null;
    parsed.push([mapKey, mapValue]);
  }
  return Object.fromEntries(parsed);
}

/**
 * Reads an optional enum member: absent stays absent, and a present value must be
 * one of the allowed strings. `null`, a non-string or an unknown value is rejected
 * rather than silently dropped.
 */
function readOptionalEnum<T extends string>(
  record: JsonRecord,
  key: string,
  allowed: ReadonlySet<T>,
): T | undefined | null {
  const value = record[key];
  if (value === undefined) return undefined;
  return typeof value === 'string' && allowed.has(value as T) ? value as T : null;
}

function parseWriteGroupMember(value: unknown, profileKind?: WriteGroupProfileKind): WriteGroupMember | null {
  if (!isRecord(value)) return null;
  const deviceId = readRequiredText(value, 'device_id');
  const pointId = readRequiredText(value, 'point_id');
  const tagId = readRequiredText(value, 'tag_id');
  const sourceRevision = readRequiredText(value, 'source_revision');
  const mappingRevision = readRequiredText(value, 'mapping_revision');
  const targetColumn = readRequiredText(value, 'target_column');
  if (!deviceId || !pointId || !tagId || !sourceRevision || !mappingRevision || !targetColumn ||
    typeof value.required !== 'boolean') {
    return null;
  }

  const measurementId = value.measurement_id === undefined || value.measurement_id === null
    ? undefined
    : boundedString(value.measurement_id);
  if (value.measurement_id !== undefined && value.measurement_id !== null && !measurementId) return null;

  const entityKey = value.entity_key === undefined || value.entity_key === null
    ? undefined
    : boundedString(value.entity_key);
  if (value.entity_key !== undefined && value.entity_key !== null && !entityKey) return null;

  const maxAgeSeconds = readOptionalInteger(value, 'max_age_seconds');
  const circuitId = readOptionalText(value, 'circuit_id');
  const circuitName = readOptionalText(value, 'circuit_name');
  const valueRole = profileKind === 'cumulative' && value.value_role === ''
    ? undefined : readOptionalEnum(value, 'value_role', WRITE_GROUP_VALUE_ROLES);
  const sourceValueType = readOptionalEnum(value, 'source_value_type', WRITE_GROUP_SOURCE_VALUE_TYPES);
  const resultValueType = readOptionalEnum(value, 'result_value_type', WRITE_GROUP_RESULT_VALUE_TYPES);
  const exactSource = value.exact_recording_source === undefined
    ? undefined : parseWriteGroupExactRecordingSource(value.exact_recording_source);
  if (maxAgeSeconds === null || circuitId === null || circuitName === null || (circuitName && !circuitId) ||
    valueRole === null || sourceValueType === null || resultValueType === null ||
    (value.exact_recording_source !== undefined && !exactSource) ||
    (exactSource && (exactSource.measurement_id !== measurementId || exactSource.original_data_type !== sourceValueType ||
      (valueRole !== undefined && valueRole !== 'counter_raw') || resultValueType !== 'decimal'))) return null;

  return {
    device_id: deviceId,
    point_id: pointId,
    tag_id: tagId,
    ...(entityKey ? { entity_key: entityKey } : {}),
    source_revision: sourceRevision,
    mapping_revision: mappingRevision,
    ...(measurementId ? { measurement_id: measurementId } : {}),
    target_column: targetColumn,
    required: value.required,
    ...(maxAgeSeconds !== undefined ? { max_age_seconds: maxAgeSeconds } : {}),
    ...(circuitId ? { circuit_id: circuitId } : {}),
    ...(circuitName ? { circuit_name: circuitName } : {}),
    ...(valueRole !== undefined ? { value_role: valueRole } : {}),
    ...(sourceValueType !== undefined ? { source_value_type: sourceValueType } : {}),
    ...(resultValueType !== undefined ? { result_value_type: resultValueType } : {}),
    ...(exactSource ? { exact_recording_source: exactSource } : {}),
  };
}

function parseWriteGroupDestination(value: unknown): WriteGroupDestination | null {
  if (!isRecord(value)) return null;
  const connectorId = readRequiredText(value, 'connector_id');
  const connectorRevision = readRequiredText(value, 'connector_revision');
  const database = readPresentText(value, 'database');
  const tableSchema = readPresentText(value, 'table_schema');
  const tableName = readRequiredText(value, 'table_name');
  const strategy = readRequiredText(value, 'storage_strategy');
  if (!connectorId || !connectorRevision || database === undefined || database === null ||
    tableSchema === undefined || tableSchema === null || !tableName || !strategy ||
    !WRITE_GROUP_STORAGE_STRATEGIES.has(strategy as WriteGroupStorageStrategy)) {
    return null;
  }

  const schemaRevision = readOptionalText(value, 'schema_revision');
  const schemaDigest = readOptionalText(value, 'schema_digest');
  if (schemaRevision === null || schemaDigest === null) return null;

  return {
    connector_id: connectorId,
    connector_revision: connectorRevision,
    database,
    table_schema: tableSchema,
    table_name: tableName,
    storage_strategy: strategy as WriteGroupStorageStrategy,
    ...(schemaRevision ? { schema_revision: schemaRevision } : {}),
    ...(schemaDigest ? { schema_digest: schemaDigest } : {}),
  };
}

function parseWriteGroupRowPolicy(value: unknown, status: WriteGroupStatus): WriteGroupRowPolicy | null {
  if (!isRecord(value)) return null;
  const minimumInterval = status === 'ready' || status === 'running' ? 1 : 0;
  const intervalSeconds = readRequiredInteger(value, 'interval_seconds', minimumInterval);
  const allowedLatenessSeconds = readRequiredInteger(value, 'allowed_lateness_seconds');
  if (intervalSeconds === null || allowedLatenessSeconds === null) return null;

  const optionalKeys = [
    'incomplete_policy', 'entity_key_column', 'value_column', 'quality_column', 'provenance_column',
    'record_key_column', 'bucket_start_column', 'group_id_column', 'device_id_column',
  ] as const;
  const optionalValues = optionalKeys.map((key) => readOptionalText(value, key));
  const groupKeyColumns = parseOptionalStringArray(value, 'group_key_columns');
  const uniqueKeyColumns = parseOptionalStringArray(value, 'unique_key_columns');
  if (optionalValues.some((item) => item === null) || groupKeyColumns === null || uniqueKeyColumns === null) return null;
  return {
    interval_seconds: intervalSeconds,
    allowed_lateness_seconds: allowedLatenessSeconds,
    ...Object.fromEntries(optionalKeys.flatMap((key, index) => {
      const item = optionalValues[index];
      return item ? [[key, item]] : [];
    })),
    ...(groupKeyColumns !== undefined ? { group_key_columns: groupKeyColumns } : {}),
    ...(uniqueKeyColumns !== undefined ? { unique_key_columns: uniqueKeyColumns } : {}),
  } as WriteGroupRowPolicy;
}

function parseWriteGroupWritePolicy(value: unknown): WriteGroupWritePolicy | null {
  if (!isRecord(value)) return null;
  const mode = readOptionalText(value, 'mode');
  const dedupeCapability = readOptionalText(value, 'dedupe_capability');
  const recreateMissing = readOptionalBoolean(value, 'recreate_missing');
  if (mode === null || dedupeCapability === null || recreateMissing === null) return null;
  return {
    ...(mode ? { mode } : {}),
    ...(dedupeCapability ? { dedupe_capability: dedupeCapability } : {}),
    ...(recreateMissing !== undefined ? { recreate_missing: recreateMissing } : {}),
  };
}

function parseWriteGroupMigration(value: unknown): WriteGroupMigration | null {
  if (value === undefined || value === null) return {};
  if (!isRecord(value)) return null;
  const sourceKind = readOptionalText(value, 'source_kind');
  const sourceRevision = readOptionalText(value, 'source_revision');
  const adapterVersion = readOptionalText(value, 'adapter_version');
  const reviewResult = readOptionalText(value, 'review_result');
  const sourceIds = parseOptionalStringArray(value, 'source_ids');
  const legacyRowGroupID = readOptionalText(value, 'legacy_row_group_id');
  const targetMappingPoints = parseOptionalStringMap(value, 'target_mapping_points');
  if (sourceKind === null || sourceRevision === null || adapterVersion === null || reviewResult === null || sourceIds === null ||
    legacyRowGroupID === null || targetMappingPoints === null) {
    return null;
  }
  return {
    ...(sourceKind ? { source_kind: sourceKind } : {}),
    ...(sourceIds ? { source_ids: sourceIds } : {}),
    ...(sourceRevision ? { source_revision: sourceRevision } : {}),
    ...(adapterVersion ? { adapter_version: adapterVersion } : {}),
    ...(reviewResult ? { review_result: reviewResult } : {}),
    ...(legacyRowGroupID ? { legacy_row_group_id: legacyRowGroupID } : {}),
    ...(targetMappingPoints !== undefined ? { target_mapping_points: targetMappingPoints } : {}),
  };
}

export function parseStudioV2WriteGroup(value: unknown): WriteGroup | null {
  if (!isRecord(value)) return null;
  const id = readRequiredText(value, 'id');
  const workspaceId = readRequiredText(value, 'workspace_id');
  const revision = readRequiredText(value, 'revision');
  const appliedRevision = readPresentText(value, 'applied_revision');
  const name = readRequiredText(value, 'name');
  const status = readRequiredText(value, 'status');
  const createdAt = readRequiredText(value, 'created_at');
  const updatedAt = readRequiredText(value, 'updated_at');
  const rawMembers = value.members;
  const basicDevice = value.basic_managed_device_id === undefined
    ? undefined : boundedString(value.basic_managed_device_id);
  if (value.basic_managed_device_id !== undefined && !basicDevice) return null;
  if (!id || !workspaceId || !revision || appliedRevision === undefined || appliedRevision === null || !name ||
    !status || !WRITE_GROUP_STATUSES.has(status as WriteGroupStatus) || !createdAt || !updatedAt ||
    !Array.isArray(rawMembers) || rawMembers.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }

  const profile = parseRecordingProfile(value.recording_profile);
  const profileKind = profile?.kind;
  const members = rawMembers.map((member) => parseWriteGroupMember(member, profileKind));
  if (members.some((member) => member === null)) return null;
  const destination = parseWriteGroupDestination(value.destination);
  const rowPolicy = parseWriteGroupRowPolicy(value.row_policy, status as WriteGroupStatus);
  const writePolicy = parseWriteGroupWritePolicy(value.write_policy);
  const migration = parseWriteGroupMigration(value.migration);
  if (!destination || !rowPolicy || !writePolicy || !migration) return null;
  const layoutVersion = readOptionalText(value, 'layout_version');
  const storageMode = readOptionalText(value, 'storage_mode');
  // An unknown layout or mode, or an interval whose unit was misread, must never reach the editor.
  if (layoutVersion === null || (layoutVersion !== undefined && !WRITE_GROUP_LAYOUT_VERSIONS.has(layoutVersion)) ||
    storageMode === null || (storageMode !== undefined && storageMode !== 'plain' && storageMode !== 'timescale' && storageMode !== 'native_monthly') ||
    profile === null || (profile && recordingProfileIntervalSeconds(profile) !== rowPolicy.interval_seconds)) {
    return null;
  }
  // The long table serves profile groups with a mode; a wide group carries neither a mode nor a long layout.
  // The reduction-aware wide profile layout exists only as a managed profile table, never with a storage mode.
  if ((layoutVersion === 'managed_long_v1') !== (storageMode !== undefined) ||
    (layoutVersion === 'managed_long_v1' && !profile) ||
    (layoutVersion === 'managed_wide_profile_v1' && (!profile || destination.storage_strategy !== 'managed'))) {
    return null;
  }
  const parsedMembers = members as WriteGroupMember[];
  if (parsedMembers.some((member) => member.exact_recording_source !== undefined) &&
    (layoutVersion !== 'managed_long_v1' || destination.storage_strategy !== 'managed')) return null;
  if (profile?.kind === 'cumulative' && parsedMembers.some((member) => member.value_role === 'counter_raw')) return null;
  if (parsedMembers.some((member) => member.exact_recording_source !== undefined &&
    (profile?.kind === 'instantaneous' ? member.value_role !== 'counter_raw' :
      profile?.kind === 'cumulative' ? member.value_role !== undefined : true))) return null;
  if (layoutVersion === 'managed_long_v1' && parsedMembers.some((member) =>
    member.result_value_type === 'decimal' && member.source_value_type !== 'float32' &&
    member.source_value_type !== 'float64' && member.exact_recording_source === undefined &&
    ((profile?.kind === 'cumulative' && member.value_role === undefined) ||
      (profile?.kind === 'instantaneous' && member.value_role === 'counter_raw')))) return null;

  return {
    id,
    workspace_id: workspaceId,
    revision,
    applied_revision: appliedRevision,
    ...(basicDevice ? { basic_managed_device_id: basicDevice } : {}),
    name,
    status: status as WriteGroupStatus,
    members: members as WriteGroupMember[],
    destination,
    row_policy: rowPolicy,
    write_policy: writePolicy,
    migration,
    ...(layoutVersion ? { layout_version: layoutVersion as WriteGroupLayoutVersion } : {}),
    ...(profile ? { recording_profile: profile } : {}),
    ...(storageMode ? { storage_mode: storageMode } : {}),
    created_at: createdAt,
    updated_at: updatedAt,
  };
}

function parseWorkspaceRevision(value: unknown): string | null {
  if (typeof value !== 'string' || value.length > 256) return null;
  return value;
}

export function parseWriteGroupListData(value: unknown): WriteGroupListResponse | null {
  if (!isRecord(value)) return null;
  const workspaceId = readRequiredText(value, 'workspace_id');
  const workspaceRevision = parseWorkspaceRevision(value.workspace_revision);
  const rawGroups = value.groups;
  if (!workspaceId || workspaceRevision === null || !Array.isArray(rawGroups) ||
    rawGroups.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const groups = rawGroups.map(parseStudioV2WriteGroup);
  const seenGroupIds = new Set<string>();
  for (const group of groups) {
    if (!group || group.workspace_id !== workspaceId || seenGroupIds.has(group.id)) return null;
    seenGroupIds.add(group.id);
  }
  return {
    workspace_id: workspaceId,
    workspace_revision: workspaceRevision,
    groups: groups as WriteGroup[],
  };
}

export function parseWriteGroupMutationData(value: unknown): WriteGroupMutationResponse | null {
  if (!isRecord(value)) return null;
  const workspaceRevision = parseWorkspaceRevision(value.workspace_revision);
  const group = parseStudioV2WriteGroup(value.group);
  if (workspaceRevision === null || !group) return null;
  return { workspace_revision: workspaceRevision, group };
}

function parseWriteGroupReadinessIssue(value: unknown): StudioV2WorkspaceReadinessIssue | null {
  if (!isRecord(value)) return null;
  const code = readRequiredText(value, 'code');
  const severity = readRequiredText(value, 'severity');
  const step = readRequiredText(value, 'step');
  const scope = readPresentText(value, 'scope');
  const message = readRequiredText(value, 'message');
  if (!code || !severity || !READINESS_SEVERITIES.has(severity as StudioV2WorkspaceReadinessSeverity) ||
    !step || !READINESS_STEPS.has(step as StudioV2WorkspaceReadinessStep) || scope === undefined || scope === null ||
    !message) {
    return null;
  }
  return {
    code,
    severity: severity as StudioV2WorkspaceReadinessSeverity,
    step: step as StudioV2WorkspaceReadinessStep,
    scope,
    message,
  };
}

export function parseWriteGroupReadinessData(value: unknown): WriteGroupReadiness | null {
  if (!isRecord(value) || typeof value.config_ready !== 'boolean' ||
    typeof value.schema_ready !== 'boolean' || typeof value.ready !== 'boolean') {
    return null;
  }
  const workspaceId = readRequiredText(value, 'workspace_id');
  const workspaceRevision = readRequiredText(value, 'workspace_revision');
  const groupId = readRequiredText(value, 'group_id');
  const groupRevision = readRequiredText(value, 'group_revision');
  const appliedRevision = readPresentText(value, 'applied_revision');
  const schemaDigest = readOptionalText(value, 'schema_digest');
  const rawIssues = value.issues;
  if (!workspaceId || !workspaceRevision || !groupId || !groupRevision || appliedRevision === undefined ||
    appliedRevision === null || schemaDigest === null || !Array.isArray(rawIssues) ||
    rawIssues.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const issues = rawIssues.map(parseWriteGroupReadinessIssue);
  if (issues.some((issue) => issue === null)) return null;
  const hasBlocking = issues.some((issue) => issue?.severity === 'blocking');
  if (value.ready !== (value.config_ready && value.schema_ready && !hasBlocking) ||
    (value.schema_ready && !value.config_ready) ||
    (value.schema_ready && !schemaDigest) ||
    (!value.schema_ready && schemaDigest)) {
    return null;
  }
  return {
    workspace_id: workspaceId,
    workspace_revision: workspaceRevision,
    group_id: groupId,
    group_revision: groupRevision,
    applied_revision: appliedRevision,
    config_ready: value.config_ready,
    schema_ready: value.schema_ready,
    ready: value.ready,
    ...(schemaDigest ? { schema_digest: schemaDigest } : {}),
    issues: issues as StudioV2WorkspaceReadinessIssue[],
  };
}

/** Error for a nominal success response that cannot safely be used by the UI. */
export class WriteGroupResponseError extends Error {
  readonly code?: string;
  readonly action?: string;
  readonly retryable?: boolean;
  readonly request_id?: string;
  readonly outcome: ResponseOutcome;

  constructor(operation: string, source?: unknown, outcome: ResponseOutcome = 'unconfirmed') {
    super(`${operation} response was invalid`);
    this.name = 'WriteGroupResponseError';
    this.outcome = outcome;
    const envelope = normalizeTypedEnvelope(source);
    this.code = envelope.code;
    this.action = envelope.action;
    this.retryable = envelope.retryable;
    this.request_id = envelope.requestId;
  }
}

export function parseWriteGroupAPIData<T>(
  value: unknown, operation: string, parseData: (data: unknown) => T | null, maxDepth: number = MAX_SAFE_JSON_DEPTH,
): T {
  const bounded = parseBoundedJson(value, maxDepth);
  if (!isRecord(bounded) || typeof bounded.success !== 'boolean') {
    throw new WriteGroupResponseError(operation, bounded);
  }
  if (!bounded.success) {
    throw new WriteGroupResponseError(operation, bounded, 'failed');
  }
  if (!Object.prototype.hasOwnProperty.call(bounded, 'data')) {
    throw new WriteGroupResponseError(operation, bounded);
  }
  const parsed = parseData(bounded.data);
  if (parsed === null) throw new WriteGroupResponseError(operation, bounded);
  return parsed;
}
