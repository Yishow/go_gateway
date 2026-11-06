import { configReadContext, configReadKey } from '@/hooks/datalink/recordingConfigRefresh';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { listRecordingTestWriteRecovery } from '@/services/recordingTestWriteRecovery';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { pointAPI } from '@/services/datalink';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
vi.mock('@/pages/datalink/workbench-v2/useStudioV2AutosaveState', () => ({ useStudioV2AutosaveState: vi.fn(() => { throw new Error('Home must not mount autosave owners'); }) }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', async original => ({
  ...await original<typeof import('@/services/studioV2Workspace')>(), studioV2WorkspaceAPI: { get: vi.fn() },
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), schemaApply: vi.fn(), remove: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: vi.fn() }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
const emptySources = (workspace_id: string) => ({ workspace_id, state: 'available' as const, source_digest: 'a'.repeat(64), devices: [], candidates: [], excluded: [] });
const workspace = { id: 'W', kind: 'single' as const, status: 'empty' as const, ordered_device_ids: [], created_at: '', updated_at: '' };
const clients: QueryClient[] = [];
beforeEach(() => {
  vi.clearAllMocks();
  vi.mocked(recordingHomeSourcesAPI.read).mockReset().mockImplementation(async id => emptySources(id));
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockReset().mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(listRecordingTestWriteRecovery).mockReset().mockResolvedValue({ workspace_id: 'W1', complete: true, operations: [] });
  vi.mocked(studioV2WorkspaceAPI.get).mockReset().mockResolvedValue(workspace);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockReset().mockResolvedValue({ workspace_id: 'W', workspace_revision: '', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockReset().mockResolvedValue([]); vi.mocked(studioV2RulesAPI.list).mockReset().mockResolvedValue([]);
  vi.mocked(studioV2MappingsAPI.list).mockReset().mockResolvedValue([]); vi.mocked(pointAPI.list).mockReset().mockResolvedValue([]);
});
afterEach(() => { cleanup(); clients.forEach(client => client.clear()); clients.length = 0; });
function mount() {
  const client = new QueryClient(); clients.push(client); const navigate = vi.fn();
  render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings&section=database', hash: '#saved' }} navigateTo={navigate} /></QueryClientProvider>);
  return navigate;
}
describe('actual optional home entry', () => {
  it('keeps the same home and filters mounted during a same-workspace background refresh', async () => {
    mount(); await screen.findByText('reason.missing_devices');
    const home = screen.getByTestId('recording-home');
    fireEvent.change(screen.getByLabelText('keyword'), { target: { value: 'keep my filter' } });
    let resolve!: (value: typeof workspace) => void;
    vi.mocked(studioV2WorkspaceAPI.get).mockImplementationOnce(() => new Promise(done => { resolve = done; }));
    let pending!: Promise<void>;
    act(() => { pending = clients[0].invalidateQueries({ queryKey: configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext()), exact: true }); });
    await waitFor(() => expect(resolve).toBeTypeOf('function'));
    expect(screen.getByTestId('recording-home')).toBe(home);
    expect(screen.getByLabelText('keyword')).toHaveValue('keep my filter');
    await act(async () => { resolve(workspace); await pending; });
    expect(screen.getByTestId('recording-home')).toBe(home);
    expect(screen.getByLabelText('keyword')).toHaveValue('keep my filter');
  });
  it('renders the initial empty database setup without Share and opens no autosave/mutation owners', async () => {
    const navigate = mount(); expect(screen.getByTestId('recording-home')).toBeInTheDocument();
    await screen.findByText('reason.missing_devices'); expect(useStudioV2AutosaveState).not.toHaveBeenCalled();
    for (const action of ['create', 'update', 'apply', 'schemaApply', 'remove'] as const) expect(studioV2WorkspaceWriteGroupsAPI[action]).not.toHaveBeenCalled();
    for (const read of [studioV2WorkspaceDevicesAPI.list, studioV2RulesAPI.list, studioV2MappingsAPI.list, pointAPI.list]) expect(read).not.toHaveBeenCalled();
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled(); expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
    fireEvent.click(screen.getByText('advanced')); expect(navigate).toHaveBeenLastCalledWith('/studio/v2?section=database&view=advanced#saved');
  });
  it('failed scope bootstrap is safe and preserves query/hash on recovery navigation', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValue(new Error('password=secret private.host')); const navigate = mount();
    await screen.findByRole('alert'); expect(screen.getByText('error')).toBeInTheDocument();
    expect(document.body.textContent).not.toContain('secret'); expect(screen.queryByText('reason.missing_devices')).not.toBeInTheDocument();
    fireEvent.click(screen.getByText('advanced')); expect(navigate).toHaveBeenCalledWith('/studio/v2?section=database&view=advanced#saved');
  });
  it('source read failure never becomes no points or an empty healthy workspace', async () => {
    vi.mocked(recordingHomeSourcesAPI.read).mockRejectedValue(new Error('private backend')); mount();
    await screen.findByText('reason.evidence_unavailable');
    await waitFor(() => expect(screen.getByText('error')).toBeInTheDocument());
    expect(screen.queryByText('reason.missing_devices')).not.toBeInTheDocument();
    expect(document.body.textContent).not.toContain('private backend');
  });
  it('missing DB revision is not silently coerced to the initial empty revision', async () => {
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', groups: [] } as unknown as Awaited<ReturnType<typeof studioV2WorkspaceWriteGroupsAPI.list>>);
    mount(); await screen.findByText('unknown'); expect(screen.queryByText('reason.missing_devices')).not.toBeInTheDocument();
  });
  it('late sources from the previous workspace cannot overwrite the newly selected workspace', async () => {
    let resolveOld!: (value: Awaited<ReturnType<typeof recordingHomeSourcesAPI.read>>) => void;
    vi.mocked(recordingHomeSourcesAPI.read).mockImplementationOnce(() => new Promise(resolve => { resolveOld = resolve; }));
    mount(); await waitFor(() => expect(recordingHomeSourcesAPI.read).toHaveBeenCalledTimes(1));
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, id: 'NEXT' });
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'NEXT', workspace_revision: '', groups: [] });
    await act(async () => { await Promise.all([
      clients[0].invalidateQueries({ queryKey: configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext()), exact: true }),
      clients[0].invalidateQueries({ queryKey: studioV2WorkspaceKeys.writeGroups(), exact: true }),
    ]); });
    await screen.findByText('reason.missing_devices');
    await act(async () => resolveOld(emptySources('W')));
    expect(screen.getByText('reason.missing_devices')).toBeInTheDocument();
    expect(screen.queryByText('unknown')).not.toBeInTheDocument();
    expect(clients[0].getQueryCache().getAll().some(query => query.queryKey.includes('recording-home-sources') && query.queryKey.includes('NEXT'))).toBe(true);
  });
  it('holds one stable reading frame until first-load sources settle, then reveals once', async () => {
    let resolveSources!: (value: Awaited<ReturnType<typeof recordingHomeSourcesAPI.read>>) => void;
    vi.mocked(recordingHomeSourcesAPI.read).mockImplementation(() => new Promise(resolve => { resolveSources = resolve; }));
    mount(); await waitFor(() => expect(recordingHomeSourcesAPI.read).toHaveBeenCalled());
    const home = screen.getByTestId('recording-home');
    expect(home).toHaveAttribute('data-home-settling', 'true'); expect(home).toHaveAttribute('aria-busy', 'true');
    expect(home.querySelector('.gw-home-settling-status')).toHaveTextContent('loading');
    await act(async () => resolveSources(emptySources('W')));
    expect(recordingHomeSourcesAPI.read).toHaveBeenCalledTimes(1);
    await waitFor(() => expect(screen.getByTestId('recording-home')).toHaveAttribute('data-home-settling', 'false'));
    expect(screen.getByTestId('recording-home')).not.toHaveAttribute('aria-busy');
    expect(screen.getByText('reason.missing_devices')).toBeInTheDocument();
  });
});
import { group, fixture, delivery } from './recording-home-guidance-fixtures';
import { rememberHomePresentation } from '@/features/datalink/workbench-v2/recording-home/recordingHomePresentation';
import { initialCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
it('opens a removed home item through original management, never a deleted editor, while alive remains canonical', async () => {
  const groups = [group('REMOVED', 'deleted'), group('ALIVE', 'ready')];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, id: 'W1' });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockImplementation(async id => {
    const found = groups.find(item => item.id === id);if (!found) throw new Error('fixture');return delivery(found);
  });
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async id => {
    const found = groups.find(item => item.id === id);if (!found) throw new Error('fixture');return {
      workspace_id: 'W1', workspace_revision: 'W2', group_id: id, group_revision: found.revision, applied_revision: found.applied_revision,
      status: found.status, connector_id: 'DB1', connector_revision: 'DB2', actions: [], history: { canonical: true, current_group_id: id, generations: [] },
    };
  });
  const navigate = mount(), table = await screen.findByTestId('home-management-table');
  await screen.findAllByText('result.lastCommit');
  fireEvent.click(within(within(table).getByText('Recording REMOVED').closest('tr')!).getByRole('button', { name: 'openRecord' }));
  expect(navigate).toHaveBeenLastCalledWith('/studio/v2?step=output&target=database&view=advanced#saved');
  fireEvent.click(within(within(table).getByText('Recording ALIVE').closest('tr')!).getByRole('button', { name: 'openRecord' }));
  expect(navigate.mock.calls.at(-1)?.[0]).toMatch(/^\/studio\/v2\?entry=locator&workspace_id=W1&step=4&group_id=ALIVE&group_revision=.+&view=recording#saved$/);
  for (const action of ['create', 'update', 'apply', 'schemaApply', 'remove'] as const) expect(studioV2WorkspaceWriteGroupsAPI[action]).not.toHaveBeenCalled();
  expect(useStudioV2AutosaveState).not.toHaveBeenCalled();
});
it('restored observed-only page does not auto scan all recovered groups', async () => {
 const groups = Array.from({length: 75}, (_,i)=>group(`G${String(i).padStart(3,'0')}`));
 const data = fixture([], false, groups);
 vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({...workspace,id:'W1'});
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({workspace_id:'W1',workspace_revision:'W2',groups});
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockImplementation(async id => (data.evidence.delivery![id] as { value: unknown }).value as never);
 vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async id => (data.evidence.lifecycle![id] as { value: unknown }).value as never);
 const client = new QueryClient(); clients.push(client);
 rememberHomePresentation(client,'W1',{state:{...initialCollectionState(),observedOnly:true},workspaceRevision:'W2',scrollTop:0,focusId:null,observedIds:groups.map(g=>g.id)});
 render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{kind:'absent'}} presentationLocation={{search:'?view=recordings',hash:''}} /></QueryClientProvider>);
 await waitFor(()=>expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(25));
 await act(async () => { await new Promise(resolve => setTimeout(resolve, 80)); });
 expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(25);
 fireEvent.click(screen.getByText('recheckVisible'));
 await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(50));
 await act(async () => { await new Promise(resolve => setTimeout(resolve, 80)); });
 expect(new Set(vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mock.calls.map(([id]) => id)).size).toBe(25);
 fireEvent.click(screen.getByText('next'));
 await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(75));
 await act(async () => { await new Promise(resolve => setTimeout(resolve, 80)); });
 const ids = new Set(vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mock.calls.map(([id]) => id));
 expect(ids.size).toBe(50); expect(ids.has('G025')).toBe(true); expect(ids.has('G050')).toBe(false);
});

it('actual home retains a SQL fact through source faults and failed refresh without claiming new configuration success', async () => {
  const saved = { ...group('G'), revision: 'C', applied_revision: 'B' };
  const view = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' });
  const fixtureData = fixture(['D1'], true, [saved]);
  const deliveryValue = fixtureData.evidence.delivery!.G;
  if (deliveryValue.state !== 'available') throw new Error('fixture');
  const response = { ...deliveryValue.value, stages: { ...deliveryValue.value.stages, blocked: 20 }, last_sql_committed_effect: { group_revision: 'B', connector_revision: 'DB2',
    record_id: 'R', effect_key: 'E', payload_digest: 'a'.repeat(64), committed_at: '2026-10-30T10:32:00Z' } };
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, id: 'W1', ordered_device_ids: ['D1'] });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ ...emptySources('W1'), devices: [{ id: 'D1', name: 'Device', protocol: 'modbus_tcp', availability: 'unknown' }] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups: [saved] });
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async () => (fixtureData.evidence.lifecycle!.G as { value: unknown }).value as never);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockResolvedValue(response as never);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: true, value: view });
  mount();
  await screen.findByText('result.sql.previous_revision_receipt');
  expect(screen.getByText('2026-10-30T10:32:00Z')).toBeInTheDocument();
  expect(screen.getAllByText('result.effective.requested_not_effective').length).toBeGreaterThan(0);
  expect(screen.getByText('result.owners.unknown')).toBeInTheDocument();
  expect(screen.queryByText('reason.delivery_attention')).not.toBeInTheDocument();
  view.observed_at = '2026-10-30T10:34:00Z'; view.groups[0].as_of = view.observed_at;
  view.groups[0].delivery = { oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 0, counter_replay_quarantines: null,
    backlog: [], test_write: null, current: { binding: view.groups[0].effective!.binding, totals_scope: 'group', shared: false,
      intake: { group_revision: 'A', state: 'active', reason: null }, stages: { ...response.stages, sql_committed: 0, blocked: 0 }, last_commit: null,
      sources: [{ device_id: 'D1', point_id: 'P', tag_id: 'T', last_committed: null, latest: { state: 'missing', window_start: '2026-10-30T10:32:30Z',
        window_end: '2026-10-30T10:33:00Z', delta_known: null, coverage: null, receipt: null } }], source_state: 'available',
      next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null } };
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: true, value: structuredClone(view) });
  await act(async () => { await clients[0].invalidateQueries({ predicate: query => query.queryKey.includes('home-rule-evidence') }); });
  await screen.findByText('result.current.source_attention');
  expect(screen.getByText('reason.source_attention')).toBeInTheDocument();
  const calls = vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mock.calls.length;
  fireEvent.change(screen.getByLabelText('keyword'), { target: { value: 'no-matching-record' } });
  await act(async () => { await new Promise(resolve => setTimeout(resolve, 80)); });
  expect(screen.getByRole('complementary', { name: 'observedAttentionTitle' })).toBeInTheDocument();
  expect(screen.getByText('reason.source_attention')).toBeInTheDocument();
  expect(studioV2RuntimeContextAPI.getRuleEvidence).toHaveBeenCalledTimes(calls);
  fireEvent.change(screen.getByLabelText('keyword'), { target: { value: '' } });
  await screen.findByText('result.current.source_attention');
  expect(screen.getByText('2026-10-30T10:32:00Z')).toBeInTheDocument();
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockRejectedValue(new Error('PRIVATE-ERROR'));
  await act(async () => { await clients[0].invalidateQueries({ predicate: query => query.queryKey.includes('delivery') }); });
  await screen.findByText('result.sql.historical_commit');
  expect(screen.getByText('2026-10-30T10:32:00Z')).toBeInTheDocument();
  expect(document.body.textContent).not.toContain('PRIVATE-ERROR');
  for (const action of ['create', 'update', 'apply', 'schemaApply', 'remove'] as const) expect(studioV2WorkspaceWriteGroupsAPI[action]).not.toHaveBeenCalled();
  expect(useStudioV2AutosaveState).not.toHaveBeenCalled();
});


it('warm home return keeps observed devices, filter and page during refresh and failure without fresh evidence', async () => {
  const devices = Array.from({ length: 30 }, (_, index) => ({ id: `D${String(index).padStart(2, '0')}`,
    name: `Observed device ${String(index).padStart(2, '0')}`, protocol: 'modbus_tcp', availability: 'unknown' }));
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, ordered_device_ids: devices.map(device => device.id) });
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ ...emptySources('W'), devices });
  const client = new QueryClient(); clients.push(client);
  const content = <QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings"
    locatorIntent={{ kind: 'absent' }} presentationLocation={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>;
  const first = render(content);
  await screen.findByRole('option', { name: 'Observed device 29' });
  fireEvent.click(screen.getByRole('button', { name: 'devicesTab' }));
  fireEvent.change(screen.getByLabelText('deviceKeyword'), { target: { value: 'Observed device' } });
  fireEvent.click(screen.getByRole('button', { name: 'next' }));
  expect(screen.getByText('Observed device 29')).toBeInTheDocument();
  expect(within(screen.getByTestId('home-management-table')).queryByText('Observed device 00')).not.toBeInTheDocument();
  first.unmount();
  let rejectRefresh!: (reason: Error) => void;
  vi.mocked(recordingHomeSourcesAPI.read).mockImplementationOnce(() => new Promise((_resolve, reject) => { rejectRefresh = reject; }));
  render(content);
  await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveTextContent('sources.refreshing'));
  expect(screen.getByLabelText('deviceKeyword')).toHaveValue('Observed device');
  expect(screen.getByText('Observed device 29')).toBeInTheDocument();
  expect(within(screen.getByTestId('home-management-table')).queryByText('Observed device 00')).not.toBeInTheDocument();
  expect(screen.queryByText('reason.missing_points')).not.toBeInTheDocument();
  await act(async () => rejectRefresh(new Error('PRIVATE-REFRESH-ERROR')));
  await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveTextContent('sources.error'));
  expect(screen.getByLabelText('deviceKeyword')).toHaveValue('Observed device');
  expect(screen.getByText('Observed device 29')).toBeInTheDocument();
  expect(within(screen.getByTestId('home-management-table')).queryByText('Observed device 00')).not.toBeInTheDocument();
  expect(screen.queryByText('reason.missing_points')).not.toBeInTheDocument();
  expect(document.body.textContent).not.toContain('PRIVATE-REFRESH-ERROR');
  expect(recordingHomeSourcesAPI.read).toHaveBeenCalledTimes(2);
  expect(useStudioV2AutosaveState).not.toHaveBeenCalled();
});


it.each(['API base URL', 'QueryClient'] as const)('isolates observed home and presentation when %s changes in place', async change => {
  const originalBaseURL = studioV2DatalinkApi.defaults.baseURL;
  const devices = Array.from({ length: 30 }, (_, index) => ({ id: `OLD${index}`,
    name: `Previous context ${String(index).padStart(2, '0')}`, protocol: 'modbus_tcp', availability: 'unknown' }));
  const scope = { ...workspace, ordered_device_ids: devices.map(device => device.id) };
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(scope);
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue({ ...emptySources('W'), devices });
  const firstClient = new QueryClient(), nextClient = new QueryClient(); clients.push(firstClient, nextClient);
  const content = (client: QueryClient) => <QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings"
    locatorIntent={{ kind: 'absent' }} presentationLocation={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>;
  const view = render(content(firstClient));
  try {
    await screen.findByRole('option', { name: 'Previous context 29' });
    fireEvent.click(screen.getByRole('button', { name: 'devicesTab' }));
    fireEvent.change(screen.getByLabelText('deviceKeyword'), { target: { value: 'Previous context' } });
    fireEvent.click(screen.getByRole('button', { name: 'next' }));
    expect(screen.getByText('Previous context 29')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'previous' })).toBeEnabled();
    let resolveNext!: (value: Awaited<ReturnType<typeof recordingHomeSourcesAPI.read>>) => void;
    vi.mocked(recordingHomeSourcesAPI.read).mockImplementationOnce(() => new Promise(resolve => { resolveNext = resolve; }));
    if (change === 'API base URL') studioV2DatalinkApi.defaults.baseURL = '/different-backend/api/v1/datalink';
    else {
      nextClient.setQueryData(configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext()), scope);
      nextClient.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'W', workspace_revision: '', groups: [] });
    }
    view.rerender(content(change === 'QueryClient' ? nextClient : firstClient));
    await waitFor(() => expect(resolveNext).toBeTypeOf('function'));
    expect(screen.getByTestId('home-source-observation')).toHaveTextContent('sources.loading');
    expect(screen.queryByText(/Previous context/)).not.toBeInTheDocument();
    expect(screen.getByLabelText('keyword')).toHaveValue('');
    expect(screen.getByRole('button', { name: 'recordingsTab' })).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByRole('button', { name: 'previous' })).toBeDisabled();
    await act(async () => resolveNext(emptySources('W')));
    await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveTextContent('sources.fresh'));
    expect(screen.queryByText(/Previous context/)).not.toBeInTheDocument();
    expect(screen.getByLabelText('keyword')).toHaveValue('');
    expect(recordingHomeSourcesAPI.read).toHaveBeenCalledTimes(2);
  } finally {
    view.unmount(); studioV2DatalinkApi.defaults.baseURL = originalBaseURL;
  }
});

it('revalidates 1000 groups with one summary GET and never replenishes the captured 25-detail demand', async () => {
  const groups = Array.from({ length: 1000 }, (_, index) => group(`G${String(index).padStart(4, '0')}`));
  const data = fixture([], false, groups);
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, id: 'W1' });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockImplementation(async id => (data.evidence.delivery![id] as { value: unknown }).value as never);
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async id => (data.evidence.lifecycle![id] as { value: unknown }).value as never);
  const actual = await vi.importActual<typeof import('@/services/recordingHomeSources')>('@/services/recordingHomeSources');
  vi.mocked(recordingHomeSourcesAPI.read).mockImplementation(actual.recordingHomeSourcesAPI.read);
  const summaryPath = '/studio-v2/recording-home/sources';
  const get = vi.spyOn(studioV2DatalinkApi, 'get').mockImplementation(async path => {
    if (path !== summaryPath) throw new Error('Unmocked readonly endpoint');
    return { data: { success: true, data: emptySources('W1') } };
  });
  const client = new QueryClient(); clients.push(client);
  rememberHomePresentation(client, 'W1', { state: { ...initialCollectionState(), observedOnly: true, keyword: 'Recording' },
    workspaceRevision: 'W2', scrollTop: 0, focusId: null, observedIds: groups.map(item => item.id) });
  try {
    render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings"
      locatorIntent={{ kind: 'absent' }} presentationLocation={{ search: '?view=recordings', hash: '' }} /></QueryClientProvider>);
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(25));
    await waitFor(() => expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(25));
    const home = screen.getByTestId('recording-home');
    const filter = screen.getByLabelText('keyword'); filter.focus();
    const panel = screen.getByTestId('home-management-table').parentElement!; panel.scrollTop = 120;
    for (let refresh = 0; refresh < 3; refresh++) {
      let resolveSummary!: (value: { data: { success: boolean; data: ReturnType<typeof emptySources> } }) => void;
      get.mockImplementationOnce(() => new Promise(resolve => { resolveSummary = resolve; }));
      let pending!: Promise<void>;
      act(() => { pending = client.invalidateQueries({ predicate: query => query.queryKey.includes('recording-home-sources') }); });
      await waitFor(() => expect(resolveSummary).toBeTypeOf('function'));
      expect(screen.getByTestId('home-source-observation')).toHaveTextContent('sources.refreshing');
      expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(25);
      await act(async () => { resolveSummary({ data: { success: true, data: emptySources('W1') } }); await pending; });
      await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveTextContent('sources.fresh'));
      await act(async () => { await new Promise(resolve => setTimeout(resolve, 80)); });
      expect(get.mock.calls.filter(([path]) => path === summaryPath)).toHaveLength(refresh + 2);
      expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(25);
      expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(25);
      expect(screen.getByTestId('recording-home')).toBe(home);
      expect(filter).toHaveValue('Recording'); expect(filter).toHaveFocus(); expect(panel.scrollTop).toBe(120);
      expect(screen.getByRole('button', { name: 'previous' })).toBeDisabled();
    }
    const ids = new Set(vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mock.calls.map(([id]) => id));
    expect(ids.size).toBe(25); expect(ids.has('G0024')).toBe(true); expect(ids.has('G0025')).toBe(false);
    for (const read of [studioV2WorkspaceDevicesAPI.list, studioV2RulesAPI.list, studioV2MappingsAPI.list, pointAPI.list]) expect(read).not.toHaveBeenCalled();
    for (const action of ['create', 'update', 'apply', 'schemaApply', 'remove'] as const) expect(studioV2WorkspaceWriteGroupsAPI[action]).not.toHaveBeenCalled();
    expect(useStudioV2AutosaveState).not.toHaveBeenCalled();
  } finally { cleanup(); get.mockRestore(); }
}, 15000);
