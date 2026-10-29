import type { DeliveryAttentionItem, WriteGroupSourceEvidence, WriteGroupStorage } from '../types/studioV2WriteGroupDelivery';
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
import { parsePreciseInstant } from './preciseInstant';
import { safeRuntimeIdentity, safeRuntimeLabel } from './runtimeOperatorSafety';
import { parseWriteGroupSourceEvidence } from './studioV2WriteGroupSourceEvidenceJson';
import { parseWriteGroupStorage } from './studioV2WriteGroupStorageJson';

type JsonRecord = Record<string, unknown>;

const INTAKE_STATES = new Set<WriteGroupIntakeState>(['active', 'retiring', 'blocked', 'not_running']);
const QUOTA_STATES = new Set<WriteGroupQuotaState>(['unconfigured', 'ok', 'warning', 'hard_limit']);
const BUCKET_ISSUE_KINDS = new Set<WriteGroupDeliveryBucketIssueKind>(['skipped', 'no_data', 'truncated']);
const BUCKET_ISSUE_CAUSES = new Set<WriteGroupDeliveryBucketCause>([
  'missing', 'bad', 'stale', 'invalid', 'no_data', 'unavailable', 'boundary',
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

// Mirror sender/reconcile, runtimeview and storage-restore producer codes. Never
// mistake a syntactically plausible unknown driver code for a trusted diagnostic.
const REASONS = new Set(`unavailable none snapshot-unavailable snapshot-config-invalid destination-unavailable destination-missing
 destination-revision-changed destination-unsupported destination-table-missing destination-table-unavailable
 tag-unavailable interval-invalid counter-semantics-unavailable insert-failed storage-restore-in-progress
 storage-recovery-unavailable commit-ambiguous destination-identity-conflict target-blocked destination-rejected-row
 target-unusable target-unavailable payload-digest-mismatch payload-undecodable schema_repair_blocked storage_owner_mismatch
 confirmation-missing ownership-mismatch repair-needs-confirmation connector-revision-changed table-occupied
 recovery-disabled recovery-policy-unconfirmed receipts-missing timescale-unavailable timescale-extension-moved
 permission-denied statement-failed verification-mismatch unknown-ddl-result repair_attempt_abandoned
 stale-reference foreign-scope deleted unsupported effective-proof-missing counter-replay-proof-missing
 counter-replay-proof-conflict counter-replay-proof-unsupported recording-ownership-unavailable recording-ownership-conflict
 recording-handoff-proof-conflict recording-cutover-finalized-overlap
 group-missing sink-missing group-not-applied member-identity-incomplete tag-type-unsupported layout-blocked
 destination-scope-invalid restore-failed dedupe-unsupported table-rebuilt table-created-again`.split(/\s+/));
const SOURCE_REASONS = new Set(`reason_unavailable no-samples type-mismatch counter-unavailable missing-value max-age-exceeded mapping-failed
 quality-bad quality-missing quality-stale quality-invalid quality-uncertain coverage-unverified bad invalid stale non-numeric
 first-reading binding-changed no-eligible-reading decrease-unproven contradictory-same-time non-monotonic-time
 mixed-numeric-basis rollover-evidence-insufficient rollover-range-invalid rollover-multiple-wraps-possible
 rollover-rate-bound-violated non-numeric-reading delta-not-representable`.split(/\s+/));
const QUOTA_NOTICES = new Set(['new readings for this scope are not recorded until space is available; already accepted data is kept',
  'capacity is running low; intake will be refused at the hard limit']);
const diagnostic = (value: string) => value === '' || REASONS.has(value) ? value : 'unavailable';
const sourceReason = (value: string) => value === '' || SOURCE_REASONS.has(value) ? value : 'reason_unavailable';
const identity = (value: unknown): value is string => typeof value === 'string' && value !== '' && safeRuntimeIdentity(value) === value;
const optionalIdentity = (value: unknown) => value === undefined || value === '' || identity(value);
const label = (value: unknown): value is string => typeof value === 'string' && value !== '' && safeRuntimeLabel(value) === value;
const schemaLabel = (value: unknown): value is string => value === '' || label(value);

const readIdentity = (value: unknown) => identity(value) ? value : null;
const readLabel = (value: unknown) => label(value) ? value : null;
const PROOF_FIELDS = new Set(['group_id', 'operation_id', 'group_revision', 'connector_id', 'connector_revision',
  'device_id', 'point_id', 'tag_id', 'circuit_id', 'record_id', 'effect_key', 'payload_digest']);
const TIME_FIELDS = new Set(['last_reconstructed_at', 'occurred_at', 'window_start', 'window_end', 'committed_at']);

/** Inspect only detached parser output, never arbitrary input keys. Native parsers
 * may trim display text; proof identities must instead match the original bytes. */
function unchangedProof(raw: unknown, parsed: unknown): boolean {
  if (Array.isArray(parsed)) return Array.isArray(raw) && parsed.every((value, index) => unchangedProof(raw[index], value));
  if (!isRecord(parsed)) return true;
  if (!isRecord(raw)) return false;
  return Object.entries(parsed).every(([key, value]) => PROOF_FIELDS.has(key)
    ? identity(value) && raw[key] === value
    : TIME_FIELDS.has(key) ? parsePreciseInstant(value) !== null && raw[key] === value : unchangedProof(raw[key], value));
}

function parseAttention(value: unknown): DeliveryAttentionItem | null {
  if (!isRecord(value)) return null;
  const { effect_key, payload_digest, group_revision, connector_id, connector_revision, table_name, table_schema, error_code, state } = value;
  const revision = count(value.state_revision);
  if (!identity(effect_key) || !identity(payload_digest) || !identity(group_revision) || !identity(connector_id) ||
    !identity(connector_revision) || !label(table_name) || !schemaLabel(table_schema) || revision === null ||
    (state !== 'blocked' && state !== 'quarantined' && state !== 'unknown') || typeof error_code !== 'string' || error_code.length > 128) return null;
  // Named fields only: no configuration, unknown payload or borrowed object crosses this boundary.
  return { effect_key, payload_digest, group_revision, connector_id, connector_revision, table_name, table_schema,
    state, state_revision: revision, error_code: diagnostic(error_code) };
}

/** The existing parsers detach these DTOs. Validate proof, then normalize only diagnostic leaves. */
function safeNestedEvidence(storage: WriteGroupStorage | undefined, evidence: WriteGroupSourceEvidence | null | undefined): boolean {
  if (storage) {
    if (storage.repair_reason !== undefined) storage.repair_reason = diagnostic(storage.repair_reason);
    for (const event of storage.events) {
      if (!optionalIdentity(event.group_id) || !optionalIdentity(event.operation_id)) return false;
      if (event.code !== undefined) event.code = diagnostic(event.code);
      if (event.reason !== undefined) event.reason = diagnostic(event.reason);
    }
  }
  if (evidence) {
    if (!identity(evidence.group_revision) || !optionalIdentity(evidence.connector_id) || !optionalIdentity(evidence.connector_revision)) return false;
    for (const source of evidence.sources) {
      if (!identity(source.device_id) || !identity(source.point_id) || !identity(source.tag_id) || !optionalIdentity(source.circuit_id)) return false;
      for (const record of [source.latest, source.last_sql_committed]) {
        if (!record?.available) continue;
        if (![record.group_revision, record.connector_revision, record.record_id, record.effect_key, record.payload_digest].every(identity)) return false;
        if (record.reason !== undefined) record.reason = sourceReason(record.reason);
        if (record.aggregation?.reason !== undefined) record.aggregation.reason = sourceReason(record.aggregation.reason);
      }
    }
  }
  return true;
}

function parseCommittedEffect(value: unknown): WriteGroupCommittedEffect | null {
  if (!isRecord(value)) return null;
  const groupRevision = readIdentity(value.group_revision);
  const connectorRevision = readIdentity(value.connector_revision);
  const recordId = readIdentity(value.record_id);
  const effectKey = readIdentity(value.effect_key);
  const digest = readIdentity(value.payload_digest);
  const committedAt = value.committed_at;
  if (!identity(groupRevision) || !identity(connectorRevision) || !identity(recordId) || !identity(effectKey) || !identity(digest) ||
    !/^[a-f0-9]{64}$/.test(digest) || typeof committedAt !== 'string' || parsePreciseInstant(committedAt) === null) return null;
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
  const groupRevision = readIdentity(value.group_revision);
  const connectorId = readIdentity(value.connector_id);
  const connectorRevision = readIdentity(value.connector_revision);
  const tableName = readLabel(value.table_name);
  const pending = count(value.pending);
  const schema = value.table_schema;
  const codes = value.error_codes;
  if (!identity(groupRevision) || !identity(connectorId) || !identity(connectorRevision) || !label(tableName) || !schemaLabel(schema) || pending === null ||
    typeof schema !== 'string' || schema.length > 256 || !Array.isArray(codes) ||
    codes.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const errorCodes: string[] = [];
  for (const code of codes) {
    const text = boundedString(code);
    if (!text) return null;
    errorCodes.push(diagnostic(text));
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
    ...(value.loss_risk_notice ? { loss_risk_notice: QUOTA_NOTICES.has(value.loss_risk_notice) ? value.loss_risk_notice : 'unavailable' } : {}),
  };
}

function parseBucketIssue(value: unknown): WriteGroupDeliveryBucketIssue | null {
  if (!isRecord(value)) return null;
  const groupRevision = readIdentity(value.group_revision);
  const bucketStart = value.bucket_start;
  const kind = boundedString(value.kind);
  const rawCauses = value.causes;
  if (!identity(groupRevision) || typeof bucketStart !== 'string' || parsePreciseInstant(bucketStart) === null ||
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
  const groupId = readIdentity(value.group_id);
  const intakeState = value.intake.state;
  const stages = parseStages(value.stages);
  const quota = parseQuota(value.quota);
  const oldest = value.oldest_pending_seconds;
  const noData = count(value.no_data_buckets);
  const skipped = count(value.skipped_buckets);
  const rawRecentIssues = value.recent_bucket_issues;
  const rawBacklog = value.backlog;
  const lastCommitted = value.last_sql_committed_at;
  if (!identity(groupId) || typeof intakeState !== 'string' || !INTAKE_STATES.has(intakeState as WriteGroupIntakeState) ||
    !stages || !quota || typeof oldest !== 'number' || !Number.isFinite(oldest) || oldest < 0 ||
    noData === null || skipped === null ||
    (rawRecentIssues !== undefined && (!Array.isArray(rawRecentIssues) || rawRecentIssues.length > MAX_RECENT_BUCKET_ISSUES)) ||
    !Array.isArray(rawBacklog) || rawBacklog.length > MAX_SAFE_JSON_ARRAY_LENGTH ||
    (value.intake.reason !== undefined && typeof value.intake.reason !== 'string')) {
    return null;
  }
  if (lastCommitted !== null && (typeof lastCommitted !== 'string' || parsePreciseInstant(lastCommitted) === null)) {
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
    const groupRevision = readIdentity(rawRevision.group_revision);
    const revisionStages = parseStages(rawRevision.stages);
    if (!identity(groupRevision) || !revisionStages || (effect &&
      (effect.group_revision !== groupRevision || revisionStages.sql_committed === 0))) return null;
    revision = { group_revision: groupRevision, stages: revisionStages };
  }
  const rawAttention = value.attention;
  const attention: DeliveryAttentionItem[] = [];
  if (rawAttention !== undefined) {
    if (!Array.isArray(rawAttention) || rawAttention.length > 20) return null;
    for (const item of rawAttention) {
      const parsed = parseAttention(item);
      if (!parsed) return null;
      attention.push(parsed);
    }
  }
  const backlog = rawBacklog.map(parseBacklog);
  if (backlog.some((entry) => entry === null)) return null;
  const recentIssues = rawRecentIssues?.map(parseBucketIssue);
  if (recentIssues?.some((entry) => entry === null)) return null;
  const storage = parseWriteGroupStorage(value.storage);
  if (storage === null) return null;
  const evidence = parseWriteGroupSourceEvidence(value.source_evidence);
  if (!evidence.ok || !unchangedProof(value.storage, storage) || !unchangedProof(value.source_evidence, evidence.value) ||
    !safeNestedEvidence(storage, evidence.value)) return null;
  return {
    group_id: groupId,
    ...(rawAttention !== undefined ? { attention } : {}),
    intake: { state: intakeState as WriteGroupIntakeState, ...(value.intake.reason ? { reason: diagnostic(value.intake.reason) } : {}) },
    stages, last_sql_committed_at: lastCommitted as string | null, oldest_pending_seconds: oldest,
    ...(rawEffect !== undefined ? { last_sql_committed_effect: effect } : {}),
    ...(rawRevision !== undefined ? { revision_stages: revision } : {}),
    no_data_buckets: noData, skipped_buckets: skipped,
    ...(recentIssues ? { recent_bucket_issues: recentIssues as WriteGroupDeliveryBucketIssue[] } : {}),
    backlog: backlog as WriteGroupDeliveryBacklog[], quota,
    ...(storage ? { storage } : {}),
    ...(evidence.value !== undefined ? { source_evidence: evidence.value } : {}),
  };
}
