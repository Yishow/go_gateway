export type RuntimeLogLevel = 'debug' | 'info' | 'warn' | 'error';
export interface RuntimeLogFilters { level: RuntimeLogLevel; source: string; q: string }
export interface RuntimeLogRecord {
  instance_id: string;
  sequence: string;
  timestamp: string;
  level: RuntimeLogLevel;
  source: string;
  code: string;
  message: string;
  fields: Record<string, string | number | boolean>;
  truncated: boolean;
}
export interface RuntimeLogMetadata {
  instance_id: string;
  oldest_cursor: string;
  latest_cursor: string;
  capture_dropped: number;
  evicted: number;
  subscriber_overflow?: number;
  file_dropped?: number;
  disk_errors?: number;
  file_recovered_gaps?: number;
  sink_health: 'healthy' | 'degraded' | 'disabled';
}
export interface RuntimeLogSnapshot extends RuntimeLogMetadata { records: RuntimeLogRecord[] }
export type RuntimeLogEvent =
  | { type: 'log'; record: RuntimeLogRecord; cursor: string }
  | { type: 'handshake' | 'heartbeat'; metadata: RuntimeLogMetadata; cursor: string }
  | { type: 'gap' | 'reset'; reason: string; cursor: string };
export interface RuntimeLogTransport {
  snapshot(filters: RuntimeLogFilters, signal: AbortSignal): Promise<RuntimeLogSnapshot>;
  stream(filters: RuntimeLogFilters, cursor: string, signal: AbortSignal,
    onEvent: (event: RuntimeLogEvent) => void): Promise<void>;
}
