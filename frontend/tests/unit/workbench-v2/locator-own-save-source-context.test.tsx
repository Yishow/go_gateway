import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceAPI, studioV2WorkspaceDevicesAPI, studioV2RulesAPI, studioV2MappingsAPI } from './device-autosave-page.testHarness';
import { useState, type ComponentProps } from 'react';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import type { GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { deviceAPI } from '@/services/datalink';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import { locatorRead } from './studio-locator-fixtures';
import { TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/datalink', async original => {
  const actual = await original<typeof import('@/services/datalink')>();
  return { ...actual, deviceAPI: { ...actual.deviceAPI, checkReadiness: vi.fn() } };
});
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
} }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', async original => {
  const actual = await original<typeof import('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor')>();
  return { GroupEditor: (props: GroupEditorProps) => { editorProps = props; return <actual.GroupEditor {...props} />; } };
});
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar', async original => {
  const actual = await original<typeof import('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar')>();
  return { GroupLifecycleBar: (props: ComponentProps<typeof actual.GroupLifecycleBar>) => <><actual.GroupLifecycleBar {...props} />
    <button onClick={() => props.onChanged(current.entry.write_groups!.groups[0], current.entry.workspace_revision)}>Generic lifecycle update</button></> };
});
vi.mock('@/features/datalink/workbench-v2/recording-home/recordingHomeSources', () => ({
  readRecordingHomeSources: async (scope: object) => ({ state: 'available', scope, value: { candidates: discoveredCandidates, excluded: [], devices: [device] } }),
}));
const device = deviceFixture({ id: 'dev-1', name: 'PLC', last_test_success: true });
function editorRead(workspaceRevision = 'W2', groupRevision = 'G3'): StudioLocatorEntryRead {
  const read = locatorRead(4);
  read.workspace.ordered_device_ids = [device.id];
  read.entry.workspace_revision = workspaceRevision;
  read.entry.group_heads.G = groupRevision;
  read.entry.write_groups!.workspace_revision = workspaceRevision;
  read.entry.write_groups!.groups[0].revision = groupRevision;
  if (read.selection.state === 'selected') Object.assign(read.selection.target, { workspaceRevision, groupRevision });
  Object.assign(read.entry, { editor_hydrated: true, devices: [], source_rules: [], mappings: [], database_targets: [],
    database_config: { id: 'conn-1', identity_revision: 'crev-1', workspace_id: 'W', setup_revision: workspaceRevision, kind: 'sqlite',
      name: 'Saved database', host: '', port: 0, database: '/d.db', username: '', schema: '', table: 'readings',
      write_mode: 'insert', write_interval_seconds: 15, timestamp_column: 'at', status: 'ready', created_at: 'now', updated_at: 'now' },
  });
  return read;
}
let current: StudioLocatorEntryRead;
let editorProps: GroupEditorProps;
let discoveredCandidates: GroupCandidate[] = [];
const clients: QueryClient[] = [];
function Harness({ route }: { route: StudioLocatorEntryRead }) {
  const [editorMounted, setEditorMounted] = useState(false);
  const owner = useStudioV2AutosaveState(false, route, editorMounted, undefined, 'W');
  return <><output data-testid="source-receipts">{JSON.stringify(owner.deviceSource.range?.observations())}</output><Step4Database state={owner.state} dispatch={owner.actions.dispatch} workspaceId="W" locatorOnly locatorEntry={route} deviceSource={owner.deviceSource}
    locatorSuspended={!owner.locatorHydration.ready} onLocatorEditorMountedChange={setEditorMounted} /></>;
}
function mount(route: StudioLocatorEntryRead) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const content = (next: StudioLocatorEntryRead) => <QueryClientProvider client={client}><Harness route={next} /></QueryClientProvider>;
  const view = render(content(route)); return { ...view, navigate: (next: StudioLocatorEntryRead) => view.rerender(content(next)) };
}
beforeEach(() => {
  resetDeviceAutosaveMocks(); vi.clearAllMocks(); current = editorRead(); discoveredCandidates = [];
  studioV2WorkspaceAPI.get.mockImplementation(async () => ({ ...current.workspace, modbus_share: { workspace_id: 'W', workspace_revision: 'SOURCE-REV' } }));
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockRejectedValue(new Error('not ready'));
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockRejectedValue(new Error('unavailable'));
  studioV2WorkspaceDevicesAPI.list.mockResolvedValue([device]);
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: device.id, status: 'ready', connect_status: 'success', probe_status: 'success',
    activation_allowed: true, planning_allowed: true, apply_allowed: true, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockImplementation(async () => current.entry.write_groups!);
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async () => ({ ok: true, value: structuredClone(current) }));
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementation(async (_id, request) => {
    current = editorRead('W3', 'G4'); current.entry.write_groups!.groups[0].name = request.group.name;
    return { workspace_revision: 'W3', group: current.entry.write_groups!.groups[0] };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
async function openEditor(route = editorRead()) {
  mount(route); fireEvent.click(await screen.findByTestId('group-open-G'));
  await screen.findByTestId('group-name');
}
async function openUnsubmittedRange() {
  fireEvent.click(screen.getByTestId('recording-points-open'));
  await screen.findByRole('option', { name: 'PLC' });
  fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: device.id } });
  fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
  fireEvent.change(await screen.findByLabelText('range.name'), { target: { value: 'Unsent after own save' } });
  fireEvent.change(screen.getByLabelText('range.address'), { target: { value: '41301' } });
  fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '2' } });
}
it('first source range acquisition after own group Save uses acknowledged canonical readback, not the old locator hint', async () => {
  const route = editorRead(); const before = JSON.stringify(route);
  await openEditor(route);
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own saved name' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  await act(async () => { await Promise.resolve(); });
  expect(screen.getByTestId('group-revisions')).toHaveTextContent('G4');
  await openUnsubmittedRange();
  await waitFor(() => expect(screen.queryByText('device.checkingScope')).not.toBeInTheDocument());
  expect(screen.queryByText('device.scopeBlocked')).not.toBeInTheDocument();
  expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled();
  expect(JSON.stringify(route)).toBe(before);
  expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
});
it('fresh current locator remains a working control without any group Save', async () => {
  current = editorRead('W3', 'G4'); await openEditor(structuredClone(current));
  await openUnsubmittedRange();
  await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});

it('keeps an unsent range and newer recording draft mounted when own-save editor readback arrives late', async () => {
  await openEditor();
  await openUnsubmittedRange();
  await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
  fireEvent.click(screen.getByTestId('recording-range-back'));
  fireEvent.click(screen.getByRole('button', { name: 'device.discard' }));
  fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
  let finish!: (value: Awaited<ReturnType<typeof readStudioLocatorEditorEntry>>) => void;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own saved name' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Newer unsaved recording name' } });
  await openUnsubmittedRange();
  await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
  const rangeName = screen.getByLabelText('range.name'), picker = screen.getByTestId('recording-point-selection');
  await act(async () => { finish({ ok: true, value: structuredClone(current) }); });
  expect(screen.queryByLabelText('range.name')).toBe(rangeName);
  expect(rangeName).toHaveValue('Unsent after own save');
  expect(screen.getByTestId('recording-point-selection')).toBe(picker);
  expect(screen.getByTestId('group-name')).toHaveValue('Newer unsaved recording name');
  expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('recording-range-back'));
  fireEvent.click(screen.getByRole('button', { name: 'device.discard' }));
  fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
  expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled();
  expect(studioV2RulesAPI.remove).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.remove).not.toHaveBeenCalled();
});

it('keeps pending point selection through own-save canonical readback without silently adding members', async () => {
  discoveredCandidates = [{ key: 'pt-extra|tag-extra', point_id: 'pt-extra', tag_id: 'tag-extra', device_id: device.id,
    device_name: 'PLC', label: 'Extra point', tag_key: 'extra.point', address: '41301', target_type: 'float64' }];
  await openEditor();
  let finish!: (value: Awaited<ReturnType<typeof readStudioLocatorEditorEntry>>) => void;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own saved name' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const selected = await screen.findByRole('checkbox', { name: 'Extra point' }); fireEvent.click(selected);
  expect(selected).toBeChecked();
  await act(async () => { finish({ ok: true, value: structuredClone(current) }); });
  expect(screen.getByRole('checkbox', { name: 'Extra point' })).toBe(selected); expect(selected).toBeChecked();
  fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
  expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
  expect(current.entry.write_groups!.groups[0].members).toHaveLength(2);
});

it('manual reload after a foreign group revision cannot borrow an own-save source acknowledgement', async () => {
  await openEditor(); current = editorRead('W3', 'G4');
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unsent conflicting draft' } });
  fireEvent.click(screen.getByTestId('group-save'));
  fireEvent.click(await screen.findByTestId('group-conflict-reload'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  await openUnsubmittedRange();
  await screen.findByText('device.scopeBlocked');
  expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
});

it.each(['foreign-ack', 'newer-readback'])('does not authorize source acquisition from %s after a nominal group Save', async kind => {
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementationOnce(async (_id, request) => {
    current = editorRead(kind === 'newer-readback' ? 'W4' : 'W3', 'G4'); current.entry.write_groups!.groups[0].name = request.group.name;
    return { workspace_revision: 'W3', group: { ...current.entry.write_groups!.groups[0], workspace_id: kind === 'foreign-ack' ? 'OTHER' : 'W' } };
  });
  await openEditor(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own submitted name' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  await openUnsubmittedRange(); await screen.findByText('device.scopeBlocked');
  expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
  expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
});

it.each(['ack', 'readback'])('late own %s cannot refresh source authority after locator A -> B -> A navigation', async stage => {
  const route = editorRead(), view = mount(route);
  fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-name');
  let finish!: () => void;
  if (stage === 'ack') vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementationOnce(async (_id, request) => {
    current = editorRead('W3', 'G4'); current.entry.write_groups!.groups[0].name = request.group.name;
    await new Promise<void>(resolve => { finish = resolve; });
    return { workspace_revision: 'W3', group: current.entry.write_groups!.groups[0] };
  });
  else vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(async () => {
    await new Promise<void>(resolve => { finish = resolve; }); return { ok: true, value: structuredClone(current) };
  });
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own submitted name' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(finish).toBeTypeOf('function'));
  view.navigate(locatorRead(4, 'OTHER')); view.navigate(route);
  await act(async () => { finish(); });
  await openUnsubmittedRange(); await screen.findByText('device.scopeBlocked');
  expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
  expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
});

it('generic lifecycle two-argument notifications never become an own Group Save acknowledgement', async () => {
  await openEditor(); current = editorRead('W3', 'G4');
  fireEvent.click(screen.getByRole('button', { name: 'Generic lifecycle update' }));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  await openUnsubmittedRange(); await screen.findByText('device.scopeBlocked');
  expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});

it('own group acknowledgement preserves the same pending source request and its later unknown receipt', async () => {
  await openEditor();
  await openUnsubmittedRange(); await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
  fireEvent.click(screen.getByTestId('recording-range-back'));
  fireEvent.click(screen.getByRole('button', { name: 'device.discard' })); fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
  let finishGroup!: () => void, failSource!: (error: Error) => void;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementationOnce(async (_id, request) => {
    current = editorRead('W3', 'G4'); current.entry.write_groups!.groups[0].name = request.group.name;
    await new Promise<void>(resolve => { finishGroup = resolve; });
    return { workspace_revision: 'W3', group: current.entry.write_groups!.groups[0] };
  });
  studioV2RulesAPI.create.mockImplementationOnce(() => new Promise((_resolve, reject) => { failSource = reject; }));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own submitted name' } }); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(finishGroup).toBeTypeOf('function'));
  await openUnsubmittedRange(); await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
  fireEvent.click(screen.getByRole('button', { name: 'range.save' }));
  await waitFor(() => expect(studioV2RulesAPI.create).toHaveBeenCalledOnce());
  const rangeName = screen.getByLabelText('range.name'), pending = screen.getByTestId('source-receipts').textContent;
  expect(pending).toContain('"status":"pending"');
  await act(async () => { finishGroup(); });
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  expect(screen.getByLabelText('range.name')).toBe(rangeName);
  expect(screen.getByTestId('source-receipts')).toHaveTextContent(pending!);
  await act(async () => { failSource(new Error('uncertain source result')); });
  await screen.findByTestId('recording-range-unknown');
  const unknown = screen.getByTestId('source-receipts').textContent;
  expect(JSON.parse(unknown!)[0].request).toEqual(JSON.parse(pending!)[0].request);
  expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
  expect(screen.getByLabelText('range.name')).toBe(rangeName);
  expect(studioV2RulesAPI.create).toHaveBeenCalledOnce(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
});

it('a prior own acknowledgement cannot preserve a later foreign picker session on workspace revision A -> B -> A', async () => {
  discoveredCandidates = [{ key: 'pt-extra|tag-extra', point_id: 'pt-extra', tag_id: 'tag-extra', device_id: device.id,
    device_name: 'PLC', label: 'Extra point', tag_key: 'extra.point', address: '41301', target_type: 'float64' }];
  await openEditor();
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Own saved name' } }); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const first = await screen.findByRole('checkbox', { name: 'Extra point' }); fireEvent.click(first);
  current = editorRead('W4', 'G4');
  fireEvent.click(screen.getByRole('button', { name: 'Generic lifecycle update', hidden: true }));
  await waitFor(() => expect(first).not.toBeInTheDocument());
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const foreign = await screen.findByRole('checkbox', { name: 'Extra point' }); fireEvent.click(foreign);
  current = editorRead('W3', 'G4');
  fireEvent.click(screen.getByRole('button', { name: 'Generic lifecycle update', hidden: true }));
  await waitFor(() => expect(foreign).not.toBeInTheDocument());
  expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
});

it('serializes repeated own Saves through canonical readback without losing the next unsent draft', async () => {
  await openEditor(); let finishFirst!: (value: Awaited<ReturnType<typeof readStudioLocatorEditorEntry>>) => void, saves = 0;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(resolve => { finishFirst = resolve; }));
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementation(async (_id, request) => {
    saves += 1; current = editorRead(saves === 1 ? 'W3' : 'W4', saves === 1 ? 'G4' : 'G5');
    current.entry.write_groups!.groups[0].name = request.group.name;
    return { workspace_revision: current.entry.workspace_revision, group: current.entry.write_groups!.groups[0] };
  });
  const cachedFetchLatest = editorProps.fetchLatest;
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'First own save' } }); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  const first = structuredClone(current);
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Second own save' } });
  // Otherwise readback1 can invalidate an already-sent Save2 lease, or ACK2 can
  // race the previous-basis handoff. This is only a save fence, not readonly.
  expect(screen.getByTestId('group-save')).toBeDisabled();
  const readsBeforeCachedSave = vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mock.calls.length;
  await expect(cachedFetchLatest()).rejects.toThrow('Recording scope is unavailable.');
  expect(studioV2WorkspaceWriteGroupsAPI.list).toHaveBeenCalledTimes(readsBeforeCachedSave);
  fireEvent.click(screen.getByTestId('group-save'));
  expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
  expect(screen.getByTestId('group-name')).toBeEnabled();
  expect(screen.getByTestId('recording-points-open')).toBeEnabled();
  await act(async () => { finishFirst({ ok: true, value: first }); });
  expect(screen.getByTestId('group-name')).toHaveValue('Second own save');
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledTimes(2));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenLastCalledWith(expect.anything(), { workspaceRevision: 'W4', groupRevision: 'G5' }, expect.any(AbortSignal)));
  expect(screen.getByTestId('group-revisions')).toHaveTextContent('G5');
  expect(screen.getByTestId('group-name')).toHaveValue('Second own save');
  await openUnsubmittedRange(); await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
  expect(screen.queryByText('device.scopeBlocked')).not.toBeInTheDocument();
});

it('releases the save-only fence on failed readback and exposes the existing error/reload path', async () => {
  await openEditor(); let finish!: (value: Awaited<ReturnType<typeof readStudioLocatorEditorEntry>>) => void;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Saved before failed read' } }); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep unsent edits' } });
  expect(screen.getByTestId('group-save')).toBeDisabled();
  await act(async () => { finish({ ok: false, reason: 'unavailable' }); });
  expect(screen.getByTestId('studio-locator-open-error')).toHaveTextContent('locatorEntry.unavailable');
  expect(screen.getByTestId('group-save')).toBeEnabled();
  expect(screen.getByTestId('group-name')).toHaveValue('Keep unsent edits');
  await act(async () => { editorProps.onReload(); });
  await waitFor(() => expect(screen.queryByTestId('studio-locator-open-error')).not.toBeInTheDocument());
  expect(screen.getByTestId('group-save')).toBeEnabled();
  expect(screen.getByTestId('group-name')).toHaveValue('Keep unsent edits');
  expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
});
