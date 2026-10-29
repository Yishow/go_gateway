import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { GroupTestWritePanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel';
import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { listRecordingTestWriteRecovery } from '@/services/recordingTestWriteRecovery';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroupTestWriteOperation, WriteGroupTestWritePreview } from '@/types/studioV2WriteGroupTestWrite';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const list = vi.mocked(listRecordingTestWriteRecovery);
const preview: WriteGroupTestWritePreview = {
  token: 'original-token', operation_id: 'O1', action: 'test_write', group_id: 'G1', group_revision: 'R1',
  expires_at: '2026-10-29T00:00:00Z', target: { connector_id: 'C1', dialect: 'sqlserver', database: '', schema: '', table: 'readings' },
  owner_column: 'owner', owner_value: 'owned-row', dedupe: '', values: [{ column: 'value', type: 'float', value: '1' }], cleanup: 'owned row only',
};
const failedCleanup: WriteGroupTestWriteOperation = { operation_id: 'O1', action: 'test_write', status: 'succeeded',
  write_outcome: 'written_verified', cleanup_status: 'failed', created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z' };
const cleaned: WriteGroupTestWriteOperation = { ...failedCleanup, cleanup_status: 'cleaned' };
const clients: QueryClient[] = [];
type FenceProps = { newOperationBlocked?: boolean; isNewOperationBlocked?: () => boolean; enabled?: boolean; scopeKey?: string };
function mount(initial: FenceProps = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  clients.push(client);
  const content = (props = initial) => <QueryClientProvider client={client}><GroupTestWritePanel workspaceId="W1" groupId="G1"
    scopeKey={props.scopeKey ?? 'G1|R1'} enabled={props.enabled ?? true}
    newOperationBlocked={props.newOperationBlocked} isNewOperationBlocked={props.isNewOperationBlocked} /></QueryClientProvider>;
  const view = render(content());
  return { ...view, rerenderPanel: (props: FenceProps) => view.rerender(content(props)) };
}
async function listReady() {
  await waitFor(() => expect(screen.queryByTestId('group-test-write-list-unavailable')).not.toBeInTheDocument());
}
async function showPreview() {
  await listReady();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  await screen.findByTestId('group-test-write-preview-card');
}
function retain() {
  recordingTestWriteRecovery.retain({ workspace_id: 'W1', group_id: 'G1', group_revision: 'R1', operation_id: 'O1',
    operation_kind: 'test_write', created_at: failedCleanup.created_at });
}
beforeEach(() => {
  vi.resetAllMocks(); localStorage.clear();
  for (const ref of recordingTestWriteRecovery.read().references) {
    recordingTestWriteRecovery.resolve({ ...cleaned, operation_id: ref.operation_id }, ref.workspace_id);
  }
  list.mockResolvedValue({ workspace_id: 'W1', complete: true, operations: [] });
  api.testWritePreview.mockResolvedValue(preview); api.testWrite.mockResolvedValue(failedCleanup);
  api.testWriteOperation.mockResolvedValue(failedCleanup);
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); });

it('disables new test-write preview while the selected schema/start owner reports an unresolved operation', async () => {
  mount({ newOperationBlocked: true }); await listReady();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  expect(api.testWritePreview).not.toHaveBeenCalled(); expect(api.testWrite).not.toHaveBeenCalled();
});
it.each(['retained', 'inspection failure'] as const)('fresh %s denial blocks preview without needing a parent rerender', async mode => {
  let blocked = false;
  const deny = () => { if (blocked && mode === 'inspection failure') throw new Error('private inspection error'); return blocked; };
  mount({ isNewOperationBlocked: deny }); await listReady();
  expect(screen.getByTestId('group-test-write-preview')).toBeEnabled();
  blocked = true; fireEvent.click(screen.getByTestId('group-test-write-preview'));
  await act(async () => {});
  expect(api.testWritePreview).not.toHaveBeenCalled(); expect(api.testWrite).not.toHaveBeenCalled();
  expect(screen.queryByText('private inspection error')).not.toBeInTheDocument();
});
it('fresh denial after preview blocks the initial insert while retaining its original preview', async () => {
  let blocked = false;
  mount({ isNewOperationBlocked: () => blocked }); await showPreview();
  blocked = true; fireEvent.click(screen.getByTestId('group-test-write-confirm'));
  await act(async () => {});
  expect(api.testWritePreview).toHaveBeenCalledOnce(); expect(api.testWrite).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-test-write-preview-card')).toBeInTheDocument();
  expect(recordingTestWriteRecovery.read().references).toEqual([]);
});
it('a newly rendered denial disables an already visible initial confirmation', async () => {
  const view = mount(); await showPreview();
  view.rerenderPanel({ newOperationBlocked: true });
  expect(screen.getByTestId('group-test-write-confirm')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-confirm'));
  expect(api.testWrite).not.toHaveBeenCalled();
});
it('a false callback cannot override an explicit denial snapshot', async () => {
  mount({ newOperationBlocked: true, isNewOperationBlocked: () => false }); await listReady();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  await act(async () => {});
  expect(api.testWritePreview).not.toHaveBeenCalled();
});
it.each(['readonly', 'list failure', 'retained operation'] as const)('a false callback never bypasses the existing %s guard', async mode => {
  if (mode === 'list failure') list.mockRejectedValue(new Error('offline'));
  if (mode === 'retained operation') retain();
  mount({ isNewOperationBlocked: () => false, enabled: mode !== 'readonly' });
  await waitFor(() => expect(list).toHaveBeenCalledOnce());
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  expect(api.testWritePreview).not.toHaveBeenCalled(); expect(api.testWrite).not.toHaveBeenCalled();
});
it('keeps readonly original-operation queries usable under the new-operation fence', async () => {
  retain(); mount({ newOperationBlocked: true, isNewOperationBlocked: () => true });
  await listReady();
  await waitFor(() => expect(api.testWriteOperation).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(screen.getByTestId('group-test-write-recovery-query-O1')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  await waitFor(() => expect(api.testWriteOperation).toHaveBeenCalledTimes(2));
  expect(api.testWriteOperation.mock.calls).toEqual([['O1'], ['O1']]);
  expect(api.testWritePreview).not.toHaveBeenCalled(); expect(api.testWrite).not.toHaveBeenCalled();
});
it('preserves SQLServer same-O1 cleanup repair while fencing new preview/insert operations', async () => {
  const view = mount(); await showPreview();
  fireEvent.click(screen.getByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-cleanup-repair');
  view.rerenderPanel({ newOperationBlocked: true, isNewOperationBlocked: () => true });
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(screen.getByTestId('group-test-write-cleanup-repair')).toBeEnabled();
  api.testWrite.mockResolvedValueOnce(cleaned);
  fireEvent.click(screen.getByTestId('group-test-write-cleanup-repair'));
  await waitFor(() => expect(api.testWrite).toHaveBeenCalledTimes(2));
  expect(api.testWrite).toHaveBeenNthCalledWith(1, 'G1', { token: 'original-token', operation_id: 'O1' });
  expect(api.testWrite).toHaveBeenNthCalledWith(2, 'G1', { token: 'original-token', operation_id: 'O1', action: 'cleanup_repair' });
  expect(api.testWritePreview).toHaveBeenCalledOnce();
  await waitFor(() => expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_cleaned'));
});
it('does not bypass existing readonly cleanup guards even when the new-operation callback allows', async () => {
  const view = mount(); await showPreview();
  fireEvent.click(screen.getByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-cleanup-repair');
  view.rerenderPanel({ enabled: false, isNewOperationBlocked: () => false });
  expect(screen.getByTestId('group-test-write-cleanup-repair')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-cleanup-repair'));
  expect(api.testWrite).toHaveBeenCalledOnce();
});
