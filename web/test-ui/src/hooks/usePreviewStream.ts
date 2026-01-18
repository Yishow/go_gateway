import { useState, useEffect, useCallback, useRef } from "react";

/**
 * SSE 預覽事件類型
 */
export interface PreviewEvent {
  type: "connected" | "preview" | "heartbeat" | "error";
  mapping_id?: string;
  raw_value?: number | string | boolean;
  final_value?: number | string | boolean;
  steps?: StepResult[];
  quality?: number;
  timestamp: string;
  error?: string;
}

/**
 * 步驟結果
 */
export interface StepResult {
  step_index: number;
  step_type: string;
  input: number | string | boolean;
  output: number | string | boolean;
  error?: string;
}

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
    autoConnect = false,
    reconnectDelay = 3000,
    maxReconnectAttempts = 5,
    baseUrl = "/api/v1",
  } = options;

  const [latestEvent, setLatestEvent] = useState<PreviewEvent | null>(null);
  const [connectionState, setConnectionState] =
    useState<SSEConnectionState>("disconnected");
  const [error, setError] = useState<string | null>(null);

  const eventSourceRef = useRef<EventSource | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

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
    // 如果已經連接，先斷開
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
    }

    setConnectionState("connecting");
    setError(null);

    const url = `${baseUrl}/datalink/preview/stream?mapping_id=${encodeURIComponent(mappingId)}`;
    const eventSource = new EventSource(url);
    eventSourceRef.current = eventSource;

    eventSource.onopen = () => {
      setConnectionState("connected");
      reconnectAttemptsRef.current = 0;
    };

    eventSource.onmessage = (event) => {
      try {
        const data: PreviewEvent = JSON.parse(event.data);
        setLatestEvent(data);

        // 如果收到 connected 事件，更新狀態
        if (data.type === "connected") {
          setConnectionState("connected");
        }

        // 如果收到 error 事件，記錄錯誤
        if (data.type === "error" && data.error) {
          setError(data.error);
        }
      } catch (err) {
        console.error("Failed to parse SSE event:", err);
      }
    };

    eventSource.onerror = () => {
      setConnectionState("error");
      eventSource.close();
      eventSourceRef.current = null;

      // 嘗試重連
      if (reconnectAttemptsRef.current < maxReconnectAttempts) {
        reconnectAttemptsRef.current += 1;
        setError(
          `連線中斷，${reconnectDelay / 1000} 秒後重連 (${reconnectAttemptsRef.current}/${maxReconnectAttempts})`,
        );

        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, reconnectDelay);
      } else {
        setError("連線失敗，已達最大重連次數");
        setConnectionState("disconnected");
      }
    };
  }, [mappingId, baseUrl, reconnectDelay, maxReconnectAttempts]);

  // 自動連接
  useEffect(() => {
    if (autoConnect && mappingId) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [autoConnect, mappingId, connect, disconnect]);

  return {
    latestEvent,
    connectionState,
    error,
    connect,
    disconnect,
    isConnected: connectionState === "connected",
  };
}

export default usePreviewStream;
