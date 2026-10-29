import type { RecordingOwnerObservations, RecordingOwnerScope } from '@/features/datalink/workbench-v2/state/recordingOwnerObservation';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { RecordingNextAction, type RecordingNextActionLabels } from '@/features/datalink/workbench-v2/steps/step4/RecordingNextAction';
import type { RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import type { WriteGroup, WriteGroupDraft, WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import type { SchemaPreviewToken } from '@/types/recordingPlan';
import type { CandidateCounterStatus } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { recordingNextAction, recordingWorkflowView, type RecordingNextActionInput, type RecordingReadonly,
  type RecordingWorkflowEvidence, type RecordingWorkflowInput, type RecordingWorkflowOperation,
} from '@/features/datalink/workbench-v2/state/recordingWorkflowView';

const scope = { workspaceId: 'W1', deviceId: 'A', circuitId: 'C1', groupId: 'G1',
  workspaceRevision: 'W2', groupRevision: 'R2', appliedRevision: 'R1' };
function observed(selected: RecordingOwnerScope): RecordingOwnerObservations {
  return ['schema', 'start', 'test_write'].map(owner => ({ owner: owner as 'schema' | 'start' | 'test_write', scope: selected, state: 'observed', result: 'none' }));
}
const nowMs = Date.parse('2026-10-27T00:00:00Z');
const labels: RecordingNextActionLabels = {
  heading: 'Next safe step', scopes: { draft: 'Unsaved draft', group: 'Shared or custom group', managed_device: 'Managed device' },
  actions: { choose_sources: 'Choose sources', resolve_semantics: 'Review counter confirmation', save_prepare: 'Review and save preparation',
    preview_schema: 'Preview table changes', recover_operation: 'Inspect retained operation', review_schema: 'Review schema', confirm_schema: 'Confirm reviewed schema', review_scoped_start: 'Review scoped start',
    query_operation: 'Check existing operation', refresh_evidence: 'Refresh evidence', monitor: 'Open monitor', inspect: 'Inspect this scope' },
  reasons: { sources_empty: 'Select sources for this scope.', counter_confirmation: 'Review the selected counter definitions.',
    save_required: 'Save the draft before preparing its schema.', offline: 'Reconnect before this request; inspection remains available.',
    operation_recovery: 'Check the existing operation before a new request.', evidence_loading: 'Evidence is still loading.',
    evidence_error: 'Refresh evidence after the query error.', evidence_unavailable: 'Fetch current evidence.',
    schema_review: 'Review or refresh the schema preview.', schema_confirmation: 'Confirm only the exact reviewed schema.',
    schema_conflict: 'Inspect and repair the blocked schema.', configuration_blocked: 'Review the configuration blockers.',
    readiness_blocked: 'Refresh the unresolved server readiness.', lifecycle_unavailable: 'Inspect this scope until server action availability is known.',
    lifecycle_blocked: 'Inspect the server action blockers.', start_review: 'Review this scope before starting.',
    monitor_available: 'Open the confirmed activated scope.', scope_unavailable: 'Inspect the unavailable selection without changing it.' },
};
afterEach(() => { cleanup(); vi.restoreAllMocks(); });
const gauge = { device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1',
  source_revision: 'S1', mapping_revision: 'M1', target_column: 'temperature', required: true, value_role: 'gauge' as const };

function draft(): WriteGroupDraft {
  return { workspace_id: 'W1', name: 'Draft', members: [],
    destination: { connector_id: 'DB1', connector_revision: 'D1', table_schema: 'public', table_name: 'readings', storage_strategy: 'managed' },
    recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second' },
    row_policy: { interval_seconds: 30, allowed_lateness_seconds: 2 }, write_policy: {} };
}
function offlineDraft(): Extract<RecordingNextActionInput, { kind: 'draft' }> {
  return { kind: 'draft', scope: { workspaceId: 'W1', deviceId: 'A', circuitId: 'C1', workspaceRevision: 'W2' },
    draft: draft(), connection: 'offline', nowMs, ownerObservations: observed({ workspaceId: 'W1', deviceId: 'A', circuitId: 'C1', workspaceRevision: 'W2' }) };
}
function group(): WriteGroup {
  return { ...draft(), id: 'G1', revision: 'R2', applied_revision: 'R1', basic_managed_device_id: 'A',
    status: 'running', members: [{ ...gauge }], migration: {},
    destination: { ...draft().destination, database: 'telemetry' }, created_at: '2026-10-01', updated_at: '2026-10-02' };
}
function runtime(): RuleEvidenceView {
  return { contract_version: 1, workspace_id: 'W1', workspace_revision: 'W2', device_id: 'A',
    observed_at: '2026-10-27T00:00:00Z', freshness: 'fresh', source: null, locator: null,
    groups: [{ group_id: 'G1', saved_revision: 'R2', applied_revision: 'R2', effective_revision: 'R2', as_of: '2026-10-27T00:00:00Z',
      applied_effective_at: null, state: 'effective', proof: 'applied-snapshot-and-layout', saved: null,
      effective: { recording_profile: null, allowed_lateness_seconds: null, binding: { workspace_id: 'W1', workspace_revision: 'W2', entity_revision: 'R2' }, effective_at: '2026-10-26T00:00:00Z', ends_at: null, enabled: true,
        destination: { connector_id: 'DB1', connector_revision: 'D1', alias: null, kind: 'postgres', schema: 'public', table: 'readings' },
        members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', source_value_type: null, column: 'temperature' }], layout: 'managed_wide_profile_v1',
        storage_mode: 'plain', interval_seconds: 30, reduction: 'latest', write_mode: 'append', incomplete_policy: 'partial' },
      delivery: { oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, backlog: [], test_write: null, counter_replay_quarantines: null,
        current: { binding: { workspace_id: 'W1', workspace_revision: 'W2', entity_revision: 'R2' }, intake: { group_revision: 'R2', state: 'active', reason: null },
          totals_scope: 'group', shared: false, stages: { collecting: 1, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
          last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null } } }] };
}
function readiness(): WriteGroupReadiness {
  return { workspace_id: 'W1', workspace_revision: 'W2', group_id: 'G1', group_revision: 'R2', applied_revision: 'R1',
    config_ready: true, schema_ready: true, ready: true, issues: [] };
}
function preview(): SchemaPreviewToken {
  return { token: 'token', operation_id: 'DDL1', digest: 'digest', workspace_id: 'W1', workspace_revision: 'W2',
    plan_id: 'G1', plan_revision: 'R2', connector_id: 'DB1', connector_revision: 'D1', table_prefix: '', statements: [], tables: [],
    expires_at: '2026-10-27T01:00:00Z', created_at: '2026-10-27T00:00:00Z' };
}
function bind<T>(value: T, selected = scope): RecordingWorkflowEvidence<T> {
  return { state: 'available', scope: selected, value: value as RecordingReadonly<T> };
}
function saved(patch: Partial<RecordingWorkflowInput> = {}): Extract<RecordingNextActionInput, { kind: 'saved' }> {
  const source: RecordingWorkflowInput = { scope, savedGroup: group(), readiness: bind(readiness()),
    operation: bind<RecordingWorkflowOperation>({ kind: 'none' }),
    lifecycle: bind({ edit: { state: 'available' as const, reasonCodes: [] } }), ...patch };
  return { kind: 'saved', scope: source.scope, workflow: recordingWorkflowView(source), connection: 'online', nowMs, ownerObservations: observed(source.scope) };
}
function counterStatus(): CandidateCounterStatus {
  return { confirmed: true,
    candidate: { key: 'P1|T1', device_id: 'A', point_id: 'P1', tag_id: 'T1', tag_key: 'count',
      label: 'Count', device_name: 'Device A', address: '1', target_type: 'uint64' },
    measurement: { id: 'M1', workspace_id: 'W1', device_id: 'A', point_id: 'P1', tag_id: 'T1', equipment_id: 'C1',
      definition_revision: 'D1', source_binding_revision: 'B1', series_epoch: 'E1', name: 'Count', quantity: 'energy',
      unit: 'kWh', semantic_kind: 'counter', counter_policy: { allow_negative: false } } };
}
function pending(status: 'pending' | 'running' | 'unknown' | 'partial'): RecordingWorkflowOperation {
  return { kind: 'schema', operation: { operation_id: 'O1', status, executed_statements: 0, created_at: '2026-10-27', updated_at: '2026-10-27' } };
}

describe('guided next safe intent', () => {
  it('keeps an offline empty draft editable without inventing saved group identity', () => {
    const input = offlineDraft();
    expect(recordingNextAction(input)).toEqual({
      intent: { kind: 'choose_sources', scope: { workspaceId: 'W1', deviceId: 'A', circuitId: 'C1', workspaceRevision: 'W2' } },
      enabled: true, reasons: ['sources_empty'], configurationScope: { kind: 'draft' },
    });
    expect(input.scope).not.toHaveProperty('groupId');
    expect(input.draft?.members).toEqual([]);
  });

  it.each(['online', 'offline'] as const)('ordinary non-counter draft reaches save/prepare without counter semantics: %s', (connection) => {
    const value = { ...draft(), members: [{ ...gauge }] };
    const result = recordingNextAction({ ...offlineDraft(), kind: 'draft', draft: value, connection });
    expect(result).toMatchObject({ intent: { kind: 'save_prepare' }, enabled: connection === 'online',
      reasons: connection === 'online' ? ['save_required'] : ['save_required', 'offline'] });
    expect(value.members[0].value_role).toBe('gauge');
  });
  it.each(['counter_raw', 'cumulative'] as const)('directs selected %s to explicit counter confirmation, not start', (mode) => {
    const value = { ...draft(), members: [{ ...gauge, value_role: mode === 'counter_raw' ? 'counter_raw' as const : 'gauge' as const }] };
    if (mode === 'cumulative') value.recording_profile = { kind: 'cumulative', interval: 10, unit: 'minute' };
    expect(recordingNextAction({ ...offlineDraft(), kind: 'draft', draft: value }))
      .toMatchObject({ intent: { kind: 'resolve_semantics' }, enabled: true, reasons: ['counter_confirmation'] });
  });
  it.each(['valid', 'unconfirmed', 'foreign candidate', 'foreign measurement', 'missing proof', 'stale binding',
    'wrong measurement id', 'wrong proof device', 'wrong proof point', 'wrong proof tag'] as const)(
    'consumes existing counter confirmation proof without creating readiness: %s', (caseName) => {
      const status = counterStatus();
      if (caseName === 'unconfirmed') status.confirmed = false;
      if (caseName === 'foreign candidate') status.candidate.tag_id = 'T9';
      if (caseName === 'foreign measurement') status.measurement!.workspace_id = 'W9';
      if (caseName === 'missing proof') delete status.measurement;
      if (caseName === 'stale binding') status.measurement!.source_binding_revision = '';
      if (caseName === 'wrong proof device') status.measurement!.device_id = 'B';
      if (caseName === 'wrong proof point') status.measurement!.point_id = 'P9';
      if (caseName === 'wrong proof tag') status.measurement!.tag_id = 'T9';
      const value = { ...draft(), members: [{ ...gauge, value_role: 'counter_raw' as const,
        measurement_id: caseName === 'wrong measurement id' ? 'M9' : undefined }] };
      const result = recordingNextAction({ ...offlineDraft(), kind: 'draft', draft: value, connection: 'online', counterStatuses: [status] });
      expect(result.intent.kind).toBe(caseName === 'valid' ? 'save_prepare' : 'resolve_semantics');
      expect(result).not.toHaveProperty('ready');
    });
  it('does not reuse saved server readiness for a changed counter draft', () => {
    expect(recordingNextAction({ ...saved(), draft: { ...draft(), members: [{ ...gauge, value_role: 'counter_raw' }] } }).intent.kind)
      .toBe('resolve_semantics');
  });
  it('sends ordinary dirty saved changes back to save/prepare instead of start', () => {
    expect(recordingNextAction({ ...saved(), draft: { ...draft(), members: [{ ...gauge }] } }).intent.kind).toBe('save_prepare');
  });
  it.each(['missing', 'blocked', 'allowed'] as const)('does not offer saving existing role changes without lifecycle edit evidence: %s', (state) => {
    const lifecycle = state === 'missing' ? undefined : bind({ edit: {
      state: state === 'allowed' ? 'available' as const : 'blocked' as const, reasonCodes: [] } });
    const result = recordingNextAction({ ...saved({ lifecycle }), draft: { ...draft(), members: [{ ...gauge }] } });
    expect(result.intent.kind).toBe('save_prepare');
    expect(result.enabled).toBe(state === 'allowed');
    expect(result.reasons).toEqual([state === 'allowed' ? 'save_required' : state === 'blocked' ? 'lifecycle_blocked' : 'lifecycle_unavailable']);
  });
  it('uses canonical final readiness for an unchanged saved counter, without forcing duplicate confirmation', () => {
    const value = group();
    value.members[0].value_role = 'counter_raw';
    expect(recordingNextAction(saved({ savedGroup: value })).intent.kind).toBe('review_scoped_start');
  });
  it.each([true, false])('prioritizes an explicit counter blocker even when config_ready=true and schema_ready=%s', (schema_ready) => {
    const value = group();
    value.members[0].value_role = 'counter_raw';
    for (const counterStatuses of [undefined, [{ ...counterStatus(), confirmed: false }], [counterStatus()]]) {
      const input = saved({ savedGroup: value, readiness: bind({ ...readiness(), config_ready: true, schema_ready, ready: false,
        issues: [{ code: 'counter-measurement-unconfirmed', severity: 'blocking', step: 'Step 4', scope: 'G1', message: 'Counter is unconfirmed' }] }) });
      expect(recordingNextAction({ ...input, counterStatuses })).toMatchObject({
        intent: { kind: 'resolve_semantics' }, reasons: ['counter_confirmation'], enabled: true,
      });
      expect(input.workflow).toMatchObject({ readiness: { value: { config_ready: true, schema_ready, ready: false } } });
    }
  });
  it.each(['counter-semantics-required', 'counter-measurement-not-counter', 'counter-measurement-foreign',
    'counter-measurement-unavailable', 'counter-type-unsupported'])('uses the backend counter repair reason %s', (code) => {
    const value = group();
    value.recording_profile = { kind: 'cumulative', interval: 10, unit: 'minute' };
    expect(recordingNextAction(saved({ savedGroup: value, readiness: bind({ ...readiness(), ready: false,
      issues: [{ code, severity: 'blocking', step: 'Step 4', scope: 'G1', message: 'Counter repair required' }] }) })).intent.kind)
      .toBe('resolve_semantics');
  });
  it.each([
    ['schema-table-unavailable', true, false, 'review_schema', 'schema_conflict'],
    ['connector-unavailable', false, false, 'save_prepare', 'configuration_blocked'],
    ['runtime-unavailable', true, true, 'refresh_evidence', 'readiness_blocked'],
  ] as const)('keeps genuinely confirmed counter on the correct storage/runtime repair path: %s', (code, config_ready, schema_ready, kind, reason) => {
    const value = group();
    value.members[0].value_role = 'counter_raw';
    expect(recordingNextAction({ ...saved({ savedGroup: value, readiness: bind({ ...readiness(), config_ready, schema_ready, ready: false,
      issues: [{ code, severity: 'blocking', step: 'Step 4', scope: 'G1', message: 'Repair storage/runtime' }] }) }),
      counterStatuses: [counterStatus()] })).toMatchObject({ intent: { kind }, reasons: [reason] });
  });
  it('does not treat config/schema readiness alone as counter proof when no confirmation evidence is available', () => {
    const value = group();
    value.members[0].value_role = 'counter_raw';
    expect(recordingNextAction(saved({ savedGroup: value, readiness: bind({ ...readiness(), ready: false }) })).intent.kind)
      .toBe('resolve_semantics');
  });
  it('does not let counter repair hide an unknown operation or override a blocking server counter issue', () => {
    const value = group();
    value.members[0].value_role = 'counter_raw';
    const input = saved({ savedGroup: value, readiness: bind({ ...readiness(), ready: true,
      issues: [{ code: 'counter-measurement-unconfirmed', severity: 'blocking', step: 'Step 4', scope: 'G1', message: 'Still blocked' }] }) });
    expect(recordingNextAction({ ...input, counterStatuses: [counterStatus()] }).intent.kind).toBe('resolve_semantics');
    expect(recordingNextAction(saved({ savedGroup: value, operation: bind(pending('unknown')), readiness: bind({ ...readiness(), ready: false,
      issues: [{ code: 'counter-measurement-unconfirmed', severity: 'blocking', step: 'Step 4', scope: 'G1', message: 'Still blocked' }] }) })).intent.kind)
      .toBe('query_operation');
  });

  it.each(['pending', 'running', 'unknown', 'partial'] as const)('queries O1 %s before any otherwise-ready new start', (status) => {
    expect(recordingNextAction(saved({ operation: bind(pending(status)) })))
      .toMatchObject({ intent: { kind: 'query_operation', operationId: 'O1' }, enabled: true, reasons: ['operation_recovery'] });
  });
  it.each(['failed', 'succeeded'] as const)('retains cleanup recovery even after outer test-write %s', (status) => {
    expect(recordingNextAction(saved({ operation: bind<RecordingWorkflowOperation>({ kind: 'test_write', operation: {
      operation_id: 'TW1', action: 'test_write', status, write_outcome: 'written_verified', cleanup_status: 'failed',
      created_at: '2026-10-27', updated_at: '2026-10-27',
    } }) }))).toMatchObject({ intent: { kind: 'query_operation', operationId: 'TW1' }, reasons: ['operation_recovery'] });
  });
  it('gives retained unconfirmed operation priority over empty offline draft navigation', () => {
    const value = offlineDraft();
    expect(recordingNextAction({ ...value, unconfirmedOperation: { scope: value.scope, operationId: 'O1' } }))
      .toMatchObject({ intent: { kind: 'query_operation', operationId: 'O1' }, enabled: false, reasons: ['operation_recovery', 'offline'] });
  });
  it('retains unknown outcome without inventing an operation ID', () => {
    expect(recordingNextAction(saved({ operation: { state: 'unknown' } })).intent).toEqual({ kind: 'recover_operation', scope });
  });
  it.each(['workspaceId', 'deviceId', 'circuitId', 'groupId'] as const)('never queries a retained operation with foreign %s', (key) => {
    const result = recordingNextAction({ ...saved(), unconfirmedOperation: {
      scope: { ...scope, [key]: 'foreign' }, operationId: 'foreign-secret' } });
    expect(result.intent).toEqual({ kind: 'inspect', scope });
    expect(result.reasons).toEqual(['scope_unavailable']);
  });
  it('retains a same-scope unconfirmed operation even while current group evidence is unavailable', () => {
    expect(recordingNextAction({ ...saved(), workflow: { state: 'unavailable', reason: 'error' },
      unconfirmedOperation: { scope, operationId: 'O1' } }))
      .toMatchObject({ intent: { kind: 'query_operation', operationId: 'O1' }, reasons: ['operation_recovery'] });
  });
  it('keeps retained recovery without an ID and loading operation lookup fail-closed', () => {
    expect(recordingNextAction({ ...saved(), unconfirmedOperation: { scope } }).intent).toEqual({ kind: 'recover_operation', scope });
    expect(recordingNextAction(saved({ operation: { state: 'loading' } })))
      .toMatchObject({ intent: { kind: 'recover_operation' }, enabled: true, reasons: ['evidence_loading'] });
  });

  it.each(['loading', 'error', 'unknown', 'unavailable'] as const)('keeps inspection available while readiness is %s', (state) => {
    expect(recordingNextAction(saved({ readiness: { state } }))).toMatchObject({ intent: { kind: 'refresh_evidence' },
      enabled: state !== 'loading', reasons: [state === 'loading' ? 'evidence_loading' : state === 'error' ? 'evidence_error' : 'evidence_unavailable'] });
  });
  it('keeps configuration blockers distinct from schema preparation', () => {
    expect(recordingNextAction(saved({ readiness: bind({ ...readiness(), config_ready: false, ready: false }) })))
      .toMatchObject({ intent: { kind: 'save_prepare' }, reasons: ['configuration_blocked'] });
  });
  it('does not override the server ready:false even when config/schema booleans are true', () => {
    expect(recordingNextAction(saved({ readiness: bind({ ...readiness(), ready: false }) })))
      .toMatchObject({ intent: { kind: 'refresh_evidence' }, reasons: ['readiness_blocked'] });
  });
  it('requires review when schema evidence is absent, without confirming or creating DDL', () => {
    expect(recordingNextAction(saved({ readiness: bind({ ...readiness(), schema_ready: false, ready: false }) })))
      .toMatchObject({ intent: { kind: 'preview_schema' }, reasons: ['schema_review'] });
  });
  it('does not confirm a conflicting schema even if a preview had been reviewed', () => {
    const value = saved({ readiness: bind({ ...readiness(), schema_ready: false, ready: false,
      issues: [{ code: 'destination-column-type-mismatch', severity: 'blocking', step: 'Step 4', scope: 'G1', message: 'raw database detail' }] }),
      schema: bind(preview()) });
    expect(recordingNextAction({ ...value, schemaReview: { scope, stage: 'reviewed' as const, operationId: 'DDL1' } }))
      .toMatchObject({ intent: { kind: 'review_schema' }, reasons: ['schema_conflict'] });
  });
  it.each(['idle', 'previewed', 'reviewed'] as const)('uses token-free owner schema stage: %s', stage => {
    const result = recordingNextAction({ ...saved({ readiness: bind({ ...readiness(), schema_ready: false, ready: false }) }),
      schemaReview: { scope, stage, ...(stage === 'idle' ? {} : { operationId: 'DDL1' }) } });
    expect(result.intent.kind).toBe(stage === 'idle' ? 'preview_schema' : stage === 'previewed' ? 'review_schema' : 'confirm_schema');
    expect(result.intent).not.toHaveProperty('token'); expect(result.intent).not.toHaveProperty('digest');
  });
  it.each(['missing', 'blocked', 'allowed'] as const)('does not confuse lifecycle edit availability with local start-review navigation: %s', state => {
    const lifecycle = state === 'missing' ? undefined : bind({ edit: {
      state: state === 'allowed' ? 'available' as const : 'blocked' as const, reasonCodes: ['pending_delivery'] } });
    const result = recordingNextAction({ ...saved({ lifecycle }), connection: 'offline' });
    expect(result).toMatchObject({ intent: { kind: 'review_scoped_start' }, enabled: true, reasons: ['start_review'] });
  });
  it('does not equate a running status or applied_revision string with a monitor handoff', () => {
    const result = recordingNextAction(saved());
    expect(result.intent.kind).toBe('review_scoped_start');
    expect(result.intent).not.toHaveProperty('activatedDeviceIds');
  });
  it('does not bypass blocking issues or fabricate an absent lifecycle action', () => {
    expect(recordingNextAction(saved({ lifecycle: bind({}) })).reasons).toEqual(['start_review']);
    expect(recordingNextAction(saved({ readiness: bind({ ...readiness(), issues: [
      { severity: 'blocking', code: 'other', scope: 'G1', step: 'Step 4', message: 'blocked' },
    ] }) })).intent.kind).toBe('refresh_evidence');
  });
  it('allows ordinary legacy groups without a profile and respects non-conflicting schema confirmation issues', () => {
    const value = group();
    delete value.recording_profile;
    expect(recordingNextAction(saved({ savedGroup: value })).intent.kind).toBe('review_scoped_start');
    expect(recordingNextAction({ ...saved({ schema: bind(preview()), readiness: bind({ ...readiness(), schema_ready: false, ready: false,
      issues: [{ severity: 'blocking', code: 'schema-confirmation-required', scope: 'G1', step: 'Step 4', message: 'Review' }] }) }),
      schemaReview: { scope, stage: 'reviewed' as const, operationId: 'DDL1' } }).intent.kind).toBe('confirm_schema');
  });
  it.each(['proven', 'draft newer', 'snapshot missing', 'not activated', 'group not applied',
    'disabled group', 'deleted group', 'draft state', 'server ready', 'failed start', 'offline', 'readiness error'] as const)(
    'hands off only the selected confirmed activated device: %s', (caseName) => {
      const selected = { ...scope, appliedRevision: 'R2' };
      const statuses: Partial<Record<typeof caseName, WriteGroup['status']>> = {
        'disabled group': 'disabled', 'deleted group': 'deleted', 'draft state': 'draft', 'server ready': 'ready',
      };
      const current = { ...group(), applied_revision: 'R2', status: statuses[caseName] ?? 'running' };
      const snapshot = runtime();
      const source = saved({ scope: selected, savedGroup: current,
        readiness: caseName === 'readiness error' ? { state: 'error' } : bind({ ...readiness(), applied_revision: 'R2' }, selected),
        runtimeEvidence: caseName === 'snapshot missing' ? undefined : bind(snapshot, selected),
        lifecycle: bind({ edit: { state: 'available' as const, reasonCodes: [] } }, selected),
        operation: bind<RecordingWorkflowOperation>({ kind: 'start', operation: {
          operation_id: 'START1', action: 'recording_start', status: caseName === 'failed start' ? 'failed' : 'succeeded', intent_digest: 'intent',
          workspace_id: 'W1', setup_revision: 'W2', device_ids: ['A', 'B'],
          groups: [{ group_id: 'G1', group_revision: 'R2', applied_revision: 'R2', saved: true, ready: true, applied: caseName !== 'group not applied' }],
          devices: [{ device_id: 'A', activated: caseName !== 'not activated' }, { device_id: 'B', activated: true }],
          stage: 'complete', created_at: '2026-10-27', updated_at: '2026-10-27',
        } }, selected) });
      const result = recordingNextAction({ ...source, connection: caseName === 'offline' ? 'offline' : 'online',
        ...(caseName === 'draft newer' ? { draft: { ...draft(), members: [{ ...gauge }] } } : {}) });
      const monitorCases = ['proven', 'server ready', 'offline', 'readiness error'];
      expect(result.intent.kind).toBe(monitorCases.includes(caseName) ? 'monitor' : caseName === 'draft newer' ? 'save_prepare' : 'review_scoped_start');
      if (monitorCases.includes(caseName)) expect(result.intent.activatedDeviceIds).toEqual(['A']);
      else expect(result.intent).not.toHaveProperty('activatedDeviceIds');
    });
  it('retains shared-group scope rather than converting it to the selected device', () => {
    const shared = group();
    delete shared.basic_managed_device_id;
    shared.members.push({ ...gauge, device_id: 'B', circuit_id: 'C2', point_id: 'P2', tag_id: 'T2' });
    expect(recordingNextAction(saved({ savedGroup: shared })).configurationScope).toEqual({ kind: 'group', groupId: 'G1' });
  });
  it('fails closed on unavailable or mismatched scope while preserving a local inspect intent', () => {
    expect(recordingNextAction({ ...saved(), workflow: { state: 'unavailable', reason: 'revision_mismatch' } }))
      .toMatchObject({ intent: { kind: 'inspect', scope }, reasons: ['scope_unavailable'], enabled: true });
    expect(recordingNextAction({ ...saved(), scope: { ...scope, groupId: 'other' } }).intent.kind).toBe('inspect');
    expect(recordingNextAction({ ...offlineDraft(), kind: 'draft', draft: { ...draft(), workspace_id: 'W9' } }).intent.kind).toBe('inspect');
  });
  it.each(['workspaceId', 'deviceId', 'circuitId'] as const)('rejects empty draft %s without selecting another scope', (key) => {
    const value = offlineDraft();
    expect(recordingNextAction({ ...value, scope: { ...value.scope, [key]: '' } }).intent.kind).toBe('inspect');
  });
  it('refuses unrelated draft sources but preserves a legitimate multi-device draft containing the selected device', () => {
    const unrelated = { ...draft(), members: [{ ...gauge, device_id: 'B' }] };
    expect(recordingNextAction({ ...offlineDraft(), kind: 'draft', draft: unrelated }).intent.kind).toBe('inspect');
    const shared = { ...draft(), members: [{ ...gauge }, { ...gauge, device_id: 'B' }] };
    expect(recordingNextAction({ ...offlineDraft(), kind: 'draft', draft: shared }).intent.kind).toBe('save_prepare');
    expect(recordingNextAction({ ...saved(), draft: unrelated }).intent.kind).toBe('inspect');
  });
  it('does not report newer saved configuration as running even after a successful historical start', () => {
    const proof = runtime(); const active = proof.groups[0];
    active.applied_revision = 'R1'; active.effective_revision = 'R1'; active.effective!.binding.entity_revision = 'R1';
    active.delivery!.current!.binding.entity_revision = 'R1'; active.delivery!.current!.intake!.group_revision = 'R1';
    const result = recordingNextAction(saved({ runtimeEvidence: bind(proof),
      operation: bind<RecordingWorkflowOperation>({ kind: 'start', operation: {
        operation_id: 'START1', action: 'recording_start', status: 'succeeded', intent_digest: 'intent',
        workspace_id: 'W1', setup_revision: 'W2', device_ids: ['A'], devices: [{ device_id: 'A', activated: true }],
        groups: [{ group_id: 'G1', group_revision: 'R2', applied_revision: 'R1', saved: true, ready: true, applied: true }],
        stage: 'complete', created_at: '2026-10-27', updated_at: '2026-10-27',
      } }) }));
    expect(result.intent.kind).toBe('review_scoped_start');
    expect(result.intent).not.toHaveProperty('activatedDeviceIds');
  });
  it('copies the accepted intent scope so later caller edits cannot retarget it', () => {
    const selected = { ...scope };
    const result = recordingNextAction({ ...saved(), scope: selected });
    selected.deviceId = 'B';
    expect(result.intent.scope.deviceId).toBe('A');
    expect(Object.isFrozen(result.intent.scope)).toBe(true);
    expect(Object.isFrozen(result.reasons)).toBe(true);
    expect(Object.isFrozen(selected)).toBe(false);
  });
  it('permits offline inspection of a preview but not a network request for a missing preview', () => {
    const input = saved({ readiness: bind({ ...readiness(), schema_ready: false, ready: false }) });
    expect(recordingNextAction({ ...input, connection: 'offline' })).toMatchObject({ enabled: false, reasons: ['schema_review', 'offline'] });
    expect(recordingNextAction({ ...saved({ schema: bind(preview()),
      readiness: bind({ ...readiness(), schema_ready: false, ready: false }) }), connection: 'offline', schemaReview: { scope, stage: 'previewed', operationId: 'DDL1' } }))
      .toMatchObject({ intent: { kind: 'review_schema' }, enabled: true, reasons: ['schema_review'] });
  });
});

describe('unmounted next-action presentation', () => {
  it('renders one accessible primary recovery action and performs no action on render', () => {
    const onIntent = vi.fn();
    const { container } = render(<RecordingNextAction input={saved({ operation: bind(pending('unknown')) })}
      labels={labels} scopeLabel="Workspace W1 / Device A / Circuit C1" onIntent={onIntent} />);
    expect(onIntent).not.toHaveBeenCalled();
    const primary = screen.getByRole('button', { name: 'Check existing operation' });
    expect(primary).toHaveAccessibleDescription('Check the existing operation before a new request.');
    expect(container.querySelectorAll('[data-primary="true"]')).toHaveLength(1);
    expect(screen.getByRole('status')).toHaveTextContent(labels.reasons.operation_recovery);
    expect(screen.queryByRole('button', { name: 'Review scoped start' })).not.toBeInTheDocument();
    fireEvent.click(primary);
    expect(onIntent).toHaveBeenCalledExactlyOnceWith({ kind: 'query_operation', scope, operationId: 'O1', operationOwner: 'schema' });
  });
  it('keeps safe inspection accessible while offline, without saving or changing a draft on render/navigation', () => {
    const onIntent = vi.fn();
    const fetchSpy = vi.spyOn(globalThis, 'fetch');
    const storageSpy = vi.spyOn(Storage.prototype, 'setItem');
    const input = { ...offlineDraft(), kind: 'draft' as const, draft: { ...draft(), members: [{ ...gauge }] } };
    const { rerender, unmount } = render(<RecordingNextAction input={input} labels={labels} scopeLabel="Device A" onIntent={onIntent} />);
    const primary = screen.getByRole('button', { name: 'Review and save preparation' });
    expect(primary).toBeDisabled();
    expect(primary).toHaveAccessibleDescription(`${labels.reasons.save_required} ${labels.reasons.offline}`);
    fireEvent.click(primary);
    expect(onIntent).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Inspect this scope' }));
    expect(onIntent).toHaveBeenCalledExactlyOnceWith({ kind: 'inspect', scope: input.scope });
    onIntent.mockClear();
    rerender(<RecordingNextAction input={{ ...input, scope: { ...input.scope, deviceId: 'B' } }} labels={labels} scopeLabel="Device B" onIntent={onIntent} />);
    rerender(<RecordingNextAction input={input} labels={labels} scopeLabel="Device A" onIntent={onIntent} />);
    unmount();
    expect(onIntent).not.toHaveBeenCalled();
    expect(input.draft.recording_profile?.interval).toBe(30);
    expect(fetchSpy).not.toHaveBeenCalled();
    expect(storageSpy).not.toHaveBeenCalled();
  });
  it('uses safe supplied labels for shared scope and schema conflict, never raw backend details', () => {
    const shared = group();
    delete shared.basic_managed_device_id;
    shared.members.push({ ...gauge, device_id: 'B' });
    render(<RecordingNextAction input={saved({ savedGroup: shared, readiness: bind({ ...readiness(), schema_ready: false, ready: false,
      issues: [{ code: 'destination-column-type-mismatch', severity: 'blocking', step: 'Step 4', scope: 'G1', message: 'dsn=password-secret' }] }) })}
      labels={labels} scopeLabel="Group G1 viewed from A" onIntent={vi.fn()} />);
    expect(screen.getByText(/Shared or custom group/)).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveTextContent('Inspect and repair the blocked schema.');
    expect(document.body).not.toHaveTextContent('password-secret');
    expect(screen.getAllByRole('button')).toHaveLength(2);
  });
  it('renders only one inspect action for an unavailable scope', () => {
    render(<RecordingNextAction input={{ ...saved(), workflow: { state: 'unavailable', reason: 'error' } }}
      labels={labels} scopeLabel="Group G1" onIntent={vi.fn()} />);
    expect(screen.getAllByRole('button')).toHaveLength(1);
    expect(screen.getByRole('button', { name: 'Inspect this scope' })).toBeEnabled();
  });
  it('does not confirm DDL offline or perform any action when a reviewed preview renders', () => {
    const onIntent = vi.fn();
    const input = { ...saved({ schema: bind(preview()), readiness: bind({ ...readiness(), schema_ready: false, ready: false }) }),
      schemaReview: { scope, stage: 'reviewed' as const, operationId: 'DDL1' }, connection: 'offline' as const };
    const { rerender } = render(<RecordingNextAction input={input} labels={labels} scopeLabel="Group G1" onIntent={onIntent} />);
    expect(screen.getByRole('button', { name: 'Confirm reviewed schema' })).toBeDisabled();
    expect(onIntent).not.toHaveBeenCalled();
    rerender(<RecordingNextAction input={{ ...input, connection: 'online' }} labels={labels} scopeLabel="Group G1" onIntent={onIntent} />);
    expect(onIntent).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Confirm reviewed schema' }));
    expect(onIntent).toHaveBeenCalledExactlyOnceWith({ kind: 'confirm_schema', scope, operationId: 'DDL1', operationOwner: 'schema' });
  });
});
