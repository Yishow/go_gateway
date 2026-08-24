import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import { AppRoutes } from '../../../src/App';

const routeStateMock = vi.hoisted(() => ({
  routeState: 'live',
  selectedDeviceId: 'device-A',
  selectedDevice: {
    device_id: 'device-A',
    name: 'Mixer PLC',
    protocol: 'modbus_tcp',
    running: true,
    availability_status: 'available',
    availability_reason: null,
  },
  devices: [{
    device_id: 'device-A',
    name: 'Mixer PLC',
    protocol: 'modbus_tcp',
    running: true,
    availability_status: 'available',
    availability_reason: null,
  }],
  snapshot: null,
  snapshotError: null,
  liveValues: {},
  streamState: 'connected',
  logs: [],
  onSelectDevice: vi.fn(),
  onRetrySnapshot: vi.fn(),
}));

vi.mock('../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({
  default: () => <div data-testid="workbench-v2-root">Mock Workbench V2</div>,
}));

vi.mock('../../../src/features/datalink/runtime-dashboard/useRuntimeDashboardState', () => ({
  useRuntimeDashboardState: () => routeStateMock,
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, fallback?: unknown) =>
      typeof fallback === 'string' ? fallback : key,
  }),
}));

describe('runtime dashboard route', () => {
  function renderRoutes(initialEntries: string[]) {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
      },
    });

    return render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={initialEntries}>
          <AppRoutes />
        </MemoryRouter>
      </QueryClientProvider>,
    );
  }

  it('renders the dedicated runtime dashboard route when device_id is present', async () => {
    renderRoutes(['/studio/runtime?device_id=device-A']);

    expect(await screen.findByTestId('runtime-dashboard-route')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live');
    expect(screen.queryByTestId('workbench-v2-root')).not.toBeInTheDocument();
  });

  it('renders the dedicated runtime dashboard route when device_id is absent', async () => {
    renderRoutes(['/studio/runtime']);

    expect(await screen.findByTestId('runtime-dashboard-route')).toBeInTheDocument();
    expect(screen.queryByTestId('workbench-v2-root')).not.toBeInTheDocument();
  });

  it('keeps /studio/v2 on the setup flow', async () => {
    renderRoutes(['/studio/v2']);

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('runtime-dashboard-route')).not.toBeInTheDocument();
  });
});
