import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import { AppRoutes } from '../../../src/App';

const routeStateMock = vi.hoisted(() => ({
  routeState: 'live',
  selectedDeviceId: 'device-A',
  selectedDevice: {
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
  },
  devices: [],
  snapshot: null,
  liveValues: {},
  streamState: 'connected',
  onSelectDevice: vi.fn(),
}));

vi.mock('../../../src/pages/datalink/workbench/DatalinkWorkbenchPage', () => ({
  default: () => <div data-testid="legacy-workbench">Legacy DatalinkWorkbenchPage</div>,
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

  it('renders the dedicated runtime dashboard route when device_id is present', () => {
    renderRoutes(['/studio/runtime?device_id=device-A']);

    expect(screen.getByTestId('runtime-dashboard-route')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live');
    expect(screen.queryByTestId('workbench-v2-root')).not.toBeInTheDocument();
  });

  it('keeps /studio/v2 on the setup flow', () => {
    renderRoutes(['/studio/v2']);

    expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('runtime-dashboard-route')).not.toBeInTheDocument();
  });
});
