import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useRecordingHomeFacts } from '@/features/datalink/workbench-v2/recording-home/useRecordingHomeFacts';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { listRecordingTestWriteRecovery } from '@/services/recordingTestWriteRecovery';
import { group, scope } from './recording-home-guidance-fixtures';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { useCircuitProfileStart } from '@/features/datalink/workbench-v2/steps/step4/useCircuitProfileStart';
import { basicRecordingIntentKey, loadBasicRecordingIntent, saveBasicRecordingIntent, type BasicRecordingIntentRecord } from '@/features/datalink/workbench-v2/state/basicRecording';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture } from './recording-home-guidance-fixtures';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: vi.fn() }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), ensureBasic: vi.fn(), lookup: vi.fn(), operation: vi.fn() } }));
function wrapper(client: QueryClient) { return function Wrapper({ children }: { children: ReactNode }) { return <QueryClientProvider client={client}>{children}</QueryClientProvider>; }; }
beforeEach(() => {
  vi.clearAllMocks();
  vi.mocked(listRecordingTestWriteRecovery).mockResolvedValue({ workspace_id: 'W1', complete: true, operations: [] });
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
});

describe('Home after original Start recovery', () => {
  const clients: QueryClient[] = [], ownerKey = 'W1/saved-record/G';
  const saved = { ...group('G'), revision: 'G2', applied_revision: 'G2' };
  const original: BasicRecordingIntentRecord = { scope_key: ownerKey, owner: { kind: 'saved_record', groupId: 'G', deviceIds: ['D1'] },
    request: { request_id: 'original-request', workspace_id: 'W1', expected_workspace_revision: 'BEFORE-START', device_ids: ['D1'],
      groups: [{ group_id: 'G', expected_group_revision: 'G2', expected_connector_revision: 'DB2' }] } };
  const terminal: RecordingStartOperation = { operation_id: 'O1', action: 'recording_start', status: 'succeeded', stage: 'complete', intent_digest: 'a'.repeat(64),
    workspace_id: 'W1', setup_revision: 'W2', device_ids: ['D1'], devices: [{ device_id: 'D1', activated: true }],
    groups: [{ group_id: 'G', group_revision: 'G2', applied_revision: 'G2', saved: true, ready: true, applied: true }],
    created_at: '2026-11-01T18:42:00Z', updated_at: '2026-11-01T18:43:00Z' };
  const operationKey = [...studioV2WorkspaceKeys.all, 'recording-start', 'O1'];
  beforeEach(() => { sessionStorage.clear(); });
  afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); sessionStorage.clear(); });
  function setup() {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
    saveBasicRecordingIntent(original);
    vi.mocked(studioV2RecordingStartAPI.lookup).mockResolvedValue(terminal);
    vi.mocked(studioV2RecordingStartAPI.operation).mockResolvedValue(terminal);
    const owner = renderHook(() => useCircuitProfileStart({ scopeKey: ownerKey, canStart: false, readonly: true, workspaceId: 'W1', deviceId: 'D1',
      savedRecord: original.owner, connectorRevision: 'DB2', prepare: vi.fn(async () => undefined) }), { wrapper: wrapper(client) });
    return { client, owner };
  }
  async function recover(value: ReturnType<typeof setup>) {
    await act(async () => { await value.owner.result.current.check(); });
    await waitFor(() => expect(value.owner.result.current.completion).toMatchObject({ state: 'available', status: 'succeeded', operationId: 'O1' }));
    await waitFor(() => expect(value.client.getQueryState(operationKey)).toMatchObject({ status: 'success', fetchStatus: 'idle', data: terminal }));
    expect(studioV2RecordingStartAPI.lookup).toHaveBeenCalledWith(original.request);
    expect(loadBasicRecordingIntent(ownerKey)).toMatchObject({ request: original.request, operation_id: 'O1' });
    expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled(); expect(studioV2RecordingStartAPI.ensureBasic).not.toHaveBeenCalled();
  }
  function startObservation(value: ReturnType<typeof useRecordingHomeFacts>) {
    const evidence = value.owners.G;
    return evidence?.state === 'available' ? evidence.value.find(item => item.owner === 'start') : undefined;
  }
  it.each(['remounted', 'already mounted'])('withdraws only unfinished guidance after original Start lookup with Home %s', async mode => {
    const value = setup();
    const before = mode === 'already mounted' ? renderHook(() => useRecordingHomeFacts(scope, [saved], ['G']), { wrapper: wrapper(value.client) }) : undefined;
    if (before) expect(startObservation(before.result.current)).toMatchObject({ state: 'unresolved', reason: 'lost_ack' });
    await recover(value); value.owner.unmount();
    const home = before ?? renderHook(() => useRecordingHomeFacts(scope, [saved], ['G']), { wrapper: wrapper(value.client) });
    await waitFor(() => expect(startObservation(home.result.current)).toMatchObject({ state: 'observed', result: 'resolved', status: 'succeeded', operationId: 'O1' }));
    expect(recordingHomeGuidance({ ...fixture(['D1'], true, [saved]), owners: home.result.current.owners }).reason).toBe('evidence_unavailable');
    expect(home.result.current.runtimeGroups.G).toBeUndefined();
    expect(loadBasicRecordingIntent(ownerKey)).toMatchObject({ request: original.request, operation_id: 'O1' });
    expect(studioV2RecordingStartAPI.operation).toHaveBeenCalledTimes(1);
    expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled(); expect(studioV2RecordingStartAPI.ensureBasic).not.toHaveBeenCalled();
  });
  it.each(['foreign workspace', 'stale setup revision', 'stale group revision', 'newer original request'])('does not resolve retained Start from %s evidence', async mismatch => {
    const value = setup(); await recover(value); value.owner.unmount();
    if (mismatch === 'newer original request') saveBasicRecordingIntent({ ...original, operation_id: 'O2', request: { ...original.request, request_id: 'newer-request' } });
    else value.client.setQueryData(operationKey, mismatch === 'foreign workspace' ? { ...terminal, workspace_id: 'OTHER' }
      : mismatch === 'stale setup revision' ? { ...terminal, setup_revision: 'OLDER' }
        : { ...terminal, groups: [{ ...terminal.groups[0], group_revision: 'OLD', applied_revision: 'OLD' }] });
    const home = renderHook(() => useRecordingHomeFacts(scope, [saved], ['G']), { wrapper: wrapper(value.client) });
    await waitFor(() => expect(listRecordingTestWriteRecovery).toHaveBeenCalled());
    expect(startObservation(home.result.current)).toMatchObject({ state: 'unresolved' });
    expect(studioV2RecordingStartAPI.operation).toHaveBeenCalledTimes(1); expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled();
  });
  it('does not resolve retained Start from invalidated or missing cache and stored terminal words', async () => {
    const value = setup(); await recover(value); value.owner.unmount();
    await value.client.invalidateQueries({ queryKey: operationKey, refetchType: 'none' });
    const home = renderHook(() => useRecordingHomeFacts(scope, [saved], ['G']), { wrapper: wrapper(value.client) });
    await waitFor(() => expect(listRecordingTestWriteRecovery).toHaveBeenCalled());
    expect(startObservation(home.result.current)).toMatchObject({ state: 'unresolved' });
    const key = basicRecordingIntentKey(ownerKey), retained = sessionStorage.getItem(key)!;
    sessionStorage.setItem(key, JSON.stringify({ ...JSON.parse(retained), status: 'succeeded', stage: 'complete' }));
    act(() => value.client.removeQueries({ queryKey: operationKey })); home.rerender();
    expect(startObservation(home.result.current)).toMatchObject({ state: 'unresolved' });
    expect(studioV2RecordingStartAPI.operation).toHaveBeenCalledTimes(1); expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled();
  });
  it('withdraws completed Start observation when its original query is fetching, failed or removed', async () => {
    const value = setup(); await recover(value); value.owner.unmount();
    const home = renderHook(() => useRecordingHomeFacts(scope, [saved], ['G']), { wrapper: wrapper(value.client) });
    await waitFor(() => expect(startObservation(home.result.current)).toMatchObject({ state: 'observed', result: 'resolved' }));
    let reject!: (reason: Error) => void;
    let pending!: Promise<unknown>;
    act(() => { pending = value.client.fetchQuery({ queryKey: operationKey, queryFn: () => new Promise<RecordingStartOperation>((_, fail) => { reject = fail; }) }).catch(() => undefined); });
    await waitFor(() => expect(startObservation(home.result.current)).toMatchObject({ state: 'unresolved' }));
    await act(async () => { reject(new Error('read unavailable')); await pending; });
    expect(value.client.getQueryState(operationKey)?.status).toBe('error');
    expect(startObservation(home.result.current)).toMatchObject({ state: 'unresolved' });
    act(() => value.client.setQueryData(operationKey, terminal));
    await waitFor(() => expect(startObservation(home.result.current)).toMatchObject({ state: 'observed', result: 'resolved' }));
    act(() => value.client.removeQueries({ queryKey: operationKey }));
    await waitFor(() => expect(startObservation(home.result.current)).toMatchObject({ state: 'unresolved' }));
    expect(loadBasicRecordingIntent(ownerKey)).toMatchObject({ request: original.request, operation_id: 'O1' });
    expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled(); expect(studioV2RecordingStartAPI.ensureBasic).not.toHaveBeenCalled();
  });
});
describe('home facts actual readonly service adapter', () => {
  it('only fixed page IDs read runtime; 300 groups require25 device reads plus one recovery read', async () => {
    const groups = Array.from({ length: 300 }, (_, i) => ({ ...group(`G${i}`), members: [{ ...group('G').members[0], device_id: `D${i % 100}` }] }));
    const client = new QueryClient();
    const { unmount, rerender } = renderHook(ids => useRecordingHomeFacts(scope, groups, ids), { initialProps: groups.slice(0, 25).map(g => g.id), wrapper: wrapper(client) });
    await waitFor(() => expect(studioV2RuntimeContextAPI.getRuleEvidence).toHaveBeenCalledTimes(25));
    expect(listRecordingTestWriteRecovery).toHaveBeenCalledTimes(1);
    expect(listRecordingTestWriteRecovery).toHaveBeenCalledWith('W1', { workspaceId: 'W1', workspaceRevision: 'W2' });
    rerender([]); await act(async () => { await Promise.resolve(); });
    expect(studioV2RuntimeContextAPI.getRuleEvidence).toHaveBeenCalledTimes(25);
    unmount(); client.clear();
  });
  it('uses one member to select complete group evidence, including the second member and shared totals', async () => {
    const saved = { ...group('G'), revision: 'C', applied_revision: 'B', members: [group('G').members[0], { ...group('G').members[0], device_id: 'D2' }] };
    const view = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' });
    view.groups[0].effective!.members.push({ ...view.groups[0].effective!.members[0], device_id: 'D2' });
    vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: true, value: view });
    const client = new QueryClient(); const { result, unmount } = renderHook(() => useRecordingHomeFacts(scope, [saved], ['G']), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.runtimeGroups.G?.effective_revision).toBe('A'));
    expect(result.current.runtimeGroups.G?.effective?.members.map(m => m.device_id)).toEqual(['D1', 'D2']);
    expect(studioV2RuntimeContextAPI.getRuleEvidence).toHaveBeenCalledTimes(1);
    expect(result.current.owners.G.state === 'available' && result.current.owners.G.value.filter(o => o.owner === 'start').every(o => o.state === 'unavailable')).toBe(true);
    unmount(); client.clear();
  });
  it('late previous context cannot publish runtime facts or recovery ownership', async () => {
    let finish!: (value: Awaited<ReturnType<typeof studioV2RuntimeContextAPI.getRuleEvidence>>) => void;
    vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
    const saved = { ...group('G'), revision: 'C', applied_revision: 'B' }; const client = new QueryClient();
    const { result, rerender, unmount } = renderHook(selected => useRecordingHomeFacts(selected, [saved], ['G']), { initialProps: scope, wrapper: wrapper(client) });
    await waitFor(() => expect(studioV2RuntimeContextAPI.getRuleEvidence).toHaveBeenCalledTimes(1));
    rerender({ ...scope, epoch: 99 });
    await act(async () => finish({ ok: true, value: ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' }) }));
    expect(result.current.runtimeGroups.G).toBeUndefined(); unmount(); client.clear();
  });
});
it('retains only historical SQL metadata when a formerly valid runtime read becomes unavailable', async () => {
  const saved = { ...group('G'), revision: 'C', applied_revision: 'B' };
  const view = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' });
  view.observed_at = '2026-10-30T10:34:00Z'; view.groups[0].as_of = view.observed_at;
  view.groups[0].delivery = { oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 0, counter_replay_quarantines: null,
    backlog: [], test_write: null, current: { binding: view.groups[0].effective!.binding, totals_scope: 'group', shared: true,
      intake: { group_revision: 'A', state: 'active', reason: null }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, skipped: 0, sql_committed: 1 },
      last_commit: { group_revision: 'A', connector_revision: '7', record_id: 'R', effect_key: 'E', payload_digest: 'a'.repeat(64), committed_at: '2026-10-30T10:32:00Z', provenance: 'runtime' },
      sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null } };
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: true, value: view });
  const client = new QueryClient(); const { result, unmount, rerender } = renderHook(selected => useRecordingHomeFacts(selected, [saved], ['G']), { initialProps: scope, wrapper: wrapper(client) });
  await waitFor(() => expect(result.current.histories.G?.committedAt).toBe('2026-10-30T10:32:00Z'));
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  await act(async () => { await result.current.retry(); });
  await waitFor(() => expect(result.current.runtimeGroups.G).toBeUndefined());
  expect(result.current.histories.G).toEqual({ workspaceId: 'W1', groupId: 'G', committedAt: '2026-10-30T10:32:00Z', destination: 'meter_old' });
  expect(JSON.stringify(result.current.histories)).not.toMatch(/payload|effect_key|record_id/);
  rerender({ ...scope, epoch: 999 });
  expect(result.current.histories.G).toBeUndefined(); unmount(); client.clear();
});
it('off-page exact cache retains the observation without fetching, while invalidation/revision removal withdraw it', async () => {
  const saved = { ...group('G'), revision: 'C', applied_revision: 'B' };
  const view = ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' });
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: true, value: view });
  const client = new QueryClient();
  const { result, rerender, unmount } = renderHook(props => useRecordingHomeFacts(scope, props.groups, props.ids), {
    initialProps: { groups: [saved], ids: ['G'] }, wrapper: wrapper(client),
  });
  await waitFor(() => expect(result.current.runtimeGroups.G?.effective_revision).toBe('A'));
  const reads = vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mock.calls.length;
  rerender({ groups: [saved], ids: [] });
  expect(result.current.runtimeGroups.G?.effective_revision).toBe('A');
  expect(studioV2RuntimeContextAPI.getRuleEvidence).toHaveBeenCalledTimes(reads);
  await act(async () => { await client.invalidateQueries({ predicate: query => query.queryKey.includes('home-rule-evidence'), refetchType: 'none' }); });
  await waitFor(() => expect(result.current.runtimeGroups.G).toBeUndefined());
  expect(studioV2RuntimeContextAPI.getRuleEvidence).toHaveBeenCalledTimes(reads);
  rerender({ groups: [{ ...saved, revision: 'NEW' }], ids: [] });
  expect(result.current.runtimeGroups.G).toBeUndefined();
  rerender({ groups: [], ids: [] }); expect(result.current.runtimeGroups).toEqual({});
  unmount(); client.clear();
});
