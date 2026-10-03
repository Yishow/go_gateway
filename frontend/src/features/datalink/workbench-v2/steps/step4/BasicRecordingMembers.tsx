import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type { GroupCandidate } from '../../state/writeGroup/candidates';

interface BasicRecordingMembersProps {
  candidates: GroupCandidate[];
  group?: WriteGroup;
}

function memberKey(pointId: string, tagId: string): string {
  return `${pointId}|${tagId}`;
}

export const BasicRecordingMembers: React.FC<BasicRecordingMembersProps> = ({ candidates, group }) => {
  const { t } = useTranslation('workbench-v2');
  const candidateByMember = new Map(candidates.map((candidate) => [memberKey(candidate.point_id, candidate.tag_id), candidate]));
  const rows = group
    ? group.members.map((member) => ({ key: memberKey(member.point_id, member.tag_id), candidate: candidateByMember.get(memberKey(member.point_id, member.tag_id)) }))
    : candidates.map((candidate) => ({ key: candidate.key, candidate }));

  return (
    <section className="space-y-2 rounded border border-slate-800/80 bg-slate-950/20 p-3" aria-label={t('step4.basic.selected_data')} data-testid="basic-recording-selected-members">
      <h4 className="text-xs font-semibold text-slate-200">{t('step4.basic.selected_data')}</h4>
      <ul className="space-y-1 text-[11px] text-slate-300">
        {rows.map((row) => (
          <li key={row.key} className="flex flex-wrap items-baseline justify-between gap-x-3 gap-y-1">
            {row.candidate ? (
              <>
                <span className="font-medium">{row.candidate.label || row.candidate.tag_key}</span>
                <span className="text-slate-500">{row.candidate.device_name} · {row.candidate.address} · {row.candidate.target_type}</span>
              </>
            ) : (
              <>
                <span className="font-medium">{t('step4.basic.saved_member_unavailable')}</span>
                <span className="text-slate-500">{t('step4.basic.saved_member_review')}</span>
              </>
            )}
          </li>
        ))}
      </ul>
    </section>
  );
};
