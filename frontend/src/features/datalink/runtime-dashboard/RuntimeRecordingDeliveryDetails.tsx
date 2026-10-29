import { useTranslation } from 'react-i18next';
import type { RecordingGroupView, RuntimeReceipt, SourceResult } from '../../../types/runtimeRuleEvidence';
import { RuntimeSetupFact } from './RuntimeSetupHelpers';

/** Read-only details from the same validated group observation; never allocate group success to a source. */
export function RuntimeRecordingDeliveryDetails({ group, deviceId }: { group: RecordingGroupView; deviceId: string }) {
  const { t } = useTranslation('runtime-dashboard');
  const current = group.delivery?.current;
  const sources = current?.source_state === 'available' ? current.sources.filter((source) => source.device_id === deviceId) : null;
  const testWrite = group.delivery?.test_write;
  const unavailable = t('recordingDetails.unavailable');
  return <section aria-label={t('recordingDetails.title')} className="mt-4 rounded-xl border border-slate-700 p-3">
    <h4 className="font-medium text-slate-100">{t('recordingDetails.title')}</h4>
    <details className="mt-3 rounded border border-slate-800 p-2">
      <summary className="cursor-pointer text-sm text-slate-200 focus-visible:outline focus-visible:outline-sky-400">{t('recordingDetails.groupReceipt')}</summary>
      <p className="mt-2 text-sm text-slate-400">{t('recordingDetails.groupWide')}</p>
      <Receipt receipt={current?.last_commit ?? null} />
    </details>
    <section aria-label={t('recordingDetails.selectedSources', { device: deviceId })} className="mt-4">
      <h5 className="font-medium text-slate-200">{t('recordingDetails.selectedSources', { device: deviceId })}</h5>
      <p className="mt-1 text-sm text-slate-400">{t('recordingDetails.sourceNote')}</p>
      {sources === null ? <p className="mt-2 text-sm text-slate-400">{t('recordingDetails.sourcesUnavailable')}</p> : !sources.length ?
        <p className="mt-2 text-sm text-slate-400">{t('recordingDetails.noSources')}</p> : sources.map((source) => <details key={JSON.stringify([source.device_id, source.point_id, source.tag_id])} className="mt-2 rounded border border-slate-800 p-2">
          <summary className="cursor-pointer text-sm text-slate-200 focus-visible:outline focus-visible:outline-sky-400">{t('recordingDetails.source', { point: source.point_id, tag: source.tag_id })}</summary>
          <dl className="mt-2 grid gap-2 sm:grid-cols-3">
            <RuntimeSetupFact label={t('recordingRules.device')} value={source.device_id} />
            <RuntimeSetupFact label={t('recordingRules.point')} value={source.point_id} />
            <RuntimeSetupFact label={t('recordingRules.tag')} value={source.tag_id} />
          </dl>
          <div className="mt-3 grid gap-3 xl:grid-cols-2">
            <section aria-label={t('recordingDetails.latest')}><h6 className="font-medium text-slate-200">{t('recordingDetails.latest')}</h6><SourceResultDetails result={source.latest} /></section>
            <section aria-label={t('recordingDetails.lastCommitted')}><h6 className="font-medium text-slate-200">{t('recordingDetails.lastCommitted')}</h6><SourceResultDetails result={source.last_committed} /></section>
          </div>
        </details>)}
    </section>
    <section aria-label={t('recordingDetails.testWrite')} className="mt-4">
      <h5 className="font-medium text-slate-200">{t('recordingDetails.testWrite')}</h5>
      <p className="mt-1 text-sm text-slate-400">{t('recordingDetails.independentPhases')}</p>
      <dl className="mt-2 grid gap-3 sm:grid-cols-2">
        <RuntimeSetupFact label={t('recordingDetails.operation')} value={testWrite?.operation_id ?? unavailable} />
        <RuntimeSetupFact label={t('recordingDetails.provenance')} value={testWrite ? t('recordingDetails.testWriteProvenance') : unavailable} />
        <RuntimeSetupFact label={t('recordingDetails.writtenAt')} value={testWrite?.written_at ?? unavailable} />
        <RuntimeSetupFact label={t('recordingDetails.readbackAt')} value={testWrite?.readback_at ?? unavailable} />
        <RuntimeSetupFact label={t('recordingDetails.cleanupAt')} value={testWrite?.cleanup_at ?? unavailable} />
      </dl>
    </section>
  </section>;
}

function SourceResultDetails({ result }: { result: SourceResult | null }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingDetails.unavailable');
  return <>
    <dl className="mt-2 grid gap-2 sm:grid-cols-2">
      <RuntimeSetupFact label={t('recordingDetails.state')} value={result ? t(`recordingDetails.states.${result.state}`) : unavailable} />
      <RuntimeSetupFact label={t('recordingDetails.windowStart')} value={result?.window_start ?? unavailable} />
      <RuntimeSetupFact label={t('recordingDetails.windowEnd')} value={result?.window_end ?? unavailable} />
      <RuntimeSetupFact label={t('recordingDetails.coverage')} value={result?.coverage ? t(`recordingDetails.coverageValues.${result.coverage}`) : unavailable} />
      <RuntimeSetupFact label={t('recordingDetails.delta')} value={result?.delta_known == null ? unavailable : t(result.delta_known ? 'recordingDetails.deltaKnown' : 'recordingDetails.deltaUnknown')} />
    </dl>
    <h6 className="mt-3 text-sm text-slate-300">{t('recordingDetails.sourceReceipt')}</h6>
    <Receipt receipt={result?.receipt ?? null} />
  </>;
}

function Receipt({ receipt }: { receipt: RuntimeReceipt | null }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('recordingDetails.unavailable');
  return <dl className="mt-2 grid gap-2 sm:grid-cols-2">
    <RuntimeSetupFact label={t('recordingDetails.record')} value={receipt?.record_id ?? unavailable} />
    <RuntimeSetupFact label={t('recordingDetails.effect')} value={receipt?.effect_key ?? unavailable} />
    <RuntimeSetupFact label={t('recordingDetails.digest')} value={receipt?.payload_digest ?? unavailable} />
    <RuntimeSetupFact label={t('recordingDetails.groupRevision')} value={receipt?.group_revision ?? unavailable} />
    <RuntimeSetupFact label={t('recordingDetails.connectorRevision')} value={receipt?.connector_revision ?? unavailable} />
    <RuntimeSetupFact label={t('recordingDetails.committedAt')} value={receipt?.committed_at ?? unavailable} />
    <RuntimeSetupFact label={t('recordingDetails.provenance')} value={receipt ? t('recordingDetails.runtime') : unavailable} />
  </dl>;
}
