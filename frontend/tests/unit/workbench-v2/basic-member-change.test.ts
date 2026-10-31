import { expect, it } from 'vitest';
import { basicManagedGroupDraft } from '@/features/datalink/workbench-v2/state/basicRecording';
import { basicMemberChange } from '@/features/datalink/workbench-v2/state/basicRecordingMembers';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
const managed = () => savedGroup({ members: savedGroup().members.map((m,i) => ({ ...m, target_column: `v_${String(i).repeat(24)}` })), basic_managed_device_id: 'dev-1', destination: { ...savedGroup().destination, storage_strategy: 'managed' }, row_policy: { ...savedGroup().row_policy, provenance_column: 'provenance', record_key_column: 'record_id' }, write_policy: { mode: 'append', dedupe_capability: 'receipt' } });
it('changes only membership and preserves canonical columns, optional required, age, dedupe and provenance', () => {
  const group = managed(); group.members[0] = { ...group.members[0], required: false, max_age_seconds: 30, measurement_id: 'measurement' };
  const { candidates } = buildGroupCandidates(savedGroupState());
  const change = basicMemberChange(group, 'dev-1', [candidates[0], candidates[2]])!;
  expect(change.added.map(c => c.tag_id)).toEqual(['tag-2']);
  expect(change.removed.map(c => c.tag_id)).toEqual(['tag-1']);
  expect(change.draft.members[0]).toEqual(basicManagedGroupDraft(group, group.row_policy.interval_seconds!, group.row_policy.incomplete_policy! ).members[0]);
  expect(change.draft.members[0]).not.toHaveProperty('source_revision');
  expect(change.draft.members[0]).not.toHaveProperty('mapping_revision');
  expect(change.draft.members[1]).toMatchObject({ point_id: 'pt-2', tag_id: 'tag-2', target_column: '', required: true });
  expect(change.draft.row_policy).toEqual(group.row_policy);
  expect(change.draft.write_policy).toEqual(group.write_policy);
  expect(group.members).toHaveLength(2);
});
it('refuses custom, entity groups and cross-device or empty proposals', () => {
  const candidates = buildGroupCandidates(savedGroupState()).candidates;
  expect(basicMemberChange(savedGroup(), 'dev-1', candidates)).toBeUndefined();
  const group = managed(); group.members[0].entity_key = 'machine';
  expect(basicMemberChange(group, 'dev-1', candidates)).toBeUndefined();
  expect(basicMemberChange(managed(), 'dev-1', [])).toBeUndefined();
  expect(basicMemberChange(managed(), 'dev-1', [{ ...candidates[0], device_id: 'other' }])).toBeUndefined();
});
