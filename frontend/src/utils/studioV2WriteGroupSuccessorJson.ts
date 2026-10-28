import type { WriteGroupSuccessorReviewRequest, WriteGroupSuccessorEnvelope, WriteGroupSuccessorReview, WriteGroupSuccessorCurrentCAS, WriteGroupSuccessorStageRequest, WriteGroupSuccessorPreparation, WriteGroupSuccessorBlocker, WriteGroupSuccessorSchemaConfirmation } from '../types/studioV2WriteGroupSuccessor';
import { parseWriteGroupSuccessorCandidate, parseStudioV2WriteGroup } from './studioV2WriteGroupDefinitionJson';
import { parseWriteGroupLifecycleHistory } from './studioV2WriteGroupLifecycleJson';
import { isRecord, texts, text, digest, whole, instant, list, code } from './studioV2LifecycleJsonPrimitives';
export { parseWriteGroupSuccessorActivationEnvelope, parseWriteGroupSuccessorActivationReview, parseWriteGroupSuccessorActivationReadiness, parseWriteGroupSuccessorActivationConfirmation, parseWriteGroupSuccessorFirstStartResult } from './studioV2WriteGroupSuccessorActivationJson';
export { parseWriteGroupSuccessorOperationStatus } from './studioV2WriteGroupSuccessorStatusJson';

export function parseWriteGroupSuccessorReviewRequest(value: unknown): WriteGroupSuccessorReviewRequest | null {
  if (!isRecord(value)) return null;
  const fields = texts(value, ['workspace_id', 'expected_workspace_revision', 'lineage_id', 'expected_predecessor_revision',
    'expected_predecessor_connector_revision', 'expected_target_connector_revision', 'operation_id']);
  const lineage = whole(value.expected_lineage_revision, 1), boundary = instant(value.boundary), group = parseWriteGroupSuccessorCandidate(value.group);
  if (!fields || lineage === null || !boundary || !group || group.workspace_id !== fields.workspace_id ||
    group.destination.connector_revision !== fields.expected_target_connector_revision) return null;
  return { ...fields, expected_lineage_revision: lineage, boundary, group };
}
export function parseWriteGroupSuccessorEnvelope(value: unknown): WriteGroupSuccessorEnvelope | null {
  if (!isRecord(value) || value.version !== 'successor-review-v1') return null;
  const predecessor = text(value.predecessor_group_id), generation = whole(value.generation, 1), applied = text(value.source_applied_revision, true);
  const source = digest(value.source_definition_digest), request = parseWriteGroupSuccessorReviewRequest(value.request);
  return predecessor && generation !== null && applied !== null && source && request
    ? { version: value.version, predecessor_group_id: predecessor, generation, source_applied_revision: applied, source_definition_digest: source, request } : null;
}
export function parseWriteGroupSuccessorBlocker(value: unknown): WriteGroupSuccessorBlocker | null {
  if (!isRecord(value)) return null;
  const reason = code(value.reason_code), next = code(value.next_action);
  return reason && next ? { reason_code: reason, next_action: next } : null;
}
export function parseWriteGroupSuccessorReview(value: unknown): WriteGroupSuccessorReview | null {
  if (!isRecord(value)) return null;
  const envelope = parseWriteGroupSuccessorEnvelope(value.envelope), hash = digest(value.review_digest), history = parseWriteGroupLifecycleHistory(value.history);
  const blockers = list(value.blockers, parseWriteGroupSuccessorBlocker), checks = list(value.required_checks, parseWriteGroupSuccessorBlocker);
  if (!envelope || !hash || !history || !blockers || !checks || history.lineage_id !== envelope.request.lineage_id ||
    !history.generations.some(g => g.group_id === envelope.predecessor_group_id && g.generation === envelope.generation)) return null;
  return { envelope, review_digest: hash, history, blockers, required_checks: checks };
}
export function parseWriteGroupSuccessorCurrentCAS(value: unknown): WriteGroupSuccessorCurrentCAS | null {
  if (!isRecord(value)) return null;
  const fields = texts(value, ['workspace_revision', 'predecessor_revision', 'predecessor_connector_revision', 'group_id', 'group_revision', 'target_connector_revision']);
  const lineage = whole(value.lineage_revision, 1);
  return fields && lineage !== null ? { ...fields, lineage_revision: lineage } : null;
}
export function parseWriteGroupSuccessorStageRequest(value: unknown): WriteGroupSuccessorStageRequest | null {
  if (!isRecord(value)) return null;
  const envelope = parseWriteGroupSuccessorEnvelope(value.envelope), hash = digest(value.review_digest), current = parseWriteGroupSuccessorCurrentCAS(value.current);
  if (!envelope || !hash || !current || current.group_id === envelope.predecessor_group_id) return null;
  return { envelope, review_digest: hash, current };
}
export function parseWriteGroupSuccessorSchemaConfirmation(value: unknown): WriteGroupSuccessorSchemaConfirmation | null {
  if (!isRecord(value)) return null;
  const stage = parseWriteGroupSuccessorStageRequest(value), token = text(value.token), hash = digest(value.digest), operation = text(value.schema_operation_id);
  return stage && token && hash && operation ? { ...stage, token, digest: hash, schema_operation_id: operation } : null;
}
export function parseWriteGroupSuccessorPreparation(value: unknown): WriteGroupSuccessorPreparation | null {
  if (!isRecord(value)) return null;
  const fields = texts(value, ['workspace_revision', 'lineage_id', 'predecessor_group_id', 'operation_id']);
  const lineage = whole(value.lineage_revision, 1), generation = whole(value.generation, 2), intent = digest(value.intent_digest), group = parseStudioV2WriteGroup(value.group);
  if (!fields || lineage === null || generation === null || !intent || !group || group.status !== 'draft' || group.applied_revision !== '' ||
    group.id === fields.predecessor_group_id || group.destination.schema_revision !== undefined || group.destination.schema_digest !== undefined) return null;
  return { ...fields, lineage_revision: lineage, generation, intent_digest: intent, group };
}
