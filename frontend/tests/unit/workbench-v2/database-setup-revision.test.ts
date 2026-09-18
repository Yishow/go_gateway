import { describe, expect, it } from 'vitest';
import {
  hydrateStudioV2DatabaseConnector,
  toStudioV2DatabaseConfigRequest,
  toStudioV2DatabaseTargetRequest,
} from '../../../src/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import type { DbConnector, DbTarget } from '../../../src/features/datalink/workbench-v2/state/types';
import type { StudioV2WorkspaceDatabaseConfigRecord } from '../../../src/types/datalink';

const record: StudioV2WorkspaceDatabaseConfigRecord = {
  id: 'db-1', identity_revision: 'identity-1', setup_revision: 'setup-7', workspace_id: 'ws-1',
  kind: 'sqlite', name: 'Line A', host: '', port: 0, database: '/tmp/line-a.db', username: '',
  schema: 'main', table: 'sensor_values', write_mode: 'insert', write_interval_seconds: 5,
  timestamp_column: 'ts', status: 'ready', row_groups: [], created_at: '2026-09-16T00:00:00Z',
  updated_at: '2026-09-16T00:00:00Z',
};

describe('AtomicSetupSave setup revision', () => {
  it('keeps the saved setup revision on the hydrated connector', () => {
    expect(hydrateStudioV2DatabaseConnector(record).setup_revision).toBe('setup-7');
  });

  it('sends the expected setup revision with configuration saves', () => {
    const connector: DbConnector = { ...hydrateStudioV2DatabaseConnector(record), setup_revision: 'setup-9' };
    expect(toStudioV2DatabaseConfigRequest(connector).expected_setup_revision).toBe('setup-9');
  });

  it('sends the expected setup revision with target saves', () => {
    const target: DbTarget = { tag_id: 'tag-1', column_name: 'line_a', enabled: true };
    expect(toStudioV2DatabaseTargetRequest(target, 'setup-10')).toEqual({
      column_name: 'line_a', enabled: true, row_group_id: undefined, expected_setup_revision: 'setup-10',
    });
  });
});
