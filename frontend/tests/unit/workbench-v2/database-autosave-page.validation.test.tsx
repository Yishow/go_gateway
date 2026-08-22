import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  resetDatabaseAutosaveMocks,
  runtimeApplied,
  sqliteConfigFixture,
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
  it('saves one valid connector edit and keeps invalid edits local', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockResolvedValueOnce(
      runtimeApplied(sqliteConfigFixture({ table: 'sensor_values_v2' })),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('make-connector-valid'));

    await waitFor(() => {
      expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledWith(expect.objectContaining({
        table: 'sensor_values_v2',
      }));
      expect(screen.getByTestId('connector-save-state')).toHaveTextContent('saved');
    });

    const savedCallCount = vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mock.calls.length;

    fireEvent.click(screen.getByTestId('make-connector-invalid'));

    await waitFor(() => {
      expect(screen.getByTestId('connector-table')).toHaveTextContent('');
      expect(screen.getByTestId('connector-save-state')).toHaveTextContent('draft-invalid');
      expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(savedCallCount);
    });
  });
});
