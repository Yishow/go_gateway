import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, expect, it, vi } from 'vitest';
import { beginConfigSave, confirmConfigSave, configReadContext, fencedConfigRead } from '@/hooks/datalink/recordingConfigRefresh';
import { useStudioV2WorkspaceQuery } from '@/hooks/datalink/useStudioV2Workspace';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { captureRecordingHomeSourceMutation } from '@/hooks/datalink/recordingHomeSourceMutation';
import type { StudioV2Workspace } from '@/types/datalink';
const baseURL = studioV2DatalinkApi.defaults.baseURL;
afterEach(() => { studioV2DatalinkApi.defaults.baseURL = baseURL; vi.restoreAllMocks(); });
it('keeps another workspace read independent of a configuration mutation', async () => {
  const client = new QueryClient(), context = configReadContext('other-workspace');
  let finish!: (value: string) => void;
  const pending = fencedConfigRead(client, context, new AbortController().signal, () => new Promise<string>(resolve => { finish = resolve; }));
  beginConfigSave(client, 'ws-1');
  finish('other workspace observation');
  await expect(pending).resolves.toBe('other workspace observation');
});
it('isolates bootstrap cache by captured API origin and preserves generic home mutation owner lookup', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const get = vi.spyOn(studioV2WorkspaceAPI, 'get').mockResolvedValue({ id: 'W1' } as StudioV2Workspace);
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const view = renderHook(() => useStudioV2WorkspaceQuery(), { wrapper });
  await waitFor(() => expect(view.result.current.data?.id).toBe('W1'));
  expect(captureRecordingHomeSourceMutation(client)?.workspaceId).toBe('W1');
  get.mockResolvedValue({ id: 'W2' } as StudioV2Workspace);
  act(() => { studioV2DatalinkApi.defaults.baseURL = '/different/api/datalink'; view.rerender(); });
  await waitFor(() => expect(view.result.current.data?.id).toBe('W2'));
  expect(captureRecordingHomeSourceMutation(client)?.workspaceId).toBe('W2');
  expect(get.mock.calls[0]?.[0]).toEqual(expect.objectContaining({ signal: expect.any(AbortSignal) }));
});
it('rejects ignored-cancellation late reads after a newer save and does not cross API origins', async () => {
  const client = new QueryClient(), context = configReadContext('ws-1');
  let finish!: (value: string) => void;
  const old = fencedConfigRead(client, context, new AbortController().signal, () => new Promise<string>(resolve => { finish = resolve; }));
  beginConfigSave(client, 'ws-1'); finish('old');
  await expect(old).rejects.toThrow('superseded');
  const originRead = fencedConfigRead(client, context, new AbortController().signal, async () => {
    studioV2DatalinkApi.defaults.baseURL = '/new/api'; return 'wrong-origin';
  });
  await expect(originRead).rejects.toThrow('superseded');
});

it('fences a read started during the write before its successful ACK', async () => {
  const client = new QueryClient(), context = configReadContext('ws-1');
  const captured = beginConfigSave(client, 'ws-1');
  let finish!: (value: string) => void;
  const pending = fencedConfigRead(client, context, new AbortController().signal, () => new Promise<string>(resolve => { finish = resolve; }));
  confirmConfigSave(client, captured); finish('pre-commit read');
  await expect(pending).rejects.toThrow('superseded');
});
it('coalesces superseded same-turn refresh demand without delaying either ACK', async () => {
  const client = new QueryClient(), refresh = vi.spyOn(client, 'refetchQueries').mockResolvedValue();
  confirmConfigSave(client, beginConfigSave(client, 'ws-1'));
  confirmConfigSave(client, beginConfigSave(client, 'ws-1'));
  await new Promise(resolve => setTimeout(resolve, 0));
  expect(refresh.mock.calls.filter(([filter]) => filter && 'predicate' in filter)).toHaveLength(1);
});
it('uses only the active scoped bootstrap owner, never cached inactive or failing authority', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const { QueryObserver } = await import('@tanstack/react-query');
  const { configReadKey } = await import('@/hooks/datalink/recordingConfigRefresh');
  const { studioV2WorkspaceKeys } = await import('@/hooks/datalink/keys');
  const oldKey = configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext('A'));
  client.setQueryData(oldKey, { id: 'A' });
  const activeKey = configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext('B'));
  const observer = new QueryObserver(client, { queryKey: activeKey, initialData: { id: 'B' }, staleTime: Infinity, queryFn: async () => { throw new Error('read unavailable'); }, retry: false });
  const unsubscribe = observer.subscribe(() => undefined);
  expect(captureRecordingHomeSourceMutation(client)?.workspaceId).toBe('B');
  await observer.refetch();
  expect(captureRecordingHomeSourceMutation(client)).toBeUndefined();
  unsubscribe();
  expect(captureRecordingHomeSourceMutation(client)).toBeUndefined();
});
