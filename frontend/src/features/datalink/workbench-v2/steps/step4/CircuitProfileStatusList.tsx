import type { SavedRecordingManagementProps } from './SavedRecordingManagement';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from '../../state/writeGroup/candidates';
import { CircuitProfileStatus } from './CircuitProfileStatus';

const ORDER = { instantaneous: 0, cumulative: 1, mixed: 2 } as const;

/** The status of each canonical profile group of the selected circuit, instantaneous first. */
export const CircuitProfileStatusList: React.FC<{ lifecycle?: Omit<SavedRecordingManagementProps, 'group'>; groups: WriteGroup[]; candidates: GroupCandidate[]; onOpenAdvanced?: (groupId: string) => void; draftEnabled?: Record<'instantaneous' | 'cumulative', boolean> }> = ({ groups, candidates, onOpenAdvanced, draftEnabled, lifecycle }) => {
  const { t } = useTranslation('workbench-v2');
  const labels = React.useMemo(() => Object.fromEntries(candidates.map((candidate) => [candidate.tag_id, candidate.label.trim() || candidate.address])), [candidates]);
  const ordered = [...groups].sort((a, b) => ORDER[a.recording_profile?.kind ?? 'instantaneous'] - ORDER[b.recording_profile?.kind ?? 'instantaneous']);
  return <>{ordered.map((group) => group.recording_profile?.kind === 'mixed' ? <section key={group.id} data-testid="mixed-recording-readonly">
    <span>{group.name}</span><p role="status">{t('step4.group_kind_unsupported')}</p>
    {onOpenAdvanced && <button type="button" onClick={() => onOpenAdvanced(group.id)}>{t('step4.basic.status.open_advanced')}</button>}
  </section> : <CircuitProfileStatus lifecycle={lifecycle} key={group.id} group={group} sourceLabels={labels} onOpenAdvanced={onOpenAdvanced} draftOff={draftEnabled ? !draftEnabled[group.recording_profile?.kind ?? 'instantaneous'] : false} />)}</>;
};
