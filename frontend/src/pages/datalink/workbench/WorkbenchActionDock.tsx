import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import { WORKBENCH_STEP_META } from './workbenchTypes';

export function WorkbenchActionDock() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId } = useWorkbench();

  return (
    <aside
      aria-label={t('workbench.actionDock.ariaLabel')}
      className="rounded-2xl border border-slate-800 bg-slate-900/70 p-5 text-sm text-slate-200"
    >
      <div className="space-y-5">
        <div className="space-y-2">
          <h2 className="text-lg font-semibold text-slate-50">{t('workbench.actionDock.title')}</h2>
          <p className="text-slate-300">{t('workbench.actionDock.description')}</p>
        </div>

        <dl className="space-y-3">
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.actionDock.activeStep')}
            </dt>
            <dd className="font-medium text-slate-50">
              {t(WORKBENCH_STEP_META[activeStep].labelKey)}
            </dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.actionDock.selectedDevice')}
            </dt>
            <dd className="font-medium text-slate-50">
              {selectedDeviceId ?? t('workbench.actionDock.noDevice')}
            </dd>
          </div>
        </dl>
      </div>
    </aside>
  );
}
