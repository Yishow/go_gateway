import { configReadContext, refreshConfigReads, useRecordingConfigRefresh } from '@/hooks/datalink/recordingConfigRefresh';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, expect, it, vi } from 'vitest';
import { useStudioV2DatabaseConfigQuery, useStudioV2DatabaseTargetsQuery, useUpdateStudioV2DatabaseConfigMutation } from '@/hooks/datalink/useStudioV2WorkspaceDatabase';
import { useStudioV2WorkspaceQuery } from '@/hooks/datalink/useStudioV2Workspace';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import type { StudioV2WorkspaceDatabaseConfigRecord, StudioV2Workspace } from '@/types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '@/types/studioV2RuntimeApply';

vi.mock('@/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: { getConfig: vi.fn(), listTargets: vi.fn(), updateConfig: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', async importOriginal => ({ ...await importOriginal<typeof import('@/services/studioV2Workspace')>(), studioV2WorkspaceAPI: { get: vi.fn() } }));
export const savedConfig = (patch: Partial<StudioV2WorkspaceDatabaseConfigRecord> = {}): StudioV2RuntimeAppliedRecord<StudioV2WorkspaceDatabaseConfigRecord> => ({
  id: 'db-1', workspace_id: 'ws-1', identity_revision: 'identity-2', setup_revision: 'setup-2', kind: 'sqlite', name: 'Saved', host: '', port: 0,
  database: 'gateway.db', username: '', schema: '', table: 'values', write_mode: 'insert', write_interval_seconds: 5, timestamp_column: 'ts',
  status: 'ready', created_at: '', updated_at: '2', runtime_apply_status: 'apply_failed', ...patch,
});
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(yes => { resolve = yes; }); return { promise, resolve }; }
const workspace = { id: 'ws-1' } as StudioV2Workspace;
beforeEach(() => {
  vi.resetAllMocks();
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(savedConfig());
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockResolvedValue(savedConfig());
});
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: 1 } } });
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const hook = renderHook(() => {
    const config = useStudioV2DatabaseConfigQuery(true), targets = useStudioV2DatabaseTargetsQuery(true), workspace = useStudioV2WorkspaceQuery();
    return { config, targets, workspace, fetching: config.isFetching && targets.isFetching && workspace.isFetching,
      refresh: useRecordingConfigRefresh('ws-1'), mutation: useUpdateStudioV2DatabaseConfigMutation() };
  }, { wrapper });
  return { ...hook, client };
}
it('settles valid configuration ACK while all three independent refresh reads remain pending', async () => {
  const { result } = mount();
  await waitFor(() => expect(result.current.workspace.isSuccess && result.current.config.isSuccess && result.current.targets.isSuccess).toBe(true));
  const config = deferred<ReturnType<typeof savedConfig>>(), targets = deferred<[]>(), bootstrap = deferred<StudioV2Workspace>();
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockReturnValue(config.promise);
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockReturnValue(targets.promise);
  vi.mocked(studioV2WorkspaceAPI.get).mockReturnValue(bootstrap.promise);
  let completed = false;
  act(() => { void result.current.mutation.mutateAsync(savedConfig()).then(() => { completed = true; }); });
  await waitFor(() => expect(completed).toBe(true));
  expect(studioV2WorkspaceDatabaseAPI.getConfig).toHaveBeenCalledTimes(2);
  expect(studioV2WorkspaceDatabaseAPI.listTargets).toHaveBeenCalledTimes(2);
  expect(studioV2WorkspaceAPI.get).toHaveBeenCalledTimes(2);
  await waitFor(() => expect(result.current.fetching).toBe(true));
  expect(result.current.mutation.data?.runtime_apply_status).toBe('apply_failed');
  await act(async () => { config.resolve(savedConfig()); targets.resolve([]); bootstrap.resolve(workspace); });
});
it.each([{ identity_revision: undefined }, { setup_revision: undefined }, { workspace_id: '' }])('does not acknowledge HTTP success without confirmation evidence: %j', async patch => {
  const { result } = mount();
  await waitFor(() => expect(result.current.config.isSuccess).toBe(true));
  vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockResolvedValue(savedConfig(patch));
  const outcome = await act(async () => result.current.mutation.mutateAsync(savedConfig()).then(() => 'saved', () => 'unknown'));
  expect(outcome).toBe('unknown');
  expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1);
});

it('coalesces refresh retry and exposes failure without failing the confirmed save or replaying PUT', async () => {
  const { result, client } = mount();
  await waitFor(() => expect(result.current.workspace.isSuccess && result.current.targets.isSuccess).toBe(true));
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockRejectedValue(new Error('secret backend detail'));
  await act(async () => { await result.current.mutation.mutateAsync(savedConfig()); });
  await waitFor(() => expect(result.current.refresh.status).toBe('stale-error'));
  expect(result.current.mutation.isSuccess).toBe(true);
  expect(result.current.refresh.confirmed).toBe(true);
  const config = deferred<ReturnType<typeof savedConfig>>();
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockReturnValue(config.promise);
  let first!: Promise<void>, second!: Promise<void>;
  act(() => { first = refreshConfigReads(client, configReadContext('ws-1')); second = result.current.refresh.retry(); });
  expect(first).toBe(second);
  await act(async () => { config.resolve(savedConfig()); await first; });
  await waitFor(() => expect(result.current.refresh.status).toBe('fresh'));
  expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1);
});
it('cannot establish the expected owner from an otherwise valid ACK', async () => {
  const client = new QueryClient();
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const { result } = renderHook(() => useUpdateStudioV2DatabaseConfigMutation(), { wrapper });
  await act(async () => { await expect(result.current.mutateAsync(savedConfig())).rejects.toThrow('confirmation unavailable'); });
});
it.each([true, false])('clears stale refresh status after independent successful reads (confirmed save: %s)', async confirmed => {
  const { result, client } = mount();
  await waitFor(() => expect(result.current.workspace.isSuccess && result.current.config.isSuccess && result.current.targets.isSuccess).toBe(true));
  if (confirmed) vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockRejectedValue(new Error('refresh unavailable'));
  else vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockRejectedValue(new Error('ACK lost'));
  await act(async () => { await result.current.mutation.mutateAsync(savedConfig()).catch(() => undefined); });
  await waitFor(() => expect(result.current.refresh.status).toBe('stale-error'));
  await act(async () => { await result.current.workspace.refetch(); });
  expect(result.current.refresh.unavailable).toBe(true); // Bootstrap alone is insufficient.
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(savedConfig());
  await act(async () => { await client.refetchQueries({ type: 'active' }); });
  await waitFor(() => expect(result.current.refresh.unavailable).toBe(false));
  expect(result.current.refresh.confirmed).toBe(confirmed);
  expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1);
});
it('accepts the synthetic 1 ms ACK before three independent 200 ms reads, not a 600 ms chain', async () => {
  const hook = mount();
  await waitFor(() => expect(hook.result.current.workspace.isSuccess && hook.result.current.config.isSuccess && hook.result.current.targets.isSuccess).toBe(true));
  vi.useFakeTimers(); vi.setSystemTime(0);
  const started: number[] = [];
  const delayedRead = <T,>(value: T) => { started.push(Date.now()); return new Promise<T>(resolve => setTimeout(() => resolve(value), 200)); };
  vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockImplementation(() => new Promise(resolve => setTimeout(() => resolve(savedConfig()), 1)));
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockImplementation(() => delayedRead(savedConfig()));
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockImplementation(() => delayedRead([]));
  vi.mocked(studioV2WorkspaceAPI.get).mockImplementation(() => delayedRead(workspace));
  let acknowledgedAt: number | undefined;
  try {
    act(() => { void hook.result.current.mutation.mutateAsync(savedConfig()).then(() => { acknowledgedAt = Date.now(); }); });
    await act(async () => { await vi.advanceTimersByTimeAsync(1); });
    expect(acknowledgedAt).toBe(1);
    await act(async () => { await vi.advanceTimersByTimeAsync(1); });
    expect(started).toEqual([2, 2, 2]);
    expect(hook.result.current.refresh.status).toBe('refreshing');
    await act(async () => { await vi.advanceTimersByTimeAsync(199); });
    expect(hook.result.current.refresh.unavailable).toBe(true);
    await act(async () => { await vi.advanceTimersByTimeAsync(2); });
    expect(hook.result.current.refresh.status).toBe('fresh');
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(1);
  } finally { hook.unmount(); hook.client.clear(); vi.useRealTimers(); }
});
it('does not clear the pending write guard when independent reads finish before its ACK', async () => {
  const { result, client } = mount();
  await waitFor(() => expect(result.current.workspace.isSuccess && result.current.config.isSuccess && result.current.targets.isSuccess).toBe(true));
  const ack = deferred<ReturnType<typeof savedConfig>>();
  vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockReturnValue(ack.promise);
  let save!: Promise<unknown>;
  act(() => { save = result.current.mutation.mutateAsync(savedConfig()); });
  await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledOnce());
  await act(async () => { await client.refetchQueries({ type: 'active' }); });
  expect(result.current.refresh.unavailable).toBe(true);
  await act(async () => { ack.resolve(savedConfig()); await save; });
});
