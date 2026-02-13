import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import type { ReactNode } from 'react';
import { beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';
import SmartDashboard from '../SmartDashboard';

const { mockDevicesState, mockDatalinkState, mockMutations, mockModbusShareAPI } = vi.hoisted(() => ({
  mockDevicesState: {
    devices: [] as Array<{
      id: string;
      name: string;
      description: string;
      protocol: 'modbus_tcp';
      status: 'active' | 'disabled' | 'draft';
      connection_config: string;
      last_test_at: string | null;
      last_test_success: boolean | null;
      last_test_error: string;
      created_at: string;
      updated_at: string;
    }>,
  },
  mockDatalinkState: {
    pollingGroups: [] as unknown[],
    points: [] as unknown[],
    mappings: [] as unknown[],
    tags: [] as unknown[],
  },
  mockMutations: {
    createPoint: { mutateAsync: vi.fn(), isPending: false },
    validatePipeline: { mutateAsync: vi.fn(), isPending: false },
    createMapping: { mutateAsync: vi.fn(), isPending: false },
    updateMapping: { mutateAsync: vi.fn(), isPending: false },
    createTag: { mutateAsync: vi.fn(), isPending: false },
    updateTag: { mutateAsync: vi.fn(), isPending: false },
  },
  mockModbusShareAPI: {
    status: vi.fn(),
    upsertMapping: vi.fn(),
    writeTagValue: vi.fn(),
    sync: vi.fn(),
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, params?: Record<string, string>) => {
      if (key === 'smartDashboard.legacyMigration.description') {
        return `legacy moved ${params?.route ?? ''}`.trim();
      }
      return key;
    },
  }),
}));

vi.mock('../../../services/datalink', () => ({
  modbusShareAPI: mockModbusShareAPI,
}));

vi.mock('../../../hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevicesState.devices }),
}));

vi.mock('../../../hooks/datalink/usePollingGroups', () => ({
  usePollingGroupsQuery: () => ({ data: mockDatalinkState.pollingGroups }),
}));

vi.mock('../../../hooks/datalink/usePoints', () => ({
  usePointsQuery: () => ({ data: mockDatalinkState.points }),
  useCreatePointMutation: () => mockMutations.createPoint,
}));

vi.mock('../../../hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockDatalinkState.mappings }),
  useValidatePipelineMutation: () => mockMutations.validatePipeline,
  useCreateMappingMutation: () => mockMutations.createMapping,
  useUpdateMappingMutation: () => mockMutations.updateMapping,
}));

vi.mock('../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mockDatalinkState.tags }),
  useCreateTagMutation: () => mockMutations.createTag,
  useUpdateTagMutation: () => mockMutations.updateTag,
}));

vi.mock('../../../hooks/useKeyboardShortcuts', () => ({
  useSmartDashboardShortcuts: () => [],
}));

vi.mock('../../../hooks/useHistory', () => ({
  usePointHistory: () => ({
    canUndo: false,
    canRedo: false,
    undo: vi.fn(),
    redo: vi.fn(),
    push: vi.fn(),
    getUndoAction: vi.fn(() => null),
    getRedoAction: vi.fn(() => null),
  }),
}));

vi.mock('../../../features/flow/stateMachine', () => ({
  useFlowLifecycle: () => ({
    state: {
      status: 'draft',
      sinkTarget: '',
      diagnostics: {
        source: { latestValue: '-', quality: 'unknown', timestamp: '-', error: '' },
        grid: { latestValue: '-', quality: 'unknown', timestamp: '-', error: '' },
        tag: { latestValue: '-', quality: 'unknown', timestamp: '-', error: '' },
        sink: { latestValue: '-', quality: 'unknown', timestamp: '-', error: '' },
      },
    },
    canValidate: true,
    canActivate: true,
    hasError: false,
    setSource: vi.fn(),
    setTag: vi.fn(),
    setSink: vi.fn(),
    setDiagnostics: vi.fn(),
    markValidated: vi.fn(),
    markActive: vi.fn(),
    markError: vi.fn(),
    resetDraft: vi.fn(),
  }),
}));

vi.mock('../../../components/datalink/DeviceTreeNav', () => ({
  DeviceTreeNav: ({ devices, onSelectDevice }: { devices: Array<{ id: string; name: string }>; onSelectDevice: (id: string) => void }) => (
    <div data-testid="device-tree-nav-mock">
      {devices.map((device) => (
        <button key={device.id} type="button" onClick={() => onSelectDevice(device.id)}>
          切換-{device.name}
        </button>
      ))}
    </div>
  ),
}));

vi.mock('../../../components/datalink/MemoryGrid', () => ({
  MemoryGrid: ({ onSelect }: { onSelect: (addresses: string[]) => void }) => (
    <div data-testid="memory-grid-mock">
      <button type="button" onClick={() => onSelect(['40001'])}>
        mock-select-address
      </button>
    </div>
  ),
}));

vi.mock('../../../components/datalink/QuickActions', () => ({
  QuickActions: () => <div data-testid="quick-actions-mock" />,
}));

vi.mock('../../../components/datalink/SlidePanel', () => ({
  SlidePanel: ({ isOpen, children }: { isOpen: boolean; children: ReactNode }) => (isOpen ? <div>{children}</div> : null),
}));

vi.mock('../../../components/datalink/BatchPointCreator', () => ({
  BatchPointCreator: () => null,
}));

vi.mock('../../../components/datalink/PointDetailPanel', () => ({
  PointDetailPanel: () => null,
}));

vi.mock('../../../components/datalink/ImportExportDialog', () => ({
  ImportDialog: () => null,
  ExportDialog: () => null,
}));

vi.mock('../../../components/datalink/wizard/DeviceOnboardingWizard', () => ({
  default: () => null,
}));

function renderDashboard(entry = '/datalink') {
  return render(
    <MemoryRouter initialEntries={[entry]}>
      <Routes>
        <Route path="/datalink" element={<SmartDashboard />} />
      </Routes>
    </MemoryRouter>
  );
}

describe('SmartDashboard interactions', () => {
  beforeAll(() => {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: vi.fn().mockImplementation(() => ({
        matches: false,
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
      })),
    });
  });

  beforeEach(() => {
    vi.clearAllMocks();
    mockDatalinkState.pollingGroups = [];
    mockDatalinkState.points = [];
    mockDatalinkState.mappings = [];
    mockDatalinkState.tags = [];
    mockDevicesState.devices = [
      {
        id: 'device-active',
        name: 'Device Active',
        description: '',
        protocol: 'modbus_tcp',
        status: 'active',
        connection_config: '',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'device-disabled',
        name: 'Device Disabled',
        description: '',
        protocol: 'modbus_tcp',
        status: 'disabled',
        connection_config: '',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
      {
        id: 'device-draft',
        name: 'Device Draft',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '',
        updated_at: '',
      },
    ];

    mockModbusShareAPI.status.mockResolvedValue({
      enabled: false,
      port: 5020,
      address: '',
      bind_state: 'fail',
      mapping_count: 0,
    });
    mockModbusShareAPI.upsertMapping.mockResolvedValue({});
    mockModbusShareAPI.writeTagValue.mockResolvedValue({});
    mockModbusShareAPI.sync.mockResolvedValue({ updated: 0, skipped: 0, errors: [] });
  });

  it('shows section intent notice and can dismiss it', async () => {
    renderDashboard('/datalink?section=devices');

    const sectionNotice = await screen.findByText(/已導向/);
    expect(sectionNotice).toBeInTheDocument();
    fireEvent.click(within(sectionNotice.closest('section') as HTMLElement).getByRole('button', { name: '關閉' }));
    await waitFor(() => {
      expect(screen.queryByText(/已導向/)).not.toBeInTheDocument();
    });
  });

  it('shows modal intent notice and can dismiss it', async () => {
    renderDashboard('/datalink?modal=devices');

    const modalNotice = await screen.findByText(/Dashboard modal 流程/);
    expect(modalNotice).toBeInTheDocument();
    fireEvent.click(within(modalNotice.closest('section') as HTMLElement).getByRole('button', { name: '關閉' }));
    await waitFor(() => {
      expect(screen.queryByText(/Dashboard modal 流程/)).not.toBeInTheDocument();
    });
  });

  it('switches device successfully and updates context bar', async () => {
    renderDashboard();

    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Active' }));

    await waitFor(() => {
      expect(screen.getAllByText('Device Active').length).toBeGreaterThan(0);
    });
    expect(screen.getByText(/最近切換:/)).toBeInTheDocument();
    expect(screen.queryByText('最近切換: -')).not.toBeInTheDocument();
  });

  it('guards unsaved changes and supports cancel then discard-switch to read-only', async () => {
    renderDashboard();

    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Active' }));
    fireEvent.click(await screen.findByRole('button', { name: 'mock-select-address' }));

    fireEvent.click(screen.getByRole('button', { name: '切換設備' }));
    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Disabled' }));

    expect(await screen.findByText('有未儲存變更')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: '取消' }));
    await waitFor(() => {
      expect(screen.queryByText('有未儲存變更')).not.toBeInTheDocument();
    });
    expect(screen.getAllByText('Device Active').length).toBeGreaterThan(0);

    fireEvent.click(screen.getByRole('button', { name: '切換設備' }));
    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Disabled' }));
    fireEvent.click(await screen.findByRole('button', { name: '放棄並切換' }));

    await waitFor(() => {
      expect(screen.getAllByText('Device Disabled').length).toBeGreaterThan(0);
    });
    expect(screen.getByText('Read-only')).toBeInTheDocument();
  });

  it('shows switch failure path with details and retry controls for offline/draft device', async () => {
    renderDashboard();

    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Active' }));
    fireEvent.click(screen.getByRole('button', { name: '切換設備' }));
    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Draft' }));
    fireEvent.click(await screen.findByRole('button', { name: '放棄並切換' }));

    expect(await screen.findByText(/目前不可切換/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '重試' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '查看詳情' }));
    expect(await screen.findByText(/若設備狀態為離線/)).toBeInTheDocument();
  });

  it('covers e2e-like mainline: success, unsaved intercept, read-only switch, then failure', async () => {
    renderDashboard();

    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Active' }));
    await waitFor(() => {
      expect(screen.getAllByText('Device Active').length).toBeGreaterThan(0);
    });

    fireEvent.click(await screen.findByRole('button', { name: 'mock-select-address' }));
    fireEvent.click(screen.getByRole('button', { name: '切換設備' }));
    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Disabled' }));
    expect(await screen.findByText('有未儲存變更')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: '放棄並切換' }));
    expect(await screen.findByText('Read-only')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '切換設備' }));
    fireEvent.click(await screen.findByRole('button', { name: '切換-Device Draft' }));
    fireEvent.click(await screen.findByRole('button', { name: '放棄並切換' }));
    expect(await screen.findByText(/目前不可切換/)).toBeInTheDocument();
  });
});
