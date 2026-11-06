import { QueryClient } from '@tanstack/react-query';
import { recordingHomeSourceContext } from '@/hooks/datalink/recordingHomeSourceCache';
import { describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { useStudioV2WorkspaceQuery } from '@/hooks/datalink/useStudioV2Workspace';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';

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

vi.mock('@/services/studioV2Workspace', async importOriginal => ({
  ...await importOriginal<typeof import('@/services/studioV2Workspace')>(),
  studioV2WorkspaceAPI: {
    get: vi.fn(),
  },
}));

describe('useStudioV2WorkspaceQuery', () => {
  it('uses the bootstrap key and returns workspace metadata', async () => {
    useQueryMock.mockReturnValue({ data: null });

    useStudioV2WorkspaceQuery();

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: [...studioV2WorkspaceKeys.bootstrap(), 'api-context', expect.any(Number), recordingHomeSourceContext().apiOrigin, ''],
        queryFn: expect.any(Function),
        retry: false,
        refetchOnWindowFocus: false,
      }),
    );

    const options = useQueryMock.mock.calls[0]?.[0] as {
      queryFn: (context: { signal: AbortSignal; client: QueryClient }) => Promise<unknown>;
    };

    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'ws-1',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });

    const signal = new AbortController().signal;
    await expect(options.queryFn({ signal, client: new QueryClient() })).resolves.toEqual({
      id: 'ws-1',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    expect(studioV2WorkspaceAPI.get).toHaveBeenCalledTimes(1);
    expect(studioV2WorkspaceAPI.get).toHaveBeenCalledWith({ signal, baseURL: recordingHomeSourceContext().apiOrigin });
  });
});
