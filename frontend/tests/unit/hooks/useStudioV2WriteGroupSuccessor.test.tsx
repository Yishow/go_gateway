import type { PropsWithChildren } from 'react';
import { act, renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useWriteGroupSuccessorActions, useCheckSuccessorOperation } from '@/hooks/datalink/useStudioV2WriteGroupSuccessor';
import { useWriteGroupLifecycleQuery } from '@/hooks/datalink/useStudioV2WriteGroupLifecycle';
import { studioV2WriteGroupLifecycleKeys as keys } from '@/hooks/datalink/studioV2WriteGroupLifecycleKeys';
import { studioV2WriteGroupSuccessorAPI as api } from '@/services/studioV2WriteGroupSuccessor';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { lifecycle, operationStatus, envelope, preparation } from '../../fixtures/writeGroupSuccessor';
import { DIGEST } from '../../fixtures/cutoverFixtures';
vi.mock('@/services/studioV2WriteGroupSuccessor', () => ({ studioV2WriteGroupSuccessorAPI: {
  review: vi.fn(), reserve: vi.fn(), schemaPreview: vi.fn(), schemaConfirm: vi.fn(), activationReview: vi.fn(), activationReadiness: vi.fn(), confirmCutover: vi.fn(), confirmFirstStart: vi.fn(), operationStatus: vi.fn(),
} }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
function setup() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: 3, retryDelay: 1 }, mutations: { retry: 3, retryDelay: 1 } } });
  const wrapper = ({ children }: PropsWithChildren) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return { client, wrapper };
}
beforeEach(() => vi.resetAllMocks());
describe('scoped lifecycle queries', () => {
  it('puts workspace, original predecessor, lineage and root ID in separate cache dimensions', () => {
    const base = keys.operation('ws-1', 'g1', 'l1', 'op1');
    expect(base).not.toEqual(keys.operation('ws-2', 'g1', 'l1', 'op1'));
    expect(base).not.toEqual(keys.operation('ws-1', 'g2', 'l1', 'op1'));
    expect(base).not.toEqual(keys.operation('ws-1', 'g1', 'l2', 'op1'));
    expect(base).not.toEqual(keys.operation('ws-1', 'g1', 'l1', 'op2'));
  });
  it('does not cache a lifecycle from a different workspace', async () => {
    const { wrapper } = setup(); vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValue(lifecycle());
    const { result } = renderHook(() => useWriteGroupLifecycleQuery('foreign', 'wide-1'), { wrapper });
    await waitFor(() => expect(result.current.isError).toBe(true));
    expect(result.current.data).toBeUndefined();
  });
  it('reads the same original operation on demand even after unrelated cached operations', async () => {
    const { client, wrapper } = setup(); client.setQueryData(keys.operation('ws-1', 'wide-1', 'lineage-1', 'new-op'), operationStatus());
    vi.mocked(api.operationStatus).mockResolvedValue(operationStatus());
    const { result } = renderHook(() => useCheckSuccessorOperation(), { wrapper });
    await act(() => result.current('wide-1', { workspace_id: 'ws-1', lineage_id: 'lineage-1', operation_id: 'root-op-1' }));
    expect(api.operationStatus).toHaveBeenCalledWith('wide-1', { workspace_id: 'ws-1', lineage_id: 'lineage-1', operation_id: 'root-op-1' });
  });
  it('does not retry a lost reservation despite inherited retry defaults', async () => {
    const { wrapper } = setup(); vi.mocked(api.reserve).mockRejectedValue(new Error('lost'));
    const { result } = renderHook(() => useWriteGroupSuccessorActions(), { wrapper });
    await act(async () => { await expect(result.current.reserve.mutateAsync({ groupId: 'wide-1', request: { envelope: envelope(), review_digest: DIGEST } })).rejects.toThrow(); });
    expect(api.reserve).toHaveBeenCalledTimes(1);
  });
  it('invalidates grouped lifecycle and bootstrap after reservation', async () => {
    const { client, wrapper } = setup(); const invalidate = vi.spyOn(client, 'invalidateQueries');
    vi.mocked(api.reserve).mockResolvedValue(preparation());
    const { result } = renderHook(() => useWriteGroupSuccessorActions(), { wrapper });
    await act(() => result.current.reserve.mutateAsync({ groupId: 'wide-1', request: { envelope: envelope(), review_digest: DIGEST } }));
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ['studio-v2-workspace', 'write-groups'] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ['studio-v2-workspace', 'bootstrap'] });
  });
  it('invalidates exact schema operation evidence as well as metadata after confirmation', async () => {
    const { client, wrapper } = setup(); const invalidate = vi.spyOn(client, 'invalidateQueries');
    vi.mocked(api.schemaConfirm).mockResolvedValue({ operation_id: 'schema-1', status: 'succeeded', executed_statements: 1, created_at: 'now', updated_at: 'now' });
    const { result } = renderHook(() => useWriteGroupSuccessorActions(), { wrapper });
    await act(() => result.current.schemaConfirm.mutateAsync({ groupId: 'wide-1', request: { envelope: envelope(), review_digest: DIGEST,
      current: { workspace_revision: 'ws-r2', lineage_revision: 2, predecessor_revision: 'rev-wide', predecessor_connector_revision: 'crev-1', group_id: 'new-2', group_revision: 'g-r2', target_connector_revision: 'crev-1' }, token: 'token-1', digest: DIGEST, schema_operation_id: 'schema-1' } }));
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ['studio-v2-workspace', 'write-group-schema', 'schema-1'] });
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ['studio-v2-workspace', 'database-metadata'] });
  });
});
