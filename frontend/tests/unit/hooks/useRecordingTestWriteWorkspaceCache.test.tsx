import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { useWriteGroupTestWriteOperationQuery } from '@/hooks/datalink/useStudioV2WriteGroupTestWrite';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { testWriteOperation: vi.fn() } }));
it('does not reuse W1 operation results under W2 even when the opaque ID repeats', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.testWriteOperation).mockResolvedValue({
    operation_id: 'O1', action: 'test_write', status: 'running', created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z',
  });
  const view = renderHook(({ workspace }) => ({ ...useWriteGroupTestWriteOperationQuery('O1', workspace) }), { initialProps: { workspace: 'W1' }, wrapper });
  await act(async () => { await view.result.current.refetch(); });
  await waitFor(() => expect(view.result.current.data?.operation_id).toBe('O1'));
  view.rerender({ workspace: 'W2' });
  expect(view.result.current.data).toBeUndefined();
  expect(studioV2WorkspaceWriteGroupsAPI.testWriteOperation).toHaveBeenCalledTimes(1);
});
