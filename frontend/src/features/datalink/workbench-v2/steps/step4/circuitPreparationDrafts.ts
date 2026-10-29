import type { WriteGroup, WriteGroupProfileKind } from '../../../../../types/studioV2WriteGroup';
import { buildCircuitProfileGroupDraft, canonicalProfileGroupMatchesDraft, purposeOfLayout } from '../../state/circuitProfileBasic';
import type { GroupCandidate } from '../../state/writeGroup/candidates';

type DraftInput = Parameters<typeof buildCircuitProfileGroupDraft>[0];
/** Reuses the native draft builder and the existing saved-group dirty comparison unchanged. */
export function circuitPreparationDrafts(groups: readonly WriteGroup[], profiles: readonly {
  profile: WriteGroupProfileKind; intervalValue: number; intervalSeconds: number; candidates: GroupCandidate[];
}[], input: Omit<DraftInput, 'profile' | 'intervalValue' | 'intervalSeconds' | 'selectedCandidates' | 'destination'> & {
  destination?: DraftInput['destination'];
}, connectorKind: Parameters<typeof canonicalProfileGroupMatchesDraft>[2]) {
  const { destination } = input;
  const drafts = destination ? profiles.map(profile => buildCircuitProfileGroupDraft({ ...input, destination,
    profile: profile.profile, intervalValue: profile.intervalValue, intervalSeconds: profile.intervalSeconds,
    selectedCandidates: profile.candidates })) : [];
  const dirtyGroupIds = new Set(groups.filter(group => {
    const profile = profiles.find(item => item.profile === group.recording_profile?.kind);
    if (!profile || !destination) return true;
    const draft = buildCircuitProfileGroupDraft({ ...input, destination, profile: profile.profile,
      intervalValue: profile.intervalValue, intervalSeconds: profile.intervalSeconds, selectedCandidates: profile.candidates,
      purpose: purposeOfLayout(group.layout_version) ?? input.purpose, reduction: group.recording_profile?.sample_reduction ?? 'latest' });
    return !canonicalProfileGroupMatchesDraft(group, draft, connectorKind);
  }).map(group => group.id));
  return { drafts, dirtyGroupIds };
}
