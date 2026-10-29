import type { StudioV2LocatorReadScope } from '../types/studioV2LocatorEntry';
import { captureStudioLocatorReadScope, studioLocatorReadParams, assertStudioLocatorWorkspaceBinding, throwStudioLocatorReadError } from './studioEntityLocator';
import type { DeliveryDecision, DeliveryDecisionResult } from '../types/studioV2WriteGroupDelivery';
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
  parseWriteGroupTestWritePreviewAPIData,
} from '../utils/studioV2WriteGroupTestWriteJson';
import { parseWriteGroupDeliveryData } from '../utils/studioV2WriteGroupDeliveryJson';
import { cloneWriteGroupExactRecordingSource } from '../utils/studioV2ExactRecordingSourceJson';
import {
  parseRecordingAPIData,
  parseRecordingSchemaOperation,
  parseRecordingSchemaPreviewToken,
} from '../utils/recordingPlanJson';
import type { WriteGroupCutover, WriteGroupCutoverRequest, WriteGroupCutoverReview } from '../types/studioV2WriteGroupCutover';
import { cutoverBody, parseWriteGroupCutover, parseWriteGroupCutoverList, parseWriteGroupCutoverReview } from '../utils/studioV2WriteGroupCutoverJson';
import { studioV2DatalinkApi } from './studioV2Workspace';

const WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH = 8;

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
      ...(member.circuit_id ? { circuit_id: member.circuit_id } : {}),
      ...(member.circuit_name ? { circuit_name: member.circuit_name } : {}),
      ...(member.value_role !== undefined ? { value_role: member.value_role } : {}),
      ...(member.source_value_type !== undefined ? { source_value_type: member.source_value_type } : {}),
      ...(member.result_value_type !== undefined ? { result_value_type: member.result_value_type } : {}),
      ...(member.exact_recording_source !== undefined ? { exact_recording_source: cloneWriteGroupExactRecordingSource(member.exact_recording_source) } : {}),
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
      ...(draft.write_policy.recreate_missing !== undefined ? { recreate_missing: draft.write_policy.recreate_missing } : {}),
    },
    ...(draft.layout_version ? { layout_version: draft.layout_version } : {}),
    ...(draft.recording_profile ? { recording_profile: { ...draft.recording_profile } } : {}),
    ...(draft.storage_mode ? { storage_mode: draft.storage_mode } : {}),
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
    return parseWriteGroupAPIData(res.data, 'write-group list', parseWriteGroupListData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
  },

  async get(id: string): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group get', parseWriteGroupMutationData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
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

  /** Read-only difference of a forward cutover (or its rollback) at an exact instant. */
  async reviewCutover(sourceId: string, request: WriteGroupCutoverRequest): Promise<WriteGroupCutoverReview> {
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(sourceId)}/cutover-review`,
      cutoverBody(request),
    );
    return parseWriteGroupAPIData(res.data, 'write-group cutover review', parseWriteGroupCutoverReview);
  },

  /** Applies exactly what was reviewed: the reviewed instant and digest are required. */
  async applyCutover(sourceId: string, request: WriteGroupCutoverRequest): Promise<WriteGroupCutover> {
    if (!request.boundary || !request.review_digest) {
      throw new Error('A cutover is applied only with the reviewed boundary and review digest.');
    }
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(sourceId)}/cutover`,
      cutoverBody(request),
    );
    return parseWriteGroupAPIData(res.data, 'write-group cutover', parseWriteGroupCutover);
  },

  /** Reads one recorded cutover, so an operator whose answer was lost can see that it happened. */
  async getCutover(sourceId: string, cutoverId: string): Promise<WriteGroupCutover> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(sourceId)}/cutovers/${encodeURIComponent(cutoverId)}`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group cutover', parseWriteGroupCutover);
  },

  /** Lists the cutovers and rollbacks recorded for a source group. */
  async listCutovers(sourceId: string): Promise<WriteGroupCutover[]> {
    const res = await studioV2DatalinkApi.get<unknown>(`/studio-v2/workspace/write-groups/${encodeURIComponent(sourceId)}/cutovers`);
    return parseWriteGroupAPIData(res.data, 'write-group cutovers', parseWriteGroupCutoverList);
  },

  /** Durable delivery truth: only `sql_committed` rows were confirmed by the destination. */
  async delivery(id: string): Promise<WriteGroupDelivery> {
    const res = await studioV2DatalinkApi.get<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/delivery`,
    );
    return parseWriteGroupAPIData(res.data, 'write-group delivery', parseWriteGroupDeliveryData, 7);
  },

  async resolveDelivery(id: string, request: DeliveryDecision): Promise<DeliveryDecisionResult> {
    const res = await studioV2DatalinkApi.post<unknown>(`/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/delivery/resolve`, request);
    return parseWriteGroupAPIData(res.data, 'write-group delivery decision', (data) => {
      if (!data || typeof data !== 'object') return null;
      const result = data as Record<string, unknown>;
      if (result.decision_id !== request.decision_id || result.effect_key !== request.effect_key || typeof result.state !== 'string' || !['pending', 'operator_skipped', 'sending', 'retrying', 'blocked', 'quarantined', 'unknown', 'sql_committed'].includes(result.state) || typeof result.duplicate !== 'boolean') return null;
      if (!result.duplicate && result.state !== (request.resolution === 'retry' ? 'pending' : 'operator_skipped')) return null;
      return result as unknown as DeliveryDecisionResult;
    });
  },

  /** Describes the exact test row and its cleanup; the target is not written. */
  async testWritePreview(id: string): Promise<WriteGroupTestWritePreview> {
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/test-write-preview`,
    );
    return parseWriteGroupTestWritePreviewAPIData(res.data);
  },

  /**
   * Confirms a preview. The result is the recorded operation (202 while it is
   * still running); write verification and cleanup are separate fields.
   */
  async testWrite(id: string, request: WriteGroupTestWriteConfirmation): Promise<WriteGroupTestWriteOperation> {
    const res = await studioV2DatalinkApi.post<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}/test-write`,
      { token: request.token, operation_id: request.operation_id, action: request.action },
    );
    return parseWriteGroupAPIData(res.data, 'write-group test write', parseWriteGroupTestWriteOperation);
  },

  /** Reads a test-write operation from the shared workspace operation endpoint. */
  async testWriteOperation(operationId: string, locatorScope?: StudioV2LocatorReadScope): Promise<WriteGroupTestWriteOperation> {
    const selected = captureStudioLocatorReadScope(locatorScope);
    const url = `/studio-v2/workspace/database-operations/${encodeURIComponent(operationId)}`;
    const res = await (selected ? studioV2DatalinkApi.get<unknown>(url, { params: studioLocatorReadParams(selected) }) : studioV2DatalinkApi.get<unknown>(url))
      .catch(error => throwStudioLocatorReadError(error, selected));
    const operation = parseWriteGroupAPIData(res.data, 'write-group test write operation', value => {
      assertStudioLocatorWorkspaceBinding(value, selected); return parseWriteGroupTestWriteOperation(value);
    });
    if (selected && operation.operation_id !== operationId) throw new Error('Scoped settings read unavailable.');
    return operation;
  },

  async create(request: WriteGroupCreateRequest): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups',
      cloneCreateRequest(request),
    );
    return parseWriteGroupAPIData(res.data, 'write-group create', parseWriteGroupMutationData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
  },

  /** Creates or reuses the canonical group of one circuit profile; a resend returns the saved group. */
  async ensureCircuitProfile(request: WriteGroupCreateRequest): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.post<unknown>(
      '/studio-v2/workspace/write-groups/circuit-profiles/ensure',
      cloneCreateRequest(request),
    );
    return parseWriteGroupAPIData(res.data, 'write-group circuit profile', parseWriteGroupMutationData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
  },

  async update(id: string, request: WriteGroupUpdateRequest): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.put<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`,
      cloneUpdateRequest(request),
    );
    return parseWriteGroupAPIData(res.data, 'write-group update', parseWriteGroupMutationData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
  },

  async remove(id: string, request: WriteGroupDeleteRequest): Promise<WriteGroupMutationResponse> {
    const res = await studioV2DatalinkApi.delete<unknown>(
      `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`,
      { data: { ...request } },
    );
    return parseWriteGroupAPIData(res.data, 'write-group delete', parseWriteGroupMutationData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
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
    return parseWriteGroupAPIData(res.data, 'write-group disable', parseWriteGroupMutationData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
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
    return parseWriteGroupAPIData(res.data, 'write-group apply', parseWriteGroupMutationData, WRITE_GROUP_OBJECT_RESPONSE_MAX_DEPTH);
  },
};
