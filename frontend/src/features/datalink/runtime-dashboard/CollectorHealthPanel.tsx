import { useTranslation } from 'react-i18next';
import type { RuntimeCollectorStatus } from '../../../types/datalink';

interface CollectorHealthPanelProps {
  collector: RuntimeCollectorStatus | null;
}

export function CollectorHealthPanel({ collector }: CollectorHealthPanelProps) {
  const { t } = useTranslation('runtime-dashboard');

  if (!collector) {
    return (
      <section
        className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
        data-testid="runtime-dashboard-health-panel"
      >
        <h2 className="text-lg font-semibold text-slate-50">
          {t('health.title', 'Collector health')}
        </h2>
        <p className="mt-3 text-sm text-slate-400">
          {t('health.empty', 'Collector status will appear after the first runtime snapshot.')}
        </p>
      </section>
    );
  }

  const rows = [
    { label: t('health.status', 'Status'), value: collector.status },
    { label: t('health.pointsTotal', 'Points total'), value: String(collector.points_total) },
    { label: t('health.pointsHealthy', 'Points healthy'), value: String(collector.points_healthy) },
    { label: t('health.pointsStale', 'Points stale'), value: String(collector.points_stale) },
    { label: t('health.pointsError', 'Points error'), value: String(collector.points_error) },
    { label: t('health.lastReadAt', 'Last read at'), value: collector.last_read_at ?? '-' },
    { label: t('health.breakerState', 'Breaker state'), value: collector.breaker_state },
  ];

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
      data-testid="runtime-dashboard-health-panel"
    >
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-slate-50">
          {t('health.title', 'Collector health')}
        </h2>
        <p className="mt-1 text-sm text-slate-400">
          {t('health.description', 'Current collector counters for the selected device.')}
        </p>
      </div>

      <dl className="grid gap-3 sm:grid-cols-2">
        {rows.map((row) => (
          <div
            key={row.label}
            className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4"
          >
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-500">{row.label}</dt>
            <dd className="mt-2 text-lg font-semibold text-slate-50">{row.value}</dd>
          </div>
        ))}
      </dl>
    </section>
  );
}
