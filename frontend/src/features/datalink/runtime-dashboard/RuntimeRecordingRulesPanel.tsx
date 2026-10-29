import { useTranslation } from 'react-i18next';
import type { EvidenceScope } from '../../../types/runtimeRuleEvidence';
import { useRuntimeRuleEvidence } from '../../../hooks/datalink/useRuntimeRuleEvidence';
import { safeRuntimeIdentity } from '../../../utils/runtimeOperatorSafety';
import { RuntimeRecordingGroupCard } from './RuntimeRecordingGroupCard';

interface Props { scope: EvidenceScope | null; selectedDeviceId: string | null; navigateTo?: (target: string) => void }

/** Configuration evidence is independent from collector health and never grants intake. */
export function RuntimeRecordingRulesPanel({ scope, selectedDeviceId, navigateTo }: Props) {
  const { t } = useTranslation('runtime-dashboard');
  const selected = scope && selectedDeviceId === scope.device_id &&
    [scope.workspace_id, scope.workspace_revision, scope.device_id].every((value) => Boolean(value) && safeRuntimeIdentity(value) === value)
    ? { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id } : null;
  return <section aria-label={t('recordingRules.title')} data-testid="runtime-recording-rules" className="rounded-3xl border border-slate-800 bg-slate-900/60 p-5">
    <h2 className="text-base font-semibold text-slate-100">{t('recordingRules.title')}</h2>
    <p className="mt-1 text-xs text-slate-400">{t('recordingRules.notIntake')}</p>
    {selected ? <ScopedRecordingRules key={JSON.stringify([selected.workspace_id, selected.workspace_revision, selected.device_id])} scope={selected} navigateTo={navigateTo} /> :
      <p role="status" className="mt-3 text-sm text-slate-400">{t('recordingRules.scopeUnavailable')}</p>}
  </section>;
}

function ScopedRecordingRules({ scope, navigateTo }: { scope: EvidenceScope; navigateTo?: (target: string) => void }) {
  const { t } = useTranslation('runtime-dashboard');
  const query = useRuntimeRuleEvidence(scope);
  const result = query.data;
  const view = !query.isError && result?.ok ? result.value : null;
  const failure = result && !result.ok ? result.reason : 'unavailable';
  return <>
    <button type="button" onClick={() => void query.refetch()} disabled={query.isFetching} className="mt-3 rounded border border-slate-700 px-3 py-1 text-xs text-slate-200 disabled:opacity-50">
      {t('recordingRules.refresh')}
    </button>
    {query.isPending ? <p role="status" className="mt-3 text-sm text-slate-400">{t('recordingRules.loading')}</p> : !view ?
      <p role="alert" className="mt-3 text-sm text-amber-100">{t(`recordingRules.errors.${failure}`)}</p> : <>
        <dl className="mt-3 flex flex-wrap gap-x-6 gap-y-2 text-xs text-slate-300">
          <div><dt className="text-slate-400">{t('recordingRules.observedAt')}</dt><dd>{view.observed_at ? <time dateTime={view.observed_at}>{view.observed_at}</time> : t('recordingRules.unavailable')}</dd></div>
          <div><dt className="text-slate-400">{t('recordingRules.freshness')}</dt><dd>{t(`recordingRules.freshnessValues.${view.freshness}`)}</dd></div>
        </dl>
        {view.groups.length ? <div className="mt-4 space-y-4">{view.groups.map((group) => <RuntimeRecordingGroupCard key={group.group_id} group={group} deviceId={scope.device_id} scope={!query.isFetching && view.freshness === 'fresh' ? scope : undefined} navigateTo={navigateTo} />)}</div> :
          view.freshness === 'fresh' && view.observed_at !== null ? <p className="mt-4 text-sm text-slate-400">{t('recordingRules.empty')}</p> :
            <p role="alert" className="mt-4 text-sm text-amber-100">{t('recordingRules.errors.unavailable')}</p>}
      </>}
  </>;
}
