import { useRuntimeDashboardState } from './useRuntimeDashboardState';

export function RuntimeDashboardRoute() {
  const state = useRuntimeDashboardState();

  if (state.routeState === 'missing-device-context') {
    return <div data-testid="runtime-dashboard-missing-device-context" />;
  }

  if (state.routeState === 'error') {
    return (
      <div data-testid="runtime-dashboard-route">
        <div data-testid="runtime-dashboard-route-state">error</div>
        <div data-testid="runtime-dashboard-selected-device">
          {state.selectedDevice?.name ?? state.selectedDeviceId ?? ''}
        </div>
        <div data-testid="runtime-dashboard-snapshot-error">
          {state.snapshotError ?? 'snapshot error'}
        </div>
        <button type="button" onClick={() => void state.onRetrySnapshot()}>
          Retry snapshot
        </button>
      </div>
    );
  }

  return (
    <div data-testid="runtime-dashboard-route">
      <div data-testid="runtime-dashboard-route-state">{state.routeState}</div>
      <div data-testid="runtime-dashboard-selected-device">
        {state.selectedDevice?.name ?? state.selectedDeviceId ?? ''}
      </div>
      <div data-testid="runtime-dashboard-collector-count">
        {state.snapshot?.collectors.length ?? 0}
      </div>
      <div data-testid="runtime-dashboard-stream-state">{state.streamState}</div>
      <div>
        {state.devices.map((device) => (
          <button
            key={device.id}
            type="button"
            onClick={() => state.onSelectDevice(device.id)}
          >
            {device.name}
          </button>
        ))}
      </div>
    </div>
  );
}
