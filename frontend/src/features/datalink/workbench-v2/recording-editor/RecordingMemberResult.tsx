import type { GroupCandidate } from '../state/writeGroup/candidates';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { WriteGroupDelivery, WriteGroupSourceEvidenceItem } from '../../../../types/studioV2WriteGroupDelivery';
export interface RecordingMemberResultProps { group: WriteGroup; delivery?: WriteGroupDelivery; candidates?: readonly GroupCandidate[] }
/** Only the canonical frozen result projection is displayed; no client reducer, window arithmetic or receipt inference. */
export function RecordingMemberResult({ group, delivery, candidates = [] }: RecordingMemberResultProps) {
  const { t } = useTranslation('recording-editor');
  const evidence = delivery?.group_id === group.id ? delivery.source_evidence : undefined;
  if (group.recording_profile?.kind !== 'mixed') return null;
  const available = evidence?.state === 'available' && Boolean(group.applied_revision) && evidence.group_revision === group.applied_revision;
  const sourceLabel = (source: WriteGroupSourceEvidenceItem) => {
    const found = candidates.filter(candidate => candidate.device_id === source.device_id && candidate.point_id === source.point_id && candidate.tag_id === source.tag_id);
    return found.length === 1 ? `${found[0].device_name} · ${found[0].label} · ${found[0].address}` : t('result.sourceUnavailable');
  };
  const result = (source: WriteGroupSourceEvidenceItem, kind: 'latest' | 'confirmed') => {
    const value = kind === 'latest' ? source.latest : source.last_sql_committed;
    if (!value) return null;
    if (!value.available || value.profile_kind !== 'mixed' || value.group_revision !== group.applied_revision ||
      value.connector_revision !== evidence?.connector_revision) return <p key={kind}>{t('result.unavailable')}</p>;
    const calculation = value.calculation;
    const committed = value.delivery === 'sql_committed' && Boolean(value.committed_at);
    const window = calculation.status !== 'not_due' && calculation.aligned_start && calculation.aligned_end;
    return <div key={kind} className="space-y-1" data-testid={`member-result-${source.point_id}-${kind}`}>
      <h5>{t(`result.${kind}`)}</h5>
      <p>{t(`calculation.${calculation.method}`)} · {calculation.interval_seconds} {t('calculation.second')} · {t(`result.${calculation.status}`)}</p>
      <p>{t(calculation.result_present ? 'result.present' : 'result.null')} · {t(`result.${calculation.outcome_kind}`)}</p>
      {calculation.reason === 'first-reading' && <p>{t('result.firstBaseline')}</p>}
      {calculation.reason && !['first-reading', 'not_due'].includes(calculation.reason) && <p>{t('result.reason')}: {t(`workbench-v2:step4.basic.status.reason.${calculation.reason}`, { defaultValue: t('result.reasonUnknown') })}</p>}
      <p>{t('result.due')}: <time>{calculation.due_at}</time></p>
      {calculation.next_due && <p>{t('result.nextDue')}: <time>{calculation.next_due}</time></p>}
      {window && <p data-testid={`member-result-window-${source.point_id}-${kind}`}>{t('result.aligned')}: <time>{calculation.aligned_start}</time> → <time>{calculation.aligned_end}</time> UTC</p>}
      {calculation.actual_start && calculation.actual_end && <p>{t('result.actual')}: <time>{calculation.actual_start}</time> → <time>{calculation.actual_end}</time> UTC</p>}
      <p>{calculation.source_type} → {calculation.result_type}{calculation.rounded && ` · ${t('result.rounded')}`}</p>
      <p>{t(calculation.raw_present ? 'result.raw' : 'result.noRaw')}</p>
      <p>{t(committed ? 'result.committed' : `result.delivery.${value.delivery}`)}</p>
    </div>;
  };
  return <section className="space-y-2 text-xs text-slate-400" data-testid="recording-member-results">
    <h4 className="text-sm font-medium text-slate-200">{t('result.title')}</h4>
    {!available ? <p>{t('result.unavailable')}</p> : evidence.sources.map(source => <div key={JSON.stringify([source.device_id, source.point_id, source.tag_id])} className="space-y-2 border-t border-slate-800 pt-2">
      <h5 className="font-medium text-slate-200">{sourceLabel(source)}</h5>{result(source, 'latest')}{result(source, 'confirmed')}
    </div>)}
  </section>;
}
