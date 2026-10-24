import { describe, expect, it } from 'vitest';
import { comparePreciseInstants, formatPreciseUtc, parsePreciseInstant } from '@/utils/preciseInstant';

const at = (text: string) => parsePreciseInstant(text) as NonNullable<ReturnType<typeof parsePreciseInstant>>;

describe('Precise RFC 3339 instants', () => {
  it('keeps microseconds and nanoseconds that Date.parse would drop', () => {
    expect(comparePreciseInstants(at('2026-10-24T10:00:00.000001Z'), at('2026-10-24T10:00:00Z'))).toBeGreaterThan(0);
    expect(comparePreciseInstants(at('2026-10-24T10:00:00.000000001Z'), at('2026-10-24T10:00:00.000000002Z'))).toBeLessThan(0);
    expect(comparePreciseInstants(at('2026-10-24T10:00:00.5Z'), at('2026-10-24T10:00:00.500000000Z'))).toBe(0);
  });

  it('orders instants written with different offsets by the moment they name', () => {
    expect(comparePreciseInstants(at('2026-10-24T18:00:00+08:00'), at('2026-10-24T10:00:00Z'))).toBe(0);
    expect(comparePreciseInstants(at('2026-10-24T18:00:00.000001+08:00'), at('2026-10-24T10:00:00Z'))).toBeGreaterThan(0);
  });

  it('rejects anything that is not a complete RFC 3339 instant', () => {
    for (const text of ['yesterday', '2026-10-24 10:00:00Z', '2026-10-24T10:00:00', '2026-10-24T10:00:00.1234567890Z', '2026-13-40T10:00:00Z', '']) {
      expect(parsePreciseInstant(text), text).toBeNull();
    }
    expect(parsePreciseInstant(undefined)).toBeNull();
  });

  it('formats UTC seconds plainly and shows a fraction only when there is one', () => {
    expect(formatPreciseUtc('2026-10-24T10:00:15Z')).toBe('2026-10-24 10:00:15');
    expect(formatPreciseUtc('2026-10-24T10:00:00.000001Z')).toBe('2026-10-24 10:00:00.000001');
    expect(formatPreciseUtc('2026-10-24T10:00:00.500000000Z')).toBe('2026-10-24 10:00:00.5');
    expect(formatPreciseUtc('2026-10-24T18:00:00+08:00')).toBe('2026-10-24 10:00:00');
    expect(formatPreciseUtc('nonsense')).toBe('');
  });
});
