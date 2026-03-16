import { useTranslation } from 'react-i18next';
import { useWorkbench } from './WorkbenchProvider';
import {
  buildDeviceCapabilitySummary,
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
  parseDeviceConnectionConfig,
} from './workbenchDeviceFormModel';
import { useWorkbenchSummary } from './useWorkbenchSummary';

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
  const { recentDeviceTests, selectedDeviceId, setActiveStep } = useWorkbench();
  const { selectedDevice } = useWorkbenchSummary();
  const capabilitySummary = selectedDevice
    ? buildDeviceCapabilitySummary(
        selectedDevice.protocol,
        parseDeviceConnectionConfig(selectedDevice.connection_config),
        t,
      )
    : [];
  const latestTestEntry = selectedDevice ? recentDeviceTests[selectedDevice.id]?.[0] : undefined;

  const deviceLabel = selectedDevice?.name ?? t('workbench.contextBar.noDevice');
  const lastTestLabel = selectedDevice
    ? latestTestEntry
      ? latestTestEntry.message
      : selectedDevice.last_test_success === true
        ? t('workbench.contextBar.testPassed')
        : selectedDevice.last_test_success === false
          ? selectedDevice.last_test_error || t('workbench.contextBar.testFailed')
          : t('workbench.contextBar.notTested')
    : null;
  const lastTestToneClass = selectedDevice
    ? selectedDevice.last_test_success === true
      ? 'text-emerald-300'
      : selectedDevice.last_test_success === false
        ? 'text-rose-300'
        : 'text-slate-400'
    : 'text-slate-400';

  return (
    <header
      aria-label={t('workbench.contextBar.ariaLabel')}
      className="flex items-center gap-4 rounded-2xl border border-slate-800 bg-slate-900/80 px-4 py-2.5"
      data-testid="workbench-context-bar"
    >
      {/* Device identity */}
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

      {capabilitySummary.length > 0 ? (
        <div className="hidden min-w-0 flex-1 items-center gap-2 overflow-hidden xl:flex">
          {capabilitySummary.map((item) => (
            <span
              key={item.id}
              className="truncate rounded-full border border-slate-700/70 bg-slate-950/60 px-2 py-0.5 text-[11px] text-slate-300"
              data-testid={`context-bar-capability-${item.id}`}
              title={`${t(item.labelKey)}: ${item.value}`}
            >
              {item.value}
            </span>
          ))}
        </div>
      ) : (
        <div className="flex-1" />
      )}

      {/* Connection state */}
      {lastTestLabel ? (
        <span
          className={`hidden shrink-0 text-xs lg:inline ${lastTestToneClass}`}
          data-testid="context-bar-test-status"
        >
          {lastTestLabel}
        </span>
      ) : null}

      {/* Quick actions */}
      <div className="flex items-center gap-2">
        <button
          className="rounded-lg border border-slate-700/60 bg-slate-950/50 px-2.5 py-1.5 text-xs font-medium text-slate-300 transition hover:border-slate-600 hover:text-slate-100"
          onClick={() => setActiveStep('device')}
          type="button"
        >
          {selectedDeviceId
            ? t('workbench.contextBar.actions.switchDevice')
            : t('workbench.contextBar.actions.selectDevice')}
        </button>
        <button
          className="rounded-lg border border-slate-700/60 bg-slate-950/50 px-2.5 py-1.5 text-xs font-medium text-slate-300 transition hover:border-slate-600 hover:text-slate-100"
          disabled={!selectedDeviceId}
          onClick={() => setActiveStep('source')}
          type="button"
        >
          {t('workbench.contextBar.actions.gotoSource')}
        </button>
        <button
          className="rounded-lg border border-slate-700/60 bg-slate-950/50 px-2.5 py-1.5 text-xs font-medium text-slate-300 transition hover:border-slate-600 hover:text-slate-100"
          disabled={!selectedDeviceId}
          onClick={() => setActiveStep('output')}
          type="button"
        >
          {t('workbench.contextBar.actions.gotoOutput')}
        </button>
      </div>
    </header>
  );
}
