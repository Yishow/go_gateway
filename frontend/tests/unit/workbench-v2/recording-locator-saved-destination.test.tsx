import { createMemoryRouter, RouterProvider, useLocation, useNavigate } from 'react-router-dom';
import { useStudioLocatorEntry } from '@/hooks/datalink/useStudioLocatorEntry';
import { RecordingRouterNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation';
import { WorkbenchV2Shell } from '@/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { buildSimpleRecordingTarget, parseStudioPresentation } from '@/services/studioRecordingView';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { RecordingDraftNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { readStudioLocatorEditorEntry, readStudioLocatorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { dbTargetAPI, modbusShareAPI } from '@/services/datalink';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import type { SchemaPreviewToken } from '@/types/recordingPlan';
import type { WriteGroupCalculationCapabilities } from '@/types/studioV2WriteGroupCalculation';
import { locatorRead, nativeDevice, nativeMapping, nativeRule } from './studio-locator-fixtures';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn(), readStudioLocatorEntry: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceDatabase', async original => ({ ...await original<typeof import('@/services/studioV2WorkspaceDatabase')>(), studioV2WorkspaceDatabaseAPI: { getMetadata: vi.fn(), updateConfig: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) } }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: vi.fn().mockResolvedValue({ workspace_id: 'W', operations: [] }) }));
vi.mock('@/services/datalink', async original => ({ ...await original<typeof import('@/services/datalink')>(), dbTargetAPI: { listConnectors: vi.fn(), getConnector: vi.fn(), updateConnector: vi.fn() }, modbusShareAPI: { status: vi.fn() } }));

const capabilities: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: true, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] };
function receipt(): StudioLocatorEntryRead {
  const read = locatorRead(4); const group = read.entry.write_groups!.groups[0];
  Object.assign(group, { name: 'Saved successor', status: 'disabled', record_id: '11111111-1111-4111-8111-111111111111', layout_version: 'managed_long_v1', storage_mode: 'plain', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' },
    row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial', record_key_column: 'record_id', bucket_start_column: 'bucket_start', provenance_column: 'provenance', group_id_column: 'group_id', device_id_column: 'device_id' },
    destination: { connector_id: 'saved-target', connector_revision: 'saved-target-rev', database: 'successor.db', table_schema: 'main', table_name: 'saved_long_table', storage_strategy: 'managed' },
    members: [{ device_id: 'D', point_id: 'P-native', tag_id: 'T', source_revision: 'R2', mapping_revision: 'M2', target_column: 'value', required: true, calculation: 'latest', cadence: { interval: 15, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'int16', result_value_type: 'float64' }] });
  read.entry.write_groups!.calculation_capabilities = capabilities;
  Object.assign(read.entry, { editor_hydrated: true, devices: [nativeDevice()], source_rules: [nativeRule()], mappings: [nativeMapping()], database_targets: [],
    hydration_scope: { complete_workspace: false, device_ids: ['D'], rule_ids: ['R'], mapping_ids: ['M'], group_ids: ['G'] }, rule_heads: { R: 'R2' }, revalidated_mapping_ids: ['M'],
    database_config: { id: 'saved-target', identity_revision: 'saved-target-rev', setup_revision: 'W2', workspace_id: 'W', kind: 'sqlite', name: 'Saved target label', host: '', port: 0, database: 'successor.db', username: '', schema: 'main', table: 'saved_long_table', write_mode: 'insert', write_interval_seconds: 15, timestamp_column: 'bucket_start', status: 'ready', save_state: 'saved', created_at: '', updated_at: '' } });
  return read;
}
let current = receipt();
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI), clients: QueryClient[] = [];
function schemaPreview(): SchemaPreviewToken {
  return { ...managedSchemaPreviewToken, workspace_id: 'W', workspace_revision: 'W2', plan_id: 'G', plan_revision: 'G3', connector_id: 'saved-target', connector_revision: 'saved-target-rev', database: 'successor.db', schema: 'main', expires_at: '2099-01-01T00:00:00Z',
    statements: [], tables: [{ name: 'saved_long_table', action: 'unchanged', columns: ['value'] }], group_layout: { ...managedSchemaPreviewToken.group_layout!, table_name: 'saved_long_table' } };
}
function mount(locatorOnly = true) {
  const route = locatorRead(4), state = savedGroupState();
  state.db.connector = { ...state.db.connector, workspace_id: 'W', connector_id: 'workspace-default', identity_revision: 'default-rev', name: 'Workspace default label', database: 'default.db', password: 'PRIVATE_DEFAULT_CREDENTIAL' };
  const before = JSON.stringify(state), dispatch = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const content = (suspended = false) => <QueryClientProvider client={client}><RecordingDraftNavigation><Step4Database state={state} dispatch={dispatch} workspaceId="W"
    presentation="simple" initialSelection="located" locatorOnly={locatorOnly} locatorEntry={route} locatorSuspended={suspended} /></RecordingDraftNavigation></QueryClientProvider>;
  return { ...render(content()), content, state, before, dispatch };
}
function noWrites() {
  expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled(); expect(dbTargetAPI.updateConnector).not.toHaveBeenCalled();
  expect(api.create).not.toHaveBeenCalled(); expect(api.update).not.toHaveBeenCalled(); expect(api.apply).not.toHaveBeenCalled();
  expect(api.schemaApply).not.toHaveBeenCalled(); expect(api.testWrite).not.toHaveBeenCalled(); expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled();
}
beforeEach(() => {
  vi.clearAllMocks(); sessionStorage.clear(); current = receipt();
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async () => ({ ok: true, value: current }));
  api.readiness.mockImplementation(async () => ({ workspace_id: 'W', workspace_revision: 'W2', group_id: 'G', group_revision: 'G3', applied_revision: '', config_ready: true, schema_ready: false, ready: false, issues: [{ code: 'schema-proof-stale', severity: 'blocking', step: 'Step 4', scope: 'G', message: '' }] }));
  api.list.mockImplementation(async () => current.entry.write_groups!);
  api.delivery.mockRejectedValue(new Error('no delivery evidence')); api.schemaPreview.mockImplementation(async () => schemaPreview());
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async () => savedLifecycle(current.entry.write_groups!.groups[0], 'W2'));
  vi.mocked(studioV2WorkspaceDatabaseAPI.getMetadata).mockResolvedValue({ workspace_id: 'W', connector_id: 'saved-target', connector_revision: 'saved-target-rev', database: 'successor.db', schema: 'main', table: 'saved_long_table', inspection_status: 'exists', columns: [{ name: 'value', data_type: 'REAL', nullable: true, primary_key: false }] });
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'RUNTIME', settings_revision: 'SETTINGS', readiness_token: 'barrier', readiness: true, enabled: false, port: 502, address: '', bind_state: 'disabled', mapping_count: 0 });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

it('uses the saved target DTO for real editor metadata and explicit schema Check without changing the workspace default', async () => {
  const view = mount(); await screen.findByTestId('group-editor');
  await waitFor(() => expect(screen.getByTestId('recording-preparation-check')).toBeEnabled());
  expect(screen.queryByTestId('group-destination-mismatch')).not.toBeInTheDocument();
  expect(screen.getByTestId('recording-editor-summary')).toHaveTextContent('Saved target label');
  expect(screen.getByTestId('recording-editor-summary')).toHaveTextContent('saved_long_table');
  await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.getMetadata).toHaveBeenCalledWith('saved-target-rev', { groupId: 'G', groupRevision: 'G3' }, { workspaceId: 'W', workspaceRevision: 'W2' }));
  expect(api.schemaPreview).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('recording-preparation-check'));
  await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledExactlyOnceWith('G', { workspace_id: 'W', expected_workspace_revision: 'W2', expected_group_revision: 'G3', expected_connector_revision: 'saved-target-rev' }));
  expect(await screen.findByTestId('group-schema-preview-result')).toHaveTextContent('successor.db');
  expect(screen.getByTestId('group-schema-preview-result')).toHaveTextContent('saved_long_table');
  expect(screen.getByTestId('recording-preparation-confirm')).toBeEnabled();
  expect(JSON.stringify(view.state)).toBe(view.before); expect(view.dispatch).not.toHaveBeenCalled();
  expect(dbTargetAPI.getConnector).not.toHaveBeenCalled(); expect(dbTargetAPI.listConnectors).not.toHaveBeenCalled();
  expect(api.list).not.toHaveBeenCalled(); expect(document.body.textContent).not.toMatch(/PRIVATE_DEFAULT_CREDENTIAL|Workspace default label|default\.db/); noWrites();
});
it('opens the original saved Start review with the same target and only its selected source, without starting', async () => {
  api.readiness.mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', group_id: 'G', group_revision: 'G3', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] });
  const view = mount(); await screen.findByTestId('group-editor');
  await waitFor(() => expect(screen.getByTestId('group-apply')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-apply'));
  const review = await screen.findByTestId('saved-recording-start-review');
  expect(review).toHaveTextContent('Saved successor'); expect(review).toHaveTextContent('Device D');
  expect(screen.getByTestId('saved-recording-start-confirm')).toBeEnabled();
  expect(JSON.stringify(view.state)).toBe(view.before); expect(view.dispatch).not.toHaveBeenCalled(); noWrites();
});
it.each(['connector', 'revision'] as const)('does not borrow a wrong saved target %s for metadata, schema or Start', async kind => {
  const config = current.entry.database_config!;
  if (kind === 'connector') config.id = 'workspace-default'; else config.identity_revision = 'stale-revision';
  mount(); await screen.findByTestId('group-destination-mismatch');
  expect(screen.queryByTestId('recording-saved-destination')).not.toBeInTheDocument();
  expect(document.body.textContent).not.toMatch(/Workspace default label|default\.db|PRIVATE_DEFAULT_CREDENTIAL/);
  await waitFor(() => expect(screen.getByTestId('recording-preparation-check')).toBeDisabled());
  expect(screen.getByTestId('group-apply')).toBeDisabled();
  fireEvent.click(screen.getByTestId('recording-preparation-check'));
  expect(studioV2WorkspaceDatabaseAPI.getMetadata).not.toHaveBeenCalled(); expect(api.schemaPreview).not.toHaveBeenCalled();
  expect(screen.queryByTestId('saved-recording-start-review')).not.toBeInTheDocument(); noWrites();
});
it('does not adopt a late saved target editor response after locator suspension', async () => {
  let resolve!: (value: Awaited<ReturnType<typeof readStudioLocatorEditorEntry>>) => void;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(done => { resolve = done; }));
  const view = mount(); await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
  view.rerender(view.content(true)); await act(async () => { resolve({ ok: true, value: current }); });
  expect(screen.queryByTestId('group-editor')).not.toBeInTheDocument(); expect(view.dispatch).not.toHaveBeenCalled();
  expect(studioV2WorkspaceDatabaseAPI.getMetadata).not.toHaveBeenCalled(); expect(api.schemaPreview).not.toHaveBeenCalled(); noWrites();
});
it('does not publish late saved target schema consent after the locator lease is suspended', async () => {
  let resolve!: (value: SchemaPreviewToken) => void;
  api.schemaPreview.mockImplementationOnce(() => new Promise(done => { resolve = done; }));
  const view = mount(); await screen.findByTestId('group-editor');
  await waitFor(() => expect(screen.getByTestId('recording-preparation-check')).toBeEnabled());
  fireEvent.click(screen.getByTestId('recording-preparation-check')); await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledOnce());
  view.rerender(view.content(true)); await act(async () => { resolve(schemaPreview()); });
  expect(screen.queryByTestId('recording-preparation-confirm')).not.toBeInTheDocument();
  expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument(); expect(view.dispatch).not.toHaveBeenCalled(); noWrites();
});

it('shows only the saved P1 destination after a hydrated workspace enters its canonical editor', async () => {
  const group = current.entry.write_groups!.groups[0];
  group.destination = { ...group.destination, database: 'qa_p1', table_schema: 'public' };
  current.entry.write_groups!.calculation_capabilities = { ...capabilities, schema_dialects: ['sqlite', 'postgres'] };
  Object.assign(current.entry.database_config!, { kind: 'postgres', name: 'MAIN PG P1', database: 'qa_p1', schema: 'public', host: 'PRIVATE_TARGET_HOST', username: 'PRIVATE_TARGET_USER' });
  vi.mocked(studioV2WorkspaceDatabaseAPI.getMetadata).mockResolvedValue({ workspace_id: 'W', connector_id: 'saved-target', connector_revision: 'saved-target-rev', database: 'qa_p1', schema: 'public', table: 'saved_long_table', inspection_status: 'exists', columns: [] });
  api.schemaPreview.mockResolvedValue({ ...schemaPreview(), database: 'qa_p1', schema: 'public' });
  const view = mount(false); await screen.findByTestId('group-editor');
  await waitFor(() => expect(screen.getByTestId('recording-preparation-check')).toBeEnabled());
  expect(screen.queryByTestId('database-connection-card')).not.toBeInTheDocument();
  const destination = screen.getByTestId('recording-saved-destination');
  expect(destination).toHaveTextContent('MAIN PG P1'); expect(destination).toHaveTextContent('step4.kind_postgres'); expect(destination).toHaveTextContent('public');
  expect(within(destination).queryByRole('button')).not.toBeInTheDocument(); expect(within(destination).queryByRole('combobox')).not.toBeInTheDocument();
  expect(document.body.textContent).not.toMatch(/Workspace default label|default\.db|PRIVATE_DEFAULT_CREDENTIAL|PRIVATE_TARGET_HOST|PRIVATE_TARGET_USER/);
  fireEvent.click(screen.getByTestId('recording-preparation-check'));
  await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledExactlyOnceWith('G', { workspace_id: 'W', expected_workspace_revision: 'W2', expected_group_revision: 'G3', expected_connector_revision: 'saved-target-rev' }));
  expect(await screen.findByTestId('group-schema-preview-result')).toHaveTextContent('qa_p1');
  expect(JSON.stringify(view.state)).toBe(view.before); expect(view.dispatch).not.toHaveBeenCalled(); noWrites();
});

function mountLocalRecord(options: { foreignHead?: boolean; fullAdvanced?: boolean; custom?: boolean; unlockedManaged?: boolean } = {}) {
  const route = receipt();
  if (options.foreignHead) route.entry.group_heads = { G: 'FOREIGN' };
  if (options.unlockedManaged) {
    const prior = current.entry.write_groups!.groups[0];
    current.entry.write_groups!.groups[0] = savedGroup({ id: prior.id, workspace_id: prior.workspace_id, revision: prior.revision, name: prior.name, destination: prior.destination,
      members: [{ device_id: 'D', point_id: 'P-native', tag_id: 'T', source_revision: 'R2', mapping_revision: 'M2', target_column: 'value', required: true }] });
  }
  if (options.custom) current.entry.write_groups!.groups[0].destination.storage_strategy = 'custom';
  const navigate = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  function Caller() {
    const actions = useWorkbenchV2State();
    return <WorkbenchV2Shell state={{ ...savedGroupState(), view: 'flow', current: 4 }} actions={actions} workspaceId="W"
      deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [] }} locatorEntry={route} locatorOnly={!options.fullAdvanced} locatorCanLeave navigateTo={navigate}
      presentationLocation={{ search: '?view=advanced', hash: '#record-position' }} />;
  }
  render(<QueryClientProvider client={client}><RecordingDraftNavigation><Caller /></RecordingDraftNavigation></QueryClientProvider>);
  return navigate;
}
it('scoped guided record presents its confirmed name and returns through the original dirty guard to that record', async () => {
  const navigate = mountLocalRecord(); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Saved successor');
  expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('locator.recordTitle');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unsent local name' } });
  fireEvent.click(screen.getByRole('button', { name: 'locator.backRecord' }));
  expect(await screen.findByRole('alertdialog')).toBeInTheDocument(); expect(navigate).not.toHaveBeenCalled();
  fireEvent.click(screen.getByText('navigation.keep')); expect(screen.getByTestId('group-name')).toHaveValue('Unsent local name');
  fireEvent.click(screen.getByRole('button', { name: 'locator.backRecord' })); fireEvent.click(await screen.findByText('navigation.discard'));
  expect(navigate).toHaveBeenCalledExactlyOnceWith(buildSimpleRecordingTarget({ kind: 'edit', workspaceId: 'W', groupId: 'G', groupRevision: 'G3' }, '#record-position'));
  noWrites();
});
it('scoped guided record collapses technical versions and preserves exact managed identities without editable assignments', async () => {
  mountLocalRecord(); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  const details = screen.getByTestId('group-revision-details'); expect(details).not.toHaveAttribute('open');
  expect(within(details).getByTestId('locator-version-details')).toBeInTheDocument();
  expect(screen.getByTestId('group-managed-table')).toHaveTextContent('saved_long_table');
  expect(screen.queryByTestId('group-table')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-managed-column-P-native')).toHaveTextContent('value');
  expect(screen.queryByTestId('group-member-column-P-native')).not.toBeInTheDocument();
  expect(screen.getByRole('columnheader', { name: 'compact.actions' })).toBeInTheDocument();
  expect(screen.getByRole('button', { name: 'picker.removeMember' })).toHaveTextContent('compact.remove'); noWrites();
});
it('scoped guided custom record keeps the original table and member assignment controls', async () => {
  mountLocalRecord({ custom: true }); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  expect(screen.getByTestId('group-table')).toBeInTheDocument();
  expect(screen.getByTestId('group-member-column-P-native')).toBeInTheDocument();
  expect(screen.queryByTestId('group-managed-column-P-native')).not.toBeInTheDocument(); noWrites();
});
it('a mismatched group head does not receive a same-record title or return target', async () => {
  const navigate = mountLocalRecord({ foreignHead: true });
  expect(await screen.findByRole('heading', { level: 1 })).toHaveTextContent('locator.localTitle');
  expect(screen.queryByRole('button', { name: 'locator.backRecord' })).not.toBeInTheDocument(); expect(navigate).not.toHaveBeenCalled(); noWrites();
});

it('scoped guided unlocked managed record retains its existing editable table name', async () => {
  mountLocalRecord({ unlockedManaged: true }); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  expect(screen.getByTestId('group-table')).toBeEnabled();
  fireEvent.change(screen.getByTestId('group-table'), { target: { value: 'operator_named_table' } });
  expect(screen.getByTestId('group-table')).toHaveValue('operator_named_table'); noWrites();
});
it('full advanced retains its original uncollapsed revisions and assignment controls', async () => {
  mountLocalRecord({ fullAdvanced: true }); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  expect(screen.queryByTestId('group-revision-details')).not.toBeInTheDocument(); expect(screen.getByTestId('group-revisions')).toBeInTheDocument();
  expect(screen.getByTestId('group-member-column-P-native')).toBeInTheDocument(); expect(screen.getByTestId('group-table')).toBeInTheDocument(); noWrites();
});

function mountSavedHandoff(guided = false) {
  function RoutedRecord() {
    const location = useLocation(), navigate = useNavigate(), intent = parseStudioPresentation(location.search), actions = useWorkbenchV2State();
    const read = useStudioLocatorEntry(intent.locatorIntent.kind === 'locator' ? intent.locatorIntent.locator : null);
    const entry = read.data?.ok ? read.data.value : undefined;
    return <><output data-testid="record-url">{location.search}</output>{entry && <WorkbenchV2Shell
      state={{ ...savedGroupState(), view: 'flow', current: 4 }} actions={actions} workspaceId="W" locatorEntry={entry} locatorOnly locatorCanLeave
      recordingIntent={intent.presentation === 'recording' ? 'located' : undefined} navigateTo={navigate}
      presentationLocation={{ search: location.search, hash: location.hash }} />}</>;
  }
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  const router = createMemoryRouter([{ path: '/studio/v2', element: <RecordingRouterNavigation><RoutedRecord /></RecordingRouterNavigation> }],
    { initialEntries: [buildSimpleRecordingTarget({ kind: 'edit', workspaceId: 'W', groupId: 'G', groupRevision: 'G3' })!.replace('view=recording', guided ? 'view=advanced' : 'view=recording')] });
  render(<QueryClientProvider client={client}><RouterProvider router={router} /></QueryClientProvider>); return router;
}
function savedHandoffReads(foreignRevision = false) {
  vi.mocked(readStudioLocatorEntry).mockImplementation(async locator => ({ ok: true, value: { ...structuredClone(current), locator } }));
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async (_locator, expected) => expected.groupRevision === current.entry.group_heads.G && expected.workspaceRevision === current.entry.workspace_revision
    ? { ok: true, value: structuredClone(current) } : { ok: false, reason: 'stale-reference' });
  api.update.mockImplementation(async (id, request) => {
    expect(id).toBe('G'); expect(request.expected_group_revision).toBe('G3');
    current.entry.write_groups!.groups[0] = { ...current.entry.write_groups!.groups[0], name: request.group.name, revision: 'G4' };
    current.entry.workspace_revision = 'W3'; current.entry.write_groups!.workspace_revision = 'W3'; current.entry.group_heads.G = 'G4';
    if (current.entry.database_config) current.entry.database_config.setup_revision = 'W3';
    if (current.selection.state === 'selected') current = { ...current, selection: { ...current.selection, target: { ...current.selection.target, workspaceRevision: 'W3', groupRevision: 'G4' } } };
    const result = { workspace_revision: 'W3', group: structuredClone(current.entry.write_groups!.groups[0]) };
    if (foreignRevision) {
      current.entry.write_groups!.groups[0].revision = 'G5'; current.entry.group_heads.G = 'G5';
      if (current.selection.state === 'selected') current = { ...current, selection: { ...current.selection, target: { ...current.selection.target, groupRevision: 'G5' } } };
    }
    return result;
  });
}
it('replaces a located record own-save URL only after its canonical read, then opens field settings with the new exact revisions', async () => {
  savedHandoffReads(); const router = mountSavedHandoff(); await screen.findByTestId('group-editor');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own saved rename' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(new URLSearchParams(router.state.location.search).get('group_revision')).toBe('G4'));
  await screen.findByTestId('group-editor'); fireEvent.click(screen.getByRole('button', { name: 'simple.advancedPurpose' }));
  fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  expect(vi.mocked(readStudioLocatorEditorEntry).mock.calls.at(-1)?.[1]).toEqual({ workspaceRevision: 'W3', groupRevision: 'G4' });
  expect(screen.getByTestId('group-name')).toHaveValue('Own saved rename'); expect(api.update).toHaveBeenCalledTimes(1);
  expect(screen.queryByTestId('studio-locator-open-error')).not.toBeInTheDocument(); expect(api.create).not.toHaveBeenCalled();
  fireEvent.click(screen.getByRole('button', { name: 'locator.backRecord' }));
  await screen.findByTestId('simple-recording-shell'); await screen.findByTestId('group-editor');
  expect(new URLSearchParams(router.state.location.search).get('group_revision')).toBe('G4'); expect(api.update).toHaveBeenCalledTimes(1);
});
it('does not replace the located URL with a later foreign revision after an own save', async () => {
  savedHandoffReads(true); const router = mountSavedHandoff(); await screen.findByTestId('group-editor');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own saved rename' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(vi.mocked(readStudioLocatorEditorEntry).mock.calls.some(([,expected]) => expected.groupRevision === 'G5')).toBe(true));
  expect(new URLSearchParams(router.state.location.search).get('group_revision')).toBe('G3'); expect(api.create).not.toHaveBeenCalled();
});

it('does not replace a later Home navigation when the located own-save readback arrives after leaving', async () => {
  savedHandoffReads();
  let finish: (() => void) | undefined;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async (_locator, expected) => {
    if (expected.groupRevision === 'G4') await new Promise<void>(resolve => { finish = resolve; });
    return expected.groupRevision === current.entry.group_heads.G && expected.workspaceRevision === current.entry.workspace_revision
      ? { ok: true, value: structuredClone(current) } : { ok: false, reason: 'stale-reference' };
  });
  const router = mountSavedHandoff(); await screen.findByTestId('group-editor');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own saved rename' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(finish).toBeDefined());
  fireEvent.click(screen.getByRole('button', { name: 'simple.back' }));
  if (screen.queryByRole('alertdialog')) fireEvent.click(screen.getByText('navigation.discard'));
  await waitFor(() => expect(router.state.location.search).toBe('?view=recordings'));
  await act(async () => { finish!(); });
  expect(router.state.location.search).toBe('?view=recordings'); expect(api.update).toHaveBeenCalledTimes(1); expect(api.create).not.toHaveBeenCalled();
});

it.each(['rejected', 'unknown'] as const)('does not update a located URL after a %s Save outcome', async outcome => {
  savedHandoffReads(); api.update.mockRejectedValue(outcome === 'rejected' ? { response: { status: 422 } } : new Error('Connection interrupted'));
  const router = mountSavedHandoff(); await screen.findByTestId('group-editor');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unconfirmed rename' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await screen.findByTestId(outcome === 'unknown' ? 'group-save-unknown' : 'group-save-error');
  expect(new URLSearchParams(router.state.location.search).get('group_revision')).toBe('G3');
  expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(1); expect(api.update).toHaveBeenCalledTimes(1); expect(api.create).not.toHaveBeenCalled();
});

it('folds only confirmed scoped group list diagnostics and leaves Open failures visible', async () => {
  mountLocalRecord(); await screen.findByTestId('group-open-G');
  const details = screen.getByTestId('group-list-diagnostics'); expect(details).not.toHaveAttribute('open');
  expect(details).toHaveTextContent('G3'); expect(within(details).getByText('locatorEntry.historicalHint')).toBeInTheDocument();
  vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: false, reason: 'unavailable' });
  fireEvent.click(screen.getByTestId('group-open-G'));
  expect(await screen.findByTestId('studio-locator-open-error')).toBeVisible();
  expect(screen.queryByTestId('group-list-diagnostics')).not.toBeInTheDocument(); noWrites();
});

it('does not fold a dirty blocked guided group warning when another locator is requested', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const content = (route: StudioLocatorEntryRead) => <QueryClientProvider client={client}><RecordingDraftNavigation>
    <Step4Database state={savedGroupState()} dispatch={vi.fn()} workspaceId="W" presentation="guided" locatorOnly locatorEntry={route} />
  </RecordingDraftNavigation></QueryClientProvider>;
  const view = render(content(receipt())); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Retain this local draft' } });
  view.rerender(content(locatorRead(4, 'OTHER-G')));
  expect(await screen.findByText('locatorEntry.draftBlocked')).toBeVisible();
  expect(screen.queryByTestId('group-list-diagnostics')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-name')).toHaveValue('Retain this local draft'); noWrites();
});

it('returns a guided own Save through the canonical local list and back to the same simple record at the new version', async () => {
  savedHandoffReads(); const router = mountSavedHandoff(true);
  fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Guided own saved rename' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(new URLSearchParams(router.state.location.search).get('group_revision')).toBe('G4'));
  expect(new URLSearchParams(router.state.location.search).get('view')).toBe('advanced');
  await screen.findByTestId('group-open-G'); expect(screen.queryByTestId('group-editor')).not.toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'locator.backRecord' }));
  await screen.findByTestId('simple-recording-shell'); await screen.findByTestId('group-editor');
  expect(screen.getByTestId('group-name')).toHaveValue('Guided own saved rename');
  expect(new URLSearchParams(router.state.location.search).get('group_revision')).toBe('G4');
  expect(vi.mocked(readStudioLocatorEditorEntry).mock.calls.at(-1)?.[1]).toEqual({ workspaceRevision: 'W3', groupRevision: 'G4' });
  expect(api.update).toHaveBeenCalledTimes(1); expect(api.create).not.toHaveBeenCalled();
});
it.each(['rejected', 'unknown'] as const)('keeps the guided locator unchanged after a %s Save', async outcome => {
  savedHandoffReads(); api.update.mockRejectedValue(outcome === 'rejected' ? { response: { status: 422 } } : new Error('Connection interrupted'));
  const router = mountSavedHandoff(true); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-editor');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unconfirmed guided rename' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await screen.findByTestId(outcome === 'unknown' ? 'group-save-unknown' : 'group-save-error');
  expect(new URLSearchParams(router.state.location.search).get('group_revision')).toBe('G3');
  expect(new URLSearchParams(router.state.location.search).get('view')).toBe('advanced'); expect(api.update).toHaveBeenCalledTimes(1); expect(api.create).not.toHaveBeenCalled();
});
