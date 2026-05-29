import { describe, it, expect } from 'vitest';
import {
  workbenchV2Reducer,
  cascadeRemoveDevice,
  INITIAL_STATE,
} from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Rule } from '../../../src/features/datalink/workbench-v2/state/types';

const mockRule: Rule = {
  id: 'rule-test-1',
  device_id: 'dev-01',
  name: 'Test Rule 1',
  start_address: '40001',
  count: 10,
  data_type: 'int16',
  naming_prefix: 'TEST_',
  enabled: true,
  scale_multiplier: 1,
  scale_offset: 0,
  data_format: '',
  skipped_addresses: ['40003'],
  share_enabled: true,
  share_start_register: 40001,
  share_stride: null,
};

describe('WorkbenchV2 Reducer - Step 2 Rule Actions', () => {
  it('selectRule: 應正確切換 selectedRuleId', () => {
    const state = { ...INITIAL_STATE, selectedRuleId: 'rule-01' };
    const nextState = workbenchV2Reducer(state, { type: 'selectRule', ruleId: 'rule-test-1' });
    expect(nextState.selectedRuleId).toBe('rule-test-1');
  });

  it('addRule: 應加入規則並將 selectedRuleId 切為新規則', () => {
    const state = { ...INITIAL_STATE, rules: [] as Rule[], selectedRuleId: null };
    const nextState = workbenchV2Reducer(state, { type: 'addRule', rule: mockRule });
    expect(nextState.rules).toHaveLength(1);
    expect(nextState.rules[0].id).toBe('rule-test-1');
    expect(nextState.selectedRuleId).toBe('rule-test-1');
  });

  it('removeRule: 應刪除規則，若被刪的是當前 selectedRuleId，自動切換至剩餘的第一個', () => {
    const rule2 = { ...mockRule, id: 'rule-test-2', name: 'Test Rule 2' };
    const state = {
      ...INITIAL_STATE,
      rules: [mockRule, rule2],
      selectedRuleId: 'rule-test-1',
    };

    // 刪除選中的 rule-test-1，應切到 rule-test-2
    let nextState = workbenchV2Reducer(state, { type: 'removeRule', ruleId: 'rule-test-1' });
    expect(nextState.rules).toHaveLength(1);
    expect(nextState.rules[0].id).toBe('rule-test-2');
    expect(nextState.selectedRuleId).toBe('rule-test-2');

    // 再次刪除剩餘唯一的 rule-test-2，selectedRuleId 應切為 null
    nextState = workbenchV2Reducer(nextState, { type: 'removeRule', ruleId: 'rule-test-2' });
    expect(nextState.rules).toHaveLength(0);
    expect(nextState.selectedRuleId).toBeNull();
  });

  describe('updateRule', () => {
    it('應正常更新 patch 欄位且不影響 skipped_addresses', () => {
      const state = { ...INITIAL_STATE, rules: [mockRule] };
      const nextState = workbenchV2Reducer(state, {
        type: 'updateRule',
        ruleId: 'rule-test-1',
        patch: { name: 'Updated Name', scale_multiplier: 2 },
      });
      const updated = nextState.rules[0];
      expect(updated.name).toBe('Updated Name');
      expect(updated.scale_multiplier).toBe(2);
      expect(updated.skipped_addresses).toEqual(['40003']); // 保留 skipped
    });

    it('更新 start_address, count 或 data_type 時應清空 skipped_addresses 且限制 count 範圍', () => {
      const state = { ...INITIAL_STATE, rules: [mockRule] };

      // 更新 start_address 觸發 reset
      let nextState = workbenchV2Reducer(state, {
        type: 'updateRule',
        ruleId: 'rule-test-1',
        patch: { start_address: '40010' },
      });
      expect(nextState.rules[0].skipped_addresses).toEqual([]);

      // 更新 count 並限制在 1-64 區間，超限 clamp 為 64
      nextState = workbenchV2Reducer(state, {
        type: 'updateRule',
        ruleId: 'rule-test-1',
        patch: { count: 100 },
      });
      expect(nextState.rules[0].count).toBe(64);
      expect(nextState.rules[0].skipped_addresses).toEqual([]); // count 改變亦 reset

      // count 低於 1 應 clamp 為 1
      nextState = workbenchV2Reducer(state, {
        type: 'updateRule',
        ruleId: 'rule-test-1',
        patch: { count: -5 },
      });
      expect(nextState.rules[0].count).toBe(1);
    });
  });

  it('renameRule: 應改名', () => {
    const state = { ...INITIAL_STATE, rules: [mockRule] };
    const nextState = workbenchV2Reducer(state, {
      type: 'renameRule',
      ruleId: 'rule-test-1',
      name: 'New Fancy Name',
    });
    expect(nextState.rules[0].name).toBe('New Fancy Name');
  });

  it('toggleRuleEnabled: 應切換 enabled 狀態', () => {
    const state = { ...INITIAL_STATE, rules: [mockRule] };
    let nextState = workbenchV2Reducer(state, { type: 'toggleRuleEnabled', ruleId: 'rule-test-1' });
    expect(nextState.rules[0].enabled).toBe(false);
    nextState = workbenchV2Reducer(nextState, { type: 'toggleRuleEnabled', ruleId: 'rule-test-1' });
    expect(nextState.rules[0].enabled).toBe(true);
  });

  it('updateRuleSkipped: 應整批更新 skippedAddresses', () => {
    const state = { ...INITIAL_STATE, rules: [mockRule] };
    const nextState = workbenchV2Reducer(state, {
      type: 'updateRuleSkipped',
      ruleId: 'rule-test-1',
      skippedAddresses: ['40001', '40002'],
    });
    expect(nextState.rules[0].skipped_addresses).toEqual(['40001', '40002']);
  });

  it('toggleRuleSkippedAddress: 應 toggle 單一 address 狀態', () => {
    const state = { ...INITIAL_STATE, rules: [mockRule] };
    // 原本 skipped 含有 '40003'，toggle '40003' 應移除
    let nextState = workbenchV2Reducer(state, {
      type: 'toggleRuleSkippedAddress',
      ruleId: 'rule-test-1',
      address: '40003',
    });
    expect(nextState.rules[0].skipped_addresses).not.toContain('40003');

    // 再次 toggle '40003' 應加回
    nextState = workbenchV2Reducer(nextState, {
      type: 'toggleRuleSkippedAddress',
      ruleId: 'rule-test-1',
      address: '40003',
    });
    expect(nextState.rules[0].skipped_addresses).toContain('40003');
  });

  it('toggleRuleShareEnabled: 應切換 share_enabled 狀態', () => {
    const state = { ...INITIAL_STATE, rules: [mockRule] };
    const nextState = workbenchV2Reducer(state, { type: 'toggleRuleShareEnabled', ruleId: 'rule-test-1' });
    expect(nextState.rules[0].share_enabled).toBe(false);
  });

  it('updateRuleShareStart / updateRuleShareStride: 應正常改值', () => {
    const state = { ...INITIAL_STATE, rules: [mockRule] };
    let nextState = workbenchV2Reducer(state, {
      type: 'updateRuleShareStart',
      ruleId: 'rule-test-1',
      shareStart: 50001,
    });
    expect(nextState.rules[0].share_start_register).toBe(50001);

    nextState = workbenchV2Reducer(nextState, {
      type: 'updateRuleShareStride',
      ruleId: 'rule-test-1',
      shareStride: 2,
    });
    expect(nextState.rules[0].share_stride).toBe(2);
  });

  it('cascadeRemoveDevice: 刪除設備時連動刪除 rule，若選中的 rule 被刪除需重新選擇', () => {
    const state = {
      ...INITIAL_STATE,
      devices: [
        { id: 'dev-1', name: 'PLC 1', description: '', protocol: 'modbus_tcp' as const, config: {}, status: 'draft' as const, test: null },
      ],
      rules: [
        { ...mockRule, id: 'rule-dev1', device_id: 'dev-1' },
      ],
      selectedRuleId: 'rule-dev1',
    };

    const nextState = cascadeRemoveDevice(state, 'dev-1');
    expect(nextState.rules).toHaveLength(0);
    expect(nextState.selectedRuleId).toBeNull();
  });
});
