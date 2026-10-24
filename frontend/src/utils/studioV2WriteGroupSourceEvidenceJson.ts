import type {
  WriteGroupSourceDelivery,
  WriteGroupSourceEvidence,
  WriteGroupSourceEvidenceItem,
  WriteGroupSourceEvidenceUnavailableReason,
  WriteGroupSourceRecord,
  WriteGroupSourceRecordUnavailable,
  WriteGroupSourceStatus,
  WriteGroupSourceUnavailableCause,
} from '../types/studioV2WriteGroupDelivery';
import { comparePreciseInstants, parsePreciseInstant } from './preciseInstant';
import { boundedString, MAX_SAFE_JSON_ARRAY_LENGTH } from './safeJson';

type JsonRecord = Record<string, unknown>;
type Parsed<T> = { ok: true; value: T } | { ok: false };

const STATUSES = new Set<WriteGroupSourceStatus>(['ok', 'partial', 'baseline', 'missing', 'bad', 'stale', 'invalid', 'reset_or_uncertain', 'gap']);
const DELIVERIES = new Set<WriteGroupSourceDelivery>(['queued', 'retrying', 'blocked', 'quarantined', 'unknown', 'sql_committed', 'operator_skipped']);
const CAUSES = new Set<WriteGroupSourceUnavailableCause>(['payload_invalid', 'digest_mismatch', 'identity_mismatch', 'status_unknown', 'window_invalid']);
const UNAVAILABLE_REASONS = new Set<WriteGroupSourceEvidenceUnavailableReason>(['not_long_layout', 'applied_snapshot_unavailable', 'revision_mismatch']);
/** Statuses whose record stores NULL in every value slot, and those that always store a reading. */
const NULL_STATUSES = new Set<WriteGroupSourceStatus>(['missing', 'bad', 'stale', 'invalid']);
const READING_STATUSES = new Set<WriteGroupSourceStatus>(['ok', 'partial', 'baseline']);
const STABLE_REASON = /^[a-z0-9_-]{1,64}$/;
const REASON_UNAVAILABLE = 'reason_unavailable';

const isRecord = (value: unknown): value is JsonRecord => typeof value === 'object' && value !== null && !Array.isArray(value);
const instant = (value: unknown) => parsePreciseInstant(boundedString(value));
const oneOf = <T extends string>(set: Set<T>, value: unknown): T | null => (typeof value === 'string' && set.has(value as T) ? value as T : null);

function parseRecord(value: unknown): WriteGroupSourceRecord | WriteGroupSourceRecordUnavailable | null {
  if (!isRecord(value)) return null;
  const delivery = oneOf(DELIVERIES, value.delivery);
  if (!delivery) return null;
  if (value.available === false) {
    const cause = oneOf(CAUSES, value.cause);
    return cause ? { available: false, cause, delivery } : null;
  }
  const status = oneOf(STATUSES, value.status);
  const start = instant(value.window_start);
  const end = instant(value.window_end);
  const digest = boundedString(value.payload_digest);
  const strings = ['group_revision', 'connector_revision', 'record_id', 'effect_key'].map((key) => boundedString(value[key]));
  const profile = value.profile_kind;
  const cumulative = profile === 'cumulative';
  if (value.available !== true || !status || start === null || end === null || comparePreciseInstants(end, start) <= 0 || !digest || !/^[a-f0-9]{64}$/.test(digest) ||
    strings.some((text) => !text) || (profile !== 'instantaneous' && !cumulative) || typeof value.raw_present !== 'boolean' ||
    (cumulative ? typeof value.delta_known !== 'boolean' : value.delta_known !== undefined)) return null;
  if ((NULL_STATUSES.has(status) && value.raw_present) || (READING_STATUSES.has(status) && !value.raw_present) ||
    (status === 'baseline' && value.delta_known === true)) return null;
  const committed = delivery === 'sql_committed';
  if (committed !== (value.committed_at !== undefined) || (committed && instant(value.committed_at) === null)) return null;
  const reason = typeof value.reason === 'string' && value.reason !== ''
    ? (STABLE_REASON.test(value.reason) ? value.reason : REASON_UNAVAILABLE)
    : undefined;
  return {
    available: true, group_revision: strings[0] as string, connector_revision: strings[1] as string, record_id: strings[2] as string,
    effect_key: strings[3] as string, payload_digest: digest, profile_kind: profile as WriteGroupSourceRecord['profile_kind'],
    window_start: value.window_start as string, window_end: value.window_end as string, status,
    ...(reason ? { reason } : {}), raw_present: value.raw_present, ...(cumulative ? { delta_known: value.delta_known as boolean } : {}),
    delivery, ...(committed ? { committed_at: value.committed_at as string } : {}),
  };
}

function parseSource(value: unknown): WriteGroupSourceEvidenceItem | null {
  if (!isRecord(value)) return null;
  const device = boundedString(value.device_id);
  const point = boundedString(value.point_id);
  const tag = boundedString(value.tag_id);
  const circuit = value.circuit_id;
  if (!device || !point || !tag || (circuit !== undefined && typeof circuit !== 'string')) return null;
  const parsed = [value.latest, value.last_sql_committed].map((record) => (record === null ? null : parseRecord(record)));
  if ((value.latest !== null && !parsed[0]) || (value.last_sql_committed !== null && !parsed[1])) return null;
  // The confirmed record is destination evidence, so it must be an available SQL committed one.
  const confirmed = parsed[1];
  if (confirmed && (!confirmed.available || confirmed.delivery !== 'sql_committed')) return null;
  return { ...(circuit ? { circuit_id: circuit } : {}), device_id: device, point_id: point, tag_id: tag, latest: parsed[0], last_sql_committed: parsed[1] };
}

/** Parses the per source evidence of a delivery answer; any contradiction rejects it rather than guessing success. */
export function parseWriteGroupSourceEvidence(value: unknown): Parsed<WriteGroupSourceEvidence | null | undefined> {
  if (value === undefined || value === null) return { ok: true, value };
  if (!isRecord(value) || !Array.isArray(value.sources) || value.sources.length > MAX_SAFE_JSON_ARRAY_LENGTH) return { ok: false };
  const revision = boundedString(value.group_revision);
  if (!revision) return { ok: false };
  const connectorId = value.connector_id === undefined ? undefined : boundedString(value.connector_id);
  const connectorRevision = value.connector_revision === undefined ? undefined : boundedString(value.connector_revision);
  if ((value.connector_id !== undefined && !connectorId) || (value.connector_revision !== undefined && !connectorRevision)) return { ok: false };
  if (value.state === 'unavailable') {
    const reason = oneOf(UNAVAILABLE_REASONS, value.unavailable_reason);
    return reason && value.sources.length === 0
      ? { ok: true, value: { state: 'unavailable', unavailable_reason: reason, group_revision: revision, sources: [] } }
      : { ok: false };
  }
  // Records are proof only for the destination the applied revision froze, so an available answer must name it.
  if (value.state !== 'available' || value.unavailable_reason !== undefined || !connectorId || !connectorRevision) return { ok: false };
  const sources = value.sources.map(parseSource);
  return sources.some((source) => !source)
    ? { ok: false }
    : { ok: true, value: { state: 'available', group_revision: revision, connector_id: connectorId, connector_revision: connectorRevision, sources: sources as WriteGroupSourceEvidenceItem[] } };
}
