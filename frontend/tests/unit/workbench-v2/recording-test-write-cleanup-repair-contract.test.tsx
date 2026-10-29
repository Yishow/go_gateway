import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupTestWritePanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel';
import { RecordingTestWriteRecoveryPanel } from '@/features/datalink/workbench-v2/steps/step4/RecordingTestWriteRecoveryPanel';
import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroupTestWriteOperation } from '@/types/studioV2WriteGroupTestWrite';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useRecordingTestWriteRecovery', () => ({ useRecordingTestWriteRecovery: () => ({ available: true, isFetching: false, refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn() } }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const failedCleanup: WriteGroupTestWriteOperation = { operation_id: 'O1', action: 'test_write', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'failed', created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z', completed_at: '2026-10-28T00:00:01Z' };
const cleaned = { ...failedCleanup, cleanup_status: 'cleaned' as const };
const preview = { token: 'original-token', operation_id: 'O1', action: 'test_write' as const, group_id: 'G1', group_revision: 'R1', expires_at: '', target: { connector_id: 'C1', dialect: 'sqlserver', database: '', schema: '', table: 'test' }, owner_column: '', owner_value: '', dedupe: '', values: [], cleanup: '' };
function retain() { recordingTestWriteRecovery.retain({ workspace_id: 'W1', group_id: 'G1', group_revision: 'R1', operation_id: 'O1', operation_kind: 'test_write', created_at: failedCleanup.created_at }); }
function mount(recoveryOnly = false) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const content = (workspace = 'W1', revision = 'R1') => <QueryClientProvider client={client}>
    {recoveryOnly ? <RecordingTestWriteRecoveryPanel workspaceId={workspace} /> : <GroupTestWritePanel workspaceId={workspace} groupId="G1" scopeKey={`G1|${revision}`} enabled />}
  </QueryClientProvider>;
  return { ...render(content()), content };
}
beforeEach(() => {
  vi.resetAllMocks(); localStorage.clear();
  for (const ref of recordingTestWriteRecovery.read().references) recordingTestWriteRecovery.resolve({ ...cleaned, operation_id: ref.operation_id }, ref.workspace_id);
  api.testWritePreview.mockResolvedValue(preview);
  api.testWrite.mockResolvedValueOnce(failedCleanup);
  api.testWriteOperation.mockResolvedValue(failedCleanup);
});
async function initialWrite() {
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-cleanup-repair');
}
it('shows queried phase facts and manual guidance without creating a token or repair action', async () => {
  retain(); mount(true);
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  expect(await screen.findByTestId('recovery-O1-readback')).toHaveTextContent('readback_verified');
  expect(screen.getByTestId('recovery-O1-cleanup')).toHaveTextContent('cleanup_failed');
  expect(screen.getByTestId('recovery-O1-manual-guidance')).toBeInTheDocument();
  expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  expect(api.testWrite).not.toHaveBeenCalled(); expect(api.testWritePreview).not.toHaveBeenCalled();
});
it.each([['written_verified', 'succeeded', 'readback_verified'], ['written_unverified', 'partial', 'readback_unverified']] as const)('keeps final %s history visible but removes pending truth; scope/unmount clears display', async (write_outcome, status, readback) => {
  retain(); api.testWriteOperation.mockResolvedValue({ ...cleaned, write_outcome, status });
  const view = mount(true);
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  expect(await screen.findByTestId('recovery-O1-cleanup')).toHaveTextContent('cleanup_cleaned');
  expect(screen.getByTestId('recovery-O1-readback')).toHaveTextContent(readback);
  expect(recordingTestWriteRecovery.read().references).toEqual([]);
  view.rerender(view.content('W2'));
  expect(screen.queryByTestId('recovery-O1-result')).not.toBeInTheDocument();
  view.rerender(view.content());
  expect(screen.queryByTestId('recovery-O1-result')).not.toBeInTheDocument();
  view.unmount(); mount(true);
  expect(screen.queryByTestId('recovery-O1-result')).not.toBeInTheDocument();
});
it('does not show a late recovery answer after W1 to W2 to W1', async () => {
  retain(); let finish!: (value: WriteGroupTestWriteOperation) => void;
  api.testWriteOperation.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  const view = mount(true); fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  view.rerender(view.content('W2')); view.rerender(view.content());
  await act(async () => { finish(failedCleanup); });
  expect(screen.queryByTestId('recovery-O1-result')).not.toBeInTheDocument();
  expect(api.testWrite).not.toHaveBeenCalled();
});
it.each(['failed', 'unknown'] as const)('repairs %s cleanup with original token/O1 and no extra insert or preview', async (cleanup_status) => {
  api.testWrite.mockReset().mockResolvedValueOnce({ ...failedCleanup, cleanup_status });
  const view = mount(); await initialWrite();
  let finish!: (value: WriteGroupTestWriteOperation) => void;
  api.testWrite.mockImplementationOnce(() => new Promise((resolve) => { finish = resolve; }));
  const repair = screen.getByTestId('group-test-write-cleanup-repair');
  fireEvent.click(repair); fireEvent.click(repair);
  await waitFor(() => expect(api.testWrite).toHaveBeenCalledTimes(2));
  expect(api.testWrite).toHaveBeenNthCalledWith(1, 'G1', { token: 'original-token', operation_id: 'O1' });
  expect(api.testWrite).toHaveBeenNthCalledWith(2, 'G1', { token: 'original-token', operation_id: 'O1', action: 'cleanup_repair' });
  await act(async () => { finish(cleaned); });
  expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_cleaned');
  expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  expect(api.testWritePreview).toHaveBeenCalledTimes(1);
  view.unmount();
});
it('lost repair ACK only queries O1 and does not blindly repair or insert again', async () => {
  mount(); await initialWrite(); api.testWrite.mockRejectedValueOnce(new Error('lost repair ACK'));
  fireEvent.click(screen.getByTestId('group-test-write-cleanup-repair'));
  await screen.findByTestId('group-test-write-unconfirmed');
  expect(screen.getByTestId('group-test-write-cleanup-repair')).toBeDisabled();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  api.testWriteOperation.mockResolvedValue(cleaned);
  fireEvent.click(screen.getByTestId('group-test-write-check-unconfirmed'));
  await waitFor(() => expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_cleaned'));
  expect(api.testWriteOperation).toHaveBeenCalledTimes(1); expect(api.testWriteOperation).toHaveBeenCalledWith('O1');
  expect(api.testWrite).toHaveBeenCalledTimes(2); expect(api.testWritePreview).toHaveBeenCalledTimes(1);
});
it('after revision navigation, recovery is query-only because original confirmation is discarded', async () => {
  const view = mount(); await initialWrite(); view.rerender(view.content('W1', 'R2'));
  await waitFor(() => expect(api.testWriteOperation).toHaveBeenCalledWith('O1'));
  expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  expect(await screen.findByTestId('group-test-write-manual-guidance')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
  expect(api.testWrite).toHaveBeenCalledTimes(1);
});
it('rejects a preview returned for a different group instead of offering its confirmation', async () => {
  api.testWritePreview.mockResolvedValue({ ...preview, group_id: 'FOREIGN-GROUP' });
  mount(); fireEvent.click(screen.getByTestId('group-test-write-preview'));
  await waitFor(() => expect(api.testWritePreview).toHaveBeenCalledTimes(1));
  await screen.findByTestId('group-test-write-error');
  expect(screen.queryByTestId('group-test-write-confirm')).not.toBeInTheDocument();
  expect(api.testWrite).not.toHaveBeenCalled();
});
it('does not credit a late repair ACK after R1 to R2 to R1 navigation', async () => {
  const view = mount(); await initialWrite();
  let finish!: (value: WriteGroupTestWriteOperation) => void;
  api.testWrite.mockImplementationOnce(() => new Promise((resolve) => { finish = resolve; }));
  fireEvent.click(screen.getByTestId('group-test-write-cleanup-repair'));
  await waitFor(() => expect(api.testWrite).toHaveBeenCalledTimes(2));
  view.rerender(view.content('W1', 'R2')); view.rerender(view.content());
  await waitFor(() => expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_failed'));
  await act(async () => { finish(cleaned); });
  expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_failed');
  expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  expect(api.testWritePreview).toHaveBeenCalledTimes(1);
});

it('keeps prior checked facts labelled as history when the next query fails', async () => {
  retain(); mount(true);
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  await screen.findByTestId('recovery-O1-readback');
  api.testWriteOperation.mockRejectedValueOnce(new Error('private endpoint credential'));
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  expect(await screen.findByRole('alert')).toHaveTextContent('unconfirmed');
  expect(screen.getByText('step4.group.test_write.last_checked_result')).toBeInTheDocument();
  expect(screen.getByTestId('recovery-O1-cleanup')).toHaveTextContent('cleanup_failed');
  expect(screen.queryByText(/private endpoint/)).not.toBeInTheDocument();
  expect(recordingTestWriteRecovery.read().references).toHaveLength(1);
  expect(api.testWrite).not.toHaveBeenCalled();
});
it('ignores an exact-query response with another operation ID', async () => {
  retain(); api.testWriteOperation.mockResolvedValue({ ...cleaned, operation_id: 'O2' }); mount(true);
  fireEvent.click(screen.getByTestId('group-test-write-recovery-query-O1'));
  await waitFor(() => expect(screen.getByTestId('group-test-write-recovery-query-O1')).toBeEnabled());
  expect(screen.queryByTestId('recovery-O1-result')).not.toBeInTheDocument();
  expect(recordingTestWriteRecovery.read().references[0].operation_id).toBe('O1');
  expect(api.testWrite).not.toHaveBeenCalled();
});
it('never repairs preview O1 while recovered evidence belongs to another same-group operation O2', async () => {
  api.testWrite.mockReset().mockRejectedValueOnce(new Error('lost O1 ACK'));
  mount();
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-unconfirmed');
  await act(async () => {
    recordingTestWriteRecovery.resolve(cleaned, 'W1');
    recordingTestWriteRecovery.retain({ workspace_id: 'W1', group_id: 'G1', group_revision: 'R1', operation_id: 'O2', operation_kind: 'test_write', created_at: failedCleanup.created_at });
  });
  api.testWriteOperation.mockResolvedValue({ ...failedCleanup, operation_id: 'O2' });
  fireEvent.click(screen.getByTestId('group-test-write-check-unconfirmed'));
  await waitFor(() => expect(screen.getByTestId('group-test-write-result')).toHaveTextContent('O2'));
  expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-manual-guidance')).toBeInTheDocument();
  expect(api.testWriteOperation).toHaveBeenLastCalledWith('O2');
  expect(api.testWrite).toHaveBeenCalledTimes(1);
});
it('treats a mismatched repair acknowledgement as unconfirmed and only queries original O1', async () => {
  mount(); await initialWrite(); api.testWrite.mockResolvedValueOnce({ ...cleaned, operation_id: 'O2' });
  fireEvent.click(screen.getByTestId('group-test-write-cleanup-repair'));
  expect(await screen.findByTestId('group-test-write-unconfirmed')).toHaveTextContent('operation_status_unconfirmed');
  expect(screen.getByTestId('group-test-write-cleanup-repair')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-test-write-cleanup-repair'));
  expect(api.testWrite).toHaveBeenCalledTimes(2);
  api.testWriteOperation.mockResolvedValue(cleaned);
  fireEvent.click(screen.getByTestId('group-test-write-check-unconfirmed'));
  await waitFor(() => expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_cleaned'));
  expect(api.testWriteOperation).toHaveBeenLastCalledWith('O1');
  expect(api.testWrite).toHaveBeenCalledTimes(2);
});
