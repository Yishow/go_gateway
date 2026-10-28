import type { WriteGroup } from '../../../../../../../types/studioV2WriteGroup';
import type { SavedRecordingActionIntent, WriteGroupLifecycle } from '../../../../../../../types/studioV2WriteGroupLifecycle';
import type * as S from '../../../../../../../types/studioV2WriteGroupSuccessor';
import { WriteGroupResponseError } from '../../../../../../../utils/studioV2WriteGroupJson';
import { RecordingPlanResponseError } from '../../../../../../../utils/recordingPlanJson';
import { parsePreciseInstant, comparePreciseInstants } from '../../../../../../../utils/preciseInstant';
import { normalizeTypedEnvelope } from '../../../../../../../utils/safeJson';
import type { RecordingChangeRecovery } from './recordingChangeTypes';

export function activationAction(envelope: S.WriteGroupSuccessorActivationEnvelope): 'confirmCutover' | 'confirmFirstStart' {
  if (envelope.version === 'successor-activation-review-v1' && envelope.activation_kind === undefined) return 'confirmCutover';
  if (envelope.version === 'successor-activation-review-v2') {
    if (envelope.activation_kind === 'cutover') return 'confirmCutover';
    if (envelope.activation_kind === 'first_start') return 'confirmFirstStart';
  }
  throw new Error('Invalid activation kind');
}
export function matchesIntent(lifecycle: WriteGroupLifecycle, intent: SavedRecordingActionIntent): boolean {
  return lifecycle.workspace_id === intent.workspaceId && lifecycle.group_id === intent.groupId &&
    lifecycle.workspace_revision === intent.workspaceRevision && lifecycle.group_revision === intent.groupRevision;
}
export function operationScope(record: RecordingChangeRecovery): S.WriteGroupSuccessorOperationStatusRequest {
  return { workspace_id: record.envelope.request.workspace_id, lineage_id: record.envelope.request.lineage_id, operation_id: record.envelope.request.operation_id };
}
export function currentStage(record: RecordingChangeRecovery, source: WriteGroupLifecycle, target: WriteGroupLifecycle): S.WriteGroupSuccessorStageRequest {
  const root = record.envelope;
  if (source.workspace_id !== root.request.workspace_id || source.group_id !== root.predecessor_group_id ||
    target.workspace_id !== source.workspace_id || target.workspace_revision !== source.workspace_revision ||
    source.history.lineage_id !== root.request.lineage_id || target.history.lineage_id !== source.history.lineage_id ||
    target.group_id !== record.preparation?.group.id || source.history.reserved_group_id !== target.group_id ||
    !source.history.lineage_revision || target.history.lineage_revision !== source.history.lineage_revision) throw new Error('Stale scope');
  return { envelope: root, review_digest: record.reviewDigest, current: {
    workspace_revision: source.workspace_revision, lineage_revision: source.history.lineage_revision,
    predecessor_revision: source.group_revision, predecessor_connector_revision: source.connector_revision,
    group_id: target.group_id, group_revision: target.group_revision, target_connector_revision: target.connector_revision,
  } };
}
export function reviewRequest(lifecycle: WriteGroupLifecycle, target: WriteGroup, operationId: string): S.WriteGroupSuccessorReviewRequest {
  if (!lifecycle.history.lineage_id || !lifecycle.history.lineage_revision) throw new Error('Missing lineage');
  return { workspace_id: lifecycle.workspace_id, expected_workspace_revision: lifecycle.workspace_revision,
    lineage_id: lifecycle.history.lineage_id, expected_lineage_revision: lifecycle.history.lineage_revision,
    expected_predecessor_revision: lifecycle.group_revision, expected_predecessor_connector_revision: lifecycle.connector_revision,
    expected_target_connector_revision: target.destination.connector_revision, operation_id: operationId,
    // The initial review has no activation authority. The operator chooses a fresh B only at step 3.
    boundary: new Date(Date.now() + 60 * 60 * 1000).toISOString(), group: structuredClone(target) };
}
/** Only a structured application 4xx proves rejection. Status-only/proxy/5xx/malformed remain unknown. */
export function explicitRejection(cause: unknown): boolean {
  if (cause instanceof WriteGroupResponseError || cause instanceof RecordingPlanResponseError) return cause.outcome === 'failed';
  if (!cause || typeof cause !== 'object') return false;
  const error = cause as { response?: { status?: number }; status?: number };
  const status = error.response?.status ?? error.status;
  return status !== undefined && status >= 400 && status < 500 && status !== 408 && Boolean(normalizeTypedEnvelope(cause).code);
}
export function isFutureBoundary(value: string): boolean { return Number.isFinite(Date.parse(value)) && Date.parse(value) > Date.now(); }

export function sameBoundary(left: string, right: string): boolean {
  const a = parsePreciseInstant(left), b = parsePreciseInstant(right);
  return !!a && !!b && comparePreciseInstants(a, b) === 0;
}
/** Match the existing profile authoring/backend contract without rebuilding any source semantics. */
export function changeTargetLayout(target: WriteGroup, layout: WriteGroup['layout_version']): WriteGroup {
  if ((target.storage_mode && target.storage_mode !== 'plain') || !target.recording_profile || (layout !== 'managed_long_v1' && layout !== 'managed_wide_profile_v1')) return target;
  const next = structuredClone(target);
  next.layout_version = layout; next.destination.table_name = '';
  if (layout === 'managed_long_v1') next.storage_mode = 'plain';
  else { delete next.storage_mode; delete next.write_policy.recreate_missing; }
  return next;
}

export function matchesRootStatus(record: RecordingChangeRecovery, status: S.WriteGroupSuccessorOperationStatus): boolean {
  const root = operationScope(record), result = status.reserve.result;
  return status.workspace_id === root.workspace_id && status.lineage_id === root.lineage_id && status.operation_id === root.operation_id &&
    status.predecessor_group_id === record.envelope.predecessor_group_id && result.operation_id === root.operation_id &&
    result.lineage_id === root.lineage_id && result.group.workspace_id === root.workspace_id && result.predecessor_group_id === status.predecessor_group_id &&
    result.generation === record.envelope.generation + 1 && (!record.preparation ||
      (record.preparation.group.id === result.group.id && record.preparation.intent_digest === result.intent_digest));
}
export function canRepreviewAfterStatus(record: RecordingChangeRecovery, status: S.WriteGroupSuccessorOperationStatus): boolean {
  return !record.unreadable && record.submitted.action === 'schemaPreview' && !record.preview && !record.schemaConfirmationSubmitted &&
    matchesRootStatus(record, status) && status.cutover.status === 'not_submitted' && status.first_start.status === 'not_submitted' &&
    status.schema_operations.every(item => item.status === 'preview_not_submitted');
}
