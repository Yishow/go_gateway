import * as React from 'react';
import type { WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';

export interface StudioV2AutosaveBarrier {
  pending_saves: number;
  save_error?: string;
}

export function isUnpersistedNonIdle(saveState?: string, persisted?: boolean): boolean {
  return !persisted && saveState !== 'idle';
}

export function hasPendingSave(saveState?: string, persisted?: boolean): boolean {
  return isUnpersistedNonIdle(saveState, persisted) || saveState === 'saving' || saveState === 'draft-invalid' || saveState === 'save-error';
}

function hasSaveError(saveState?: string): boolean {
  return saveState === 'save-error';
}

function hasAutosaveState(state: WorkbenchV2State): boolean {
  return Boolean(state?.devices && state.rules && state.mappings && state.db && state.settings);
}

function hasAutosaveInFlight(state: WorkbenchV2State): boolean {
  return [
    ...state.devices.map((item) => item.save_state),
    ...state.rules.map((item) => item.save_state),
    ...Object.values(state.mappings).map((item) => item.save_state),
    state.db.connector.save_state,
    ...Object.values(state.db.targets).map((item) => item.save_state),
    state.settings.modbus_share.save_state,
  ].some((saveState) => saveState === 'saving');
}

/** Returns the durable-save barrier without exposing backend or exception text. */
export function getStudioV2AutosaveBarrier(state: WorkbenchV2State): StudioV2AutosaveBarrier {
  const saveStates = [
    ...state.devices.map((item) => [item.save_state, item.persisted] as const),
    ...state.rules.map((item) => [item.save_state, item.persisted] as const),
    ...Object.values(state.mappings).map((item) => [item.save_state, item.persisted] as const),
    [state.db.connector.save_state, state.db.connector.persisted] as const,
    ...Object.values(state.db.targets).map((item) => [item.save_state, item.persisted] as const),
    [state.settings.modbus_share.save_state, true] as const,
  ];
  const pendingSaves = saveStates.filter(([saveState, persisted]) => hasPendingSave(saveState, persisted)).length;
  const hasError = saveStates.some(([saveState]) => hasSaveError(saveState)) || Boolean(state.settings.modbus_share.save_error);

  return {
    pending_saves: pendingSaves,
    ...(hasError ? { save_error: 'autosave_failed' } : {}),
  };
}

interface StudioV2AutosaveSettlement {
  state: WorkbenchV2State;
  barrier: StudioV2AutosaveBarrier;
}

interface StudioV2AutosaveSettlementWaiter {
  resolve: (settlement: StudioV2AutosaveSettlement) => void;
  timer: ReturnType<typeof setTimeout>;
}

/** Upper bound for one activation wait; a stuck save must not hold activation open forever. */
export const AUTOSAVE_SETTLEMENT_TIMEOUT_MS = 10_000;

function settlementOf(state: WorkbenchV2State): StudioV2AutosaveSettlement {
  return { state, barrier: getStudioV2AutosaveBarrier(state) };
}

/**
 * Waits only for active saves; invalid drafts and save errors remain activation blockers.
 * The wait is bounded: once the bound elapses it resolves with the current barrier snapshot,
 * which still reports the pending saves, so activation fails as retryable instead of hanging.
 */
export function useStudioV2AutosaveSettlement(
  state: WorkbenchV2State,
  timeoutMs: number = AUTOSAVE_SETTLEMENT_TIMEOUT_MS,
) {
  const stateRef = React.useRef(state);
  const waitersRef = React.useRef<StudioV2AutosaveSettlementWaiter[]>([]);
  stateRef.current = state;

  React.useEffect(() => {
    if (!hasAutosaveState(state)) {
      return;
    }
    if (hasAutosaveInFlight(state)) {
      return;
    }
    const waiters = waitersRef.current.splice(0);
    const settlement = settlementOf(state);
    waiters.forEach((waiter) => {
      clearTimeout(waiter.timer);
      waiter.resolve(settlement);
    });
  }, [state]);

  React.useEffect(() => () => {
    // 卸載時清除計時器；此時畫面已不存在，不再對已消失的呼叫端回報 barrier 結果。
    waitersRef.current.splice(0).forEach((waiter) => clearTimeout(waiter.timer));
  }, []);

  return React.useCallback((): Promise<StudioV2AutosaveSettlement> => {
    const currentState = stateRef.current;
    // 與 effect 使用同一個守衛：狀態尚未成形時不得對其取 barrier，
    // 否則啟用會丟出 TypeError 而不是可辨識的 barrier 失敗。
    if (!hasAutosaveState(currentState)) {
      return Promise.resolve({
        state: currentState,
        barrier: { pending_saves: 1, save_error: 'autosave_failed' },
      });
    }
    if (!hasAutosaveInFlight(currentState)) {
      return Promise.resolve(settlementOf(currentState));
    }
    return new Promise((resolve) => {
      const timer = setTimeout(() => {
        waitersRef.current = waitersRef.current.filter((waiter) => waiter.timer !== timer);
        resolve(settlementOf(stateRef.current));
      }, timeoutMs);
      waitersRef.current.push({ resolve, timer });
    });
  }, [timeoutMs]);
}
