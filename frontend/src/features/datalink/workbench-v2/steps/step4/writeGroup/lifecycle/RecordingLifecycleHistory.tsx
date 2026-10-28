import { useTranslation } from 'react-i18next';
import type { WriteGroupLifecycleHistory } from '@/types/studioV2WriteGroupLifecycle';

/** Counts retain their source units; samples and delivery rows must never be added. */
export function RecordingLifecycleHistory({ history }: { history?: WriteGroupLifecycleHistory }) {
  const { t } = useTranslation('recording-lifecycle');
  if (!history) return <p role="status">{t('history.unavailable')}</p>;
  return <div className="space-y-2" data-testid="recording-lifecycle-history">
    <p>{t('history.kept')}</p>
    {history.generations.map((generation) => <section key={generation.group_id} className="space-y-1 rounded border border-slate-700 p-2">
      <h4>{t('history.generation', { number: generation.generation ?? '—' })} · {generation.group_id}</h4>
      <p>{t('history.saved_destination')}: {generation.saved_destination.connector_id} · {generation.saved_destination.connector_revision}</p>
      {generation.revisions.length === 0 && <p>{t('history.no_applied')}</p>}
      {generation.revisions.map((revision) => <div key={revision.group_revision} className="space-y-1">
        <p>{t('history.revision')} · {revision.group_revision}</p>
        <p>{t('history.original_destination')}: {revision.destination.connector_id} · {revision.destination.connector_revision}</p>
        <dl className="grid grid-cols-2 gap-x-3">
          <dt>{t('history.samples')}</dt><dd data-testid="lifecycle-pending-samples">{revision.pending.samples}</dd>
          <dt>{t('history.outbox')}</dt><dd data-testid="lifecycle-pending-outbox">{revision.pending.outbox.outstanding}</dd>
          <dt>{t('history.blocked')}</dt><dd data-testid="lifecycle-pending-blocked">{revision.pending.outbox.blocked}</dd>
          <dt>{t('history.unknown')}</dt><dd data-testid="lifecycle-pending-unknown">{revision.pending.outbox.unknown}</dd>
          <dt>{t('history.quarantined')}</dt><dd>{revision.pending.outbox.quarantined}</dd>
        </dl>
        {revision.pending.handoffs.map((handoff, index) => <p key={`${handoff.origin.group_id}-${handoff.origin.group_revision}-${index}`} data-testid="lifecycle-handoff">
          {t('history.handoff', { count: handoff.samples })} · {handoff.origin.group_id}/{handoff.origin.group_revision} → {handoff.target.group_id}/{handoff.target.group_revision}
        </p>)}
      </div>)}
    </section>)}
  </div>;
}
