import { memberKey, type EditorMember } from '../state/writeGroup/draft';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import { createPicker, pendingAddition, restoreMember, sourceKey, type CandidateEvidence, type MemberHistory, type PickerProjection, type SourceIdentity } from './recordingPickerState';
/** Candidates already came from canonical saved-source projection. Eligible means
 * selectable local intent, never source revision/CAS or a measurement semantic. */
export function recordingPickerProjection(scope: string, candidates: readonly GroupCandidate[], members: readonly SourceIdentity[], evidence: Readonly<Record<string, CandidateEvidence | undefined>> = {}): PickerProjection {
  return { scope, candidates, members, evidence: Object.fromEntries(candidates.map(candidate => [sourceKey(candidate),
    evidence[sourceKey(candidate)] ?? { availability: 'eligible' as const }])) };
}
/** Atomic local intent merge into the existing GroupEditor owner, not another draft. */
export function addRecordingMembers(members: readonly EditorMember[], history: MemberHistory<EditorMember>, identities: readonly SourceIdentity[], projection: PickerProjection): EditorMember[] | undefined {
  if (history.scope !== projection.scope || new Set(identities.map(sourceKey)).size !== identities.length) return;
  const current = { ...projection, members };
  if (!pendingAddition({ ...createPicker(history.scope), pending: identities }, current)) return;
  let next: readonly EditorMember[] = members;
  for (const identity of identities) next = restoreMember(next, history, identity, current, source => ({
    device_id: source.device_id, point_id: source.point_id, tag_id: source.tag_id, key: memberKey(source.point_id, source.tag_id), entity_key: '', target_column: '', required: true,
  }));
  return [...next];
}
