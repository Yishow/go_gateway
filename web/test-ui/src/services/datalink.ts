/**
 * Datalink API 服務層
 *
 * 封裝所有 Datalink REST API 呼叫
 */

import axios from 'axios';
import type {
  Device,
  CreateDeviceRequest,
  UpdateDeviceRequest,
  ConnectionTestResult,
  ProtocolInfo,
  Point,
  CreatePointRequest,
  UpdatePointRequest,
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
} from '../types/datalink';

// API 基礎路徑
const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1';
const DATALINK_BASE = `${API_BASE}/datalink`;

// 建立 axios 實例
const api = axios.create({
  baseURL: DATALINK_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
});

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

export const tagAPI = {
  /** 列出標籤 */
  async list(params?: {
    status?: string;
    data_type?: string;
    key_prefix?: string;
    limit?: number;
    offset?: number;
  }): Promise<Tag[]> {
    const res = await api.get<APIResponse<Tag[]>>('/tags', { params });
    return res.data.data ?? [];
  },

  /** 取得標籤 */
  async get(id: string): Promise<Tag> {
    const res = await api.get<APIResponse<Tag>>(`/tags/${id}`);
    return res.data.data!;
  },

  /** 建立標籤 */
  async create(data: CreateTagRequest): Promise<Tag> {
    const res = await api.post<APIResponse<Tag>>('/tags', data);
    return res.data.data!;
  },

  /** 更新標籤 */
  async update(id: string, data: UpdateTagRequest): Promise<Tag> {
    const res = await api.put<APIResponse<Tag>>(`/tags/${id}`, data);
    return res.data.data!;
  },

  /** 刪除標籤 */
  async delete(id: string): Promise<void> {
    await api.delete(`/tags/${id}`);
  },

  /** 啟用標籤 */
  async activate(id: string): Promise<Tag> {
    const res = await api.post<APIResponse<Tag>>(`/tags/${id}/activate`);
    return res.data.data!;
  },

  /** 退役標籤 */
  async retire(id: string): Promise<Tag> {
    const res = await api.post<APIResponse<Tag>>(`/tags/${id}/retire`);
    return res.data.data!;
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
    const updates: Promise<any>[] = [];

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
