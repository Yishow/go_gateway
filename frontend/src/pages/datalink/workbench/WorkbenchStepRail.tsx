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

/**
 * 依就緒狀態回傳 readiness 圓點的底色類名。
 *
 * @param status - 就緒、部分、阻擋等狀態
 * @returns Tailwind 背景色類名
 */
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

/**
 * 工作台步驟導覽（水平）：由 WorkbenchContextBar 嵌入左側，含「步驟」標籤與四步驟按鈕。
 */
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
      className="flex min-w-0 flex-[1_1_100%] flex-wrap items-stretch gap-2 sm:flex-[1_1_auto] lg:max-w-none"
      data-testid="workbench-step-rail"
    >
      <div
        className={[
          'flex min-w-0 flex-1 flex-wrap items-center gap-2 rounded-xl border border-slate-800/60',
          'bg-slate-950/50 p-1 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.04)]',
          'sm:flex-nowrap sm:gap-0',
        ].join(' ')}
      >
        <div className="flex shrink-0 items-center gap-2.5 px-2 py-1 sm:border-r sm:border-slate-800/55 sm:py-0 sm:pr-3">
          <span
            aria-hidden
            className="hidden h-6 w-0.5 shrink-0 rounded-full bg-gradient-to-b from-cyan-300/90 to-cyan-500/50 sm:block"
          />
          <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-slate-500">
            {t('workbench.stepRail.heading')}
          </p>
        </div>
        <ol className="flex min-w-0 flex-1 flex-wrap items-center gap-0.5 sm:flex-nowrap sm:pl-0.5">
          {WORKBENCH_STEPS.map((step, index) => {
            const isActive = step === activeStep;
            const readiness = stepReadiness[step];
            return (
              <li key={step} className="shrink-0">
                <button
                  type="button"
                  onClick={() => setActiveStep(step)}
                  aria-current={isActive ? 'step' : undefined}
                  aria-label={`${t(WORKBENCH_STEP_META[step].labelKey)} - ${t(`workbench.readiness.${readiness.status}`)}`}
                  data-readiness={readiness.status}
                  className={[
                    'flex min-h-9 items-center gap-2 rounded-lg px-2 py-1.5 text-left text-xs font-medium',
                    'transition-[color,background-color,border-color,box-shadow] duration-200',
                    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950',
                    isActive
                      ? 'border border-cyan-500/45 bg-slate-800/95 text-cyan-50 shadow-sm shadow-black/25'
                      : 'border border-transparent text-slate-400 hover:border-slate-700/60 hover:bg-slate-800/40 hover:text-slate-200',
                  ].join(' ')}
                >
                  <span
                    aria-hidden="true"
                    className={[
                      'flex h-6 w-6 shrink-0 items-center justify-center rounded-md text-[11px] font-semibold tabular-nums',
                      isActive
                        ? 'bg-cyan-500/20 text-cyan-100 ring-1 ring-cyan-400/35'
                        : 'bg-slate-800/90 text-slate-500 ring-1 ring-slate-700/50',
                    ].join(' ')}
                  >
                    {index + 1}
                  </span>
                  <span className="max-w-[10rem] truncate sm:max-w-[7.5rem] lg:max-w-[9rem]">
                    {t(WORKBENCH_STEP_META[step].labelKey)}
                  </span>
                  <span
                    aria-hidden="true"
                    className={[
                      'inline-block h-2 w-2 shrink-0 rounded-full ring-2 ring-slate-950/80',
                      getReadinessColor(readiness.status),
                    ].join(' ')}
                    data-testid={`step-readiness-${step}`}
                  />
                </button>
              </li>
            );
          })}
        </ol>
      </div>
    </nav>
  );
}
