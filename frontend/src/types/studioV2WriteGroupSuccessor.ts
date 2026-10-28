import type { WriteGroup } from './studioV2WriteGroup';
import type { WriteGroupReadiness } from './studioV2WriteGroup';
import type { WriteGroupCutover } from './studioV2WriteGroupCutover';
import type { WriteGroupLifecycleHistory } from './studioV2WriteGroupLifecycle';
import type { SchemaOperationStatus } from './recordingPlan';

export interface WriteGroupSuccessorReviewRequest {
  workspace_id: string;
  expected_workspace_revision: string;
  lineage_id: string;
  expected_lineage_revision: number;
  expected_predecessor_revision: string;
  expected_predecessor_connector_revision: string;
  expected_target_connector_revision: string;
  operation_id: string;
  boundary: string;
  group: WriteGroup;
}
export interface WriteGroupSuccessorEnvelope {
  version: 'successor-review-v1';
  predecessor_group_id: string;
  generation: number;
  source_applied_revision: string;
  source_definition_digest: string;
  request: WriteGroupSuccessorReviewRequest;
}
export interface WriteGroupSuccessorBlocker { reason_code: string; next_action: string }
export interface WriteGroupSuccessorReview {
  envelope: WriteGroupSuccessorEnvelope;
  review_digest: string;
  history: WriteGroupLifecycleHistory;
  blockers: WriteGroupSuccessorBlocker[];
  required_checks: WriteGroupSuccessorBlocker[];
}
export interface WriteGroupSuccessorCurrentCAS {
  workspace_revision: string;
  lineage_revision: number;
  predecessor_revision: string;
  predecessor_connector_revision: string;
  group_id: string;
  group_revision: string;
  target_connector_revision: string;
}
export interface WriteGroupSuccessorReservationRequest {
  envelope: WriteGroupSuccessorEnvelope;
  review_digest: string;
}
export interface WriteGroupSuccessorStageRequest extends WriteGroupSuccessorReservationRequest {
  current: WriteGroupSuccessorCurrentCAS;
}
export interface WriteGroupSuccessorSchemaConfirmation extends WriteGroupSuccessorStageRequest {
  token: string;
  digest: string;
  schema_operation_id: string;
}
export interface WriteGroupSuccessorPreparation {
  workspace_revision: string;
  lineage_id: string;
  lineage_revision: number;
  generation: number;
  predecessor_group_id: string;
  operation_id: string;
  intent_digest: string;
  group: WriteGroup;
}
export interface WriteGroupSuccessorActivationRequest {
  stage: WriteGroupSuccessorStageRequest;
  boundary: string;
}
interface WriteGroupSuccessorActivationEvidence {
  request: WriteGroupSuccessorActivationRequest;
  source_applied_revision: string;
  source_definition_digest: string;
  target_definition_digest: string;
  counter_proof_digest: string;
  schema_revision: string;
  schema_digest: string;
}
/** Legacy v1 can only mean cutover; v2 always binds an explicit kind. */
export type WriteGroupSuccessorActivationEnvelope = WriteGroupSuccessorActivationEvidence & (
  | { version: 'successor-activation-review-v1'; activation_kind?: never }
  | { version: 'successor-activation-review-v2'; activation_kind: 'cutover' | 'first_start' }
);
export interface WriteGroupSuccessorActivationReview {
  envelope: WriteGroupSuccessorActivationEnvelope;
  review_digest: string;
  history: WriteGroupLifecycleHistory;
  blockers: WriteGroupSuccessorBlocker[];
}
export interface WriteGroupSuccessorActivationReadinessRequest {
  envelope: WriteGroupSuccessorActivationEnvelope;
  review_digest: string;
}
export interface WriteGroupSuccessorActivationConfirmation extends WriteGroupSuccessorActivationReadinessRequest {
  readiness_digest: string;
}
export interface WriteGroupSuccessorActivationReadiness {
  review: WriteGroupSuccessorActivationReview;
  readiness: WriteGroupReadiness;
  readiness_digest: string;
}
export interface WriteGroupSuccessorFirstStartResult {
  operation_id: string;
  root_operation_id: string;
  workspace_revision: string;
  lineage_id: string;
  lineage_revision: number;
  predecessor_group_id: string;
  group_id: string;
  applied_revision: string;
  boundary: string;
  generation: number;
  status: 'succeeded';
}
export interface WriteGroupSuccessorOperationStatusRequest {
  workspace_id: string;
  lineage_id: string;
  operation_id: string;
}
export type WriteGroupSuccessorRecordedStage<T> =
  | { status: 'not_submitted'; result?: never }
  | { status: 'succeeded'; result: T };
export interface WriteGroupSuccessorOperationStatus {
  workspace_id: string;
  lineage_id: string;
  predecessor_group_id: string;
  operation_id: string;
  reserve: { status: 'succeeded'; result: WriteGroupSuccessorPreparation };
  schema_operations: Array<{
    operation_id: string;
    status: SchemaOperationStatus | 'preview_not_submitted';
    reason_code?: string;
    next_action?: string;
  }>;
  cutover: WriteGroupSuccessorRecordedStage<WriteGroupCutover>;
  first_start: WriteGroupSuccessorRecordedStage<WriteGroupSuccessorFirstStartResult>;
}
