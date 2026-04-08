import { fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import type {
  Device,
  Mapping,
  Point,
  SourceRuleRecord,
  Tag,
} from '../../../../types/datalink';

const {
  mockDevices,
  mockPoints,
  mockMappings,
  mockTags,
  mockSourceRules,
  mockSourceRulesError,
  mockSourceRulesLoading,
  mockSourceRulesRefetch,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockMappings: [] as Mapping[],
  mockTags: [] as Tag[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockSourceRulesError: { value: null as Error | null },
  mockSourceRulesLoading: { value: false },
  mockSourceRulesRefetch: vi.fn(),
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockDevices,
    isLoading: false,
  }),
  useCreateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestDraftConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useToggleDeviceStatusMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: () => ({ data: mockPoints, isLoading: false }),
  useCreatePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeletePointMutation: () => ({ mutate: vi.fn(), mutateAsync: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useSourceRules', () => ({
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

vi.mock('../../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockMappings, isLoading: false }),
  useCreateMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteMappingMutation: () => ({ mutateAsync: vi.fn(), mutate: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mockTags, isLoading: false }),
  useCreateTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
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

function openSourceStepForDevice() {
  fireEvent.click(screen.getByRole('button', { name: /workbench\.steps\.source/ }));
  fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
}

describe('SourceCanvasSection status states', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    mockSourceRulesError.value = null;
    mockSourceRulesLoading.value = false;
    mockSourceRulesRefetch.mockReset();
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
    } satisfies Device);
    mockPoints.splice(0, mockPoints.length);
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
    mockSourceRules.splice(0, mockSourceRules.length);
  });

  it('shows a loading surface before source planning data is ready', () => {
    mockSourceRulesLoading.value = true;
    renderPage();
    openSourceStepForDevice();
    expect(screen.getByTestId('source-loading-state')).toHaveTextContent(
      'workbench.source.loading',
    );
    expect(screen.queryByLabelText('workbench.source.planner.startAddress')).not.toBeInTheDocument();
  });

  it('shows an error surface when source planning data fails to load', () => {
    mockSourceRulesError.value = new Error('source rules exploded');
    renderPage();
    openSourceStepForDevice();
    expect(screen.getByTestId('source-error-state')).toHaveTextContent(
      'workbench.source.loadFailed',
    );
    expect(screen.getByTestId('source-error-state')).toHaveTextContent('source rules exploded');
    expect(screen.queryByLabelText('workbench.source.planner.startAddress')).not.toBeInTheDocument();
  });

  it('retries source planning data loading from the error surface', () => {
    mockSourceRulesError.value = new Error('source rules exploded');
    renderPage();
    openSourceStepForDevice();
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.retry' }));
    expect(mockSourceRulesRefetch).toHaveBeenCalledTimes(1);
  });
});
