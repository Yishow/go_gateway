import { fireEvent, render, screen, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { SOURCE_TEMPLATE_STORAGE_KEY } from '../../../../features/datalink/sourceTemplateStorage';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import { WORKBENCH_EXPERIMENT_PHASES } from '../workbenchExperimentContract';
import * as experimentContract from '../workbenchExperimentContract';
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
  mockCreatePointMutation,
  mockCreateSourceRuleMutation,
  mockUpdateSourceRuleMutation,
  mockDeleteSourceRuleMutation,
  mockEnableSourceRuleMutation,
  mockDisableSourceRuleMutation,
  mockToggleDeviceStatusMutation,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockMappings: [] as Mapping[],
  mockTags: [] as Tag[],
  mockSourceRules: [] as SourceRuleRecord[],
  mockCreatePointMutation: { mutateAsync: vi.fn(), isPending: false },
  mockCreateSourceRuleMutation: { mutateAsync: vi.fn(), isPending: false },
  mockUpdateSourceRuleMutation: { mutateAsync: vi.fn(), isPending: false },
  mockDeleteSourceRuleMutation: { mutateAsync: vi.fn(), isPending: false },
  mockEnableSourceRuleMutation: { mutateAsync: vi.fn(), isPending: false },
  mockDisableSourceRuleMutation: { mutateAsync: vi.fn(), isPending: false },
  mockToggleDeviceStatusMutation: { mutateAsync: vi.fn(), isPending: false },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('../../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevices, isLoading: false }),
  useCreateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestDraftConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useToggleDeviceStatusMutation: () => ({
    mutateAsync: mockToggleDeviceStatusMutation.mutateAsync,
    isPending: mockToggleDeviceStatusMutation.isPending,
  }),
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => mockCreatePointMutation,
  useDeletePointMutation: () => ({ mutateAsync: vi.fn(), mutate: vi.fn(), isPending: false }),
}));

vi.mock('../../../../hooks/datalink/useSourceRules', () => ({
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

function openSourceForMixer() {
  renderPage();
  fireEvent.click(screen.getByRole('tab', { name: /workbench.steps.source/ }));
  fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
}

describe('MuiSourceCommandDeck reopened Phase 2', () => {
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
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
    mockSourceRules.splice(0, mockSourceRules.length);
    mockCreatePointMutation.mutateAsync.mockResolvedValue(undefined);
    mockCreateSourceRuleMutation.mutateAsync.mockImplementation(async (payload) => payload);
    mockUpdateSourceRuleMutation.mutateAsync.mockImplementation(async ({ data }) => data);
    mockDeleteSourceRuleMutation.mutateAsync.mockResolvedValue(undefined);
    mockEnableSourceRuleMutation.mutateAsync.mockResolvedValue(undefined);
    mockDisableSourceRuleMutation.mutateAsync.mockResolvedValue(undefined);
    mockToggleDeviceStatusMutation.mutateAsync.mockResolvedValue(undefined);
    mockToggleDeviceStatusMutation.isPending = false;
  });

  it('aligns phase2 scenario focus with the reopened source flow', () => {
    const phase2 = WORKBENCH_EXPERIMENT_PHASES.find((phase) => phase.id === 'phase2');

    expect(phase2?.scenarioFocus).toEqual([
      'create-rule',
      'apply-template',
      'plan-live-link',
      'stale-preview-recovery',
      'handoff-tag',
    ]);
  });

  it('re-exports the source compare contract so phase2 compare gates can consume it', () => {
    expect('WORKBENCH_SOURCE_COMPARE_SCENARIOS' in experimentContract).toBe(true);
    expect('WORKBENCH_SOURCE_COMPARE_ACCEPTANCE' in experimentContract).toBe(true);
    expect('WORKBENCH_SOURCE_COMPARE_ARCHETYPES' in experimentContract).toBe(true);
    expect('WORKBENCH_SOURCE_COMPARE_CRITICAL_TASK' in experimentContract).toBe(true);
  });

  it('moves desk controls and handoff into a dedicated incident command surface', () => {
    openSourceForMixer();

    const commandPanel = screen.getByTestId('source-incident-command-panel');
    const handoffPanel = screen.getByTestId('source-incident-handoff-panel');
    const workboard = screen.getByTestId('source-incident-workboard');

    expect(within(commandPanel).getByTestId('source-desk-tab-inspect')).toBeInTheDocument();
    expect(within(commandPanel).getByTestId('source-desk-tab-build')).toBeInTheDocument();
    expect(within(commandPanel).getByTestId('source-desk-tab-triage')).toBeInTheDocument();
    expect(within(handoffPanel).getByTestId('source-handoff-tag-btn')).toBeInTheDocument();
    expect(within(workboard).getByTestId('sentry-source-workspace')).toBeInTheDocument();
    expect(screen.queryByTestId('source-mode-strip')).not.toBeInTheDocument();
    expect(screen.queryByTestId('source-handoff-dock')).not.toBeInTheDocument();
  });

  it('surfaces stale template auto-recovery in the source workspace', () => {
    localStorage.setItem(
      SOURCE_TEMPLATE_STORAGE_KEY,
      JSON.stringify([
        {
          id: 'legacy-template',
          name: 'Legacy Template',
          dataType: 'int16',
          count: 2,
          startAddress: '40001',
          updatedAt: '2026-03-01T00:00:00Z',
          version: 1,
        },
      ]),
    );

    openSourceForMixer();

    expect(screen.getByTestId('source-template-warning')).toHaveTextContent(
      'workbench.source.templates.recoveredLegacy',
    );
  });

  it('prioritizes template apply warnings over mount-time recovery warnings', () => {
    localStorage.setItem(
      SOURCE_TEMPLATE_STORAGE_KEY,
      JSON.stringify([
        {
          id: 'legacy-template',
          name: 'Legacy Template',
          dataType: 'int16',
          count: 2,
          startAddress: '40001',
          updatedAt: '2026-03-01T00:00:00Z',
          version: 1,
        },
        {
          id: 'mismatch-template',
          name: 'Mismatch Template',
          dataType: 'int16',
          count: 2,
          startAddress: '40001',
          preferredViewMode: 'plan',
          capabilitySnapshot: {
            protocol: 'mqtt',
            addressBase: 'topic-based',
            wordOrder: 'not-applicable',
          },
          updatedAt: '2026-04-01T00:00:00Z',
          lastUsedAt: '2026-04-01T00:00:00Z',
          version: 3,
        },
      ]),
    );

    openSourceForMixer();

    expect(screen.getByTestId('source-template-warning')).toHaveTextContent(
      'workbench.source.templates.recoveredLegacy',
    );

    fireEvent.click(screen.getByTestId('source-toolbar-more-trigger'));
    fireEvent.click(
      within(screen.getByTestId('source-toolbar-more-menu')).getByRole('menuitem', {
        name: 'workbench.source.toolbar.loadTemplate',
      }),
    );
    fireEvent.click(screen.getByRole('button', { name: 'Mismatch Template' }));

    expect(screen.getByTestId('source-template-warning')).toHaveTextContent(
      'workbench.source.templates.warning.message',
    );
  });
});
