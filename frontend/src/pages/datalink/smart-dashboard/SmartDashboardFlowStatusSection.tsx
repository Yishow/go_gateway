import type { TFunction } from 'i18next';
import type { FlowState, FlowSegment } from '../../../features/flow/stateMachine';

interface SmartDashboardFlowStatusSectionProps {
  flowSegments: FlowSegment[];
  flowState: FlowState;
  statusStyle: Record<string, string>;
  hasError: boolean;
  t: TFunction;
}

export default function SmartDashboardFlowStatusSection({
  flowSegments,
  flowState,
  statusStyle,
  hasError,
  t,
}: SmartDashboardFlowStatusSectionProps) {
  return (
    <section className="mb-6 rounded-2xl border border-white/10 bg-slate-900/50 p-4">
      <div className="mb-3 flex items-center justify-between gap-2">
        <h3 className="text-sm font-semibold text-slate-100">{t('smartDashboard.flowTitle')}</h3>
        <span className={`rounded-lg border px-2 py-1 text-xs ${statusStyle[flowState.status]}`}>
          {t(`smartDashboard.flowStatus.${flowState.status}`)}
        </span>
      </div>
      {hasError && (
        <p className="mb-3 rounded-lg border border-red-500/40 bg-red-500/10 px-3 py-2 text-xs text-red-200">
          {t('smartDashboard.flowErrorHint')}
        </p>
      )}
      <div className="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-4">
        {flowSegments.map((segment) => {
          const diag = flowState.diagnostics[segment];
          return (
            <article key={segment} className="rounded-xl border border-white/10 bg-slate-800/40 p-3">
              <div className="flex items-center justify-between">
                <p className="text-xs uppercase tracking-wider text-slate-400">
                  {t(`smartDashboard.flowSegments.${segment}.title`)}
                </p>
                <span
                  className={`h-2 w-2 rounded-full ${
                    diag.quality === 'good'
                      ? 'bg-emerald-400'
                      : diag.quality === 'warning'
                        ? 'bg-yellow-400'
                        : diag.quality === 'bad'
                          ? 'bg-red-400'
                          : 'bg-slate-500'
                  }`}
                />
              </div>
              <p className="mt-1 text-xs text-slate-500">{t(`smartDashboard.flowSegments.${segment}.subtitle`)}</p>
              <p className="mt-2 truncate font-mono text-sm text-slate-100">{diag.latestValue}</p>
              <p className="mt-1 truncate text-[11px] text-slate-400">{diag.timestamp}</p>
              {diag.error && <p className="mt-1 truncate text-[11px] text-red-300">{diag.error}</p>}
            </article>
          );
        })}
      </div>
    </section>
  );
}
