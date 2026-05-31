import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useCreateStudioV2MappingMutation,
  useStudioV2MappingsQuery,
  useUpdateStudioV2MappingMutation,
} from '@/hooks/datalink/useStudioV2Mappings';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';

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

vi.mock('@/services/studioV2Mappings', () => ({
  studioV2MappingsAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

describe('useStudioV2Mappings hooks', () => {
  beforeEach(() => {
    useMutationMock.mockReset();
    useQueryMock.mockReset();
    invalidateQueriesMock.mockReset();
  });

  it('queries workspace mappings with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: [] });

    useStudioV2MappingsQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.mappings(),
        queryFn: expect.any(Function),
        enabled: true,
        retry: false,
        refetchOnWindowFocus: false,
      }),
    );

    const options = useQueryMock.mock.calls[0]?.[0] as {
      queryFn: () => Promise<unknown>;
    };

    vi.mocked(studioV2MappingsAPI.list).mockResolvedValueOnce([]);
    await expect(options.queryFn()).resolves.toEqual([]);
  });

  it('creates one workspace mapping row and invalidates the mapping query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useCreateStudioV2MappingMutation();

    const options = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (payload: unknown) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2MappingsAPI.create).mockResolvedValueOnce({
      id: 'mapping-01',
      runtime_apply_status: 'applied',
    } as any);
    await expect(options.mutationFn({ rule_id: 'rule-01' })).resolves.toEqual({
      id: 'mapping-01',
      runtime_apply_status: 'applied',
    });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.mappings(),
    });
  });

  it('updates one workspace mapping row and preserves not_running runtime status', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useUpdateStudioV2MappingMutation();

    const options = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (payload: { mappingId: string; request: unknown }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2MappingsAPI.update).mockResolvedValueOnce({
      id: 'mapping-01',
      runtime_apply_status: 'not_running',
    } as any);
    await expect(
      options.mutationFn({ mappingId: 'mapping-01', request: { tag_key: 'line.a.temp' } }),
    ).resolves.toEqual({ id: 'mapping-01', runtime_apply_status: 'not_running' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.mappings(),
    });
  });
});
