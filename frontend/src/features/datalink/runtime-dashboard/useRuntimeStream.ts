import { useCallback, useEffect, useMemo, useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { runtimeAPI } from '../../../services/datalink';
import type {
  RuntimeDeviceStatusEvent,
  RuntimeStreamConnectionState,
  RuntimeStreamRecovery,
  RuntimeValueEvent,
} from '../../../types/datalink';
import {
  parseMessageEventRecord,
  parseRuntimeStatusRecord,
  parseRuntimeStreamStateRecord,
  parseRuntimeValueRecord,
} from '../../../utils/safeJson';

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

interface UseRuntimeStreamResult {
  connectionState: RuntimeDashboardStreamConnectionState;
  liveValues: Record<string, RuntimeValueEvent>;
  latestStatus: RuntimeDeviceStatusEvent | null;
  logs: RuntimeDashboardLog[];
  streamRecovery: RuntimeStreamRecovery | null;
  reconnect: () => void;
}

export function useRuntimeDashboardStream({
  deviceId,
  pointIds = [],
}: UseRuntimeStreamOptions): UseRuntimeStreamResult {
  const { t } = useTranslation('runtime-dashboard');
  const translateRef = useRef(t);
  translateRef.current = t;
  const [connectionState, setConnectionState] =
    useState<RuntimeDashboardStreamConnectionState>('disconnected');
  const [liveValues, setLiveValues] = useState<Record<string, RuntimeValueEvent>>(
    {},
  );
  const [latestStatus, setLatestStatus] =
    useState<RuntimeDeviceStatusEvent | null>(null);
  const [logs, setLogs] = useState<RuntimeDashboardLog[]>([]);
  const [streamRecovery, setStreamRecovery] = useState<RuntimeStreamRecovery | null>(null);
  const [reconnectVersion, setReconnectVersion] = useState(0);
  const reconnect = useCallback(() => {
    setReconnectVersion((current) => current + 1);
  }, []);

  // 1. 緩衝與先前狀態記錄 Refs，以實現更新節流與狀態去重
  const bufferRef = useRef<Record<string, RuntimeValueEvent>>({});
  const streamGenerationRef = useRef(0);
  const lastErrorRef = useRef<string | null>(null);
  const pointStaleStatesRef = useRef<Record<string, boolean>>({});

  // 2. 進行冪等排序以防 EventSource 因點位列表順序不同而無效重連
  const pointIdsKey = useMemo(() => [...pointIds].sort().join(','), [pointIds]);
  const normalizedPointIds = useMemo(
    () => (pointIdsKey ? pointIdsKey.split(',') : []),
    [pointIdsKey],
  );

  useEffect(() => {
    let active = true;
    const generation = streamGenerationRef.current + 1;
    streamGenerationRef.current = generation;
    bufferRef.current = {};
    setLiveValues({});
    setLatestStatus(null);
    setStreamRecovery(null);
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
          message: translateRef.current('errors.browser_event_source_unsupported'),
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
        message: translateRef.current('logs.stream_connecting'),
      },
    ]);

    const eventSource = new EventSource(
      runtimeAPI.getStreamUrl(deviceId, normalizedPointIds),
    );

    // 高頻事件不直接更新 state，而是緩衝在 Ref 中
    const handleValueEvent = (event: Event) => {
      if (!active) return;
      const payload = parseMessageEventRecord(event);
      const value = parseRuntimeValueRecord(payload);
      if (!active || streamGenerationRef.current !== generation) return;
      if (!value) {
        setConnectionState('degraded');
        return;
      }
      setConnectionState('connected');
      bufferRef.current[value.point_id] = value;

      // 避免重複或高頻的 stale 警報，只在點位 stale 狀態改變時寫入日誌，並記錄恢復狀態
      const wasStale = pointStaleStatesRef.current[value.point_id] || false;
      const isStale = value.stale;

      if (isStale !== wasStale) {
        pointStaleStatesRef.current[value.point_id] = isStale;
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: new Date().toISOString(),
            level: isStale ? 'warn' : 'info',
          message: isStale
            ? translateRef.current('logs.point_stale')
            : translateRef.current('logs.point_recovered'),
          },
        ]);
      }
    };

    const handleStatusEvent = (event: Event) => {
      if (!active) return;
      const rawPayload = parseMessageEventRecord(event);
      const payload = parseRuntimeStatusRecord(rawPayload);
      if (!active || streamGenerationRef.current !== generation) return;
      if (!payload) {
        setConnectionState('degraded');
        return;
      }
      setConnectionState('connected');
      const { last_error: _lastError, ...safeStatus } = payload;
      setLatestStatus(safeStatus);

      if (payload.last_error && payload.last_error !== lastErrorRef.current) {
        lastErrorRef.current = payload.last_error;
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: new Date().toISOString(),
            level: 'error',
            message: translateRef.current('logs.device_error', 'Device reported an error. Review the runtime status and retry.'),
          },
        ]);
      } else if (!payload.last_error && lastErrorRef.current) {
        lastErrorRef.current = null;
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: new Date().toISOString(),
            level: 'info',
            message: translateRef.current('logs.device_error_cleared', 'Device error cleared.'),
          },
        ]);
      }
    };

    const handleStreamStateEvent = (event: Event) => {
      if (!active) return;
      const rawPayload = parseMessageEventRecord(event);
      const payload = parseRuntimeStreamStateRecord(rawPayload);
      if (!active || streamGenerationRef.current !== generation) return;
      if (!payload) {
        setConnectionState('degraded');
        return;
      }
      if (payload.device_id && payload.device_id !== deviceId) {
        return;
      }

      const streamTruth = payload.stream_state;
      if (!streamTruth) {
        return;
      }

      if (streamTruth.state === 'ready') {
        setStreamRecovery(null);
        setConnectionState('connected');
        return;
      }

      setStreamRecovery({
        code: payload.code,
        action: payload.action,
        requestId: payload.request_id,
        retryable: payload.retryable === true,
      });

      if (streamTruth.unavailable || streamTruth.state === 'unavailable') {
        setConnectionState('unavailable');
        setLogs((current) => [
          ...current.slice(-99),
          {
            timestamp: payload.timestamp ?? new Date().toISOString(),
            level: 'warn',
            message: translateRef.current(
              'errors.runtime_stream_unavailable',
              'Live runtime stream is unavailable. Retry the runtime view.',
            ),
          },
        ]);
      }
    };

    eventSource.onopen = () => {
      if (!active || streamGenerationRef.current !== generation) return;
      setConnectionState('connecting');
      setLogs((current) => [
        ...current.slice(-99),
        {
          timestamp: new Date().toISOString(),
          level: 'info',
          message: translateRef.current('logs.stream_transport_open', 'Runtime stream transport opened; waiting for server readiness.'),
        },
      ]);
    };

    eventSource.onerror = () => {
      if (!active || streamGenerationRef.current !== generation) return;
      setConnectionState('error');
      setLogs((current) => [
        ...current.slice(-99),
        {
          timestamp: new Date().toISOString(),
          level: 'error',
          message: translateRef.current('logs.stream_connection_lost'),
        },
      ]);
    };

    eventSource.addEventListener('value', handleValueEvent);
    eventSource.addEventListener('status', handleStatusEvent);
    eventSource.addEventListener('stream_state', handleStreamStateEvent);

    // 每 200ms 批量合併並寫入 React 狀態，防止 DOM 重繪頻率過高造成卡頓
    const interval = setInterval(() => {
      if (!active || streamGenerationRef.current !== generation) return;
      if (Object.keys(bufferRef.current).length === 0) return;
      setLiveValues((current) => {
        const next = { ...current, ...bufferRef.current };
        bufferRef.current = {};
        return next;
      });
    }, 200);

    return () => {
      active = false;
      if (streamGenerationRef.current === generation) streamGenerationRef.current += 1;
      bufferRef.current = {};
      clearInterval(interval);
      eventSource.removeEventListener('value', handleValueEvent);
      eventSource.removeEventListener('status', handleStatusEvent);
      eventSource.removeEventListener('stream_state', handleStreamStateEvent);
      eventSource.close();
      setConnectionState('disconnected');
    };
  }, [deviceId, normalizedPointIds, pointIdsKey, reconnectVersion]);

  return {
    connectionState,
    liveValues,
    latestStatus,
    logs,
    streamRecovery,
    reconnect,
  };
}
