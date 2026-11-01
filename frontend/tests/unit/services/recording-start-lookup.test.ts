import { beforeEach, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { studioV2RecordingStartAPI } from '../../../src/services/studioV2RecordingStart';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn() } }));
const request = { request_id: 'original', workspace_id: 'W', expected_workspace_revision: 'database-head', workspace_revision: 'share-head', settings_revision: 'settings-head', readiness_token: 'original-proof', device_ids: ['A', 'B'], groups: [{ group_id: 'G', expected_group_revision: 'g1', expected_connector_revision: 'c1' }] };
const operation = { operation_id: 'O', action: 'recording_start', status: 'unknown', intent_digest: 'a'.repeat(64), workspace_id: 'W', setup_revision: 'database-later', device_ids: ['B', 'A'],
 groups: [{ group_id: 'G', group_revision: 'g1', saved: true, applied: false, ready: true }], devices: [{ device_id: 'A', activated: false }, { device_id: 'B', activated: false }], stage: 'activation', created_at: '2026-10-31T00:00:00Z', updated_at: '2026-10-31T00:00:01Z' };
beforeEach(() => vi.clearAllMocks());
it('looks up the unchanged original request at the read-only facade with separate revision domains', async () => {
 vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: operation } });
 expect(await studioV2RecordingStartAPI.lookup(request)).toMatchObject({ operation_id: 'O', status: 'unknown' });
 expect(studioV2DatalinkApi.post).toHaveBeenCalledExactlyOnceWith('/studio-v2/workspace/recording-start/lookup', request);
 expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
});
it.each([{ workspace_id: 'foreign' }, { device_ids: ['A'] }, { groups: [{ ...operation.groups[0], group_id: 'foreign' }] }, { action: 'schema_apply' }])('rejects a mismatched original operation %j', async changed => {
 vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: { ...operation, ...changed } } });
 await expect(studioV2RecordingStartAPI.lookup(request)).rejects.toBeDefined();
 expect(studioV2DatalinkApi.post).toHaveBeenCalledTimes(1);
});
