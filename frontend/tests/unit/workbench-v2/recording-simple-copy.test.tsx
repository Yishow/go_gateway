import { dbTargetAPI } from '@/services/datalink';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createMemoryRouter, RouterProvider, useLocation, useNavigate } from 'react-router-dom';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import type { ComponentProps } from 'react';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { RecordingRouterNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { parseStudioPresentation } from '@/services/studioRecordingView';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { GroupEditorPreparation } from '@/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation';
import type { LocatorRecordingLifecycleHost } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/LocatorRecordingLifecycleHost';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
import type { StudioLocatorEntryRead, StudioLocatorEntryResult } from '@/types/studioV2LocatorEntry';
import type { WriteGroupCalculationCapabilities } from '@/types/studioV2WriteGroupCalculation';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { locatorRead } from './studio-locator-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { name?: string }) => key === 'copy.name' ? `${options?.name} copy` : key }) }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor', async original => ({ ...await original<typeof import('@/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor')>(), RecordingMixedResults: () => null }));
vi.mock('@/features/datalink/workbench-v2/recording-editor/GroupEditorPreparation', () => ({ GroupEditorPreparation: ({ save }: ComponentProps<typeof GroupEditorPreparation>) => <button data-testid="group-save" disabled={save.disabled} onClick={() => void save.invoke()}>Original Save</button> }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar', () => ({ GroupLifecycleBar: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupTestWritePanel', () => ({ GroupTestWritePanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel', () => ({ GroupReadinessPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupDeliveryStrip', () => ({ GroupDeliveryStrip: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/LocatorRecordingLifecycleHost', () => ({ LocatorRecordingLifecycleHost: ({ children }: ComponentProps<typeof LocatorRecordingLifecycleHost>) => children({ onAction: () => undefined }, false, () => true) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/SavedRecordingCard', () => ({ SavedRecordingCard: () => null }));
vi.mock('@/hooks/datalink/useStudioV2Workspace', () => ({ useStudioV2WorkspaceQuery: () => ({ data: source.workspace, isSuccess: true, isLoading: false, isError: false, refetch: vi.fn() }) }));
vi.mock('@/hooks/datalink/useModbusShareStatus', () => ({ useModbusShareStatusQuery: () => ({ data: undefined, refetch: vi.fn() }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceAuditHistory', () => ({ useStudioV2WorkspaceAuditHistoryQuery: () => ({ data: { entries: [] }, isError: false }) }));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceActivation', () => ({ useActivateStudioV2WorkspaceMutation: () => ({ mutateAsync: vi.fn() }) }));
vi.mock('@/hooks/datalink/useStudioLocatorEntry', () => ({ useStudioLocatorEntry: (locator: EntityLocator | null) => ({ data: { ok: true, value: locator?.group_id === 'COPY' ? copied : source }, isFetching: false, isError: false, refetch: vi.fn() }) }));
vi.mock('@/pages/datalink/workbench-v2/studioV2LocatorHydration', async original => ({ ...await original<typeof import('@/pages/datalink/workbench-v2/studioV2LocatorHydration')>(), studioV2LocatorHydration: () => savedGroupState() }));
vi.mock('@/pages/datalink/workbench-v2/useStudioV2AutosaveState', () => ({ useStudioV2AutosaveState: () => {
  const actions = useWorkbenchV2State(savedGroupState()); const query = { isLoading: false, isError: false, refetch: vi.fn() };
  return { state: actions.state, actions, workspaceHydrated: false, locatorHydration: { ready: true, blocked: false },
    devicesQuery: query, rulesQuery: query, mappingsQuery: query, databaseConfigQuery: query, databaseTargetsQuery: query };
} }));
const capabilities: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: true, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] };
const sourceGroup = savedGroup({ id: 'G', workspace_id: 'W', revision: 'G3', record_id: '11111111-1111-4111-8111-111111111111', name: 'Original recording',
 destination: { ...savedGroup().destination, storage_strategy: 'managed' }, recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, layout_version: 'managed_wide_profile_v1', storage_mode: 'plain',
 row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial' }, members: savedGroup().members.map(member => ({ ...member, calculation: 'latest', cadence: { interval: 15, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'float64' })) });
function receipt(group = sourceGroup, revision = 'W2'): StudioLocatorEntryRead {
 const value = locatorRead(4, group.id);
 return { ...value, locator: { ...value.locator, group_id: group.id, group_revision: group.revision }, entry: { ...value.entry, editor_hydrated: true, database_config: null, database_targets: [], workspace_revision: revision,
  group_heads: { [group.id]: group.revision }, write_groups: { workspace_id: 'W', workspace_revision: revision, groups: [group], calculation_capabilities: capabilities } },
  selection: { state: 'selected', target: { workspaceId: 'W', workspaceRevision: revision, step: 4, groupId: group.id, groupRevision: group.revision } } };
}
const source = receipt(); let copied = receipt({ ...sourceGroup, id: 'COPY', revision: 'C1', record_id: '22222222-2222-4222-8222-222222222222' }, 'W3');
const sourceURL = '/studio/v2?entry=locator&workspace_id=W&step=4&group_id=G&group_revision=G3&view=recording';
const clients: QueryClient[] = [];
function mount(url = `${sourceURL}&intent=copy#position`) {
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
 function Entry() { const location = useLocation(), navigate = useNavigate(), intent = parseStudioPresentation(location.search);
  return <DatalinkWorkbenchV2Page presentation={intent.presentation} locatorIntent={intent.locatorIntent} presentationLocation={{ search: location.search, hash: location.hash }} navigateTo={navigate} />;
 }
 const router = createMemoryRouter([{ path: '/prior', element: <div>Prior page</div> }, { path: '/studio/v2', element: <QueryClientProvider client={client}><RecordingRouterNavigation><Entry /></RecordingRouterNavigation></QueryClientProvider> }], { initialEntries: ['/prior', url], initialIndex: 1 });
 return { router, ...render(<RouterProvider router={router} />) };
}
beforeEach(() => {
 vi.clearAllMocks(); sessionStorage.clear();
 vi.spyOn(dbTargetAPI, 'listConnectors').mockResolvedValue([{ id: 'conn-1', identity_revision: 'crev-1', name: 'Line DB', kind: 'sqlite', enabled: true, status: 'ready', connection_config: { path: savedGroupState().db.connector.database }, last_check_error: '', default_write_interval_seconds: 5, created_at: '', updated_at: '' }]);
 vi.spyOn(window, 'scrollTo').mockImplementation(() => undefined);
 vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async locator => ({ ok: true, value: locator.group_id === 'COPY' ? copied : source }));
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(source.entry.write_groups!);
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.get).mockImplementation(async id => ({ group: id === 'COPY' ? copied.entry.write_groups!.groups[0] : sourceGroup, workspace_revision: id === 'COPY' ? 'W3' : 'W2' }));
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', group_id: 'G', group_revision: 'G3', applied_revision: '', config_ready: false, schema_ready: false, ready: false, issues: [] });
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockResolvedValue(undefined as never);
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mockImplementation(async request => {
  const group = { ...sourceGroup, ...request.group, destination: { ...sourceGroup.destination, ...request.group.destination, table_name: 'new_server_table' },
   members: request.group.members.map((member, i) => ({ ...member, source_revision: `fresh-source-${i}`, mapping_revision: `fresh-map-${i}` })), id: 'COPY', revision: 'C1', applied_revision: '' };
  copied = receipt(group, 'W3'); return { group, workspace_revision: 'W3' };
 });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
it('opens one fresh local copy through the real page, shell and editor without saving or starting', async () => {
 const view = mount(sourceURL); expect(await screen.findByTestId('group-name')).toHaveValue('Original recording');
 const copy = await screen.findByTestId('recording-copy'); fireEvent.click(copy); fireEvent.click(copy);
 await waitFor(() => expect(view.router.state.location.search).toContain('intent=copy'));
 await waitFor(() => expect(screen.getByTestId('group-name')).toHaveValue('Original recording copy'));
 expect(screen.getAllByTestId('group-name')).toHaveLength(1);
 expect(readStudioLocatorEditorEntry).toHaveBeenCalledWith(source.locator, { workspaceRevision: 'W2', groupRevision: 'G3' }, expect.any(AbortSignal));
 expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
 expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled();
});
it('saves a copy using only the original create owner then replaces the URL with the new canonical identity', async () => {
 const view = mount(); const save = await screen.findByTestId('group-save'); await waitFor(() => expect(save).toBeEnabled()); fireEvent.click(save); fireEvent.click(save);
 await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce());
 const request = vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mock.calls[0][0];
 expect(request.group.record_id).not.toBe(sourceGroup.record_id); expect(request.group).not.toHaveProperty('id'); expect(request.group).not.toHaveProperty('revision');
 expect(request.group.destination.table_name).toBe(''); expect(request.group.members[0]).not.toHaveProperty('source_revision');
 await waitFor(() => expect(view.router.state.location.search).toContain('group_id=COPY'));
 expect(view.router.state.location.search).not.toContain('intent=copy'); expect(view.router.state.location.hash).toBe('#position'); expect(view.router.state.historyAction).toBe('REPLACE');
 expect(studioV2WorkspaceWriteGroupsAPI.get).toHaveBeenCalledWith('COPY');
 expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled();
 const url = view.router.state.location.pathname + view.router.state.location.search; view.unmount(); mount(url);
 await waitFor(() => expect(screen.getByTestId('group-name')).toHaveValue('Original recording copy'));
 expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce();
});
it.each(['foreign-scope', 'unavailable'] as const)('does not turn a %s source read into an empty or editable original draft', async reason => {
 vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: false, reason }); mount();
 await screen.findByTestId('studio-locator-open-error'); expect(screen.queryByTestId('group-save')).not.toBeInTheDocument();
 expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('abandons a late copy source reply after Back without creating a draft or changing the new location', async () => {
 let finish!: (value: StudioLocatorEntryResult) => void;
 vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(resolve => { finish = resolve; }));
 const view = mount(); await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
 const signal = vi.mocked(readStudioLocatorEditorEntry).mock.calls[0][2]!;
 await act(async () => { await view.router.navigate(-1); }); expect(signal.aborted).toBe(true);
 await act(async () => finish({ ok: true, value: source }));
 expect(screen.getByText('Prior page')).toBeInTheDocument(); expect(screen.queryByTestId('group-save')).not.toBeInTheDocument(); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
});
it('leaves a guided canonical locator on the original surface without auto-opening a copy', async () => {
 mount(sourceURL.replace('&view=recording', '&view=advanced'));
 expect(await screen.findByTestId('workbench-v2-shell')).toBeInTheDocument(); expect(screen.queryByTestId('recording-copy')).not.toBeInTheDocument();
 expect(readStudioLocatorEditorEntry).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
});

it('retains the saved copy on foreign readback without returning to or modifying the source record', async () => {
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.get).mockImplementation(async id => ({ group: id === 'COPY' ? { ...copied.entry.write_groups!.groups[0], workspace_id: 'OTHER' } : sourceGroup, workspace_revision: id === 'COPY' ? 'W3' : 'W2' }));
 const view = mount(); const save = await screen.findByTestId('group-save'); await waitFor(() => expect(save).toBeEnabled()); fireEvent.click(save);
 await screen.findByTestId('studio-locator-open-error');
 expect(view.router.state.location.search).toContain('intent=copy'); expect(view.router.state.location.search).toContain('group_id=G');
 expect(screen.getByTestId('group-name')).toHaveValue('Original recording copy');
 expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce(); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('keeps the original dirty guard and does not overwrite an edited copy during rerender', async () => {
 const view = mount(sourceURL); const name = await screen.findByTestId('group-name');
 fireEvent.change(name, { target: { value: 'Unsaved original input' } }); fireEvent.click(await screen.findByTestId('recording-copy'));
 expect(await screen.findByRole('button', { name: 'navigation.keep' })).toBeInTheDocument();
 fireEvent.click(screen.getByRole('button', { name: 'navigation.keep' }));
 expect(view.router.state.location.search).not.toContain('intent=copy'); expect(name).toHaveValue('Unsaved original input');
 fireEvent.click(screen.getByTestId('recording-copy')); fireEvent.click(await screen.findByRole('button', { name: 'navigation.discard' }));
 await waitFor(() => expect(screen.getByTestId('group-name')).toHaveValue('Original recording copy'));
 const copyName = screen.getByTestId('group-name'); fireEvent.change(copyName, { target: { value: 'My unsaved copy' } });
 view.rerender(<RouterProvider router={view.router} />);
 expect(screen.getByTestId('group-name')).toBe(copyName); expect(copyName).toHaveValue('My unsaved copy');
 expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('never updates the source after a create ACK incorrectly returns the source identity', async () => {
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.create).mockResolvedValue({ group: sourceGroup, workspace_revision: 'W2' });
 const view = mount(); const save = await screen.findByTestId('group-save'); await waitFor(() => expect(save).toBeEnabled()); fireEvent.click(save);
 await screen.findByTestId('studio-locator-open-error');
 // Even a synthetic change after this invalid ACK cannot turn the copy into an update owner for the source.
 fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Must not overwrite source' } });
 await act(async () => { fireEvent.click(screen.getByTestId('group-save')); });
 expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
 expect(screen.getByTestId('group-name')).toBeDisabled(); expect(screen.getByTestId('group-save')).toBeDisabled();
 expect(view.router.state.location.search).toContain('intent=copy'); expect(studioV2WorkspaceWriteGroupsAPI.create).toHaveBeenCalledOnce();
});
