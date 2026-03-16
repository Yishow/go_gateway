import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import {
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
} from './workbenchDeviceFormModel';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import type { WorkbenchStep } from './workbenchTypes';

function joinClasses(...classNames: Array<string | false | null | undefined>) {
  return classNames.filter(Boolean).join(' ');
}

function getStatusClasses(status: 'draft' | 'active' | 'disabled') {
  switch (status) {
    case 'active':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-300';
    case 'disabled':
      return 'border-rose-500/40 bg-rose-500/10 text-rose-300';
    case 'draft':
      return 'border-amber-500/40 bg-amber-500/10 text-amber-300';
  }
}

export function WorkbenchContextBar() {
  const { t } = useTranslation();
  const { activeStep, selectedDeviceId, setActiveStep } = useWorkbench();
  const { selectedDevice, sourceReady, tagReady, outputReady } = useWorkbenchSummary();
  const deviceLabel = selectedDevice?.name ?? t('workbench.contextBar.noDevice');
  const primaryAction = getPrimaryAction({
    activeStep,
    hasSelectedDevice: Boolean(selectedDeviceId),
    sourceReady,
    tagReady,
    outputReady,
    t,
    setActiveStep,
  });

  return (
    <header
      aria-label={t('workbench.contextBar.ariaLabel')}
      className="flex items-center justify-between gap-4 rounded-2xl border border-slate-800 bg-slate-900/80 px-4 py-2.5"
      data-testid="workbench-context-bar"
    >
      <div className="min-w-0 space-y-1" data-testid="context-bar-step-summary">
        <p className="text-[11px] font-semibold uppercase tracking-[0.22em] text-slate-500">
          {t(`workbench.steps.${activeStep}`)}
        </p>
        <div className="flex min-w-0 items-center gap-2.5">
          <span className="truncate text-sm font-semibold text-slate-100" data-testid="context-bar-device-name">
            {deviceLabel}
          </span>
          {selectedDevice ? (
            <>
              <span className="shrink-0 rounded-full border border-slate-700 bg-slate-950/70 px-2 py-0.5 text-[11px] text-slate-300">
                {t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))}
              </span>
              <span
                className={joinClasses(
                  'shrink-0 rounded-full border px-2 py-0.5 text-[11px] font-medium',
                  getStatusClasses(selectedDevice.status),
                )}
              >
                {t(getWorkbenchDeviceStatusLabelKey(selectedDevice.status))}
              </span>
            </>
          ) : null}
        </div>
      </div>

      <button
        className="rounded-lg border border-cyan-500/40 bg-cyan-500/10 px-3 py-1.5 text-xs font-semibold text-cyan-100 transition hover:border-cyan-400/60 hover:bg-cyan-500/15 disabled:cursor-not-allowed disabled:border-slate-700 disabled:bg-slate-950/40 disabled:text-slate-500"
        disabled={primaryAction.disabled}
        onClick={primaryAction.onClick}
        type="button"
      >
        {primaryAction.label}
      </button>
    </header>
  );
}

function getPrimaryAction(
  input: {
    activeStep: WorkbenchStep;
    hasSelectedDevice: boolean;
    sourceReady: boolean;
    tagReady: boolean;
    outputReady: boolean;
    t: (key: string) => string;
    setActiveStep: (step: WorkbenchStep) => void;
  },
) {
  if (!input.hasSelectedDevice) {
    return {
      label: input.t('workbench.contextBar.actions.selectDevice'),
      onClick: () => input.setActiveStep('device'),
      disabled: false,
    };
  }

  if (input.activeStep === 'device') {
    return {
      label: input.t('workbench.contextBar.actions.gotoSource'),
      onClick: () => input.setActiveStep('source'),
      disabled: false,
    };
  }

  if (input.activeStep === 'source') {
    return {
      label: input.t('workbench.contextBar.actions.gotoTag'),
      onClick: () => input.setActiveStep('tag'),
      disabled: !input.sourceReady,
    };
  }

  if (input.activeStep === 'tag') {
    return {
      label: input.t('workbench.contextBar.actions.gotoOutput'),
      onClick: () => input.setActiveStep('output'),
      disabled: !input.tagReady,
    };
  }

  return {
    label: input.outputReady
      ? input.t('workbench.actionDock.nextAction.configureOutput')
      : input.t('workbench.actionDock.nextAction.completeTagBinding'),
    onClick: input.outputReady
      ? () => {
          const anchor = document.querySelector<HTMLElement>('[data-testid="output-primary-anchor"]');
          anchor?.focus();
          anchor?.scrollIntoView?.({ block: 'nearest' });
        }
      : () => input.setActiveStep('tag'),
    disabled: false,
  };
}
