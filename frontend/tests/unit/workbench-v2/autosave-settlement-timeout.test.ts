import { act, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import {
  AUTOSAVE_SETTLEMENT_TIMEOUT_MS,
  useStudioV2AutosaveSettlement,
} from '../../../src/pages/datalink/workbench-v2/studioV2AutosaveBarrier';
import { useStep4Activation } from '../../../src/features/datalink/workbench-v2/steps/step4/useStep4Activation';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';

const SAVING_STATE: WorkbenchV2State = {
  ...INITIAL_STATE,
  devices: [{ ...INITIAL_STATE.devices[0], persisted: true, save_state: 'saving' }],
};

describe('Activation autosave settlement is time-bounded', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('resolves immediately when no autosave is in flight', async () => {
    const { result } = renderHook(() => useStudioV2AutosaveSettlement(INITIAL_STATE));

    const settlement = await result.current();

    expect(settlement.barrier.pending_saves).toBe(0);
    expect(settlement.state).toBe(INITIAL_STATE);
  });

  it('resolves with the current barrier once the bound elapses', async () => {
    const { result } = renderHook(() => useStudioV2AutosaveSettlement(SAVING_STATE, 50));

    let settled = false;
    const pending = result.current().then((settlement) => {
      settled = true;
      return settlement;
    });

    await act(async () => { await vi.advanceTimersByTimeAsync(49); });
    expect(settled).toBe(false);

    await act(async () => { await vi.advanceTimersByTimeAsync(1); });
    const settlement = await pending;

    expect(settlement.barrier.pending_saves).toBeGreaterThanOrEqual(1);
  });

  it('resolves as soon as the saves settle, without waiting for the bound', async () => {
    const { result, rerender } = renderHook(
      ({ state }: { state: WorkbenchV2State }) => useStudioV2AutosaveSettlement(state, 50),
      { initialProps: { state: SAVING_STATE } },
    );

    const pending = result.current();
    await act(async () => { rerender({ state: INITIAL_STATE }); });

    const settlement = await pending;

    expect(settlement.barrier.pending_saves).toBe(0);
    expect(vi.getTimerCount()).toBe(0);
  });

  it('leaves no pending timer after a timed-out wait or an unmount', async () => {
    const { result, unmount } = renderHook(() => useStudioV2AutosaveSettlement(SAVING_STATE, 50));

    const pending = result.current();
    await act(async () => { await vi.advanceTimersByTimeAsync(50); });
    await pending;
    expect(vi.getTimerCount()).toBe(0);

    void result.current();
    expect(vi.getTimerCount()).toBe(1);
    unmount();
    expect(vi.getTimerCount()).toBe(0);
  });

  it('fails closed when the workspace state has not formed yet', async () => {
    const { result } = renderHook(() => useStudioV2AutosaveSettlement({} as WorkbenchV2State, 50));

    const settlement = await result.current();

    expect(settlement.barrier.pending_saves).toBeGreaterThanOrEqual(1);
    expect(settlement.barrier.save_error).toBe('autosave_failed');
  });

  it('exposes a named default bound', () => {
    expect(AUTOSAVE_SETTLEMENT_TIMEOUT_MS).toBe(10_000);
  });
});

describe('Activation leaves the activating phase when the barrier reports save-incomplete', () => {
  it('reaches a terminal phase carrying the retryable save-incomplete code', async () => {
    class BarrierError extends Error {
      code = 'modbus_share_save_incomplete';
      retryable = true;
    }
    const activateWorkspace = vi.fn().mockRejectedValue(new BarrierError('barrier'));
    const { result } = renderHook(() => useStep4Activation(activateWorkspace, vi.fn()));

    await act(async () => { await result.current.start(); });

    expect(result.current.phase).toBe('done');
    expect(result.current.phase).not.toBe('activating');
    expect(result.current.response?.code).toBe('modbus_share_save_incomplete');
    expect(result.current.response?.retryable).toBe(true);
    expect(result.current.logs[0]?.code).toBe('modbus_share_save_incomplete');
  });
});
