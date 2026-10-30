import { describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { delivery, group } from './recording-home-guidance-fixtures';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
describe('home read transport cancellation', () => {
  it('forwards the Query-owned signal without changing legacy service signatures', async () => {
    const signal = new AbortController().signal; const value = group('G');
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: delivery(value) } });
    await studioV2WorkspaceWriteGroupsAPI.delivery('G', signal);
    expect(studioV2DatalinkApi.get).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/G/delivery', { signal });
    await studioV2WorkspaceWriteGroupsAPI.delivery('G');
    expect(studioV2DatalinkApi.get).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/G/delivery');
    const lifecycle = savedLifecycle(value, 'W2');
    lifecycle.actions = lifecycle.actions.map(action => action.action === 'recreate' || action.action === 'transition'
      ? { ...action, allowed: false, reason_code: 'not_available' } : action);
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: lifecycle } });
    await studioV2WriteGroupLifecycleAPI.lifecycle('G', signal);
    expect(studioV2DatalinkApi.get).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/G/lifecycle', { signal });
    await studioV2WriteGroupLifecycleAPI.lifecycle('G');
    expect(studioV2DatalinkApi.get).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/G/lifecycle');
  });
});
