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

  constructor() {
    mockEventSources.push(this);
  }

  addEventListener() {}
  removeEventListener() {}
  close() {}

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
        failure_reason: 'permission denied',
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
    expect(panel).toHaveTextContent('permission denied');
    expect(panel).toHaveTextContent('2026-05-29T10:12:00Z');
  });
});
