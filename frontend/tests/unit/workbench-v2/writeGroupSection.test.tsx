import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && 'count' in options ? `${key}:${options.count}` : key) }),
}));

let metadata: { status: string; columns: typeof TABLE_COLUMNS } = { status: 'exists', columns: TABLE_COLUMNS };
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: metadata.columns, status: metadata.status, refetch: vi.fn() }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
  },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);

function list(groups: WriteGroup[], revision = 'wrev-1') {
  return { workspace_id: 'ws-1', workspace_revision: revision, groups };
}

function renderSection(props: Partial<React.ComponentProps<typeof WriteGroupSection>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const utils = render(
    <QueryClientProvider client={client}>
      <WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} {...props} />
    </QueryClientProvider>,
  );
  return { client, ...utils };
}

beforeEach(() => {
  vi.clearAllMocks();
  metadata = { status: 'exists', columns: TABLE_COLUMNS };
  api.list.mockResolvedValue(list([]));
  api.delivery.mockResolvedValue({
    group_id: 'group-1', intake: { state: 'not_running' },
    stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  });
});

describe('EmptyErrorMismatchPlanRecovery', () => {
  it('names the missing prerequisite instead of leaving a dead control', async () => {
    const onNavigateStep = vi.fn();
    renderSection({ state: { ...savedGroupState(), devices: [], points: [], mappings: {} }, onNavigateStep });
    expect(await screen.findByTestId('group-prerequisite-device')).toBeInTheDocument();
    expect(screen.getByTestId('group-create')).toBeDisabled();
    fireEvent.click(screen.getByTestId('group-prerequisite-go'));
    expect(onNavigateStep).toHaveBeenCalledWith(1);
  });

  it('says saved Tags are missing, counts what is excluded, and points at the step that fixes it', async () => {
    const onNavigateStep = vi.fn();
    const state = savedGroupState();
    state.mappings = Object.fromEntries(Object.entries(state.mappings).map(([key, mapping]) => [key, { ...mapping, persisted: false, save_state: 'idle' as const }]));
    renderSection({ state, onNavigateStep });
    expect(await screen.findByTestId('group-prerequisite-tags')).toBeInTheDocument();
    expect(screen.getByTestId('group-prerequisite-excluded')).toHaveTextContent(':3');
    fireEvent.click(screen.getByTestId('group-prerequisite-go'));
    expect(onNavigateStep).toHaveBeenCalledWith(3);
  });

  it('blocks creating a group until the destination is saved', async () => {
    const state = savedGroupState();
    state.db = { ...state.db, connector: { ...state.db.connector, persisted: false, save_state: 'idle', connector_id: undefined, identity_revision: undefined } };
    renderSection({ state });
    expect(await screen.findByTestId('group-prerequisite-destination')).toBeInTheDocument();
    expect(screen.getByTestId('group-create')).toBeDisabled();
  });

  it('keeps a failed load apart from an empty list and retries', async () => {
    api.list.mockRejectedValueOnce(new Error('network down'));
    renderSection();
    expect(await screen.findByTestId('group-section-error')).toBeInTheDocument();
    expect(screen.queryByTestId('group-empty')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-section-retry'));
    expect(await screen.findByTestId('group-empty')).toBeInTheDocument();
    expect(api.list).toHaveBeenCalledTimes(2);
  });

  it('offers creating a group from an empty list once prerequisites are met', async () => {
    renderSection();
    expect(await screen.findByTestId('group-empty')).toBeInTheDocument();
    expect(screen.getByTestId('group-create')).not.toBeDisabled();
  });

  it('allows a blank connector table for a managed group and lets the server assign it', async () => {
    const state = savedGroupState({
      db: { ...savedGroupState().db, connector: { ...savedGroupState().db.connector, table: '' } },
    });
    const canonical = savedGroup({
      destination: { ...savedGroup().destination, table_name: 'gw_group_server', storage_strategy: 'managed' },
      members: [{ ...savedGroup().members[0], target_column: 'v_server' }],
      row_policy: { ...savedGroup().row_policy, record_key_column: 'record_id', group_id_column: 'group_id' },
      write_policy: { mode: 'append', dedupe_capability: 'receipt' },
    });
    api.create.mockResolvedValue({ workspace_revision: 'wrev-2', group: canonical });
    renderSection({ state });
    expect(await screen.findByTestId('group-empty')).toBeInTheDocument();
    expect(screen.getByTestId('group-create')).not.toBeDisabled();

    fireEvent.click(screen.getByTestId('group-create'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Managed line' } });
    fireEvent.click(screen.getByTestId('group-storage-managed'));
    fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
    fireEvent.click(screen.getByTestId('group-save'));

    await waitFor(() => expect(api.create).toHaveBeenCalledTimes(1));
    expect(api.create.mock.calls[0][0].group.destination).toMatchObject({ table_name: '', storage_strategy: 'managed' });
    expect(await screen.findByTestId('group-table')).toHaveValue('gw_group_server');
  });

  it('keeps custom groups strict when the connector table is blank', async () => {
    const state = savedGroupState({
      db: { ...savedGroupState().db, connector: { ...savedGroupState().db.connector, table: '' } },
    });
    renderSection({ state });
    fireEvent.click(await screen.findByTestId('group-create'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Custom line' } });
    fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
    fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
    expect(await screen.findByTestId('group-issues')).toHaveTextContent('step4.group.issue.table_required');
    expect(screen.getByTestId('group-save')).toBeDisabled();
    expect(api.create).not.toHaveBeenCalled();
  });

  it('shows groups for another destination as a mismatch, not as a match and not as nothing', async () => {
    api.list.mockResolvedValue(list([savedGroup({ destination: { ...savedGroup().destination, connector_revision: 'old-rev' } })]));
    renderSection();
    expect(await screen.findByTestId('group-list-mismatch-group-1')).toBeInTheDocument();
    expect(screen.getByTestId('group-none-match')).toBeInTheDocument();
    expect(screen.getByTestId('group-create')).not.toBeDisabled();
    fireEvent.click(screen.getByTestId('group-open-group-1'));
    expect(await screen.findByTestId('group-destination-mismatch')).toBeInTheDocument();
    // Re-pointing is an explicit edit that stays unsaved until the operator saves it.
    fireEvent.click(screen.getByTestId('group-use-current-destination'));
    expect(screen.queryByTestId('group-destination-mismatch')).not.toBeInTheDocument();
    expect(await screen.findByTestId('group-dirty-note')).toBeInTheDocument();
  });

  it('lists removed groups separately and opens them read-only', async () => {
    api.list.mockResolvedValue(list([savedGroup({ id: 'gone', status: 'deleted', name: 'Old group' })]));
    renderSection();
    expect(await screen.findByTestId('group-removed')).toBeInTheDocument();
    fireEvent.click(screen.getByText(/step4.group.list.removed/));
    fireEvent.click(screen.getByTestId('group-open-gone'));
    expect(await screen.findByTestId('group-editor')).toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toBeDisabled();
    expect(screen.getByTestId('group-save')).toBeDisabled();
    expect(screen.getByTestId('group-apply')).toBeDisabled();
    expect(screen.getByTestId('group-delete')).toBeDisabled();
  });
});

describe('WriteGroupEditorManagedAndCustom', () => {
  it('shows the persisted members, destination and revisions of a saved group', async () => {
    api.list.mockResolvedValue(list([savedGroup({ applied_revision: 'rev-0' })]));
    renderSection();
    fireEvent.click(await screen.findByTestId('group-open-group-1'));
    expect(await screen.findByTestId('group-name')).toHaveValue('Line A');
    expect(screen.getByTestId('group-table')).toHaveValue('readings');
    expect(screen.getByTestId('group-applied-revision')).toHaveTextContent('rev-0');
    expect(screen.getByTestId('group-member-include-pt-0')).toBeChecked();
    expect(screen.getByTestId('group-member-column-pt-0')).toHaveValue('temperature');
    expect(screen.getByTestId('group-member-include-pt-2')).not.toBeChecked();
    expect(screen.getByTestId('group-saved-note')).toBeInTheDocument();
  });

  it('adopts server-assigned managed table and columns after saving', async () => {
    const initial = savedGroup({
      destination: { ...savedGroup().destination, table_name: '', storage_strategy: 'managed' },
      members: [{ ...savedGroup().members[0], target_column: '' }],
      row_policy: { ...savedGroup().row_policy, record_key_column: 'record_id', group_id_column: 'group_id' },
      write_policy: { mode: 'append', dedupe_capability: 'receipt' },
    });
    const canonical = {
      ...initial,
      destination: { ...initial.destination, table_name: 'gw_group_server' },
      members: [{ ...initial.members[0], target_column: 'v_server' }],
    };
    api.list.mockResolvedValue(list([initial]));
    api.update.mockResolvedValue({ workspace_revision: 'wrev-2', group: canonical });
    renderSection();
    fireEvent.click(await screen.findByTestId('group-open-group-1'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Managed line' } });
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    expect(await screen.findByTestId('group-table')).toHaveValue('gw_group_server');
    expect(screen.getByTestId('group-member-column-pt-0')).toHaveValue('v_server');
    expect(screen.getByTestId('group-saved-note')).toBeInTheDocument();
  });

  it('switches managed and custom on the same group without asking for manual targets', async () => {
    const state = savedGroupState();
    expect(Object.keys(state.db.targets)).toHaveLength(0);
    api.list.mockResolvedValue(list([savedGroup()]));
    renderSection({ state });
    fireEvent.click(await screen.findByTestId('group-open-group-1'));
    fireEvent.click(await screen.findByTestId('group-storage-managed'));
    expect(screen.getByTestId('group-storage-note')).toHaveTextContent('managed_note');
    expect(screen.getByTestId('group-dirty-note')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-storage-custom'));
    await waitFor(() => expect(screen.getByTestId('group-saved-note')).toBeInTheDocument());
  });

  it('creates a group from saved Tags with persisted IDs and the revisions it was shown with', async () => {
    api.create.mockResolvedValue({ workspace_revision: 'wrev-2', group: savedGroup() });
    renderSection();
    // The refetch after saving has not returned the new group yet: the editor must not drop out meanwhile.
    fireEvent.click(await screen.findByTestId('group-create'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Line A' } });
    fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
    fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
    fireEvent.click(screen.getByTestId('group-member-include-pt-1'));
    fireEvent.change(screen.getByTestId('group-member-column-pt-1'), { target: { value: 'pressure' } });
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(api.create).toHaveBeenCalledTimes(1));
    const request = api.create.mock.calls[0][0];
    expect(request).toMatchObject({ workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_connector_revision: 'crev-1' });
    expect(request.group.members.map((m) => [m.device_id, m.point_id, m.tag_id, m.target_column])).toEqual([
      ['dev-1', 'pt-0', 'tag-0', 'temperature'], ['dev-1', 'pt-1', 'tag-1', 'pressure'],
    ]);
    expect(request.group.destination).toMatchObject({ connector_id: 'conn-1', connector_revision: 'crev-1', table_name: 'readings', storage_strategy: 'custom' });
    // After saving, the editor shows the saved group, not the create form.
    expect(await screen.findByTestId('group-revisions')).toBeInTheDocument();
  });

  it('renames and edits members through an update that carries the group revision', async () => {
    api.list.mockResolvedValue(list([savedGroup()]));
    api.update.mockResolvedValue({ workspace_revision: 'wrev-2', group: savedGroup({ revision: 'rev-2', name: 'Line B' }) });
    renderSection();
    fireEvent.click(await screen.findByTestId('group-open-group-1'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Line B' } });
    fireEvent.click(screen.getByTestId('group-member-include-pt-1'));
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    const [id, request] = api.update.mock.calls[0];
    expect(id).toBe('group-1');
    expect(request).toMatchObject({ expected_group_revision: 'rev-1', expected_workspace_revision: 'wrev-1', expected_connector_revision: 'crev-1' });
    expect(request.group.name).toBe('Line B');
    expect(request.group.members).toHaveLength(1);
  });

  it('keeps the draft and offers a reload when the saved configuration changed elsewhere', async () => {
    api.list.mockResolvedValue(list([savedGroup()]));
    api.update.mockRejectedValue(Object.assign(new Error('conflict'), { response: { status: 409, data: { success: false, error: { code: 'revision_mismatch', message: 'x', retryable: false, request_id: 'r', action: 'reload' } } } }));
    renderSection();
    fireEvent.click(await screen.findByTestId('group-open-group-1'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Mine' } });
    fireEvent.click(screen.getByTestId('group-save'));
    expect(await screen.findByTestId('group-conflict')).toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toHaveValue('Mine');
    fireEvent.click(screen.getByTestId('group-conflict-reload'));
    await waitFor(() => expect(api.list.mock.calls.length).toBeGreaterThan(1));
  });

  it('does not send a second save while the first is in flight', async () => {
    api.list.mockResolvedValue(list([savedGroup()]));
    let release: (value: unknown) => void = () => undefined;
    api.update.mockReturnValue(new Promise((resolve) => { release = resolve; }) as never);
    renderSection();
    fireEvent.click(await screen.findByTestId('group-open-group-1'));
    fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'Once' } });
    const save = screen.getByTestId('group-save');
    fireEvent.click(save);
    fireEvent.click(save);
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    expect(api.update).toHaveBeenCalledTimes(1);
    api.list.mockResolvedValue(list([savedGroup({ revision: 'rev-2', name: 'Once' })], 'wrev-2'));
    release({ workspace_revision: 'wrev-2', group: savedGroup({ revision: 'rev-2', name: 'Once' }) });
    await waitFor(() => expect(screen.queryByTestId('group-dirty-note')).not.toBeInTheDocument());
  });

  it('blocks every edit and mutation while the session is read-only', async () => {
    api.list.mockResolvedValue(list([savedGroup()]));
    renderSection({ readonly: true });
    expect(await screen.findByTestId('group-create')).toBeDisabled();
    fireEvent.click(screen.getByTestId('group-open-group-1'));
    expect(await screen.findByTestId('group-name')).toBeDisabled();
    for (const id of ['group-save', 'group-apply', 'group-disable', 'group-delete', 'group-test-write-preview', 'group-bulk-include']) {
      expect(screen.getByTestId(id), id).toBeDisabled();
    }
    expect(screen.getByTestId('group-member-include-pt-0')).toBeDisabled();
  });
});
