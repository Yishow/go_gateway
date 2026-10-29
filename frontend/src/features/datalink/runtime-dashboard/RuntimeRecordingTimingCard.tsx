import { useTranslation } from 'react-i18next';
import type { RecordingGroupView } from '../../../types/runtimeRuleEvidence';
import { RuntimeSetupFact } from './RuntimeSetupHelpers';

/** Snapshot timing only: no local clock, synthesized deadline or inferred absence of pending work. */
export function RuntimeRecordingTimingCard({ group }: { group: RecordingGroupView }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingTiming.unavailable');
  const age = group.delivery?.oldest_pending_seconds;
  const current = group.delivery?.current;
  const duration = current?.attempt_duration_ms;
  return <section aria-label={t('recordingTiming.title')} className="mt-4 rounded-xl border border-slate-700 p-3">
    <h4 className="font-medium text-slate-100">{t('recordingTiming.title')}</h4>
    <p className="mt-1 text-sm text-slate-400">{t('recordingTiming.observedNote')}</p>
    <p className="mt-1 text-sm text-slate-400">{t('recordingTiming.scopeNote')}</p>
    <dl className="mt-3 grid gap-3 sm:grid-cols-2">
      <RuntimeSetupFact label={t('recordingTiming.oldestAge')} value={age == null ? unavailable : t('recordingTiming.seconds', { value: age })} />
      <RuntimeSetupFact label={t('recordingTiming.nextDue')} value={current?.next_due_at ?? unavailable} />
      <RuntimeSetupFact label={t('recordingTiming.nextRetry')} value={current?.next_retry_at ?? unavailable} />
      <RuntimeSetupFact label={t('recordingTiming.duration')} value={duration == null ? unavailable : t('recordingTiming.milliseconds', { value: duration })} />
    </dl>
  </section>;
}
