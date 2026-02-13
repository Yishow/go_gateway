import { describe, expect, it } from 'vitest';
import { getSpanByDataType, resolveConflictSeverity } from '../typedOccupancy';

describe('typedOccupancy', () => {
  it('returns expected span by data type', () => {
    expect(getSpanByDataType('int16')).toBe(1);
    expect(getSpanByDataType('float32')).toBe(2);
    expect(getSpanByDataType('float64')).toBe(4);
  });

  it('resolves hard conflict when planned overlaps existing point', () => {
    expect(
      resolveConflictSeverity({
        hasUsedPoint: true,
        hasLinkedAddress: false,
      })
    ).toBe('hard');
  });

  it('resolves soft conflict when planned overlaps linked address only', () => {
    expect(
      resolveConflictSeverity({
        hasUsedPoint: false,
        hasLinkedAddress: true,
      })
    ).toBe('soft');
  });

  it('returns none when no overlap exists', () => {
    expect(
      resolveConflictSeverity({
        hasUsedPoint: false,
        hasLinkedAddress: false,
      })
    ).toBe('none');
  });
});
