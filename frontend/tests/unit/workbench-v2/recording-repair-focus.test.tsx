import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import type { StudioRecordingRepair } from '@/services/studioRecordingView';
import type { StudioLocatorEntryRead, StudioLocatorEntryResult } from '@/types/studioV2LocatorEntry';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { locatorRead, nativeDevice, nativeMapping, nativeRule } from './studio-locator-fixtures';
import { WorkbenchV2Shell } from '@/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { buildRecordingRepairTarget } from '@/services/studioRecordingView';
const native = vi.hoisted(() => ({ query: vi.fn(), mutation: vi.fn(), schemaLookup: true }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', () => ({ GroupEditor: (props: GroupEditorProps) => <div data-testid="group-editor">
  <section id="recording-editor-sources" tabIndex={-1}>Sources {props.group?.id}</section>
  <section id="recording-editor-destination" tabIndex={-1}>Destination</section>
  <input aria-label="local draft" />
  {native.schemaLookup && <button data-testid="recording-preparation-lookup" onClick={native.query}>Original schema lookup</button>}
  <button data-testid="group-schema-check" onClick={native.mutation}>New schema check</button>
  <button data-testid="group-test-write-recovery-query-O1" onClick={native.query}>Original test-write lookup</button>
  <button onClick={native.mutation}>Original command</button>
  {props.repairContent}
</div> }));
vi.mock('@/features/datalink/workbench-v2/recording-editor/LocatorSavedRecordingStartReview', () => ({ LocatorSavedRecordingStartReview: ({ entry }: { entry: StudioLocatorEntryRead }) =>
  <section data-testid="original-saved-start-review">{entry.locator.group_id}<button onClick={native.query}>Original Start query</button></section> }));
const clients: QueryClient[] = [];
const repair: StudioRecordingRepair = { workspaceId: 'W', workspaceRevision: 'W2', groupId: 'G', groupRevision: 'G3', section: 'sources' };
function entry(id = 'G'): StudioLocatorEntryRead {
  const read = locatorRead(4, id); read.locator.group_revision = 'G3';
  return { ...read, entry: { ...read.entry, editor_hydrated: true, database_config: null, database_targets: [] } };
}
function startEntry(strategy: 'managed' | 'custom' = 'managed') {
  const read = entry(), group = read.entry.write_groups!.groups[0];
  group.destination.storage_strategy = strategy;
  group.members = [{ ...group.members[0], device_id: 'D', point_id: 'P-native', tag_id: 'T' }];
  read.entry.devices = [nativeDevice()]; read.entry.source_rules = [nativeRule()]; read.entry.mappings = [nativeMapping()];
  read.entry.hydration_scope = { complete_workspace: false, device_ids: ['D'], rule_ids: ['R'], mapping_ids: ['M'], group_ids: ['G'] };
  read.entry.rule_heads = { R: 'R2' };
  read.entry.database_config = { id: group.destination.connector_id, identity_revision: group.destination.connector_revision, workspace_id: 'W', setup_revision: 'W2',
    kind: 'sqlite', name: 'Saved destination', host: '', port: 0, database: group.destination.database, username: '', schema: group.destination.table_schema,
    table: group.destination.table_name, write_mode: 'insert', write_interval_seconds: 15, timestamp_column: '', status: 'ready', save_state: 'saved', created_at: '', updated_at: '' };
  return read;
}
function retainStart(read: StudioLocatorEntryRead) {
  const group = read.entry.write_groups!.groups[0];
  saveBasicRecordingIntent({ scope_key: 'W/saved-record/G', operation_id: 'O1', owner: { kind: 'saved_record', groupId: 'G', deviceIds: ['D'] },
    request: { request_id: 'original-request', workspace_id: 'W', expected_workspace_revision: 'W2', device_ids: ['D'],
    groups: [{ group_id: 'G', expected_group_revision: 'G3', expected_connector_revision: group.destination.connector_revision }] } });
}
function mount(intent = repair, initial = entry()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
  const advanced = vi.fn();
  const content = (read = initial, next = intent, suspended = false) => <QueryClientProvider client={client}>
    <Step4Database state={studioV2LocatorHydration(read)!} dispatch={vi.fn()} workspaceId="W" presentation="simple" initialSelection="located"
      locatorOnly locatorEntry={read} locatorSuspended={suspended} repairIntent={next} onRepairAdvanced={advanced} />
  </QueryClientProvider>;
  return { ...render(content()), content, advanced };
}
beforeEach(() => { vi.clearAllMocks(); native.schemaLookup = true; sessionStorage.clear(); localStorage.clear(); vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: entry() }); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); sessionStorage.clear(); localStorage.clear(); });
describe('recording repair reaches only its fresh original editor', () => {
  it('passes the bounded repair from the real shell through Step4 to the fresh editor', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
    const read = entry(), location = new URL(buildRecordingRepairTarget(repair)!, 'http://localhost');
    function Shell() {
      const actions = useWorkbenchV2State();
      return <WorkbenchV2Shell state={actions.state} actions={actions} workspaceId="W" recordingIntent="located"
        locatorEntry={read} locatorOnly presentationLocation={{ search: location.search, hash: '#same-row' }} />;
    }
    render(<QueryClientProvider client={client}><Shell /></QueryClientProvider>);
    await waitFor(() => expect(document.getElementById('recording-editor-sources')).toHaveFocus());
    expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce(); expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
  });
  it.each(['sources', 'destination'] as const)('focuses %s only after the canonical Open read and never invokes a command', async section => {
    let finish!: (value: StudioLocatorEntryResult) => void;
    vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(resolve => { finish = resolve; }));
    mount({ ...repair, section });
    await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
    expect(screen.queryByTestId('group-editor')).toBeNull();
    await act(async () => finish({ ok: true, value: entry() }));
    await waitFor(() => expect(document.getElementById(`recording-editor-${section}`)).toHaveFocus());
    expect(native.mutation).not.toHaveBeenCalled(); expect(native.query).not.toHaveBeenCalled();
    screen.getByLabelText('local draft').focus();
    fireEvent.change(screen.getByLabelText('local draft'), { target: { value: 'kept' } });
    expect(screen.getByLabelText('local draft')).toHaveFocus();
  });
  it.each(['schema', 'test_write'] as const)('locates the original %s query without automatically invoking it', async owner => {
    if (owner === 'schema') retainSchemaOperation('W', 'G', 'O1');
    else recordingTestWriteRecovery.retain({ workspace_id: 'W', group_id: 'G', group_revision: 'G2', operation_id: 'O1', operation_kind: 'test_write', created_at: '2026-11-01T00:00:00Z' });
    mount({ ...repair, section: 'operation', operation: { kind: 'query_operation', owner, deviceId: 'dev-1', operationId: 'O1' } });
    const control = await screen.findByTestId(owner === 'schema' ? 'recording-preparation-lookup' : 'group-test-write-recovery-query-O1');
    await waitFor(() => expect(control).toHaveFocus());
    expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
    fireEvent.click(control); expect(native.query).toHaveBeenCalledOnce(); expect(native.mutation).not.toHaveBeenCalled();
  });
  it('never substitutes a new schema check when the original operation lookup has not mounted', async () => {
    retainSchemaOperation('W', 'G', 'O1'); native.schemaLookup = false;
    const view = mount({ ...repair, section: 'operation', operation: { kind: 'query_operation', owner: 'schema', deviceId: 'dev-1', operationId: 'O1' } });
    expect(await screen.findByTestId('group-schema-check')).not.toHaveFocus();
    expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
    native.schemaLookup = true; view.rerender(view.content());
    await waitFor(() => expect(screen.getByTestId('recording-preparation-lookup')).toHaveFocus());
    expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
  });
  it('opens only the retained saved-record Start in the original review after an explicit click', async () => {
    const read = startEntry(); retainStart(read);
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: read });
    mount({ ...repair, section: 'operation', operation: { kind: 'query_operation', owner: 'start', deviceId: 'D', operationId: 'O1' } }, read);
    const open = await screen.findByTestId('recording-repair-start'); expect(open).toHaveFocus();
    expect(screen.getByTestId('group-editor')).toContainElement(open);
    expect(screen.getAllByTestId('recording-repair-start')).toHaveLength(1);
    expect(screen.queryByTestId('original-saved-start-review')).toBeNull(); expect(native.query).not.toHaveBeenCalled();
    fireEvent.click(open); expect(await screen.findByTestId('original-saved-start-review')).toHaveTextContent('G');
    expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
    fireEvent.click(screen.getByText('Original Start query')); expect(native.query).toHaveBeenCalledOnce(); expect(native.mutation).not.toHaveBeenCalled();
  });
  it('shows a safe handoff instead of a no-op Start lookup for a custom saved destination', async () => {
    const read = startEntry('custom'); retainStart(read);
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: read });
    const view = mount({ ...repair, section: 'operation', operation: { kind: 'query_operation', owner: 'start', deviceId: 'D', operationId: 'O1' } }, read);
    expect(await screen.findByTestId('recording-repair-unavailable')).toHaveTextContent('start.recoveryUnavailable');
    expect(screen.queryByTestId('recording-repair-start')).toBeNull(); expect(screen.queryByTestId('original-saved-start-review')).toBeNull();
    fireEvent.click(screen.getByText('advanced')); expect(view.advanced).toHaveBeenCalledOnce();
    expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
  });
  it('does not reinterpret a legacy Start owner or arbitrary storage text as a saved-record owner', async () => {
    saveBasicRecordingIntent({ scope_key: 'W/dev-1/managed-recording', operation_id: 'O1', request: { request_id: 'request-1', workspace_id: 'W', expected_workspace_revision: 'W2',
      device_ids: ['dev-1'], groups: [{ group_id: 'G', expected_group_revision: 'G3', expected_connector_revision: 'C1' }] } });
    sessionStorage.setItem('wbv2.basic-recording.W%2Fsaved-record%2FG', 'succeeded');
    const view = mount({ ...repair, section: 'operation', operation: { kind: 'query_operation', owner: 'start', deviceId: 'dev-1', operationId: 'O1' } });
    expect(await screen.findByTestId('recording-repair-unavailable')).toBeVisible();
    expect(screen.queryByTestId('recording-repair-start')).toBeNull();
    fireEvent.click(screen.getByText('advanced')); expect(view.advanced).toHaveBeenCalledOnce();
    expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
  });
  it.each(['revision', 'operation'])('does not focus a repair with mismatching %s', async reason => {
    retainSchemaOperation('W', 'G', reason === 'operation' ? 'O2' : 'O1');
    mount({ ...repair, workspaceRevision: reason === 'revision' ? 'old' : 'W2', section: 'operation', operation: { kind: 'query_operation', owner: 'schema', deviceId: 'dev-1', operationId: 'O1' } });
    expect(await screen.findByTestId('recording-repair-unavailable')).toBeVisible();
    expect(screen.getByTestId('recording-preparation-lookup')).not.toHaveFocus();
    expect(native.query).not.toHaveBeenCalled(); expect(native.mutation).not.toHaveBeenCalled();
  });
  it('drops a late source repair read on a new scope rather than focusing the wrong editor', async () => {
    let finish!: (value: StudioLocatorEntryResult) => void;
    vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }))
      .mockResolvedValueOnce({ ok: true, value: entry('B') });
    const view = mount(); await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
    view.rerender(view.content(entry('B'), { ...repair, groupId: 'B' }));
    await waitFor(() => expect(screen.getByText('Sources B')).toHaveFocus());
    await act(async () => finish({ ok: true, value: entry() }));
    expect(screen.getByText('Sources B')).toHaveFocus(); expect(screen.queryByText('Sources G')).toBeNull();
    expect(native.mutation).not.toHaveBeenCalled();
  });
});
// The only mocked child is the presentation seam; Step4/WGS canonical read and lease are real.
