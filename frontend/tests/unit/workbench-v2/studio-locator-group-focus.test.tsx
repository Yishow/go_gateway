import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { locatorRead } from './studio-locator-fixtures';
import { GroupLifecycleBar } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { recordingTestWriteRecoveryKey } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import type { SavedRecordingActionIntent } from '@/types/studioV2WriteGroupLifecycle';
import type { StudioLocatorEntryResult } from '@/types/studioV2LocatorEntry';
const captured = vi.hoisted(() => ({ lifecycle: undefined as ((intent: SavedRecordingActionIntent) => void) | undefined }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions', async importOriginal => {
  const actual = await importOriginal<typeof import('@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions')>();
  return { ...actual, SavedRecordingActions: (props: React.ComponentProps<typeof actual.SavedRecordingActions>) => {
    captured.lifecycle = props.onAction; return <actual.SavedRecordingActions {...props} />;
  } };
});
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', async importOriginal => {
  const actual = await importOriginal<typeof import('@/services/studioV2WorkspaceWriteGroups')>();
  return { studioV2WorkspaceWriteGroupsAPI: { ...actual.studioV2WorkspaceWriteGroupsAPI, list: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn() } };
});
const native = vi.hoisted(() => ({ editor: false }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', async importOriginal => {
  const actual = await importOriginal<typeof import('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor')>();
  return { GroupEditor: (props: GroupEditorProps) => native.editor ? <actual.GroupEditor {...props} /> :
    <div data-testid="native-editor"><input aria-label="Existing draft" defaultValue={props.group?.name} /><span data-testid="editor-revision">{props.group?.revision}</span>
      <button onClick={props.onReload}>Reload native group</button><button onClick={props.onClose}>Close native group</button>
      <button onClick={() => props.group && props.onSaved({ ...props.group, revision: 'G4' })}>Accepted native save</button></div> };
});
const clients: QueryClient[] = [];
const entry = () => locatorRead(4);
function mount(read = entry()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
  const renderContent = (next: typeof read, suspended = false) => <QueryClientProvider client={client}><Step4Database state={studioV2LocatorHydration(read)!}
    dispatch={vi.fn()} workspaceId="W" locatorOnly locatorEntry={next} locatorSuspended={suspended} /></QueryClientProvider>;
  return { ...render(renderContent(read)), renderContent };
}
beforeEach(() => { vi.resetAllMocks(); captured.lifecycle = undefined; native.editor = false; vi.mocked(studioV2DatalinkApi.get).mockRejectedValue({ response: { status: 503 } }); vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockRejectedValue(new Error('unavailable')); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); });
describe('Step4 locator focuses native saved card, never mounts an editor automatically', () => {
  it('focuses exact group with historical hint and unknown destination without mounting Basic or recovery', async () => {
    sessionStorage.setItem('wbv2.basic-legacy-draft.W|D|database', '{"intervalSeconds":31}'); sessionStorage.setItem('wbv2.circuit-profile-draft.W%2FD', '{"circuit_name":"Unsaved circuit"}');
    const before = JSON.stringify(sessionStorage); mount();
    const card = await screen.findByTestId('group-list-item-G'); expect(card).toHaveAttribute('data-locator-selected', 'true'); expect(card).toHaveFocus();
    expect(screen.queryByTestId('native-editor')).not.toBeInTheDocument(); expect(screen.queryByTestId('basic-recording-panel')).not.toBeInTheDocument();
    expect(screen.queryByTestId('workspace-test-write-recovery')).not.toBeInTheDocument(); expect(screen.queryByTestId('group-list-mismatch-G')).not.toBeInTheDocument();
    expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled(); expect(readStudioLocatorEditorEntry).not.toHaveBeenCalled(); expect(JSON.stringify(sessionStorage)).toBe(before);
  });
  it('explicit Open is singleflight and uses current W2/G3, never historical G1 as CAS', async () => {
    let resolve!: (value: StudioLocatorEntryResult) => void; vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(done => { resolve = done; }));
    const read = entry(); mount(read); fireEvent.click(screen.getByTestId('group-open-G')); fireEvent.click(screen.getByTestId('group-open-G'));
    expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce(); expect(readStudioLocatorEditorEntry).toHaveBeenCalledWith(read.locator,
      { workspaceRevision: 'W2', groupRevision: 'G3' }, expect.any(AbortSignal));
    await act(async () => { resolve({ ok: true, value: { ...read, entry: { ...read.entry, editor_hydrated: true, database_config: null, database_targets: [], devices: [], source_rules: [], mappings: [] } } }); });
    expect(await screen.findByTestId('native-editor')).toBeInTheDocument(); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
  it.each(['foreign-scope', 'stale-reference', 'unavailable'] as const)('failed explicit Open %s does not create or substitute another group', async reason => {
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: false, reason }); mount(); fireEvent.click(screen.getByTestId('group-open-G'));
    await waitFor(() => expect(screen.getByTestId('studio-locator-open-error')).toBeInTheDocument()); expect(screen.queryByTestId('native-editor')).not.toBeInTheDocument();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
  it('query-only navigation cannot replace an already open native draft', async () => {
    const read = entry(); vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: { ...read, entry: { ...read.entry, editor_hydrated: true, database_config: null, database_targets: [] } } });
    const view = mount(read); fireEvent.click(screen.getByTestId('group-open-G')); const input = await screen.findByLabelText('Existing draft');
    fireEvent.change(input, { target: { value: 'Unsaved editor draft' } }); view.rerender(view.renderContent(locatorRead(4, 'OTHER')));
    expect(screen.getByLabelText('Existing draft')).toHaveValue('Unsaved editor draft'); expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce();
  });
  it('ordinary list still opens its native editor directly', async () => {
    const read = entry(); vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(read.entry.write_groups!);
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
    render(<QueryClientProvider client={client}><WriteGroupSection state={studioV2LocatorHydration(read)!} workspaceId="W" readonly={false} /></QueryClientProvider>);
    fireEvent.click(await screen.findByTestId('group-open-G')); expect(screen.getByTestId('native-editor')).toBeInTheDocument(); expect(readStudioLocatorEditorEntry).not.toHaveBeenCalled();
  });
  it('aborts an unfinished Open on a new locator and never mounts its late editor', async () => {
    let resolve!: (value: StudioLocatorEntryResult) => void;
    vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(done => { resolve = done; }));
    const read = entry(), view = mount(read); fireEvent.click(screen.getByTestId('group-open-G'));
    const signal = vi.mocked(readStudioLocatorEditorEntry).mock.calls[0][2]!;
    view.rerender(view.renderContent(locatorRead(4, 'OTHER'))); expect(signal.aborted).toBe(true);
    await act(async () => { resolve({ ok: true, value: { ...read, entry: { ...read.entry, editor_hydrated: true, database_config: null, database_targets: [] } } }); });
    expect(screen.queryByTestId('native-editor')).not.toBeInTheDocument(); expect(screen.getByTestId('group-list-item-OTHER')).toBeInTheDocument();
  });
  it('rejects a deleted exact group in editor response rather than same-name replacement', async () => {
    const read = entry(), deleted = { ...read.entry.write_groups!.groups[0], status: 'deleted' as const };
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: { ...read, entry: { ...read.entry, editor_hydrated: true,
      database_config: null, database_targets: [], write_groups: { ...read.entry.write_groups!, groups: [deleted, { ...deleted, id: 'G2', status: 'ready' }] } } } });
    mount(read); fireEvent.click(screen.getByTestId('group-open-G')); await screen.findByTestId('studio-locator-open-error');
    expect(screen.queryByTestId('native-editor')).not.toBeInTheDocument(); expect(screen.queryByTestId('group-list-item-G2')).not.toBeInTheDocument();
  });
  it('adopts accepted Save and Reload using fresh same-W/G proof, never the URL hint', async () => {
    const read = entry(); const opened = { ...read, entry: { ...read.entry, editor_hydrated: true as const, database_config: null, database_targets: [] } };
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValueOnce({ ok: true, value: opened });
    const next = structuredClone(opened); next.entry.workspace_revision = 'W3'; next.entry.write_groups!.workspace_revision = 'W3'; next.entry.write_groups!.groups[0].revision = 'G4';
    if (next.selection.state === 'selected') next.selection = { ...next.selection, target: { ...next.selection.target, workspaceRevision: 'W3', groupRevision: 'G4' } };
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValueOnce({ ok: true, value: next });
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(next.entry.write_groups!);
    mount(read); fireEvent.click(screen.getByTestId('group-open-G')); await screen.findByTestId('native-editor');
    fireEvent.click(screen.getByText('Accepted native save')); await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
    expect(readStudioLocatorEditorEntry).toHaveBeenLastCalledWith(read.locator, { workspaceRevision: 'W3', groupRevision: 'G4' }, expect.any(AbortSignal));
    expect(screen.getByTestId('editor-revision')).toHaveTextContent('G4');
  });
  it.each(['deleted', 'foreign'] as const)('post-Save %s current proof does not recreate or substitute the editor', async kind => {
    const read = entry(); vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: { ...read, entry: { ...read.entry, editor_hydrated: true, database_config: null, database_targets: [] } } });
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: kind === 'foreign' ? 'OTHER' : 'W', workspace_revision: 'W3', groups: [] });
    mount(read); fireEvent.click(screen.getByTestId('group-open-G')); const input = await screen.findByLabelText('Existing draft');
    fireEvent.change(input, { target: { value: 'Retain local edits' } }); fireEvent.click(screen.getByText('Accepted native save')); await screen.findByTestId('studio-locator-open-error');
    expect(input).toHaveValue('Retain local edits'); expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce(); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
  it('canceling a late Reload never reopens or replaces the closed editor', async () => {
    const read = entry(), opened = { ...read, entry: { ...read.entry, editor_hydrated: true as const, database_config: null, database_targets: [] } };
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValueOnce({ ok: true, value: opened });
    let resolve!: (value: StudioLocatorEntryResult) => void;
    vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(done => { resolve = done; }));
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(read.entry.write_groups!);
    mount(read); fireEvent.click(screen.getByTestId('group-open-G')); await screen.findByTestId('native-editor'); fireEvent.click(screen.getByText('Reload native group'));
    await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2)); fireEvent.click(screen.getByText('Close native group'));
    expect(vi.mocked(readStudioLocatorEditorEntry).mock.calls[1][2]!.aborted).toBe(true);
    await act(async () => { resolve({ ok: true, value: opened }); }); expect(screen.queryByTestId('native-editor')).not.toBeInTheDocument();
  });
  it('hot normal Step4 keeps the actual Basic circuit draft subtree and pending recovery bytes on locator navigation', async () => {
    const read = entry(), state = savedGroupState();
    state.db.connector.workspace_id = 'W';
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(read.entry.write_groups!);
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
    const dispatch = vi.fn();
    const content = (locatorEntry?: typeof read) => <QueryClientProvider client={client}><Step4Database state={state} dispatch={dispatch}
      workspaceId="W" locatorEntry={locatorEntry} /></QueryClientProvider>;
    const view = render(content());
    const field = await screen.findByTestId('circuit-profile-circuit-name');
    fireEvent.change(field, { target: { value: 'Keep my unsaved circuit' } });
    await waitFor(() => expect(field).toHaveValue('Keep my unsaved circuit'));
    const basic = screen.getByTestId('step4-basic-recording');
    const sessionBefore = JSON.stringify(sessionStorage), recoveryBefore = localStorage.getItem(recordingTestWriteRecoveryKey);
    dispatch.mockClear(); vi.mocked(studioV2DatalinkApi.post).mockClear(); vi.mocked(studioV2DatalinkApi.put).mockClear();
    view.rerender(content(read));
    expect(screen.getByTestId('step4-basic-recording')).toBe(basic); expect(screen.getByTestId('circuit-profile-circuit-name')).toBe(field);
    expect(field).toHaveValue('Keep my unsaved circuit'); expect(JSON.stringify(sessionStorage)).toBe(sessionBefore);
    expect(localStorage.getItem(recordingTestWriteRecoveryKey)).toBe(recoveryBefore); expect(dispatch).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
    expect(readStudioLocatorEditorEntry).not.toHaveBeenCalled();
  });
  it.each(['suspend', 'roundtrip', 'unchanged'] as const)('deferred native Save uses a non-revivable locator lease: %s', async transition => {
    native.editor = true;
    const read = entry(), group = read.entry.write_groups!.groups[0];
    const opened = { ...read, entry: { ...read.entry, editor_hydrated: true as const, devices: [], source_rules: [], mappings: [], database_targets: [],
      database_config: { id: 'conn-1', identity_revision: 'crev-1', workspace_id: 'W', setup_revision: 'W2', kind: 'sqlite' as const,
        name: 'Saved database', host: '', port: 0, database: '/d.db', username: '', schema: '', table: 'readings',
        write_mode: 'insert' as const, write_interval_seconds: 15, timestamp_column: 'at', status: 'ready', created_at: 'now', updated_at: 'now' } } };
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: opened });
    let finish!: (value: NonNullable<typeof read.entry.write_groups>) => void;
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; })).mockRejectedValue({ response: { status: 503 } });
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockResolvedValue({ workspace_revision: 'W3', group: { ...group, name: 'Explicit draft name', revision: 'G4' } });
    const view = mount(read); fireEvent.click(screen.getByTestId('group-open-G'));
    const input = await screen.findByTestId('group-name'); fireEvent.change(input, { target: { value: 'Explicit draft name' } });
    expect(screen.getByTestId('group-save')).toBeEnabled();
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.list).toHaveBeenCalledOnce()); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
    if (transition === 'suspend' || transition === 'roundtrip') view.rerender(view.renderContent(locatorRead(4, 'OTHER'), true));
    if (transition === 'roundtrip') view.rerender(view.renderContent(read, false));
    await act(async () => { finish(read.entry.write_groups!); });
    if (transition === 'unchanged') await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
    else expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
    expect(input).toHaveValue('Explicit draft name'); expect(screen.getByTestId('group-editor-close')).toBeEnabled();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it.each(['stop-suspend', 'stop-roundtrip', 'stop-unchanged', 'remove-suspend', 'remove-roundtrip', 'remove-unchanged'])('deferred lifecycle refetch uses a non-revivable locator lease: %s', async scenario => {
    const read = entry(), group = read.entry.write_groups!.groups[0], removing = scenario.startsWith('remove');
    const lifecycle = savedLifecycle(group, 'W2'); let reads = 0;
    let finish!: (value: typeof lifecycle) => void;
    vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockImplementation(async () => {
      if (++reads === (removing ? 3 : 2)) return new Promise(resolve => { finish = resolve; });
      return lifecycle;
    });
    const accepted = { data: { success: true, data: { workspace_revision: 'W3', group } } };
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue(accepted); vi.mocked(studioV2DatalinkApi.delete).mockResolvedValue(accepted);
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
    const onChanged = vi.fn();
    const content = (readonly = false) => <QueryClientProvider client={client}><GroupLifecycleBar group={group} workspaceId="W" workspaceRevision="W2"
      locatorScope={{ workspaceId: 'W', workspaceRevision: 'W2' }} readonly={readonly} canApply={false} onChanged={onChanged} onConflict={vi.fn()} /></QueryClientProvider>;
    const view = render(content());
    const action = await screen.findByTestId(removing ? 'group-delete' : 'group-disable'); await waitFor(() => expect(action).toBeEnabled()); fireEvent.click(action);
    if (removing) { const confirm = await screen.findByTestId('group-delete-confirm-button'); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled(); fireEvent.click(confirm); }
    await waitFor(() => expect(reads).toBe(removing ? 3 : 2));
    if (!scenario.endsWith('unchanged')) view.rerender(content(true));
    if (scenario.endsWith('roundtrip')) view.rerender(content(false));
    await act(async () => { finish(lifecycle); });
    const mutation = removing ? studioV2DatalinkApi.delete : studioV2DatalinkApi.post;
    if (scenario.endsWith('unchanged')) await waitFor(() => expect(mutation).toHaveBeenCalledOnce());
    else { expect(mutation).not.toHaveBeenCalled(); expect(onChanged).not.toHaveBeenCalled(); }
  });
  it.each([false, true])('captured lifecycle callback cannot borrow a newer lease; suspended=%s', async suspended => {
    const read = entry(), group = read.entry.write_groups!.groups[0], lifecycle = savedLifecycle(group, 'W2');
    const get = vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(lifecycle);
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: { workspace_revision: 'W3', group } } });
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
    const content = (readonly = false) => <QueryClientProvider client={client}><GroupLifecycleBar group={group} workspaceId="W" workspaceRevision="W2"
      locatorScope={{ workspaceId: 'W', workspaceRevision: 'W2' }} readonly={readonly} canApply={false} onChanged={vi.fn()} onConflict={vi.fn()} /></QueryClientProvider>;
    const view = render(content()); await waitFor(() => expect(screen.getByTestId('group-disable')).toBeEnabled());
    const oldAction = captured.lifecycle; expect(oldAction).toBeTypeOf('function');
    if (suspended) view.rerender(content(true));
    await act(async () => { oldAction!({ action: 'stop', workspaceId: 'W', workspaceRevision: 'W2', groupId: 'G', groupRevision: 'G3' }); });
    if (suspended) { expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(get).toHaveBeenCalledOnce(); }
    else expect(studioV2DatalinkApi.post).toHaveBeenCalledOnce();
  });
  it.each(['incoming', 'roundtrip', 'ordinary'] as const)('hot ordinary editor pending Save cannot cross locator entry: %s', async transition => {
    native.editor = true;
    const read = entry(), state = savedGroupState(), group = read.entry.write_groups!.groups[0]; state.db.connector.workspace_id = 'W';
    let finish!: (value: NonNullable<typeof read.entry.write_groups>) => void;
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValueOnce(read.entry.write_groups!)
      .mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockResolvedValue({ workspace_revision: 'W3', group: { ...group, name: 'Keep hot draft', revision: 'G4' } });
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
    const content = (incoming = false) => <QueryClientProvider client={client}><WriteGroupSection state={state} workspaceId="W" readonly={false}
      locatorEntry={incoming ? locatorRead(4, 'OTHER') : undefined} locatorSuspended={incoming} /></QueryClientProvider>;
    const view = render(content()); fireEvent.click(await screen.findByTestId('group-open-G'));
    const input = await screen.findByTestId('group-name'); fireEvent.change(input, { target: { value: 'Keep hot draft' } });
    fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.list).toHaveBeenCalledTimes(2));
    if (transition !== 'ordinary') view.rerender(content(true));
    if (transition === 'roundtrip') view.rerender(content(false));
    await act(async () => { finish(read.entry.write_groups!); });
    if (transition === 'ordinary') expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
    else expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
    expect(input).toHaveValue('Keep hot draft'); expect(readStudioLocatorEditorEntry).not.toHaveBeenCalled();
  });
  it.each(['mysql', 'sqlserver'] as const)('actual cumulative %s editor forwards validated scope to every automatic mutation-bearing GET, not cold card', async kind => {
    native.editor = true;
    const read = entry(), group = read.entry.write_groups!.groups[0];
    group.layout_version = 'managed_long_v1'; group.destination.storage_strategy = 'managed';
    group.recording_profile = { kind: 'cumulative', interval: 15, unit: 'minute' };
    const editor = { ...read, entry: { ...read.entry, editor_hydrated: true as const, devices: [], source_rules: [], mappings: [], database_targets: [],
      database_config: { id: 'conn-1', identity_revision: 'crev-1', workspace_id: 'W', setup_revision: 'W2', kind,
        name: 'Saved database', host: 'db.example', port: 3306, database: '/d.db', username: 'reader', schema: '', table: 'readings',
        write_mode: 'insert' as const, write_interval_seconds: 15, timestamp_column: 'at', status: 'ready', created_at: 'now', updated_at: 'now' } } };
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: editor });
    const retained = JSON.stringify([{ workspace_id: 'W', group_id: 'G', group_revision: 'G3', operation_id: 'pending-G',
      operation_kind: 'test_write', created_at: '2026-10-29T00:00:00Z' }]);
    localStorage.setItem(recordingTestWriteRecoveryKey, retained);
    sessionStorage.setItem('wbv2.circuit-profile-draft.W%2FD', 'unsaved-circuit');
    const before = JSON.stringify(sessionStorage); mount(read);
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled(); expect(localStorage.getItem(recordingTestWriteRecoveryKey)).toBe(retained);
    fireEvent.click(screen.getByTestId('group-open-G')); await screen.findByTestId('group-editor');
    const guarded = ['/studio-v2/workspace/recording-plans/capabilities', '/studio-v2/workspace/database-metadata',
      '/studio-v2/workspace/measurements', '/studio-v2/workspace/database-operations/test-writes', '/studio-v2/workspace/database-operations/pending-G'];
    await waitFor(() => guarded.forEach(path => expect(studioV2DatalinkApi.get).toHaveBeenCalledWith(path,
      expect.objectContaining({ params: expect.objectContaining({ locator_workspace_id: 'W', locator_workspace_revision: 'W2' }) }))));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/write-groups/G/readiness');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/write-groups/G/lifecycle');
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
    expect(localStorage.getItem(recordingTestWriteRecoveryKey)).toBe(retained); expect(JSON.stringify(sessionStorage)).toBe(before);
    expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled();
  });

});
