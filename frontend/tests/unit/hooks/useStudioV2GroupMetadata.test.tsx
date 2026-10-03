import type { PropsWithChildren } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useStep4TargetColumns } from '@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('@/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: { getMetadata: vi.fn() } }));

describe('saved group metadata query scope', () => {
  beforeEach(() => vi.clearAllMocks());

  function wrapper({ children }: PropsWithChildren) {
    return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  }
  let client: QueryClient;
  beforeEach(() => { client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); });

  it('requests group B and refetches after the saved group revision changes', async () => {
    const connector = { ...savedGroupState().db.connector, schema: 'main', table: 'group_b' };
    vi.mocked(studioV2WorkspaceDatabaseAPI.getMetadata).mockResolvedValue({
      workspace_id: 'ws-1', connector_id: 'conn-1', connector_revision: 'crev-1',
      database: connector.database, schema: 'main', table: 'group_b', inspection_status: 'exists',
      columns: [{ name: 'counter', data_type: 'TEXT', nullable: false, primary_key: false }],
    });
    const { result, rerender } = renderHook(({ revision }) => useStep4TargetColumns(connector,
      { groupId: 'group-B', groupRevision: revision }), { wrapper, initialProps: { revision: 'rev-1' } });
    await waitFor(() => expect(result.current.status).toBe('exists'));
    expect(studioV2WorkspaceDatabaseAPI.getMetadata).toHaveBeenLastCalledWith('crev-1', { groupId: 'group-B', groupRevision: 'rev-1' });
    expect(result.current.columns).toEqual([{ name: 'counter', type: 'TEXT', nullable: false, primary_key: false }]);
    rerender({ revision: 'rev-2' });
    await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.getMetadata).toHaveBeenLastCalledWith('crev-1', { groupId: 'group-B', groupRevision: 'rev-2' }));
  });

  it('does not inspect an unsaved table or expose columns from the connector default table', async () => {
    const connector = { ...savedGroupState().db.connector, persisted: false, table: 'unsaved_group_b' };
    const { result } = renderHook(() => useStep4TargetColumns(connector, { groupId: 'group-B', groupRevision: 'rev-1' }), { wrapper });
    expect(result.current.status).toBe('not_checked');
    expect(result.current.columns).toEqual([]);
    expect(studioV2WorkspaceDatabaseAPI.getMetadata).not.toHaveBeenCalled();
  });
});
