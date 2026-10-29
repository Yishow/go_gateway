import type { SchemaOperation, SchemaPreviewToken } from '../../../src/types/recordingPlan';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
import type { RecordingStartRequest } from '../../../src/types/studioV2RecordingStart';
import { isRecord } from '../../../src/utils/studioV2LifecycleJsonPrimitives';
import { parseRecordingStartOperation } from '../../../src/utils/studioV2RecordingStartJson';
import { parseRecordingSchemaOperation, parseRecordingSchemaPreviewToken } from '../../../src/utils/recordingPlanJson';
import { AUTHORING_HASH as HASH, AUTHORING_TIME as TIME } from './recording-authoring-state';
import { MockContractError } from './recording-lifecycle-successor';
import { advanceGuidedHead, type GuidedMockState } from './recording-guided-state';

export const rejectGuided = (code = 'MOCK_UNSUPPORTED_ROUTE', status = 404): never => { throw new MockContractError(code, status); };
export function guidedBody(body: unknown, keys: readonly string[]): Record<string, unknown> {
  if (!isRecord(body) || Object.keys(body).some(key => !keys.includes(key))) return rejectGuided('MOCK_INVALID_REQUEST', 400);
  return body;
}
export const GUIDED_CAS_KEYS = ['workspace_id', 'expected_workspace_revision', 'expected_group_revision', 'expected_connector_revision'];
export function guidedCAS(state: GuidedMockState, group: WriteGroup, body: Record<string, unknown>) {
  if (body.workspace_id !== 'ws-1' || body.expected_workspace_revision !== state.workspaceRevision ||
    body.expected_group_revision !== group.revision || body.expected_connector_revision !== group.destination.connector_revision ||
    group.status === 'deleted') rejectGuided('MOCK_REVISION_MISMATCH', 409);
}
export type GuidedStartOutcome = 'succeeded' | 'partial' | 'unknown' | 'pending';

/** Immutable original preview/start receipts; all canonical groups/head remain in the supplied state. */
export class GuidedRecordingOperations {
  private previews = new Map<string, SchemaPreviewToken>();
  private schemaOperations = new Map<string, SchemaOperation>();
  private outcome: GuidedStartOutcome;
  constructor(private readonly state: GuidedMockState) {
    this.outcome = state.scenario === 'partial' || state.scenario === 'unknown' ? state.scenario : 'succeeded';
  }
  get previewCount() { return this.previews.size; }
  setStartOutcome(value: GuidedStartOutcome) { this.outcome = value; }
  preview(group: WriteGroup, input: unknown) {
    const body = guidedBody(input, GUIDED_CAS_KEYS); guidedCAS(this.state, group, body);
    if (group.destination.storage_strategy !== 'managed') rejectGuided('MOCK_CUSTOM_SCHEMA_FORBIDDEN', 409);
    const number = this.previews.size + 1, operationId = `SCHEMA-${number}`;
    const existing = Boolean(group.destination.schema_digest);
    const columns = group.members.map((member, index) => member.target_column || `value_${index}`);
    const preview = parseRecordingSchemaPreviewToken({ token: `MOCK-SCHEMA-TOKEN-${number}`, operation_id: operationId, action: 'schema_apply',
      workspace_id: 'ws-1', workspace_revision: this.state.workspaceRevision, plan_id: group.id, plan_revision: group.revision,
      connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision, table_prefix: '',
      statements: existing ? [] : [`CREATE TABLE ${group.destination.table_name} (__mock_owner TEXT)`],
      tables: [{ name: group.destination.table_name, action: existing ? 'unchanged' : 'create', columns }],
      group_layout: { table_name: group.destination.table_name, owner_column: '__mock_owner',
        columns: columns.map(name => ({ name, sql_type: 'REAL', nullable: true, primary_key: false })) },
      digest: HASH, ...(existing ? { schema_revision: group.destination.schema_revision, schema_digest: HASH, no_change_reason: 'schema_current' } : {}),
      expires_at: '2099-01-01T00:00:00Z', created_at: TIME,
    }, { requireGroupLayout: true });
    if (!preview) throw new Error('Malformed guided schema fixture');
    this.previews.set(operationId, preview); return preview;
  }
  apply(group: WriteGroup, input: unknown) {
    const body = guidedBody(input, [...GUIDED_CAS_KEYS, 'token', 'operation_id']); guidedCAS(this.state, group, body);
    const preview = typeof body.operation_id === 'string' ? this.previews.get(body.operation_id) : undefined;
    if (!preview || body.token !== preview.token || preview.workspace_revision !== this.state.workspaceRevision ||
      preview.plan_id !== group.id || preview.plan_revision !== group.revision || preview.connector_revision !== group.destination.connector_revision ||
      this.schemaOperations.has(preview.operation_id)) return rejectGuided('MOCK_SCHEMA_CONFIRMATION_MISMATCH', 409);
    const operation = parseRecordingSchemaOperation({ operation_id: preview.operation_id, action: 'schema_apply', status: 'succeeded',
      executed_statements: preview.statements.length, verified_digest: HASH, created_at: TIME, updated_at: TIME, completed_at: TIME });
    if (!operation) throw new Error('Malformed guided schema operation fixture');
    advanceGuidedHead(this.state);
    group.revision = `R${this.state.workspaceRevision.slice(1)}`;
    group.destination.schema_revision = group.revision; group.destination.schema_digest = HASH;
    group.status = 'ready'; this.schemaOperations.set(operation.operation_id, operation); return operation;
  }
  schemaOperation(id: string) { return this.schemaOperations.get(id) ?? rejectGuided('MOCK_OPERATION_NOT_FOUND'); }
  start(input: unknown) {
    const body = guidedBody(input, ['request_id', 'workspace_id', 'expected_workspace_revision', 'device_ids', 'groups']);
    if (this.state.scenario === 'blocked-wide') rejectGuided('MOCK_CONNECTOR_STALE', 409);
    if (this.state.originalStart) rejectGuided('MOCK_ORIGINAL_REQUEST_RETAINED', 409);
    if (typeof body.request_id !== 'string' || !body.request_id.trim() || body.workspace_id !== 'ws-1' || body.expected_workspace_revision !== this.state.workspaceRevision ||
      !Array.isArray(body.device_ids) || body.device_ids.length !== 1 || body.device_ids[0] !== 'dev-1' || !Array.isArray(body.groups) || body.groups.length === 0 ||
      new Set(body.groups.map(item => isRecord(item) ? item.group_id : undefined)).size !== body.groups.length) return rejectGuided('MOCK_START_SCOPE_MISMATCH', 409);
    const groups = body.groups.map(item => {
      const intent = guidedBody(item, ['group_id', 'expected_group_revision', 'expected_connector_revision']);
      const group = this.state.groups.find(candidate => candidate.id === intent.group_id);
      if (!group || group.members.some(member => member.device_id !== 'dev-1') || !group.members.length || group.status === 'disabled' || group.status === 'deleted' ||
        group.revision !== intent.expected_group_revision || group.destination.connector_revision !== intent.expected_connector_revision ||
        group.destination.storage_strategy !== 'managed' || group.destination.schema_revision !== group.revision || !group.destination.schema_digest) return rejectGuided('MOCK_START_PREREQUISITE', 409);
      return group;
    });
    const request = structuredClone(body) as unknown as RecordingStartRequest;
    const succeeded = this.outcome === 'succeeded', partial = this.outcome === 'partial';
    const operation = parseRecordingStartOperation({ operation_id: 'O1', action: 'recording_start', status: this.outcome,
      stage: succeeded ? 'complete' : 'activation', intent_digest: HASH, workspace_id: 'ws-1', setup_revision: this.state.workspaceRevision,
      device_ids: ['dev-1'], groups: groups.map(group => ({ group_id: group.id, group_revision: group.revision, saved: true, ready: true,
        applied: succeeded, ...(succeeded ? { applied_revision: group.revision } : {}) })),
      devices: [{ device_id: 'dev-1', activated: succeeded || partial }], created_at: TIME, updated_at: TIME });
    if (!operation) throw new Error('Malformed guided start fixture');
    this.state.originalStart = request; this.state.operation = operation;
    if (succeeded) for (const group of groups) { group.applied_revision = group.revision; group.status = 'running'; }
    return operation;
  }
  operation(id: string) {
    if (!this.state.operation || this.state.operation.operation_id !== id) rejectGuided('MOCK_OPERATION_NOT_FOUND');
    return this.state.operation;
  }
}
