import { describe, expect, expectTypeOf, it, vi } from 'vitest';
import type { RuleEvidenceView, RecordingGroupFacts } from '@/types/runtimeRuleEvidence';
import type { WriteGroup, WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import type { SchemaPreviewToken } from '@/types/recordingPlan';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroupTestWriteOperation } from '@/types/studioV2WriteGroupTestWrite';
import {
  recordingWorkflowView,
  type EffectiveRecordingFacts, type RecordingWorkflowEvidence, type RecordingWorkflowInput,
  type RecordingWorkflowOperation, type RecordingWorkflowScope, type SavedRecordingFacts,
} from '@/features/datalink/workbench-v2/state/recordingWorkflowView';

const scope: RecordingWorkflowScope = {
  workspaceId: 'W1', deviceId: 'A', circuitId: 'C1', groupId: 'G1',
  workspaceRevision: 'W2', groupRevision: 'R2', appliedRevision: 'R1',
};

function draft(): WriteGroup {
  return {
    id: 'G1', workspace_id: 'W1', revision: 'R2', applied_revision: 'R1', name: 'Recording', status: 'running',
    basic_managed_device_id: 'A',
    members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1',
      source_revision: 'S2', mapping_revision: 'M2', target_column: 'temperature', required: true }],
    destination: { connector_id: 'DB1', connector_revision: 'D2', database: 'telemetry',
      table_schema: 'public', table_name: 'new_readings', storage_strategy: 'managed' },
    recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' },
    row_policy: { interval_seconds: 30, allowed_lateness_seconds: 2, group_key_columns: ['device'] },
    write_policy: {}, migration: {}, created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-02T00:00:00Z',
  };
}

function runtime(): RuleEvidenceView {
  return { contract_version: 1, workspace_id: 'W1', workspace_revision: 'W2', device_id: 'A',
    observed_at: '2026-10-02T00:00:00Z', freshness: 'fresh', source: null, locator: null,
    groups: [{ group_id: 'G1', saved_revision: 'R2', applied_revision: 'R1', effective_revision: 'R1',
      as_of: '2026-10-02T00:00:00Z', applied_effective_at: null, state: 'effective', proof: 'applied-snapshot-and-layout', saved: null,
      effective: { recording_profile: null, allowed_lateness_seconds: null, binding: { workspace_id: 'W1', workspace_revision: 'W2', entity_revision: 'R1' },
        effective_at: '2026-10-01T00:00:00Z', ends_at: null, enabled: true,
        members: [{ device_id: 'A', circuit_id: 'C1', point_id: 'P1', tag_id: 'T1', source_value_type: null, column: 'temperature' }],
        destination: { connector_id: 'DB1', connector_revision: 'D1', alias: 'Archive', kind: 'postgres', schema: 'public', table: 'old_readings' },
        layout: 'managed_wide_profile_v1', storage_mode: 'plain', interval_seconds: 15, reduction: 'latest', write_mode: 'append', incomplete_policy: 'partial' }, delivery: null }] };
}

function sharedGroup(): WriteGroup {
  return { id: 'G1', workspace_id: 'W1', revision: 'R2', applied_revision: 'R1', name: 'Shared custom group', status: 'running',
    members: [
      { device_id: 'A', circuit_id: 'C1', point_id: 'PA', tag_id: 'TA', source_revision: 'SA', mapping_revision: 'MA', target_column: 'a', required: true },
      { device_id: 'B', circuit_id: 'C2', point_id: 'PB', tag_id: 'TB', source_revision: 'SB', mapping_revision: 'MB', target_column: 'b', required: true },
    ],
    destination: { connector_id: 'DB1', connector_revision: 'D2', database: 'telemetry',
      table_schema: 'public', table_name: 'shared_readings', storage_strategy: 'custom' },
    row_policy: { interval_seconds: 30, allowed_lateness_seconds: 2 }, write_policy: {}, migration: {},
    created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-02T00:00:00Z' };
}

function evidence<T>(value: T): RecordingWorkflowEvidence<T> {
  return { state: 'available', scope: { ...scope }, value } as RecordingWorkflowEvidence<T>;
}
function input(): RecordingWorkflowInput {
  return { scope: { ...scope }, savedGroup: draft(), runtimeEvidence: evidence(runtime()) };
}
function readiness(): WriteGroupReadiness {
  return { workspace_id: 'W1', workspace_revision: 'W2', group_id: 'G1', group_revision: 'R2',
    applied_revision: 'R1', config_ready: true, schema_ready: true, ready: false, issues: [] };
}
function schema(): SchemaPreviewToken {
  return { token: 'preview-token', operation_id: 'DDL1', workspace_id: 'W1', workspace_revision: 'W2',
    plan_id: 'G1', plan_revision: 'R2', connector_id: 'DB1', connector_revision: 'D2',
    table_prefix: '', statements: [], tables: [], digest: 'digest', created_at: '2026-10-02T00:00:00Z',
    expires_at: '2026-10-02T01:00:00Z' };
}
function operation(status: 'pending' | 'running' | 'succeeded' | 'partial' | 'failed' | 'unknown'): RecordingWorkflowOperation {
  return { kind: 'schema', operation: { operation_id: 'O1', status, executed_statements: 0,
    created_at: '2026-10-02T00:00:00Z', updated_at: '2026-10-02T00:00:00Z' } };
}
function start(): RecordingStartOperation {
  return { operation_id: 'O1', action: 'recording_start', status: 'unknown', intent_digest: 'intent',
    workspace_id: 'W1', setup_revision: 'W2', device_ids: ['A'], devices: [{ device_id: 'A', activated: false }],
    groups: [{ group_id: 'G1', group_revision: 'R2', applied_revision: 'R1', saved: true, ready: true, applied: true }],
    stage: 'activation', created_at: '2026-10-02T00:00:00Z', updated_at: '2026-10-02T00:00:00Z' };
}
function withLifecycle(): RecordingWorkflowInput {
  return { ...input(), operation: evidence<RecordingWorkflowOperation>({ kind: 'none' }),
    lifecycle: evidence({ edit: { state: 'available' as const, reasonCodes: [] },
      remove: { state: 'blocked' as const, reasonCodes: ['pending_delivery'] } }) };
}
function available(value = input()) {
  const view = recordingWorkflowView(value);
  expect(view.state).toBe('available');
  if (view.state !== 'available') throw new Error('Expected available view');
  return view;
}

describe('isolated guided recording caller contract', () => {
  it('keeps saved R2/30 seconds separate from proven applied R1/15 seconds and its old destination', () => {
    const view = available();
    expect(view.saved).toMatchObject({ kind: 'saved', group: { revision: 'R2',
      recording_profile: { interval: 30, unit: 'second' }, destination: { table_name: 'new_readings' } } });
    expect(view.effective).toMatchObject({ state: 'available', value: { kind: 'effective', revision: 'R1',
      facts: { interval_seconds: 15, destination: { table: 'old_readings' } } } });
    if (view.effective.state !== 'available') throw new Error('Expected effective');
    expect(view.effective.value.facts.destination).not.toHaveProperty('database');
  });

  it.each([undefined, { state: 'unavailable' } as const])('does not invent applied facts without a snapshot: %j', (snapshot) => {
    expect(available({ ...input(), runtimeEvidence: snapshot }).effective)
      .toEqual({ state: 'unavailable', reason: 'effective_snapshot_unavailable' });
  });

  it('does not use a draft as a snapshot even when it advertises applied_revision R1', () => {
    expect(available({ ...input(), runtimeEvidence: evidence(draft()) as unknown as RecordingWorkflowEvidence<RuleEvidenceView> }).effective)
      .toEqual({ state: 'unavailable', reason: 'invalid_evidence' });
  });

  it('requires a separate applied snapshot even when saved and applied revision strings match', () => {
    expect(available({ scope: { ...scope, groupRevision: 'R1' }, savedGroup: { ...draft(), revision: 'R1' } }).effective)
      .toEqual({ state: 'unavailable', reason: 'effective_snapshot_unavailable' });
  });

  it.each(['workspace_id', 'id', 'basic_managed_device_id'] as const)('refuses same-name saved groups with a foreign %s', (key) => {
    expect(recordingWorkflowView({ ...input(), savedGroup: { ...draft(), [key]: 'foreign' } }))
      .toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
  });
  it.each(['device_id', 'circuit_id'] as const)('validates member %s rather than display names', (key) => {
    const group = draft();
    group.members[0][key] = 'foreign';
    expect(recordingWorkflowView({ ...input(), savedGroup: group, runtimeEvidence: undefined }))
      .toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
  });

  it.each([['A', 'C1'], ['B', 'C2']] as const)('retains the full custom A+B group viewed from %s/%s without exclusive ownership', (deviceId, circuitId) => {
    const view = available({ ...input(), scope: { ...scope, deviceId, circuitId },
      savedGroup: sharedGroup(), runtimeEvidence: undefined });
    expect(view.saved).toMatchObject({ configurationScope: { kind: 'group', groupId: 'G1' } });
    expect(view.scope).toMatchObject({ deviceId, circuitId });
    expect(view.saved.group.members.map((member) => member.device_id)).toEqual(['A', 'B']);
    expect(view.saved.group.destination.storage_strategy).toBe('custom');
    expect(view.saved.group).not.toHaveProperty('basic_managed_device_id');
  });
  it.each(['A', 'B'])('does not treat a shared group as exclusive when basic_managed_device_id is %s', (owner) => {
    expect(recordingWorkflowView({ ...input(), savedGroup: { ...sharedGroup(), basic_managed_device_id: owner } }))
      .toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
  });
  it('labels an actual managed-device group explicitly rather than inferring ownership from selection', () => {
    expect(available().saved).toMatchObject({ configurationScope: { kind: 'managed_device', deviceId: 'A' } });
  });
  it.each(['valid', 'foreign workspace', 'foreign group', 'wrong revision', 'missing', 'unrelated device'] as const)(
    'retains selected D from applied R1 when saved R2 removes D for E: %s snapshot', (caseName) => {
      const selected = { ...scope, deviceId: 'D' };
      const saved = sharedGroup();
      saved.members = [{ device_id: 'E', circuit_id: 'C2', point_id: 'PE', tag_id: 'TE',
        source_revision: 'SE', mapping_revision: 'ME', target_column: 'e', required: true }];
      const snapshot = runtime();
      snapshot.device_id = 'D';
      snapshot.groups[0].effective!.members = [{ device_id: 'D', circuit_id: 'C1', point_id: 'PD', tag_id: 'TD', source_value_type: null, column: 'd' }];
      if (caseName === 'foreign workspace') snapshot.workspace_id = 'W9';
      if (caseName === 'foreign group') snapshot.groups[0].group_id = 'G9';
      if (caseName === 'wrong revision') snapshot.groups[0].effective!.binding.entity_revision = 'R0';
      if (caseName === 'unrelated device') snapshot.groups[0].effective!.members[0].device_id = 'X';
      const result = recordingWorkflowView({ scope: selected, savedGroup: saved,
        runtimeEvidence: caseName === 'missing' ? undefined : { state: 'available', scope: selected, value: snapshot } });
      if (caseName !== 'valid') {
        expect(result).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
        return;
      }
      expect(result).toMatchObject({ state: 'available', scope: { deviceId: 'D' },
        saved: { configurationScope: { kind: 'group', groupId: 'G1' }, group: { members: [{ device_id: 'E' }] } },
        effective: { state: 'available', value: { configurationScope: { kind: 'group', groupId: 'G1' }, facts: { members: [{ device_id: 'D' }] } } } });
    });
  it.each(['workspaceId', 'deviceId', 'groupId', 'circuitId'] as const)('rejects empty scope %s', (key) => {
    expect(recordingWorkflowView({ ...input(), scope: { ...scope, [key]: '' } }))
      .toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
  });
  it.each(['workspaceRevision', 'groupRevision'] as const)('rejects absent %s', (key) => {
    expect(recordingWorkflowView({ ...input(), scope: { ...scope, [key]: '' } }))
      .toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
  });
  it.each(['revision', 'applied_revision'] as const)('refuses a stale saved %s', (key) => {
    expect(recordingWorkflowView({ ...input(), savedGroup: { ...draft(), [key]: 'old' } }))
      .toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
  });

  it('accepts device-only empty managed groups without fabricating sources or lifecycle availability', () => {
    const group = draft();
    group.members = [];
    const view = available({ scope: { ...scope, circuitId: undefined }, savedGroup: group });
    expect(view.saved.group.members).toEqual([]);
    expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'missing' });
    expect(recordingWorkflowView({ scope: { ...scope, circuitId: undefined },
      savedGroup: { ...group, basic_managed_device_id: undefined } })).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
    expect(recordingWorkflowView({ scope, savedGroup: group })).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
  });

  it('does not infer readiness or lifecycle from running/applied/config/schema booleans', () => {
    const view = available({ ...input(), readiness: evidence(readiness()) });
    expect(view.readiness).toEqual({ state: 'available', value: { workspace_id: 'W1', workspace_revision: 'W2',
      group_id: 'G1', group_revision: 'R2', applied_revision: 'R1',
      config_ready: true, schema_ready: true, ready: false, issues: [] } });
    expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'missing' });
    expect(view).not.toHaveProperty('ready');
    expect(view).not.toHaveProperty('canStart');
    expect(view).not.toHaveProperty('canDelete');
  });
  it('preserves schema evidence without inventing DDL, delivery or readback success', () => {
    const view = available({ ...input(), schema: evidence(schema()) });
    expect(view.schema).toMatchObject({ state: 'available', value: { statements: [], operation_id: 'DDL1' } });
    for (const field of ['sql_committed', 'readback_verified', 'cleanup_completed', 'next_due', 'metrics', 'primaryAction']) {
      expect(view).not.toHaveProperty(field);
    }
  });

  it.each(['loading', 'error', 'unknown', 'unavailable'] as const)('preserves distinct evidence degradation: %s', (state) => {
    const missing = { state: 'unavailable', reason: state === 'unavailable' ? 'missing' : state };
    const view = available({ ...input(), readiness: { state }, schema: { state }, operation: { state }, lifecycle: { state } });
    expect(view.readiness).toEqual(missing);
    expect(view.schema).toEqual(missing);
    expect(view.operation).toEqual(missing);
    expect(view.lifecycle).toEqual(missing);
    expect(available({ ...input(), runtimeEvidence: { state } }).effective).toEqual({ ...missing,
      reason: state === 'unavailable' ? 'effective_snapshot_unavailable' : state });
  });

  it.each(['workspaceId', 'deviceId', 'circuitId', 'groupId', 'workspaceRevision', 'groupRevision', 'appliedRevision'] as const)(
    'refuses foreign or stale evidence bound by %s', (key) => {
      const foreignScope = { ...scope, [key]: 'foreign' };
      const reason = key.endsWith('Revision') ? 'revision_mismatch' : 'scope_mismatch';
      const view = available({ ...withLifecycle(),
        runtimeEvidence: { state: 'available', scope: foreignScope, value: runtime() },
        readiness: { state: 'available', scope: foreignScope, value: readiness() },
        schema: { state: 'available', scope: foreignScope, value: schema() },
        operation: { state: 'available', scope: foreignScope, value: operation('unknown') },
        lifecycle: { state: 'available', scope: foreignScope, value: {} },
      });
      for (const field of ['effective', 'readiness', 'schema', 'operation', 'lifecycle'] as const) {
        expect(view[field]).toEqual({ state: 'unavailable', reason });
      }
    });

  it.each(['workspace_id', 'group_id', 'workspace_revision', 'group_revision', 'applied_revision'] as const)(
    'checks canonical readiness %s inside the envelope', (key) => {
      const view = available({ ...input(), readiness: evidence({ ...readiness(), [key]: 'foreign' }) });
      expect(view.readiness).toEqual({ state: 'unavailable', reason: key.endsWith('_id') ? 'scope_mismatch' : 'revision_mismatch' });
    });
  it.each(['workspace_id', 'plan_id', 'connector_id', 'workspace_revision', 'plan_revision', 'connector_revision'] as const)(
    'checks canonical schema %s inside the envelope', (key) => {
      expect(available({ ...input(), schema: evidence({ ...schema(), [key]: 'foreign' }) }).schema)
        .toEqual({ state: 'unavailable', reason: key.endsWith('_id') ? 'scope_mismatch' : 'revision_mismatch' });
    });

  it.each(['pending', 'running', 'partial', 'unknown'] as const)('retains operation O1 %s and blocks all new-action eligibility', (status) => {
    const view = available({ ...withLifecycle(), operation: evidence(operation(status)) });
    expect(view.operation).toMatchObject({ state: 'available', value: { kind: 'schema', operation: { operation_id: 'O1', status } } });
    expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'operation_unresolved' });
  });
  it.each([undefined, { state: 'unknown' } as const])('does not turn unqueried operation evidence into idle: %j', (operationEvidence) => {
    expect(available({ ...withLifecycle(), operation: operationEvidence }).lifecycle)
      .toEqual({ state: 'unavailable', reason: 'operation_unresolved' });
  });
  it.each(['succeeded', 'failed'] as const)('only forwards server lifecycle availability after a terminal %s result', (status) => {
    const view = available({ ...withLifecycle(), operation: evidence(operation(status)) });
    expect(view.lifecycle).toEqual({ state: 'available', value: {
      edit: { state: 'available', reasonCodes: [] },
      remove: { state: 'blocked', reasonCodes: ['pending_delivery'] },
    } });
  });
  it('never fills missing actions when the server supplies an empty availability record', () => {
    expect(available({ ...withLifecycle(), lifecycle: evidence({}) }).lifecycle).toEqual({ state: 'available', value: {} });
  });

  it('retains unknown start O1 for device A without treating progress booleans as applied facts', () => {
    const view = available({ ...withLifecycle(), runtimeEvidence: undefined,
      operation: evidence({ kind: 'start', operation: start() }) });
    expect(view.operation).toMatchObject({ state: 'available', value: { kind: 'start',
      operation: { operation_id: 'O1', status: 'unknown', device_ids: ['A'] } } });
    expect(view.effective).toEqual({ state: 'unavailable', reason: 'effective_snapshot_unavailable' });
    expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'operation_unresolved' });
  });
  it.each([
    ['workspace', { workspace_id: 'W2' }, 'scope_mismatch'],
    ['device', { device_ids: ['B'] }, 'scope_mismatch'],
    ['group', { groups: [] }, 'scope_mismatch'],
    ['workspace revision', { setup_revision: 'old' }, 'revision_mismatch'],
    ['group revision', { groups: [{ group_id: 'G1', group_revision: 'old', saved: true, ready: true, applied: true }] }, 'revision_mismatch'],
    ['applied revision', { groups: [{ group_id: 'G1', group_revision: 'R2', applied_revision: 'old', saved: true, ready: true, applied: true }] }, 'revision_mismatch'],
  ] as const)('validates start %s inside the caller envelope', (_label, patch, reason) => {
    const invalid = { ...start(), ...patch };
    const view = available({ ...withLifecycle(), operation: evidence({ kind: 'start', operation: invalid }) });
    expect(view.operation).toEqual({ state: 'unavailable', reason });
    expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'operation_unresolved' });
  });
  it('allows start progress without applied_revision but never manufactures that missing fact', () => {
    const pending = start();
    delete pending.groups[0].applied_revision;
    const view = available({ ...input(), operation: evidence({ kind: 'start', operation: pending }) });
    expect(view.operation).toMatchObject({ state: 'available', value: { operation: { status: 'unknown' } } });
    if (view.operation.state !== 'available' || view.operation.value.kind !== 'start') throw new Error('Expected start');
    expect(view.operation.value.operation.groups[0]).not.toHaveProperty('applied_revision');
  });
  it('retains unknown test-write independently of write and cleanup outcomes', () => {
    const view = available({ ...withLifecycle(), operation: evidence({ kind: 'test_write', operation: {
      operation_id: 'T1', action: 'test_write', status: 'unknown', write_outcome: 'written_verified', cleanup_status: 'unknown',
      created_at: '2026-10-02T00:00:00Z', updated_at: '2026-10-02T00:00:00Z',
    } }) });
    expect(view.operation).toMatchObject({ state: 'available', value: { operation: {
      status: 'unknown', write_outcome: 'written_verified', cleanup_status: 'unknown',
    } } });
    expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'operation_unresolved' });
  });
  it.each(['failed', 'succeeded'] as const)('does not use outer test-write %s to erase unresolved nested results', (status) => {
    const unresolved: Array<Pick<WriteGroupTestWriteOperation, 'write_outcome' | 'cleanup_status'>> = [
      { write_outcome: 'unknown', cleanup_status: 'cleaned' },
      { write_outcome: 'written_verified', cleanup_status: 'unknown' },
      { write_outcome: 'written_verified', cleanup_status: 'failed' },
      { write_outcome: 'written_verified', cleanup_status: 'not_attempted' },
      { cleanup_status: 'cleaned' }, { write_outcome: 'written_verified' },
    ];
    for (const nested of unresolved) {
      const record: WriteGroupTestWriteOperation = { operation_id: 'TW1', action: 'test_write', status, ...nested,
        created_at: '2026-10-02T00:00:00Z', updated_at: '2026-10-02T00:00:00Z' };
      const view = available({ ...withLifecycle(), readiness: evidence(readiness()),
        operation: evidence({ kind: 'test_write', operation: record }) });
      expect(view.operation).toEqual({ state: 'available', value: { kind: 'test_write', operation: record } });
      expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'operation_unresolved' });
      expect(view.readiness).toMatchObject({ state: 'available', value: { ready: false } });
    }
  });
  it.each([
    ['written_verified', 'cleaned'], ['written_unverified', 'cleaned'], ['failed', 'cleaned'], ['failed', 'not_attempted'],
  ] as const)('only forwards lifecycle evidence after known test-write %s and cleanup %s', (write_outcome, cleanup_status) => {
    const view = available({ ...withLifecycle(), operation: evidence({ kind: 'test_write', operation: {
      operation_id: 'TW1', action: 'test_write', status: 'failed', write_outcome, cleanup_status,
      created_at: '2026-10-02T00:00:00Z', updated_at: '2026-10-02T00:00:00Z',
    } }) });
    expect(view.lifecycle).toMatchObject({ state: 'available', value: { remove: { state: 'blocked' } } });
    expect(view.readiness).toEqual({ state: 'unavailable', reason: 'missing' });
  });
  it('rejects an operation with no lookup identity', () => {
    const value = start();
    value.operation_id = '';
    expect(available({ ...input(), operation: evidence({ kind: 'start', operation: value }) }).operation)
      .toEqual({ state: 'unavailable', reason: 'missing' });
  });

  it('requires a lifecycle revision instead of accepting bare can-start/delete booleans', () => {
    const view = available({ ...withLifecycle(), lifecycle: {
      state: 'available', scope: { ...scope, groupRevision: undefined } as unknown as RecordingWorkflowScope,
      value: { edit: { state: 'available', reasonCodes: [] } },
    } });
    expect(view.lifecycle).toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
    expectTypeOf<boolean>().not.toMatchTypeOf<NonNullable<RecordingWorkflowInput['lifecycle']>>();
  });
  it.each(['workspace', 'group', 'effective revision', 'saved revision', 'requested revision'] as const)(
    'checks safe runtime %s independently of its envelope', (field) => {
      const proof = runtime();
      if (field === 'workspace') proof.groups[0].effective!.binding.workspace_id = 'foreign';
      if (field === 'group') proof.groups[0].group_id = 'foreign';
      if (field === 'effective revision') proof.groups[0].effective!.binding.entity_revision = 'foreign';
      if (field === 'saved revision') proof.groups[0].saved_revision = 'foreign';
      if (field === 'requested revision') proof.groups[0].applied_revision = 'foreign';
      expect(available({ ...input(), runtimeEvidence: evidence(proof) }).effective)
        .toEqual({ state: 'unavailable', reason: ['workspace', 'group'].includes(field) ? 'scope_mismatch' : 'revision_mismatch' });
    });
  it('does not promote a supplied snapshot when the canonical group has no applied revision', () => {
    expect(available({ ...input(), scope: { ...scope, appliedRevision: '' }, savedGroup: { ...draft(), applied_revision: '' },
      runtimeEvidence: { state: 'available', scope: { ...scope, appliedRevision: '' }, value: runtime() } }).effective)
      .toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
  });
  it('accepts legacy member ownership without inventing a managed role or circuit identity', () => {
    const group = draft();
    delete group.basic_managed_device_id;
    delete group.members[0].circuit_id;
    const view = available({ scope: { ...scope, circuitId: undefined }, savedGroup: group });
    expect(view.saved.group).not.toHaveProperty('basic_managed_device_id');
    expect(view.saved.group.members[0]).not.toHaveProperty('circuit_id');
  });
  it('rejects schema evidence without connector revision', () => {
    const preview = schema();
    delete preview.connector_revision;
    expect(available({ ...input(), schema: evidence(preview) }).schema)
      .toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
  });

  it('copies and deeply freezes accepted caller facts without freezing or mutating the caller', () => {
    const saved = draft();
    const snapshot = runtime();
    const selected = { ...scope };
    const action = { state: 'blocked' as const, reasonCodes: ['pending_delivery'] };
    const view = available({ ...withLifecycle(), scope: selected, savedGroup: saved,
      runtimeEvidence: evidence(snapshot), lifecycle: evidence({ remove: action }) });
    saved.recording_profile!.interval = 120;
    saved.row_policy.group_key_columns!.push('changed');
    snapshot.groups[0].effective!.destination.table = 'changed';
    action.reasonCodes.push('changed');
    selected.deviceId = 'B';
    expect(view.scope.deviceId).toBe('A');
    expect(view.saved.group.recording_profile?.interval).toBe(30);
    expect(view.saved.group.row_policy.group_key_columns).toEqual(['device']);
    expect(view.effective).toMatchObject({ value: { facts: { destination: { table: 'old_readings' } } } });
    expect(view.lifecycle).toMatchObject({ value: { remove: { reasonCodes: ['pending_delivery'] } } });
    expect(Object.isFrozen(view)).toBe(true);
    expect(Object.isFrozen(view.saved.group.members[0])).toBe(true);
    expect(Object.isFrozen(view.saved.group.row_policy.group_key_columns)).toBe(true);
    expect(Object.isFrozen(saved)).toBe(false);
  });

  it('has no I/O or mutable state when viewing A, unavailable B, then A again', () => {
    const fetchSpy = vi.spyOn(globalThis, 'fetch');
    const storageSpy = vi.spyOn(Storage.prototype, 'setItem');
    try {
      const before = input();
      const first = recordingWorkflowView(before);
      expect(recordingWorkflowView({ ...before, scope: { ...scope, deviceId: 'B' } })).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
      expect(recordingWorkflowView(before)).toEqual(first);
      expect(before.savedGroup.recording_profile?.interval).toBe(30);
      expect(fetchSpy).not.toHaveBeenCalled();
      expect(storageSpy).not.toHaveBeenCalled();
    } finally { fetchSpy.mockRestore(); storageSpy.mockRestore(); }
  });

  it('keeps saved/applied types distinct and exposes no second readiness field or mutation callbacks', () => {
    expectTypeOf<SavedRecordingFacts>().not.toMatchTypeOf<EffectiveRecordingFacts>();
    expectTypeOf<EffectiveRecordingFacts>().not.toMatchTypeOf<SavedRecordingFacts>();
    expectTypeOf<EffectiveRecordingFacts['facts']>().not.toMatchTypeOf<WriteGroup>();
    expectTypeOf<EffectiveRecordingFacts['facts']['binding']['entity_revision']>().toEqualTypeOf<RecordingGroupFacts['binding']['entity_revision']>();
    expectTypeOf<Extract<keyof RecordingWorkflowInput,
      'ready' | 'canStart' | 'canDelete' | 'onSave' | 'onApply' | 'onStart'>>().toEqualTypeOf<never>();
    expectTypeOf<SavedRecordingFacts['group']['members']>().not.toMatchTypeOf<WriteGroup['members']>();
    expectTypeOf<SavedRecordingFacts['group']['row_policy']['group_key_columns']>()
      .toEqualTypeOf<readonly string[] | undefined>();
  });
});
