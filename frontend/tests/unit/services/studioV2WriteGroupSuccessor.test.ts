import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WriteGroupSuccessorAPI as api } from '@/services/studioV2WriteGroupSuccessor';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { lifecycle, operationStatus } from '../../fixtures/writeGroupSuccessor';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { post: vi.fn(), get: vi.fn() } }));
const request = { workspace_id: 'ws-1', lineage_id: 'lineage-1', operation_id: 'root-op-1' };
beforeEach(() => vi.resetAllMocks());
describe('successor exact scope and safe failure', () => {
  it('posts the original root scope, never a latest lookup', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: operationStatus() } });
    expect(await api.operationStatus('wide-1', request)).toEqual(operationStatus());
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith('/studio-v2/workspace/write-groups/wide-1/successor', { action: 'operation-status', ...request });
  });
  it('rejects a valid but different original operation as unknown', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: operationStatus() } });
    await expect(api.operationStatus('wide-1', { ...request, operation_id: 'other-root' })).rejects.toMatchObject({ outcome: 'unconfirmed' });
  });
  it('rejects a foreign lifecycle identity', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: lifecycle() } });
    await expect(studioV2WriteGroupLifecycleAPI.lifecycle('other-group')).rejects.toMatchObject({ outcome: 'unconfirmed' });
  });
  it.each([200, 408, 502])('never returns raw response or transport details for HTTP %s', async (status) => {
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValue({ message: 'postgres://secret-token', config: { secret: 'secret-token' }, response: { status, data: '<html>secret-token</html>' } });
    try { await api.operationStatus('wide-1', request); expect.fail('expected rejection'); }
    catch (error) { expect(error).toMatchObject({ outcome: 'unconfirmed' }); expect(JSON.stringify(error)).not.toContain('secret-token'); }
  });
  it('keeps a typed HTTP timeout unknown but safely classifies a typed validation rejection', async () => {
    const data = { success: false, error: { code: 'WRITE_GROUP_INVALID_REQUEST', message: 'private DSN postgres://token' } };
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValueOnce({ response: { status: 408, data } });
    await expect(api.operationStatus('wide-1', request)).rejects.toMatchObject({ outcome: 'unconfirmed' });
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValueOnce({ response: { status: 400, data } });
    await expect(api.operationStatus('wide-1', request)).rejects.toMatchObject({ outcome: 'failed', code: 'WRITE_GROUP_INVALID_REQUEST' });
  });
  it('treats malformed nominal success as unknown, never completed', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: { status: 'succeeded', rawDetails: 'secret-token' } } });
    await expect(api.operationStatus('wide-1', request)).rejects.toMatchObject({ outcome: 'unconfirmed' });
  });
});
