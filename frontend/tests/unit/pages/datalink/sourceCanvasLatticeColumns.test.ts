import { describe, expect, it } from 'vitest';
import {
  clampLatticeColumns,
  LATTICE_COLUMNS_DEFAULT,
  LATTICE_COLUMNS_MAX,
  LATTICE_COLUMNS_MIN,
} from '../../../../src/pages/datalink/workbench/sourceCanvasLatticeColumns';

describe('sourceCanvasLatticeColumns', () => {
  it('clamps lattice column count to the allowed range', () => {
    expect(clampLatticeColumns(NaN)).toBe(LATTICE_COLUMNS_DEFAULT);
    expect(clampLatticeColumns(3)).toBe(LATTICE_COLUMNS_MIN);
    expect(clampLatticeColumns(10)).toBe(10);
    expect(clampLatticeColumns(99)).toBe(LATTICE_COLUMNS_MAX);
  });
});
