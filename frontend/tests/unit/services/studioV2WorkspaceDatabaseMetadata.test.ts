import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

const envelope = {
  workspace_id: 'ws-1', connector_id: 'db-1', connector_revision: 'identity-1', database: '/tmp/line-a.db',
  schema: 'main', table: 'sensor_values', inspection_status: 'exists',
  columns: [{ name: 'recorded_at', data_type: 'TEXT', nullable: false, primary_key: true }],
};

describe('studio V2 workspace database metadata service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.get).mockReset();
  });

  it('requests the saved target metadata with the expected connector revision', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: envelope } } as never);

    await expect(studioV2WorkspaceDatabaseAPI.getMetadata('identity-1')).resolves.toEqual(envelope);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-metadata', {
      params: { expected_connector_revision: 'identity-1' },
    });
  });

  it('rejects an unknown inspection status, a missing column list or a missing envelope', async () => {
    for (const data of [{ ...envelope, inspection_status: 'maybe' }, { ...envelope, columns: undefined }, null]) {
      vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data } } as never);
      await expect(studioV2WorkspaceDatabaseAPI.getMetadata('identity-1')).rejects.toThrow();
    }
  });

  it('never returns columns for a table that was not confirmed to exist', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({
      data: { success: true, data: { ...envelope, inspection_status: 'forbidden', reason: 'permission_denied' } },
    } as never);

    await expect(studioV2WorkspaceDatabaseAPI.getMetadata('identity-1')).resolves.toEqual(
      expect.objectContaining({ inspection_status: 'forbidden', reason: 'permission_denied', columns: [] }),
    );
  });
});
