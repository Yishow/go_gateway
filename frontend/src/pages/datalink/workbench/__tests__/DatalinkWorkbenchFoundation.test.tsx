import { fireEvent, render, screen } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import App from '../../../../App';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

vi.mock('@/contexts/ThemeContext', () => ({
  ThemeProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

vi.mock('@/contexts/ToastContext', () => ({
  ToastProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

vi.mock('@/components/CardMinimizeProvider', () => ({
  CardMinimizeProvider: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

vi.mock('@/components/Layout', () => ({
  default: ({ children }: { children: ReactNode }) => <div data-testid="layout-mock">{children}</div>,
}));

vi.mock('@/pages/TestPage', () => ({
  default: () => <div data-testid="test-page-mock">test-page</div>,
}));

vi.mock('@/pages/TemplatesPage', () => ({
  default: () => <div data-testid="templates-page-mock">templates-page</div>,
}));

vi.mock('@/pages/HistoryPage', () => ({
  default: () => <div data-testid="history-page-mock">history-page</div>,
}));

vi.mock('@/pages/ComparePage', () => ({
  default: () => <div data-testid="compare-page-mock">compare-page</div>,
}));

vi.mock('@/pages/AnalyzerPage', () => ({
  default: () => <div data-testid="analyzer-page-mock">analyzer-page</div>,
}));

vi.mock('@/pages/datalink/SmartDashboard', () => ({
  default: () => <div data-testid="smart-dashboard-mock">smart-dashboard</div>,
}));

vi.mock('@/pages/datalink/LocalModbusWorkbenchPage', () => ({
  default: () => <div data-testid="local-modbus-workbench-mock">local-modbus-workbench</div>,
}));

vi.mock('@/router/gateway', () => ({
  GatewayCreateEntryRedirect: () => <div data-testid="gateway-create-entry-redirect-mock" />,
  GatewayEntryRoute: () => <div data-testid="gateway-entry-route-mock" />,
  GatewayExpertWorkbenchRoute: () => <div data-testid="gateway-expert-workbench-route-mock" />,
  GatewayQuickSetupRoute: () => <div data-testid="gateway-quick-setup-route-mock" />,
}));

vi.mock('@/features/datalink/legacyRoutes', () => ({
  buildDashboardModalRedirect: (intent: string) => `/mock/dashboard/${intent}`,
  buildLegacyMigrationRedirect: (intent: string) => `/mock/legacy/${intent}`,
}));

describe('DatalinkWorkbench foundation route', () => {
  beforeEach(() => {
    window.history.pushState({}, '', '/datalink/workbench');
  });

  it('renders the new foundation shell through /datalink/workbench', () => {
    render(<App />);

    expect(screen.getByRole('heading', { name: 'workbench.title' })).toBeInTheDocument();
    expect(
      screen.getByRole('navigation', { name: 'workbench.stepNavigator.ariaLabel' }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole('button', { name: 'workbench.steps.device' }),
    ).toHaveAttribute('aria-current', 'step');
    expect(screen.getByText('workbench.placeholders.device')).toBeInTheDocument();
    expect(
      screen.getByRole('complementary', { name: 'workbench.actionDock.ariaLabel' }),
    ).toBeInTheDocument();
  });

  it('switches the active step from the step navigator', () => {
    render(<App />);

    fireEvent.click(screen.getByRole('button', { name: 'workbench.steps.source' }));

    expect(
      screen.getByRole('button', { name: 'workbench.steps.source' }),
    ).toHaveAttribute('aria-current', 'step');
    expect(screen.getByRole('heading', { name: 'workbench.steps.source' })).toBeInTheDocument();
    expect(screen.getByText('workbench.placeholders.source')).toBeInTheDocument();
  });
});
