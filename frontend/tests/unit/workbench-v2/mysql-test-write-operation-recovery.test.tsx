import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
vi.mock('@/hooks/datalink/useRecordingTestWriteRecovery', () => ({ useRecordingTestWriteRecovery: () => ({ available: true, isFetching: false, refetch: vi.fn() }) }));
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { CircuitProfileSchemaPreparation } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileSchemaPreparation';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { savedGroup } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({
  GroupSchemaPanel: () => null,
}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupReadinessQuery: () => ({
    data: { group_revision: 'rev-1', schema_ready: true }, isFetching: false, isError: false,
  }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
  },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const preview = {
  token: 'token-original', operation_id: 'operation-original', action: 'test_write' as const,
  group_id: 'group-1', group_revision: 'rev-1', expires_at: '2026-10-27T12:10:00Z',
  target: { connector_id: 'conn-1', dialect: 'mysql', database: 'recording', schema: 'recording', table: 'gw_long_test' },
  owner_column: '_gw_owner_test', owner_value: 'owned-test', dedupe: 'receipt', cleanup: 'scoped cleanup',
  values: [{ column: 'value_decimal', type: 'decimal', value: '12000.4' }],
};
const operation = {
  operation_id: preview.operation_id, action: 'test_write' as const, status: 'succeeded' as const,
  write_outcome: 'written_verified' as const, cleanup_status: 'cleaned' as const,
  created_at: '2026-10-27T12:00:00Z', updated_at: '2026-10-27T12:00:01Z',
};

function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const group = savedGroup({
    layout_version: 'managed_long_v1', storage_mode: 'plain',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
  });
  const content = (available = true, revision = 'workspace-1') => (
    <QueryClientProvider client={client}>
      <CircuitProfileSchemaPreparation
        workspaceId="ws-1" workspaceRevision={revision} preparedGroups={[group]}
        disabled={false} readonly={false} isPreparing={false} destinationMatches
        onPrepare={vi.fn()} onApplied={vi.fn()} isMySQL
        mysqlTestWriteAvailable={available}
      />
    </QueryClientProvider>
  );
  const view = render(content());
  return { view, content };
}

async function confirmWithLostReply() {
  fireEvent.click(screen.getByTestId('group-test-write-preview'));
  fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
  expect(await screen.findByTestId('group-test-write-unconfirmed')).toBeInTheDocument();
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
  api.testWrite.mockRejectedValue(new Error('lost acknowledgement'));
  api.testWriteOperation.mockResolvedValue(operation);
});

describe('MySQL test-write operation recovery', () => {
  it('keeps the original operation query during capability refetch, including while new writes are unavailable', async () => {
    const { view, content } = mount();
    await confirmWithLostReply();
    view.rerender(content(false));

    expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
    expect(screen.getByTestId('group-test-write-confirm')).toBeDisabled();
    fireEvent.click(screen.getByTestId('group-test-write-confirm'));
    expect(api.testWrite).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByTestId('group-test-write-check-unconfirmed'));
    expect(await screen.findByTestId('group-test-write-outcome')).toHaveTextContent('outcome_written_verified');
    expect(api.testWriteOperation).toHaveBeenCalledWith(preview.operation_id);
    expect(api.testWrite).toHaveBeenCalledTimes(1);

    view.rerender(content());
    expect(screen.getByTestId('group-test-write-result')).toHaveTextContent(preview.operation_id);
  });

  it('keeps the original operation when an unrelated setup save moves only the workspace revision', async () => {
    const { view, content } = mount();
    await confirmWithLostReply();
    view.rerender(content(true, 'workspace-2'));

    fireEvent.click(screen.getByTestId('group-test-write-check-unconfirmed'));
    await waitFor(() => expect(api.testWriteOperation).toHaveBeenCalledWith(preview.operation_id));
    expect(await screen.findByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_cleaned');
    expect(api.testWrite).toHaveBeenCalledTimes(1);
  });
});
