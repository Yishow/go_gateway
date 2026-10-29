import type { RecordingOwnerObservations, RecordingOwnerScope } from '@/features/datalink/workbench-v2/state/recordingOwnerObservation';
import { describe, expect, expectTypeOf, it } from 'vitest';
import type { RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import type { WriteGroup, WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import { recordingWorkflowView, recordingNextAction, type RecordingWorkflowEvidence,
  type RecordingReadonly, type RecordingWorkflowInput, type RecordingWorkflowScope } from '@/features/datalink/workbench-v2/state/recordingWorkflowView';

const scope: RecordingWorkflowScope = { workspaceId: 'W1', deviceId: 'A', groupId: 'G1',
  workspaceRevision: 'W3', groupRevision: 'C', appliedRevision: 'B' };
const bind = <T>(value: T, selected = scope): RecordingWorkflowEvidence<T> => ({ state: 'available', scope: selected, value: value as RecordingReadonly<T> });
function saved(): WriteGroup {
  return { id: 'G1', workspace_id: 'W1', revision: 'C', applied_revision: 'B', name: 'Editable saved C', status: 'ready',
    members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', source_revision: 'S3', mapping_revision: 'M3', target_column: 'new_value', required: true }],
    destination: { connector_id: 'DB3', connector_revision: 'D3', database: '/private/saved.sqlite', table_schema: 'public', table_name: 'saved_c', storage_strategy: 'managed' },
    recording_profile: { kind: 'instantaneous', interval: 60, unit: 'second' }, row_policy: { interval_seconds: 60, allowed_lateness_seconds: 1 },
    write_policy: {}, migration: {}, created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-27T10:00:00Z' };
}
function runtime(): RuleEvidenceView {
  const member = { device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', source_value_type: null, column: 'old_value' };
  return { contract_version: 1, workspace_id: 'W1', workspace_revision: 'W3', device_id: 'A',
    observed_at: '2026-10-27T10:00:00.000000001Z', freshness: 'fresh', source: null, locator: null,
    groups: [{ group_id: 'G1', saved_revision: 'C', applied_revision: 'B', effective_revision: 'A',
      applied_effective_at: '2026-10-27T10:00:01Z', as_of: '2026-10-27T10:00:00.000000001Z', state: 'effective', proof: 'applied-snapshot-and-layout', saved: null,
      effective: { recording_profile: null, allowed_lateness_seconds: null, binding: { workspace_id: 'W1', workspace_revision: 'W3', entity_revision: 'A' }, enabled: true,
        effective_at: '2026-10-27T09:00:00Z', ends_at: null,
        destination: { connector_id: 'DB1', connector_revision: 'D1', alias: 'Archive', kind: 'postgres', schema: 'public', table: 'effective_a' },
        members: [member], layout: 'managed_wide_profile_v1', storage_mode: 'plain', interval_seconds: 15, reduction: 'mean', write_mode: 'append', incomplete_policy: 'partial' },
      delivery: { current: { binding: { workspace_id: 'W1', workspace_revision: 'W3', entity_revision: 'A' },
        intake: { group_revision: 'A', state: 'active', reason: null }, totals_scope: 'group', shared: false,
        stages: { collecting: 1, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
        last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null },
        oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, backlog: [], test_write: null, counter_replay_quarantines: null } }] };
}
function observed(selected: RecordingOwnerScope): RecordingOwnerObservations {
  return ['schema', 'start', 'test_write'].map(owner => ({ owner: owner as 'schema' | 'start' | 'test_write', scope: selected, state: 'observed', result: 'none' }));
}
function input(proof = runtime()): RecordingWorkflowInput {
  const readiness: WriteGroupReadiness = { workspace_id: 'W1', workspace_revision: 'W3', group_id: 'G1', group_revision: 'C', applied_revision: 'B',
    config_ready: true, schema_ready: true, ready: true, issues: [] };
  return { scope, savedGroup: saved(), runtimeEvidence: bind(proof), readiness: bind(readiness),
    operation: bind({ kind: 'none' as const }), lifecycle: bind({ edit: { state: 'available' as const, reasonCodes: [] } }) };
}
function available(source = input()) {
  const result = recordingWorkflowView(source);
  expect(result.state).toBe('available');
  if (result.state !== 'available') throw new Error('Expected saved view');
  return result;
}
function next(source: RecordingWorkflowInput) {
  return recordingNextAction({ kind: 'saved', scope: source.scope, workflow: recordingWorkflowView(source), connection: 'online', nowMs: 0, ownerObservations: observed(source.scope) });
}
function activatedB(): RecordingWorkflowInput {
  const selected = { ...scope, groupRevision: 'B' }; const proof = runtime(); const group = proof.groups[0];
  proof.observed_at = '2026-10-27T10:00:01.000000001Z'; group.as_of = '2026-10-27T10:00:01.000000001Z';
  group.saved_revision = 'B'; group.effective_revision = 'B'; group.effective!.binding.entity_revision = 'B';
  group.effective!.effective_at = '2026-10-27T10:00:01Z'; group.effective!.destination.table = 'effective_b'; group.effective!.interval_seconds = 30;
  group.delivery!.current!.binding.entity_revision = 'B'; group.delivery!.current!.intake!.group_revision = 'B';
  return { scope: selected, savedGroup: { ...saved(), revision: 'B' }, runtimeEvidence: bind(proof, selected),
    readiness: bind({ workspace_id: 'W1', workspace_revision: 'W3', group_id: 'G1', group_revision: 'B', applied_revision: 'B',
      config_ready: true, schema_ready: true, ready: true, issues: [] }, selected),
    lifecycle: bind({ edit: { state: 'blocked', reasonCodes: ['already_active'] } }, selected),
    operation: bind({ kind: 'start', operation: { operation_id: 'START_B', action: 'recording_start', status: 'succeeded', intent_digest: 'digest',
      workspace_id: 'W1', setup_revision: 'W3', device_ids: ['A'], devices: [{ device_id: 'A', activated: true }],
      groups: [{ group_id: 'G1', group_revision: 'B', applied_revision: 'B', saved: true, ready: true, applied: true }],
      stage: 'complete', created_at: '2026-10-27T10:00:00Z', updated_at: '2026-10-27T10:00:00Z' } }, selected) };
}
function activeProof(source: RecordingWorkflowInput): RuleEvidenceView {
  if (source.runtimeEvidence?.state !== 'available') throw new Error('Expected runtime fixture');
  return structuredClone(source.runtimeEvidence.value) as RuleEvidenceView;
}

describe('safe as-of effective evidence compensation', () => {
  it('keeps saved C, requested scheduled B and actual effective A separate without reconstructing a WriteGroup', () => {
    const result = recordingWorkflowView(input());
    expect(result).toMatchObject({ state: 'available', saved: { group: { revision: 'C', applied_revision: 'B' } },
      runtime: { state: 'available', value: { saved_revision: 'C', applied_revision: 'B', effective_revision: 'A', applied_effective_at: '2026-10-27T10:00:01Z' } },
      effective: { state: 'available', value: { kind: 'effective', revision: 'A', asOf: '2026-10-27T10:00:00.000000001Z',
        configurationScope: { kind: 'group', groupId: 'G1' }, facts: { interval_seconds: 15, destination: { table: 'effective_a' } } } } });
    expect(result).not.toHaveProperty('applied');
  });
  it('does not promote a succeeded scheduled B start into an active B monitor handoff before its boundary', () => {
    const selected = { ...scope, groupRevision: 'B' }; const source = input(); const proof = runtime(); proof.groups[0].saved_revision = 'B';
    const value = recordingWorkflowView({ ...source, scope: selected, savedGroup: { ...saved(), revision: 'B' }, runtimeEvidence: bind(proof, selected),
      readiness: bind({ workspace_id: 'W1', workspace_revision: 'W3', group_id: 'G1', group_revision: 'B', applied_revision: 'B',
        config_ready: true, schema_ready: true, ready: true, issues: [] }, selected),
      lifecycle: bind({ edit: { state: 'blocked', reasonCodes: ['scheduled'] } }, selected),
      operation: bind({ kind: 'start', operation: { operation_id: 'START_B', action: 'recording_start', status: 'succeeded', intent_digest: 'digest',
        workspace_id: 'W1', setup_revision: 'W3', device_ids: ['A'], devices: [{ device_id: 'A', activated: true }],
        groups: [{ group_id: 'G1', group_revision: 'B', applied_revision: 'B', saved: true, ready: true, applied: true }],
        stage: 'complete', created_at: '2026-10-27T10:00:00Z', updated_at: '2026-10-27T10:00:00Z' } }, selected) });
    expect(value).toMatchObject({ effective: { state: 'available', value: { revision: 'A' } } });
    const next = recordingNextAction({ kind: 'saved', scope: selected, workflow: value, connection: 'online', nowMs: 0, ownerObservations: observed(selected) });
    expect(next).toMatchObject({ intent: { kind: 'review_scoped_start' }, enabled: true, reasons: ['start_review'] });
    expect(next.intent).not.toHaveProperty('activatedDeviceIds');
  });
  it('uses new server as-of evidence after B becomes effective, independent of the browser clock', () => {
    const source = activatedB(); const view = available(source);
    expect(view.effective).toMatchObject({ state: 'available', value: { revision: 'B', asOf: '2026-10-27T10:00:01.000000001Z',
      facts: { effective_at: '2026-10-27T10:00:01Z', interval_seconds: 30, destination: { table: 'effective_b' } } } });
    for (const nowMs of [0, NaN, Date.parse('2030-01-01')]) {
      expect(recordingNextAction({ kind: 'saved', scope: source.scope, workflow: view, connection: 'online', nowMs, ownerObservations: observed(source.scope) }))
        .toMatchObject({ intent: { kind: 'monitor', activatedDeviceIds: ['A'] }, enabled: true });
    }
    expect(view.runtime).toMatchObject({ value: { delivery: { current: { stages: { sql_committed: 0 }, last_commit: null, next_due_at: null } } } });
  });
  it('uses updated immutable facts across the boundary while saved C remains distinct from effective B', () => {
    const proof = runtime(); const group = proof.groups[0];
    proof.observed_at = '2026-10-27T10:00:01Z'; group.as_of = '2026-10-27T10:00:01Z'; group.effective_revision = 'B';
    group.effective!.binding.entity_revision = 'B'; group.effective!.effective_at = '2026-10-27T10:00:01Z';
    group.effective!.destination.table = 'effective_b'; group.effective!.interval_seconds = 30;
    group.delivery!.current!.binding.entity_revision = 'B'; group.delivery!.current!.intake!.group_revision = 'B';
    expect(available(input(proof))).toMatchObject({ saved: { group: { revision: 'C', destination: { table_name: 'saved_c' } } },
      effective: { value: { revision: 'B', facts: { interval_seconds: 30, destination: { table: 'effective_b' } } } } });
    expect(next(input(proof)).intent.kind).not.toBe('monitor');
  });
  it.each(['scope', 'saved', 'requested', 'binding', 'as_of', 'future_snapshot', 'intake_revision'] as const)(
    'refuses stale/mixed runtime proof: %s', (kind) => {
      const proof = runtime(); const group = proof.groups[0];
      if (kind === 'scope') proof.workspace_revision = 'W2';
      if (kind === 'saved') group.saved_revision = 'OLD';
      if (kind === 'requested') group.applied_revision = 'OTHER';
      if (kind === 'binding') group.effective!.binding.entity_revision = 'B';
      if (kind === 'as_of') group.as_of = '2026-10-27T10:00:00.000000002Z';
      if (kind === 'future_snapshot') group.effective!.effective_at = '2026-10-27T10:00:00.000000002Z';
      if (kind === 'intake_revision') group.delivery!.current!.intake!.group_revision = 'B';
      const view = available(input(proof));
      expect(view.effective).toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
      expect(view.saved.group.revision).toBe('C');
    });
  it.each(['effective_revision', 'as_of', 'effective', 'effective_at', 'layout', 'interval', 'proof'] as const)(
    'retains explicit safe backlog but no effective/current facts with missing %s', (key) => {
      const proof = runtime(); const group = proof.groups[0];
      group.delivery!.backlog = [{ binding: { workspace_id: 'W1', workspace_revision: 'W3', entity_revision: 'OLD' },
        destination: { connector_id: 'DB0', connector_revision: 'D0', alias: null, kind: 'postgres', schema: 'public', table: 'old_backlog' }, pending: 2, reasons: ['none'] }];
      if (key === 'effective_revision') group.effective_revision = null;
      if (key === 'as_of') group.as_of = null;
      if (key === 'effective') group.effective = null;
      if (key === 'effective_at') group.effective!.effective_at = null;
      if (key === 'layout') group.effective!.layout = 'unavailable';
      if (key === 'interval') group.effective!.interval_seconds = null;
      if (key === 'proof') group.proof = 'unavailable';
      const view = available(input(proof));
      expect(view.effective).toEqual({ state: 'unavailable', reason: 'effective_snapshot_unavailable' });
      expect(view.runtime).toMatchObject({ state: 'available', value: { state: 'unavailable', effective: null,
        delivery: { current: null, backlog: [{ binding: { entity_revision: 'OLD' }, pending: 2 }] } } });
    });
  it.each(['stale', 'unavailable'] as const)('does not treat %s observation as fresh authority', (freshness) => {
    const proof = runtime(); proof.freshness = freshness;
    expect(available(input(proof)).effective).toEqual({ state: 'unavailable', reason: freshness === 'stale' ? 'stale' : 'effective_snapshot_unavailable' });
  });
  it('does not infer effective evidence from ready booleans, saved status, or requested revision when no runtime proof exists', () => {
    const source = { ...input() }; delete source.runtimeEvidence;
    const view = available(source);
    expect(view.effective).toEqual({ state: 'unavailable', reason: 'effective_snapshot_unavailable' });
    expect(view.runtime).toEqual({ state: 'unavailable', reason: 'missing' });
    expect(view).not.toHaveProperty('intake'); expect(view).not.toHaveProperty('ready');
  });
  it.each(['ended', 'disabled', 'tombstone'] as const)('retains %s snapshot for history without granting intake or start', (kind) => {
    const source = { ...activatedB() }; const proof = activeProof(source); const group = proof.groups[0];
    if (kind === 'ended') group.effective!.ends_at = '2026-10-27T10:00:01.000000001Z';
    else group.effective!.enabled = false;
    if (kind === 'tombstone') source.savedGroup = { ...source.savedGroup, status: 'deleted' };
    group.delivery!.current!.intake!.state = 'retiring';
    const history = { ...source, runtimeEvidence: bind(proof, source.scope), lifecycle: undefined };
    expect(available(history).effective.state).toBe('available');
    expect(next(history)).toMatchObject({ intent: { kind: 'review_scoped_start' }, enabled: true, reasons: ['start_review'] });
    expect(next(history).intent).not.toHaveProperty('activatedDeviceIds');
  });
  it.each(['ended', 'disabled'] as const)('rejects contradictory active intake on an %s retained snapshot', (kind) => {
    const proof = runtime();
    if (kind === 'ended') proof.groups[0].effective!.ends_at = proof.observed_at;
    else proof.groups[0].effective!.enabled = false;
    expect(available(input(proof)).effective).toEqual({ state: 'unavailable', reason: 'invalid_evidence' });
  });
  it.each(['missing delivery', 'missing current', 'missing intake', 'not_running', 'blocked', 'retiring', 'unavailable'] as const)(
    'requires independent active intake for monitor, not succeeded start alone: %s', (kind) => {
      const source = { ...activatedB() }; const proof = activeProof(source); const group = proof.groups[0];
      if (kind === 'missing delivery') group.delivery = null;
      else if (kind === 'missing current') group.delivery!.current = null;
      else if (kind === 'missing intake') group.delivery!.current!.intake = null;
      else group.delivery!.current!.intake!.state = kind;
      expect(next({ ...source, runtimeEvidence: bind(proof, source.scope) }))
        .toMatchObject({ intent: { kind: 'review_scoped_start' }, enabled: true, reasons: ['start_review'] });
    });
  it.each(['unknown schema', 'unknown start', 'unknown test write', 'cleanup failed'] as const)(
    'retains unresolved-operation priority even with current effective active intake: %s', (kind) => {
      const source = { ...activatedB() }; const selected = source.scope;
      if (kind === 'unknown start') {
        if (source.operation?.state !== 'available' || source.operation.value.kind !== 'start') throw new Error('Expected start fixture');
        source.operation = bind({ kind: 'start', operation: { ...source.operation.value.operation, status: 'unknown' } }, selected);
      } else if (kind === 'unknown schema') source.operation = bind({ kind: 'schema', operation: {
        operation_id: 'O_UNKNOWN', status: 'unknown', executed_statements: 0, created_at: '', updated_at: '' } }, selected);
      else source.operation = bind({ kind: 'test_write', operation: { operation_id: 'O_UNKNOWN', action: 'test_write',
        status: kind === 'cleanup failed' ? 'succeeded' : 'unknown', write_outcome: 'written_verified', cleanup_status: 'failed', created_at: '', updated_at: '' } }, selected);
      expect(next(source)).toMatchObject({ intent: { kind: 'query_operation' }, reasons: ['operation_recovery'] });
      expect(available(source).lifecycle).toEqual({ state: 'unavailable', reason: 'operation_unresolved' });
    });
  it.each(['C1', null, 'C2'] as const)('requires exact frozen circuit proof when saved C removed A/C1: %s', (circuitId) => {
    const selected = { ...scope, circuitId: 'C1' }; const source = input(); const proof = runtime();
    proof.groups[0].effective!.members[0].circuit_id = circuitId;
    const value = recordingWorkflowView({ ...source, scope: selected, savedGroup: { ...saved(), members: [
      { device_id: 'B', circuit_id: 'C2', point_id: 'P2', tag_id: 'T2', source_revision: 'S2', mapping_revision: 'M2', target_column: 'b', required: true }] },
      runtimeEvidence: bind(proof, selected) });
    if (circuitId !== 'C1') expect(value).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
    else expect(value).toMatchObject({ state: 'available', effective: { value: { configurationScope: { kind: 'group', groupId: 'G1' },
      facts: { members: [{ device_id: 'A', circuit_id: 'C1' }] } } }, saved: { group: { members: [{ device_id: 'B', circuit_id: 'C2' }] } } });
  });
  it('does not borrow an effective managed-device role or private destination data from the newer editable saved group', () => {
    const view = available({ ...input(), savedGroup: { ...saved(), basic_managed_device_id: 'A' } });
    expect(view.saved.configurationScope).toEqual({ kind: 'managed_device', deviceId: 'A' });
    expect(view.effective).toMatchObject({ value: { configurationScope: { kind: 'group', groupId: 'G1' } } });
    expect(JSON.stringify(view.effective)).not.toMatch(/database|private|saved_c|basic_managed_device_id|applied_revision/);
    expectTypeOf<Extract<keyof RecordingWorkflowInput, 'appliedSnapshot'>>().toEqualTypeOf<never>();
  });
  it('keeps accepted-input quarantine separate from finalized delivery and copied immutable evidence', () => {
    const proof = runtime(); proof.groups[0].delivery!.counter_replay_quarantines = [{ workspace_id: 'W1', group_id: 'G1', group_revision: 'OLD',
      pending_samples: 3, reason: 'counter-replay-proof-missing', action: 'preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance' }];
    const view = available(input(proof));
    proof.groups[0].effective!.destination.table = 'changed'; proof.groups[0].delivery!.counter_replay_quarantines![0].pending_samples = 99;
    expect(view.effective).toMatchObject({ value: { facts: { destination: { table: 'effective_a' } } } });
    expect(view.runtime).toMatchObject({ value: { delivery: { counter_replay_quarantines: [{ pending_samples: 3 }],
      current: { stages: { queued: 0, sql_committed: 0 }, last_commit: null }, backlog: [] } } });
    if (view.effective.state !== 'available' || view.runtime.state !== 'available') throw new Error('Expected proof');
    expect(Object.isFrozen(view.effective.value.facts.members[0])).toBe(true);
    expect(Object.isFrozen(view.runtime.value.delivery?.counter_replay_quarantines)).toBe(true);
    expect(Object.isFrozen(proof)).toBe(false);
  });
});
