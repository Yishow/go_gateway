import { describe, it, expect } from 'vitest';
import { renderHook } from '@testing-library/react';
import type { Rule, Point, Mapping } from '../../../src/features/datalink/workbench-v2/state/types';
import {
  useAllPoints,
  useShareLayout,
  useConflictAddrs,
  useRulePoints,
  useMappingValidation,
  useSelectedMapping,
} from '../../../src/features/datalink/workbench-v2/state/selectors';

const mockRules: Rule[] = [
  {
    id: 'rule-1',
    device_id: 'dev-1',
    name: 'Coils',
    start_address: '00001',
    count: 2,
    data_type: 'bool',
    naming_prefix: 'C_',
    enabled: true,
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
  },
];

describe('useAllPoints', () => {
  it('應衍生出正確的點位列表，且當引數未變時回傳同一個快取參考', () => {
    let rules = [...mockRules];
    const devId = 'dev-1';

    const { result, rerender } = renderHook(() => useAllPoints(rules, devId));
    const firstResult = result.current;

    expect(firstResult).toHaveLength(2);
    expect(firstResult[0].address).toBe('1');

    // 重新渲染，引數 reference 未改變
    rerender();
    expect(result.current).toBe(firstResult); // 記憶體參考應相同 (cached)

    // 重新渲染，引數內容相同但 reference 變了 (但 useMemo 依賴 values，若陣列參考變了，因為 rules 依賴陣列 [rules, fallbackDeviceId]，Array reference 變了通常會觸發重算)
    // 我們需要確認 useMemo 依賴 [rules, fallbackDeviceId] 的機制，重新傳入新的 Rules Array 應觸發重新計算
    rules = [...mockRules];
    rerender();
    expect(result.current).not.toBe(firstResult); // 由於 rules reference 變了，觸發重算
  });
});

describe('useShareLayout', () => {
  it('應計算出正確的 layout，且引數未變時回傳快取的 layout 參考', () => {
    const rules = [
      {
        ...mockRules[0],
        share_enabled: true,
        share_start_register: 40001,
      },
    ];
    const baseReg = 40001;

    const { result, rerender } = renderHook(() => useShareLayout(rules, baseReg));
    const firstResult = result.current;

    expect(firstResult['rule-1']).toEqual({
      start: 40001,
      stride: 1,
      end: 40003,
      auto: false,
    });

    rerender();
    expect(result.current).toBe(firstResult); // cached
  });
});

describe('useConflictAddrs', () => {
  it('應檢測出衝突地址，且當點位陣列未變時回傳相同的 Set 參考', () => {
    const mockPoints: Point[] = [
      {
        id: 'p-1',
        device_id: 'dev-1',
        rule_id: 'rule-1',
        rule_name: 'Rule 1',
        name: 'TAG_0',
        address: '40001',
        data_type: 'int16',
        function: 'holding_register',
        width: 1,
        enabled: true,
        skipped: false,
        _rule_scale: 1,
        _rule_offset: 0,
      },
      {
        id: 'p-2',
        device_id: 'dev-1',
        rule_id: 'rule-2',
        rule_name: 'Rule 2',
        name: 'TAG_1',
        address: '40001', // 重複
        data_type: 'int16',
        function: 'holding_register',
        width: 1,
        enabled: true,
        skipped: false,
        _rule_scale: 1,
        _rule_offset: 0,
      },
    ];

    const points = [...mockPoints];
    const { result, rerender } = renderHook(() => useConflictAddrs(points));
    const firstResult = result.current;

    expect(firstResult.has('40001')).toBe(true);

    rerender();
    expect(result.current).toBe(firstResult); // cached
  });
});

describe('useRulePoints', () => {
  it('應衍生出正確的 rule 點位，且引數未變時快取', () => {
    const rule = mockRules[0];
    const devId = 'dev-1';
    const skipped = new Set<string>();

    const { result, rerender } = renderHook(() => useRulePoints(rule, devId, skipped));
    const firstResult = result.current;

    expect(firstResult).toHaveLength(2);

    rerender();
    expect(result.current).toBe(firstResult); // cached
  });
});

describe('useMappingValidation', () => {
  it('應正確評估 mappings 狀態是否可繼續', () => {
    const mappings: Record<string, Mapping> = {
      'p-1': {
        point_id: 'p-1',
        tag_key: 'line01.temp.inlet',
        display_name: 'n1',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
      },
      'p-2': {
        point_id: 'p-2',
        tag_key: '   ',
        display_name: 'n2',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
      },
    };

    const { result } = renderHook(() => useMappingValidation(mappings));
    expect(result.current.canContinue).toBe(false);
    expect(result.current.emptyTagCount).toBe(1);
    expect(result.current.enabledCount).toBe(2);

    // 禁用 p-2
    mappings['p-2'].enabled = false;
    const { result: result2 } = renderHook(() => useMappingValidation(mappings));
    expect(result2.current.canContinue).toBe(true);
    expect(result2.current.emptyTagCount).toBe(0);
    expect(result2.current.enabledCount).toBe(1);
  });
});

describe('useSelectedMapping', () => {
  it('應返回正確的點位、映射與原始模擬值', () => {
    const points: Point[] = [
      {
        id: 'p-1',
        device_id: 'dev-1',
        rule_id: 'rule-1',
        rule_name: 'Rule 1',
        name: 'TAG_0',
        address: '40001',
        data_type: 'int16',
        function: 'holding_register',
        width: 1,
        enabled: true,
        skipped: false,
        _rule_scale: 1,
        _rule_offset: 0,
      },
    ];

    const mappings: Record<string, Mapping> = {
      'p-1': {
        point_id: 'p-1',
        tag_key: 'line01.temp.inlet',
        display_name: 'n1',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
      },
    };

    const { result } = renderHook(() => useSelectedMapping(0, points, mappings));
    expect(result.current).not.toBeNull();
    expect(result.current?.point).toEqual(points[0]);
    expect(result.current?.mapping).toEqual(mappings['p-1']);

    const { result: resultNull } = renderHook(() => useSelectedMapping(5, points, mappings));
    expect(resultNull.current).toBeNull();
  });
});
