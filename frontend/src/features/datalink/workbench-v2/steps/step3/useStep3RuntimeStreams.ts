import * as React from 'react';
import { runtimeAPI } from '../../../../../services/datalink';
import type {
  RuntimeStreamConnectionState,
  RuntimeStreamRecovery,
  RuntimeValueEvent,
} from '../../../../../types/datalink';
import {
  parseMessageEventRecord,
  parseRuntimeStreamStateRecord,
  parseRuntimeValueRecord,
} from '../../../../../utils/safeJson';
import type { DevicePointGroup } from './step3LiveSubscription';

export const STEP3_RETRY_DELAYS_MS = [1000, 2000, 4000, 8000, 16000, 30000] as const;

interface RuntimeStreamsState {
  connectionByDevice: Record<string, RuntimeStreamConnectionState>;
  liveValues: Record<string, RuntimeValueEvent>;
  lastSuccessAtByDevice: Record<string, string>;
  recoveryByDevice: Record<string, RuntimeStreamRecovery>;
}

export function useStep3RuntimeStreams(
  groups: DevicePointGroup[],
  persistedToLocalPointId: Record<string, string>,
  pointAddressToLocalPointId: Record<string, string>,
): RuntimeStreamsState {
  const [connectionByDevice, setConnectionByDevice] = React.useState<Record<string, RuntimeStreamConnectionState>>({});
  const [liveValues, setLiveValues] = React.useState<Record<string, RuntimeValueEvent>>({});
  const [lastSuccessAtByDevice, setLastSuccessAtByDevice] = React.useState<Record<string, string>>({});
  const [recoveryByDevice, setRecoveryByDevice] = React.useState<Record<string, RuntimeStreamRecovery>>({});
  const lastSuccessAtRef = React.useRef<Record<string, string>>({});

  React.useEffect(() => {
    setLiveValues({});
    setLastSuccessAtByDevice({});
    lastSuccessAtRef.current = {};
    setRecoveryByDevice({});

    if (groups.length === 0) {
      setConnectionByDevice({});
      return undefined;
    }

    if (typeof EventSource === 'undefined') {
      setConnectionByDevice(groups.reduce<Record<string, RuntimeStreamConnectionState>>((acc, group) => {
        acc[group.deviceId] = 'degraded';
        return acc;
      }, {}));
      return undefined;
    }

    setConnectionByDevice(groups.reduce<Record<string, RuntimeStreamConnectionState>>((acc, group) => {
      acc[group.deviceId] = 'connecting';
      return acc;
    }, {}));

    let disposed = false;
    type StreamSource = {
      group: DevicePointGroup;
      eventSource: EventSource | null;
      retryIndex: number;
      retryTimer: number | null;
      handleValueEvent?: (event: Event) => void;
      handleStreamStateEvent?: (event: Event) => void;
    };
    const sources: StreamSource[] = groups.map((group) => ({
      group,
      eventSource: null,
      retryIndex: 0,
      retryTimer: null,
    }));

    const connect = (source: StreamSource): void => {
      if (disposed) return;
      setConnectionByDevice((current) => ({
        ...current,
        [source.group.deviceId]: source.retryIndex === 0 ? 'connecting' : 'reconnecting',
      }));
      const eventSource = new EventSource(runtimeAPI.getStreamUrl(source.group.deviceId, source.group.streamPointIds));
      source.eventSource = eventSource;

      const scheduleRetry = (): void => {
        if (disposed || source.retryTimer !== null) return;
        if (source.retryIndex >= STEP3_RETRY_DELAYS_MS.length) {
          setConnectionByDevice((current) => ({ ...current, [source.group.deviceId]: 'degraded' }));
          return;
        }
        const delay = STEP3_RETRY_DELAYS_MS[source.retryIndex];
        source.retryIndex += 1;
        setConnectionByDevice((current) => ({
          ...current,
          [source.group.deviceId]: lastSuccessAtRef.current[source.group.deviceId] ? 'stale' : 'reconnecting',
        }));
        source.retryTimer = window.setTimeout(() => {
          source.retryTimer = null;
          connect(source);
        }, delay);
      };

      const closeAndRetry = (): void => {
        eventSource.removeEventListener('value', handleValueEvent);
        eventSource.removeEventListener('stream_state', handleStreamStateEvent);
        eventSource.close();
        source.eventSource = null;
        scheduleRetry();
      };

      const handleValueEvent = (event: Event): void => {
        const payload = parseRuntimeValueRecord(parseMessageEventRecord(event));
        if (!payload) {
          closeAndRetry();
          return;
        }
        const localPointId = persistedToLocalPointId[payload.point_id] ??
          pointAddressToLocalPointId[`${payload.device_id}::${payload.address}`];
        if (!localPointId) return;
        source.retryIndex = 0;
        lastSuccessAtRef.current[source.group.deviceId] = payload.timestamp;
        setLastSuccessAtByDevice((current) => ({ ...current, [source.group.deviceId]: payload.timestamp }));
        setConnectionByDevice((current) => ({ ...current, [source.group.deviceId]: 'live' }));
        setRecoveryByDevice((current) => {
          if (!current[source.group.deviceId]) return current;
          const next = { ...current };
          delete next[source.group.deviceId];
          return next;
        });
        setLiveValues((current) => ({ ...current, [localPointId]: payload }));
      };
      source.handleValueEvent = handleValueEvent;
      eventSource.addEventListener('value', handleValueEvent);

      const handleStreamStateEvent = (event: Event): void => {
        const payload = parseRuntimeStreamStateRecord(parseMessageEventRecord(event));
        if (!payload ||
          (payload.device_id && payload.device_id !== source.group.deviceId)) {
          if (!payload) closeAndRetry();
          return;
        }
        const streamState = payload.stream_state;
        if (streamState.state === 'ready') {
          setConnectionByDevice((current) => ({ ...current, [source.group.deviceId]: 'connected' }));
          setRecoveryByDevice((current) => {
            if (!current[source.group.deviceId]) return current;
            const next = { ...current };
            delete next[source.group.deviceId];
            return next;
          });
          return;
        }
        const connectionState: RuntimeStreamConnectionState = streamState.state === 'stale' ? 'stale' : 'degraded';
        setConnectionByDevice((current) => ({ ...current, [source.group.deviceId]: connectionState }));
        setRecoveryByDevice((current) => ({
          ...current,
          [source.group.deviceId]: {
            code: payload.code,
            action: payload.action,
            requestId: payload.request_id,
            retryable: payload.retryable === true,
          },
        }));
      };
      source.handleStreamStateEvent = handleStreamStateEvent;
      eventSource.addEventListener('stream_state', handleStreamStateEvent);
      eventSource.onopen = () => {
        if (!disposed) setConnectionByDevice((current) => ({ ...current, [source.group.deviceId]: 'connected' }));
      };
      eventSource.onerror = closeAndRetry;
    };

    sources.forEach(connect);
    return () => {
      disposed = true;
      sources.forEach((source) => {
        if (source.retryTimer !== null) window.clearTimeout(source.retryTimer);
        if (source.eventSource) {
          if (source.handleValueEvent) source.eventSource.removeEventListener('value', source.handleValueEvent);
          if (source.handleStreamStateEvent) source.eventSource.removeEventListener('stream_state', source.handleStreamStateEvent);
          source.eventSource.close();
        }
      });
    };
  }, [groups, persistedToLocalPointId, pointAddressToLocalPointId]);

  return { connectionByDevice, liveValues, lastSuccessAtByDevice, recoveryByDevice };
}
