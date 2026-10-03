import type {
  WriteGroupCreateRequest,
  WriteGroupDeleteRequest,
  WriteGroupDraft,
  WriteGroupListResponse,
  WriteGroupMutationResponse,
  WriteGroupReadiness,
  WriteGroupSchemaApplyRequest,
  WriteGroupSchemaPreviewRequest,
  WriteGroupUpdateRequest,
} from '../types/studioV2WriteGroup';
import type { SchemaOperation, SchemaPreviewToken } from '../types/recordingPlan';
import {
  parseWriteGroupAPIData,
  parseWriteGroupListData,
  parseWriteGroupMutationData,
  parseWriteGroupReadinessData,
} from '../utils/studioV2WriteGroupJson';
import type { WriteGroupDelivery } from '../types/studioV2WriteGroupDelivery';
import type {
  WriteGroupTestWriteConfirmation,
  WriteGroupTestWriteOperation,
  WriteGroupTestWritePreview,
} from '../types/studioV2WriteGroupTestWrite';
import {
  parseWriteGroupTestWriteOperation,
  parseWriteGroupTestWritePreview,
} from '../utils/studioV2WriteGroupTestWriteJson';
import { parseWriteGroupDeliveryData } from '../utils/studioV2WriteGroupDeliveryJson';
import {
  parseRecordingAPIData,
  parseRecordingSchemaOperation,
  parseRecordingSchemaPreviewToken,
} from '../utils/recordingPlanJson';
import { studioV2DatalinkApi } from './studioV2Workspace';

function cloneWriteGroupDraft(draft: WriteGroupDraft): WriteGroupDraft {
  return {
    workspace_id: draft.workspace_id,
    name: draft.name,
    members: draft.members.map((member) => ({
      device_id: member.device_id,
      point_id: member.point_id,
      tag_id: member.tag_id,
      ...(member.entity_key ? { entity_key: member.entity_key } : {}),
      ...(member.source_revision ? { source_revision: member.source_revision } : {}),
      ...(member.mapping_revision ? { mapping_revision: member.mapping_revision } : {}),
      ...(member.measurement_id ? { measurement_id: member.measurement_id } : {}),
      target_column: member.target_column?.trim() ?? '',
      required: member.required,
      ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
    })),
    destination: {
      connector_id: draft.destination.connector_id,
      connector_revision: draft.destination.connector_revision,
      table_schema: draft.destination.table_schema,
      table_name: draft.destination.table_name,
      storage_strategy: draft.destination.storage_strategy,
    },
    row_policy: {
      interval_seconds: draft.row_policy.interval_seconds,
      allowed_lateness_seconds: draft.row_policy.allowed_lateness_seconds,
      ...(draft.row_policy.incomplete_policy ? { incomplete_policy: draft.row_policy.incomplete_policy } : {}),
      ...(draft.row_policy.entity_key_column ? { entity_key_column: draft.row_policy.entity_key_column } : {}),
      ...(draft.row_policy.group_key_columns ? { group_key_columns: [...draft.row_policy.group_key_columns] } : {}),
      ...(draft.row_policy.unique_key_columns ? { unique_key_columns: [...draft.row_policy.unique_key_columns] } : {}),
      ...(draft.row_policy.value_column ? { value_column: draft.row_policy.value_column } : {}),
      ...(draft.row_policy.quality_column ? { quality_column: draft.row_policy.quality_column } : {}),
      ...(draft.row_policy.provenance_column ? { provenance_column: draft.row_policy.provenance_column } : {}),
      ...(draft.row_policy.record_key_column ? { record_key_column: draft.row_policy.record_key_column } : {}),
      ...(draft.row_policy.bucket_start_column ? { bucket_start_column: draft.row_policy.bucket_start_column } : {}),
      ...(draft.row_policy.group_id_column ? { group_id_column: draft.row_policy.group_id_column } : {}),
      ...(draft.row_policy.device_id_column ? { device_id_column: draft.row_policy.device_id_column } : {}),
    },
    write_policy: {
      ...(draft.write_policy.mode ? { mode: draft.write_policy.mode } : {}),
      ...(draft.write_policy.dedupe_capability ? { dedupe_capability: draft.write_policy.dedupe_capability } : {}),
    },
  };
}

function cloneCreateRequest(request: WriteGroupCreateRequest): WriteGroupCreateRequest {
  return { ...request, group: cloneWriteGroupDraft(request.group) };
}

function cloneUpdateRequest(request: WriteGroupUpdateRequest): WriteGroupUpdateRequest {
  return { ...request, group: cloneWriteGroupDraft(request.group) };
}

export const studioV2WorkspaceWriteGroupsAPI = {
  async list(): Promise<WriteGroupListResponse> {
    const res = await studioV2DatalinkApi.get<unknown>('/studio-v2/workspace/write-groups');
    return parseWriteGroupAPIData(res.data, 'write-group list', parseWriteGroupListData);
  },

  async get(id: string): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group get', parseWriteGroupMutationData);
  },

  async readiness(id: string): Promise<WriteGroupReadiness> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/readiness`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group readiness', parseWriteGroupReadinessData);
  },

  /** Read-only managed schema preview; the server resolves table and SQL from the saved group. */
  async schemaPreview(id: string, request: WriteGroupSchemaPreviewRequest): Promise<SchemaPreviewToken> {
    const payload: WriteGroupSchemaPreviewRequest = {
      workspace_id: request.workspace_id,
      expected_workspace_revision: request.expected_workspace_revision,
      expected_group_revision: request.expected_group_revision,
      expected_connector_revision: request.expected_connector_revision,
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/schema-preview`,
      payload,
    );
    return parseRecordingAPIData(
      res.data,
      'write-group schema preview',
      (data) => parseRecordingSchemaPreviewToken(data, { requireGroupLayout: true }),
    );
  },

  /** Explicitly confirms the stored preview operation; no client table or SQL is accepted. */
  async schemaApply(id: string, request: WriteGroupSchemaApplyRequest): Promise<SchemaOperation> {
    const payload: WriteGroupSchemaApplyRequest = {
      workspace_id: request.workspace_id,
      expected_workspace_revision: request.expected_workspace_revision,
      expected_group_revision: request.expected_group_revision,
      expected_connector_revision: request.expected_connector_revision,
      token: request.token,
      operation_id: request.operation_id,
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/schema-apply`,
      payload,
    );
    return parseRecordingAPIData(res.data, 'write-group schema apply', parseRecordingSchemaOperation);
  },

  /** Reads the durable operation once; callers must never retry DDL automatically. */
  async schemaOperation(operationId: string): Promise<SchemaOperation> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/database-operations/${encodeURIComponent(operationId)}`,
    );
    return parseRecordingAPIData(res.data, 'write-group schema operation', parseRecordingSchemaOperation);
  },

  /** Durable delivery truth: only `sql_committed` rows were confirmed by the destination. */
  async delivery(id: string): Promise<WriteGroupDelivery> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/delivery`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group delivery', parseWriteGroupDeliveryData);
  },

  /** Describes the exact test row and its cleanup; the target is not written. */
  async testWritePreview(id: string): Promise<WriteGroupTestWritePreview> {
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/test-write-preview`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group test write preview', parseWriteGroupTestWritePreview);
  },

  /**
   * Confirms a preview. The result is the recorded operation (202 while it is
   * still running); write verification and cleanup are separate fields.
   */
  async testWrite(id: string, request: WriteGroupTestWriteConfirmation): Promise<WriteGroupTestWriteOperation> {
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/test-write`,
      { token: request.token, operation_id: request.operation_id },
    );
    return parseWriteGroupAPIData(res.data, 'write-group test write', parseWriteGroupTestWriteOperation);
  },

  /** Reads a test-write operation from the shared workspace operation endpoint. */
  async testWriteOperation(operationId: string): Promise<WriteGroupTestWriteOperation> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/database-operations/${encodeURIComponent(operationId)}`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group test write operation', parseWriteGroupTestWriteOperation);
  },

  async create(request: WriteGroupCreateRequest): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups',
      cloneCreateRequest(request),
    );
    return parseWriteGroupAPIData(res.data, 'write-group create', parseWriteGroupMutationData);
  },

  async update(id: string, request: WriteGroupUpdateRequest): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.put<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`,
      cloneUpdateRequest(request),
    );
    return parseWriteGroupAPIData(res.data, 'write-group update', parseWriteGroupMutationData);
  },

  async remove(id: string, request: WriteGroupDeleteRequest): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.delete<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`,
      { data: { ...request } },
    );
    return parseWriteGroupAPIData(res.data, 'write-group delete', parseWriteGroupMutationData);
  },

  async disable(id: string, request: WriteGroupDeleteRequest): Promise<WriteGroupMutationResponse> {
    const payload: WriteGroupDeleteRequest = {
      workspace_id: request.workspace_id,
      expected_workspace_revision: request.expected_workspace_revision,
      expected_group_revision: request.expected_group_revision,
      expected_connector_revision: request.expected_connector_revision,
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/disable`,
      payload,
    );
    return parseWriteGroupAPIData(res.data, 'write-group disable', parseWriteGroupMutationData);
  },

  /**
   * Schedules the saved draft for the next UTC bucket once the server says it is
   * ready. It carries only expected revisions, never a group payload.
   */
  async apply(id: string, request: WriteGroupDeleteRequest): Promise<WriteGroupMutationResponse> {
    const payload: WriteGroupDeleteRequest = {
      workspace_id: request.workspace_id,
      expected_workspace_revision: request.expected_workspace_revision,
      expected_group_revision: request.expected_group_revision,
      expected_connector_revision: request.expected_connector_revision,
    };
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/apply`,
      payload,
    );
    return parseWriteGroupAPIData(res.data, 'write-group apply', parseWriteGroupMutationData);
  },
};
