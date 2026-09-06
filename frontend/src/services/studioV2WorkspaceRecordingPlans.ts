import type { APIResponse } from '../types/datalink';
import type {
  ConnectorCapability,
  RecordingPlan,
  SchemaPreviewToken,
  TestWriteResult,
} from '../types/recordingPlan';
import { studioV2DatalinkApi } from './studioV2Workspace';

export interface SchemaPreviewRequest {
  plan_id: string;
  connector_id: string;
  table_prefix?: string;
  dialect: string;
}

export interface TestWriteRequest {
  plan_id: string;
  stream_id?: string;
  table_prefix?: string;
}

export const studioV2WorkspaceRecordingPlansAPI = {
  async list(deviceId?: string): Promise<RecordingPlan[]> {
    const params = deviceId ? { device_id: deviceId } : undefined;
    const res = await studioV2DatalinkApi.get<APIResponse<RecordingPlan[]>>(
      '/studio-v2/workspace/recording-plans',
      { params }
    );
    return res.data.data ?? [];
  },

  async get(id: string): Promise<RecordingPlan> {
    const res = await studioV2DatalinkApi.get<APIResponse<RecordingPlan>>(
      `/studio-v2/workspace/recording-plans/${id}`
    );
    return res.data.data!;
  },

  async create(plan: Partial<RecordingPlan>): Promise<RecordingPlan> {
    const res = await studioV2DatalinkApi.post<APIResponse<RecordingPlan>>(
      '/studio-v2/workspace/recording-plans',
      plan
    );
    return res.data.data!;
  },

  async update(id: string, plan: Partial<RecordingPlan>): Promise<RecordingPlan> {
    const res = await studioV2DatalinkApi.put<APIResponse<RecordingPlan>>(
      `/studio-v2/workspace/recording-plans/${id}`,
      plan
    );
    return res.data.data!;
  },

  async remove(id: string): Promise<void> {
    await studioV2DatalinkApi.delete(`/studio-v2/workspace/recording-plans/${id}`);
  },

  async capabilities(kind?: string): Promise<ConnectorCapability[]> {
    const params = kind ? { kind } : undefined;
    const res = await studioV2DatalinkApi.get<APIResponse<ConnectorCapability[] | ConnectorCapability>>(
      '/studio-v2/workspace/recording-plans/capabilities',
      { params }
    );
    const data = res.data.data;
    if (Array.isArray(data)) {
      return data;
    }
    return data ? [data] : [];
  },

  async schemaPreview(request: SchemaPreviewRequest): Promise<SchemaPreviewToken> {
    const res = await studioV2DatalinkApi.post<APIResponse<SchemaPreviewToken>>(
      '/studio-v2/workspace/recording-plans/schema-preview',
      request
    );
    return res.data.data!;
  },

  async schemaApply(token: string): Promise<{ applied: boolean; message: string }> {
    const res = await studioV2DatalinkApi.post<APIResponse<{ applied: boolean; message: string }>>(
      '/studio-v2/workspace/recording-plans/schema-apply',
      { token }
    );
    return res.data.data!;
  },

  async testWrite(request: TestWriteRequest): Promise<TestWriteResult> {
    const res = await studioV2DatalinkApi.post<APIResponse<TestWriteResult>>(
      '/studio-v2/workspace/recording-plans/test-write',
      request
    );
    return res.data.data!;
  },
};
