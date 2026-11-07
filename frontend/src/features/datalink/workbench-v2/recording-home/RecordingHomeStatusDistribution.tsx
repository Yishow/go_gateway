import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import { RecordingStatusBadge } from './RecordingHomeBadges';
import { countRecordingHomeStatuses, recordingHomeStatuses } from './recordingHomeStatusCounts';

/** Saved-status counts over every loaded saved record (never the filtered page). It counts saved
 * configuration only, never runtime, delivery or SQL outcomes. */
export function RecordingHomeStatusDistribution({ groups }: { groups: readonly RecordingReadonly<WriteGroup>[] }) {
  const { t } = useTranslation('recording-home');
  const scopeId = useId();
  const totals = countRecordingHomeStatuses(groups);
  const counts = recordingHomeStatuses.map(status => [status, totals[status]] as const).filter(([, count]) => count > 0);
  if (!counts.length) return null;
  return <div className="gw-home-status-distribution">
    {/* Decorative proportion of saved statuses only; the counts below carry the meaning. */}
    <div className="gw-home-status-bar" data-testid="home-status-bar" aria-hidden="true">{counts.map(([status, count]) =>
      <span key={status} data-status-segment={status} style={{ flexGrow: count }} />)}</div>
    <p id={scopeId} className="gw-home-status-scope">{t('statusDistributionScope', { count: groups.length })}</p>
    <ul aria-label={t('statusDistribution')} aria-describedby={scopeId}>{counts.map(([status, count]) =>
      <li key={status} data-status-count={status}><RecordingStatusBadge status={status} /><span className="gw-home-status-count">{t('statusCount', { count })}</span></li>)}</ul>
  </div>;
}
