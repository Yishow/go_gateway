import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useStudioV2DatabaseConfigQuery,
  useStudioV2DatabaseTargetsQuery,
  useUpdateStudioV2DatabaseConfigMutation,
  useUpsertStudioV2DatabaseTargetMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceDatabase';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';

const useMutationMock = vi.fn();
const useQueryMock = vi.fn();
const invalidateQueriesMock = vi.fn();

vi.mock('@tanstack/react-query', async () => {
  const actual =
    await vi.importActual<typeof import('@tanstack/react-query')>(
      '@tanstack/react-query',
    );

  return {
    ...actual,
    useQuery: (options: unknown) => useQueryMock(options),
    useMutation: (options: unknown) => useMutationMock(options),
    useQueryClient: () => ({
      invalidateQueries: invalidateQueriesMock,
    }),
  };
});

vi.mock('@/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
}));

describe('useStudioV2WorkspaceDatabase hooks', () => {
  beforeEach(() => {
    useMutationMock.mockReset();
    useQueryMock.mockReset();
    invalidateQueriesMock.mockReset();
  });

  it('queries workspace database config with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: null });

    useStudioV2DatabaseConfigQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(expect.objectContaining({
      queryKey: studioV2WorkspaceKeys.databaseConfig(),
      queryFn: expect.any(Function),
      enabled: true,
      retry: false,
      refetchOnWindowFocus: false,
    }));

    const options = useQueryMock.mock.calls[0]?.[0] as { queryFn: () => Promise<unknown> };
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValueOnce(null);
    await expect(options.queryFn()).resolves.toBeNull();
  });

  it('queries workspace database targets with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: [] });

    useStudioV2DatabaseTargetsQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(expect.objectContaining({
      queryKey: studioV2WorkspaceKeys.databaseTargets(),
      queryFn: expect.any(Function),
      enabled: true,
      retry: false,
      refetchOnWindowFocus: false,
    }));

    const options = useQueryMock.mock.calls[0]?.[0] as { queryFn: () => Promise<unknown> };
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValueOnce([]);
    await expect(options.queryFn()).resolves.toEqual([]);
  });

  it('updates workspace database config and preserves apply_failed runtime status', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useUpdateStudioV2DatabaseConfigMutation();

    const options = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (payload: unknown) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockResolvedValueOnce({
      id: 'db-01',
      runtime_apply_status: 'apply_failed',
      runtime_apply_message: 'runtime sync failed',
    } as any);
    await expect(options.mutationFn({ name: 'Line A SQLite' })).resolves.toEqual({
      id: 'db-01',
      runtime_apply_status: 'apply_failed',
      runtime_apply_message: 'runtime sync failed',
    });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.databaseConfig() });
    expect(invalidateQueriesMock).toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.databaseTargets() });
  });

  it('upserts one workspace database target row and invalidates targets query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useUpsertStudioV2DatabaseTargetMutation();

    const options = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (payload: { pointId: string; request: unknown }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockResolvedValueOnce({
      id: 'row-01',
      runtime_apply_status: 'not_running',
    } as any);
    await expect(
      options.mutationFn({ pointId: 'point-01', request: { column_name: 'line_a', enabled: true } }),
    ).resolves.toEqual({ id: 'row-01', runtime_apply_status: 'not_running' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.databaseTargets() });
  });
});
