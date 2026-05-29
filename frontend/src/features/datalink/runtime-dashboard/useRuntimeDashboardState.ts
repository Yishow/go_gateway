import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import type {
  Device,
  RuntimeStatus,
  RuntimeStreamConnectionState,
  RuntimeValueEvent,
} from '../../../types/datalink';
import { useRuntimeStatus } from './useRuntimeStatus';
import { useRuntimeDashboardStream } from './useRuntimeStream';
import type { RuntimeDashboardLog } from './useRuntimeStream';

const degradedPollingIntervalMs = 5000;

export type RuntimeDashboardRouteState =
  | 'missing-device-context'
  | 'loading'
  | 'live'
  | 'degraded'
  | 'error';

export interface RuntimeDashboardState {
  routeState: RuntimeDashboardRouteState;
  selectedDeviceId: string | null;
  selectedDevice: Device | null;
  devices: Device[];
  snapshot: RuntimeStatus | null;
  snapshotError: string | null;
  liveValues: Record<string, RuntimeValueEvent>;
  streamState: RuntimeStreamConnectionState;
  logs?: RuntimeDashboardLog[];
  onSelectDevice: (deviceId: string) => void;
  onRetrySnapshot: () => Promise<unknown>;
}

const emptySnapshot: RuntimeStatus = {
  running: false,
  uptime_seconds: 0,
  metrics: {
    collected_total: 0,
    write_success_total: 0,
    write_error_total: 0,
    mapping_error_total: 0,
    point_state_error_total: 0,
  },
  collectors: [],
};

export function useRuntimeDashboardState(): RuntimeDashboardState {
  const [searchParams, setSearchParams] = useSearchParams();
  const selectedDeviceId = searchParams.get('device_id');
  const [lastSnapshot, setLastSnapshot] = useState<RuntimeStatus | null>(null);
  const [isDegraded, setIsDegraded] = useState(false);

  const devicesQuery = useDevicesQuery();
  const devices = useMemo(() => devicesQuery.data ?? [], [devicesQuery.data]);
  const selectedDevice = useMemo(
    () => devices.find((device) => device.id === selectedDeviceId) ?? null,
    [devices, selectedDeviceId],
  );

  const pointsQuery = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const pointIds = useMemo(
    () => (pointsQuery.data ?? []).map((point) => point.id),
    [pointsQuery.data],
  );

  const snapshotQuery = useRuntimeStatus({
    deviceId: selectedDeviceId,
    pollingIntervalMs: isDegraded ? degradedPollingIntervalMs : false,
  });

  const stream = useRuntimeDashboardStream({
    deviceId: selectedDeviceId,
    pointIds,
  });

  useEffect(() => {
    setLastSnapshot(null);
    setIsDegraded(false);
  }, [selectedDeviceId]);

  useEffect(() => {
    if (snapshotQuery.data) {
      setLastSnapshot(snapshotQuery.data);
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

  const routeState = useMemo<RuntimeDashboardRouteState>(() => {
    if (!selectedDeviceId) {
      return 'missing-device-context';
    }

    if (!lastSnapshot && snapshotQuery.isError) {
      return 'error';
    }

    if (!lastSnapshot && (snapshotQuery.isLoading || snapshotQuery.isFetching)) {
      return 'loading';
    }

    if (isDegraded || stream.connectionState === 'error') {
      return 'degraded';
    }

    if (stream.connectionState === 'connected') {
      return 'live';
    }

    return 'loading';
  }, [
    isDegraded,
    lastSnapshot,
    selectedDeviceId,
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
    snapshot: lastSnapshot ?? snapshotQuery.data ?? emptySnapshot,
    snapshotError: snapshotQuery.error instanceof Error ? snapshotQuery.error.message : null,
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
