import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';

const dbPoints = [
  {
    id: 'point-A',
    device_id: 'dev-A',
    rule_id: 'rule-A',
    rule_name: 'Line A Registers',
    name: 'A_0',
    address: '40001',
    data_type: 'float64',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
  {
    id: 'point-B',
    device_id: 'dev-A',
    rule_id: 'rule-A',
    rule_name: 'Line A Registers',
    name: 'B_0',
    address: '40002',
    data_type: 'float64',
    function: 'holding_register',
    width: 1,
    enabled: true,
    skipped: false,
    _rule_scale: 1,
    _rule_offset: 0,
  },
];

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ state, actions }: any) => (
    <div data-testid="database-autosave-shell">
      <button
        type="button"
        data-testid="prime-mappings"
        onClick={() => actions.dispatch({ type: 'initMappingsForPoints', points: dbPoints })}
      />
      <button
        type="button"
        data-testid="make-connector-valid"
        onClick={() => actions.dispatch({ type: 'updateDbConnector', patch: { table: 'sensor_values_v2' } })}
      />
      <button
        type="button"
        data-testid="make-connector-invalid"
        onClick={() => actions.dispatch({ type: 'updateDbConnector', patch: { table: '' } })}
      />
      <button
        type="button"
        data-testid="make-point-a-fail"
        onClick={() => actions.dispatch({ type: 'updateDbTarget', pointId: 'point-A', patch: { column_name: 'missing_column' } })}
      />
      <button
        type="button"
        data-testid="make-point-b-valid"
        onClick={() => actions.dispatch({ type: 'updateDbTarget', pointId: 'point-B', patch: { column_name: 'line_b_saved' } })}
      />
      <button
        type="button"
        data-testid="create-row-group"
        onClick={() => {
          const pointIds = state.points.slice(0, 2).map((point: { id: string }) => point.id);
          actions.dispatch({
            type: 'setDbRowGroups',
            rowGroups: [{
              id: 'row-group-delayed',
              connector_id: state.db.connector.connector_id,
              table_schema: state.db.connector.schema,
              table_name: state.db.connector.table,
              member_point_ids: pointIds,
              group_key_columns: [state.db.connector.timestamp_column || 'ts'],
            }],
          });
          pointIds.forEach((pointId: string) => {
            actions.dispatch({ type: 'updateDbTarget', pointId, patch: { row_group_id: 'row-group-delayed' } });
          });
        }}
      />
      <button
        type="button"
        data-testid="change-scope-with-row-group-reset"
        onClick={() => {
          const pointIds = state.points.slice(0, 2).map((point: { id: string }) => point.id);
          actions.dispatch({ type: 'setDbRowGroups', rowGroups: [] });
          pointIds.forEach((pointId: string) => {
            actions.dispatch({ type: 'updateDbTarget', pointId, patch: { row_group_id: undefined } });
          });
          actions.dispatch({ type: 'updateDbConnector', patch: { table: 'sensor_values_v2' } });
        }}
      />
      <div data-testid="connector-table">{state.db.connector.table}</div>
      <div data-testid="connector-save-state">{state.db.connector.save_state}</div>
      <div data-testid="point-count">{state.points.length}</div>
      <div data-testid="row-group-count">{state.db.row_groups?.length ?? 0}</div>
      <div data-testid="target-column-point-A">{state.db.targets['point-A']?.column_name ?? ''}</div>
      <div data-testid="target-save-state-point-A">{state.db.targets['point-A']?.save_state ?? ''}</div>
      <div data-testid="target-save-error-point-A">{state.db.targets['point-A']?.save_error ?? ''}</div>
      <div data-testid="target-column-point-B">{state.db.targets['point-B']?.column_name ?? ''}</div>
      <div data-testid="target-save-state-point-B">{state.db.targets['point-B']?.save_state ?? ''}</div>
      <div data-testid="target-save-error-point-B">{state.db.targets['point-B']?.save_error ?? ''}</div>
    </div>
  ),
}));

vi.mock('../../../src/services/studioV2Workspace', () => ({
  studioV2WorkspaceAPI: { get: vi.fn() },
}));

vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({
  studioV2WorkspaceDevicesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
    updateOrder: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2Rules', () => ({
  studioV2RulesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2Mappings', () => ({
  studioV2MappingsAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
}));

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

describe('DatalinkWorkbenchV2Page database autosave orchestration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-A'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-A',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([
      {
        id: 'rule-A',
        device_id: 'dev-A',
        workspace_id: 'workspace-1',
        start_address: '40001',
        count: 2,
        data_type: 'float64',
        naming_prefix: 'LINE_',
        enabled: true,
        locked: false,
        origin: 'manual',
        skipped_addresses: [],
        revision_id: 'rev-A',
        scale_multiplier: 1,
        scale_offset: 0,
        data_format: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([
      {
        id: 'mapping-A',
        workspace_id: 'workspace-1',
        point_id: 'persisted-point-A',
        rule_id: 'rule-A',
        device_id: 'dev-A',
        address: '40001',
        tag_id: 'tag-A',
        tag_key: 'line.a.temp',
        display_name: 'Line A Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'mapping-B',
        workspace_id: 'workspace-1',
        point_id: 'persisted-point-B',
        rule_id: 'rule-A',
        device_id: 'dev-A',
        address: '40002',
        tag_id: 'tag-B',
        tag_key: 'line.b.temp',
        display_name: 'Line B Temp',
        unit: 'C',
        target_type: 'float64',
        scale: 1,
        offset: 0,
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue({
      id: 'db-1',
      workspace_id: 'workspace-1',
      kind: 'sqlite',
      name: 'Line A SQLite',
      host: '',
      port: 0,
      database: '/tmp/target.db',
      username: '',
      schema: 'main',
      table: 'sensor_values',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    } as any);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([
      {
        id: 'row-A',
        workspace_id: 'workspace-1',
        point_id: 'persisted-point-A',
        tag_id: 'tag-A',
        column_name: 'line_a',
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'row-B',
        workspace_id: 'workspace-1',
        point_id: 'persisted-point-B',
        tag_id: 'tag-B',
        column_name: 'line_b',
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ] as any);
  });

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

  it('saves one valid connector edit and keeps invalid edits local', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockResolvedValueOnce({
      id: 'db-1',
      workspace_id: 'workspace-1',
      kind: 'sqlite',
      name: 'Line A SQLite',
      host: '',
      port: 0,
      database: '/tmp/target.db',
      username: '',
      schema: 'main',
      table: 'sensor_values_v2',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    } as any);

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

  it('isolates one failing target row from another valid row', async () => {
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockImplementation(async (pointId, request) => {
      if (pointId === 'persisted-point-A') {
        throw new Error('找不到資料欄位: missing_column');
      }
      return {
        id: 'row-B',
        workspace_id: 'workspace-1',
        point_id: 'persisted-point-B',
        tag_id: 'tag-B',
        column_name: request.column_name,
        enabled: request.enabled,
        save_state: 'saved',
        runtime_apply_status: 'not_running',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      } as any;
    });

    renderPage();

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
      expect(screen.getByTestId('target-save-error-point-A')).toHaveTextContent('找不到資料欄位: missing_column');
      expect(screen.getByTestId('target-save-state-point-B')).toHaveTextContent('saved');
      expect(screen.getByTestId('target-save-error-point-B')).toHaveTextContent('');
    });
  });

  it('waits for row group config save before saving target row-group refs', async () => {
    let resolveConfigSave: ((value: Awaited<ReturnType<typeof studioV2WorkspaceDatabaseAPI.updateConfig>>) => void) | undefined;
    const savedConfig: Awaited<ReturnType<typeof studioV2WorkspaceDatabaseAPI.updateConfig>> = {
      id: 'db-1',
      workspace_id: 'workspace-1',
      kind: 'sqlite',
      name: 'Line A SQLite',
      host: '',
      port: 0,
      database: '/tmp/target.db',
      username: '',
      schema: 'main',
      table: 'sensor_values',
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    };
    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockImplementation((request) => {
      if (request.row_groups?.some((group) => group.id === 'row-group-delayed')) {
        return new Promise((resolve) => {
          resolveConfigSave = resolve;
        });
      }
      return Promise.resolve({ ...savedConfig, row_groups: request.row_groups ?? [] });
    });
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockResolvedValue({
      id: 'row-A',
      workspace_id: 'workspace-1',
      point_id: 'persisted-point-A',
      tag_id: 'tag-A',
      column_name: 'line_a',
      enabled: true,
      row_group_id: 'row-group-delayed',
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });

    renderPage();

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
    const buildSavedConfig = (table: string, rowGroups: Array<Record<string, unknown>> = []) => ({
      id: 'db-1',
      workspace_id: 'workspace-1',
      kind: 'sqlite',
      name: 'Line A SQLite',
      host: '',
      port: 0,
      database: '/tmp/target.db',
      username: '',
      schema: 'main',
      table,
      write_mode: 'insert',
      write_interval_seconds: 5,
      timestamp_column: 'ts',
      status: 'ready',
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      row_groups: rowGroups,
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });

    vi.mocked(studioV2WorkspaceDatabaseAPI.updateConfig).mockImplementation((request) => {
      if (request.table === 'sensor_values_v2') {
        return new Promise((resolve) => {
          resolveScopeChangeSave = resolve;
        });
      }
      return Promise.resolve(buildSavedConfig(request.table, request.row_groups ?? []) as any);
    });
    vi.mocked(studioV2WorkspaceDatabaseAPI.upsertTarget).mockImplementation(async (pointId, request) => ({
      id: pointId === 'persisted-point-A' ? 'row-A' : 'row-B',
      workspace_id: 'workspace-1',
      point_id: pointId,
      tag_id: pointId === 'persisted-point-A' ? 'tag-A' : 'tag-B',
      column_name: request.column_name,
      enabled: request.enabled,
      row_group_id: request.row_group_id,
      save_state: 'saved',
      runtime_apply_status: 'not_running',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    }) as any);

    renderPage();

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

    resolveScopeChangeSave?.(buildSavedConfig('sensor_values_v2') as any);

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
