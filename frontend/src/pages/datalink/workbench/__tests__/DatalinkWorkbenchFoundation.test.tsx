import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import App from '../../../../App';

const {
  mockDevices,
  mockCreateDeviceMutation,
  mockTestDraftConnectionMutation,
  mockUpdateDeviceMutation,
  mockTestConnectionMutation,
} = vi.hoisted(() => ({
  mockDevices: [] as Array<{
    id: string;
    name: string;
    description: string;
    protocol: string;
    status: 'draft' | 'active' | 'disabled';
    connection_config: string;
    last_test_at: string | null;
    last_test_success: boolean | null;
    last_test_error: string;
    created_at: string;
    updated_at: string;
  }>,
  mockCreateDeviceMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockTestDraftConnectionMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockUpdateDeviceMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockTestConnectionMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('@/contexts/ThemeContext', () => ({
  ThemeProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

vi.mock('@/contexts/ToastContext', () => ({
  ToastProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

vi.mock('@/components/CardMinimizeProvider', () => ({
  CardMinimizeProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

vi.mock('@/components/Layout', () => ({
  default: ({ children }: { children: ReactNode }) => <div data-testid="layout-mock">{children}</div>,
}));

vi.mock('@/pages/TestPage', () => ({
  default: () => <div data-testid="test-page-mock">test-page</div>,
}));

vi.mock('@/pages/TemplatesPage', () => ({
  default: () => <div data-testid="templates-page-mock">templates-page</div>,
}));

vi.mock('@/pages/HistoryPage', () => ({
  default: () => <div data-testid="history-page-mock">history-page</div>,
}));

vi.mock('@/pages/ComparePage', () => ({
  default: () => <div data-testid="compare-page-mock">compare-page</div>,
}));

vi.mock('@/pages/AnalyzerPage', () => ({
  default: () => <div data-testid="analyzer-page-mock">analyzer-page</div>,
}));

vi.mock('@/pages/datalink/SmartDashboard', () => ({
  default: () => <div data-testid="smart-dashboard-mock">smart-dashboard</div>,
}));

vi.mock('@/pages/datalink/LocalModbusWorkbenchPage', () => ({
  default: () => <div data-testid="local-modbus-workbench-mock">local-modbus-workbench</div>,
}));

vi.mock('@/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockDevices,
    isLoading: false,
  }),
  useCreateDeviceMutation: () => mockCreateDeviceMutation,
  useTestDraftConnectionMutation: () => mockTestDraftConnectionMutation,
  useUpdateDeviceMutation: () => mockUpdateDeviceMutation,
  useTestConnectionMutation: () => mockTestConnectionMutation,
}));

vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: () => ({
    data: [],
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

vi.mock('@/router/gateway', () => ({
  GatewayCreateEntryRedirect: () => <div data-testid="gateway-create-entry-redirect-mock" />,
  GatewayEntryRoute: () => <div data-testid="gateway-entry-route-mock" />,
  GatewayExpertWorkbenchRoute: () => <div data-testid="gateway-expert-workbench-route-mock" />,
  GatewayQuickSetupRoute: () => <div data-testid="gateway-quick-setup-route-mock" />,
}));

vi.mock('@/features/datalink/legacyRoutes', () => ({
  buildDashboardModalRedirect: (intent: string) => `/mock/dashboard/${intent}`,
  buildLegacyMigrationRedirect: (intent: string) => `/mock/legacy/${intent}`,
  buildLocalModbusCompatRedirect: (section?: string | null) =>
    `/datalink/workbench?step=output&target=modbus${section ? `&section=${section}` : ''}`,
}));

vi.mock('@/services/datalink', async () => {
  const actual = await vi.importActual<typeof import('@/services/datalink')>(
    '@/services/datalink',
  );
  return {
    ...actual,
    protocolAPI: {
      list: vi.fn().mockResolvedValue([]),
    },
    settingsAPI: {
      get: vi.fn().mockResolvedValue({}),
    },
  };
});

describe('DatalinkWorkbench foundation route', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/datalink/workbench');
    mockDevices.splice(
      0,
      mockDevices.length,
      {
        id: 'device-1',
        name: 'Mixer PLC',
        description: 'Main line',
        protocol: 'modbus_tcp',
        status: 'active',
        connection_config: '{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: true,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'device-2',
        name: 'Backup PLC',
        description: 'Fallback',
        protocol: 'modbus_rtu',
        status: 'draft',
        connection_config: '{"serial_port":"COM3","baud_rate":9600,"data_bits":8,"stop_bits":1,"parity":"none","slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
    );
    mockCreateDeviceMutation.mutateAsync.mockReset();
    mockCreateDeviceMutation.mutateAsync.mockResolvedValue({
      id: 'device-new',
      name: 'Browser Smoke PLC',
      description: 'Smoke flow',
      protocol: 'modbus_tcp',
      status: 'draft',
      connection_config: '{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    });
    mockTestDraftConnectionMutation.mutateAsync.mockReset();
    mockTestDraftConnectionMutation.isPending = false;
    mockUpdateDeviceMutation.mutateAsync.mockReset();
    mockUpdateDeviceMutation.mutateAsync.mockResolvedValue(undefined);
    mockTestConnectionMutation.mutateAsync.mockReset();
    mockTestConnectionMutation.isPending = false;
  });

  function renderApp() {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });

    return render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>,
    );
  }

  it('renders a real device setup step through /datalink/workbench', () => {
    renderApp();

    expect(
      screen.getByRole('navigation', { name: 'workbench.stepRail.ariaLabel' }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: /workbench\.steps\.device/ }),
    ).toHaveAttribute('aria-current', 'step');
    expect(
      screen.getByRole('textbox', { name: 'workbench.device.search.label' }),
    ).toBeInTheDocument();
    expect(screen.getByText('Mixer PLC')).toBeInTheDocument();
    expect(screen.getByText('Backup PLC')).toBeInTheDocument();
    expect(screen.queryByText('workbench.placeholders.device')).not.toBeInTheDocument();
    expect(screen.getByTestId('workbench-inspector-panel')).toBeInTheDocument();
  });

  it('keeps Step 1 search, protocol filter, and create action inside one primary toolbar', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const toolbar = screen.getByTestId('device-primary-toolbar');

    expect(
      within(toolbar).getByRole('textbox', { name: 'workbench.device.search.label' }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('combobox', { name: 'workbench.device.filters.protocol' }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', { name: 'workbench.device.actions.create' }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).queryByRole('button', { name: 'workbench.device.actions.clone' }),
    ).not.toBeInTheDocument();
    expect(
      within(toolbar).queryByRole('button', { name: 'workbench.device.actions.continue' }),
    ).not.toBeInTheDocument();
  });

  it('drops the decorative Step 1 hero block once the compact toolbar is available', () => {
    renderApp();

    expect(screen.queryByText('workbench.device.title')).not.toBeInTheDocument();
    expect(screen.queryByText('workbench.device.description')).not.toBeInTheDocument();
  });

  it('renders compact device rows with endpoint, health, and two primary capability hints', () => {
    renderApp();

    const row = screen.getByTestId('device-row-device-1');

    expect(within(row).getByText('192.168.1.10:502')).toBeInTheDocument();
    expect(within(row).getByTestId('device-health-device-1')).toHaveTextContent(
      'workbench.device.card.testPassed',
    );
    expect(
      within(row).getByText('workbench.device.capability.labels.unitId'),
    ).toBeInTheDocument();
    expect(
      within(row).getByText('workbench.device.capability.labels.addressBase'),
    ).toBeInTheDocument();
    expect(
      within(row).queryByText('workbench.device.capability.labels.wordOrder'),
    ).not.toBeInTheDocument();
    expect(
      within(row).queryByText('workbench.device.capability.labels.protocolTraits'),
    ).not.toBeInTheDocument();
    expect(within(row).queryByText('Main line')).not.toBeInTheDocument();
  });

  it('uses protocol-specific connection targets for compact row endpoints', () => {
    mockDevices.push(
      {
        id: 'device-3',
        name: 'Telemetry Broker',
        description: 'MQTT edge',
        protocol: 'mqtt',
        status: 'active',
        connection_config:
          '{"broker_url":"mqtt://broker.internal:1883","client_id":"edge-gateway","topics":["plant/telemetry"],"use_tls":true,"qos":1}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'device-4',
        name: 'Packaging PLC',
        description: 'MC 3E line',
        protocol: 'mc_3e',
        status: 'active',
        connection_config:
          '{"host":"10.0.0.20","port":5000,"network_no":1,"station_no":2,"data_format":"binary"}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
    );

    renderApp();

    const mqttRow = screen.getByTestId('device-row-device-3');
    expect(within(mqttRow).getByText('mqtt://broker.internal:1883')).toBeInTheDocument();
    expect(within(mqttRow).queryByText(/QoS 1/)).not.toBeInTheDocument();

    const mcRow = screen.getByTestId('device-row-device-4');
    expect(within(mcRow).getByText('10.0.0.20:5000')).toBeInTheDocument();
    expect(within(mcRow).queryByText('1/2 · binary')).not.toBeInTheDocument();
  });

  it('supports workbench deep links for step and output target', async () => {
    window.history.pushState({}, '', '/datalink/workbench?step=output&target=database');

    renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-current', 'step');
    });
    expect(screen.getByTestId('active-output-target')).toHaveTextContent(
      'workbench.bottomSummary.targets.database',
    );
  });

  it('does not lock step navigation after a deep-link is applied', async () => {
    window.history.pushState({}, '', '/datalink/workbench?step=output&target=database');

    renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-current', 'step');
    });

    // User clicks the device step in the step rail — must NOT be locked back to output.
    fireEvent.click(
      screen.getByRole('button', { name: /workbench\.steps\.device/ }),
    );

    expect(
      screen.getByRole('button', { name: /workbench\.steps\.device/ }),
    ).toHaveAttribute('aria-current', 'step');
    expect(
      screen.getByRole('button', { name: /workbench\.steps\.output/ }),
    ).not.toHaveAttribute('aria-current', 'step');
  });

  it('redirects legacy local modbus entry into the new workbench output step', async () => {
    window.history.pushState({}, '', '/datalink/local-modbus?section=settings');

    renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-current', 'step');
    });
    expect(screen.queryByTestId('local-modbus-workbench-mock')).not.toBeInTheDocument();
    expect(screen.getByTestId('active-output-target')).toHaveTextContent(
      'workbench.bottomSummary.targets.modbus',
    );
  });

  it('selects a device from the device step and advances to source planning', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

    expect(
      screen.getByRole('button', { name: /workbench\.steps\.source/ }),
    ).toHaveAttribute('aria-current', 'step');
    expect(
      screen.getByLabelText('workbench.source.planner.startAddress'),
    ).toBeInTheDocument();
    expect(
      screen.getByLabelText('workbench.source.planner.count'),
    ).toBeInTheDocument();
  });

  it('shows a master-detail device panel after a device is selected', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const detailPanel = screen.getByTestId('device-detail-panel');
    expect(detailPanel).toBeInTheDocument();
    expect(within(detailPanel).getByText('Mixer PLC')).toBeInTheDocument();
    expect(within(detailPanel).getByTestId('device-detail-endpoint')).toHaveTextContent(
      '192.168.1.10:502',
    );
    expect(
      within(detailPanel).getByRole('button', {
        name: 'workbench.device.actions.continue',
      }),
    ).toBeInTheDocument();
  });

  it('context bar keeps a single primary action after selecting a device on step 1', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const contextBar = screen.getByTestId('workbench-context-bar');
    const actions = within(contextBar).getAllByRole('button');

    expect(actions).toHaveLength(1);
    expect(actions[0]).toHaveTextContent('workbench.contextBar.actions.gotoSource');
  });

  it('opens the create device form inline when there are no devices yet', () => {
    mockDevices.splice(0, mockDevices.length);

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));

    expect(screen.getByTestId('device-inline-editor')).toBeInTheDocument();
    expect(screen.queryByTestId('device-panel-overlay')).not.toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.createTitle' }),
    ).toBeInTheDocument();
  });

  it('localizes protocol-specific connection option labels in the create drawer', () => {
    mockDevices.splice(0, mockDevices.length);

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));

    fireEvent.change(screen.getByLabelText('workbench.device.fields.protocol'), {
      target: { value: 'modbus_rtu' },
    });

    const paritySelect = screen.getByLabelText('workbench.device.connection.parity');
    expect(within(paritySelect).getByRole('option', { name: 'device.parityNone' })).toBeInTheDocument();
    expect(within(paritySelect).getByRole('option', { name: 'device.parityEven' })).toBeInTheDocument();
    expect(within(paritySelect).getByRole('option', { name: 'device.parityOdd' })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('workbench.device.fields.protocol'), {
      target: { value: 'mc_3e' },
    });

    const dataFormatSelect = screen.getByLabelText('workbench.device.connection.dataFormat');
    expect(
      within(dataFormatSelect).getByRole('option', {
        name: 'workbench.device.connection.dataFormats.binary',
      }),
    ).toBeInTheDocument();
    expect(
      within(dataFormatSelect).getByRole('option', {
        name: 'workbench.device.connection.dataFormats.ascii',
      }),
    ).toBeInTheDocument();
  });

  it('keeps the created device selected while the list refreshes', async () => {
    mockDevices.splice(0, mockDevices.length);

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));
    fireEvent.change(screen.getByLabelText('workbench.device.fields.name'), {
      target: { value: 'Browser Smoke PLC' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.host'), {
      target: { value: '127.0.0.1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.port'), {
      target: { value: '502' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.slaveId'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.device.connection.timeout'), {
      target: { value: '5' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(
        screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }),
      ).toBeEnabled();
    });
  });

  it('allows clearing device description when editing', async () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.device.actions.edit' }),
    );
    fireEvent.change(
      screen.getByLabelText('workbench.device.fields.description'),
      { target: { value: '' } },
    );
    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(mockUpdateDeviceMutation.mutateAsync).toHaveBeenCalledWith(
        expect.objectContaining({
          id: 'device-1',
          data: expect.objectContaining({
            name: 'Mixer PLC',
            description: '',
          }),
        }),
      );
    });
  });

  it('shows the inline editor in the detail column while keeping the device list visible', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.device.actions.edit' }),
    );

    expect(screen.getByTestId('device-inline-editor')).toBeInTheDocument();
    expect(screen.queryByTestId('device-panel-overlay')).not.toBeInTheDocument();
    expect(screen.queryByTestId('device-detail-panel')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Mixer PLC' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Backup PLC' })).toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.editTitle' }),
    ).toBeInTheDocument();
  });

  it('tests the current draft config from the inline editor and explains that the backend host runs the dial', async () => {
    mockTestDraftConnectionMutation.mutateAsync.mockResolvedValueOnce({
      success: true,
      error: '',
      latency_ms: 18,
      can_activate: true,
      can_collect: true,
      connect: {
        status: 'success',
        message: 'connect ok',
        latency_ms: 7,
      },
      probe: {
        status: 'success',
        message: 'probe ok',
        latency_ms: 11,
      },
    });

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.device.actions.edit' }),
    );

    expect(
      screen.getByText('workbench.device.connection.backendHostHint'),
    ).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('workbench.device.connection.host'), {
      target: { value: '10.0.0.77' },
    });
    fireEvent.click(
      screen.getByRole('button', {
        name: 'workbench.device.actions.testDraftConnection',
      }),
    );

    await waitFor(() => {
      expect(mockTestDraftConnectionMutation.mutateAsync).toHaveBeenCalledWith({
        protocol: 'modbus_tcp',
        connection_config: {
          host: '10.0.0.77',
          port: 502,
          slave_id: 1,
          timeout: 5,
        },
      });
    });
  });

  it('hides zero-value test timestamps in the inspector', () => {
    mockDevices[0].last_test_at = '0001-01-01T00:00:00Z';

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(
      screen.getByText('workbench.device.inspector.unknownTestTime'),
    ).toBeInTheDocument();
    expect(screen.queryByText('0001-01-01T00:00:00Z')).not.toBeInTheDocument();
  });

  it('renders selected device details in the right-side inspector', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const inspector = screen.getByTestId('workbench-inspector-panel');
    expect(within(inspector).getByText('Mixer PLC')).toBeInTheDocument();
    expect(within(inspector).getByText('workbench.device.inspector.connectionSummary')).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.edit' })).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.testConnection' })).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.clone' })).toBeInTheDocument();
  });

  it('opens a clone drawer with connection defaults but requires a new device name', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-inspector-panel')).getByRole('button', {
        name: 'workbench.device.actions.clone',
      }),
    );

    expect(screen.getByTestId('device-inline-editor')).toBeInTheDocument();
    expect(screen.queryByTestId('device-panel-overlay')).not.toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.cloneTitle' }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText('workbench.device.fields.name')).toHaveValue('');
    expect(screen.getByLabelText('workbench.device.connection.host')).toHaveValue('192.168.1.10');
    expect(screen.getByLabelText('workbench.device.connection.port')).toHaveValue('502');
    expect(screen.getByLabelText('workbench.device.connection.slaveId')).toHaveValue('1');
  });

  it('prevents saving a clone with the same source device name', async () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-inspector-panel')).getByRole('button', {
        name: 'workbench.device.actions.clone',
      }),
    );

    fireEvent.change(screen.getByLabelText('workbench.device.fields.name'), {
      target: { value: 'Mixer PLC' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.save' }));

    await waitFor(() => {
      expect(
        screen.getByText('workbench.device.validation.cloneNameDistinct'),
      ).toBeInTheDocument();
    });
    expect(mockCreateDeviceMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('disables inspector test action while a connection test is already pending', () => {
    mockTestConnectionMutation.isPending = true;

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    const inspector = screen.getByTestId('workbench-inspector-panel');
    expect(
      within(inspector).getByRole('button', {
        name: 'workbench.device.actions.testing',
      }),
    ).toBeDisabled();
  });

  it('keeps the latest session test result in the inspector while the context bar stays compact', async () => {
    mockTestConnectionMutation.mutateAsync
      .mockReset()
      .mockResolvedValueOnce({ success: false, error: 'timeout-latest', latency_ms: 0 });

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-inspector-panel')).getByRole('button', {
        name: 'workbench.device.actions.testConnection',
      }),
    );

    await waitFor(() => {
      const inspectorEntry = within(screen.getByTestId('workbench-inspector-panel')).getByText(
        'timeout-latest',
      );
      expect(inspectorEntry).toBeInTheDocument();
      expect(inspectorEntry).toHaveClass('text-rose-300');
      expect(screen.queryByTestId('context-bar-test-status')).not.toBeInTheDocument();
    });
  });

  it('keeps only the three most recent connection tests in the inspector timeline', async () => {
    mockTestConnectionMutation.mutateAsync
      .mockReset()
      .mockResolvedValueOnce({ success: false, error: 'timeout-1', latency_ms: 0 })
      .mockResolvedValueOnce({ success: true, error: '', latency_ms: 14 })
      .mockResolvedValueOnce({ success: false, error: 'crc-2', latency_ms: 0 })
      .mockResolvedValueOnce({ success: false, error: 'offline-3', latency_ms: 0 });

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    const inspector = screen.getByTestId('workbench-inspector-panel');
    const testButton = within(inspector).getByRole('button', {
      name: 'workbench.device.actions.testConnection',
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(within(inspector).getByText('timeout-1')).toBeInTheDocument();
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(
        within(inspector).getByText('workbench.device.messages.testSuccess'),
      ).toBeInTheDocument();
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(within(inspector).getByText('crc-2')).toBeInTheDocument();
    });

    fireEvent.click(testButton);
    await waitFor(() => {
      expect(within(inspector).getByText('offline-3')).toBeInTheDocument();
      expect(within(inspector).getAllByRole('listitem')).toHaveLength(3);
    });

    expect(within(inspector).queryByText('timeout-1')).not.toBeInTheDocument();
  });

  it('shows connect and probe phase diagnostics separately in the inspector timeline', async () => {
    mockTestConnectionMutation.mutateAsync
      .mockReset()
      .mockResolvedValueOnce({
        success: false,
        error: 'probe failed',
        latency_ms: 21,
        can_activate: false,
        can_collect: false,
        connect: {
          status: 'success',
          message: 'TCP ready',
          latency_ms: 8,
        },
        probe: {
          status: 'failed',
          error: 'CRC mismatch',
          latency_ms: 13,
        },
      });

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    const inspector = screen.getByTestId('workbench-inspector-panel');

    fireEvent.click(
      within(inspector).getByRole('button', {
        name: 'workbench.device.actions.testConnection',
      }),
    );

    await waitFor(() => {
      const timelineEntry = within(inspector).getByText('probe failed').closest('li');
      expect(timelineEntry).not.toBeNull();
      expect(within(timelineEntry as HTMLElement).getByText('workbench.device.inspector.phases.connect')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('TCP ready')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('workbench.device.inspector.phases.probe')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('CRC mismatch')).toBeInTheDocument();
      expect(within(timelineEntry as HTMLElement).getByText('workbench.device.inspector.activationBlocked')).toBeInTheDocument();
    });
  });

  it('announces device notices through a polite live region', async () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.refresh' }));

    const statusRegions = await screen.findAllByRole('status');
    const liveNotice = statusRegions.find((element) =>
      element.textContent?.includes('workbench.device.messages.refreshed'),
    );

    expect(liveNotice).toHaveAttribute('aria-live', 'polite');
  });
});
