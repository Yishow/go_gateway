import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import { WORKBENCH_STEP_META } from './workbenchTypes';
import { useWorkbenchSummary } from './useWorkbenchSummary';

function getNextActionKey(input: {
  activeStep: ReturnType<typeof useWorkbench>['activeStep'];
  hasDevice: boolean;
  sourceReady: boolean;
  tagReady: boolean;
  outputReady: boolean;
}) {
  if (!input.hasDevice) {
    return 'workbench.actionDock.nextAction.selectDevice';
  }

  if (input.activeStep === 'device') {
    return 'workbench.actionDock.nextAction.advanceToSourcePlanning';
  }

  if (input.activeStep === 'source') {
    return input.sourceReady
      ? 'workbench.actionDock.nextAction.advanceToTag'
      : 'workbench.actionDock.nextAction.createPoints';
  }

  if (input.activeStep === 'tag') {
    return input.tagReady
      ? 'workbench.actionDock.nextAction.advanceToOutput'
      : 'workbench.actionDock.nextAction.bindTags';
  }

  if (input.activeStep === 'output') {
    return input.outputReady
      ? 'workbench.actionDock.nextAction.configureOutput'
      : 'workbench.actionDock.nextAction.completeTagBinding';
  }

  return 'workbench.actionDock.nextAction.advanceToSourcePlanning';
}

export function WorkbenchActionDock() {
  const { t } = useTranslation();
  const { activeStep } = useWorkbench();
  const {
    selectedDevice,
    pointCount,
    linkedTagCount,
    outputCandidateCount,
    sourceReady,
    tagReady,
    outputReady,
  } = useWorkbenchSummary();
  const selectedDeviceLabel = selectedDevice?.name ?? t('workbench.actionDock.noDevice');
  const nextActionKey = getNextActionKey({
    activeStep,
    hasDevice: selectedDevice !== null,
    sourceReady,
    tagReady,
    outputReady,
  });

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
              <dd className="font-medium text-slate-50">{selectedDeviceLabel}</dd>
            </div>
            <div className="space-y-1">
              <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.actionDock.nextActionLabel')}
              </dt>
              <dd className="font-medium text-cyan-100">
                {t(nextActionKey)}
              </dd>
            </div>
          </dl>

          <div className="grid gap-3 sm:grid-cols-3 xl:grid-cols-1">
            <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.actionDock.metrics.points')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">{pointCount}</p>
            </div>
            <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.actionDock.metrics.tags')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">{linkedTagCount}</p>
            </div>
            <div className="rounded-xl border border-slate-800 bg-slate-950/60 p-3">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.actionDock.metrics.outputs')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">{outputCandidateCount}</p>
            </div>
          </div>

          <div className="space-y-3">
            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.actionDock.readiness.title')}
            </p>
            <ul className="space-y-2 text-sm text-slate-200">
              <li className="flex items-center justify-between rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-2">
                <span>{t('workbench.actionDock.readiness.source')}</span>
                <span className={sourceReady ? 'text-emerald-200' : 'text-slate-400'}>
                  {sourceReady ? t('workbench.actionDock.readiness.ready') : t('workbench.actionDock.readiness.pending')}
                </span>
              </li>
              <li className="flex items-center justify-between rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-2">
                <span>{t('workbench.actionDock.readiness.tag')}</span>
                <span className={tagReady ? 'text-emerald-200' : 'text-slate-400'}>
                  {tagReady ? t('workbench.actionDock.readiness.ready') : t('workbench.actionDock.readiness.pending')}
                </span>
              </li>
              <li className="flex items-center justify-between rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-2">
                <span>{t('workbench.actionDock.readiness.output')}</span>
                <span className={outputReady ? 'text-emerald-200' : 'text-slate-400'}>
                  {outputReady ? t('workbench.actionDock.readiness.ready') : t('workbench.actionDock.readiness.pending')}
                </span>
              </li>
            </ul>
          </div>
        </div>
    </aside>
  );
}
