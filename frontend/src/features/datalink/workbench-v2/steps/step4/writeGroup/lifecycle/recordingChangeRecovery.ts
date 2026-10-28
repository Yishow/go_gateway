import { parseWriteGroupSuccessorEnvelope, parseWriteGroupSuccessorStageRequest, parseWriteGroupSuccessorSchemaConfirmation, parseWriteGroupSuccessorActivationConfirmation, parseWriteGroupSuccessorPreparation, parseWriteGroupSuccessorFirstStartResult, parseWriteGroupSuccessorOperationStatus } from '../../../../../../../utils/studioV2WriteGroupSuccessorJson';
import { parseWriteGroupCutover } from '../../../../../../../utils/studioV2WriteGroupCutoverJson';
import { parseRecordingSchemaPreviewToken } from '../../../../../../../utils/recordingPlanJson';
import type { RecordingChangeRecovery, UnreadableRecordingChangeScope } from './recordingChangeTypes';

const prefix = 'wbv2.recording-change.';
const memory = new Map<string, RecordingChangeRecovery>();
export const recordingChangeKey = (record: RecordingChangeRecovery) => prefix + JSON.stringify([
  record.envelope.request.workspace_id, record.envelope.predecessor_group_id,
  record.envelope.request.lineage_id, record.envelope.request.operation_id,
]);
const forbidden = /^(?:password|credential|credentials|secret|dsn|private_?detail|raw_?error|statements)$/i;
function safe(value: unknown): boolean {
  if (Array.isArray(value)) return value.every(safe);
  return !value || typeof value !== 'object' || Object.entries(value).every(([key, child]) => !forbidden.test(key) && safe(child));
}
function canonical(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(canonical).join(',')}]`;
  if (value && typeof value === 'object') return `{${Object.entries(value).sort(([a], [b]) => a.localeCompare(b)).map(([key, item]) => `${JSON.stringify(key)}:${canonical(item)}`).join(',')}}`;
  return JSON.stringify(value);
}
function valid(value: unknown): value is RecordingChangeRecovery {
  if (!value || typeof value !== 'object' || !safe(value)) return false;
  const record = value as RecordingChangeRecovery;
  const envelope = parseWriteGroupSuccessorEnvelope(record.envelope);
  const statuses = new Set(['preview_not_submitted', 'pending', 'running', 'succeeded', 'failed', 'partial', 'unknown']);
  const opaque = (id: unknown): id is string => typeof id === 'string' && /^[A-Za-z0-9][A-Za-z0-9:_-]{0,255}$/.test(id);
  if (record.attempt !== undefined && (!Number.isSafeInteger(record.attempt) || record.attempt < 0)) return false;
  if ([record.schemaConfirmationSubmitted, record.stageRejected, record.unreadable].some(flag => flag !== undefined && typeof flag !== 'boolean')) return false;
  if (record.priorPreviewIds !== undefined && (!Array.isArray(record.priorPreviewIds) || record.priorPreviewIds.length > 256 || !record.priorPreviewIds.every(opaque))) return false;
  if (record.priorPreviewStates !== undefined && (!record.priorPreviewStates || typeof record.priorPreviewStates !== 'object' || Array.isArray(record.priorPreviewStates) ||
    Object.entries(record.priorPreviewStates).some(([id, status]) => !opaque(id) || !statuses.has(status)))) return false;
  if (record.previewStage !== undefined && !parseWriteGroupSuccessorStageRequest(record.previewStage)) return false;
  if (record.terminalStatus !== undefined && !parseWriteGroupSuccessorOperationStatus(record.terminalStatus)) return false;
  if (record.version !== 1 || !envelope || !/^[a-f0-9]{64}$/.test(record.reviewDigest) ||
    !['choose', 'impact', 'reserved', 'schema-plan', 'schema-confirmed', 'activation-ready', 'scheduled'].includes(record.phase) ||
    !['pending', 'recorded', 'unknown', 'rejected'].includes(record.outcome)) return false;
  const submitted = record.submitted;
  if (!submitted || !submitted.request) return false;
  let root;
  if (submitted.action === 'reserve') {
    root = parseWriteGroupSuccessorEnvelope(submitted.request.envelope);
    if (submitted.request.review_digest !== record.reviewDigest) return false;
  } else if (submitted.action === 'schemaPreview' || submitted.action === 'schemaConfirm') {
    const stage = submitted.action === 'schemaConfirm' ? parseWriteGroupSuccessorSchemaConfirmation(submitted.request) : parseWriteGroupSuccessorStageRequest(submitted.request);
    if (!stage || stage.review_digest !== record.reviewDigest) return false;
    root = stage.envelope;
  } else if (submitted.action === 'confirmCutover' || submitted.action === 'confirmFirstStart') {
    const confirmation = parseWriteGroupSuccessorActivationConfirmation(submitted.request);
    if (!confirmation) return false;
    root = confirmation.envelope.request.stage.envelope;
  } else return false;
  if (!root || JSON.stringify(root) !== JSON.stringify(envelope)) return false;
  if (record.preparation) {
    const prepared = parseWriteGroupSuccessorPreparation(record.preparation);
    if (!prepared || prepared.operation_id !== envelope.request.operation_id || prepared.lineage_id !== envelope.request.lineage_id ||
      prepared.predecessor_group_id !== envelope.predecessor_group_id || prepared.group.workspace_id !== envelope.request.workspace_id) return false;
  }
  if (record.preview) {
    const preview = parseRecordingSchemaPreviewToken({ ...record.preview, statements: [] }, { requireGroupLayout: true });
    const stage = parseWriteGroupSuccessorStageRequest(record.previewStage);
    if (!preview || !stage || preview.plan_id !== record.preparation?.group.id || preview.workspace_id !== envelope.request.workspace_id ||
      JSON.stringify(stage.envelope) !== JSON.stringify(envelope)) return false;
  }
  if (record.phase === 'scheduled') {
    if (record.outcome !== 'recorded' || !record.result) return false;
    if (record.terminalStatus) {
      const status = parseWriteGroupSuccessorOperationStatus(record.terminalStatus);
      if (!status || status.workspace_id !== envelope.request.workspace_id || status.lineage_id !== envelope.request.lineage_id ||
        status.operation_id !== envelope.request.operation_id || status.predecessor_group_id !== envelope.predecessor_group_id) return false;
      const result = status.first_start.status === 'succeeded' ? status.first_start.result : status.cutover.status === 'succeeded' ? status.cutover.result : undefined;
      return !!result && canonical(result) === canonical(record.result);
    }
    if (submitted.action !== 'confirmCutover' && submitted.action !== 'confirmFirstStart') return false;
    const confirmation = submitted.request;
    if (submitted.action === 'confirmFirstStart') {
      const result = parseWriteGroupSuccessorFirstStartResult(record.result);
      if (!result || result.root_operation_id !== envelope.request.operation_id || result.lineage_id !== envelope.request.lineage_id ||
        result.predecessor_group_id !== envelope.predecessor_group_id || result.group_id !== confirmation.envelope.request.stage.current.group_id ||
        result.boundary !== confirmation.envelope.request.boundary || result.applied_revision !== confirmation.envelope.request.stage.current.group_revision ||
        result.generation !== envelope.generation + 1 || result.lineage_revision !== confirmation.envelope.request.stage.current.lineage_revision + 1) return false;
    } else {
      const result = parseWriteGroupCutover(record.result);
      if (!result || result.source_group_id !== envelope.predecessor_group_id || result.workspace_id !== envelope.request.workspace_id ||
        result.boundary !== confirmation.envelope.request.boundary || result.kind !== 'forward' || result.review_digest !== confirmation.review_digest ||
        result.ended.length !== 1 || result.ended[0].group_id !== envelope.predecessor_group_id || result.ended[0].revision !== confirmation.envelope.source_applied_revision ||
        result.started.length !== 1 || result.started[0].group_id !== confirmation.envelope.request.stage.current.group_id || result.started[0].revision !== confirmation.envelope.request.stage.current.group_revision) return false;
    }
  } else if (record.result !== undefined) return false;
  return true;
}
/** A failure to persist the exact request blocks a new mutation, rather than silently losing reload recovery. */
export function saveRecordingChange(record: RecordingChangeRecovery): boolean {
  if (!valid(record)) return false;
  const copy = structuredClone(record);
  const key = recordingChangeKey(copy);
  const existing = memory.get(key);
  if ((existing?.attempt ?? 0) > (copy.attempt ?? 0) || (existing?.phase === 'scheduled' && copy.phase !== 'scheduled')) return true;
  memory.set(key, copy);
  try { sessionStorage.setItem(key, JSON.stringify(copy)); return true; } catch { return false; }
}
export function loadRecordingChanges(workspaceId: string): RecordingChangeRecovery[] {
  const found = new Map<string, RecordingChangeRecovery>();
  try {
    for (let index = 0; index < sessionStorage.length; index++) {
      const key = sessionStorage.key(index);
      if (!key?.startsWith(prefix)) continue;
      try {
        const item: unknown = JSON.parse(sessionStorage.getItem(key) ?? 'null');
        if (valid(item) && key === recordingChangeKey(item) && item.envelope.request.workspace_id === workspaceId) found.set(key, item);
        else if (item && typeof item === 'object') {
          // Retain only verified root identity from an unreadable entry. It authorizes status, never retry or a new preview.
          const raw = item as RecordingChangeRecovery;
          const envelope = parseWriteGroupSuccessorEnvelope(raw.envelope);
          if (envelope && /^[a-f0-9]{64}$/.test(raw.reviewDigest)) {
            const fence: RecordingChangeRecovery = { version: 1, envelope, reviewDigest: raw.reviewDigest,
              submitted: { action: 'reserve', request: { envelope, review_digest: raw.reviewDigest } },
              phase: 'impact', outcome: 'unknown', unreadable: true };
            if (key === recordingChangeKey(fence) && envelope.request.workspace_id === workspaceId) found.set(key, fence);
          }
        }
      } catch { /* A corrupt entry cannot become mutation authority. */ }
    }
  } catch { /* Keep the in-memory fence if browser storage is unavailable. */ }
  for (const [key, record] of memory) if (record.envelope.request.workspace_id === workspaceId) found.set(key, record);
  return [...found.values()].map((item) => structuredClone(item));
}

/** A recognized root key remains a read-only fence even when its JSON/envelope cannot be recovered. */
export function loadUnreadableRecordingScopes(workspaceId: string): UnreadableRecordingChangeScope[] {
  const readable = new Set(loadRecordingChanges(workspaceId).map(recordingChangeKey));
  const scopes: UnreadableRecordingChangeScope[] = [];
  try {
    for (let index = 0; index < sessionStorage.length; index++) {
      const key = sessionStorage.key(index);
      if (!key?.startsWith(prefix) || readable.has(key)) continue;
      try {
        const ids: unknown = JSON.parse(key.slice(prefix.length));
        if (Array.isArray(ids) && ids.length === 4 && ids.every(id => typeof id === 'string' && id.length > 0 && id.length <= 256 && id.trim() === id) && ids[0] === workspaceId)
          scopes.push({ workspaceId: ids[0], groupId: ids[1], lineageId: ids[2], operationId: ids[3] });
      } catch { /* Unrecognizable keys cannot authorize any operation. */ }
    }
  } catch { /* The existing in-memory records remain fenced separately. */ }
  return scopes;
}
