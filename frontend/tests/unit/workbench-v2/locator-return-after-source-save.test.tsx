import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceAPI, studioV2WorkspaceDevicesAPI } from './device-autosave-page.testHarness';
import { useRef, useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, renderHook, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { useStudioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import { locatorRead } from './studio-locator-fixtures';
import { openNewRecordingDevice } from './recording-editor-test-actions';
import { hydrateStudioV2Device } from '@/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
} }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/features/datalink/workbench-v2/recording-editor/recordingSourceDiscovery', () => ({ readRecordingSourceDiscovery: async (scope: object) => ({ state: 'available', scope, value: { candidates: [], excluded: [], devices: [] } }) }));
const groups = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const clients: QueryClient[] = [];
const existing = deviceFixture({ id: 'dev-1' });
let route: StudioLocatorEntryRead;
let revision: string;
let owner: ReturnType<typeof useStudioV2AutosaveState>;
function Harness({ entry }: { entry: StudioLocatorEntryRead }) {
  const [mounted, setMounted] = useState(false);
  owner = useStudioV2AutosaveState(false, entry, mounted, undefined, 'W');
  return <><output data-testid="hydration">{owner.locatorHydration.ready ? 'ready' : 'blocked'}</output>
    <output data-testid="receipts">{JSON.stringify(owner.deviceSource.observations())}</output>
    <div hidden={!owner.locatorHydration.ready && !mounted}>
      <Step4Database state={owner.state} dispatch={owner.actions.dispatch} workspaceId="W" locatorOnly locatorEntry={entry}
        deviceSource={owner.deviceSource} locatorSuspended={!owner.locatorHydration.ready} onLocatorEditorMountedChange={setMounted} />
    </div></>;
}
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const content = (entry: StudioLocatorEntryRead) => <QueryClientProvider client={client}><Harness entry={entry} /></QueryClientProvider>;
  const view = render(content(route)); return { ...view, navigate: (entry: StudioLocatorEntryRead) => view.rerender(content(entry)) };
}
beforeEach(() => {
  resetDeviceAutosaveMocks(); vi.clearAllMocks(); route = locatorRead(4); revision = route.entry.workspace_revision;
  groups.list.mockImplementation(async () => ({ ...route.entry.write_groups!, workspace_revision: revision }));
  groups.readiness.mockRejectedValue(new Error('unavailable')); groups.delivery.mockRejectedValue(new Error('unavailable'));
  studioV2WorkspaceDevicesAPI.list.mockResolvedValue([existing]);
  studioV2WorkspaceAPI.get.mockResolvedValue({ ...route.workspace, ordered_device_ids: [existing.id] });
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async (_locator, expected) => {
    if (expected.workspaceRevision !== revision) return { ok: false, reason: 'stale-reference' };
    const value = structuredClone(route); value.entry.workspace_revision = expected.workspaceRevision;
    value.entry.write_groups!.workspace_revision = expected.workspaceRevision;
    if (value.selection.state === 'selected') value.selection = { ...value.selection, target: { ...value.selection.target, workspaceRevision: expected.workspaceRevision } };
    return { ok: true, value: { ...value, entry: { ...value.entry, editor_hydrated: true, database_config: null, database_targets: [] } } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
async function saveSource(advanceRevision = false) {
  const view = mount(); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-name');
  await openNewRecordingDevice();
  fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Confirmed PLC' } });
  fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.20' } });
  studioV2WorkspaceDevicesAPI.create.mockImplementation(async payload => {
    const saved = deviceFixture({ id: payload.id, name: payload.name, protocol: payload.protocol, connection_config: JSON.stringify(payload.connection_config) });
    studioV2WorkspaceDevicesAPI.list.mockResolvedValue([existing, saved]);
    studioV2WorkspaceAPI.get.mockResolvedValue({ ...route.workspace, ordered_device_ids: [existing.id, saved.id] });
    if (advanceRevision) revision = 'W-after-source';
    return saved;
  });
  fireEvent.click(screen.getByTestId('recording-source-save'));
  await screen.findByTestId('recording-source-confirmed');
  await waitFor(() => expect(owner.state.devices.find(device => device.name === 'Confirmed PLC')).toMatchObject({ persisted: true, save_state: 'saved' }));
  fireEvent.click(screen.getByTestId('recording-range-back'));
  fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
  return view;
}
it('real confirmed source survives same-entry Back to groups and repeated editor open/close', async () => {
  const beforeGroup = JSON.stringify(route.entry.write_groups);
  await saveSource(); const saved = structuredClone(owner.state.devices); const receipts = owner.deviceSource.observations();
  fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(screen.getByTestId('group-open-G')).toBeVisible(); expect(screen.getByTestId('hydration')).toHaveTextContent('ready');
  expect(owner.state.devices).toEqual(saved); expect(owner.deviceSource.observations()).toEqual(receipts);
  fireEvent.click(screen.getByTestId('group-open-G')); await screen.findByTestId('group-name');
  fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(screen.getByTestId('group-open-G')).toBeVisible(); expect(owner.state.devices).toEqual(saved);
  expect(JSON.stringify(route.entry.write_groups)).toBe(beforeGroup);
  expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce(); expect(studioV2WorkspaceDevicesAPI.remove).not.toHaveBeenCalled();
  expect(groups.create).not.toHaveBeenCalled(); expect(groups.update).not.toHaveBeenCalled(); expect(groups.remove).not.toHaveBeenCalled();
});
it('reopening after source advances workspace revision requires fresh canonical editor proof', async () => {
  await saveSource(true); fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(screen.getByTestId('group-open-G')).toBeVisible();
  fireEvent.click(screen.getByTestId('group-open-G'));
  expect(await screen.findByTestId('group-name')).toBeVisible();
  expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce(); expect(groups.update).not.toHaveBeenCalled();
});

it.each(['saved', 'saving', 'save-error', 'draft-invalid'] as const)('same-entry mount transitions preserve %s local state but never clear the new-intent guard', saveState => {
  const first = locatorRead(4);
  const { result, rerender } = renderHook(({ entry, mounted }: { entry: StudioLocatorEntryRead; mounted: boolean }) => {
    const actions = useWorkbenchV2State(); const state = useRef(actions.state); state.current = actions.state;
    const hydration = useStudioV2LocatorHydration(actions, state, false, entry, mounted);
    return { actions, state, hydration };
  }, { initialProps: { entry: first, mounted: false } });
  rerender({ entry: first, mounted: true });
  act(() => result.current.actions.dispatch({ type: 'SET_STATE', payload: { db: { ...result.current.state.current.db,
    connector: { ...result.current.state.current.db.connector, name: 'Keep local draft', save_state: saveState } } } }));
  const retained = result.current.state.current;
  rerender({ entry: first, mounted: false });
  expect(result.current.hydration.ready).toBe(true); expect(result.current.state.current).toBe(retained);
  rerender({ entry: first, mounted: true });
  expect(result.current.hydration.ready).toBe(true); expect(result.current.state.current).toBe(retained);
  rerender({ entry: locatorRead(1, 'B'), mounted: true });
  expect(result.current.hydration.blocked).toBe(true); expect(result.current.state.current).toBe(retained);
  rerender({ entry: locatorRead(1, 'B'), mounted: false });
  expect(result.current.hydration.blocked).toBe(true); expect(result.current.state.current).toBe(retained);
});

it.each(['pending', 'unknown'] as const)('same-entry close retains real %s source receipt without cancellation or duplicate writes', async status => {
  mount(); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-name');
  const buffer = hydrateStudioV2Device(existing);
  let resolve!: (device: ReturnType<typeof deviceFixture>) => void;
  studioV2WorkspaceDevicesAPI.create.mockImplementation(() => new Promise(yes => { resolve = yes; }));
  await act(async () => { expect(await owner.deviceSource.prepare?.('W')).toBe(true); });
  let outcome!: ReturnType<typeof owner.deviceSource.submit>;
  act(() => { outcome = owner.deviceSource.submit({ workspaceId: 'W', sessionId: 'retained', requestId: 'retained-1', generation: 1,
    buffer: { name: 'Retained source', description: buffer.description, protocol: buffer.protocol, config: buffer.config } }); });
  await waitFor(() => expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce());
  if (status === 'unknown') {
    const payload = studioV2WorkspaceDevicesAPI.create.mock.calls[0][0];
    await act(async () => { resolve(deviceFixture({ id: payload.id, name: payload.name, protocol: payload.protocol, connection_config: JSON.stringify(payload.connection_config) })); await outcome; });
  }
  expect(owner.deviceSource.observations()[0].status).toBe(status);
  const before = owner.deviceSource.observations(); const state = structuredClone(owner.state.devices);
  fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.pending');
  fireEvent.click(screen.getByRole('button', { name: 'navigation.discard' }));
  expect(screen.getByTestId('group-open-G')).toBeVisible();
  expect(owner.deviceSource.observations()).toEqual(before); expect(owner.state.devices).toEqual(state);
  expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce(); expect(studioV2WorkspaceDevicesAPI.remove).not.toHaveBeenCalled();
  expect(groups.update).not.toHaveBeenCalled();
});

it.each(['list', 'editor'] as const)('late reopen %s response cannot replace a newer locator intent', async stage => {
  const view = await saveSource(true); fireEvent.click(screen.getByTestId('group-editor-close'));
  let resolve!: () => void;
  if (stage === 'list') groups.list.mockImplementationOnce(() => new Promise(yes => {
    resolve = () => yes({ ...route.entry.write_groups!, workspace_revision: revision });
  }));
  else vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(yes => {
    const value = structuredClone(route); value.entry.workspace_revision = revision; value.entry.editor_hydrated = true;
    value.entry.write_groups!.workspace_revision = revision;
    if (value.selection.state === 'selected') value.selection = { ...value.selection, target: { ...value.selection.target, workspaceRevision: revision } };
    resolve = () => yes({ ok: true, value });
  }));
  fireEvent.click(screen.getByTestId('group-open-G'));
  await waitFor(() => expect(resolve).toBeTypeOf('function'));
  const retained = structuredClone(owner.state.devices); const receipts = owner.deviceSource.observations();
  view.navigate(locatorRead(1, 'B'));
  await act(async () => { resolve(); });
  expect(owner.locatorHydration.blocked).toBe(true); expect(screen.queryByTestId('group-name')).not.toBeInTheDocument();
  expect(owner.state.devices).toEqual(retained); expect(owner.deviceSource.observations()).toEqual(receipts);
  expect(groups.update).not.toHaveBeenCalled(); expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
});
it.each(['foreign', 'removed', 'conflict'] as const)('reopen %s evidence remains safely unavailable without restoring stale editor', async reason => {
  await saveSource(true); fireEvent.click(screen.getByTestId('group-editor-close'));
  if (reason === 'foreign') groups.list.mockResolvedValueOnce({ ...route.entry.write_groups!, workspace_id: 'OTHER' });
  if (reason === 'removed') groups.list.mockResolvedValueOnce({ ...route.entry.write_groups!, groups: [] });
  if (reason === 'conflict') vi.mocked(readStudioLocatorEditorEntry).mockResolvedValueOnce({ ok: false, reason: 'stale-reference' });
  fireEvent.click(screen.getByTestId('group-open-G'));
  expect(await screen.findByTestId('studio-locator-open-error')).toBeVisible();
  expect(screen.getByTestId('group-open-G')).toBeVisible(); expect(screen.queryByTestId('group-name')).not.toBeInTheDocument();
  expect(owner.state.devices[0].name).toBe('Confirmed PLC'); expect(groups.update).not.toHaveBeenCalled();
});
it('preserves an explicit empty canonical workspace revision on reopen, without inventing a revision', async () => {
  mount(); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-name');
  fireEvent.click(screen.getByTestId('group-editor-close')); revision = '';
  fireEvent.click(screen.getByTestId('group-open-G'));
  expect(await screen.findByTestId('group-name')).toBeVisible();
  expect(readStudioLocatorEditorEntry).toHaveBeenLastCalledWith(route.locator,
    { workspaceRevision: '', groupRevision: 'G3' }, expect.any(AbortSignal));
});
