import { describe, expect, it } from 'vitest';
import { parseWriteGroupExactRecordingSource } from '../../../src/utils/studioV2ExactRecordingSourceJson';

const proof = {
  measurement_id: 'measurement-1',
  definition_revision: 'definition-rev-1',
  source_binding_revision: 'binding-rev-1',
  series_epoch: 'epoch-1',
  original_data_type: 'uint64',
  scale: '0.1',
  unit: 'kWh',
  semantic_kind: 'counter',
  counter_policy: { allow_negative: false },
};

function parse(patch: Record<string, unknown> = {}) {
  return parseWriteGroupExactRecordingSource({ ...proof, ...patch });
}

describe('backend exact-recording-source decimal contract', () => {
  it.each(['+0.1', '00.1', '.1', '1.'])('accepts and preserves positive backend decimal spelling %s', (scale) => {
    expect(parse({ scale })?.scale).toBe(scale);
  });

  it('accepts exactly 64 scale characters and rejects longer values without numeric conversion', () => {
    const limit = '1'.padEnd(64, '0');
    const tooLong = '1'.padEnd(65, '0');
    expect(parse({ scale: limit })?.scale).toBe(limit);
    expect(parse({ scale: tooLong })).toBeNull();
  });

  it.each(['0', '00.000', '+0.0', '1e-1', ' 0.1', '0.1 ', '-0.1', '', '.'])('rejects non-positive or non-plain scale spelling %j', (scale) => {
    expect(parse({ scale })).toBeNull();
  });

  it('allows a missing policy only for a non-negative semantic counter', () => {
    const withoutPolicy: Record<string, unknown> = { ...proof };
    delete withoutPolicy.counter_policy;
    expect(parseWriteGroupExactRecordingSource(withoutPolicy)?.counter_policy).toBeUndefined();
    expect(parse({ semantic_kind: 'counter', counter_policy: { allow_negative: true } })).toBeNull();
    const signedWithoutPolicy = { ...withoutPolicy, semantic_kind: 'signed_counter' };
    expect(parseWriteGroupExactRecordingSource(signedWithoutPolicy)).toBeNull();
  });

  it('requires signed counters to allow negatives and forbids rollover configuration', () => {
    const signed = {
      semantic_kind: 'signed_counter',
      original_data_type: 'int64',
      counter_policy: { allow_negative: true },
    };
    expect(parse(signed)?.semantic_kind).toBe('signed_counter');
    expect(parse({ ...signed, counter_policy: { allow_negative: false } })).toBeNull();
    expect(parse({ ...signed, counter_policy: { allow_negative: true, rollover_modulus: '10' } })).toBeNull();
    expect(parse({ ...signed, counter_policy: { allow_negative: true, max_rate_per_second: '1' } })).toBeNull();
    expect(parse({ ...signed, counter_policy: { allow_negative: true, max_roll_over: 10 } })?.counter_policy)
      .toMatchObject({ allow_negative: true, max_roll_over: 10 });
  });

  it('validates and preserves positive rollover and rate decimals and requires a modulus for a rate', () => {
    const modulus = '+00.1';
    const rate = '1.';
    const parsed = parse({ counter_policy: { allow_negative: false, rollover_modulus: modulus, max_rate_per_second: rate } });
    expect(parsed?.counter_policy).toMatchObject({ rollover_modulus: modulus, max_rate_per_second: rate });
    expect(parse({ counter_policy: { allow_negative: false, max_rate_per_second: rate } })).toBeNull();
    expect(parse({ counter_policy: { allow_negative: false, rollover_modulus: ' 10' } })).toBeNull();
    expect(parse({ counter_policy: { allow_negative: false, rollover_modulus: '1e2' } })).toBeNull();
    expect(parse({ counter_policy: { allow_negative: false, rollover_modulus: '0' } })).toBeNull();
  });

  it('applies the same 64-character boundary to rollover and rate values', () => {
    const limit = '1'.padEnd(64, '0');
    const tooLong = '1'.padEnd(65, '0');
    expect(parse({ counter_policy: { allow_negative: false, rollover_modulus: limit, max_rate_per_second: limit } })?.counter_policy)
      .toMatchObject({ rollover_modulus: limit, max_rate_per_second: limit });
    expect(parse({ counter_policy: { allow_negative: false, rollover_modulus: tooLong } })).toBeNull();
    expect(parse({ counter_policy: { allow_negative: false, rollover_modulus: limit, max_rate_per_second: tooLong } })).toBeNull();
  });

  it('continues to reject an empty unit', () => {
    expect(parse({ unit: '' })).toBeNull();
  });
});
