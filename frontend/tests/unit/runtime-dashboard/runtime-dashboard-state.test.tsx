import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Routes, Route, useLocation } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { RuntimeDashboardRoute } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute';
import type { Device, Point, RuntimeStatus } from '../../../src/types/datalink';

const {
  mockDevices,
  mockPoints,
  mockRuntimeStatus,
  mockEventSources,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockRuntimeStatus: vi.fn<() => Promise<RuntimeStatus>>(),
  mockEventSources: [] as MockEventSource[],
}));

type EventHandler = (event: MessageEvent<string>) => void;

class MockEventSource {
  url: string;
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

  close() {}

  emit(type: string, data: unknown) {
    const payload = { data: JSON.stringify(data) } as MessageEvent<string>;
    this.listeners.get(type)?.forEach((listener) => listener(payload));
  }

  fail() {
    this.onerror?.(new Event('error'));
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

vi.mock('../../../src/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({
    data: mockDevices,
    isLoading: false,
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
    getStreamUrl: vi.fn((deviceId: string, pointIds?: string[]) => {
      const params = new URLSearchParams({ device_id: deviceId });
      if (pointIds && pointIds.length > 0) {
        params.set('point_ids', pointIds.join(','));
      }
      return `/api/v1/datalink/runtime/stream?${params.toString()}`;
    }),
  },
}));

function LocationProbe() {
  const location = useLocation();

  return <div data-testid="runtime-dashboard-location">{location.search}</div>;
}

function renderRoute(initialEntry: string) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialEntry]}>
        <Routes>
          <Route
            path="/studio/runtime"
            element={
              <>
                <RuntimeDashboardRoute />
                <LocationProbe />
              </>
            }
          />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('runtime dashboard route state', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockEventSources.splice(0, mockEventSources.length);
    vi.stubGlobal('EventSource', MockEventSource);

    mockDevices.splice(0, mockDevices.length, {
      id: 'device-A',
      name: 'Mixer PLC',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '2026-05-29T00:00:00Z',
      updated_at: '2026-05-29T00:00:00Z',
    }, {
      id: 'device-B',
      name: 'Filler PLC',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '2026-05-29T00:00:00Z',
      updated_at: '2026-05-29T00:00:00Z',
    });

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
    } as Point, {
      id: 'point-2',
      device_id: 'device-B',
      name: 'Level',
      description: '',
      address: '40002',
      data_type: 'float32',
      enabled: true,
      polling_group_id: 'group-B',
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
      metrics: {
        collected_total: 10,
        write_success_total: 8,
        write_error_total: 0,
        mapping_error_total: 0,
        point_state_error_total: 0,
      },
      collectors: [{
        device_id: 'device-A',
        device_name: 'Mixer PLC',
        protocol: 'modbus_tcp',
        status: 'running',
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

  it('enters missing-device-context when device_id is absent', () => {
    renderRoute('/studio/runtime');

    expect(screen.getByTestId('runtime-dashboard-missing-device-context')).toBeInTheDocument();
    expect(mockRuntimeStatus).not.toHaveBeenCalled();
  });

  it('updates query state when switching device', async () => {
    renderRoute('/studio/runtime?device_id=device-A');

    await screen.findByTestId('runtime-dashboard-route');

    fireEvent.click(screen.getByRole('button', { name: 'Filler PLC' }));

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-location')).toHaveTextContent('?device_id=device-B');
    });
  });

  it('transitions from loading to live after snapshot and stream attach', async () => {
    renderRoute('/studio/runtime?device_id=device-A');

    expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('loading');

    await waitFor(() => {
      expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A');
    });

    expect(mockEventSources).toHaveLength(1);

    act(() => {
      mockEventSources[0].open();
    });

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live');
    });
  });

  it('keeps the last snapshot and enters degraded when the stream fails', async () => {
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
      mockEventSources[0].fail();
    });

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('degraded');
    });

    expect(screen.getByTestId('runtime-dashboard-selected-device')).toHaveTextContent('Mixer PLC');
    expect(screen.getByTestId('runtime-dashboard-collector-count')).toHaveTextContent('1');
  });

  it('surfaces a recoverable error state when the snapshot load fails', async () => {
    mockRuntimeStatus.mockReset();
    mockRuntimeStatus.mockRejectedValueOnce(new Error('snapshot failed'));
    mockRuntimeStatus.mockResolvedValueOnce({
      running: true,
      uptime_seconds: 20,
      metrics: {
        collected_total: 12,
        write_success_total: 9,
        write_error_total: 0,
        mapping_error_total: 0,
        point_state_error_total: 0,
      },
      collectors: [{
        device_id: 'device-A',
        device_name: 'Mixer PLC',
        protocol: 'modbus_tcp',
        status: 'running',
        points_total: 1,
        points_healthy: 1,
        points_stale: 0,
        points_error: 0,
        last_read_at: '2026-05-29T00:01:00Z',
        last_error: null,
        breaker_state: 'closed',
      }],
    });

    renderRoute('/studio/runtime?device_id=device-A');

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('error');
    });

    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' }));

    await waitFor(() => {
      expect(mockRuntimeStatus).toHaveBeenCalledTimes(2);
    });
  });

  it('updates logs panel when points go stale or recover without flooding', async () => {
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
      mockEventSources[0].emit('value', {
        device_id: 'device-A',
        point_id: 'point-1',
        address: '40001',
        raw_value: 100,
        transformed_value: 100,
        quality: 'bad',
        stale: true,
        timestamp: '2026-05-29T10:00:00Z',
      });
    });

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-logs-panel')).toHaveTextContent('Point 40001 went stale.');
    });

    act(() => {
      mockEventSources[0].emit('value', {
        device_id: 'device-A',
        point_id: 'point-1',
        address: '40001',
        raw_value: 101,
        transformed_value: 101,
        quality: 'bad',
        stale: true,
        timestamp: '2026-05-29T10:00:01Z',
      });
    });

    const logsText = screen.getByTestId('runtime-dashboard-logs-panel').innerHTML;
    const occurrences = (logsText.match(/went stale/g) || []).length;
    expect(occurrences).toBe(1);

    act(() => {
      mockEventSources[0].emit('value', {
        device_id: 'device-A',
        point_id: 'point-1',
        address: '40001',
        raw_value: 102,
        transformed_value: 102,
        quality: 'good',
        stale: false,
        timestamp: '2026-05-29T10:00:02Z',
      });
    });

    await waitFor(() => {
      expect(screen.getByTestId('runtime-dashboard-logs-panel')).toHaveTextContent('Point 40001 recovered.');
    });
  });
});
