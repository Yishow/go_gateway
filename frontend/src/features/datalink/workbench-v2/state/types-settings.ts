import type { DatabaseDeliveryOutcomeFields } from '../../../../types/databaseDelivery';

export type DatabaseSaveState = 'idle' | 'draft-invalid' | 'saving' | 'saved' | 'save-error';

export interface DbConnector extends DatabaseDeliveryOutcomeFields {
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  name: string;
  host: string;
  port: number;
  database: string;
  username: string;
  password?: string;
  /** 連線身分被換掉、既有密碼已不適用時為 true，需重新輸入才能存檔。 */
  password_required?: boolean;
  schema: string;
  table: string;
  write_mode: 'insert' | 'upsert';
  write_interval_seconds: number;
  timestamp_column: string;
  status: string;
  connector_id?: string;
  workspace_id?: string;
  persisted?: boolean;
  save_state?: DatabaseSaveState;
  save_error?: string | null;
}

export interface DbRowGroup {
  id: string;
  connector_id?: string;
  table_schema?: string;
  table_name: string;
  member_point_ids: string[];
  group_key_columns?: string[];
  unique_key_columns?: string[];
}

export interface DbTarget {
  tag_id: string;
  column_name: string;
  enabled: boolean;
  row_group_id?: string;
  point_id?: string;
  row_id?: string;
  workspace_id?: string;
  persisted?: boolean;
  save_state?: DatabaseSaveState;
  save_error?: string | null;
}

export interface Connector {
  id: string;
  name: string;
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  config: Record<string, unknown>;
}

export interface SettingsConnector {
  id: string;
  name: string;
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  host: string;
  port: number;
  database: string;
  username: string;
  password?: string;
  /** 連線身分被換掉、既有密碼已不適用時為 true，需重新輸入才會送出憑證。 */
  password_required?: boolean;
  schema: string;
  table: string;
  enabled: boolean;
  status: 'unknown' | 'testing' | 'ready' | 'unreachable' | 'auth_failed' | 'error';
  last_check_at?: string;
  last_check_error?: string;
  default_write_interval_seconds: number;
}

export interface TimeseriesSettings {
  write_precision: 'second' | 'millisecond';
  partition_interval: 'daily' | 'weekly' | 'monthly';
  batch_size: number;
  retention_days: number;
}

export interface SchedulerSettings {
  default_interval_ms: number;
  default_retry_count: number;
  default_retry_delay_ms: number;
  breaker_threshold: number;
  auto_start: boolean;
}

export interface ModbusShareSettings {
  enabled: boolean;
  bind_address: string;
  port: number;
  slave_id: number;
  capacity_registers: number;
  settings_revision: string;
  expected_settings_revision?: string;
  save_state?: 'idle' | 'saving' | 'saved' | 'save-error';
  save_error?: string | null;
  base_register: number;
}

export interface GeneralSettings {
  theme: 'dark' | 'light' | 'auto';
  locale: 'zh-TW' | 'en';
  addr_format: 'modbus' | 'hex' | 'raw';
  api_base: string;
  api_version: 'v1' | 'v2';
  timeout_seconds: number;
  log_level: 'trace' | 'debug' | 'info' | 'warn' | 'error';
  sse_heartbeat_seconds: number;
  enable_debug_panel: boolean;
  enable_audit_log: boolean;
}

export interface Settings {
  connectors: SettingsConnector[];
  timeseries: TimeseriesSettings;
  scheduler: SchedulerSettings;
  modbus_share: ModbusShareSettings;
  general: GeneralSettings;
}

export interface CommitLog {
  label: string;
  detail: string;
  status?: 'pending' | 'running' | 'success' | 'failed' | 'skipped';
  code?: string;
  retryable?: boolean;
  action?: string;
  request_id?: string;
}

export interface CommitState {
  status: 'idle' | 'committing' | 'success' | 'failed';
  logs: CommitLog[];
  started_at?: string;
  finished_at?: string;
}
