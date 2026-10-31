import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from './writeGroup/candidates';
import { basicManagedGroupDraft, isBasicManagedWriteGroup } from './basicRecording';

const key = (value: { device_id: string; point_id: string; tag_id: string }) => JSON.stringify([value.device_id, value.point_id, value.tag_id]);

/** Keep canonical bindings and policies; only membership changes in this draft. */
export function basicMemberChange(group: WriteGroup, deviceId: string, candidates: GroupCandidate[]) {
  if (group.recording_profile?.kind === 'mixed' || group.row_policy.interval_seconds === undefined || !isBasicManagedWriteGroup(group, deviceId) || !candidates.length || candidates.some(c => c.device_id !== deviceId)) return undefined;
  // The managed backend reasserts these bindings. Customized variants must
  // stay in Advanced rather than silently normalize an operator's settings.
  if ((group.write_policy.mode && group.write_policy.mode !== 'append') ||
    (group.write_policy.dedupe_capability && group.write_policy.dedupe_capability !== 'receipt') ||
    (group.row_policy.provenance_column && group.row_policy.provenance_column !== 'provenance') ||
    group.members.some(m => m.target_column && !/^v_[a-f0-9]{24}$/.test(m.target_column))) return undefined;
  const keys = new Set(candidates.map(key));
  if (keys.size !== candidates.length) return undefined;
  const existing = new Map(group.members.map(m => [key(m), m]));
  const added = candidates.filter(c => !existing.has(key(c)));
  const removed = group.members.filter(m => !keys.has(key(m)));
  const draft = basicManagedGroupDraft(group, group.row_policy.interval_seconds, group.row_policy.incomplete_policy ?? 'skip_row');
  // Keep editable bindings. The server revalidates source/mapping revisions
  // under the expected workspace + group revision, as in the existing editor.
  draft.members = [...draft.members.filter(m => keys.has(key(m))), ...added.map(c => ({ device_id: c.device_id, point_id: c.point_id, tag_id: c.tag_id, target_column: '', required: true }))];
  return { added, removed, draft };
}
