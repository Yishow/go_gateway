import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, renderHook, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { useCircuitProfileStart } from '@/features/datalink/workbench-v2/steps/step4/useCircuitProfileStart';
import { basicRecordingScopeKey, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { studioV2RecordingStartAPI as startAPI } from '@/services/studioV2RecordingStart';
import { studioV2WorkspaceWriteGroupsAPI as api } from '@/services/studioV2WorkspaceWriteGroups';
import type { SchemaOperation, SchemaPreviewToken } from '@/types/recordingPlan';
import type { WriteGroup, WriteGroupCreateRequest } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

// Real hooks, the real Step 4 host and the real schema panel; only the service boundary is simulated (frontend
// evidence, not SQL). The backend keeps its own groups so that what Basic shows after a save is what was saved.
const backend = vi.hoisted(() => ({ revision: 1, groups: [] as WriteGroup[], tokens: new Map<string, SchemaPreviewToken>() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(), ensureCircuitProfile: vi.fn(),
  },
}));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), ensureBasic: vi.fn() } }));

const now = '2026-10-24T00:00:00Z';
const revision = () => `workspace-${backend.revision}`;
/** Every call that would change the running writer, apply a group or start recording. Schema calls are checked on their own. */
const WRITES = [api.create, api.update, api.apply, api.disable, api.remove, startAPI.start, startAPI.ensureBasic];
const expectNoWriterWrites = () => WRITES.forEach((write) => expect(write).not.toHaveBeenCalled());

/** The group Basic has always owned: a managed wide group of one device, applied and recording every 60 seconds. */
const legacyWide = (): WriteGroup => savedGroup({
  id: 'wide-1', name: 'PLC A raw', basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1', status: 'running',
  revision: 'rev-wide', applied_revision: 'applied-wide',
  members: ['temperature', 'pressure', 'batch'].map((column, index) => ({
    device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}`, source_revision: `s${index}`, mapping_revision: `m${index}`, target_column: column, required: index !== 2,
  })),
  destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: '', table_name: 'gw_basic_dev_1', storage_strategy: 'managed' },
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row' },
});
const quietDelivery = {
  group_id: 'wide-1', intake: { state: 'active' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
  last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
};

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  backend.revision = 1; backend.groups = [legacyWide()]; backend.tokens.clear();
  vi.mocked(api.list).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: revision(), groups: structuredClone(backend.groups) }));
  vi.mocked(api.delivery).mockResolvedValue(quietDelivery as never);
  vi.mocked(api.readiness).mockResolvedValue({ workspace_id: 'ws-1', group_id: 'x', group_revision: 'r', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] } as never);
  vi.mocked(api.ensureCircuitProfile).mockImplementation(async (request: WriteGroupCreateRequest) => {
    expect(request.expected_workspace_revision).toBe(revision());
    const kind = request.group.recording_profile?.kind;
    const group = savedGroup({
      ...request.group, id: `group-${kind}`, revision: 'group-rev-1', status: 'draft', applied_revision: '',
      members: request.group.members.map((member) => ({ ...member, source_revision: 'source-1', mapping_revision: 'mapping-1' })),
      destination: { ...request.group.destination, database: '/fixture.db', table_name: 'gw_long_fixture' },
    });
    backend.groups.push(group); backend.revision += 1;
    return { workspace_revision: revision(), group: structuredClone(group) };
  });
  vi.mocked(api.schemaPreview).mockImplementation(async (id, request) => {
    const group = backend.groups.find((item) => item.id === id);
    if (!group) throw new Error('missing fixture group');
    const token: SchemaPreviewToken = {
      token: `tok-${id}`, operation_id: `op-${id}`, workspace_id: 'ws-1', workspace_revision: request.expected_workspace_revision, plan_id: id, plan_revision: group.revision,
      connector_id: 'conn-1', table_prefix: 'gw_long_fixture', statements: ['CREATE TABLE fixture (...)'], tables: [], digest: 'a'.repeat(64), created_at: now,
      expires_at: '2026-10-25T00:00:00Z', recreate_missing: true,
    };
    backend.tokens.set(token.token, token);
    return token;
  });
  vi.mocked(api.schemaApply).mockImplementation(async (id, request) => {
    const group = backend.groups.find((item) => item.id === id);
    if (!group) throw new Error('missing fixture group');
    group.destination.schema_revision = group.revision; group.destination.schema_digest = 'a'.repeat(64); backend.revision += 1;
    return { operation_id: request.operation_id, status: 'succeeded', executed_statements: 1, created_at: now, updated_at: now } as SchemaOperation;
  });
});

function renderHost(over: { readonly?: boolean } = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const element = (
    <QueryClientProvider client={client}>
      {over.readonly
        ? <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly />
        : <Step4Database state={savedGroupState()} dispatch={vi.fn()} workspaceId="ws-1" />}
    </QueryClientProvider>
  );
  const view = render(element);
  return { ...view, reload: () => { view.unmount(); return render(element); } };
}
const openPrepare = async () => {
  fireEvent.click(await screen.findByTestId('legacy-migration-open'));
  return screen.findByTestId('circuit-profile-basic-panel');
};
const chooseExplicitly = async (name = '遷移迴路') => {
  fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: name } });
  fireEvent.click(screen.getByTestId('circuit-profile-toggle-instantaneous'));
  fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
  fireEvent.click(screen.getByTestId('circuit-profile-member-inst-tag-0'));
};

describe('A saved legacy wide group offers an explicit way to prepare long targets', () => {
  it('shows the entry only inside the legacy panel, and opens the same profile panel in a prepare-only context', async () => {
    renderHost();
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(60);
    expect(screen.queryByTestId('circuit-profile-basic-panel')).not.toBeInTheDocument();
    const entry = screen.getByTestId('legacy-migration-open');
    expect(entry).toBeEnabled();
    expect(entry).toHaveTextContent('step4.basic.migration.open');

    fireEvent.click(entry);
    expect(await screen.findByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare-only')).toHaveTextContent('step4.basic.migration.banner');
    // Two effective panels are never shown together, and the prepare context has no recording or legacy period control.
    expect(screen.queryByTestId('basic-recording-panel')).not.toBeInTheDocument();
    for (const id of ['basic-recording-interval', 'basic-recording-completeness', 'circuit-profile-start', 'basic-recording-start', 'basic-recording-retry', 'basic-recording-check', 'basic-recording-continue-runtime']) {
      expect(screen.queryByTestId(id), id).not.toBeInTheDocument();
    }
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toBeInTheDocument();
    expectNoWriterWrites();
  });

  it('starts from an empty explicit choice: no source, no profile and no period is taken from the old group', async () => {
    renderHost();
    await openPrepare();
    expect(screen.getByTestId('circuit-profile-toggle-instantaneous')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-toggle-cumulative')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('');
    expect(screen.getByTestId('circuit-profile-interval-cumulative')).toHaveValue('');
    expandSources();
    for (const box of screen.getAllByRole('checkbox', { name: /^line\./ })) expect(box).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(api.ensureCircuitProfile).not.toHaveBeenCalled();
  });

  it('saves only the chosen canonical target, confirms its schema on request and never touches the old writer', async () => {
    const seeded = structuredClone(backend.groups[0]);
    renderHost();
    await openPrepare();
    expandSources();
    await chooseExplicitly();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    const prepared = within(await screen.findByTestId('circuit-prepared-group-instantaneous'));

    expect(api.ensureCircuitProfile).toHaveBeenCalledTimes(1);
    const sent = vi.mocked(api.ensureCircuitProfile).mock.calls[0][0].group;
    expect(sent.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' });
    expect(sent.row_policy.interval_seconds).toBe(30);
    expect(sent.members.map((member) => member.tag_id)).toEqual(['tag-0']);
    expect(sent.layout_version).toBe('managed_long_v1');
    expect(sent).not.toHaveProperty('basic_managed_device_id');
    expect(sent.destination).toMatchObject({ connector_id: 'conn-1', connector_revision: 'crev-1', storage_strategy: 'managed' });
    expect(api.schemaApply).not.toHaveBeenCalled();

    fireEvent.click(await prepared.findByTestId('group-schema-apply'));
    await waitFor(() => expect(api.schemaApply).toHaveBeenCalledTimes(1));
    expect(vi.mocked(api.schemaApply).mock.calls[0][0]).toBe('group-instantaneous');
    expect(vi.mocked(api.schemaPreview).mock.calls.every(([id]) => id !== 'wide-1')).toBe(true);
    expectNoWriterWrites();
    expect(startAPI.operation).not.toHaveBeenCalled();
    expect(backend.groups.find((group) => group.id === 'wide-1')).toEqual(seeded);
  });

  it('returns to the unchanged legacy panel on close, keeps the saved targets as drafts and writes nothing', async () => {
    renderHost();
    await openPrepare();
    expandSources();
    await chooseExplicitly();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await screen.findByTestId('circuit-prepared-group-instantaneous');

    fireEvent.click(screen.getByTestId('circuit-profile-prepare-close'));
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-basic-panel')).not.toBeInTheDocument();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(60);
    expect(backend.groups.map((group) => group.id)).toEqual(['wide-1', 'group-instantaneous']);
    expect(backend.groups[1].applied_revision).toBe('');
    expect(backend.groups[1].status).toBe('draft');
    expectNoWriterWrites();
    expect(api.ensureCircuitProfile).toHaveBeenCalledTimes(1);
  });
});

describe('The prepare context can also prepare a new wide mean target', () => {
  it('offers the purpose choice and saves a managed_wide_profile_v1 group with an explicit mean, still starting and applying nothing', async () => {
    const seeded = structuredClone(backend.groups[0]);
    renderHost();
    await openPrepare();
    expect(screen.getByTestId('circuit-profile-purpose-wide'), 'the one purpose choice is the same as outside the prepare context').toBeInTheDocument();
    fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    expandSources();
    await chooseExplicitly();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await screen.findByTestId('circuit-prepared-group-instantaneous');

    expect(api.ensureCircuitProfile).toHaveBeenCalledTimes(1);
    const sent = vi.mocked(api.ensureCircuitProfile).mock.calls[0][0].group;
    expect(sent.layout_version).toBe('managed_wide_profile_v1');
    expect(sent.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' });
    expect(sent.members.map((member) => [member.tag_id, member.value_role])).toEqual([['tag-0', 'gauge']]);
    expect(sent).not.toHaveProperty('basic_managed_device_id');
    expectNoWriterWrites();
    expect(startAPI.operation).not.toHaveBeenCalled();
    expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
    expect(backend.groups.find((group) => group.id === 'wide-1')).toEqual(seeded);
  });
});

describe('Reloading or choosing a cumulative profile cannot change the old writer or invent a counter', () => {
  it('returns to the unchanged legacy panel after a reload and keeps the prepared target as a draft', async () => {
    const host = renderHost();
    await openPrepare();
    expandSources();
    await chooseExplicitly();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await screen.findByTestId('circuit-prepared-group-instantaneous');

    host.reload();
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-basic-panel')).not.toBeInTheDocument();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(60);
    expect(backend.groups.map((group) => [group.id, group.status, group.applied_revision])).toEqual([['wide-1', 'running', 'applied-wide'], ['group-instantaneous', 'draft', '']]);
    expectNoWriterWrites();
    expect(startAPI.operation).not.toHaveBeenCalled();
  });

  it('does not prepare a cumulative target from a raw Tag without a confirmed counter definition', async () => {
    renderHost();
    await openPrepare();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: '累積迴路' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.change(screen.getByTestId('circuit-profile-interval-cumulative'), { target: { value: '10' } });
    expandSources();
    fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-0'));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-interval-cumulative')).toHaveValue('10'));
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(api.ensureCircuitProfile).not.toHaveBeenCalled();
    expectNoWriterWrites();
  });
});

describe('The prepare context never restores or resends a recording start', () => {
  it('ignores an unresolved start intent left by an earlier session and still lets the operator prepare', async () => {
    const circuitId = '遷移迴路';
    const stale = {
      request_id: 'req-stale', workspace_id: 'ws-1', expected_workspace_revision: 'workspace-0', device_ids: ['dev-1'],
      groups: [{ group_id: 'group-instantaneous', expected_group_revision: 'group-rev-1', expected_connector_revision: 'crev-1' }],
    };
    for (const scope of [circuitId, '', 'cir-stale']) {
      saveBasicRecordingIntent({ scope_key: basicRecordingScopeKey('ws-1', 'dev-1', encodeURIComponent(JSON.stringify([scope, 'conn-1', 'crev-1', '']))), request: stale, operation_id: 'op-stale' });
    }
    saveBasicRecordingIntent({ scope_key: basicRecordingScopeKey('ws-1', 'dev-1', 'managed-recording'), request: stale, operation_id: 'op-stale' });
    renderHost();
    await openPrepare();
    expandSources();
    await chooseExplicitly(circuitId);
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await screen.findByTestId('circuit-prepared-group-instantaneous');
    expect(startAPI.operation).not.toHaveBeenCalled();
    expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
    expect(screen.queryByTestId('basic-recording-check')).not.toBeInTheDocument();
    expectNoWriterWrites();
  });

  it('keeps the start hook inert in prepare-only mode even when a stored intent exists for the very same scope', async () => {
    const scopeKey = basicRecordingScopeKey('ws-1', 'dev-1', 'scope-x');
    saveBasicRecordingIntent({
      scope_key: scopeKey, operation_id: 'op-stale',
      request: { request_id: 'req', workspace_id: 'ws-1', expected_workspace_revision: 'w', device_ids: ['dev-1'], groups: [] },
    });
    const prepare = vi.fn(async () => ({ groups: [], workspaceRevision: 'w' }));
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    const wrapper = ({ children }: { children: React.ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
    const { result } = renderHook(() => useCircuitProfileStart({
      scopeKey, canStart: true, readonly: false, workspaceId: 'ws-1', deviceId: 'dev-1', connectorRevision: 'crev-1', prepare, prepareOnly: true,
    }), { wrapper });
    expect(result.current.unresolved).toBe(false);
    expect(result.current.operationId).toBeUndefined();
    expect(result.current.canRetrySameRequest).toBe(false);
    await result.current.start(false);
    await result.current.start(true);
    expect(startAPI.start).not.toHaveBeenCalled();
    expect(startAPI.operation).not.toHaveBeenCalled();
    await result.current.prepare();
    expect(prepare).toHaveBeenCalledTimes(1);
  });
});

describe('The entry respects the same locks as the legacy panel', () => {
  it('is disabled while the page is read-only', async () => {
    renderHost({ readonly: true });
    expect(await screen.findByTestId('legacy-migration-open')).toBeDisabled();
  });

  it.each([
    ['only Modbus Share is configured', { enabled: true, configured_enabled: true, port: 5020, address: '127.0.0.1', bind_state: 'pass' as const, mapping_count: 3 }],
    ['no destination is saved and Share is off', undefined],
  ])('is not offered when %s', async (_name, shareStatus) => {
    const state = savedGroupState();
    const unsaved = { ...state, db: { ...state.db, connector: { ...state.db.connector, persisted: false, save_state: 'idle' as const, connector_id: undefined, identity_revision: undefined } } };
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <BasicRecordingPanel state={unsaved} workspaceId="ws-1" readonly={false} shareStatus={shareStatus} />
      </QueryClientProvider>,
    );
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    // The saved wide group must actually have arrived, otherwise "absent" proves nothing.
    await waitFor(() => expect(api.list).toHaveBeenCalled());
    await act(async () => { await Promise.all(vi.mocked(api.list).mock.results.map((result) => result.value)); });
    expect(screen.queryByTestId('legacy-migration-open')).not.toBeInTheDocument();
    expectNoWriterWrites();
  });
});
