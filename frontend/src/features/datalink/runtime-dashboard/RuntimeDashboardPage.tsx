import { useTranslation } from 'react-i18next';
import { CollectorHealthPanel } from './CollectorHealthPanel';
import { FocusedDeviceHeader } from './FocusedDeviceHeader';
import { LivePointsTable } from './LivePointsTable';
import { LiveStateBanner } from './LiveStateBanner';
import { RuntimeDiagnosticsPanel } from './RuntimeDiagnosticsPanel';
import { RuntimeSummaryPanel } from './RuntimeSummaryPanel';
import { RealtimeLogsPanel } from './RealtimeLogsPanel';
import type { RuntimeDashboardState } from './useRuntimeDashboardState';

export interface RuntimeDashboardPageProps extends RuntimeDashboardState {}

export function RuntimeDashboardPage({
  routeState,
  selectedDeviceId,
  selectedDevice,
  devices,
  snapshot,
  snapshotError,
  liveValues,
  streamState,
  logs = [],
  onSelectDevice,
  onRetrySnapshot,
}: RuntimeDashboardPageProps) {
  const { t } = useTranslation('runtime-dashboard');
  const collector =
    snapshot?.collectors.find((item) => item.device_id === selectedDeviceId) ?? null;

  if (routeState === 'empty-workspace') {
    return (
      <div
        className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100"
        data-testid="runtime-dashboard-empty-workspace"
      >
        <div className="mx-auto max-w-6xl space-y-6">
          <div className="sr-only" data-testid="runtime-dashboard-route-state">
            empty-workspace
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-selected-device">
            {selectedDevice?.name ?? selectedDeviceId ?? ''}
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-collector-count">
            {snapshot?.collectors.length ?? 0}
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-stream-state">{streamState}</div>
          <FocusedDeviceHeader
            selectedDeviceId={selectedDeviceId}
            selectedDevice={selectedDevice}
            devices={devices}
            onSelectDevice={onSelectDevice}
          />
          <section className="rounded-3xl border border-dashed border-slate-700 bg-slate-900/60 p-8">
            <h2 className="text-xl font-semibold text-slate-50">
              {t('emptyWorkspace.title', 'No available runtime device yet')}
            </h2>
            <p className="mt-2 text-sm text-slate-400">
              {t(
                'emptyWorkspace.description',
                'No available workspace device can be observed right now. Return to /studio/v2 to fix availability or activate another device.',
              )}
            </p>
            <a
              href="/studio/v2"
              className="mt-5 inline-flex rounded-xl border border-slate-700 bg-slate-950/40 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-slate-500"
            >
              {t('emptyWorkspace.cta', 'Return to Studio V2')}
            </a>
          </section>
        </div>
      </div>
    );
  }

  if (routeState === 'missing-device-context') {
    return (
      <div
        className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100"
        data-testid="runtime-dashboard-missing-device-context"
      >
        <div className="mx-auto max-w-6xl space-y-6">
          <div className="sr-only" data-testid="runtime-dashboard-route-state">
            missing-device-context
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-selected-device">
            {selectedDevice?.name ?? selectedDeviceId ?? ''}
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-collector-count">
            {snapshot?.collectors.length ?? 0}
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-stream-state">{streamState}</div>
          <FocusedDeviceHeader
            selectedDeviceId={selectedDeviceId}
            selectedDevice={selectedDevice}
            devices={devices}
            onSelectDevice={onSelectDevice}
          />
          <section className="rounded-3xl border border-dashed border-slate-700 bg-slate-900/60 p-8">
            <h2 className="text-xl font-semibold text-slate-50">
              {t('empty.title', 'Choose a runtime device')}
            </h2>
            <p className="mt-2 text-sm text-slate-400">
              {t(
                'empty.description',
                'This page stays focused on one committed device and does not switch to a fleet-wide dashboard.',
              )}
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
          <div className="sr-only" data-testid="runtime-dashboard-route-state">
            empty
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-selected-device">
            {selectedDevice?.name ?? selectedDeviceId ?? ''}
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-collector-count">
            {snapshot?.collectors.length ?? 0}
          </div>
          <div className="sr-only" data-testid="runtime-dashboard-stream-state">{streamState}</div>
          <FocusedDeviceHeader
            selectedDeviceId={selectedDeviceId}
            selectedDevice={selectedDevice}
            devices={devices}
            onSelectDevice={onSelectDevice}
          />
          <section
            className="rounded-3xl border border-dashed border-slate-700 bg-slate-900/60 p-8"
            data-testid="runtime-dashboard-empty-snapshot"
          >
            <h2 className="text-xl font-semibold text-slate-50">
              {t('emptySnapshot.title', 'Runtime data is not available yet')}
            </h2>
            <p className="mt-2 text-sm text-slate-400">
              {t(
                'emptySnapshot.description',
                'The backend has not produced a runtime snapshot for the selected device yet.',
              )}
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
        <div data-testid="runtime-dashboard-route-state">loading</div>
        <div className="sr-only" data-testid="runtime-dashboard-selected-device">
          {selectedDevice?.name ?? selectedDeviceId ?? ''}
        </div>
        <div className="sr-only" data-testid="runtime-dashboard-collector-count">0</div>
        <div className="sr-only" data-testid="runtime-dashboard-stream-state">{streamState}</div>
        <FocusedDeviceHeader
          selectedDeviceId={selectedDeviceId}
          selectedDevice={selectedDevice}
            devices={devices}
            onSelectDevice={onSelectDevice}
          />
          <section className="rounded-3xl border border-slate-800 bg-slate-900/70 p-8 text-sm text-slate-400">
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
        <div data-testid="runtime-dashboard-route-state">error</div>
        <div className="sr-only" data-testid="runtime-dashboard-selected-device">
          {selectedDevice?.name ?? selectedDeviceId ?? ''}
        </div>
        <div className="sr-only" data-testid="runtime-dashboard-collector-count">
          {snapshot?.collectors.length ?? 0}
        </div>
        <div className="sr-only" data-testid="runtime-dashboard-stream-state">{streamState}</div>
        <FocusedDeviceHeader
          selectedDeviceId={selectedDeviceId}
          selectedDevice={selectedDevice}
            devices={devices}
            onSelectDevice={onSelectDevice}
          />
          <section className="rounded-3xl border border-rose-500/30 bg-rose-500/10 p-8">
            <h2 className="text-xl font-semibold text-rose-100">
              {t('error.title', 'Runtime snapshot unavailable')}
            </h2>
            <p
              className="mt-2 text-sm text-rose-100/90"
              data-testid="runtime-dashboard-snapshot-error"
            >
              {snapshotError ?? t('error.description', 'The runtime snapshot failed to load.')}
            </p>
            <button
              type="button"
              onClick={() => void onRetrySnapshot()}
              className="mt-5 rounded-xl border border-rose-300/30 bg-slate-950/40 px-4 py-2 text-sm font-medium text-rose-50 transition hover:border-rose-200/50"
            >
              {t('error.retry', 'Retry snapshot')}
            </button>
          </section>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 px-6 py-8 text-slate-100" data-testid="runtime-dashboard-route">
      <div className="mx-auto max-w-6xl space-y-6">
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
        <FocusedDeviceHeader
          selectedDeviceId={selectedDeviceId}
          selectedDevice={selectedDevice}
          devices={devices}
          onSelectDevice={onSelectDevice}
        />
        <LiveStateBanner routeState={routeState} />
        {snapshot ? (
          <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
            <RuntimeSummaryPanel snapshot={snapshot} />
            <CollectorHealthPanel collector={collector} />
            <RuntimeDiagnosticsPanel
              diagnostics={snapshot.diagnostics}
              selectedDeviceId={selectedDeviceId}
            />
          </div>
        ) : null}
        <div className="grid gap-6 lg:grid-cols-12 items-start">
          <div className="lg:col-span-8">
            <LivePointsTable liveValues={liveValues} />
          </div>
          <div className="lg:col-span-4">
            <RealtimeLogsPanel logs={logs} />
          </div>
        </div>
      </div>
    </div>
  );
}
