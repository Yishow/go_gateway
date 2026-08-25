import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { RuntimeDashboardRoute } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute';
import type { RuntimeStatusWithDiagnostics } from '../../../src/types/runtimeDiagnostics';

const {
  mockRuntimeStatus,
  mockEventSources,
} = vi.hoisted(() => ({
  mockRuntimeStatus: vi.fn<() => Promise<RuntimeStatusWithDiagnostics>>(),
  mockEventSources: [] as MockEventSource[],
}));

class MockEventSource {
  onopen: ((event: Event) => void) | null = null;
  private listeners = new Map<string, Set<(event: MessageEvent<string>) => void>>();

  constructor() {
    mockEventSources.push(this);
  }

  addEventListener(type: string, listener: (event: MessageEvent<string>) => void) {
    const listeners = this.listeners.get(type) ?? new Set();
    listeners.add(listener);
    this.listeners.set(type, listeners);
  }
  removeEventListener(type: string, listener: (event: MessageEvent<string>) => void) {
    this.listeners.get(type)?.delete(listener);
  }
  close() {}

  emit(type: string, data: unknown) {
    const event = { data: JSON.stringify(data) } as MessageEvent<string>;
    this.listeners.get(type)?.forEach((listener) => listener(event));
  }

  open() {
    this.onopen?.(new Event('open'));
    this.emit('stream_state', {
      stream_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false },
    });
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
    data: {
      workspace_id: 'workspace-1',
      devices: [{
        device_id: 'device-A',
        name: 'Mixer PLC',
        protocol: 'modbus_tcp',
        running: true,
        availability_status: 'available',
        availability_reason: null,
      }],
      default_device_id: 'device-A',
    },
    isLoading: false,
    isError: false,
    error: null,
  }),
}));

vi.mock('../../../src/hooks/datalink/usePoints', () => ({
  usePointsQuery: () => ({
    data: [],
    isLoading: false,
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  runtimeAPI: {
    getStatus: mockRuntimeStatus,
    getStreamUrl: vi.fn(() => '/api/v1/datalink/runtime/stream?device_id=device-A'),
  },
}));

function renderRoute() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={['/studio/runtime?device_id=device-A']}>
        <Routes>
          <Route path="/studio/runtime" element={<RuntimeDashboardRoute />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe('runtime dashboard diagnostics state', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockEventSources.splice(0, mockEventSources.length);
    vi.stubGlobal('EventSource', MockEventSource);
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
        last_read_at: '2026-05-29T10:12:00Z',
        last_error: null,
        breaker_state: 'closed',
      }],
      diagnostics: [{
        scope: 'device:device-A',
        device_id: 'device-A',
        point_id: 'point-A',
        tag_id: 'tag-A',
        last_failure_at: '2026-05-29T10:12:00Z',
        latest_successful_stage: 'runtime_projection',
        failure_stage: 'database_delivery',
        failure_code: 'runtime_snapshot_unavailable',
        failure_reason: 'permission denied: secret backend detail',
        request_id: 'req-diagnostics-1',
        stages: [],
      }],
    });
  });

  it('preserves runtime diagnostics from the backend snapshot into the page state', async () => {
    renderRoute();

    await waitFor(() => {
      expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A');
    });
    act(() => {
      mockEventSources[0].open();
    });

    const panel = await screen.findByTestId('runtime-dashboard-diagnostics-panel');
    expect(panel).toHaveTextContent('database_delivery');
    expect(panel).toHaveTextContent('Runtime snapshot is currently unavailable.');
    expect(panel).toHaveTextContent('Request ID: req-diagnostics-1');
    expect(panel).toHaveTextContent('Retry snapshot');
    expect(panel).not.toHaveTextContent('permission denied');
    expect(panel).toHaveTextContent('2026-05-29T10:12:00Z');
  });

  it('uses Modbus Share delivery copy for its own failure stage without exposing backend detail', async () => {
    mockRuntimeStatus.mockResolvedValueOnce({
      running: false,
      uptime_seconds: 12,
      collectors: [],
      diagnostics: [{
        scope: 'device:device-A',
        device_id: 'device-A',
        last_failure_at: '2026-05-29T10:12:00Z',
        failure_stage: 'modbus_share_delivery',
        failure_code: 'modbus_share_delivery',
        failure_reason: 'secret backend failure detail',
        request_id: 'req-share-delivery-1',
        stages: [],
      }],
    });

    renderRoute();

    await waitFor(() => {
      expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A');
    });
    act(() => {
      mockEventSources[0].open();
    });

    const panel = await screen.findByTestId('runtime-dashboard-diagnostics-panel');
    expect(panel).toHaveTextContent('modbus_share_delivery');
    expect(panel).not.toHaveTextContent('database_delivery');
    expect(panel).toHaveTextContent('Modbus Share delivery is unavailable.');
    expect(panel).toHaveTextContent('Retry Modbus Share delivery');
    expect(panel).toHaveTextContent('Request ID: req-share-delivery-1');
    expect(panel).not.toHaveTextContent('secret backend failure detail');
  });

  it('merges database and Modbus Share delivery outcomes into diagnostics', async () => {
    mockRuntimeStatus.mockResolvedValueOnce({
      running: true,
      uptime_seconds: 12,
      collectors: [],
      diagnostics: [],
      database_delivery: [{
        device_id: 'device-A',
        point_id: 'point-db',
        tag_id: 'tag-db',
        status: 'failed',
        stages: ['collected', 'mapped', 'db_write_failed'],
        failed_stage: 'db_write',
        observed_at: '2026-05-29T10:12:01Z',
        last_failure_at: '2026-05-29T10:12:01Z',
      }],
      modbus_share_delivery: [{
        device_id: 'device-A',
        point_id: 'point-share',
        tag_id: 'tag-share',
        stage: 'share_write',
        observed_at: '2026-05-29T10:12:02Z',
      }],
    });

    renderRoute();
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A'));
    act(() => {
      mockEventSources[0].open();
    });

    const panel = await screen.findByTestId('runtime-dashboard-diagnostics-panel');
    expect(panel).toHaveTextContent('database_delivery');
    expect(panel).toHaveTextContent('modbus_share_delivery');
    expect(panel).toHaveTextContent('tag-db');
    expect(panel).toHaveTextContent('tag-share');
  });

  it('does not present disabled Modbus Share status as a delivery failure', async () => {
    mockRuntimeStatus.mockResolvedValueOnce({
      running: true,
      uptime_seconds: 12,
      collectors: [],
      diagnostics: [],
      modbus_share_delivery: [{
        device_id: 'device-A',
        point_id: 'point-share',
        tag_id: 'tag-share',
        status: 'disabled',
        stage: 'share_write',
        error: 'modbus share is disabled',
        observed_at: '2026-05-29T10:12:02Z',
      }],
    });

    renderRoute();
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A'));
    act(() => {
      mockEventSources[0].open();
    });

    const panel = await screen.findByTestId('runtime-dashboard-diagnostics-panel');
    expect(panel).toHaveTextContent('Diagnostics unavailable');
    expect(panel).not.toHaveTextContent('delivery');
    expect(panel).not.toHaveTextContent('modbus share is disabled');
  });
});
