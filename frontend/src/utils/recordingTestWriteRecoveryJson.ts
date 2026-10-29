import type { RecordingTestWriteRecoveryList, RecordingTestWriteRecoveryOperation } from '../types/recordingTestWriteRecovery';
import { parseWriteGroupTestWriteOperation } from './studioV2WriteGroupTestWriteJson';
import { testWriteResolved } from '../features/datalink/workbench-v2/state/recordingTestWriteRecovery';
const identityFields = ['operation_id', 'workspace_id', 'workspace_revision', 'group_id', 'group_revision', 'connector_id', 'connector_revision', 'target_fingerprint'] as const;
function record(value: unknown): value is Record<string, unknown> {
  return value !== null && typeof value === 'object' && !Array.isArray(value);
}
function timestamp(value: unknown): value is string {
  if (typeof value !== 'string' || value.length > 64) return false;
  const match = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})(?:\.\d{1,9})?(?:Z|[+-](\d{2}):(\d{2}))$/.exec(value);
  if (!match || !Number.isFinite(Date.parse(value))) return false;
  const [, year, month, day, hour, minute, second, offsetHour = '0', offsetMinute = '0'] = match;
  const days = new Date(Date.UTC(Number(year), Number(month), 0)).getUTCDate();
  return Number(month) >= 1 && Number(month) <= 12 && Number(day) >= 1 && Number(day) <= days &&
    Number(hour) < 24 && Number(minute) < 60 && Number(second) < 60 && Number(offsetHour) < 24 && Number(offsetMinute) < 60;
}
/** A bounded allowlist: nothing from raw diagnostics, target names or credentials escapes. */
export function parseRecordingTestWriteRecovery(value: unknown, workspaceId: string): RecordingTestWriteRecoveryList {
  const invalid = () => new Error('Test-write recovery is unavailable.');
  if (!workspaceId || !record(value) || value.success !== true || !record(value.data)) throw invalid();
  const data = value.data;
  if (data.workspace_id !== workspaceId || data.complete !== true || !Array.isArray(data.operations) || data.operations.length > 1000) throw invalid();
  const seen = new Set<string>();
  const operations: RecordingTestWriteRecoveryOperation[] = [];
  for (const raw of data.operations) {
    if (!record(raw) || identityFields.some((field) => typeof raw[field] !== 'string' ||
      !raw[field] || raw[field].length > 256 || raw[field].trim() !== raw[field]) || raw.workspace_id !== workspaceId ||
      typeof raw.target_fingerprint !== 'string' || !/^[a-f0-9]{64}$/.test(raw.target_fingerprint)) throw invalid();
    if (!timestamp(raw.created_at) || !timestamp(raw.updated_at) || Date.parse(raw.updated_at) < Date.parse(raw.created_at) ||
      (raw.completed_at !== undefined && (!timestamp(raw.completed_at) || Date.parse(raw.completed_at) < Date.parse(raw.created_at) || Date.parse(raw.completed_at) > Date.parse(raw.updated_at))) ||
      raw.write_outcome === '' || raw.cleanup_status === '') throw invalid();
    const operation = parseWriteGroupTestWriteOperation(raw);
    if (!operation || testWriteResolved(operation) || seen.has(operation.operation_id)) throw invalid();
    seen.add(operation.operation_id);
    // Explicitly omit reason fields even if an unexpected server sends them.
    const { reason: _reason, cleanup_reason: _cleanupReason, ...safeOperation } = operation;
    operations.push({ ...safeOperation,
      workspace_id: workspaceId, workspace_revision: raw.workspace_revision as string,
      group_id: raw.group_id as string, group_revision: raw.group_revision as string,
      connector_id: raw.connector_id as string, connector_revision: raw.connector_revision as string,
      target_fingerprint: raw.target_fingerprint,
    });
  }
  return { workspace_id: workspaceId, complete: true, operations };
}
