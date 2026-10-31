import { describe, expect, it } from 'vitest';
import { parseWriteGroupSourceEvidence } from '../../../src/utils/studioV2WriteGroupSourceEvidenceJson';

const digest = 'a'.repeat(64);
const start = '2026-10-31T00:00:00Z', due = '2026-10-31T00:00:15Z';
function record() {
  return { available: true, group_revision: 'revision', connector_revision: 'connector-rev', record_id: 'record', effect_key: 'effect', payload_digest: digest,
    profile_kind: 'mixed', window_start: start, window_end: due, status: 'ok', raw_present: true, delivery: 'queued',
    calculation: { result_shape: 'member-calculation-v1', method: 'latest', interval_seconds: 15, alignment: 'utc_epoch', due_at: due,
      outcome_digest: digest, outcome_kind: 'row', status: 'ok', result_present: true, row_emitted: true, raw_present: true,
      source_type: 'int64', result_type: 'int64', aligned_start: start, aligned_end: due, actual_start: start, actual_end: due } };
}
function envelope(latest: unknown, committed: unknown = null) {
  return { state: 'available', group_revision: 'revision', connector_id: 'connector', connector_revision: 'connector-rev',
    sources: [{ device_id: 'device', point_id: 'point', tag_id: 'tag', latest, last_sql_committed: committed }] };
}
function notDue() {
  const r: Record<string, unknown> = record();
  r.window_start = r.window_end = '0001-01-01T00:00:00Z'; r.status = 'not_due'; r.raw_present = false;
  const c = { ...record().calculation } as Record<string, unknown>;
  c.status = 'not_due'; c.reason = 'not_due'; c.result_present = c.raw_present = false; c.next_due = '2026-10-31T00:01:00Z'; c.interval_seconds = 60;
  for (const key of ['aligned_start', 'aligned_end', 'actual_start', 'actual_end']) delete c[key];
  r.calculation = c; return r;
}
describe('mixed source evidence wire contract', () => {
  it('keeps selected result separate from raw presence and row emission', () => {
    const r = record(); r.delivery = 'window_skipped'; r.effect_key = r.payload_digest = ''; r.calculation.outcome_kind = 'skipped'; r.calculation.row_emitted = false;
    const parsed = parseWriteGroupSourceEvidence(envelope(r)); expect(parsed.ok).toBe(true);
    if (parsed.ok) expect(parsed.value?.sources[0].latest).toMatchObject({ effect_key: '', calculation: { result_present: true, raw_present: true, row_emitted: false } });
  });
  it('accepts a receipt-confirmed not-due NULL slot without inventing a data window', () => {
    const r = notDue(); r.delivery = 'sql_committed'; r.committed_at = '2026-10-31T00:00:16Z';
    const parsed = parseWriteGroupSourceEvidence(envelope(r, r)); expect(parsed.ok).toBe(true);
    if (parsed.ok) expect(parsed.value?.sources[0].last_sql_committed).toMatchObject({ window_start: null, window_end: null, calculation: { status: 'not_due', result_present: false, raw_present: false, row_emitted: true } });
  });
  it('keeps retention absence as absence', () => { expect(parseWriteGroupSourceEvidence(envelope(null))).toMatchObject({ ok: true, value: { sources: [{ latest: null, last_sql_committed: null }] } }); });
  it.each(['method', 'alignment', 'result_shape', 'source_type', 'result_type', 'status'])('rejects unknown calculation %s', (key) => {
    const r = record(); (r.calculation as Record<string, unknown>)[key] = 'future'; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  });
  it.each(['raw_value', 'sum', 'future'])('rejects unadvertised calculation field %s', (key) => {
    const r = record(); (r.calculation as Record<string, unknown>)[key] = 7; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  });
  it('rejects null or omitted typed calculation and contradictory skip receipt', () => {
    for (const calculation of [null, undefined]) expect(parseWriteGroupSourceEvidence(envelope({ ...record(), calculation })).ok).toBe(false);
    const r = record(); r.delivery = 'window_skipped'; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  });
  it('rejects not-due with a fabricated actual window or selected value', () => {
    for (const patch of [{ actual_start: start }, { result_present: true }, { raw_present: true }]) {
      const r = notDue(); r.calculation = { ...(r.calculation as object), ...patch }; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
    }
  });
});

describe('mixed method-specific evidence', () => {
  it('accepts mean rounding/coverage without confusing the raw observation with the selected result', () => {
    const r: Record<string, unknown> = record();
    r.calculation = { ...record().calculation, method: 'mean', result_type: 'decimal', rounded: true };
    r.aggregation = { method: 'arithmetic-mean', version: 'arithmetic-mean-v1', basis: 'exact', source_type: 'int64', result_type: 'decimal', valid_count: 3, rejected_count: 0, coverage: 'unknown' };
    expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(true);
    r.aggregation = { ...(r.aggregation as object), coverage: 'complete' };
    expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  });
  it('keeps a first counter baseline raw-present but result-unavailable and row-skipped', () => {
    const r: Record<string, unknown> = { ...record(), delivery: 'window_skipped', status: 'baseline', effect_key: '', payload_digest: '', delta_known: false };
    r.calculation = { ...record().calculation, method: 'delta', result_type: 'decimal', status: 'invalid', reason: 'counter-unavailable', outcome_kind: 'skipped', row_emitted: false, result_present: false };
    const parsed = parseWriteGroupSourceEvidence(envelope(r)); expect(parsed.ok).toBe(true);
    if (parsed.ok) expect(parsed.value?.sources[0].latest).toMatchObject({ delta_known: false, raw_present: true, calculation: { result_present: false } });
    r.delta_known = true; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  });
  it.each([0, 1.5, 9_223_372_037, Number.MAX_SAFE_INTEGER, Infinity])('rejects nonrepresentable cadence %s', interval_seconds => {
    const r = record(); r.calculation.interval_seconds = interval_seconds; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  });
  it.each(['actual_start', 'actual_end', 'due_at'])('rejects invalid/out-of-domain %s', key => {
    for (const value of ['2026-02-30T00:00:00Z', '2300-01-01T00:00:00Z', null]) {
      const r = record(); (r.calculation as Record<string, unknown>)[key] = value; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
    }
  });
  it('rejects unexpected mixed intent on a legacy record', () => {
    const r = { ...record(), profile_kind: 'instantaneous' }; expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  });
});

it('rejects mixed evidence from a foreign revision or destination and duplicate source identity', () => {
  for (const patch of [{ group_revision: 'other' }, { connector_revision: 'other' }]) expect(parseWriteGroupSourceEvidence(envelope({ ...record(), ...patch })).ok).toBe(false);
  const e = envelope(record()); e.sources.push({ ...e.sources[0] }); expect(parseWriteGroupSourceEvidence(e).ok).toBe(false);
});
it('accepts no-data closure without value/receipt and keeps a stable missing reason', () => {
  const r: Record<string, unknown> = { ...record(), delivery: 'window_skipped', status: 'missing', reason: 'no-samples', raw_present: false, effect_key: '', payload_digest: '' };
  r.calculation = { ...record().calculation, status: 'missing', reason: 'no-samples', outcome_kind: 'no_data', result_present: false, row_emitted: false, raw_present: false };
  expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(true);
});
it('never discards present calculation data to downgrade a malformed record into unavailable', () => {
  expect(parseWriteGroupSourceEvidence(envelope({ available: false, delivery: 'unknown', cause: 'calculation_projection_unsupported' })).ok).toBe(true);
  expect(parseWriteGroupSourceEvidence(envelope({ available: false, delivery: 'unknown', cause: 'payload_invalid', calculation: {} })).ok).toBe(false);
});
it.each([5, 15, 60, 600, 9_223_372_036])('accepts a representable independently aligned %s-second result window', seconds => {
  const r = record(), windowStart = '1970-01-01T00:00:00Z';
  const windowEnd = new Date(seconds * 1000).toISOString();
  r.window_start = windowStart; r.window_end = windowEnd;
  Object.assign(r.calculation, { interval_seconds: seconds, due_at: windowEnd, aligned_start: windowStart, aligned_end: windowEnd, actual_start: windowStart, actual_end: windowEnd });
  expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(true);
});
it('keeps an exhausted not-due member without a fabricated next window; present next_due must be future and aligned', () => {
  const ended = notDue(); delete (ended.calculation as Record<string, unknown>).next_due;
  expect(parseWriteGroupSourceEvidence(envelope(ended)).ok).toBe(true);
  for (const next_due of [due, start, '2026-10-31T00:00:59Z', null]) {
    const r = notDue(); (r.calculation as Record<string, unknown>).next_due = next_due;
    expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
  }
});
it('requires raw provenance for a usable selected result without treating raw presence as sufficient', () => {
  const r = record(); r.raw_present = r.calculation.raw_present = false;
  expect(parseWriteGroupSourceEvidence(envelope(r)).ok).toBe(false);
});
