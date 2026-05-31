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
      <div data-testid="connector-table">{state.db.connector.table}</div>
      <div data-testid="connector-save-state">{state.db.connector.save_state}</div>
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
        point_id: 'point-A',
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
        point_id: 'point-B',
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
        point_id: 'point-A',
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
        point_id: 'point-B',
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
      if (pointId === 'point-A') {
        throw new Error('找不到資料欄位: missing_column');
      }
      return {
        id: 'row-B',
        workspace_id: 'workspace-1',
        point_id: 'point-B',
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
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('point-A', expect.objectContaining({
        column_name: 'missing_column',
        enabled: true,
      }));
      expect(studioV2WorkspaceDatabaseAPI.upsertTarget).toHaveBeenCalledWith('point-B', expect.objectContaining({
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
});
