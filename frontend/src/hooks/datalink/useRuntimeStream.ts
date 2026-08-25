import { useEffect, useMemo, useState } from 'react';
import { runtimeAPI } from '../../services/datalink';
import type {
  RuntimeStreamConnectionState,
  RuntimeValueEvent,
} from '../../types/datalink';
import { parseMessageEventRecord, parseRuntimeValueRecord } from '../../utils/safeJson';

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

  const pointIdsKey = useMemo(() => [...pointIds].sort().join(','), [pointIds]);
  const normalizedPointIds = useMemo(
    () => (pointIdsKey ? pointIdsKey.split(',') : []),
    [pointIdsKey],
  );

  useEffect(() => {
    let active = true;
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
      if (!active) return;
      const payload = parseMessageEventRecord(event);
      const value = parseRuntimeValueRecord(payload);
      if (!value) {
        setConnectionState('degraded');
        return;
      }
      setConnectionState('connected');
      setLiveValues((current) => ({
        ...current,
        [value.point_id]: value,
      }));
    };
    const handleValueEvent = (event: Event) => {
      handleValue(event as MessageEvent<string>);
    };

    eventSource.onopen = () => {
      if (!active) return;
      setConnectionState('connected');
    };
    eventSource.onerror = () => {
      if (!active) return;
      setConnectionState('error');
    };
    eventSource.addEventListener('value', handleValueEvent);

    return () => {
      active = false;
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
