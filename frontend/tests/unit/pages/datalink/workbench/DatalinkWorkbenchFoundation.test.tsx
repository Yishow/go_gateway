import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { useLocation } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import App from '@/App';
// DatalinkWorkbenchPage 在 App 路由中為 lazy chunk；靜態導入讓模組預先載入快取，
// renderApp 內的 act flush 才能在同步斷言前完成首次 lazy 渲染。
import '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type { Mapping, Point, SourceRuleRecord, Tag } from '@/types/datalink';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockSourceRules,
  mockSourceRulesError,
  mockSourceRulesLoading,
  mockSourceRulesRefetch,
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
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockSourceRulesError: { value: null as Error | null },
  mockSourceRulesLoading: { value: false },
  mockSourceRulesRefetch: vi.fn(),
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

vi.mock('@/pages/TestPage', () => ({
  default: () => <div data-testid="test-page-mock">test-page</div>,
}));

vi.mock('@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({
  default: function MockWorkbenchV2Root() {
    const location = useLocation();

    return (
      <div data-testid="workbench-v2-root">
        <span data-testid="workbench-v2-location-pathname">{location.pathname}</span>
        <span data-testid="workbench-v2-location-search">{location.search}</span>
      </div>
    );
  },
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
    refetch: vi.fn().mockResolvedValue({ data: mockTags }),
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

vi.mock('@/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockSourceRules.filter((rule) => rule.device_id === filters.device_id)
      : [],
    error: mockSourceRulesError.value,
    isError: mockSourceRulesError.value !== null,
    isLoading: mockSourceRulesLoading.value,
    isSuccess: !mockSourceRulesLoading.value && mockSourceRulesError.value === null,
    refetch: mockSourceRulesRefetch,
  }),
  useCreateSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useEnableSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDisableSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('@/router/gateway', () => ({
  GatewayCreateEntryRedirect: () => <div data-testid="gateway-create-entry-redirect-mock" />,
  GatewayEntryRoute: () => <div data-testid="gateway-entry-route-mock" />,
  GatewayExpertWorkbenchRoute: () => <div data-testid="gateway-expert-workbench-route-mock" />,
  GatewayQuickSetupRoute: () => <div data-testid="gateway-quick-setup-route-mock" />,
}));

vi.mock('@/features/datalink/legacyRoutes', () => ({
  buildWorkbenchV2EntryRedirect: () => '/studio/v2',
  buildWorkbenchRedirect: () => '/studio',
  buildDashboardModalRedirect: (intent: string) => `/mock/dashboard/${intent}`,
  buildLegacyMigrationRedirect: (intent: string) => `/mock/legacy/${intent}`,
  buildLocalModbusCompatRedirect: (section?: string | null) =>
    `/studio?step=output&target=modbus${section ? `&section=${section}` : ''}`,
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
    window.history.pushState({}, '', '/studio');
    mockSourceRules.splice(0, mockSourceRules.length);
    mockSourceRulesError.value = null;
    mockSourceRulesLoading.value = false;
    mockSourceRulesRefetch.mockReset();
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
    mockPoints.splice(
      0,
      mockPoints.length,
      {
        id: 'point-1',
        device_id: 'device-1',
        name: 'Flow Sensor',
        description: '',
        data_type: 'int16',
        address: '40001',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
      {
        id: 'point-2',
        device_id: 'device-1',
        name: 'Pressure Sensor',
        description: '',
        data_type: 'int16',
        address: '40002',
        enabled: true,
        polling_group_id: '',
        last_value: null,
        last_read_at: '',
        last_error: '',
        error_count: 0,
        created_at: '',
        updated_at: '',
      },
    );
    mockTags.splice(0, mockTags.length, {
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
    });
    mockMappings.splice(0, mockMappings.length, {
      id: 'mapping-1',
      point_id: 'point-1',
      tag_id: 'tag-1',
      enabled: true,
      transform_pipeline: '',
      created_at: '',
      updated_at: '',
    });
  });

  async function renderApp() {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });

    const { container } = render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>,
    );

    // 等待 lazy 路由 chunk（已預載入模組快取）完成首次渲染
    await act(async () => {
      expect(container).toBeInTheDocument();
    });

    return { container };
  }

  it('redirects /datalink/workbench into /studio/v2 and renders the V2 entry', async () => {
    window.history.pushState({}, '', '/datalink/workbench');

    await renderApp();

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.getByTestId('workbench-v2-location-pathname')).toHaveTextContent('/studio/v2');
    expect(window.location.pathname).toBe('/studio/v2');
  });

  it('keeps Step 1 search, protocol filter, and create action inside one primary toolbar', async () => {
    await renderApp();

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

  it('drops the decorative Step 1 hero block once the compact toolbar is available', async () => {
    await renderApp();

    expect(screen.queryByText('workbench.device.title')).not.toBeInTheDocument();
    expect(screen.queryByText('workbench.device.description')).not.toBeInTheDocument();
  });

  it('renders compact device rows with endpoint, health, and two primary capability hints', async () => {
    await renderApp();

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

  it('uses protocol-specific connection targets for compact row endpoints', async () => {
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

    await renderApp();

    const mqttRow = screen.getByTestId('device-row-device-3');
    expect(within(mqttRow).getByText('mqtt://broker.internal:1883')).toBeInTheDocument();
    expect(within(mqttRow).queryByText(/QoS 1/)).not.toBeInTheDocument();

    const mcRow = screen.getByTestId('device-row-device-4');
    expect(within(mcRow).getByText('10.0.0.20:5000')).toBeInTheDocument();
    expect(within(mcRow).queryByText('1/2 · binary')).not.toBeInTheDocument();
  });

  it('preserves deep-link query context when /datalink/workbench redirects into /studio/v2', async () => {
    window.history.pushState({}, '', '/datalink/workbench?step=output&target=database');

    await renderApp();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });
    expect(screen.getByTestId('workbench-v2-location-pathname')).toHaveTextContent('/studio/v2');
    expect(screen.getByTestId('workbench-v2-location-search')).toHaveTextContent(
      '?step=output&target=database',
    );
    expect(window.location.pathname).toBe('/studio/v2');
    expect(window.location.search).toContain('step=output');
    expect(window.location.search).toContain('target=database');
  });

  it('does not lock step navigation after a deep-link is applied', async () => {
    window.history.pushState({}, '', '/studio?step=output&target=database');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    // User clicks the device step in the step rail — must NOT be locked back to output.
    fireEvent.click(
      screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
    );

    expect(
      screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
    ).toHaveAttribute('aria-selected', 'true');
    expect(
      screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
    ).not.toHaveAttribute('aria-selected', 'true');
  });

  it('redirects legacy local modbus entry into the new workbench output step', async () => {
    window.history.pushState({}, '', '/datalink/local-modbus?section=settings');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(screen.getByTestId('active-output-target')).toHaveTextContent(
      'workbench.bottomSummary.targets.modbus',
    );
    expect(window.location.pathname).toBe('/studio');
    expect(window.location.search).toContain('step=output');
  });

  it('redirects nested legacy workbench entries into /studio while preserving query context', async () => {
    window.history.pushState({}, '', '/datalink/workbench/legacy-output?step=output&target=database');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(screen.getByTestId('active-output-target')).toHaveTextContent(
      'workbench.bottomSummary.targets.database',
    );
    expect(window.location.pathname).toBe('/studio');
    expect(window.location.search).toContain('step=output');
    expect(window.location.search).toContain('target=database');
  });

  it('renders the main operator flow through /studio', async () => {
    window.history.pushState({}, '', '/studio');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    expect(screen.getByText('Mixer PLC')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/studio');
  });

  it('keeps the normal operator workflow inside /studio from device through output', async () => {
    window.history.pushState({}, '', '/studio');

    await renderApp();

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.device/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    expect(window.location.pathname).toBe('/studio');

    const contextBar = screen.getByTestId('workbench-context-bar');
    fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));
    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.source/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeInTheDocument();
    });
    expect(window.location.pathname).toBe('/studio');

    fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }));
    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.tag/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(window.location.pathname).toBe('/studio');

    fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoOutput' }));
    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.output/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });
    expect(screen.getByTestId('active-output-target')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/studio');
  });

  it('keeps source engineering tools inside /studio instead of opening another workflow route', async () => {
    window.history.pushState({}, '', '/studio');

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-context-bar')).getByRole('button', {
        name: 'workbench.contextBar.actions.gotoSource',
      }),
    );

    await waitFor(() => {
      expect(
        screen.getByRole('tab', { name: /workbench\.steps\.source/ }),
      ).toHaveAttribute('aria-selected', 'true');
    });

    await waitFor(() => {
      expect(screen.getByTestId('source-toolbar-more-trigger')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('source-toolbar-more-trigger'));

    expect(screen.getByTestId('source-toolbar-more-menu')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/studio');
  });

  it('redirects /datalink into /studio', async () => {
    window.history.pushState({}, '', '/datalink');

    await renderApp();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    expect(screen.getByTestId('workbench-v2-location-pathname')).toHaveTextContent('/studio/v2');
    expect(window.location.pathname).toBe('/studio/v2');
  });

  it('renders /test without the legacy sidebar layout shell', async () => {
    window.history.pushState({}, '', '/test');

    await renderApp();

    expect(await screen.findByTestId('test-page-mock')).toBeInTheDocument();
    expect(window.location.pathname).toBe('/test');
  });

  it('redirects legacy test utility routes into /test', async () => {
    window.history.pushState({}, '', '/templates');

    await renderApp();

    await waitFor(() => {
      expect(screen.getByTestId('test-page-mock')).toBeInTheDocument();
    });

    expect(window.location.pathname).toBe('/test');
  });

  it('selects a device from the device step and advances to source planning', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(
      within(screen.getByTestId('workbench-context-bar')).getByRole('button', {
        name: 'workbench.contextBar.actions.gotoSource',
      }),
    );

    expect(screen.getByRole('tab', { name: /workbench\.steps\.source/ })).toHaveAttribute(
      'aria-selected',
      'true',
    );
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeInTheDocument();
      expect(screen.getByLabelText('workbench.source.planner.count')).toBeInTheDocument();
    });
  });

  it('shows a master-detail device panel after a device is selected', async () => {
    await renderApp();

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

  it('context bar includes step rail plus primary CTA after selecting a device on step 1', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const contextBar = screen.getByTestId('workbench-context-bar');
    const tabs = within(contextBar).getAllByRole('tab');
    expect(tabs.length).toBe(4);
    expect(
      within(contextBar).getByRole('button', {
        name: 'workbench.contextBar.actions.gotoSource',
      }),
    ).toBeInTheDocument();
  });

  it('opens the create device form inline when there are no devices yet', async () => {
    mockDevices.splice(0, mockDevices.length);

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.create' }));

    expect(screen.getByTestId('device-inline-editor')).toBeInTheDocument();
    expect(screen.queryByTestId('device-panel-overlay')).not.toBeInTheDocument();
    expect(
      screen.getByRole('heading', { name: 'workbench.device.panel.createTitle' }),
    ).toBeInTheDocument();
  });

  it('localizes protocol-specific connection option labels in the create drawer', async () => {
    mockDevices.splice(0, mockDevices.length);

    await renderApp();

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
        name: 'workbench.device.connection.dataFormats.abcd',
      }),
    ).toBeInTheDocument();
    expect(
      within(dataFormatSelect).getByRole('option', {
        name: 'workbench.device.connection.dataFormats.cdab',
      }),
    ).toBeInTheDocument();
  });

  it('keeps the created device selected while the list refreshes', async () => {
    mockDevices.splice(0, mockDevices.length);

    await renderApp();

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
    await renderApp();

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

  it('shows the inline editor in the detail column while keeping the device list visible', async () => {
    await renderApp();

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

    await renderApp();

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

  it('hides zero-value test timestamps in the inspector', async () => {
    mockDevices[0].last_test_at = '0001-01-01T00:00:00Z';

    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(
      screen.getByText('workbench.device.inspector.unknownTestTime'),
    ).toBeInTheDocument();
    expect(screen.queryByText('0001-01-01T00:00:00Z')).not.toBeInTheDocument();
  });

  it('renders selected device details in the right-side inspector', async () => {
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const inspector = screen.getByTestId('workbench-inspector-panel');
    expect(within(inspector).getByText('Mixer PLC')).toBeInTheDocument();
    expect(within(inspector).getByText('workbench.device.inspector.connectionSummary')).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.edit' })).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.testConnection' })).toBeInTheDocument();
    expect(within(inspector).getByRole('button', { name: 'workbench.device.actions.clone' })).toBeInTheDocument();
  });

  it('opens a clone drawer with connection defaults but requires a new device name', async () => {
    await renderApp();

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
    await renderApp();

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

  it('disables inspector test action while a connection test is already pending', async () => {
    mockTestConnectionMutation.isPending = true;

    await renderApp();

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

    await renderApp();

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

    await renderApp();

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

    await renderApp();

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
    await renderApp();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.device.actions.refresh' }));

    const statusRegions = await screen.findAllByRole('status');
    const liveNotice = statusRegions.find((element) =>
      element.textContent?.includes('workbench.device.messages.refreshed'),
    );

    expect(liveNotice).toHaveAttribute('aria-live', 'polite');
  });
});
