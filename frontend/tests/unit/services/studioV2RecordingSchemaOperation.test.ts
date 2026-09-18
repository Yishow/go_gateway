import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceRecordingPlansAPI } from '@/services/studioV2WorkspaceRecordingPlans';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

const operation = {
  operation_id: 'op-1',
  action: 'schema_apply',
  status: 'succeeded',
  executed_statements: 9,
  verified_digest: 'a'.repeat(64),
  created_at: '2026-09-16T00:50:00Z',
  updated_at: '2026-09-16T00:51:00Z',
  completed_at: '2026-09-16T00:51:00Z',
};

const confirmation = {
  token: 'tok-1',
  operation_id: 'op-1',
  expected_workspace_revision: 'setup-9',
  expected_plan_revision: 'rev-2',
  expected_connector_revision: 'identity-7',
};

describe('studioV2WorkspaceRecordingPlansAPI schema operations', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.get).mockReset();
  });

  it('sends the whole confirmation and returns the recorded operation', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: { success: true, data: operation } });

    await expect(studioV2WorkspaceRecordingPlansAPI.schemaApplyConfirmed(confirmation)).resolves.toEqual(operation);

    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/recording-plans/schema-apply',
      confirmation,
    );
  });

  it('reads one operation by its identity', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: { ...operation, status: 'partial' } } });

    const result = await studioV2WorkspaceRecordingPlansAPI.schemaOperation('op-1');

    expect(result.status).toBe('partial');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/op-1');
  });

  it('refuses a malformed operation instead of reporting a result', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: { success: true, data: { ...operation, status: 'done' } } });

    await expect(studioV2WorkspaceRecordingPlansAPI.schemaApplyConfirmed(confirmation)).rejects.toThrow();
  });
});
