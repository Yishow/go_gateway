import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { guidedReadiness, guidedLifecycle, guidedPreview, startGuided } from '../../fixtures/guidedRecordingPreparation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { ensureOrUpdateCircuitProfileGroups, type EnsureOrUpdateCircuitProfilesInput } from '@/features/datalink/workbench-v2/steps/step4/circuitProfileActions';
import type { CircuitProfileCapabilities } from '@/features/datalink/workbench-v2/state/circuitProfileCapabilities';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { StudioV2MySQLRecordingCapability } from '@/types/studioV2MySQLRecordingCapability';
import type { StudioV2SQLServerRecordingCapability } from '@/types/studioV2SQLServerRecordingCapability';
import type { WriteGroup, WriteGroupCreateRequest, WriteGroupUpdateRequest } from '@/types/studioV2WriteGroup';
import productionI18n from '@/i18n/config';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] },
  isLoading: false, isError: false, refetch: vi.fn(),
};
const ensure = vi.fn();
const update = vi.fn();
const start = vi.fn();
const applySchema = vi.fn();
const previewSchema = vi.fn();
let canonicalClient: QueryClient;
const evidenceRefetch = vi.fn();
const evidence = new Map<string, { readiness: ReturnType<typeof guidedReadiness>; lifecycle: ReturnType<typeof guidedLifecycle> }>();
function currentEvidence(id: string) {
  const group = groupsQuery.data.groups.find(value => value.id === id);
  if (!group) return undefined;
  const key = JSON.stringify([groupsQuery.data.workspace_revision, group]);
  if (!evidence.has(key)) evidence.set(key, { readiness: guidedReadiness(group, groupsQuery.data.workspace_revision), lifecycle: guidedLifecycle(group, groupsQuery.data.workspace_revision) });
  return evidence.get(key)!;
}
function publishGroups() {
  canonicalClient?.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data);
  return { data: groupsQuery.data };
}
vi.mock('@/hooks/datalink/useStudioV2WriteGroupLifecycle', () => ({
  useWriteGroupLifecycleQuery: (_workspace: string, id: string) => ({ data: currentEvidence(id)?.lifecycle, isLoading: false, isFetching: false, isError: false, refetch: evidenceRefetch }),
}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useWriteGroupReadinessQuery: (id: string) => ({ data: currentEvidence(id)?.readiness, isLoading: false, isFetching: false, isError: false, refetch: evidenceRefetch }),
  useWriteGroupSchemaPreviewMutation: () => ({ mutateAsync: previewSchema, isPending: false }),
  useWriteGroupSchemaOperationQuery: () => ({ data: undefined, isError: false, refetch: evidenceRefetch }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensure, isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: update, isPending: false }),
  useWriteGroupSchemaApplyMutation: () => ({ mutateAsync: applySchema, isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => ({ data: [], isLoading: false, isError: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: start, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isLoading: false, isError: false }),
}));
function mysqlCapability(family = 'mysql'): StudioV2MySQLRecordingCapability {
  return {
    connector_id: 'conn-1', connector_revision: 'crev-1', database: 'recording_test',
    server_family: family, server_version: 'verified-test-version', capability_revision: 'cap-1',
    connected: true, metadata_status: 'verified', transactional_engine: true,
    supports_managed_long: true, supports_native_partition: true,
    supports_group_test_write: false, supports_missing_repair: false, reasons: [],
  };
}
function sqlserverCapability(): StudioV2SQLServerRecordingCapability {
  return {
    connector_id: 'conn-1', connector_revision: 'crev-1', database: 'recording_test', schema: 'dbo',
    server_version: '16.0.4165.4', server_edition: 'Developer Edition', engine_edition: 3, compatibility_level: 160,
    connected: true, metadata_status: 'verified', capability_revision: 'cap-1',
    supports_managed_long: true, supports_native_partition: true, supports_group_test_write: false,
    supports_missing_repair: false, reasons: [], permissions: {
      verified: true, database_view_definition: true, schema_view_definition: true, create_table: true,
      alter_schema: true, select_schema: true, insert_schema: true, delete_schema: true, alter_any_dataspace: true,
    },
  };
}
function stateFor(kind: 'mysql' | 'sqlserver' | 'sqlite' | 'postgres') {
  const state = savedGroupState();
  state.db.connector = { ...state.db.connector, kind, connector_id: 'conn-1', identity_revision: 'crev-1' };
  return state;
}
async function show(state = stateFor('mysql'), readonly = false, lng = 'en') {
  const i18n = createInstance();
  await i18n.init({ lng, resources: productionI18n.options.resources });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  canonicalClient = client; client.setQueryDefaults(studioV2WorkspaceKeys.writeGroups(), { structuralSharing: false }); publishGroups();
  const panel = (next = state) => <I18nextProvider i18n={i18n}><QueryClientProvider client={client}>
    <CircuitProfileBasicPanel state={next} workspaceId="ws-1" readonly={readonly} />
  </QueryClientProvider></I18nextProvider>;
  const view = render(panel());
  return { ...view, client, changeState: (next: typeof state) => view.rerender(panel(next)) };
}
function actionInput(kind = 'mysql'): EnsureOrUpdateCircuitProfilesInput {
  return {
    workspaceId: 'ws-1', initialWorkspaceRevision: 'wrev-1', deviceName: 'PLC A', connectorKind: kind,
    enabledProfiles: [{ profile: 'instantaneous' as const, intervalValue: 15, intervalSeconds: 15,
      candidates: buildGroupCandidates(savedGroupState()).candidates }],
    storageMode: 'plain' as const, recreateMissing: true, purpose: 'wide' as const,
    circuitId: 'cir-a', circuitName: 'Draft A',
    destination: { connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: '', storage_strategy: 'managed' as const },
    measurementsByPointKey: new Map(), savedLongGroups: [], preparedGroups: [],
    ensureCircuitProfile: ensure, updateWriteGroup: update,
  };
}
function capabilitiesFor(kind: string): CircuitProfileCapabilities {
  const snapshot = { isLoading: false, isFetching: false, isError: false };
  return kind === 'sqlserver' ? { sqlserver: { ...snapshot, data: sqlserverCapability() } }
    : { mysql: { ...snapshot, data: mysqlCapability() } };
}
function seedDraft(purpose = 'long', storageMode = 'plain') {
  window.sessionStorage.setItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1', JSON.stringify({
    circuit_id: 'cir-a', circuit_name: 'Draft A', instantaneous_interval_input: '37', purpose, storage_mode: storageMode,
  }));
}
function deferred() {
  let resolve!: (value: Awaited<ReturnType<typeof studioV2DatalinkApi.get>>) => void;
  let reject!: (error: Error) => void;
  const promise = new Promise<Awaited<ReturnType<typeof studioV2DatalinkApi.get>>>((yes, no) => { resolve = yes; reject = no; });
  return { promise, resolve, reject };
}
// Recovery observation is independent of the capability refresh being deferred.
function mockCapabilityRefresh(refresh: ReturnType<typeof deferred>) {
  let capabilityReads = 0;
  vi.spyOn(studioV2DatalinkApi, 'get').mockImplementation((url, config) => {
    if (url === '/studio-v2/workspace/database-operations/test-writes') {
      return Promise.resolve({ data: { success: true, data: { workspace_id: 'ws-1', complete: true, operations: [] } } });
    }
    if (url === '/studio-v2/workspace/recording-plans/capabilities') {
      expect(config?.params).toEqual({ kind: 'mysql', connector_id: 'conn-1', expected_connector_revision: 'crev-1' });
      capabilityReads += 1;
      if (capabilityReads === 1) return Promise.resolve({ data: { data: mysqlCapability() } });
      if (capabilityReads === 2) return refresh.promise;
      return Promise.reject(new Error('Unexpected extra capability GET'));
    }
    return Promise.reject(new Error(`Unexpected GET: ${url}`));
  });
}
function noWrites() {
  for (const fn of [ensure, update, start, applySchema, previewSchema]) expect(fn).not.toHaveBeenCalled();
}
beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockImplementation(async () => publishGroups());
  ensure.mockImplementation(async ({ request }: { request: WriteGroupCreateRequest }) => { const result = {
    workspace_revision: 'wrev-2', group: savedGroup({ ...request.group, id: 'new-profile', revision: 'new-revision',
      destination: { ...request.group.destination, database: 'recording_test' },
      members: request.group.members.map((member) => ({ ...member, source_revision: 's1', mapping_revision: 'm1' })) }),
  }; groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: result.workspace_revision, groups: [result.group] }; publishGroups(); return result; });
  previewSchema.mockImplementation(async ({ id }) => guidedPreview(groupsQuery.data.groups.find(group => group.id === id)!, groupsQuery.data.workspace_revision));
  applySchema.mockImplementation(async ({ id, request }) => {
    const group = groupsQuery.data.groups.find(value => value.id === id)!; const preview = guidedPreview(group, groupsQuery.data.workspace_revision);
    expect(request).toMatchObject({ expected_workspace_revision: groupsQuery.data.workspace_revision, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision, token: preview.token, operation_id: preview.operation_id });
    groupsQuery.data = { ...groupsQuery.data, groups: groupsQuery.data.groups.map(value => value.id === id ? { ...value, destination: { ...value.destination, schema_revision: value.revision, schema_digest: 'a'.repeat(64) } } : value) }; publishGroups();
    return { operation_id: request.operation_id, status: 'succeeded', executed_statements: 1, created_at: '', updated_at: '' };
  });
});
afterEach(() => vi.restoreAllMocks());
describe('destination layout choices use the saved connector capability', () => {
  it.each(['mysql', 'mariadb', 'sqlserver'] as const)('does not offer executable wide for verified %s managed-long support', async (family) => {
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValue({ data: { data: family === 'sqlserver' ? sqlserverCapability() : mysqlCapability(family) } });
    await show(stateFor(family === 'sqlserver' ? 'sqlserver' : 'mysql'));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    const wide = screen.getByTestId('circuit-profile-purpose-wide');
    expect(wide).toBeDisabled();
    expect(wide).toHaveAccessibleDescription('This destination supports selected-measurement recording only; wide tables are unavailable.');
    fireEvent.click(wide);
    expect(screen.getByTestId('circuit-profile-purpose-long')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-purpose-reason-wide')).toHaveTextContent('This destination supports selected-measurement recording only; wide tables are unavailable.');
    noWrites();
  });
  it.each(['mysql', 'sqlserver'])('rejects a programmatic %s wide prepare before any mutation', async (kind) => {
    await expect(ensureOrUpdateCircuitProfileGroups({ ...actionInput(kind), capabilities: capabilitiesFor(kind) }))
      .rejects.toThrow('Recording layout is unavailable for this destination.');
    noWrites();
  });
  it.each(['mysql', 'mariadb', 'sqlserver'] as const)('prepares only the explicit supported long draft for %s, without starting or applying schema', async (family) => {
    const kind = family === 'sqlserver' ? 'sqlserver' : 'mysql';
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValue({ data: { data: kind === 'sqlserver' ? sqlserverCapability() : mysqlCapability(family) } });
    await show(stateFor(kind));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expect(screen.getByText('An available layout does not mean the connection, schema, or recording is ready.')).toBeVisible();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensure).toHaveBeenCalledOnce());
    expect(ensure.mock.calls[0][0].request).toMatchObject({ expected_connector_revision: 'crev-1', group: {
      layout_version: 'managed_long_v1', storage_mode: 'plain', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' },
    } });
    expect(start).not.toHaveBeenCalled(); expect(applySchema).not.toHaveBeenCalled();
  });
  it.each(['draft', 'saved'] as const)('keeps an unsupported %s wide purpose visible and unchanged with an actionable explanation', async (origin) => {
    seedDraft('wide');
    const group = longProfileGroup({ id: 'saved-wide', circuitId: 'cir-a', circuitName: 'Draft A', interval: 37, points: [0, 1] });
    group.layout_version = 'managed_wide_profile_v1';
    if (origin === 'saved') groupsQuery.data.groups = [group];
    const original = structuredClone(group);
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValue({ data: { data: mysqlCapability() } });
    const view = await show();
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-layout-blocked')).toHaveTextContent('wide tables are unavailable');
    expect(screen.getByRole('button', { name: 'Open destination settings' })).toBeEnabled();
    fireEvent.click(screen.getByRole('button', { name: 'Open destination settings' }));
    // One guided primary: unavailable scope may be inspected, never apply or start.
    expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-schema-apply')).not.toBeInTheDocument();
    const prepare = screen.queryByTestId('circuit-profile-prepare'); if (prepare) expect(prepare).toBeDisabled();
    const inspect = screen.queryByTestId('recording-next-action'); if (inspect) { expect(inspect).toHaveTextContent('Inspect this scope'); fireEvent.click(inspect); }
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('37');
    expect(group).toEqual(original);
    expect(JSON.parse(window.sessionStorage.getItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1')!).purpose).toBe('wide');
    noWrites();
  });
  it.each(['id', 'revision'] as const)('ignores a late old capability after the connector %s changes', async (change) => {
    const old = deferred(), current = deferred();
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockImplementationOnce(() => old.promise).mockImplementationOnce(() => current.promise);
    seedDraft();
    const view = await show();
    const next = stateFor('mysql');
    if (change === 'id') next.db.connector.connector_id = 'conn-2';
    else next.db.connector.identity_revision = 'crev-2';
    view.changeState(next);
    await waitFor(() => expect(get).toHaveBeenCalledTimes(2));
    expect(get.mock.calls[1][1]).toEqual({ params: { kind: 'mysql', connector_id: change === 'id' ? 'conn-2' : 'conn-1', expected_connector_revision: change === 'revision' ? 'crev-2' : 'crev-1' } });
    await act(async () => old.resolve({ data: { data: mysqlCapability() } }));
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-purpose-long')).toBeDisabled();
    await act(async () => current.resolve({ data: { data: { ...mysqlCapability(), connector_id: next.db.connector.connector_id, connector_revision: next.db.connector.identity_revision } } }));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('37');
    noWrites();
  });
  it('blocks cached capabilities during refetch and error, then recovers by read only', async () => {
    const refresh = deferred();
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce({ data: { data: mysqlCapability() } }).mockImplementationOnce(() => refresh.promise);
    seedDraft();
    const view = await show();
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    act(() => { void view.client.invalidateQueries({ queryKey: ['studio-v2-workspace', 'mysql-recording-capability'] }); });
    await waitFor(() => expect(get).toHaveBeenCalledTimes(2));
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    await act(async () => refresh.reject(new Error('password=secret; private-db.internal')));
    const error = await screen.findByTestId('circuit-profile-mysql-capability-error');
    expect(screen.getByTestId('circuit-profile-purpose-long')).toBeDisabled();
    expect(document.body).not.toHaveTextContent('private-db.internal');
    get.mockResolvedValueOnce({ data: { data: mysqlCapability() } });
    fireEvent.click(within(error).getByRole('button'));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Draft A');
    noWrites();
  });
  it.each(['mysql', 'sqlserver'] as const)('does not downgrade a cached native mode when %s native capability is unavailable', async (kind) => {
    seedDraft('long', 'native_monthly');
    const capability = kind === 'mysql' ? mysqlCapability() : sqlserverCapability();
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValue({ data: { data: { ...capability, supports_native_partition: false } } });
    const view = await show(stateFor(kind));
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    expect(screen.getByTestId('circuit-profile-storage-native-monthly')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-storage-native-monthly')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-layout-blocked')).toHaveTextContent('no saved mode has been changed');
    expect(JSON.parse(window.sessionStorage.getItem('wbv2.circuit-profile-draft.ws-1%2Fdev-1')!).storage_mode).toBe('native_monthly');
    noWrites();
  });
  it('keeps PostgreSQL Timescale long-only without declaring its backend extension ready', async () => {
    seedDraft('long', 'timescale');
    await show(stateFor('postgres'));
    expect(screen.getByTestId('circuit-profile-storage-timescale')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeDisabled();
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensure).toHaveBeenCalledOnce());
    expect(ensure.mock.calls[0][0].request.group).toMatchObject({ layout_version: 'managed_long_v1', storage_mode: 'timescale' });
    expect(start).not.toHaveBeenCalled(); expect(applySchema).not.toHaveBeenCalled();
  });
  it('keeps unsupported settings readable in readonly Traditional Chinese without mutations', async () => {
    seedDraft('wide');
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValue({ data: { data: mysqlCapability() } });
    const view = await show(stateFor('mysql'), true, 'zh-TW');
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-layout-blocked')).toHaveTextContent('此目的地僅支援所選量測記錄，不支援寬表。');
    expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
    noWrites();
  });
  it.each(['mysql', 'sqlserver'] as const)('rejects every unverified %s capability snapshot at the action boundary', async (kind) => {
    for (const condition of ['missing', 'loading', 'fetching', 'error', 'id', 'revision', 'capability-revision', 'metadata', 'disconnected', 'unsupported'] as const) {
      const capabilities = capabilitiesFor(kind);
      const snapshot = capabilities.mysql ?? capabilities.sqlserver!;
      if (condition === 'missing') snapshot.data = undefined;
      else if (condition === 'loading') snapshot.isLoading = true;
      else if (condition === 'fetching') snapshot.isFetching = true;
      else if (condition === 'error') snapshot.isError = true;
      else if (snapshot.data) {
        if (condition === 'id') snapshot.data.connector_id = 'old-connector';
        if (condition === 'revision') snapshot.data.connector_revision = 'old-revision';
        if (condition === 'capability-revision') snapshot.data.capability_revision = '';
        if (condition === 'metadata') snapshot.data.metadata_status = 'unknown';
        if (condition === 'disconnected') snapshot.data.connected = false;
        if (condition === 'unsupported') snapshot.data.supports_managed_long = false;
      }
      await expect(ensureOrUpdateCircuitProfileGroups({ ...actionInput(kind), purpose: 'long', capabilities }), condition)
        .rejects.toThrow('Recording layout is unavailable for this destination.');
      noWrites();
    }
  });
  it.each([
    ['sqlite', 'wide', 'plain', 'managed_wide_profile_v1'], ['postgres', 'wide', 'plain', 'managed_wide_profile_v1'],
    ['sqlite', 'long', 'plain', 'managed_long_v1'], ['postgres', 'long', 'timescale', 'managed_long_v1'],
    ['mysql', 'long', 'native_monthly', 'managed_long_v1'], ['sqlserver', 'long', 'native_monthly', 'managed_long_v1'],
  ] as const)('preserves supported %s/%s/%s request shape', async (kind, purpose, storageMode, layout) => {
    await ensureOrUpdateCircuitProfileGroups({ ...actionInput(kind), purpose, storageMode, capabilities: capabilitiesFor(kind) });
    expect(ensure).toHaveBeenCalledOnce();
    expect(ensure.mock.calls[0][0].request.group.layout_version).toBe(layout);
    if (purpose === 'wide') expect(ensure.mock.calls[0][0].request.group).not.toHaveProperty('storage_mode');
    else expect(ensure.mock.calls[0][0].request.group.storage_mode).toBe(storageMode);
    expect(start).not.toHaveBeenCalled(); expect(applySchema).not.toHaveBeenCalled();
  });
  it.each([
    ['postgres', 'wide', 'timescale'], ['sqlite', 'long', 'timescale'], ['sqlite', 'long', 'native_monthly'],
    ['postgres', 'long', 'native_monthly'], ['mysql', 'long', 'timescale'], ['unknown', 'long', 'plain'],
  ] as const)('rejects unsupported %s/%s/%s before building a mutation', async (kind, purpose, storageMode) => {
    await expect(ensureOrUpdateCircuitProfileGroups({ ...actionInput(kind), purpose, storageMode, capabilities: capabilitiesFor(kind) }))
      .rejects.toThrow('Recording layout is unavailable for this destination.');
    noWrites();
  });
  it('validates all saved profile layouts before any earlier compatible profile can mutate', async () => {
    const saved = longProfileGroup({ id: 'saved-cumulative', circuitId: 'cir-a', circuitName: 'Draft A', interval: 10, points: [0, 1], kind: 'cumulative' });
    saved.layout_version = 'managed_wide_profile_v1';
    const input = actionInput('mysql');
    await expect(ensureOrUpdateCircuitProfileGroups({ ...input, purpose: 'long', capabilities: capabilitiesFor('mysql'),
      savedLongGroups: [saved], enabledProfiles: [...input.enabledProfiles, { ...input.enabledProfiles[0], profile: 'cumulative', intervalValue: 10, intervalSeconds: 600 }],
    })).rejects.toThrow('Recording layout is unavailable for this destination.');
    noWrites();
  });
  it.each(['missing-kind', 'missing-id', 'missing-revision', 'missing-snapshot', 'nontransactional', 'native-unsupported'] as const)(
    'rejects the %s default/error path without inventing destination authority', async (condition) => {
      const input = { ...actionInput('mysql'), purpose: 'long' as const, capabilities: capabilitiesFor('mysql') };
      if (condition === 'missing-kind') input.connectorKind = undefined;
      if (condition === 'missing-id') input.destination.connector_id = '';
      if (condition === 'missing-revision') input.destination.connector_revision = '';
      if (condition === 'missing-snapshot') input.capabilities = {};
      if (condition === 'nontransactional') input.capabilities.mysql!.data!.transactional_engine = false;
      if (condition === 'native-unsupported') {
        input.storageMode = 'native_monthly';
        input.capabilities.mysql!.data!.supports_native_partition = false;
      }
      await expect(ensureOrUpdateCircuitProfileGroups(input)).rejects.toThrow('Recording layout is unavailable for this destination.');
      noWrites();
    },
  );
  it('allows batch inspection and close but not batch preparation while capability is unknown', async () => {
    groupsQuery.data.groups = ['a', 'b'].map((id) => longProfileGroup({
      id: `group-${id}`, circuitId: `cir-${id}`, circuitName: `Circuit ${id}`, interval: 15, points: [0, 1],
    }));
    vi.spyOn(studioV2DatalinkApi, 'get').mockRejectedValue(new Error('capability unavailable'));
    const view = await show();
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    expect(screen.getByTestId('circuit-profile-batch-preview-btn')).toBeDisabled();
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-close'));
    expect(screen.queryByTestId('circuit-profile-batch-editor')).not.toBeInTheDocument();
    noWrites();
  });
  it('does not expose an executable start retry while destination capability refreshes', async () => {
    const refresh = deferred();
    mockCapabilityRefresh(refresh);
    start.mockRejectedValueOnce(new TypeError('Failed to fetch'));
    const view = await show();
    await startGuided('circuit-profile');
    await waitFor(() => expect(start).toHaveBeenCalledOnce());
    const original = structuredClone(start.mock.calls[0][0]);
    // Lost ACK without O1 permits manual inspection, never a fabricated retry identity.
    await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('Inspect operation recovery'));
    expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
    act(() => { void view.client.invalidateQueries({ queryKey: ['studio-v2-workspace', 'mysql-recording-capability'] }); });
    await waitFor(() => expect(view.client.isFetching()).toBeGreaterThan(0));
    fireEvent.click(screen.getByTestId('recording-next-action'));
    expect(start).toHaveBeenCalledOnce();
    expect(ensure).toHaveBeenCalledOnce();
    expect(update).not.toHaveBeenCalled(); expect(applySchema).toHaveBeenCalledOnce(); expect(previewSchema).toHaveBeenCalledOnce();
    await act(async () => refresh.resolve({ data: { data: mysqlCapability() } }));
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
    expect(start).toHaveBeenCalledOnce(); expect(start.mock.calls[0][0]).toEqual(original);
  });
  it('blocks the whole selected batch when inactive circuit B has unsupported native storage', async () => {
    const a = longProfileGroup({ id: 'group-a', circuitId: 'cir-a', circuitName: 'A', interval: 15, points: [0, 1] });
    const b = longProfileGroup({ id: 'group-b', circuitId: 'cir-b', circuitName: 'B', interval: 15, points: [0, 1] });
    b.storage_mode = 'native_monthly';
    groupsQuery.data.groups = [a, b];
    const original = structuredClone(groupsQuery.data.groups);
    update.mockImplementation(async ({ id, request }: { id: string; request: WriteGroupUpdateRequest }) => ({
      workspace_revision: 'wrev-2', group: { ...request.group, id, revision: 'rev-next' },
    }));
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValue({ data: { data: { ...mysqlCapability(), supports_native_partition: false } } });
    const view = await show();
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    expect(screen.getByTestId('circuit-profile-storage-plain')).toBeChecked();
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    await act(async () => fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm')));
    noWrites();
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-batch-layout-blocked')).toHaveTextContent('selected storage mode is unavailable');
    expect(groupsQuery.data.groups).toEqual(original);
    fireEvent.click(screen.getByTestId('circuit-profile-batch-cancel'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-select-cir-b'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(update).toHaveBeenCalledOnce());
    expect(update.mock.calls[0][0].id).toBe('group-a');
    expect(update.mock.calls[0][0].request.expected_connector_revision).toBe('crev-1');
    expect(b.storage_mode).toBe('native_monthly');
    expect(ensure).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled(); expect(applySchema).not.toHaveBeenCalled();
  });
  it('rechecks all frozen batch targets against current capability after a preview', async () => {
    const a = longProfileGroup({ id: 'group-a', circuitId: 'cir-a', circuitName: 'A', interval: 15, points: [0, 1] });
    const b = longProfileGroup({ id: 'group-b', circuitId: 'cir-b', circuitName: 'B', interval: 15, points: [0, 1] });
    b.storage_mode = 'native_monthly';
    groupsQuery.data.groups = [a, b];
    const refresh = deferred();
    mockCapabilityRefresh(refresh);
    const view = await show();
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    const confirm = screen.getByTestId('circuit-profile-batch-confirm');
    expect(confirm).toBeEnabled();
    act(() => { void view.client.invalidateQueries({ queryKey: ['studio-v2-workspace', 'mysql-recording-capability'] }); });
    await waitFor(() => expect(confirm).toBeDisabled());
    await act(async () => refresh.resolve({ data: { data: { ...mysqlCapability(), supports_native_partition: false } } }));
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    expect(confirm).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-batch-layout-blocked')).toHaveTextContent('selected storage mode is unavailable');
    fireEvent.click(confirm);
    expect(b.storage_mode).toBe('native_monthly');
    noWrites();
  });
  it('lets the operator exclude the unsupported active circuit and prepare only a compatible circuit', async () => {
    const a = longProfileGroup({ id: 'group-a', circuitId: 'cir-a', circuitName: 'A', interval: 15, points: [0, 1] });
    const b = longProfileGroup({ id: 'group-b', circuitId: 'cir-b', circuitName: 'B', interval: 15, points: [0, 1] });
    a.storage_mode = 'native_monthly';
    groupsQuery.data.groups = [a, b];
    update.mockImplementation(async ({ id, request }: { id: string; request: WriteGroupUpdateRequest }) => ({
      workspace_revision: 'wrev-2', group: { ...request.group, id, revision: 'rev-next' },
    }));
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValue({ data: { data: { ...mysqlCapability(), supports_native_partition: false } } });
    const view = await show();
    await waitFor(() => expect(view.client.isFetching()).toBe(0));
    expect(screen.getByTestId('circuit-profile-storage-native-monthly')).toBeChecked();
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    const selectedA = screen.getByTestId('circuit-profile-batch-select-cir-a');
    expect(selectedA).toBeEnabled();
    fireEvent.click(selectedA);
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeEnabled();
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(update).toHaveBeenCalledOnce());
    expect(update.mock.calls[0][0].id).toBe('group-b');
    expect(a.storage_mode).toBe('native_monthly');
    expect(ensure).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled(); expect(applySchema).not.toHaveBeenCalled();
  });
});
