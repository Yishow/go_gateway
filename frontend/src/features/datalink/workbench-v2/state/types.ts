import type { StudioV2AvailabilityStatus } from '../../../../types/studioV2Availability';

/**
 * Workbench V2 資料狀態型別定義
 * 
 * 落地設計決策：「State 管理：本地 React state + 可序列化形狀」
 * 嚴格對齊「new_prototype/docs/_overview.md」的 State shape，禁用 any。
 */

export type ProtocolId = 'modbus_tcp' | 'modbus_rtu' | 'modbus_udp' | 'fatek_fbs' | 'mc_3e' | 'mqtt';

export interface ReadinessStage {
  id: string; // e.g. 'resolve', 'connect', 'probe'
  label: string; // label or key
  group: 'connect' | 'probe';
  status: 'pending' | 'running' | 'success' | 'failed' | 'skipped';
  latency_ms?: number;
  message?: string;
}

export interface DeviceTest {
  status: 'running' | 'success' | 'failed';
  latency_ms?: number;
  stages: Record<string, {
    status: 'pending' | 'running' | 'success' | 'failed' | 'skipped';
    latency_ms?: number;
    message?: string;
  }>;
  tested_at?: string; // ISO timestamp
}

export type DeviceSaveState = 'idle' | 'draft-invalid' | 'saving' | 'saved' | 'save-error';
export type RuleSaveState = 'idle' | 'draft-invalid' | 'saving' | 'saved' | 'save-error';
export type MappingSaveState = 'idle' | 'draft-invalid' | 'saving' | 'saved' | 'save-error';
export type DatabaseSaveState = 'idle' | 'draft-invalid' | 'saving' | 'saved' | 'save-error';

export interface Device {
  id: string; // dev-xxx
  name: string;
  description: string;
  protocol: ProtocolId;
  config: Record<string, unknown>; // 依協定而異，禁止 any
  status: 'draft' | 'tested' | 'active';
  test: DeviceTest | null;
  persisted?: boolean;
  save_state?: DeviceSaveState;
  save_error?: string | null;
  runtime_apply_status?: StudioV2RuntimeApplyStatus | null;
  runtime_apply_message?: string | null;
  availability_status?: StudioV2AvailabilityStatus | null;
  availability_reason?: string | null;
  running?: boolean;
}

export interface Rule {
  id: string; // rule-xxx
  device_id: string; // 所屬裝置 ID
  workspace_id?: string;
  revision_id?: string;
  name: string; // 規則名稱 (例如 "Holding Registers")
  start_address: string; // 位址 (例如 "40001")
  count: number; // 點位數量
  data_type: 'bool' | 'int16' | 'int32' | 'int64' | 'uint16' | 'uint32' | 'uint64' | 'float32' | 'float64' | 'string';
  naming_prefix: string; // 點位名稱前綴
  enabled: boolean;
  scale_multiplier: number;
  scale_offset: number;
  data_format: '' | 'ABCD' | 'BADC' | 'CDAB' | 'DCBA'; // byte order
  skipped_addresses: string[]; // 被略過的位址
  share_enabled: boolean;
  share_start_register: number | null;
  share_stride: number | null;
  persisted?: boolean;
  save_state?: RuleSaveState;
  save_error?: string | null;
}

export interface ShareLayout {
  start: number;
  stride: number;
  end: number;
  auto: boolean;
}

export interface Point {
  id: string;
  device_id: string;
  rule_id: string;
  rule_name: string;
  name: string;
  address: string;
  data_type: string;
  function: 'coil' | 'discrete_input' | 'input_register' | 'holding_register';
  width: number;
  enabled: boolean;
  skipped: boolean;
  display?: string;
  unit?: string;
  tag_key_suggest?: string;
  _rule_scale: number;
  _rule_offset: number;
}

export type TargetType = 'bool' | 'int16' | 'int32' | 'int64' | 'uint16' | 'uint32' | 'uint64' | 'float32' | 'float64' | 'string';

export interface MappingValue {
  tag_key: string;
  display_name: string;
  unit: string;
  target_type: TargetType;
  scale: number;
  offset: number;
  enabled: boolean;
}

export interface Mapping {
  point_id: string;
  tag_key: string; // e.g., "line01.temp.inlet"
  display_name: string;
  unit: string;
  target_type: TargetType;
  scale: number;
  offset: number;
  enabled: boolean;
  mapping_id?: string;
  workspace_id?: string;
  rule_id?: string;
  device_id?: string;
  address?: string;
  tag_id?: string;
  persisted_point_id?: string;
  persisted?: boolean;
  local_value?: MappingValue;
  persisted_value?: MappingValue;
  save_state?: MappingSaveState;
  save_error?: string | null;
}

export interface DbConnector {
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  name: string;
  host: string;
  port: number;
  database: string;
  username: string;
  password?: string;
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

export interface DbTarget {
  tag_id: string;
  column_name: string;
  enabled: boolean;
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
  status: 'pending' | 'running' | 'success' | 'failed';
}

export interface CommitState {
  status: 'idle' | 'committing' | 'success' | 'failed';
  logs: CommitLog[];
  started_at?: string;
  finished_at?: string;
}

export interface WorkbenchV2State {
  view: 'flow' | 'settings';
  current: 1 | 2 | 3 | 4;
  completed: Set<number>;
  sidebarCollapsed: boolean;
  showSummaryRail: boolean;
  devices: Device[];
  rules: Rule[];
  selectedRuleId: string | null;
  points: Point[];
  mappings: Record<string, Mapping>; // key: pointId
  db: {
    connector: DbConnector;
    targets: Record<string, DbTarget>; // key: pointId
  };
  settings: Settings;
  commit?: CommitState;
  committed: boolean;
}
import type { StudioV2RuntimeApplyStatus } from '../../../../types/studioV2RuntimeApply';
