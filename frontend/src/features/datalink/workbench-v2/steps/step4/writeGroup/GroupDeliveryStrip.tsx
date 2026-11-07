import { GroupDeliveryRecovery } from './GroupDeliveryRecovery';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupDeliveryQuery } from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroupDeliveryBucketCause } from '../../../../../../types/studioV2WriteGroupDelivery';
import type { WriteGroup } from '../../../../../../types/studioV2WriteGroup';

const SAFE_BUCKET_CAUSES = new Set<WriteGroupDeliveryBucketCause>([
  'missing', 'bad', 'stale', 'invalid', 'no_data', 'unavailable',
]);

/**
 * Where the group's accepted data stands, one fact per cell: saved and applied
 * revisions come from the group, the rest from the delivery status. Only
 * `sql_committed` is destination evidence; nothing is summed into a success.
 */
export const GroupDeliveryStrip: React.FC<{ group: WriteGroup; readonly?: boolean; compact?: boolean }> = ({ group, readonly = false, compact = false }) => {
  const { t } = useTranslation('workbench-v2');
  const { t: summaryText } = useTranslation('recording-home');
  const delivery = useWriteGroupDeliveryQuery(group.id, true, 5000);
  const confirmed = delivery.isError || delivery.data?.group_id !== group.id ? undefined : delivery.data;
  const stages = confirmed?.stages;
  const pending = stages ? stages.queued + stages.retrying : undefined;
  const committed = stages?.sql_committed;
  const hasHistory = Boolean(committed || confirmed?.last_sql_committed_at);
  const intakeLabel = !confirmed ? 'unconfirmed' : confirmed.intake.state === 'active' ? 'active' : confirmed.intake.state === 'not_running' ? 'notRunning' : confirmed.intake.state === 'blocked' || confirmed.intake.state === 'retiring' ? 'blocked' : 'unconfirmed';
  const recentIssues = confirmed?.recent_bucket_issues ?? [];
  const destinationUnavailable = confirmed?.backlog.some(item => item.pending > 0 && item.error_codes.includes('target-unavailable'));
  const cell = (label: string, value: string, testId: string) => (
    <div key={testId} className="min-w-0">
      <dt className="text-xs text-[var(--gw-muted)]">{label}</dt>
      <dd className="text-sm text-[var(--gw-text)]" data-testid={testId}>{value}</dd>
    </div>
  );
  const unknown = t('step4.group.delivery.unconfirmed');
  const num = (value: number | undefined) => (value === undefined ? unknown : String(value));
  const causeLabel = (cause: WriteGroupDeliveryBucketCause) => {
    const safeCause = SAFE_BUCKET_CAUSES.has(cause) ? cause : 'unavailable';
    return t(`step4.group.delivery.cause.${safeCause}`);
  };
  const facts = (
    <dl className={compact ? 'recording-kpi-strip' : 'grid grid-cols-2 gap-3 sm:grid-cols-4'}>
        {cell(t('step4.group.delivery.saved'), group.status === 'deleted' ? t('step4.group.status.deleted') : t('step4.group.delivery.yes'), 'delivery-saved')}
        {cell(t('step4.group.delivery.applied'), group.applied_revision ? t('step4.group.delivery.yes') : t('step4.group.delivery.no'), 'delivery-applied')}
        {cell(t('step4.group.delivery.collecting'), num(stages?.collecting), 'delivery-collecting')}
        {cell(t('step4.group.delivery.skipped'), num(confirmed?.skipped_buckets), 'delivery-skipped-buckets')}
        {cell(t('step4.group.delivery.no_data'), num(confirmed?.no_data_buckets), 'delivery-no-data-buckets')}
        {cell(t('step4.group.delivery.local'), stages ? String(stages.queued + stages.retrying) : unknown, 'delivery-local')}
        {cell(t('step4.group.delivery.attention'), stages ? String(stages.blocked + stages.quarantined + stages.unknown) : unknown, 'delivery-attention')}
        {cell(t('step4.group.delivery.committed'), num(stages?.sql_committed), 'delivery-committed')}
      </dl>
  );
  return (
    <section aria-label={t('step4.group.delivery.title')} className="space-y-2" data-testid="group-delivery-strip">
      <h4 className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.group.delivery.title')}{compact && !confirmed && <span className="font-normal text-[var(--gw-muted)]"> · {unknown}</span>}</h4>
      {delivery.isError && <p role="alert" className="text-xs text-[var(--gw-warning-fg)]" data-testid="group-delivery-failed">{t('step4.group.delivery.failed')}</p>}
      {destinationUnavailable && <p role="alert" className="text-sm text-[var(--gw-warning-fg)]" data-testid="group-destination-unavailable">{summaryText('deliverySummary.destinationUnavailable')}</p>}
      {compact && stages && stages.blocked + stages.quarantined + stages.unknown > 0 && <p role="status" data-testid="group-delivery-attention-summary">{t('step4.group.delivery.attention')}: {stages.blocked + stages.quarantined + stages.unknown}</p>}
      {compact && <div data-testid="recording-delivery-summary" className="space-y-1 text-sm">
        <p>{t('step4.group.delivery.saved')}: {group.status === 'deleted' ? t('step4.group.status.deleted') : t('step4.group.delivery.yes')} · {summaryText(`deliverySummary.${intakeLabel}`)}</p>
        {hasHistory && <p data-testid="recording-delivery-history">{t('step4.group.delivery.committed')}: {num(committed)}{confirmed?.last_sql_committed_at && <> · <time dateTime={confirmed.last_sql_committed_at}>{confirmed.last_sql_committed_at}</time></>}</p>}
        {pending !== undefined && pending > 0 && <p data-testid="recording-delivery-pending">{t('step4.group.delivery.local')}: {pending}</p>}
        {confirmed?.backlog.some(item => item.pending > 0) && <p>{summaryText('deliverySummary.backlog')}</p>}
      </div>}
      {/* A record that never ran shows its one-line summary only; eight zero cells would read as progress. */}
      {compact ? (group.applied_revision || hasHistory || (pending ?? 0) > 0) && <div data-testid="group-delivery-diagnostics">{facts}</div> : facts}
      {recentIssues.length > 0 && (
        <div data-testid="delivery-bucket-causes" className="space-y-1">
          <p className="text-xs font-medium text-[var(--gw-text)]">{t('step4.group.delivery.recent_issues')}</p>
          <ul className="space-y-1 text-xs text-[var(--gw-muted)]">
            {recentIssues.map((issue, index) => (
              <li key={`${issue.bucket_start}-${index}`} data-testid={`delivery-recent-issue-${index}`}>
                <span>{t(`step4.group.delivery.issue.${issue.kind}`)}</span>{' '}
                <time dateTime={issue.bucket_start}>{issue.bucket_start}</time>{': '}
                {issue.causes.map((cause) => causeLabel(cause)).join(', ')}
              </li>
            ))}
          </ul>
        </div>
      )}
      {group.applied_revision && <p className="text-xs text-[var(--gw-muted)]" data-testid="delivery-applied-revision">{t('step4.group.delivery.recovery.diagnostics')}：<span className="font-mono" title={group.applied_revision}>{group.applied_revision}</span></p>}
      <GroupDeliveryRecovery groupId={group.id} items={confirmed?.attention ?? []} readonly={readonly} onResolved={() => void delivery.refetch()} />
      {(!compact || group.applied_revision || hasHistory || (pending ?? 0) > 0) && <p className="text-xs text-[var(--gw-muted)]">{t('step4.group.delivery.note')}</p>}
    </section>
  );
};
