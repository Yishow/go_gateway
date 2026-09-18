import { act, cleanup, fireEvent, render, renderHook, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { useStep4Activation } from '../../../src/features/datalink/workbench-v2/steps/step4/useStep4Activation';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2ActivationResponse } from '../../../src/types/studioV2Activation';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/WorkspaceRecordingPlanSetupSection', () => ({
  WorkspaceRecordingPlanSetupSection: () => null,
}));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels', () => ({
  Step4SupportPanels: () => null,
}));

const baseState: WorkbenchV2State = {
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
        host: 'database.invalid',
        port: 5432,
        database: 'fixture_metrics',
        username: 'fixture_user',
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

const savedState: WorkbenchV2State = {
  ...baseState,
  devices: baseState.devices.map((device) => ({ ...device, save_state: 'saved' })),
  rules: baseState.rules.map((rule) => ({ ...rule, save_state: 'saved' })),
  mappings: { 'p-1': { ...baseState.mappings['p-1'], save_state: 'saved' } },
  db: {
    ...baseState.db,
    connector: { ...baseState.db.connector, save_state: 'saved' },
    targets: { 'p-1': { ...baseState.db.targets['p-1'], save_state: 'saved' } },
  },
};

function deferredResponse() {
  let resolve!: (response: StudioV2ActivationResponse) => void;
  const promise = new Promise<StudioV2ActivationResponse>((done) => { resolve = done; });
  return { promise, resolve };
}

afterEach(() => { cleanup(); vi.useRealTimers(); });

describe('Step 4 evidence boundaries', () => {
  it('does not infer success when all presentation timers elapse before a response', async () => {
    vi.useFakeTimers();
    const pending = deferredResponse();
    const activate = vi.fn().mockReturnValue(pending.promise);
    const { result, unmount } = renderHook(() => useStep4Activation(activate));
    let attempt!: Promise<void>;
    act(() => { attempt = result.current.start(); });
    await act(async () => { await vi.advanceTimersByTimeAsync(60_000); });
    expect(result.current.phase).toBe('activating');
    expect(result.current.response).toBeNull();
    expect(result.current.canContinue).toBe(false);
    expect(result.current.logs).toEqual([]);
    expect(activate).toHaveBeenCalledTimes(1);
    unmount();
    await act(async () => {
      pending.resolve({ workspace_id: 'workspace-1', results: [] });
      await attempt;
    });
  });

  it('reads server state on remount without replaying activation or accepting the old acknowledgement', async () => {
    const pending = deferredResponse();
    const activate = vi.fn().mockReturnValue(pending.promise);
    const recover = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-1', devices: [{ device_id: 'd-1', running: true }], operation_id: 'op-existing',
    });
    const first = renderHook(() => useStep4Activation(activate, recover, 'workspace-1'));
    await waitFor(() => expect(first.result.current.recovery?.operation_id).toBe('op-existing'));
    let attempt!: Promise<void>;
    act(() => { attempt = first.result.current.start(); });
    first.unmount();
    const second = renderHook(() => useStep4Activation(activate, recover, 'workspace-1'));
    await waitFor(() => expect(second.result.current.recovery?.operation_id).toBe('op-existing'));
    await act(async () => {
      pending.resolve({ workspace_id: 'workspace-1', results: [{ device_id: 'd-1', status: 'success' }] });
      await attempt;
    });
    expect(recover).toHaveBeenCalledTimes(2);
    expect(activate).toHaveBeenCalledTimes(1);
    expect(second.result.current.phase).toBe('idle');
    expect(second.result.current.response).toBeNull();
    expect(second.result.current.canContinue).toBe(false);
  });

  it('labels saved configuration without claiming activation or database delivery', () => {
    render(<Step4Database state={savedState} dispatch={vi.fn()} />);
    expect(screen.getByTestId('step4-configuration-saved')).toHaveTextContent('step4.configuration_saved');
    expect(screen.queryByTestId('activation-results-card')).not.toBeInTheDocument();
    expect(screen.queryByTestId('activation-delivery-status')).not.toBeInTheDocument();
  });

  it('keeps an earlier active device navigable without confirming an empty new attempt', async () => {
    const onCommit = vi.fn();
    const activate = vi.fn().mockResolvedValue({ workspace_id: 'workspace-1', results: [] });
    const state = { ...savedState, devices: savedState.devices.map((device) => ({ ...device, running: true })) };
    render(<Step4Database state={state} dispatch={vi.fn()} activateWorkspace={activate} onCommit={onCommit} />);
    fireEvent.click(screen.getByText('step4.activate_btn'));
    expect(await screen.findByTestId('activation-empty-message')).toHaveAttribute('data-outcome', 'unconfirmed');
    expect(screen.getByTestId('activation-delivery-status')).toHaveTextContent('step4.activation_delivery_unconfirmed');
    expect(screen.queryByTestId('activation-results-card')).not.toBeInTheDocument();
    const handoff = screen.getByText('step4.go_to_dashboard_btn');
    fireEvent.click(handoff);
    fireEvent.click(handoff);
    expect(onCommit).toHaveBeenCalledTimes(1);
  });

  it('keeps delivery unconfirmed after an acknowledgement and does not rewrite saved device status', async () => {
    const dispatch = vi.fn();
    const activate = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-1', results: [{ device_id: 'd-1', status: 'success' }],
    });
    render(<Step4Database state={savedState} dispatch={dispatch} activateWorkspace={activate} />);
    dispatch.mockClear();
    fireEvent.click(screen.getByText('step4.activate_btn'));
    expect(await screen.findByTestId('activation-results-card')).toHaveAttribute('data-outcome', 'confirmed');
    expect(screen.getByTestId('activation-delivery-status')).toHaveTextContent('step4.activation_delivery_unconfirmed');
    expect(dispatch).not.toHaveBeenCalledWith(expect.objectContaining({ type: 'updateDevice' }));
  });

  it('retains each rejection when all devices fail and offers no success handoff', async () => {
    const activate = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-1', results: [
        { device_id: 'd-1', status: 'failed' }, { device_id: 'd-2', status: 'failed' },
      ],
    });
    render(<Step4Database state={savedState} dispatch={vi.fn()} activateWorkspace={activate} onCommit={vi.fn()} />);
    fireEvent.click(screen.getByText('step4.activate_btn'));
    expect(await screen.findByTestId('activation-empty-message')).toHaveAttribute('data-outcome', 'failed');
    expect(screen.getByTestId('activation-neutral-result-d-1')).toHaveTextContent('step4.result_status.failed');
    expect(screen.getByTestId('activation-neutral-result-d-2')).toHaveTextContent('step4.result_status.failed');
    expect(screen.queryByTestId('activation-results-card')).not.toBeInTheDocument();
    expect(screen.queryByText('step4.go_to_dashboard_btn')).not.toBeInTheDocument();
  });
});
