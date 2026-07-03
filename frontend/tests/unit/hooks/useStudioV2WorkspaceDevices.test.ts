import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useCreateStudioV2WorkspaceDeviceMutation,
  useDeleteStudioV2WorkspaceDeviceMutation,
  useStudioV2WorkspaceDevicesQuery,
  useUpdateStudioV2WorkspaceDeviceAvailabilityMutation,
  useUpdateStudioV2WorkspaceDeviceMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceDevices';
import {
  type StudioV2WorkspaceDeviceRecord,
  studioV2WorkspaceDevicesAPI,
} from '@/services/studioV2WorkspaceDevices';

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
    updateAvailability: vi.fn(),
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

    const createdDevice = { id: 'dev-01' } as StudioV2WorkspaceDeviceRecord;
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockResolvedValueOnce(createdDevice);
    await expect(options.mutationFn({ id: 'dev-01' })).resolves.toEqual({ id: 'dev-01' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.bootstrap(),
    });
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

    const restartRequiredDevice = {
      id: 'dev-01',
      runtime_apply_status: 'restart-required',
    } as StudioV2WorkspaceDeviceRecord;
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockResolvedValueOnce(restartRequiredDevice);
    await expect(
      options.mutationFn({
        deviceId: 'dev-01',
        request: { connection_config: { host: '192.168.10.20' } },
      }),
    ).resolves.toEqual({ id: 'dev-01', runtime_apply_status: 'restart-required' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.bootstrap(),
    });
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.devices(),
    });
  });

  it('updates availability and invalidates both the workspace bootstrap and device query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useUpdateStudioV2WorkspaceDeviceAvailabilityMutation();

    const options = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (payload: { deviceId: string; request: unknown }) => Promise<unknown>;
      onSuccess: () => Promise<void>;
    };

    const unavailableDevice = {
      id: 'dev-01',
      availability_status: 'unavailable',
    } as StudioV2WorkspaceDeviceRecord;
    vi.mocked(studioV2WorkspaceDevicesAPI.updateAvailability).mockResolvedValueOnce(unavailableDevice);
    await expect(
      options.mutationFn({
        deviceId: 'dev-01',
        request: { availability_status: 'unavailable', availability_reason: 'device form is invalid' },
      }),
    ).resolves.toEqual({ id: 'dev-01', availability_status: 'unavailable' });

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.bootstrap(),
    });
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.devices(),
    });
  });

  it('deletes one workspace device and invalidates both the workspace bootstrap and device query', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useDeleteStudioV2WorkspaceDeviceMutation();

    const options = useMutationMock.mock.calls.at(-1)?.[0] as {
      mutationFn: (deviceId: string) => Promise<void>;
      onSuccess: () => Promise<void>;
    };

    vi.mocked(studioV2WorkspaceDevicesAPI.remove).mockResolvedValueOnce();
    await expect(options.mutationFn('dev-01')).resolves.toBeUndefined();

    await options.onSuccess();
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.bootstrap(),
    });
    expect(invalidateQueriesMock).toHaveBeenCalledWith({
      queryKey: studioV2WorkspaceKeys.devices(),
    });
  });
});
