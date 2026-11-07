import { useTranslation } from 'react-i18next';
import { CircleHelp, History, Hourglass } from 'lucide-react';
import type { recordingHomeResultSummary } from './recordingHomeResultSummary';
import { RecordingHomeWriteStages } from './RecordingHomeWriteStages';
/** Compact result facts; secondary timestamps/owner/readback details are opt-in. */
export function RecordingHomeResultFacts({ id, name, summary, showSourceState = true }: { id: string; name: string; showSourceState?: boolean; summary: ReturnType<typeof recordingHomeResultSummary> }) {
  const { t } = useTranslation('recording-home');
  const elapsed = summary.lastCommittedAt ? (Date.now() - Date.parse(summary.lastCommittedAt)) / 1000 : NaN;
  const relative = !Number.isFinite(elapsed) || elapsed < 0 ? summary.lastCommittedAt : elapsed < 60 ? t('rowEvidence.secondsAgo', { count: Math.floor(elapsed) }) :
    elapsed < 3600 ? t('rowEvidence.minutesAgo', { count: Math.floor(elapsed / 60) }) :
      elapsed < 86400 ? t('rowEvidence.hoursAgo', { count: Math.floor(elapsed / 3600) }) : t('rowEvidence.daysAgo', { count: Math.floor(elapsed / 86400) });
  return <div data-testid={`home-result-${id}`} className="gw-result-facts">
    <RecordingHomeWriteStages name={name} summary={summary} />
    {/* Receipts are historical SQL facts; no icon here claims current success. */}
    <p className="gw-home-sql-fact" data-sql-fact={summary.sql}>{(summary.sql === 'unknown' ? <CircleHelp aria-hidden="true" size={14} /> : summary.sql === 'awaiting_first_write'
      ? <Hourglass aria-hidden="true" size={14} /> : <History aria-hidden="true" size={14} />)}<span><span>{t(`result.sql.${summary.sql}`)}</span>{summary.lastCommittedAt && <> · <time dateTime={summary.lastCommittedAt} title={summary.lastCommittedAt}>{relative}</time></>}</span></p>
    {(showSourceState || summary.current === 'delivery_attention') && <p>{t(`result.current.${summary.current}`)}</p>}
    {showSourceState && summary.configuration === 'not_requested' && <p>{t('result.configuration.not_requested')}</p>}
    {summary.effective === 'requested_not_effective' && <p>{t('result.effective.requested_not_effective')}</p>}
  </div>;
}
