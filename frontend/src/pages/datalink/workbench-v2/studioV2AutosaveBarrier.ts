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
