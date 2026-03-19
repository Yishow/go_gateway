import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchPage from '../DatalinkWorkbenchPage';
import type { Device, Mapping, Point, Tag } from '../../../../types/datalink';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockRuntimeStatus,
  mockEventSources,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockRuntimeStatus: vi.fn(),
  mockEventSources: [] as MockEventSource[],
}));

type EventHandler = (event: MessageEvent<string>) => void;

class MockEventSource {
  url: string;
  readyState = 0;
  onopen: ((event: Event) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  private listeners = new Map<string, Set<EventHandler>>();

  constructor(url: string) {
    this.url = url;
    mockEventSources.push(this);
  }

  addEventListener(type: string, listener: EventHandler) {
    const existing = this.listeners.get(type) ?? new Set<EventHandler>();
    existing.add(listener);
    this.listeners.set(type, existing);
  }

  removeEventListener(type: string, listener: EventHandler) {
    this.listeners.get(type)?.delete(listener);
  }

  close() {
    this.readyState = 2;
  }

  emit(type: string, data: unknown) {
    const payload = { data: JSON.stringify(data) } as MessageEvent<string>;
    this.listeners.get(type)?.forEach((listener) => listener(payload));
  }
}

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

vi.mock('../../../../hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: () => ({
    data: [],
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

vi.mock('../../../../hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({
    data: mockTags,
    isLoading: false,
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
  runtimeAPI: {
    getStatus: mockRuntimeStatus,
    getStreamUrl: vi.fn((deviceId: string, pointIds?: string[]) => {
      const params = new URLSearchParams({ device_id: deviceId });
      if (pointIds && pointIds.length > 0) {
        params.set('point_ids', pointIds.join(','));
      }
      return `/api/v1/datalink/runtime/stream?${params.toString()}`;
    }),
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

describe('DatalinkWorkbench runtime phase surface', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockEventSources.splice(0, mockEventSources.length);
    vi.stubGlobal('EventSource', MockEventSource);

    mockDevices.splice(0, mockDevices.length, {
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
    });

    mockPoints.splice(0, mockPoints.length, {
      id: 'point-1',
      device_id: 'device-1',
      name: 'Flow Sensor',
      description: '',
      data_type: 'int16',
      address: '40001',
      enabled: true,
      polling_group_id: '',
      last_value: 37.5,
      last_read_at: '2026-03-16T06:00:00Z',
      last_error: '',
      error_count: 0,
      created_at: '',
      updated_at: '',
    });

    mockTags.splice(0, mockTags.length);
    mockMappings.splice(0, mockMappings.length);

    mockRuntimeStatus.mockResolvedValue({
      running: true,
      uptime_seconds: 120,
      collectors: [
        {
          device_id: 'device-1',
          device_name: 'Mixer PLC',
          protocol: 'modbus_tcp',
          status: 'running',
          points_total: 1,
          points_healthy: 1,
          points_stale: 0,
          points_error: 0,
          last_read_at: '2026-03-16T06:00:00Z',
          last_error: null,
          breaker_state: 'closed',
        },
      ],
    });
  });

  it('shows runtime health summary and current point value inside the source view', async () => {
    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.live' }));

    await waitFor(() => {
      expect(screen.getByTestId('workbench-runtime-status')).toHaveTextContent(
        'workbench.runtime.summary.running',
      );
    });

    expect(screen.getByTestId('address-cell-40001')).toHaveTextContent('37.5');
    expect(screen.getByText('workbench.runtime.summary.title')).toBeInTheDocument();
  });

  it('updates the source cell when runtime stream emits a new live value', async () => {
    mockPoints[0].last_value = null;

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.live' }));

    await waitFor(() => {
      expect(mockEventSources).toHaveLength(1);
    });

    act(() => {
      mockEventSources[0].emit('value', {
        device_id: 'device-1',
        point_id: 'point-1',
        address: '40001',
        raw_value: 42.25,
        transformed_value: 42.25,
        quality: 'good',
        stale: false,
        timestamp: '2026-03-16T06:01:00Z',
      });
    });

    await waitFor(() => {
      expect(screen.getByTestId('address-cell-40001')).toHaveTextContent('42.25');
    });
  });

  it('keeps source canvas on the raw point value when transformed value differs', async () => {
    mockPoints[0].last_value = null;

    renderPage();

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));
    fireEvent.click(screen.getByRole('button', { name: 'Mixer PLC' }));
    fireEvent.click(screen.getByRole('button', { name: 'workbench.source.view.live' }));

    await waitFor(() => {
      expect(mockEventSources).toHaveLength(1);
    });

    act(() => {
      mockEventSources[0].emit('value', {
        device_id: 'device-1',
        point_id: 'point-1',
        address: '40001',
        raw_value: 99,
        transformed_value: 12,
        quality: 'good',
        stale: false,
        timestamp: '2026-03-16T06:02:00Z',
      });
    });

    await waitFor(() => {
      expect(screen.getByTestId('address-cell-40001')).toHaveTextContent('99');
    });
    expect(screen.getByTestId('address-cell-40001')).not.toHaveTextContent('12');
  });
});
