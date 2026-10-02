import type {
  WriteGroupMigrationCandidate,
  WriteGroupMigrationFinding,
  WriteGroupMigrationIntent,
  WriteGroupMigrationPreview,
  WriteGroupMigrationPreviewItem,
  WriteGroupMigrationReviewResponse,
  WriteGroupMigrationPreviewStatus,
} from '../types/studioV2WriteGroupMigration';
import type {
  WriteGroupDestination,
  WriteGroupMemberDraft,
  WriteGroupMigration,
  WriteGroupRowPolicy,
  WriteGroupStorageStrategy,
  WriteGroupWritePolicy,
  WriteGroup,
} from '../types/studioV2WriteGroup';
import {
  boundedString,
  MAX_SAFE_JSON_ARRAY_LENGTH,
  MAX_SAFE_JSON_OBJECT_KEYS,
} from './safeJson';
import {
  parseStudioV2WriteGroup,
  parseWriteGroupAPIData as parsePersistedWriteGroupAPIData,
  WriteGroupResponseError,
} from './studioV2WriteGroupJson';

const PREVIEW_STATUSES = new Set<WriteGroupMigrationPreviewStatus>(['needs_review', 'blocked']);
const STORAGE_STRATEGIES = new Set<WriteGroupStorageStrategy>(['managed', 'custom']);
const ZERO_GO_TIME = '0001-01-01T00:00:00Z';
const PROTOTYPE_KEYS = new Set(['__proto__', 'constructor', 'prototype']);

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === 'object' && value !== null && !Array.isArray(value) &&
    Object.keys(value).length <= MAX_SAFE_JSON_OBJECT_KEYS;
}

function readPresentText(record: JsonRecord, key: string): string | null | undefined {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  return typeof value === 'string' && value.length <= 256 ? value : null;
}

function readRequiredText(record: JsonRecord, key: string): string | null {
  return boundedString(record[key]) ?? null;
}

function readOptionalText(record: JsonRecord, key: string): string | undefined | null {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  return boundedString(value) ?? null;
}

function readOptionalIdentity(record: JsonRecord, key: string): string | undefined | null {
  const value = readPresentText(record, key);
  if (value === undefined || value === null || value === '') return value === null ? null : undefined;
  return boundedString(value) ?? null;
}

function readOptionalInteger(record: JsonRecord, key: string, minimum = 0): number | undefined | null {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= minimum
    ? value
    : null;
}

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

function parseFinding(value: unknown): WriteGroupMigrationFinding | null {
  if (!isRecord(value)) return null;
  const code = readRequiredText(value, 'code');
  const message = readRequiredText(value, 'message');
  return code && message ? { code, message } : null;
}

function parseMigrationIntent(value: unknown, allowIncompleteSource = false): WriteGroupMigrationIntent | null {
  if (!isRecord(value)) return null;
  const sourceId = readRequiredText(value, 'source_id');
  const deviceId = allowIncompleteSource ? readPresentText(value, 'device_id') : readRequiredText(value, 'device_id');
  const pointId = readRequiredText(value, 'point_id');
  const tagId = readRequiredText(value, 'tag_id');
  const connectorId = readRequiredText(value, 'connector_id');
  const connectorRevision = readRequiredText(value, 'connector_revision');
  const database = readPresentText(value, 'database');
  const tableSchema = readPresentText(value, 'table_schema');
  const tableName = readPresentText(value, 'table_name');
  const columnName = readPresentText(value, 'column_name');
  const writeMode = readPresentText(value, 'write_mode');
  const timestampColumn = readPresentText(value, 'timestamp_column');
  const groupKey = readPresentText(value, 'group_key');
  const writeIntervalSeconds = readOptionalInteger(value, 'write_interval_seconds', Number.MIN_SAFE_INTEGER);
  const intervalSource = readOptionalText(value, 'interval_source');
  if (!sourceId || (deviceId === undefined || deviceId === null) || !pointId || !tagId || !connectorId || !connectorRevision ||
    database === null || tableSchema === undefined || tableSchema === null ||
    tableName === undefined || tableName === null || columnName === undefined || columnName === null ||
    writeMode === undefined || writeMode === null ||
    timestampColumn === null || groupKey === null || writeIntervalSeconds === null || intervalSource === null ||
    typeof value.enabled !== 'boolean') {
    return null;
  }
  return {
    source_id: sourceId,
    device_id: deviceId,
    point_id: pointId,
    tag_id: tagId,
    connector_id: connectorId,
    connector_revision: connectorRevision,
    ...(database !== undefined ? { database } : {}),
    table_schema: tableSchema,
    table_name: tableName,
    column_name: columnName,
    write_mode: writeMode,
    ...(timestampColumn !== undefined ? { timestamp_column: timestampColumn } : {}),
    ...(groupKey !== undefined ? { group_key: groupKey } : {}),
    ...(writeIntervalSeconds !== undefined ? { write_interval_seconds: writeIntervalSeconds } : {}),
    ...(intervalSource ? { interval_source: intervalSource } : {}),
    enabled: value.enabled,
  };
}

function parseRowGroupStringArray(record: JsonRecord, key: string): string[] | null {
  const value = record[key];
  if (value === undefined || value === null) return [];
  if (!Array.isArray(value) || value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const result = value.map((item) => boundedString(item));
  return result.some((item) => !item) ? null : result as string[];
}

function parseRowGroupMigrationIntent(
  value: unknown,
  allowUnresolvedSource = false,
): WriteGroupMigrationPreviewItem['before_row_group_intent'] | null {
  if (!isRecord(value)) return null;
  const sourceId = readRequiredText(value, 'source_id');
  const connectorId = readRequiredText(value, 'connector_id');
  const connectorRevision = readRequiredText(value, 'connector_revision');
  const rowGroupValue = value.row_group;
  const rawMembers = value.members;
  if (!sourceId || !connectorId || !connectorRevision || !isRecord(rowGroupValue) ||
    !Array.isArray(rawMembers) || rawMembers.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;

  const rowGroupID = readRequiredText(rowGroupValue, 'id');
  const rowGroupConnectorID = readRequiredText(rowGroupValue, 'connector_id');
  const tableSchema = readPresentText(rowGroupValue, 'table_schema');
  const tableName = readPresentText(rowGroupValue, 'table_name');
  const memberPointIDs = parseRowGroupStringArray(rowGroupValue, 'member_point_ids');
  const groupKeyColumns = parseRowGroupStringArray(rowGroupValue, 'group_key_columns');
  const uniqueKeyColumns = parseRowGroupStringArray(rowGroupValue, 'unique_key_columns');
  if (!rowGroupID || rowGroupID !== sourceId || !rowGroupConnectorID || rowGroupConnectorID !== connectorId ||
    tableSchema === null || tableName === null || memberPointIDs === null ||
    groupKeyColumns === null || uniqueKeyColumns === null) return null;

  const members = rawMembers.map((member) => parseMigrationIntent(member, allowUnresolvedSource));
  if (members.some((member) => member === null) || members.some((member) =>
    member?.connector_id !== connectorId || member.connector_revision !== connectorRevision)) return null;
  return {
    source_id: sourceId,
    connector_id: connectorId,
    connector_revision: connectorRevision,
    row_group: {
      id: rowGroupID,
      connector_id: rowGroupConnectorID,
      table_schema: tableSchema ?? '',
      table_name: tableName ?? '',
      member_point_ids: memberPointIDs,
      group_key_columns: groupKeyColumns,
      unique_key_columns: uniqueKeyColumns,
    },
    members: members as WriteGroupMigrationIntent[],
  };
}

function parseCandidateMember(value: unknown): WriteGroupMemberDraft | null {
  if (!isRecord(value)) return null;
  const deviceId = readRequiredText(value, 'device_id');
  const pointId = readRequiredText(value, 'point_id');
  const tagId = readRequiredText(value, 'tag_id');
  const sourceRevision = readOptionalIdentity(value, 'source_revision');
  const mappingRevision = readOptionalIdentity(value, 'mapping_revision');
  const measurementId = readOptionalIdentity(value, 'measurement_id');
  const entityKey = readOptionalIdentity(value, 'entity_key');
  const targetColumn = readRequiredText(value, 'target_column');
  const maxAgeSeconds = readOptionalInteger(value, 'max_age_seconds');
  if (!deviceId || !pointId || !tagId || sourceRevision === null || mappingRevision === null ||
    measurementId === null || entityKey === null || !targetColumn || maxAgeSeconds === null || typeof value.required !== 'boolean') {
    return null;
  }
  return {
    device_id: deviceId,
    point_id: pointId,
    tag_id: tagId,
    ...(entityKey ? { entity_key: entityKey } : {}),
    ...(sourceRevision ? { source_revision: sourceRevision } : {}),
    ...(mappingRevision ? { mapping_revision: mappingRevision } : {}),
    ...(measurementId ? { measurement_id: measurementId } : {}),
    target_column: targetColumn,
    required: value.required,
    ...(maxAgeSeconds !== undefined ? { max_age_seconds: maxAgeSeconds } : {}),
  };
}

function parseCandidateDestination(value: unknown): WriteGroupDestination | null {
  if (!isRecord(value)) return null;
  const connectorId = readRequiredText(value, 'connector_id');
  const connectorRevision = readRequiredText(value, 'connector_revision');
  const database = readPresentText(value, 'database');
  const tableSchema = readPresentText(value, 'table_schema');
  const tableName = readRequiredText(value, 'table_name');
  const storageStrategy = readRequiredText(value, 'storage_strategy');
  const schemaRevision = readOptionalIdentity(value, 'schema_revision');
  const schemaDigest = readOptionalIdentity(value, 'schema_digest');
  if (!connectorId || !connectorRevision || database === undefined || database === null ||
    tableSchema === undefined || tableSchema === null || !tableName || !storageStrategy ||
    !STORAGE_STRATEGIES.has(storageStrategy as WriteGroupStorageStrategy) ||
    schemaRevision === null || schemaDigest === null) {
    return null;
  }
  return {
    connector_id: connectorId,
    connector_revision: connectorRevision,
    database,
    table_schema: tableSchema,
    table_name: tableName,
    storage_strategy: storageStrategy as WriteGroupStorageStrategy,
    ...(schemaRevision ? { schema_revision: schemaRevision } : {}),
    ...(schemaDigest ? { schema_digest: schemaDigest } : {}),
  };
}

function parseCandidateRowPolicy(value: unknown): WriteGroupRowPolicy | null {
  if (!isRecord(value)) return null;
  const intervalSeconds = readOptionalInteger(value, 'interval_seconds', 0);
  const allowedLatenessSeconds = readOptionalInteger(value, 'allowed_lateness_seconds', 0);
  const optionalKeys = [
    'incomplete_policy', 'entity_key_column', 'value_column', 'quality_column', 'provenance_column',
  ] as const;
  const optionalValues = optionalKeys.map((key) => readOptionalText(value, key));
  const groupKeyColumns = parseOptionalStringArray(value, 'group_key_columns');
  const uniqueKeyColumns = parseOptionalStringArray(value, 'unique_key_columns');
  if (intervalSeconds === undefined || intervalSeconds === null || allowedLatenessSeconds === undefined ||
    allowedLatenessSeconds === null || optionalValues.some((item) => item === null) ||
    groupKeyColumns === null || uniqueKeyColumns === null) return null;
  return {
    interval_seconds: intervalSeconds,
    allowed_lateness_seconds: allowedLatenessSeconds,
    ...(optionalValues[0] ? { incomplete_policy: optionalValues[0] } : {}),
    ...(optionalValues[1] ? { entity_key_column: optionalValues[1] } : {}),
    ...(groupKeyColumns !== undefined ? { group_key_columns: groupKeyColumns } : {}),
    ...(uniqueKeyColumns !== undefined ? { unique_key_columns: uniqueKeyColumns } : {}),
    ...(optionalValues[2] ? { value_column: optionalValues[2] } : {}),
    ...(optionalValues[3] ? { quality_column: optionalValues[3] } : {}),
    ...(optionalValues[4] ? { provenance_column: optionalValues[4] } : {}),
  };
}

function parseCandidateWritePolicy(value: unknown): WriteGroupWritePolicy | null {
  if (!isRecord(value)) return null;
  const mode = readOptionalText(value, 'mode');
  const dedupeCapability = readOptionalText(value, 'dedupe_capability');
  if (mode === null || dedupeCapability === null) return null;
  return {
    ...(mode ? { mode } : {}),
    ...(dedupeCapability ? { dedupe_capability: dedupeCapability } : {}),
  };
}

function parseCandidateMigration(value: unknown): WriteGroupMigration | null {
  if (value === undefined || value === null) return {};
  if (!isRecord(value)) return null;
  const sourceKind = readOptionalText(value, 'source_kind');
  const sourceIds = parseOptionalStringArray(value, 'source_ids');
  const sourceRevision = readOptionalText(value, 'source_revision');
  const adapterVersion = readOptionalText(value, 'adapter_version');
  const reviewResult = readOptionalText(value, 'review_result');
  const legacyRowGroupID = readOptionalText(value, 'legacy_row_group_id');
  const targetMappingPoints = parseOptionalStringMap(value, 'target_mapping_points');
  if (sourceKind === null || sourceIds === null || sourceRevision === null || adapterVersion === null || reviewResult === null ||
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

function parseCandidateGroup(value: unknown): WriteGroupMigrationCandidate | null {
  if (!isRecord(value)) return null;
  const id = readPresentText(value, 'id');
  const workspaceId = readRequiredText(value, 'workspace_id');
  const revision = readPresentText(value, 'revision');
  const appliedRevision = readPresentText(value, 'applied_revision');
  const name = readRequiredText(value, 'name');
  const status = readRequiredText(value, 'status');
  const createdAt = readPresentText(value, 'created_at');
  const updatedAt = readPresentText(value, 'updated_at');
  const rawMembers = value.members;
  if (id === undefined || id === null || !workspaceId || revision === undefined || revision === null ||
    appliedRevision === undefined || appliedRevision === null || appliedRevision !== '' || !name || status !== 'draft' ||
    createdAt !== ZERO_GO_TIME || updatedAt !== ZERO_GO_TIME || !Array.isArray(rawMembers) ||
    rawMembers.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const members = rawMembers.map(parseCandidateMember);
  if (members.some((member) => member === null)) return null;
  const destination = parseCandidateDestination(value.destination);
  const rowPolicy = parseCandidateRowPolicy(value.row_policy);
  const writePolicy = parseCandidateWritePolicy(value.write_policy);
  const migration = parseCandidateMigration(value.migration);
  if (!destination || !rowPolicy || !writePolicy || !migration) return null;
  return {
    id,
    workspace_id: workspaceId,
    revision,
    applied_revision: appliedRevision,
    name,
    status: 'draft',
    members: members as WriteGroupMemberDraft[],
    destination,
    row_policy: rowPolicy,
    write_policy: writePolicy,
    migration,
    created_at: createdAt,
    updated_at: updatedAt,
  };
}

function parsePreviewItem(
  value: unknown,
  workspaceId: string,
  connectorRevision: string,
): WriteGroupMigrationPreviewItem | null {
  if (!isRecord(value)) return null;
  const sourceId = readRequiredText(value, 'source_id');
  const sourceRevision = readPresentText(value, 'source_revision');
  const status = readRequiredText(value, 'status');
  const rawDifferences = value.differences;
  const rawIssues = value.issues;
  if (!sourceId || sourceRevision === undefined || sourceRevision === null || !status || !PREVIEW_STATUSES.has(status as WriteGroupMigrationPreviewStatus) ||
    !Array.isArray(rawDifferences) || rawDifferences.length > MAX_SAFE_JSON_ARRAY_LENGTH ||
    !Array.isArray(rawIssues) || rawIssues.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const differences = rawDifferences.map(parseFinding);
  const issues = rawIssues.map(parseFinding);
  if (differences.some((finding) => finding === null) || issues.some((finding) => finding === null)) return null;
  if ((status === 'needs_review' && (!sourceRevision || differences.length === 0 || issues.length > 0)) ||
    (status === 'blocked' && (issues.length === 0 || (value.candidate_group !== undefined && value.candidate_group !== null)))) return null;
  const beforeIntent = value.before_intent === undefined || value.before_intent === null
    ? undefined
    : parseMigrationIntent(value.before_intent);
  const beforeRowGroupIntent = value.before_row_group_intent === undefined || value.before_row_group_intent === null
    ? undefined
    : parseRowGroupMigrationIntent(value.before_row_group_intent, status === 'blocked');
  const candidateGroup = value.candidate_group === undefined || value.candidate_group === null
    ? undefined
    : parseCandidateGroup(value.candidate_group);
  if ((value.before_intent !== undefined && value.before_intent !== null && !beforeIntent) ||
    (value.before_row_group_intent !== undefined && value.before_row_group_intent !== null && !beforeRowGroupIntent) ||
    (status === 'needs_review' && (!candidateGroup || (!beforeIntent && !beforeRowGroupIntent))) ||
    (beforeRowGroupIntent && beforeRowGroupIntent.source_id !== sourceId) ||
    (beforeRowGroupIntent && beforeRowGroupIntent.connector_revision !== connectorRevision) ||
    (candidateGroup && (candidateGroup.workspace_id !== workspaceId ||
      candidateGroup.destination.connector_revision !== connectorRevision)) ||
    (value.candidate_group !== undefined && value.candidate_group !== null && !candidateGroup)) return null;
  return {
    source_id: sourceId,
    source_revision: sourceRevision,
    status: status as WriteGroupMigrationPreviewStatus,
    differences: differences as WriteGroupMigrationFinding[],
    issues: issues as WriteGroupMigrationFinding[],
    ...(beforeIntent ? { before_intent: beforeIntent } : {}),
    ...(beforeRowGroupIntent ? { before_row_group_intent: beforeRowGroupIntent } : {}),
    ...(candidateGroup ? { candidate_group: candidateGroup } : {}),
  };
}

/** Parse the read-only migration preview without treating candidates as persisted groups. */
export function parseWriteGroupMigrationPreviewData(value: unknown): WriteGroupMigrationPreview | null {
  if (!isRecord(value)) return null;
  const workspaceId = readRequiredText(value, 'workspace_id');
  const workspaceRevision = readPresentText(value, 'workspace_revision');
  const connectorRevision = readRequiredText(value, 'connector_revision');
  const adapterVersion = readRequiredText(value, 'adapter_version');
  const reviewDigest = readRequiredText(value, 'review_digest');
  const rawItems = value.items;
  if (!workspaceId || workspaceRevision === undefined || workspaceRevision === null || !connectorRevision ||
    !adapterVersion || !reviewDigest || !Array.isArray(rawItems) || rawItems.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const items = rawItems.map((item) => parsePreviewItem(item, workspaceId, connectorRevision));
  if (items.some((item) => item === null)) return null;
  return {
    workspace_id: workspaceId,
    workspace_revision: workspaceRevision,
    connector_revision: connectorRevision,
    adapter_version: adapterVersion,
    review_digest: reviewDigest,
    items: items as WriteGroupMigrationPreviewItem[],
  };
}

function parseWorkspaceRevision(value: unknown): string | null {
  return typeof value === 'string' && value.length <= 256 ? value : null;
}

/** Parse canonical persisted groups returned by a reviewed migration save. */
export function parseWriteGroupMigrationReviewData(value: unknown): WriteGroupMigrationReviewResponse | null {
  if (!isRecord(value)) return null;
  const workspaceId = readRequiredText(value, 'workspace_id');
  const workspaceRevision = parseWorkspaceRevision(value.workspace_revision);
  const connectorRevision = readRequiredText(value, 'connector_revision');
  const rawGroups = value.groups;
  if (!workspaceId || workspaceRevision === null || !connectorRevision || !Array.isArray(rawGroups) ||
    rawGroups.length === 0 || rawGroups.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }

  const groups = rawGroups.map(parseStudioV2WriteGroup);
  const seenGroupIDs = new Set<string>();
  for (const group of groups) {
    if (!group || seenGroupIDs.has(group.id) || group.workspace_id !== workspaceId ||
      group.destination.connector_revision !== connectorRevision) {
      return null;
    }
    seenGroupIDs.add(group.id);
  }
  return {
    workspace_id: workspaceId,
    workspace_revision: workspaceRevision,
    connector_revision: connectorRevision,
    groups: groups as WriteGroup[],
  };
}

/** Validate the standard envelope without imposing canonical-group depth on a candidate draft. */
export function parseWriteGroupMigrationAPIData<T>(
  value: unknown,
  operation: string,
  parseData: (data: unknown) => T | null,
): T {
  // Preview candidates are intentionally nested one level deeper than
  // persisted groups; keep their established parser depth while reusing the
  // canonical parser's bounded error surface for review responses.
  if (operation === 'write-group migration review' || operation === 'write-group row-group migration review') {
    return parsePersistedWriteGroupAPIData(value, operation, parseData);
  }
  if (!isRecord(value) || typeof value.success !== 'boolean') {
    throw new WriteGroupResponseError(operation, value);
  }
  if (!value.success) {
    throw new WriteGroupResponseError(operation, value, 'failed');
  }
  if (!Object.prototype.hasOwnProperty.call(value, 'data')) {
    throw new WriteGroupResponseError(operation, value);
  }
  const parsed = parseData(value.data);
  if (parsed === null) throw new WriteGroupResponseError(operation, value);
  return parsed;
}
