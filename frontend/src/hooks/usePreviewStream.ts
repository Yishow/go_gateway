import { useState, useEffect, useCallback, useRef } from "react";
import { logger } from "../utils/logger";
import { VITE_API_BASE_URL } from "../env";
import { parsePreviewEvent, type PreviewEvent } from './previewStreamEvents';
import { parseBoundedJson } from '../utils/safeJson';

export { parsePreviewEvent } from './previewStreamEvents';
export type { PreviewErrorEnvelope, PreviewEvent, StepResult } from './previewStreamEvents';

/**
 * SSE 連接狀態
 */
export type SSEConnectionState =
  | "connecting"
  | "connected"
  | "disconnected"
  | "error";

/**
 * usePreviewStream Hook 配置
 */
interface UsePreviewStreamOptions {
  /** Mapping ID */
  mappingId: string;
  /** Workspace scope required by the preview stream */
  workspaceId: string;
  /** 是否自動連接 */
  autoConnect?: boolean;
  /** 重連延遲 (ms) */
  reconnectDelay?: number;
  /** 最大重連次數 */
  maxReconnectAttempts?: number;
  /** API 基礎 URL */
  baseUrl?: string;
}

/**
 * usePreviewStream Hook 返回值
 */
interface UsePreviewStreamReturn {
  /** 最新的預覽事件 */
  latestEvent: PreviewEvent | null;
  /** 連接狀態 */
  connectionState: SSEConnectionState;
  /** 錯誤訊息 */
  error: string | null;
  requestId: string | null;
  /** 開始連接 */
  connect: () => void;
  /** 斷開連接 */
  disconnect: () => void;
  /** 是否已連接 */
  isConnected: boolean;
}

/**
 * 即時預覽 SSE Hook
 *
 * 連接到後端 SSE 端點，接收即時預覽資料
 *
 * @example
 * ```tsx
 * const { latestEvent, connectionState, connect, disconnect } = usePreviewStream({
 *   mappingId: 'mapping-123',
 *   workspaceId: 'workspace-123',
 *   autoConnect: true,
 * });
 *
 * useEffect(() => {
 *   if (latestEvent?.type === 'preview') {
 *     console.log('Raw:', latestEvent.raw_value);
 *     console.log('Final:', latestEvent.final_value);
 *   }
 * }, [latestEvent]);
 * ```
 */
export function usePreviewStream(
  options: UsePreviewStreamOptions,
): UsePreviewStreamReturn {
  const {
    mappingId,
    workspaceId,
    autoConnect = false,
    reconnectDelay = 3000,
    maxReconnectAttempts = 5,
    baseUrl = VITE_API_BASE_URL,
  } = options;

  const [latestEvent, setLatestEvent] = useState<PreviewEvent | null>(null);
  const [connectionState, setConnectionState] =
    useState<SSEConnectionState>("disconnected");
  const [error, setError] = useState<string | null>(null);
  const [requestId, setRequestId] = useState<string | null>(null);

  const eventSourceRef = useRef<EventSource | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const connectionGenerationRef = useRef(0);

  /**
   * 清理重連計時器
   */
  const clearReconnectTimeout = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
  }, []);

  /**
   * 斷開連接
   */
  const disconnect = useCallback(() => {
    connectionGenerationRef.current += 1;
    clearReconnectTimeout();

    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }

    setConnectionState("disconnected");
    reconnectAttemptsRef.current = 0;
  }, [clearReconnectTimeout]);

  /**
   * 連接到 SSE 端點
   */
  const connect = useCallback(() => {
    if (!mappingId || !workspaceId) {
      setConnectionState("disconnected");
      return;
    }
    connectionGenerationRef.current += 1;
    const generation = connectionGenerationRef.current;
    clearReconnectTimeout();
    // 如果已經連接，先斷開
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }

    setConnectionState("connecting");
    setError(null);
    setRequestId(null);

    const url = `${baseUrl}/datalink/preview/stream?workspace_id=${encodeURIComponent(workspaceId)}&mapping_id=${encodeURIComponent(mappingId)}`;
    const eventSource = new EventSource(url);
    eventSourceRef.current = eventSource;
    const isCurrent = () => connectionGenerationRef.current === generation;

    eventSource.onopen = () => {
      if (!isCurrent()) return;
      setConnectionState("connected");
      reconnectAttemptsRef.current = 0;
    };

    eventSource.onmessage = (event) => {
      if (!isCurrent()) return;
      try {
        const data = parsePreviewEvent(parseBoundedJson(event.data));
        if (!data) {
          logger.warn("Ignored malformed preview SSE event");
          return;
        }
        setLatestEvent(data);

        // 如果收到 connected 事件，更新狀態
        if (data.type === "connected") {
          setConnectionState("connected");
        }

        // 如果收到 error 事件，記錄錯誤
        if (data.type === "error") {
          setError(data.code ?? "preview_unavailable");
          setRequestId(data.request_id ?? null);
        }

        if (data.type === "close") {
          setError(data.code ?? "preview_unavailable");
          setRequestId(data.request_id ?? null);
          setConnectionState("disconnected");
          eventSource.close();
          if (isCurrent()) eventSourceRef.current = null;
        }
      } catch (err) {
        logger.error("Failed to parse SSE event:", err);
      }
    };

    eventSource.onerror = () => {
      if (!isCurrent()) return;
      setConnectionState("error");
      eventSource.close();
      eventSourceRef.current = null;

      // 嘗試重連
      if (reconnectAttemptsRef.current < maxReconnectAttempts) {
        reconnectAttemptsRef.current += 1;
        setError(
          "preview_unavailable",
        );

        reconnectTimeoutRef.current = setTimeout(() => {
          if (isCurrent()) connect();
        }, reconnectDelay);
      } else {
        setError("preview_unavailable");
        setConnectionState("disconnected");
      }
    };
  }, [baseUrl, clearReconnectTimeout, mappingId, maxReconnectAttempts, reconnectDelay, workspaceId]);

  // 自動連接
  useEffect(() => {
    if (autoConnect && mappingId && workspaceId) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [autoConnect, mappingId, workspaceId, connect, disconnect]);

  return {
    latestEvent,
    connectionState,
    error,
    requestId,
    connect,
    disconnect,
    isConnected: connectionState === "connected",
  };
}

export default usePreviewStream;
