import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import {
  WORKBENCH_STEPS,
  WORKBENCH_STEP_META,
  type StepReadinessState,
  type WorkbenchReadiness,
  type WorkbenchStep,
} from './workbenchTypes';

function getReadinessColor(status: WorkbenchReadiness): string {
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

export function WorkbenchStepRail() {
  const { t } = useTranslation();
  const { activeStep, setActiveStep } = useWorkbench();
  const {
    deviceReadiness,
    sourceReadiness,
    tagReadiness,
    outputReadiness,
  } = useWorkbenchSummary();

  const stepReadiness: Record<WorkbenchStep, StepReadinessState> = {
    device: deviceReadiness,
    source: sourceReadiness,
    tag: tagReadiness,
    output: outputReadiness,
  };

  return (
    <nav
      aria-label={t('workbench.stepRail.ariaLabel')}
      className="flex flex-col gap-1 rounded-2xl border border-slate-800 bg-slate-900/80 p-2"
      data-testid="workbench-step-rail"
    >
      <p className="px-3 pb-1 pt-2 text-[10px] font-bold uppercase tracking-[0.22em] text-slate-500">
        {t('workbench.stepRail.heading')}
      </p>
      <ol className="flex flex-col gap-1">
        {WORKBENCH_STEPS.map((step, index) => {
          const isActive = step === activeStep;
          const readiness = stepReadiness[step];
          return (
            <li key={step}>
              <button
                type="button"
                onClick={() => setActiveStep(step)}
                aria-current={isActive ? 'step' : undefined}
                data-readiness={readiness.status}
                className={[
                  'flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left text-sm font-medium transition-all',
                  isActive
                    ? 'border border-cyan-500/30 bg-cyan-500/10 text-cyan-100 shadow-sm shadow-cyan-900/20'
                    : 'border border-transparent text-slate-400 hover:bg-slate-800/60 hover:text-slate-200',
                ].join(' ')}
              >
                <span
                  aria-hidden="true"
                  className={[
                    'flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-xs font-semibold',
                    isActive
                      ? 'bg-cyan-500/20 text-cyan-200'
                      : 'bg-slate-800/80 text-slate-500',
                  ].join(' ')}
                >
                  {index + 1}
                </span>
                <span className="truncate">{t(WORKBENCH_STEP_META[step].labelKey)}</span>
                <span
                  aria-hidden="true"
                  className={[
                    'ml-auto inline-block h-2 w-2 rounded-full',
                    getReadinessColor(readiness.status),
                  ].join(' ')}
                  data-testid={`step-readiness-${step}`}
                />
              </button>
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
