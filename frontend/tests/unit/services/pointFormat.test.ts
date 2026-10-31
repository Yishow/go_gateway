import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '@/services/datalinkClient';
import { createPointFormatOperation, pointFormatAPI } from '@/services/pointFormat';
import type { PointFormatSnapshot } from '@/types/pointFormat';
vi.mock('@/services/datalinkClient', () => ({ api: { get: vi.fn(), put: vi.fn() } }));
const scope = { workspace_id: 'w', device_id: 'd', point_id: 'p' };
const snapshot = (): PointFormatSnapshot => ({ ...scope, rule_id: 'r', workspace_revision: 'share-7', workspace_revision_origin: 'persisted', rule_revision: 'rule-2', point_revision: '', connection_revision: 'conn-4', editable: true, resolution: { resolved: true, format_mode: 'legacy', configured_format: '', effective_format: 'CDAB', effective_source: 'connection', effective_policy: 'legacy_mc3e', native_compatibility: false, applicable: true, format_version: 0, point_revision: '', workspace_revision: 'share-7', rule_revision: 'rule-2', connection_revision: 'conn-4' } });
const operation = () => createPointFormatOperation(snapshot(), { format_mode: 'override', configured_format: 'DCBA' }, 'edit-1');
const result = () => ({ point: { id: 'p', device_id: 'd', revision_id: 'p-2', format_mode: 'override', format_override: 'DCBA', format_version: 1, data_format: 'DCBA' }, resolution: { ...snapshot().resolution, point_revision: 'p-2', format_mode: 'override', configured_format: 'DCBA', effective_format: 'DCBA', effective_source: 'override', effective_policy: 'override_scalar', format_version: 1 }, workspace_revision: 'share-7', runtime_apply_status: 'pending', invalidation_status: 'invalidated' });
const reply = (data: unknown) => ({ data: { success: true, data } });
beforeEach(() => vi.clearAllMocks());
describe('point format canonical operation', () => {
  it('reads source revisions and preserves explicitly observed empty legacy revision', async () => {
    vi.mocked(api.get).mockResolvedValueOnce(reply(snapshot()));
    expect(await pointFormatAPI.read(scope)).toEqual(snapshot());
  });
  it('freezes initial revisions and sends one PUT without a latest read or rebase', async () => {
    const initial = snapshot(); const op = createPointFormatOperation(initial, { format_mode: 'override', configured_format: 'DCBA' }, 'edit-1');
    initial.workspace_revision = 'DATABASE-SETUP-NOT-SOURCE';
    vi.mocked(api.put).mockResolvedValueOnce(reply(result()));
    expect(await pointFormatAPI.save(op)).toMatchObject({ runtime_apply_status: 'pending' });
    expect(api.get).not.toHaveBeenCalled();
    expect(api.put).toHaveBeenCalledExactlyOnceWith('/points/p/format', expect.objectContaining({ expected_workspace_revision: 'share-7', expected_point_revision: '', request_id: 'edit-1' }));
    expect(Object.isFrozen(op)).toBe(true);
  });
  it('reset explicitly sends inheritance with empty configured format', () => {
    expect(createPointFormatOperation(snapshot(), { format_mode: 'inherit', configured_format: '' }, 'reset-1')).toMatchObject({ format_mode: 'inherit', configured_format: '' });
  });
  it('does not allow missing revision, unsupported override, or noneditable scope', () => {
    const missing = snapshot(); delete (missing as Partial<PointFormatSnapshot>).point_revision;
    expect(() => createPointFormatOperation(missing, { format_mode: 'inherit', configured_format: '' }, 'x')).toThrow();
    const unsupported = snapshot(); unsupported.resolution.applicable = false;
    expect(() => createPointFormatOperation(unsupported, { format_mode: 'override', configured_format: 'ABCD' }, 'x')).toThrow();
    unsupported.editable = false;
    expect(() => createPointFormatOperation(unsupported, { format_mode: 'inherit', configured_format: '' }, 'x')).toThrow();
  });
  it('sanitizes conflict and invalid errors without raw backend secret or config', async () => {
    for (const [status, code] of [[409, 'point_format_conflict'], [422, 'point_format_invalid']] as const) {
      vi.mocked(api.put).mockRejectedValueOnce({ message: 'SECRET', config: { password: 'SECRET' }, response: { status, data: { success: false, error: { code, message: 'SECRET' } } } });
      const error = await pointFormatAPI.save(operation()).catch((e: unknown) => e);
      expect(error).toMatchObject({ code, outcome: 'rejected' });
      expect(JSON.stringify(error)).not.toContain('SECRET');
    }
  });
  it('unknown outcome is inspected with the original request ID, never resent', async () => {
    vi.mocked(api.put).mockRejectedValueOnce(new Error('SECRET timeout'));
    const op = operation();
    await expect(pointFormatAPI.save(op)).rejects.toMatchObject({ outcome: 'unknown' });
    vi.mocked(api.get).mockResolvedValueOnce(reply(result()));
    await pointFormatAPI.inspect(op);
    expect(api.put).toHaveBeenCalledTimes(1);
    expect(api.get).toHaveBeenCalledWith('/points/p/format/operations/edit-1', { params: { workspace_id: 'w', device_id: 'd', rule_id: 'r' } });
  });
  it.each(['foreign-point', 'wrong-intent', 'mixed-revision', 'false-success'])('rejects malformed success %s as unknown', async (kind) => {
    const data = result();
    if (kind === 'foreign-point') data.point.id = 'foreign';
    if (kind === 'wrong-intent') data.resolution.configured_format = 'ABCD';
    if (kind === 'mixed-revision') data.resolution.point_revision = 'wrong';
    vi.mocked(api.put).mockResolvedValueOnce(kind === 'false-success' ? { data: { success: false, data } } : reply(data));
    await expect(pointFormatAPI.save(operation())).rejects.toMatchObject({ outcome: 'unknown' });
  });
  it('permits native inherited empty effective format without inventing a format', async () => {
    const data = snapshot(); Object.assign(data.resolution, { format_mode: 'inherit', format_version: 1, effective_format: '', effective_source: 'native', effective_policy: 'legacy_fatek_components', native_compatibility: true, applicable: false });
    vi.mocked(api.get).mockResolvedValueOnce(reply(data));
    expect((await pointFormatAPI.read(scope)).resolution.effective_format).toBe('');
  });
});

describe('safe operation recovery evidence', () => {
  it('retains only the original request ID on malformed-success errors', async () => {
    vi.mocked(api.put).mockResolvedValueOnce(reply({}));
    await expect(pointFormatAPI.save(operation())).rejects.toMatchObject({ outcome: 'unknown', requestId: 'edit-1' });
  });
  it.each(['point_revision', 'connection_revision', 'rule_revision'])('rejects missing canonical %s even when empty is legal', async (key) => {
    const data: Record<string, unknown> = { ...snapshot() }; delete data[key];
    vi.mocked(api.get).mockResolvedValueOnce(reply(data));
    await expect(pointFormatAPI.read(scope)).rejects.toMatchObject({ outcome: 'unknown' });
  });
  it('rejects a foreign requested rule and strips extra secret fields', async () => {
    vi.mocked(api.get).mockResolvedValueOnce(reply({ ...snapshot(), connection_config: 'SECRET' }));
    expect(JSON.stringify(await pointFormatAPI.read(scope))).not.toContain('SECRET');
    vi.mocked(api.get).mockResolvedValueOnce(reply(snapshot()));
    await expect(pointFormatAPI.read({ ...scope, rule_id: 'other' })).rejects.toMatchObject({ outcome: 'unknown' });
  });
  it('does not permit patches to inject format or scope authority', () => {
    expect(() => createPointFormatOperation(snapshot(), { format_mode: 'inherit', configured_format: '' }, 'x', { data_format: 'ABCD' } as never)).toThrow();
  });
  it('preserves pending and unavailable as unknown rather than saved/applied', async () => {
    vi.mocked(api.put).mockRejectedValueOnce({ response: { status: 503, data: { success: false, error: { code: 'point_format_unavailable', message: 'SECRET' } } } });
    await expect(pointFormatAPI.save(operation())).rejects.toMatchObject({ code: 'point_format_unavailable', outcome: 'unknown' });
  });
});

describe('versioned evidence and operation identity', () => {
  it.each(['format_version', 'effective_policy'])('fails closed on unknown %s', async (field) => {
    const data = snapshot(); Object.assign(data.resolution, { [field]: field === 'format_version' ? 2 : 'future_policy' });
    vi.mocked(api.get).mockResolvedValueOnce(reply(data));
    await expect(pointFormatAPI.read(scope)).rejects.toMatchObject({ outcome: 'unknown' });
  });
  it.each(['', '../foreign', 'x'.repeat(129)])('rejects unsafe operation ID %s before transport', (id) => {
    expect(() => createPointFormatOperation(snapshot(), { format_mode: 'inherit', configured_format: '' }, id)).toThrow();
    expect(api.put).not.toHaveBeenCalled();
  });
  it('deep freezes a detached ordinary patch', () => {
    const patch = { name: 'before' };
    const op = createPointFormatOperation(snapshot(), { format_mode: 'inherit', configured_format: '' }, 'edit-patch', patch);
    patch.name = 'after';
    expect(op.patch?.name).toBe('before');
    expect(Object.isFrozen(op.patch)).toBe(true);
  });
  it('returns a safe inspect-not-found without another write', async () => {
    vi.mocked(api.get).mockRejectedValueOnce({ response: { status: 404, data: { success: false, error: { code: 'point_format_operation_not_found' } } } });
    await expect(pointFormatAPI.inspect(operation())).rejects.toMatchObject({ code: 'point_format_operation_not_found', outcome: 'unknown', requestId: 'edit-1' });
    expect(api.put).not.toHaveBeenCalled();
  });
});

 it('rejects unresolved read evidence instead of allowing an inferred default', async () => {
   const data = snapshot(); data.resolution.resolved = false;
   vi.mocked(api.get).mockResolvedValueOnce(reply(data));
   await expect(pointFormatAPI.read(scope)).rejects.toMatchObject({ outcome: 'unknown' });
 });
 it('rejects a mutation result bound to a different source-rule revision', async () => {
   const data = result(); data.resolution.rule_revision = 'foreign-rule-revision';
   vi.mocked(api.put).mockResolvedValueOnce(reply(data));
   await expect(pointFormatAPI.save(operation())).rejects.toMatchObject({ outcome: 'unknown' });
 });
 it.each([['legacy', 1], ['inherit', 0]] as const)('accepts server-resolved %s version %s without a client policy inference', async (format_mode, format_version) => {
   const data = snapshot(); Object.assign(data.resolution, { format_mode, format_version, effective_format: '', effective_source: 'native', effective_policy: 'legacy_fatek_components', native_compatibility: true });
   vi.mocked(api.get).mockResolvedValueOnce(reply(data));
   expect((await pointFormatAPI.read(scope)).resolution).toEqual(data.resolution);
 });
