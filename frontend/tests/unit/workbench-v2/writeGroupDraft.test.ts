import { describe, expect, it } from 'vitest';
import {
  draftIsDirty,
  draftToRequestGroup,
  groupToDraft,
  newDraft,
  validateDraft,
} from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';

const group: WriteGroup = {
  id: 'group-1', workspace_id: 'ws-1', revision: 'rev-2', applied_revision: 'rev-1', name: 'Line A', status: 'ready',
  members: [
    { device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', entity_key: 'line-1', source_revision: 's1', mapping_revision: 'm1', target_column: 'temperature', required: true },
    { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', source_revision: 's2', mapping_revision: 'm2', target_column: 'pressure', required: false, max_age_seconds: 20 },
  ],
  destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: 'main', table_name: 'readings', storage_strategy: 'custom' },
  row_policy: { interval_seconds: 10, allowed_lateness_seconds: 2, incomplete_policy: 'skip_row', entity_key_column: 'entity' },
  write_policy: { dedupe_capability: 'receipt' },
  migration: {}, created_at: 'a', updated_at: 'b',
};

describe('write group editor draft', () => {
  it('round-trips a saved group through the editor without losing fields', () => {
    const draft = groupToDraft(group);
    expect(draft.members.map((m) => [m.point_id, m.target_column, m.required, m.entity_key])).toEqual([
      ['pt-1', 'temperature', true, 'line-1'], ['pt-2', 'pressure', false, ''],
    ]);
    const request = draftToRequestGroup(draft, 'ws-1');
    expect(request.workspace_id).toBe('ws-1');
    expect(request.destination).toEqual({
      connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: 'main', table_name: 'readings', storage_strategy: 'custom',
    });
    expect(request.row_policy).toMatchObject({ interval_seconds: 10, allowed_lateness_seconds: 2, incomplete_policy: 'skip_row', entity_key_column: 'entity' });
    expect(request.write_policy).toEqual({ dedupe_capability: 'receipt' });
    expect(request.members[1]).toMatchObject({ max_age_seconds: 20, required: false });
    expect(request.members[0]).not.toHaveProperty('max_age_seconds');
    // Source and mapping revisions are server-owned: the draft sends none, so the server pins the current ones.
    expect(request.members[0]).not.toHaveProperty('source_revision');
  });

  it('preserves hidden migrated contracts when the editor only changes visible fields', () => {
    const migrated = {
      ...group,
      members: group.members.map((member, index) => ({
        ...member,
        ...(index === 0 ? { measurement_id: 'measurement-1' } : {}),
      })),
      row_policy: {
        ...group.row_policy,
        group_key_columns: ['entity'],
        unique_key_columns: ['entity', 'observed_at'],
        value_column: 'value',
        quality_column: 'quality',
      },
    };

    const draft = groupToDraft(migrated);
    const request = draftToRequestGroup({ ...draft, name: 'Renamed' }, 'ws-1');

    expect(request.name).toBe('Renamed');
    expect(request.row_policy).toMatchObject({
      group_key_columns: ['entity'],
      unique_key_columns: ['entity', 'observed_at'],
      value_column: 'value',
      quality_column: 'quality',
    });
    expect(request.members[0]).toMatchObject({ measurement_id: 'measurement-1' });
  });

  it('is dirty only when something the server stores actually changed', () => {
    const draft = groupToDraft(group);
    expect(draftIsDirty(draft, group)).toBe(false);
    expect(draftIsDirty({ ...draft, name: 'Line B' }, group)).toBe(true);
    expect(draftIsDirty({ ...draft, members: draft.members.slice(1) }, group)).toBe(true);
    expect(draftIsDirty({ ...draft, members: draft.members.map((m, i) => (i === 0 ? { ...m, target_column: 'temp' } : m)) }, group)).toBe(true);
  });

  it('starts a new draft with explicit, visible defaults and nothing preselected', () => {
    const draft = newDraft({ connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: 'readings' });
    expect(draft.members).toEqual([]);
    expect(draft.incomplete_policy).toBe('skip_row');
    expect(draft.storage_strategy).toBe('custom');
    expect(draft.interval_seconds).toBeGreaterThan(0);
  });

  describe('validateDraft', () => {
    const base = groupToDraft(group);
    const codes = (draft = base, ctx = {}) => validateDraft(draft, { destinationSaved: true, ...ctx }).map((issue) => issue.code);

    it('accepts a complete draft', () => expect(codes()).toEqual([]));
    it('requires a name, members, a saved destination and a table', () => {
      expect(codes({ ...base, name: ' ' })).toContain('name-required');
      expect(codes({ ...base, members: [] })).toContain('members-required');
      expect(codes(base, { destinationSaved: false })).toContain('destination-unsaved');
      expect(codes({ ...base, table_name: '' })).toContain('table-required');
    });
    it('requires a column for each member and rejects shared columns without distinct entity keys', () => {
      expect(validateDraft({ ...base, members: [{ ...base.members[0], target_column: '' }] }, { destinationSaved: true })
        .find((issue) => issue.code === 'column-required')?.member_key).toBe(base.members[0].key);
      const shared = { ...base, members: base.members.map((m) => ({ ...m, target_column: 'v', entity_key: '' })) };
      expect(codes(shared)).toContain('column-conflict');
      const distinct = { ...base, members: base.members.map((m, i) => ({ ...m, target_column: 'v', entity_key: `line-${i}` })) };
      expect(codes(distinct)).not.toContain('column-conflict');
    });
    it('checks the row policy', () => {
      expect(codes({ ...base, interval_seconds: 0 })).toContain('interval-invalid');
      expect(codes({ ...base, allowed_lateness_seconds: -1 })).toContain('lateness-invalid');
      expect(codes({ ...base, incomplete_policy: 'partial', provenance_column: '' })).toContain('partial-needs-provenance');
      expect(codes({ ...base, incomplete_policy: 'partial', provenance_column: 'prov' })).not.toContain('partial-needs-provenance');
    });
    it('rejects equivalent SQL columns in the same persisted entity before Save or Apply', () => {
      const shared = { ...base, members: base.members.map((m, i) => ({
        ...m, target_column: i === 0 ? 'Value' : 'value', entity_key: 'A',
      })) };
      expect(codes(shared)).toContain('column-conflict');
      const distinct = { ...shared, members: shared.members.map((m, i) => ({ ...m, entity_key: i === 0 ? 'A' : 'B' })) };
      expect(codes(distinct)).not.toContain('column-conflict');
    });
    it('treats a member with an entity key as needing the entity key column', () => {
      expect(codes({ ...base, entity_key_column: '' })).toContain('entity-column-required');
    });
  });
});
