import type { APIResponse } from '../types/datalink';
import type {
  ConnectorCapability,
  RecordingPlan,
  SchemaOperation,
  SchemaPreviewToken,
  TestWriteResult,
} from '../types/recordingPlan';
import { studioV2DatalinkApi } from './studioV2Workspace';
import { MAX_SAFE_JSON_ARRAY_LENGTH } from '../utils/safeJson';
import {
  parseRecordingAPIData,
  parseRecordingConnectorCapability,
  parseRecordingSchemaOperation,
  parseRecordingSchemaPreviewToken,
  parseRecordingTestWriteResult,
} from '../utils/recordingPlanJson';

export interface SchemaPreviewRequest {
  plan_id: string;
  connector_id: string;
  expected_connector_revision?: string;
  expected_workspace_revision: string;
  expected_plan_revision: string;
  table_prefix?: string;
  dialect: string;
}

/** A schema apply carries the issued operation and the revisions it was previewed against. */
export interface SchemaApplyRequest {
  token: string;
  operation_id: string;
  expected_workspace_revision: string;
  expected_plan_revision: string;
  expected_connector_revision: string;
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
    return parseRecordingAPIData(res.data, 'recording capability list', (data) => {
      if (Array.isArray(data)) {
        if (data.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
        const capabilities = data.map(parseRecordingConnectorCapability);
        return capabilities.some((capability) => capability === null)
          ? null
          : capabilities as ConnectorCapability[];
      }
      const capability = parseRecordingConnectorCapability(data);
      return capability ? [capability] : null;
    });
  },

  async schemaPreview(request: SchemaPreviewRequest): Promise<SchemaPreviewToken> {
    const res = await studioV2DatalinkApi.post<APIResponse<SchemaPreviewToken>>(
      '/studio-v2/workspace/recording-plans/schema-preview',
      request
    );
    return parseRecordingAPIData(res.data, 'recording schema preview', parseRecordingSchemaPreviewToken);
  },

  /** Confirms a preview; the response is the recorded operation, never a bare success flag. */
  async schemaApplyConfirmed(request: SchemaApplyRequest): Promise<SchemaOperation> {
    const res = await studioV2DatalinkApi.post<APIResponse<unknown>>(
      '/studio-v2/workspace/recording-plans/schema-apply',
      request
    );
    return parseRecordingAPIData(res.data, 'recording schema apply', parseRecordingSchemaOperation);
  },

  /** Reads the durable state of one operation of the current workspace. */
  async schemaOperation(operationId: string): Promise<SchemaOperation> {
    const res = await studioV2DatalinkApi.get<APIResponse<unknown>>(
      `/studio-v2/workspace/database-operations/${operationId}`
    );
    return parseRecordingAPIData(res.data, 'recording schema operation', parseRecordingSchemaOperation);
  },

  async testWrite(request: TestWriteRequest): Promise<TestWriteResult> {
    const res = await studioV2DatalinkApi.post<APIResponse<TestWriteResult>>(
      '/studio-v2/workspace/recording-plans/test-write',
      request
    );
    return parseRecordingAPIData(res.data, 'recording test write', parseRecordingTestWriteResult);
  },
};
