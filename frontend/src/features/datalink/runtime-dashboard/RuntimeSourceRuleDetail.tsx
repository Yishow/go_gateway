import { RuntimeSettingsLink } from './RuntimeSettingsLink';
import { useTranslation } from 'react-i18next';
import { useRuntimeSourceEvidence } from '../../../hooks/datalink/useRuntimeSourceEvidence';
import type { EvidenceScope, SourceFacts, SafeTransformStep } from '../../../types/runtimeRuleEvidence';
import type { RuntimeSourceSelection } from '../../../types/studioV2RuntimeContext';
import { RuntimeSetupFact } from './RuntimeSetupHelpers';

interface Props { scope: EvidenceScope; selection: RuntimeSourceSelection; onClose: () => void; navigateTo?: (target: string) => void }

/** Single explicit source selection. No source snapshot is inferred from a healthy collector. */
export function RuntimeSourceRuleDetail({ scope, selection, onClose, navigateTo }: Props) {
  const { t } = useTranslation('runtime-dashboard');
  const query = useRuntimeSourceEvidence(scope, selection);
  const result = query.data;
  const view = !query.isFetching && result?.ok ? result.value : null;
  const source = view?.source;
  const locator = { workspace_id: scope.workspace_id, device_id: scope.device_id, rule_id: null, mapping_id: null, group_id: null, group_revision: null };
  const canLocate = view?.freshness === 'fresh' && source?.rule_id === selection.rule_id && source.point_id === selection.point_id && source.saved_revision === selection.rule_revision && (source.mapping_id ?? undefined) === selection.mapping_id;
  const unavailable = t('sourceDetail.unavailable');
  const recording = view?.groups.filter((group) => group.state === 'effective' && group.proof === 'applied-snapshot-and-layout' && source?.tag_id &&
    group.effective?.members.some((member) => member.device_id === scope.device_id && member.point_id === source.point_id && member.tag_id === source.tag_id)) ?? [];
  const failure = result && !result.ok ? result.reason : 'unavailable';
  return (
    <section aria-label={t('sourceDetail.title')} className="rounded-2xl border border-cyan-800 bg-slate-950/70 p-4 text-xs">
      <div className="flex items-center justify-between gap-3">
        <h3 className="font-semibold text-slate-100">{t('sourceDetail.title')}</h3>
        <button type="button" onClick={onClose} className="rounded border border-slate-700 px-2 py-1 text-slate-200">{t('sourceDetail.close')}</button>
      </div>
      {query.isPending || query.isFetching ? <p role="status" className="mt-3 text-slate-300">{t('sourceDetail.loading')}</p> : !source ? (
        <div role="alert" className="mt-3 text-amber-100">
          <p>{t(`sourceDetail.errors.${failure}`)}</p>
          <button type="button" onClick={() => void query.refetch()} className="mt-2 rounded border border-amber-700 px-2 py-1">{t('sourceDetail.retry')}</button>
        </div>
      ) : (
        <>
          {canLocate && <div className="mt-3 flex flex-wrap gap-2">
            <RuntimeSettingsLink locator={{ ...locator, step: 2, rule_id: source.rule_id }} navigateTo={navigateTo}>{t('settingsLink.rule')}</RuntimeSettingsLink>
            {source.mapping_id && <RuntimeSettingsLink locator={{ ...locator, step: 3, mapping_id: source.mapping_id }} navigateTo={navigateTo}>{t('settingsLink.mapping')}</RuntimeSettingsLink>}
          </div>}
          <dl className="mt-3 grid gap-2 sm:grid-cols-2">
            <RuntimeSetupFact label={t('sourceDetail.rule')} value={source.rule_id} />
            <RuntimeSetupFact label={t('sourceDetail.point')} value={source.point_id} />
            <RuntimeSetupFact label={t('sourceDetail.mapping')} value={source.mapping_id ?? unavailable} />
            <RuntimeSetupFact label="Tag" value={source.tag_id ?? unavailable} />
            <RuntimeSetupFact label={t('sourceDetail.savedRevision')} value={source.saved_revision ?? unavailable} />
          </dl>
          <h4 className="mt-4 font-semibold text-amber-100">{t('sourceDetail.savedOnly')}</h4>
          <p className="mt-1 text-slate-400">{t('sourceDetail.savedNote')}</p>
          <SourceFactsDetail facts={source.saved} />
          <h4 className="mt-4 font-semibold text-slate-100">{t('sourceDetail.effectiveTitle')}</h4>
          {source.state === 'effective' && source.proof === 'immutable-source-facts' && source.effective ? <SourceFactsDetail facts={source.effective} /> : <p className="mt-1 text-slate-400">{t('sourceDetail.effectiveUnavailable')}</p>}
          <h4 className="mt-4 font-semibold text-slate-100">{t('sourceDetail.recording')}</h4>
          <p className="mt-1 text-slate-400">{t('sourceDetail.recordingNote')}</p>
          {recording.length ? <ul className="mt-2 space-y-2 text-slate-200">{recording.map((group) => <li key={group.group_id}>
            <span className="mr-2 font-mono">{group.group_id}</span>
            <span>{group.effective?.interval_seconds != null && group.effective.reduction !== 'unavailable' ? t('sourceDetail.recordingInterval', { reduction: group.effective.reduction, seconds: group.effective.interval_seconds }) : unavailable}</span>
          </li>)}</ul> : <p className="mt-2 text-slate-400">{unavailable}</p>}
        </>
      )}
    </section>
  );
}

function SourceFactsDetail({ facts }: { facts: SourceFacts | null }) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('sourceDetail.unavailable');
  if (!facts) return <p className="mt-2 text-slate-400">{unavailable}</p>;
  const decoding = facts.transforms.filter((step) => step.kind === 'byte-order');
  const transforms = facts.transforms.filter((step) => step.kind !== 'byte-order');
  return <>
    <dl className="mt-2 grid gap-2 sm:grid-cols-2">
      <RuntimeSetupFact label={t('sourceDetail.enabled')} value={t(facts.enabled ? 'sourceDetail.enabledYes' : 'sourceDetail.enabledNo')} />
      <RuntimeSetupFact label={t('sourceDetail.address')} value={facts.address ?? unavailable} />
      <RuntimeSetupFact label={t('sourceDetail.function')} value={facts.protocol_function ?? unavailable} />
      <RuntimeSetupFact label={t('sourceDetail.readType')} value={facts.read_type ?? unavailable} />
      <RuntimeSetupFact label={t('sourceDetail.resultType')} value={facts.result_type ?? unavailable} />
      <RuntimeSetupFact label={t('sourceDetail.byteOrder')} value={decoding.length === 1 ? decoding[0].order ?? unavailable : unavailable} />
      <RuntimeSetupFact label={t('sourceDetail.polling')} value={facts.polling ? t('sourceDetail.pollingInterval', { interval: facts.polling.interval_ms, group: facts.polling.group_id }) : unavailable} />
    </dl>
    <h5 className="mt-3 text-slate-300">{t('sourceDetail.transforms')}</h5>
    {transforms.length ? <ol className="mt-1 space-y-1 text-slate-200">{transforms.map((step, index) => <li key={index}><Transform step={step} /></li>)}</ol> : <p className="text-slate-400">{t('sourceDetail.noTransforms')}</p>}
  </>;
}

function Transform({ step }: { step: SafeTransformStep }) {
  const { t } = useTranslation('runtime-dashboard');
  switch (step.kind) {
    case 'byte-order': return <>{t('sourceDetail.byteOrder')}: <span>{step.order}</span></>;
    case 'scale': return <>{t('sourceDetail.scale')}: <span>{`× ${step.multiplier} + ${step.offset}`}</span></>;
    case 'cast': return <>{t('sourceDetail.cast')}: <span>{step.target_type}</span></>;
    default: return <>{t('sourceDetail.unsupportedTransform')}</>;
  }
}
