import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupTestWritePanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel';
import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroupTestWriteOperation } from '@/types/studioV2WriteGroupTestWrite';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useRecordingTestWriteRecovery', () => ({ useRecordingTestWriteRecovery: () => ({ available: true, isFetching: false, refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn() } }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const base: WriteGroupTestWriteOperation = { operation_id: 'O1', action: 'test_write', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'failed', created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z', completed_at: '2026-10-28T00:00:01Z' };
beforeEach(() => {
  vi.clearAllMocks(); localStorage.clear();
  for (const ref of recordingTestWriteRecovery.read().references) recordingTestWriteRecovery.resolve({ ...base, operation_id: ref.operation_id, cleanup_status: 'cleaned' }, ref.workspace_id);
  api.testWritePreview.mockResolvedValue({ token: 'original-token', operation_id: 'O1', action: 'test_write', group_id: 'G1', group_revision: 'R1', expires_at: '', target: { connector_id: 'C1', dialect: 'sqlite', database: '', schema: '', table: 'test' }, owner_column: '', owner_value: '', dedupe: '', values: [], cleanup: '' });
});
async function confirm(operation: WriteGroupTestWriteOperation) {
  api.testWrite.mockResolvedValue(operation);
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><GroupTestWritePanel workspaceId="W1" groupId="G1" scopeKey="G1|R1" enabled /></QueryClientProvider>);
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  await screen.findByTestId('group-test-write-result');
}
it.each([
  [{ ...base, cleanup_status: 'failed' }, 'write_committed', 'readback_verified', 'cleanup_failed'],
  [{ ...base, cleanup_status: 'unknown' }, 'write_committed', 'readback_verified', 'cleanup_unknown'],
  [{ ...base, status: 'partial', write_outcome: 'written_unverified', cleanup_status: 'cleaned' }, 'write_committed', 'readback_unverified', 'cleanup_cleaned'],
  [{ ...base, status: 'failed', write_outcome: 'failed', cleanup_status: 'not_attempted' }, 'write_not_written', 'readback_not_attempted', 'cleanup_not_attempted'],
  [{ ...base, status: 'unknown', write_outcome: 'unknown', cleanup_status: 'unknown' }, 'write_unconfirmed', 'readback_unconfirmed', 'cleanup_unknown'],
  [{ ...base, status: 'pending', write_outcome: undefined, cleanup_status: undefined, completed_at: undefined }, 'write_unconfirmed', 'readback_unconfirmed', 'cleanup_unconfirmed'],
  [{ ...base, status: 'running', write_outcome: undefined, cleanup_status: undefined, completed_at: undefined }, 'write_unconfirmed', 'readback_unconfirmed', 'cleanup_unconfirmed'],
] as const)('separates write/readback/cleanup facts for %j', async (operation, write, readback, cleanup) => {
  await confirm(operation);
  expect(screen.getByTestId('group-test-write-write-phase')).toHaveTextContent(write);
  expect(screen.getByTestId('group-test-write-readback')).toHaveTextContent(readback);
  expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent(cleanup);
  expect(api.testWrite).toHaveBeenCalledTimes(1);
});
it('shows a running cleanup claim with prior verified evidence, never as current cleanup success', async () => {
  await confirm({ ...base, status: 'running' });
  expect(screen.getByTestId('group-test-write-write-phase')).toHaveTextContent('write_committed');
  expect(screen.getByTestId('group-test-write-readback')).toHaveTextContent('readback_verified');
  expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_repair_running');
  expect(screen.getByTestId('group-test-write-prior-cleanup')).toHaveTextContent('cleanup_failed');
  expect(screen.queryByTestId('group-test-write-manual-guidance')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-check')).toBeInTheDocument();
  expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
});
it('gives query-only/manual guidance for unsupported cleanup repair', async () => {
  await confirm(base);
  expect(screen.getByTestId('group-test-write-manual-guidance')).toBeInTheDocument();
  expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  await act(async () => { fireEvent.click(screen.getByTestId('group-test-write-preview')); });
  expect(api.testWritePreview).toHaveBeenCalledTimes(1);
  expect(api.testWrite).toHaveBeenCalledTimes(1);
});
