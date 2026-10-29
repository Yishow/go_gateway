import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { useRecordingWorkspaceGuard } from '@/features/datalink/workbench-v2/steps/step4/useRecordingWorkspaceGuard';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import type { WriteGroupListResponse } from '@/types/studioV2WriteGroup';
const key = studioV2WorkspaceKeys.writeGroups();
const clients: QueryClient[] = [];
const data = (workspaceId = 'W1', revision = 'R1'): WriteGroupListResponse => ({ workspace_id: workspaceId, workspace_revision: revision, groups: [] });
function host(initial?: WriteGroupListResponse) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
  if (initial) client.setQueryData(key, initial);
  const accepted = client.getQueryData<WriteGroupListResponse>(key);
  const hook = renderHook(({ workspaceId, snapshot }: { workspaceId?: string; snapshot?: WriteGroupListResponse }) =>
    useRecordingWorkspaceGuard(workspaceId, snapshot), { initialProps: { workspaceId: 'W1', snapshot: accepted },
    wrapper: ({ children }) => <QueryClientProvider client={client}>{children}</QueryClientProvider> });
  return { ...hook, client, accepted };
}
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
it('only removes its denial for the exact healthy accepted canonical cache entry and never queries', () => {
  const view = host(data()); const fetch = vi.spyOn(view.client, 'fetchQuery');
  expect(view.result.current()).toBe(false); expect(fetch).not.toHaveBeenCalled();
});
it('denies absent cache even with a caller-supplied plausible snapshot', () => {
  const view = host(); view.rerender({ workspaceId: 'W1', snapshot: data() });
  expect(view.result.current()).toBe(true);
});
it.each([data('W2'), data('W1', ''), data('W1', 'R2')])('synchronously denies a changed current head before observer render', next => {
  const view = host(data()); const captured = view.result.current;
  act(() => { view.client.setQueryData(key, next); });
  expect(captured()).toBe(true);
});
it('requires the actual cached object rather than a cloned or foreign claimed acceptance', () => {
  const view = host(data());
  view.rerender({ workspaceId: 'W1', snapshot: data() }); expect(view.result.current()).toBe(true);
  view.rerender({ workspaceId: 'W2', snapshot: view.accepted }); expect(view.result.current()).toBe(true);
});
it('keeps old A callbacks denied after accepted A to B to A, and after unmount', () => {
  const view = host(data()); const stale = view.result.current;
  view.rerender({ workspaceId: 'W2', snapshot: view.accepted }); view.rerender({ workspaceId: 'W1', snapshot: view.accepted });
  expect(stale()).toBe(true); expect(view.result.current()).toBe(false);
  const latest = view.result.current; view.unmount(); expect(latest()).toBe(true);
});
it('denies fetching and error states even when the last canonical data is still present', async () => {
  const view = host(data()); let reject!: (cause: Error) => void;
  const fetch = view.client.fetchQuery({ queryKey: key, queryFn: () => new Promise<WriteGroupListResponse>((_, fail) => { reject = fail; }), retry: false });
  expect(view.result.current()).toBe(true);
  await act(async () => { reject(new Error('unavailable')); await fetch.catch(() => undefined); });
  expect(view.client.getQueryData(key)).toBe(view.accepted); expect(view.result.current()).toBe(true);
});
