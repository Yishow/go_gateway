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
    id: 'device-A',
    name: 'Mixer PLC',
    description: '',
    protocol: 'modbus_tcp' as const,
    status: 'active' as const,
    connection_config: '{}',
    last_test_at: null,
    last_test_success: true,
    last_test_error: '',
    created_at: '2026-05-29T00:00:00Z',
    updated_at: '2026-05-29T00:00:00Z',
  },
  devices: [{
    id: 'device-A',
    name: 'Mixer PLC',
    description: '',
    protocol: 'modbus_tcp' as const,
    status: 'active' as const,
    connection_config: '{}',
    last_test_at: null,
    last_test_success: true,
    last_test_error: '',
    created_at: '2026-05-29T00:00:00Z',
    updated_at: '2026-05-29T00:00:00Z',
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
    expect(screen.queryByText(/diagnostic logs/i)).not.toBeInTheDocument();
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
});
