import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn() } }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({ useWriteGroupsQuery: () => ({ data: { workspace_id: 'W1', workspace_revision: 'WR2', groups: [] }, refetch: vi.fn() }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel', () => ({ BasicRecordingPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection', () => ({ WriteGroupSection: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/Step4SupportPanels', () => ({ Step4SupportPanels: () => null }));
const get = vi.mocked(studioV2DatalinkApi.get);
const operation = {
  operation_id: 'orphan-O1', workspace_id: 'W1', workspace_revision: 'WR1', group_id: 'removed-G1', group_revision: 'old-R1',
  connector_id: 'C1', connector_revision: 'CR1', target_fingerprint: 'a'.repeat(64), action: 'test_write', status: 'running',
  created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z',
};

/** Dispose independent test fixtures explicitly: storage deletion is no longer resolution. */
function resetRecoveryFixture() {
  localStorage.clear();
  for (const ref of recordingTestWriteRecovery.read().references) recordingTestWriteRecovery.resolve({
    operation_id: ref.operation_id, action: 'test_write', status: 'succeeded', write_outcome: 'written_verified',
    cleanup_status: 'cleaned', created_at: '', updated_at: '',
  }, ref.workspace_id);
}

beforeEach(() => { vi.clearAllMocks(); resetRecoveryFixture(); });
it('shows and queries unresolved operations in Step4 even with no current groups, destination or editor', async () => {
  get.mockImplementation(async (path) => ({ data: { success: true, data: path.endsWith('/test-writes')
    ? { workspace_id: 'W1', complete: true, operations: [operation] } : operation } }));
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><Step4Database workspaceId="W1" state={INITIAL_STATE} dispatch={vi.fn()} /></QueryClientProvider>);
  await waitFor(() => expect(screen.getByTestId('workspace-test-write-recovery')).toHaveTextContent('orphan-O1'));
  await act(async () => { localStorage.clear(); window.dispatchEvent(new StorageEvent('storage', { key: null })); });
  expect(screen.getByTestId('workspace-test-write-recovery')).toHaveTextContent('orphan-O1');
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-orphan-O1'));
  await waitFor(() => expect(get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/orphan-O1'));
  expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
});
