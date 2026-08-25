import { describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceActivationAPI } from '../../../src/services/studioV2WorkspaceActivation';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';

describe('Studio V2 workspace activation request', () => {
  it('sends the durable revision and readiness barrier as one request body', async () => {
    const post = vi.spyOn(studioV2DatalinkApi, 'post').mockResolvedValue({
      data: { success: true, data: { workspace_id: 'workspace-1', results: [] } },
    });

    await studioV2WorkspaceActivationAPI.activate({
      workspace_revision: 'workspace-rev-3',
      settings_revision: 'settings-rev-7',
      readiness_token: 'ready-token-1',
      pending_saves: 0,
    });

    expect(post).toHaveBeenCalledWith('/studio-v2/workspace/activate', {
      workspace_revision: 'workspace-rev-3',
      settings_revision: 'settings-rev-7',
      readiness_token: 'ready-token-1',
      pending_saves: 0,
    });
    post.mockRestore();
  });

  it('blocks pending or stale activation before making a network request', async () => {
    const post = vi.spyOn(studioV2DatalinkApi, 'post');

    await expect(studioV2WorkspaceActivationAPI.activate({
      workspace_revision: 'workspace-rev-3',
      settings_revision: 'settings-rev-7',
      readiness_token: 'ready-token-1',
      pending_saves: 1,
    })).rejects.toMatchObject({ code: 'modbus_share_save_incomplete' });
    await expect(studioV2WorkspaceActivationAPI.activate({
      workspace_revision: 'workspace-rev-3',
      settings_revision: 'settings-rev-7',
      readiness_token: 'ready-token-1',
      stale: true,
    })).rejects.toMatchObject({ code: 'modbus_share_revision_conflict' });
    expect(post).not.toHaveBeenCalled();
    post.mockRestore();
  });
});
