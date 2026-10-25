import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { CommitProgress } from '../../../src/features/datalink/workbench-v2/steps/step4/CommitProgress';
import { useStep4Activation } from '../../../src/features/datalink/workbench-v2/steps/step4/useStep4Activation';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
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

describe('Step 4 overview and explicit legacy activation', () => {
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

  it('keeps the overview free of any workspace-wide activation, even with several unstarted devices', () => {
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>();
    render(
      <Step4Database
        state={{
          ...mockState,
          devices: [
            ...mockState.devices,
            { id: 'd-2', name: 'PLC-2', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null, running: false },
          ],
        }}
        dispatch={() => { }}
        activateWorkspace={activateWorkspace}
      />
    );

    expect(screen.getByTestId('step4-group-summary-note')).toBeInTheDocument();
    expect(screen.queryByText('step4.activate_btn')).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'step4.retry_action' })).not.toBeInTheDocument();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });

  // The activation hook and its progress rows stay compatible for an explicit legacy caller of the global API.
  describe('explicit legacy activation pieces', () => {
    function LegacyCaller({ activate }: { activate: () => Promise<StudioV2ActivationResponse> }) {
      const activation = useStep4Activation(activate);
      return (
        <>
          <button type="button" data-testid="legacy-activate" onClick={() => void activation.start()}>go</button>
          <CommitProgress logs={activation.logs} status="failed" onRetry={activation.start} />
        </>
      );
    }

    it('renders a typed activation barrier error in CommitProgress without raw details', async () => {
      const activate = vi.fn<() => Promise<StudioV2ActivationResponse>>()
        .mockRejectedValue({
          code: 'modbus_share_revision_conflict',
          action: 'hydrate again with server state',
          request_id: 'req-activation-7',
          retryable: true,
          message: 'raw backend revision details',
        });
      render(<LegacyCaller activate={activate} />);

      fireEvent.click(screen.getByTestId('legacy-activate'));

      const failedRow = await screen.findByTestId('commit-log-row-0');
      expect(failedRow).toHaveAttribute('data-status', 'failed');
      expect(failedRow).toHaveTextContent('errors.modbus_share_revision_conflict');
      expect(failedRow).toHaveTextContent('req-activation-7');
      expect(failedRow).not.toHaveTextContent('raw backend revision details');
      expect(failedRow).not.toHaveTextContent('hydrate again with server state');
    });

    it('wires the localized failed-row retry to the explicit activation callback', async () => {
      const activate = vi.fn<() => Promise<StudioV2ActivationResponse>>()
        .mockResolvedValueOnce({
          workspace_id: 'workspace-1',
          results: [{ device_id: 'd-1', status: 'failed', message: 'activation failed', retryable: true }],
        })
        .mockResolvedValueOnce({
          workspace_id: 'workspace-1',
          results: [{ device_id: 'd-1', status: 'success', message: 'activated' }],
        });
      render(<LegacyCaller activate={activate} />);

      fireEvent.click(screen.getByTestId('legacy-activate'));
      const retry = await screen.findByRole('button', { name: 'step4.retry_action' });
      fireEvent.click(retry);

      await waitFor(() => expect(activate).toHaveBeenCalledTimes(2));
      expect(screen.queryByRole('button', { name: 'step4.retry_action' })).not.toBeInTheDocument();
    });
  });

  it('shows readiness blockers as status with a way to fix them, and no activation action', () => {
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

    expect(screen.getByTestId('step4-readiness-panel')).toHaveTextContent('device-probe-required');
    expect(screen.getByTestId('step4-blocker-resolution')).toHaveTextContent('step4.blocker_resolution_title');
    fireEvent.click(screen.getByTestId('step4-blocker-action-device-probe-required'));
    expect(onNavigateStep).toHaveBeenCalledWith(1);
    expect(screen.queryByText('step4.activate_btn')).not.toBeInTheDocument();
    expect(activateWorkspace).not.toHaveBeenCalled();
  });

  it('shows readiness warnings as status and no activation action', () => {
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

    expect(screen.getByTestId('step4-readiness-panel')).toHaveTextContent('database-connector-unreachable');
    expect(screen.queryByTestId('step4-blocker-resolution')).not.toBeInTheDocument();
    expect(screen.queryByText('step4.activate_btn')).not.toBeInTheDocument();
  });
});
