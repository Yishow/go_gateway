import { describe, expect, it } from 'vitest';
import { buildBatchNamePreview } from '../batchNaming';

describe('buildBatchNamePreview', () => {
  it('creates sequential names with zero padding', () => {
    const result = buildBatchNamePreview('LINEA', 3, []);
    expect(result.map((item) => item.name)).toEqual(['LINEA-001', 'LINEA-002', 'LINEA-003']);
  });

  it('marks conflicts against existing names (case-insensitive)', () => {
    const result = buildBatchNamePreview('linea', 3, ['LINEA-002']);
    expect(result[0].conflict).toBe(false);
    expect(result[1].conflict).toBe(true);
    expect(result[2].conflict).toBe(false);
  });

  it('falls back to SRC prefix for blank input', () => {
    const result = buildBatchNamePreview('   ', 2, []);
    expect(result.map((item) => item.name)).toEqual(['SRC-001', 'SRC-002']);
  });
});
