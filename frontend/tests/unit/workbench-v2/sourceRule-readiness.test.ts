import { describe, it, expect } from 'vitest';
import {
  deriveAllPoints,
  getRuleReadinessIssues,
  getRuleReadinessReason,
  isStep2Ready,
  computeShareLayout,
} from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import type { Rule } from '../../../src/features/datalink/workbench-v2/state/types';
describe('SourceRule readiness and derived points', () => {
  const makeRule = (overrides: Partial<Rule> = {}): Rule => ({
    id: 'rule-1',
    device_id: 'dev-1',
    name: 'Holding Rule',
    start_address: '40001',
    count: 1,
    data_type: 'int16',
    naming_prefix: 'TAG_',
    enabled: true,
    scale_multiplier: 1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: false,
    share_start_register: null,
    share_stride: null,
    ...overrides,
  });
  describe('deriveAllPoints', () => {
    it('應能合併多條規則衍生出的全部點位', () => {
      const rules: Rule[] = [
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
        {
          id: 'rule-2',
          device_id: 'dev-2', // 不同裝置
          name: 'Inputs',
          start_address: '30001',
          count: 1,
          data_type: 'int16',
          naming_prefix: 'I_',
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

      const all = deriveAllPoints(rules, { 'dev-1': 'modbus_tcp', 'dev-2': 'modbus_tcp' });
      expect(all).toHaveLength(3);
      expect(all[0].device_id).toBe('dev-1');
      expect(all[2].device_id).toBe('dev-2');
    });
    it('指向不存在或未對應協議的設備時不應衍生點位，不得自動回退至備用設備', () => {
      const orphanRule: Rule = {
        id: 'rule-orphan',
        device_id: 'deleted-device',
        name: 'Orphan Rule',
        start_address: '40001',
        count: 2,
        data_type: 'int16',
        naming_prefix: 'O_',
        enabled: true,
        scale_multiplier: 1,
        scale_offset: 0,
        data_format: '',
        skipped_addresses: [],
        share_enabled: false,
        share_start_register: null,
        share_stride: null,
      };

      const all = deriveAllPoints([orphanRule], { 'dev-valid': 'modbus_tcp' });
      expect(all).toHaveLength(0);
    });
    it('多設備具備不同協議時，各規則應以自身所屬設備的協議進行解析', () => {
      const rules: Rule[] = [
        {
          id: 'rule-modbus',
          device_id: 'dev-modbus',
          name: 'Modbus Rule',
          start_address: '40001',
          count: 2,
          data_type: 'int16',
          naming_prefix: 'M_',
          enabled: true,
          scale_multiplier: 1,
          scale_offset: 0,
          data_format: '',
          skipped_addresses: [],
          share_enabled: false,
          share_start_register: null,
          share_stride: null,
        },
        {
          id: 'rule-mc',
          device_id: 'dev-mc',
          name: 'MC Rule',
          start_address: 'D0',
          count: 2,
          data_type: 'int16',
          naming_prefix: 'MC_',
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

      const all = deriveAllPoints(rules, {
        'dev-modbus': 'modbus_tcp',
        'dev-mc': 'mc_3e',
      });
      expect(all).toHaveLength(4);
      expect(all[0].address).toBe('40001');
      expect(all[1].address).toBe('40002');
      expect(all[2].address).toBe('D0');
      expect(all[3].address).toBe('D1');
      expect(all[2].function).toBe('D (Word)');
    });
    it('規則非空但設備協議對照表為空時應 fail closed，不得套用 fallback 協議', () => {
      const rule = makeRule();

      expect(deriveAllPoints([rule], {})).toEqual([]);
    });
    it('應以同一個 typed readiness 判斷區分 unknown、deleted 與 invalid address，並攜帶 identity', () => {
      const baseRule = makeRule();
      const unknownRule = { ...baseRule, id: 'rule-unknown', device_id: 'unknown-device', persisted: false };
      const emptyDeviceRule = { ...baseRule, id: 'rule-empty-device', device_id: '', persisted: false };
      const deletedRule = { ...baseRule, id: 'rule-deleted', device_id: 'deleted-device', persisted: true };
      const invalidRule = { ...baseRule, id: 'rule-invalid', start_address: 'Z999', persisted: true };
      const disabledInvalidRule = {
        ...invalidRule,
        id: 'rule-disabled-invalid',
        enabled: false,
      };
      const protocolMap = { 'dev-1': 'modbus_tcp' as const };

      expect(getRuleReadinessReason(unknownRule, protocolMap)).toBe('unknown_device');
      expect(getRuleReadinessReason(emptyDeviceRule, protocolMap)).toBe('unknown_device');
      expect(getRuleReadinessReason(deletedRule, protocolMap)).toBe('deleted_device');
      expect(getRuleReadinessReason(invalidRule, protocolMap)).toBe('invalid_address');
      expect(getRuleReadinessIssues(
        [unknownRule, deletedRule, invalidRule, disabledInvalidRule],
        protocolMap,
      ).map(({ ruleId, deviceId, startAddress, reason }) => ({
        ruleId,
        deviceId,
        startAddress,
        reason,
      }))).toEqual([
        { ruleId: 'rule-unknown', deviceId: 'unknown-device', startAddress: '40001', reason: 'unknown_device' },
        { ruleId: 'rule-deleted', deviceId: 'deleted-device', startAddress: '40001', reason: 'deleted_device' },
        { ruleId: 'rule-invalid', deviceId: 'dev-1', startAddress: 'Z999', reason: 'invalid_address' },
      ]);
    });

    it('isStep2Ready 應要求每個啟用規則都有有效點位，且忽略停用的無效規則', () => {
      const validRule = makeRule({ id: 'rule-valid' });
      const disabledInvalidRule = makeRule({
        id: 'rule-disabled-invalid',
        enabled: false,
        start_address: 'Z999',
      });
      const protocolMap = { 'dev-1': 'modbus_tcp' as const };

      expect(isStep2Ready([validRule, disabledInvalidRule], protocolMap)).toBe(true);
      expect(isStep2Ready([{ ...validRule, skipped_addresses: ['40001'] }], protocolMap)).toBe(false);
      expect(isStep2Ready([{ ...validRule, start_address: 'Z999' }], protocolMap)).toBe(false);
      expect(isStep2Ready([{ ...validRule, device_id: 'deleted-device', persisted: true }], protocolMap)).toBe(false);
    });
  });

  describe('computeShareLayout', () => {
    it('應正確計算手動、自動、未啟用的混合佈局，且自動起點依 cursor 取最大值', () => {
      const rules: Rule[] = [
        {
          id: 'rule-1',
          device_id: 'dev-1',
          name: 'Auto Share 1',
          start_address: '40001',
          count: 5, // 5 * 2 = 10 registers
          data_type: 'int32',
          naming_prefix: 'A_',
          enabled: true,
          scale_multiplier: 1,
          scale_offset: 0,
          data_format: '',
          skipped_addresses: ['40003'], // skipped 1 point, so count = 4, 4 * 2 = 8 registers
          share_enabled: true,
          share_start_register: null, // 自動
          share_stride: null,
        },
        {
          id: 'rule-2',
          device_id: 'dev-1',
          name: 'Disabled Share',
          start_address: '40010',
          count: 5,
          data_type: 'int16',
          naming_prefix: 'B_',
          enabled: true,
          scale_multiplier: 1,
          scale_offset: 0,
          data_format: '',
          skipped_addresses: [],
          share_enabled: false, // 停用
          share_start_register: null,
          share_stride: null,
        },
        {
          id: 'rule-3',
          device_id: 'dev-1',
          name: 'Manual Override',
          start_address: '40020',
          count: 2, // 2 * 4 = 8 registers
          data_type: 'float64',
          naming_prefix: 'C_',
          enabled: true,
          scale_multiplier: 1,
          scale_offset: 0,
          data_format: '',
          skipped_addresses: [],
          share_enabled: true,
          share_start_register: 40100, // 手動指定起點 40100
          share_stride: null,
        },
        {
          id: 'rule-4',
          device_id: 'dev-1',
          name: 'Auto Share 2',
          start_address: '40030',
          count: 2, // 2 * 1 = 2 registers
          data_type: 'int16',
          naming_prefix: 'D_',
          enabled: true,
          scale_multiplier: 1,
          scale_offset: 0,
          data_format: '',
          skipped_addresses: [],
          share_enabled: true,
          share_start_register: null, // 自動，接續前一個的 end (40100 + 2*4 = 40108)
          share_stride: null,
        },
      ];

      const layouts = computeShareLayout(rules, 40001);

      // rule-1: start 40001, end 40001 + 4*2 = 40009
      expect(layouts['rule-1']).toEqual({
        start: 40001,
        stride: 2,
        end: 40009,
        auto: true,
      });

      // rule-2: 未啟用 share 應為 null
      expect(layouts['rule-2']).toBeNull();

      // rule-3: 手動指定為 40100, end 40100 + 2*4 = 40108
      expect(layouts['rule-3']).toEqual({
        start: 40100,
        stride: 4,
        end: 40108,
        auto: false,
      });

      // rule-4: 自動，因前一條 end 是 40108，所以 start 應為 40108，end 為 40110
      expect(layouts['rule-4']).toEqual({
        start: 40108,
        stride: 1,
        end: 40110,
        auto: true,
      });
    });
  });
});
