import { useTranslation } from 'react-i18next';
import type { RuntimeFlowDiagnostic } from '../../../types/runtimeDiagnostics';

interface RuntimeDiagnosticsPanelProps {
  diagnostics?: RuntimeFlowDiagnostic[];
  selectedDeviceId: string | null;
}

export function RuntimeDiagnosticsPanel({
  diagnostics = [],
  selectedDeviceId,
}: RuntimeDiagnosticsPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const diagnostic = selectLatestDiagnostic(diagnostics, selectedDeviceId);

  if (!diagnostic) {
    return (
      <section
        className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
        data-testid="runtime-dashboard-diagnostics-panel"
      >
        <h2 className="text-lg font-semibold text-slate-50">
          {t('diagnostics.title', 'Runtime diagnostics')}
        </h2>
        <p className="mt-3 text-sm text-slate-400">
          {t('diagnostics.unavailable', 'Diagnostics unavailable')}
        </p>
      </section>
    );
  }

  const latestFailureAt = diagnostic.last_failure_at ?? t('diagnostics.unavailable', 'Unavailable');
  const rows = [
    {
      label: t('diagnostics.failureStage', 'Failure stage'),
      value: diagnostic.failure_stage ?? t('diagnostics.noFailureStage', 'No recent failure'),
    },
    {
      label: t('diagnostics.failureReason', 'Failure reason'),
      value: diagnostic.failure_reason ?? t('diagnostics.noFailureReason', 'No recent failure reason'),
    },
    {
      label: t('diagnostics.lastFailureAt', 'Last failure at'),
      value: latestFailureAt,
    },
    {
      label: t('diagnostics.latestSuccessfulStage', 'Latest successful stage'),
      value: diagnostic.latest_successful_stage ?? t('diagnostics.unavailable', 'Unavailable'),
    },
  ];

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
      data-testid="runtime-dashboard-diagnostics-panel"
    >
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-slate-50">
          {t('diagnostics.title', 'Runtime diagnostics')}
        </h2>
        <p className="mt-1 text-sm text-slate-400">
          {t('diagnostics.description', 'Latest selected-scope failure context.')}
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

function selectLatestDiagnostic(
  diagnostics: RuntimeFlowDiagnostic[],
  selectedDeviceId: string | null,
): RuntimeFlowDiagnostic | null {
  const scoped = diagnostics.filter((diagnostic) => {
    if (!selectedDeviceId) {
      return true;
    }
    return diagnostic.device_id === selectedDeviceId || diagnostic.scope === `device:${selectedDeviceId}`;
  });
  const sorted = [...scoped].sort((left, right) =>
    diagnosticTimestamp(right).localeCompare(diagnosticTimestamp(left)),
  );
  return sorted[0] ?? null;
}

function diagnosticTimestamp(diagnostic: RuntimeFlowDiagnostic): string {
  return diagnostic.last_failure_at ?? diagnostic.last_success_at ?? '';
}
