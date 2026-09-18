import { describe, expect, it } from 'vitest';
import { isRowGroupScopeChange } from '../../../src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { DbConnector } from '../../../src/features/datalink/workbench-v2/state/types';

const bound: DbConnector = {
  ...INITIAL_STATE.db.connector,
  kind: 'postgres', connector_id: 'C1', identity_revision: 'R1', schema: 'public', table: 'sensor_readings',
};

describe('row group scope', () => {
  it('treats another saved connection with the same table as a new row-group scope', () => {
    expect(isRowGroupScopeChange(bound, {
      connector_id: 'C2', identity_revision: 'R9', kind: 'postgres', schema: 'public', table: 'sensor_readings',
    })).toBe(true);
  });

  it('keeps row groups when the same connection only changes identity revision or endpoint fields', () => {
    expect(isRowGroupScopeChange(bound, { connector_id: 'C1', identity_revision: 'R2' })).toBe(false);
    expect(isRowGroupScopeChange(bound, { host: 'db-b.internal' })).toBe(false);
  });
});
