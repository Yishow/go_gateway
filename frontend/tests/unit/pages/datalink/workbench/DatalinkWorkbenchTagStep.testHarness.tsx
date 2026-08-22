import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { vi } from 'vitest';
import DatalinkWorkbenchPage from '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type { Device, Mapping, Point, Tag } from '@/types/datalink';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockCreateTagMutation,
  mockDeleteTagMutation,
  mockCreateMappingMutation,
  mockDeleteMappingMutation,
  mockDeletePointMutation,
  mockBatchCreate,
  mockMutation,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockCreateTagMutation: { mutateAsync: vi.fn(), isPending: false },
  mockDeleteTagMutation: { mutateAsync: vi.fn(), isPending: false },
  mockCreateMappingMutation: { mutateAsync: vi.fn(), isPending: false },
  mockDeleteMappingMutation: { mutateAsync: vi.fn(), isPending: false },
  mockDeletePointMutation: { mutateAsync: vi.fn(), isPending: false },
  mockBatchCreate: vi.fn(),
  mockMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

export {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockCreateTagMutation,
  mockDeleteTagMutation,
  mockCreateMappingMutation,
  mockDeleteMappingMutation,
  mockDeletePointMutation,
  mockBatchCreate,
};

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock('@/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevices, isLoading: false }),
  useCreateDeviceMutation: () => mockMutation(),
  useUpdateDeviceMutation: () => mockMutation(),
  useTestConnectionMutation: () => mockMutation(),
  useTestDraftConnectionMutation: () => mockMutation(),
  useToggleDeviceStatusMutation: () => mockMutation(),
}));

vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => mockMutation(),
  useDeletePointMutation: () => mockDeletePointMutation,
}));

vi.mock('@/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({ data: [], isLoading: false, isSuccess: true }),
  useCreateSourceRuleMutation: () => mockMutation(),
  useUpdateSourceRuleMutation: () => mockMutation(),
  useDeleteSourceRuleMutation: () => mockMutation(),
  useEnableSourceRuleMutation: () => mockMutation(),
  useDisableSourceRuleMutation: () => mockMutation(),
}));

vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
    refetch: vi.fn().mockResolvedValue({ data: mockTags }),
  }),
  useCreateTagMutation: () => mockCreateTagMutation,
  useDeleteTagMutation: () => mockDeleteTagMutation,
}));

vi.mock('@/services/datalink', () => ({
  tagAPI: { batchCreate: mockBatchCreate },
}));

vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockMappings, isLoading: false }),
  useCreateMappingMutation: () => mockCreateMappingMutation,
  useDeleteMappingMutation: () => mockDeleteMappingMutation,
}));

const defaultDevice: Device = {
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
};

const defaultPoints: Point[] = [
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
];

export function resetTagStepMocks() {
  vi.clearAllMocks();
  mockDevices.splice(0, mockDevices.length, { ...defaultDevice });
  mockPoints.splice(0, mockPoints.length, ...defaultPoints.map((point) => ({ ...point })));
  mockTags.splice(0, mockTags.length);
  mockMappings.splice(0, mockMappings.length);
  mockBatchCreate.mockImplementation(async (tags: Array<{ key: string }>) => {
    const created = tags.map((tag) => `tag-${tag.key}`);
    for (const tag of tags) {
      mockTags.push({
        id: `tag-${tag.key}`,
        key: tag.key,
        display_name: tag.key,
        description: '',
        data_type: 'int16',
        unit: '',
        labels: null,
        status: 'draft',
        created_at: '',
        updated_at: '',
      });
    }
    return { created, errors: [] };
  });
  mockCreateTagMutation.mutateAsync.mockImplementation(async (request: { key: string }) => ({
    id: `tag-${request.key}`,
    key: request.key,
    display_name: request.key,
    description: '',
    data_type: 'int16',
    unit: '',
    labels: null,
    status: 'draft',
    created_at: '',
    updated_at: '',
  }));
  mockDeleteTagMutation.mutateAsync.mockResolvedValue(undefined);
  mockCreateMappingMutation.mutateAsync.mockResolvedValue({
    id: 'mapping-created',
    point_id: 'point-1',
    tag_id: 'tag-created',
    enabled: true,
    transform_pipeline: '',
    created_at: '',
    updated_at: '',
  });
  mockDeletePointMutation.mutateAsync.mockResolvedValue(undefined);
}

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
