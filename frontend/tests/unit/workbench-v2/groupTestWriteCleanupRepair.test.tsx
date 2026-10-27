import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { GroupTestWritePanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) =>
      options && 'count' in options ? `${key}:${options.count}` : key,
  }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    testWritePreview: vi.fn(),
    testWrite: vi.fn(),
    testWriteOperation: vi.fn(),
  },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);

function mountPanel() {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={client}>
      <GroupTestWritePanel
        groupId="group-1"
        scopeKey="group-1|rev-1|ws-1|cap-1"
        enabled={true}
      />
    </QueryClientProvider>
  );
}

describe('GroupTestWritePanel - Cleanup Repair Behavioral Regression', () => {
  const preview = {
    token: 'tok-1',
    operation_id: 'op-1',
    action: 'test_write' as const,
    group_id: 'group-1',
    group_revision: 'rev-1',
    expires_at: '2026-10-02T12:10:00Z',
    target: { connector_id: 'conn-1', dialect: 'sqlserver', database: 'db', schema: 'dbo', table: 'readings' },
    owner_column: 'entity',
    owner_value: 'gw-test-1',
    dedupe: '',
    cleanup: 'remove only the owned row',
    values: [{ column: 'temp', type: 'float', value: '25.0' }],
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders cleanup repair button and triggers cleanup repair when write succeeded but cleanup failed', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    api.testWrite.mockResolvedValueOnce({
      operation_id: 'op-1',
      action: 'test_write',
      status: 'succeeded',
      write_outcome: 'written_verified',
      cleanup_status: 'failed',
      cleanup_reason: 'cleanup-denied',
      created_at: 'a',
      updated_at: 'b',
      completed_at: 'c',
    });

    mountPanel();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    expect(await screen.findByTestId('group-test-write-preview-card')).toBeInTheDocument();

    await act(async () => {
      fireEvent.click(screen.getByTestId('group-test-write-confirm'));
    });
    expect(await screen.findByTestId('group-test-write-outcome')).toHaveTextContent('step4.group.test_write.outcome_written_verified');
    expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('step4.group.test_write.cleanup_failed');

    const repairBtn = screen.getByTestId('group-test-write-cleanup-repair');
    expect(repairBtn).toBeInTheDocument();

    // Now click cleanup repair
    api.testWrite.mockResolvedValueOnce({
      operation_id: 'op-1',
      action: 'test_write',
      status: 'succeeded',
      write_outcome: 'written_verified',
      cleanup_status: 'cleaned',
      created_at: 'a',
      updated_at: 'b',
      completed_at: 'd',
    });

    await act(async () => {
      fireEvent.click(repairBtn);
    });

    expect(api.testWrite).toHaveBeenCalledTimes(2);
    expect(api.testWrite).toHaveBeenLastCalledWith('group-1', {
      token: 'tok-1',
      operation_id: 'op-1',
      action: 'cleanup_repair',
    });

    expect(await screen.findByTestId('group-test-write-cleanup')).toHaveTextContent('step4.group.test_write.cleanup_cleaned');
    expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  });

  it('does not offer cleanup repair when write failed or cleanup already cleaned', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    api.testWrite.mockResolvedValueOnce({
      operation_id: 'op-1',
      action: 'test_write',
      status: 'failed',
      write_outcome: 'failed',
      cleanup_status: 'not_attempted',
      created_at: 'a',
      updated_at: 'b',
    });

    mountPanel();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    expect(await screen.findByTestId('group-test-write-preview-card')).toBeInTheDocument();

    await act(async () => {
      fireEvent.click(screen.getByTestId('group-test-write-confirm'));
    });
    expect(await screen.findByTestId('group-test-write-outcome')).toHaveTextContent('step4.group.test_write.outcome_failed');
    expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
  });

  it('does not offer cleanup repair for non-sqlserver dialects (e.g. mysql or sqlite)', async () => {
    for (const dialect of ['mysql', 'sqlite']) {
      api.testWritePreview.mockResolvedValueOnce({
        ...preview,
        target: { ...preview.target, dialect },
      });
      api.testWrite.mockResolvedValueOnce({
        operation_id: 'op-1',
        action: 'test_write',
        status: 'succeeded',
        write_outcome: 'written_verified',
        cleanup_status: 'failed',
        cleanup_reason: 'cleanup-denied',
        created_at: 'a',
        updated_at: 'b',
      });

      const { unmount } = mountPanel();
      fireEvent.click(screen.getByTestId('group-test-write-preview'));
      expect(await screen.findByTestId('group-test-write-preview-card')).toBeInTheDocument();

      await act(async () => {
        fireEvent.click(screen.getByTestId('group-test-write-confirm'));
      });
      expect(await screen.findByTestId('group-test-write-outcome')).toHaveTextContent('step4.group.test_write.outcome_written_verified');
      expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('step4.group.test_write.cleanup_failed');

      // Crucial: repair button is NOT shown for non-sqlserver dialect
      expect(screen.queryByTestId('group-test-write-cleanup-repair')).not.toBeInTheDocument();
      unmount();
    }
  });
});
