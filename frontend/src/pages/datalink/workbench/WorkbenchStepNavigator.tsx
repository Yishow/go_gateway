import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import { WORKBENCH_STEPS, WORKBENCH_STEP_META } from './workbenchTypes';

export function WorkbenchStepNavigator() {
  const { t } = useTranslation();
  const { activeStep, setActiveStep } = useWorkbench();

  return (
    <nav
      aria-label={t('workbench.stepNavigator.ariaLabel')}
      className="rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
    >
      <ol className="grid gap-3 sm:grid-cols-2 xl:grid-cols-1">
        {WORKBENCH_STEPS.map((step, index) => {
          const isActive = step === activeStep;

          return (
            <li key={step}>
              <button
                type="button"
                onClick={() => setActiveStep(step)}
                aria-current={isActive ? 'step' : undefined}
                className={`flex w-full items-center gap-3 rounded-xl border px-4 py-3 text-left transition ${
                  isActive
                    ? 'border-cyan-400 bg-cyan-500/10 text-cyan-100'
                    : 'border-slate-800 bg-slate-950/60 text-slate-200 hover:border-slate-700'
                }`}
              >
                <span
                  aria-hidden="true"
                  className="flex h-8 w-8 items-center justify-center rounded-full border border-current text-xs font-semibold"
                >
                  {index + 1}
                </span>
                <span className="font-medium">{t(WORKBENCH_STEP_META[step].labelKey)}</span>
              </button>
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
