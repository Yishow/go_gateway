import { describe, expect, it } from 'vitest';
import type { WriteGroupSourceRecord } from '@/types/studioV2WriteGroupDelivery';
import { parseWriteGroupDeliveryData } from '@/utils/studioV2WriteGroupDeliveryJson';

const stages = { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 };
const record = (over: Record<string, unknown> = {}): Record<string, unknown> => ({
  available: true, group_revision: 'applied-1', connector_revision: 'crev-1', record_id: 'rec-1', effect_key: 'eff-1',
  payload_digest: 'b'.repeat(64), profile_kind: 'cumulative', window_start: '2026-10-24T10:00:00Z', window_end: '2026-10-24T10:10:00Z',
  status: 'baseline', reason: 'first-reading', raw_present: true, delta_known: false, delivery: 'sql_committed',
  committed_at: '2026-10-24T10:10:03Z', ...over,
});
const evidence = (latest: unknown, over: Record<string, unknown> = {}, confirmed: unknown = latest) => ({
  state: 'available', group_revision: 'applied-1', connector_id: 'conn-1', connector_revision: 'crev-1',
  sources: [{ circuit_id: 'cir-a', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest, last_sql_committed: confirmed }], ...over,
});
const payload = (source_evidence?: unknown) => ({
  group_id: 'grp-1', intake: { state: 'active' }, stages, last_sql_committed_at: '2026-10-24T10:10:03Z', oldest_pending_seconds: 0,
  no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  ...(source_evidence === undefined ? {} : { source_evidence }),
});

describe('The delivery parser carries per source evidence without inventing success', () => {
  it('keeps every verified field of a baseline record, including an unknown difference', () => {
    const parsed = parseWriteGroupDeliveryData(payload(evidence(record())));
    const kept = parsed?.source_evidence?.sources[0].latest as WriteGroupSourceRecord;
    expect(parsed?.source_evidence?.state).toBe('available');
    expect(kept).toMatchObject({ available: true, status: 'baseline', reason: 'first-reading', raw_present: true, delta_known: false, delivery: 'sql_committed' });
    expect(kept.window_start).toBe('2026-10-24T10:00:00Z');
  });

  it('does not fill evidence into a legacy answer that has none', () => {
    expect(parseWriteGroupDeliveryData(payload())?.source_evidence).toBeUndefined();
    expect(parseWriteGroupDeliveryData(payload(null))?.source_evidence).toBeNull();
  });

  it('keeps an unverifiable record as unavailable and never as available', () => {
    const parsed = parseWriteGroupDeliveryData(payload(evidence({ available: false, cause: 'digest_mismatch', delivery: 'unknown' }, {}, null)));
    expect(parsed?.source_evidence?.sources[0].latest).toEqual({ available: false, cause: 'digest_mismatch', delivery: 'unknown' });
  });

  it('never accepts an unverifiable or not yet committed record as the SQL confirmed one', () => {
    expect(parseWriteGroupDeliveryData(payload(evidence(null, {}, { available: false, cause: 'digest_mismatch', delivery: 'sql_committed' })))).toBeNull();
    expect(parseWriteGroupDeliveryData(payload(evidence(null, {}, record({ delivery: 'queued', committed_at: undefined }))))).toBeNull();
  });

  it('replaces a reason that is not a stable code with reason_unavailable', () => {
    const parsed = parseWriteGroupDeliveryData(payload(evidence(record({ reason: 'driver: password=hunter2 failed' }))));
    expect((parsed?.source_evidence?.sources[0].latest as WriteGroupSourceRecord).reason).toBe('reason_unavailable');
    expect(JSON.stringify(parsed)).not.toContain('hunter2');
  });

  it.each([
    ['an unknown status', record({ status: 'great' })],
    ['an unknown delivery stage', record({ delivery: 'done' })],
    ['a SQL committed record without a receipt time', record({ committed_at: undefined })],
    ['a receipt time on a record that is not SQL committed', record({ delivery: 'queued' })],
    ['a NULL status that claims a stored reading', record({ status: 'missing', raw_present: true })],
    ['a stored status that claims no reading', record({ status: 'ok', raw_present: false })],
    ['a baseline that claims a known difference', record({ delta_known: true })],
    ['a malformed digest', record({ payload_digest: 'xyz' })],
    ['a non UTC-parsable window', record({ window_start: 'yesterday' })],
    ['a window that ends before it starts', record({ window_end: '2026-10-24T09:00:00Z' })],
    ['a missing raw flag', record({ raw_present: undefined })],
  ])('rejects the whole answer for %s', (_name, bad) => {
    expect(parseWriteGroupDeliveryData(payload(evidence(bad)))).toBeNull();
  });

  it('rejects contradictory evidence state and too many sources', () => {
    expect(parseWriteGroupDeliveryData(payload(evidence(record(), { state: 'unavailable' })))).toBeNull();
    expect(parseWriteGroupDeliveryData(payload({ state: 'unavailable', unavailable_reason: 'revision_mismatch', group_revision: 'x', sources: [] })))
      .toMatchObject({ source_evidence: { state: 'unavailable', unavailable_reason: 'revision_mismatch', sources: [] } });
    expect(parseWriteGroupDeliveryData(payload({ state: 'available', group_revision: 'x', sources: Array.from({ length: 300 }, () => ({})) }))).toBeNull();
  });
});

describe('The parser keeps the frozen applied destination and compares windows exactly', () => {
  const unavailable = { state: 'unavailable', unavailable_reason: 'revision_mismatch', group_revision: 'x', sources: [] };

  it('requires an available answer to name the destination the applied revision froze', () => {
    const { connector_id: _id, ...withoutId } = evidence(record());
    const { connector_revision: _rev, ...withoutRevision } = evidence(record());
    expect(parseWriteGroupDeliveryData(payload(withoutId))).toBeNull();
    expect(parseWriteGroupDeliveryData(payload(withoutRevision))).toBeNull();
    expect(parseWriteGroupDeliveryData(payload(evidence(record(), { connector_revision: '' })))).toBeNull();
    expect(parseWriteGroupDeliveryData(payload(evidence(record())))?.source_evidence).toMatchObject({ connector_id: 'conn-1', connector_revision: 'crev-1' });
  });

  it('accepts an unavailable answer without a destination and still rejects a malformed one', () => {
    expect(parseWriteGroupDeliveryData(payload(unavailable))?.source_evidence).toMatchObject({ state: 'unavailable' });
    expect(parseWriteGroupDeliveryData(payload({ ...unavailable, connector_id: 7 }))).toBeNull();
  });

  it.each([
    ['a one microsecond window', '2026-10-24T10:00:00Z', '2026-10-24T10:00:00.000001Z', true],
    ['a one nanosecond window', '2026-10-24T10:00:00Z', '2026-10-24T10:00:00.000000001Z', true],
    ['an empty window', '2026-10-24T10:00:00.000001Z', '2026-10-24T10:00:00.000001Z', false],
    ['a window that ends a microsecond early', '2026-10-24T10:00:00.000002Z', '2026-10-24T10:00:00.000001Z', false],
    ['offsets naming the same moment', '2026-10-24T18:00:00+08:00', '2026-10-24T10:00:00Z', false],
    ['an offset end one microsecond later', '2026-10-24T10:00:00Z', '2026-10-24T18:00:00.000001+08:00', true],
  ])('treats %s exactly', (_name, start, end, accepted) => {
    const parsed = parseWriteGroupDeliveryData(payload(evidence(record({ window_start: start, window_end: end, committed_at: '2026-10-24T10:10:03Z' }))));
    expect(parsed !== null).toBe(accepted);
    if (accepted) expect((parsed?.source_evidence?.sources[0].latest as WriteGroupSourceRecord).window_end).toBe(end);
  });
});
