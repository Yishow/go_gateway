import { describe, expect, it, vi } from 'vitest';
import { deviceKeys } from '@/hooks/datalink/keys';
import { useDevicesQuery } from '@/hooks/datalink/useDevices';
import { deviceAPI } from '@/services/datalink';

const useQueryMock = vi.fn();

vi.mock('@tanstack/react-query', async () => {
  const actual =
    await vi.importActual<typeof import('@tanstack/react-query')>(
      '@tanstack/react-query',
    );

  return {
    ...actual,
    useQuery: (options: unknown) => useQueryMock(options),
  };
});

vi.mock('@/services/datalink', () => ({
  deviceAPI: {
    list: vi.fn(),
  },
}));

describe('useDevicesQuery', () => {
  it('keeps previous data and disables window-focus refetching', async () => {
    useQueryMock.mockReturnValue({ data: [] });

    useDevicesQuery({ status: 'active' });

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: deviceKeys.list({ status: 'active' }),
        queryFn: expect.any(Function),
        refetchOnWindowFocus: false,
        placeholderData: expect.any(Function),
      }),
    );

    const options = useQueryMock.mock.calls[0]?.[0] as {
      queryFn: () => Promise<unknown>;
      placeholderData: <T>(previousData: T) => T;
    };

    vi.mocked(deviceAPI.list).mockResolvedValueOnce([]);

    await expect(options.queryFn()).resolves.toEqual([]);
    expect(deviceAPI.list).toHaveBeenCalledWith({ status: 'active' });
    expect(options.placeholderData(['cached-device'])).toEqual(['cached-device']);
  });
});
