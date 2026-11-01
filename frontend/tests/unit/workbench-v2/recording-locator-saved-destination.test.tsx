import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { RecordingDraftNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { dbTargetAPI, modbusShareAPI } from '@/services/datalink';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import type { SchemaPreviewToken } from '@/types/recordingPlan';
import type { WriteGroupCalculationCapabilities } from '@/types/studioV2WriteGroupCalculation';
import { locatorRead, nativeDevice, nativeMapping, nativeRule } from './studio-locator-fixtures';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
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
function mount() {
  const route = locatorRead(4), state = savedGroupState();
  state.db.connector = { ...state.db.connector, workspace_id: 'W', connector_id: 'workspace-default', identity_revision: 'default-rev', name: 'Workspace default label', database: 'default.db', password: 'PRIVATE_DEFAULT_CREDENTIAL' };
  const before = JSON.stringify(state), dispatch = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const content = (suspended = false) => <QueryClientProvider client={client}><RecordingDraftNavigation><Step4Database state={state} dispatch={dispatch} workspaceId="W"
    presentation="simple" initialSelection="located" locatorOnly locatorEntry={route} locatorSuspended={suspended} /></RecordingDraftNavigation></QueryClientProvider>;
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
