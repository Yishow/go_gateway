import { QueryClient } from '@tanstack/react-query';
import { recordingHomeSourceContext } from '@/hooks/datalink/recordingHomeSourceCache';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useStudioV2DatabaseConfigQuery,
  useStudioV2DatabaseTargetsQuery,
  useUpdateStudioV2DatabaseConfigMutation,
  useUpsertStudioV2DatabaseTargetMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceDatabase';
import { studioV2WorkspaceDatabaseAPI, type StudioV2WorkspaceDatabaseConfigRequest } from '@/services/studioV2WorkspaceDatabase';

const useMutationMock = vi.fn();
const useQueryMock = vi.fn();
let queryClient: QueryClient;

vi.mock('@tanstack/react-query', async () => {
  const actual =
    await vi.importActual<typeof import('@tanstack/react-query')>(
      '@tanstack/react-query',
    );

  return {
    ...actual,
    useQuery: (options: unknown) => useQueryMock(options),
    useMutation: (options: unknown) => useMutationMock(options),
    useQueryClient: () => queryClient,
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
    queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  });

  it('queries workspace database config with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: null });

    useStudioV2DatabaseConfigQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(expect.objectContaining({
      queryKey: [...studioV2WorkspaceKeys.databaseConfig(), 'api-context', expect.any(Number), recordingHomeSourceContext().apiOrigin, ''],
      queryFn: expect.any(Function),
      enabled: true,
      retry: false,
      refetchOnWindowFocus: false,
    }));

    const options = useQueryMock.mock.calls[0]?.[0] as { queryFn: (context: { signal: AbortSignal; client: QueryClient }) => Promise<unknown> };
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValueOnce(null);
    const signal = new AbortController().signal;
    await expect(options.queryFn({ signal, client: queryClient })).resolves.toBeNull();
    expect(studioV2WorkspaceDatabaseAPI.getConfig).toHaveBeenCalledWith({ signal, baseURL: recordingHomeSourceContext().apiOrigin });
  });

  it('queries workspace database targets with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: [] });

    useStudioV2DatabaseTargetsQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(expect.objectContaining({
      queryKey: [...studioV2WorkspaceKeys.databaseTargets(), 'api-context', expect.any(Number), recordingHomeSourceContext().apiOrigin, ''],
      queryFn: expect.any(Function),
      enabled: true,
      retry: false,
      refetchOnWindowFocus: false,
    }));

    const options = useQueryMock.mock.calls[0]?.[0] as { queryFn: (context: { signal: AbortSignal; client: QueryClient }) => Promise<unknown> };
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValueOnce([]);
    const signal = new AbortController().signal;
    await expect(options.queryFn({ signal, client: queryClient })).resolves.toEqual([]);
    expect(studioV2WorkspaceDatabaseAPI.listTargets).toHaveBeenCalledWith({ signal, baseURL: recordingHomeSourceContext().apiOrigin });
  });

  it('updates workspace database config and preserves apply_failed runtime status', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });
    const keys = [studioV2WorkspaceKeys.databaseConfig(), studioV2WorkspaceKeys.databaseTargets(), studioV2WorkspaceKeys.bootstrap()];
    queryClient.setQueryData(keys[0], null);
    queryClient.setQueryData(keys[1], []);
    queryClient.setQueryData(keys[2], { id: 'ws-1' });
    useUpdateStudioV2DatabaseConfigMutation('ws-1');

    const options = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (payload: StudioV2WorkspaceDatabaseConfigRequest) => Promise<unknown>;
    };
    const request: StudioV2WorkspaceDatabaseConfigRequest = {
      connector_id: 'db-01', kind: 'sqlite', name: 'Line A SQLite', host: '', port: 0,
      database: '/tmp/line-a.sqlite', username: '', schema: 'main', table: 'sensor_values',
      write_mode: 'insert', write_interval_seconds: 5, timestamp_column: 'ts',
    };
    const saved = {
      ...request, id: 'db-01', workspace_id: 'ws-1', identity_revision: 'identity-1', setup_revision: 'setup-2',
      status: 'ready' as const, save_state: 'saved' as const, created_at: '2026-05-30T00:00:00Z', updated_at: '2026-05-30T00:00:01Z',
      runtime_apply_status: 'apply_failed' as const, runtime_apply_message: 'runtime sync failed',
    };
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockResolvedValueOnce(saved);
    await expect(options.mutationFn(request)).resolves.toEqual(saved);
    expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledWith(request, { baseURL: recordingHomeSourceContext().apiOrigin });
    keys.forEach(key => expect(queryClient.getQueryState(key)?.isInvalidated).toBe(true));
  });

  it('upserts one workspace database target row and invalidates targets query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    const invalidateQueriesMock = vi.spyOn(queryClient, 'invalidateQueries');
    useUpsertStudioV2DatabaseTargetMutation();

    const options = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (payload: { pointId: string; request: unknown }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockResolvedValueOnce({
      id: 'row-01',
      runtime_apply_status: 'not_running',
    } as unknown as Awaited<ReturnType<typeof studioV2WorkspaceDatabaseAPI.upsertTarget>>);
    await expect(
      options.mutationFn({ pointId: 'point-01', request: { column_name: 'line_a', enabled: true } }),
    ).resolves.toEqual({ id: 'row-01', runtime_apply_status: 'not_running' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.databaseTargets() });
    expect(invalidateQueriesMock).toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.bootstrap() });
  });
});
