import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import type { RecordingStartRequest } from '@/types/studioV2RecordingStart';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn() } }));

const request: RecordingStartRequest = {
  request_id: 'request-A', workspace_id: 'workspace-A', expected_workspace_revision: 'workspace-rev-A',
  device_ids: ['device-A'],
  groups: [{ group_id: 'group-A', expected_group_revision: 'revision-A', expected_connector_revision: 'connector-A' }],
};

const operation = {
  operation_id: 'operation-A', action: 'recording_start', status: 'succeeded', intent_digest: 'a'.repeat(64),
  workspace_id: 'workspace-A', setup_revision: 'workspace-rev-applied', device_ids: ['device-A'],
  groups: [{ group_id: 'group-A', group_revision: 'revision-A', applied_revision: 'revision-A', saved: true, ready: true, applied: true }],
  devices: [{ device_id: 'device-A', activated: true }], stage: 'complete',
  created_at: '2026-10-04T00:00:00Z', updated_at: '2026-10-04T00:00:01Z',
};

const envelope = (data: unknown) => ({ data: { success: true, data } });

describe('scoped recording start SDK', () => {
  beforeEach(() => vi.clearAllMocks());

  it.each([
    { workspace_id: 'workspace-B' },
    { device_ids: ['device-B'] },
    { groups: [{ ...operation.groups[0], group_id: 'group-B' }] },
  ])('rejects success for a different requested scope %j', async (changed) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({ ...operation, ...changed }) as never);
    await expect(studioV2RecordingStartAPI.start(request)).rejects.toBeDefined();
  });

  it('rejects a lookup response for a different operation', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(operation) as never);
    await expect(studioV2RecordingStartAPI.operation('operation-B')).rejects.toBeDefined();
  });

  it('retains separate saved, applied and activated facts without inventing SQL evidence', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(operation) as never);
    const result = await studioV2RecordingStartAPI.start(request);
    expect(result).toMatchObject({ action: 'recording_start', setup_revision: 'workspace-rev-applied', groups: [{ saved: true, applied: true }], devices: [{ activated: true }] });
    expect(result).not.toHaveProperty('sql_committed');
    expect(result).not.toHaveProperty('verified');
    expect(studioV2DatalinkApi.post).toHaveBeenCalledTimes(1);
  });

  it('accepts an empty setup revision only for a Share-only operation', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      ...operation, setup_revision: '', groups: [],
    }) as never);
    const result = await studioV2RecordingStartAPI.start({ ...request, groups: [] });
    expect(result).toMatchObject({ setup_revision: '', groups: [] });
  });

  it('rejects a missing setup revision for a Share-only operation', async () => {
    const { setup_revision: _setupRevision, ...missingRevision } = operation;
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      ...missingRevision, groups: [],
    }) as never);
    await expect(studioV2RecordingStartAPI.start({ ...request, groups: [] })).rejects.toBeDefined();
  });

  it.each(['schema_apply', 'test_write'])('rejects a different action %s', async (action) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({ ...operation, action }) as never);
    await expect(studioV2RecordingStartAPI.start(request)).rejects.toBeDefined();
  });

  it.each([undefined, '', 'x'.repeat(1025)])('rejects missing or malformed setup progress revision', async (setup_revision) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({ ...operation, setup_revision }) as never);
    await expect(studioV2RecordingStartAPI.start(request)).rejects.toBeDefined();
  });
});
