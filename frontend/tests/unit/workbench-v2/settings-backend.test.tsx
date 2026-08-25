import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { SettingsPage } from '../../../src/features/datalink/workbench-v2/settings/SettingsPage';
import { INITIAL_STATE, workbenchV2Reducer } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, fallback?: string) => fallback || key,
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  settingsAPI: {
    listItems: vi.fn(),
    updateKey: vi.fn(),
  },
  dbTargetAPI: {
    listConnectors: vi.fn(),
    createConnector: vi.fn(),
    updateConnector: vi.fn(),
    deleteConnector: vi.fn(),
    testConnector: vi.fn(),
  },
}));

import { settingsAPI, dbTargetAPI } from '../../../src/services/datalink';

function SettingsHarness() {
  const [state, dispatch] = React.useReducer(workbenchV2Reducer, {
    ...INITIAL_STATE,
    view: 'settings' as const,
  });
  return <SettingsPage state={state} dispatch={dispatch} />;
}

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  render(
    <QueryClientProvider client={queryClient}>
      <SettingsHarness />
    </QueryClientProvider>,
  );
}

describe('SettingsPage backend wiring', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(settingsAPI.listItems).mockResolvedValue([
      {
        key: 'write_precision',
        value: 'second',
        description: '',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        key: 'partition_interval',
        value: 'weekly',
        description: '',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        key: 'batch_size',
        value: 777,
        description: '',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        key: 'retention_days',
        value: 180,
        description: '',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        key: 'default_retry_count',
        value: 5,
        description: '',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        key: 'api_base',
        value: 'http://backend-gateway:9090',
        description: '',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        key: 'modbus_share',
        value: {
          enabled: false,
          bind_address: '127.0.0.1',
          port: 15020,
          slave_id: 9,
          capacity_registers: 4096,
          settings_revision: 'settings-rev-7',
        },
        description: '',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([
      {
        id: 'conn-backend',
        name: 'Backend SQLite',
        kind: 'sqlite',
        connection_config: {
          path: '/tmp/backend.db',
          table: 'sensor_values',
        },
        status: 'ready',
        last_check_at: '2026-05-30T00:00:00Z',
        last_check_error: '',
        enabled: true,
        default_write_interval_seconds: 8,
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
  });

  it('boots from backend instead of silently using local defaults', async () => {
    renderPage();

    expect(screen.getByTestId('settings-backend-loading')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByDisplayValue('/tmp/backend.db')).toBeInTheDocument();
    });

    expect(screen.getByDisplayValue('http://backend-gateway:9090')).toBeInTheDocument();
    expect(
      screen.getByText('分區間隔').closest('div')?.querySelector('select'),
    ).toHaveValue('weekly');
    expect(screen.getByDisplayValue('777')).toBeInTheDocument();
    expect(screen.getByTestId('connector-row-conn-backend')).toBeInTheDocument();
  });

  it('persists listener settings as one backend payload when save is clicked', async () => {
    vi.mocked(settingsAPI.updateKey).mockImplementation(async (key, value) => ({
      key,
      value,
      description: '',
      updated_at: '2026-05-30T00:01:00Z',
    }));

    renderPage();

    await waitFor(() => {
      expect(screen.getByDisplayValue('777')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByDisplayValue('777'), { target: { value: '888' } });
    fireEvent.change(screen.getByDisplayValue('http://backend-gateway:9090'), {
      target: { value: 'http://backend-gateway:9191' },
    });
    fireEvent.click(screen.getByText('儲存所有設定'));

    await waitFor(() => {
      expect(settingsAPI.updateKey).toHaveBeenCalledWith('batch_size', 888);
    });
    expect(settingsAPI.updateKey).toHaveBeenCalledWith('api_base', 'http://backend-gateway:9191');
    expect(settingsAPI.updateKey).toHaveBeenCalledWith('modbus_share', {
      enabled: false,
      bind_address: '127.0.0.1',
      port: 15020,
      slave_id: 9,
      capacity_registers: 4096,
      settings_revision: 'settings-rev-7',
      expected_settings_revision: 'settings-rev-7',
    });
    expect(settingsAPI.updateKey).not.toHaveBeenCalledWith('modbus_share_enabled', expect.anything());
    expect(settingsAPI.updateKey).not.toHaveBeenCalledWith('modbus_share_bind_address', expect.anything());
    expect(settingsAPI.updateKey).not.toHaveBeenCalledWith('modbus_share_port', expect.anything());
    expect(settingsAPI.updateKey).not.toHaveBeenCalledWith('modbus_share_slave_id', expect.anything());
    expect(settingsAPI.updateKey).not.toHaveBeenCalledWith('modbus_share_base_register', expect.anything());
  });

  it('round-trips connector create, update, delete, and test through backend APIs', async () => {
    vi.mocked(dbTargetAPI.createConnector).mockResolvedValue({
      id: 'conn-new',
      name: '新連線 2',
      kind: 'postgres',
      connection_config: {
        host: '127.0.0.1',
        port: 5432,
        database: 'gateway_metrics',
        user: 'postgres',
        schema: 'public',
        table: 'sensor_readings',
      },
      status: 'ready',
      last_check_at: '2026-05-30T00:02:00Z',
      last_check_error: '',
      enabled: true,
      default_write_interval_seconds: 5,
      created_at: '2026-05-30T00:02:00Z',
      updated_at: '2026-05-30T00:02:00Z',
    });
    vi.mocked(dbTargetAPI.updateConnector).mockResolvedValue({
      id: 'conn-backend',
      name: 'Backend SQLite',
      kind: 'sqlite',
      connection_config: {
        path: '/tmp/backend-v2.db',
        table: 'sensor_values',
      },
      status: 'ready',
      last_check_at: '2026-05-30T00:03:00Z',
      last_check_error: '',
      enabled: true,
      default_write_interval_seconds: 8,
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:03:00Z',
    });
    vi.mocked(dbTargetAPI.testConnector).mockResolvedValue({
      id: 'conn-backend',
      name: 'Backend SQLite',
      kind: 'sqlite',
      connection_config: {
        path: '/tmp/backend-v2.db',
        table: 'sensor_values',
      },
      status: 'unreachable',
      last_check_at: '2026-05-30T00:04:00Z',
      last_check_error: 'connection refused',
      enabled: true,
      default_write_interval_seconds: 8,
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:04:00Z',
    });
    vi.mocked(dbTargetAPI.deleteConnector).mockResolvedValue();

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('connector-row-conn-backend')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('新增連接器'));
    await waitFor(() => {
      expect(dbTargetAPI.createConnector).toHaveBeenCalled();
    });

    const sqlitePathInput = screen.getByDisplayValue('/tmp/backend.db');
    fireEvent.change(sqlitePathInput, { target: { value: '/tmp/backend-v2.db' } });
    await waitFor(() => {
      expect(dbTargetAPI.updateConnector).toHaveBeenCalledWith(
        'conn-backend',
        expect.objectContaining({
          connection_config: expect.objectContaining({
            path: '/tmp/backend-v2.db',
          }),
        }),
      );
    });

    const backendRow = screen.getByTestId('connector-row-conn-backend');
    fireEvent.click(
      backendRow.querySelectorAll('button')[1] as HTMLButtonElement,
    );
    await waitFor(() => {
      expect(dbTargetAPI.testConnector).toHaveBeenCalledWith('conn-backend');
    });
    await waitFor(() => {
      expect(screen.getByText('connection refused')).toBeInTheDocument();
    });

    fireEvent.click(
      backendRow.querySelectorAll('button')[2] as HTMLButtonElement,
    );
    await waitFor(() => {
      expect(dbTargetAPI.deleteConnector).toHaveBeenCalledWith('conn-backend');
    });
  });
});
