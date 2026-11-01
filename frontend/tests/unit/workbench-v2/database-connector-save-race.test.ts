import { workbenchV2Reducer } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
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

it('does not overwrite an A→B→A field edit with an earlier canonical normalization', () => {
  const state = savedGroupState();
  const start = { ...state, db: { ...state.db, connector: hydrateStudioV2DatabaseConnector({ ...record, name: ' Linux ' }) } };
  const submittedState = workbenchV2Reducer(start, { type: 'updateDbConnector', patch: { name: ' Linux ' } });
  const edited = workbenchV2Reducer(submittedState, { type: 'updateDbConnector', patch: { name: 'other' } });
  const latest = workbenchV2Reducer(edited, { type: 'updateDbConnector', patch: { name: ' Linux ' } });
  expect(mergeSavedStudioV2DatabaseConnector(latest.db.connector, submittedState.db.connector, { ...record, name: 'Linux' }))
    .toMatchObject({ name: ' Linux ', save_state: 'saving' });
});

it('accepts canonical untouched fields without replacing a newer field edit', () => {
  const submitted = hydrateStudioV2DatabaseConnector({ ...record, name: ' Linux ', database: ' target.sqlite ' });
  const current = { ...submitted, name: 'newer name' };
  expect(mergeSavedStudioV2DatabaseConnector(current, submitted, { ...record, name: 'Linux', database: 'target.sqlite' }))
    .toMatchObject({ name: 'newer name', database: 'target.sqlite', save_state: 'saving' });
});

it('ignores an old workspace ACK rather than adopting its setup revision or identity', () => {
  const submitted = hydrateStudioV2DatabaseConnector(record);
  const current = { ...submitted, workspace_id: 'ws-2', connector_id: 'other', setup_revision: 'other-setup' };
  expect(mergeSavedStudioV2DatabaseConnector(current, submitted, record)).toEqual(current);
});
