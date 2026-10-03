import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) =>
      options && 'count' in options ? `${key}:${options.count}` : key,
  }),
}));

const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] } as { workspace_id?: string; workspace_revision?: string; groups?: WriteGroup[] } | undefined,
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};
const deliveryQuery = { data: undefined, isError: false };
const ensureAsync = vi.fn();
const startAsync = vi.fn();
const operationQuery = {
  data: undefined as RecordingStartOperation | undefined,
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: () => deliveryQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: ensureAsync, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));

function renderPanel(props: Partial<React.ComponentProps<typeof BasicRecordingPanel>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} {...props} />
    </QueryClientProvider>,
  );
}

type ShareStatus = NonNullable<React.ComponentProps<typeof BasicRecordingPanel>['shareStatus']>;
function shareOnlyState() {
  const state = savedGroupState();
  return savedGroupState({ db: {
    ...state.db,
    connector: { ...state.db.connector, persisted: false, save_state: 'idle', connector_id: undefined, identity_revision: undefined },
  } });
}
function shareOnlyStatus(overrides: Partial<ShareStatus> = {}): ShareStatus {
  return {
    workspace_id: 'ws-1', enabled: true, configured_enabled: true, port: 502, address: '127.0.0.1',
    bind_state: 'pass', mapping_count: 1, readiness_token: 'share-token', settings_revision: 'settings-1', workspace_revision: 'share-wrev-1',
    ...overrides,
  };
}

const canonicalManaged = savedGroup({
  basic_managed_device_id: 'dev-1',
  destination: {
    ...savedGroup().destination,
    storage_strategy: 'managed', table_name: 'gw_group_server', database: '/recording.db',
  },
  members: [
    ...savedGroup().members.map((member) => ({ ...member, target_column: '' })),
    { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', source_revision: 's2', mapping_revision: 'm2', target_column: '', required: true },
  ],
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 },
  write_policy: { mode: 'append', dedupe_capability: 'receipt' },
});

function operation(overrides: Partial<RecordingStartOperation> = {}): RecordingStartOperation {
  return {
    operation_id: 'operation-1', action: 'recording_start', status: 'succeeded', stage: 'complete',
    intent_digest: 'a'.repeat(64), workspace_id: 'ws-1', setup_revision: 'wrev-1', device_ids: ['dev-1'],
    groups: [{ group_id: canonicalManaged.id, group_revision: canonicalManaged.revision, saved: true, ready: true, applied: true }],
    devices: [{ device_id: 'dev-1', activated: true }],
    created_at: '2026-10-04T00:00:00Z', updated_at: '2026-10-04T00:00:01Z', ...overrides,
  };
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockResolvedValue({ data: groupsQuery.data });
  ensureAsync.mockResolvedValue({ workspace_revision: 'wrev-2', group: canonicalManaged });
  startAsync.mockResolvedValue(operation({ status: 'pending', stage: 'save' }));
  operationQuery.data = undefined;
});

describe('BasicRecordingPanel regressions', () => {
  it('shows saved, readiness and applied checkpoints plus device activation', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording',
      request: {
        request_id: '99999999-9999-4999-8999-999999999999', workspace_id: 'ws-1',
        expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
        groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
      }, operation_id: 'operation-1',
    });
    operationQuery.data = operation({
      status: 'partial', stage: 'activation',
      groups: [{ group_id: canonicalManaged.id, group_revision: canonicalManaged.revision, saved: true, ready: false, applied: false }],
      devices: [{ device_id: 'dev-1', activated: false, reason: 'activation_failed' }],
    });
    renderPanel();
    expect(await screen.findByTestId('basic-recording-group-progress')).toHaveTextContent('step4.basic.progress.saved');
    expect(screen.getByTestId('basic-recording-group-progress')).toHaveTextContent('step4.basic.progress.readiness_unverified');
    expect(screen.getByTestId('basic-recording-group-progress')).toHaveTextContent('step4.basic.progress.not_applied');
    expect(screen.getByTestId('basic-recording-device-progress')).toHaveTextContent('step4.basic.progress.not_activated');
  });

  it('offers the original request while the ledger is still running', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    const originalRequest = {
      request_id: '55555555-5555-4555-8555-555555555555', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1',
    });
    operationQuery.data = operation({ status: 'running', stage: 'activation' });

    renderPanel();

    const retry = await screen.findByTestId('basic-recording-retry');
    expect(retry).not.toBeDisabled();
    fireEvent.click(retry);

    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync.mock.calls[0][0]).toEqual(originalRequest);
  });

  it('keeps a running retry disabled in external readonly mode', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    const originalRequest = {
      request_id: '56565656-5656-4565-8565-565656565656', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1',
    });
    operationQuery.data = operation({ status: 'running', stage: 'activation' });

    renderPanel({ readonly: true });

    const retry = await screen.findByTestId('basic-recording-retry');
    expect(retry).toBeDisabled();
    fireEvent.click(retry);
    expect(startAsync).not.toHaveBeenCalled();
  });

  it('single-flights a running retry while the original request replay is pending', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    const originalRequest = {
      request_id: '57575757-5757-4575-8575-575757575757', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1',
    });
    operationQuery.data = operation({ status: 'running', stage: 'activation' });
    let release!: (value: RecordingStartOperation) => void;
    startAsync.mockReturnValue(new Promise<RecordingStartOperation>((resolve) => { release = resolve; }));

    renderPanel();

    const retry = await screen.findByTestId('basic-recording-retry');
    fireEvent.click(retry);
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(retry).toBeDisabled();
    fireEvent.click(retry);
    expect(startAsync).toHaveBeenCalledTimes(1);
    await act(async () => {
      release(operation({ status: 'running', stage: 'activation' }));
    });
  });

  it('disables a running retry when the group revision no longer matches', async () => {
    const current = { ...canonicalManaged, revision: 'rev-current' };
    groupsQuery.data = { ...groupsQuery.data, groups: [current] };
    const originalRequest = {
      request_id: '58585858-5858-4585-8585-585858585858', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: current.id, expected_group_revision: 'rev-before', expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({
      scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1',
    });
    operationQuery.data = operation({
      status: 'running', stage: 'activation',
      groups: [{ group_id: current.id, group_revision: 'rev-before', saved: true, ready: true, applied: false }],
    });

    renderPanel();

    const retry = await screen.findByTestId('basic-recording-retry');
    expect(retry).toBeDisabled();
    fireEvent.click(retry);
    expect(startAsync).not.toHaveBeenCalled();
  });

  it('shows persisted members and blocks newly mapped values until Advanced saves them', async () => {
    const savedSubset = { ...canonicalManaged, members: canonicalManaged.members.slice(0, 2) };
    groupsQuery.data = { ...groupsQuery.data, groups: [savedSubset] };
    const onOpenAdvanced = vi.fn();
    renderPanel({ onOpenAdvanced });
    const members = await screen.findByTestId('basic-recording-selected-members');
    expect(members).toHaveTextContent('line.temperature');
    expect(members).toHaveTextContent('line.pressure');
    expect(members).not.toHaveTextContent('line.batch');
    expect(await screen.findByTestId('basic-recording-blocked')).toHaveTextContent('step4.basic.blocked.members_changed');
    expect(screen.getByTestId('basic-recording-start')).toBeDisabled();
    fireEvent.click(screen.getByTestId('basic-recording-members-review'));
    expect(onOpenAdvanced).toHaveBeenCalledTimes(1);
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
  });

  it('keeps the start fence when canonical group data appears during ensure', async () => {
    let release: (value: { workspace_revision: string; group: WriteGroup }) => void = () => undefined;
    ensureAsync.mockReturnValue(new Promise((resolve) => { release = resolve; }));
    const rendered = renderPanel();
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(ensureAsync).toHaveBeenCalledTimes(1));
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    rendered.rerender(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} />
      </QueryClientProvider>,
    );
    const start = await screen.findByTestId('basic-recording-start');
    expect(start).toBeDisabled();
    fireEvent.click(start);
    expect(ensureAsync).toHaveBeenCalledTimes(1);
    release({ workspace_revision: 'wrev-2', group: canonicalManaged });
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
  });

  it('keeps activation revisions when Share is disabled but hydration supplied a barrier token', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    renderPanel({
      shareStatus: {
        workspace_id: 'ws-1', enabled: false, configured_enabled: false, port: 502, address: '127.0.0.1',
        bind_state: 'disabled', mapping_count: 0, hydration_state: 'ready', readiness: true,
        readiness_token: 'hydrated-token', settings_revision: 'settings-2', workspace_revision: 'wrev-hydrated',
      },
    });
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toMatchObject({
      readiness_token: 'hydrated-token', settings_revision: 'settings-2', workspace_revision: 'wrev-hydrated',
    });
  });

  it('creates a new start intent after external schema preparation advances workspace revision', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    const originalRequest = {
      request_id: 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1' });
    operationQuery.data = operation({ status: 'failed', stage: 'readiness', reason: 'preparation_required', next_action: 'prepare_schema', setup_revision: 'wrev-1' });
    const rendered = renderPanel();
    groupsQuery.data = { ...groupsQuery.data, workspace_revision: 'wrev-2' };
    rendered.rerender(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} />
      </QueryClientProvider>,
    );
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toMatchObject({
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-2',
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision }],
    });
    expect(startAsync.mock.calls[0][0].request_id).not.toBe(originalRequest.request_id);
  });

  it('reuses an own Apply operation identity after its workspace and group revisions advance', async () => {
    const current = { ...canonicalManaged, revision: 'rev-applied', applied_revision: 'applied-1' };
    groupsQuery.data = { ...groupsQuery.data, workspace_revision: 'wrev-2', groups: [current] };
    const originalRequest = {
      request_id: 'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: current.id, expected_group_revision: 'rev-before', expected_connector_revision: 'crev-1' }],
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1' });
    operationQuery.data = operation({
      status: 'succeeded', setup_revision: 'wrev-2',
      groups: [{ group_id: current.id, group_revision: current.revision, applied_revision: current.applied_revision, saved: true, ready: true, applied: true }],
    });
    renderPanel();
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toEqual(originalRequest);
  });

  it('keeps the fresh submitted result when a stale operation GET arrives', async () => {
    groupsQuery.data = { ...groupsQuery.data, workspace_revision: 'wrev-2', groups: [canonicalManaged] };
    const fresh = operation({ status: 'succeeded', setup_revision: 'wrev-2', updated_at: '2026-10-04T00:00:02Z' });
    startAsync.mockResolvedValueOnce(fresh);
    const onCommit = vi.fn();
    const rendered = renderPanel({ onCommit });
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(screen.getByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.succeeded'));
    operationQuery.data = operation({ status: 'pending', setup_revision: 'wrev-1', updated_at: '2026-10-04T00:00:00Z' });
    rendered.rerender(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} onCommit={onCommit} />
      </QueryClientProvider>,
    );
    expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.succeeded');
    expect(screen.getByTestId('basic-recording-continue-runtime')).toBeInTheDocument();
  });

  it('replays an unresolved persisted request after own Apply changed group and workspace revisions', async () => {
    const current = { ...canonicalManaged, revision: 'rev-applied', applied_revision: 'applied-1' };
    groupsQuery.data = { ...groupsQuery.data, workspace_revision: 'wrev-2', groups: [current] };
    const originalRequest = {
      request_id: 'cccccccc-cccc-4ccc-8ccc-cccccccccccc', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: current.id, expected_group_revision: 'rev-before', expected_connector_revision: 'crev-1' }],
      readiness_token: 'share-before', settings_revision: 'settings-before', workspace_revision: 'share-wrev-before',
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest });
    renderPanel({
      shareStatus: {
        workspace_id: 'ws-1', enabled: false, configured_enabled: false, port: 502, address: '127.0.0.1',
        bind_state: 'disabled', mapping_count: 0, hydration_state: 'ready', readiness: true,
        readiness_token: 'share-after', settings_revision: 'settings-after', workspace_revision: 'share-wrev-after',
      },
    });
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toEqual(originalRequest);
  });

  it('replays an unresolved Share-only request after a restart token rotates', async () => {
    const originalRequest = {
      request_id: 'dddddddd-dddd-4ddd-8ddd-dddddddddddd', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'], groups: [],
      readiness_token: 'old-token', settings_revision: 'settings-old', workspace_revision: 'wrev-old',
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest });
    renderPanel({
      state: shareOnlyState(),
      shareStatus: shareOnlyStatus({ readiness_token: 'new-token', settings_revision: 'settings-new', workspace_revision: 'wrev-new' }),
    });
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toEqual(originalRequest);
  });

  it('keeps Share-only start free of database group and ensure calls', async () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    renderPanel({
      state: shareOnlyState(),
      shareStatus: shareOnlyStatus({ workspace_revision: 'wrev-1' }),
    });
    expect(await screen.findByTestId('basic-recording-destination')).toHaveTextContent('step4.basic.share_only_destination');
    fireEvent.click(screen.getByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync.mock.calls[0][0]).toMatchObject({ device_ids: ['dev-1'], groups: [], readiness_token: 'share-token' });
  });

  it('allows fresh Share-only start with a loaded empty workspace revision', async () => {
    groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: '', groups: [] };
    renderPanel({
      state: shareOnlyState(),
      shareStatus: shareOnlyStatus(),
    });
    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync.mock.calls[0][0]).toMatchObject({ expected_workspace_revision: '', device_ids: ['dev-1'], groups: [] });
  });

  it('keeps the database path blocked when the loaded workspace revision is empty', async () => {
    groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: '', groups: [] };
    renderPanel();
    expect(await screen.findByTestId('basic-recording-blocked')).toHaveTextContent('step4.basic.blocked.workspace_revision');
    expect(screen.getByTestId('basic-recording-start')).toBeDisabled();
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
  });

  it.each([
    undefined,
    { workspace_id: 'other-workspace', workspace_revision: '', groups: [] },
  ])('blocks Share-only start when the groups snapshot is missing or belongs to another workspace', async (snapshot) => {
    groupsQuery.data = snapshot;
    renderPanel({
      state: shareOnlyState(),
      shareStatus: shareOnlyStatus(),
    });
    expect(await screen.findByTestId('basic-recording-blocked')).toHaveTextContent('step4.basic.blocked.workspace_revision');
    expect(screen.getByTestId('basic-recording-start')).toBeDisabled();
    expect(ensureAsync).not.toHaveBeenCalled();
    expect(startAsync).not.toHaveBeenCalled();
  });

  it.each([
    { settings_revision: 'settings-new', workspace_revision: 'share-old' },
    { settings_revision: 'settings-old', workspace_revision: 'share-new' },
  ])('does not keep Share-only activation evidence after one persisted revision changes (%s)', async ({ settings_revision, workspace_revision }) => {
    groupsQuery.data = { ...groupsQuery.data, workspace_revision: '', groups: [] };
    const originalRequest = {
      request_id: 'eeeeeeee-eeee-4eee-8eee-eeeeeeeeeeee', workspace_id: 'ws-1',
      expected_workspace_revision: '', device_ids: ['dev-1'], groups: [],
      readiness_token: 'old-token', settings_revision: 'settings-old', workspace_revision: 'share-old',
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1' });
    operationQuery.data = operation({
      status: 'succeeded', setup_revision: '', groups: [],
    });
    startAsync.mockResolvedValueOnce(operation({ status: 'pending', setup_revision: '', groups: [] }));
    renderPanel({
      state: shareOnlyState(),
      shareStatus: shareOnlyStatus({ readiness_token: 'new-token', settings_revision, workspace_revision }),
    });

    expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.succeeded');
    expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toMatchObject({
      request_id: expect.not.stringMatching(originalRequest.request_id),
      expected_workspace_revision: '', readiness_token: 'new-token',
      settings_revision, workspace_revision, groups: [],
    });
  });

  it.each(['succeeded', 'partial'] as const)('reuses the original Share-only request when only the readiness token rotates (%s)', async (status) => {
    groupsQuery.data = { ...groupsQuery.data, workspace_revision: '', groups: [] };
    const originalRequest = {
      request_id: 'fefefefe-fefe-4fef-8fef-fefefefefefe', workspace_id: 'ws-1',
      expected_workspace_revision: '', device_ids: ['dev-1'], groups: [],
      readiness_token: 'old-token', settings_revision: 'settings-same', workspace_revision: 'share-same',
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1' });
    operationQuery.data = operation({ status, setup_revision: '', groups: [] });
    startAsync.mockResolvedValueOnce(operation({ status: 'pending', setup_revision: '', groups: [] }));
    renderPanel({
      state: shareOnlyState(),
      shareStatus: shareOnlyStatus({ readiness_token: 'new-token', settings_revision: 'settings-same', workspace_revision: 'share-same' }),
    });

    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toEqual(originalRequest);
  });

  it.each(['succeeded', 'partial'] as const)('reuses the original database request when only the readiness token rotates (%s)', async (status) => {
    groupsQuery.data = { ...groupsQuery.data, groups: [canonicalManaged] };
    const originalRequest = {
      request_id: 'abababab-abab-4bab-8bab-abababababab', workspace_id: 'ws-1',
      expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'],
      groups: [{ group_id: canonicalManaged.id, expected_group_revision: canonicalManaged.revision, expected_connector_revision: 'crev-1' }],
      readiness_token: 'old-token', settings_revision: 'settings-same', workspace_revision: 'share-same',
    };
    saveBasicRecordingIntent({ scope_key: 'ws-1/dev-1/managed-recording', request: originalRequest, operation_id: 'operation-1' });
    operationQuery.data = operation({ status, setup_revision: 'wrev-1' });
    startAsync.mockResolvedValueOnce(operation({ status: 'pending', setup_revision: 'wrev-1' }));
    renderPanel({
      shareStatus: {
        workspace_id: 'ws-1', enabled: false, configured_enabled: false, port: 502, address: '127.0.0.1',
        bind_state: 'disabled', mapping_count: 0, hydration_state: 'ready', readiness: true,
        readiness_token: 'new-token', settings_revision: 'settings-same', workspace_revision: 'share-same',
      },
    });

    fireEvent.click(await screen.findByTestId('basic-recording-start'));
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0]).toEqual(originalRequest);
  });
});
