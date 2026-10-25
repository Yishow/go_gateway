import { describe, expect, it } from 'vitest';
import { parseWriteGroupSourceEvidence } from '@/utils/studioV2WriteGroupSourceEvidenceJson';

const aggregation = {
  method: 'arithmetic-mean', version: 'arithmetic-mean-v1', basis: 'exact', source_type: 'uint64', result_type: 'decimal',
  valid_count: 2, rejected_count: 0, coverage: 'unknown',
};

const record = (over: Record<string, unknown> = {}) => ({
  available: true, delivery: 'sql_committed', group_revision: 'rev-1', connector_revision: 'crev-1', record_id: 'r0', effect_key: 'e0',
  payload_digest: 'a'.repeat(64), profile_kind: 'instantaneous', window_start: '2026-01-01T00:00:00Z', window_end: '2026-01-01T00:00:15Z',
  status: 'ok', raw_present: true, committed_at: '2026-01-01T00:00:16Z', aggregation, ...over,
});

const evidence = (latest: unknown) => ({
  state: 'available', group_revision: 'rev-1', connector_id: 'conn-1', connector_revision: 'crev-1',
  sources: [{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest, last_sql_committed: null }],
});

describe('the delivery decoder and the mean evidence of a source', () => {
  it('keeps the verified aggregation of a record as it was sent', () => {
    const parsed = parseWriteGroupSourceEvidence(evidence(record()));

    expect(parsed.ok).toBe(true);
    const latest = parsed.ok && parsed.value && parsed.value.state === 'available' ? parsed.value.sources[0].latest : null;
    expect(latest && latest.available ? (latest as { aggregation?: unknown }).aggregation : undefined).toEqual(aggregation);
  });

  it('keeps a record without a mean exactly as before', () => {
    const parsed = parseWriteGroupSourceEvidence(evidence(record({ aggregation: undefined })));

    expect(parsed.ok).toBe(true);
    const latest = parsed.ok && parsed.value && parsed.value.state === 'available' ? parsed.value.sources[0].latest : null;
    expect(latest && latest.available ? 'aggregation' in latest : true).toBe(false);
  });

  it.each([
    ['a coverage of complete, which no stored proof can back', { coverage: 'complete' }],
    ['a coverage that is not one of the three observed', { coverage: 'bogus' }],
    ['a negative valid count', { valid_count: -1 }],
    ['a fractional rejected count', { rejected_count: 0.5 }],
    ['a method that is not the arithmetic mean', { method: 'median' }],
    ['a basis that is neither exact nor float64', { basis: 'approximate' }],
    ['a missing result type', { result_type: '' }],
  ])('rejects the whole answer for %s', (_name, change) => {
    expect(parseWriteGroupSourceEvidence(evidence(record({ aggregation: { ...aggregation, ...change } }))).ok).toBe(false);
  });

  it.each([
    ['an exact basis stored as a float64 result', { basis: 'exact', result_type: 'float64' }],
    ['a float64 basis stored as a decimal result', { basis: 'float64', result_type: 'decimal' }],
    ['a version this contract does not know', { version: 'arithmetic-mean-v9' }],
    ['a source that is not a number', { source_type: 'string' }],
    ['no valid observation but a coverage other than missing, on an ok record', { valid_count: 0 }],
    ['a missing coverage with valid observations', { coverage: 'missing' }],
  ])('rejects the whole answer for %s', (_name, change) => {
    expect(parseWriteGroupSourceEvidence(evidence(record({ aggregation: { ...aggregation, ...change } }))).ok).toBe(false);
  });

  it('rejects valid observations on a window that stored nothing', () => {
    expect(parseWriteGroupSourceEvidence(evidence(record({ status: 'missing', raw_present: false }))).ok).toBe(false);
  });

  it('accepts what the projection can really send: a silent window and a float source', () => {
    const silent = record({
      status: 'missing', raw_present: false, aggregation: { ...aggregation, valid_count: 0, coverage: 'missing' },
    });
    expect(parseWriteGroupSourceEvidence(evidence(silent)).ok).toBe(true);
    const floating = record({ aggregation: { ...aggregation, basis: 'float64', source_type: 'float32', result_type: 'float64', coverage: 'partial' } });
    expect(parseWriteGroupSourceEvidence(evidence(floating)).ok).toBe(true);
  });

  it('rejects an aggregation on a cumulative record, which has no mean', () => {
    expect(parseWriteGroupSourceEvidence(evidence(record({ profile_kind: 'cumulative', delta_known: true }))).ok).toBe(false);
  });
});
