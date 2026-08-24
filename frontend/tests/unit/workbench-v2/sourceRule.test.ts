import { describe, it, expect } from 'vitest';
import {
  dataTypeWidth,
  fnFromAddr,
  formatAddr,
  derivePoints,
  detectAddressConflicts,
} from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import type { Rule } from '../../../src/features/datalink/workbench-v2/state/types';

describe('SourceRule State Pure Helpers', () => {
  describe('dataTypeWidth', () => {
    it('應為 10 種不同型別回傳正確的暫存器寬度', () => {
      expect(dataTypeWidth('bool')).toBe(1);
      expect(dataTypeWidth('int16')).toBe(1);
      expect(dataTypeWidth('uint16')).toBe(1);
      expect(dataTypeWidth('int32')).toBe(2);
      expect(dataTypeWidth('uint32')).toBe(2);
      expect(dataTypeWidth('float32')).toBe(2);
      expect(dataTypeWidth('int64')).toBe(4);
      expect(dataTypeWidth('uint64')).toBe(4);
      expect(dataTypeWidth('float64')).toBe(4);
      expect(dataTypeWidth('string')).toBe(10);
    });
  });

  describe('fnFromAddr', () => {
    it('應根據前綴正確推斷 Modbus Function Type', () => {
      expect(fnFromAddr('00001')).toBe('coil');
      expect(fnFromAddr('10001')).toBe('discrete_input');
      expect(fnFromAddr('30001')).toBe('input_register');
      expect(fnFromAddr('40001')).toBe('holding_register');
      expect(fnFromAddr(' 40010 ')).toBe('holding_register'); // 容忍空格
      expect(() => fnFromAddr('invalid')).toThrow();
    });
  });

  describe('formatAddr', () => {
    it('應正確格式化地址為字串', () => {
      expect(formatAddr(40001)).toBe('40001');
      expect(formatAddr('30005')).toBe('30005');
    });
  });

  describe('derivePoints', () => {
    const mockRule: Rule = {
      id: 'rule-1',
      device_id: 'dev-1',
      name: 'Holding Rule',
      start_address: '40001',
      count: 4,
      data_type: 'int32',
      naming_prefix: 'TAG_',
      enabled: true,
      scale_multiplier: 0.5,
      scale_offset: -1,
      data_format: '',
      skipped_addresses: ['40003'],
      share_enabled: false,
      share_start_register: null,
      share_stride: null,
    };

    it('應能依據 count 與 stride 正確產生點位陣列，且支援已略過(skipped)狀態', () => {
      const skippedSet = new Set(['40003']);
      const points = derivePoints(mockRule, 'dev-1', skippedSet);

      expect(points).toHaveLength(4);
      
      // 第一個點 (index 0)：絕對位址 40001
      expect(points[0]).toEqual({
        id: 'rule-1-p-0',
        device_id: 'dev-1',
        rule_id: 'rule-1',
        rule_name: 'Holding Rule',
        name: 'TAG_0',
        address: '40001',
        data_type: 'int32',
        function: 'holding_register',
        width: 2,
        enabled: true,
        skipped: false,
        _rule_scale: 0.5,
        _rule_offset: -1,
      });

      // 第二個點 (index 1)：位址為 40001 + 1*stride = 40003，此點被 skipped 故 enabled 應為 false
      expect(points[1].address).toBe('40003');
      expect(points[1].skipped).toBe(true);
      expect(points[1].enabled).toBe(false);

      // 第三個點 (index 2)：位址為 40001 + 2*2 = 40005
      expect(points[2].address).toBe('40005');
      expect(points[2].skipped).toBe(false);
      expect(points[2].enabled).toBe(true);
    });

    it('應正確衍生 Mitsubishi MC 3E (D0) 點位位址且保留 D 前綴', () => {
      const mcRule: Rule = {
        ...mockRule,
        start_address: 'D0',
        count: 4,
        data_type: 'int16',
      };
      const points = derivePoints(mcRule, 'dev-mc', new Set(), 'mc_3e');
      expect(points.map((p) => p.address)).toEqual(['D0', 'D1', 'D2', 'D3']);
      expect(points[0].function).toBe('D (Word)');
    });

    it('應正確處理 MC 3E 32-bit stride 步進 (D100 -> D102)', () => {
      const mcRule: Rule = {
        ...mockRule,
        start_address: 'D100',
        count: 2,
        data_type: 'int32',
      };
      const points = derivePoints(mcRule, 'dev-mc', new Set(), 'mc_3e');
      expect(points.map((p) => p.address)).toEqual(['D100', 'D102']);
    });

    it('應正確衍生 MC 3E Bit 設備位址 (X0)', () => {
      const mcRule: Rule = {
        ...mockRule,
        start_address: 'X0',
        count: 3,
        data_type: 'bool',
      };
      const points = derivePoints(mcRule, 'dev-mc', new Set(), 'mc_3e');
      expect(points.map((p) => p.address)).toEqual(['X0', 'X1', 'X2']);
      expect(points[0].function).toBe('X (Bit)');
    });

    it('應正確處理 MC 3E 十六進位接點跨位與進位步進 (XF -> X10 -> X11)', () => {
      const mcRuleXF: Rule = {
        ...mockRule,
        start_address: 'XF',
        count: 3,
        data_type: 'bool',
      };
      const points = derivePoints(mcRuleXF, 'dev-mc', new Set(), 'mc_3e');
      expect(points.map((p) => p.address)).toEqual(['XF', 'X10', 'X11']);
      expect(points[0].function).toBe('X (Bit)');
    });

    it('應正確衍生 FATEK FBs 設備位址 (R0 與 D100)', () => {
      const fatekRule: Rule = {
        ...mockRule,
        start_address: 'R0',
        count: 3,
        data_type: 'int16',
      };
      const points = derivePoints(fatekRule, 'dev-fatek', new Set(), 'fatek_fbs');
      expect(points.map((p) => p.address)).toEqual(['R0', 'R1', 'R2']);
      expect(points[0].function).toBe('R (Word)');
    });

    it('遇到不符合設備協議的位址時應 fail-closed，不得產生原字串或 Modbus 點位', () => {
      const invalidRule: Rule = {
        ...mockRule,
        start_address: 'Z999',
        count: 1,
      };

      expect(derivePoints(invalidRule, 'dev-mc', new Set(), 'mc_3e')).toEqual([]);
    });

    it('遇到空位址時應 fail-closed，不得套用協議預設位址', () => {
      const emptyRule: Rule = {
        ...mockRule,
        start_address: '',
        count: 1,
      };

      expect(derivePoints(emptyRule, 'dev-mc', new Set(), 'mc_3e')).toEqual([]);
    });
  });

  describe('detectAddressConflicts', () => {
    it('跨規則同址且皆啟用時判定為衝突，skipped 點位應排除不衝突', () => {
      const allPoints = [
        {
          id: 'r1-p0',
          device_id: 'dev-1',
          rule_id: 'r1',
          rule_name: 'Rule 1',
          name: 'TAG_0',
          address: '40001',
          data_type: 'int16',
          function: 'holding_register' as const,
          width: 1,
          enabled: true,
          skipped: false,
          _rule_scale: 1,
          _rule_offset: 0,
        },
        {
          id: 'r2-p0',
          device_id: 'dev-1',
          rule_id: 'r2',
          rule_name: 'Rule 2',
          name: 'TAG_A',
          address: '40001', // 同址衝突
          data_type: 'int16',
          function: 'holding_register' as const,
          width: 1,
          enabled: true,
          skipped: false,
          _rule_scale: 1,
          _rule_offset: 0,
        },
        {
          id: 'r3-p0',
          device_id: 'dev-1',
          rule_id: 'r3',
          rule_name: 'Rule 3',
          name: 'TAG_B',
          address: '40005',
          data_type: 'int16',
          function: 'holding_register' as const,
          width: 1,
          enabled: true,
          skipped: false,
          _rule_scale: 1,
          _rule_offset: 0,
        },
        {
          id: 'r4-p0',
          device_id: 'dev-1',
          rule_id: 'r4',
          rule_name: 'Rule 4',
          name: 'TAG_C',
          address: '40005', // 同址，但此點被 skipped，不應列入衝突
          data_type: 'int16',
          function: 'holding_register' as const,
          width: 1,
          enabled: false,
          skipped: true,
          _rule_scale: 1,
          _rule_offset: 0,
        },
      ];

      const conflicts = detectAddressConflicts(allPoints);
      expect(conflicts.size).toBe(1);
      expect(conflicts.has('40001')).toBe(true);
      expect(conflicts.has('40005')).toBe(false);
    });
  });
});
