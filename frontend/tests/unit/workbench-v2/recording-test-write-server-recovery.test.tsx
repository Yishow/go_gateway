import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, expect, it, vi } from 'vitest';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { GroupTestWritePanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn() } }));
const get = vi.mocked(studioV2DatalinkApi.get);
const post = vi.mocked(studioV2DatalinkApi.post);
const operation = {
  operation_id: 'O1', workspace_id: 'W1', workspace_revision: 'WR1', group_id: 'G1', group_revision: 'R1',
  connector_id: 'C1', connector_revision: 'CR1', target_fingerprint: 'a'.repeat(64), action: 'test_write',
  status: 'running', created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z',
};
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const content = (groupId = 'G1', revision = 'R2', workspaceId = 'W1') => <QueryClientProvider client={client}>
    <GroupTestWritePanel workspaceId={workspaceId} groupId={groupId} scopeKey={`${groupId}|${revision}`} enabled />
  </QueryClientProvider>;
  return { ...render(content()), content, client };
}

/** Dispose independent test fixtures explicitly: storage deletion is no longer resolution. */
function resetRecoveryFixture() {
  localStorage.clear();
  for (const ref of recordingTestWriteRecovery.read().references) recordingTestWriteRecovery.resolve({
    operation_id: ref.operation_id, action: 'test_write', status: 'succeeded', write_outcome: 'written_verified',
    cleanup_status: 'cleaned', created_at: '', updated_at: '',
  }, ref.workspace_id);
}

beforeEach(() => { vi.clearAllMocks(); resetRecoveryFixture(); });
it('rediscovers O1 from the production listing after browser references were lost', async () => {
  get.mockImplementation(async (path) => ({ data: { success: true, data: path.endsWith('/test-writes')
    ? { workspace_id: 'W1', complete: true, operations: [operation] } : operation } }));
  mount();
  expect(await screen.findByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  expect(get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/test-writes');
  expect(get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/O1');
  expect(post).not.toHaveBeenCalled();
});
it.each([
  ['foreign workspace', { workspace_id: 'W2', complete: true, operations: [operation] }],
  ['incomplete', { workspace_id: 'W1', complete: false, operations: [] }],
  ['missing completion', { workspace_id: 'W1', operations: [] }],
  ['foreign entry', { workspace_id: 'W1', complete: true, operations: [{ ...operation, workspace_id: 'W2' }] }],
  ['duplicate', { workspace_id: 'W1', complete: true, operations: [operation, operation] }],
  ['invalid fingerprint', { workspace_id: 'W1', complete: true, operations: [{ ...operation, target_fingerprint: 'bad' }] }],
  ['invalid timestamp', { workspace_id: 'W1', complete: true, operations: [{ ...operation, created_at: 'yesterday' }] }],
  ['invalid outcome', { workspace_id: 'W1', complete: true, operations: [{ ...operation, write_outcome: 'written_verified' }] }],
])('fails closed on %s listing', async (_label, data) => {
  get.mockResolvedValue({ data: { success: true, data } }); mount();
  await waitFor(() => expect(screen.getByTestId('group-test-write-list-retry')).toBeEnabled());
  expect(screen.getByTestId('group-test-write-list-unavailable')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(post).not.toHaveBeenCalled();
});
it('fails closed on list transport errors without displaying raw diagnostics', async () => {
  get.mockRejectedValue(new Error('private endpoint / credential')); mount();
  await waitFor(() => expect(screen.getByTestId('group-test-write-list-retry')).toBeEnabled());
  expect(screen.getByTestId('group-test-write-list-unavailable')).toBeInTheDocument();
  expect(screen.queryByText(/private endpoint/)).not.toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
});
it('accepts the repository cleanup claim as unresolved, without tokenless repair', async () => {
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [{
    ...operation, status: 'running', write_outcome: 'written_verified', cleanup_status: 'failed',
    completed_at: '2026-10-28T00:00:01Z', updated_at: '2026-10-28T00:00:02Z',
  }] } } }); mount();
  expect(await screen.findByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  expect(post).not.toHaveBeenCalled();
});
it('retains a local in-flight request absent from a complete empty list', async () => {
  localStorage.setItem('studio.recording-test-write-recovery.v1', JSON.stringify([{
    workspace_id: 'W1', group_id: 'G1', group_revision: 'R1', operation_id: 'O-local', operation_kind: 'test_write', created_at: operation.created_at,
  }]));
  get.mockImplementation(async (path) => {
    if (!path.endsWith('/test-writes')) throw new Error('not visible yet');
    return { data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [] } } };
  }); mount();
  await waitFor(() => expect(screen.queryByTestId('group-test-write-list-unavailable')).not.toBeInTheDocument());
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O-local');
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(post).not.toHaveBeenCalled();
});
it('does not treat cached empty success as proof of no write after refresh fails', async () => {
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [] } } });
  const view = mount();
  await waitFor(() => expect(screen.getByTestId('group-test-write-preview')).toBeEnabled());
  get.mockRejectedValue(new Error('unavailable'));
  await act(async () => { await view.client.refetchQueries(); });
  await waitFor(() => expect(screen.getByTestId('group-test-write-list-retry')).toBeEnabled());
  expect(screen.getByTestId('group-test-write-list-unavailable')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(post).not.toHaveBeenCalled();
});
it('keeps historical revisions and different groups separately, without another preview or write', async () => {
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [
    operation, { ...operation, operation_id: 'O2', group_revision: 'R2' }, { ...operation, operation_id: 'O3', group_id: 'G2' },
  ] } } });
  const view = mount();
  expect(await screen.findByTestId('group-test-write-recovery')).toHaveTextContent('O3');
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('R1');
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('R2');
  view.rerender(view.content('G2')); view.rerender(view.content('G1'));
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(post).not.toHaveBeenCalled();
});
it('rejects a server identity conflicting with a retained local group revision', async () => {
  localStorage.setItem('studio.recording-test-write-recovery.v1', JSON.stringify([{
    workspace_id: 'W1', group_id: 'G1', group_revision: 'R-other', operation_id: 'O1', operation_kind: 'test_write', created_at: operation.created_at,
  }]));
  get.mockImplementation(async (path) => ({ data: { success: true, data: path.endsWith('/test-writes')
    ? { workspace_id: 'W1', complete: true, operations: [operation] } : operation } })); mount();
  await waitFor(() => expect(screen.getByTestId('group-test-write-list-retry')).toBeEnabled());
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(localStorage.getItem('studio.recording-test-write-recovery.v1')).toContain('R-other');
});
it('rejects immutable connector changes for an operation across list refreshes', async () => {
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [operation] } } });
  const view = mount(); await screen.findByTestId('group-test-write-recovery');
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [{ ...operation, connector_revision: 'CR-OTHER' }] } } });
  await act(async () => { await view.client.refetchQueries(); });
  await waitFor(() => expect(screen.getByTestId('group-test-write-list-retry')).toBeEnabled());
  expect(screen.getByTestId('group-test-write-list-unavailable')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
});
it('recovers the same server operation after lost ACK, reload and local-reference loss without resending', async () => {
  let listed = false;
  get.mockImplementation(async (path) => ({ data: { success: true, data: path.endsWith('/test-writes')
    ? { workspace_id: 'W1', complete: true, operations: listed ? [operation] : [] } : operation } }));
  post.mockImplementation(async (path) => {
    if (path.endsWith('/test-write')) { listed = true; throw new Error('lost ACK'); }
    return { data: { success: true, data: {
      token: 'SECRET', operation_id: 'O1', action: 'test_write', group_id: 'G1', group_revision: 'R1',
      expires_at: '2026-10-28T00:05:00Z', target: { connector_id: 'C1', dialect: 'sqlite', database: 'private', schema: '', table: 'private' },
      owner_column: 'owner', owner_value: 'owned', dedupe: 'receipt', values: [{ column: 'owner', type: 'text', value: 'owned' }], cleanup: 'owned row only',
    } } };
  });
  const first = mount();
  await waitFor(() => expect(screen.getByTestId('group-test-write-preview')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-unconfirmed');
  first.unmount(); localStorage.clear(); mount();
  expect(await screen.findByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  await waitFor(() => expect(get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/O1'));
  expect(post).toHaveBeenCalledTimes(2);
  expect(post).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/G1/test-write', { token: 'SECRET', operation_id: 'O1', action: undefined });
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
});
it('never exposes a late W1 listing while W2 is current, then restores W1 recovery on return', async () => {
  let resolveFirst!: (value: unknown) => void;
  const pending = new Promise((resolve) => { resolveFirst = resolve; });
  get.mockImplementationOnce(() => pending as ReturnType<typeof get>);
  const view = mount();
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W2', complete: true, operations: [] } } });
  view.rerender(view.content('G2', 'R1', 'W2'));
  await waitFor(() => expect(screen.getByTestId('group-test-write-preview')).toBeEnabled());
  await act(async () => { resolveFirst({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [operation] } } }); });
  expect(screen.queryByTestId('group-test-write-recovery')).not.toBeInTheDocument();
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [operation] } } });
  view.rerender(view.content());
  expect(await screen.findByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  expect(post).not.toHaveBeenCalled();
});
it('keeps a locally pending operation visible through a mounted storage-deletion event and cached empty list', async () => {
  localStorage.setItem('studio.recording-test-write-recovery.v1', JSON.stringify([{
    workspace_id: 'W1', group_id: 'G1', group_revision: 'R1', operation_id: 'O-local', operation_kind: 'test_write', created_at: operation.created_at,
  }]));
  get.mockImplementation(async (path) => {
    if (!path.endsWith('/test-writes')) throw new Error('request not yet visible');
    return { data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [] } } };
  }); mount();
  await waitFor(() => expect(screen.queryByTestId('group-test-write-list-unavailable')).not.toBeInTheDocument());
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O-local');
  const reads = get.mock.calls.length;
  await act(async () => { localStorage.clear(); window.dispatchEvent(new StorageEvent('storage', { key: null })); });
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O-local');
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(get).toHaveBeenCalledTimes(reads);
  expect(post).not.toHaveBeenCalled();
});
// Last: deliberately corrupt storage retains a memory-only safety locator.
it('still displays server discovery with corrupt browser storage, without overwriting it or allowing writes', async () => {
  localStorage.setItem('studio.recording-test-write-recovery.v1', '{corrupt');
  get.mockResolvedValue({ data: { success: true, data: { workspace_id: 'W1', complete: true, operations: [operation] } } }); mount();
  expect(await screen.findByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  expect(screen.getByTestId('group-test-write-storage-error')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(localStorage.getItem('studio.recording-test-write-recovery.v1')).toBe('{corrupt');
  expect(post).not.toHaveBeenCalled();
});
