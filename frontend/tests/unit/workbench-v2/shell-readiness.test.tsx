import { describe, it, expect, vi, beforeEach } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { WorkbenchV2Shell } from '../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2WorkspaceReadinessSummary } from '../../../src/types/studioV2WorkspaceReadiness';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) => {
      if (key === 'runtime_return.point_missing_action') {
        return 'Re-save this source rule to rebuild the missing derived point and link.';
      }
      if (key === 'runtime_return.scope') {
        return `Scope: ${String(options?.scope ?? '')}`;
      }
      return key;
    },
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

  it('keeps readiness compact in the shell banner and summary rail', () => {
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

    expect(screen.getByTestId('workbench-v2-readiness-banner')).toHaveTextContent('workspace_readiness.blocked');
    expect(screen.getByTestId('workbench-v2-readiness-banner')).toHaveTextContent('workspace_readiness.blocking_count');
    expect(screen.getByTestId('workbench-v2-readiness-banner')).not.toHaveTextContent('device-probe-required');
    expect(screen.getByTestId('summary-rail-readiness')).toHaveTextContent('workspace_readiness.blocked');
    expect(screen.getByTestId('summary-rail-readiness')).not.toHaveTextContent('database-connector-unreachable');
  });

  it('focuses the runtime-reported readiness issue without resetting the saved flow', () => {
    const readinessSummary: StudioV2WorkspaceReadinessSummary = {
      ready: false,
      blocking_count: 1,
      warning_count: 0,
      issues: [{
        code: 'database-target-missing',
        severity: 'blocking',
        step: 'Step 4',
        scope: 'point-1',
        message: 'derived point is missing its persisted database target',
      }],
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
          runtimeReturnFocus={{ step: 4, issueCode: 'database-target-missing' }}
        />
      </QueryClientProvider>,
    );

    expect(actions.setView).toHaveBeenCalledWith('flow');
    expect(actions.setCurrent).toHaveBeenCalledWith(4);
    expect(screen.getByTestId('workbench-v2-runtime-return-focus')).toHaveTextContent('database-target-missing');
    expect(screen.getByTestId('workbench-v2-runtime-return-focus')).toHaveTextContent(
      'derived point is missing its persisted database target',
    );
  });

  it('explains Step 2 point-missing runtime returns with the repair target', () => {
    const readinessSummary: StudioV2WorkspaceReadinessSummary = {
      ready: false,
      blocking_count: 1,
      warning_count: 0,
      issues: [{
        code: 'point-missing',
        severity: 'blocking',
        step: 'Step 2',
        scope: 'rule-A',
        message: 'source rule link is missing its derived point',
      }],
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
          runtimeReturnFocus={{ step: 2, issueCode: 'point-missing' }}
        />
      </QueryClientProvider>,
    );

    const focusPanel = screen.getByTestId('workbench-v2-runtime-return-focus');
    expect(actions.setCurrent).toHaveBeenCalledWith(2);
    expect(focusPanel).toHaveTextContent('point-missing');
    expect(focusPanel).toHaveTextContent('source rule link is missing its derived point');
    expect(focusPanel).toHaveTextContent(
      'Re-save this source rule to rebuild the missing derived point and link.',
    );
    expect(focusPanel).toHaveTextContent('Scope: rule-A');
  });

  it('offers a direct runtime entry when devices are already running', () => {
    const state: WorkbenchV2State = {
      ...INITIAL_STATE,
      devices: [{
        id: 'dev-A',
        name: 'Mixer PLC',
        description: '',
        protocol: 'modbus_tcp',
        config: {},
        status: 'active',
        test: null,
        running: true,
      }],
    };
    const actions = {
      state,
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
    const navigateTo = vi.fn();
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <WorkbenchV2Shell
          state={state}
          actions={actions}
          navigateTo={navigateTo}
        />
      </QueryClientProvider>,
    );

    fireEvent.click(screen.getByTestId('summary-rail-runtime-link'));

    expect(navigateTo).toHaveBeenCalledWith('/studio/runtime');
  });

  it('renders recent activation and audit history in the summary rail', () => {
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
          workspaceAuditHistory={[{
            id: 'audit-1',
            workspace_id: 'workspace-1',
            event_type: 'workspace_activation',
            result: 'partial_success',
            scope: 'devices:dev-A,dev-B',
            details: '{"message":"dev-B failed"}',
            occurred_at: '2026-05-29T10:12:00Z',
            created_at: '2026-05-29T10:12:01Z',
          }, {
            id: 'audit-2',
            workspace_id: 'workspace-1',
            event_type: 'database_config_saved',
            result: 'success',
            scope: 'database_connector:sqlite',
            details: '{}',
            occurred_at: '2026-05-29T10:10:00Z',
            created_at: '2026-05-29T10:10:01Z',
          }]}
        />
      </QueryClientProvider>,
    );

    const auditPanel = screen.getByTestId('summary-rail-audit-history');
    expect(auditPanel).toHaveTextContent('workspace_activation');
    expect(auditPanel).toHaveTextContent('partial_success');
    expect(auditPanel).toHaveTextContent('database_config_saved');
    expect(auditPanel).toHaveTextContent('2026-05-29T10:12:00Z');
  });
});
