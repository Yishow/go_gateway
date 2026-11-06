import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { GroupEditorPreparation } from '@/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation';
import { newMixedDraft } from '@/features/datalink/workbench-v2/recording-editor/recordingMixedDraft';
import { createBasicRecordingRequestId } from '@/features/datalink/workbench-v2/state/basicRecording';
import { newDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { GroupPolicyFields } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields';
import { RecordingDestinationFields } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationFields';
import { readRecordingDestinationCatalog } from '@/services/studioV2DatabaseCatalog';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('@/services/studioV2DatabaseCatalog', () => ({
  readRecordingDestinationCatalog: vi.fn(),
}));

const destination = { connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: 'public', table_name: '' };

beforeEach(async () => {
  vi.clearAllMocks();
  await i18n.changeLanguage('zh-TW');
});

afterEach(cleanup);

describe('Task 3.1: Hidden and pending areas state why', () => {
  it('explains why row-rule fields are hidden for a per-point (mixed) recording in GroupPolicyFields', () => {
    const mixedDraft = newMixedDraft(destination, createBasicRecordingRequestId());
    const { rerender } = render(
      <GroupPolicyFields
        expanded
        draft={mixedDraft}
        readonly={false}
        profile={true}
        onChange={vi.fn()}
        onDedupe={vi.fn()}
      />,
    );

    const note = screen.getByTestId('group-policy-mixed-note');
    expect(note).toHaveTextContent(
      '間隔在各點位設定；不完整的列固定為部分結果；受管資料表的重複保護、實體鍵與來源紀錄欄位由系統決定',
    );
    expect(screen.queryByTestId('group-interval')).toBeNull();
    expect(screen.queryByTestId('group-incomplete')).toBeNull();

    const legacyDraft = newDraft(destination);
    rerender(
      <GroupPolicyFields
        expanded
        draft={legacyDraft}
        readonly={false}
        profile={false}
        onChange={vi.fn()}
        onDedupe={vi.fn()}
      />,
    );
    expect(screen.queryByTestId('group-policy-mixed-note')).toBeNull();
    expect(screen.getByTestId('group-interval')).toBeInTheDocument();
  });

  it('shows that table structure can be checked after saving when a managed recording is unsaved', () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const baseSave = { saved: false, dirty: true, busy: false, disabled: false, unknown: false, invoke: vi.fn() };
    const baseReadiness = { loading: false, failed: false, updatedAt: 0, check: vi.fn() };

    const { rerender } = render(
      <QueryClientProvider client={client}>
        <GroupEditorPreparation
          compact
          managed={true}
          group={undefined}
          workspaceId="ws-1"
          workspaceRevision="w1"
          destinationMatches={true}
          dirty={true}
          readonly={false}
          consentIdentity="unsaved-managed"
          save={baseSave}
          blocked={false}
          readiness={baseReadiness}
        />
      </QueryClientProvider>,
    );

    expect(screen.getByTestId('group-preparation-unsaved-note')).toHaveTextContent('保存設定後才能檢查資料表結構');

    rerender(
      <QueryClientProvider client={client}>
        <GroupEditorPreparation
          compact
          managed={false}
          group={undefined}
          workspaceId="ws-1"
          workspaceRevision="w1"
          destinationMatches={true}
          dirty={true}
          readonly={false}
          consentIdentity="unsaved-custom"
          save={baseSave}
          blocked={false}
          readiness={baseReadiness}
        />
      </QueryClientProvider>,
    );

    expect(screen.queryByTestId('group-preparation-unsaved-note')).toBeNull();
  });

  it('does not show loading when the destination catalog is not being read, and shows unreachable database error with retry on failure', async () => {
    const connector = savedGroupState().db.connector;
    const lockedDraft = {
      ...newMixedDraft(destination, createBasicRecordingRequestId()),
      table_name: 'gw_locked',
    };
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    // 1. Not reading catalog (locked destination, even without scopedPresentation) -> no request, no loading message
    const { unmount } = render(
      <QueryClientProvider client={client}>
        <RecordingDestinationFields
          workspaceId="ws-1"
          workspaceRevision="w1"
          connector={connector}
          group={savedGroup()}
          draft={lockedDraft}
          readonly={false}
          locked={true}
          storageLocked={true}
          onChange={vi.fn()}
          onStorageStrategy={vi.fn()}
        />
      </QueryClientProvider>,
    );
    await act(async () => {});
    expect(readRecordingDestinationCatalog).not.toHaveBeenCalled();
    expect(screen.queryByText(/正在讀取 Schema 與資料表/)).toBeNull();
    unmount();

    // 2. Catalog read fails -> shows that destination database could not be reached + retry action
    vi.mocked(readRecordingDestinationCatalog)
      .mockRejectedValueOnce(new Error('DB unreachable'))
      .mockResolvedValueOnce({
        workspace_id: 'ws-1',
        workspace_revision: 'w1',
        connector_id: 'conn-1',
        connector_revision: 'crev-1',
        kind: 'postgres',
        default_schema: 'public',
        schemas: [{ name: 'public', tables: ['t1'] }],
      });

    const editableDraft = newMixedDraft(destination, createBasicRecordingRequestId());
    render(
      <QueryClientProvider client={client}>
        <RecordingDestinationFields
          workspaceId="ws-1"
          workspaceRevision="w1"
          connector={connector}
          draft={editableDraft}
          readonly={false}
          locked={false}
          storageLocked={true}
          onChange={vi.fn()}
          onStorageStrategy={vi.fn()}
        />
      </QueryClientProvider>,
    );

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent('無法連到目的地資料庫');
    expect(screen.queryByText(/正在讀取 Schema 與資料表/)).toBeNull();

    const retryBtn = screen.getByRole('button', { name: '重新讀取資料表' });
    fireEvent.click(retryBtn);

    await waitFor(() => expect(readRecordingDestinationCatalog).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(screen.queryByRole('alert')).toBeNull());
    expect(screen.getByRole('option', { name: 't1' })).toBeInTheDocument();
  });
});
