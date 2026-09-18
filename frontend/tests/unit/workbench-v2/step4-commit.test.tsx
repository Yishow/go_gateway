import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { WorkbenchV2Shell, type WorkbenchV2ShellProps } from '../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { INITIAL_STATE, workbenchV2Reducer } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import type { WorkbenchV2Action } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { StudioV2ActivationResponse } from '../../../src/types/studioV2Activation';
import type { StudioV2WorkspaceReadinessSummary } from '../../../src/types/studioV2WorkspaceReadiness';

// Mock react-i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: { interval?: number }) => {
      if (options && options.interval !== undefined) {
        return `${key}_interval_${options.interval}`;
      }
      return key;
    }
  })
}));

describe('Step 4 first activation flow integration', () => {
  function LocationProbe() {
    const location = useLocation();
    return <div data-testid="step4-shell-location">{`${location.pathname}${location.search}`}</div>;
  }

  function ShellRouterHarness({
    state,
    actions,
    activateWorkspace,
  }: {
    state: WorkbenchV2State;
    actions: {
      state: WorkbenchV2State;
      setView: ReturnType<typeof vi.fn>;
      setCurrent: ReturnType<typeof vi.fn>;
      completeStep: ReturnType<typeof vi.fn>;
      toggleSidebar: ReturnType<typeof vi.fn>;
      toggleSummaryRail: ReturnType<typeof vi.fn>;
      setSidebarCollapsed: ReturnType<typeof vi.fn>;
      setShowSummaryRail: ReturnType<typeof vi.fn>;
      resetFlow: ReturnType<typeof vi.fn>;
      selectRule: ReturnType<typeof vi.fn>;
      dispatch: ReturnType<typeof vi.fn>;
    };
    activateWorkspace?: WorkbenchV2ShellProps['activateWorkspace'];
  }) {
    const navigate = useNavigate();

    return (
      <>
        <WorkbenchV2Shell
          state={state}
          actions={actions}
          navigateTo={navigate}
          activateWorkspace={activateWorkspace}
        />
        <LocationProbe />
      </>
    );
  }

  const mockState: WorkbenchV2State = {
    view: 'flow',
    current: 4,
    completed: new Set([1, 2, 3]),
    sidebarCollapsed: false,
    showSummaryRail: true,
    devices: [
      { id: 'd-1', name: 'PLC-1', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null }
    ],
    rules: [
      { id: 'r-1', device_id: 'd-1', name: 'Holding Registers', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 't_', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null }
    ],
    selectedRuleId: null,
    points: [
      { id: 'p-1', device_id: 'd-1', rule_id: 'r-1', rule_name: 'Holding Registers', name: 't_1', address: '40001', data_type: 'int16', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0 }
    ],
    mappings: {
      'p-1': { point_id: 'p-1', tag_key: 'line1.t_1', display_name: 'T1', unit: 'C', target_type: 'float64', scale: 1, offset: 0, enabled: true }
    },
    db: {
      connector: {
        kind: 'postgres',
        name: 'PostgreSQL Connector',
        host: 'tsdb.internal',
        port: 5432,
        database: 'gateway_metrics',
        username: 'postgres',
        schema: 'public',
        table: 'sensor_readings',
        write_mode: 'insert',
        write_interval_seconds: 5,
        timestamp_column: 'ts',
        status: 'unknown'
      },
      targets: {
        'p-1': { tag_id: 'tag.line1.t_1', column_name: 'temp_in_c', enabled: true }
      }
    },
    settings: INITIAL_STATE.settings,
    committed: false,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('activates all eligible devices and shows per-device partial success results', async () => {
    let state = { ...mockState };
    const dispatch = vi.fn((action: WorkbenchV2Action) => {
      state = workbenchV2Reducer(state, action);
    });
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [
        { device_id: 'd-1', status: 'success', message: 'activated' },
        { device_id: 'd-2', status: 'failed', message: 'activation timeout' },
      ],
    });

    const { rerender } = render(
      <Step4Database
        state={{
          ...state,
          devices: [
            ...state.devices,
            { id: 'd-2', name: 'PLC-2', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null, running: false },
          ],
        }}
        dispatch={dispatch}
        onCommit={() => {}}
        activateWorkspace={activateWorkspace}
      />
    );

    fireEvent.click(screen.getByText('step4.activate_btn'));

    await waitFor(() => {
      expect(activateWorkspace).toHaveBeenCalledTimes(1);
    });

    rerender(
      <Step4Database
        state={state}
        dispatch={dispatch}
        onCommit={() => {}}
        activateWorkspace={activateWorkspace}
      />
    );

    await waitFor(() => {
      expect(screen.getByTestId('activation-result-d-1')).toHaveTextContent('success');
      expect(screen.getByTestId('activation-result-d-2')).toHaveTextContent('failed');
      expect(screen.getByText('step4.go_to_dashboard_btn')).toBeInTheDocument();
    });
  });

  it('shows an actionable empty result when no devices are eligible', async () => {
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [],
      message: '目前沒有可啟動設備',
    });

    render(
      <Step4Database
        state={mockState}
        dispatch={() => { }}
        activateWorkspace={activateWorkspace}
      />
    );

    fireEvent.click(screen.getByText('step4.activate_btn'));

    await waitFor(() => {
      expect(screen.getByTestId('activation-empty-message')).toHaveTextContent('step4.activation_unconfirmed');
      expect(screen.getByTestId('activation-delivery-status')).toHaveTextContent('step4.activation_delivery_unconfirmed');
      expect(screen.getByText('step4.reset_activation_btn')).toBeInTheDocument();
    });
  });

  it('renders a typed activation barrier error in CommitProgress without raw details', async () => {
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>()
      .mockRejectedValue({
        code: 'modbus_share_revision_conflict',
        action: 'hydrate again with server state',
        request_id: 'req-activation-7',
        retryable: true,
        message: 'raw backend revision details',
      });

    render(
      <Step4Database
        state={mockState}
        dispatch={() => { }}
        activateWorkspace={activateWorkspace}
      />,
    );

    fireEvent.click(screen.getByText('step4.activate_btn'));

    const failedRow = await screen.findByTestId('commit-log-row-0');
    expect(failedRow).toHaveAttribute('data-status', 'failed');
    expect(failedRow).toHaveTextContent('errors.modbus_share_revision_conflict');
    expect(failedRow).toHaveTextContent('req-activation-7');
    expect(failedRow).not.toHaveTextContent('raw backend revision details');
    expect(failedRow).not.toHaveTextContent('hydrate again with server state');
  });

  it('keeps activation acknowledgement out of the draft device reducer path', async () => {
    const dispatch = vi.fn();
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [{ device_id: 'd-1', status: 'success' }],
    });

    render(
      <Step4Database
        state={mockState}
        dispatch={dispatch}
        onCommit={() => {}}
        activateWorkspace={activateWorkspace}
      />,
    );
    dispatch.mockClear();

    fireEvent.click(screen.getByText('step4.activate_btn'));
    await waitFor(() => expect(screen.getByTestId('activation-results-card')).toBeInTheDocument());

    expect(dispatch).not.toHaveBeenCalledWith(expect.objectContaining({ type: 'updateDevice' }));
  });

  it('wires the localized failed-row retry to the Step 4 activation callback', async () => {
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>()
      .mockResolvedValueOnce({
        workspace_id: 'workspace-1',
        results: [{ device_id: 'd-1', status: 'failed', message: 'activation failed', retryable: true }],
      })
      .mockResolvedValueOnce({
        workspace_id: 'workspace-1',
        results: [{ device_id: 'd-1', status: 'success', message: 'activated' }],
      });

    render(
      <Step4Database
        state={mockState}
        dispatch={() => { }}
        activateWorkspace={activateWorkspace}
      />,
    );

    fireEvent.click(screen.getByText('step4.activate_btn'));
    const retry = await screen.findByRole('button', { name: 'step4.retry_action' });
    fireEvent.click(retry);

    await waitFor(() => expect(activateWorkspace).toHaveBeenCalledTimes(2));
    expect(screen.queryByRole('button', { name: 'step4.retry_action' })).not.toBeInTheDocument();
  });

  it('allows runtime navigation after partial activation failure when at least one device succeeds', async () => {
    const onCommit = vi.fn();
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [
        { device_id: 'd-1', status: 'success', message: 'activated' },
        { device_id: 'd-2', status: 'failed', message: 'activation timeout' },
      ],
    });
    let state: WorkbenchV2State = {
      ...mockState,
      devices: [
        { id: 'd-1', name: 'PLC-1', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null, running: false },
        { id: 'd-2', name: 'PLC-2', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null, running: false },
      ],
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
      dispatch: vi.fn((action: WorkbenchV2Action) => {
        state = workbenchV2Reducer(state, action);
        actions.state = state;
      }),
    };

    render(
      <MemoryRouter initialEntries={['/studio/v2']}>
        <ShellRouterHarness
          state={state}
          actions={actions}
          activateWorkspace={activateWorkspace}
        />
      </MemoryRouter>,
    );

    fireEvent.click(screen.getByText('step4.activate_btn'));

    await waitFor(() => {
      expect(screen.getByText('step4.go_to_dashboard_btn')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('step4.go_to_dashboard_btn'));

    expect(screen.getByTestId('step4-shell-location')).toHaveTextContent('/studio/runtime');
    expect(onCommit).not.toHaveBeenCalled();
  });

  it('hands off once to the confirmed backend device route', async () => {
    const confirmedDeviceId = '550e8400-e29b-41d4-a716-446655440010';
    const navigateTo = vi.fn();
    const state: WorkbenchV2State = {
      ...mockState,
      devices: [{ ...mockState.devices[0], id: confirmedDeviceId }],
      rules: [{ ...mockState.rules[0], device_id: confirmedDeviceId }],
      points: [{ ...mockState.points[0], device_id: confirmedDeviceId }],
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
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [{ device_id: confirmedDeviceId, status: 'success' }],
    });

    render(
      <MemoryRouter initialEntries={['/studio/v2']}>
        <WorkbenchV2Shell
          state={state}
          actions={actions}
          navigateTo={navigateTo}
          activateWorkspace={activateWorkspace}
        />
      </MemoryRouter>,
    );

    fireEvent.click(screen.getByText('step4.activate_btn'));
    const handoff = await screen.findByText('step4.go_to_dashboard_btn');
    fireEvent.click(handoff);
    fireEvent.click(handoff);

    expect(navigateTo).toHaveBeenCalledTimes(1);
    expect(navigateTo).toHaveBeenCalledWith(`/studio/runtime?device_id=${confirmedDeviceId}`);
  });

  it('shows readiness blockers and disables activation before start', () => {
    const readinessSummary: StudioV2WorkspaceReadinessSummary = {
      ready: false,
      blocking_count: 1,
      warning_count: 0,
      issues: [
        {
          code: 'device-probe-required',
          severity: 'blocking',
          step: 'Step 1',
          scope: 'd-1',
          message: 'probe diagnostics failed',
        },
      ],
    };
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>();
    const onNavigateStep = vi.fn();

    render(
      <Step4Database
        state={mockState}
        dispatch={() => { }}
        activateWorkspace={activateWorkspace}
        workspaceReadiness={readinessSummary}
        onNavigateStep={onNavigateStep}
      />,
    );

    const activateButton = screen.getByText('step4.activate_btn').closest('button');
    expect(screen.getByTestId('step4-readiness-panel')).toHaveTextContent('device-probe-required');
    expect(screen.getByTestId('step4-blocker-resolution')).toHaveTextContent('step4.blocker_resolution_title');
    fireEvent.click(screen.getByTestId('step4-blocker-action-device-probe-required'));
    expect(onNavigateStep).toHaveBeenCalledWith(1);
    expect(activateButton).toBeDisabled();
    fireEvent.click(screen.getByText('step4.activate_btn'));
    expect(activateWorkspace).not.toHaveBeenCalled();
  });

  it('shows readiness warnings but keeps activation enabled', () => {
    const readinessSummary: StudioV2WorkspaceReadinessSummary = {
      ready: true,
      blocking_count: 0,
      warning_count: 1,
      issues: [
        {
          code: 'database-connector-unreachable',
          severity: 'warning',
          step: 'Step 4',
          scope: 'db-main',
          message: 'database connector requires attention',
        },
      ],
    };

    render(
      <Step4Database
        state={mockState}
        dispatch={() => { }}
        workspaceReadiness={readinessSummary}
      />,
    );

    const activateButton = screen.getByText('step4.activate_btn').closest('button');
    expect(screen.getByTestId('step4-readiness-panel')).toHaveTextContent('database-connector-unreachable');
    expect(activateButton).not.toBeDisabled();
  });
});
