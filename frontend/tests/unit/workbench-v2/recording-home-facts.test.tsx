import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useRecordingHomeFacts } from '@/features/datalink/workbench-v2/recording-home/useRecordingHomeFacts';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { listRecordingTestWriteRecovery } from '@/services/recordingTestWriteRecovery';
import { group, scope } from './recording-home-guidance-fixtures';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: vi.fn() }));
function wrapper(client: QueryClient) { return function Wrapper({ children }: { children: ReactNode }) { return <QueryClientProvider client={client}>{children}</QueryClientProvider>; }; }
beforeEach(() => {
  vi.clearAllMocks();
  vi.mocked(listRecordingTestWriteRecovery).mockResolvedValue({ workspace_id: 'W1', complete: true, operations: [] });
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
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
