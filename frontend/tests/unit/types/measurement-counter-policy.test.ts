import { describe, expect, it } from 'vitest';
import type { CounterPolicy, MeasurementDefinition } from '../../../src/types/measurement';

// A saved definition travels through JSON unchanged: the exact rollover evidence
// stays text (digits above 2^53 survive) and the legacy float field is untouched.
describe('CounterPolicy exact rollover evidence', () => {
  it('keeps exact digits and the legacy float side by side through a JSON round trip', () => {
    const policy: CounterPolicy = {
      allow_negative: false,
      max_roll_over: 65536,
      rollover_modulus: '9007199254740993',
      max_rate_per_second: '0.125',
    };
    const roundTripped = JSON.parse(JSON.stringify(policy)) as CounterPolicy;
    expect(roundTripped.rollover_modulus).toBe('9007199254740993');
    expect(roundTripped.max_rate_per_second).toBe('0.125');
    expect(roundTripped.max_roll_over).toBe(65536);
  });

  it('is optional, so a definition saved before the evidence existed still fits the type', () => {
    const legacy: Pick<MeasurementDefinition, 'counter_policy'> = {
      counter_policy: JSON.parse('{"allow_negative":false,"max_roll_over":65536}') as CounterPolicy,
    };
    expect(legacy.counter_policy?.rollover_modulus).toBeUndefined();
    expect(legacy.counter_policy?.max_roll_over).toBe(65536);
  });
});
