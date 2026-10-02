import type {
  RecordingPlanMigrationIntent,
  RecordingPlanMigrationPreview,
  RecordingPlanMigrationPreviewItem,
  RecordingPlanMigrationRawPlan,
  RecordingPlanMigrationRawValue,
  RecordingPlanMigrationSource,
  WriteGroupMigrationFinding,
  WriteGroupMigrationPreviewItem,
  WriteGroupMigrationReviewResponse,
} from '../types/studioV2WriteGroupMigration';
import {
  boundedString,
  MAX_SAFE_JSON_ARRAY_LENGTH,
  MAX_SAFE_JSON_DEPTH,
  MAX_SAFE_JSON_OBJECT_KEYS,
  MAX_SAFE_JSON_STRING_LENGTH,
  parseBoundedJson,
} from './safeJson';
import { parseWriteGroupMigrationPreviewData } from './studioV2WriteGroupMigrationJson';

const RECORDING_PLAN_ADAPTER_VERSION = 'recording-plan-v1';
const OPEN_WRITE_GROUPS_ACTION = 'open_write_groups';
const PROTOTYPE_KEYS = new Set(['__proto__', 'constructor', 'prototype']);
const REQUIRED_PLAN_COLLECTIONS = ['members', 'streams', 'destinations'] as const;
const INTENT_KEYS = new Set(['source_id', 'plan', 'sources']);
const SOURCE_KEYS = new Set([
  'measurement_id', 'device_id', 'point_id', 'tag_id', 'definition_revision',
  'source_binding_revision', 'series_epoch', 'source_revision', 'mapping_revision', 'status',
]);

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === 'object' && value !== null && !Array.isArray(value) &&
    Object.keys(value).length <= MAX_SAFE_JSON_OBJECT_KEYS;
}

function isSafeRecord(value: unknown): value is JsonRecord {
  return isRecord(value) && Object.keys(value).every((key) => !PROTOTYPE_KEYS.has(key));
}

function readRawString(record: JsonRecord, key: string): string | undefined | null {
  const value = record[key];
  if (value === undefined) return undefined;
  return typeof value === 'string' && value.length <= MAX_SAFE_JSON_STRING_LENGTH ? value : null;
}

function isBoundedRawValue(value: unknown): value is RecordingPlanMigrationRawValue {
  if (value === null || typeof value === 'boolean') return true;
  if (typeof value === 'number') return Number.isFinite(value);
  if (typeof value === 'string') return value.length <= MAX_SAFE_JSON_STRING_LENGTH;
  if (Array.isArray(value)) {
    return value.length <= MAX_SAFE_JSON_ARRAY_LENGTH && value.every(isBoundedRawValue);
  }
  if (!isSafeRecord(value)) return false;
  return Object.values(value).every(isBoundedRawValue);
}

/** Detects prototype keys and cycles before safeJson clones an object. */
function hasUnsafeRawShape(value: unknown, depth = 0, seen = new WeakSet<object>()): boolean {
  if (value === null || typeof value === 'boolean' || typeof value === 'number' || typeof value === 'string') {
    return false;
  }
  if (typeof value !== 'object' || depth > MAX_SAFE_JSON_DEPTH || seen.has(value)) return true;
  seen.add(value);
  if (Array.isArray(value)) {
    return value.length > MAX_SAFE_JSON_ARRAY_LENGTH || value.some((item) => hasUnsafeRawShape(item, depth + 1, seen));
  }
  const record = value as JsonRecord;
  const keys = Object.keys(record);
  return keys.length > MAX_SAFE_JSON_OBJECT_KEYS || keys.some((key) => PROTOTYPE_KEYS.has(key)) ||
    keys.some((key) => hasUnsafeRawShape(record[key], depth + 1, seen));
}

function parseRawPlan(value: unknown): RecordingPlanMigrationRawPlan | null {
  if (!isRecord(value) || hasUnsafeRawShape(value)) return null;
  const bounded = parseBoundedJson(value);
  if (!isSafeRecord(bounded) || !isBoundedRawValue(bounded)) return null;

  const planID = readRawString(bounded, 'id');
  const workspaceID = readRawString(bounded, 'workspace_id');
  const revision = readRawString(bounded, 'revision');
  if (planID === undefined || planID === null || !boundedString(planID) ||
    workspaceID === undefined || workspaceID === null || !boundedString(workspaceID) ||
    revision === undefined || revision === null) return null;

  for (const key of REQUIRED_PLAN_COLLECTIONS) {
    if (!Object.prototype.hasOwnProperty.call(bounded, key)) return null;
    const collection = bounded[key];
    if (collection !== null && (!Array.isArray(collection) || collection.length > MAX_SAFE_JSON_ARRAY_LENGTH)) {
      return null;
    }
  }
  return bounded as RecordingPlanMigrationRawPlan;
}

function parseSource(value: unknown): RecordingPlanMigrationSource | null {
  if (!isSafeRecord(value) || Object.keys(value).some((key) => !SOURCE_KEYS.has(key))) return null;
  const requiredKeys = [
    'measurement_id', 'device_id', 'point_id', 'tag_id', 'definition_revision',
    'source_binding_revision', 'series_epoch', 'status',
  ] as const;
  const required = requiredKeys.map((key) => readRawString(value, key));
  if (required.some((item) => item === undefined || item === null)) return null;
  const status = required[7];
  if (status !== 'resolved' && status !== 'blocked') return null;
  const sourceRevision = readRawString(value, 'source_revision');
  const mappingRevision = readRawString(value, 'mapping_revision');
  if (sourceRevision === null || mappingRevision === null) return null;
  return {
    measurement_id: required[0] as string,
    device_id: required[1] as string,
    point_id: required[2] as string,
    tag_id: required[3] as string,
    definition_revision: required[4] as string,
    source_binding_revision: required[5] as string,
    series_epoch: required[6] as string,
    ...(sourceRevision !== undefined ? { source_revision: sourceRevision } : {}),
    ...(mappingRevision !== undefined ? { mapping_revision: mappingRevision } : {}),
    status,
  };
}

function parseIntent(value: unknown, sourceID: string, workspaceID: string): RecordingPlanMigrationIntent | null {
  if (!isSafeRecord(value) || Object.keys(value).some((key) => !INTENT_KEYS.has(key))) return null;
  const intentSourceID = readRawString(value, 'source_id');
  const rawPlan = parseRawPlan(value.plan);
  const sources = value.sources;
  if (intentSourceID === undefined || intentSourceID === null || !boundedString(intentSourceID) ||
    intentSourceID !== sourceID || !rawPlan ||
    readRawString(rawPlan, 'id') !== sourceID ||
    readRawString(rawPlan, 'workspace_id') !== workspaceID ||
    !Array.isArray(sources) || sources.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const parsedSources = sources.map(parseSource);
  if (parsedSources.some((source) => source === null)) return null;
  return {
    source_id: intentSourceID,
    plan: rawPlan,
    sources: parsedSources as RecordingPlanMigrationSource[],
  };
}

function parsePlanPreviewItem(
  value: unknown,
  base: WriteGroupMigrationPreviewItem,
  workspaceID: string,
): RecordingPlanMigrationPreviewItem | null {
  if (!isSafeRecord(value) || base.status !== 'blocked') return null;
  if (value.before_intent !== undefined && value.before_intent !== null) return null;
  if (value.before_row_group_intent !== undefined && value.before_row_group_intent !== null) return null;
  if (value.candidate_group !== undefined && value.candidate_group !== null) return null;
  if (readRawString(value, 'repair_action') !== OPEN_WRITE_GROUPS_ACTION) return null;
  const intent = parseIntent(value.before_recording_plan_intent, base.source_id, workspaceID);
  if (!intent) return null;
  return {
    source_id: base.source_id,
    source_revision: base.source_revision,
    status: 'blocked',
    differences: base.differences as WriteGroupMigrationFinding[],
    issues: base.issues as WriteGroupMigrationFinding[],
    repair_action: 'open_write_groups',
    before_recording_plan_intent: intent,
  };
}

/** Parse the blocked recording-plan preview without inventing a candidate group. */
export function parseRecordingPlanMigrationPreviewData(value: unknown): RecordingPlanMigrationPreview | null {
  const preview = parseWriteGroupMigrationPreviewData(value);
  if (!preview || preview.adapter_version !== RECORDING_PLAN_ADAPTER_VERSION || !isRecord(value) ||
    !Array.isArray(value.items) || value.items.length !== preview.items.length) return null;
  const items = value.items.map((item, index) => parsePlanPreviewItem(item, preview.items[index]!, preview.workspace_id));
  if (items.some((item) => item === null)) return null;
  return {
    workspace_id: preview.workspace_id,
    workspace_revision: preview.workspace_revision,
    connector_revision: preview.connector_revision,
    adapter_version: 'recording-plan-v1',
    review_digest: preview.review_digest,
    items: items as RecordingPlanMigrationPreviewItem[],
  };
}

/** Recording-plan conversion has no successful review response in this adapter. */
export function parseRecordingPlanMigrationReviewData(_value: unknown): WriteGroupMigrationReviewResponse | null {
  return null;
}
