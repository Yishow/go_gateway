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
    it('renders all five shell regions', () => {
      renderPage();

      expect(screen.getByTestId('workbench-frame')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-context-bar')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-step-rail')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-inspector-panel')).toBeInTheDocument();
      expect(screen.getByTestId('workbench-bottom-summary-bar')).toBeInTheDocument();
    });

    it('locks desktop layout sizing and keeps overflow inside the work area', () => {
      renderPage();

      const frame = screen.getByTestId('workbench-frame');
      const workArea = screen.getByTestId('workbench-primary-work-area');

      expect(frame).toHaveStyle({
        gridTemplateRows: 'auto 1fr auto',
        gridTemplateColumns: '200px 1fr 280px',
      });
      expect(frame).toHaveClass('overflow-hidden');
      expect(workArea).toHaveClass('overflow-auto');
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
      expect(within(rail).getByRole('button', { name: /workbench\.steps\.device/ })).toBeInTheDocument();
      expect(within(rail).getByRole('button', { name: /workbench\.steps\.source/ })).toBeInTheDocument();
      expect(within(rail).getByRole('button', { name: /workbench\.steps\.tag/ })).toBeInTheDocument();
      expect(within(rail).getByRole('button', { name: /workbench\.steps\.output/ })).toBeInTheDocument();
    });

    it('highlights the active step with aria-current', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      const deviceBtn = within(rail).getByRole('button', { name: /workbench\.steps\.device/ });
      expect(deviceBtn).toHaveAttribute('aria-current', 'step');

      fireEvent.click(within(rail).getByRole('button', { name: /workbench\.steps\.source/ }));
      expect(deviceBtn).not.toHaveAttribute('aria-current');
      expect(within(rail).getByRole('button', { name: /workbench\.steps\.source/ })).toHaveAttribute('aria-current', 'step');
    });
  });

  describe('WorkbenchContextBar', () => {
    it('shows device name when a device is selected', () => {
      renderPage();

      // Select the device first by clicking it in the device step
      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('context-bar-device-name')).toHaveTextContent('Mixer PLC');
    });

    it('shows device capability chips when a device is selected', () => {
      renderPage();

      fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));

      expect(screen.getByTestId('context-bar-capability-unit-id')).toBeInTheDocument();
      expect(screen.getByTestId('context-bar-capability-address-base')).toBeInTheDocument();
      expect(screen.getByTestId('context-bar-capability-word-order')).toBeInTheDocument();
      expect(screen.getByTestId('context-bar-capability-protocol-traits')).toBeInTheDocument();
    });

    it('shows no-device label when no device is selected', () => {
      mockDevices.splice(0, mockDevices.length);
      renderPage();

      expect(screen.getByTestId('context-bar-device-name')).toHaveTextContent('workbench.contextBar.noDevice');
    });

    it('provides quick-action buttons for navigation', () => {
      renderPage();

      const contextBar = screen.getByTestId('workbench-context-bar');
      expect(within(contextBar).getByText('workbench.contextBar.actions.selectDevice')).toBeInTheDocument();
      expect(within(contextBar).getByText('workbench.contextBar.actions.gotoSource')).toBeInTheDocument();
      expect(within(contextBar).getByText('workbench.contextBar.actions.gotoOutput')).toBeInTheDocument();
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

      expect(
        screen.getByText('workbench.device.inspector.emptyTitle'),
      ).toBeInTheDocument();
      expect(
        screen.getByText('workbench.device.inspector.emptyDescription'),
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
  });

  describe('Step switching wires content into PrimaryWorkArea', () => {
    it('switches step content when step rail buttons are clicked', () => {
      renderPage();

      const rail = screen.getByTestId('workbench-step-rail');
      const workArea = screen.getByTestId('workbench-primary-work-area');

      // Device step is the default — device step content should be in work area
      expect(workArea).toBeInTheDocument();

      // Switch to source step
      fireEvent.click(within(rail).getByRole('button', { name: /workbench\.steps\.source/ }));
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();

      // Switch to tag step
      fireEvent.click(within(rail).getByRole('button', { name: /workbench\.steps\.tag/ }));
      expect(screen.getByTestId('workbench-primary-work-area')).toBeInTheDocument();
    });
  });
});
