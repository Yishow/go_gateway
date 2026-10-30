import { useEffect, useId, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import LanguageSwitcher from '../../components/LanguageSwitcher';
import type { RuntimeLogFilters, RuntimeLogLevel } from '../../types/runtimeLogs';
import type { EvidenceStage } from '../../types/runtimeRuleEvidence';
import { safeRuntimeIdentity } from '../../utils/runtimeOperatorSafety';
import { useRuntimeLogs } from './useRuntimeLogs';
import { runtimeLogCursor } from './runtimeLogBuffer';
import RuntimeLogDomainEvent from './RuntimeLogDomainEvent';
import RuntimeLogLegacyEvent from './RuntimeLogLegacyEvent';
import { projectRuntimeLogDomain } from '../../services/runtimeLogDomainProjection';
const control = 'rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-600 disabled:cursor-not-allowed disabled:opacity-60';
const checkbox = 'accent-blue-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-600';
const stages: EvidenceStage[] = ['unavailable', 'acquisition', 'accepted', 'collecting', 'queued', 'retrying',
  'blocked', 'quarantined', 'unknown', 'skipped', 'sql-committed', 'readback-verified', 'cleanup-completed', 'transport-event'];
const scopeFields = ['workspaceId', 'deviceId', 'ruleId', 'groupId'] as const;
const initialFilters: RuntimeLogFilters = { level: 'debug', source: '', q: '' };
export default function RuntimeLogsPage() {
  const { t } = useTranslation('runtime-logs');
  const [filters, setFilters] = useState(initialFilters);
  const [draft, setDraft] = useState(initialFilters);
  const [tail, setTail] = useState(true);
  const [scopeOpen, setScopeOpen] = useState(false);
  const [invalidIDs, setInvalidIDs] = useState(false);
  const scopeID = useId();
  const appliedScopeCount = scopeFields.filter((field) => Boolean(filters[field])).length;
  const view = useRuntimeLogs(filters);
  const rows = useRef<HTMLDivElement>(null);
  useEffect(() => { if (tail && rows.current) rows.current.scrollTop = rows.current.scrollHeight; }, [tail, view.records]);
  const denied = view.status === 'denied'; const failed = view.status === 'failed';
  const incompleteHistory = view.gap || view.fileGap || view.reset || view.truncated || (view.metadata?.evicted ?? 0) > 0;
  return (
    <main className="min-h-screen bg-slate-50 p-4 text-slate-900 sm:p-8">
      <div className="mx-auto max-w-7xl space-y-5">
        <header className="flex flex-wrap items-start justify-between gap-4">
          <div><p className="text-xs uppercase tracking-widest text-blue-700">Studio</p>
            <h1 className="mt-1 text-3xl font-semibold">{t('title')}</h1>
            <p className="mt-2 max-w-3xl text-sm leading-6 text-slate-600">{t('scope')}</p></div>
          <LanguageSwitcher />
        </header>
        <nav aria-label={t('navigation')} className="flex flex-wrap gap-3">
          <a href="/studio/v2" className={control}>{t('setup')}</a>
          <a href="/studio/runtime" className={control}>{t('runtime')}</a>
        </nav>
        <form className="flex flex-wrap items-end gap-4 rounded-xl border border-slate-200 bg-white p-4"
          onSubmit={(event) => {
            event.preventDefault();
            const invalid = scopeFields.some((field) => {
              const value = draft[field];
              return value !== undefined && value !== '' && safeRuntimeIdentity(value) !== value;
            });
            setInvalidIDs(invalid);
            if (!invalid) setFilters({ ...draft });
          }}>
          <label className="grid gap-2 text-sm">{t('level')}
            <select className={control} value={draft.level} onChange={(event) => setDraft({ ...draft, level: event.target.value as RuntimeLogLevel })}>
              {(['debug', 'info', 'warn', 'error'] as const).map((level) => <option key={level} value={level}>{t(`levels.${level}`)}</option>)}
            </select>
          </label>
          <label className="grid gap-2 text-sm">{t('source')}
            <select className={control} value={draft.source} onChange={(event) => setDraft({ ...draft, source: event.target.value })}>
              <option value="">{t('allSources')}</option>
              {['startup', 'runtime', 'shutdown', 'http', 'standard', 'slog', 'application'].map((source) => <option key={source} value={source}>{source}</option>)}
            </select>
          </label>
          <label className="grid gap-2 text-sm">{t('stage')}
            <select className={control} value={draft.stage ?? ''} onChange={(event) => setDraft({ ...draft,
              stage: event.target.value === '' ? undefined : event.target.value as EvidenceStage })}>
              <option value="">{t('allStages')}</option>
              {stages.map((stage) => <option key={stage} value={stage}>{t(`stages.${stage}`)}</option>)}
            </select>
          </label>
          <label className="grid min-w-48 flex-1 gap-2 text-sm">{t('search')}
            <input className={control} type="search" maxLength={256} value={draft.q} onChange={(event) => setDraft({ ...draft, q: event.target.value })} />
          </label>
          <label className="flex items-center gap-2 py-2 text-sm">
            <input className={checkbox} type="checkbox" checked={draft.actionable ?? false} aria-describedby="runtime-log-actionable-help"
              onChange={(event) => setDraft({ ...draft, actionable: event.target.checked })} />{t('actionable')}
          </label>
          <button className={control} type="submit">{t('apply')}</button>
          <div className="w-full">
            <button type="button" className={`${control}${scopeOpen ? ' ring-2 ring-blue-600' : ''}`}
              aria-expanded={scopeOpen} aria-controls={scopeID} onClick={() => setScopeOpen(!scopeOpen)}>{t('scopeFilters')}
              {appliedScopeCount > 0 && <span className="ml-2 rounded bg-blue-100 px-2 py-0.5 text-blue-900">
                {t('appliedScopeCount', { count: appliedScopeCount })}
              </span>}
            </button>
            {scopeOpen && <div id={scopeID} className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <p className="text-sm text-slate-600 sm:col-span-2 lg:col-span-4">{t('scopeHelp')}</p>
              {scopeFields.map((field) => <label key={field} className="grid gap-2 text-sm">{t(`scopeFields.${field}`)}
                <input className={control} value={draft[field] ?? ''} maxLength={128} autoComplete="off" spellCheck={false}
                  onChange={(event) => setDraft({ ...draft, [field]: event.target.value })} />
              </label>)}
            </div>}
          </div>
          {invalidIDs && <p role="alert" className="w-full text-sm text-rose-800">{t('invalidId')}</p>}
        </form>
        <div className="flex flex-wrap items-center gap-3">
          <button type="button" className={control} disabled={denied || failed} onClick={view.status === 'paused' ? view.resume : view.pause}>
            {t(view.status === 'paused' ? 'resume' : 'pause')}
          </button>
          <button type="button" className={control} onClick={view.clear}>{t('clear')}</button>
          <label className="flex items-center gap-2 text-sm"><input className={checkbox} type="checkbox" checked={tail} onChange={(event) => setTail(event.target.checked)} />{t('tail')}</label>
          {(denied || failed) && <button type="button" className={control} onClick={view.retry}>{t('retry')}</button>}
          <span role="status" aria-live="polite" className="ml-auto text-sm text-blue-800">{t(`status.${view.status}`)}</span>
        </div>
        <p className="text-xs text-slate-600">{t('controlsHelp')}</p>
        <p id="runtime-log-actionable-help" className="text-sm text-slate-600">{t('actionableHelp')}</p>
        <p className="text-xs text-slate-600">{t('legacy.transportOnly')}</p>
        {denied && <p role="alert" className="rounded-lg border border-amber-300 bg-amber-50 p-4 text-amber-900">{t('deniedHelp')}</p>}
        {failed && <p role="alert" className="rounded-lg border border-rose-300 bg-rose-50 p-4 text-rose-900">{t('failedHelp')}</p>}
        <div aria-live="polite" className="space-y-2 text-sm text-amber-800">
          {incompleteHistory && <p>{t('historyUnknown')}</p>}
          {view.reset && <p>{t('reset')}</p>}
          {view.gap && <p>{t('gap')}</p>}
          {view.fileGap && <p>{t('fileGap')}</p>}
          {view.truncated && <p>{t('truncated')}</p>}
          {view.metadata?.sink_health === 'degraded' && <p>{t('degraded')}</p>}
          {view.metadata?.sink_health === 'disabled' && <p className="text-slate-600">{t('storageDisabled')}</p>}
        </div>
        <div className="flex flex-wrap gap-4 text-xs text-slate-600">
          <span>{t('count', { count: view.records.length })}</span>
          {view.metadata && <><span>{t('counters', { dropped: view.metadata.capture_dropped, evicted: view.metadata.evicted,
            overflow: view.metadata.subscriber_overflow ?? 0 })}</span>
            <span>{t('fileCounters', { dropped: view.metadata.file_dropped ?? 0, errors: view.metadata.disk_errors ?? 0 })}</span></>}
          {view.lastUpdate && <span>{t('updated')} <time dateTime={view.lastUpdate}>{new Date(view.lastUpdate).toLocaleTimeString()}</time></span>}
        </div>
        <div ref={rows} role="region" aria-label={t('records')} tabIndex={0}
          className="max-h-[60vh] overflow-auto rounded-xl border border-slate-200 bg-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-600">
          <table className="w-full text-left text-xs">
            <thead className="sticky top-0 bg-slate-100 text-slate-700"><tr>
              {['time', 'severity', 'sourceColumn', 'message'].map((key) => <th key={key} scope="col" className="p-3">{t(key)}</th>)}
            </tr></thead>
            <tbody className="divide-y divide-slate-200 font-mono">
              {view.records.map((record) => projectRuntimeLogDomain(record).kind !== 'legacy'
                ? <RuntimeLogDomainEvent key={runtimeLogCursor(record)} record={record} />
                : <RuntimeLogLegacyEvent key={runtimeLogCursor(record)} record={record} />)}
            </tbody>
          </table>
          {view.status === 'connected' && view.records.length === 0 && <p className="p-8 text-center text-slate-600">{t('empty')}</p>}
        </div>
      </div>
    </main>
  );
}
