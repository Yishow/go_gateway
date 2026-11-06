import { afterEach, describe, expect, it, vi } from 'vitest';
import { readStudioLocatorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const locator: EntityLocator = { step: 1, workspace_id: 'W', device_id: 'D', rule_id: null, mapping_id: null, group_id: null, group_revision: null };
const at = '2026-11-05T00:00:00Z';
function payload(acquisition: Record<string, unknown>) {
  return { id: 'W', kind: 'single', status: 'ready', readiness_summary: null, ordered_device_ids: ['D'], created_at: at, updated_at: at,
    locator_entry: { workspace_revision: '', hydration_scope: { complete_workspace: false, device_ids: ['D'], rule_ids: [], mapping_ids: [], group_ids: [] },
      devices: [{ id: 'D', name: 'PLC', description: '', protocol: 'mc_3e', status: 'draft', connection_config: '{}', last_test_at: null, last_test_success: null,
        last_test_error: '', created_at: at, updated_at: at, ...acquisition }], rule_heads: {}, group_heads: {}, revalidated_mapping_ids: [] } };
}
const read = (acquisition: Record<string, unknown>) => {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: payload(acquisition) } });
  return readStudioLocatorEntry(locator);
};
afterEach(() => vi.resetAllMocks());

describe('locator entry device acquisition settings', () => {
  it('accepts the acquisition fields the device API now returns', async () => {
    const result = await read({ acquisition_interval_ms: 500, acquisition_auto_interval: true, stale_threshold_ms: null, frozen_threshold_ms: 0 });
    expect(result.ok).toBe(true);
    if (!result.ok) throw new Error('entry unavailable');
    expect(result.value.entry.devices![0]).toMatchObject({ acquisition_interval_ms: 500, acquisition_auto_interval: true, stale_threshold_ms: null, frozen_threshold_ms: 0 });
    expect((await read({ acquisition_auto_interval: false, stale_threshold_ms: 2500, frozen_threshold_ms: 60000 })).ok).toBe(true);
  });
  it.each([
    { acquisition_interval_ms: 0 }, { acquisition_interval_ms: 1.5 }, { acquisition_interval_ms: '500' }, { acquisition_auto_interval: 'true' },
    { stale_threshold_ms: 0 }, { frozen_threshold_ms: -1 }, { frozen_threshold_ms: null },
  ])('fails closed for malformed %o', async (acquisition) => {
    expect(await read(acquisition)).toEqual({ ok: false, reason: 'unavailable' });
  });
});
