import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
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
  useCreateDeviceMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useUpdateDeviceMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useTestConnectionMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
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

    await screen.findByLabelText('workbench.output.mapping.register');
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

  it('renders only the active target mapping field inside the shared candidate rows', async () => {
    mockModbusShareAPI.listMappings.mockResolvedValue([
      {
        tag_id: 'tag-1',
        register: 12,
        data_type: 'int16',
        updated_at: '',
      },
    ]);
    mockDBTargetAPI.listMappings.mockResolvedValue([
      {
        id: 'db-mapping-2',
        tag_id: 'tag-2',
        connector_id: 'connector-1',
        table_schema: 'main',
        table_name: 'sensor_values',
        column_name: 'value',
        write_mode: 'insert',
        timestamp_column: null,
        enabled: true,
        created_at: '',
        updated_at: '',
      },
    ]);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.modbus',
      }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    ).toBeInTheDocument();
    expect(screen.getByTestId('output-candidate-tag-1')).toHaveTextContent('TAG_40001');
    expect(
      within(screen.getByTestId('output-candidate-tag-1')).getByTestId(
        'output-modbus-status-tag-1',
      ),
    ).toHaveTextContent('HR12');
    expect(
      within(screen.getByTestId('output-candidate-tag-1')).queryByTestId(
        'output-db-status-tag-1',
      ),
    ).not.toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    expect(
      within(screen.getByTestId('output-candidate-tag-2')).getByTestId(
        'output-db-status-tag-2',
      ),
    ).toHaveTextContent(
      'main.sensor_values.value',
    );
    expect(
      within(screen.getByTestId('output-candidate-tag-2')).queryByTestId(
        'output-modbus-status-tag-2',
      ),
    ).not.toBeInTheDocument();

    expect(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    ).toHaveAttribute('aria-pressed', 'true');
    await screen.findByLabelText('workbench.output.database.mapping.table');
  });

  it('marks modbus operational panels as supporting sections', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByTestId('modbus-secondary-panels')).toHaveAttribute(
      'data-emphasis',
      'supporting',
    );
  });

  it('uses the shared candidate board as the only tag selection surface for both targets', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByTestId('output-candidate-tag-1');

    expect(screen.queryByLabelText('workbench.output.mapping.tag')).not.toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    await screen.findByLabelText('workbench.output.database.mapping.table');

    expect(
      screen.queryByLabelText('workbench.output.database.mapping.tag'),
    ).not.toBeInTheDocument();
  });

  it('binds a linked tag to a local modbus register', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.register');
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

    await screen.findByLabelText('workbench.output.mapping.register');
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
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

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

  it('keeps the database mapping form aligned with the shared output candidate selection', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    await screen.findByLabelText('workbench.output.database.mapping.table');
    fireEvent.click(screen.getByTestId('output-candidate-tag-2'));

    await waitFor(() => {
      expect(screen.getByTestId('database-selected-tag')).toHaveTextContent('TAG_40002');
    });
    await waitFor(() => {
      expect(screen.getByTestId('trace-tag-key')).toHaveTextContent('TAG_40002');
    });
  });

  // ---------------------------------------------------------------------------
  // Modbus studio deepening (redesign-modbus-studio)
  // ---------------------------------------------------------------------------

  describe('modbus studio register map canvas', () => {
    it('renders a register map canvas showing allocated register ranges', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
        { tag_id: 'tag-2', register: 2, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(
        await screen.findByTestId('register-map-canvas'),
      ).toBeInTheDocument();
      expect(screen.getByTestId('register-slot-0')).toHaveTextContent('TAG_40001');
      expect(screen.getByTestId('register-slot-2')).toHaveTextContent('TAG_40002');
    });

    it('highlights conflicting register slots', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 10, data_type: 'int16', updated_at: '' },
        { tag_id: 'tag-2', register: 10, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const conflictSlot = await screen.findByTestId('register-slot-10');
      expect(conflictSlot).toHaveAttribute('data-conflict', 'true');
    });

    it('treats overlapping multi-word ranges as conflicting slots', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 10, data_type: 'int32', updated_at: '' },
        { tag_id: 'tag-2', register: 11, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const overlapSlot = await screen.findByTestId('register-slot-11');
      expect(overlapSlot).toHaveAttribute('data-conflict', 'true');
    });

    it('renders auto-map strategy selector with three strategies', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('register-map-canvas');

      expect(
        screen.getByRole('button', { name: 'workbench.output.modbusStudio.autoMap.sequential' }),
      ).toBeInTheDocument();
      expect(
        screen.getByRole('button', { name: 'workbench.output.modbusStudio.autoMap.gapAware' }),
      ).toBeInTheDocument();
      expect(
        screen.getByRole('button', { name: 'workbench.output.modbusStudio.autoMap.aligned' }),
      ).toBeInTheDocument();
    });

    it('sequential auto-map appends after the current occupied range', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('register-map-canvas');

      fireEvent.click(
        screen.getByRole('button', {
          name: 'workbench.output.modbusStudio.autoMap.sequential',
        }),
      );

      await waitFor(() => {
        expect(mockModbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-2', 1);
      });
    });

    it('renders a dry-run validation button and displays results', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('register-map-canvas');

      const dryRunButton = screen.getByRole('button', {
        name: 'workbench.output.modbusStudio.actions.dryRun',
      });
      expect(dryRunButton).toBeInTheDocument();

      fireEvent.click(dryRunButton);

      await waitFor(() => {
        expect(screen.getByTestId('dry-run-results')).toBeInTheDocument();
      });
    });
  });

  // ---------------------------------------------------------------------------
  // Database studio deepening (redesign-database-studio)
  // ---------------------------------------------------------------------------

  describe('database studio schema snapshot', () => {
    it('collapses connector fields until the operator expands connector setup', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      expect(
        screen.queryByLabelText('workbench.output.database.connector.name'),
      ).not.toBeInTheDocument();

      fireEvent.click(
        screen.getByRole('button', {
          name: 'workbench.output.database.actions.configureConnector',
        }),
      );

      expect(
        await screen.findByLabelText('workbench.output.database.connector.name'),
      ).toBeInTheDocument();
    });

    it('groups schema snapshot and write preview into supporting secondary panels', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const secondaryPanels = await screen.findByTestId('database-secondary-panels');
      expect(secondaryPanels).toHaveAttribute('data-emphasis', 'supporting');
      expect(await screen.findByTestId('schema-snapshot')).toBeInTheDocument();
      expect(within(secondaryPanels).getByTestId('write-row-preview')).toBeInTheDocument();
    });

    it('renders a schema snapshot with column type badges', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      expect(
        await screen.findByTestId('schema-snapshot'),
      ).toBeInTheDocument();
      expect(screen.getByTestId('schema-column-ts')).toHaveTextContent('datetime');
      expect(screen.getByTestId('schema-column-value')).toHaveTextContent('real');
    });

    it('marks primary key columns with a PK badge', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const tsColumn = await screen.findByTestId('schema-column-ts');
      expect(tsColumn).toHaveTextContent('PK');
    });

    it('renders a write-row preview panel', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      expect(
        await screen.findByTestId('write-row-preview'),
      ).toBeInTheDocument();
    });

    it('highlights unmapped required columns in the schema snapshot', async () => {
      mockDBTargetAPI.listTables.mockResolvedValue([
        {
          schema: 'main',
          name: 'sensor_values',
          columns: [
            { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true, unique: true },
            { name: 'value', data_type: 'real', nullable: false, primary_key: false, unique: false },
            { name: 'source_id', data_type: 'text', nullable: false, primary_key: false, unique: false },
          ],
        },
      ]);
      mockDBTargetAPI.listMappings.mockResolvedValue([]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const requiredColumn = await screen.findByTestId('schema-column-source_id');
      expect(requiredColumn).toHaveAttribute('data-required', 'true');
    });
  });

  // ---------------------------------------------------------------------------
  // Output inspector traceability (redesign-output-inspector)
  // ---------------------------------------------------------------------------

  describe('output inspector traceability', () => {
    it('shows source→tag→output trace when an output candidate is selected', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const candidate = await screen.findByTestId('output-candidate-tag-1');
      fireEvent.click(candidate);

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-source-address')).toHaveTextContent('40001');
      expect(screen.getByTestId('trace-tag-key')).toHaveTextContent('TAG_40001');
      expect(screen.getByTestId('trace-output-modbus')).toHaveTextContent('HR0');
    });

    it('shows readiness reasons for partial output candidate', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);
      mockDBTargetAPI.listMappings.mockResolvedValue([]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const candidate = await screen.findByTestId('output-candidate-tag-1');
      fireEvent.click(candidate);

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-readiness')).toHaveTextContent(
        'workbench.output.inspector.readiness.partial',
      );
    });

    it('shows ready readiness when both modbus and database are mapped', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);
      mockDBTargetAPI.listMappings.mockResolvedValue([
        {
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
        },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const candidate = await screen.findByTestId('output-candidate-tag-1');
      fireEvent.click(candidate);

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-readiness')).toHaveTextContent(
        'workbench.output.inspector.readiness.ready',
      );
    });

    it('shows database path in trace when database mapping exists', async () => {
      mockDBTargetAPI.listMappings.mockResolvedValue([
        {
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
        },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const candidate = await screen.findByTestId('output-candidate-tag-1');
      fireEvent.click(candidate);

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-output-database')).toHaveTextContent(
        'main.sensor_values.value',
      );
    });

    it('announces inspector mapping load failures through a polite live region', async () => {
      mockModbusShareAPI.listMappings.mockRejectedValue(new Error('boom'));
      mockDBTargetAPI.listMappings.mockResolvedValue([]);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.output' }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const candidate = await screen.findByTestId('output-candidate-tag-1');
      fireEvent.click(candidate);

      const tracePanel = await screen.findByTestId('inspector-trace-panel');
      const loadError = await within(tracePanel).findByText('boom');

      expect(loadError).toHaveAttribute('role', 'status');
      expect(loadError).toHaveAttribute('aria-live', 'polite');
    });
  });
});
