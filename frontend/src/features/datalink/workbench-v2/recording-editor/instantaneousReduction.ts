import type { WriteGroup, WriteGroupSampleReduction, WriteGroupValueRole } from '../../../../types/studioV2WriteGroup';
import { normalizeTypedEnvelope } from '../../../../utils/safeJson';
import { isValidRecordingProfile } from '../../../../utils/studioV2RecordingProfile';
import { defaultMemberRole } from '../state/circuitProfileBasic';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import { memberKey, type EditorDraft, type EditorMember } from '../state/writeGroup/draft';

/** True only for a saved-shape instantaneous profile record; mixed and cumulative records are never edited here. */
export function isInstantaneousProfileDraft(draft: EditorDraft): boolean {
  return draft.recording_profile?.kind === 'instantaneous' &&
    (draft.layout_version === 'managed_long_v1' || draft.layout_version === 'managed_wide_profile_v1');
}

export function draftReduction(draft: EditorDraft): WriteGroupSampleReduction {
  const profile = draft.recording_profile;
  return profile && profile.kind !== 'mixed' ? profile.sample_reduction ?? 'latest' : 'latest';
}

function withoutResultType(member: EditorMember): EditorMember {
  const next = { ...member };
  delete next.result_value_type;
  return next;
}

/**
 * The role a member takes once the operator confirmed gauge use: an explicit saved role is kept (a counter_raw
 * binding stays a counter), a known non-numeric source keeps its latest value as state, and every other source is
 * the gauge the operator confirmed. The backend still verifies the real tag type and rejects what cannot be averaged.
 */
function meanRole(member: EditorMember, candidates: GroupCandidate[]): WriteGroupValueRole {
  if (member.value_role) return member.value_role;
  const type = member.source_value_type ?? candidates.find(candidate => candidate.point_id === member.point_id && candidate.tag_id === member.tag_id)?.target_type;
  return type === undefined ? 'gauge' : defaultMemberRole({ target_type: type });
}

/**
 * Changes only the reduction of an instantaneous draft. The result type is always re-derived by the backend after a
 * reduction change; returning to latest restores the saved member metadata when the saved record was latest, so an
 * undone change leaves no new shape behind. Cumulative and mixed drafts are returned unchanged.
 */
export function applyInstantaneousReduction(draft: EditorDraft, reduction: WriteGroupSampleReduction, saved: WriteGroup | undefined, candidates: GroupCandidate[]): EditorDraft {
  const profile = draft.recording_profile;
  if (!profile || profile.kind !== 'instantaneous' || draftReduction(draft) === reduction) return draft;
  if (reduction === 'mean') {
    return {
      ...draft,
      recording_profile: { ...profile, sample_reduction: 'mean' },
      members: draft.members.map(member => ({ ...withoutResultType(member), value_role: meanRole(member, candidates) })),
    };
  }
  const savedProfile = saved?.recording_profile;
  const savedLatest = savedProfile !== undefined && savedProfile.kind === 'instantaneous' && (savedProfile.sample_reduction ?? 'latest') === 'latest';
  const nextProfile = { ...profile };
  delete nextProfile.sample_reduction;
  if (savedLatest && savedProfile.sample_reduction === 'latest') nextProfile.sample_reduction = 'latest';
  const savedMembers = new Map((saved?.members ?? []).map(member => [memberKey(member.point_id, member.tag_id), member]));
  return {
    ...draft,
    recording_profile: nextProfile,
    members: draft.members.map(member => {
      const original = savedLatest ? savedMembers.get(member.key) : undefined;
      if (!original) return withoutResultType(member);
      const next = withoutResultType(member);
      delete next.value_role;
      return {
        ...next,
        ...(original.value_role !== undefined ? { value_role: original.value_role } : {}),
        ...(original.result_value_type !== undefined ? { result_value_type: original.result_value_type } : {}),
      };
    }),
  };
}

/** Text typed into the interval field; anything but a whole number becomes 0 so validation blocks the save. */
export function parseIntervalText(text: string): number {
  const trimmed = text.trim();
  return /^\d+$/.test(trimmed) ? Number(trimmed) : 0;
}

export function applyInstantaneousInterval(draft: EditorDraft, interval: number): EditorDraft {
  const profile = draft.recording_profile;
  if (!profile || profile.kind !== 'instantaneous') return draft;
  return { ...draft, recording_profile: { ...profile, interval } };
}

export function instantaneousIntervalValid(draft: EditorDraft): boolean {
  return isValidRecordingProfile(draft.recording_profile);
}

/** The saved record's reduction or interval differs from the draft: saving creates a revision that needs a schema check. */
export function instantaneousMethodChanged(draft: EditorDraft, saved: WriteGroup | undefined): boolean {
  const profile = draft.recording_profile;
  const savedProfile = saved?.recording_profile;
  if (!saved || !profile || profile.kind !== 'instantaneous' || !savedProfile || savedProfile.kind !== 'instantaneous') return false;
  return profile.interval !== savedProfile.interval || draftReduction(draft) !== (savedProfile.sample_reduction ?? 'latest');
}

const REJECTION_CODES = new Set(['WRITE_GROUP_INVALID', 'member_metadata_invalid', 'profile_interval_invalid']);

/** A save failure the backend attributes to the recording configuration, shown at the method fields. */
export function instantaneousSaveRejected(error: unknown): boolean {
  if (error === null || error === undefined) return false;
  const code = normalizeTypedEnvelope(error).code;
  return code !== undefined && REJECTION_CODES.has(code);
}
