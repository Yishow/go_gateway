import { recordingReadKey } from '@/hooks/datalink/studioV2RecordingReadScope';
import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useRecordingHomeEvidence } from '@/features/datalink/workbench-v2/recording-home/useRecordingHomeEvidence';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { fixture, group, scope } from './recording-home-guidance-fixtures';
import type { RecordingHomeScope } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), delivery: vi.fn(), create: vi.fn(), apply: vi.fn(), schemaApply: vi.fn(), update: vi.fn(), remove: vi.fn(),
} }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
function setup(selected: RecordingHomeScope = scope, groups = [group('G')]) {
  const data = fixture(['D1'], true, groups);
  if (data.evidence.workspace?.state !== 'available' || data.evidence.groups?.state !== 'available') throw new Error('fixture');
  const workspace = { ...data.evidence.workspace.value, id: selected.workspaceId,
    modbus_share: { hydration_state: 'ready', readiness: false, workspace_id: selected.workspaceId, workspace_revision: 'share-source-7' } };
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(structuredClone(workspace) as Awaited<ReturnType<typeof studioV2WorkspaceAPI.get>>);
  api.list.mockResolvedValue({ ...data.evidence.groups.value, workspace_id: selected.workspaceId, workspace_revision: selected.workspaceRevision, groups });
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async id => {
    const value = data.evidence.lifecycle?.[id]; if (value?.state !== 'available') throw new Error('fixture'); return structuredClone(value.value) as Awaited<ReturnType<typeof studioV2WriteGroupLifecycleAPI.lifecycle>>;
  });
  api.delivery.mockImplementation(async id => {
    const value = data.evidence.delivery?.[id]; if (value?.state !== 'available') throw new Error('fixture'); return structuredClone(value.value) as Awaited<ReturnType<typeof api.delivery>>;
  });
  return { scope: selected, bootstrapRevision: 'share-source-7', sources: data.evidence.sources, owners: data.owners };
}
function wrapper(client: QueryClient) { return function Wrapper({ children }: { children: React.ReactNode }) {
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}; }
beforeEach(() => { vi.resetAllMocks(); });
describe('recording home readonly query adapter', () => {
  it('reads actual canonical services, not mutations or a separate workspace store', async () => {
    const input = setup(); const client = new QueryClient();
    const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    expect(result.current.guidance.availability.state).toBe('unavailable');
    await waitFor(() => expect(result.current.guidance.reason).toBe('recordings'));
    expect(api.list).toHaveBeenCalledTimes(1); expect(api.delivery).toHaveBeenCalledWith('G', expect.any(AbortSignal));
    for (const action of ['create', 'apply', 'schemaApply', 'update', 'remove'] as const) expect(api[action]).not.toHaveBeenCalled();
    expect(result.current.error).toBeNull(); client.clear();
  });
  it('does not expose raw bootstrap errors or portray failure as empty', async () => {
    const input = setup(); vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValue(new Error('password=secret host=private'));
    const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.error).toEqual({ code: 'recording_home_read_failed', action: 'retry' }));
    expect(JSON.stringify(result.current)).not.toContain('secret');
    expect(result.current.guidance.reason).toBe('evidence_unavailable'); client.clear();
  });
  it('missing bootstrap Share revision remains unknown, never a fabricated empty workspace', async () => {
    const input = setup(scope, []); vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'empty', ordered_device_ids: [], created_at: '', updated_at: '' });
    const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.projection.workspace).toEqual({ state: 'unavailable', reason: 'unknown' }));
    expect(result.current.guidance.reason).toBe('evidence_unavailable'); client.clear();
  });
  it.each([{ workspaceId: '' }, { workspaceRevision: undefined as unknown as string }, { epoch: -1 }])('invalid scope %j starts no reads', patch => {
    const input = setup(); const client = new QueryClient();
    const { result } = renderHook(() => useRecordingHomeEvidence({ ...input, scope: { ...scope, ...patch } }), { wrapper: wrapper(client) });
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled(); expect(api.list).not.toHaveBeenCalled();
    expect(result.current.guidance.primaryAction).toBeNull(); client.clear();
  });
  it('legacy unscoped cache cannot become fresh scoped evidence', () => {
    const input = setup(); const client = new QueryClient();
    client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { id: 'W1', ordered_device_ids: [] });
    client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'W1', workspace_revision: 'W2', groups: [] });
    const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    expect(result.current.projection.workspace.state).toBe('unavailable');
    expect(result.current.guidance.reason).toBe('evidence_unavailable'); client.clear();
  });
  it('late A completion after B selection cannot overwrite B or enable commands', async () => {
    const input = setup(scope, []); let resolveA!: (value: Awaited<ReturnType<typeof api.list>>) => void;
    api.list.mockImplementationOnce(() => new Promise(resolve => { resolveA = resolve; }));
    const client = new QueryClient(); const { result, rerender } = renderHook(props => useRecordingHomeEvidence(props), { initialProps: input, wrapper: wrapper(client) });
    const bScope = { workspaceId: 'B', workspaceRevision: 'B2', epoch: 3 }; const b = setup(bScope, []);
    rerender({ ...b, sources: { state: 'available', scope: bScope, value: { candidates: [], excluded: [] } } });
    await waitFor(() => expect(result.current.projection.groups.state).toBe('available'));
    await act(async () => resolveA({ workspace_id: 'W1', workspace_revision: 'W2', groups: [] }));
    expect(result.current.projection.scope).toEqual(bScope);
    expect(result.current.guidance.primaryAction?.scope).toEqual(bScope); client.clear();
  });
  it('refetch error invalidates old success instead of keeping green evidence', async () => {
    const input = setup(); const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.guidance.reason).toBe('recordings'));
    api.list.mockRejectedValue(new Error('secret'));
    await act(async () => { await result.current.retry(); });
    await waitFor(() => expect(result.current.projection.groups).toEqual({ state: 'unavailable', reason: 'error' }));
    expect(result.current.guidance.reason).toBe('evidence_unavailable'); client.clear();
  });
  it('leave and return to the same workspace rejects the original request epoch', async () => {
    const input = setup(scope, []); let resolveOld!: (value: Awaited<ReturnType<typeof api.list>>) => void;
    api.list.mockImplementationOnce(() => new Promise(resolve => { resolveOld = resolve; }));
    const client = new QueryClient(); const { result, rerender } = renderHook(props => useRecordingHomeEvidence(props), { initialProps: input, wrapper: wrapper(client) });
    const b = setup({ workspaceId: 'B', workspaceRevision: 'B2', epoch: 3 }, []); rerender(b);
    await waitFor(() => expect(result.current.projection.groups.state).toBe('available'));
    const returned = setup({ ...scope, epoch: 4 }, []);
    returned.sources = { state: 'available', scope: returned.scope, value: { candidates: [], excluded: [] } };
    rerender(returned);
    await waitFor(() => expect(result.current.projection.groups.state).toBe('available'));
    await act(async () => resolveOld({ workspace_id: 'W1', workspace_revision: 'W2', groups: [group('OLD')] }));
    expect(result.current.projection.groups).toEqual({ state: 'available', value: [] });
    expect(api.list).toHaveBeenCalledTimes(3); client.clear();
  });
  it('late delivery for a previous applied revision cannot bless the new head', async () => {
    const old = group('G'); const input = setup(); let resolveOld!: (value: Awaited<ReturnType<typeof api.delivery>>) => void;
    api.delivery.mockImplementationOnce(() => new Promise(resolve => { resolveOld = resolve; }));
    const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(api.delivery).toHaveBeenCalledTimes(1));
    const updated = { ...old, revision: 'G3', applied_revision: 'G2' };
    const fresh = fixture(['D1'], true, [updated]);
    api.list.mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups: [updated] });
    api.delivery.mockImplementation(async () => { throw new Error('new head unavailable'); });
    vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async () => {
      const item = fresh.evidence.lifecycle?.G;
      if (item?.state !== 'available') throw new Error('fixture');
      return structuredClone(item.value) as Awaited<ReturnType<typeof studioV2WriteGroupLifecycleAPI.lifecycle>>;
    });
    await act(async () => { await client.invalidateQueries({ queryKey: recordingReadKey(studioV2WorkspaceKeys.writeGroups(), scope), exact: true }); });
    await waitFor(() => expect(api.delivery).toHaveBeenCalledTimes(2));
    const oldValue = fixture(['D1'], true, [old]).evidence.delivery?.G;
    if (oldValue?.state !== 'available') throw new Error('fixture');
    await act(async () => resolveOld(structuredClone(oldValue.value) as Awaited<ReturnType<typeof api.delivery>>));
    await waitFor(() => expect(result.current.projection.groups.state === 'available' && result.current.projection.groups.value[0].group.revision).toBe('G3'));
    const projected = result.current.projection.groups;
    expect(projected.state === 'available' && projected.value[0].delivery).toEqual({ state: 'unavailable', reason: 'error' }); client.clear();
  });
  it('does not relabel old source evidence after revision change', async () => {
    const input = setup({ ...scope, workspaceRevision: 'NEW', epoch: 3 }, []);
    const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.projection.groups.state).toBe('available'));
    expect(result.current.projection.sources).toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
    expect(result.current.guidance.reason).toBe('evidence_unavailable'); client.clear();
  });
  it('disabled readers do not promote cached data or refetch', async () => {
    const input = setup(); const client = new QueryClient();
    const { result, rerender } = renderHook(props => useRecordingHomeEvidence(props), { initialProps: { ...input, enabled: true }, wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.guidance.reason).toBe('recordings'));
    const calls = api.list.mock.calls.length;
    rerender({ ...input, enabled: false });
    expect(result.current.guidance.reason).toBe('evidence_unavailable');
    await act(async () => { await result.current.retry(); });
    expect(api.list).toHaveBeenCalledTimes(calls); client.clear();
  });
  it('explicit retry recovers after failure without submitting any operation', async () => {
    const input = setup(); vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValueOnce(new Error('failed'));
    const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.error).not.toBeNull());
    await act(async () => { await result.current.retry(); });
    await waitFor(() => expect(result.current.guidance.reason).toBe('recordings'));
    expect(result.current.error).toBeNull(); expect(api.apply).not.toHaveBeenCalled(); client.clear();
  });

  it('keeps the real backend Share and database setup revisions independent', async () => {
    const input = setup(); const client = new QueryClient();
    const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.guidance.reason).toBe('recordings'));
    expect(input.bootstrapRevision).not.toBe(input.scope.workspaceRevision);
    expect(result.current.projection.workspace.state).toBe('available');
    expect(result.current.projection.groups.state).toBe('available'); client.clear();
  });
  it('a new Share revision does not accept a late old bootstrap even when database setup is unchanged', async () => {
    const input = setup(); let resolveOld!: (value: Awaited<ReturnType<typeof studioV2WorkspaceAPI.get>>) => void;
    vi.mocked(studioV2WorkspaceAPI.get).mockImplementationOnce(() => new Promise(resolve => { resolveOld = resolve; }));
    const client = new QueryClient(); const { result, rerender } = renderHook(props => useRecordingHomeEvidence(props), { initialProps: input, wrapper: wrapper(client) });
    const fresh = { id: 'W1', kind: 'single' as const, status: 'ready' as const, ordered_device_ids: ['D1'], created_at: '', updated_at: '',
      modbus_share: { hydration_state: 'ready', readiness: false, workspace_id: 'W1', workspace_revision: 'share-source-8' } };
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(fresh);
    rerender({ ...input, bootstrapRevision: 'share-source-8' });
    await waitFor(() => expect(result.current.projection.workspace.state).toBe('available'));
    await act(async () => resolveOld({ ...fresh, modbus_share: { ...fresh.modbus_share, workspace_revision: 'share-source-7' } }));
    const projected = result.current.projection.workspace;
    expect(projected.state === 'available' && projected.value.modbus_share?.workspace_revision).toBe('share-source-8'); client.clear();
  });
  it('stale Share revision does not borrow the matching database setup revision', async () => {
    const input = setup(); const client = new QueryClient();
    const { result } = renderHook(() => useRecordingHomeEvidence({ ...input, bootstrapRevision: 'share-source-8' }), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.projection.workspace).toEqual({ state: 'unavailable', reason: 'unknown' }));
    expect(result.current.guidance.reason).toBe('evidence_unavailable'); client.clear();
  });
  it('stale database setup is rejected while the independent Share revision still matches', async () => {
    const input = setup(); api.list.mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'old-database-setup', groups: [] });
    const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.projection.groups).toEqual({ state: 'unavailable', reason: 'revision_mismatch' }));
    expect(result.current.projection.workspace.state).toBe('available'); client.clear();
  });
  it('bootstrap without Share is still an identity read, not a requirement to configure Share', async () => {
    const input = setup(); vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'ready', ordered_device_ids: ['D1'], created_at: '', updated_at: '' });
    const client = new QueryClient();
    const { result } = renderHook(() => useRecordingHomeEvidence({ ...input, bootstrapRevision: undefined }), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.guidance.reason).toBe('recordings'));
    expect(result.current.projection.workspace.state).toBe('available'); client.clear();
  });
  it('newRecord with explicitly observed empty database revision can guide the first device', async () => {
    // Backend: service_database_setup.go newRecord omits DatabaseSetupRevision;
    // write_group_service.go List emits the explicit empty workspace_revision.
    const selected = { workspaceId: 'W1', workspaceRevision: '', epoch: 1 };
    const input = setup(selected, []);
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W1', kind: 'single', status: 'empty', ordered_device_ids: [], created_at: '', updated_at: '' });
    const client = new QueryClient();
    const { result } = renderHook(() => useRecordingHomeEvidence({ ...input, bootstrapRevision: undefined,
      sources: { state: 'available', scope: selected, value: { candidates: [], excluded: [] } } }), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.guidance.reason).toBe('missing_devices'));
    expect(result.current.guidance.primaryAction).toMatchObject({ target: 'connect_device', scope: selected }); client.clear();
  });
  it('missing returned database revision is not repaired into a legal empty revision', async () => {
    const input = setup({ ...scope, workspaceRevision: '' }, []);
    api.list.mockResolvedValue({ workspace_id: 'W1', groups: [] } as unknown as Awaited<ReturnType<typeof api.list>>);
    const client = new QueryClient(); const { result } = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(result.current.projection.groups).toEqual({ state: 'unavailable', reason: 'revision_mismatch' }));
    expect(result.current.guidance.reason).toBe('evidence_unavailable'); client.clear();
  });

});

describe('recording home bounded page demand', () => {
  it.each([25, 50])('300 saved groups request only the visible %s, never every group', async size => {
    const groups = Array.from({ length: 300 }, (_, index) => group(`G${index}`));
    const input = setup(scope, groups); const client = new QueryClient();
    const { unmount } = renderHook(() => useRecordingHomeEvidence({ ...input,
      detailGroupIds: groups.slice(0, size).map(value => value.id) }), { wrapper: wrapper(client) });
    await waitFor(() => expect(api.delivery).toHaveBeenCalledTimes(size));
    expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(size);
    expect(api.delivery.mock.calls.every(([id]) => groups.slice(0, size).some(value => value.id === id))).toBe(true);
    unmount(); client.clear();
  });
  it('caps even malformed oversized demand at 50 canonical IDs', async () => {
    const groups = Array.from({ length: 300 }, (_, index) => group(`G${index}`));
    const input = setup(scope, groups); const client = new QueryClient();
    const { unmount } = renderHook(() => useRecordingHomeEvidence({ ...input, detailGroupIds: groups.map(value => value.id) }), { wrapper: wrapper(client) });
    await waitFor(() => expect(api.delivery).toHaveBeenCalledTimes(50));
    expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(50);
    unmount(); client.clear();
  });
  it('rapid paging cancels queued old-page reads and caps actual active transport at six', async () => {
    const groups = Array.from({ length: 100 }, (_, index) => group(`G${index}`));
    const input = setup(scope, groups); const client = new QueryClient();
    let active = 0; let peak = 0;
    const lifecycle = vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle);
    const actual = lifecycle.getMockImplementation()!;
    lifecycle.mockImplementation((id, signal) => new Promise((resolve, reject) => {
      active++; peak = Math.max(peak, active);
      const cancel = () => { active--; reject(new DOMException('cancelled', 'AbortError')); };
      signal?.addEventListener('abort', cancel, { once: true });
      if (Number(id.slice(1)) >= 25) {
        signal?.removeEventListener('abort', cancel); active--; void actual(id, signal).then(resolve, reject);
      }
    }));
    const { rerender, unmount } = renderHook(props => useRecordingHomeEvidence(props), {
      initialProps: { ...input, detailGroupIds: groups.slice(0, 25).map(value => value.id) }, wrapper: wrapper(client),
    });
    await waitFor(() => expect(lifecycle).toHaveBeenCalledTimes(6));
    rerender({ ...input, detailGroupIds: groups.slice(25, 50).map(value => value.id) });
    await waitFor(() => expect(api.delivery).toHaveBeenCalledTimes(25));
    expect(lifecycle.mock.calls.filter(([id]) => Number(id.slice(1)) < 25)).toHaveLength(6);
    expect(api.delivery.mock.calls.every(([id]) => Number(id.slice(1)) >= 25)).toBe(true);
    expect(peak).toBeLessThanOrEqual(6); expect(active).toBe(0);
    unmount(); client.clear();
  });
  it('unmounting one shared observer does not abort the other observer transport', async () => {
    const input = setup(); const client = new QueryClient(); let captured: AbortSignal | undefined;
    let finish!: () => void;
    const actual = api.delivery.getMockImplementation()!;
    api.delivery.mockImplementation((id, signal) => new Promise((resolve, reject) => {
      captured = signal; finish = () => { void actual(id, signal).then(resolve, reject); };
    }));
    const first = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    const second = renderHook(() => useRecordingHomeEvidence(input), { wrapper: wrapper(client) });
    await waitFor(() => expect(api.delivery).toHaveBeenCalledTimes(1));
    first.unmount(); expect(captured?.aborted).toBe(false);
    await act(async () => finish());
    await waitFor(() => expect(second.result.current.guidance.reason).toBe('recordings'));
    second.unmount(); client.clear();
  });
});
