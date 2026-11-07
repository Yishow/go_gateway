import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { group } from '../workbench-v2/recording-home-guidance-fixtures';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));
function response(count: number) {
  return { success: true, data: {
    workspace_id: 'W1', workspace_revision: 'W2',
    groups: Array.from({ length: count }, (_, index) => ({ ...group(`G${index}`),
      created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' })),
  } };
}
function serve(value: unknown) {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: value } as never);
  return studioV2WorkspaceWriteGroupsAPI.list();
}
const invalid = { name: 'WriteGroupResponseError', message: 'write-group list response was invalid' };
describe('bounded write-group LIST transport', () => {
  beforeEach(() => vi.clearAllMocks());
  it.each(['object', 'JSON'] as const)('loads all 1000 groups plus four exclusions from a %s envelope', async format => {
    const value = response(1004);
    const result = await serve(format === 'JSON' ? JSON.stringify(value) : value);
    expect(result).toEqual(value.data);
    expect(result.groups).toHaveLength(1004);
    expect(result.groups[1003].id).toBe('G1003');
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
    result.groups[0].members[0].tag_id = 'changed';
    expect(value.data.groups[0].members[0].tag_id).toBe('T1');
  });
  it('accepts 2048 groups but rejects 2049 instead of truncating', async () => {
    expect((await serve(response(2048))).groups).toHaveLength(2048);
    await expect(serve(response(2049))).rejects.toMatchObject(invalid);
  });
  it.each([
    ['duplicate identity', (value: ReturnType<typeof response>) => { value.data.groups[1003].id = 'G0'; }],
    ['foreign workspace', (value: ReturnType<typeof response>) => { value.data.groups[1003].workspace_id = 'other'; }],
    ['nested array overflow', (value: ReturnType<typeof response>) => { Object.assign(value.data.groups[1003], { extra: Array(257).fill(null) }); }],
    ['nested string overflow', (value: ReturnType<typeof response>) => { Object.assign(value.data.groups[1003], { extra: 'x'.repeat(257) }); }],
    ['nested object overflow', (value: ReturnType<typeof response>) => { Object.assign(value.data.groups[1003], { extra: Object.fromEntries(Array.from({ length: 65 }, (_, i) => [`k${i}`, null])) }); }],
    ['nested depth overflow', (value: ReturnType<typeof response>) => { Object.assign(value.data.groups[1003], { extra: { a: { b: { c: { d: { e: true } } } } } }); }],
  ] as const)('rejects %s even in a late list entry', async (_, corrupt) => {
    const value = response(1004);
    corrupt(value);
    await expect(serve(value)).rejects.toMatchObject(invalid);
  });

  it('rejects JSON-encoded groups instead of accepting a second decoding layer', async () => {
    const value = response(1);
    Object.assign(value.data, { groups: [JSON.stringify(value.data.groups[0])] });
    await expect(serve(value)).rejects.toMatchObject(invalid);
  });

  it.each(['object', 'JSON'] as const)('rejects %s responses above the total 8 MiB budget', async format => {
    const value = response(600);
    value.data.groups.forEach(entry => Object.assign(entry, { extra: Array(64).fill('界'.repeat(256)) }));
    await expect(serve(format === 'JSON' ? JSON.stringify(value) : value)).rejects.toMatchObject(invalid);
  });

  it('keeps typed failure metadata and rejects malformed JSON', async () => {
    await expect(serve({ success: false, error: { code: 'revision_mismatch', action: 'reload', retryable: true } }))
      .rejects.toMatchObject({ ...invalid, code: 'revision_mismatch', action: 'reload', retryable: true, outcome: 'failed' });
    await expect(serve('{broken')).rejects.toMatchObject(invalid);
  });

});
