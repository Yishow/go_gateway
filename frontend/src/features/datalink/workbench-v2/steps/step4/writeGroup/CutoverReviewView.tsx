import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';
import type { CutoverGroupView, WriteGroupCutoverReview } from '../../../../../../types/studioV2WriteGroupCutover';
import { CUTOVER_ISSUE_REPAIR } from '../../../state/cutoverMessages';
import type { GroupCandidate } from '../../../state/writeGroup/candidates';

function formatMemberRef(ref: { point_id: string; tag_id: string; device_id?: string }, candidates?: GroupCandidate[], source?: WriteGroup): string {
  const candidate = candidates?.find((c) => (ref.point_id && c.point_id === ref.point_id) || (ref.tag_id && c.tag_id === ref.tag_id));
  if (candidate?.tag_key) return candidate.tag_key;
  if (candidate?.label) return candidate.label;
  const member = source?.members.find((m) => (ref.point_id && m.point_id === ref.point_id) || (ref.tag_id && m.tag_id === ref.tag_id));
  if (member?.target_column) return member.target_column;
  return ref.tag_id || ref.point_id || 'unavailable';
}

const GroupRow: React.FC<{ view: CutoverGroupView; role: 'source' | 'target' }> = ({ view, role }) => {
  const { t } = useTranslation('workbench-v2');
  const long = view.layout === 'managed_long_v1';
  return (
    <tr data-testid={`cutover-review-${role}-${view.group_id}`} className="border-t border-slate-800 align-top">
      <th scope="row" className="py-1 pr-3 text-left font-medium text-slate-100">{view.name}</th>
      <td className="pr-3">{t(long ? 'step4.cutover.layout_long' : 'step4.cutover.layout_wide')}</td>
      <td className="pr-3">{view.storage_mode ? t(`step4.basic.circuit.storage_${view.storage_mode}`) : '—'}</td>
      <td className="pr-3">{t('step4.cutover.period_seconds', { count: view.interval_seconds })}</td>
      <td className="pr-3">{view.table_schema ? `${view.table_schema}.${view.table_name}` : view.table_name}</td>
      <td className="pr-3">{t(view.schema_confirmed ? 'step4.cutover.schema_confirmed' : 'step4.cutover.schema_pending')}</td>
      <td>{view.member_count}</td>
    </tr>
  );
};

export interface CutoverReviewViewProps {
  review: WriteGroupCutoverReview;
  /** The selection, boundary or a saved revision changed after this review. */
  stale: boolean;
  onOpenGroup?: (groupId: string) => void;
  candidates?: GroupCandidate[];
  sourceGroup?: WriteGroup;
}

/** The read-only difference an operator confirms: what moves, what stays, which window is cut and why it cannot be applied. */
export const CutoverReviewView: React.FC<CutoverReviewViewProps> = ({ review, stale, onOpenGroup, candidates, sourceGroup }) => {
  const { t } = useTranslation('workbench-v2');
  const movedNames = review.moved.map((ref) => formatMemberRef(ref, candidates, sourceGroup));
  const unselectedNames = review.unselected.map((ref) => formatMemberRef(ref, candidates, sourceGroup));
  return (
    <div data-testid="cutover-review" className="space-y-3 rounded border border-slate-800 bg-slate-950/40 p-3 text-xs text-slate-300">
      <h5 className="text-xs font-semibold text-slate-100">{t('step4.cutover.review_title')}</h5>
      {stale && <p role="status" data-testid="cutover-review-stale" className="text-amber-300">{t('step4.cutover.stale')}</p>}
      <p>{t('step4.cutover.boundary_reviewed')} <code data-testid="cutover-review-boundary" className="text-slate-100">{review.boundary}</code></p>
      <div className="overflow-x-auto">
        <table className="w-full min-w-[34rem] text-xs">
          <thead>
            <tr className="text-left text-slate-500">
              {['group', 'layout', 'storage', 'period', 'destination', 'schema', 'sources'].map((column) => <th key={column} scope="col" className="pr-3 font-normal">{t(`step4.cutover.col_${column}`)}</th>)}
            </tr>
          </thead>
          <tbody>
            <GroupRow view={review.source} role="source" />
            {review.targets.map((target) => <GroupRow key={target.group_id} view={target} role="target" />)}
          </tbody>
        </table>
      </div>
      <ul className="space-y-1">
        <li data-testid="cutover-moved">
          <span>{t('step4.cutover.moved', { count: review.moved.length })}</span>
          {movedNames.length > 0 && <span className="ml-1 text-slate-200">({movedNames.join(', ')})</span>}
        </li>
        {review.unselected.length > 0 && (
          <li data-testid="cutover-unselected" className="text-amber-300">
            <span>{t('step4.cutover.unselected', { count: review.unselected.length })}</span>
            {unselectedNames.length > 0 && <span className="ml-1 font-medium text-amber-200">({unselectedNames.join(', ')})</span>}
            {' '}<span>{t('step4.cutover.remaining', { count: review.remaining.length })}</span>
          </li>
        )}
      </ul>
      <ul className="list-disc space-y-1 pl-4 text-slate-400">
        <li>{t('step4.cutover.explain_truncation')}</li>
        <li>{t('step4.cutover.explain_drain')}</li>
        <li>{t('step4.cutover.explain_baseline')}</li>
      </ul>
      {review.issues.length > 0 && (
        <div role="alert" data-testid="cutover-issues" className="space-y-2 text-amber-200">
          <p className="font-semibold">{t('step4.cutover.issues_title')}</p>
          <ul className="space-y-1">
            {review.issues.map((code) => {
              const repair = CUTOVER_ISSUE_REPAIR[code];
              return (
                <li key={code} data-testid={`cutover-issue-${code}`}>
                  <span>{t(`step4.cutover.issues.${code}`)}</span>
                  {repair === 'source' && onOpenGroup && <button type="button" data-testid="cutover-repair-source" onClick={() => onOpenGroup(review.source.group_id)} className="ml-2 underline">{t('step4.cutover.repair_source')}</button>}
                  {repair === 'target' && onOpenGroup && review.targets.map((target) => (
                    <button key={target.group_id} type="button" data-testid={`cutover-repair-target-${target.group_id}`} onClick={() => onOpenGroup(target.group_id)} className="ml-2 underline">{t('step4.cutover.repair_target', { name: target.name })}</button>
                  ))}
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
};
