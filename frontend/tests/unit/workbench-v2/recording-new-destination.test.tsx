import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import type { ComponentProps } from 'react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor, type GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { RecordingDraftNavigation, useRecordingNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import type { GroupEditorPreparation } from '@/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation';
import { recordingCopyPrefill } from '@/features/datalink/workbench-v2/recording-editor/recordingCopyDraft';
import { dbTargetAPI } from '@/services/datalink';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { dbTargetKeys } from '@/hooks/datalink/keys';
import type { DatabaseConnector } from '@/types/datalink';
import type { WriteGroupCalculationCapabilities } from '@/types/studioV2WriteGroupCalculation';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), readiness: vi.fn(), delivery: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) } }));
vi.mock('@/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation', () => ({ GroupEditorPreparation: ({ save }: ComponentProps<typeof GroupEditorPreparation>) => <button data-testid="group-save" disabled={save.disabled} onClick={() => void save.invoke()}>Original Save</button> }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/RecordingTestWriteRecoveryPanel', () => ({ RecordingTestWriteRecoveryPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar', () => ({ GroupLifecycleBar: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel', () => ({ GroupTestWritePanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel', () => ({ GroupReadinessPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip', () => ({ GroupDeliveryStrip: () => null }));
const capabilities: WriteGroupCalculationCapabilities = {
  version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: true,
  reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'],
  cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true },
  schema_dialects: ['sqlite', 'postgres'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [],
};
const original = savedGroup({ status: 'running', record_id: '11111111-1111-4111-8111-111111111111', layout_version: 'managed_long_v1', storage_mode: 'plain',
  recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, row_policy: { incomplete_policy: 'partial', allowed_lateness_seconds: 0 },
  destination: { ...savedGroup().destination, storage_strategy: 'managed' },
  members: [{ ...savedGroup().members[0], calculation: 'latest', cadence: { interval: 15, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'float64' }],
});
function savedConnector(id: string, kind: 'sqlite' | 'postgres' = 'postgres'): DatabaseConnector {
  return { id, name: id === 'conn-1' ? 'Original S1' : `MAIN PG ${id}`, kind, enabled: true, status: 'ready', identity_revision: id === 'conn-1' ? 'crev-1' : `${id}-rev`,
    connection_config: kind === 'sqlite' ? { path: '/d.db' } : { database: `qa_${id.toLowerCase()}`, schema: 'public', host: '127.0.0.1', port: 55441, user: 'qa_gateway', password: 'DO_NOT_RENDER_SECRET' },
    last_check_error: '', default_write_interval_seconds: 5, created_at: '', updated_at: '' };
}
const pool = () => [savedConnector('conn-1', 'sqlite'), savedConnector('P1'), savedConnector('P2')];
const clients: QueryClient[] = [];
const dispatch = vi.fn(), leave = vi.fn();
function Back() { const navigation = useRecordingNavigation(); return <button onClick={() => navigation.request(leave)}>Back</button>; }
function mount(options: { host?: boolean; props?: Partial<GroupEditorProps> } = {}) {
  const state = savedGroupState(); state.db.connector = { ...state.db.connector, save_state: 'save-error', save_error: 'WRITE_GROUP_LEGACY_WRITE_CONFLICT' };
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  const props: GroupEditorProps = { group: undefined, presentation: 'simple', workspaceId: 'ws-1', workspaceRevision: 'db-setup-9', connector: state.db.connector,
    calculationCapabilities: capabilities, candidates: buildGroupCandidates(state).candidates, excluded: [], readonly: false,
    onSaved: vi.fn(), onClose: leave, onReload: vi.fn(), fetchLatest: vi.fn().mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'fresh-10', groups: [original] }),
    destinationContent: <button data-testid="original-setup">Original setup owner</button>, ...options.props };
  const tree = () => <QueryClientProvider client={client}><RecordingDraftNavigation><Back />{options.host
    ? <Step4Database presentation="simple" initialSelection="new" state={state} dispatch={dispatch} workspaceId="ws-1" />
    : <GroupEditor {...props} />}</RecordingDraftNavigation></QueryClientProvider>;
  const view = render(tree());
  return { ...view, client, state, props, update: (patch: Partial<GroupEditorProps>) => { Object.assign(props, patch); view.rerender(tree()); } };
}
async function fillDraft() {
  fireEvent.change(await screen.findByTestId('group-name'), { target: { value: 'MAIN PG Mean Counter' } });
  fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
  const member = await screen.findByTestId('member-calculation-pt-0');
  fireEvent.change(within(member).getByLabelText('calculation.interval'), { target: { value: '5' } });
}
async function choose(id: string) {
  const select = await screen.findByRole('combobox', { name: 'newDestination.label' });
  await waitFor(() => expect(within(select).getByRole('option', { name: `MAIN PG ${id}` })).toBeEnabled());
  fireEvent.change(select, { target: { value: id } });
  return select;
}
beforeEach(() => {
  vi.clearAllMocks(); sessionStorage.clear();
  vi.spyOn(dbTargetAPI, 'listConnectors').mockResolvedValue(pool());
  vi.spyOn(studioV2WorkspaceDatabaseAPI, 'updateConfig').mockRejectedValue(new Error('WRITE_GROUP_LEGACY_WRITE_CONFLICT'));
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'fresh-10', groups: [original], calculation_capabilities: capabilities });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'fresh-10', group_id: original.id, group_revision: original.revision, config_ready: false, schema_ready: false, ready: false, issues: [], applied_revision: '' });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockResolvedValue(undefined as never);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mockImplementation(async request => ({ workspace_revision: 'created-11', group: {
    ...original, ...request.group, id: 'NEW-B', revision: 'B1', status: 'ready', applied_revision: '',
    destination: { ...original.destination, ...request.group.destination, table_name: 'server_new_table' },
    members: request.group.members.map(member => ({ ...member, source_revision: 'fresh-source', mapping_revision: 'fresh-map' })),
  } }));
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); });

it('selects a fresh saved P1 for the real new host without rewriting the failed legacy default or running A', async () => {
  const view = mount({ host: true }), before = structuredClone(view.state.db.connector), groupBefore = structuredClone(original);
  await fillDraft(); await choose('P1');
  expect(screen.queryByTestId('group-use-current-destination')).not.toBeInTheDocument();
  expect(screen.queryByLabelText('step4.field_password')).not.toBeInTheDocument();
  expect(document.body).not.toHaveTextContent('DO_NOT_RENDER_SECRET');
  expect(within(screen.getByTestId('member-calculation-pt-0')).getByLabelText('calculation.interval')).toHaveValue('5');
  fireEvent.click(screen.getByText('Back')); fireEvent.click(await screen.findByRole('button', { name: 'navigation.keep' }));
  expect(screen.getByTestId('group-name')).toHaveValue('MAIN PG Mean Counter'); expect(leave).not.toHaveBeenCalled();
  const save = screen.getByTestId('group-save'); await waitFor(() => expect(save).toBeEnabled()); fireEvent.click(save); fireEvent.click(save);
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
  const request = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0];
  expect(request.expected_connector_revision).toBe('P1-rev');
  expect(request.group.destination).toMatchObject({ connector_id: 'P1', connector_revision: 'P1-rev', table_schema: 'public', table_name: '' });
  expect(request.group.record_id).not.toBe(original.record_id); expect(request.group).not.toHaveProperty('id');
  expect(request.group.members[0]).toMatchObject({ point_id: 'pt-0', calculation: 'latest', cadence: { interval: 5 } });
  expect(view.state.db.connector).toEqual(before); expect(original).toEqual(groupBefore);
  expect(dispatch.mock.calls.filter(([action]) => action.type === 'updateDbConnector')).toHaveLength(0);
  expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});
it('copies to saved P2 without retaining the source table or changing copied methods and record identity', async () => {
  const prefill = recordingCopyPrefill(original, capabilities, 'sqlite'); expect(prefill).toBeDefined();
  mount({ props: { copyPrefill: prefill } }); await choose('P2');
  const save = screen.getByTestId('group-save'); await waitFor(() => expect(save).toBeEnabled()); fireEvent.click(save);
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
  const wire = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0].group;
  expect(wire.destination).toMatchObject({ connector_id: 'P2', connector_revision: 'P2-rev', table_schema: 'public', table_name: '' });
  expect(wire.record_id).not.toBe(original.record_id); expect(wire.layout_version).toBe('managed_long_v1');
  expect(wire.members[0]).toMatchObject({ point_id: 'pt-0', calculation: 'latest', cadence: { interval: 15 } });
  await waitFor(() => expect(screen.queryByRole('combobox', { name: 'newDestination.label' })).not.toBeInTheDocument());
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
});
it('keeps local input but never substitutes the legacy pool while the canonical list is pending or fails', async () => {
  let reject!: (error: Error) => void;
  vi.mocked(dbTargetAPI.listConnectors).mockImplementation(() => new Promise((_resolve, fail) => { reject = fail; }));
  mount(); await fillDraft();
  expect(screen.getByTestId('new-destination-status')).toHaveTextContent('newDestination.loading');
  expect(screen.queryByTestId('original-setup')).not.toBeInTheDocument(); expect(screen.getByTestId('group-save')).toBeDisabled();
  await act(async () => { reject(new Error('DO_NOT_RENDER_SECRET')); });
  expect(await screen.findByTestId('new-destination-status')).toHaveTextContent('newDestination.error');
  expect(screen.getByTestId('group-name')).toHaveValue('MAIN PG Mean Counter'); expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(document.body).not.toHaveTextContent('DO_NOT_RENDER_SECRET'); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
});
it('preserves the original initial setup session only for a fresh empty pool and absent connector identity', async () => {
  vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([]);
  const connector = { ...savedGroupState().db.connector, persisted: false, save_state: 'idle' as const, connector_id: undefined, identity_revision: undefined };
  const view = mount({ props: { connector } });
  expect(await screen.findByTestId('original-setup')).toBeVisible();
  view.update({ connector: savedGroupState().db.connector });
  expect(screen.getByTestId('original-setup')).toBeVisible();
  expect(screen.queryByRole('combobox', { name: 'newDestination.label' })).not.toBeInTheDocument();
});
it.each(['disabled', 'revision', 'removed'] as const)('denies saving a chosen target when fresh canonical evidence is %s', async reason => {
  const view = mount(); await fillDraft(); await choose('P1'); await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled());
  const changed = pool().flatMap(connector => connector.id !== 'P1' ? [connector] : reason === 'removed' ? [] : [{ ...connector, enabled: reason !== 'disabled', identity_revision: reason === 'revision' ? 'new-revision' : connector.identity_revision }]);
  await act(async () => { view.client.setQueryData(dbTargetKeys.connectors(), changed); });
  fireEvent.click(screen.getByTestId('group-save')); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeDisabled()); expect(screen.getByTestId('new-destination-status')).toHaveTextContent('newDestination.stale');
  fireEvent.click(screen.getByTestId('group-save')); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-name')).toHaveValue('MAIN PG Mean Counter');
  if (reason === 'revision') {
    expect(screen.getByRole('combobox', { name: 'newDestination.label' })).toHaveValue('');
    await choose('P1'); await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled());
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
    expect(vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0].expected_connector_revision).toBe('new-revision');
  }
});
it('rechecks selected connector evidence after the original fresh workspace read before sending Save', async () => {
  let finish!: () => void;
  const fetchLatest = vi.fn(() => new Promise<{ workspace_revision: string; groups: typeof original[] }>(resolve => { finish = () => resolve({ workspace_revision: 'fresh-10', groups: [original] }); }));
  const view = mount({ props: { fetchLatest } }); await fillDraft(); await choose('P1'); await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(fetchLatest).toHaveBeenCalledOnce());
  await act(async () => { await view.client.invalidateQueries({ queryKey: dbTargetKeys.connectors(), refetchType: 'none' }); finish(); });
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-save')).toBeDisabled();
});
it('does not send a delayed Save after workspace A to B to A even with the same canonical connector', async () => {
  let finish!: () => void;
  const fetchLatest = vi.fn(() => new Promise<{ workspace_revision: string; groups: typeof original[] }>(resolve => { finish = () => resolve({ workspace_revision: 'fresh-10', groups: [original] }); }));
  const view = mount({ props: { fetchLatest } }); await fillDraft(); await choose('P1'); await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(fetchLatest).toHaveBeenCalledOnce());
  view.update({ workspaceId: 'OTHER' }); view.update({ workspaceId: 'ws-1' });
  await act(async () => { finish(); });
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
});
it.each(['saved', 'guided'] as const)('does not mount the new destination query for %s editing', async mode => {
  mount({ props: mode === 'saved' ? { group: original, connector: savedGroupState().db.connector } : { presentation: 'guided' } });
  expect(screen.queryByRole('combobox', { name: 'newDestination.label' })).not.toBeInTheDocument();
  expect(dbTargetAPI.listConnectors).not.toHaveBeenCalled();
  expect(screen.getByTestId('original-setup')).toBeInTheDocument();
});
it('retains saved MySQL and SQL Server choices while preserving their original capability and mixed-dialect denial', async () => {
  const mysql = { ...savedConnector('MYSQL'), kind: 'mysql' as const };
  const sqlserver = { ...savedConnector('SQLSERVER'), kind: 'sqlserver' as const };
  vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([...pool(), mysql, sqlserver]);
  const mysqlRead = vi.spyOn(studioV2WorkspaceDatabaseAPI, 'getMySQLRecordingCapability').mockRejectedValue(new Error('unavailable'));
  const sqlserverRead = vi.spyOn(studioV2WorkspaceDatabaseAPI, 'getSQLServerRecordingCapability').mockRejectedValue(new Error('unavailable'));
  mount(); await fillDraft(); await choose('MYSQL');
  await waitFor(() => expect(mysqlRead).toHaveBeenCalledWith('MYSQL', 'MYSQL-rev'));
  expect(screen.getByTestId('group-save')).toBeDisabled();
  await choose('SQLSERVER'); await waitFor(() => expect(sqlserverRead).toHaveBeenCalledWith('SQLSERVER', 'SQLSERVER-rev'));
  expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
});
it('does not present the legacy initial setup for a nonempty pool containing only disabled connectors', async () => {
  vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([{ ...savedConnector('P1'), enabled: false }]);
  const connector = { ...savedGroupState().db.connector, persisted: false, save_state: 'idle' as const, connector_id: undefined, identity_revision: undefined };
  mount({ props: { connector } });
  await waitFor(() => expect(screen.getByTestId('new-destination-status')).toHaveTextContent('newDestination.empty'));
  expect(screen.queryByTestId('original-setup')).not.toBeInTheDocument(); expect(screen.getByTestId('group-save')).toBeDisabled();
});
