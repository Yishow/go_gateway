import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { basicRecordingScopeKey, loadBasicRecordingIntent, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { buildBasicStartRequest } from '@/features/datalink/workbench-v2/steps/step4/basicStartRequest';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { guidedLifecycle, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

// Real Step4, Basic/Circuit owners, schema panels, dialogs and React Query hooks.
// Only read/write service boundaries are simulated; no fixture grants a start permission.
const api = vi.hoisted(() => ({ list: vi.fn(), readiness: vi.fn(), lifecycle: vi.fn(), delivery: vi.fn(), listCutovers: vi.fn(),
  create: vi.fn(), update: vi.fn(), ensureCircuitProfile: vi.fn(), ensureBasic: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
  schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(), start: vi.fn(), operation: vi.fn(),
  runtime: vi.fn(), metadata: vi.fn(), measurements: vi.fn(), measurementCreate: vi.fn(), measurementUpdate: vi.fn(),
  testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), recoveryList: vi.fn(), reviewCutover: vi.fn(), applyCutover: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: api }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: api.lifecycle } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: api.start, operation: api.operation, ensureBasic: api.ensureBasic } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: api.runtime } }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: api.measurements, create: api.measurementCreate, update: api.measurementUpdate } }));
vi.mock('@/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: { getMetadata: api.metadata } }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: api.recoveryList }));

let groups: WriteGroup[];
const clients: QueryClient[] = [];
const dispatch = vi.fn(), activate = vi.fn();
function canonical(id: string) {
  const group = groups.find(value => value.id === id);
  if (!group) throw new Error('Missing canonical fixture group');
  return group;
}
function legacyGroup(deviceId = 'dev-1'): WriteGroup {
  const base = savedGroup();
  return savedGroup({ id: `saved-${deviceId}`, basic_managed_device_id: deviceId, layout_version: 'managed_wide_v1',
    members: [0, 1, 2].map(index => ({ device_id: deviceId, point_id: `${deviceId === 'dev-1' ? 'pt' : 'b-pt'}-${index}`,
      tag_id: `${deviceId === 'dev-1' ? 'tag' : 'b-tag'}-${index}`, source_revision: `s${index}`, mapping_revision: `m${index}`,
      target_column: ['temperature', 'pressure', 'batch'][index], required: true })),
    destination: { ...base.destination, storage_strategy: 'managed', schema_revision: base.revision, schema_digest: 'verified-saved-schema' } });
}
function twoDevices(): WorkbenchV2State {
  const state = savedGroupState();
  state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'PLC B' }, { ...state.devices[0], id: 'dev-empty', name: 'Empty device' });
  state.rules.push({ ...state.rules[0], id: 'rule-2', device_id: 'dev-2' });
  state.points.slice().forEach((point, index) => {
    const id = `rule-2-p-${index}`;
    state.points.push({ ...point, id, device_id: 'dev-2', rule_id: 'rule-2', name: `B ${point.name}` });
    state.mappings[id] = { ...state.mappings[point.id], point_id: id, device_id: 'dev-2', tag_key: `b.${index}`,
      display_name: `B source ${index}`, persisted_point_id: `b-pt-${index}`, tag_id: `b-tag-${index}` };
  });
  return state;
}
beforeEach(() => {
  vi.clearAllMocks(); sessionStorage.clear(); localStorage.clear(); onlineManager.setOnline(true);
  Object.defineProperty(navigator, 'onLine', { configurable: true, value: true });
  groups = [];
  api.list.mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: 'W1', groups: structuredClone(groups) }));
  api.readiness.mockImplementation(async (id: string) => guidedReadiness(canonical(id), 'W1'));
  api.lifecycle.mockImplementation(async (id: string) => guidedLifecycle(canonical(id), 'W1'));
  api.runtime.mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  api.delivery.mockRejectedValue(new Error('Delivery evidence unavailable')); api.listCutovers.mockResolvedValue([]);
  api.measurements.mockResolvedValue([]);
  api.metadata.mockResolvedValue({ workspace_id: 'ws-1', connector_id: 'conn-1', connector_revision: 'crev-1', database: '', schema: '',
    table: 'readings', inspection_status: 'exists', columns: [] });
  api.recoveryList.mockImplementation(async (workspaceId: string) => ({ workspace_id: workspaceId, complete: true, operations: [] }));
  api.operation.mockRejectedValue(new Error('Operation outcome remains unavailable'));
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); onlineManager.setOnline(true); });
function show(state = twoDevices()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const tree = () => <QueryClientProvider client={client}><Step4Database state={state} workspaceId="ws-1" dispatch={dispatch} activateWorkspace={activate} /></QueryClientProvider>;
  const view = render(tree());
  return { ...view, client, rerenderHost: () => view.rerender(tree()) };
}
function noMutation() {
  for (const command of [api.create, api.update, api.ensureCircuitProfile, api.ensureBasic, api.apply, api.disable, api.remove,
    api.schemaPreview, api.schemaApply, api.start, api.measurementCreate, api.measurementUpdate, api.testWritePreview, api.testWrite,
    api.reviewCutover, api.applyCutover, dispatch, activate]) expect(command).not.toHaveBeenCalled();
}
async function selectDevice(name: string) {
  const tile = await screen.findByRole('button', { name }); fireEvent.click(tile);
  await waitFor(() => expect(screen.getByRole('button', { name })).toHaveAttribute('aria-pressed', 'true'));
}
async function editLegacy() { fireEvent.click(await screen.findByTestId('device-recording-edit-existing')); await screen.findByTestId('basic-recording-panel'); }
async function openReview() {
  const start = await screen.findByTestId('basic-recording-start'); await waitFor(() => expect(start).toBeEnabled());
  start.focus(); fireEvent.click(start); expect(await screen.findByRole('dialog')).toBeVisible(); return start;
}

describe('mounted Step4 navigation preserves native drafts without commands', () => {
  it('retains a new Circuit 30-second draft and exact native source selection across A → B → A and reload', async () => {
    const view = show(); await screen.findByTestId('circuit-profile-interval-instantaneous');
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
    expandSources(); fireEvent.click(screen.getByLabelText('line.pressure')); fireEvent.click(screen.getByLabelText('line.batch'));
    await selectDevice('PLC B'); expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('15');
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '45' } });
    await selectDevice('PLC A'); expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30');
    expandSources(); expect(screen.getByLabelText('line.temperature')).toBeChecked(); expect(screen.getByLabelText('line.pressure')).not.toBeChecked();
    expect(screen.getByLabelText('line.batch')).not.toBeChecked(); noMutation();
    view.unmount(); show(); await waitFor(() => expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30'));
    expandSources(); expect(screen.getByLabelText('line.temperature')).toBeChecked(); expect(screen.getByLabelText('line.pressure')).not.toBeChecked(); noMutation();
  });
  it('cancels a genuinely ready review, then retains unsaved legacy 30 seconds and completeness through B → A', async () => {
    groups = [legacyGroup(), legacyGroup('dev-2')]; show(); await editLegacy(); await openReview();
    fireEvent.click(screen.getByTestId('recording-start-cancel')); expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
    fireEvent.change(screen.getByTestId('basic-recording-completeness'), { target: { value: 'partial' } });
    await selectDevice('PLC B'); await editLegacy(); expect(screen.getByTestId('basic-recording-interval')).toHaveValue(10);
    await selectDevice('PLC A'); await editLegacy(); expect(screen.getByTestId('basic-recording-interval')).toHaveValue(30);
    expect(screen.getByTestId('basic-recording-completeness')).toHaveValue('partial'); expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    expect(groups[0].row_policy.interval_seconds).toBe(10); noMutation();
  });
  it('invalidates an existing review on unsaved 30-second edit and never revives its old confirmation after A → B → A', async () => {
    groups = [legacyGroup(), legacyGroup('dev-2')]; show(); await editLegacy(); await openReview();
    const retiredConfirm = screen.getByTestId('recording-start-confirm');
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); await selectDevice('PLC B'); await selectDevice('PLC A'); await editLegacy();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(30); fireEvent.click(retiredConfirm);
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); noMutation();
  });
  it('keeps the unsaved legacy draft through actual Advanced disclosure, Back to setup and remount', async () => {
    groups = [legacyGroup()]; const view = show(); await editLegacy();
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
    const advanced = screen.getByTestId('step4-advanced-recording'); const summary = within(advanced).getByText('management.advancedTitle');
    fireEvent.click(summary); expect(advanced).toHaveAttribute('open'); fireEvent.click(summary);
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(30);
    fireEvent.click(screen.getByTestId('device-recording-back-setup')); expect(screen.getByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
    await editLegacy(); expect(screen.getByTestId('basic-recording-interval')).toHaveValue(30); noMutation();
    view.unmount(); show(); await editLegacy(); expect(screen.getByTestId('basic-recording-interval')).toHaveValue(30); noMutation();
  });
  it('retains an explicitly selected empty device through disclosure and remount instead of borrowing A sources', async () => {
    const view = show(); await selectDevice('Empty device');
    expect(screen.getByTestId('recording-scope-summary')).toHaveTextContent('Empty device');
    expandSources(); expect(screen.queryByLabelText('line.temperature')).not.toBeInTheDocument();
    fireEvent.click(screen.getByText('management.advancedTitle')); noMutation(); view.unmount(); show();
    await waitFor(() => expect(screen.getByRole('button', { name: 'Empty device' })).toHaveAttribute('aria-pressed', 'true'));
    expect(screen.getByRole('button', { name: 'PLC A' })).toHaveAttribute('aria-pressed', 'false');
    expect(screen.queryByTestId('circuit-profile-prepare')).not.toBeInTheDocument(); noMutation();
  });
  it('retains the exact unresolved original O1 across Advanced, another device and remount without resending', async () => {
    groups = [legacyGroup(), legacyGroup('dev-2')]; const owner = basicRecordingScopeKey('ws-1', 'dev-1', 'managed-recording');
    const request = buildBasicStartRequest({ workspaceId: 'ws-1', deviceId: 'dev-1', connectorRevision: 'crev-1' }, groups[0], 'W1', 'original-request');
    saveBasicRecordingIntent({ scope_key: owner, operation_id: 'O1', request }); const view = show(); await editLegacy();
    await waitFor(() => expect(api.operation).toHaveBeenCalledWith('O1')); fireEvent.click(screen.getByText('management.advancedTitle'));
    await selectDevice('PLC B'); await selectDevice('PLC A'); await editLegacy();
    expect(loadBasicRecordingIntent(owner)).toMatchObject({ operation_id: 'O1', request }); noMutation();
    view.unmount(); show(); await editLegacy(); expect(loadBasicRecordingIntent(owner)).toMatchObject({ operation_id: 'O1', request });
    expect(screen.getByTestId('basic-recording-check')).toBeInTheDocument(); noMutation();
  });
});

describe('mounted navigation accessibility contracts (not browser geometry)', () => {
  it.each(['Cancel', 'Escape'] as const)('restores the actual review trigger after %s without changing saved settings', async dismissal => {
    groups = [legacyGroup()]; show(); await editLegacy(); await openReview();
    if (dismissal === 'Cancel') fireEvent.click(screen.getByTestId('recording-start-cancel'));
    else fireEvent.keyDown(screen.getByRole('dialog'), { key: 'Escape', code: 'Escape' });
    await waitFor(() => expect(screen.getByTestId('basic-recording-start')).toHaveFocus()); noMutation();
  });
  it.each(['native draft', 'device scope'] as const)('does not return focus to a retired opener when %s invalidates review', async change => {
    groups = [legacyGroup(), legacyGroup('dev-2')]; const state = twoDevices(); const view = show(state); await editLegacy();
    const trigger = await openReview(), cancel = screen.getByTestId('recording-start-cancel'); const focus = vi.spyOn(trigger, 'focus');
    if (change === 'native draft') fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
    else { state.devices = state.devices.filter(device => device.id !== 'dev-1'); view.rerenderHost(); }
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    const current = await screen.findByRole('button', { name: change === 'native draft' ? 'PLC A' : 'PLC B' });
    current.focus(); fireEvent.click(cancel); expect(current).toHaveFocus(); expect(focus).not.toHaveBeenCalled(); noMutation();
  });
  it('Escape closes complete rule details and restores their real trigger without another evidence query', async () => {
    groups = [legacyGroup()]; show(); await editLegacy(); const trigger = await screen.findByRole('button', { name: 'details' });
    await waitFor(() => expect(api.runtime).toHaveBeenCalledOnce()); trigger.focus(); fireEvent.click(trigger);
    const dialog = await screen.findByRole('dialog'); expect(dialog).toHaveAccessibleName('detailsTitle');
    expect(dialog).toHaveClass('overflow-y-auto'); fireEvent.keyDown(dialog, { key: 'Escape', code: 'Escape' });
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument()); expect(trigger).toHaveFocus();
    expect(api.runtime).toHaveBeenCalledOnce(); noMutation();
  });
  it('uses native focusable source controls, labelled disclosure and bounded list scrolling without replacing selection', async () => {
    show(); await screen.findByTestId('circuit-profile-interval-instantaneous');
    const toggle = screen.getByTestId('circuit-profile-sources-toggle'); toggle.focus(); expect(toggle).toHaveFocus();
    expect(toggle.tagName).toBe('BUTTON'); expect(toggle).toHaveAttribute('aria-expanded', 'false');
    fireEvent.click(toggle); expect(toggle).toHaveAttribute('aria-expanded', 'true');
    expect(document.getElementById(toggle.getAttribute('aria-controls') ?? '')).toBeVisible();
    expect(screen.getByTestId('circuit-profile-member-list')).toHaveClass('max-h-72', 'overflow-y-auto');
    const source = screen.getByRole('checkbox', { name: 'line.pressure' }); source.focus(); expect(source).toHaveFocus();
    fireEvent.click(source); expect(source).not.toBeChecked(); fireEvent.click(toggle); fireEvent.click(toggle);
    expect(screen.getByRole('checkbox', { name: 'line.pressure' })).not.toBeChecked(); noMutation();
  });
  it('keeps cumulative source repair local, opens the disclosure and focuses the real toggle', async () => {
    show(); await screen.findByTestId('circuit-profile-interval-instantaneous'); fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    const toggle = screen.getByTestId('circuit-profile-sources-toggle'); expect(toggle).toHaveAttribute('aria-expanded', 'true');
    fireEvent.click(toggle); const repair = screen.getByTestId('circuit-profile-cumulative-choose-sources'); repair.focus();
    fireEvent.click(repair); await waitFor(() => expect(toggle).toHaveFocus()); expect(toggle).toHaveAttribute('aria-expanded', 'true');
    const error = screen.getByTestId('circuit-profile-cumulative-sources-error'); expect(error).toHaveAttribute('role', 'alert');
    expect(screen.getByRole('checkbox', { name: 'line.temperature step4.basic.circuit.cumulative_label' })).not.toBeChecked(); noMutation();
  });
});
