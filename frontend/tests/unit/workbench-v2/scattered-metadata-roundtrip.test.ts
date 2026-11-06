import { describe, expect, it } from 'vitest';
import { toStudioV2RuleUpdateRequest } from '@/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import type { Rule } from '@/features/datalink/workbench-v2/state/types';

describe('scattered range save', () => {
  it('preserves a saved item’s measurement metadata when the rule is edited', () => {
    const rule: Rule = {
      id: 'rule-1', device_id: 'device-1', name: 'energy', naming_prefix: 'energy',
      start_address: 'D100', count: 1, data_type: 'uint32', enabled: true,
      scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [],
      share_enabled: false, share_start_register: null, share_stride: null,
      layout_mode: 'mixed', persisted: true, save_state: 'saved', save_error: null,
      mixed_items: [{
        item_id: 'energy', name: 'Energy', address: 'D100', register_offset: 0,
        data_type: 'uint32', semantic_kind: 'counter', quantity: 'energy', unit: 'kWh',
      }],
    };

    expect(toStudioV2RuleUpdateRequest(rule).mixed_items).toEqual([
      expect.objectContaining({ item_id: 'energy', quantity: 'energy', unit: 'kWh' }),
    ]);
    expect(toStudioV2RuleUpdateRequest({ ...rule, layout_mode: 'homogeneous' }).layout_mode).toBe('homogeneous');
  });
});
