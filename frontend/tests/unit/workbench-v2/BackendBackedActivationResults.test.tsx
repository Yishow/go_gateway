import * as React from 'react';
import { act, render, renderHook, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CommitProgress } from '../../../src/features/datalink/workbench-v2/steps/step4/CommitProgress';
import { CommitSuccessCard } from '../../../src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard';
import { ActivationNeutralSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/ActivationNeutralSummary';
import { useStep4Activation } from '../../../src/features/datalink/workbench-v2/steps/step4/useStep4Activation';
import { parseStudioV2ActivationResponse } from '../../../src/utils/safeJson';
import type { StudioV2ActivationResponse } from '../../../src/types/studioV2Activation';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

describe('Backend backed activation results', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('keeps a partial activation card amber and hands off only confirmed device ids once', () => {
    const onCommit = vi.fn();

    render(
      <CommitSuccessCard
        response={{
          workspace_id: 'workspace-1',
          results: [
            { device_id: '550e8400-e29b-41d4-a716-446655440000', status: 'success' },
            { device_id: '550e8400-e29b-41d4-a716-446655440001', status: 'failed' },
          ],
        }}
        canContinue
        onCommit={onCommit}
        onReset={() => {}}
      />,
    );

    const card = screen.getByTestId('activation-results-card');
    expect(card).toHaveAttribute('data-outcome', 'partial');
    expect(card).toHaveClass('border-amber-500/30');
    expect(card).toHaveTextContent('step4.result_status.success');
    expect(card).toHaveTextContent('step4.result_status.failed');

    const handoff = screen.getByRole('button', { name: 'step4.go_to_dashboard_btn' });
    handoff.click();
    handoff.click();
    expect(onCommit).toHaveBeenCalledTimes(1);
    expect(onCommit).toHaveBeenCalledWith(['550e8400-e29b-41d4-a716-446655440000']);
  });

  it('does not describe an empty response as an already-running successful workspace', () => {
    render(
      <CommitSuccessCard
        response={{ workspace_id: 'workspace-1', results: [] }}
        canContinue={false}
        onReset={() => {}}
      />,
    );

    expect(screen.queryByText('step4.already_active_info')).not.toBeInTheDocument();
  });

  it('keeps pending and skipped device results visible when no device succeeded', () => {
    render(
      <ActivationNeutralSummary
        response={{
          workspace_id: 'workspace-1',
          results: [
            { device_id: 'device-pending', status: 'pending' },
            { device_id: 'device-skipped', status: 'skipped' },
          ],
        }}
        canContinue={false}
        onReset={() => {}}
      />,
    );

    expect(screen.getByTestId('activation-neutral-results')).toBeInTheDocument();
    expect(screen.getByTestId('activation-empty-message')).toHaveAttribute('data-outcome', 'unconfirmed');
    expect(screen.getByTestId('activation-neutral-result-device-pending'))
      .toHaveTextContent('step4.result_status.pending');
    expect(screen.getByTestId('activation-neutral-result-device-skipped'))
      .toHaveTextContent('step4.result_status.skipped');
  });

  it('keeps a fresh definite rejection separate from an older recovery failure', () => {
    render(<ActivationNeutralSummary
      response={{ workspace_id: 'workspace-1', results: [], outcome: 'failed' }}
      recoveryUnavailable
      onReset={() => {}}
    />);
    expect(screen.getByTestId('activation-empty-message')).toHaveAttribute('data-outcome', 'failed');
    expect(screen.queryByTestId('activation-recovery-unavailable')).not.toBeInTheDocument();
  });

  it('labels backend-confirmed rows without inventing an HTTP 200 badge', () => {
    render(
      <CommitProgress
        logs={[{ label: 'Activate device', detail: '', status: 'success' }]}
        status="success"
      />,
    );

    const row = screen.getByTestId('commit-log-row-0');
    expect(row).toHaveTextContent('step4.progress_confirmed');
    expect(row).not.toHaveTextContent('200');
  });

  it('keeps activation acknowledgement local instead of updating draft device state', async () => {
    const activateWorkspace = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [{ device_id: 'device-1', status: 'success' }],
    });
    const { result } = renderHook(() => useStep4Activation(activateWorkspace));

    await act(async () => {
      await result.current.start();
    });

    expect(result.current.canContinue).toBe(true);
  });

  it('processes an acknowledgement after React StrictMode remount setup', async () => {
    const activateWorkspace = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-1',
      results: [{ device_id: 'device-1', status: 'success' }],
    });
    const wrapper = ({ children }: { children: React.ReactNode }) => (
      <React.StrictMode>{children}</React.StrictMode>
    );
    const { result } = renderHook(() => useStep4Activation(activateWorkspace), { wrapper });

    await act(async () => {
      await result.current.start();
    });

    expect(result.current.response?.results[0]?.status).toBe('success');
    expect(result.current.canContinue).toBe(true);
  });

  it('blocks a second activation click while the first backend attempt is pending', async () => {
    let resolveActivation!: (value: StudioV2ActivationResponse) => void;
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>(() => new Promise((resolve) => {
      resolveActivation = resolve;
    }));
    const { result } = renderHook(() => useStep4Activation(activateWorkspace));

    await act(async () => {
      void result.current.start();
      void result.current.start();
      await Promise.resolve();
    });

    expect(activateWorkspace).toHaveBeenCalledTimes(1);

    await act(async () => {
      resolveActivation({ workspace_id: 'workspace-1', results: [] });
    });
  });

  it('keeps the in-flight guard when reset is pressed before the backend settles', async () => {
    let resolveActivation!: (value: StudioV2ActivationResponse) => void;
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>(() => new Promise((resolve) => {
      resolveActivation = resolve;
    }));
    const { result } = renderHook(() => useStep4Activation(activateWorkspace));

    await act(async () => {
      void result.current.start();
      await Promise.resolve();
    });
    act(() => result.current.reset());
    await act(async () => {
      void result.current.start();
      await Promise.resolve();
    });
    expect(activateWorkspace).toHaveBeenCalledTimes(1);

    await act(async () => {
      resolveActivation({ workspace_id: 'workspace-1', results: [] });
      await Promise.resolve();
    });
    await act(async () => {
      void result.current.start();
      await Promise.resolve();
    });
    expect(activateWorkspace).toHaveBeenCalledTimes(2);
    await act(async () => {
      resolveActivation({ workspace_id: 'workspace-1', results: [] });
    });
  });

  it('keeps an untyped transport rejection unconfirmed and without retry evidence', async () => {
    const activateWorkspace = vi.fn().mockRejectedValue(new Error('network details must stay hidden'));
    const { result } = renderHook(() => useStep4Activation(activateWorkspace));

    await act(async () => {
      await result.current.start();
    });

    const response = result.current.response as unknown as {
      outcome?: string;
      retryable?: boolean;
    };
    expect(response.outcome).toBe('unconfirmed');
    expect(response.retryable).toBe(false);
    expect(result.current.logs[0]?.status).toBe('pending');
  });

  it('keeps an operation id from a malformed acknowledgement while remaining unconfirmed', async () => {
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>().mockResolvedValue({
      workspace_id: 'workspace-1',
      operation_id: 'activation-op-malformed',
      results: 'not-an-array' as unknown as StudioV2ActivationResponse['results'],
    });
    const { result } = renderHook(() => useStep4Activation(activateWorkspace));

    await act(async () => {
      await result.current.start();
    });

    expect(result.current.response?.outcome).toBe('unconfirmed');
    expect(result.current.response?.operation_id).toBe('activation-op-malformed');
  });

  it('retains a bounded operation id supplied by the backend response', () => {
    const response = parseStudioV2ActivationResponse({
      workspace_id: 'workspace-1',
      operation_id: 'activation-op-1',
      results: [{ device_id: 'device-1', status: 'success' }],
    });

    expect((response as unknown as { operation_id?: string })?.operation_id).toBe('activation-op-1');
  });

  it('reads existing runtime recovery on mount without submitting activation', async () => {
    const activateWorkspace = vi.fn();
    const recoverActivationStatus = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-1',
      devices: [{ device_id: 'device-1', running: true }],
      operation_id: 'recovery-op-1',
    });
    const { result } = renderHook(() => useStep4Activation(
      activateWorkspace,
      recoverActivationStatus,
      'workspace-1',
    ));

    await waitFor(() => expect(result.current.recovery?.devices[0]?.running).toBe(true));

    expect(recoverActivationStatus).toHaveBeenCalledTimes(1);
    expect(activateWorkspace).not.toHaveBeenCalled();
    expect(result.current.recovery?.operation_id).toBe('recovery-op-1');
  });

  it('does not treat a foreign recovery workspace as current activation evidence', async () => {
    const recoverActivationStatus = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-foreign',
      devices: [{ device_id: 'device-1', running: true }],
    });
    const { result } = renderHook(() => useStep4Activation(
      undefined,
      recoverActivationStatus,
      'workspace-1',
    ));

    await waitFor(() => expect(result.current.recoveryUnavailable).toBe(true));

    expect(result.current.recovery).toBeNull();
  });

  it('drops an old workspace acknowledgement after switching workspace while pending', async () => {
    let resolveActivation!: (value: StudioV2ActivationResponse) => void;
    const activateWorkspace = vi.fn<() => Promise<StudioV2ActivationResponse>>(() => new Promise((resolve) => {
      resolveActivation = resolve;
    }));
    const { result, rerender } = renderHook(
      ({ workspaceId }: { workspaceId: string }) => useStep4Activation(activateWorkspace, undefined, workspaceId),
      { initialProps: { workspaceId: 'workspace-A' } },
    );

    await act(async () => {
      void result.current.start();
      await Promise.resolve();
    });
    await act(async () => {
      rerender({ workspaceId: 'workspace-B' });
      await Promise.resolve();
    });

    await act(async () => {
      resolveActivation({
        workspace_id: 'workspace-A',
        results: [{ device_id: 'device-A', status: 'success' }],
      });
      await Promise.resolve();
    });

    expect(result.current.response).toBeNull();
    expect(result.current.canContinue).toBe(false);
  });

  it('clears a completed acknowledgement when another workspace becomes current', async () => {
    const activateWorkspace = vi.fn().mockResolvedValue({
      workspace_id: 'workspace-A', results: [{ device_id: 'device-A', status: 'success' }],
    });
    const { result, rerender } = renderHook(
      ({ workspaceId }) => useStep4Activation(activateWorkspace, undefined, workspaceId),
      { initialProps: { workspaceId: 'workspace-A' } },
    );
    await act(async () => { await result.current.start(); });
    expect(result.current.canContinue).toBe(true);
    rerender({ workspaceId: 'workspace-B' });
    expect(result.current.phase).toBe('idle');
    expect(result.current.response).toBeNull();
    expect(result.current.canContinue).toBe(false);
    expect(activateWorkspace).toHaveBeenCalledTimes(1);
  });
});
