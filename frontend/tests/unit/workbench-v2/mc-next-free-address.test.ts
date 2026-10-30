import { describe, expect, it } from 'vitest';
import { nextFreeStartAddress } from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import type { Rule } from '../../../src/features/datalink/workbench-v2/state/types';

function rule(start: string, count: number, dataType: Rule['data_type']): Rule {
  return {
    id: 'mc-rule', device_id: 'mc-device', name: 'MC', start_address: start,
    count, data_type: dataType, naming_prefix: 'MC_', enabled: true,
    scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [],
    share_enabled: false, share_start_register: null, share_stride: null,
  };
}

describe('MC next free source address', () => {
  it.each([
    ['X0', 8, 'bool', 'X10'],
    ['Y7', 1, 'bool', 'Y10'],
    ['W0', 17, 'int16', 'W11'],
    ['WA', 1, 'int16', 'WB'],
    ['BFF', 1, 'bool', 'B100'],
    ['D100', 2, 'int32', 'D104'],
    ['M0', 8, 'bool', 'M8'],
  ] as const)('advances %s by its physical source span', (start, count, type, expected) => {
    expect(nextFreeStartAddress([rule(start, count, type)], 'mc_3e', 'D0')).toBe(expected);
  });

  it('uses persisted legacy addresses when choosing an unoccupied start', () => {
    const saved = rule('W0', 17, 'int16');
    saved.persisted = true;
    saved.mc_persisted_points = {
      status: 'conflict',
      points: ['W0', 'W1', 'W2', 'W3', 'W4', 'W5', 'W6', 'W7', 'W8', 'W9',
        'W10', 'W11', 'W12', 'W13', 'W14', 'W15', 'W16']
        .map((address, index) => ({ point_id: `legacy-${index}`, address })),
    };
    expect(nextFreeStartAddress([saved], 'mc_3e', 'D0')).toBe('W17');
  });

  it('does not invent a free MC address before persisted evidence is available', () => {
    const saved = { ...rule('W0', 1, 'int16'), persisted: true };
    expect(nextFreeStartAddress([saved], 'mc_3e', 'D0')).toBe('');
  });

  it('keeps first-rule defaults and other protocol allocation unchanged', () => {
    expect(nextFreeStartAddress([], 'mc_3e', 'D0')).toBe('D0');
    expect(nextFreeStartAddress([rule('40001', 3, 'float32')], 'modbus_tcp', '40001')).toBe('40007');
    expect(nextFreeStartAddress([rule('R0', 3, 'int16')], 'fatek_fbs', 'R0')).toBe('R3');
  });
});
