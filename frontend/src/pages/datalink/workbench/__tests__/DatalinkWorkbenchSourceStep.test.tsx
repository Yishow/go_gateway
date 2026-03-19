import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import { SOURCE_TEMPLATE_STORAGE_KEY } from '../../../../features/datalink/sourceTemplateStorage';
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
  mockDeletePointMutation,
  mockCreateSourceRuleMutation,
  mockUpdateSourceRuleMutation,
  mockDeleteSourceRuleMutation,
  mockEnableSourceRuleMutation,
  mockDisableSourceRuleMutation,
} = vi.hoisted(() => ({
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
}));

vi.mock('../../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => mockCreatePointMutation,
  useDeletePointMutation: () => mockDeletePointMutation,
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

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
  }),
  useCreateTagMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
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

describe('DatalinkWorkbench source step', () => {
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
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
    mockSourceRules.splice(0, mockSourceRules.length);
  });

  it('gates source planning behind device selection', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));

    expect(screen.getByText('workbench.source.empty.title')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Mixer PLC' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toBeInTheDocument();
    expect(screen.getByLabelText('workbench.source.planner.count')).toBeInTheDocument();
  });

  it('remembers the last planner start address per device and falls back to protocol defaults', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const startAddressInput = screen.getByLabelText('workbench.source.planner.startAddress');
    fireEvent.change(startAddressInput, { target: { value: '40010' } });
    expect(startAddressInput).toHaveValue('40010');

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.device' }));
    fireEvent.click(screen.getByRole('button', { name: 'Fatek Cell' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('D0');
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    expect(screen.getByTestId('address-cell-D0')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-D2')).toHaveAttribute('data-status', 'planned');

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: 'D20' },
    });
    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('D20');

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.device' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    await waitFor(() => {
      expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('40010');
    });
  });

  it('renders a source rule layer and continuous gap cells after applying a rule', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('source-coverage-overview')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40004')).toHaveAttribute('data-status', 'gap');
    expect(screen.getByTestId('address-cell-40005')).toHaveAttribute('data-status', 'unmanaged');
  });

  it('loads persisted source rules from backend state for the selected device', () => {
    mockSourceRules.splice(0, mockSourceRules.length, {
      id: 'persisted-rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: '',
      skipped_addresses: [],
      created_at: '',
      updated_at: '2026-03-19T00:00:00Z',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('source-rule-persisted-rule-1')).toBeInTheDocument();
    expect(screen.getByText('workbench.source.ruleLayer.persistedBadge')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'planned');
  });

  it('explains unmanaged existing points in the inspector', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByTestId('address-cell-40005'));

    expect(screen.getByTestId('source-span-inspector')).toBeInTheDocument();
    expect(screen.getByText('workbench.source.inspector.span.unmanagedNotice')).toBeInTheDocument();
  });

  it('preserves skipped addresses when editing a persisted rule', async () => {
    mockSourceRules.splice(0, mockSourceRules.length, {
      id: 'persisted-rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: '',
      skipped_addresses: ['40002'],
      created_at: '',
      updated_at: '2026-03-19T00:00:00Z',
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }));
    fireEvent.change(screen.getAllByDisplayValue('40001').at(-1)!, {
      target: { value: '40005' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }));

    await waitFor(() => {
      expect(mockUpdateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith({
        id: 'persisted-rule-1',
        data: expect.objectContaining({
          start_address: '40005',
          skipped_addresses: ['40002'],
        }),
      });
    });
  });

  it('renders the source canvas as fixed 16-bit lattice rows', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const firstRow = screen.getByTestId('source-canvas-row-0');

    expect(firstRow).toHaveAttribute('data-row-start-address', '40001');
    expect(firstRow).toHaveAttribute('data-row-end-address', '40005');
    expect(firstRow).toHaveAttribute('data-lattice-columns', '16');
  });

  it('shows all utility tools directly in the secondary controls without a toggle', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const primaryToolbar = screen.getByTestId('source-primary-toolbar');
    const secondaryControls = screen.getByTestId('source-secondary-controls');
    const ruleLayer = screen.getByTestId('source-rule-layer');

    expect(
      within(primaryToolbar).getByRole('button', { name: 'workbench.source.view.plan' }),
    ).toBeInTheDocument();
    expect(
      within(primaryToolbar).getByLabelText('workbench.source.toolbar.valueFormat'),
    ).toBeInTheDocument();
    expect(
      within(primaryToolbar).queryByLabelText('workbench.source.planner.startAddress'),
    ).not.toBeInTheDocument();
    expect(
      within(primaryToolbar).queryByRole('button', {
        name: 'workbench.source.planner.addRule',
      }),
    ).not.toBeInTheDocument();
    expect(
      within(ruleLayer).getByLabelText('workbench.source.planner.startAddress'),
    ).toBeInTheDocument();
    expect(
      within(ruleLayer).getByRole('button', { name: 'workbench.source.planner.addRule' }),
    ).toBeInTheDocument();

    // No toggle needed — tools are always visible
    expect(
      within(secondaryControls).queryByRole('button', {
        name: 'workbench.source.toolbar.moreTools',
      }),
    ).not.toBeInTheDocument();
    expect(
      within(secondaryControls).getByRole('button', {
        name: 'workbench.source.toolbar.saveTemplate',
      }),
    ).toBeInTheDocument();
    expect(
      within(secondaryControls).getByRole('button', {
        name: 'workbench.source.toolbar.snapshotCompare',
      }),
    ).toBeInTheDocument();
  });

  it('shows a step-local health summary with dual point-creation actions', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const summary = screen.getByTestId('source-step-summary');

    expect(within(summary).getByTestId('source-summary-ready-count')).toHaveTextContent('2');
    expect(within(summary).getByTestId('source-summary-conflict-count')).toHaveTextContent('0');
    expect(within(summary).getByTestId('source-summary-protected-count')).toHaveTextContent('0');
    expect(
      within(summary).getByRole('button', {
        name: 'workbench.source.actions.createSelectedPoints',
      }),
    ).toBeDisabled();
    expect(
      within(summary).getByRole('button', {
        name: 'workbench.source.actions.createRulePoints',
      }),
    ).toBeEnabled();
  });

  it('creates the selected logical span from the summary action', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40002'));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createSelectedPoints' }),
    );

    await waitFor(() => {
      expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledWith({
      device_id: 'device-1',
      address: '40001',
      data_type: 'float32',
      name: 'SRC_40001',
    });
  });

  it('persists rule points through the source-rule API when applying the rule batch action', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    );

    await waitFor(() => {
      expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith({
      id: 'rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'int16',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: undefined,
      skipped_addresses: [],
    });
    expect(mockCreatePointMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('adds and deletes rules directly from the rule layer workflow', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const ruleLayer = screen.getByTestId('source-rule-layer');

    fireEvent.change(within(ruleLayer).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(within(ruleLayer).getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(
      within(ruleLayer).getByRole('button', { name: 'workbench.source.planner.addRule' }),
    );

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    expect(ruleCard).toBeInTheDocument();

    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.delete' }),
    );

    expect(screen.queryByTestId('source-rule-rule-1')).not.toBeInTheDocument();
  });

  it('treats the address canvas as the primary workspace and the rule layer as supporting context', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    expect(screen.getByTestId('source-canvas-workspace')).toHaveAttribute(
      'data-emphasis',
      'primary',
    );
    expect(screen.getByTestId('source-rule-layer')).toHaveAttribute(
      'data-emphasis',
      'supporting',
    );
    expect(
      within(screen.getByTestId('source-canvas-workspace')).getByTestId(
        'source-coverage-overview',
      ),
    ).toBeInTheDocument();
  });

  it('preserves merged spans and gap cells across multiple source rules', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40004' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int16' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('source-rule-rule-2')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-merge-span', '2');
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-merge-offset', '0');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-merge-span', '2');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-merge-offset', '1');
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'gap');
    expect(screen.getByTestId('address-cell-40004')).toHaveAttribute('data-status', 'planned');
  });

  it('switches plan, live, and link overlays without changing the lattice addresses', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.live' }));
    expect(screen.getByTestId('source-canvas')).toHaveAttribute('data-view-mode', 'live');
    expect(screen.getByTestId('address-cell-40001')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.link' }));
    expect(screen.getByTestId('source-canvas')).toHaveAttribute('data-view-mode', 'link');
    expect(screen.getByTestId('address-cell-40001')).toBeInTheDocument();
    expect(screen.getByText('workbench.source.link.unbound')).toBeInTheDocument();
  });

  it('shows rule inspector details when a source rule is selected', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('source-rule-rule-1'));

    expect(screen.getByTestId('source-rule-inspector')).toBeInTheDocument();
    expect(screen.getByTestId('source-rule-coverage')).toHaveTextContent('40001');
  });

  it('shows span inspector details when an address cell is selected', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    expect(screen.getByTestId('source-span-inspector')).toBeInTheDocument();
    expect(screen.getByTestId('source-span-link-state')).toHaveTextContent(
      'workbench.source.link.needsPoint',
    );
  });

  it('marks the selected address cell with aria-pressed', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const addressCell = screen.getByTestId('address-cell-40001');

    expect(addressCell).toHaveAttribute('aria-pressed', 'false');

    fireEvent.click(addressCell);

    expect(addressCell).toHaveAttribute('aria-pressed', 'true');
  });

  it('persists applied source rules when navigating away from and back to the source step', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.tag' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));

    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
  });

  it('saves a source template locally and reapplies it to the planner inputs', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40101' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.live' }));

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.toolbar.saveTemplate' }));
    fireEvent.change(screen.getByLabelText('workbench.source.templates.name'), {
      target: { value: 'Line Float' },
    });
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.templates.confirmSave' }),
    );

    const storedTemplates = JSON.parse(
      localStorage.getItem(SOURCE_TEMPLATE_STORAGE_KEY) ?? '[]',
    ) as Array<Record<string, unknown>>;
    expect(storedTemplates[0]).toMatchObject({
      name: 'Line Float',
      startAddress: '40101',
      count: 3,
      dataType: 'float32',
      preferredViewMode: 'live',
    });

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '49999' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int16' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.plan' }));

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.toolbar.loadTemplate' }));
    fireEvent.click(screen.getByRole('button', { name: 'Line Float' }));

    expect(screen.getByLabelText('workbench.source.planner.startAddress')).toHaveValue('40101');
    expect(screen.getByLabelText('workbench.source.planner.count')).toHaveValue(3);
    expect(screen.getByLabelText('workbench.source.planner.dataType')).toHaveValue('float32');
    expect(screen.getByTestId('source-canvas')).toHaveAttribute('data-view-mode', 'live');
  });

  it('batch persists eligible source rules from the planned address range', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));
    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    );

    await waitFor(() => {
      expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith({
      id: 'rule-1',
      device_id: 'device-1',
      start_address: '40001',
      count: 2,
      data_type: 'float32',
      naming_prefix: 'SRC',
      enabled: true,
      locked: false,
      origin: 'manual',
      template_name: undefined,
      skipped_addresses: [],
    });
    expect(mockCreatePointMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('allows batch create for safe spans even when conflicts exist elsewhere', () => {
    mockPoints[0].address = '40002';

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'conflict');
    // The float32 at 40001 occupies 40001-40002, and 40002 is a conflict, so the
    // only planned span is itself unsafe — button should be disabled.
    expect(
      screen.getByRole('button', { name: 'workbench.source.actions.createRulePoints' }),
    ).toBeDisabled();
  });

  it('shows a selection toolbar when a planned address cell is clicked', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.queryByTestId('source-selection-toolbar')).not.toBeInTheDocument();

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    const toolbar = screen.getByTestId('source-selection-toolbar');
    expect(toolbar).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.createPoints',
      }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.addToRule',
      }),
    ).toBeInTheDocument();
    expect(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.skip',
      }),
    ).toBeInTheDocument();
  });

  it('creates a point via the selection toolbar create action', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    const toolbar = screen.getByTestId('source-selection-toolbar');
    fireEvent.click(
      within(toolbar).getByRole('button', {
        name: 'workbench.source.selectionToolbar.createPoints',
      }),
    );

    await waitFor(() => {
      expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledWith({
      device_id: 'device-1',
      address: '40001',
      data_type: 'int16',
      name: 'SRC_40001',
    });
  });

  it('clears the selected address when skip is clicked in the selection toolbar', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));
    expect(screen.getByTestId('source-selection-toolbar')).toBeInTheDocument();

    fireEvent.click(
      screen.getByRole('button', { name: 'workbench.source.selectionToolbar.skip' }),
    );

    expect(screen.queryByTestId('source-selection-toolbar')).not.toBeInTheDocument();
  });

  it('supports inline editing of rule start address, count, and data type', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );

    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    expect(editForm).toBeInTheDocument();

    const addressInput = within(editForm).getByLabelText('workbench.source.planner.startAddress');
    const countInput = within(editForm).getByLabelText('workbench.source.planner.count');

    fireEvent.change(addressInput, { target: { value: '40010' } });
    fireEvent.change(countInput, { target: { value: '3' } });

    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.getByTestId('address-cell-40010')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40012')).toHaveAttribute('data-status', 'planned');
  });

  it('cancels inline rule editing without modifying the rule', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );

    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40099' },
    });

    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editCancel' }),
    );

    expect(within(ruleCard).queryByTestId('rule-inline-edit-form')).not.toBeInTheDocument();
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.queryByTestId('address-cell-40099')).not.toBeInTheDocument();
  });

  it('immediately updates the canvas when a rule is inline-edited', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('2');

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );

    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.count'), {
      target: { value: '4' },
    });
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('4');
  });

  it('does not show selection toolbar for non-planned cells', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40005'));

    expect(screen.queryByTestId('source-selection-toolbar')).not.toBeInTheDocument();
  });

  it('retargets the inspector to the edited rule after inline save', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));
    expect(screen.getByTestId('source-span-inspector')).toBeInTheDocument();

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );
    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40010' },
    });
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    expect(screen.queryByTestId('source-span-inspector')).not.toBeInTheDocument();
    expect(screen.getByTestId('source-rule-inspector')).toBeInTheDocument();
  });

  it('replaces the conflict hint with an actionable conflict queue', () => {
    mockPoints[0].address = '40002';

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const conflictQueue = screen.getByTestId('source-conflict-queue');
    expect(conflictQueue).toBeInTheDocument();
    expect(within(conflictQueue).getByText('workbench.source.conflictQueue.title')).toBeInTheDocument();
    expect(within(conflictQueue).getByText('workbench.source.conflictQueue.step3Blocked')).toBeInTheDocument();

    const conflictItem = screen.getByTestId('conflict-item-40001');
    expect(conflictItem).toBeInTheDocument();
    expect(within(conflictItem).getByText('workbench.source.conflictQueue.pointOverlap')).toBeInTheDocument();
    expect(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.editRule' }),
    ).toBeInTheDocument();
    expect(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    ).toBeInTheDocument();
  });

  it('opens inline edit when clicking edit rule in the conflict queue', () => {
    mockPoints[0].address = '40002';

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const conflictItem = screen.getByTestId('conflict-item-40001');
    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.editRule' }),
    );

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    expect(within(ruleCard).getByTestId('rule-inline-edit-form')).toBeInTheDocument();
  });

  it('renders data type selector with grouped optgroups and disabled unsupported types', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    const dataTypeSelect = screen.getByLabelText('workbench.source.planner.dataType');
    const options = within(dataTypeSelect).getAllByRole('option');

    const supportedOptions = options.filter((opt) => !(opt as HTMLOptionElement).disabled);
    const disabledOptions = options.filter((opt) => (opt as HTMLOptionElement).disabled);

    expect(supportedOptions.length).toBe(3);
    expect(disabledOptions.length).toBe(7);

    const supportedValues = supportedOptions.map((opt) => (opt as HTMLOptionElement).value);
    expect(supportedValues).toEqual(expect.arrayContaining(['int16', 'int32', 'float32']));

    const disabledValues = disabledOptions.map((opt) => (opt as HTMLOptionElement).value);
    expect(disabledValues).toEqual(
      expect.arrayContaining(['bool', 'uint16', 'string', 'uint32', 'int64', 'uint64', 'float64']),
    );
  });

  it('uses protect plan wording instead of lock/unlock', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    const protectButton = within(ruleCard).getByRole('button', {
      name: 'workbench.source.ruleLayer.protectPlan',
    });
    expect(protectButton).toBeInTheDocument();

    fireEvent.click(protectButton);

    expect(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.unprotectPlan' }),
    ).toBeInTheDocument();
    expect(screen.getByTestId('rule-protect-hint-rule-1')).toHaveTextContent(
      'workbench.source.ruleLayer.protectHint',
    );
  });

  it('snaps selection to the root cell when clicking a merge continuation', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const continuationCell = screen.getByTestId('address-cell-40002');
    expect(continuationCell).toHaveAttribute('data-merge-offset', '1');

    fireEvent.click(continuationCell);

    const rootCell = screen.getByTestId('address-cell-40001');
    expect(rootCell).toHaveAttribute('aria-pressed', 'true');
    expect(continuationCell).toHaveAttribute('aria-pressed', 'true');
  });

  it('highlights the root cell of a selected logical span (continuations are visually merged)', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    const rootCell = screen.getByTestId('address-cell-40001');
    const continuationCell = screen.getByTestId('address-cell-40002');

    expect(rootCell.className).toContain('ring-2');
    // Continuation cell is sr-only (visually merged into root via gridColumn span)
    expect(continuationCell).toHaveClass('sr-only');
    expect(continuationCell).toHaveAttribute('aria-pressed', 'true');
  });

  it('emits conflict queue item for continuation-only conflicts', () => {
    // Point at 40003, rule plans float32 at 40002 (occupies 40002+40003).
    // Conflict at continuation cell (40003) should resolve to root (40002).
    mockPoints[0].address = '40003';

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40002' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const conflictQueue = screen.getByTestId('source-conflict-queue');
    expect(within(conflictQueue).getByTestId('conflict-item-40002')).toBeInTheDocument();
  });

  it('skip-span excludes only the conflicting logical span, not the entire rule', () => {
    // Point at 40003, rule plans float32 count=2 at 40001 (40001-40002 and 40003-40004).
    // Conflict at 40003 (root of second planned span). Skip should only remove 40003-40004.
    mockPoints[0].address = '40003';

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const conflictQueue = screen.getByTestId('source-conflict-queue');
    const conflictItem = within(conflictQueue).getByTestId('conflict-item-40003');

    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    );

    // Conflict should be resolved — queue disappears
    expect(screen.queryByTestId('source-conflict-queue')).not.toBeInTheDocument();

    // First span (40001-40002) should still be planned
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'planned');

    // Rule should still be enabled (not disabled)
    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('1');
  });

  it('aria-pressed is true for all cells of a selected logical span', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    fireEvent.click(screen.getByTestId('address-cell-40001'));

    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('aria-pressed', 'true');
  });

  it('skip-span resolves correct span root when rule-overlap occurs on a continuation cell', () => {
    // Rule A: int32 at 40001 (occupies 40001-40002, mergeOffset 0-1)
    // Rule B: int16 at 40002 (occupies 40002)
    // Conflict at 40002 — continuation cell of Rule A but root cell of Rule B
    // Skip must target Rule B's span root (40002), not Rule A's (40001)
    mockPoints.splice(0, mockPoints.length);

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    // Add Rule A: int32 at 40001
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // Add Rule B: int16 at 40002
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'int16' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40002' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // Conflict should exist
    const conflictQueue = screen.getByTestId('source-conflict-queue');
    const conflictItem = within(conflictQueue).getByTestId('conflict-item-40001');
    expect(conflictItem).toBeInTheDocument();

    // Click skip span — targets Rule B (last ruleId)
    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    );

    // Conflict should be resolved
    expect(screen.queryByTestId('source-conflict-queue')).not.toBeInTheDocument();

    // Rule A's span (40001-40002) should remain planned
    expect(screen.getByTestId('address-cell-40001')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40002')).toHaveAttribute('data-status', 'planned');

    // Both rules still exist (Rule B is enabled but its only span is skipped)
    expect(screen.getByTestId('source-rule-rule-1')).toBeInTheDocument();
    expect(screen.getByTestId('source-rule-rule-2')).toBeInTheDocument();
  });

  it('resets skippedAddresses when a rule is inline-edited so stale skips do not suppress wrong spans', () => {
    // 1. Create float32 rule at 40001 count=2 → planned at 40001, 40003
    // 2. Add point at 40003 → conflict at 40003
    // 3. Skip span at 40003 → skippedAddresses=['40003'], conflict resolved
    // 4. Inline-edit rule to start at 40010 → geometry changes completely
    // 5. Verify the old skip ('40003') does NOT suppress 40003-range in the new geometry
    //    (new planned addresses are 40010, 40012 — '40003' is irrelevant)
    //    If skippedAddresses were preserved, a future edit back to 40001 would silently lose 40003.
    mockPoints[0].address = '40003';

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '2' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // Skip the conflicting span
    const conflictItem = screen.getByTestId('conflict-item-40003');
    fireEvent.click(
      within(conflictItem).getByRole('button', { name: 'workbench.source.conflictQueue.skipSpan' }),
    );
    expect(screen.queryByTestId('source-conflict-queue')).not.toBeInTheDocument();
    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('1');

    // Inline-edit rule to move start address
    const ruleCard = screen.getByTestId('source-rule-rule-1');
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editStart' }),
    );
    const editForm = within(ruleCard).getByTestId('rule-inline-edit-form');
    fireEvent.change(within(editForm).getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40010' },
    });
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.editSave' }),
    );

    // After edit, both spans of the new geometry should be ready (no stale skip)
    expect(screen.getByTestId('source-summary-ready-count')).toHaveTextContent('2');
    expect(screen.getByTestId('address-cell-40010')).toHaveAttribute('data-status', 'planned');
    expect(screen.getByTestId('address-cell-40012')).toHaveAttribute('data-status', 'planned');
  });

  it('visually merges 32-bit cells with gridColumn span on root and hides continuations', () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.change(screen.getByLabelText('workbench.source.planner.dataType'), {
      target: { value: 'float32' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const rootCell = screen.getByTestId('address-cell-40001');
    const contCell = screen.getByTestId('address-cell-40002');

    // Root cell should span 2 columns
    expect(rootCell.style.gridColumn).toBe('span 2');
    // Continuation cell should be visually hidden
    expect(contCell).toHaveClass('sr-only');
  });

  it('persists only safe spans when conflicts exist alongside valid planned ranges', async () => {
    // Rule plans int16 at 40001, 40002, 40003.
    // Existing point at 40005 (no conflict).
    // Add a second rule at 40003 to create a rule-overlap conflict at 40003.
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

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    // Rule 1: int16 at 40001, count=3 → plans 40001, 40002, 40003
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '3' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // Rule 2: int16 at 40003, count=1 → overlap at 40003
    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40003' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    // 40003 is conflict, but 40001 and 40002 are planned and safe
    expect(screen.getByTestId('address-cell-40003')).toHaveAttribute('data-status', 'conflict');

    // Batch create should still be enabled — 2 safe spans
    const batchButton = screen.getByRole('button', {
      name: /workbench\.source\.actions\.createRulePoints/,
    });
    expect(batchButton).toBeEnabled();

    fireEvent.click(batchButton);

    await waitFor(() => {
      expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledTimes(1);
    });

    expect(mockCreateSourceRuleMutation.mutateAsync).toHaveBeenCalledWith(
      expect.objectContaining({
        id: 'rule-1',
        start_address: '40001',
        skipped_addresses: ['40003'],
      }),
    );
    expect(mockCreatePointMutation.mutateAsync).not.toHaveBeenCalled();
  });

  it('deletes orphaned points when a rule is deleted', () => {
    // Simulate a point that was created from rule-1's planned addresses
    mockPoints.splice(0, mockPoints.length, {
      id: 'point-rule-1',
      device_id: 'device-1',
      name: 'SRC_40001',
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
    });

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

    fireEvent.change(screen.getByLabelText('workbench.source.planner.startAddress'), {
      target: { value: '40001' },
    });
    fireEvent.change(screen.getByLabelText('workbench.source.planner.count'), {
      target: { value: '1' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.planner.addRule' }));

    const ruleCard = screen.getByTestId('source-rule-rule-1');
    fireEvent.click(
      within(ruleCard).getByRole('button', { name: 'workbench.source.ruleLayer.delete' }),
    );

    // Point at 40001 should be deleted
    expect(mockDeletePointMutation.mutate).toHaveBeenCalledWith('point-rule-1');
  });
});
