import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { WorkbenchV2Shell } from '../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { StudioV2WorkspaceReadinessSummary } from '../../../src/types/studioV2WorkspaceReadiness';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
    i18n: { changeLanguage: () => Promise.resolve() },
  }),
}));

describe('Workbench V2 workspace readiness surfaces', () => {
  beforeEach(() => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1440,
    });
  });

  it('renders the same workspace readiness summary in shell banner and summary rail', () => {
    const readinessSummary: StudioV2WorkspaceReadinessSummary = {
      ready: false,
      blocking_count: 1,
      warning_count: 1,
      issues: [
        {
          code: 'device-probe-required',
          severity: 'blocking',
          step: 'Step 1',
          scope: 'dev-A',
          message: 'probe diagnostics failed',
        },
        {
          code: 'database-connector-unreachable',
          severity: 'warning',
          step: 'Step 4',
          scope: 'db-main',
          message: 'database connector requires attention',
        },
      ],
    };
    const actions = {
      state: INITIAL_STATE,
      setView: vi.fn(),
      setCurrent: vi.fn(),
      completeStep: vi.fn(),
      toggleSidebar: vi.fn(),
      toggleSummaryRail: vi.fn(),
      setSidebarCollapsed: vi.fn(),
      setShowSummaryRail: vi.fn(),
      resetFlow: vi.fn(),
      selectRule: vi.fn(),
      dispatch: vi.fn(),
    };
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <WorkbenchV2Shell
          state={INITIAL_STATE}
          actions={actions}
          workspaceReadiness={readinessSummary}
        />
      </QueryClientProvider>,
    );

    expect(screen.getByTestId('workbench-v2-readiness-banner')).toHaveTextContent('device-probe-required');
    expect(screen.getByTestId('workbench-v2-readiness-banner')).toHaveTextContent('database-connector-unreachable');
    expect(screen.getByTestId('summary-rail-readiness')).toHaveTextContent('device-probe-required');
    expect(screen.getByTestId('summary-rail-readiness')).toHaveTextContent('database-connector-unreachable');
  });
});
