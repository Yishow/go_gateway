/**
 * Datalink 模組 TypeScript 型別定義
 *
 * 對應後端 internal/datalink/schema/models.go
 */

// =============================================================================
// 基礎型別
// =============================================================================

/** 協議類型 */
export type ProtocolType =
  | "modbus_tcp"
  | "modbus_rtu"
  | "modbus_udp"
  | "fatek_fbs"
  | "mc_3e"
  | "mqtt";

/** 設備狀態 */
export type DeviceStatus = "draft" | "active" | "disabled";

/** 資料類型 */
export type DataType =
  | "bool"
  | "int16"
  | "int32"
  | "int64"
  | "uint16"
  | "uint32"
  | "uint64"
  | "float32"
  | "float64"
  | "string";

/** 標籤狀態 */
export type TagStatus = "draft" | "active" | "retired";

/** 轉換類型 */
export type TransformType =
  | "decode"
  | "cast"
  | "scale"
  | "lookup"
  | "conditional"
  | "formula";

/** 時間精度 */
export type TimePrecision = "second" | "millisecond";

/** 分區間隔 */
export type PartitionInterval = "daily" | "weekly" | "monthly";

// =============================================================================
// 設備相關
// =============================================================================

/** 設備 */
export interface Device {
  id: string;
  name: string;
  description: string;
  protocol: ProtocolType;
  status: DeviceStatus;
  connection_config: string;
  last_test_at: string | null;
  last_test_success: boolean | null;
  last_test_error: string;
  readiness_status?: "ready" | "not_ready" | "error";
  created_at: string;
  updated_at: string;
}

/** 建立設備請求 */
export interface CreateDeviceRequest {
  name: string;
  description?: string;
  protocol: ProtocolType;
  connection_config: Record<string, unknown>;
}

/** 更新設備請求 */
export interface UpdateDeviceRequest {
  name?: string;
  description?: string;
  connection_config?: Record<string, unknown>;
}

/** 連線測試結果 */
export interface ConnectionTestResult {
  success: boolean;
  error: string;
  latency_ms: number;
}

/** 協議資訊 */
export interface ProtocolInfo {
  type: string;
  name: string;
  description: string;
  config_schema: string;
}

// =============================================================================
// 點位相關
// =============================================================================

/** 點位 */
export interface Point {
  id: string;
  device_id: string;
  name: string;
  description: string;
  data_type: DataType;
  address: string;
  enabled: boolean;
  polling_group_id: string;
  last_value: unknown;
  last_read_at: string;
  last_error: string;
  error_count: number;
  created_at: string;
  updated_at: string;
}

/** 建立點位請求 */
export interface CreatePointRequest {
  device_id: string;
  name: string;
  description?: string;
  data_type: DataType;
  address: string;
  enabled?: boolean;
  polling_group_id?: string;
}

/** 更新點位請求 */
export interface UpdatePointRequest {
  name?: string;
  description?: string;
  data_type?: DataType;
  address?: string;
  enabled?: boolean;
  polling_group_id?: string;
}

/** 輪詢群組 */
export interface PollingGroup {
  id: string;
  name: string;
  description?: string;
  interval_ms: number;
  priority: number;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

/** 建立輪詢群組請求 */
export interface CreatePollingGroupRequest {
  name: string;
  interval_ms: number;
  priority?: number;
  enabled?: boolean;
}

/** 更新輪詢群組請求 */
export interface UpdatePollingGroupRequest {
  name?: string;
  interval_ms?: number;
  priority?: number;
  enabled?: boolean;
}

// =============================================================================
// 標籤相關
// =============================================================================

/** 標籤 */
export interface Tag {
  id: string;
  key: string;
  display_name: string;
  description: string;
  data_type: DataType;
  unit: string;
  labels: Record<string, string> | null;
  status: TagStatus;
  created_at: string;
  updated_at: string;
}

/** 建立標籤請求 */
export interface CreateTagRequest {
  key: string;
  display_name?: string;
  description?: string;
  data_type: DataType;
  unit?: string;
  labels?: Record<string, string>;
}

/** 更新標籤請求 */
export interface UpdateTagRequest {
  display_name?: string;
  description?: string;
  unit?: string;
  labels?: Record<string, string>;
}

// =============================================================================
// 映射相關
// =============================================================================

/** 轉換步驟 */
export interface TransformStep {
  type: TransformType;
  order: number;
  params: Record<string, unknown>;
}

/** 映射 */
export interface Mapping {
  id: string;
  point_id: string;
  tag_id: string;
  enabled: boolean;
  transform_pipeline: string;
  created_at: string;
  updated_at: string;
}

/** 建立映射請求 */
export interface CreateMappingRequest {
  point_id: string;
  tag_id: string;
  enabled?: boolean;
  transform_pipeline?: TransformStep[];
}

/** 更新映射請求 */
export interface UpdateMappingRequest {
  enabled?: boolean;
  transform_pipeline?: TransformStep[];
}

/** 映射預覽請求 */
export interface MappingPreviewRequest {
  raw_value?: unknown;
  point_id?: string;
  mapping_id?: string;
  transform_pipeline?: TransformStep[];
}

/** 步驟結果 */
export interface StepResult {
  step_index: number;
  step_type: TransformType;
  input_value: unknown;
  output_value: unknown;
  error: string;
}

/** 映射預覽回應 */
export interface MappingPreviewResponse {
  raw_value: unknown;
  final_value: unknown;
  step_results: StepResult[];
  pipeline?: TransformStep[];
  error?: string;
}

// =============================================================================
// 系統設定相關
// =============================================================================

/** 系統設定 (聚合視圖) */
export interface SystemSettings {
  write_precision: TimePrecision;
  partition_interval: PartitionInterval;
  batch_size: number;
  default_retry_count?: number;
  default_retry_delay?: number;
}

/** 更新系統設定請求 */
export interface UpdateSystemSettingsRequest {
  write_precision?: TimePrecision;
  partition_interval?: PartitionInterval;
  batch_size?: number;
  default_retry_count?: number;
  default_retry_delay?: number;
}

/** 設定項目 (原始 API 回應) */
export interface SettingItem {
  key: string;
  value: unknown;
  description: string;
  updated_at: string;
}

// =============================================================================
// API 回應
// =============================================================================

/** 標準 API 回應 */
export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: string;
  };
  meta?: {
    total?: number;
    limit?: number;
    offset?: number;
  };
}

// =============================================================================
// 輪詢結果
// =============================================================================

/** 輪詢結果 */
export interface PollResult {
  point_id: string;
  value: unknown;
  timestamp: string;
  quality: number;
  error: string;
}
