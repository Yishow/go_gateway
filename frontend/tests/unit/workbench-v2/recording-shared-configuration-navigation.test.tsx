import type { ComponentProps } from 'react';
import type { CircuitRecordingPreparation } from '../../../src/features/datalink/workbench-v2/steps/step4/CircuitRecordingPreparation';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { basicRecordingScopeKey, loadBasicRecordingIntent, saveBasicRecordingIntent } from '../../../src/features/datalink/workbench-v2/state/basicRecording';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI } from '../../../src/services/studioV2RecordingStart';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { savedGroup, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
 ensureCircuitProfile: vi.fn(), list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
vi.mock('../../../src/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('../../../src/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), get: vi.fn() } }));
let preparation: ComponentProps<typeof CircuitRecordingPreparation>;
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/CircuitRecordingPreparation', async original => {
 const actual = await original<typeof import('../../../src/features/datalink/workbench-v2/steps/step4/CircuitRecordingPreparation')>();
 return { CircuitRecordingPreparation: (props: ComponentProps<typeof CircuitRecordingPreparation>) => { preparation = props; return <actual.CircuitRecordingPreparation {...props} />; } };
});
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const group = longProfileGroup({ id: 'profile-1', circuitId: 'circuit-1', circuitName: 'Saved circuit', interval: 15, points: [0, 1] });
function mount(groups: WriteGroup[], secondDevice = false, share = false) {
 let state = batchPanelState(); if (secondDevice) state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'Other offline PLC' }); const before = structuredClone(state); const dispatch = vi.fn();
 api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'w1', groups });
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 const tree = () => <QueryClientProvider client={client}><RecordingDraftNavigation><Step4Database state={state} dispatch={dispatch} workspaceId="ws-1" shareStatus={share ? { enabled: true, configured_enabled: true, running: false, port: 5020, address: '127.0.0.1', bind_state: 'pass', mapping_count: 0 } : undefined} /></RecordingDraftNavigation></QueryClientProvider>;
 const view = render(tree());
 return { ...view, state, before, dispatch, client, setState: (next: WorkbenchV2State) => { state = next; view.rerender(tree()); } };
}
beforeEach(() => { vi.restoreAllMocks(); vi.clearAllMocks(); sessionStorage.clear(); api.readiness.mockResolvedValue(null as never); api.delivery.mockResolvedValue(null as never); });
function expectNoMutation() { expect(api.ensureCircuitProfile).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled(); expect(api.update).not.toHaveBeenCalled(); expect(api.create).not.toHaveBeenCalled(); expect(api.apply).not.toHaveBeenCalled(); expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled(); }
it('Circuit raw invalid intent after failed storage is retained on Advanced Keep and explicitly discarded without converting the saved identity', async () => {
 const view = mount([group]); const interval = await screen.findByTestId('circuit-profile-interval-instantaneous');
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(interval, { target: { value: '1e' } });
 fireEvent.click(screen.getByTestId('circuit-status-open-advanced-instantaneous'));
 expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
 fireEvent.click(screen.getByText('navigation.keep')); expect(interval).toHaveValue('1e'); expect(screen.queryByTestId('group-editor')).not.toBeInTheDocument();
 fireEvent.click(screen.getByTestId('circuit-status-open-advanced-instantaneous')); fireEvent.click(await screen.findByText('navigation.discard'));
 expect(await screen.findByTestId('group-name')).toHaveValue(group.name); expectNoMutation(); expect(view.state).toEqual(view.before); expect(view.dispatch).not.toHaveBeenCalled();
});
it('legacy Basic interval intent after failed storage cannot be silently lost when returning to setup', async () => {
 const legacy = savedGroup({ basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1', destination: { ...savedGroup().destination, storage_strategy: 'managed' }, row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 } });
 mount([legacy]); fireEvent.click(await screen.findByTestId('device-recording-edit-existing'));
 await screen.findByTestId('basic-recording-interval');
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '37' } });
 fireEvent.click(screen.getByTestId('device-recording-back-setup')); expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
 fireEvent.click(screen.getByText('navigation.keep')); expect(screen.getByTestId('basic-recording-interval')).toHaveValue(37);
 fireEvent.click(screen.getByTestId('device-recording-back-setup')); fireEvent.click(await screen.findByText('navigation.discard'));
 fireEvent.click(await screen.findByTestId('device-recording-edit-existing')); expect(await screen.findByTestId('basic-recording-interval')).toHaveValue(60); expectNoMutation();
});
it('clean same-identity Basic to Advanced navigation preserves roles, mappings and layout without mutation', async () => {
 const view = mount([group]);
 fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-instantaneous'));
 expect(await screen.findByTestId('group-name')).toHaveValue(group.name);
 expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
 expect(screen.getByTestId('group-profile-kind')).toHaveTextContent('instantaneous');
 expect(view.state).toEqual(view.before); expect(view.dispatch).not.toHaveBeenCalled(); expectNoMutation();
 fireEvent.click(screen.getByTestId('group-editor-close'));
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('15'); expectNoMutation();
});
it('circuit selection keeps invalid lexical input until an explicit discard, and does not switch canonical roles', async () => {
 const other = longProfileGroup({ id: 'profile-2', circuitId: 'circuit-2', circuitName: 'Other circuit', interval: 30, points: [2, 3] });
 mount([group, other]); await screen.findByTestId('circuit-profile-interval-instantaneous');
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(await screen.findByTestId('circuit-profile-interval-instantaneous'), { target: { value: '015.' } });
 fireEvent.change(screen.getByTestId('circuit-profile-circuit-chooser'), { target: { value: 'circuit-2' } });
 expect(await screen.findByRole('alertdialog')).toBeInTheDocument(); fireEvent.click(screen.getByText('navigation.keep'));
 expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue('circuit-1'); expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('015.');
 fireEvent.change(screen.getByTestId('circuit-profile-circuit-chooser'), { target: { value: 'circuit-2' } }); fireEvent.click(await screen.findByText('navigation.discard'));
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30'); expectNoMutation();
});

it('offline raw Circuit draft proven by original scoped storage can navigate and reload without pretending it is saved configuration', async () => {
 const view = mount([group]); fireEvent.change(await screen.findByTestId('circuit-profile-interval-instantaneous'), { target: { value: '1e' } });
 fireEvent.click(screen.getByTestId('circuit-status-open-advanced-instantaneous'));
 expect(await screen.findByTestId('group-name')).toHaveValue(group.name); expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
 expectNoMutation(); view.unmount(); mount([group]);
 expect(await screen.findByTestId('circuit-profile-interval-instantaneous')).toHaveValue('1e'); expectNoMutation();
});

it('failed-storage Circuit intent blocks a device scope switch but a confirmed discard leaves device identities unchanged', async () => {
 const view = mount([group], true); await screen.findByTestId('circuit-profile-interval-instantaneous');
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '2e' } });
 const other = screen.getAllByTestId('basic-recording-device-tile').find(tile => tile.getAttribute('data-device-id') === 'dev-2')!;
 fireEvent.click(other); expect(await screen.findByRole('alertdialog')).toBeInTheDocument(); fireEvent.click(screen.getByText('navigation.keep'));
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('2e'); expect(other).toHaveAttribute('aria-pressed', 'false');
 fireEvent.click(other); fireEvent.click(await screen.findByText('navigation.discard'));
 expect(screen.getAllByTestId('basic-recording-device-tile').find(tile => tile.getAttribute('data-device-id') === 'dev-2')).toHaveAttribute('aria-pressed', 'true');
 expect(view.state).toEqual(view.before); expectNoMutation();
});
it('a failed background group read cannot destroy Circuit fields whose storage write failed', async () => {
 const view = mount([group]); await screen.findByTestId('circuit-profile-interval-instantaneous');
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '3e' } });
 api.list.mockRejectedValue(new Error('offline'));
 await act(async () => { await view.client.invalidateQueries({ queryKey: ['studio-v2-workspace'] }); });
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('3e');
 expect(screen.getByTestId('circuit-profile-groups-error')).toBeInTheDocument();
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toBeDisabled();
 api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'w2', groups: [{ ...group, revision: 'g2' }] });
 fireEvent.click(within(screen.getByTestId('circuit-profile-groups-error')).getByRole('button', { name: 'step4.group.retry' }));
 await waitFor(() => expect(screen.queryByTestId('circuit-profile-groups-error')).not.toBeInTheDocument());
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('3e'); expectNoMutation();
});

it('leaving an unknown existing Circuit request does not cancel, delete, rebase or replay its original operation intent', async () => {
 const scope = basicRecordingScopeKey('ws-1', 'dev-1', encodeURIComponent(JSON.stringify(['circuit-1', 'conn-1', 'crev-1', ''])));
 const original = { scope_key: scope, request: { request_id: 'unknown-original', workspace_id: 'ws-1', expected_workspace_revision: 'original-w', device_ids: ['dev-1'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: 'crev-1' }] } };
 saveBasicRecordingIntent(original); const before = loadBasicRecordingIntent(scope);
 mount([group], true); await screen.findByTestId('circuit-profile-interval-instantaneous');
 const other = screen.getAllByTestId('basic-recording-device-tile').find(tile => tile.getAttribute('data-device-id') === 'dev-2')!;
 fireEvent.click(other); expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.pending');
 fireEvent.click(screen.getByText('navigation.discard'));
 expect(loadBasicRecordingIntent(scope)).toEqual(before); expectNoMutation();
});
it('discard after a later storage failure keeps the earlier successfully retained Circuit draft', async () => {
 mount([group]); fireEvent.change(await screen.findByTestId('circuit-profile-interval-instantaneous'), { target: { value: '017' } });
 const stored = Object.fromEntries(Object.keys(sessionStorage).map(key => [key, sessionStorage.getItem(key)]));
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '2e' } });
 fireEvent.click(screen.getByTestId('circuit-status-open-advanced-instantaneous')); fireEvent.click(await screen.findByText('navigation.discard'));
 expect(await screen.findByTestId('group-name')).toHaveValue(group.name);
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('017');
 expect(Object.fromEntries(Object.keys(sessionStorage).map(key => [key, sessionStorage.getItem(key)]))).toEqual(stored); expectNoMutation();
});

it('ordinary destination authority changes keep the same Circuit draft owner and invalidate preparation without a per-keystroke leave prompt', async () => {
 const view = mount([group]); await screen.findByTestId('circuit-profile-interval-instantaneous');
 const previous = preparation; expect(previous.canPrepare).toBe(true);
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '4e' } });
 view.setState({ ...view.state, db: { ...view.state.db, connector: { ...view.state.db.connector, name: 'Editing connector', persisted: false } } });
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('4e');
 expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument(); expectNoMutation();
 expect(preparation.ownerKey).not.toBe(previous.ownerKey); expect(preparation.canPrepare).toBe(false); expect(preparation.schema.disabled).toBe(true);
 await act(async () => { await previous.flow.prepare(); await previous.flow.start(); }); expectNoMutation();
 view.setState({ ...view.state, db: { ...view.state.db, connector: { ...view.state.db.connector, identity_revision: 'changed-identity' } } });
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('4e');
 expect(preparation.ownerKey).not.toBe(previous.ownerKey); expect(preparation.canPrepare).toBe(false); expect(preparation.schema.disabled).toBe(true);
 await act(async () => { await previous.flow.prepare(); await previous.flow.start(); }); expectNoMutation();
});

it('fresh canonical legacy interval cannot overwrite same-owner failed-storage input', async () => {
 const legacy = savedGroup({ basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1', destination: { ...savedGroup().destination, storage_strategy: 'managed' }, row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 } });
 const view = mount([legacy]); fireEvent.click(await screen.findByTestId('device-recording-edit-existing'));
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '37' } });
 api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'w2', groups: [{ ...legacy, revision: 'g2', row_policy: { interval_seconds: 80, allowed_lateness_seconds: 0 } }] });
 await act(async () => { await view.client.invalidateQueries({ queryKey: ['studio-v2-workspace'] }); });
 expect(screen.getByTestId('basic-recording-interval')).toHaveValue(37); expectNoMutation();
});
it('configured Share does not replace an admitted failed-storage Circuit form when destination authority temporarily disappears', async () => {
 const view = mount([group], false, true); await screen.findByTestId('circuit-profile-interval-instantaneous');
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage full'); });
 fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '5e' } });
 view.setState({ ...view.state, db: { ...view.state.db, connector: { ...view.state.db.connector, persisted: false } } });
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('5e');
 expect(preparation.canPrepare).toBe(false); expectNoMutation();
});
it('one Advanced intent asks once for an unknown operation and leaves its receipt untouched', async () => {
 const scope = basicRecordingScopeKey('ws-1', 'dev-1', encodeURIComponent(JSON.stringify(['circuit-1', 'conn-1', 'crev-1', ''])));
 saveBasicRecordingIntent({ scope_key: scope, request: { request_id: 'unknown-advanced', workspace_id: 'ws-1', expected_workspace_revision: 'old-w', device_ids: ['dev-1'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: 'crev-1' }] } });
 const original = loadBasicRecordingIntent(scope); mount([group]);
 fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-instantaneous'));
 expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.pending'); fireEvent.click(screen.getByText('navigation.discard'));
 expect(await screen.findByTestId('group-name')).toHaveValue(group.name); expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
 expect(loadBasicRecordingIntent(scope)).toEqual(original); expectNoMutation();
});

it('canonical Save readback making an implicit gauge explicit does not block recoverable raw Circuit navigation', async () => {
 const initial = { ...group, members: group.members.map((member, index) => index === 0 ? { ...member, value_role: 'state' as const } : member) };
 const view = mount([initial]);
 fireEvent.change(await screen.findByTestId('circuit-profile-interval-instantaneous'), { target: { value: '1e' } });
 const canonical = { ...initial, revision: 'saved-role-revision', members: initial.members.map((member, index) => index === 1 ? { ...member, value_role: 'gauge' as const } : member) };
 api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'saved-role-workspace', groups: [canonical] });
 await act(async () => { await view.client.invalidateQueries({ queryKey: ['studio-v2-workspace'] }); });
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('1e');
 fireEvent.click(screen.getByTestId('circuit-status-open-advanced-instantaneous'));
 expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
 expect(await screen.findByTestId('group-name')).toHaveValue(group.name); expectNoMutation();
 view.unmount(); mount([canonical]);
 expect(await screen.findByTestId('circuit-profile-interval-instantaneous')).toHaveValue('1e'); expectNoMutation();
});

it('a genuinely different canonical role is not normalized away when retaining native Circuit intent', async () => {
 const view = mount([group]);
 fireEvent.change(await screen.findByTestId('circuit-profile-interval-instantaneous'), { target: { value: '1e' } });
 const changed = { ...group, revision: 'different-role', members: group.members.map(member => ({ ...member, value_role: 'state' as const })) };
 api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'different-role-workspace', groups: [changed] });
 await act(async () => { await view.client.invalidateQueries({ queryKey: ['studio-v2-workspace'] }); });
 fireEvent.click(screen.getByTestId('circuit-status-open-advanced-instantaneous'));
 expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
 fireEvent.click(screen.getByText('navigation.keep'));
 expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('1e'); expectNoMutation();
});
