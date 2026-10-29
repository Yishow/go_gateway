import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useLocation, useSearchParams } from 'react-router-dom';
import { safeRuntimeIdentity } from '../../../utils/runtimeOperatorSafety';
import { parseRuntimeDeviceSelection } from './runtimeDeviceSelection';
import { useStudioV2RuntimeContextQuery } from '../../../hooks/datalink/useStudioV2RuntimeContext';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import type { RuntimeValueEvent } from '../../../types/datalink';
import type {
  StudioV2RuntimeContextDevice,
  StudioV2RuntimeSetupContext,
} from '../../../types/studioV2RuntimeContext';
import type { EvidenceScope } from '../../../types/runtimeRuleEvidence';
import type { RuntimeStatusWithDiagnostics } from '../../../types/runtimeDiagnostics';
import { getSafeErrorMessage, type SafeErrorMessage } from '../../../utils/typedErrors';
import { useRuntimeStatus } from './useRuntimeStatus';
import { useRuntimeDashboardStream } from './useRuntimeStream';
import type {
  RuntimeDashboardLog,
  RuntimeDashboardStreamConnectionState,
} from './useRuntimeStream';
import type { RuntimeStreamRecovery } from '../../../types/datalink';

const degradedPollingIntervalMs = 5000;
// SSE 只推送點位數值；collector 健康與累計指標仍須靠 snapshot，live 時以較低頻率更新。
const liveSnapshotRefreshIntervalMs = 15000;

export type RuntimeDashboardRouteState =
  | 'missing-device-context'
  | 'empty-workspace'
  | 'empty'
  | 'loading'
  | 'live'
  | 'degraded'
  | 'error';

export interface RuntimeDashboardState {
  routeState: RuntimeDashboardRouteState;
  selectedDeviceId: string | null;
  selectedDevice: StudioV2RuntimeContextDevice | null;
  devices: StudioV2RuntimeContextDevice[];
  setupContext: StudioV2RuntimeSetupContext | null;
  sourceEvidenceScope?: EvidenceScope | null;
  snapshot: RuntimeStatusWithDiagnostics | null;
  snapshotReceivedAt?: number | null;
  snapshotStale?: boolean;
  snapshotError: SafeErrorMessage | null;
  liveValues: Record<string, RuntimeValueEvent>;
  streamState: RuntimeDashboardStreamConnectionState;
  logs?: RuntimeDashboardLog[];
  streamRecovery?: RuntimeStreamRecovery | null;
  onSelectDevice: (deviceId: string) => void;
  onRetrySnapshot: () => Promise<unknown>;
  onReconnectStream: () => void;
}

export function useRuntimeDashboardState(): RuntimeDashboardState {
  const { t: translateError } = useTranslation('workbench-v2');
  const [searchParams, setSearchParams] = useSearchParams();
  const { search } = useLocation();
  const deviceSelection = useMemo(() => parseRuntimeDeviceSelection(search), [search]);
  const [isDegraded, setIsDegraded] = useState(false);

  const runtimeContextQuery = useStudioV2RuntimeContextQuery();
  const devices = useMemo(
    () => (runtimeContextQuery.data?.devices ?? []).filter((device) => Boolean(safeRuntimeIdentity(device.device_id))),
    [runtimeContextQuery.data],
  );
  const setupContext = runtimeContextQuery.data?.setup ?? null;
  const requestedDeviceId = deviceSelection.kind === 'valid' ? deviceSelection.deviceId
    : deviceSelection.kind === 'absent' ? safeRuntimeIdentity(runtimeContextQuery.data?.default_device_id) || null : null;
  const selectedDevice = useMemo(
    () => devices.find((device) => device.device_id === requestedDeviceId) ?? null,
    [devices, requestedDeviceId],
  );
  // A syntactically safe selector still is not an observed workspace device.
  const selectedDeviceId = selectedDevice?.device_id ?? null;
  const selectedDeviceContextMissing = deviceSelection.kind === 'invalid' || Boolean(
    runtimeContextQuery.data && requestedDeviceId && !selectedDevice,
  );
  const observableDeviceId = selectedDeviceId;

  const pointsQuery = usePointsQuery(
    observableDeviceId ? { device_id: observableDeviceId } : undefined,
    { enabled: Boolean(observableDeviceId) },
  );
  const pointIds = useMemo(
    () => observableDeviceId ? (pointsQuery.data ?? []).map((point) => point.id) : [],
    [observableDeviceId, pointsQuery.data],
  );

  const snapshotQuery = useRuntimeStatus({
    deviceId: observableDeviceId,
    pollingIntervalMs: isDegraded ? degradedPollingIntervalMs : liveSnapshotRefreshIntervalMs,
  });
  const snapshotError = useMemo(() => {
    const error = runtimeContextQuery.error ?? snapshotQuery.error;
    return error ? getSafeErrorMessage(error, translateError) : null;
  }, [runtimeContextQuery.error, snapshotQuery.error, translateError]);

  const stream = useRuntimeDashboardStream({
    deviceId: observableDeviceId,
    pointIds,
  });

  // React Query retains the last successful data and its original receipt on a
  // failed refetch. Its device-keyed observer supplies an exact synchronous pair;
  // do not copy an unkeyed previous snapshot through a device-switch effect.
  const candidate = snapshotQuery.data as RuntimeStatusWithDiagnostics | undefined;
  const ownedSnapshot = candidate && candidate.collectors.every((collector) => collector.device_id === observableDeviceId);
  const currentSnapshot = observableDeviceId && ownedSnapshot ? candidate : null;
  const snapshotReceivedAt = currentSnapshot && Number.isFinite(snapshotQuery.dataUpdatedAt) && snapshotQuery.dataUpdatedAt > 0
    ? snapshotQuery.dataUpdatedAt : null;
  const failure = useRef({ deviceId: observableDeviceId, stream: stream.failureVersion, error: 0, receiptFloor: null as number | null });
  if (failure.current.deviceId !== observableDeviceId) {
    failure.current = { deviceId: observableDeviceId, stream: stream.failureVersion, error: 0, receiptFloor: null };
  }
  if (failure.current.stream !== stream.failureVersion ||
    (snapshotQuery.isError && failure.current.error !== snapshotQuery.errorUpdatedAt)) {
    failure.current = { deviceId: observableDeviceId, stream: stream.failureVersion,
      error: snapshotQuery.errorUpdatedAt, receiptFloor: snapshotQuery.dataUpdatedAt };
  }
  const snapshotState = currentSnapshot?.snapshot_state;
  const snapshotReady = Boolean(currentSnapshot?.collectors.length) && snapshotState?.state === 'ready' && snapshotState.empty === false &&
    snapshotState.degraded === false && snapshotState.unavailable === false && snapshotState.stale === false;
  // Same-millisecond receipts are not unique request IDs: conservatively require
  // a later successful receipt after failure, without inventing elapsed-age TTLs.
  const retainedAfterFailure = failure.current.receiptFloor !== null &&
    (snapshotReceivedAt === null || snapshotReceivedAt <= failure.current.receiptFloor);
  const snapshotStale = Boolean(currentSnapshot) && (!snapshotQuery.isFetchedAfterMount || snapshotReceivedAt === null || !snapshotReady || !stream.serverReady ||
    snapshotQuery.isError || retainedAfterFailure);
  useEffect(() => { setIsDegraded(snapshotStale); }, [snapshotStale]);

  const routeState = useMemo<RuntimeDashboardRouteState>(() => {
    if (!runtimeContextQuery.data && (runtimeContextQuery.isLoading || runtimeContextQuery.isFetching)) {
      return 'loading';
    }

    if (!runtimeContextQuery.data && runtimeContextQuery.isError) {
      return 'error';
    }

    if (selectedDeviceContextMissing) {
      return 'missing-device-context';
    }

    if (!selectedDeviceId) {
      return 'empty-workspace';
    }

    if (!currentSnapshot && snapshotError?.code === 'runtime_device_not_found') {
      return 'missing-device-context';
    }

    if (candidate && !ownedSnapshot && !snapshotQuery.isFetching) return 'error';

    if (!currentSnapshot && snapshotQuery.isError) {
      return 'error';
    }

    if (snapshotState?.empty) {
      return 'empty';
    }

    if (!currentSnapshot && (snapshotQuery.isLoading || snapshotQuery.isFetching)) {
      return 'loading';
    }

    if (
      snapshotState?.degraded ||
      snapshotState?.unavailable ||
      snapshotState?.stale ||
      snapshotStale ||
      stream.connectionState === 'unavailable' ||
      stream.connectionState === 'degraded' ||
      stream.connectionState === 'error'
    ) {
      return 'degraded';
    }

    if (stream.serverReady && snapshotReady && snapshotReceivedAt !== null && !snapshotStale) {
      return 'live';
    }

    return 'loading';
  }, [
    candidate,
    ownedSnapshot,
    currentSnapshot,
    snapshotReady,
    snapshotReceivedAt,
    snapshotStale,
    stream.serverReady,
    runtimeContextQuery.data,
    runtimeContextQuery.isError,
    runtimeContextQuery.isFetching,
    runtimeContextQuery.isLoading,
    selectedDeviceContextMissing,
    selectedDeviceId,
    snapshotError?.code,
    snapshotState?.degraded,
    snapshotState?.empty,
    snapshotState?.stale,
    snapshotState?.unavailable,
    snapshotQuery.isError,
    snapshotQuery.isFetching,
    snapshotQuery.isLoading,
    stream.connectionState,
  ]);

  return {
    routeState,
    selectedDeviceId,
    selectedDevice,
    devices,
    setupContext,
    sourceEvidenceScope: runtimeContextQuery.data?.workspace_revision && observableDeviceId ? {
      workspace_id: runtimeContextQuery.data.workspace_id,
      workspace_revision: runtimeContextQuery.data.workspace_revision,
      device_id: observableDeviceId,
    } : null,
    snapshot: currentSnapshot,
    snapshotReceivedAt,
    snapshotStale,
    snapshotError,
    liveValues: stream.liveValues,
    streamState: stream.connectionState,
    logs: stream.logs,
    streamRecovery: stream.streamRecovery,
    onSelectDevice: (deviceId: string) => {
      if (!safeRuntimeIdentity(deviceId) || !devices.some((device) => device.device_id === deviceId)) return;
      const nextParams = new URLSearchParams(searchParams);
      nextParams.set('device_id', deviceId);
      setSearchParams(nextParams);
    },
    onRetrySnapshot: () => observableDeviceId ? snapshotQuery.refetch() : Promise.resolve(undefined),
    onReconnectStream: () => { if (observableDeviceId) stream.reconnect(); },
  };
}
