import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppRoutes } from '../../../src/App';

// Mock 舊版頁面以避免複雜的 Material-UI 和 Provider 依賴
vi.mock('../../../src/pages/datalink/workbench/DatalinkWorkbenchPage', () => ({
  default: () => <div data-testid="legacy-workbench">Legacy DatalinkWorkbenchPage</div>,
}));

vi.mock('../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page', () => ({
  default: () => <div data-testid="workbench-v2-root">Workbench V2</div>,
}));

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
    i18n: { changeLanguage: () => Promise.resolve() },
  }),
}));

describe('App Routing Integration', () => {
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
        </MemoryRouter>
      </QueryClientProvider>
    );
  }

  it('renders DatalinkWorkbenchV2Page on /', () => {
    renderRoute('/');

    expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on an unknown route fallback', () => {
    renderRoute('/does-not-exist');

    expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on /datalink', () => {
    renderRoute('/datalink');

    expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on /datalink/workbench', () => {
    renderRoute('/datalink/workbench');

    expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('renders Legacy DatalinkWorkbenchPage on /studio', () => {
    renderRoute('/studio');

    expect(screen.getByTestId('legacy-workbench')).toBeInTheDocument();
    expect(screen.queryByTestId('workbench-v2-root')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on /studio/v2', () => {
    renderRoute('/studio/v2');

    expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });
});
