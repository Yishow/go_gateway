import { useTranslation } from 'react-i18next';
import type { RuntimeLogRecord } from '../../types/runtimeLogs';
import { projectRuntimeLogLegacy } from '../../services/runtimeLogLegacyProjection';

/** Fixed process/HTTP history only; no entity, operation or delivery authority. */
export default function RuntimeLogLegacyEvent({ record }: { record: RuntimeLogRecord }) {
  const { t } = useTranslation('runtime-logs');
  const view = projectRuntimeLogLegacy(record);
  const unavailable = t('legacy.unavailable');
  const known = view.kind === 'legacy';
  const facts: [label: string, value: string | null][] = [[t('legacy.reference'), view.cursor]];
  if (known) {
    facts.push([t('legacy.code'), view.code]);
    if (view.category === 'http') facts.push([t('legacy.method'), view.method],
      [t('legacy.status'), view.status === null ? null : String(view.status)],
      [t('legacy.duration'), view.durationMs === null ? null : String(view.durationMs)],
      [t('legacy.requestID'), view.requestID]);
    if (view.code.startsWith('shutdown.')) facts.push([t('legacy.phase'), view.phase === null ? null : t(`legacy.phases.${view.phase}`)]);
  }
  return <tr role="row" data-log-cursor={view.cursor ?? undefined} className="block align-top font-sans sm:table-row">
    <td role="cell" className="block whitespace-nowrap p-3 text-slate-600 sm:table-cell">
      {view.timestamp ? <time dateTime={view.timestamp}>{view.timestamp}</time> : unavailable}
    </td>
    <td role="cell" className="block p-3 sm:table-cell">{known ? t(`levels.${view.level}`) : unavailable}</td>
    <td role="cell" className="block p-3 sm:table-cell">{known ? t(`legacy.sources.${view.source}`) : t('legacy.suppressedSource')}</td>
    <td role="cell" className="block min-w-0 max-w-3xl space-y-3 break-words p-3 sm:table-cell">
      <h2 className="font-semibold text-slate-900">{known ? t(`legacy.problems.${view.code}`) : t('legacy.suppressedProblem')}</h2>
      <p className="text-slate-700">{known ? t(`legacy.impacts.${view.code}`) : t('legacy.suppressedImpact')}</p>
      <p className="text-slate-600">{t(`legacy.categories.${known ? view.category : 'suppressed'}`)}</p>
      <details>
        <summary className="cursor-pointer rounded text-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-600">{t('legacy.details')}</summary>
        <div className="mt-3 space-y-3 rounded-lg border border-slate-200 bg-slate-50 p-3">
          {known && view.category === 'http' && <p className="text-slate-600">{t('legacy.safeRequestSummary')}</p>}
          <dl className="grid gap-x-6 gap-y-2 sm:grid-cols-2">
            {facts.map(([label, value]) => <div key={label} className="min-w-0">
              <dt className="text-slate-600">{label}</dt>
              <dd className="break-words font-mono text-slate-900">{value ?? unavailable}</dd>
            </div>)}
          </dl>
        </div>
      </details>
      {view.truncated && <p className="text-amber-800">{t('recordTruncated')}</p>}
    </td>
  </tr>;
}
