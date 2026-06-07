import { useEffect, useMemo, useState, useRef } from 'react';
import { runtimeAPI } from '../../../services/datalink';
import type {
  RuntimeDeviceStatusEvent,
  RuntimeStreamConnectionState,
  RuntimeValueEvent,
} from '../../../types/datalink';
import type { RuntimeTruthState } from '../../../types/runtimeTruth';

interface UseRuntimeStreamOptions {
  deviceId: string | null;
  pointIds?: string[];
}

export interface RuntimeDashboardLog {
  timestamp: string;
  level: 'info' | 'warn' | 'error';
  message: string;
}

export type RuntimeDashboardStreamConnectionState =
  | RuntimeStreamConnectionState
  | 'unavailable';

interface RuntimeStreamStateEvent {
  device_id?: string;
  stream_state?: RuntimeTruthState;
  timestamp?: string;
}

interface UseRuntimeStreamResult {
  connectionState: RuntimeDashboardStreamConnectionState;
  liveValues: Record<string, RuntimeValueEvent>;
  latestStatus: RuntimeDeviceStatusEvent | null;
  logs: RuntimeDashboardLog[];
}

export function useRuntimeDashboardStream({
  deviceId,
  pointIds = [],
}: UseRuntimeStreamOptions): UseRuntimeStreamResult {
  const [connectionState, setConnectionState] =
    useState<RuntimeDashboardStreamConnectionState>('disconnected');
  const [liveValues, setLiveValues] = useState<Record<string, RuntimeValueEvent>>(
    {},
  );
  const [latestStatus, setLatestStatus] =
    useState<RuntimeDeviceStatusEvent | null>(null);
  const [logs, setLogs] = useState<RuntimeDashboardLog[]>([]);

  // 1. 緩衝與先前狀態記錄 Refs，以實現更新節流與狀態去重
  const bufferRef = useRef<Record<string, RuntimeValueEvent>>({});
  const lastErrorRef = useRef<string | null>(null);
  const pointStaleStatesRef = useRef<Record<string, boolean>>({});

  // 2. 進行冪等排序以防 EventSource 因點位列表順序不同而無效重連
  const pointIdsKey = useMemo(() => [...pointIds].sort().join(','), [pointIds]);
  const normalizedPointIds = useMemo(
    () => (pointIdsKey ? pointIdsKey.split(',') : []),
    [pointIdsKey],
  );

  useEffect(() => {
    setLiveValues({});
    setLatestStatus(null);
    lastErrorRef.current = null;
    pointStaleStatesRef.current = {};

    if (!deviceId) {
      setConnectionState('disconnected');
      return undefined;
    }

    if (typeof EventSource === 'undefined') {
      setConnectionState('error');
      setLogs((current) => [
        ...current.slice(-99),
        {
          timestamp: new Date().toISOString(),
          level: 'error',
          message: 'Browser does not support EventSource stream.',
        },
      ]);
      return undefined;
    }

    setConnectionState('connecting');
    setLogs((current) => [
      ...current.slice(-99),
      {
        timestamp: new Date().toISOString(),
        level: 'info',
        message: `Connecting to SSE stream for device ${deviceId}...`,
      },
    ]);

    const eventSource = new EventSource(
      runtimeAPI.getStreamUrl(deviceId, normalizedPointIds),
    );

    // 高頻事件不直接更新 state，而是緩衝在 Ref 中
    const handleValueEvent = (event: Event) => {
      const payload = JSON.parse((event as MessageEvent<string>).data) as RuntimeValueEvent;
      bufferRef.current[payload.point_id] = payload;

      // 避免重複或高頻的 stale 警報，只在點位 stale 狀態改變時寫入日誌，並記錄恢復狀態
      const wasStale = pointStaleStatesRef.current[payload.point_id] || false;
      const isStale = payload.stale;

      if (isStale !== wasStale) {
        pointStaleStatesRef.current[payload.point_id] = isStale;
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: new Date().toISOString(),
            level: isStale ? 'warn' : 'info',
            message: isStale
              ? `Point ${payload.address} went stale.`
              : `Point ${payload.address} recovered.`,
          },
        ]);
      }
    };

    const handleStatusEvent = (event: Event) => {
      const payload = JSON.parse((event as MessageEvent<string>).data) as RuntimeDeviceStatusEvent;
      setLatestStatus(payload);

      if (payload.last_error && payload.last_error !== lastErrorRef.current) {
        lastErrorRef.current = payload.last_error;
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: new Date().toISOString(),
            level: 'error',
            message: `Device error: ${payload.last_error}`,
          },
        ]);
      } else if (!payload.last_error && lastErrorRef.current) {
        lastErrorRef.current = null;
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: new Date().toISOString(),
            level: 'info',
            message: 'Device error cleared.',
          },
        ]);
      }
    };

    const handleStreamStateEvent = (event: Event) => {
      const payload = JSON.parse((event as MessageEvent<string>).data) as RuntimeStreamStateEvent;
      if (payload.device_id && payload.device_id !== deviceId) {
        return;
      }

      const streamTruth = payload.stream_state;
      if (!streamTruth) {
        return;
      }

      if (streamTruth.state === 'ready') {
        setConnectionState('connected');
        return;
      }

      if (streamTruth.unavailable || streamTruth.state === 'unavailable') {
        setConnectionState('unavailable');
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: payload.timestamp ?? new Date().toISOString(),
            level: 'warn',
            message: streamTruth.reason
              ? `Runtime stream unavailable: ${streamTruth.reason}`
              : 'Runtime stream unavailable.',
          },
        ]);
      }
    };

    eventSource.onopen = () => {
      setConnectionState('connected');
      setLogs((current) => [
        ...current.slice(-99),
        {
          timestamp: new Date().toISOString(),
          level: 'info',
          message: 'SSE stream connected successfully.',
        },
      ]);
    };

    eventSource.onerror = () => {
      setConnectionState('error');
      setLogs((current) => [
        ...current.slice(-99),
        {
          timestamp: new Date().toISOString(),
          level: 'error',
          message: 'SSE stream connection lost. Attempting reconnect...',
        },
      ]);
    };

    eventSource.addEventListener('value', handleValueEvent);
    eventSource.addEventListener('status', handleStatusEvent);
    eventSource.addEventListener('stream_state', handleStreamStateEvent);

    // 每 200ms 批量合併並寫入 React 狀態，防止 DOM 重繪頻率過高造成卡頓
    const interval = setInterval(() => {
      if (Object.keys(bufferRef.current).length === 0) return;
      setLiveValues((current) => {
        const next = { ...current, ...bufferRef.current };
        bufferRef.current = {};
        return next;
      });
    }, 200);

    return () => {
      clearInterval(interval);
      eventSource.removeEventListener('value', handleValueEvent);
      eventSource.removeEventListener('status', handleStatusEvent);
      eventSource.removeEventListener('stream_state', handleStreamStateEvent);
      eventSource.close();
      setConnectionState('disconnected');
    };
  }, [deviceId, normalizedPointIds, pointIdsKey]);

  return {
    connectionState,
    liveValues,
    latestStatus,
    logs,
  };
}
