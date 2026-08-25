import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RuntimeDashboardPage } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardPage';
import { createRuntimeSetupFixture } from './runtimeSetupFixture';

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
  setupContext: createRuntimeSetupFixture(),
  onSelectDevice: vi.fn(),
  onRetrySnapshot: vi.fn(),
  onReconnectStream: vi.fn(),
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

  it('shows configured workspace setup even when runtime data is only live counters', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        selectedDevice={{
          ...baseState.selectedDevice,
          projection_alignment: 'stale',
          runtime_projection_version: 'projection-v12',
          workspace_projection_version: 'projection-v13',
        }}
      />,
    );

    const panel = screen.getByTestId('runtime-dashboard-setup-context');
    expect(panel).toHaveTextContent('Workspace readiness');
    expect(panel).toHaveTextContent('2 blockers');
    expect(panel).toHaveTextContent('Step 4');
    expect(panel).toHaveTextContent('database-target-missing');
    expect(panel).toHaveTextContent('40001');
    expect(panel).toHaveTextContent('8');
    expect(panel).toHaveTextContent('line01.temp.inlet');
    expect(panel).toHaveTextContent('入口溫度');
    expect(panel).toHaveTextContent('PostgreSQL Connector');
    expect(panel).toHaveTextContent('gateway_metrics');
    expect(panel).toHaveTextContent('sensor_readings');
    expect(panel).toHaveTextContent('Projection stale');
    expect(panel).toHaveTextContent('projection-v12');
    expect(panel).toHaveTextContent('projection-v13');
  });

  it('links readiness blockers back to their Studio V2 owner step', () => {
    render(<RuntimeDashboardPage {...baseState} />);

    expect(screen.getByTestId('runtime-dashboard-setup-return-link')).toHaveAttribute(
      'href',
      '/studio/v2?step=4&focus=readiness&issue=database-target-missing',
    );
    expect(screen.getByTestId('runtime-dashboard-setup-fix-database-target-missing')).toHaveAttribute(
      'href',
      '/studio/v2?step=4&focus=readiness&issue=database-target-missing',
    );
    expect(screen.getByTestId('runtime-dashboard-setup-fix-tag-missing')).toHaveAttribute(
      'href',
      '/studio/v2?step=3&focus=readiness&issue=tag-missing',
    );
  });

  it('normalizes tag-missing blockers to Step 3 even if the backend step is stale', () => {
    const setupContext = createRuntimeSetupFixture();
    setupContext.readiness_summary!.issues = [{
      code: 'tag-missing',
      severity: 'blocking',
      step: 'Step 2',
      scope: 'point-2',
      message: 'derived point is missing its persisted tag',
    }];

    render(<RuntimeDashboardPage {...baseState} setupContext={setupContext} />);

    const panel = screen.getByTestId('runtime-dashboard-setup-context');
    expect(panel).toHaveTextContent('Step 3');
    expect(panel).not.toHaveTextContent('Step 2');
    expect(screen.getByTestId('runtime-dashboard-setup-return-link')).toHaveAttribute(
      'href',
      '/studio/v2?step=3&focus=readiness&issue=tag-missing',
    );
    expect(screen.getByTestId('runtime-dashboard-setup-fix-tag-missing')).toHaveTextContent(
      'Fix in Step 3',
    );
  });

  it('explains point-missing blockers as a Step 2 source rule re-save', () => {
    const setupContext = createRuntimeSetupFixture();
    setupContext.readiness_summary!.issues = [{
      code: 'point-missing',
      severity: 'blocking',
      step: 'Step 2',
      scope: 'rule-A',
      message: 'source rule link is missing its derived point',
    }];

    render(<RuntimeDashboardPage {...baseState} setupContext={setupContext} />);

    const panel = screen.getByTestId('runtime-dashboard-setup-context');
    expect(panel).toHaveTextContent('Step 2');
    expect(panel).toHaveTextContent('point-missing');
    expect(panel).toHaveTextContent(
      'Fix in Step 2: re-save this source rule to rebuild the missing derived point.',
    );
    expect(screen.getByTestId('runtime-dashboard-setup-fix-point-missing')).toHaveAttribute(
      'href',
      '/studio/v2?step=2&focus=readiness&issue=point-missing',
    );
  });

  it('uses SPA navigation for Studio V2 readiness fixes when available', () => {
    const navigateTo = vi.fn();

    render(<RuntimeDashboardPage {...baseState} navigateTo={navigateTo} />);

    fireEvent.click(screen.getByTestId('runtime-dashboard-setup-fix-database-target-missing'));

    expect(navigateTo).toHaveBeenCalledWith(
      '/studio/v2?step=4&focus=readiness&issue=database-target-missing',
    );
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
            failure_code: 'runtime_snapshot_unavailable',
            failure_reason: 'permission denied: secret backend detail',
            request_id: 'req-diagnostics-2',
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
    expect(panel).toHaveTextContent('Runtime snapshot is currently unavailable.');
    expect(panel).toHaveTextContent('Request ID: req-diagnostics-2');
    expect(panel).toHaveTextContent('Retry snapshot');
    expect(panel).not.toHaveTextContent('permission denied');
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

  it('reconnects the SSE stream for a retryable runtime stream failure', () => {
    const onReconnectStream = vi.fn();
    const onRetrySnapshot = vi.fn();
    render(
      <RuntimeDashboardPage
        {...baseState}
        routeState="degraded"
        streamRecovery={{
          code: 'runtime_stream_unavailable',
          requestId: 'req-stream-7',
          retryable: true,
        }}
        onReconnectStream={onReconnectStream}
        onRetrySnapshot={onRetrySnapshot}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Reconnect live stream' }));

    expect(onReconnectStream).toHaveBeenCalledTimes(1);
    expect(onRetrySnapshot).not.toHaveBeenCalled();
  });

  it('renders typed runtime error copy and request id without raw diagnostics', () => {
    render(
      <RuntimeDashboardPage
        {...baseState}
        routeState="error"
        snapshotError={{
          code: 'runtime_snapshot_unavailable',
          title: 'runtime_snapshot_unavailable',
          message: 'errors.runtime_snapshot_unavailable',
          requestId: 'req-runtime-42',
          retryable: true,
        }}
      />,
    );

    expect(screen.getByTestId('runtime-dashboard-error-code')).toHaveTextContent(
      'runtime_snapshot_unavailable',
    );
    expect(screen.getByTestId('runtime-dashboard-error-request-id')).toHaveTextContent(
      'req-runtime-42',
    );
    expect(screen.getByTestId('runtime-dashboard-snapshot-error')).toHaveTextContent(
      'errors.runtime_snapshot_unavailable',
    );
    expect(screen.getByTestId('runtime-dashboard-snapshot-error')).not.toHaveTextContent(
      'dial tcp',
    );
  });

  it('keeps snapshot failures on the snapshot retry action', () => {
    const onRetrySnapshot = vi.fn();
    render(
      <RuntimeDashboardPage
        {...baseState}
        routeState="error"
        snapshotError={{
          code: 'runtime_snapshot_unavailable',
          title: 'runtime_snapshot_unavailable',
          message: 'errors.runtime_snapshot_unavailable',
          retryable: true,
        }}
        onRetrySnapshot={onRetrySnapshot}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' }));

    expect(onRetrySnapshot).toHaveBeenCalledTimes(1);
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

  it('keeps unavailable devices visible with safe localized recovery guidance', () => {
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
    expect(screen.getByTestId('runtime-dashboard-device-switcher')).toHaveTextContent('Review the device setup in Studio V2 and retry.');
    expect(screen.getByTestId('runtime-dashboard-device-switcher')).not.toHaveTextContent('invalid Step 1 configuration');
  });
});
