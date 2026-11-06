import { dbTargetAPI } from '@/services/datalink';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { RecordingDestinationLockNotice } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/RecordingDestinationLockNotice';
import { GroupEditor } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { newMixedDraft } from '@/features/datalink/workbench-v2/recording-editor/recordingMixedDraft';
import { createBasicRecordingRequestId } from '@/features/datalink/workbench-v2/state/basicRecording';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { EditorDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroupCalculationCapabilities } from '@/types/studioV2WriteGroupCalculation';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    readiness: vi.fn().mockResolvedValue(null),
    delivery: vi.fn().mockResolvedValue(null),
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    apply: vi.fn(),
    disable: vi.fn(),
    remove: vi.fn(),
    testWritePreview: vi.fn(),
    testWrite: vi.fn(),
    testWriteOperation: vi.fn(),
  },
}));

vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({
  studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) },
}));

const capabilities: WriteGroupCalculationCapabilities = {
  version: 1,
  result_shape: 'member-calculation-v1',
  editor_supported: true,
  draft_supported: true,
  activation_supported: false,
  reason: 'record_readiness_required',
  methods: ['latest', 'mean', 'delta'],
  cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true },
  schema_dialects: ['sqlite'],
  schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'],
  unsupported_combinations: [],
};

const destination = { connector_id: 'c1', connector_revision: 'r1', table_schema: 'public', table_name: 't' };
const base = newMixedDraft(destination, createBasicRecordingRequestId(), 'managed_wide_profile_v1');
const profileWide = { ...base, recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' } } as unknown as EditorDraft;
const legacyLong = { ...base, layout_version: 'managed_long_v1', recording_profile: undefined } as unknown as EditorDraft;
const legacyWide = { ...base, layout_version: undefined, recording_profile: undefined } as unknown as EditorDraft;

beforeEach(async () => {
  vi.clearAllMocks();
  await i18n.changeLanguage('zh-TW');
  vi.spyOn(dbTargetAPI, 'listConnectors').mockResolvedValue([
    {
      id: 'conn-1',
      identity_revision: 'crev-1',
      name: 'Line DB',
      kind: 'sqlite',
      enabled: true,
      status: 'ready',
      connection_config: { path: savedGroupState().db.connector.database },
      last_check_error: '',
      default_write_interval_seconds: 5,
      created_at: '',
      updated_at: '',
    },
  ]);
});
afterEach(cleanup);

describe('saved destination lock notice', () => {
  it('is shown only for saved locked records', () => {
    const { rerender } = render(<RecordingDestinationLockNotice draft={profileWide} saved={false} onCopy={vi.fn()} />);
    expect(screen.queryByTestId('recording-destination-lock-notice')).toBeNull();
    rerender(<RecordingDestinationLockNotice draft={legacyWide} saved onCopy={vi.fn()} />);
    expect(screen.queryByTestId('recording-destination-lock-notice')).toBeNull();
    rerender(<RecordingDestinationLockNotice draft={base} saved={false} onCopy={vi.fn()} />);
    expect(screen.queryByTestId('recording-destination-lock-notice')).toBeNull();
    rerender(<RecordingDestinationLockNotice draft={profileWide} saved onCopy={vi.fn()} />);
    expect(screen.getByTestId('recording-destination-lock-notice')).toBeInTheDocument();
  });

  it('uses different reasons for long-table and profile records', () => {
    const { rerender } = render(<RecordingDestinationLockNotice draft={legacyLong} saved onCopy={vi.fn()} />);
    const long = screen.getByTestId('recording-destination-lock-reason').textContent;
    expect(long).toContain('長表記錄與同一 Schema 的其他長表記錄共用資料表');
    rerender(<RecordingDestinationLockNotice draft={profileWide} saved onCopy={vi.fn()} />);
    const bound = screen.getByTestId('recording-destination-lock-reason').textContent;
    expect(bound).toContain('這份記錄的資料表已綁定');
    expect(bound).not.toBe(long);
  });

  it('calls the existing copy action, renders visible secondary button classes, and hides it when copy is unavailable', () => {
    const onCopy = vi.fn();
    const { rerender } = render(<RecordingDestinationLockNotice draft={profileWide} saved onCopy={onCopy} />);
    const btn = screen.getByRole('button', { name: '複製成新記錄' });
    expect(btn.className).toContain('border-[var(--gw-control)]');
    expect(btn.className).toContain('bg-[var(--gw-surface)]');
    fireEvent.click(btn);
    expect(onCopy).toHaveBeenCalledTimes(1);
    rerender(<RecordingDestinationLockNotice draft={profileWide} saved onCopy={onCopy} copyDisabled />);
    expect(screen.getByRole('button', { name: '複製成新記錄' })).toBeDisabled();
    rerender(<RecordingDestinationLockNotice draft={profileWide} saved />);
    expect(screen.queryByRole('button')).toBeNull();
    expect(screen.getByTestId('recording-destination-lock-reason')).toBeInTheDocument();
  });

  it.each([
    { layout_version: 'managed_long_v1' as const, kind: 'instantaneous' as const, interval: 15, unit: 'second' as const },
    { layout_version: 'managed_wide_profile_v1' as const, kind: 'instantaneous' as const, interval: 30, unit: 'second' as const },
    { layout_version: 'managed_long_v1' as const, kind: 'cumulative' as const, interval: 10, unit: 'minute' as const },
  ])('provides a copy entry in GroupEditor for saved locked $layout_version $kind profile records that opens a new prefilled recording without mutating the original', async ({ layout_version, kind, interval, unit }) => {
    const state = savedGroupState();
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    const group = savedGroup({
      id: 'prof-1',
      name: 'Profile Record',
      layout_version,
      storage_mode: 'plain',
      recording_profile: { kind, interval, unit },
      destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'profile_tbl' },
    });

    render(
      <I18nextProvider i18n={i18n}>
        <QueryClientProvider client={client}>
          <GroupEditor
            presentation="simple"
            group={group}
            calculationCapabilities={capabilities}
            workspaceId="ws-1"
            workspaceRevision="db-setup-9"
            connector={state.db.connector}
            candidates={buildGroupCandidates(state).candidates}
            excluded={[]}
            readonly={false}
            onSaved={vi.fn()}
            onClose={vi.fn()}
            onReload={vi.fn()}
            fetchLatest={vi.fn().mockResolvedValue({ workspace_revision: 'db-fresh-10', groups: [group] })}
          />
        </QueryClientProvider>
      </I18nextProvider>,
    );

    const copyBtn = screen.getByTestId('recording-destination-lock-copy');
    expect(copyBtn).toHaveTextContent('複製成新記錄');
    expect(copyBtn.className).toContain('border-[var(--gw-control)]');

    fireEvent.click(copyBtn);

    expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
    await waitFor(() => expect(screen.getByTestId('group-name')).toHaveValue('Profile Record（副本）'));
    expect(screen.queryByTestId('recording-destination-lock-notice')).toBeNull();
    expect(screen.getByTestId('recording-method-groups')).toBeVisible();
    expect(await screen.findByTestId('member-calculation-pt-0')).toBeVisible();
    expect(await screen.findByTestId('member-calculation-pt-1')).toBeVisible();
  });
});
