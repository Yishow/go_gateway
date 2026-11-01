import { describe, expect, it } from 'vitest';
import { workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { hydrateStudioV2DatabaseConnector, toStudioV2DatabaseConfigRequest } from '@/features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import type { StudioV2WorkspaceDatabaseConfigRecord } from '@/types/datalink';
import { savedGroupState } from '../../fixtures/writeGroupState';

const record: StudioV2WorkspaceDatabaseConfigRecord = {
  id: 'db-1', workspace_id: 'ws-1', identity_revision: 'identity-1', setup_revision: 'setup-2',
  kind: 'postgres', name: 'Linux', host: 'localhost', port: 5432, database: 'gateway', username: 'writer', schema: 'public',
  table: 'sensor_readings', write_mode: 'insert', write_interval_seconds: 5, timestamp_column: 'ts',
  status: 'ready', save_state: 'saved', created_at: '2026-10-31T08:00:00Z', updated_at: '2026-10-31T08:00:00Z',
};
function initial() {
  const state = savedGroupState();
  return { ...state, db: { ...state.db, connector: hydrateStudioV2DatabaseConnector(record) } };
}
const edit = (patch: Extract<WorkbenchV2Action, { type: 'updateDbConnector' }>['patch']): WorkbenchV2Action => ({ type: 'updateDbConnector', patch });

describe('database connector request-bound reducer acknowledgements', () => {
  it('preserves later raw invalid input, password intent, targets and row groups', () => {
    const state = initial();
    const submitted = workbenchV2Reducer(state, edit({ name: 'Linux ' }));
    let current = workbenchV2Reducer(submitted, edit({ name: '  ', password: '  exact secret  ', password_required: false, clear_password: false }));
    current = workbenchV2Reducer(current, { type: 'upsertDbTarget', pointId: 'new-point', target: { column_name: 'later_column', enabled: true, tag_id: 'tag-new' } });
    current = workbenchV2Reducer(current, { type: 'setDbRowGroups', rowGroups: [{ id: 'new-group', table_name: 'new_table', member_point_ids: ['new-point'] }] });
    const next = workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: submitted.db.connector, rowGroups: [], record });
    expect(next.db.connector).toMatchObject({ name: '  ', password: '  exact secret  ', password_required: false, clear_password: false, save_state: 'draft-invalid' });
    expect(next.db.targets).toEqual(current.db.targets);
    expect(next.db.row_groups).toEqual(current.db.row_groups);
    expect(toStudioV2DatabaseConfigRequest(next.db.connector)).not.toHaveProperty('field_generations');
  });

  it.each(['saving', 'save-error', 'draft-invalid'] as const)('does not let a GET acknowledge %s local intent', save_state => {
    const state = initial();
    const current = workbenchV2Reducer(state, edit({ name: 'raw local text', save_state }));
    expect(workbenchV2Reducer(current, { type: 'hydrateDbConnector', record })).toBe(current);
  });

  it('rejects old GET and ACK after workspace switch', () => {
    const state = initial();
    const current = workbenchV2Reducer(state, edit({ workspace_id: 'ws-2', connector_id: 'db-2', setup_revision: 'other-setup' }));
    expect(workbenchV2Reducer(current, { type: 'hydrateDbConnector', record })).toBe(current);
    expect(workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: state.db.connector, record })).toBe(current);
  });

  it('does not resurrect an old connector after selection is cleared or changed A→B→A', () => {
    const state = initial();
    let current = workbenchV2Reducer(state, edit({ connector_id: undefined, identity_revision: undefined, name: 'new unsent connector' }));
    let next = workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: state.db.connector, record });
    expect(next.db.connector).toMatchObject({ connector_id: undefined, name: 'new unsent connector', setup_revision: 'setup-2', save_state: 'saving' });
    current = workbenchV2Reducer(current, edit({ connector_id: 'db-1', identity_revision: 'identity-new', name: 'new selection' }));
    next = workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: state.db.connector, record });
    expect(next.db.connector).toMatchObject({ connector_id: 'db-1', identity_revision: 'identity-new', name: 'new selection', save_state: 'saving' });
  });

  it('does not import a masked response password or trim a newly typed one', () => {
    const state = initial();
    const current = workbenchV2Reducer(state, edit({ password: '  exact  ', clear_password: false }));
    const masked = { ...record, password: '********' };
    const next = workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: state.db.connector, record: masked });
    expect(next.db.connector.password).toBe('  exact  ');
    expect(toStudioV2DatabaseConfigRequest(next.db.connector).password).toBe('  exact  ');
  });

  it('ignores an ACK from an earlier workspace visit even after returning to the same workspace', () => {
    const state = initial();
    const other = workbenchV2Reducer(state, edit({ workspace_id: 'ws-2' }));
    const current = workbenchV2Reducer(other, edit({ workspace_id: 'ws-1', setup_revision: 'new-visit' }));
    expect(workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: state.db.connector, record })).toBe(current);
  });
  it('does not restore old canonical row groups into a newly selected connector', () => {
    const state = initial();
    const current = workbenchV2Reducer(state, edit({ connector_id: 'db-2', identity_revision: 'identity-2' }));
    const saved = { ...record, row_groups: [{ id: 'old-group', connector_id: 'db-1', table_name: 'old_table', member_point_ids: ['old-point'] }] };
    const next = workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: state.db.connector, rowGroups: [], record: saved });
    expect(next.db.row_groups ?? []).toEqual([]);
  });

  it('rejects a same-workspace response for another saved connector ID', () => {
    const state = initial();
    const current = workbenchV2Reducer(state, edit({ name: 'new local name' }));
    expect(workbenchV2Reducer(current, { type: 'acknowledgeDbConnectorSave', submitted: state.db.connector, record: { ...record, id: 'wrong-connector' } })).toBe(current);
  });

});
