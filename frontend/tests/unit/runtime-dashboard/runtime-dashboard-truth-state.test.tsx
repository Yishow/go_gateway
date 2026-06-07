import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { RuntimeDashboardRoute } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute';
import type { Point, RuntimeStatus } from '../../../src/types/datalink';

const {
  mockRuntimeContext,
  mockPoints,
  mockRuntimeStatus,
  mockEventSources,
} = vi.hoisted(() => ({
  mockRuntimeContext: {
    workspace_id: 'workspace-1',
    devices: [] as Array<{
      device_id: string;
      name: string;
      protocol: string;
      running: boolean;
      availability_status: 'available' | 'unavailable';
      availability_reason: string | null;
    }>,
    default_device_id: null as string | null,
  },
  mockPoints: [] as Point[],
  mockRuntimeStatus: vi.fn<() => Promise<RuntimeStatus>>(),
  mockEventSources: [] as MockEventSource[],
}));

type EventHandler = (event: MessageEvent<string>) => void;

class MockEventSource {
  onopen: ((event: Event) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;
  private listeners = new Map<string, Set<EventHandler>>();

  constructor() {
    mockEventSources.push(this);
  }

  addEventListener(type: string, listener: EventHandler) {
    const listeners = this.listeners.get(type) ?? new Set<EventHandler>();
    listeners.add(listener);
    this.listeners.set(type, listeners);
  }

  removeEventListener(type: string, listener: EventHandler) {
    this.listeners.get(type)?.delete(listener);
  }

  close() {}

  emit(type: string, data: unknown) {
    const payload = { data: JSON.stringify(data) } as MessageEvent<string>;
    this.listeners.get(type)?.forEach((listener) => listener(payload));
  }

  open() {
    this.onopen?.(new Event('open'));
  }
}

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, fallback?: unknown) =>
      typeof fallback === 'string' ? fallback : key,
  }),
}));

vi.mock('../../../src/hooks/datalink/useStudioV2RuntimeContext', () => ({
  useStudioV2RuntimeContextQuery: () => ({
    data: mockRuntimeContext,
    isLoading: false,
    isError: false,
    error: null,
  }),
}));

vi.mock('../../../src/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id
      ? mockPoints.filter((point) => point.device_id === filters.device_id)
      : [],
    isLoading: false,
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  runtimeAPI: {
    getStatus: mockRuntimeStatus,
    getStreamUrl: vi.fn((deviceId: string) =>
      `/api/v1/datalink/runtime/stream?device_id=${deviceId}`,
    ),
  },
}));

function renderRoute(initialEntry: string) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialEntry]}>
        <Routes>
          <Route path="/studio/runtime" element={<RuntimeDashboardRoute />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('runtime dashboard truth states', () => {
  beforeEach(() => {
    mockRuntimeStatus.mockReset();
    mockEventSources.splice(0, mockEventSources.length);
    vi.stubGlobal('EventSource', MockEventSource);

    mockRuntimeContext.devices.splice(0, mockRuntimeContext.devices.length, {
      device_id: 'device-A',
      name: 'Mixer PLC',
      protocol: 'modbus_tcp',
      running: true,
      availability_status: 'available',
      availability_reason: null,
    });
    mockRuntimeContext.default_device_id = 'device-A';
    mockPoints.splice(0, mockPoints.length, {
      id: 'point-1',
      device_id: 'device-A',
      name: 'Flow',
      description: '',
      address: '40001',
      data_type: 'float32',
      enabled: true,
      polling_group_id: 'group-A',
      created_at: '2026-05-29T00:00:00Z',
      updated_at: '2026-05-29T00:00:00Z',
      last_read_at: '',
      last_value: null,
      last_error: '',
      error_count: 0,
    } as Point);
    mockRuntimeStatus.mockResolvedValue({
      running: true,
      uptime_seconds: 12,
      collectors: [{
        device_id: 'device-A',
        device_name: 'Mixer PLC',
        protocol: 'modbus_tcp',
        status: 'running',
        availability_status: 'available',
        availability_reason: null,
        running: true,
        points_total: 1,
        points_healthy: 1,
        points_stale: 0,
        points_error: 0,
        last_read_at: '2026-05-29T00:00:00Z',
        last_error: null,
        breaker_state: 'closed',
      }],
    });
  });

  it('renders explicit backend empty snapshot state instead of a synthetic ready dashboard', async () => {
    mockRuntimeStatus.mockResolvedValueOnce({
      running: false,
      uptime_seconds: 0,
      snapshot_state: {
        state: 'empty',
        empty: true,
        degraded: false,
        unavailable: false,
        stale: false,
        reason: 'selected device has no runtime snapshot',
      },
      collectors: [],
    } as RuntimeStatus);

    renderRoute('/studio/runtime?device_id=device-A');

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('empty');
    });
    expect(screen.getByTestId('runtime-dashboard-empty-snapshot')).toBeInTheDocument();
    expect(screen.queryByTestId('runtime-dashboard-summary-panel')).not.toBeInTheDocument();
  });

  it('does not replace an unknown query device with the workspace default device', async () => {
    renderRoute('/studio/runtime?device_id=device-Z');

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-missing-device-context')).toBeInTheDocument();
    });
    expect(screen.getByTestId('runtime-dashboard-selected-device')).toHaveTextContent('device-Z');
    expect(mockRuntimeStatus).not.toHaveBeenCalledWith('device-A');
  });

  it('keeps the last snapshot and marks backend stream_state unavailable as degraded live state', async () => {
    renderRoute('/studio/runtime?device_id=device-A');

    await waitFor(() => {
      expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A');
    });
    act(() => {
      mockEventSources[0].open();
    });
    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live');
    });

    act(() => {
      mockEventSources[0].emit('stream_state', {
        device_id: 'device-A',
        stream_state: {
          state: 'unavailable',
          empty: false,
          degraded: false,
          unavailable: true,
          stale: false,
          reason: 'runtime status stream closed',
        },
        timestamp: '2026-05-29T10:00:05Z',
      });
    });

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('degraded');
    });
    expect(screen.getByTestId('runtime-dashboard-stream-state')).toHaveTextContent('unavailable');
    expect(screen.getByTestId('runtime-dashboard-collector-count')).toHaveTextContent('1');
    expect(screen.getByTestId('runtime-dashboard-logs-panel')).toHaveTextContent('runtime status stream closed');
  });
});
