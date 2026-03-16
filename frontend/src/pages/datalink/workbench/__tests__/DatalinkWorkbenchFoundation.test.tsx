import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import App from '../../../../App';

const {
  mockDevices,
  mockCreateDeviceMutation,
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

  it('selects a device from the device step and advances to source planning', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.device.actions.continue' }),
    );

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

  it('context bar enables quick-action buttons after selecting a device on step 1', () => {
    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(
      screen.getByText('workbench.contextBar.actions.switchDevice'),
    ).toBeInTheDocument();
  });

  it('opens the create device form when there are no devices yet', () => {
    mockDevices.splice(0, mockDevices.length);

    renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));

    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.createTitle' }),
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

  it('uses the latest session test result tone in the context bar', async () => {
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
      expect(screen.getByTestId('context-bar-test-status')).toHaveTextContent('timeout-latest');
      expect(screen.getByTestId('context-bar-test-status')).toHaveClass('text-rose-300');
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
});
