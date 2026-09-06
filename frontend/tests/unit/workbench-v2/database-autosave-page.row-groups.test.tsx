import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { StudioV2WorkspaceDatabaseRowGroupRecord } from '../../../src/types/datalink';
import {
  dbTargetFixture,
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
  it('isolates one failing target row from another valid row', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockImplementation(async (pointId, request) => {
      if (pointId === 'persisted-point-A') {
        throw new Error('找不到資料欄位: missing_column');
      }
      return runtimeApplied(dbTargetFixture({
        point_id: 'persisted-point-B',
        tag_id: 'tag-B',
        column_name: request.column_name,
        enabled: request.enabled,
      }));
    });

    await renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));

    await waitFor(() => {
      expect(screen.getByTestId('target-column-point-A')).toHaveTextContent('line_a');
    });

    fireEvent.click(screen.getByTestId('make-point-a-fail'));
    fireEvent.click(screen.getByTestId('make-point-b-valid'));

    await waitFor(() => {
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-A', expect.objectContaining({
        column_name: 'missing_column',
        enabled: true,
      }));
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-B', expect.objectContaining({
        column_name: 'line_b_saved',
        enabled: true,
      }));
    });

    await waitFor(() => {
      expect(screen.getByTestId('target-save-state-point-A')).toHaveTextContent('save-error');
      expect(screen.getByTestId('target-save-error-point-A')).toHaveTextContent('errors.autosave_failed');
      expect(screen.getByTestId('target-save-state-point-B')).toHaveTextContent('saved');
      expect(screen.getByTestId('target-save-error-point-B')).toHaveTextContent('');
    });
  });

  it('waits for row group config save before saving target row-group refs', async () => {
    let resolveConfigSave: ((value: Awaited<ReturnType<typeof studioV2WorkspaceDatabaseAPI.updateConfig>>) => void) | undefined;
    const savedConfig = runtimeApplied(sqliteConfigFixture());
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockImplementation((request) => {
      if (request.row_groups?.some((group: StudioV2WorkspaceDatabaseRowGroupRecord) => group.id === 'row-group-delayed')) {
        return new Promise((resolve) => {
          resolveConfigSave = resolve;
        });
      }
      return Promise.resolve({ ...savedConfig, row_groups: request.row_groups ?? [] });
    });
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockResolvedValue(
      runtimeApplied(dbTargetFixture({
        point_id: 'persisted-point-A',
        row_group_id: 'row-group-delayed',
      })),
    );

    await renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
      expect(screen.getByTestId('point-count')).toHaveTextContent('2');
    });
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockClear();
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockClear();

    fireEvent.click(screen.getByTestId('create-row-group'));

    await waitFor(() => {
      expect(screen.getByTestId('row-group-count')).toHaveTextContent('1');
    });

    await waitFor(() => {
      expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledWith(expect.objectContaining({
        row_groups: [expect.objectContaining({ id: 'row-group-delayed' })],
      }));
    });
    expect(studioV2WorkspaceDatabaseAPI.upsertTarget).not.toHaveBeenCalled();

    resolveConfigSave?.({
      ...savedConfig,
      row_groups: [{
        id: 'row-group-delayed',
        connector_id: 'db-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        member_point_ids: ['rule-A-p-0', 'rule-A-p-1'],
        group_key_columns: ['ts'],
      }],
    });

    await waitFor(() => {
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-A', expect.objectContaining({
        row_group_id: 'row-group-delayed',
      }));
    });
  });

  it('clears stale row groups before saving a table change and defers target updates until the connector save completes', async () => {
    let resolveScopeChangeSave: ((value: Awaited<ReturnType<typeof studioV2WorkspaceDatabaseAPI.updateConfig>>) => void) | undefined;
    const buildSavedConfig = (
      table: string,
      rowGroups: StudioV2WorkspaceDatabaseRowGroupRecord[] = [],
    ) => runtimeApplied(sqliteConfigFixture({ table, row_groups: rowGroups }));

    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockImplementation((request) => {
      if (request.table === 'sensor_values_v2') {
        return new Promise((resolve) => {
          resolveScopeChangeSave = resolve;
        });
      }
      return Promise.resolve(buildSavedConfig(request.table, request.row_groups ?? []));
    });
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockImplementation(async (pointId, request) =>
      runtimeApplied(dbTargetFixture({
        id: pointId === 'persisted-point-A' ? 'row-A' : 'row-B',
        point_id: pointId,
        tag_id: pointId === 'persisted-point-A' ? 'tag-A' : 'tag-B',
        column_name: request.column_name,
        enabled: request.enabled,
        row_group_id: request.row_group_id,
      })));

    await renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('database-autosave-shell')).toBeInTheDocument();
      expect(screen.getByTestId('point-count')).toHaveTextContent('2');
    });

    fireEvent.click(screen.getByTestId('prime-mappings'));

    await waitFor(() => {
      expect(screen.getByTestId('target-column-point-A')).toHaveTextContent('line_a');
      expect(screen.getByTestId('target-column-point-B')).toHaveTextContent('line_b');
    });

    fireEvent.click(screen.getByTestId('create-row-group'));

    await waitFor(() => {
      expect(screen.getByTestId('row-group-count')).toHaveTextContent('1');
      expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledWith(expect.objectContaining({
        row_groups: [expect.objectContaining({ id: 'row-group-delayed' })],
      }));
    });

    await waitFor(() => {
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-A', expect.objectContaining({
        row_group_id: 'row-group-delayed',
      }));
    });

    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockClear();
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockClear();

    fireEvent.click(screen.getByTestId('change-scope-with-row-group-reset'));

    await waitFor(() => {
      expect(screen.getByTestId('row-group-count')).toHaveTextContent('0');
      expect(studioV2WorkspaceDatabaseAPI.updateConfig).toHaveBeenCalledWith(expect.objectContaining({
        table: 'sensor_values_v2',
        row_groups: [],
      }));
    });

    expect(studioV2WorkspaceDatabaseAPI.upsertTarget).not.toHaveBeenCalled();

    resolveScopeChangeSave?.(buildSavedConfig('sensor_values_v2'));

    await waitFor(() => {
      expect(screen.getByTestId('connector-table')).toHaveTextContent('sensor_values_v2');
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-A', expect.objectContaining({
        row_group_id: undefined,
      }));
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('persisted-point-B', expect.objectContaining({
        row_group_id: undefined,
      }));
    });
  });
});
