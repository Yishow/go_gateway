import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { expect, it, vi } from 'vitest';
import { useUpsertStudioV2DatabaseTargetMutation } from '@/hooks/datalink/useStudioV2WorkspaceDatabase';
import { useUpdateWriteGroupMutation } from '@/hooks/datalink/useStudioV2WriteGroups';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
function deferred() { let resolve!: () => void; const promise = new Promise<void>(yes => { resolve = yes; }); return { promise, resolve }; }
it.each(['target', 'write-group'] as const)('%s save still waits for both existing reconciliation reads', async kind => {
  const client = new QueryClient(), first = deferred(), second = deferred();
  const invalidate = vi.spyOn(client, 'invalidateQueries').mockReturnValueOnce(first.promise).mockReturnValueOnce(second.promise);
  vi.spyOn(studioV2WorkspaceDatabaseAPI, 'upsertTarget').mockResolvedValue({ id: 'row' } as never);
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'update').mockResolvedValue({ id: 'group' } as never);
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const { result } = renderHook(() => ({ target: useUpsertStudioV2DatabaseTargetMutation(), group: useUpdateWriteGroupMutation() }), { wrapper });
  let completed = false;
  act(() => { void (kind === 'target' ? result.current.target.mutateAsync({ pointId: 'p', request: { column_name: 'v', enabled: true } })
    : result.current.group.mutateAsync({ id: 'group', request: {} as never })).then(() => { completed = true; }); });
  await act(async () => {});
  expect(invalidate).toHaveBeenCalledTimes(1);
  expect(completed).toBe(false);
  await act(async () => { first.resolve(); });
  expect(invalidate).toHaveBeenCalledTimes(2);
  expect(completed).toBe(false);
  await act(async () => { second.resolve(); });
  expect(completed).toBe(true);
  vi.restoreAllMocks();
});
