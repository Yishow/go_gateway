import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RuntimeDashboardPage } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardPage';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, fallback?: unknown) =>
      typeof fallback === 'string' ? fallback : key,
  }),
}));

const baseState = {
  routeState: 'live' as const,
  selectedDeviceId: 'device-A',
  selectedDevice: {
    device_id: 'device-A',
    name: 'Mixer PLC',
    protocol: 'modbus_tcp' as const,
    running: true,
    availability_status: 'available' as const,
    availability_reason: null,
  },
  devices: [{
    device_id: 'device-A',
    name: 'Mixer PLC',
    protocol: 'modbus_tcp' as const,
    running: true,
    availability_status: 'available' as const,
    availability_reason: null,
  }],
  snapshot: {
    running: true,
    uptime_seconds: 18,
    metrics: {
      collected_total: 128,
      write_success_total: 128,
      write_error_total: 0,
      mapping_error_total: 0,
      point_state_error_total: 0,
    },
    collectors: [{
      device_id: 'device-A',
      device_name: 'Mixer PLC',
      protocol: 'modbus_tcp',
      status: 'running' as const,
      availability_status: 'available' as const,
      availability_reason: null,
      running: true,
      points_total: 4,
      points_healthy: 3,
      points_stale: 1,
      points_error: 0,
      last_read_at: '2026-05-29T10:00:00Z',
      last_error: null,
      breaker_state: 'closed',
    }],
  },
  snapshotError: null,
  liveValues: {
    'pt-1': {
      device_id: 'device-A',
      point_id: 'pt-1',
      address: '40001',
      raw_value: 150,
      transformed_value: 20,
      quality: 'good' as const,
      stale: false,
      timestamp: '2026-05-29T10:00:05Z',
    },
  },
  streamState: 'connected' as const,
  onSelectDevice: vi.fn(),
  onRetrySnapshot: vi.fn(),
};

describe('RuntimeDashboardPage', () => {
  it('renders the focused runtime dashboard surface for the selected device', () => {
    render(<RuntimeDashboardPage {...baseState} />);

    expect(screen.getByTestId('runtime-dashboard-route')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-header')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-summary-panel')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-health-panel')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-live-points-table')).toBeInTheDocument();
  });

  it('shows only backend-supported summary and health contract fields', () => {
    render(<RuntimeDashboardPage {...baseState} />);

    expect(screen.getByTestId('runtime-dashboard-summary-panel')).toHaveTextContent('Collected total');
    expect(screen.getByTestId('runtime-dashboard-summary-panel')).toHaveTextContent('128');
    expect(screen.getByTestId('runtime-dashboard-health-panel')).toHaveTextContent('closed');
    expect(screen.queryByText(/queue backlog/i)).not.toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-logs-panel')).toBeInTheDocument();
  });

  it('shows latest runtime diagnostics failure context for the selected device', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        snapshot={{
          ...baseState.snapshot,
          diagnostics: [{
            scope: 'device:device-A',
            device_id: 'device-A',
            point_id: 'pt-1',
            tag_id: 'tag-pressure',
            last_success_at: '2026-05-29T10:10:00Z',
            last_failure_at: '2026-05-29T10:12:00Z',
            latest_successful_stage: 'runtime_projection',
            failure_stage: 'database_delivery',
            failure_reason: 'permission denied',
            stages: [
              { stage: 'collector', status: 'success', observed_at: '2026-05-29T10:12:00Z' },
              { stage: 'mapping', status: 'success', observed_at: '2026-05-29T10:12:00Z' },
              {
                stage: 'runtime_projection',
                status: 'success',
                observed_at: '2026-05-29T10:12:00Z',
              },
              {
                stage: 'database_delivery',
                status: 'failed',
                observed_at: '2026-05-29T10:12:00Z',
                reason: 'permission denied',
              },
            ],
          }],
        }}
      />,
    );

    const panel = screen.getByTestId('runtime-dashboard-diagnostics-panel');
    expect(panel).toHaveTextContent('database_delivery');
    expect(panel).toHaveTextContent('permission denied');
    expect(panel).toHaveTextContent('2026-05-29T10:12:00Z');
    expect(panel).toHaveTextContent('runtime_projection');
  });

  it('marks unsupported summary metrics unavailable instead of filling zero defaults', () => {
    const snapshotWithoutMetrics = {
      ...baseState.snapshot,
      metrics: undefined,
    };

    render(
      <RuntimeDashboardPage
        {...baseState}
        snapshot={snapshotWithoutMetrics}
      />,
    );

    const summary = screen.getByTestId('runtime-dashboard-summary-panel');
    expect(summary).toHaveTextContent('Write success total');
    expect(summary).toHaveTextContent('Unavailable');
  });

  it('keeps the last successful data visible while showing a degraded banner', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        routeState="degraded"
        streamState="error"
      />,
    );

    expect(screen.getByTestId('runtime-dashboard-live-state-banner')).toHaveTextContent('Live stream degraded');
    expect(screen.getByTestId('runtime-dashboard-header')).toHaveTextContent('Mixer PLC');
    expect(screen.getByText('2026-05-29T10:00:05Z')).toBeInTheDocument();
  });

  it('shows a live-points placeholder when the snapshot is ready but no live values have arrived', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        liveValues={{}}
      />,
    );

    expect(screen.getByTestId('runtime-dashboard-live-points-placeholder')).toBeInTheDocument();
  });

  it('renders a focused empty state when device context is missing', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        routeState="missing-device-context"
        selectedDeviceId={null}
        selectedDevice={null}
        snapshot={null}
      />,
    );

    expect(screen.getByTestId('runtime-dashboard-missing-device-context')).toBeInTheDocument();
    expect(screen.queryByTestId('runtime-dashboard-summary-panel')).not.toBeInTheDocument();
  });

  it('renders backend empty snapshot state without showing synthetic summary panels', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        routeState="empty"
        snapshot={null}
        liveValues={{}}
      />,
    );

    expect(screen.getByTestId('runtime-dashboard-empty-snapshot')).toBeInTheDocument();
    expect(screen.queryByTestId('runtime-dashboard-summary-panel')).not.toBeInTheDocument();
    expect(screen.queryByTestId('runtime-dashboard-health-panel')).not.toBeInTheDocument();
  });

  it('does not borrow another device collector when selected device has no collector', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        snapshot={{
          ...baseState.snapshot,
          collectors: [{
            ...baseState.snapshot.collectors[0],
            device_id: 'device-B',
            device_name: 'Filler PLC',
            points_total: 9,
          }],
        }}
      />,
    );

    const healthPanel = screen.getByTestId('runtime-dashboard-health-panel');
    expect(healthPanel).toHaveTextContent('Collector status will appear after the first runtime snapshot.');
    expect(healthPanel).not.toHaveTextContent('9');
  });

  it('keeps unavailable devices visible in the switcher with their reason', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        devices={[
          ...baseState.devices,
          {
            device_id: 'device-B',
            name: 'Filler PLC',
            protocol: 'modbus_tcp',
            running: false,
            availability_status: 'unavailable',
            availability_reason: 'invalid Step 1 configuration',
          },
        ]}
      />,
    );

    expect(screen.getByTestId('runtime-dashboard-device-switcher')).toHaveTextContent('Filler PLC');
    expect(screen.getByTestId('runtime-dashboard-device-switcher')).toHaveTextContent('invalid Step 1 configuration');
  });
});
