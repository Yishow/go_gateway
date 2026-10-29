import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { loadBasicRecordingIntent, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { guidedLifecycle, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';

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
function noOtherMutation() {
  for (const command of [api.create, api.update, api.ensureCircuitProfile, api.ensureBasic, api.apply, api.disable, api.remove,
    api.schemaPreview, api.schemaApply, api.measurementCreate, api.measurementUpdate, api.testWritePreview, api.testWrite,
    api.reviewCutover, api.applyCutover, dispatch, activate]) expect(command).not.toHaveBeenCalled();
}
async function selectDevice(name: string) {
  const tile = await screen.findByRole('button', { name }); fireEvent.click(tile);
  await waitFor(() => expect(screen.getByRole('button', { name })).toHaveAttribute('aria-pressed', 'true'));
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
function deferred<T>() {
  let resolve!: (value: T) => void;
  let reject!: (error: unknown) => void;
  const promise = new Promise<T>((yes, no) => { resolve = yes; reject = no; });
  return { promise, resolve, reject };
}
function retained(request: RecordingStartRequest) {
  const entries = Object.values(sessionStorage).map(value => { try { return JSON.parse(value); } catch { return null; } });
  return entries.find(value => value?.request?.request_id === request.request_id);
}
function result(request: RecordingStartRequest): RecordingStartOperation {
  return { operation_id: 'original-O1', action: 'recording_start', status: 'unknown', stage: 'activation',
    intent_digest: 'a'.repeat(64), workspace_id: request.workspace_id, setup_revision: 'W1', device_ids: request.device_ids,
    groups: request.groups.map(group => ({ group_id: group.group_id, group_revision: group.expected_group_revision, saved: true, ready: true, applied: true })),
    devices: request.device_ids.map(device_id => ({ device_id, activated: false })), created_at: '2026-10-28T12:00:00Z', updated_at: '2026-10-28T12:00:01Z' };
}
async function confirmStart(mode: 'circuit' | 'legacy') {
  if (mode === 'legacy') await editLegacy();
  const trigger = await screen.findByTestId(mode === 'circuit' ? 'circuit-profile-start' : 'basic-recording-start');
  await waitFor(() => expect(trigger).toBeEnabled()); fireEvent.click(trigger);
  expect(screen.getByRole('dialog')).toBeVisible(); expect(api.start).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('recording-start-confirm'));
  await waitFor(() => expect(api.start).toHaveBeenCalledOnce());
  return structuredClone(api.start.mock.calls[0][0]) as RecordingStartRequest;
}

it.each(['post', 'invalidation'] as const)('adopts the original Circuit O1 after A → B → A remount while %s is pending, without resending', async pending => {
  groups = [profileGroup()];
  const response = deferred<RecordingStartOperation>(); api.start.mockReturnValue(response.promise);
  const view = show(); const request = await confirmStart('circuit');
  const record = retained(request); expect(record.operation_id).toBeUndefined();
  expect(request.groups).toEqual([{ group_id: 'profile-A', expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1' }]);
  const read = deferred<{ workspace_id: string; workspace_revision: string; groups: WriteGroup[] }>();
  if (pending === 'invalidation') {
    api.list.mockReturnValue(read.promise);
    await act(async () => response.resolve(result(request)));
    await waitFor(() => expect(view.client.getQueryState(studioV2WorkspaceKeys.writeGroups())?.fetchStatus).toBe('fetching'));
    expect(retained(request).operation_id).toBeUndefined();
  }
  await selectDevice('PLC B'); expect(screen.queryByTestId('basic-recording-status')).not.toBeInTheDocument();
  await selectDevice('PLC A');
  expect(screen.getByTestId('basic-recording-status')).toHaveTextContent('step4.basic.operation.unknown');
  api.operation.mockResolvedValue(result(request));
  if (pending === 'post') await act(async () => response.resolve(result(request)));
  else await act(async () => read.resolve({ workspace_id: 'ws-1', workspace_revision: 'W1', groups: structuredClone(groups) }));
  await waitFor(() => expect(retained(request).operation_id).toBe('original-O1'));
  expect(loadBasicRecordingIntent(record.scope_key)).toMatchObject({ operation_id: 'original-O1', request });
  await waitFor(() => expect(api.operation).toHaveBeenCalledWith('original-O1'));
  const check = await screen.findByTestId('basic-recording-check'); await waitFor(() => expect(check).toBeEnabled());
  expect(api.operation).toHaveBeenCalledOnce(); fireEvent.click(check);
  await waitFor(() => expect(api.operation).toHaveBeenCalledTimes(2));
  expect(api.operation.mock.calls.every(([id]) => id === 'original-O1')).toBe(true);
  expect(api.start).toHaveBeenCalledTimes(1); expect(api.start.mock.calls[0][0]).toEqual(request); noOtherMutation();
});

it('recovers legacy O1 after a full host reload before the original POST reply, with exactly the retained original request', async () => {
  groups = [legacyGroup()]; const response = deferred<RecordingStartOperation>(); api.start.mockReturnValue(response.promise);
  const first = show(); const request = await confirmStart('legacy'); const record = retained(request);
  first.unmount(); show(); await editLegacy();
  expect(await screen.findByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation');
  api.operation.mockResolvedValue(result(request)); await act(async () => response.resolve(result(request)));
  await waitFor(() => expect(retained(request).operation_id).toBe('original-O1'));
  await waitFor(() => expect(api.operation).toHaveBeenCalledWith('original-O1'));
  expect(loadBasicRecordingIntent(record.scope_key)).toMatchObject({ operation_id: 'original-O1', request });
  expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it('never overwrites a newer retained O2 when an old unmounted A POST later completes', async () => {
  groups = [profileGroup()]; const response = deferred<RecordingStartOperation>(); api.start.mockReturnValue(response.promise);
  show(); const request = await confirmStart('circuit'); const original = retained(request);
  await selectDevice('PLC B');
  const replacement = { scope_key: original.scope_key, operation_id: 'newer-O2', request: { ...request, request_id: 'newer-request-O2' } };
  saveBasicRecordingIntent(replacement);
  await act(async () => response.resolve(result(request)));
  await waitFor(() => expect(api.list.mock.calls.length).toBeGreaterThan(1));
  expect(loadBasicRecordingIntent(original.scope_key)).toMatchObject(replacement);
  expect(screen.queryByTestId('basic-recording-status')).not.toBeInTheDocument();
  expect(api.operation).not.toHaveBeenCalled(); expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});

it('never erases newer O2 when an old unmounted Circuit request later reports a definitive rejection', async () => {
  groups = [profileGroup()]; const response = deferred<RecordingStartOperation>(); api.start.mockReturnValue(response.promise);
  show(); const request = await confirmStart('circuit'); const original = retained(request);
  await selectDevice('PLC B');
  const replacement = { scope_key: original.scope_key, operation_id: 'newer-O2', request: { ...request, request_id: 'newer-request-O2' } };
  saveBasicRecordingIntent(replacement);
  await act(async () => response.reject({ code: 'revision-conflict', status: 409 }));
  expect(loadBasicRecordingIntent(original.scope_key)).toMatchObject(replacement);
  expect(api.start).toHaveBeenCalledOnce(); expect(api.operation).not.toHaveBeenCalled(); noOtherMutation();
});

it.each(['circuit', 'legacy'] as const)('keeps lost-ACK %s manual after navigation and reload, never inventing O1 or resending', async mode => {
  groups = [mode === 'circuit' ? profileGroup() : legacyGroup()];
  const response = deferred<RecordingStartOperation>(); api.start.mockReturnValue(response.promise);
  const first = show(); const request = await confirmStart(mode); const original = retained(request);
  await act(async () => response.reject(new TypeError('Failed to fetch')));
  first.unmount(); show(); if (mode === 'legacy') await editLegacy();
  const manual = await screen.findByTestId('recording-next-action');
  await waitFor(() => expect(manual).toHaveTextContent('next.actions.recover_operation'));
  fireEvent.click(manual); expect(api.operation).not.toHaveBeenCalled();
  expect(loadBasicRecordingIntent(original.scope_key)).toEqual(original);
  expect(screen.queryByTestId('recording-start-confirm')).not.toBeInTheDocument();
  expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});
it('leaves a remounted no-ID owner manual when durable completion fails, with no auto query or start', async () => {
  groups = [profileGroup()]; const response = deferred<RecordingStartOperation>(); api.start.mockReturnValue(response.promise);
  show(); const request = await confirmStart('circuit'); const original = retained(request);
  await selectDevice('PLC B'); await selectDevice('PLC A');
  const write = Storage.prototype.setItem;
  vi.spyOn(Storage.prototype, 'setItem').mockImplementation(function (this: Storage, key, value) {
    if (key.startsWith('wbv2.basic-recording.')) throw new Error('Storage write unavailable');
    return write.call(this, key, value);
  });
  await act(async () => response.resolve(result(request)));
  expect(await screen.findByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation');
  expect(loadBasicRecordingIntent(original.scope_key)).toEqual(original);
  expect(api.operation).not.toHaveBeenCalled(); expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});
it('keeps original O1 queryable after saved revision advances, and never displays A under B', async () => {
  groups = [profileGroup()]; const response = deferred<RecordingStartOperation>(); api.start.mockReturnValue(response.promise);
  show(); const request = await confirmStart('circuit'); const original = retained(request);
  await selectDevice('PLC B'); groups = [{ ...groups[0], revision: 'new-group-revision' }];
  api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'new-workspace-head', groups: structuredClone(groups) });
  api.operation.mockResolvedValue(result(request)); await act(async () => response.resolve(result(request)));
  expect(screen.queryByTestId('basic-recording-status')).not.toBeInTheDocument(); expect(api.operation).not.toHaveBeenCalled();
  await selectDevice('PLC A'); await waitFor(() => expect(api.operation).toHaveBeenCalledWith('original-O1'));
  expect(loadBasicRecordingIntent(original.scope_key)).toMatchObject({ operation_id: 'original-O1', request });
  expect(request.expected_workspace_revision).toBe('W1'); expect(request.groups[0].expected_group_revision).toBe('rev-1');
  expect(api.start).toHaveBeenCalledOnce(); noOtherMutation();
});
