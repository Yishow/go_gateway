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
export const GroupDeliveryStrip: React.FC<{ group: WriteGroup; readonly?: boolean }> = ({ group, readonly = false }) => {
  const { t } = useTranslation('workbench-v2');
  const delivery = useWriteGroupDeliveryQuery(group.id, true, 5000);
  const confirmed = delivery.isError ? undefined : delivery.data;
  const stages = confirmed?.stages;
  const recentIssues = confirmed?.recent_bucket_issues ?? [];
  const cell = (label: string, value: string, testId: string) => (
    <div key={testId} className="min-w-0">
      <dt className="text-[11px] text-slate-500">{label}</dt>
      <dd className="text-sm text-slate-200" data-testid={testId}>{value}</dd>
    </div>
  );
  const unknown = t('step4.group.delivery.unconfirmed');
  const num = (value: number | undefined) => (value === undefined ? unknown : String(value));
  const causeLabel = (cause: WriteGroupDeliveryBucketCause) => {
    const safeCause = SAFE_BUCKET_CAUSES.has(cause) ? cause : 'unavailable';
    return t(`step4.group.delivery.cause.${safeCause}`);
  };
  return (
    <section aria-label={t('step4.group.delivery.title')} className="space-y-2" data-testid="group-delivery-strip">
      <h4 className="text-sm font-semibold text-slate-200">{t('step4.group.delivery.title')}</h4>
      {delivery.isError && <p role="alert" className="text-xs text-amber-300" data-testid="group-delivery-failed">{t('step4.group.delivery.failed')}</p>}
      <dl className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        {cell(t('step4.group.delivery.saved'), group.status === 'deleted' ? t('step4.group.status.deleted') : t('step4.group.delivery.yes'), 'delivery-saved')}
        {cell(t('step4.group.delivery.applied'), group.applied_revision ? t('step4.group.delivery.yes') : t('step4.group.delivery.no'), 'delivery-applied')}
        {cell(t('step4.group.delivery.collecting'), num(stages?.collecting), 'delivery-collecting')}
        {cell(t('step4.group.delivery.skipped'), num(confirmed?.skipped_buckets), 'delivery-skipped-buckets')}
        {cell(t('step4.group.delivery.no_data'), num(confirmed?.no_data_buckets), 'delivery-no-data-buckets')}
        {cell(t('step4.group.delivery.local'), stages ? String(stages.queued + stages.retrying) : unknown, 'delivery-local')}
        {cell(t('step4.group.delivery.attention'), stages ? String(stages.blocked + stages.quarantined + stages.unknown) : unknown, 'delivery-attention')}
        {cell(t('step4.group.delivery.committed'), num(stages?.sql_committed), 'delivery-committed')}
      </dl>
      {recentIssues.length > 0 && (
        <div data-testid="delivery-bucket-causes" className="space-y-1">
          <p className="text-xs font-medium text-slate-300">{t('step4.group.delivery.recent_issues')}</p>
          <ul className="space-y-1 text-xs text-slate-400">
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
      {group.applied_revision && <details className="text-[11px] text-slate-500"><summary>{t('step4.group.delivery.recovery.diagnostics')}</summary>{group.applied_revision}</details>}
      <GroupDeliveryRecovery groupId={group.id} items={confirmed?.attention ?? []} readonly={readonly} onResolved={() => void delivery.refetch()} />
      <p className="text-[11px] text-slate-500">{t('step4.group.delivery.note')}</p>
    </section>
  );
};
