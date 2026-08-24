import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, useLocation } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppRoutes } from '../../../src/App';

vi.mock('../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({
  default: (props: { runtimeReturnFocus?: { step?: number; issueCode?: string | null } | null }) => (
    <div
      data-testid="workbench-v2-root"
      data-focus-step={props.runtimeReturnFocus?.step ?? ''}
      data-focus-issue={props.runtimeReturnFocus?.issueCode ?? ''}
    >
      Workbench V2
    </div>
  ),
}));

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
    i18n: { changeLanguage: () => Promise.resolve() },
  }),
}));

describe('App Routing Integration', () => {
  function LocationProbe() {
    const location = useLocation();
    return <output data-testid="location">{location.pathname}{location.search}{location.hash}</output>;
  }

  function renderRoute(initialEntry: string) {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
        },
      },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={[initialEntry]}>
          <AppRoutes />
          <LocationProbe />
        </MemoryRouter>
      </QueryClientProvider>
    );
  }

  it('renders DatalinkWorkbenchV2Page on /', async () => {
    renderRoute('/');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on an unknown route fallback', async () => {
    renderRoute('/does-not-exist');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on /datalink', async () => {
    renderRoute('/datalink');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on /datalink/workbench', async () => {
    renderRoute('/datalink/workbench');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('converges nested legacy datalink routes directly to V2 with query and hash', async () => {
    renderRoute('/datalink/workbench/legacy-output?step=output&target=database#legacy');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.getByTestId('location')).toHaveTextContent(
      '/studio/v2?step=output&target=database#legacy',
    );
  });

  it('projects only the supported Local Modbus section into V2', async () => {
    renderRoute('/datalink/local-modbus?section=settings&step=device&target=database#legacy');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.getByTestId('location')).toHaveTextContent(
      '/studio/v2?step=output&target=modbus&section=settings',
    );
  });

  it('uses the generic V2 fallback on the retired /studio path', async () => {
    renderRoute('/studio');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on /studio/v2', async () => {
    renderRoute('/studio/v2');

    expect(await screen.findByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('passes runtime readiness focus query into /studio/v2', async () => {
    renderRoute('/studio/v2?step=4&focus=readiness&issue=database-target-missing');

    expect(await screen.findByTestId('workbench-v2-root')).toHaveAttribute('data-focus-step', '4');
    expect(screen.getByTestId('workbench-v2-root')).toHaveAttribute(
      'data-focus-issue',
      'database-target-missing',
    );
  });
});
