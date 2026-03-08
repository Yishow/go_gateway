import { describe, it, expect } from 'vitest';
import { toNumericValue } from '@/components/monitoring/valueParser';

describe('toNumericValue', () => {
  it('returns numeric values as-is', () => {
    expect(toNumericValue(0)).toBe(0);
    expect(toNumericValue(42.5)).toBe(42.5);
  });

  it('parses numeric strings including zero', () => {
    expect(toNumericValue('0')).toBe(0);
    expect(toNumericValue('  12.75 ')).toBe(12.75);
  });

  it('returns null for invalid values', () => {
    expect(toNumericValue('')).toBeNull();
    expect(toNumericValue('not-a-number')).toBeNull();
    expect(toNumericValue(undefined)).toBeNull();
    expect(toNumericValue(null)).toBeNull();
  });
});
