import type { DataType } from "./datalink";

/** 點位 */
export interface Point {
  id: string;
  device_id: string;
  name: string;
  description: string;
  data_type: DataType;
  /** Compatibility projection only; use the scoped format read API for effective policy. */
  data_format?: string;
  revision_id?: string;
  format_mode?: "legacy" | "inherit" | "override";
  format_override?: string;
  format_version?: number;
  function?: string;
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
  data_format?: string;
  address: string;
  enabled?: boolean;
  polling_group_id?: string;
}

/** 更新點位請求 */
export interface UpdatePointRequest {
  name?: string;
  description?: string;
  data_type?: DataType;
  data_format?: string;
  address?: string;
  enabled?: boolean;
  polling_group_id?: string;
}

