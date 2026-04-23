import { fireEvent, render, screen, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import type { Device, Mapping, Point, Tag } from '../../../../types/datalink';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
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
  useTestDraftConnectionMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useToggleDeviceStatusMutation: () => ({
    mutateAsync: vi.fn().mockResolvedValue(undefined),
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
  useCreatePointMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
  useDeletePointMutation: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('../../../../hooks/datalink/useTags', () => ({
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
    isPending: false,
  }),
}));

vi.mock('../../../../services/datalink', () => ({
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

describe('DatalinkWorkbench five-region shell', () => {
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

  describe('WorkbenchFrame layout', () => {
    it('renders shell regions (step rail embedded in context bar)', () => {
      renderPage();

      expect(screen.getByTestId('workbench-frame')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-context-bar')).toBeInTheDocument();
      expect(
        within(screen.getByTestId('workbench-context-bar')).getByTestId('workbench-step-rail'),
      ).toBeInTheDocument();
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-inspector-panel')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-bottom-summary-bar')).toBeInTheDocument();
    });

    it('locks desktop layout sizing (v2-mui uses sx-based grid)', () => {
      renderPage();

      const frame = screen.getByTestId('workbench-frame');
      const workArea = screen.getByTestId('workbench-primary-work-area');

      expect(frame).toHaveAttribute('data-variant', 'v2-mui');
      expect(frame).toBeInTheDocument();
      expect(workArea).toBeInTheDocument();
    });

    it('does not render the old ActionDock or HeaderBar', () => {
      renderPage();

      expect(screen.queryByText('workbench.actionDock.title')).not.toBeInTheDocument();
      expect(screen.queryByText('workbench.header.eyebrow')).not.toBeInTheDocument();
    });
  });

  describe('WorkbenchStepRail', () => {
    it('renders all four steps in the rail', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.device/ })).toBeInTheDocument();
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ })).toBeInTheDocument();
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.tag/ })).toBeInTheDocument();
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.output/ })).toBeInTheDocument();
    });

    it('highlights the active step with aria-selected (MUI Tabs)', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      const deviceTab = within(rail).getByRole('tab', { name: /workbench\.steps\.device/ });
      expect(deviceTab).toHaveAttribute('aria-selected', 'true');

      fireEvent.click(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ }));
      expect(deviceTab).toHaveAttribute('aria-selected', 'false');
      expect(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ })).toHaveAttribute('aria-selected', 'true');
    });
  });

  describe('WorkbenchContextBar', () => {
    it('shows device name when a device is selected', () => {
      renderPage();

      // Select the device first by clicking it in the device step
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('context-bar-device-name')).toHaveTextContent('Mixer PLC');
    });

    it('renders a compact step summary instead of capability chips when a device is selected', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const summary = screen.getByTestId('context-bar-step-summary');
      expect(summary).toHaveAttribute('data-active-step', 'device');
      expect(summary).toHaveTextContent('Mixer PLC');
      expect(screen.queryByTestId('context-bar-capability-unit-id')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-capability-address-base')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-capability-word-order')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-capability-protocol-traits')).not.toBeInTheDocument();
      expect(screen.queryByTestId('context-bar-test-status')).not.toBeInTheDocument();
    });

    it('shows no-device label when no device is selected', () => {
      mockDevices.splice(0, mockDevices.length);
      renderPage();

      expect(screen.getByTestId('context-bar-device-name')).toHaveTextContent('workbench.contextBar.noDevice');
    });

    it('shows a single primary action instead of a quick-action cluster', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      const contextBar = screen.getByTestId('workbench-context-bar');
      expect(
        within(contextBar).getByRole('button', {
          name: 'workbench.contextBar.actions.gotoSource',
        }),
      ).toBeInTheDocument();
    });

    it('keeps the primary action aligned with the next step instead of skipping from source to output', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      const contextBar = screen.getByTestId('workbench-context-bar');
      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'source');
      expect(
        within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }),
      ).toBeInTheDocument();
    });

    it('disables the source-step primary action until source points exist', () => {
      mockPoints.splice(0, mockPoints.length);
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'source');
      expect(
        within(screen.getByTestId('workbench-context-bar')).getByRole('button', {
          name: 'workbench.contextBar.actions.gotoTag',
        }),
      ).toBeDisabled();
    });

    it('renders the source command banner when the source step is active', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      fireEvent.click(screen.getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));

      expect(screen.getByTestId('source-command-banner')).toBeInTheDocument();
      expect(screen.getByTestId('source-command-banner')).toHaveTextContent(
        'workbench.source.sentryBanner.title.build',
      );
    });

    it('returns to tag review on step 4 when output scope has no active rule handoff', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      const contextBar = screen.getByTestId('workbench-context-bar');
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoOutput' }));

      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'output');
      const outputCta = within(contextBar).getByTestId('context-bar-primary-action');
      expect(outputCta).toBeInTheDocument();
      expect(outputCta).toHaveTextContent('workbench.shell.actions.returnTag');
      expect(outputCta).not.toHaveTextContent('workbench.contextBar.actions.switchDevice');
    });

    it('routes back to tag review when output scope must recover the active rule handoff', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
      const contextBar = screen.getByTestId('workbench-context-bar');
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoSource' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoTag' }));
      fireEvent.click(within(contextBar).getByRole('button', { name: 'workbench.contextBar.actions.gotoOutput' }));

      const action = within(contextBar).getByTestId('context-bar-primary-action');
      expect(action).toHaveTextContent('workbench.shell.actions.returnTag');

      fireEvent.click(action);

      expect(screen.getByTestId('context-bar-step-summary')).toHaveAttribute('data-active-step', 'tag');
    });
  });

  describe('WorkbenchInspectorPanel', () => {
    it('renders the inspector with step-dependent heading', () => {
      renderPage();

      const inspector = screen.getByTestId('workbench-inspector-panel');
      expect(inspector).toBeInTheDocument();
      expect(
        within(inspector).getByText('workbench.device.inspector.emptyTitle'),
      ).toBeInTheDocument();
    });

    it('shows device-step empty state when no item is selected', () => {
      renderPage();

      const inspector = screen.getByTestId('workbench-inspector-panel');

      expect(
        within(inspector).getByText('workbench.device.inspector.emptyTitle'),
      ).toBeInTheDocument();
      expect(
        within(inspector).getByText('workbench.device.inspector.emptyDescription'),
      ).toBeInTheDocument();
      expect(screen.queryByTestId('inspector-selection-context')).not.toBeInTheDocument();
    });
  });

  describe('WorkbenchBottomSummaryBar', () => {
    it('shows summary metrics after device is selected', () => {
      renderPage();

      // Select device to populate counters
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('summary-point-count')).toHaveTextContent('2');
      expect(screen.getByTestId('summary-tag-count')).toHaveTextContent('1');
      expect(screen.getByTestId('summary-output-count')).toHaveTextContent('1');
    });

    it('shows readiness indicators with four steps', () => {
      renderPage();

      const summaryBar = screen.getByTestId('workbench-bottom-summary-bar');
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.device')).toBeInTheDocument();
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.source')).toBeInTheDocument();
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.tag')).toBeInTheDocument();
      expect(within(summaryBar).getByText('workbench.bottomSummary.readiness.output')).toBeInTheDocument();
    });

    it('shows the active output target badge', () => {
      renderPage();

      expect(screen.getByTestId('active-output-target')).toBeInTheDocument();
      expect(screen.getByText('workbench.bottomSummary.targets.modbus')).toBeInTheDocument();
    });

    it('propagates rich readiness status via data attributes', () => {
      renderPage();

      // No device selected yet → device readiness should be "draft"
      const deviceIndicator = screen.getByTestId('readiness-device');
      expect(deviceIndicator).toHaveAttribute('data-readiness', 'draft');
    });

    it('marks tag readiness as blocked when source points are missing', () => {
      mockPoints.splice(0, mockPoints.length);
      mockTags.splice(0, mockTags.length);
      mockMappings.splice(0, mockMappings.length);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-readiness', 'draft');
      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-reason', 'no-points');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-readiness', 'blocked');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-reason', 'no-source-points');
    });

    it('marks output readiness as blocked until tags are linked', () => {
      mockTags.splice(0, mockTags.length);
      mockMappings.splice(0, mockMappings.length);

      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-readiness', 'draft');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-reason', 'no-tags-linked');
      expect(screen.getByTestId('readiness-output')).toHaveAttribute('data-readiness', 'blocked');
      expect(screen.getByTestId('readiness-output')).toHaveAttribute('data-reason', 'no-tags-linked');
    });

    it('emphasizes only the active step and compacts the rest of the readiness strip', () => {
      renderPage();

      expect(screen.getByTestId('readiness-device')).toHaveAttribute('data-emphasis', 'active');
      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-emphasis', 'compact');
      expect(screen.getByTestId('readiness-tag')).toHaveAttribute('data-emphasis', 'compact');
      expect(screen.getByTestId('readiness-output')).toHaveAttribute('data-emphasis', 'compact');

      fireEvent.click(screen.getByRole('tab', { name: /workbench\.steps\.source/ }));

      expect(screen.getByTestId('readiness-device')).toHaveAttribute('data-emphasis', 'compact');
      expect(screen.getByTestId('readiness-source')).toHaveAttribute('data-emphasis', 'active');
    });
  });

  describe('Step switching wires content into PrimaryWorkArea', () => {
    it('switches step content when step rail tabs are clicked', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      const workArea = screen.getByTestId('workbench-primary-work-area');

      // Device step is the default — device step content should be in work area
      expect(workArea).toBeInTheDocument();

      // Switch to source step
      fireEvent.click(within(rail).getByRole('tab', { name: /workbench\.steps\.source/ }));
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();

      // Switch to tag step
      fireEvent.click(within(rail).getByRole('tab', { name: /workbench\.steps\.tag/ }));
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();
    });
  });
});
