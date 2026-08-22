import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, vi } from 'vitest';
import DatalinkWorkbenchPage from '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type { Device, Mapping, Point, SourceRuleRecord, Tag } from '@/types/datalink';

const shellMocks = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockSourceRules: [] as SourceRuleRecord[],
}));

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockSourceRules,
} = shellMocks;

export function getShellUiMocks() {
  return shellMocks;
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
    mutateAsync: vi.fn().mockResolvedValue(undefined),
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
    isLoading: false,
    isSuccess: true,
  }),
  useCreateSourceRuleMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useUpdateSourceRuleMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDeleteSourceRuleMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useEnableSourceRuleMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDisableSourceRuleMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('@/services/datalink', () => ({
  modbusShareAPI: {
    status: vi.fn().mockResolvedValue({
      enabled: false,
      port: 5020,
      address: '',
      bind_state: 'fail',
      mapping_count: 0,
    }),
    listMappings: vi.fn().mockResolvedValue([]),
    start: vi.fn(),
    stop: vi.fn(),
    upsertMapping: vi.fn(),
    deleteMapping: vi.fn(),
    writeTagValue: vi.fn(),
    sync: vi.fn(),
  },
  tagAPI: {
    batchCreate: vi.fn().mockResolvedValue({ created: [], errors: [] }),
  },
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

export function registerShellUiFixtures() {
  beforeEach(() => {
    vi.clearAllMocks();

    mockDevices.splice(0, mockDevices.length, {
      id: 'device-1',
      name: 'Mixer PLC',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{"host":"192.168.1.10","port":502,"slave_id":7,"timeout":5}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '',
      updated_at: '',
    });

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
}
