import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useCreateStudioV2WorkspaceDeviceMutation,
  useStudioV2WorkspaceDevicesQuery,
  useUpdateStudioV2WorkspaceDeviceMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceDevices';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';

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

vi.mock('@/services/studioV2WorkspaceDevices', () => ({
  studioV2WorkspaceDevicesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
    updateOrder: vi.fn(),
  },
}));

describe('useStudioV2WorkspaceDevices hooks', () => {
  beforeEach(() => {
    useMutationMock.mockReset();
    useQueryMock.mockReset();
    invalidateQueriesMock.mockReset();
  });

  it('queries the workspace device list with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: [] });

    useStudioV2WorkspaceDevicesQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.devices(),
        queryFn: expect.any(Function),
        enabled: true,
        retry: false,
        refetchOnWindowFocus: false,
      }),
    );

    const options = useQueryMock.mock.calls[0]?.[0] as {
      queryFn: () => Promise<unknown>;
    };

    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValueOnce([]);
    await expect(options.queryFn()).resolves.toEqual([]);
  });

  it('creates one workspace device and invalidates the workspace device query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useCreateStudioV2WorkspaceDeviceMutation();

    const options = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (payload: unknown) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockResolvedValueOnce({ id: 'dev-01' } as any);
    await expect(options.mutationFn({ id: 'dev-01' })).resolves.toEqual({ id: 'dev-01' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.devices(),
    });
  });

  it('updates one workspace device and invalidates the workspace device query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useUpdateStudioV2WorkspaceDeviceMutation();

    const options = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (payload: { deviceId: string; request: unknown }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockResolvedValueOnce({ id: 'dev-01' } as any);
    await expect(
      options.mutationFn({ deviceId: 'dev-01', request: { name: 'Saved' } }),
    ).resolves.toEqual({ id: 'dev-01' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.devices(),
    });
  });
});
