import { useTranslation } from 'react-i18next';
import { Activity, Database, CheckCircle2, AlertOctagon, Clock } from 'lucide-react';
import type { RuntimeStatus } from '../../../types/datalink';

interface RuntimeSummaryPanelProps {
  snapshot: RuntimeStatus;
}

function formatRunningState(running: boolean) {
  return running ? 'running' : 'stopped';
}

function formatMetric(value: number | undefined, unavailable: string): string {
  return typeof value === 'number' ? String(value) : unavailable;
}

function formatUptime(seconds: number): string {
  if (seconds < 60) return `${seconds}s`;
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  if (m < 60) return `${m}m ${s}s`;
  const h = Math.floor(m / 60);
  const remM = m % 60;
  return `${h}h ${remM}m`;
}

export function RuntimeSummaryPanel({ snapshot }: RuntimeSummaryPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const metrics = snapshot.metrics;
  const unavailable = t('summary.unavailable', 'Unavailable');

  const collectedVal = metrics?.collected_total;
  const writeSuccessVal = metrics?.write_success_total;
  const writeErrorVal = metrics?.write_error_total ?? 0;
  const mappingErrorVal = metrics?.mapping_error_total ?? 0;
  const pointStateErrorVal = metrics?.point_state_error_total ?? 0;
  const totalErrors = writeErrorVal + mappingErrorVal + pointStateErrorVal;

  const successRate =
    typeof writeSuccessVal === 'number' && typeof collectedVal === 'number' && collectedVal > 0
      ? ((writeSuccessVal / Math.max(writeSuccessVal + writeErrorVal, collectedVal)) * 100).toFixed(1)
      : typeof writeSuccessVal === 'number' && writeSuccessVal > 0
      ? '100.0'
      : null;

  return (
    <section
      className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm"
      data-testid="runtime-dashboard-summary-panel"
    >
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3 border-b border-slate-800/80 pb-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-cyan-500/30 bg-cyan-500/10 text-cyan-300">
            <Activity className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-base font-semibold text-slate-50">
              {t('summary.title', 'Runtime summary')}
            </h2>
            <p className="text-xs text-slate-400">
              {t('summary.description', 'Only backend-supported counters are shown here.')}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <span className="relative flex h-2.5 w-2.5">
            {snapshot.running && (
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
            )}
            <span
              className={`relative inline-flex h-2.5 w-2.5 rounded-full ${
                snapshot.running ? 'bg-emerald-500' : 'bg-amber-500'
              }`}
            />
          </span>
          <span
            className={`rounded-lg px-2.5 py-1 text-xs font-semibold uppercase tracking-wider ${
              snapshot.running
                ? 'border border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                : 'border border-amber-500/30 bg-amber-500/10 text-amber-300'
            }`}
          >
            {formatRunningState(snapshot.running)}
          </span>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {/* KPI 1: 交付成功率 */}
        <div className="relative overflow-hidden rounded-2xl border border-slate-800/90 bg-gradient-to-br from-slate-950/90 to-slate-900/50 p-4">
          <div className="flex items-center justify-between text-slate-400">
            <span className="text-xs font-semibold uppercase tracking-wider">交付成功率</span>
            <CheckCircle2 className="h-4 w-4 text-emerald-400" />
          </div>
          <div className="mt-2 text-2xl font-bold font-mono text-emerald-400">
            {successRate ? `${successRate}%` : unavailable}
          </div>
          <div className="mt-2 text-[11px] text-slate-400 flex items-center justify-between">
            <span>{t('summary.writeSuccessTotal', 'Write success total')}:</span>
            <span className="font-mono text-slate-200">
              {formatMetric(metrics?.write_success_total, unavailable)}
            </span>
          </div>
          <div className="mt-2 h-1.5 w-full rounded-full bg-slate-800 overflow-hidden">
            <div
              className="h-full rounded-full bg-emerald-500 transition-all duration-500"
              style={{ width: `${successRate ?? (metrics ? 100 : 0)}%` }}
            />
          </div>
        </div>

        {/* KPI 2: 採集總數 */}
        <div className="relative overflow-hidden rounded-2xl border border-slate-800/90 bg-gradient-to-br from-slate-950/90 to-slate-900/50 p-4">
          <div className="flex items-center justify-between text-slate-400">
            <span className="text-xs font-semibold uppercase tracking-wider">
              {t('summary.collectedTotal', 'Collected total')}
            </span>
            <Database className="h-4 w-4 text-cyan-400" />
          </div>
          <div className="mt-2 text-2xl font-bold font-mono text-cyan-300">
            {formatMetric(metrics?.collected_total, unavailable)}
          </div>
          <p className="mt-2 text-[11px] text-slate-400">
            即時採集累計數據點
          </p>
        </div>

        {/* KPI 3: 運轉時間 */}
        <div className="relative overflow-hidden rounded-2xl border border-slate-800/90 bg-gradient-to-br from-slate-950/90 to-slate-900/50 p-4">
          <div className="flex items-center justify-between text-slate-400">
            <span className="text-xs font-semibold uppercase tracking-wider">
              {t('summary.runningState', 'Runtime state')}
            </span>
            <Clock className="h-4 w-4 text-sky-400" />
          </div>
          <div className="mt-2 text-2xl font-bold font-mono text-slate-100">
            {formatUptime(snapshot.uptime_seconds)}
          </div>
          <div className="mt-2 text-[11px] text-slate-400 flex items-center justify-between">
            <span>{t('summary.uptimeSeconds', 'Uptime seconds')}:</span>
            <span className="font-mono text-slate-300">{String(snapshot.uptime_seconds)}</span>
          </div>
        </div>

        {/* KPI 4: 管線異常檢測 */}
        <div className="relative overflow-hidden rounded-2xl border border-slate-800/90 bg-gradient-to-br from-slate-950/90 to-slate-900/50 p-4">
          <div className="flex items-center justify-between text-slate-400">
            <span className="text-xs font-semibold uppercase tracking-wider">管線異常監控</span>
            <AlertOctagon
              className={`h-4 w-4 ${totalErrors > 0 ? 'text-rose-400' : 'text-emerald-400'}`}
            />
          </div>
          <div
            className={`mt-2 text-2xl font-bold font-mono ${
              totalErrors > 0 ? 'text-rose-400' : 'text-emerald-400'
            }`}
          >
            {totalErrors === 0 ? '0 異常' : `${totalErrors} 異常`}
          </div>
          <div className="mt-2 space-y-0.5 text-[10px] text-slate-400">
            <div className="flex justify-between">
              <span>{t('summary.writeErrorTotal', 'Write error total')}:</span>
              <span className="font-mono text-slate-300">{formatMetric(metrics?.write_error_total, unavailable)}</span>
            </div>
            <div className="flex justify-between">
              <span>{t('summary.mappingErrorTotal', 'Mapping error total')}:</span>
              <span className="font-mono text-slate-300">{formatMetric(metrics?.mapping_error_total, unavailable)}</span>
            </div>
            <div className="flex justify-between">
              <span>{t('summary.pointStateErrorTotal', 'Point state error total')}:</span>
              <span className="font-mono text-slate-300">{formatMetric(metrics?.point_state_error_total, unavailable)}</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
