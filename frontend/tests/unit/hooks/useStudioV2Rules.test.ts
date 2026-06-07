import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useCreateStudioV2RuleMutation,
  useStudioV2RulesQuery,
  useUpdateStudioV2RuleMutation,
} from '@/hooks/datalink/useStudioV2Rules';
import { studioV2RulesAPI } from '@/services/studioV2Rules';

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

vi.mock('@/services/studioV2Rules', () => ({
  studioV2RulesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

describe('useStudioV2Rules hooks', () => {
  beforeEach(() => {
    useMutationMock.mockReset();
    useQueryMock.mockReset();
    invalidateQueriesMock.mockReset();
  });

  it('queries workspace source rules with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: [] });

    useStudioV2RulesQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.sourceRules(),
        queryFn: expect.any(Function),
        enabled: true,
        retry: false,
        refetchOnWindowFocus: false,
      }),
    );

    const options = useQueryMock.mock.calls[0]?.[0] as {
      queryFn: () => Promise<unknown>;
    };

    vi.mocked(studioV2RulesAPI.list).mockResolvedValueOnce([]);
    await expect(options.queryFn()).resolves.toEqual([]);
  });

  it('creates one workspace source rule and preserves applied runtime status', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useCreateStudioV2RuleMutation();

    const options = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (payload: unknown) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2RulesAPI.create).mockResolvedValueOnce({
      id: 'rule-01',
      runtime_apply_status: 'applied',
    } as any);
    await expect(options.mutationFn({ id: 'rule-01' })).resolves.toEqual({
      id: 'rule-01',
      runtime_apply_status: 'applied',
    });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.sourceRules(),
    });
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.bootstrap(),
    });
  });

  it('updates one workspace source rule and invalidates the rule query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useUpdateStudioV2RuleMutation();

    const options = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (payload: { ruleId: string; request: unknown }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2RulesAPI.update).mockResolvedValueOnce({
      id: 'rule-01',
      runtime_apply_status: 'not_running',
    } as any);
    await expect(
      options.mutationFn({ ruleId: 'rule-01', request: { naming_prefix: 'LINE_' } }),
    ).resolves.toEqual({ id: 'rule-01', runtime_apply_status: 'not_running' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.sourceRules(),
    });
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.bootstrap(),
    });
  });
});
