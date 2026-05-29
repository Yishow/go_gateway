import { describe, it, expect } from 'vitest';
import { POINT_SEMANTIC, RAW_VALUE_SEEDS, buildDefaultMapping } from '../../../src/features/datalink/workbench-v2/state/mappingDefaults';
import type { Point } from '../../../src/features/datalink/workbench-v2/state/types';

describe('mappingDefaults', () => {
  it('should have 8 entries in POINT_SEMANTIC', () => {
    expect(POINT_SEMANTIC).toHaveLength(8);
  });

  it('should have 8 seeds in RAW_VALUE_SEEDS', () => {
    expect(RAW_VALUE_SEEDS).toHaveLength(8);
  });

  it('should build default mapping correctly based on index', () => {
    const point: Point = {
      id: 'p-01',
      device_id: 'dev-01',
      rule_id: 'rule-01',
      rule_name: 'Holding Registers',
      name: 'SENSOR_1',
      address: '40001',
      data_type: 'int16',
      function: 'holding_register',
      width: 1,
      enabled: true,
      skipped: false,
      _rule_scale: 0.5,
      _rule_offset: 10,
    };

    // Index 0 -> line01.temp.inlet
    const mapping0 = buildDefaultMapping(point, 0);
    expect(mapping0.point_id).toBe('p-01');
    expect(mapping0.tag_key).toBe('line01.temp.inlet');
    expect(mapping0.display_name).toBe('入口溫度');
    expect(mapping0.unit).toBe('°C');
    expect(mapping0.target_type).toBe('float64');
    expect(mapping0.scale).toBe(0.5);
    expect(mapping0.offset).toBe(10);
    expect(mapping0.enabled).toBe(true);

    // Index 9 -> (9 % 8 = 1) -> line01.temp.outlet
    const mapping9 = buildDefaultMapping(point, 9);
    expect(mapping9.tag_key).toBe('line01.temp.outlet');
    expect(mapping9.display_name).toBe('出口溫度');
  });
});
