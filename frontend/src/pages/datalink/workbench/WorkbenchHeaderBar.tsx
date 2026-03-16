import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import { WORKBENCH_STEP_META } from './workbenchTypes';

export function WorkbenchHeaderBar() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId } = useWorkbench();

  return (
    <header
      className="mx-auto flex w-full max-w-7xl flex-col gap-6 px-6 py-8 lg:flex-row lg:items-end lg:justify-between"
      aria-label={t('workbench.header.ariaLabel')}
    >
      <div className="space-y-3">
        <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
          {t('workbench.header.eyebrow')}
        </p>
        <div className="space-y-2">
          <h1 className="text-3xl font-semibold text-slate-50">{t('workbench.title')}</h1>
          <p className="max-w-3xl text-sm text-slate-300">{t('workbench.subtitle')}</p>
        </div>
      </div>

      <dl className="grid gap-4 rounded-2xl border border-slate-800 bg-slate-900/70 p-4 text-sm text-slate-200 sm:grid-cols-2">
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.header.currentStep')}
          </dt>
          <dd className="font-medium text-slate-50">
            {t(WORKBENCH_STEP_META[activeStep].labelKey)}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.header.selectedDevice')}
          </dt>
          <dd className="font-medium text-slate-50">
            {selectedDeviceId ?? t('workbench.header.noDevice')}
          </dd>
        </div>
      </dl>
    </header>
  );
}
