import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import {
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
} from './workbenchDeviceFormModel';
import { WORKBENCH_STEP_META } from './workbenchTypes';
import { useWorkbenchSummary } from './useWorkbenchSummary';

function joinClasses(...classNames: Array<string | false | null | undefined>) {
  return classNames.filter(Boolean).join(' ');
}

function getStatusClasses(status: 'draft' | 'active' | 'disabled') {
  switch (status) {
    case 'active':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-200';
    case 'disabled':
      return 'border-rose-500/40 bg-rose-500/10 text-rose-200';
    case 'draft':
      return 'border-amber-500/40 bg-amber-500/10 text-amber-200';
  }
}

export function WorkbenchHeaderBar() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId, setActiveStep } = useWorkbench();
  const {
    selectedDevice,
    pointCount,
    linkedTagCount,
    outputCandidateCount,
  } = useWorkbenchSummary();
  const selectedDeviceLabel =
    selectedDevice?.name ?? t('workbench.header.noDevice');
  const lastTestLabel = selectedDevice
    ? selectedDevice.last_test_success === true
      ? t('workbench.header.testPassed')
      : selectedDevice.last_test_success === false
        ? selectedDevice.last_test_error || t('workbench.header.testFailed')
        : t('workbench.header.notTested')
    : t('workbench.header.noDevice');

  return (
    <header
      className="mx-auto flex w-full max-w-7xl flex-col gap-6 px-6 py-8"
      aria-label={t('workbench.header.ariaLabel')}
    >
      <div className="flex flex-col gap-6 xl:flex-row xl:items-end xl:justify-between">
        <div className="space-y-3">
          <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
            {t('workbench.header.eyebrow')}
          </p>
          <div className="space-y-2">
            <h1 className="text-3xl font-semibold text-slate-50">
              {t('workbench.title')}
            </h1>
            <p className="max-w-3xl text-sm text-slate-300">
              {t('workbench.subtitle')}
            </p>
          </div>
        </div>

        <div className="grid gap-4 rounded-3xl border border-slate-800 bg-slate-900/70 p-4 text-sm text-slate-200 xl:min-w-[520px] xl:grid-cols-[minmax(0,1fr)_auto]">
          <div className="space-y-3">
            <div className="flex flex-wrap items-center gap-2">
              <span className="text-lg font-semibold text-slate-50">
                {selectedDeviceLabel}
              </span>
              {selectedDevice ? (
                <>
                  <span className="rounded-full border border-slate-700 bg-slate-950/70 px-2 py-1 text-xs text-slate-300">
                    {t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))}
                  </span>
                  <span
                    className={joinClasses(
                      'rounded-full border px-2 py-1 text-xs font-medium',
                      getStatusClasses(selectedDevice.status),
                    )}
                  >
                    {t(getWorkbenchDeviceStatusLabelKey(selectedDevice.status))}
                  </span>
                </>
              ) : null}
            </div>
            <div className="grid gap-3 text-xs text-slate-400 sm:grid-cols-3">
              <div className="space-y-1">
                <p className="uppercase tracking-[0.18em]">
                  {t('workbench.header.currentStep')}
                </p>
                <p className="text-sm font-medium text-slate-100">
                  {t(WORKBENCH_STEP_META[activeStep].labelKey)}
                </p>
              </div>
              <div className="space-y-1">
                <p className="uppercase tracking-[0.18em]">
                  {t('workbench.header.selectedDevice')}
                </p>
                <p className="text-sm font-medium text-slate-100">
                  {selectedDeviceLabel}
                </p>
              </div>
              <div className="space-y-1">
                <p className="uppercase tracking-[0.18em]">
                  {t('workbench.header.lastTest')}
                </p>
                <p className="text-sm font-medium text-slate-100">{lastTestLabel}</p>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2 xl:justify-end">
            <button
              className="inline-flex items-center justify-center rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
              onClick={() => setActiveStep('device')}
              type="button"
            >
              {selectedDeviceId
                ? t('workbench.header.actions.switchDevice')
                : t('workbench.header.actions.selectDevice')}
            </button>
            <button
              className="inline-flex items-center justify-center rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
              disabled={!selectedDeviceId}
              onClick={() => setActiveStep('source')}
              type="button"
            >
              {t('workbench.header.actions.gotoSource')}
            </button>
            <button
              className="inline-flex items-center justify-center rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
              disabled={!selectedDeviceId}
              onClick={() => setActiveStep('output')}
              type="button"
            >
              {t('workbench.header.actions.gotoOutput')}
            </button>
          </div>
        </div>
      </div>

      <dl className="grid gap-4 rounded-2xl border border-slate-800 bg-slate-900/70 p-4 text-sm text-slate-200 sm:grid-cols-2 xl:grid-cols-3">
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.header.pointCount')}
          </dt>
          <dd
            className="font-medium text-slate-50"
            data-testid="workbench-header-point-count"
          >
            {pointCount}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.header.linkedTagCount')}
          </dt>
          <dd
            className="font-medium text-slate-50"
            data-testid="workbench-header-tag-count"
          >
            {linkedTagCount}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.header.outputCandidateCount')}
          </dt>
          <dd
            className="font-medium text-slate-50"
            data-testid="workbench-header-output-count"
          >
            {outputCandidateCount}
          </dd>
        </div>
      </dl>
    </header>
  );
}
