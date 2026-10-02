import type {
  WriteGroupTestCleanupStatus,
  WriteGroupTestWriteOperation,
  WriteGroupTestWriteOutcome,
  WriteGroupTestWritePreview,
  WriteGroupTestWriteStatus,
  WriteGroupTestWriteValue,
} from '../types/studioV2WriteGroupTestWrite';
import { boundedString, MAX_SAFE_JSON_ARRAY_LENGTH } from './safeJson';

type JsonRecord = Record<string, unknown>;

const OUTCOMES = new Set<WriteGroupTestWriteOutcome>(['written_verified', 'written_unverified', 'failed', 'unknown']);
const CLEANUPS = new Set<WriteGroupTestCleanupStatus>(['not_attempted', 'cleaned', 'failed', 'unknown']);
const STATUSES = new Set<WriteGroupTestWriteStatus>(['pending', 'running', 'succeeded', 'partial', 'failed', 'unknown']);

/** The operation status a write outcome is allowed to end in. */
const STATUS_FOR_OUTCOME: Record<WriteGroupTestWriteOutcome, WriteGroupTestWriteStatus> = {
  written_verified: 'succeeded',
  written_unverified: 'partial',
  failed: 'failed',
  unknown: 'unknown',
};

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function parseValues(value: unknown): WriteGroupTestWriteValue[] | null {
  if (!Array.isArray(value) || value.length === 0 || value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const values: WriteGroupTestWriteValue[] = [];
  for (const item of value) {
    if (!isRecord(item)) return null;
    const column = boundedString(item.column);
    const type = boundedString(item.type);
    if (!column || !type || typeof item.value !== 'string' || item.value.length > 1024) return null;
    values.push({ column, type, value: item.value });
  }
  return values;
}

/** Parses a test-write preview; a token without its operation or owner marker is unusable. */
export function parseWriteGroupTestWritePreview(value: unknown): WriteGroupTestWritePreview | null {
  if (!isRecord(value) || !isRecord(value.target) || value.action !== 'test_write') return null;
  const token = boundedString(value.token);
  const operationId = boundedString(value.operation_id);
  const groupId = boundedString(value.group_id);
  const groupRevision = boundedString(value.group_revision);
  const expiresAt = boundedString(value.expires_at);
  const ownerColumn = boundedString(value.owner_column);
  const ownerValue = boundedString(value.owner_value);
  const cleanup = boundedString(value.cleanup);
  const connectorId = boundedString(value.target.connector_id);
  const dialect = boundedString(value.target.dialect);
  const table = boundedString(value.target.table);
  const values = parseValues(value.values);
  if (!token || !operationId || !groupId || !groupRevision || !expiresAt || Number.isNaN(Date.parse(expiresAt)) ||
    !ownerColumn || !ownerValue || !cleanup || !connectorId || !dialect || !table || !values ||
    typeof value.dedupe !== 'string' || typeof value.target.database !== 'string' || typeof value.target.schema !== 'string') {
    return null;
  }
  // The owner marker is the only thing cleanup may match, so it must be what the preview lists for that column.
  if (!values.some((entry) => entry.column === ownerColumn && entry.value === ownerValue)) return null;
  return {
    token, operation_id: operationId, action: 'test_write', group_id: groupId, group_revision: groupRevision,
    expires_at: expiresAt,
    target: { connector_id: connectorId, dialect, database: value.target.database, schema: value.target.schema, table },
    owner_column: ownerColumn, owner_value: ownerValue, dedupe: value.dedupe, values, cleanup,
  };
}

function optionalCode(value: JsonRecord, key: string): string | undefined | null {
  if (value[key] === undefined || value[key] === '') return undefined;
  return boundedString(value[key]) ?? null;
}

/**
 * Parses a test-write operation. Contradictions are rejected rather than shown:
 * a running operation cannot already carry a result, a verified write cannot
 * claim its cleanup was never attempted, and a write that failed wrote nothing
 * to clean.
 */
export function parseWriteGroupTestWriteOperation(value: unknown): WriteGroupTestWriteOperation | null {
  if (!isRecord(value) || value.action !== 'test_write') return null;
  const operationId = boundedString(value.operation_id);
  const status = boundedString(value.status);
  const createdAt = boundedString(value.created_at);
  const updatedAt = boundedString(value.updated_at);
  const outcome = optionalCode(value, 'write_outcome');
  const cleanup = optionalCode(value, 'cleanup_status');
  const reason = optionalCode(value, 'reason');
  const cleanupReason = optionalCode(value, 'cleanup_reason');
  const completedAt = optionalCode(value, 'completed_at');
  if (!operationId || !status || !STATUSES.has(status as WriteGroupTestWriteStatus) || !createdAt || !updatedAt ||
    outcome === null || cleanup === null || reason === null || cleanupReason === null || completedAt === null ||
    (outcome !== undefined && !OUTCOMES.has(outcome as WriteGroupTestWriteOutcome)) ||
    (cleanup !== undefined && !CLEANUPS.has(cleanup as WriteGroupTestCleanupStatus))) {
    return null;
  }
  const active = status === 'pending' || status === 'running';
  if (active) {
    if (outcome !== undefined || cleanup !== undefined || completedAt !== undefined) return null;
  } else {
    if (outcome === undefined || cleanup === undefined || completedAt === undefined) return null;
    if (STATUS_FOR_OUTCOME[outcome as WriteGroupTestWriteOutcome] !== status) return null;
    if (outcome === 'failed' && cleanup !== 'not_attempted') return null;
    if (outcome === 'written_verified' && cleanup === 'not_attempted') return null;
  }
  return {
    operation_id: operationId, action: 'test_write', status: status as WriteGroupTestWriteStatus,
    ...(outcome ? { write_outcome: outcome as WriteGroupTestWriteOutcome } : {}),
    ...(cleanup ? { cleanup_status: cleanup as WriteGroupTestCleanupStatus } : {}),
    ...(reason ? { reason } : {}),
    ...(cleanupReason ? { cleanup_reason: cleanupReason } : {}),
    created_at: createdAt, updated_at: updatedAt,
    ...(completedAt ? { completed_at: completedAt } : {}),
  };
}
