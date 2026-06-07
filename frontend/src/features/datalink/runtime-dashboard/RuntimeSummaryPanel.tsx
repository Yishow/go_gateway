import { useTranslation } from 'react-i18next';
import type { RuntimeStatus } from '../../../types/datalink';

interface RuntimeSummaryPanelProps {
  snapshot: RuntimeStatus;
}

function formatRunningState(running: boolean) {
  return running ? 'running' : 'stopped';
}

export function RuntimeSummaryPanel({ snapshot }: RuntimeSummaryPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const metrics = snapshot.metrics;
  const unavailable = t('summary.unavailable', 'Unavailable');

  const cards = [
    {
      label: t('summary.runningState', 'Runtime state'),
      value: formatRunningState(snapshot.running),
    },
    {
      label: t('summary.uptimeSeconds', 'Uptime seconds'),
      value: String(snapshot.uptime_seconds),
    },
    {
      label: t('summary.collectedTotal', 'Collected total'),
      value: formatMetric(metrics?.collected_total, unavailable),
    },
    {
      label: t('summary.writeSuccessTotal', 'Write success total'),
      value: formatMetric(metrics?.write_success_total, unavailable),
    },
    {
      label: t('summary.writeErrorTotal', 'Write error total'),
      value: formatMetric(metrics?.write_error_total, unavailable),
    },
    {
      label: t('summary.mappingErrorTotal', 'Mapping error total'),
      value: formatMetric(metrics?.mapping_error_total, unavailable),
    },
    {
      label: t('summary.pointStateErrorTotal', 'Point state error total'),
      value: formatMetric(metrics?.point_state_error_total, unavailable),
    },
  ];

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
      data-testid="runtime-dashboard-summary-panel"
    >
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-slate-50">
          {t('summary.title', 'Runtime summary')}
        </h2>
        <p className="mt-1 text-sm text-slate-400">
          {t('summary.description', 'Only backend-supported counters are shown here.')}
        </p>
      </div>

      <dl className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
        {cards.map((card) => (
          <div
            key={card.label}
            className="rounded-2xl border border-slate-800 bg-slate-950/70 p-4"
          >
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-500">{card.label}</dt>
            <dd className="mt-2 text-2xl font-semibold text-slate-50">{card.value}</dd>
          </div>
        ))}
      </dl>
    </section>
  );
}

function formatMetric(value: number | undefined, unavailable: string): string {
  return typeof value === 'number' ? String(value) : unavailable;
}
