import { useTranslation } from 'react-i18next';
import type { CurrentDelivery, HistoricalBacklog, RecordingGroupView } from '../../../types/runtimeRuleEvidence';
import { actionForEvidenceReason, parseEvidenceReason } from '../../../services/runtimeRuleEvidenceParser';
import { RuntimeRecordingResourceEvidence } from './RuntimeRecordingResourceEvidence';
import { RuntimeSetupFact } from './RuntimeSetupHelpers';

/** Presentation of the existing validated delivery projection; never a delivery authority. */
export function RuntimeRecordingDeliveryCard({ group }: { group: RecordingGroupView }) {
  const { t } = useTranslation('runtime-dashboard');
  const delivery = group.delivery;
  return <section aria-label={t('recordingDelivery.title')} className="mt-4 rounded-xl border border-slate-700 p-3">
    <h4 className="font-medium text-slate-100">{t('recordingDelivery.title')}</h4>
    {!delivery ? <p className="mt-2 text-sm text-slate-400">{t('recordingDelivery.deliveryUnavailable')}</p> : <>
      <section aria-label={t('recordingDelivery.currentTitle')} className="mt-3">
        <h5 className="font-medium text-slate-200">{t('recordingDelivery.currentTitle')}</h5>
        {delivery.current ? <CurrentStages current={delivery.current} /> : <p className="mt-2 text-sm text-slate-400">{t('recordingDelivery.currentUnavailable')}</p>}
      </section>
      <section aria-label={t('recordingDelivery.backlogTitle')} className="mt-4">
        <h5 className="font-medium text-slate-200">{t('recordingDelivery.backlogTitle')}</h5>
        <p className="mt-1 text-sm text-slate-400">{t('recordingDelivery.backlogNote')}</p>
        {delivery.backlog.length ? delivery.backlog.map((backlog) => <Backlog key={JSON.stringify([backlog.binding.entity_revision, backlog.destination.connector_id, backlog.destination.connector_revision, backlog.destination.schema, backlog.destination.table])} backlog={backlog} />) : <p className="mt-2 text-sm text-slate-400">{t('recordingDelivery.noBacklog')}</p>}
      </section>
      <section aria-label={t('recordingDelivery.replayTitle')} className="mt-4">
        <h5 className="font-medium text-slate-200">{t('recordingDelivery.replayTitle')}</h5>
        <p className="mt-1 text-sm text-slate-400">{t('recordingDelivery.replayNote')}</p>
        {delivery.counter_replay_quarantines === null ? <p className="mt-2 text-sm text-slate-400">{t('recordingDelivery.replayUnavailable')}</p> :
          delivery.counter_replay_quarantines.length ? delivery.counter_replay_quarantines.map((replay) => <details key={replay.group_revision} className="mt-2 rounded border border-slate-800 p-2">
            <summary className="cursor-pointer text-sm text-slate-200 focus-visible:outline focus-visible:outline-sky-400">{t('recordingDelivery.replayRevision', { revision: replay.group_revision })}</summary>
            <dl className="mt-2 grid gap-2 sm:grid-cols-2">
              <RuntimeSetupFact label={t('recordingDelivery.group')} value={replay.group_id} />
              <RuntimeSetupFact label={t('recordingDelivery.revision')} value={replay.group_revision} />
              <RuntimeSetupFact label={t('recordingDelivery.pendingSamples')} value={t('recordingDelivery.samples', { count: replay.pending_samples })} />
            </dl>
            <Reason reason={replay.reason} guidance={false} />
            <p className="mt-2 text-sm text-amber-100">{t('recordingDelivery.preserveInputs')}</p>
          </details>) : <p className="mt-2 text-sm text-slate-400">{t('recordingDelivery.noReplay')}</p>}
      </section>
    </>}
    <RuntimeRecordingResourceEvidence current={delivery?.current ?? null} />
  </section>;
}

function CurrentStages({ current }: { current: CurrentDelivery }) {
  const { t } = useTranslation('runtime-dashboard');
  const outputStages = ['queued', 'retrying', 'blocked', 'quarantined', 'unknown', 'sql_committed', 'skipped'] as const;
  return <>
    <p className="mt-1 text-sm text-slate-300">{t(current.shared ? 'recordingDelivery.shared' : 'recordingDelivery.groupTotals')}</p>
    <p className="mt-1 text-sm text-slate-400">{t('recordingDelivery.unitsNote')}</p>
    <dl className="mt-3 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      <RuntimeSetupFact label={t('recordingDelivery.revision')} value={current.binding.entity_revision} />
      <RuntimeSetupFact label={t('recordingDelivery.intakeState')} value={current.intake ? t(`recordingDelivery.intake.${current.intake.state}`) : t('recordingDelivery.unavailable')} />
      <RuntimeSetupFact label={t('recordingDelivery.intakeRevision')} value={current.intake?.group_revision ?? t('recordingDelivery.unavailable')} />
      <RuntimeSetupFact label={t('recordingDelivery.stages.collecting')} value={t('recordingDelivery.samples', { count: current.stages.collecting })} />
      {outputStages.map((stage) => <RuntimeSetupFact key={stage} label={t(`recordingDelivery.stages.${stage}`)} value={t('recordingDelivery.rows', { count: current.stages[stage] })} />)}
      <RuntimeSetupFact label={t('recordingDelivery.outputReason')} value={t('recordingDelivery.unavailable')} />
    </dl>
    {current.intake?.reason && <section aria-label={t('recordingDelivery.intakeReason')} className="mt-2">
      <h6 className="text-sm text-slate-300">{t('recordingDelivery.intakeReason')}</h6><Reason reason={current.intake.reason} />
    </section>}
  </>;
}

function Backlog({ backlog }: { backlog: HistoricalBacklog }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingDelivery.unavailable');
  const destination = backlog.destination;
  return <details className="mt-2 rounded border border-slate-800 p-2">
    <summary className="cursor-pointer text-sm text-slate-200 focus-visible:outline focus-visible:outline-sky-400">{t('recordingDelivery.backlogRevision', { revision: backlog.binding.entity_revision })}</summary>
    <dl className="mt-2 grid gap-2 sm:grid-cols-2">
      <RuntimeSetupFact label={t('recordingDelivery.revision')} value={backlog.binding.entity_revision} />
      <RuntimeSetupFact label={t('recordingRules.connectorId')} value={destination.connector_id} />
      <RuntimeSetupFact label={t('recordingRules.connectorRevision')} value={destination.connector_revision} />
      <RuntimeSetupFact label={t('recordingRules.alias')} value={destination.alias ?? unavailable} />
      <RuntimeSetupFact label={t('recordingRules.kind')} value={destination.kind === 'unavailable' ? unavailable : destination.kind} />
      <RuntimeSetupFact label={t('recordingRules.schema')} value={destination.schema ?? unavailable} />
      <RuntimeSetupFact label={t('recordingRules.table')} value={destination.table} />
      <RuntimeSetupFact label={t('recordingDelivery.pendingRows')} value={t('recordingDelivery.rows', { count: backlog.pending })} />
    </dl>
    {backlog.reasons.length ? <ul className="mt-2 space-y-2">{backlog.reasons.map((reason, index) => <li key={index}><Reason reason={reason} /></li>)}</ul> : <p className="mt-2 text-sm text-slate-400">{t('recordingDelivery.reasons.unavailable')}</p>}
  </details>;
}

function Reason({ reason, guidance = true }: { reason: string; guidance?: boolean }) {
  const { t } = useTranslation('runtime-dashboard');
  const safeReason = parseEvidenceReason(reason);
  return <div className="mt-2 text-sm text-slate-300">
    <p>{t(`recordingDelivery.reasons.${safeReason}`)}</p>
    {guidance && <p className="mt-1 text-slate-400">{t(`recordingDelivery.guidance.${actionForEvidenceReason(safeReason)}`)}</p>}
  </div>;
}
