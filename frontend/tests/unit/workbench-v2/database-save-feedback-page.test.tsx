import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { resetDatabaseAutosaveMocks, runtimeApplied, sqliteConfigFixture, studioV2WorkspaceAPI, studioV2WorkspaceDatabaseAPI } from './database-autosave-page.testHarness';
beforeEach(resetDatabaseAutosaveMocks);
async function mount() {
  const { default: Page } = await import('@/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page');
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><Page /></QueryClientProvider>);
  return client;
}
it('retains a confirmed same-workspace saved editor after failed bootstrap refresh and preserves later editing', async () => {
  await mount();
  await screen.findByTestId('database-autosave-shell');
  vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValue(new Error('private backend detail'));
  vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockResolvedValue(runtimeApplied(sqliteConfigFixture({ table: 'sensor_values_v2', identity_revision: 'identity-2', setup_revision: 'setup-2', updated_at: 'later' })));
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(sqliteConfigFixture({ table: 'sensor_values_v2', identity_revision: 'identity-2', setup_revision: 'setup-2', updated_at: 'later' }));
  fireEvent.click(screen.getByTestId('make-connector-valid'));
  await waitFor(() => expect(screen.getByTestId('connector-save-state')).toHaveTextContent('saved'));
  expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
  expect(screen.getByTestId('recording-draft-refresh-failed')).toBeInTheDocument();
  expect(document.body.textContent).not.toContain('private backend detail');
  let finish!: (value: ReturnType<typeof runtimeApplied<ReturnType<typeof sqliteConfigFixture>>>) => void;
  vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockReturnValueOnce(new Promise(resolve => { finish = resolve; }));
  fireEvent.click(screen.getByTestId('make-connector-valid'));
  await waitFor(() => expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledTimes(2));
  expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
  await act(async () => { finish(runtimeApplied(sqliteConfigFixture({ identity_revision: 'identity-3', setup_revision: 'setup-3' }))); });
  await waitFor(() => expect(screen.getByTestId('connector-save-state')).toHaveTextContent('saved'));
});
it('does not retain an editor for malformed initial bootstrap', async () => {
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(undefined as never);
  const client = await mount();
  await act(async () => {});
  await waitFor(() => expect(client.isFetching()).toBe(0));
  expect(screen.queryByTestId('database-autosave-shell')).not.toBeInTheDocument();
  expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
});
