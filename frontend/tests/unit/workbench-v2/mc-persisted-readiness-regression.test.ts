import { describe, expect, it } from 'vitest';
import { getRuleReadinessReason, projectPersistedMCPoint } from '../../../src/features/datalink/workbench-v2/state/sourceRule';
import type { Rule } from '../../../src/features/datalink/workbench-v2/state/types';

function mcRule(patch: Partial<Rule> = {}): Rule {
  return {
    id: 'rule-mc', device_id: 'dev-mc', name: 'MC', start_address: 'W0', count: 2,
    data_type: 'int16', naming_prefix: 'MC_', enabled: true,
    scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [],
    share_enabled: false, share_start_register: null, share_stride: null,
    persisted: true, save_state: 'saved', mc_persisted_points: {
      status: 'compatible', points: [{ point_id: 'p0', address: 'W0' }, { point_id: 'p1', address: 'W1' }],
    },
    ...patch,
  };
}

describe('MC persisted projection and current draft readiness', () => {
  it.each([
    ['WA', 'W (Word)'], ['WF', 'W (Word)'], ['BFF', 'B (Bit)'], ['X8', 'X (Bit)'],
  ])('keeps the single MC area separate from the suffix in %s', (address, label) => {
    expect(projectPersistedMCPoint(mcRule(), { point_id: 'native-point', address })).toMatchObject({
      id: 'native-point', address, function: label, data_type: 'int16', width: 1,
    });
  });

  it.each([
    { start_address: 'WFFFFFFFF', count: 2 },
    { start_address: 'WFFFFFF', count: 2 },
  ])('blocks the current invalid draft %o despite compatible saved evidence', (patch) => {
    expect(getRuleReadinessReason(mcRule(patch), { 'dev-mc': 'mc_3e' })).toBe('invalid_address');
    expect(getRuleReadinessReason(mcRule({ ...patch, persisted: false }), { 'dev-mc': 'mc_3e' })).toBe('invalid_address');
  });

  it('accepts a valid maximum range and keeps legacy conflict evidence authoritative', () => {
    expect(getRuleReadinessReason(mcRule({ start_address: 'WFFFFFE' }), { 'dev-mc': 'mc_3e' })).toBeNull();
    const legacy = mcRule({ start_address: 'X8', mc_persisted_points: { status: 'conflict', points: [{ point_id: 'legacy', address: 'X8' }] } });
    expect(getRuleReadinessReason(legacy, { 'dev-mc': 'mc_3e' })).toBe('mc_points_conflict');
    expect(projectPersistedMCPoint(legacy, { point_id: 'legacy', address: 'X8' }).id).toBe('legacy');
  });
});
