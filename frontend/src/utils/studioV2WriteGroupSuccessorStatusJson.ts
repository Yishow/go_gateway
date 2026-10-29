import type { WriteGroupSuccessorOperationStatus, WriteGroupSuccessorRecordedStage } from '../types/studioV2WriteGroupSuccessor';
import { parseWriteGroupSuccessorPreparation } from './studioV2WriteGroupSuccessorJson';
import { parseWriteGroupSuccessorFirstStartResult } from './studioV2WriteGroupSuccessorActivationJson';
import { parseWriteGroupCutover } from './studioV2WriteGroupCutoverJson';
import { isRecord, text, texts, list, optionalCode } from './studioV2LifecycleJsonPrimitives';

function stage<T>(value: unknown, parse: (item: unknown) => T | null): WriteGroupSuccessorRecordedStage<T> | null {
  if (!isRecord(value)) return null;
  if (value.status === 'not_submitted') return value.result === undefined ? { status: 'not_submitted' } : null;
  const result = value.status === 'succeeded' ? parse(value.result) : null;
  return result ? { status: 'succeeded', result } : null;
}
function schema(value: unknown): WriteGroupSuccessorOperationStatus['schema_operations'][number] | null {
  if (!isRecord(value)) return null;
  const operation = text(value.operation_id), status = value.status, reason = optionalCode(value.reason_code), next = optionalCode(value.next_action);
  if (!operation || reason === null || next === null || (status !== 'preview_not_submitted' && status !== 'pending' && status !== 'running' &&
    status !== 'succeeded' && status !== 'partial' && status !== 'failed' && status !== 'unknown')) return null;
  return { operation_id: operation, status, ...(reason ? { reason_code: reason } : {}), ...(next ? { next_action: next } : {}) };
}
export function parseWriteGroupSuccessorOperationStatus(value: unknown): WriteGroupSuccessorOperationStatus | null {
  if (!isRecord(value) || !isRecord(value.reserve) || value.reserve.status !== 'succeeded') return null;
  const fields = texts(value, ['workspace_id', 'lineage_id', 'predecessor_group_id', 'operation_id']), reserve = parseWriteGroupSuccessorPreparation(value.reserve.result);
  const operations = list(value.schema_operations, schema), cutover = stage(value.cutover, parseWriteGroupCutover), first = stage(value.first_start, parseWriteGroupSuccessorFirstStartResult);
  if (!fields || !reserve || !operations || !cutover || !first || reserve.group.workspace_id !== fields.workspace_id ||
    reserve.lineage_id !== fields.lineage_id || reserve.predecessor_group_id !== fields.predecessor_group_id || reserve.operation_id !== fields.operation_id ||
    new Set(operations.map(op => op.operation_id)).size !== operations.length || (cutover.status === 'succeeded' && first.status === 'succeeded')) return null;
  if (cutover.status === 'succeeded' && (cutover.result.workspace_id !== fields.workspace_id || cutover.result.source_group_id !== fields.predecessor_group_id ||
    !cutover.result.started.some(ref => ref.group_id === reserve.group.id))) return null;
  if (first.status === 'succeeded' && (first.result.root_operation_id !== fields.operation_id || first.result.lineage_id !== fields.lineage_id ||
    first.result.predecessor_group_id !== fields.predecessor_group_id || first.result.group_id !== reserve.group.id || first.result.generation !== reserve.generation)) return null;
  return { ...fields, reserve: { status: 'succeeded', result: reserve }, schema_operations: operations, cutover, first_start: first };
}
