import { describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { useStudioV2WorkspaceAuditHistoryQuery } from '@/hooks/datalink/useStudioV2WorkspaceAuditHistory';
import { studioV2WorkspaceAuditAPI } from '@/services/studioV2WorkspaceAudit';

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

vi.mock('@/services/studioV2WorkspaceAudit', () => ({
  studioV2WorkspaceAuditAPI: {
    list: vi.fn(),
  },
}));

describe('useStudioV2WorkspaceAuditHistoryQuery', () => {
  it('queries recent workspace audit history with the dedicated key', async () => {
    useQueryMock.mockReturnValue({ data: { entries: [] } });

    useStudioV2WorkspaceAuditHistoryQuery(true);

    expect(useQueryMock).toHaveBeenCalledWith(
      expect.objectContaining({
        queryKey: studioV2WorkspaceKeys.auditHistory(),
        queryFn: expect.any(Function),
        enabled: true,
        retry: false,
        refetchOnWindowFocus: false,
      }),
    );

    const options = useQueryMock.mock.calls[0]?.[0] as {
      queryFn: () => Promise<unknown>;
    };

    vi.mocked(studioV2WorkspaceAuditAPI.list).mockResolvedValueOnce({
      entries: [{
        id: 'audit-1',
        workspace_id: 'ws-1',
        event_type: 'workspace_activation',
        result: 'partial_success',
        scope: 'devices:dev-A,dev-B',
        occurred_at: '2026-05-29T10:12:00Z',
        created_at: '2026-05-29T10:12:01Z',
      }],
    });

    await expect(options.queryFn()).resolves.toEqual({
      entries: [{
        id: 'audit-1',
        workspace_id: 'ws-1',
        event_type: 'workspace_activation',
        result: 'partial_success',
        scope: 'devices:dev-A,dev-B',
        occurred_at: '2026-05-29T10:12:00Z',
        created_at: '2026-05-29T10:12:01Z',
      }],
    });
    expect(studioV2WorkspaceAuditAPI.list).toHaveBeenCalledWith(10);
  });
});
