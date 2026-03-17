import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import { SOURCE_TEMPLATE_STORAGE_KEY } from '../../../../features/datalink/sourceTemplateStorage';
import type { Device, Mapping, Point, Tag } from '../../../../types/datalink';

const { mockDevices, mockPoints, mockMappings, mockTags, mockCreatePointMutation } = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockMappings: [] as Mapping[],
  mockTags: [] as Tag[],
  mockCreatePointMutation: {
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
    mockMappings.splice(0, mockMappings.length);
    mockTags.splice(0, mockTags.length);
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
    expect(screen.getByTestId('address-cell-40005')).toHaveAttribute('data-status', 'used');
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

  it('keeps only view controls in the primary toolbar and reveals utility tools on demand', () => {
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
    expect(
      within(primaryToolbar).queryByRole('button', {
        name: 'workbench.source.toolbar.saveTemplate',
      }),
    ).not.toBeInTheDocument();
    expect(
      within(primaryToolbar).queryByRole('button', {
        name: 'workbench.source.toolbar.snapshotCompare',
      }),
    ).not.toBeInTheDocument();

    expect(
      within(secondaryControls).getByRole('button', {
        name: 'workbench.source.toolbar.moreTools',
      }),
    ).toBeInTheDocument();
    expect(
      within(secondaryControls).queryByRole('button', {
        name: 'workbench.source.toolbar.saveTemplate',
      }),
    ).not.toBeInTheDocument();
    expect(
      within(secondaryControls).queryByRole('button', {
        name: 'workbench.source.toolbar.loadTemplate',
      }),
    ).not.toBeInTheDocument();
    expect(
      within(secondaryControls).queryByRole('button', {
        name: 'workbench.source.toolbar.snapshotCompare',
      }),
    ).not.toBeInTheDocument();

    fireEvent.click(
      within(secondaryControls).getByRole('button', {
        name: 'workbench.source.toolbar.moreTools',
      }),
    );

    expect(
      within(secondaryControls).getByRole('button', {
        name: 'workbench.source.toolbar.saveTemplate',
      }),
    ).toBeInTheDocument();
    expect(
      within(secondaryControls).getByRole('button', {
        name: 'workbench.source.toolbar.loadTemplate',
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

    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.toolbar.moreTools' }));
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

  it('batch creates points from the planned address range', async () => {
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
      expect(mockCreatePointMutation.mutateAsync).toHaveBeenCalledTimes(2);
    });

    expect(mockCreatePointMutation.mutateAsync).toHaveBeenNthCalledWith(1, {
      device_id: 'device-1',
      address: '40001',
      data_type: 'float32',
      name: 'SRC_40001',
    });
    expect(mockCreatePointMutation.mutateAsync).toHaveBeenNthCalledWith(2, {
      device_id: 'device-1',
      address: '40003',
      data_type: 'float32',
      name: 'SRC_40003',
    });
  });

  it('blocks batch create when planned cells overlap existing points', () => {
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
});
