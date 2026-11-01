import { afterEach, expect, it, vi } from 'vitest';
import { readAcquisitionPlanPreview } from '@/services/acquisitionPlanPreview';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
afterEach(() => vi.resetAllMocks());
const scope = { workspace_id: 'w', device_id: 'd', expected_workspace_revision: '' };
function fixture() {
  return { workspace_id: 'w', device_id: 'd', workspace_revision: '', status: 'available', preview_revision: 'head', estimate: true,
    groups: [{ polling_group_id: 'g', poll_interval_ms: 1000, plan_revision: 'plan', declared_plan_revision: 'plan', capability_revision: 'cap', basis: 'declared_current_connection', acquisition_point_count: 1, planned_requests: 1,
      unique_wire_units: [{ area: 'D', function: '', address_unit_bits: 16 as unknown, unique_address_units: 1 }], point_errors: [], fallback_reasons: [],
      request_limits: [{ kind: 'contiguous', area: 'D', function: '', native_unit_limit: 125 as unknown, request_frame_bytes: 12, response_frame_bytes: 12, frame_cost_known: true }],
      actual: { status: 'available', acquisition_id: 'scan', observed_at: '2026-10-31T00:00:00Z', declared_plan_revision: 'plan', initial_plan_revision: 'plan', phases: [{ phase: 'initial', plan_revision: 'plan' }], wire_attempts: 2, retry_wire_attempts: 1, fallback_wire_attempts: 2, logical_requests: { status: 'observed', succeeded: 1, failed: 0, canceled: 0 }, cycle_latency_ms: 1, overrun: false, last_good_age_ms: 1 } }],
    output_loss: { status: 'unavailable', reason: 'no_matching_retained_detail' } };
}
async function read(value = fixture()) { vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: value } }); return readAcquisitionPlanPreview(scope); }
it('accepts legacy-empty workspace and overlapping retry/fallback subsets', async () => {
  expect((await read()).ok).toBe(true);
  expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/acquisition-plan-preview', expect.objectContaining({ params: scope }));
});
it.each(['16', true, null])('rejects nonnumeric physical unit %s', async value => { const data = fixture(); data.groups[0].unique_wire_units[0].address_unit_bits = value; expect((await read(data)).ok).toBe(false); });
it.each(['125', -1, null])('rejects malformed request limit %s', async value => { const data = fixture(); data.groups[0].request_limits[0].native_unit_limit = value; expect((await read(data)).ok).toBe(false); });
it.each(['retry_wire_attempts', 'fallback_wire_attempts'] as const)('rejects %s exceeding total', async key => { const data = fixture(); data.groups[0].actual[key] = 3; expect((await read(data)).ok).toBe(false); });
it('does not fetch an explicitly empty preview head', async () => { expect((await readAcquisitionPlanPreview({ ...scope, expected_preview_revision: '' })).ok).toBe(false); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled(); });
it('drops aborted late responses', async () => { const controller = new AbortController(); controller.abort(); vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: fixture() } }); expect((await readAcquisitionPlanPreview(scope, controller.signal)).ok).toBe(false); });
it.each([403, 409, 500])('maps safe HTTP reason %s', async status => { vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { status }, message: 'secret' }); expect(await readAcquisitionPlanPreview(scope)).toEqual({ ok: false, reason: status === 403 ? 'foreign-scope' : status === 409 ? 'stale-reference' : 'unavailable' }); });
