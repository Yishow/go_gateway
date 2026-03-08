import { describe, expect, it } from 'vitest';
import {
  getSpanByDataType,
  resolveConflictSeverity,
  validateTypedOccupancyPlan,
} from '@/features/datalink/typedOccupancy';

describe('typedOccupancy', () => {
  it('returns expected span by data type', () => {
    expect(getSpanByDataType('int16')).toBe(1);
    expect(getSpanByDataType('float32')).toBe(2);
    expect(getSpanByDataType('int64')).toBe(4);
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

  it('validates typed plan count and total cell calculation', () => {
    const valid = validateTypedOccupancyPlan('float32', 10);
    expect(valid.valid).toBe(true);
    expect(valid.totalCells).toBe(20);
    expect(valid.errors).toEqual([]);
  });

  it('validates int64 as 4 contiguous cells per source', () => {
    const valid = validateTypedOccupancyPlan('int64', 5);
    expect(valid.valid).toBe(true);
    expect(valid.span).toBe(4);
    expect(valid.totalCells).toBe(20);
  });

  it('rejects count out of supported range', () => {
    const invalid = validateTypedOccupancyPlan('int16', 0);
    expect(invalid.valid).toBe(false);
    expect(invalid.errors).toContain('invalid_count');
  });
});
