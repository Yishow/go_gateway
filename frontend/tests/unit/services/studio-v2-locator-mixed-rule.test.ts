import { afterEach, describe, expect, it, vi } from 'vitest';
import { readStudioLocatorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));

// 後端 5.3 之後，規則回應一律帶 layout_mode，mixed 規則另帶 mixed_items；
// locator entry 不得因此把所有已存規則判成無法取得（驗收 Bug B）。
const locator: EntityLocator = { step: 2, workspace_id: 'W', device_id: 'D', rule_id: 'R', mapping_id: null, group_id: null, group_revision: null };
const at = '2026-11-05T00:00:00Z';
const baseRule = {
  id: 'R', workspace_id: 'W', device_id: 'D', revision_id: 'REV', start_address: 'D100', count: 2, data_type: 'int16', naming_prefix: 'MIX_',
  enabled: true, locked: false, origin: 'manual', skipped_addresses: [], created_at: at, updated_at: at,
};
const item = (id: string, address: string, dataType: string) => ({ item_id: id, name: id, address, register_offset: 0, data_type: dataType, semantic_kind: 'gauge' });
function payload(rule: Record<string, unknown>) {
  return { id: 'W', kind: 'single', status: 'ready', readiness_summary: null, ordered_device_ids: ['D'], created_at: at, updated_at: at,
    locator_entry: { workspace_revision: '', hydration_scope: { complete_workspace: false, device_ids: ['D'], rule_ids: ['R'], mapping_ids: [], group_ids: [] },
      devices: [{ id: 'D', name: 'PLC', description: '', protocol: 'mc_3e', status: 'draft', connection_config: '{}', last_test_at: null, last_test_success: null,
        last_test_error: '', created_at: at, updated_at: at }],
      source_rules: [rule], rule_heads: { R: 'REV' }, group_heads: {}, revalidated_mapping_ids: [] } };
}
const read = (rule: Record<string, unknown>) => {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: payload(rule) } });
  return readStudioLocatorEntry(locator);
};
afterEach(() => vi.resetAllMocks());

describe('locator entry source rule layout', () => {
  it('accepts homogeneous rules that carry layout_mode', async () => {
    const result = await read({ ...baseRule, layout_mode: 'homogeneous', scale_multiplier: 1, scale_offset: 0 });
    expect(result.ok).toBe(true);
  });
  it('accepts mixed rules with their items', async () => {
    const mixed = { ...baseRule, layout_mode: 'mixed', mixed_items: [item('item-1', 'D100', 'int16'), item('item-2', 'D200', 'int32')] };
    const result = await read(mixed);
    expect(result.ok).toBe(true);
    if (!result.ok) throw new Error('entry unavailable');
    expect(result.value.entry.source_rules![0]).toMatchObject({ layout_mode: 'mixed', mixed_items: mixed.mixed_items });
  });
  it('accepts mixed rules with mixed_items_error so the operator can open and repair them', async () => {
    const corruptMixed = { ...baseRule, layout_mode: 'mixed', mixed_items_error: 'invalid stored mixed_items' };
    const result = await read(corruptMixed);
    expect(result.ok).toBe(true);
    if (!result.ok) throw new Error('entry unavailable');
    expect(result.value.entry.source_rules![0]).toMatchObject({
      layout_mode: 'mixed',
      mixed_items_error: 'invalid stored mixed_items',
    });
  });
  it.each([
    { layout_mode: 'mixed' },
    { layout_mode: 'mixed', mixed_items: [item('a', 'D1', 'int16')], mixed_items_error: 'corrupt' },
    { layout_mode: 'homogeneous', mixed_items: [item('a', 'D1', 'int16')] },
    { layout_mode: 'homogeneous', mixed_items_error: 'corrupt' },
    { layout_mode: 'mixed', mixed_items: [] },
    { layout_mode: 'mixed', mixed_items: [item('', 'D1', 'int16')] },
    { layout_mode: 'mixed', mixed_items: [item('a', '   ', 'int16')] },
    { layout_mode: 'mixed', mixed_items: [item('a', 'D1', 'string')] },
    { layout_mode: 'mixed', mixed_items: [{ ...item('a', 'D1', 'int16'), state_map: { 0: 'off' } }] },
    { layout_mode: 'mixed', mixed_items: [{ ...item('a', 'D1', 'int16'), register_offset: 2 }] },
    { layout_mode: 'mixed', mixed_items: [item('a', 'D1', 'not-a-type')] },
    { layout_mode: 'sparse' },
  ])('fails closed for inconsistent layout %o', async (layout) => {
    expect(await read({ ...baseRule, ...layout })).toEqual({ ok: false, reason: 'unavailable' });
  });
});
