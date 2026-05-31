import type { Settings, SettingsConnector } from './types';

/**
 * 預設系統設定值
 * 落地設計決策：「Settings 預設值集中管理」
 */
export const DEFAULT_SETTINGS: Settings = {
  connectors: [
    {
      id: 'conn-prod',
      name: 'TimeSeries Prod',
      kind: 'postgres',
      host: 'tsdb.internal',
      port: 5432,
      database: 'gateway_metrics',
      username: 'gw_writer',
      password: '',
      schema: 'public',
      table: 'sensor_readings',
      enabled: true,
      status: 'ready',
      last_check_at: '2026-05-29T03:00:00.000Z',
      default_write_interval_seconds: 5
    }
  ],
  timeseries: {
    write_precision: 'millisecond',
    partition_interval: 'daily',
    batch_size: 500,
    retention_days: 90
  },
  scheduler: {
    default_interval_ms: 1000,
    default_retry_count: 3,
    default_retry_delay_ms: 500,
    breaker_threshold: 10,
    auto_start: true
  },
  modbus_share: {
    enabled: true, // 預設啟用以利 Step 2 進行 Modbus share 運算展示
    bind_address: '0.0.0.0',
    port: 5020,
    slave_id: 1,
    base_register: 40001
  },
  general: {
    theme: 'dark',
    locale: 'zh-TW',
    addr_format: 'modbus',
    api_base: 'http://localhost:8080',
    api_version: 'v1',
    timeout_seconds: 30,
    log_level: 'info',
    sse_heartbeat_seconds: 15,
    enable_debug_panel: false,
    enable_audit_log: true
  }
};

/**
 * 產生預設的連接器實體
 * 
 * @param idx 連接器索引順序
 * @returns 預設的連接器設定物件
 */
export function makeDefaultConnector(idx: number): SettingsConnector {
  return {
    id: `conn-${Date.now()}-${idx}`,
    name: `新連線 ${idx}`,
    kind: 'postgres',
    host: '127.0.0.1',
    port: 5432,
    database: 'gateway_metrics',
    username: 'postgres',
    password: '',
    schema: 'public',
    table: 'sensor_readings',
    enabled: true,
    status: 'unknown',
    default_write_interval_seconds: 5
  };
}
