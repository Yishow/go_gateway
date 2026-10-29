import { studioV2DatalinkApi } from './studioV2Workspace';
import { parseWriteGroupAPIData } from '../utils/studioV2WriteGroupJson';
import { parseRecordingAPIData, parseRecordingSchemaOperation, parseRecordingSchemaPreviewToken } from '../utils/recordingPlanJson';
import { parseWriteGroupCutover } from '../utils/studioV2WriteGroupCutoverJson';
import { parseWriteGroupSuccessorReview, parseWriteGroupSuccessorPreparation, parseWriteGroupSuccessorOperationStatus,
  parseWriteGroupSuccessorActivationReview, parseWriteGroupSuccessorActivationReadiness, parseWriteGroupSuccessorFirstStartResult } from '../utils/studioV2WriteGroupSuccessorJson';
import type { WriteGroupSuccessorReviewRequest, WriteGroupSuccessorReservationRequest, WriteGroupSuccessorStageRequest,
  WriteGroupSuccessorSchemaConfirmation, WriteGroupSuccessorActivationRequest, WriteGroupSuccessorActivationReadinessRequest,
  WriteGroupSuccessorActivationConfirmation, WriteGroupSuccessorOperationStatusRequest } from '../types/studioV2WriteGroupSuccessor';
import { requireLifecycleMatch as match, safeLifecycleRequest, sameLifecycleJSON } from './studioV2LifecycleTransport';
import { digest } from '../utils/studioV2LifecycleJsonPrimitives';
export { isWriteGroupSuccessorUnconfirmed } from './studioV2LifecycleTransport';

const path = (id: string) => `/studio-v2/workspace/write-groups/${encodeURIComponent(id)}`;
async function post<T>(id: string, action: string, request: object, parse: (value: unknown) => T | null): Promise<T> {
  const response = await studioV2DatalinkApi.post<unknown>(`${path(id)}/successor`, { ...request, action });
  return parseWriteGroupAPIData(response.data, 'recording successor', parse, 16);
}
function rootMatches(a: WriteGroupSuccessorReviewRequest, b: WriteGroupSuccessorReviewRequest): boolean {
  return a.workspace_id === b.workspace_id && a.lineage_id === b.lineage_id && a.operation_id === b.operation_id &&
    a.expected_workspace_revision === b.expected_workspace_revision && a.expected_lineage_revision === b.expected_lineage_revision &&
    a.expected_predecessor_revision === b.expected_predecessor_revision && a.expected_predecessor_connector_revision === b.expected_predecessor_connector_revision &&
    a.expected_target_connector_revision === b.expected_target_connector_revision && a.group.destination.connector_id === b.group.destination.connector_id;
}
/** Stages always use one immutable root envelope; current stage CAS is separate. */
export const studioV2WriteGroupSuccessorAPI = {
  review(groupId: string, request: WriteGroupSuccessorReviewRequest) {
    return safeLifecycleRequest(async () => {
      const response = await studioV2DatalinkApi.post<unknown>(`${path(groupId)}/successor-review`, request);
      const result = parseWriteGroupAPIData(response.data, 'recording successor review', parseWriteGroupSuccessorReview, 16);
      return match(result, result.envelope.predecessor_group_id === groupId && rootMatches(result.envelope.request, request));
    });
  },
  reserve(groupId: string, request: WriteGroupSuccessorReservationRequest) {
    return safeLifecycleRequest(async () => {
      match(request, groupId === request.envelope.predecessor_group_id);
      const result = await post(groupId, 'reserve', request, parseWriteGroupSuccessorPreparation);
      const root = request.envelope.request;
      return match(result, result.predecessor_group_id === groupId && result.operation_id === root.operation_id &&
        result.lineage_id === root.lineage_id && result.group.workspace_id === root.workspace_id && result.generation === request.envelope.generation + 1);
    });
  },
  schemaPreview(groupId: string, request: WriteGroupSuccessorStageRequest) {
    return safeLifecycleRequest(async () => {
      match(request, groupId === request.envelope.predecessor_group_id);
      const response = await studioV2DatalinkApi.post<unknown>(`${path(groupId)}/successor`, { ...request, action: 'schema-preview' });
      const result = parseRecordingAPIData(response.data, 'successor schema preview', data => parseRecordingSchemaPreviewToken(data, { requireGroupLayout: true }));
      return match(result, result.workspace_id === request.envelope.request.workspace_id && result.workspace_revision === request.current.workspace_revision &&
        result.plan_id === request.current.group_id && result.plan_revision === request.current.group_revision &&
        result.connector_id === request.envelope.request.group.destination.connector_id && result.connector_revision === request.current.target_connector_revision);
    });
  },
  schemaConfirm(groupId: string, request: WriteGroupSuccessorSchemaConfirmation) {
    return safeLifecycleRequest(async () => {
      match(request, groupId === request.envelope.predecessor_group_id);
      const response = await studioV2DatalinkApi.post<unknown>(`${path(groupId)}/successor`, { ...request, action: 'schema-confirm' });
      const result = parseRecordingAPIData(response.data, 'successor schema confirmation', parseRecordingSchemaOperation);
      return match(result, result.operation_id === request.schema_operation_id &&
        (result.status !== 'succeeded' || digest(result.verified_digest) !== null));
    });
  },
  activationReview(groupId: string, request: WriteGroupSuccessorActivationRequest) {
    return safeLifecycleRequest(async () => {
      match(request, groupId === request.stage.envelope.predecessor_group_id);
      const result = await post(groupId, 'activation-review', request, parseWriteGroupSuccessorActivationReview);
      return match(result, result.envelope.request.stage.envelope.predecessor_group_id === groupId &&
        rootMatches(result.envelope.request.stage.envelope.request, request.stage.envelope.request) &&
        sameLifecycleJSON(result.envelope.request.stage.current, request.stage.current) &&
        sameLifecycleJSON(result.envelope.request.stage.envelope, request.stage.envelope));
    });
  },
  activationReadiness(groupId: string, request: WriteGroupSuccessorActivationReadinessRequest) {
    return safeLifecycleRequest(async () => {
      match(request, groupId === request.envelope.request.stage.envelope.predecessor_group_id);
      const result = await post(groupId, 'activation-readiness', request, parseWriteGroupSuccessorActivationReadiness);
      return match(result, result.review.review_digest === request.review_digest &&
        sameLifecycleJSON(result.review.envelope, request.envelope));
    });
  },
  confirmCutover(groupId: string, request: WriteGroupSuccessorActivationConfirmation) {
    return safeLifecycleRequest(async () => {
      match(request, groupId === request.envelope.request.stage.envelope.predecessor_group_id &&
        (request.envelope.version === 'successor-activation-review-v1' ? request.envelope.activation_kind === undefined : request.envelope.activation_kind === 'cutover'));
      const result = await post(groupId, 'activation-confirm', request, parseWriteGroupCutover);
      return match(result, result.workspace_id === request.envelope.request.stage.envelope.request.workspace_id && result.source_group_id === groupId &&
        result.boundary === request.envelope.request.boundary && result.review_digest === request.review_digest && result.kind === 'forward' &&
        result.ended.length === 1 && result.ended[0].group_id === groupId && result.ended[0].revision === request.envelope.source_applied_revision &&
        result.started.length === 1 && result.started[0].group_id === request.envelope.request.stage.current.group_id &&
        result.started[0].revision === request.envelope.request.stage.current.group_revision);
    });
  },
  confirmFirstStart(groupId: string, request: WriteGroupSuccessorActivationConfirmation) {
    return safeLifecycleRequest(async () => {
      match(request, groupId === request.envelope.request.stage.envelope.predecessor_group_id &&
        request.envelope.version === 'successor-activation-review-v2' && request.envelope.activation_kind === 'first_start');
      const result = await post(groupId, 'first-start-confirm', request, parseWriteGroupSuccessorFirstStartResult);
      const root = request.envelope.request.stage.envelope.request;
      return match(result, result.predecessor_group_id === groupId && result.root_operation_id === root.operation_id && result.lineage_id === root.lineage_id &&
        result.group_id === request.envelope.request.stage.current.group_id && result.boundary === request.envelope.request.boundary &&
        result.applied_revision === request.envelope.request.stage.current.group_revision &&
        result.generation === request.envelope.request.stage.envelope.generation + 1 &&
        result.lineage_revision === request.envelope.request.stage.current.lineage_revision + 1);
    });
  },
  operationStatus(groupId: string, request: WriteGroupSuccessorOperationStatusRequest) {
    return safeLifecycleRequest(async () => {
      const exact = { workspace_id: request.workspace_id, lineage_id: request.lineage_id, operation_id: request.operation_id };
      const result = await post(groupId, 'operation-status', exact, parseWriteGroupSuccessorOperationStatus);
      return match(result, result.workspace_id === request.workspace_id && result.predecessor_group_id === groupId &&
        result.lineage_id === request.lineage_id && result.operation_id === request.operation_id);
    });
  },
};
