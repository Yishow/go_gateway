import type {
  WriteGroupDeliveryBucketCause,
  WriteGroupDeliveryBucketIssue,
  WriteGroupDeliveryBucketIssueKind,
  WriteGroupDelivery,
  WriteGroupCommittedEffect,
  WriteGroupDeliveryBacklog,
  WriteGroupDeliveryQuota,
  WriteGroupDeliveryStages,
  WriteGroupIntakeState,
  WriteGroupQuotaState,
} from '../types/studioV2WriteGroupDelivery';
import { boundedString, MAX_SAFE_JSON_ARRAY_LENGTH } from './safeJson';

type JsonRecord = Record<string, unknown>;

const INTAKE_STATES = new Set<WriteGroupIntakeState>(['active', 'retiring', 'blocked', 'not_running']);
const QUOTA_STATES = new Set<WriteGroupQuotaState>(['unconfigured', 'ok', 'warning', 'hard_limit']);
const BUCKET_ISSUE_KINDS = new Set<WriteGroupDeliveryBucketIssueKind>(['skipped', 'no_data']);
const BUCKET_ISSUE_CAUSES = new Set<WriteGroupDeliveryBucketCause>([
  'missing', 'bad', 'stale', 'invalid', 'no_data', 'unavailable',
]);
const MAX_RECENT_BUCKET_ISSUES = 20;
const STAGE_KEYS: (keyof WriteGroupDeliveryStages)[] = [
  'collecting', 'queued', 'retrying', 'blocked', 'quarantined', 'unknown', 'sql_committed', 'skipped',
];

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function count(value: unknown): number | null {
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= 0 ? value : null;
}

function parseCommittedEffect(value: unknown): WriteGroupCommittedEffect | null {
  if (!isRecord(value)) return null;
  const groupRevision = boundedString(value.group_revision);
  const connectorRevision = boundedString(value.connector_revision);
  const recordId = boundedString(value.record_id);
  const effectKey = boundedString(value.effect_key);
  const digest = boundedString(value.payload_digest);
  const committedAt = boundedString(value.committed_at);
  if (!groupRevision || !connectorRevision || !recordId || !effectKey || !digest ||
    !/^[a-f0-9]{64}$/.test(digest) || !committedAt || Number.isNaN(Date.parse(committedAt))) return null;
  return { group_revision: groupRevision, connector_revision: connectorRevision, record_id: recordId,
    effect_key: effectKey, payload_digest: digest, committed_at: committedAt };
}

function parseStages(value: unknown): WriteGroupDeliveryStages | null {
  if (!isRecord(value)) return null;
  const stages = {} as WriteGroupDeliveryStages;
  for (const key of STAGE_KEYS) {
    const parsed = count(value[key]);
    if (parsed === null) return null;
    stages[key] = parsed;
  }
  return stages;
}

function parseBacklog(value: unknown): WriteGroupDeliveryBacklog | null {
  if (!isRecord(value)) return null;
  const groupRevision = boundedString(value.group_revision);
  const connectorId = boundedString(value.connector_id);
  const connectorRevision = boundedString(value.connector_revision);
  const tableName = boundedString(value.table_name);
  const pending = count(value.pending);
  const schema = value.table_schema;
  const codes = value.error_codes;
  if (!groupRevision || !connectorId || !connectorRevision || !tableName || pending === null ||
    typeof schema !== 'string' || schema.length > 256 || !Array.isArray(codes) ||
    codes.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const errorCodes: string[] = [];
  for (const code of codes) {
    const text = boundedString(code);
    if (!text) return null;
    errorCodes.push(text);
  }
  return {
    group_revision: groupRevision, connector_id: connectorId, connector_revision: connectorRevision,
    table_schema: schema, table_name: tableName, pending, error_codes: errorCodes,
  };
}

function parseQuota(value: unknown): WriteGroupDeliveryQuota | null {
  if (!isRecord(value) || typeof value.configured !== 'boolean' || typeof value.intake_refused !== 'boolean') {
    return null;
  }
  const state = value.state;
  const used = count(value.used_bytes);
  const max = count(value.max_bytes);
  const scope = value.scope;
  if (typeof state !== 'string' || !QUOTA_STATES.has(state as WriteGroupQuotaState) || used === null || max === null ||
    (scope !== 'global' && scope !== 'group' && scope !== '') ||
    (value.loss_risk_notice !== undefined && typeof value.loss_risk_notice !== 'string')) {
    return null;
  }
  // An unconfigured quota cannot also claim to be refusing intake.
  if ((state === 'unconfigured') === value.configured || (state !== 'hard_limit' && value.intake_refused)) {
    return null;
  }
  return {
    configured: value.configured, state: state as WriteGroupQuotaState, scope, used_bytes: used, max_bytes: max,
    intake_refused: value.intake_refused,
    ...(value.loss_risk_notice ? { loss_risk_notice: value.loss_risk_notice as string } : {}),
  };
}

function parseBucketIssue(value: unknown): WriteGroupDeliveryBucketIssue | null {
  if (!isRecord(value)) return null;
  const groupRevision = boundedString(value.group_revision);
  const bucketStart = boundedString(value.bucket_start);
  const kind = boundedString(value.kind);
  const rawCauses = value.causes;
  if (!groupRevision || !bucketStart || Number.isNaN(Date.parse(bucketStart)) ||
    !kind || !BUCKET_ISSUE_KINDS.has(kind as WriteGroupDeliveryBucketIssueKind) ||
    !Array.isArray(rawCauses) || rawCauses.length === 0 || rawCauses.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const causes: WriteGroupDeliveryBucketCause[] = [];
  for (const rawCause of rawCauses) {
    const cause = boundedString(rawCause);
    if (!cause) return null;
    causes.push(BUCKET_ISSUE_CAUSES.has(cause as WriteGroupDeliveryBucketCause)
      ? cause as WriteGroupDeliveryBucketCause
      : 'unavailable');
  }
  return {
    group_revision: groupRevision,
    bucket_start: bucketStart,
    kind: kind as WriteGroupDeliveryBucketIssueKind,
    causes,
  };
}

/**
 * Parses one group's delivery status. Contradictions are rejected rather than
 * displayed: a committed time without committed rows (or the reverse) would
 * make buffered data look written.
 */
export function parseWriteGroupDeliveryData(value: unknown): WriteGroupDelivery | null {
  if (!isRecord(value) || !isRecord(value.intake)) return null;
  const groupId = boundedString(value.group_id);
  const intakeState = value.intake.state;
  const stages = parseStages(value.stages);
  const quota = parseQuota(value.quota);
  const oldest = value.oldest_pending_seconds;
  const noData = count(value.no_data_buckets);
  const skipped = count(value.skipped_buckets);
  const rawRecentIssues = value.recent_bucket_issues;
  const rawBacklog = value.backlog;
  const lastCommitted = value.last_sql_committed_at;
  if (!groupId || typeof intakeState !== 'string' || !INTAKE_STATES.has(intakeState as WriteGroupIntakeState) ||
    !stages || !quota || typeof oldest !== 'number' || !Number.isFinite(oldest) || oldest < 0 ||
    noData === null || skipped === null ||
    (rawRecentIssues !== undefined && (!Array.isArray(rawRecentIssues) || rawRecentIssues.length > MAX_RECENT_BUCKET_ISSUES)) ||
    !Array.isArray(rawBacklog) || rawBacklog.length > MAX_SAFE_JSON_ARRAY_LENGTH ||
    (value.intake.reason !== undefined && typeof value.intake.reason !== 'string')) {
    return null;
  }
  if (lastCommitted !== null && (typeof lastCommitted !== 'string' || Number.isNaN(Date.parse(lastCommitted)))) {
    return null;
  }
  if ((stages.sql_committed > 0) !== (lastCommitted !== null)) return null;
  const rawEffect = value.last_sql_committed_effect;
  const effect = rawEffect == null ? null : parseCommittedEffect(rawEffect);
  if (rawEffect != null && (!effect || stages.sql_committed === 0)) return null;
  const rawRevision = value.revision_stages;
  let revision: WriteGroupDelivery['revision_stages'];
  if (rawRevision === null) revision = null;
  else if (rawRevision !== undefined) {
    if (!isRecord(rawRevision)) return null;
    const groupRevision = boundedString(rawRevision.group_revision);
    const revisionStages = parseStages(rawRevision.stages);
    if (!groupRevision || !revisionStages || (effect &&
      (effect.group_revision !== groupRevision || revisionStages.sql_committed === 0))) return null;
    revision = { group_revision: groupRevision, stages: revisionStages };
  }
  const backlog = rawBacklog.map(parseBacklog);
  if (backlog.some((entry) => entry === null)) return null;
  const recentIssues = rawRecentIssues?.map(parseBucketIssue);
  if (recentIssues?.some((entry) => entry === null)) return null;
  return {
    group_id: groupId,
    intake: { state: intakeState as WriteGroupIntakeState, ...(value.intake.reason ? { reason: value.intake.reason as string } : {}) },
    stages, last_sql_committed_at: lastCommitted as string | null, oldest_pending_seconds: oldest,
    ...(rawEffect !== undefined ? { last_sql_committed_effect: effect } : {}),
    ...(rawRevision !== undefined ? { revision_stages: revision } : {}),
    no_data_buckets: noData, skipped_buckets: skipped,
    ...(recentIssues ? { recent_bucket_issues: recentIssues as WriteGroupDeliveryBucketIssue[] } : {}),
    backlog: backlog as WriteGroupDeliveryBacklog[], quota,
  };
}
