import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppRoutes } from '../../../src/App';

// Mock 舊版頁面以避免複雜的 Material-UI 和 Provider 依賴
vi.mock('../../../src/pages/datalink/workbench/DatalinkWorkbenchPage', () => ({
  default: () => <div data-testid="legacy-workbench">Legacy DatalinkWorkbenchPage</div>,
}));

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
    i18n: { changeLanguage: () => Promise.resolve() },
  }),
}));

describe('App Routing Integration', () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  it('renders Legacy DatalinkWorkbenchPage on /studio', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={['/studio']}>
          <AppRoutes />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByTestId('legacy-workbench')).toBeInTheDocument();
    expect(screen.queryByTestId('workbench-v2-root')).not.toBeInTheDocument();
  });

  it('renders DatalinkWorkbenchV2Page on /studio/v2', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={['/studio/v2']}>
          <AppRoutes />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    expect(screen.queryByTestId('legacy-workbench')).not.toBeInTheDocument();
  });

  it('redirects /datalink/workbench to /studio', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={['/datalink/workbench']}>
          <AppRoutes />
        </MemoryRouter>
      </QueryClientProvider>
    );

    // 應該會被 LegacyStudioRedirect 重定向到 /studio 並加載舊版頁面
    expect(screen.getByTestId('legacy-workbench')).toBeInTheDocument();
  });
});
