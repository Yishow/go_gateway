import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type { Device, Mapping, Point, Tag } from '@/types/datalink';

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
    generateSchema: vi.fn(),
    dryRunMappings: vi.fn(),
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

vi.mock('@/services/datalink', () => ({
  modbusShareAPI: mockModbusShareAPI,
  dbTargetAPI: mockDBTargetAPI,
}));

vi.mock('@/hooks/datalink/useDevices', () => ({
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
  useTestDraftConnectionMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useToggleDeviceStatusMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('@/hooks/datalink/usePoints', () => ({
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
  useDeletePointMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
  }),
  useCreateTagMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDeleteTagMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({
    data: mockMappings,
    isLoading: false,
  }),
  useCreateMappingMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDeleteMappingMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
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
    mockDBTargetAPI.generateSchema.mockResolvedValue({
      connector_id: 'connector-1',
      dry_run: false,
      statements: ['CREATE TABLE sensor_values (...)'],
      executed: 1,
    });
    mockDBTargetAPI.dryRunMappings.mockResolvedValue({
      connector_id: 'connector-1',
      results: [
        {
          candidate_id: 'db-mapping-1',
          status: 'ready',
          mapping_id: 'db-mapping-1',
          tag_id: 'tag-1',
        },
      ],
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

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByText('workbench.output.empty.title')).toBeInTheDocument();
  });

  it('starts the local modbus server with the specified port', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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

  it('renders compact tag chips with status suffixes for the active target', async () => {
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

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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

    const tagChips = screen.getByTestId('output-tag-chips');
    expect(tagChips).toBeInTheDocument();
    expect(screen.getByTestId('output-candidate-tag-1')).toBeInTheDocument();

    // Click tag-1 chip and verify active badge appears
    fireEvent.click(screen.getByTestId('output-candidate-tag-1'));
    expect(screen.getByTestId('active-tag-badge')).toBeInTheDocument();

    // Switch to database target
    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Select tag-2 via chip click
    fireEvent.click(screen.getByTestId('output-candidate-tag-2'));

    await waitFor(() => {
      expect(screen.getByTestId('active-tag-badge')).toHaveTextContent('TAG_40002');
    });

    expect(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    ).toHaveAttribute('aria-pressed', 'true');
  });

  it('marks modbus operational panels as supporting sections', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByTestId('modbus-secondary-panels')).toHaveAttribute(
      'data-emphasis',
      'supporting',
    );
  });

  it('uses the tag chips as the unified selection surface for both targets', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByTestId('output-tag-chips');

    const tagChips = screen.getByTestId('output-tag-chips');
    expect(tagChips).toBeInTheDocument();
    expect(screen.getByTestId('output-candidate-tag-1')).toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Tag chips remain visible after switching target
    expect(tagChips).toBeInTheDocument();
    // Database panel shows selected tag context
    await waitFor(() => {
      expect(screen.getByTestId('database-selected-tag')).toBeInTheDocument();
    });
  });

  it('binds a linked tag to a local modbus register', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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
      expect(mockModbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-1', 11);
    });
  });

  it('keeps tag chip selection and register input in sync', async () => {
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

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    await screen.findByLabelText('workbench.output.mapping.register');
    await waitFor(() => {
      expect(screen.getByDisplayValue('13')).toBeInTheDocument();
    });

    const tag1Chip = screen.getByTestId('output-candidate-tag-1');
    expect(tag1Chip).toHaveAttribute('aria-pressed', 'true');

    fireEvent.click(screen.getByTestId('output-candidate-tag-2'));

    await waitFor(() => {
      expect(screen.getByDisplayValue('25')).toBeInTheDocument();
    });
    expect(screen.getByTestId('output-candidate-tag-2')).toHaveAttribute('aria-pressed', 'true');
    expect(tag1Chip).toHaveAttribute('aria-pressed', 'false');
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

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(await screen.findByRole('button', { name: 'workbench.output.actions.sync' })).toBeDisabled();
    expect(screen.getByRole('status')).toHaveTextContent('workbench.output.conflicts.summary');
  });

  it('binds a linked tag to a database column via direct surface click', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Wait for schema columns to load
    await screen.findByTestId('schema-column-surface');

    // Select a table first
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.output.database.mapping.table')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText('workbench.output.database.mapping.table'), {
      target: { value: 'main.sensor_values' },
    });

    // Click on the 'value' column to bind tag-1
    await waitFor(() => {
      expect(screen.getByTestId('schema-column-value')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('schema-column-value'));

    await waitFor(() => {
      expect(mockDBTargetAPI.createMapping).toHaveBeenCalledWith(
        expect.objectContaining({
          tag_id: 'tag-1',
          connector_id: 'connector-1',
          table_name: 'sensor_values',
          column_name: 'value',
        }),
      );
    });
    expect(screen.getByRole('status')).toHaveTextContent(
      'workbench.output.database.results.mappingSaved',
    );
  });

  it('removes a database mapping with inline feedback when clicking a bound column', async () => {
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

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    await screen.findByTestId('schema-column-value');
    fireEvent.click(screen.getByTestId('schema-column-value'));

    await waitFor(() => {
      expect(mockDBTargetAPI.deleteMapping).toHaveBeenCalledWith('db-mapping-1');
    });
    expect(screen.getByRole('status')).toHaveTextContent(
      'workbench.output.database.results.mappingDeleted',
    );
  });

  it('keeps the database selected-tag display aligned with the tag chips', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      await screen.findByRole('button', {
        name: 'workbench.output.targetSwitcher.database',
      }),
    );

    // Select tag-2 via chip
    fireEvent.click(screen.getByTestId('output-candidate-tag-2'));

    await waitFor(() => {
      expect(screen.getByTestId('database-selected-tag')).toHaveTextContent('TAG_40002');
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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(
        await screen.findByTestId('register-map-canvas'),
      ).toBeInTheDocument();
      expect(screen.getByTestId('register-slot-1')).toHaveTextContent('TAG_40001');
      expect(screen.getByTestId('register-slot-3')).toHaveTextContent('TAG_40002');
    });

    it('highlights conflicting register slots', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 10, data_type: 'int16', updated_at: '' },
        { tag_id: 'tag-2', register: 10, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const conflictSlot = await screen.findByTestId('register-slot-11');
      expect(conflictSlot).toHaveAttribute('data-conflict', 'true');
    });

    it('treats overlapping multi-word ranges as conflicting slots', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 10, data_type: 'int32', updated_at: '' },
        { tag_id: 'tag-2', register: 11, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const overlapSlot = await screen.findByTestId('register-slot-12');
      expect(overlapSlot).toHaveAttribute('data-conflict', 'true');
    });

    it('keeps high mapped registers visible in the canvas', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-2', register: 199, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(await screen.findByTestId('register-slot-200')).toHaveTextContent('TAG_40002');
    });

    it('lets the operator unbind HR1 directly from the canvas', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      fireEvent.click(await screen.findByTestId('register-slot-1'));

      await waitFor(() => {
        expect(mockModbusShareAPI.deleteMapping).toHaveBeenCalledWith('tag-1');
      });
      expect(screen.getByRole('status')).toHaveTextContent('workbench.output.results.slotUnbound');
    });

    it('binds the selected tag by clicking an empty register slot', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const targetSlot = await screen.findByTestId('register-slot-5');
      fireEvent.click(targetSlot);

      await waitFor(() => {
        expect(mockModbusShareAPI.upsertMapping).toHaveBeenCalledWith('tag-1', 4);
      });
      expect(screen.getByRole('status')).toHaveTextContent('workbench.output.results.slotBound');
    });

    it('renders auto-map strategy selector with three strategies', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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
    it('clears stale table scope when the operator switches connectors', async () => {
      mockDBTargetAPI.listConnectors.mockResolvedValue([
        {
          id: 'connector-1',
          name: 'Main SQLite',
          kind: 'sqlite',
          connection_config: { dsn: '/tmp/target.db' },
          status: 'ready',
          last_check_at: '',
          last_check_error: '',
          enabled: true,
          created_at: '',
          updated_at: '',
        },
        {
          id: 'connector-2',
          name: 'Warehouse PostgreSQL',
          kind: 'postgres',
          connection_config: {
            host: '127.0.0.1',
            port: '5432',
            user: 'postgres',
            database: 'warehouse',
            sslmode: 'disable',
          },
          status: 'ready',
          last_check_at: '',
          last_check_error: '',
          enabled: true,
          created_at: '',
          updated_at: '',
        },
      ]);
      mockDBTargetAPI.listTables.mockImplementation(async (connectorId: string) => {
        if (connectorId === 'connector-2') {
          return [
            {
              schema: 'analytics',
              name: 'metrics',
              columns: [
                { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true, unique: true },
                { name: 'reading', data_type: 'real', nullable: false, primary_key: false, unique: false },
              ],
            },
          ];
        }

        return [
          {
            schema: 'main',
            name: 'sensor_values',
            columns: [
              { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true, unique: true },
              { name: 'value', data_type: 'real', nullable: false, primary_key: false, unique: false },
            ],
          },
        ];
      });

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      await waitFor(() => {
        expect(screen.getByLabelText('workbench.output.database.mapping.table')).toHaveValue(
          'main.sensor_values',
        );
      });
      expect(screen.getByTestId('schema-column-value')).toBeInTheDocument();

      fireEvent.click(screen.getByRole('button', { name: /Warehouse PostgreSQL/i }));

      await waitFor(() => {
        expect(screen.getByLabelText('workbench.output.database.mapping.table')).toHaveValue(
          'analytics.metrics',
        );
      });
      expect(screen.getByTestId('schema-column-reading')).toBeInTheDocument();
      expect(screen.queryByTestId('schema-column-value')).not.toBeInTheDocument();
    });

    it('keeps write-mode and timestamp scope synchronized', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const writeModeSelect = await screen.findByLabelText(
        'workbench.output.database.mapping.writeMode',
      );
      fireEvent.change(writeModeSelect, {
        target: { value: 'upsert' },
      });

      await waitFor(() => {
        expect(
          screen.getByLabelText('workbench.output.database.mapping.timestampColumn'),
        ).toHaveValue('ts');
      });

      fireEvent.change(writeModeSelect, {
        target: { value: 'insert' },
      });

      await waitFor(() => {
        expect(
          screen.queryByLabelText('workbench.output.database.mapping.timestampColumn'),
        ).not.toBeInTheDocument();
      });
    });

    it('does not overwrite manual write-mode edits when a mapping already exists', async () => {
      mockDBTargetAPI.listMappings.mockResolvedValue([
        {
          id: 'db-mapping-1',
          tag_id: 'tag-1',
          connector_id: 'connector-1',
          table_schema: 'main',
          table_name: 'sensor_values',
          column_name: 'value',
          write_mode: 'upsert',
          timestamp_column: 'ts',
          enabled: true,
          created_at: '',
          updated_at: '',
        },
      ]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const writeModeSelect = await screen.findByLabelText(
        'workbench.output.database.mapping.writeMode',
      );
      await waitFor(() => {
        expect(
          screen.getByLabelText('workbench.output.database.mapping.writeMode'),
        ).toHaveValue('upsert');
      });
      expect(
        screen.getByLabelText('workbench.output.database.mapping.timestampColumn'),
      ).toHaveValue('ts');

      fireEvent.change(writeModeSelect, {
        target: { value: 'insert' },
      });

      await waitFor(() => {
        expect(
          screen.getByLabelText('workbench.output.database.mapping.writeMode'),
        ).toHaveValue('insert');
      });
      expect(
        screen.queryByLabelText('workbench.output.database.mapping.timestampColumn'),
      ).not.toBeInTheDocument();
    });

    it('collapses connector fields until the operator expands connector setup', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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

    it('keeps schema snapshot in supporting panels while the row planner stays on the main surface', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const secondaryPanels = await screen.findByTestId('database-secondary-panels');
      expect(secondaryPanels).toHaveAttribute('data-emphasis', 'supporting');
      expect(await within(secondaryPanels).findByTestId('schema-snapshot')).toBeInTheDocument();
      expect(screen.getByTestId('database-grouped-row-planner')).toBeInTheDocument();
      expect(within(secondaryPanels).queryByTestId('database-grouped-row-planner')).not.toBeInTheDocument();
    });

    it('renders a schema snapshot with column type badges', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      const tsColumn = await screen.findByTestId('schema-column-ts');
      expect(tsColumn).toHaveTextContent('PK');
    });

    it('renders the grouped row planner surface', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(
        await screen.findByRole('button', {
          name: 'workbench.output.targetSwitcher.database',
        }),
      );

      expect(
        await screen.findByTestId('database-grouped-row-planner'),
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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

      await waitFor(() => {
        expect(screen.getByTestId('inspector-trace-panel')).toBeInTheDocument();
      });

      expect(screen.getByTestId('trace-source-address')).toHaveTextContent('40001');
      expect(screen.getByTestId('trace-tag-key')).toHaveTextContent('TAG_40001');
      expect(screen.getByTestId('trace-output-modbus')).toHaveTextContent('HR1');
    });

    it('shows readiness reasons for partial output candidate', async () => {
      mockModbusShareAPI.listMappings.mockResolvedValue([
        { tag_id: 'tag-1', register: 0, data_type: 'int16', updated_at: '' },
      ]);
      mockDBTargetAPI.listMappings.mockResolvedValue([]);

      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

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

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');
      fireEvent.click(screen.getByTestId('output-candidate-tag-1'));

      const tracePanel = await screen.findByTestId('inspector-trace-panel');
      const loadError = await within(tracePanel).findByText('boom');

      expect(loadError).toHaveAttribute('role', 'status');
      expect(loadError).toHaveAttribute('aria-live', 'polite');
    });
  });

  describe('cross-step tag handoff (Step 3 → Step 4)', () => {
    it('pre-selects the focused tag from Step 3 when entering Step 4', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      await screen.findByTestId('output-tag-chips');

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.tag/ }));
      const boundRow = await screen.findByTestId('tag-candidate-point-2');
      fireEvent.click(boundRow);

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      await screen.findByTestId('output-tag-chips');

      await waitFor(() => {
        expect(screen.getByTestId('output-candidate-tag-2')).toHaveAttribute('aria-pressed', 'true');
      });
    });

    it('falls back to the first candidate when focusedTagIds do not match', async () => {
      renderPage();

      fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      await screen.findByTestId('output-tag-chips');

      await waitFor(() => {
        expect(screen.getByTestId('output-candidate-tag-1')).toHaveAttribute('aria-pressed', 'true');
      });
    });
  });
});
