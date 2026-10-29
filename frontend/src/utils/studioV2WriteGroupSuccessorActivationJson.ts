import type { WriteGroupSuccessorActivationEnvelope, WriteGroupSuccessorActivationReview, WriteGroupSuccessorActivationReadiness, WriteGroupSuccessorActivationConfirmation, WriteGroupSuccessorFirstStartResult } from '../types/studioV2WriteGroupSuccessor';
import { parseWriteGroupSuccessorStageRequest, parseWriteGroupSuccessorBlocker } from './studioV2WriteGroupSuccessorJson';
import { parseWriteGroupLifecycleHistory } from './studioV2WriteGroupLifecycleJson';
import { parseWriteGroupReadinessData } from './studioV2WriteGroupJson';
import { isRecord, texts, text, digest, whole, instant, list } from './studioV2LifecycleJsonPrimitives';

export function parseWriteGroupSuccessorActivationEnvelope(value: unknown): WriteGroupSuccessorActivationEnvelope | null {
  if (!isRecord(value) || !isRecord(value.request)) return null;
  const version = value.version, kind = value.activation_kind;
  if (version === 'successor-activation-review-v1' ? kind !== undefined :
    version !== 'successor-activation-review-v2' || (kind !== 'cutover' && kind !== 'first_start')) return null;
  const stage = parseWriteGroupSuccessorStageRequest(value.request.stage), boundary = instant(value.request.boundary);
  const sourceApplied = text(value.source_applied_revision, true), source = digest(value.source_definition_digest), target = digest(value.target_definition_digest);
  const counter = value.counter_proof_digest === '' ? '' : digest(value.counter_proof_digest);
  const schema = text(value.schema_revision, true), schemaHash = value.schema_digest === '' ? '' : digest(value.schema_digest);
  if (!stage || !boundary || sourceApplied === null || !source || !target || counter === null || schema === null || schemaHash === null ||
    (schema === '') !== (schemaHash === '') || sourceApplied !== stage.envelope.source_applied_revision || (kind === 'first_start' ? sourceApplied !== '' : sourceApplied === '')) return null;
  const evidence = { request: { stage, boundary }, source_applied_revision: sourceApplied, source_definition_digest: source,
    target_definition_digest: target, counter_proof_digest: counter, schema_revision: schema, schema_digest: schemaHash };
  return version === 'successor-activation-review-v1' ? { ...evidence, version } : { ...evidence, version: 'successor-activation-review-v2', activation_kind: kind as 'cutover' | 'first_start' };
}
export function parseWriteGroupSuccessorActivationReview(value: unknown): WriteGroupSuccessorActivationReview | null {
  if (!isRecord(value)) return null;
  const envelope = parseWriteGroupSuccessorActivationEnvelope(value.envelope), hash = digest(value.review_digest);
  const history = parseWriteGroupLifecycleHistory(value.history), blockers = list(value.blockers, parseWriteGroupSuccessorBlocker);
  if (!envelope || !hash || !history || !blockers || history.lineage_id !== envelope.request.stage.envelope.request.lineage_id ||
    (blockers.length === 0 && (!envelope.schema_revision || !envelope.schema_digest || !envelope.counter_proof_digest))) return null;
  return { envelope, review_digest: hash, history, blockers };
}
export function parseWriteGroupSuccessorActivationReadiness(value: unknown): WriteGroupSuccessorActivationReadiness | null {
  if (!isRecord(value)) return null;
  const review = parseWriteGroupSuccessorActivationReview(value.review), readiness = parseWriteGroupReadinessData(value.readiness), hash = value.readiness_digest === '' ? '' : digest(value.readiness_digest);
  if (!review || !readiness || hash === null || readiness.ready !== Boolean(hash) || readiness.workspace_id !== review.envelope.request.stage.envelope.request.workspace_id ||
    readiness.group_id !== review.envelope.request.stage.current.group_id || readiness.group_revision !== review.envelope.request.stage.current.group_revision ||
    readiness.workspace_revision !== review.envelope.request.stage.current.workspace_revision) return null;
  return { review, readiness, readiness_digest: hash };
}
export function parseWriteGroupSuccessorActivationConfirmation(value: unknown): WriteGroupSuccessorActivationConfirmation | null {
  if (!isRecord(value)) return null;
  const envelope = parseWriteGroupSuccessorActivationEnvelope(value.envelope), review = digest(value.review_digest), ready = digest(value.readiness_digest);
  return envelope && review && ready ? { envelope, review_digest: review, readiness_digest: ready } : null;
}
export function parseWriteGroupSuccessorFirstStartResult(value: unknown): WriteGroupSuccessorFirstStartResult | null {
  if (!isRecord(value) || value.status !== 'succeeded') return null;
  const fields = texts(value, ['operation_id', 'root_operation_id', 'workspace_revision', 'lineage_id', 'predecessor_group_id', 'group_id', 'applied_revision']);
  const lineage = whole(value.lineage_revision, 1), generation = whole(value.generation, 2), boundary = instant(value.boundary);
  return fields && lineage !== null && generation !== null && boundary && fields.group_id !== fields.predecessor_group_id
    ? { ...fields, lineage_revision: lineage, generation, boundary, status: 'succeeded' } : null;
}
