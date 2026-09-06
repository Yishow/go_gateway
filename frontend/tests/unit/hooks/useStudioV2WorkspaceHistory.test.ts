import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import {
  useExportHistoryCSVMutation,
  useStudioV2WorkspaceHistoryQuery,
} from '@/hooks/datalink/useStudioV2WorkspaceHistory';
import { studioV2WorkspaceHistoryAPI } from '@/services/studioV2WorkspaceHistory';
import type { HistoryQuery } from '@/types/historyReport';

const useMutationMock = vi.fn();
const useQueryMock = vi.fn();

vi.mock('@tanstack/react-query', async () => {
  const actual = await vi.importActual<typeof import('@tanstack/react-query')>('@tanstack/react-query');
  return {
    ...actual,
    useQuery: (options: unknown) => useQueryMock(options),
    useMutation: (options: unknown) => useMutationMock(options),
  };
});

vi.mock('@/services/studioV2WorkspaceHistory', () => ({
  studioV2WorkspaceHistoryAPI: {
    queryHistory: vi.fn(),
    exportCSV: vi.fn(),
  },
}));

describe('useStudioV2WorkspaceHistory hooks', () => {
  beforeEach(() => {
    useMutationMock.mockReset();
    useQueryMock.mockReset();
  });

  it('queries history with dedicated composite key', () => {
    useQueryMock.mockReturnValue({ data: null });
    const query: HistoryQuery = {
      plan_id: 'plan-1',
      measurement_ids: ['meas-kw'],
      resolution: '1m',
    };
    useStudioV2WorkspaceHistoryQuery(query, true);

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.history('plan-1:meas-kw:::1m'),
        enabled: true,
        retry: false,
      })
    );
  });

  it('exports CSV via mutation', async () => {
    useMutationMock.mockReturnValue({ mutateAsync: vi.fn() });

    useExportHistoryCSVMutation();
    const exportOptions = useMutationMock.mock.calls[0]?.[0] as {
      mutationFn: (req: HistoryQuery) => Promise<Blob>;
    };

    const mockBlob = new Blob(['observed_at,value\n']);
    vi.mocked(studioV2WorkspaceHistoryAPI.exportCSV).mockResolvedValueOnce(mockBlob);

    const query: HistoryQuery = { plan_id: 'plan-1' };
    await expect(exportOptions.mutationFn(query)).resolves.toEqual(mockBlob);
  });
});
