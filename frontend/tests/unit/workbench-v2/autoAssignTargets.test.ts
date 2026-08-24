import { describe, it, expect } from 'vitest';
import { autoAssignTargets } from '../../../src/features/datalink/workbench-v2/state/autoAssignTargets';
import type { Point, Mapping, DbTarget } from '../../../src/features/datalink/workbench-v2/state/types';

/**
 * @file autoAssignTargets.test.ts
 * @description 測試自動欄位分配演算法，覆蓋無衝突、保留現有設定、有衝突以及名稱匹配等場景。
 */

describe('autoAssignTargets', () => {
  const mockPoints: Point[] = Array.from({ length: 8 }, (_, i) => ({
    id: `p-${i + 1}`,
    device_id: 'd-1',
    rule_id: 'r-1',
    rule_name: 'Holding Registers',
    name: `Point_${i + 1}`,
    address: `${40001 + i}`,
    data_type: 'float64',
    function: 'holding_register',
    width: 2,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0
  }));

  const mockMappings: Record<string, Mapping> = {
    'p-1': { point_id: 'p-1', tag_key: 'line1.temp_in', display_name: 'Temp In', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-2': { point_id: 'p-2', tag_key: 'line1.temp_out', display_name: 'Temp Out', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-3': { point_id: 'p-3', tag_key: 'line1.pressure_main', display_name: 'Pressure Main', unit: 'kPa', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-4': { point_id: 'p-4', tag_key: 'line1.pressure_sub', display_name: 'Pressure Sub', unit: 'kPa', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-5': { point_id: 'p-5', tag_key: 'line1.flow', display_name: 'Flow', unit: 'LPM', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-6': { point_id: 'p-6', tag_key: 'line1.humidity', display_name: 'Humidity', unit: '%', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-7': { point_id: 'p-7', tag_key: 'line1.vibration', display_name: 'Vibration', unit: 'mm/s', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    'p-8': { point_id: 'p-8', tag_key: 'line1.motor', display_name: 'Motor', unit: 'rpm', target_type: 'float64', scale: 1, offset: 0, enabled: true }
  };

  const columns = [
    'temp_in_c',
    'temp_out_c',
    'pressure_main_kpa',
    'pressure_sub_kpa',
    'flow_lpm',
    'humidity_pct',
    'vibration_mms',
    'motor_rpm'
  ];

  it('8 points 與 8 columns 精確名稱匹配（如 temp_in 與 temp_in_c 透過後綴匹配）且無衝突', () => {
    const targets = autoAssignTargets(mockPoints, mockMappings, columns, {});
    expect(Object.keys(targets)).toHaveLength(8);

    expect(targets['p-1'].column_name).toBe('temp_in_c');
    expect(targets['p-2'].column_name).toBe('temp_out_c');
    expect(targets['p-3'].column_name).toBe('pressure_main_kpa');
    expect(targets['p-4'].column_name).toBe('pressure_sub_kpa');
    expect(targets['p-5'].column_name).toBe('flow_lpm');
    expect(targets['p-6'].column_name).toBe('humidity_pct');
    expect(targets['p-7'].column_name).toBe('vibration_mms');
    expect(targets['p-8'].column_name).toBe('motor_rpm');
  });

  it('現有的 targets 應被保留，其餘的依序分配', () => {
    const existing: Record<string, DbTarget> = {
      'p-1': { tag_id: 'tag.line1.temp_in', column_name: 'motor_rpm', enabled: true }
    };
    const targets = autoAssignTargets(mockPoints, mockMappings, columns, existing);
    
    // p-1 保留為 motor_rpm
    expect(targets['p-1'].column_name).toBe('motor_rpm');
    
    // 由於 motor_rpm 被 p-1 佔用，原本 p-8 (motor) 不能再選 motor_rpm，應 fallback 到其他欄位 (例如 index fallback / first unused)
    expect(targets['p-8'].column_name).not.toBe('motor_rpm');
    expect(targets['p-2'].column_name).toBe('temp_out_c');
  });

  it('當 columns 數量少於 points 時，應產生 index wrap 衝突', () => {
    const shortColumns = ['temp_in_c', 'temp_out_c', 'flow_lpm', 'motor_rpm'];
    const targets = autoAssignTargets(mockPoints, mockMappings, shortColumns, {});
    
    expect(Object.keys(targets)).toHaveLength(8);
    
    // 檢查是否有欄位被重複使用 (衝突)
    const usedColumns = Object.values(targets).map(t => t.column_name);
    const uniqueUsedColumns = new Set(usedColumns);
    
    expect(uniqueUsedColumns.size).toBeLessThan(8);
    expect(usedColumns).toContain('temp_in_c');
  });

  it('Cross-protocol Modbus Share and database target binding: matches MC/FATEK tag keys with DB columns', () => {
    const mcPoints: Point[] = [
      {
        id: 'p-mc-1',
        device_id: 'dev-mc-1',
        rule_id: 'rule-mc-1',
        rule_name: 'D Registers',
        name: 'SENSOR_D0',
        address: 'D0',
        data_type: 'int16',
        function: 'D (Word)',
        width: 1,
        enabled: true,
        skipped: false,
        _rule_scale: 1,
        _rule_offset: 0,
      },
      {
        id: 'p-mc-2',
        device_id: 'dev-mc-1',
        rule_id: 'rule-mc-1',
        rule_name: 'D Registers',
        name: 'SENSOR_D1',
        address: 'D1',
        data_type: 'int16',
        function: 'D (Word)',
        width: 1,
        enabled: true,
        skipped: false,
        _rule_scale: 1,
        _rule_offset: 0,
      },
    ];

    const mcMappings: Record<string, Mapping> = {
      'p-mc-1': { point_id: 'p-mc-1', tag_key: 'line1.sensor_d0', display_name: 'Sensor D0', unit: '', target_type: 'float64', scale: 1, offset: 0, enabled: true },
      'p-mc-2': { point_id: 'p-mc-2', tag_key: 'line1.sensor_d1', display_name: 'Sensor D1', unit: '', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    };

    const mcColumns = ['sensor_d0', 'sensor_d1', 'other_col'];
    const targets = autoAssignTargets(mcPoints, mcMappings, mcColumns, {});

    expect(targets['p-mc-1'].column_name).toBe('sensor_d0');
    expect(targets['p-mc-2'].column_name).toBe('sensor_d1');
  });

  it('matches tag keys having r prefix like rd0 with DB column sensor_d0 or d0', () => {
    const points: Point[] = [
      {
        id: 'p-1',
        device_id: 'dev-1',
        rule_id: 'rule-1',
        rule_name: 'D Registers',
        name: 'SENSOR_D0',
        address: 'D0',
        data_type: 'int16',
        function: 'D (Word)',
        width: 1,
        enabled: true,
        skipped: false,
        _rule_scale: 1,
        _rule_offset: 0,
      },
    ];

    const mappings: Record<string, Mapping> = {
      'p-1': { point_id: 'p-1', tag_key: 'dev.mc.01.sensor.rd0', display_name: 'Sensor D0', unit: '', target_type: 'float64', scale: 1, offset: 0, enabled: true },
    };

    const columns = ['sensor_d0', 'ambient_temp'];
    const targets = autoAssignTargets(points, mappings, columns, {});

    expect(targets['p-1'].column_name).toBe('sensor_d0');
  });
});
