import { useEffect, useMemo, useState } from 'react';
import { runtimeAPI } from '../../../services/datalink';
import type {
  RuntimeDeviceStatusEvent,
  RuntimeStreamConnectionState,
  RuntimeValueEvent,
} from '../../../types/datalink';

interface UseRuntimeStreamOptions {
  deviceId: string | null;
  pointIds?: string[];
}

interface UseRuntimeStreamResult {
  connectionState: RuntimeStreamConnectionState;
  liveValues: Record<string, RuntimeValueEvent>;
  latestStatus: RuntimeDeviceStatusEvent | null;
}

export function useRuntimeDashboardStream({
  deviceId,
  pointIds = [],
}: UseRuntimeStreamOptions): UseRuntimeStreamResult {
  const [connectionState, setConnectionState] =
    useState<RuntimeStreamConnectionState>('disconnected');
  const [liveValues, setLiveValues] = useState<Record<string, RuntimeValueEvent>>(
    {},
  );
  const [latestStatus, setLatestStatus] =
    useState<RuntimeDeviceStatusEvent | null>(null);

  const pointIdsKey = useMemo(() => pointIds.join(','), [pointIds]);
  const normalizedPointIds = useMemo(
    () => (pointIdsKey ? pointIdsKey.split(',') : []),
    [pointIdsKey],
  );

  useEffect(() => {
    setLiveValues({});
    setLatestStatus(null);

    if (!deviceId) {
      setConnectionState('disconnected');
      return undefined;
    }

    if (typeof EventSource === 'undefined') {
      setConnectionState('error');
      return undefined;
    }

    setConnectionState('connecting');
    const eventSource = new EventSource(
      runtimeAPI.getStreamUrl(deviceId, normalizedPointIds),
    );

    const handleValueEvent = (event: Event) => {
      const payload = JSON.parse((event as MessageEvent<string>).data) as RuntimeValueEvent;
      setLiveValues((current) => ({
        ...current,
        [payload.point_id]: payload,
      }));
    };

    const handleStatusEvent = (event: Event) => {
      const payload = JSON.parse((event as MessageEvent<string>).data) as RuntimeDeviceStatusEvent;
      setLatestStatus(payload);
    };

    eventSource.onopen = () => {
      setConnectionState('connected');
    };
    eventSource.onerror = () => {
      setConnectionState('error');
    };
    eventSource.addEventListener('value', handleValueEvent);
    eventSource.addEventListener('status', handleStatusEvent);

    return () => {
      eventSource.removeEventListener('value', handleValueEvent);
      eventSource.removeEventListener('status', handleStatusEvent);
      eventSource.close();
      setConnectionState('disconnected');
    };
  }, [deviceId, normalizedPointIds, pointIdsKey]);

  return {
    connectionState,
    liveValues,
    latestStatus,
  };
}
