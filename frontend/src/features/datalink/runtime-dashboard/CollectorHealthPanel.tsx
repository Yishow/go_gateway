import { useTranslation } from 'react-i18next';
import { Cpu, ShieldCheck, AlertCircle } from 'lucide-react';
import type { RuntimeCollectorStatus } from '../../../types/datalink';

type RuntimeCollectorProjectionStatus = RuntimeCollectorStatus & {
  projection_alignment?: string;
  runtime_projection_version?: string;
  workspace_projection_version?: string;
  projection_message?: string;
};

interface CollectorHealthPanelProps {
  collector: RuntimeCollectorProjectionStatus | null;
}

function formatRelativeTime(isoString: string | null | undefined): string {
  if (!isoString) return '-';
  try {
    const diff = Math.floor((Date.now() - new Date(isoString).getTime()) / 1000);
    if (diff < 3) return '剛剛 (Just now)';
    if (diff < 60) return `${diff} 秒前`;
    if (diff < 3600) return `${Math.floor(diff / 60)} 分鐘前`;
    return new Date(isoString).toLocaleTimeString();
  } catch {
    return isoString;
  }
}

export function CollectorHealthPanel({ collector }: CollectorHealthPanelProps) {
  const { t } = useTranslation('runtime-dashboard');

  if (!collector) {
    return (
      <section
        className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm"
        data-testid="runtime-dashboard-health-panel"
      >
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-slate-700 bg-slate-800/50 text-slate-400">
            <Cpu className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-base font-semibold text-slate-50">
              {t('health.title', 'Collector health')}
            </h2>
            <p className="text-xs text-slate-400">
              {t('health.empty', 'Collector status will appear after the first runtime snapshot.')}
            </p>
          </div>
        </div>
      </section>
    );
  }

  const totalPoints = Math.max(1, collector.points_total);
  const healthyPct = (collector.points_healthy / totalPoints) * 100;
  const stalePct = (collector.points_stale / totalPoints) * 100;
  const errorPct = (collector.points_error / totalPoints) * 100;
  const isBreakerClosed = collector.breaker_state === 'closed';

  return (
    <section
      className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm"
      data-testid="runtime-dashboard-health-panel"
    >
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3 border-b border-slate-800/80 pb-4">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-emerald-500/30 bg-emerald-500/10 text-emerald-300">
            <Cpu className="h-5 w-5" />
          </div>
          <div>
            <h2 className="text-base font-semibold text-slate-50">
              {t('health.title', 'Collector health')} • PLC 通訊與點位健康
            </h2>
            <p className="text-xs text-slate-400">
              {t('health.description', 'Current collector counters for the selected device.')} (PLC 連線狀態與點位健康度分佈)
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <span className="text-[11px] text-slate-500 uppercase tracking-wider">{t('health.status', 'Status')}:</span>
          <span className="rounded-lg border border-cyan-500/30 bg-cyan-500/10 px-2.5 py-1 text-xs font-semibold uppercase tracking-wider text-cyan-300">
            {collector.status}
          </span>
        </div>
      </div>

      {/* 點位品質比例條 */}
      <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
        <div className="flex items-center justify-between text-xs text-slate-400 mb-2">
          <span className="font-semibold uppercase tracking-wider">點位健康度分佈</span>
          <span className="font-mono text-emerald-300">
            {collector.points_healthy}/{collector.points_total} 正常 ({(healthyPct).toFixed(0)}%)
          </span>
        </div>
        <div className="flex h-2.5 w-full overflow-hidden rounded-full bg-slate-800">
          <div style={{ width: `${healthyPct}%` }} className="bg-emerald-500 transition-all duration-500" />
          <div style={{ width: `${stalePct}%` }} className="bg-amber-400 transition-all duration-500" />
          <div style={{ width: `${errorPct}%` }} className="bg-rose-500 transition-all duration-500" />
        </div>

        <div className="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
          <div className="rounded-xl border border-slate-800 bg-slate-900/60 p-2.5">
            <div className="text-[11px] text-slate-500">{t('health.pointsTotal', 'Points total')}</div>
            <div className="mt-1 text-lg font-bold font-mono text-slate-100">{String(collector.points_total)}</div>
          </div>
          <div className="rounded-xl border border-slate-800 bg-slate-900/60 p-2.5">
            <div className="text-[11px] text-emerald-400/80">{t('health.pointsHealthy', 'Points healthy')}</div>
            <div className="mt-1 text-lg font-bold font-mono text-emerald-400">{String(collector.points_healthy)}</div>
          </div>
          <div className="rounded-xl border border-slate-800 bg-slate-900/60 p-2.5">
            <div className="text-[11px] text-amber-400/80">{t('health.pointsStale', 'Points stale')}</div>
            <div className="mt-1 text-lg font-bold font-mono text-amber-400">{String(collector.points_stale)}</div>
          </div>
          <div className="rounded-xl border border-slate-800 bg-slate-900/60 p-2.5">
            <div className="text-[11px] text-rose-400/80">{t('health.pointsError', 'Points error')}</div>
            <div className="mt-1 text-lg font-bold font-mono text-rose-400">{String(collector.points_error)}</div>
          </div>
        </div>
      </div>

      {/* 通訊與對齊狀態 */}
      <div className="mt-4 grid gap-3 sm:grid-cols-2">
        {/* 斷路器與通訊 */}
        <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-3.5 flex items-center justify-between">
          <div>
            <div className="text-[11px] text-slate-500 uppercase tracking-wider">{t('health.breakerState', 'Breaker state')}</div>
            <div className="mt-1 flex items-center gap-1.5 font-mono text-sm font-semibold text-slate-200">
              {isBreakerClosed ? (
                <ShieldCheck className="h-4 w-4 text-emerald-400" />
              ) : (
                <AlertCircle className="h-4 w-4 text-rose-400" />
              )}
              <span>{collector.breaker_state}</span>
              <span className="text-xs font-sans text-slate-400 font-normal">
                ({isBreakerClosed ? '正常導通' : '保護斷開'})
              </span>
            </div>
          </div>
          <div className="text-right">
            <div className="text-[11px] text-slate-500 uppercase tracking-wider">{t('health.lastReadAt', 'Last read at')}</div>
            <div className="mt-1 text-xs font-mono text-cyan-200" title={collector.last_read_at ?? '-'}>
              {formatRelativeTime(collector.last_read_at)}
            </div>
            <span className="sr-only">{collector.last_read_at ?? '-'}</span>
          </div>
        </div>

        {/* 投影對齊資訊 */}
        <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-3.5">
          <div className="flex items-center justify-between">
            <span className="text-[11px] text-slate-500 uppercase tracking-wider">
              {t('health.projectionAlignment', 'Projection alignment')}
            </span>
            <span
              className={`rounded px-1.5 py-0.5 text-[11px] font-semibold ${
                collector.projection_alignment === 'aligned'
                  ? 'border border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                  : 'border border-amber-500/30 bg-amber-500/10 text-amber-300'
              }`}
            >
              {collector.projection_alignment ?? 'aligned'}
            </span>
          </div>
          <div className="mt-2 space-y-1 text-[11px] text-slate-400 font-mono">
            {collector.runtime_projection_version && (
              <div className="truncate">
                <span className="text-slate-500">{t('health.runtimeProjectionVersion', 'Runtime projection')}: </span>
                <span className="text-slate-300">{collector.runtime_projection_version}</span>
              </div>
            )}
            {collector.workspace_projection_version && (
              <div className="truncate">
                <span className="text-slate-500">{t('health.workspaceProjectionVersion', 'Workspace projection')}: </span>
                <span className="text-slate-300">{collector.workspace_projection_version}</span>
              </div>
            )}
            {collector.projection_message && (
              <div className="text-amber-300">
                <span className="text-slate-500">{t('health.projectionMessage', 'Projection message')}: </span>
                <span>{collector.projection_message}</span>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
