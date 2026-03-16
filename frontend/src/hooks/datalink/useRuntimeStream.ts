import { useEffect, useMemo, useState } from 'react';
import { runtimeAPI } from '../../services/datalink';
import type {
  RuntimeStreamConnectionState,
  RuntimeValueEvent,
} from '../../types/datalink';

interface UseRuntimeStreamOptions {
  deviceId: string | null;
  pointIds?: string[];
}

interface UseRuntimeStreamResult {
  connectionState: RuntimeStreamConnectionState;
  liveValues: Record<string, RuntimeValueEvent>;
}

export function useRuntimeStream({
  deviceId,
  pointIds = [],
}: UseRuntimeStreamOptions): UseRuntimeStreamResult {
  const [connectionState, setConnectionState] =
    useState<RuntimeStreamConnectionState>('disconnected');
  const [liveValues, setLiveValues] = useState<Record<string, RuntimeValueEvent>>(
    {},
  );

  const pointIdsKey = useMemo(() => pointIds.join(','), [pointIds]);
  const normalizedPointIds = useMemo(
    () => (pointIdsKey ? pointIdsKey.split(',') : []),
    [pointIdsKey],
  );

  useEffect(() => {
    setLiveValues({});

    if (!deviceId) {
      setConnectionState('disconnected');
      return undefined;
    }

    if (typeof EventSource === 'undefined') {
      setConnectionState('disconnected');
      return undefined;
    }

    setConnectionState('connecting');
    const eventSource = new EventSource(
      runtimeAPI.getStreamUrl(deviceId, normalizedPointIds),
    );

    const handleValue = (event: MessageEvent<string>) => {
      const payload = JSON.parse(event.data) as RuntimeValueEvent;
      setLiveValues((current) => ({
        ...current,
        [payload.point_id]: payload,
      }));
    };
    const handleValueEvent = (event: Event) => {
      handleValue(event as MessageEvent<string>);
    };

    eventSource.onopen = () => {
      setConnectionState('connected');
    };
    eventSource.onerror = () => {
      setConnectionState('error');
    };
    eventSource.addEventListener('value', handleValueEvent);

    return () => {
      eventSource.removeEventListener('value', handleValueEvent);
      eventSource.close();
      setConnectionState('disconnected');
    };
  }, [deviceId, normalizedPointIds, pointIdsKey]);

  return {
    connectionState,
    liveValues,
  };
}
