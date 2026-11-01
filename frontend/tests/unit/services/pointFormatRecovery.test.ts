import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '../../../src/services/datalinkClient';
import { createPointFormatOperation } from '../../../src/services/pointFormat';
import { createPointFormatRecovery, pointFormatRecoveryAPI } from '../../../src/services/pointFormatRecovery';
import type { PointFormatMutationResult, PointFormatSnapshot } from '../../../src/types/pointFormat';
vi.mock('../../../src/services/datalinkClient', () => ({ api: { get: vi.fn(), put: vi.fn(), post: vi.fn() } }));
const initial: PointFormatSnapshot = { workspace_id: 'W', device_id: 'D', point_id: 'P', rule_id: 'R', workspace_revision: 'source-7', workspace_revision_origin: 'persisted', point_revision: 'point-1', rule_revision: 'rule-2', connection_revision: 'connection-3', editable: true,
  resolution: { resolved: true, format_mode: 'inherit', configured_format: '', effective_format: 'ABCD', effective_source: 'rule', effective_policy: 'inherit_scalar', native_compatibility: false, applicable: true, format_version: 1, workspace_revision: 'source-7', point_revision: 'point-1', rule_revision: 'rule-2', connection_revision: 'connection-3' } };
const operation = createPointFormatOperation(initial, { format_mode: 'override', configured_format: 'CDAB' }, 'original-save');
const saved: PointFormatMutationResult = { point: { id: 'P', device_id: 'D', revision_id: 'point-2', format_mode: 'override', format_override: 'CDAB', format_version: 1, data_format: 'CDAB' },
  resolution: { ...initial.resolution, format_mode: 'override', configured_format: 'CDAB', effective_format: 'CDAB', effective_source: 'override', effective_policy: 'override_scalar', point_revision: 'point-2', workspace_revision: 'source-8' }, workspace_revision: 'source-8', runtime_apply_status: 'pending', invalidation_status: 'invalidated' };
const request = () => createPointFormatRecovery(initial, operation, saved, 'settings-4');
const response = () => ({ outcome: 'recovered', runtime_state: 'not_started', share_applied: false,
  hydration: { state: 'ready', workspace_id: 'W', workspace_revision: 'source-8', settings_revision: 'settings-4', readiness: true, readiness_token: 'do-not-copy-token' } });
beforeEach(() => vi.clearAllMocks());
describe('point format explicit recovery uses the existing backend contract', () => {
  it('captures only original operation ownership, saved result revisions and the selected recovery settings revision', () => {
    expect(request()).toEqual({ workspace_id: 'W', device_id: 'D', point_id: 'P', rule_id: 'R', request_id: 'original-save',
      expected_workspace_revision: 'source-8', expected_point_revision: 'point-2', expected_rule_revision: 'rule-2',
      expected_connection_revision: 'connection-3', expected_settings_revision: 'settings-4' });
    expect(Object.isFrozen(request())).toBe(true);
  });
  it('posts one nested format_recovery only, does not activate or claim a Good reading', async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { success: true, data: response() } });
    const result = await pointFormatRecoveryAPI.recover(request());
    expect(api.post).toHaveBeenCalledExactlyOnceWith('/modbus-share/reconcile', { format_recovery: request() });
    expect(result).toMatchObject({ outcome: 'recovered', runtime_state: 'not_started', share_applied: false });
    expect(JSON.stringify(result)).not.toContain('do-not-copy-token');
    expect(api.get).not.toHaveBeenCalled(); expect(api.put).not.toHaveBeenCalled();
  });
  it.each(['workspace_id', 'workspace_revision', 'settings_revision'])('rejects a successful-looking response with mismatched hydration %s', async field => {
    const value = response(); Object.assign(value.hydration, { [field]: 'foreign' });
    vi.mocked(api.post).mockResolvedValue({ data: { success: true, data: value } });
    await expect(pointFormatRecoveryAPI.recover(request())).rejects.toMatchObject({ outcome: 'unknown', requestId: 'original-save' });
  });
  it.each(['outcome', 'runtime_state', 'share_applied', 'hydration'])('fails closed on malformed recovery field %s', async field => {
    vi.mocked(api.post).mockResolvedValue({ data: { success: true, data: { ...response(), [field]: 'malformed' } } });
    await expect(pointFormatRecoveryAPI.recover(request())).rejects.toMatchObject({ outcome: 'unknown' });
  });
  it('allows explicit same-request replay without fetching or rebasing revisions', async () => {
    const value = { ...response(), outcome: 'replayed', runtime_state: 'installed', share_applied: true };
    vi.mocked(api.post).mockResolvedValue({ data: { success: true, data: value } });
    expect(await pointFormatRecoveryAPI.recover(request())).toMatchObject({ outcome: 'replayed', runtime_state: 'installed', share_applied: true });
    expect(api.get).not.toHaveBeenCalled();
  });
  it('rejects injected ordinary reconcile authority before transport', async () => {
    await expect(pointFormatRecoveryAPI.recover({ ...request(), restore: true } as never)).rejects.toMatchObject({ outcome: 'rejected' });
    expect(api.post).not.toHaveBeenCalled();
  });
  it('rejects changed result identity or source connection rather than using newer evidence', () => {
    expect(() => createPointFormatRecovery(initial, operation, { ...saved, point: { ...saved.point, id: 'foreign' } }, 'settings-4')).toThrow();
    expect(() => createPointFormatRecovery(initial, operation, { ...saved, resolution: { ...saved.resolution, connection_revision: 'changed' } }, 'settings-4')).toThrow();
    expect(() => createPointFormatRecovery(initial, operation, saved, '')).toThrow();
  });
  it('sanitizes a conflict without retrying or leaking transport secrets', async () => {
    vi.mocked(api.post).mockRejectedValue({ config: { password: 'SECRET' }, response: { status: 409, data: { error: { code: 'point_format_conflict', message: 'SECRET' } } } });
    const result = await pointFormatRecoveryAPI.recover(request()).catch((error: unknown) => error);
    expect(result).toMatchObject({ code: 'point_format_conflict', outcome: 'rejected', requestId: 'original-save' });
    expect(JSON.stringify(result)).not.toContain('SECRET'); expect(api.post).toHaveBeenCalledOnce();
  });
});
