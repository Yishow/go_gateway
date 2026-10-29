import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
vi.mock('@/hooks/datalink/useRecordingTestWriteRecovery', () => ({ useRecordingTestWriteRecovery: () => ({ available: true, isFetching: false, refetch: vi.fn() }) }));
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, expect, it, vi } from 'vitest';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { GroupTestWritePanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const preview = {
  token: 'SECRET', operation_id: 'O1', action: 'test_write' as const, group_id: 'G1', group_revision: 'R1',
  expires_at: '', target: { connector_id: 'C1', dialect: 'mysql', database: 'private-db', schema: '', table: 'private-table' },
  owner_column: '', owner_value: '', dedupe: '', values: [], cleanup: '',
};
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const content = (group = 'G1', revision = 'R1') => <QueryClientProvider client={client}>
    <GroupTestWritePanel workspaceId="W1" groupId={group} scopeKey={`${group}|${revision}`} enabled />
  </QueryClientProvider>;
  const view = render(content());
  return { ...view, content };
}

/** Dispose independent test fixtures explicitly: storage deletion is no longer resolution. */
function resetRecoveryFixture() {
  localStorage.clear();
  for (const ref of recordingTestWriteRecovery.read().references) recordingTestWriteRecovery.resolve({
    operation_id: ref.operation_id, action: 'test_write', status: 'succeeded', write_outcome: 'written_verified',
    cleanup_status: 'cleaned', created_at: '', updated_at: '',
  }, ref.workspace_id);
}

beforeEach(() => {
  vi.clearAllMocks(); resetRecoveryFixture();
  api.testWritePreview.mockResolvedValue(preview);
  api.testWrite.mockRejectedValue(new Error('lost ACK'));
});
async function loseReply() {
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-unconfirmed');
}
it('does not replace O1 after a lost response with a new G1 preview', async () => {
  mount(); await loseReply();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  expect(api.testWritePreview).toHaveBeenCalledTimes(1);
});
it('retains O1 across G1/R1 to G1/R2 and A to B to A navigation', async () => {
  const view = mount(); await loseReply();
  view.rerender(view.content('G1', 'R2'));
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  view.rerender(view.content('G2'));
  expect(screen.getByTestId('group-test-write-preview')).toBeEnabled();
  view.rerender(view.content());
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
});
it.each(['pending', 'running', 'unknown'] as const)('keeps %s queryable without another insert', async (status) => {
  api.testWrite.mockResolvedValue({ operation_id: 'O1', action: 'test_write', status, created_at: '', updated_at: '' });
  mount();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-result');
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(screen.getByTestId('group-test-write-check')).toBeInTheDocument();
  expect(api.testWrite).toHaveBeenCalledTimes(1);
});
it('reload queries exactly O1 and an absent response does not unlock another write', async () => {
  const first = mount(); await loseReply(); first.unmount();
  api.testWriteOperation.mockRejectedValue(new Error('absent or inaccessible'));
  mount();
  await screen.findByTestId('group-test-write-check-unconfirmed');
  expect(api.testWriteOperation).toHaveBeenCalledTimes(1);
  expect(api.testWriteOperation).toHaveBeenCalledWith('O1');
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(api.testWrite).toHaveBeenCalledTimes(1);
  const stored = localStorage.getItem('studio.recording-test-write-recovery.v1');
  expect(stored).toContain('O1');
  expect(stored).not.toMatch(/SECRET|private-db|private-table|token|values/);
});
it('retains multiple groups and provides a query entry while viewing another group', async () => {
  const view = mount(); await loseReply();
  view.rerender(view.content('G2'));
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  expect(screen.getByTestId('group-test-write-recovery-query-O1')).toBeInTheDocument();
  api.testWritePreview.mockResolvedValue({ ...preview, operation_id: 'O2', group_id: 'G2' });
  await loseReply();
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O2');
  expect(api.testWrite).toHaveBeenCalledTimes(2);
});
it('ignores an old completion after A to B to A, then explicitly queries O1', async () => {
  let finish!: (value: Awaited<ReturnType<typeof api.testWrite>>) => void;
  api.testWrite.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  api.testWriteOperation.mockRejectedValue(new Error('offline'));
  const view = mount();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-unconfirmed');
  view.rerender(view.content('G2')); view.rerender(view.content());
  await act(async () => { finish({ operation_id: 'O1', action: 'test_write', status: 'running', created_at: '', updated_at: '' }); });
  expect(screen.queryByTestId('group-test-write-result')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(api.testWriteOperation).toHaveBeenCalledWith('O1');
  expect(api.testWrite).toHaveBeenCalledTimes(1);
});
it('blocks dispatch with corrupt storage and reports recovery unavailable', () => {
  localStorage.setItem('studio.recording-test-write-recovery.v1', '{broken'); mount();
  expect(screen.getByTestId('group-test-write-storage-error')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(api.testWrite).not.toHaveBeenCalled();
});
it('does not expose foreign workspace references', () => {
  localStorage.setItem('studio.recording-test-write-recovery.v1', JSON.stringify([{
    workspace_id: 'W2', group_id: 'FOREIGN', group_revision: 'R1', operation_id: 'FOREIGN-O',
    operation_kind: 'test_write', created_at: '2026-10-28T00:00:00Z',
  }]));
  mount();
  expect(screen.queryByText(/FOREIGN/)).not.toBeInTheDocument();
  expect(api.testWriteOperation).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-test-write-preview')).toBeEnabled();
});
it('shows an honest query failure while a running operation remains unresolved', async () => {
  api.testWrite.mockResolvedValue({ operation_id: 'O1', action: 'test_write', status: 'running', created_at: '', updated_at: '' });
  api.testWriteOperation.mockRejectedValue(new Error('private endpoint details'));
  mount();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  fireEvent.click(await screen.findByTestId('group-test-write-check'));
  expect(await screen.findByTestId('group-test-write-unconfirmed')).toBeInTheDocument();
  expect(screen.queryByText(/private endpoint/)).not.toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(api.testWrite).toHaveBeenCalledTimes(1);
});
it('fails closed if persistence fails before POST and retains O1 in memory', async () => {
  mount();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  await screen.findByTestId('group-test-write-confirm');
  const fail = vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('quota'); });
  fireEvent.click(screen.getByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-unconfirmed');
  expect(api.testWrite).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-test-write-storage-error')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-recovery')).toHaveTextContent('O1');
  fail.mockRestore();
});
