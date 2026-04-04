import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import { WB_SHELL_SURFACE } from './workbenchShellTokens';
import {
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
} from './workbenchDeviceFormModel';
import { useWorkbenchSummary } from './useWorkbenchSummary';
import type { WorkbenchStep } from './workbenchTypes';
import { WorkbenchStepRail } from './WorkbenchStepRail';

/**
 * 合併 Tailwind 類名，略過 falsy 值。
 *
 * @param classNames - 類名字串或條件略過值
 * @returns 單一空格分隔的類名字串
 */
function joinClasses(...classNames: Array<string | false | null | undefined>) {
  return classNames.filter(Boolean).join(' ');
}

/**
 * 依裝置狀態回傳 chip 邊框／底色／文字色類名。
 *
 * @param status - 草稿、啟用或停用
 * @returns Tailwind 類名字串
 */
function getStatusClasses(status: 'draft' | 'active' | 'disabled') {
  switch (status) {
    case 'active':
      return 'border-emerald-500/35 bg-emerald-500/[0.12] text-emerald-200';
    case 'disabled':
      return 'border-rose-500/35 bg-rose-500/[0.12] text-rose-200';
    case 'draft':
      return 'border-amber-500/35 bg-amber-500/[0.12] text-amber-200';
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
      className={joinClasses(
        'flex flex-wrap items-center gap-x-4 gap-y-3 px-5 py-3.5 sm:px-6',
        WB_SHELL_SURFACE,
        'ring-1 ring-inset ring-white/[0.04]',
      )}
      data-testid="workbench-context-bar"
    >
      <WorkbenchStepRail />

      <span
        className="hidden h-10 w-px shrink-0 bg-gradient-to-b from-transparent via-slate-600/70 to-transparent lg:block"
        aria-hidden
      />

      <div
        className="min-w-0 flex-1 basis-full lg:basis-0"
        data-testid="context-bar-step-summary"
        data-active-step={activeStep}
      >
        <div className="flex min-w-0 flex-wrap items-center gap-2 sm:gap-3">
          <span
            className="truncate text-[15px] font-semibold leading-tight tracking-tight text-slate-50"
            data-testid="context-bar-device-name"
          >
            {deviceLabel}
          </span>
          {selectedDevice ? (
            <>
              <span className="inline-flex h-6 shrink-0 items-center rounded-md border border-slate-700/80 bg-slate-950/60 px-2 text-[11px] font-medium text-slate-300">
                {t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))}
              </span>
              <span
                className={joinClasses(
                  'inline-flex h-6 shrink-0 items-center rounded-md border px-2 text-[11px] font-semibold',
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
        className={joinClasses(
          'min-h-10 w-full shrink-0 rounded-xl border border-cyan-500/40 bg-gradient-to-b from-cyan-500/15 to-cyan-500/[0.08]',
          'px-5 py-2 text-xs font-semibold text-cyan-50 shadow-sm shadow-cyan-950/20',
          'transition-[border-color,box-shadow,background-color] duration-200',
          'hover:border-cyan-400/55 hover:from-cyan-500/20 hover:to-cyan-500/10 hover:shadow-md hover:shadow-cyan-950/25',
          'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/45 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900',
          'disabled:cursor-not-allowed disabled:border-slate-700/90 disabled:bg-slate-950/60 disabled:bg-none disabled:text-slate-500 disabled:shadow-none sm:w-auto',
        )}
        disabled={primaryAction.disabled}
        onClick={primaryAction.onClick}
        type="button"
      >
        {primaryAction.label}
      </button>
    </header>
  );
}

/**
 * 依目前步驟與就緒狀態決定頂欄主按鈕文案、點擊行為與是否停用。
 *
 * @param input - 步驟、裝置選取、各段就緒與 i18n／setter
 * @returns 主按鈕的 label、onClick、disabled
 */
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
