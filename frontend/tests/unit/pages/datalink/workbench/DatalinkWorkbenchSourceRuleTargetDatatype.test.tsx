import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type {
  Device,
  Mapping,
  Point,
  SourceRuleRecord,
  Tag,
} from '@/types/datalink';

const {
  mockDevices,
  mockPoints,
  mockMappings,
  mockTags,
  mockSourceRules,
  mockCreateSourceRuleMutation,
  mockUpdateSourceRuleMutation,
  mockRuntimeStatus,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockMappings: [] as Mapping[],
  mockTags: [] as Tag[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockCreateSourceRuleMutation: { mutateAsync: vi.fn(), isPending: false },
  mockUpdateSourceRuleMutation: { mutateAsync: vi.fn(), isPending: false },
  mockRuntimeStatus: vi.fn(),
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock('@/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevices, isLoading: false }),
  useCreateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestDraftConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useToggleDeviceStatusMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeletePointMutation: () => ({ mutate: vi.fn(), mutateAsync: vi.fn(), isPending: false }),
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
  useDeleteSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useEnableSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDisableSourceRuleMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockMappings, isLoading: false }),
  useCreateMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteMappingMutation: () => ({ mutateAsync: vi.fn(), mutate: vi.fn(), isPending: false }),
}));

vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mockTags, isLoading: false }),
  useCreateTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('@/hooks/datalink/useRuntimeStream', () => ({
  useRuntimeStream: () => ({ liveValues: {} }),
}));

vi.mock('@/services/datalink', () => ({
  runtimeAPI: {
    getStatus: mockRuntimeStatus,
  },
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

async function openSourceStep() {
  renderPage();
  fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
  fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
  await screen.findByLabelText('workbench.source.planner.startAddress');
}

async function switchToDatabaseTarget() {
  fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.output/ }));
  fireEvent.click(screen.getByRole('tab', { name: 'workbench.output.targetSwitcher.database' }));
  fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
  await screen.findByLabelText('workbench.source.planner.startAddress');
}

describe('DatalinkWorkbench source-rule target datatype fields', () => {
  beforeEach(() => {
    vi.clearAllMocks();
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
    });
    mockPoints.splice(0, mockPoints.length);
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
    mockSourceRules.splice(0, mockSourceRules.length);
    mockCreateSourceRuleMutation.mutateAsync.mockImplementation(async (payload) => payload);
    mockUpdateSourceRuleMutation.mutateAsync.mockImplementation(async ({ data }) => data);
    mockRuntimeStatus.mockResolvedValue({ collectors: [] });
  });

  it('persists target datatype and scale fields when creating a source rule', async () => {
    await openSourceStep();

    fireEvent.change(screen.getByLabelText('workbench.source.planner.targetDataType'), {
      target: { value: 'float64' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.scaleMultiplier'), {
      target: { value: '0.1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.scaleOffset'), {
      target: { value: '5' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    fireEvent.click(screen.getByTestId('source-desk-tab-inspect'));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    );

    await waitFor(() => {
      expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });
    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        target_data_type: 'float64',
        scale_multiplier: 0.1,
        scale_offset: 5,
      }),
    );
  });

  it('loads and saves target datatype and scale fields when editing a persisted rule', async () => {
    mockSourceRules.splice(0, mockSourceRules.length, {
      id: 'persisted-rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'int16',
      target_data_type: 'float64',
      scale_multiplier: 0.5,
      scale_offset: 2,
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: '',
      skipped_addresses: [],
      created_at: '',
      updated_at: '2026-03-19T00:00:00Z',
    });

    await openSourceStep();

    fireEvent.click(screen.getByTestId('source-rule-layer-tab-rules'));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }));

    expect(screen.getAllByLabelText('workbench.source.planner.targetDataType').at(-1)).toHaveValue(
      'float64',
    );
    expect(screen.getAllByLabelText('workbench.source.planner.scaleMultiplier').at(-1)).toHaveValue(
      0.5,
    );
    expect(screen.getAllByLabelText('workbench.source.planner.scaleOffset').at(-1)).toHaveValue(
      2,
    );

    fireEvent.change(screen.getAllByLabelText('workbench.source.planner.targetDataType').at(-1)!, {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getAllByLabelText('workbench.source.planner.scaleMultiplier').at(-1)!, {
      target: { value: '2' },
    });
    fireEvent.change(screen.getAllByLabelText('workbench.source.planner.scaleOffset').at(-1)!, {
      target: { value: '-1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }));

    await waitFor(() => {
      expect(mockUpdateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith({
        id: 'persisted-rule-1',
        data: expect.objectContaining({
          target_data_type: 'float32',
          scale_multiplier: 2,
          scale_offset: -1,
        }),
      });
    });
  });

  it('applies database-aware defaults without overriding explicit planner choices', async () => {
    await openSourceStep();
    await switchToDatabaseTarget();

    expect(screen.getByTestId('source-planner-database-defaults')).toBeInTheDocument();
    expect(screen.getByLabelText('workbench.source.planner.namingPrefix')).toHaveValue('MBT_ROW');
    expect(screen.getByLabelText('workbench.source.planner.targetDataType')).toHaveValue('float64');
    expect(screen.getByLabelText('workbench.source.planner.scaleMultiplier')).toHaveValue(1);
    expect(screen.getByLabelText('workbench.source.planner.scaleOffset')).toHaveValue(0);

    fireEvent.change(screen.getByLabelText('workbench.source.planner.namingPrefix'), {
      target: { value: 'CUSTOM_ROW' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.targetDataType'), {
      target: { value: 'bool' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.scaleMultiplier'), {
      target: { value: '2.5' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'uint32' },
    });

    expect(screen.getByLabelText('workbench.source.planner.namingPrefix')).toHaveValue('CUSTOM_ROW');
    expect(screen.getByLabelText('workbench.source.planner.targetDataType')).toHaveValue('bool');
    expect(screen.getByLabelText('workbench.source.planner.scaleMultiplier')).toHaveValue(2.5);
    expect(screen.getByLabelText('workbench.source.planner.scaleOffset')).toHaveValue(0);
  });
});
