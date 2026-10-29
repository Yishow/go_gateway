import { useTranslation } from 'react-i18next';
import type { CurrentDelivery, StorageEvidence } from '../../../types/runtimeRuleEvidence';
import { RuntimeSetupFact } from './RuntimeSetupHelpers';

/** Existing validated read evidence only: no quota inference, repair authority or additional reader. */
export function RuntimeRecordingResourceEvidence({ current }: { current: CurrentDelivery | null }) {
  const { t } = useTranslation('runtime-dashboard');
  const quota = current?.quota;
  const unavailable = t('recordingResources.unavailable');
  return <section aria-label={t('recordingResources.title')} className="mt-4 rounded border border-slate-800 p-3">
    <h5 className="font-medium text-slate-200">{t('recordingResources.title')}</h5>
    <p className="mt-1 text-sm text-slate-400">{t('recordingResources.quotaNote')}</p>
    {!quota && <p className="mt-2 text-sm text-slate-400">{t('recordingResources.quotaUnavailable')}</p>}
    <dl className="mt-3 grid gap-3 sm:grid-cols-2">
      <RuntimeSetupFact label={t('recordingResources.quotaScope')} value={quota ? t(`recordingResources.scopeValues.${quota.scope}`) : unavailable} />
      <RuntimeSetupFact label={t('recordingResources.quotaState')} value={quota ? t(`recordingResources.quotaStates.${quota.state}`) : unavailable} />
      <RuntimeSetupFact label={t('recordingResources.usedBytes')} value={quota ? String(quota.used_bytes) : unavailable} />
      <RuntimeSetupFact label={t('recordingResources.maxBytes')} value={quota ? String(quota.max_bytes) : unavailable} />
      <RuntimeSetupFact label={t('recordingResources.refused')} value={quota ? t(`recordingResources.booleans.${quota.intake_refused ? 'yes' : 'no'}`) : unavailable} />
    </dl>
    <StorageHistory storage={current?.storage ?? null} />
  </section>;
}

function StorageHistory({ storage }: { storage: StorageEvidence | null }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingResources.unavailable');
  return <section aria-label={t('recordingResources.storageTitle')} className="mt-4">
    <h6 className="font-medium text-slate-200">{t('recordingResources.storageTitle')}</h6>
    <p className="mt-1 text-sm text-slate-400">{t('recordingResources.notPermission')}</p>
    {!storage ? <p className="mt-2 text-sm text-slate-400">{t('recordingResources.storageUnavailable')}</p> : <>
      <p className="mt-1 text-sm text-slate-400">{t('recordingResources.storageScope')}</p>
      <dl className="mt-3 grid gap-3 sm:grid-cols-2">
        <RuntimeSetupFact label={t('recordingResources.generation')} value={String(storage.generation)} />
        <RuntimeSetupFact label={t('recordingResources.historyLoss')} value={t(`recordingResources.booleans.${storage.history_may_be_lost ? 'yes' : 'no'}`)} />
        <RuntimeSetupFact label={t('recordingResources.reconstructions')} value={String(storage.reconstructions)} />
        <RuntimeSetupFact label={t('recordingResources.lastReconstructed')} value={storage.last_reconstructed_at ?? unavailable} />
        <RuntimeSetupFact label={t('recordingResources.repairState')} value={t(`recordingResources.repairStates.${storage.repair_state}`)} />
      </dl>
      <h6 className="mt-3 text-sm text-slate-300">{t('recordingResources.events')}</h6>
      {storage.events.length === 0 ? <p className="mt-2 text-sm text-slate-400">{t('recordingResources.noEvents')}</p> : <ul className="mt-2 space-y-2 text-sm text-slate-300">
        {storage.events.map((event, index) => <li key={JSON.stringify([event.kind, event.at, index])} className="rounded border border-slate-800 p-2">
          <span>{t(`recordingResources.eventKinds.${event.kind}`)}</span>
          <dl className="mt-1"><RuntimeSetupFact label={t('recordingResources.eventAt')} value={event.at ?? unavailable} /></dl>
        </li>)}
      </ul>}
    </>}
  </section>;
}
