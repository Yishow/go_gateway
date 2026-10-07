import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import LanguageSwitcher from '../../components/LanguageSwitcher';
import type { RuntimeLogFilters, RuntimeLogLevel } from '../../types/runtimeLogs';
import { useRuntimeLogs } from './useRuntimeLogs';
import { runtimeLogCursor } from './runtimeLogBuffer';
const control = 'rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-cyan-300';
const initialFilters: RuntimeLogFilters = { level: 'debug', source: '', q: '' };
export default function RuntimeLogsPage() {
  const { t } = useTranslation('runtime-logs');
  const [filters, setFilters] = useState(initialFilters);
  const [draft, setDraft] = useState(initialFilters);
  const [tail, setTail] = useState(true);
  const view = useRuntimeLogs(filters);
  const rows = useRef<HTMLDivElement>(null);
  useEffect(() => { if (tail && rows.current) rows.current.scrollTop = rows.current.scrollHeight; }, [tail, view.records]);
  const denied = view.status === 'denied'; const failed = view.status === 'failed';
  return (
    <main className="min-h-screen bg-slate-950 p-4 text-slate-100 sm:p-8">
      <div className="mx-auto max-w-7xl space-y-5">
        <header className="flex flex-wrap items-start justify-between gap-4">
          <div><p className="text-xs uppercase tracking-widest text-cyan-300">Studio</p>
            <h1 className="mt-1 text-3xl font-semibold">{t('title')}</h1>
            <p className="mt-2 max-w-3xl text-sm leading-6 text-slate-300">{t('scope')}</p></div>
          <LanguageSwitcher />
        </header>
        <nav aria-label={t('navigation')} className="flex flex-wrap gap-3">
          <a href="/studio/v2" className={control}>{t('setup')}</a>
          <a href="/studio/runtime" className={control}>{t('runtime')}</a>
        </nav>
        <form className="flex flex-wrap items-end gap-4 rounded-xl border border-slate-700 bg-slate-900/60 p-4"
          onSubmit={(event) => { event.preventDefault(); setFilters({ ...draft }); }}>
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
          <label className="grid min-w-48 flex-1 gap-2 text-sm">{t('search')}
            <input className={control} type="search" maxLength={256} value={draft.q} onChange={(event) => setDraft({ ...draft, q: event.target.value })} />
          </label>
          <button className={control} type="submit">{t('apply')}</button>
        </form>
        <div className="flex flex-wrap items-center gap-3">
          <button type="button" className={control} disabled={denied || failed} onClick={view.status === 'paused' ? view.resume : view.pause}>
            {t(view.status === 'paused' ? 'resume' : 'pause')}
          </button>
          <button type="button" className={control} onClick={view.clear}>{t('clear')}</button>
          <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={tail} onChange={(event) => setTail(event.target.checked)} />{t('tail')}</label>
          {(denied || failed) && <button type="button" className={control} onClick={view.retry}>{t('retry')}</button>}
          <span role="status" aria-live="polite" className="ml-auto text-sm text-cyan-200">{t(`status.${view.status}`)}</span>
        </div>
        <p className="text-xs text-slate-400">{t('controlsHelp')}</p>
        {denied && <p role="alert" className="rounded-lg border border-amber-600 p-4 text-amber-100">{t('deniedHelp')}</p>}
        {failed && <p role="alert" className="rounded-lg border border-rose-600 p-4 text-rose-100">{t('failedHelp')}</p>}
        <div aria-live="polite" className="space-y-2 text-sm text-amber-200">
          {view.reset && <p>{t('reset')}</p>}
          {view.gap && <p>{t('gap')}</p>}
          {view.fileGap && <p>{t('fileGap')}</p>}
          {view.truncated && <p>{t('truncated')}</p>}
          {view.metadata?.sink_health === 'degraded' && <p>{t('degraded')}</p>}
          {view.metadata?.sink_health === 'disabled' && <p className="text-slate-300">{t('storageDisabled')}</p>}
        </div>
        <div className="flex flex-wrap gap-4 text-xs text-slate-400">
          <span>{t('count', { count: view.records.length })}</span>
          {view.metadata && <><span>{t('counters', { dropped: view.metadata.capture_dropped, evicted: view.metadata.evicted,
            overflow: view.metadata.subscriber_overflow ?? 0 })}</span>
            <span>{t('fileCounters', { dropped: view.metadata.file_dropped ?? 0, errors: view.metadata.disk_errors ?? 0 })}</span></>}
          {view.lastUpdate && <span>{t('updated')} <time dateTime={view.lastUpdate}>{new Date(view.lastUpdate).toLocaleTimeString()}</time></span>}
        </div>
        <div ref={rows} role="region" aria-label={t('records')} tabIndex={0}
          className="max-h-[60vh] overflow-auto rounded-xl border border-slate-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-cyan-300">
          <table className="w-full text-left text-xs">
            <thead className="sticky top-0 bg-slate-800 text-slate-200"><tr>
              {['time', 'severity', 'sourceColumn', 'message'].map((key) => <th key={key} scope="col" className="p-3">{t(key)}</th>)}
            </tr></thead>
            <tbody className="divide-y divide-slate-800 font-mono">
              {view.records.map((record) => <tr key={runtimeLogCursor(record)} className="align-top">
                <td className="whitespace-nowrap p-3 text-slate-400"><time dateTime={record.timestamp}>{record.timestamp}</time></td>
                <td className="p-3">{t(`levels.${record.level}`)}</td><td className="p-3">{record.source}</td>
                <td className="max-w-3xl whitespace-pre-wrap break-words p-3">
                  <span className="text-cyan-300">{record.code}</span><p className="mt-1">{record.message}</p>
                  {Object.keys(record.fields).length > 0 && <p className="mt-1 text-slate-400">{JSON.stringify(record.fields)}</p>}
                  {record.truncated && <span className="text-amber-200">{t('recordTruncated')}</span>}
                </td>
              </tr>)}
            </tbody>
          </table>
          {view.status === 'connected' && view.records.length === 0 && <p className="p-8 text-center text-slate-400">{t('empty')}</p>}
        </div>
      </div>
    </main>
  );
}
