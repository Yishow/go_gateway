import type { WriteGroupMember, WriteGroupDestination, WriteGroupRowPolicy, WriteGroupWritePolicy, WriteGroupMigration, WriteGroupStatus, WriteGroupProfileKind, WriteGroupStorageStrategy, WriteGroupValueRole, WriteGroupResultValueType } from '../types/studioV2WriteGroup';
import { DATALINK_DATA_TYPES, type DataType } from '../types/datalink';
import { boundedString, MAX_SAFE_JSON_ARRAY_LENGTH, MAX_SAFE_JSON_OBJECT_KEYS } from './safeJson';
import { isRecord, type JsonRecord, readOptionalBoolean, readOptionalInteger, readOptionalText, readPresentText, readRequiredInteger, readRequiredText } from './studioV2WriteGroupJsonPrimitives';
import { parseWriteGroupExactRecordingSource } from './studioV2ExactRecordingSourceJson';
const WRITE_GROUP_STORAGE_STRATEGIES = new Set<WriteGroupStorageStrategy>(['managed', 'custom']);
const WRITE_GROUP_VALUE_ROLES = new Set<WriteGroupValueRole>(['gauge', 'counter_raw', 'state']);
const WRITE_GROUP_SOURCE_VALUE_TYPES = new Set<DataType>(DATALINK_DATA_TYPES);
const WRITE_GROUP_RESULT_VALUE_TYPES = new Set<WriteGroupResultValueType>(['bool', 'text', 'int64', 'uint64', 'decimal', 'float64']);
const PROTOTYPE_KEYS = new Set(['__proto__', 'constructor', 'prototype']);

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

export function parseWriteGroupMember(value: unknown, profileKind?: WriteGroupProfileKind): WriteGroupMember | null {
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

export function parseWriteGroupDestination(value: unknown, candidate = false): WriteGroupDestination | null {
  if (!isRecord(value)) return null;
  const connectorId = readRequiredText(value, 'connector_id');
  const connectorRevision = readRequiredText(value, 'connector_revision');
  const database = readPresentText(value, 'database');
  const tableSchema = readPresentText(value, 'table_schema');
  const tableName = candidate ? readPresentText(value, 'table_name') : readRequiredText(value, 'table_name');
  const strategy = readRequiredText(value, 'storage_strategy');
  if (!connectorId || !connectorRevision || database === undefined || database === null ||
    tableSchema === undefined || tableSchema === null || tableName === undefined || tableName === null || !strategy ||
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

export function parseWriteGroupRowPolicy(value: unknown, status: WriteGroupStatus): WriteGroupRowPolicy | null {
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

export function parseWriteGroupWritePolicy(value: unknown): WriteGroupWritePolicy | null {
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

export function parseWriteGroupMigration(value: unknown): WriteGroupMigration | null {
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

