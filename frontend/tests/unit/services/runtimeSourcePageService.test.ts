import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const scope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'A' };
const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' };
function page() { return { contract_version: 1, ...scope, kind: 'rules', query: '', limit: 50, list_revision: 'L1', rules: [
  { id: 'A6', device_id: 'A', start_address: '40001', count: 1, data_type: 'float32', naming_prefix: 'Meter', enabled: true, revision_id: 'R6' },
], mappings: [], next_cursor: null }; }
const reply = (source_page: unknown = page(), workspace_id = 'W', success = true) => vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success, data: { workspace_id, source_page } } });
beforeEach(() => vi.mocked(studioV2DatalinkApi.get).mockReset());
describe('bounded runtime source page service', () => {
  it('exposes one page GET with the scoped default limit and abort signal', async () => {
    const signal = new AbortController().signal; reply();
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules' }, signal)).toEqual({ ok: true, value: page() });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...scope, source_list: 'rules', q: '', limit: 50 }, signal });
  });
  it('normalizes query and clamps requested page limit to 100', async () => {
    const source_page = { ...page(), query: 'Meter', limit: 100 }; reply(source_page);
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', query: '  Meter  ', limit: 500 })).toEqual({ ok: true, value: source_page });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...scope, source_list: 'rules', q: 'Meter', limit: 100 }, signal: undefined });
  });
  it('accepts a trimmed safe 128-codepoint query, including supplementary Unicode', async () => {
    const query = '😀'.repeat(128); const source_page = { ...page(), query }; reply(source_page);
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', query: `  ${query}  ` })).toEqual({ ok: true, value: source_page });
  });
  it('sends one explicit filtered mapping page with its cursor and parent revision', async () => {
    const source_page = { ...page(), kind: 'mappings', rules: [], mappings: [{ rule_id: 'A6', rule_revision: 'R6', device_id: 'A', point_id: 'P', address: '40001', enabled: false }] }; reply(source_page);
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'mappings', cursor: 'e30', rule_id: 'A6', rule_revision: 'R6' })).toEqual({ ok: true, value: source_page });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...scope, source_list: 'mappings', q: '', limit: 50, cursor: 'e30', rule_id: 'A6', rule_revision: 'R6' }, signal: undefined });
  });
  it('captures caller fields before waiting, forwards only allowlisted parameters', async () => {
    let finish!: (value: unknown) => void; vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(new Promise((resolve) => { finish = resolve; }));
    const observed = { ...scope, secret: '/private/secret' }; const request = { kind: 'rules' as const, query: '', limit: 50, secret: '/private/secret' };
    const pending = studioV2RuntimeContextAPI.getSourcePage(observed, request);
    observed.device_id = 'B'; request.query = 'later'; request.limit = 1;
    finish({ data: { success: true, data: { workspace_id: 'W', source_page: page() } } });
    expect(await pending).toEqual({ ok: true, value: page() });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...scope, source_list: 'rules', q: '', limit: 50 }, signal: undefined });
  });
  it.each([0, -1, 1.5, NaN, Infinity])('does not transport invalid limit %s', async (limit) => {
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', limit })).toEqual(unavailable);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it.each(['x'.repeat(129), 'password=secret', '/private/secret', 'x\ny'])('does not transport unsafe or overlong query %s', async (query) => {
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', query })).toEqual(unavailable);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it.each(['', 'abc=', 'a', 'A'.repeat(1025)])('does not transport invalid cursor %s', async (cursor) => {
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', cursor })).toEqual(unavailable);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('rejects missing scope, partial filters and rule filters in rules mode', async () => {
    expect(await studioV2RuntimeContextAPI.getSourcePage({ ...scope, workspace_revision: '' }, { kind: 'rules' })).toEqual(unavailable);
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'mappings', rule_id: 'A6' })).toEqual(unavailable);
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', rule_id: 'A6', rule_revision: 'R6' })).toEqual(unavailable);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('rejects foreign outer context and absent or failed page envelopes', async () => {
    reply(page(), 'foreign'); expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules' })).toEqual({ ok: false, reason: 'foreign-scope', action: 'open-workspace-list' });
    reply(null);
    expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules' })).toEqual(unavailable);
    reply(page(), 'W', false); expect(await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules' })).toEqual(unavailable);
  });
  it.each(['stale-reference', 'foreign-scope', 'unavailable', 'driver: secret'])('returns safe typed transport error %s', async (code) => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code, message: '/private/secret' } } } });
    const result = await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules' });
    expect(result).toEqual({ ok: false, reason: ['stale-reference', 'foreign-scope'].includes(code) ? code : 'unavailable', action: code === 'foreign-scope' ? 'open-workspace-list' : 'refresh' });
    expect(JSON.stringify(result)).not.toContain('secret');
  });
});
