import type { RuntimeSourceSelection } from '../../../types/studioV2RuntimeContext';
import { useId, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Activity, Layers, BarChart3 } from 'lucide-react';
import { CollectorHealthPanel } from './CollectorHealthPanel';
import { FocusedDeviceHeader } from './FocusedDeviceHeader';
import { LivePointsTable } from './LivePointsTable';
import { LiveStateBanner } from './LiveStateBanner';
import { RuntimeDiagnosticsPanel } from './RuntimeDiagnosticsPanel';
import { RuntimeSummaryPanel } from './RuntimeSummaryPanel';
import { RealtimeLogsPanel } from './RealtimeLogsPanel';
import { RuntimeSetupContextPanel } from './RuntimeSetupContextPanel';
import { RuntimeRecordingRulesPanel } from './RuntimeRecordingRulesPanel';
import { HistoryReportsPanel } from './HistoryReportsPanel';
import type { RuntimeDashboardState } from './useRuntimeDashboardState';

export interface RuntimeDashboardPageProps extends RuntimeDashboardState {
  navigateTo?: (target: string) => void;
  sourceSelection?: RuntimeSourceSelection | null;
  onSourceSelection?: (selection: RuntimeSourceSelection | null) => void;
}

export function RuntimeDashboardPage({
  routeState,
  selectedDeviceId,
  selectedDevice,
  devices,
  setupContext,
  sourceEvidenceScope = null,
  snapshot,
  snapshotError,
  snapshotReceivedAt = null,
  snapshotStale = false,
  liveValues,
  streamState,
  logs = [],
  streamRecovery,
  navigateTo,
  sourceSelection = null,
  onSourceSelection,
  onSelectDevice,
  onRetrySnapshot,
  onReconnectStream,
}: RuntimeDashboardPageProps) {
  const { t } = useTranslation('runtime-dashboard');
  const [reportsExpanded, setReportsExpanded] = useState(false);
  const reportsId = useId();
  const collector =
    snapshot?.collectors.find((item) => item.device_id === selectedDeviceId) ?? null;

  const header = (
    <FocusedDeviceHeader
      selectedDeviceId={selectedDeviceId}
      streamState={streamState}
      selectedDevice={selectedDevice}
      devices={devices}
      onSelectDevice={onSelectDevice}
    />
  );

  const recordingPanel = <RuntimeRecordingRulesPanel key={JSON.stringify([sourceEvidenceScope?.workspace_id, sourceEvidenceScope?.workspace_revision, sourceEvidenceScope?.device_id])} scope={sourceEvidenceScope} selectedDeviceId={selectedDeviceId} navigateTo={navigateTo} />;

  const setupPanel = (
    <RuntimeSetupContextPanel
      setupContext={setupContext}
      evidenceScope={sourceEvidenceScope}
      sourceSelection={sourceSelection}
      onSourceSelection={onSourceSelection}
      selectedDevice={selectedDevice}
      navigateTo={navigateTo}
    />
  );

  const testMeta = (
    <>
      <div className="sr-only" data-testid="runtime-dashboard-route-state">
        {routeState}
      </div>
      <div className="sr-only" data-testid="runtime-dashboard-selected-device">
        {selectedDevice?.name ?? selectedDeviceId ?? ''}
      </div>
      <div className="sr-only" data-testid="runtime-dashboard-collector-count">
        {snapshot?.collectors.length ?? 0}
      </div>
      <div className="sr-only" data-testid="runtime-dashboard-stream-state">
        {streamState}
      </div>
    </>
  );

  if (routeState === 'empty-workspace') {
    return (
      <div className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100" data-testid="runtime-dashboard-empty-workspace">
        <div className="mx-auto max-w-6xl space-y-6">
          {testMeta}
          {header}
          {setupPanel}
          {recordingPanel}
          <section className="rounded-3xl border border-dashed border-slate-800 bg-slate-900/60 p-8 text-center">
            <h2 className="text-xl font-semibold text-slate-50">{t('emptyWorkspace.title', 'No available runtime device yet')}</h2>
            <p className="mt-2 text-sm text-slate-400 max-w-lg mx-auto">
              {t('emptyWorkspace.description', 'No available workspace device can be observed right now. Return to /studio/v2 to fix availability or activate another device.')}
            </p>
            <a href="/studio/v2" className="mt-5 inline-flex rounded-xl border border-slate-700 bg-slate-950/60 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-slate-500">
              {t('emptyWorkspace.cta', 'Return to Studio V2')}
            </a>
          </section>
        </div>
      </div>
    );
  }

  if (routeState === 'missing-device-context') {
    return (
      <div className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100" data-testid="runtime-dashboard-missing-device-context">
        <div className="mx-auto max-w-6xl space-y-6">
          {testMeta}
          {header}
          {setupPanel}
          {recordingPanel}
          <section className="rounded-3xl border border-dashed border-slate-800 bg-slate-900/60 p-8 text-center">
            <h2 className="text-xl font-semibold text-slate-50">{t('empty.title', 'Choose a runtime device')}</h2>
            <p className="mt-2 text-sm text-slate-400 max-w-lg mx-auto">
              {t('empty.description', 'This page stays focused on one committed device and does not switch to a fleet-wide dashboard.')}
            </p>
          </section>
        </div>
      </div>
    );
  }

  if (routeState === 'empty') {
    return (
      <div className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100" data-testid="runtime-dashboard-route">
        <div className="mx-auto max-w-6xl space-y-6">
          {testMeta}
          {header}
          {setupPanel}
          {recordingPanel}
          <section className="rounded-3xl border border-dashed border-slate-800 bg-slate-900/60 p-8 text-center" data-testid="runtime-dashboard-empty-snapshot">
            <h2 className="text-xl font-semibold text-slate-50">{t('emptySnapshot.title', 'Runtime data is not available yet')}</h2>
            <p className="mt-2 text-sm text-slate-400 max-w-lg mx-auto">
              {t('emptySnapshot.description', 'The backend has not produced a runtime snapshot for the selected device yet.')}
            </p>
          </section>
        </div>
      </div>
    );
  }

  if (routeState === 'loading') {
    return (
      <div className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100" data-testid="runtime-dashboard-route">
        <div className="mx-auto max-w-6xl space-y-6">
          <div data-testid="runtime-dashboard-route-state" className="sr-only">loading</div>
          <div className="sr-only" data-testid="runtime-dashboard-selected-device">{selectedDevice?.name ?? selectedDeviceId ?? ''}</div>
          <div className="sr-only" data-testid="runtime-dashboard-collector-count">0</div>
          <div className="sr-only" data-testid="runtime-dashboard-stream-state">{streamState}</div>
          {header}
          {setupPanel}
          {recordingPanel}
          <section className="rounded-3xl border border-slate-800 bg-slate-900/70 p-8 text-sm text-slate-400 text-center">
            {t('loading.message', 'Loading the latest runtime snapshot...')}
          </section>
        </div>
      </div>
    );
  }

  if (routeState === 'error') {
    return (
      <div className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100" data-testid="runtime-dashboard-route">
        <div className="mx-auto max-w-6xl space-y-6">
          <div data-testid="runtime-dashboard-route-state" className="sr-only">error</div>
          <div className="sr-only" data-testid="runtime-dashboard-selected-device">{selectedDevice?.name ?? selectedDeviceId ?? ''}</div>
          <div className="sr-only" data-testid="runtime-dashboard-collector-count">{snapshot?.collectors.length ?? 0}</div>
          <div className="sr-only" data-testid="runtime-dashboard-stream-state">{streamState}</div>
          {header}
          {setupPanel}
          {recordingPanel}
          <section className="rounded-3xl border border-rose-500/30 bg-rose-500/10 p-8">
            <h2 className="text-xl font-semibold text-rose-100">{t('error.title', 'Runtime snapshot unavailable')}</h2>
            {snapshotError?.code && (
              <p className="mt-2 font-mono text-xs text-rose-200/80" data-testid="runtime-dashboard-error-code">
                {snapshotError.code}
              </p>
            )}
            <p className="mt-2 text-sm text-rose-100/90" data-testid="runtime-dashboard-snapshot-error">
              {snapshotError?.message ?? t('error.description', 'The runtime snapshot failed to load.')}
            </p>
            {snapshotError?.requestId && (
              <p className="mt-2 text-xs text-rose-100/70" data-testid="runtime-dashboard-error-request-id">
                {t('error.requestId', 'Request ID')}: {snapshotError.requestId}
              </p>
            )}
            {snapshotError?.code === 'runtime_device_not_found' ? (
              <a href="/studio/v2" className="mt-5 inline-flex rounded-xl border border-rose-300/30 bg-slate-950/40 px-4 py-2 text-sm font-medium text-rose-50 transition hover:border-rose-200/50">
                {t('error.returnToStudio', 'Return to Studio V2')}
              </a>
            ) : (
              <button type="button" onClick={() => void onRetrySnapshot()} className="mt-5 rounded-xl border border-rose-300/30 bg-slate-950/40 px-4 py-2 text-sm font-medium text-rose-50 transition hover:border-rose-200/50">
                {t('error.retry', 'Retry snapshot')}
              </button>
            )}
          </section>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100" data-testid="runtime-dashboard-route">
      <div className="mx-auto max-w-7xl space-y-6">
        {testMeta}
        {header}
        {setupPanel}
        {recordingPanel}
        <LiveStateBanner
          routeState={routeState}
          snapshotReceivedAt={snapshotReceivedAt}
          snapshotStale={snapshotStale}
          streamRecovery={streamRecovery}
          onRetry={onRetrySnapshot}
          onReconnect={onReconnectStream}
          navigateTo={navigateTo}
        />

        {/* 區塊 1：全鏈路管線監控總覽 (整合 Summary + Health + Diagnostics) */}
        {snapshot ? (
          <section className="space-y-4">
            <div className="flex items-center gap-2 px-1 text-xs font-semibold uppercase tracking-wider text-cyan-400">
              <Activity className="h-4 w-4" />
              <span>{t('focused.health', 'Collector health and runtime counters')}</span>
            </div>
            <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
              <RuntimeSummaryPanel snapshot={snapshot} />
              <CollectorHealthPanel collector={collector} />
            </div>
            <RuntimeDiagnosticsPanel
              diagnostics={snapshot.diagnostics}
              databaseDelivery={snapshot.database_delivery}
              modbusShareDelivery={snapshot.modbus_share_delivery}
              selectedDeviceId={selectedDeviceId}
            />
          </section>
        ) : null}

        {/* 區塊 2：即時數據觀測中心 (整合 LivePointsTable + RealtimeLogsPanel) */}
        <section className="space-y-4">
          <div className="flex items-center gap-2 px-1 text-xs font-semibold uppercase tracking-wider text-emerald-400">
            <Layers className="h-4 w-4" />
            <span>{t('focused.values', 'Point values and connection events')}</span>
          </div>
          <div className="grid gap-6 lg:grid-cols-12 items-start">
            <div className="lg:col-span-8">
              <LivePointsTable selectedDeviceId={selectedDeviceId} liveValues={liveValues} setupContext={setupContext} />
            </div>
            <div className="lg:col-span-4">
              <RealtimeLogsPanel logs={logs} routeState={routeState} />
            </div>
          </div>
        </section>

        <section aria-label={t('focused.reportsTitle', 'Historical reports')} className="rounded-3xl border border-slate-800/80 bg-slate-900/40 p-4 sm:p-6">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <h2 className="flex items-center gap-2 text-sm font-semibold text-slate-300">
              <BarChart3 className="h-4 w-4 text-indigo-400" />
              {t('focused.reportsTitle', 'Historical reports')}
            </h2>
            <button type="button" aria-expanded={reportsExpanded} aria-controls={reportsId} onClick={() => setReportsExpanded((expanded) => !expanded)} className="rounded-lg border border-slate-700 px-3 py-2 text-xs text-slate-200 focus-visible:outline focus-visible:outline-cyan-400">
              {reportsExpanded ? t('focused.hideReports', 'Hide historical reports') : t('focused.showReports', 'Show historical reports')}
            </button>
          </div>
          {/* Keep the existing mounted reader and its query budget; only visibility changes. */}
          <div id={reportsId} hidden={!reportsExpanded} className="mt-4">
            <HistoryReportsPanel selectedDeviceId={selectedDeviceId} />
          </div>
        </section>
      </div>
    </div>
  );
}
