import { fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, vi } from 'vitest';
import DatalinkWorkbenchPage from '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type {
  Device,
  Mapping,
  Point,
  SourceRuleRecord,
  Tag,
} from '@/types/datalink';

const sourceMocks = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockMappings: [] as Mapping[],
  mockTags: [] as Tag[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockCreatePointMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockDeletePointMutation: {
    mutate: vi.fn(),
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockCreateSourceRuleMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockUpdateSourceRuleMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockDeleteSourceRuleMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockEnableSourceRuleMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockDisableSourceRuleMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
  mockToggleDeviceStatusMutation: {
    mutateAsync: vi.fn(),
    isPending: false,
  },
}));
const {
  mockDevices,
  mockPoints,
  mockMappings,
  mockTags,
  mockSourceRules,
  mockCreatePointMutation,
  mockDeletePointMutation,
  mockCreateSourceRuleMutation,
  mockUpdateSourceRuleMutation,
  mockDeleteSourceRuleMutation,
  mockEnableSourceRuleMutation,
  mockDisableSourceRuleMutation,
  mockToggleDeviceStatusMutation,
} = sourceMocks;

export function getSourceStepMocks() {
  return sourceMocks;
}

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
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
    mutateAsync: mockToggleDeviceStatusMutation.mutateAsync,
    isPending: mockToggleDeviceStatusMutation.isPending,
  }),
}));

vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => mockCreatePointMutation,
  useDeletePointMutation: () => mockDeletePointMutation,
}));

vi.mock('@/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockSourceRules.filter((rule) => rule.device_id === filters.device_id)
      : [],
    isLoading: false,
    isSuccess: true,
  }),
  useCreateSourceRuleMutation: () => mockCreateSourceRuleMutation,
  useUpdateSourceRuleMutation: () => mockUpdateSourceRuleMutation,
  useDeleteSourceRuleMutation: () => mockDeleteSourceRuleMutation,
  useEnableSourceRuleMutation: () => mockEnableSourceRuleMutation,
  useDisableSourceRuleMutation: () => mockDisableSourceRuleMutation,
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
    mutate: vi.fn(),
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

export function renderPage() {
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

export function openSourceRuleLayerRulesTab() {
    fireEvent.click(screen.getByTestId('source-rule-layer-tab-rules'));
  }

export function registerSourceStepFixtures() {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    mockDevices.splice(0, mockDevices.length, {
      id: 'device-1',
      name: 'Mixer PLC',
      description: 'Line A mixer',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    }, {
      id: 'device-2',
      name: 'Fatek Cell',
      description: 'Line B fatek',
      protocol: 'fatek_fbs',
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
      name: 'Existing Pressure',
      description: '',
      data_type: 'int16',
      address: '40005',
      enabled: true,
      polling_group_id: '',
      last_value: 12,
      last_read_at: '',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    });
    mockCreatePointMutation.mutateAsync.mockResolvedValue(undefined);
    mockCreateSourceRuleMutation.mutateAsync.mockImplementation(async (payload) => payload);
    mockUpdateSourceRuleMutation.mutateAsync.mockImplementation(async ({ data }) => data);
    mockDeleteSourceRuleMutation.mutateAsync.mockResolvedValue(undefined);
    mockEnableSourceRuleMutation.mutateAsync.mockResolvedValue(undefined);
    mockDisableSourceRuleMutation.mutateAsync.mockResolvedValue(undefined);
    mockToggleDeviceStatusMutation.mutateAsync.mockResolvedValue(undefined);
    mockToggleDeviceStatusMutation.isPending = false;
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
    mockSourceRules.splice(0, mockSourceRules.length);
  });
}
