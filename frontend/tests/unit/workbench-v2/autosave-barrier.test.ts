import { describe, expect, it } from 'vitest';
import { getStudioV2AutosaveBarrier } from '../../../src/pages/datalink/workbench-v2/studioV2AutosaveBarrier';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

describe('Studio V2 autosave barrier', () => {
  it('counts pending durable saves and emits a typed error marker without raw diagnostics', () => {
    const state = {
      ...INITIAL_STATE,
      devices: [{ ...INITIAL_STATE.devices[0], persisted: true, save_state: 'saving' as const }],
      rules: [{ ...INITIAL_STATE.rules[0], persisted: true, save_state: 'save-error' as const, save_error: 'secret backend details' }],
      mappings: {
        'point-1': {
          ...INITIAL_STATE.mappings['point-1'],
          persisted: false,
          save_state: 'draft-invalid' as const,
        },
      },
    };

    const barrier = getStudioV2AutosaveBarrier(state);

    expect(barrier.pending_saves).toBe(3);
    expect(barrier.save_error).toBe('autosave_failed');
    expect(JSON.stringify(barrier)).not.toContain('secret backend details');
  });

  it('returns an empty barrier when every durable save is settled', () => {
    const barrier = getStudioV2AutosaveBarrier(INITIAL_STATE);

    expect(barrier).toEqual({ pending_saves: 0 });
  });
});
