import { describe, expect, it } from 'vitest';
import { metadataColumnsForScope } from '../../../src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { DbConnector } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2DatabaseMetadata } from '../../../src/types/studioV2DatabaseMetadata';

const connector: DbConnector = {
  ...INITIAL_STATE.db.connector, kind: 'sqlite', connector_id: 'db-1', identity_revision: 'identity-1',
  database: '/tmp/line-a.db', schema: 'main', table: 'sensor_values',
};
const metadata: StudioV2DatabaseMetadata = {
  workspace_id: 'ws-1', connector_id: 'db-1', connector_revision: 'identity-1', database: '/tmp/line-a.db',
  schema: 'main', table: 'sensor_values', inspection_status: 'exists',
  columns: [{ name: 'recorded_at', data_type: 'TEXT', nullable: false, primary_key: true }],
};

describe('Step 4 target metadata scope', () => {
  it('uses inspected columns only for an existing table in the current saved scope', () => {
    expect(metadataColumnsForScope(metadata, connector)).toEqual([
      { name: 'recorded_at', type: 'TEXT', nullable: false, primary_key: true },
    ]);
  });

  it('drops metadata from another connector, revision, database, schema or table', () => {
    for (const changed of [
      { connector_id: 'db-2' }, { identity_revision: 'identity-2' }, { database: '/tmp/other.db' },
      { schema: 'aux' }, { table: 'other_values' },
    ] satisfies Partial<DbConnector>[]) {
      expect(metadataColumnsForScope(metadata, { ...connector, ...changed })).toEqual([]);
    }
  });

  it('never supplies columns for missing, forbidden or failed inspections', () => {
    for (const status of ['missing', 'forbidden', 'failed'] as const) {
      expect(metadataColumnsForScope({ ...metadata, inspection_status: status }, connector)).toEqual([]);
    }
    expect(metadataColumnsForScope(undefined, connector)).toEqual([]);
  });
});
