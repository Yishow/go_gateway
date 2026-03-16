import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import type { Device, Mapping, Point, Tag } from '../../../../types/datalink';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockModbusShareAPI,
  mockDBTargetAPI,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockModbusShareAPI: {
    status: vi.fn(),
    start: vi.fn(),
    stop: vi.fn(),
    listMappings: vi.fn(),
    upsertMapping: vi.fn(),
    deleteMapping: vi.fn(),
    writeTagValue: vi.fn(),
    sync: vi.fn(),
  },
  mockDBTargetAPI: {
    listConnectors: vi.fn(),
    getConnector: vi.fn(),
    createConnector: vi.fn(),
    updateConnector: vi.fn(),
    deleteConnector: vi.fn(),
    testConnector: vi.fn(),
    listTables: vi.fn(),
    validateConnector: vi.fn(),
    listMappings: vi.fn(),
    getMapping: vi.fn(),
    createMapping: vi.fn(),
    updateMapping: vi.fn(),
    deleteMapping: vi.fn(),
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../services/datalink', () => ({
  modbusShareAPI: mockModbusShareAPI,
  dbTargetAPI: mockDBTargetAPI,
}));

vi.mock('../../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockDevices,
    isLoading: false,
  }),
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
  }),
}));

vi.mock('../../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({
    data: mockMappings,
    isLoading: false,
  }),
}));

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchPage />
    </QueryClientProvider>,
  );
}

describe('DatalinkWorkbench output step', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    mockDevices.splice(0, mockDevices.length, {
      id: 'device-1',
      name: 'Mixer PLC',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    });

    mockPoints.splice(0, mockPoints.length, {
      id: 'point-1',
      device_id: 'device-1',
      name: 'Flow Sensor',
      description: '',
      data_type: 'int16',
      address: '40001',
      enabled: true,
      polling_group_id: '',
      last_value: 12,
      last_read_at: '',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    }, {
      id: 'point-2',
      device_id: 'device-1',
      name: 'Pressure Sensor',
      description: '',
      data_type: 'int16',
      address: '40002',
      enabled: true,
      polling_group_id: '',
      last_value: 28,
      last_read_at: '',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    });

    mockTags.splice(
      0,
      mockTags.length,
      {
        id: 'tag-1',
        key: 'TAG_40001',
        display_name: 'Flow Sensor',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'draft',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'tag-2',
        key: 'TAG_40002',
        display_name: 'Pressure Sensor',
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'draft',
        created_at: '',
        updated_at: '',
      },
    );

    mockMappings.splice(
      0,
      mockMappings.length,
      {
        id: 'mapping-1',
        point_id: 'point-1',
        tag_id: 'tag-1',
        enabled: true,
        transform_pipeline: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'mapping-2',
        point_id: 'point-2',
        tag_id: 'tag-2',
        enabled: true,
        transform_pipeline: '',
        created_at: '',
        updated_at: '',
      },
    );

    mockModbusShareAPI.status.mockResolvedValue({
      enabled: false,
      port: 5020,
      address: '',
      bind_state: 'fail',
      mapping_count: 0,
    });
    mockModbusShareAPI.listMappings.mockResolvedValue([]);
    mockModbusShareAPI.start.mockResolvedValue({
      enabled: true,
      port: 5030,
      address: '127.0.0.1:5030',
      bind_state: 'pass',
      mapping_count: 0,
    });
    mockModbusShareAPI.stop.mockResolvedValue({
      enabled: false,
      port: 0,
      address: '',
      bind_state: 'fail',
      mapping_count: 0,
    });
    mockModbusShareAPI.upsertMapping.mockResolvedValue({
      tag_id: 'tag-1',
      register: 12,
      data_type: 'int16',
      updated_at: '',
    });
    mockModbusShareAPI.sync.mockResolvedValue({
      updated: 1,
      skipped: 0,
      errors: [],
    });
    mockModbusShareAPI.writeTagValue.mockResolvedValue(undefined);

    mockDBTargetAPI.listConnectors.mockResolvedValue([
      {
        id: 'connector-1',
        name: 'Main SQLite',
        kind: 'sqlite',
        connection_config: {
          dsn: '/tmp/target.db',
        },
        status: 'ready',
        last_check_at: '',
        last_check_error: '',
        enabled: true,
        created_at: '',
        updated_at: '',
      },
    ]);
    mockDBTargetAPI.listMappings.mockResolvedValue([]);
    mockDBTargetAPI.listTables.mockResolvedValue([
      {
        schema: 'main',
        name: 'sensor_values',
        columns: [
          { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true },
          { name: 'value', data_type: 'real', nullable: false, primary_key: false },
        ],
      },
    ]);
    mockDBTargetAPI.validateConnector.mockResolvedValue({
      ready: true,
      issues: [],
    });
    mockDBTargetAPI.createConnector.mockResolvedValue({
      id: 'connector-1',
      name: 'Main SQLite',
      kind: 'sqlite',
      connection_config: {
        dsn: '/tmp/target.db',
      },
      status: 'ready',
      last_check_at: '',
      last_check_error: '',
      enabled: true,
      created_at: '',
      updated_at: '',
    });
    mockDBTargetAPI.updateConnector.mockResolvedValue({
      id: 'connector-1',
      name: 'Main SQLite',
      kind: 'sqlite',
      connection_config: {
        dsn: '/tmp/target.db',
      },
      status: 'ready',
      last_check_at: '',
      last_check_error: '',
      enabled: true,
      created_at: '',
      updated_at: '',
    });
    mockDBTargetAPI.createMapping.mockResolvedValue({
      id: 'db-mapping-1',
      tag_id: 'tag-1',
      connector_id: 'connector-1',
      table_schema: 'main',
      table_name: 'sensor_values',
      column_name: 'value',
      write_mode: 'insert',
      timestamp_column: null,
      enabled: true,
      created_at: '',
      updated_at: '',
    });
  });

  it('shows an empty state when the selected device has no linked tags', async () => {
    mockMappings.splice(0, mockMappings.length);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByText('workbench.output.empty.title')).toBeInTheDocument();
  });

  it('starts the local modbus server with the specified port', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.tag');
    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: 'workbench.output.actions.startServer' }),
      ).toBeEnabled();
    });

    fireEvent.change(screen.getByLabelText('workbench.output.server.port'), {
      target: { value: '5030' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.output.actions.startServer' }));

    await waitFor(() => {
      expect(mockModbusShareAPI.start).toHaveBeenCalledWith(5030);
    });
  });

  it('binds a linked tag to a local modbus register', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.tag');
    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: 'workbench.output.actions.bind' }),
      ).toBeEnabled();
    });

    fireEvent.change(screen.getByLabelText('workbench.output.mapping.register'), {
      target: { value: '12' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.output.actions.bind' }));

    await waitFor(() => {
      expect(mockModbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-1', 12);
    });
  });

  it('keeps candidate selection and register input in sync', async () => {
    mockModbusShareAPI.listMappings.mockResolvedValue([
      {
        tag_id: 'tag-1',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
      {
        tag_id: 'tag-2',
        register: 24,
        data_type: 'int16',
        updated_at: '',
      },
    ]);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.tag');
    await waitFor(() => {
      expect(screen.getByDisplayValue('12')).toBeInTheDocument();
    });

    const firstCandidate = screen.getAllByRole('button', { name: /TAG_40001/ })[0];
    const secondCandidate = screen.getAllByRole('button', { name: /TAG_40002/ })[0];

    expect(firstCandidate).toHaveAttribute('aria-pressed', 'true');
    expect(secondCandidate).toHaveAttribute('aria-pressed', 'false');

    fireEvent.click(secondCandidate);

    await waitFor(() => {
      expect(screen.getByDisplayValue('24')).toBeInTheDocument();
    });
    expect(firstCandidate).toHaveAttribute('aria-pressed', 'false');
    expect(secondCandidate).toHaveAttribute('aria-pressed', 'true');
  });

  it('blocks sync when register conflicts exist', async () => {
    mockModbusShareAPI.status.mockResolvedValue({
      enabled: true,
      port: 5020,
      address: '127.0.0.1:5020',
      bind_state: 'pass',
      mapping_count: 2,
    });
    mockModbusShareAPI.listMappings.mockResolvedValue([
      {
        tag_id: 'tag-1',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
      {
        tag_id: 'tag-2',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
    ]);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByRole('button', { name: 'workbench.output.actions.sync' })).toBeDisabled();
    expect(screen.getByRole('status')).toHaveTextContent('workbench.output.conflicts.summary');
  });

  it('binds a linked tag to a database target mapping', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.database.mapping.table');
    await waitFor(() => {
      expect(
        screen.getByRole('button', {
          name: 'workbench.output.database.actions.saveMapping',
        }),
      ).toBeEnabled();
    });

    fireEvent.change(screen.getByLabelText('workbench.output.database.mapping.table'), {
      target: { value: 'main.sensor_values' },
    });
    fireEvent.change(screen.getByLabelText('workbench.output.database.mapping.column'), {
      target: { value: 'value' },
    });
    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.database.actions.saveMapping',
      }),
    );

    await waitFor(() => {
      expect(mockDBTargetAPI.createMapping).toHaveBeenCalledWith({
        tag_id: 'tag-1',
        connector_id: 'connector-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        column_name: 'value',
        write_mode: 'insert',
        timestamp_column: '',
        enabled: true,
      });
    });
  });
});
