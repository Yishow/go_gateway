import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { deviceKeys } from '../../../hooks/datalink/keys';
import {
  useDevicesQuery,
  useTestConnectionMutation,
} from '../../../hooks/datalink/useDevices';
import {
  buildDeviceCapabilitySummary,
  buildDeviceConnectionSummary,
  getDeviceTestTimestampLabel,
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
  parseDeviceConnectionConfig,
} from './workbenchDeviceFormModel';
import { useWorkbench } from './WorkbenchProvider';
import { WORKBENCH_STEP_META, type InspectorSelection } from './workbenchTypes';

function joinClasses(...classNames: Array<string | false | null | undefined>) {
  return classNames.filter(Boolean).join(' ');
}

function getDeviceStatusClasses(status: 'draft' | 'active' | 'disabled') {
  switch (status) {
    case 'active':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-300';
    case 'disabled':
      return 'border-rose-500/40 bg-rose-500/10 text-rose-300';
    case 'draft':
      return 'border-amber-500/40 bg-amber-500/10 text-amber-300';
  }
}

function getSelectionLabelKey(selection: InspectorSelection): string {
  switch (selection.kind) {
    case 'none':
      return 'workbench.inspector.noSelection';
    case 'device':
      return 'workbench.inspector.selectionKind.device';
    case 'rule':
      return 'workbench.inspector.selectionKind.rule';
    case 'span':
      return 'workbench.inspector.selectionKind.span';
    case 'tag':
      return 'workbench.inspector.selectionKind.tag';
    case 'outputCandidate':
      return 'workbench.inspector.selectionKind.outputCandidate';
  }
}

function getSelectionId(selection: InspectorSelection): string | null {
  switch (selection.kind) {
    case 'none':
      return null;
    case 'device':
      return selection.deviceId;
    case 'rule':
      return selection.ruleId;
    case 'span':
      return selection.spanAddress;
    case 'tag':
      return selection.tagId;
    case 'outputCandidate':
      return selection.tagId;
  }
}

function DeviceInspectorContent() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const testConnectionMutation = useTestConnectionMutation();
  const {
    openCloneDevicePanel,
    openCreateDevicePanel,
    openEditDevicePanel,
    recentDeviceTests,
    recordDeviceTest,
    selectedDeviceId,
  } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const selectedDevice =
    devices.find((device) => device.id === selectedDeviceId) ?? null;

  if (!selectedDevice) {
    return (
      <div className="flex flex-1 flex-col justify-center gap-4 rounded-xl border border-dashed border-slate-700/60 bg-slate-950/30 p-6">
      <div className="space-y-2">
        <h3 className="text-sm font-semibold text-slate-200">
          {t('workbench.device.inspector.emptyTitle')}
        </h3>
        <p className="text-xs text-slate-400">
          {t('workbench.device.inspector.emptyDescription')}
        </p>
      </div>
        {devices.length > 0 ? (
          <button
            className="rounded-xl bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950 transition hover:bg-cyan-400"
            onClick={openCreateDevicePanel}
            type="button"
          >
            {t('workbench.device.actions.create')}
          </button>
        ) : null}
      </div>
    );
  }

  const connectionConfig = parseDeviceConnectionConfig(selectedDevice.connection_config);
  const capabilitySummary = buildDeviceCapabilitySummary(
    selectedDevice.protocol,
    connectionConfig,
    t,
  );
  const connectionSummary = buildDeviceConnectionSummary(
    selectedDevice.protocol,
    connectionConfig,
  );
  const storedHistory = recentDeviceTests[selectedDevice.id] ?? [];
  const fallbackHistory =
    storedHistory.length === 0
      && (selectedDevice.last_test_success !== null || selectedDevice.last_test_at)
      ? [
          {
            id: `${selectedDevice.id}-latest`,
            testedAt: getDeviceTestTimestampLabel(selectedDevice.last_test_at, t),
            success: selectedDevice.last_test_success === true,
            message:
              selectedDevice.last_test_success === true
                ? t('workbench.device.card.testPassed')
                : selectedDevice.last_test_success === false
                  ? selectedDevice.last_test_error || t('workbench.device.card.testFailed')
                  : t('workbench.device.inspector.noTestYet'),
            latencyMs: null,
          },
        ]
      : [];
  const recentHistory = [...storedHistory, ...fallbackHistory].slice(0, 3);

  const handleTestConnection = async () => {
    try {
      const result = await testConnectionMutation.mutateAsync(selectedDevice.id);
      await queryClient.invalidateQueries({ queryKey: deviceKeys.lists() });
      recordDeviceTest(selectedDevice.id, {
        testedAt: new Date().toISOString(),
        success: result.success,
        message: result.success
          ? t('workbench.device.messages.testSuccess')
          : result.error || t('workbench.device.messages.testFailed'),
        latencyMs: result.latency_ms ?? null,
      });
    } catch (error) {
      recordDeviceTest(selectedDevice.id, {
        testedAt: new Date().toISOString(),
        success: false,
        message:
          error instanceof Error
            ? error.message
            : t('workbench.device.messages.testFailed'),
        latencyMs: null,
      });
    }
  };

  return (
    <div className="flex h-full flex-col gap-4">
      <div className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <div className="space-y-2">
          <h3 className="text-lg font-semibold text-slate-50">{selectedDevice.name}</h3>
          <p className="text-xs text-slate-400">
            {selectedDevice.description || t('workbench.device.card.noDescription')}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <span className="rounded-full border border-slate-700 bg-slate-950/70 px-2 py-0.5 text-[11px] text-slate-300">
            {t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))}
          </span>
          <span
            className={joinClasses(
              'rounded-full border px-2 py-0.5 text-[11px] font-medium',
              getDeviceStatusClasses(selectedDevice.status),
            )}
          >
            {t(getWorkbenchDeviceStatusLabelKey(selectedDevice.status))}
          </span>
        </div>
      </div>

      <section className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.device.inspector.capabilitySummary')}
        </p>
        <dl className="space-y-3 text-sm text-slate-200">
          {capabilitySummary.map((item) => (
            <div className="space-y-1" key={item.id}>
              <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
                {t(item.labelKey)}
              </dt>
              <dd className="font-medium text-slate-100">{item.value}</dd>
            </div>
          ))}
        </dl>
      </section>

      <section className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.device.inspector.connectionSummary')}
        </p>
        <dl className="space-y-3 text-sm text-slate-200">
          {connectionSummary.map((item) => (
            <div className="flex items-center justify-between gap-4" key={item.labelKey}>
              <dt className="text-slate-400">{t(item.labelKey)}</dt>
              <dd className="text-right font-medium text-slate-100">{item.value}</dd>
            </div>
          ))}
        </dl>
      </section>

      <section className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.device.inspector.recentTests')}
        </p>
        {recentHistory.length > 0 ? (
          <ol className="space-y-3">
            {recentHistory.map((entry) => (
              <li
                className="rounded-xl border border-slate-800 bg-slate-950/70 p-3"
                key={entry.id}
              >
                <p
                  className={joinClasses(
                    'text-sm font-medium',
                    entry.success ? 'text-emerald-300' : 'text-rose-300',
                  )}
                >
                  {entry.message}
                </p>
                <p className="mt-1 text-xs text-slate-500">
                  {getDeviceTestTimestampLabel(entry.testedAt, t)}
                </p>
              </li>
            ))}
          </ol>
        ) : (
          <p className="text-xs text-slate-400">{t('workbench.device.inspector.noTestYet')}</p>
        )}
      </section>

      <div className="mt-auto grid gap-2">
        <button
          className="rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
          onClick={() => openEditDevicePanel(selectedDevice.id)}
          type="button"
        >
          {t('workbench.device.actions.edit')}
        </button>
        <button
          className="rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
          onClick={() => void handleTestConnection()}
          type="button"
        >
          {testConnectionMutation.isPending
            ? t('workbench.device.actions.testing')
            : t('workbench.device.actions.testConnection')}
        </button>
        <button
          className="rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
          onClick={() => openCloneDevicePanel(selectedDevice.id)}
          type="button"
        >
          {t('workbench.device.actions.clone')}
        </button>
      </div>
    </div>
  );
}

export function WorkbenchInspectorPanel() {
  const { t } = useTranslation();
  const { activeStep, inspectorSelection } = useWorkbench();
  const stepMeta = WORKBENCH_STEP_META[activeStep];
  const hasSelection = inspectorSelection.kind !== 'none';
  const selectionId = getSelectionId(inspectorSelection);

  return (
    <aside
      aria-label={t('workbench.inspector.ariaLabel')}
      className="flex flex-col gap-4 rounded-2xl border border-slate-800 bg-slate-900/80 p-4"
      data-testid="workbench-inspector-panel"
    >
      <div className="space-y-1">
        <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-slate-500">
          {t('workbench.inspector.heading')}
        </p>
        <h3 className="text-sm font-semibold text-slate-200">
          {t(stepMeta.labelKey)}
        </h3>
      </div>

      {activeStep === 'device' ? (
        <DeviceInspectorContent />
      ) : hasSelection ? (
        <div
          className="space-y-2 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4"
          data-testid="inspector-selection-context"
        >
          <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
            {t(getSelectionLabelKey(inspectorSelection))}
          </p>
          {selectionId ? (
            <p
              className="truncate text-xs font-medium text-slate-200"
              data-testid="inspector-selection-id"
            >
              {selectionId}
            </p>
          ) : null}
        </div>
      ) : (
        <div className="flex flex-1 items-center justify-center rounded-xl border border-dashed border-slate-700/60 bg-slate-950/30 p-6">
          <p className="text-center text-xs text-slate-500" data-testid="inspector-placeholder">
            {t('workbench.inspector.placeholder')}
          </p>
        </div>
      )}
    </aside>
  );
}
