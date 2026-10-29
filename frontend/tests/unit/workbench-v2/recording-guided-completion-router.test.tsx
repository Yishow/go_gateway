import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { WorkbenchV2Shell } from '@/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { parseRecordingStartOperation } from '@/utils/studioV2RecordingStartJson';
import { basicRecordingIntentKey, basicRecordingScopeKey, loadBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { guidedLifecycle, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';

// Real Shell → Step4 → Basic/Circuit owners, native state, schema panels, dialogs and React Query hooks.
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
const activate = vi.fn(), navigateTo = vi.fn();
let nextStatus: RecordingStartOperation['status'], nextActivated: boolean, serverOperation: RecordingStartOperation | undefined;
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
  state.current = 4; state.completed = new Set([1, 2, 3]); state.selectedRuleId = 'rule-2';
  return state;
}
beforeEach(() => {
  vi.clearAllMocks(); sessionStorage.clear(); localStorage.clear(); onlineManager.setOnline(true);
  Object.defineProperty(navigator, 'onLine', { configurable: true, value: true });
  groups = []; nextStatus = 'succeeded'; nextActivated = true; serverOperation = undefined;
  api.list.mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: 'W1', groups: structuredClone(groups) }));
  api.readiness.mockImplementation(async (id: string) => guidedReadiness(canonical(id), 'W1'));
  api.lifecycle.mockImplementation(async (id: string) => guidedLifecycle(canonical(id), 'W1'));
  api.runtime.mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  api.delivery.mockRejectedValue(new Error('Delivery evidence unavailable')); api.listCutovers.mockResolvedValue([]);
  api.measurements.mockResolvedValue([]);
  api.metadata.mockResolvedValue({ workspace_id: 'ws-1', connector_id: 'conn-1', connector_revision: 'crev-1', database: '', schema: '',
    table: 'readings', inspection_status: 'exists', columns: [] });
  api.recoveryList.mockImplementation(async (workspaceId: string) => ({ workspace_id: workspaceId, complete: true, operations: [] }));
  api.start.mockImplementation(async (request: RecordingStartRequest) => {
    expect(request.workspace_id).toBe('ws-1'); expect(request.device_ids).toEqual(['dev-1']);
    expect(request.expected_workspace_revision).toBe('W1');
    for (const intent of request.groups) expect(canonical(intent.group_id).revision).toBe(intent.expected_group_revision);
    serverOperation = result(request); return checkedOperation(serverOperation);
  });
  api.operation.mockImplementation(async (id: string) => {
    if (!serverOperation || serverOperation.operation_id !== id) throw new Error('Unverified operation'); return checkedOperation(serverOperation);
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); onlineManager.setOnline(true); });
function RouterHost({ state, holdNavigation = false }: { state: WorkbenchV2State; holdNavigation?: boolean }) {
  const actions = useWorkbenchV2State(state), location = useLocation(), navigate = useNavigate();
  return <><output data-testid="completion-location">{location.pathname}{location.search}</output>
    {location.pathname === '/studio/v2' ? <WorkbenchV2Shell state={actions.state} actions={actions} workspaceId="ws-1"
      navigateTo={target => { navigateTo(target); if (!holdNavigation) navigate(target); }} activateWorkspace={activate} />
      : <section data-testid="runtime-destination"><button type="button" onClick={() => navigate(-1)}>Back to setup</button></section>}
  </>;
}
function show(holdNavigation = false, state = twoDevices()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const tree = () => <QueryClientProvider client={client}><MemoryRouter initialEntries={['/studio/v2']}>
    <RouterHost state={state} holdNavigation={holdNavigation} />
  </MemoryRouter></QueryClientProvider>;
  const view = render(tree());
  return { ...view, client, rerenderHost: () => view.rerender(tree()) };
}
function noOtherMutation() {
  for (const command of [api.create, api.update, api.ensureCircuitProfile, api.ensureBasic, api.apply, api.disable, api.remove,
    api.schemaPreview, api.schemaApply, api.measurementCreate, api.measurementUpdate, api.testWritePreview, api.testWrite,
    api.reviewCutover, api.applyCutover, activate]) expect(command).not.toHaveBeenCalled();
}
async function editLegacy() { fireEvent.click(await screen.findByTestId('device-recording-edit-existing')); await screen.findByTestId('basic-recording-panel'); }

function profileGroup(): WriteGroup {
  const base = savedGroup();
  return savedGroup({ id: 'profile-A', layout_version: 'managed_long_v1', storage_mode: 'plain',
    members: [{ ...base.members[0], circuit_id: 'circuit-A', circuit_name: 'Water', measurement_id: '', target_column: '', required: false }],
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
    destination: { ...base.destination, storage_strategy: 'managed', schema_revision: base.revision, schema_digest: 'saved-schema' } });
}

function checkedOperation(value: RecordingStartOperation): RecordingStartOperation {
  const parsed = parseRecordingStartOperation(value);
  if (!parsed) throw new Error('Malformed server operation fixture');
  return parsed;
}
function result(request: RecordingStartRequest): RecordingStartOperation {
  return { operation_id: 'original-O1', action: 'recording_start', status: nextStatus,
    stage: nextStatus === 'succeeded' ? 'complete' : 'activation', intent_digest: 'a'.repeat(64), workspace_id: request.workspace_id,
    setup_revision: 'W1', device_ids: [...request.device_ids],
    groups: request.groups.map(group => ({ group_id: group.group_id, group_revision: group.expected_group_revision,
      saved: true, ready: true, applied: nextStatus !== 'partial',
      ...(nextStatus !== 'partial' ? { applied_revision: group.expected_group_revision } : {}) })),
    devices: request.device_ids.map(device_id => ({ device_id, activated: nextActivated })),
    created_at: '2026-10-28T12:00:00Z', updated_at: '2026-10-28T12:00:01Z' };
}
function ownerKey(mode: 'circuit' | 'legacy') {
  return basicRecordingScopeKey('ws-1', 'dev-1', mode === 'legacy' ? 'managed-recording' : encodeURIComponent(JSON.stringify(['circuit-A', 'conn-1', 'crev-1', ''])));
}
async function confirmStart(mode: 'circuit' | 'legacy') {
  if (mode === 'legacy') await editLegacy();
  const trigger = await screen.findByTestId(mode === 'circuit' ? 'circuit-profile-start' : 'basic-recording-start');
  await waitFor(() => expect(trigger).toBeEnabled()); fireEvent.click(trigger);
  expect(screen.getByRole('dialog')).toBeVisible(); expect(api.start).not.toHaveBeenCalled();
  const confirm = screen.getByTestId('recording-start-confirm'); fireEvent.click(confirm); fireEvent.click(confirm);
  await waitFor(() => expect(api.start).toHaveBeenCalledOnce());
  await waitFor(() => expect(loadBasicRecordingIntent(ownerKey(mode))?.operation_id).toBe('original-O1'));
  return structuredClone(api.start.mock.calls[0][0]) as RecordingStartRequest;
}
function stillSetup() { expect(screen.getByTestId('completion-location')).toHaveTextContent(/^\/studio\/v2$/); }

it.each(['succeeded', 'partial'] as const)('keeps the real Circuit %s result in setup until explicit Continue, routes only A and preserves O1 on Back', async status => {
  groups = [profileGroup()]; nextStatus = status; show(); const request = await confirmStart('circuit');
  expect(request.device_ids).toEqual(['dev-1']); expect(request.groups.map(group => group.group_id)).toEqual(['profile-A']);
  stillSetup(); expect(navigateTo).not.toHaveBeenCalled();
  const statusPanel = await screen.findByTestId('basic-recording-status'); expect(statusPanel).toHaveTextContent(`step4.basic.operation.${status}`);
  if (status === 'partial') { expect(statusPanel).toHaveAttribute('role', 'alert'); expect(statusPanel).toHaveTextContent('step4.basic.progress.not_applied'); }
  const proceed = await screen.findByTestId('basic-recording-continue-runtime'); fireEvent.click(proceed); fireEvent.click(proceed);
  expect(screen.getByTestId('completion-location')).toHaveTextContent('/studio/runtime?device_id=dev-1');
  expect(navigateTo).toHaveBeenCalledExactlyOnceWith('/studio/runtime?device_id=dev-1'); expect(api.start).toHaveBeenCalledOnce();
  fireEvent.click(screen.getByRole('button', { name: 'Back to setup' }));
  await screen.findByTestId('basic-recording-status'); stillSetup();
  expect(loadBasicRecordingIntent(ownerKey('circuit'))).toMatchObject({ operation_id: 'original-O1', request });
  expect(navigateTo).toHaveBeenCalledOnce(); expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it('offers explicit legacy partial continuation for confirmed A while retaining the incomplete recording status', async () => {
  groups = [legacyGroup()]; nextStatus = 'partial'; show(); await confirmStart('legacy'); stillSetup();
  const status = await screen.findByTestId('basic-recording-status'); expect(status).toHaveTextContent('step4.basic.operation.partial');
  expect(status).toHaveAttribute('role', 'alert'); expect(status).toHaveTextContent('step4.basic.progress.not_applied');
  expect(navigateTo).not.toHaveBeenCalled(); fireEvent.click(await screen.findByTestId('basic-recording-continue-runtime'));
  expect(screen.getByTestId('completion-location')).toHaveTextContent('/studio/runtime?device_id=dev-1');
  expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it.each(['succeeded', 'partial'] as const)('does not expose a second reviewed-start primary alongside the Circuit %s completion', async status => {
  groups = [profileGroup()]; nextStatus = status; show(true); await confirmStart('circuit');
  expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent(`step4.basic.operation.${status}`);
  expect(screen.getByTestId('basic-recording-continue-runtime')).toBeEnabled();
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument(); expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it.each(['unknown', 'pending'] as const)('does not hand off an unfinished %s response even if its device progress has an activated flag', async status => {
  groups = [profileGroup()]; nextStatus = status; show(true); await confirmStart('circuit');
  expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent(`step4.basic.operation.${status}`);
  expect(navigateTo).not.toHaveBeenCalled(); expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
  expect(screen.getByTestId('basic-recording-check')).toBeInTheDocument(); expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it('keeps a partial result with no confirmed activated device in setup and never falls back to selected-rule B', async () => {
  groups = [profileGroup()]; nextStatus = 'partial'; nextActivated = false; show(); await confirmStart('circuit');
  expect(await screen.findByTestId('basic-recording-status')).toHaveTextContent('step4.basic.progress.not_activated');
  stillSetup(); expect(navigateTo).not.toHaveBeenCalled(); expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
  expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it('retains live partial O1 query priority when only its persisted intent disappears during the same mounted owner', async () => {
  groups = [profileGroup()]; nextStatus = 'partial'; const view = show(true); await confirmStart('circuit');
  expect(serverOperation?.next_action).toBeUndefined();
  const panel = screen.getByTestId('circuit-profile-basic-panel');
  await waitFor(() => expect(screen.getByTestId('basic-recording-check')).toBeEnabled());
  const key = basicRecordingIntentKey(ownerKey('circuit')); expect(sessionStorage.getItem(key)).not.toBeNull();
  sessionStorage.removeItem(key); expect(loadBasicRecordingIntent(ownerKey('circuit'))).toBeUndefined();
  view.rerenderHost(); expect(screen.getByTestId('circuit-profile-basic-panel')).toBe(panel);
  const query = await screen.findByTestId('basic-recording-check'); await waitFor(() => expect(query).toBeEnabled());
  expect(query).toHaveAttribute('data-primary', 'true'); expect(query).toHaveTextContent('next.actions.query_operation');
  expect(document.querySelectorAll('button[data-primary="true"]')).toHaveLength(1);
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
  const reads = api.operation.mock.calls.length; fireEvent.click(query);
  await waitFor(() => expect(api.operation).toHaveBeenCalledTimes(reads + 1));
  expect(api.operation.mock.calls.every(([id]) => id === 'original-O1')).toBe(true);
  expect(screen.getByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.partial');
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
  stillSetup(); expect(navigateTo).not.toHaveBeenCalled(); expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it.each(['foreign workspace', 'foreign device'] as const)('never lends confirmed activation from a %s GET for retained O1 to the current owner', async foreign => {
  groups = [profileGroup()]; nextStatus = 'unknown'; nextActivated = false;
  api.operation.mockImplementation(async () => {
    if (!serverOperation) throw new Error('Operation has not been accepted');
    return checkedOperation({ ...serverOperation, status: 'succeeded', stage: 'complete', workspace_id: foreign === 'foreign workspace' ? 'foreign-W' : 'ws-1',
      device_ids: foreign === 'foreign device' ? ['dev-2'] : ['dev-1'],
      devices: [{ device_id: foreign === 'foreign device' ? 'dev-2' : 'dev-1', activated: true }], updated_at: '2026-10-28T12:00:02Z' });
  });
  const view = show(true); await confirmStart('circuit');
  await waitFor(() => expect(api.operation).toHaveBeenCalledWith('original-O1'));
  await act(async () => { await view.client.refetchQueries({ queryKey: [...studioV2WorkspaceKeys.all, 'recording-start', 'original-O1'] }); });
  stillSetup(); expect(navigateTo).not.toHaveBeenCalled(); expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
  const query = await screen.findByTestId('basic-recording-check'); await waitFor(() => expect(query).toBeEnabled());
  expect(query).toHaveAttribute('data-primary', 'true'); expect(query).toHaveTextContent('next.actions.query_operation');
  expect(document.querySelectorAll('button[data-primary="true"]')).toHaveLength(1);
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
  const reads = api.operation.mock.calls.length; fireEvent.click(query);
  await waitFor(() => expect(api.operation).toHaveBeenCalledTimes(reads + 1));
  expect(api.operation.mock.calls.every(([id]) => id === 'original-O1')).toBe(true);
  expect(loadBasicRecordingIntent(ownerKey('circuit'))?.operation_id).toBe('original-O1');
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
  expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument(); expect(navigateTo).not.toHaveBeenCalled();
  expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});
