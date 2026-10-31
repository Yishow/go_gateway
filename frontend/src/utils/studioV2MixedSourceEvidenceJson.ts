import type { WriteGroupMixedSourceRecord, WriteGroupSourceAggregation, WriteGroupSourceCalculation, WriteGroupSourceDelivery, WriteGroupSourceStatus } from '../types/studioV2WriteGroupDelivery';
import { boundedString } from './safeJson';
import { parsePreciseInstant } from './preciseInstant';

type JsonRecord = Record<string, unknown>;
type AggregationParser = (value: unknown, status: WriteGroupSourceStatus) => WriteGroupSourceAggregation | null;
const isRecord = (v: unknown): v is JsonRecord => v !== null && typeof v === 'object' && !Array.isArray(v);
const DIGEST = /^[a-f0-9]{64}$/;
const REASON = /^[a-z0-9_-]{1,64}$/;
const ZERO_TIME = '0001-01-01T00:00:00Z';
const NANOSECOND = 1_000_000_000n, MAX_NANOS = (1n << 63n) - 1n;
const CALC_KEYS = new Set(['result_shape', 'method', 'interval_seconds', 'alignment', 'due_at', 'next_due', 'outcome_digest', 'outcome_kind', 'status', 'reason', 'result_present', 'row_emitted', 'raw_present', 'source_type', 'result_type', 'aligned_start', 'aligned_end', 'actual_start', 'actual_end', 'rounded']);
const RECORD_KEYS = new Set(['available', 'group_revision', 'connector_revision', 'record_id', 'effect_key', 'payload_digest', 'profile_kind', 'window_start', 'window_end', 'status', 'reason', 'raw_present', 'aggregation', 'delta_known', 'delivery', 'committed_at', 'calculation']);
const MEAN_KEYS = new Set(['method', 'version', 'basis', 'source_type', 'result_type', 'valid_count', 'rejected_count', 'coverage', 'reason']);
const DELIVERIES = new Set(['queued', 'retrying', 'blocked', 'quarantined', 'unknown', 'sql_committed', 'operator_skipped', 'window_skipped']);
const STATUSES = new Set(['ok', 'partial', 'baseline', 'missing', 'bad', 'stale', 'invalid', 'reset_or_uncertain', 'gap', 'not_due']);
const MEMBER_STATUSES = new Set(['ok', 'missing', 'bad', 'stale', 'invalid', 'not_due']);
const SOURCE_RESULT: Record<string, string> = { int16: 'int64', int32: 'int64', int64: 'int64', uint16: 'uint64', uint32: 'uint64', uint64: 'uint64', float32: 'float64', float64: 'float64', bool: 'bool', string: 'text' };
const keysAllowed = (v: JsonRecord, keys: Set<string>) => Object.keys(v).every(key => keys.has(key));
function nanos(value: unknown): bigint | null {
  const parsed = parsePreciseInstant(value);
  if (!parsed || typeof value !== 'string' || !value.endsWith('Z') || new Date(parsed.seconds * 1000).toISOString().slice(0, 19) !== value.slice(0, 19)) return null;
  const n = BigInt(parsed.seconds) * NANOSECOND + BigInt(parsed.nanos);
  return n >= -MAX_NANOS - 1n && n <= MAX_NANOS ? n : null;
}
const reasonValid = (v: unknown) => v === undefined || (typeof v === 'string' && REASON.test(v));

function calculation(value: unknown): WriteGroupSourceCalculation | null {
  if (!isRecord(value) || !keysAllowed(value, CALC_KEYS)) return null;
  const c = value;
  const method = c.method, interval = c.interval_seconds;
  if (c.result_shape !== 'member-calculation-v1' || c.alignment !== 'utc_epoch' ||
    (method !== 'latest' && method !== 'mean' && method !== 'delta') || typeof interval !== 'number' || !Number.isSafeInteger(interval) || interval < 1 ||
    BigInt(interval) > MAX_NANOS / NANOSECOND || typeof c.outcome_digest !== 'string' || !DIGEST.test(c.outcome_digest) ||
    (c.outcome_kind !== 'row' && c.outcome_kind !== 'skipped' && c.outcome_kind !== 'no_data') ||
    typeof c.status !== 'string' || !MEMBER_STATUSES.has(c.status) || !reasonValid(c.reason) ||
    typeof c.result_present !== 'boolean' || typeof c.row_emitted !== 'boolean' || typeof c.raw_present !== 'boolean' ||
    (c.rounded !== undefined && typeof c.rounded !== 'boolean')) return null;
  if ((c.status === 'not_due' && c.reason !== 'not_due') || (c.status !== 'ok' && !boundedString(c.reason))) return null;
  const due = nanos(c.due_at);
  if (due === null || c.row_emitted !== (c.outcome_kind === 'row') || c.result_present !== (c.status === 'ok') ||
    (c.outcome_kind === 'no_data' && c.result_present) || (c.result_present && !c.raw_present)) return null;
  if (typeof c.source_type !== 'string' || !Object.hasOwn(SOURCE_RESULT, c.source_type)) return null;
  const source = SOURCE_RESULT[c.source_type];
  const expected = method === 'latest' ? source : source === 'float64' ? source : ['int64', 'uint64'].includes(source) ? 'decimal' : null;
  if (!expected || c.result_type !== expected || (method !== 'mean' && c.rounded !== undefined)) return null;
  const duration = BigInt(interval) * NANOSECOND;
  if (c.next_due !== undefined) {
    const next = nanos(c.next_due);
    if (next === null || next <= due || next % duration !== 0n) return null;
  }
  if (c.status === 'not_due') {
    if (c.result_present || c.raw_present || ['aligned_start', 'aligned_end', 'actual_start', 'actual_end', 'rounded'].some(k => c[k] !== undefined)) return null;
  } else {
    const alignedStart = nanos(c.aligned_start), alignedEnd = nanos(c.aligned_end), actualStart = nanos(c.actual_start), actualEnd = nanos(c.actual_end);
    if (alignedStart === null || alignedEnd === null || actualStart === null || actualEnd === null || alignedStart % duration !== 0n ||
      alignedEnd - alignedStart !== duration || alignedEnd !== due || actualStart < alignedStart || actualEnd > alignedEnd || actualStart >= actualEnd) return null;
  }
  // All fields are scalar and validated; copy rather than retaining caller objects.
  return { ...c } as unknown as WriteGroupSourceCalculation;
}

/** Mixed-only wire decoding. The legacy parser remains responsible for legacy records. */
export function parseMixedSourceRecord(v: JsonRecord, parseAggregation: AggregationParser): WriteGroupMixedSourceRecord | null {
  if (!keysAllowed(v, RECORD_KEYS) || v.available !== true || v.profile_kind !== 'mixed') return null;
  const c = calculation(v.calculation);
  const group = boundedString(v.group_revision), connector = boundedString(v.connector_revision), id = boundedString(v.record_id);
  if (!c || !group || !connector || !id || typeof v.delivery !== 'string' || !DELIVERIES.has(v.delivery) ||
    typeof v.status !== 'string' || !STATUSES.has(v.status) || typeof v.raw_present !== 'boolean' || v.raw_present !== c.raw_present || !reasonValid(v.reason)) return null;
  const committed = v.delivery === 'sql_committed';
  if (committed !== (v.committed_at !== undefined) || (committed && nanos(v.committed_at) === null)) return null;
  if (c.row_emitted) {
    if (!boundedString(v.effect_key) || typeof v.payload_digest !== 'string' || !DIGEST.test(v.payload_digest) || v.delivery === 'window_skipped') return null;
  } else if (v.effect_key !== '' || v.payload_digest !== '' || v.delivery !== 'window_skipped' || committed) return null;
  let start: string | null = null, end: string | null = null;
  if (c.status === 'not_due') {
    if (v.status !== 'not_due' || v.window_start !== ZERO_TIME || v.window_end !== ZERO_TIME || v.aggregation !== undefined) return null;
  } else {
    if (nanos(v.window_start) !== nanos(c.actual_start) || nanos(v.window_end) !== nanos(c.actual_end)) return null;
    start = v.window_start as string; end = v.window_end as string;
  }
  if (c.method !== 'delta' && v.status !== c.status) return null;
  if (c.method === 'delta' ? typeof v.delta_known !== 'boolean' || v.delta_known !== c.result_present : v.delta_known !== undefined) return null;
  if (c.method === 'delta' && c.result_present && !['ok', 'partial'].includes(v.status)) return null;
  let aggregation: WriteGroupSourceAggregation | undefined;
  if (c.method === 'mean' && c.status !== 'not_due') {
    if (!isRecord(v.aggregation) || !keysAllowed(v.aggregation, MEAN_KEYS)) return null;
    const parsed = parseAggregation(v.aggregation, v.status as WriteGroupSourceStatus);
    if (!parsed || parsed.source_type !== c.source_type || parsed.result_type !== c.result_type || (c.result_present && c.rounded === undefined)) return null;
    aggregation = parsed;
  } else if (v.aggregation !== undefined) return null;
  return {
    available: true, group_revision: group, connector_revision: connector, record_id: id, effect_key: v.effect_key as string,
    payload_digest: v.payload_digest as string, profile_kind: 'mixed', window_start: start, window_end: end,
    status: v.status as WriteGroupMixedSourceRecord['status'], raw_present: v.raw_present, calculation: c,
    delivery: v.delivery as WriteGroupSourceDelivery, ...(v.reason === undefined ? {} : { reason: v.reason as string }),
    ...(aggregation ? { aggregation } : {}), ...(c.method === 'delta' ? { delta_known: v.delta_known as boolean } : {}),
    ...(committed ? { committed_at: v.committed_at as string } : {}),
  };
}
