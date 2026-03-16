import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import type { StepReadinessState, WorkbenchReadiness, WorkbenchStep } from './workbenchTypes';

function getReadinessDotColor(status: WorkbenchReadiness): string {
  switch (status) {
    case 'ready':
    case 'applied':
      return 'bg-emerald-400';
    case 'partial':
      return 'bg-amber-400';
    case 'blocked':
      return 'bg-rose-400';
    case 'draft':
      return 'bg-slate-600';
  }
}

function getReadinessTextColor(status: WorkbenchReadiness): string {
  switch (status) {
    case 'ready':
    case 'applied':
      return 'text-emerald-300';
    case 'partial':
      return 'text-amber-300';
    case 'blocked':
      return 'text-rose-300';
    case 'draft':
      return 'text-slate-500';
  }
}

export function WorkbenchBottomSummaryBar() {
  const { t } = useTranslation();
  const { activeOutputTarget } = useWorkbench();
  const {
    pointCount,
    linkedTagCount,
    outputCandidateCount,
    deviceReadiness,
    sourceReadiness,
    tagReadiness,
    outputReadiness,
  } = useWorkbenchSummary();

  const readinessItems: ReadonlyArray<{
    step: WorkbenchStep;
    labelKey: string;
    state: StepReadinessState;
  }> = [
    { step: 'device', labelKey: 'workbench.bottomSummary.readiness.device', state: deviceReadiness },
    { step: 'source', labelKey: 'workbench.bottomSummary.readiness.source', state: sourceReadiness },
    { step: 'tag', labelKey: 'workbench.bottomSummary.readiness.tag', state: tagReadiness },
    { step: 'output', labelKey: 'workbench.bottomSummary.readiness.output', state: outputReadiness },
  ];

  return (
    <footer
      aria-label={t('workbench.bottomSummary.ariaLabel')}
      className="flex flex-wrap items-center gap-x-6 gap-y-2 rounded-2xl border border-slate-800 bg-slate-900/80 px-4 py-2.5 text-xs text-slate-300"
      data-testid="workbench-bottom-summary-bar"
    >
      {/* Metric counters */}
      <dl className="flex items-center gap-5">
        <div className="flex items-center gap-1.5">
          <dt className="text-slate-500">{t('workbench.bottomSummary.points')}</dt>
          <dd className="font-semibold text-slate-100" data-testid="summary-point-count">
            {pointCount}
          </dd>
        </div>
        <div className="flex items-center gap-1.5">
          <dt className="text-slate-500">{t('workbench.bottomSummary.tags')}</dt>
          <dd className="font-semibold text-slate-100" data-testid="summary-tag-count">
            {linkedTagCount}
          </dd>
        </div>
        <div className="flex items-center gap-1.5">
          <dt className="text-slate-500">{t('workbench.bottomSummary.outputs')}</dt>
          <dd className="font-semibold text-slate-100" data-testid="summary-output-count">
            {outputCandidateCount}
          </dd>
        </div>
      </dl>

      {/* Separator */}
      <span className="hidden h-4 w-px bg-slate-700/60 sm:block" aria-hidden="true" />

      {/* Readiness indicators */}
      <ul className="flex items-center gap-4" data-testid="readiness-indicators">
        {readinessItems.map((item) => (
          <li
            key={item.step}
            className="flex items-center gap-1.5"
            data-testid={`readiness-${item.step}`}
            data-readiness={item.state.status}
          >
            <span
              aria-hidden="true"
              className={[
                'inline-block h-2 w-2 rounded-full',
                getReadinessDotColor(item.state.status),
              ].join(' ')}
            />
            <span className={getReadinessTextColor(item.state.status)}>
              {t(item.labelKey)}
            </span>
          </li>
        ))}
      </ul>

      {/* Separator */}
      <span className="hidden h-4 w-px bg-slate-700/60 sm:block" aria-hidden="true" />

      {/* Active output target badge */}
      <div className="flex items-center gap-1.5" data-testid="active-output-target">
        <span className="text-slate-500">{t('workbench.bottomSummary.target')}</span>
        <span className="font-semibold text-slate-100">
          {t(`workbench.bottomSummary.targets.${activeOutputTarget}`)}
        </span>
      </div>
    </footer>
  );
}
