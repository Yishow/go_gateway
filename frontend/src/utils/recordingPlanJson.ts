import {
  boundedString,
  MAX_SAFE_JSON_ARRAY_LENGTH,
  normalizeTypedEnvelope,
  parseBoundedJson,
} from './safeJson';
import type {
  ConnectorCapability,
  GroupSchemaColumn,
  GroupSchemaLayout,
  SchemaOperation,
  SchemaOperationStatus,
  SchemaPreviewTable,
  SchemaPreviewTableAction,
  SchemaPreviewToken,
  TestWriteResult,
  TestWriteResultStatus,
} from '../types/recordingPlan';

const TEST_WRITE_STATUSES = new Set<TestWriteResultStatus>([
  'written_verified',
  'written_unverified',
  'failed',
  'unknown',
]);
const MAX_RECORDING_STATEMENT_LENGTH = 16 * 1024;
const MAX_RECORDING_STATEMENTS_TOTAL_LENGTH = 64 * 1024;

export type RecordingOperationOutcome = 'failed' | 'unconfirmed';

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

/** Bounded recording text: trimmed, non-empty, at most MAX_RECORDING_STATEMENT_LENGTH (16KB) characters. */
function boundedRecordingText(value: unknown): string | undefined {
  if (typeof value !== 'string') return undefined;
  const normalized = value.trim();
  return normalized.length > 0 && normalized.length <= MAX_RECORDING_STATEMENT_LENGTH
    ? normalized
    : undefined;
}

function optionalShortText(record: Record<string, unknown>, key: string): string | undefined | null {
  if (record[key] === undefined || record[key] === null) return undefined;
  return boundedString(record[key]) ?? null;
}

/** Preview proof fields may be explicitly empty before the first schema apply. */
function optionalPreviewText(record: Record<string, unknown>, key: string): string | undefined | null {
  if (record[key] === undefined || record[key] === null) return undefined;
  if (typeof record[key] !== 'string' || record[key].length > MAX_RECORDING_STATEMENT_LENGTH) return null;
  const normalized = record[key].trim();
  return normalized || undefined;
}

/**
 * Error raised when a recording API response is not safe to interpret.
 * The response body is never retained; only allowlisted typed metadata is copied.
 */
export class RecordingPlanResponseError extends Error {
  readonly code?: string;
  readonly action?: string;
  readonly retryable?: boolean;
  readonly request_id?: string;
  readonly operation_id?: string;
  readonly outcome: RecordingOperationOutcome;

  constructor(
    operation: string,
    source?: unknown,
    outcome: RecordingOperationOutcome = 'unconfirmed',
  ) {
    super(`${operation} response was invalid`);
    this.name = 'RecordingPlanResponseError';
    this.outcome = outcome;
    const envelope = normalizeTypedEnvelope(source);
    this.code = envelope.code;
    this.action = envelope.action;
    this.retryable = envelope.retryable;
    this.request_id = envelope.requestId;
    this.operation_id = envelope.operationId;
  }
}

/** Parse the backend's recording test-write result without trusting its status string. */
export function parseRecordingTestWriteResult(value: unknown): TestWriteResult | null {
  const bounded = parseBoundedJson(value);
  if (!isRecord(bounded)) return null;
  const status = boundedString(bounded.status);
  if (!status || !TEST_WRITE_STATUSES.has(status as TestWriteResultStatus)) return null;
  const recordId = boundedString(bounded.record_id);
  const table = boundedString(bounded.table);
  const observedAt = boundedString(bounded.observed_at);
  const deliveredAt = boundedString(bounded.delivered_at);
  if (!recordId || !table || !observedAt || !deliveredAt) return null;
  const message = optionalShortText(bounded, 'message');
  const operationId = optionalShortText(bounded, 'operation_id');
  if (message === null || operationId === null) return null;
  return {
    status: status as TestWriteResultStatus,
    record_id: recordId,
    table,
    observed_at: observedAt,
    delivered_at: deliveredAt,
    ...(operationId ? { operation_id: operationId } : {}),
    ...(message ? { message } : {}),
  };
}

const SCHEMA_OPERATION_STATUSES = new Set<SchemaOperationStatus>([
  'pending', 'running', 'succeeded', 'partial', 'failed', 'unknown',
]);
const SCHEMA_OPERATION_TEXT_KEYS = ['action', 'verified_digest', 'reason', 'next_action', 'completed_at'] as const;
type SchemaOperationText = Partial<Record<typeof SCHEMA_OPERATION_TEXT_KEYS[number], string>>;

/** Optional operation text; an empty value means the backend had none to report. */
function schemaOperationText(value: Record<string, unknown>): SchemaOperationText | null {
  const parsed: SchemaOperationText = {};
  for (const key of SCHEMA_OPERATION_TEXT_KEYS) {
    if (value[key] === '') continue;
    const text = optionalShortText(value, key);
    if (text === null) return null;
    if (text) parsed[key] = text;
  }
  return parsed;
}

/** Parse a schema operation; an unknown status invalidates the whole result. */
export function parseRecordingSchemaOperation(value: unknown): SchemaOperation | null {
  const bounded = parseBoundedJson(value);
  if (!isRecord(bounded)) return null;
  const operationId = boundedString(bounded.operation_id);
  const status = boundedString(bounded.status);
  const createdAt = boundedString(bounded.created_at);
  const updatedAt = boundedString(bounded.updated_at);
  const executed = bounded.executed_statements;
  const text = schemaOperationText(bounded);
  if (!operationId || !status || !SCHEMA_OPERATION_STATUSES.has(status as SchemaOperationStatus) ||
    !createdAt || !updatedAt || !text) return null;
  if (typeof executed !== 'number' || !Number.isInteger(executed) || executed < 0) return null;
  return {
    operation_id: operationId,
    status: status as SchemaOperationStatus,
    executed_statements: executed,
    ...text,
    created_at: createdAt,
    updated_at: updatedAt,
  };
}

const SCHEMA_PREVIEW_TABLE_ACTIONS = new Set<SchemaPreviewTableAction>(['create', 'unchanged']);
const PREVIEW_SCOPE_TEXT_KEYS = ['action', 'dialect', 'database', 'schema', 'no_change_reason'] as const;
type PreviewScopeText = Partial<Record<typeof PREVIEW_SCOPE_TEXT_KEYS[number], string>>;

/** Parse the per-table preview summary; unknown actions invalidate the whole preview. */
function parseSchemaPreviewTables(value: unknown): SchemaPreviewTable[] | null {
  if (!Array.isArray(value) || value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const tables: SchemaPreviewTable[] = [];
  for (const item of value) {
    if (!isRecord(item)) return null;
    const name = boundedString(item.name);
    const action = boundedString(item.action);
    const columns = item.columns ?? [];
    if (!name || !action || !SCHEMA_PREVIEW_TABLE_ACTIONS.has(action as SchemaPreviewTableAction) ||
      !Array.isArray(columns) || columns.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
    const parsedColumns = columns.map((column) => boundedString(column));
    if (parsedColumns.some((column) => !column)) return null;
    tables.push({ name, action: action as SchemaPreviewTableAction, columns: parsedColumns as string[] });
  }
  return tables;
}

function parseGroupSchemaLayout(value: unknown): GroupSchemaLayout | undefined | null {
  if (value === undefined || value === null) return undefined;
  if (!isRecord(value)) return null;
  const tableName = boundedString(value.table_name);
  const ownerColumn = boundedString(value.owner_column);
  if (!tableName || !ownerColumn || !Array.isArray(value.columns) || value.columns.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }
  const columns: GroupSchemaColumn[] = [];
  for (const item of value.columns) {
    if (!isRecord(item)) return null;
    const name = boundedString(item.name);
    const sqlType = boundedString(item.sql_type);
    if (!name || !sqlType || typeof item.nullable !== 'boolean' || typeof item.primary_key !== 'boolean') return null;
    columns.push({ name, sql_type: sqlType, nullable: item.nullable, primary_key: item.primary_key });
  }
  return { table_name: tableName, columns, owner_column: ownerColumn };
}

/** Optional preview scope text; an empty value means the server had none to report. */
function previewScopeText(value: Record<string, unknown>): PreviewScopeText | null {
  const parsed: PreviewScopeText = {};
  for (const key of PREVIEW_SCOPE_TEXT_KEYS) {
    if (value[key] === '') continue;
    const text = optionalShortText(value, key);
    if (text === null) return null;
    if (text) parsed[key] = text;
  }
  return parsed;
}

/** Parse a bounded schema-preview token while allowing multi-line DDL statements. */
export interface RecordingSchemaPreviewParseOptions {
  /** Managed write-group previews must include the server-resolved exact layout. */
  requireGroupLayout?: boolean;
}

export function parseRecordingSchemaPreviewToken(
  value: unknown,
  options: RecordingSchemaPreviewParseOptions = {},
): SchemaPreviewToken | null {
  if (!isRecord(value) || !Array.isArray(value.statements) || value.statements.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const token = boundedString(value.token);
  const workspaceId = boundedString(value.workspace_id);
  const planId = boundedString(value.plan_id);
  const planRevision = boundedString(value.plan_revision);
  const connectorId = boundedString(value.connector_id);
  const connectorRevision = value.connector_revision === undefined ? undefined : boundedString(value.connector_revision);
  if (value.connector_revision !== undefined && !connectorRevision) return null;
  // Managed write-group previews resolve the canonical table in group_layout;
  // their legacy table prefix is intentionally an explicit empty string.
  const tablePrefix = value.table_prefix === '' ? '' : boundedString(value.table_prefix);
  const expiresAt = boundedString(value.expires_at);
  const createdAt = boundedString(value.created_at);
  const operationId = boundedString(value.operation_id);
  const workspaceRevision = boundedString(value.workspace_revision);
  const digest = boundedString(value.digest);
  const tables = parseSchemaPreviewTables(value.tables);
  const groupLayout = parseGroupSchemaLayout(value.group_layout);
  const sourceDigest = optionalPreviewText(value, 'source_digest');
  const schemaRevision = optionalPreviewText(value, 'schema_revision');
  const schemaDigest = optionalPreviewText(value, 'schema_digest');
  const scopeText = previewScopeText(value);
  if (!token || !workspaceId || !planId || !planRevision || !connectorId || tablePrefix === undefined ||
    (tablePrefix.length === 0 && groupLayout === undefined) || !expiresAt || !createdAt ||
    !operationId || !workspaceRevision || !digest || !tables || !scopeText || groupLayout === null ||
    sourceDigest === null || schemaRevision === null || schemaDigest === null ||
    (options.requireGroupLayout && groupLayout === undefined)) return null;
  let totalLength = 0;
  const statements = value.statements.map((statement) => {
    const parsed = boundedRecordingText(statement);
    if (parsed) totalLength += parsed.length;
    return parsed;
  });
  if (statements.some((statement) => !statement) || totalLength > MAX_RECORDING_STATEMENTS_TOTAL_LENGTH) return null;
  return {
    token,
    workspace_id: workspaceId,
    plan_id: planId,
    plan_revision: planRevision,
    connector_id: connectorId,
    ...(connectorRevision ? { connector_revision: connectorRevision } : {}),
    table_prefix: tablePrefix,
    statements: statements as string[],
    tables,
    ...(groupLayout ? { group_layout: groupLayout } : {}),
    ...(sourceDigest ? { source_digest: sourceDigest } : {}),
    ...(schemaRevision ? { schema_revision: schemaRevision } : {}),
    ...(schemaDigest ? { schema_digest: schemaDigest } : {}),
    ...scopeText,
    operation_id: operationId,
    workspace_revision: workspaceRevision,
    digest,
    expires_at: expiresAt,
    created_at: createdAt,
  };
}

/** Parse the capability matrix so UI flags cannot be enabled by truthy non-booleans. */
export function parseRecordingConnectorCapability(value: unknown): ConnectorCapability | null {
  if (!isRecord(value)) return null;
  const kind = boundedString(value.kind);
  const supportedModes = value.supported_modes;
  if (!kind || !Array.isArray(supportedModes) || supportedModes.length > MAX_SAFE_JSON_ARRAY_LENGTH ||
    supportedModes.some((mode) => !boundedString(mode)) ||
    typeof value.supported !== 'boolean' ||
    typeof value.supports_managed_schema !== 'boolean' ||
    typeof value.supports_transactions !== 'boolean' ||
    typeof value.supports_receipts !== 'boolean' ||
    typeof value.supports_test_writes !== 'boolean') return null;
  if (value.supports_group_test_writes !== undefined && typeof value.supports_group_test_writes !== 'boolean') return null;
  const notes = optionalShortText(value, 'notes');
  if (notes === null) return null;
  return {
    kind,
    supported: value.supported,
    supports_managed_schema: value.supports_managed_schema,
    supports_transactions: value.supports_transactions,
    supports_receipts: value.supports_receipts,
    supports_test_writes: value.supports_test_writes,
    ...(value.supports_group_test_writes !== undefined ? { supports_group_test_writes: value.supports_group_test_writes } : {}),
    supported_modes: supportedModes.map((mode) => boundedString(mode) as string),
    ...(notes ? { notes } : {}),
  };
}

/**
 * Validate a standard API success envelope and parse its data at the operation seam.
 * A nominal 200 with an invalid body is deliberately classified as unconfirmed.
 */
export function parseRecordingAPIData<T>(
  value: unknown,
  operation: string,
  parseData: (data: unknown) => T | null,
): T {
  if (!isRecord(value) || typeof value.success !== 'boolean') {
    throw new RecordingPlanResponseError(operation, value);
  }
  if (!value.success) {
    throw new RecordingPlanResponseError(operation, value, 'failed');
  }
  if (!Object.prototype.hasOwnProperty.call(value, 'data')) {
    throw new RecordingPlanResponseError(operation, value);
  }
  const parsed = parseData(value.data);
  if (parsed === null) throw new RecordingPlanResponseError(operation, value);
  return parsed;
}

/** Return the existing opaque operation identity if the server supplied one. */
export function getRecordingOperationId(value: unknown): string | undefined {
  return normalizeTypedEnvelope(value).operationId;
}

const RECORDING_NOT_IMPLEMENTED_CODES = new Set([
  'RECORDING_SCHEMA_NOT_IMPLEMENTED',
  'RECORDING_TEST_WRITE_NOT_IMPLEMENTED',
]);

/** Transport loss and malformed 200 responses need explicit unconfirmed UI state. */
export function isRecordingOperationUnconfirmed(value: unknown): boolean {
  if (value instanceof RecordingPlanResponseError) return value.outcome === 'unconfirmed';
  if (!isRecord(value)) return true;
  const response = isRecord(value.response) ? value.response : undefined;
  if (!response) return true;
  const envelope = normalizeTypedEnvelope(value);
  const status = typeof response.status === 'number' ? response.status : undefined;
  if (status === 501 && envelope.code && RECORDING_NOT_IMPLEMENTED_CODES.has(envelope.code)) return false;
  if (status === undefined || !Number.isInteger(status) || status < 100 || status > 599) return true;
  return status >= 500;
}
