import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ state }: { state: { devices: Array<{ running?: boolean }>; rules: unknown[]; points: unknown[]; mappings: Record<string, unknown>; db: { targets: Record<string, unknown> } } }) => (
    <div
      data-testid="boot-shell"
      data-device-count={state.devices.length}
      data-rule-count={state.rules.length}
      data-point-count={state.points.length}
      data-mapping-count={Object.keys(state.mappings).length}
      data-target-count={Object.keys(state.db.targets).length}
    >
      boot-shell
    </div>
  ),
}));

vi.mock('../../../src/services/studioV2Workspace', () => ({
  studioV2WorkspaceAPI: {
    get: vi.fn(),
  },
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
      queries: {
        retry: false,
      },
    },
  });

  render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

describe('DatalinkWorkbenchV2Page workspace bootstrap', () => {
  it('boots from the workspace API on first open', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-1',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    expect(screen.getByTestId('workbench-v2-bootstrap-loading')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    expect(studioV2WorkspaceAPI.get).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-id', 'workspace-1');
    expect(screen.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-status', 'empty');
    expect(screen.getByTestId('boot-shell')).toBeInTheDocument();
  });

  it('clears seeded local draft data when the workspace has no devices or rules', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-empty',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-device-count', '0');
    });

    expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-rule-count', '0');
    expect(screen.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-id', 'workspace-empty');
  });

  it('shows an actionable bootstrap error when the workspace API fails', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValueOnce(new Error('db offline'));
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-bootstrap-error')).toBeInTheDocument();
    });

    expect(screen.getByText('無法載入 Studio V2 工作區。請檢查後端資料庫連線後重試。')).toBeInTheDocument();
  });

  it('shows bootstrap error when the mappings API fails', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-mapping-error',
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
        status: 'active',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        availability_status: 'available',
        running: true,
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockRejectedValueOnce(new Error('mapping db offline'));
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-bootstrap-error')).toBeInTheDocument();
    });

    expect(screen.getByText('無法載入 Studio V2 工作區。請檢查後端資料庫連線後重試。')).toBeInTheDocument();
  });

  it('hydrates persisted cross-step setup before shell render', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-persisted',
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
        workspace_id: 'workspace-persisted',
        start_address: '40001',
        count: 1,
        data_type: 'int16',
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
        workspace_id: 'workspace-persisted',
        point_id: 'persisted-point-A',
        rule_id: 'rule-A',
        device_id: 'dev-A',
        address: '40001',
        tag_id: 'tag-A',
        tag_key: 'line.a.persisted',
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
    ]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue({
      id: 'db-1',
      workspace_id: 'workspace-persisted',
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
        id: 'target-A',
        workspace_id: 'workspace-persisted',
        point_id: 'rule-A-p-0',
        tag_id: 'tag-A',
        column_name: 'line_a',
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'target-stale',
        workspace_id: 'workspace-persisted',
        point_id: 'missing-point',
        tag_id: 'tag-stale',
        column_name: 'stale_column',
        enabled: true,
        save_state: 'saved',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ] as any);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-device-count', '1');
    });

    expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-rule-count', '1');
    expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-point-count', '1');
    expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-mapping-count', '1');
    expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-target-count', '1');
  });
  it('shows an unrecovered draft warning after reload when prior local draft was not saved', async () => {
    window.sessionStorage.setItem('wbv2_unrecovered_draft', '1');

    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-warning',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-draft-loss-warning')).toBeInTheDocument();
    });

    expect(screen.getByText('偵測到上次重新整理前有未保存的本地草稿，系統已回復為最後一次成功保存的設定。')).toBeInTheDocument();
    expect(window.sessionStorage.getItem('wbv2_unrecovered_draft')).toBeNull();
  });
});
