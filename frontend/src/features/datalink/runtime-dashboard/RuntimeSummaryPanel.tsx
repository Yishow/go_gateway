import { useTranslation } from 'react-i18next';
import { Activity } from 'lucide-react';
import type { RuntimeStatus } from '../../../types/datalink';

interface RuntimeSummaryPanelProps {
  snapshot: RuntimeStatus;
}

function formatMetric(value: number | undefined, unavailable: string): string {
  return typeof value === 'number' && Number.isFinite(value) && value >= 0 ? String(value) : unavailable;
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

/** Runtime-wide observations are not selected-device totals or durable SQL proof. */
export function RuntimeSummaryPanel({ snapshot }: RuntimeSummaryPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const unavailable = t('summary.unavailable', 'Unavailable');
  const metrics = [
    { key: 'collectedTotal', label: 'Collected total', value: snapshot.metrics?.collected_total },
    { key: 'writeSuccessTotal', label: 'Write success total', value: snapshot.metrics?.write_success_total },
    { key: 'writeErrorTotal', label: 'Write error total', value: snapshot.metrics?.write_error_total },
    { key: 'mappingErrorTotal', label: 'Mapping error total', value: snapshot.metrics?.mapping_error_total },
    { key: 'pointStateErrorTotal', label: 'Point state error total', value: snapshot.metrics?.point_state_error_total },
  ];

  return (
    <section className="rounded-3xl border border-slate-800/80 bg-slate-900/70 p-6 shadow-xl backdrop-blur-sm" data-testid="runtime-dashboard-summary-panel">
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3 border-b border-slate-800/80 pb-4">
        <div className="flex items-start gap-3">
          <Activity className="mt-1 h-5 w-5 shrink-0 text-cyan-300" />
          <div>
            <h2 className="text-base font-semibold text-slate-50">{t('summary.title', 'Global runtime observations')}</h2>
            <p className="mt-1 text-xs text-slate-400">{t('summary.description', 'Counters cover the whole runtime, not the selected device or recording group.')}</p>
            <p className="mt-1 text-xs text-slate-400">{t('summary.notDeliveryProof', 'These global counters are not SQL receipts, readback or cleanup proof.')}</p>
          </div>
        </div>
        <span className="rounded-lg border border-slate-700 px-2.5 py-1 text-xs font-semibold text-slate-300">
          {t('summary.runningState', 'Runtime state')}: {snapshot.running ? t('summary.running', 'running') : t('summary.stopped', 'stopped')}
        </span>
      </div>
      <dl className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {metrics.map((metric) => <div key={metric.key} className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
          <dt className="text-xs font-semibold text-slate-400">{t(`summary.${metric.key}`, metric.label)}</dt>
          <dd className="mt-2 font-mono text-xl text-slate-100">{formatMetric(metric.value, unavailable)}</dd>
        </div>)}
        <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
          <dt className="text-xs font-semibold text-slate-400">{t('summary.uptimeSeconds', 'Uptime seconds')}</dt>
          <dd className="mt-2 font-mono text-xl text-slate-100">{String(snapshot.uptime_seconds)}</dd>
          <dd className="mt-1 text-xs text-slate-400">{formatUptime(snapshot.uptime_seconds)}</dd>
        </div>
      </dl>
    </section>
  );
}
