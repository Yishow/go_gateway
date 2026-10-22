import type { RuntimeCollectorStatus, RuntimeStatus } from '../../../src/types/datalink';
import { act, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { mockEventSources, mockRuntimeStatus, resetRuntimeDashboardFixture, renderRoute } from './runtime-dashboard-state.fixture';

function snapshotWithHealthy(pointsHealthy: number): RuntimeStatus {
  const collector: RuntimeCollectorStatus = { device_id: 'device-A', device_name: 'Mixer PLC', protocol: 'modbus_tcp', status: 'running', availability_status: 'available', availability_reason: null, running: true, points_total: 2, points_healthy: pointsHealthy, points_stale: 0, points_error: 0, last_read_at: '2026-05-29T00:00:00Z', last_error: null, breaker_state: 'closed' };
  return { running: true, uptime_seconds: 12, collectors: [collector] };
}

describe('runtime dashboard live snapshot refresh', () => {
  beforeEach(() => {
    resetRuntimeDashboardFixture();
    vi.useFakeTimers({ shouldAdvanceTime: true });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('keeps refreshing collector health while the stream stays live', async () => {
    mockRuntimeStatus.mockReset();
    mockRuntimeStatus.mockResolvedValueOnce(snapshotWithHealthy(0));
    mockRuntimeStatus.mockResolvedValue(snapshotWithHealthy(2));
    renderRoute('/studio/runtime?device_id=device-A');
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledTimes(1));
    act(() => mockEventSources[0].open());
    await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live'));
    expect(screen.getByTestId('runtime-dashboard-health-panel')).toHaveTextContent('0/2');

    await act(async () => { await vi.advanceTimersByTimeAsync(15_000); });

    await waitFor(() => expect(mockRuntimeStatus.mock.calls.length).toBeGreaterThan(1));
    await waitFor(() => expect(screen.getByTestId('runtime-dashboard-health-panel')).toHaveTextContent('2/2'));
    expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live');
  });
});
