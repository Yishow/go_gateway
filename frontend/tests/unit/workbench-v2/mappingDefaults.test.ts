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

    const mapping0 = buildDefaultMapping(point, 0);
    expect(mapping0.point_id).toBe('p-01');
    expect(mapping0.tag_key).toBe('dev.01.sensor.1.r40001');
    expect(mapping0.display_name).toBe('SENSOR_1');
    expect(mapping0.unit).toBe('');
    expect(mapping0.target_type).toBe('float64');
    expect(mapping0.scale).toBe(0.5);
    expect(mapping0.offset).toBe(10);
    expect(mapping0.enabled).toBe(true);
  });

  it('should build unique tag keys for different point identities instead of cycling 8 templates', () => {
    const pointA: Point = {
      id: 'p-01',
      device_id: 'dev-01',
      rule_id: 'rule-01',
      rule_name: 'Holding Registers',
      name: 'LINE_A_0',
      address: '40001',
      data_type: 'int16',
      function: 'holding_register',
      width: 1,
      enabled: true,
      skipped: false,
      _rule_scale: 0.1,
      _rule_offset: 0,
    };
    const pointB: Point = {
      ...pointA,
      id: 'p-09',
      rule_id: 'rule-02',
      name: 'LINE_B_0',
      address: '40101',
    };

    const mappingA = buildDefaultMapping(pointA, 0);
    const mappingB = buildDefaultMapping(pointB, 8);

    expect(mappingA.tag_key).toBe('dev.01.line.a.0.r40001');
    expect(mappingB.tag_key).toBe('dev.01.line.b.0.r40101');
    expect(mappingA.tag_key).not.toBe(mappingB.tag_key);
  });
});
