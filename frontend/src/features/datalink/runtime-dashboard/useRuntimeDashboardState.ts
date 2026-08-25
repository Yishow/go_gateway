import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';
import { useStudioV2RuntimeContextQuery } from '../../../hooks/datalink/useStudioV2RuntimeContext';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import type { RuntimeValueEvent } from '../../../types/datalink';
import type {
  StudioV2RuntimeContextDevice,
  StudioV2RuntimeSetupContext,
} from '../../../types/studioV2RuntimeContext';
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
  snapshot: RuntimeStatusWithDiagnostics | null;
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
  const queryDeviceId = searchParams.get('device_id');
  const [lastSnapshot, setLastSnapshot] = useState<RuntimeStatusWithDiagnostics | null>(null);
  const [isDegraded, setIsDegraded] = useState(false);

  const runtimeContextQuery = useStudioV2RuntimeContextQuery();
  const devices = useMemo(
    () => runtimeContextQuery.data?.devices ?? [],
    [runtimeContextQuery.data],
  );
  const setupContext = runtimeContextQuery.data?.setup ?? null;
  const selectedDeviceId = useMemo(() => {
    if (queryDeviceId) {
      return queryDeviceId;
    }
    return runtimeContextQuery.data?.default_device_id ?? null;
  }, [queryDeviceId, runtimeContextQuery.data]);
  const selectedDevice = useMemo(
    () => devices.find((device) => device.device_id === selectedDeviceId) ?? null,
    [devices, selectedDeviceId],
  );
  const selectedDeviceContextMissing = Boolean(
    runtimeContextQuery.data && selectedDeviceId && !selectedDevice,
  );
  const observableDeviceId = selectedDeviceContextMissing ? null : selectedDeviceId;

  const pointsQuery = usePointsQuery(
    observableDeviceId ? { device_id: observableDeviceId } : undefined,
  );
  const pointIds = useMemo(
    () => (pointsQuery.data ?? []).map((point) => point.id),
    [pointsQuery.data],
  );

  const snapshotQuery = useRuntimeStatus({
    deviceId: observableDeviceId,
    pollingIntervalMs: isDegraded ? degradedPollingIntervalMs : false,
  });
  const snapshotError = useMemo(() => {
    const error = runtimeContextQuery.error ?? snapshotQuery.error;
    return error ? getSafeErrorMessage(error, translateError) : null;
  }, [runtimeContextQuery.error, snapshotQuery.error, translateError]);

  const stream = useRuntimeDashboardStream({
    deviceId: observableDeviceId,
    pointIds,
  });

  useEffect(() => {
    setLastSnapshot(null);
    setIsDegraded(false);
  }, [selectedDeviceId]);

  useEffect(() => {
    if (snapshotQuery.data) {
      setLastSnapshot(snapshotQuery.data as RuntimeStatusWithDiagnostics);
    }
  }, [snapshotQuery.data]);

  useEffect(() => {
    if (stream.connectionState === 'connected') {
      setIsDegraded(false);
      return;
    }

    if (
      (stream.connectionState === 'error' ||
        stream.connectionState === 'degraded' ||
        stream.connectionState === 'unavailable' ||
        stream.connectionState === 'reconnecting') &&
      lastSnapshot
    ) {
      setIsDegraded(true);
    }
  }, [lastSnapshot, stream.connectionState]);

  const currentSnapshot =
    lastSnapshot ?? (snapshotQuery.data as RuntimeStatusWithDiagnostics | undefined) ?? null;
  const snapshotState = currentSnapshot?.snapshot_state;

  const routeState = useMemo<RuntimeDashboardRouteState>(() => {
    if (!runtimeContextQuery.data && (runtimeContextQuery.isLoading || runtimeContextQuery.isFetching)) {
      return 'loading';
    }

    if (!runtimeContextQuery.data && runtimeContextQuery.isError) {
      return 'error';
    }

    if (!selectedDeviceId) {
      return 'empty-workspace';
    }

    if (selectedDeviceContextMissing) {
      return 'missing-device-context';
    }

    if (!lastSnapshot && snapshotError?.code === 'runtime_device_not_found') {
      return 'missing-device-context';
    }

    if (!lastSnapshot && snapshotQuery.isError) {
      return 'error';
    }

    if (snapshotState?.empty) {
      return 'empty';
    }

    if (!lastSnapshot && (snapshotQuery.isLoading || snapshotQuery.isFetching)) {
      return 'loading';
    }

    if (
      snapshotState?.degraded ||
      snapshotState?.unavailable ||
      snapshotState?.stale ||
      isDegraded ||
      stream.connectionState === 'unavailable' ||
      stream.connectionState === 'degraded' ||
      stream.connectionState === 'error'
    ) {
      return 'degraded';
    }

    if (stream.connectionState === 'connected') {
      return 'live';
    }

    return 'loading';
  }, [
    isDegraded,
    lastSnapshot,
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
    snapshot: currentSnapshot,
    snapshotError,
    liveValues: stream.liveValues,
    streamState: stream.connectionState,
    logs: stream.logs,
    streamRecovery: stream.streamRecovery,
    onSelectDevice: (deviceId: string) => {
      const nextParams = new URLSearchParams(searchParams);
      nextParams.set('device_id', deviceId);
      setSearchParams(nextParams);
    },
    onRetrySnapshot: () => snapshotQuery.refetch(),
    onReconnectStream: stream.reconnect,
  };
}
