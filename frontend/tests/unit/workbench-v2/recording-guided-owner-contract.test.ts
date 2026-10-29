import { describe, expect, expectTypeOf, it } from 'vitest';
import { inspectRecordingOwners, type RecordingOwnerObservation, type RecordingOwnerObservations, type RecordingOwnerScope, type RecordingSchemaReview } from '@/features/datalink/workbench-v2/state/recordingOwnerObservation';
import { recordingNextAction, recordingWorkflowView, type RecordingNextActionInput } from '@/features/datalink/workbench-v2/state/recordingWorkflowView';
import { savedGroup } from '../../fixtures/writeGroupState';

const scope: RecordingOwnerScope = { workspaceId: 'ws-1', workspaceRevision: 'W1', deviceId: 'dev-1', groupId: 'group-1', groupRevision: 'rev-1', appliedRevision: '' };
const deviceScope = { workspaceId: 'ws-1', workspaceRevision: 'W1', deviceId: 'dev-1' };
function observed(): RecordingOwnerObservations {
  return ['schema', 'start', 'test_write'].map(owner => ({ owner: owner as RecordingOwnerObservation['owner'],
    scope: owner === 'schema' ? scope : deviceScope, state: 'observed', result: 'none' }));
}
function input(observations = observed(), schemaReady = true): RecordingNextActionInput {
  const selected = { ...scope, groupId: 'group-1', groupRevision: 'rev-1', appliedRevision: '' };
  const workflow = recordingWorkflowView({ scope: selected, savedGroup: savedGroup(),
    readiness: { state: 'available', scope: selected, value: { workspace_id: 'ws-1', workspace_revision: 'W1', group_id: 'group-1', group_revision: 'rev-1', applied_revision: '', config_ready: true, schema_ready: schemaReady, ready: schemaReady, issues: [] } } });
  return { kind: 'saved', scope: selected, workflow, connection: 'online', nowMs: 0, ownerObservations: observations } as RecordingNextActionInput;
}
describe('token-free owner-local next-action contract', () => {
  it('accepts observed local none for every exact owner, without inventing global absence', () => {
    expect(inspectRecordingOwners(observed(), scope)).toEqual({ state: 'healthy' });
    expect(recordingNextAction(input())).toMatchObject({ intent: { kind: 'review_scoped_start' }, enabled: true, reasons: ['start_review'] });
  });
  it('keeps exact-scope start review navigation available offline with no lifecycle start permission', () => {
    expect(recordingNextAction({ ...input(), connection: 'offline' })).toMatchObject({ intent: { kind: 'review_scoped_start' }, enabled: true });
  });
  it('prioritizes known unresolved O1 over unreadable storage and other owner empty observations', () => {
    const observations: RecordingOwnerObservations = [observed()[0], { owner: 'start', scope: deviceScope, state: 'unresolved', reason: 'unknown', operationId: 'O1' },
      { owner: 'test_write', scope: deviceScope, state: 'unavailable', reason: 'unreadable' }];
    expect(recordingNextAction(input(observations))).toMatchObject({ intent: { kind: 'query_operation', operationId: 'O1', operationOwner: 'start' }, enabled: true });
    expect(recordingNextAction({ ...input(observations), connection: 'offline' })).toMatchObject({ intent: { kind: 'query_operation', operationId: 'O1' }, enabled: false });
  });
  it('routes unknown outcome without an ID to local recovery, never a fake query or automatic resend', () => {
    const observations: RecordingOwnerObservations = [observed()[0], { owner: 'start', scope: deviceScope, state: 'unresolved', reason: 'lost_ack' }, observed()[2]];
    expect(recordingNextAction({ ...input(observations), connection: 'offline' })).toMatchObject({ intent: { kind: 'recover_operation', operationOwner: 'start' }, enabled: true });
    expect(recordingNextAction(input(observations)).intent).not.toHaveProperty('operationId');
  });
  it.each(['loading', 'unreadable', 'invalid_record', 'scope_mismatch', 'revision_mismatch'] as const)('never treats %s as local none or enables a preview command', reason => {
    const observations: RecordingOwnerObservations = [observed()[0], { owner: 'start', scope: deviceScope, state: 'unavailable', reason }, observed()[2]];
    expect(recordingNextAction(input(observations, false))).toMatchObject({ intent: { kind: 'recover_operation' }, enabled: true });
  });
  it.each(['workspaceId', 'deviceId', 'circuitId', 'groupId', 'workspaceRevision', 'groupRevision', 'appliedRevision'] as const)('does not borrow schema owner evidence with a different %s', key => {
    const observations: RecordingOwnerObservations = [{ owner: 'schema', scope: { ...scope, [key]: 'foreign' }, state: 'observed', result: 'none' }, ...observed().slice(1)];
    expect(inspectRecordingOwners(observations, scope).state).toBe('unavailable');
    expect(recordingNextAction(input(observations, false)).intent.kind).toBe('recover_operation');
  });
  it('requires all owners and rejects duplicate current owner evidence', () => {
    expect(inspectRecordingOwners(observed().slice(0, 2), scope).state).toBe('unavailable');
    expect(inspectRecordingOwners([...observed(), observed()[1]], scope).state).toBe('unavailable');
  });
  it('makes missing schema an explicit preview command and reviewed schema a guarded confirm intent with no token', () => {
    const source = input(observed(), false);
    expect(recordingNextAction(source)).toMatchObject({ intent: { kind: 'preview_schema' }, enabled: true });
    expect(recordingNextAction({ ...source, connection: 'offline' })).toMatchObject({ intent: { kind: 'preview_schema' }, enabled: false });
    const schemaReview: RecordingSchemaReview = { scope, stage: 'reviewed', operationId: 'DDL1' };
    const next = recordingNextAction({ ...source, schemaReview } as RecordingNextActionInput);
    expect(next.intent).toMatchObject({ kind: 'confirm_schema', operationId: 'DDL1' });
    expect(next.intent).not.toHaveProperty('token'); expect(next.intent).not.toHaveProperty('digest');
  });
  it('keeps preview navigation local and rejects stale review references', () => {
    const source = input(observed(), false);
    expect(recordingNextAction({ ...source, connection: 'offline', schemaReview: { scope, stage: 'previewed', operationId: 'DDL1' } } as RecordingNextActionInput))
      .toMatchObject({ intent: { kind: 'review_schema' }, enabled: true });
    expect(recordingNextAction({ ...source, schemaReview: { scope: { ...scope, groupRevision: 'old' }, stage: 'reviewed', operationId: 'DDL1' } } as RecordingNextActionInput).intent.kind)
      .toBe('recover_operation');
  });
  it('owner inspection does not let legacy none bypass a known unresolved owner or invent server edit permission', () => {
    const selected = { ...scope, groupId: 'group-1', groupRevision: 'rev-1', appliedRevision: '' };
    const observation: RecordingOwnerObservation = { owner: 'start', scope: deviceScope, state: 'unresolved', reason: 'unknown', operationId: 'O1' };
    const value = recordingWorkflowView({ scope: selected, savedGroup: savedGroup(),
      operation: { state: 'available', scope: selected, value: { kind: 'none' } },
      lifecycle: { state: 'available', scope: selected, value: { edit: { state: 'available', reasonCodes: [] } } },
      ownerObservations: [observed()[0], observation, observed()[2]] });
    expect(value).toMatchObject({ lifecycle: { state: 'unavailable', reason: 'operation_unresolved' } });
    const healthy = recordingWorkflowView({ scope: selected, savedGroup: savedGroup(), ownerObservations: observed(),
      lifecycle: { state: 'available', scope: selected, value: { edit: { state: 'blocked', reasonCodes: ['pending_delivery'] } } } });
    expect(healthy).toMatchObject({ lifecycle: { state: 'available', value: { edit: { state: 'blocked', reasonCodes: ['pending_delivery'] } } } });
    const noServerAuthority = recordingWorkflowView({ scope: selected, savedGroup: savedGroup(), ownerObservations: observed() });
    expect(noServerAuthority).toMatchObject({ lifecycle: { state: 'unavailable', reason: 'missing' } });
  });
  it('does not expose token, raw operation, request, draft, or error fields in observation contracts', () => {
    expectTypeOf<RecordingOwnerObservation>().not.toHaveProperty('token'); expectTypeOf<RecordingOwnerObservation>().not.toHaveProperty('operation');
    expectTypeOf<RecordingOwnerObservation>().not.toHaveProperty('request'); expectTypeOf<RecordingOwnerObservation>().not.toHaveProperty('draft');
    expectTypeOf<RecordingSchemaReview>().not.toHaveProperty('token'); expectTypeOf<RecordingSchemaReview>().not.toHaveProperty('digest');
  });
});
