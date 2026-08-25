import type { RuntimeCollectorStatus } from '../../../src/types/datalink';
import { act, screen, waitFor, fireEvent } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { mockEventSources, mockRuntimeStatus, resetRuntimeDashboardFixture, renderRoute } from './runtime-dashboard-state.fixture';

describe('runtime dashboard stream state', () => {
  beforeEach(resetRuntimeDashboardFixture);

  it('surfaces projection drift from the runtime snapshot', async () => {
    const staleCollector: RuntimeCollectorStatus & { projection_alignment: string; runtime_projection_version: string; workspace_projection_version: string } = { device_id: 'device-A', device_name: 'Mixer PLC', protocol: 'modbus_tcp', status: 'running', availability_status: 'available', availability_reason: null, running: true, points_total: 1, points_healthy: 1, points_stale: 0, points_error: 0, last_read_at: '2026-05-29T00:00:00Z', last_error: null, breaker_state: 'closed', projection_alignment: 'stale', runtime_projection_version: 'projection-v12', workspace_projection_version: 'projection-v13' };
    mockRuntimeStatus.mockResolvedValueOnce({ running: true, uptime_seconds: 12, collectors: [staleCollector] });
    renderRoute('/studio/runtime?device_id=device-A');
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A')); act(() => mockEventSources[0].open());
    await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live'));
    expect(screen.getByTestId('runtime-dashboard-health-panel')).toHaveTextContent('stale');
    expect(screen.getByTestId('runtime-dashboard-health-panel')).toHaveTextContent('projection-v12');
    expect(screen.getByTestId('runtime-dashboard-health-panel')).toHaveTextContent('projection-v13');
  });

  it('keeps the last snapshot and enters degraded when the stream fails', async () => {
    renderRoute('/studio/runtime?device_id=device-A'); await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A'));
    act(() => mockEventSources[0].open()); await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live'));
    act(() => mockEventSources[0].fail()); await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('degraded'));
    expect(screen.getByTestId('runtime-dashboard-selected-device')).toHaveTextContent('Mixer PLC');
    expect(screen.getByTestId('runtime-dashboard-collector-count')).toHaveTextContent('1');
  });

  it('surfaces a recoverable error state when the snapshot load fails', async () => {
    mockRuntimeStatus.mockReset(); mockRuntimeStatus.mockRejectedValueOnce(new Error('snapshot failed')); mockRuntimeStatus.mockResolvedValueOnce({ running: true, uptime_seconds: 20, metrics: { collected_total: 12, write_success_total: 9, write_error_total: 0, mapping_error_total: 0, point_state_error_total: 0 }, collectors: [{ device_id: 'device-A', device_name: 'Mixer PLC', protocol: 'modbus_tcp', status: 'running', availability_status: 'available', availability_reason: null, running: true, points_total: 1, points_healthy: 1, points_stale: 0, points_error: 0, last_read_at: '2026-05-29T00:01:00Z', last_error: null, breaker_state: 'closed' }] });
    renderRoute('/studio/runtime?device_id=device-A'); await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('error'));
    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' })); await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledTimes(2));
  });

  it('updates logs panel when points go stale or recover without flooding', async () => {
    renderRoute('/studio/runtime?device_id=device-A'); await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A'));
    act(() => mockEventSources[0].open()); await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live'));
    const emitValue = (quality: string, stale: boolean, value: number, timestamp: string) => act(() => mockEventSources[0].emit('value', { device_id: 'device-A', point_id: 'point-1', address: '40001', raw_value: value, transformed_value: value, quality, stale, timestamp }));
    emitValue('bad', true, 100, '2026-05-29T10:00:00Z'); await waitFor(() => expect(screen.getByTestId('runtime-dashboard-logs-panel')).toHaveTextContent('logs.point_stale'));
    emitValue('bad', true, 101, '2026-05-29T10:00:01Z');
    const logsText = screen.getByTestId('runtime-dashboard-logs-panel').innerHTML; expect((logsText.match(/logs\.point_stale/g) || []).length).toBe(1);
    emitValue('good', false, 102, '2026-05-29T10:00:02Z'); await waitFor(() => expect(screen.getByTestId('runtime-dashboard-logs-panel')).toHaveTextContent('logs.point_recovered'));
  });
});
