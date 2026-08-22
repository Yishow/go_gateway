import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  dbTargetFixture,
  resetDatabaseAutosaveMocks,
  runtimeApplied,
  studioV2WorkspaceDatabaseAPI,
} from './database-autosave-page.testHarness';

async function renderPage() {
  const { default: DatalinkWorkbenchV2Page } = await import(
    '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page'
  );
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

beforeEach(resetDatabaseAutosaveMocks);

describe('DatalinkWorkbenchV2Page database autosave orchestration', () => {
  it('autosaves bulk database target enable toggles for every persisted row', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockImplementation(async (pointId, request) =>
      runtimeApplied(dbTargetFixture({
        id: pointId === 'persisted-point-A' ? 'row-A' : 'row-B',
        point_id: pointId,
        tag_id: pointId === 'persisted-point-A' ? 'tag-A' : 'tag-B',
        column_name: request.column_name,
        enabled: request.enabled,
        row_group_id: request.row_group_id,
      })));

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));

    await waitFor(() => {
      expect(screen.getByTestId('target-enabled-point-A')).toHaveTextContent('true');
      expect(screen.getByTestId('target-enabled-point-B')).toHaveTextContent('true');
    });

    fireEvent.click(screen.getByTestId('disable-all-db-targets'));

    await waitFor(() => {
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-A', expect.objectContaining({ enabled: false }));
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-B', expect.objectContaining({ enabled: false }));
    });

    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockClear();

    fireEvent.click(screen.getByTestId('enable-all-db-targets'));

    await waitFor(() => {
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-A', expect.objectContaining({ enabled: true }));
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-B', expect.objectContaining({ enabled: true }));
    });
  });
});
