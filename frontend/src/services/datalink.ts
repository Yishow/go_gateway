/**
 * Datalink API 服務層
 *
 * 封裝所有 Datalink REST API 呼叫
 */

import axios from 'axios';
import { resolveGatewayUiVersion } from '../features/gateway/uiVersion';
import type {
  Device,
  CreateDeviceRequest,
  UpdateDeviceRequest,
  ConnectionTestResult,
  ProtocolInfo,
  Point,
  CreatePointRequest,
  UpdatePointRequest,
  SourceRuleRecord,
  CreateSourceRuleRequest,
  UpdateSourceRuleRequest,
  PollingGroup,
  CreatePollingGroupRequest,
  UpdatePollingGroupRequest,
  Tag,
  CreateTagRequest,
  UpdateTagRequest,
  Mapping,
  CreateMappingRequest,
  UpdateMappingRequest,
  MappingPreviewRequest,
  MappingPreviewResponse,
  SystemSettings,
  UpdateSystemSettingsRequest,
  SettingItem,
  APIResponse,
  PollResult,
  ModbusShareStatus,
  ModbusShareMapping,
  RuntimeStatus,
  DatabaseConnector,
  CreateDatabaseConnectorRequest,
  UpdateDatabaseConnectorRequest,
  DatabaseTableInfo,
  DatabaseTargetValidationResult,
  DatabaseTargetMapping,
  CreateDatabaseTargetMappingRequest,
  UpdateDatabaseTargetMappingRequest,
} from '../types/datalink';
import { VITE_API_BASE_URL } from '../env';

// API 基礎路徑
const API_BASE = VITE_API_BASE_URL;
const DATALINK_BASE = `${API_BASE}/datalink`;

interface DeviceReadinessResult {
  device_id: string;
  status: 'ready' | 'warning' | 'error';
  checks: Array<{
    name: string;
    pass: boolean;
    message: string;
  }>;
}

// 建立 axios 實例
const api = axios.create({
  baseURL: DATALINK_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 追加 UI 版本埋點 Header（僅雙入口頁）
api.interceptors.request.use((config) => {
  const uiVersion = resolveGatewayUiVersion();
  if (uiVersion) {
    config.headers = config.headers ?? {};
    if (!config.headers['X-UI-Version']) {
      config.headers['X-UI-Version'] = uiVersion;
    }
  }
  return config;
});

// 添加響應攔截器以統一處理錯誤
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // 提取錯誤訊息
    if (error.response?.data?.error?.message) {
      error.message = error.response.data.error.message;
    } else if (error.response?.data?.message) {
      error.message = error.response.data.message;
    } else if (error.message) {
      // 保留原始錯誤訊息
    } else {
      error.message = 'Request failed with status code ' + (error.response?.status || 'unknown');
    }
    return Promise.reject(error);
  }
);

// =============================================================================
// 設備 API
// =============================================================================

export const deviceAPI = {
  /** 列出設備 */
  async list(params?: {
    protocol?: string;
    status?: string;
    limit?: number;
    offset?: number;
  }): Promise<Device[]> {
    const res = await api.get<APIResponse<Device[]>>('/devices', { params });
    return res.data.data ?? [];
  },

  /** 取得設備 */
  async get(id: string): Promise<Device> {
    const res = await api.get<APIResponse<Device>>(`/devices/${id}`);
    return res.data.data!;
  },

  /** 建立設備 */
  async create(data: CreateDeviceRequest): Promise<Device> {
    const res = await api.post<APIResponse<Device>>('/devices', data);
    return res.data.data!;
  },

  /** 更新設備 */
  async update(id: string, data: UpdateDeviceRequest): Promise<Device> {
    const res = await api.put<APIResponse<Device>>(`/devices/${id}`, data);
    return res.data.data!;
  },

  /** 刪除設備 */
  async delete(id: string): Promise<void> {
    await api.delete(`/devices/${id}`);
  },

  /** 測試連線 */
  async testConnection(id: string): Promise<ConnectionTestResult> {
    const res = await api.post<APIResponse<ConnectionTestResult>>(
      `/devices/${id}/test`
    );
    return res.data.data!;
  },

  /** 批量測試連線 */
  async batchTestConnections(deviceIds: string[]): Promise<ConnectionTestResult[]> {
    const res = await api.post<APIResponse<ConnectionTestResult[]>>(
      '/devices/test-batch',
      { device_ids: deviceIds }
    );
    return res.data.data ?? [];
  },

  /** 檢查設備就緒狀態 */
  async checkReadiness(id: string): Promise<DeviceReadinessResult> {
    const res = await api.post<APIResponse<DeviceReadinessResult>>(`/devices/${id}/readiness`);
    return res.data.data!;
  },

  /** 啟用設備 */
  async activate(id: string): Promise<Device> {
    const res = await api.post<APIResponse<Device>>(`/devices/${id}/activate`);
    return res.data.data!;
  },

  /** 停用設備 */
  async disable(id: string): Promise<Device> {
    const res = await api.post<APIResponse<Device>>(`/devices/${id}/disable`);
    return res.data.data!;
  },
};

// =============================================================================
// 點位 API
// =============================================================================

export const pointAPI = {
  /** 列出點位 */
  async list(params?: {
    device_id?: string;
    polling_group_id?: string;
    enabled?: boolean;
    data_type?: string;
    limit?: number;
    offset?: number;
  }): Promise<Point[]> {
    const res = await api.get<APIResponse<Point[]>>('/points', { params });
    return res.data.data ?? [];
  },

  /** 取得點位 */
  async get(id: string): Promise<Point> {
    const res = await api.get<APIResponse<Point>>(`/points/${id}`);
    return res.data.data!;
  },

  /** 建立點位 */
  async create(data: CreatePointRequest): Promise<Point> {
    const res = await api.post<APIResponse<Point>>('/points', data);
    return res.data.data!;
  },

  /** 更新點位 */
  async update(id: string, data: UpdatePointRequest): Promise<Point> {
    const res = await api.put<APIResponse<Point>>(`/points/${id}`, data);
    return res.data.data!;
  },

  /** 刪除點位 */
  async delete(id: string): Promise<void> {
    await api.delete(`/points/${id}`);
  },

  /** 立即輪詢 */
  async pollNow(id: string): Promise<PollResult> {
    const res = await api.post<APIResponse<PollResult>>(`/points/${id}/poll`);
    return res.data.data!;
  },

  /** 批量輪詢 */
  async batchPoll(pointIds: string[]): Promise<PollResult[]> {
    const res = await api.post<APIResponse<PollResult[]>>('/points/poll', {
      point_ids: pointIds,
    });
    return res.data.data ?? [];
  },
};

// =============================================================================
// 來源規則 API
// =============================================================================

export const sourceRuleAPI = {
  /** 列出來源規則 */
  async list(params?: { device_id?: string; enabled?: boolean }): Promise<SourceRuleRecord[]> {
    const res = await api.get<APIResponse<SourceRuleRecord[]>>('/source-rules', { params });
    return res.data.data ?? [];
  },

  /** 取得來源規則 */
  async get(id: string): Promise<SourceRuleRecord> {
    const res = await api.get<APIResponse<SourceRuleRecord>>(`/source-rules/${id}`);
    return res.data.data!;
  },

  /** 建立來源規則 */
  async create(data: CreateSourceRuleRequest): Promise<SourceRuleRecord> {
    const res = await api.post<APIResponse<SourceRuleRecord>>('/source-rules', data);
    return res.data.data!;
  },

  /** 更新來源規則 */
  async update(id: string, data: UpdateSourceRuleRequest): Promise<SourceRuleRecord> {
    const res = await api.put<APIResponse<SourceRuleRecord>>(`/source-rules/${id}`, data);
    return res.data.data!;
  },

  /** 刪除來源規則 */
  async delete(id: string): Promise<void> {
    await api.delete(`/source-rules/${id}`);
  },

  /** 啟用來源規則 */
  async enable(id: string): Promise<SourceRuleRecord> {
    const res = await api.post<APIResponse<SourceRuleRecord>>(`/source-rules/${id}/enable`);
    return res.data.data!;
  },

  /** 停用來源規則 */
  async disable(id: string): Promise<SourceRuleRecord> {
    const res = await api.post<APIResponse<SourceRuleRecord>>(`/source-rules/${id}/disable`);
    return res.data.data!;
  },
};

// =============================================================================
// 輪詢群組 API
// =============================================================================

export const pollingGroupAPI = {
  /** 列出輪詢群組 */
  async list(): Promise<PollingGroup[]> {
    const res = await api.get<APIResponse<PollingGroup[]>>('/polling-groups');
    return res.data.data ?? [];
  },

  /** 建立輪詢群組 */
  async create(data: CreatePollingGroupRequest): Promise<PollingGroup> {
    const res = await api.post<APIResponse<PollingGroup>>(
      '/polling-groups',
      data
    );
    return res.data.data!;
  },

  /** 更新輪詢群組 */
  async update(
    id: string,
    data: UpdatePollingGroupRequest
  ): Promise<PollingGroup> {
    const res = await api.put<APIResponse<PollingGroup>>(
      `/polling-groups/${id}`,
      data
    );
    return res.data.data!;
  },

  /** 刪除輪詢群組 */
  async delete(id: string): Promise<void> {
    await api.delete(`/polling-groups/${id}`);
  },
};

// =============================================================================
// 標籤 API
// =============================================================================

type RawTag = Omit<Tag, 'labels'> & {
  labels: Record<string, string> | string | null;
};

const normalizeTagLabels = (labels: RawTag['labels']): Record<string, string> | null => {
  if (!labels) {
    return null;
  }
  if (typeof labels === 'string') {
    try {
      const parsed = JSON.parse(labels);
      if (parsed && typeof parsed === 'object') {
        return parsed as Record<string, string>;
      }
      return null;
    } catch {
      return null;
    }
  }
  return labels;
};

const normalizeTag = (tag: RawTag): Tag => ({
  ...tag,
  labels: normalizeTagLabels(tag.labels),
});

export const tagAPI = {
  /** 列出標籤 */
  async list(params?: {
    status?: string;
    data_type?: string;
    key_prefix?: string;
    limit?: number;
    offset?: number;
  }): Promise<Tag[]> {
    const res = await api.get<APIResponse<RawTag[]>>('/tags', { params });
    return (res.data.data ?? []).map(normalizeTag);
  },

  /** 取得標籤 */
  async get(id: string): Promise<Tag> {
    const res = await api.get<APIResponse<RawTag>>(`/tags/${id}`);
    return normalizeTag(res.data.data!);
  },

  /** 建立標籤 */
  async create(data: CreateTagRequest): Promise<Tag> {
    const res = await api.post<APIResponse<RawTag>>('/tags', data);
    return normalizeTag(res.data.data!);
  },

  /** 更新標籤 */
  async update(id: string, data: UpdateTagRequest): Promise<Tag> {
    const res = await api.put<APIResponse<RawTag>>(`/tags/${id}`, data);
    return normalizeTag(res.data.data!);
  },

  /** 刪除標籤 */
  async delete(id: string): Promise<void> {
    await api.delete(`/tags/${id}`);
  },

  /** 啟用標籤 */
  async activate(id: string): Promise<Tag> {
    const res = await api.post<APIResponse<RawTag>>(`/tags/${id}/activate`);
    return normalizeTag(res.data.data!);
  },

  /** 退役標籤 */
  async retire(id: string): Promise<Tag> {
    const res = await api.post<APIResponse<RawTag>>(`/tags/${id}/retire`);
    return normalizeTag(res.data.data!);
  },

  /** 批量建立 */
  async batchCreate(
    tags: CreateTagRequest[]
  ): Promise<{ created: string[]; errors: Array<{ key: string; error: string }> }> {
    const res = await api.post<
      APIResponse<{ created: string[]; errors: Array<{ key: string; error: string }> }>
    >('/tags/batch', { tags });
    return res.data.data!;
  },

  /** 驗證標籤鍵 */
  async validateKey(
    key: string
  ): Promise<{ valid: boolean; normalized: string; exists: boolean; error?: string }> {
    const res = await api.post<
      APIResponse<{ valid: boolean; normalized: string; exists: boolean; error?: string }>
    >('/tags/validate-key', { key });
    return res.data.data!;
  },
};

// =============================================================================
// 映射 API
// =============================================================================

export const mappingAPI = {
  /** 列出映射 */
  async list(params?: {
    point_id?: string;
    tag_id?: string;
    enabled?: boolean;
    limit?: number;
    offset?: number;
  }): Promise<Mapping[]> {
    const res = await api.get<APIResponse<Mapping[]>>('/mappings', { params });
    return res.data.data ?? [];
  },

  /** 取得映射 */
  async get(id: string): Promise<Mapping> {
    const res = await api.get<APIResponse<Mapping>>(`/mappings/${id}`);
    return res.data.data!;
  },

  /** 建立映射 */
  async create(data: CreateMappingRequest): Promise<Mapping> {
    const res = await api.post<APIResponse<Mapping>>('/mappings', data);
    return res.data.data!;
  },

  /** 更新映射 */
  async update(id: string, data: UpdateMappingRequest): Promise<Mapping> {
    const res = await api.put<APIResponse<Mapping>>(`/mappings/${id}`, data);
    return res.data.data!;
  },

  /** 刪除映射 */
  async delete(id: string): Promise<void> {
    await api.delete(`/mappings/${id}`);
  },

  /** 預覽映射 */
  async preview(data: MappingPreviewRequest): Promise<MappingPreviewResponse> {
    const res = await api.post<APIResponse<MappingPreviewResponse>>(
      '/mappings/preview',
      data
    );
    return res.data.data!;
  },

  /** 驗證管線 */
  async validatePipeline(
    pipeline: MappingPreviewRequest['transform_pipeline']
  ): Promise<{ valid: boolean; error?: string }> {
    const res = await api.post<
      APIResponse<{ valid: boolean; error?: string }>
    >('/mappings/validate-pipeline', { pipeline });
    return res.data.data!;
  },
};

// =============================================================================
// 設定 API
// =============================================================================

export const settingsAPI = {
  /** 取得原始設定清單（配置中心） */
  async listItems(): Promise<SettingItem[]> {
    const res = await api.get<APIResponse<SettingItem[]>>('/settings');
    return res.data.data ?? [];
  },

  /** 取得所有設定 (聚合) */
  async get(): Promise<SystemSettings> {
    const res = await api.get<APIResponse<SettingItem[]>>('/settings');
    const items = res.data.data ?? [];
    
    // Default values
    const settings: SystemSettings = {
      write_precision: 'second',
      partition_interval: 'monthly',
      batch_size: 1000,
    };

    // Map items to settings
    items.forEach(item => {
      if (item.key === 'write_precision') settings.write_precision = String(item.value) as SystemSettings['write_precision'];
      if (item.key === 'partition_interval') {
        settings.partition_interval = String(item.value) as SystemSettings['partition_interval'];
      }
      if (item.key === 'batch_size') settings.batch_size = Number(item.value);
      if (item.key === 'default_retry_count') settings.default_retry_count = Number(item.value);
      if (item.key === 'default_retry_delay') settings.default_retry_delay = Number(item.value);
    });

    return settings;
  },

  /** 更新設定 (聚合) */
  async update(data: UpdateSystemSettingsRequest): Promise<SystemSettings> {
    const updates: Array<Promise<SettingItem>> = [];

    if (data.write_precision) {
      updates.push(this.updateKey('write_precision', data.write_precision));
    }
    if (data.partition_interval) {
      updates.push(this.updateKey('partition_interval', data.partition_interval));
    }
    if (data.batch_size) {
      updates.push(this.updateKey('batch_size', data.batch_size));
    }
    if (data.default_retry_count !== undefined) {
      updates.push(this.updateKey('default_retry_count', data.default_retry_count));
    }
    if (data.default_retry_delay !== undefined) {
      updates.push(this.updateKey('default_retry_delay', data.default_retry_delay));
    }

    await Promise.all(updates);
    return this.get();
  },

  /** 更新單一設定 (內部用) */
  async updateKey(key: string, value: unknown): Promise<SettingItem> {
    const res = await api.put<APIResponse<SettingItem>>(`/settings/${key}`, {
      value,
    });
    return res.data.data!;
  },
};

// =============================================================================
// 協議 API
// =============================================================================

export const protocolAPI = {
  /** 列出協議 */
  async list(): Promise<ProtocolInfo[]> {
    const res = await api.get<APIResponse<ProtocolInfo[]>>('/protocols');
    return res.data.data ?? [];
  },
};

// =============================================================================
// 健康檢查
// =============================================================================

export const healthAPI = {
  /** 健康檢查 */
  async check(): Promise<{ status: string; service: string }> {
    const res = await api.get<APIResponse<{ status: string; service: string }>>(
      '/health'
    );
    return res.data.data!;
  },
};

// =============================================================================
// 儀表板 API
// =============================================================================

/** 儀表板統計數據 */
export interface DashboardStats {
  total_devices: number;
  active_devices: number;
  disabled_devices: number;
  error_devices: number;
  total_points: number;
  enabled_points: number;
  total_tags: number;
  active_tags: number;
  retired_tags: number;
  errors_last_24h: number;
  estimated_throughput: number;
}

/** 設備狀態資訊 */
export interface DeviceStatus {
  id: string;
  name: string;
  protocol: string;
  status: string;
  last_test_at?: string;
  last_success?: boolean;
  last_error?: string;
}

export const dashboardAPI = {
  /** 取得儀表板統計數據 */
  async getStats(): Promise<DashboardStats> {
    const res = await api.get<APIResponse<DashboardStats>>('/dashboard/stats');
    return res.data.data!;
  },

  /** 取得設備狀態列表 */
  async getDeviceStatuses(): Promise<DeviceStatus[]> {
    const res = await api.get<APIResponse<DeviceStatus[]>>('/dashboard/device-statuses');
    return res.data.data ?? [];
  },
};

export const runtimeAPI = {
  async getStatus(deviceId?: string): Promise<RuntimeStatus> {
    const res = await api.get<APIResponse<RuntimeStatus>>('/runtime/status', {
      params: deviceId ? { device_id: deviceId } : undefined,
    });
    return res.data.data ?? {
      running: false,
      uptime_seconds: 0,
      collectors: [],
    };
  },

  getStreamUrl(deviceId: string, pointIds?: string[]): string {
    const params = new URLSearchParams({ device_id: deviceId });
    if (pointIds && pointIds.length > 0) {
      params.set('point_ids', pointIds.join(','));
    }
    return `${DATALINK_BASE}/runtime/stream?${params.toString()}`;
  },
};

// =============================================================================
// Local Modbus Share API
// =============================================================================

export const modbusShareAPI = {
  async status(): Promise<ModbusShareStatus> {
    const res = await api.get<APIResponse<ModbusShareStatus>>('/modbus-share/status');
    return res.data.data!;
  },

  async start(port?: number): Promise<ModbusShareStatus> {
    const res = await api.post<APIResponse<ModbusShareStatus>>('/modbus-share/start', {
      port: typeof port === 'number' ? port : undefined,
    });
    return res.data.data!;
  },

  async stop(): Promise<ModbusShareStatus> {
    const res = await api.post<APIResponse<ModbusShareStatus>>('/modbus-share/stop', {});
    return res.data.data!;
  },

  async listMappings(): Promise<ModbusShareMapping[]> {
    const res = await api.get<APIResponse<ModbusShareMapping[]>>('/modbus-share/mappings');
    return res.data.data ?? [];
  },

  async upsertMapping(tagId: string, register: number): Promise<ModbusShareMapping> {
    const res = await api.put<APIResponse<ModbusShareMapping>>(`/modbus-share/mappings/${tagId}`, { register });
    return res.data.data!;
  },

  async deleteMapping(tagId: string): Promise<void> {
    await api.delete(`/modbus-share/mappings/${tagId}`);
  },

  async writeTagValue(tagId: string, value: unknown): Promise<void> {
    await api.post('/modbus-share/write-tag-value', { tag_id: tagId, value });
  },

  async sync(): Promise<{ updated: number; skipped: number; errors: string[] }> {
    const res = await api.post<APIResponse<{ updated: number; skipped: number; errors: string[] }>>('/modbus-share/sync');
    return res.data.data ?? { updated: 0, skipped: 0, errors: [] };
  },
};

export const dbTargetAPI = {
  async listConnectors(params?: { enabled?: boolean }): Promise<DatabaseConnector[]> {
    const res = await api.get<APIResponse<DatabaseConnector[]>>('/db-targets/connectors', {
      params,
    });
    return res.data.data ?? [];
  },

  async getConnector(id: string): Promise<DatabaseConnector> {
    const res = await api.get<APIResponse<DatabaseConnector>>(`/db-targets/connectors/${id}`);
    return res.data.data!;
  },

  async createConnector(
    data: CreateDatabaseConnectorRequest,
  ): Promise<DatabaseConnector> {
    const res = await api.post<APIResponse<DatabaseConnector>>(
      '/db-targets/connectors',
      data,
    );
    return res.data.data!;
  },

  async updateConnector(
    id: string,
    data: UpdateDatabaseConnectorRequest,
  ): Promise<DatabaseConnector> {
    const res = await api.put<APIResponse<DatabaseConnector>>(
      `/db-targets/connectors/${id}`,
      data,
    );
    return res.data.data!;
  },

  async deleteConnector(id: string): Promise<void> {
    await api.delete(`/db-targets/connectors/${id}`);
  },

  async testConnector(id: string): Promise<DatabaseConnector> {
    const res = await api.post<APIResponse<DatabaseConnector>>(
      `/db-targets/connectors/${id}/test`,
    );
    return res.data.data!;
  },

  async listTables(id: string): Promise<DatabaseTableInfo[]> {
    const res = await api.get<APIResponse<DatabaseTableInfo[]>>(
      `/db-targets/connectors/${id}/tables`,
    );
    return res.data.data ?? [];
  },

  async validateConnector(id: string): Promise<DatabaseTargetValidationResult> {
    const res = await api.get<APIResponse<DatabaseTargetValidationResult>>(
      `/db-targets/connectors/${id}/validate`,
    );
    return res.data.data ?? { ready: false, issues: [] };
  },

  async listMappings(params?: {
    connector_id?: string;
    tag_id?: string;
    enabled?: boolean;
  }): Promise<DatabaseTargetMapping[]> {
    const res = await api.get<APIResponse<DatabaseTargetMapping[]>>(
      '/db-targets/mappings',
      { params },
    );
    return res.data.data ?? [];
  },

  async getMapping(id: string): Promise<DatabaseTargetMapping> {
    const res = await api.get<APIResponse<DatabaseTargetMapping>>(`/db-targets/mappings/${id}`);
    return res.data.data!;
  },

  async createMapping(
    data: CreateDatabaseTargetMappingRequest,
  ): Promise<DatabaseTargetMapping> {
    const res = await api.post<APIResponse<DatabaseTargetMapping>>(
      '/db-targets/mappings',
      data,
    );
    return res.data.data!;
  },

  async updateMapping(
    id: string,
    data: UpdateDatabaseTargetMappingRequest,
  ): Promise<DatabaseTargetMapping> {
    const res = await api.put<APIResponse<DatabaseTargetMapping>>(
      `/db-targets/mappings/${id}`,
      data,
    );
    return res.data.data!;
  },

  async deleteMapping(id: string): Promise<void> {
    await api.delete(`/db-targets/mappings/${id}`);
  },
};
