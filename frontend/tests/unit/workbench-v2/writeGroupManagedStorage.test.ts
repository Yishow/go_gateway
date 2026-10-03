import { describe, expect, it } from 'vitest';
import { draftToRequestGroup, newDraft, validateDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { proposeColumns } from '@/features/datalink/workbench-v2/state/writeGroup/proposal';

const destination = { connector_id: 'connector-1', connector_revision: 'connector-rev-1', table_schema: 'main', table_name: '' };

function managedDraft() {
  const draft = newDraft(destination);
  draft.name = 'Line A';
  draft.storage_strategy = 'managed';
  draft.members = [{
    key: 'point-1|tag-1', device_id: 'device-1', point_id: 'point-1', tag_id: 'tag-1',
    entity_key: '', target_column: '', required: true,
  }];
  return draft;
}

describe('managed write-group draft contract', () => {
  it('allows an empty managed table and server-generated member columns', () => {
    const draft = managedDraft();
    expect(validateDraft(draft, { destinationSaved: true, storageStrategy: 'managed' })).toEqual([]);
    expect(draftToRequestGroup(draft, 'workspace-1')).toMatchObject({
      destination: { storage_strategy: 'managed', table_name: '' },
      members: [{ target_column: '' }],
    });
  });

  it('retains server-managed metadata bindings when an existing group is edited', () => {
    const draft = managedDraft();
    draft.record_key_column = 'record_id';
    draft.group_id_column = 'group_id';
    draft.device_id_column = 'device_id';
    draft.bucket_start_column = 'bucket_start';
    draft.provenance_column = 'provenance';
    expect(draftToRequestGroup(draft, 'workspace-1').row_policy).toMatchObject({
      record_key_column: 'record_id', group_id_column: 'group_id', device_id_column: 'device_id',
      bucket_start_column: 'bucket_start', provenance_column: 'provenance',
    });
  });

  it('uses exact PostgreSQL integer and uint64 proposals', () => {
    const candidate = (target_type: 'int16' | 'int32' | 'uint64') => ({
      key: target_type, device_id: 'device-1', point_id: target_type, tag_id: target_type,
      tag_key: `line.${target_type}`, label: target_type, device_name: 'PLC A', address: '1', target_type,
    });
    expect(proposeColumns([candidate('int16'), candidate('int32'), candidate('uint64')], [], 'postgres').map((item) => item.sql_type))
      .toEqual(['BIGINT', 'BIGINT', 'NUMERIC(20,0)']);
    expect(proposeColumns([candidate('int16'), candidate('uint64')], [], 'sqlite').map((item) => item.sql_type))
      .toEqual(['INTEGER', 'TEXT']);
  });
});
