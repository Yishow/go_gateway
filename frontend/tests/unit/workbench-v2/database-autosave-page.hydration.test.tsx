import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { resetDatabaseAutosaveMocks } from './database-autosave-page.testHarness';

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
  it('hydrates persisted database config and targets onto the same point rows', async () => {
    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
      expect(screen.getByTestId('point-count')).toHaveTextContent('2');
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));

    await waitFor(() => {
      expect(screen.getByTestId('connector-table')).toHaveTextContent('sensor_values');
      expect(screen.getByTestId('target-column-point-A')).toHaveTextContent('line_a');
      expect(screen.getByTestId('target-column-point-B')).toHaveTextContent('line_b');
    });
  });
});
