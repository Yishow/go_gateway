import { expect, it } from 'vitest';
import { hydrateStudioV2DatabaseConnector, mergeSavedStudioV2DatabaseConnector } from '@/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { destinationIsSaved } from '@/features/datalink/workbench-v2/steps/step4/basicRecordingPanelHelpers';
import type { StudioV2WorkspaceDatabaseConfigRecord } from '@/types/datalink';
import { savedGroupState } from '../../fixtures/writeGroupState';

const record: StudioV2WorkspaceDatabaseConfigRecord = {
  id: 'new-sqlite', workspace_id: 'ws-1', identity_revision: 'identity-1', setup_revision: 'setup-2',
  kind: 'sqlite', name: 'SQLite', host: '', port: 0, database: '/tmp/owned-test.db', username: '', schema: '',
  table: 'readings', write_mode: 'insert', write_interval_seconds: 5, timestamp_column: 'ts',
  status: 'ready', save_state: 'saved', created_at: '2026-10-24T00:00:00Z', updated_at: '2026-10-24T00:00:00Z',
};

it('keeps the first save complete when canonical query hydration wins the response race', () => {
  const hydrated = hydrateStudioV2DatabaseConnector(record);
  const submitted = { ...hydrated, connector_id: undefined, identity_revision: undefined, persisted: false, save_state: 'saving' as const };
  const merged = mergeSavedStudioV2DatabaseConnector(hydrated, submitted, record);
  const state = savedGroupState();
  expect(destinationIsSaved({ ...state, db: { ...state.db, connector: merged } })).toBe(true);
  expect(merged.connector_id).toBe(record.id);
});

it('preserves a different connector selected while the first save is in flight', () => {
  const submitted = hydrateStudioV2DatabaseConnector(record);
  const current = { ...submitted, connector_id: 'other-sqlite', identity_revision: 'other-identity' };
  expect(mergeSavedStudioV2DatabaseConnector(current, submitted, record)).toMatchObject({ connector_id: 'other-sqlite', identity_revision: 'other-identity', save_state: 'saving' });
});

it('preserves a newer local edit while its first save response arrives', () => {
  const submitted = hydrateStudioV2DatabaseConnector(record);
  expect(mergeSavedStudioV2DatabaseConnector({ ...submitted, name: 'New local name' }, submitted, record)).toMatchObject({ name: 'New local name', save_state: 'saving' });
});
