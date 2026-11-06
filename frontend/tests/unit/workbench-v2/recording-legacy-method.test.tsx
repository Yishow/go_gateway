import { dbTargetAPI } from '@/services/datalink';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { RecordingEditorSummary } from '@/features/datalink/workbench-v2/recording-editor/RecordingEditorSummary';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { GroupEditor, type GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroupCalculationCapabilities } from '@/types/studioV2WriteGroupCalculation';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', async (importOriginal) => {
  const actual = await importOriginal<typeof import('react-i18next')>();
  return {
    ...actual,
    useTranslation: (ns?: string) => {
      if (useRealI18n.value) {
        return actual.useTranslation(ns);
      }
      return {
        t: (key: string, options?: { seconds?: number; name?: string }) => {
          if (options?.seconds !== undefined) return `${key}:${options.seconds}`;
          if (options?.name !== undefined) return `${key}:${options.name}`;
          return key;
        },
      };
    },
  };
});

const useRealI18n = { value: false };

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

afterEach(cleanup);

beforeEach(() => {
  vi.clearAllMocks();
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
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mockImplementation(async request => {
    const canonical = savedGroup();
    const members = request.group.members.map(member => {
      const source = canonical.members.find(value => value.device_id === member.device_id && value.point_id === member.point_id && value.tag_id === member.tag_id);
      if (!source) throw new Error('Fixture source identity is unavailable');
      return { ...member, source_revision: source.source_revision, mapping_revision: source.mapping_revision };
    });
    return {
      workspace_revision: 'created-1',
      group: savedGroup({
        ...request.group,
        members,
        destination: { ...canonical.destination, ...request.group.destination },
        id: 'new-group',
        revision: 'new-1',
      }),
    };
  });
});

function mountEditor(overrides: Partial<GroupEditorProps> = {}) {
  const state = savedGroupState();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const legacyGroup = savedGroup({
    id: 'legacy-1',
    name: 'Legacy Wide 60s',
    recording_profile: undefined,
    row_policy: {
      ...savedGroup().row_policy,
      interval_seconds: 60,
    },
  });
  const props: GroupEditorProps = {
    presentation: 'simple',
    group: legacyGroup,
    calculationCapabilities: capabilities,
    workspaceId: 'ws-1',
    workspaceRevision: 'db-setup-9',
    connector: state.db.connector,
    candidates: buildGroupCandidates(state).candidates,
    excluded: [],
    readonly: false,
    onSaved: vi.fn(),
    onClose: vi.fn(),
    onReload: vi.fn(),
    fetchLatest: vi.fn().mockResolvedValue({ workspace_revision: 'db-fresh-10', groups: [legacyGroup] }),
    ...overrides,
  };
  return render(
    <QueryClientProvider client={client}>
      <GroupEditor {...props} />
    </QueryClientProvider>,
  );
}

describe('legacy wide-table recording method section', () => {
  beforeEach(() => {
    useRealI18n.value = false;
  });

  it('shows latest (legacy wide table) with saved interval and opens a new recording prefilled with the same points and connector without updating the original recording', async () => {
    mountEditor();

    expect(screen.getByTestId('recording-legacy-method-text')).toHaveTextContent('legacyMethod.summary:60');
    expect(screen.getByTestId('recording-editor-summary')).toHaveTextContent('legacyMethod.label · 60 legacyMethod.secondsUnit');
    const shortcut = screen.getByTestId('recording-legacy-create-new');
    expect(shortcut).toHaveTextContent('legacyMethod.createWithPoints');
    expect(shortcut.className).toContain('border-[var(--gw-control)]');
    expect(shortcut.className).toContain('bg-[var(--gw-surface)]');

    fireEvent.click(shortcut);

    // Original recording was not updated in place.
    expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();

    // Editor now holds a new recording draft prefilled with the copy name, same points, and same connector.
    expect(screen.getByTestId('group-name')).toHaveValue('copy.name:Legacy Wide 60s');
    expect(screen.getByTestId('recording-method-groups')).toBeVisible();
    expect(await screen.findByTestId('member-calculation-pt-0')).toBeVisible();
    expect(await screen.findByTestId('member-calculation-pt-1')).toBeVisible();

    // Saving creates a new recording instead of updating the legacy recording.
    await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled());
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
    expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();

    const createCall = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0];
    expect(createCall.expected_connector_revision).toBe('crev-1');
    expect(createCall.group.destination.connector_id).toBe('conn-1');
    expect(createCall.group.members.map(m => m.point_id)).toEqual(['pt-0', 'pt-1']);
  });

  it('renders right-rail summary method and interval for legacy wide table in zh-TW and en (F3)', async () => {
    useRealI18n.value = true;
    const legacyDraft = groupToDraft(
      savedGroup({
        id: 'legacy-1',
        name: 'Legacy Wide 60s',
        recording_profile: undefined,
        row_policy: {
          ...savedGroup().row_policy,
          interval_seconds: 60,
        },
      }),
    );

    await act(async () => {
      await i18n.changeLanguage('zh-TW');
    });
    const { rerender } = render(
      <I18nextProvider i18n={i18n}>
        <RecordingEditorSummary draft={legacyDraft} saved dirty={false} />
      </I18nextProvider>,
    );
    expect(screen.getByTestId('recording-editor-summary')).toHaveTextContent(
      '最新值（舊式寬表） · 60 秒',
    );

    await act(async () => {
      await i18n.changeLanguage('en');
    });
    rerender(
      <I18nextProvider i18n={i18n}>
        <RecordingEditorSummary draft={legacyDraft} saved dirty={false} />
      </I18nextProvider>,
    );
    expect(screen.getByTestId('recording-editor-summary')).toHaveTextContent(
      'Latest (legacy wide table) · 60 s',
    );
  });
});
