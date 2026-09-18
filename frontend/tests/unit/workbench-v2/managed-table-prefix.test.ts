import { describe, expect, it } from 'vitest';
import { MANAGED_TABLE_PREFIX } from '../../../src/features/datalink/workbench-v2/steps/step4/useRecordingPlanActions';

/**
 * Locks the client-supplied managed table prefix. The backend
 * (internal/datalink/recordingplan, defaultManagedTablePrefix) accepts and
 * validates this prefix; changing it here must be a conscious, coordinated
 * change on both sides.
 */
describe('MANAGED_TABLE_PREFIX', () => {
  it('stays locked to the backend defaultManagedTablePrefix value', () => {
    expect(MANAGED_TABLE_PREFIX).toBe('gw_record_');
  });
});
