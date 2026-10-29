import type { SchemaOperation, SchemaOperationStatus, SchemaPreviewToken } from '../../../src/types/recordingPlan';
import type * as S from '../../../src/types/studioV2WriteGroupSuccessor';
import { parseWriteGroupSuccessorReviewRequest, parseWriteGroupSuccessorEnvelope, parseWriteGroupSuccessorStageRequest,
  parseWriteGroupSuccessorSchemaConfirmation, parseWriteGroupSuccessorActivationEnvelope, parseWriteGroupSuccessorActivationConfirmation } from '../../../src/utils/studioV2WriteGroupSuccessorJson';
import { isRecord } from '../../../src/utils/studioV2LifecycleJsonPrimitives';
import { sameLifecycleJSON } from '../../../src/services/studioV2LifecycleTransport';
import { fixtureCurrent, fixtureHistory, fixtureReadiness, FIXTURE_HASH as HASH, FIXTURE_TIME as TIME, type LifecycleMockState } from './recording-lifecycle-state';

// Pure TEST-ONLY contract simulator: never opens a database or executes SQL.
export class MockContractError extends Error {
  constructor(readonly code = 'revision_mismatch', readonly status = 409) { super(code); }
}
function requireMock(condition: unknown, code?: string): asserts condition {
  if (!condition) throw new MockContractError(code);
}
const same = sameLifecycleJSON;
function exactKeys(value: unknown, keys: string[]) {
  if (!isRecord(value) || Object.keys(value).some(key => !keys.includes(key))) throw new MockContractError('WRITE_GROUP_INVALID_REQUEST', 400);
}
export class LifecycleSuccessorMock {
  private previews = new Map<string, { preview: SchemaPreviewToken; stage: S.WriteGroupSuccessorStageRequest }>();
  private confirmationStatus: SchemaOperationStatus = 'succeeded';
  private activationReady = true;
  private readinessIssued = false;
  constructor(readonly state: LifecycleMockState) {}
  setActivationReady(ready: boolean) { this.activationReady = ready; this.readinessIssued = false; }
  setSchemaOperationStatus(status: SchemaOperationStatus) {
    this.confirmationStatus = status;
    for (const entry of this.state.status?.schema_operations ?? []) {
      if (entry.status !== 'preview_not_submitted') entry.status = status;
    }
    const operation = this.state.status?.schema_operations.find(entry => entry.status === 'succeeded');
    if (operation && this.state.target) Object.assign(this.state.target.destination, { schema_revision: operation.operation_id, schema_digest: HASH });
  }
  review(body: unknown): S.WriteGroupSuccessorReview {
    const request = parseWriteGroupSuccessorReviewRequest(body);
    const { state } = this;
    requireMock(request && same(body, request), 'WRITE_GROUP_INVALID_REQUEST');
    requireMock(request && request.workspace_id === 'ws-1' && request.lineage_id === 'L' && !state.target &&
      request.expected_workspace_revision === state.workspaceRevision && request.expected_lineage_revision === state.lineageRevision &&
      request.expected_predecessor_revision === state.source.revision && request.expected_predecessor_connector_revision === 'crev-1');
    requireMock(state.scenario !== 'stale-wide');
    requireMock(state.scenario !== 'blocked-wide' && (state.source.status === 'deleted' || state.source.applied_revision));
    const targetRevision = request.group.destination.connector_id === 'conn-1' ? 'crev-1' : request.group.destination.connector_id === 'conn-2' ? 'crev-2' : undefined;
    requireMock(targetRevision && request.expected_target_connector_revision === targetRevision && request.group.destination.connector_revision === targetRevision);
    requireMock(same(request.group.members, state.source.members) && same(request.group.recording_profile, state.source.recording_profile) &&
      (request.group.layout_version === 'managed_wide_profile_v1' || request.group.layout_version === 'managed_long_v1' ||
        request.group.layout_version === 'managed_wide_v1' && state.source.layout_version === 'managed_wide_v1') &&
      (!request.group.storage_mode || request.group.storage_mode === 'plain'), 'WRITE_GROUP_SUCCESSOR_UNSUPPORTED');
    state.envelope = { version: 'successor-review-v1', predecessor_group_id: 'G1', generation: 1,
      source_applied_revision: state.source.applied_revision, source_definition_digest: HASH, request: structuredClone(request) };
    return { envelope: structuredClone(state.envelope), review_digest: HASH, history: fixtureHistory(state),
      blockers: state.scenario === 'review-blocked-wide' ? [{ reason_code: 'connector_identity_stale', next_action: 'review_connector' }] : [], required_checks: [] };
  }
  private root(body: unknown) {
    requireMock(isRecord(body));
    const envelope = parseWriteGroupSuccessorEnvelope(body.envelope);
    requireMock(envelope && this.state.envelope && same(body.envelope, envelope) && same(envelope, this.state.envelope) && body.review_digest === HASH);
    return envelope;
  }
  private stage(body: unknown) {
    const parsed = parseWriteGroupSuccessorStageRequest(body);
    this.root(body);
    requireMock(parsed && isRecord(body) && same(body.current, parsed.current) && this.state.target && same(parsed.current, fixtureCurrent(this.state)));
    return parsed;
  }
  handle(action: string, body: unknown): unknown {
    const keys: Record<string, string[]> = {
      reserve: ['envelope', 'review_digest'], 'schema-preview': ['envelope', 'review_digest', 'current'],
      'schema-confirm': ['envelope', 'review_digest', 'current', 'token', 'digest', 'schema_operation_id'],
      'activation-review': ['stage', 'boundary'], 'activation-readiness': ['envelope', 'review_digest'],
      'activation-confirm': ['envelope', 'review_digest', 'readiness_digest'], 'first-start-confirm': ['envelope', 'review_digest', 'readiness_digest'],
      'operation-status': ['workspace_id', 'lineage_id', 'operation_id'],
    };
    exactKeys(body, ['action', ...(keys[action] ?? [])]);
    if (action === 'reserve') return this.reserve(body);
    if (action === 'schema-preview') return this.preview(body);
    if (action === 'schema-confirm') return this.confirmSchema(body);
    if (action === 'activation-review') return this.activationReview(body);
    if (action === 'activation-readiness') return this.readiness(body);
    if (action === 'activation-confirm' || action === 'first-start-confirm') return this.confirmActivation(action, body);
    if (action === 'operation-status') {
      requireMock(isRecord(body) && this.state.status && body.workspace_id === 'ws-1' && body.lineage_id === 'L' &&
        body.operation_id === this.state.envelope?.request.operation_id, 'WRITE_GROUP_NOT_FOUND');
      return structuredClone(this.state.status);
    }
    throw new MockContractError('MOCK_UNSUPPORTED_ACTION', 404);
  }
  private reserve(body: unknown): S.WriteGroupSuccessorPreparation {
    const envelope = this.root(body);
    const { state } = this;
    requireMock(state.scenario !== 'review-blocked-wide');
    if (state.preparation) return structuredClone(state.preparation);
    requireMock(envelope.request.expected_workspace_revision === state.workspaceRevision && envelope.request.expected_predecessor_revision === state.source.revision);
    state.workspaceRevision = 'W2'; state.lineageRevision = 2;
    state.target = { ...structuredClone(envelope.request.group), id: 'G2', revision: 'R2', created_at: TIME, updated_at: TIME,
      destination: { ...envelope.request.group.destination, database: 'mock-only-successor.db', table_name: 'mock_recording_G2' } };
    state.preparation = { workspace_revision: 'W2', lineage_id: 'L', lineage_revision: 2, generation: 2, predecessor_group_id: 'G1',
      operation_id: envelope.request.operation_id, intent_digest: HASH, group: structuredClone(state.target) };
    state.status = { workspace_id: 'ws-1', lineage_id: 'L', predecessor_group_id: 'G1', operation_id: envelope.request.operation_id,
      reserve: { status: 'succeeded', result: structuredClone(state.preparation) }, schema_operations: [], cutover: { status: 'not_submitted' }, first_start: { status: 'not_submitted' } };
    return structuredClone(state.preparation);
  }
  private preview(body: unknown): SchemaPreviewToken {
    const stage = this.stage(body), { state } = this;
    requireMock(state.status && state.target && state.status.schema_operations.every(item => item.status === 'preview_not_submitted'));
    const id = `SCHEMA-${this.previews.size + 1}`;
    const preview: SchemaPreviewToken = { token: `MOCK-TOKEN-${id}`, operation_id: id, action: 'schema_apply', workspace_id: 'ws-1',
      workspace_revision: state.workspaceRevision, plan_id: 'G2', plan_revision: state.target.revision, connector_id: state.target.destination.connector_id,
      connector_revision: state.target.destination.connector_revision, table_prefix: '', statements: ['-- MOCK ONLY: no SQL is executed'],
      tables: [{ name: 'mock_recording_G2', action: 'create', columns: ['col_0', 'col_1'] }],
      group_layout: { table_name: 'mock_recording_G2', owner_column: '__owner', columns: [0, 1].map(index => ({ name: `col_${index}`, sql_type: 'REAL', nullable: true, primary_key: false })) },
      digest: HASH, expires_at: '2099-01-01T00:00:00Z', created_at: TIME };
    this.previews.set(id, { preview, stage });
    state.status.schema_operations.push({ operation_id: id, status: 'preview_not_submitted' });
    return structuredClone(preview);
  }
  private confirmSchema(body: unknown): SchemaOperation {
    const request = parseWriteGroupSuccessorSchemaConfirmation(body);
    this.stage(body);
    requireMock(request);
    const original = this.previews.get(request.schema_operation_id);
    requireMock(original && request.token === original.preview.token && request.digest === original.preview.digest &&
      same(request.current, original.stage.current) && same(request.envelope, original.stage.envelope));
    const operation = this.state.status?.schema_operations.find(item => item.operation_id === request.schema_operation_id);
    requireMock(operation && (operation.status === 'preview_not_submitted' || operation.status === 'succeeded'));
    operation.status = this.confirmationStatus;
    if (operation.status === 'succeeded' && this.state.target) {
      Object.assign(this.state.target.destination, { schema_revision: operation.operation_id, schema_digest: HASH });
    }
    return { operation_id: operation.operation_id, status: operation.status, verified_digest: operation.status === 'succeeded' ? HASH : undefined,
      executed_statements: operation.status === 'succeeded' ? 1 : 0, created_at: TIME, updated_at: TIME };
  }
  private activationReview(body: unknown): S.WriteGroupSuccessorActivationReview {
    requireMock(isRecord(body));
    exactKeys(body.stage, ['envelope', 'review_digest', 'current']);
    const stage = this.stage(body.stage), boundary = body.boundary;
    requireMock(typeof boundary === 'string' && Date.parse(boundary) > Date.now() && this.state.target?.destination.schema_revision);
    const envelope: S.WriteGroupSuccessorActivationEnvelope = { version: 'successor-activation-review-v2',
      activation_kind: this.state.source.applied_revision ? 'cutover' : 'first_start', request: { stage, boundary },
      source_applied_revision: this.state.source.applied_revision, source_definition_digest: HASH, target_definition_digest: HASH, counter_proof_digest: HASH,
      schema_revision: this.state.target.destination.schema_revision, schema_digest: HASH };
    this.state.activation = { envelope, review_digest: HASH, history: fixtureHistory(this.state), blockers: [] };
    this.readinessIssued = false;
    return structuredClone(this.state.activation);
  }
  private readiness(body: unknown): S.WriteGroupSuccessorActivationReadiness {
    requireMock(isRecord(body));
    const envelope = parseWriteGroupSuccessorActivationEnvelope(body.envelope);
    requireMock(envelope && same(body.envelope, envelope) && this.state.activation && same(envelope, this.state.activation.envelope) && body.review_digest === HASH && this.state.target);
    this.stage(envelope.request.stage);
    const readiness = fixtureReadiness(this.state, this.state.target);
    if (!this.activationReady) { readiness.ready = false; readiness.issues = [{ code: 'mock-readiness-blocked', severity: 'blocking', step: 'Step 4', scope: 'G2', message: 'Mock readiness blocked; review saved recording.' }]; }
    this.readinessIssued = this.activationReady;
    return { review: structuredClone(this.state.activation), readiness, readiness_digest: this.activationReady ? HASH : '' };
  }
  private confirmActivation(action: string, body: unknown) {
    const request = parseWriteGroupSuccessorActivationConfirmation(body), { state } = this;
    requireMock(this.readinessIssued && this.activationReady && request && isRecord(body) && same(body.envelope, request.envelope) && state.activation && state.status && state.target && same(request.envelope, state.activation.envelope) &&
      request.review_digest === HASH && request.readiness_digest === HASH && Date.parse(request.envelope.request.boundary) > Date.now());
    const first = request.envelope.activation_kind === 'first_start';
    requireMock(action === (first ? 'first-start-confirm' : 'activation-confirm'));
    if (first && state.status.first_start.status === 'succeeded') return state.status.first_start.result;
    if (!first && state.status.cutover.status === 'succeeded') return state.status.cutover.result;
    this.stage(request.envelope.request.stage);
    const boundary = request.envelope.request.boundary;
    state.workspaceRevision = 'W3'; state.lineageRevision = 3;
    // Backend records a ready scheduled generation, not proof of a running device or SQL delivery.
    state.target.applied_revision = state.target.revision; state.target.status = 'ready';
    if (first) {
      const result: S.WriteGroupSuccessorFirstStartResult = { operation_id: 'MOCK-FIRST-START', root_operation_id: state.status.operation_id,
        workspace_revision: 'W3', lineage_id: 'L', lineage_revision: 3, predecessor_group_id: 'G1', group_id: 'G2', applied_revision: state.target.revision,
        boundary, generation: 2, status: 'succeeded' };
      state.status.first_start = { status: 'succeeded', result }; return result;
    }
    const result = { id: 'MOCK-CUTOVER', workspace_id: 'ws-1', kind: 'forward' as const, source_group_id: 'G1', boundary, review_digest: HASH,
      ended: [{ group_id: 'G1', revision: state.source.applied_revision }], started: [{ group_id: 'G2', revision: state.target.revision }],
      moved: state.target.members.map(({ device_id, point_id, tag_id }) => ({ device_id, point_id, tag_id })), created_at: TIME };
    state.status.cutover = { status: 'succeeded', result }; return result;
  }
}
