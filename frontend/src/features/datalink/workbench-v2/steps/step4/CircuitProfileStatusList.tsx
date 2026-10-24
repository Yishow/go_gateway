import * as React from 'react';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from '../../state/writeGroup/candidates';
import { CircuitProfileStatus } from './CircuitProfileStatus';

const ORDER = { instantaneous: 0, cumulative: 1 } as const;

/** The status of each canonical profile group of the selected circuit, instantaneous first. */
export const CircuitProfileStatusList: React.FC<{ groups: WriteGroup[]; candidates: GroupCandidate[]; onOpenAdvanced?: (groupId: string) => void; draftEnabled?: Record<'instantaneous' | 'cumulative', boolean> }> = ({ groups, candidates, onOpenAdvanced, draftEnabled }) => {
  const labels = React.useMemo(() => Object.fromEntries(candidates.map((candidate) => [candidate.tag_id, candidate.tag_key])), [candidates]);
  const ordered = [...groups].sort((a, b) => ORDER[a.recording_profile?.kind ?? 'instantaneous'] - ORDER[b.recording_profile?.kind ?? 'instantaneous']);
  return <>{ordered.map((group) => <CircuitProfileStatus key={group.id} group={group} sourceLabels={labels} onOpenAdvanced={onOpenAdvanced} draftOff={draftEnabled ? !draftEnabled[group.recording_profile?.kind ?? 'instantaneous'] : false} />)}</>;
};
