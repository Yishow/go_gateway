import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useStudioV2RuntimeContextQuery } from '../../../hooks/datalink/useStudioV2RuntimeContext';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import type {
  StudioV2RuntimeContextDevice,
  RuntimeValueEvent,
} from '../../../types/datalink';
import type { RuntimeStatusWithDiagnostics } from '../../../types/runtimeDiagnostics';
import { useRuntimeStatus } from './useRuntimeStatus';
import { useRuntimeDashboardStream } from './useRuntimeStream';
import type {
  RuntimeDashboardLog,
  RuntimeDashboardStreamConnectionState,
} from './useRuntimeStream';

const degradedPollingIntervalMs = 5000;

function isMissingRuntimeDeviceError(message: string | null): boolean {
  if (!message) {
    return false;
  }

  return message.includes('device not found');
}

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
  snapshot: RuntimeStatusWithDiagnostics | null;
  snapshotError: string | null;
  liveValues: Record<string, RuntimeValueEvent>;
  streamState: RuntimeDashboardStreamConnectionState;
  logs?: RuntimeDashboardLog[];
  onSelectDevice: (deviceId: string) => void;
  onRetrySnapshot: () => Promise<unknown>;
}

export function useRuntimeDashboardState(): RuntimeDashboardState {
  const [searchParams, setSearchParams] = useSearchParams();
  const queryDeviceId = searchParams.get('device_id');
  const [lastSnapshot, setLastSnapshot] = useState<RuntimeStatusWithDiagnostics | null>(null);
  const [isDegraded, setIsDegraded] = useState(false);

  const runtimeContextQuery = useStudioV2RuntimeContextQuery();
  const devices = useMemo(
    () => runtimeContextQuery.data?.devices ?? [],
    [runtimeContextQuery.data],
  );
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
  const snapshotQueryError =
    snapshotQuery.error instanceof Error ? snapshotQuery.error.message : null;
  const runtimeContextError =
    runtimeContextQuery.error instanceof Error ? runtimeContextQuery.error.message : null;
  const snapshotError = runtimeContextError ?? snapshotQueryError;

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

    if (stream.connectionState === 'error' && lastSnapshot) {
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

    if (!lastSnapshot && isMissingRuntimeDeviceError(snapshotError)) {
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
    snapshotError,
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
    snapshot: currentSnapshot,
    snapshotError,
    liveValues: stream.liveValues,
    streamState: stream.connectionState,
    logs: stream.logs,
    onSelectDevice: (deviceId: string) => {
      const nextParams = new URLSearchParams(searchParams);
      nextParams.set('device_id', deviceId);
      setSearchParams(nextParams);
    },
    onRetrySnapshot: () => snapshotQuery.refetch(),
  };
}
