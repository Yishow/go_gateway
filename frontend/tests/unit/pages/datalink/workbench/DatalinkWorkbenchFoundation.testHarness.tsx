import { act, render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { useLocation } from 'react-router-dom';
import { expect, vi } from 'vitest';
import App from '@/App';
import '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type { Mapping, Point, SourceRuleRecord, Tag } from '@/types/datalink';

type MockDevice = {
  id: string; name: string; description: string; protocol: string;
  status: 'draft' | 'active' | 'disabled'; connection_config: string;
  last_test_at: string | null; last_test_success: boolean | null;
  last_test_error: string; created_at: string; updated_at: string;
};

const {
  mockDevices, mockPoints, mockTags, mockMappings, mockSourceRules,
  mockSourceRulesError, mockSourceRulesLoading, mockSourceRulesRefetch,
  mockCreateDeviceMutation, mockTestDraftConnectionMutation,
  mockUpdateDeviceMutation, mockTestConnectionMutation, mockMutation,
} = vi.hoisted(() => ({
  mockDevices: [] as MockDevice[], mockPoints: [] as Point[], mockTags: [] as Tag[],
  mockMappings: [] as Mapping[], mockSourceRules: [] as SourceRuleRecord[],
  mockSourceRulesError: { value: null as Error | null },
  mockSourceRulesLoading: { value: false }, mockSourceRulesRefetch: vi.fn(),
  mockCreateDeviceMutation: { mutateAsync: vi.fn(), isPending: false },
  mockTestDraftConnectionMutation: { mutateAsync: vi.fn(), isPending: false },
  mockUpdateDeviceMutation: { mutateAsync: vi.fn(), isPending: false },
  mockTestConnectionMutation: { mutateAsync: vi.fn(), isPending: false },
  mockMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

export {
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
};

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/contexts/ThemeContext', () => ({
  ThemeProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));
vi.mock('@/contexts/ToastContext', () => ({
  ToastProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));
vi.mock('@/components/CardMinimizeProvider', () => ({
  CardMinimizeProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));
vi.mock('@/pages/TestPage', () => ({ default: () => <div data-testid="test-page-mock">test-page</div> }));
vi.mock('@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({
  default: function MockWorkbenchV2Root() {
    const location = useLocation();
    return <div data-testid="workbench-v2-root">
      <span data-testid="workbench-v2-location-pathname">{location.pathname}</span>
      <span data-testid="workbench-v2-location-search">{location.search}</span>
    </div>;
  },
}));
vi.mock('@/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevices, isLoading: false }),
  useCreateDeviceMutation: () => mockCreateDeviceMutation,
  useTestDraftConnectionMutation: () => mockTestDraftConnectionMutation,
  useUpdateDeviceMutation: () => mockUpdateDeviceMutation,
  useTestConnectionMutation: () => mockTestConnectionMutation,
  useToggleDeviceStatusMutation: () => mockMutation(),
}));
vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id ? mockPoints.filter((point) => point.device_id === filters.device_id) : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => mockMutation(),
  useDeletePointMutation: () => mockMutation(),
}));
vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mockTags, isLoading: false, refetch: vi.fn().mockResolvedValue({ data: mockTags }) }),
  useCreateTagMutation: () => mockMutation(),
  useDeleteTagMutation: () => mockMutation(),
}));
vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockMappings, isLoading: false }),
  useCreateMappingMutation: () => mockMutation(),
  useDeleteMappingMutation: () => mockMutation(),
}));
vi.mock('@/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id ? mockSourceRules.filter((rule) => rule.device_id === filters.device_id) : [],
    error: mockSourceRulesError.value, isError: mockSourceRulesError.value !== null,
    isLoading: mockSourceRulesLoading.value,
    isSuccess: !mockSourceRulesLoading.value && mockSourceRulesError.value === null,
    refetch: mockSourceRulesRefetch,
  }),
  useCreateSourceRuleMutation: () => mockMutation(),
  useUpdateSourceRuleMutation: () => mockMutation(),
  useDeleteSourceRuleMutation: () => mockMutation(),
  useEnableSourceRuleMutation: () => mockMutation(),
  useDisableSourceRuleMutation: () => mockMutation(),
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
  buildDashboardModalRedirect: (intent: string) => '/mock/dashboard/' + intent,
  buildLegacyMigrationRedirect: (intent: string) => '/mock/legacy/' + intent,
  buildLocalModbusCompatRedirect: (section?: string | null) =>
    '/studio?step=output&target=modbus' + (section ? '&section=' + section : ''),
}));
vi.mock('@/services/datalink', async () => {
  const actual = await vi.importActual<typeof import('@/services/datalink')>('@/services/datalink');
  return {
    ...actual,
    protocolAPI: { list: vi.fn().mockResolvedValue([]) },
    settingsAPI: { get: vi.fn().mockResolvedValue({}) },
  };
});

const defaultDevices: typeof mockDevices = [
  { id: 'device-1', name: 'Mixer PLC', description: 'Main line', protocol: 'modbus_tcp', status: 'active',
    connection_config: '{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}', last_test_at: null,
    last_test_success: true, last_test_error: '', created_at: '', updated_at: '' },
  { id: 'device-2', name: 'Backup PLC', description: 'Fallback', protocol: 'modbus_rtu', status: 'draft',
    connection_config: '{"serial_port":"COM3","baud_rate":9600,"data_bits":8,"stop_bits":1,"parity":"none","slave_id":1,"timeout":5}',
    last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' },
];
const defaultPoints: Point[] = [
  { id: 'point-1', device_id: 'device-1', name: 'Flow Sensor', description: '', data_type: 'int16', address: '40001',
    enabled: true, polling_group_id: '', last_value: null, last_read_at: '', last_error: '', error_count: 0, created_at: '', updated_at: '' },
  { id: 'point-2', device_id: 'device-1', name: 'Pressure Sensor', description: '', data_type: 'int16', address: '40002',
    enabled: true, polling_group_id: '', last_value: null, last_read_at: '', last_error: '', error_count: 0, created_at: '', updated_at: '' },
];
const defaultTags: Tag[] = [
  { id: 'tag-1', key: 'TAG_40001', display_name: 'Flow Sensor', description: '', data_type: 'int16', unit: '', labels: null,
    status: 'draft', created_at: '', updated_at: '' },
];
const defaultMappings: Mapping[] = [
  { id: 'mapping-1', point_id: 'point-1', tag_id: 'tag-1', enabled: true, transform_pipeline: '',
    created_at: '', updated_at: '' },
];
const defaultCreatedDevice = {
  id: 'device-new', name: 'Browser Smoke PLC', description: 'Smoke flow', protocol: 'modbus_tcp', status: 'draft',
  connection_config: '{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}', last_test_at: null,
  last_test_success: null, last_test_error: '', created_at: '', updated_at: '',
};

export function resetFoundationMocks() {
  window.history.pushState({}, '', '/studio');
  mockSourceRules.splice(0, mockSourceRules.length);
  mockSourceRulesError.value = null; mockSourceRulesLoading.value = false; mockSourceRulesRefetch.mockReset();
  mockDevices.splice(0, mockDevices.length, ...defaultDevices.map((device) => ({ ...device })));
  mockCreateDeviceMutation.mutateAsync.mockReset();
  mockCreateDeviceMutation.mutateAsync.mockResolvedValue({ ...defaultCreatedDevice });
  mockTestDraftConnectionMutation.mutateAsync.mockReset(); mockTestDraftConnectionMutation.isPending = false;
  mockUpdateDeviceMutation.mutateAsync.mockReset(); mockUpdateDeviceMutation.mutateAsync.mockResolvedValue(undefined);
  mockTestConnectionMutation.mutateAsync.mockReset(); mockTestConnectionMutation.isPending = false;
  mockPoints.splice(0, mockPoints.length, ...defaultPoints.map((point) => ({ ...point })));
  mockTags.splice(0, mockTags.length, ...defaultTags.map((tag) => ({ ...tag })));
  mockMappings.splice(0, mockMappings.length, ...defaultMappings.map((mapping) => ({ ...mapping })));
}

export async function renderApp() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const { container } = render(<QueryClientProvider client={queryClient}><App /></QueryClientProvider>);
  await act(async () => { expect(container).toBeInTheDocument(); });
  return { container };
}
