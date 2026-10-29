import { expect, it } from 'vitest';
import { createGuidedRecordingMock } from '../../e2e/fixtures/recording-guided-api';
import { guidedScenarioFrom, recordingGuidedHostState, initialGuidedGroups, recordingGuidedScenarios, GUIDED_PATH as P } from '../../e2e/fixtures/recording-guided-state';
import { parseRecordingStartOperation } from '@/utils/studioV2RecordingStartJson';
import { parseRecordingSchemaOperation, parseRecordingSchemaPreviewToken } from '@/utils/recordingPlanJson';
import type { WriteGroup, WriteGroupDraft } from '@/types/studioV2WriteGroup';
import type { SchemaPreviewToken } from '@/types/recordingPlan';
import type { RecordingStartRequest } from '@/types/studioV2RecordingStart';
import type { LifecycleMockReply } from '../../e2e/fixtures/recording-lifecycle-api';
import { parseRuntimeRuleEvidence } from '@/services/runtimeRuleEvidenceParser';
import type { RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import { recordingWorkflowView } from '@/features/datalink/workbench-v2/state/recordingWorkflowView';

type Mock = ReturnType<typeof createGuidedRecordingMock>;
function data<T>(reply: LifecycleMockReply): T { expect(reply.status).toBe(200); return (reply.body as { data: T }).data; }
function call(mock: Mock, method: string, path: string, body?: unknown) { return mock.handle({ method, url: `${P}${path}`, ...(body ? { body } : {}) }); }
function cas(mock: Mock, group = mock.snapshot().groups[0]) {
  return { workspace_id: 'ws-1', expected_workspace_revision: mock.snapshot().workspaceRevision,
    expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision };
}
function draft(group = initialGuidedGroups('saved-long')[0]): WriteGroupDraft {
  return { workspace_id: group.workspace_id, name: group.name, members: group.members, recording_profile: group.recording_profile,
    layout_version: group.layout_version, storage_mode: group.storage_mode, row_policy: group.row_policy, write_policy: group.write_policy,
    destination: { connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: 'main', table_name: '', storage_strategy: 'managed' } };
}
function ensure(mock: Mock) {
  return data<{ group: WriteGroup }>(call(mock, 'POST', '/write-groups/circuit-profiles/ensure', {
    workspace_id: 'ws-1', expected_workspace_revision: mock.snapshot().workspaceRevision, expected_connector_revision: 'crev-1', group: draft(),
  })).group;
}
function start(mock: Mock): RecordingStartRequest {
  return { request_id: 'REQUEST1', workspace_id: 'ws-1', expected_workspace_revision: mock.snapshot().workspaceRevision, device_ids: ['dev-1'],
    groups: mock.snapshot().groups.map(group => ({ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: 'crev-1' })) };
}

it('offers deterministic A/B/empty states and forty independently selectable sources, without claiming effective runtime', () => {
  expect(guidedScenarioFrom('invalid')).toBe('new'); expect(guidedScenarioFrom('empty-device')).toBe('empty-device');
  const state = recordingGuidedHostState('long-sources');
  expect(state.devices.map(device => device.id)).toEqual(['dev-1', 'dev-2', 'dev-empty']);
  expect(state.points).toHaveLength(40); expect(new Set(state.points.map(point => point.id)).size).toBe(40);
  expect(Object.values(state.mappings).map(mapping => mapping.persisted_point_id)).toContain('pt-39');
});
it('serves one canonical list/head to scoped readiness/lifecycle and explicit unavailable runtime evidence', () => {
  const mock = createGuidedRecordingMock('saved-long');
  expect(data(call(mock, 'GET', '/write-groups'))).toMatchObject({ workspace_id: 'ws-1', workspace_revision: 'W1', groups: [{ id: 'G1' }] });
  expect(data(call(mock, 'GET', '/write-groups/G1/readiness'))).toMatchObject({ workspace_revision: 'W1', group_revision: 'R1', ready: true });
  expect(data(call(mock, 'GET', '/write-groups/G1/lifecycle'))).toMatchObject({ workspace_revision: 'W1', group_revision: 'R1' });
  expect(data(call(mock, 'GET', '/runtime-context?workspace_id=ws-1&workspace_revision=W1&device_id=dev-1'))).toMatchObject({
    rule_evidence: { freshness: 'unavailable', observed_at: null, source: null, groups: [] },
  });
  expect(data(call(mock, 'GET', '/database-operations/test-writes?workspace_id=ws-1'))).toEqual({ workspace_id: 'ws-1', complete: true, operations: [] });
  expect(data(call(mock, 'GET', '/write-groups/G1/cutovers'))).toEqual([]);
  expect(mock.ledger.every(entry => !entry.mutation)).toBe(true);
});
it('requires explicit schema preview then exact token confirmation before the original singleton start', () => {
  const mock = createGuidedRecordingMock(); ensure(mock);
  expect(call(mock, 'POST', '/recording-start', start(mock)).status).toBe(409);
  const before = mock.snapshot();
  const preview = data<SchemaPreviewToken>(call(mock, 'POST', '/write-groups/G1/schema-preview', cas(mock)));
  expect(parseRecordingSchemaPreviewToken(preview, { requireGroupLayout: true })).not.toBeNull();
  expect(mock.snapshot().workspaceRevision).toBe(before.workspaceRevision);
  expect(mock.snapshot().groups).toEqual(before.groups);
  expect(parseRecordingSchemaOperation(data(call(mock, 'POST', '/write-groups/G1/schema-apply', { ...cas(mock), token: preview.token, operation_id: preview.operation_id })))).not.toBeNull();
  expect(mock.snapshot().workspaceRevision).not.toBe(before.workspaceRevision);
  const request = start(mock), operation = data(call(mock, 'POST', '/recording-start', request));
  expect(parseRecordingStartOperation(operation)).toMatchObject({ operation_id: 'O1', status: 'succeeded', device_ids: ['dev-1'], devices: [{ device_id: 'dev-1', activated: true }] });
  expect(data(call(mock, 'GET', '/recording-start/operations/O1'))).toEqual(operation);
  expect(mock.snapshot()).toMatchObject({ originalStart: request, operation });
  expect(mock.ledger.filter(entry => entry.action === 'schema-preview')).toHaveLength(1);
});
it.each(['workspace', 'head', 'group', 'connector', 'extra'] as const)('rejects %s preview scope without creating a token or advancing state', invalid => {
  const mock = createGuidedRecordingMock('saved-long'), before = mock.snapshot();
  const body: Record<string, unknown> = cas(mock);
  if (invalid === 'workspace') body.workspace_id = 'foreign';
  if (invalid === 'head') body.expected_workspace_revision = 'W0';
  if (invalid === 'group') body.expected_group_revision = 'old';
  if (invalid === 'connector') body.expected_connector_revision = 'old';
  if (invalid === 'extra') body.sql = 'not accepted';
  expect(call(mock, 'POST', '/write-groups/G1/schema-preview', body).status).not.toBe(200);
  expect(mock.snapshot()).toEqual(before); expect(mock.ledger.at(-1)?.mutation).toBe(true);
});
it.each(['missing', 'token', 'operation', 'stale'] as const)('rejects %s schema confirmation with no schema/head change', invalid => {
  const mock = createGuidedRecordingMock(), group = ensure(mock);
  const preview = data<SchemaPreviewToken>(call(mock, 'POST', `/write-groups/${group.id}/schema-preview`, cas(mock)));
  const before = mock.snapshot();
  const body = { ...cas(mock), ...(invalid === 'missing' ? {} : { token: invalid === 'token' ? 'foreign' : preview.token,
    operation_id: invalid === 'operation' ? 'foreign' : preview.operation_id }), ...(invalid === 'stale' ? { expected_workspace_revision: 'W0' } : {}) };
  expect(call(mock, 'POST', '/write-groups/G1/schema-apply', body).status).not.toBe(200); expect(mock.snapshot()).toEqual(before);
});
it.each(['device', 'duplicate-device', 'group', 'duplicate-group', 'head', 'connector', 'extra'] as const)('rejects %s start scope without accepting an original request', invalid => {
  const mock = createGuidedRecordingMock('saved-long'), body: Record<string, unknown> = { ...start(mock) };
  if (invalid === 'device') body.device_ids = ['dev-2'];
  if (invalid === 'duplicate-device') body.device_ids = ['dev-1', 'dev-1'];
  if (invalid === 'group') body.groups = [];
  if (invalid === 'duplicate-group') body.groups = [...start(mock).groups, ...start(mock).groups];
  if (invalid === 'head') body.expected_workspace_revision = 'W0';
  if (invalid === 'connector') body.groups = [{ ...start(mock).groups[0], expected_connector_revision: 'old' }];
  if (invalid === 'extra') body.apply_all = true;
  expect(call(mock, 'POST', '/recording-start', body).status).not.toBe(200);
  expect(mock.snapshot()).not.toHaveProperty('originalStart'); expect(mock.ledger.at(-1)?.mutation).toBe(true);
});
it.each(['partial', 'unknown'] as const)('preserves the original %s receipt after lost ACK and rejects a replacement request', outcome => {
  const mock = createGuidedRecordingMock(outcome), request = start(mock); mock.loseNext('start');
  expect(call(mock, 'POST', '/recording-start', request).transport).toBe('abort');
  const result = data(call(mock, 'GET', '/recording-start/operations/O1'));
  expect(parseRecordingStartOperation(result)).toMatchObject({ operation_id: 'O1', status: outcome, device_ids: ['dev-1'] });
  expect(call(mock, 'POST', '/recording-start', { ...request, request_id: 'replacement' }).status).toBe(409);
  expect(mock.snapshot()).toMatchObject({ originalStart: request, operation: result });
});
it.each(['saved-long', 'saved-wide', 'saved-custom'] as const)('stops/removes %s using current CAS without changing original history/destination or converting layout', scenario => {
  const mock = createGuidedRecordingMock(scenario), before = mock.snapshot();
  const first = data<{ history: unknown }>(call(mock, 'GET', '/write-groups/G1/lifecycle')).history;
  const oldCAS = cas(mock); data(call(mock, 'POST', '/write-groups/G1/disable', oldCAS));
  expect(call(mock, 'DELETE', '/write-groups/G1', oldCAS).status).toBe(409);
  data(call(mock, 'DELETE', '/write-groups/G1', cas(mock)));
  const current = mock.snapshot().groups[0]; expect(current.status).toBe('deleted');
  expect(current.destination).toEqual(before.groups[0].destination); expect(current.layout_version).toBe(before.groups[0].layout_version);
  const after = data<{ history: { generations: { revisions: unknown }[] } }>(call(mock, 'GET', '/write-groups/G1/lifecycle')).history;
  expect(after.generations[0].revisions).toEqual((first as typeof after).generations[0].revisions);
  expect(mock.ledger.filter(entry => ['ensure', 'start', 'schema-apply'].includes(entry.action))).toEqual([]);
});
it('rejects blocked removal, saved-custom conversion, unknown routes and foreign reads, counting forbidden writes', () => {
  const blocked = createGuidedRecordingMock('blocked-wide'); expect(call(blocked, 'DELETE', '/write-groups/G1', cas(blocked)).status).toBe(409);
  const mock = createGuidedRecordingMock('saved-custom');
  expect(call(mock, 'POST', '/write-groups/circuit-profiles/ensure', { workspace_id: 'ws-1', expected_workspace_revision: 'W1', expected_connector_revision: 'crev-1', group: draft() }).status).toBe(409);
  expect(call(mock, 'POST', '/activate', {}).status).toBe(404);
  expect(call(mock, 'GET', '/write-groups?workspace_id=foreign').status).not.toBe(200);
  expect(call(mock, 'GET', '/runtime-context?workspace_id=ws-1&workspace_revision=W0&device_id=dev-1').status).toBe(409);
  expect(mock.ledger.filter(entry => entry.mutation)).toHaveLength(2);
});
it('controlled rejection is distinct from lost ACK and snapshots cannot alter canonical state', () => {
  const mock = createGuidedRecordingMock('saved-long'); mock.rejectNext('start');
  expect(call(mock, 'POST', '/recording-start', start(mock)).status).toBe(503);
  expect(mock.snapshot()).not.toHaveProperty('originalStart');
  const copy = mock.snapshot(); copy.groups[0].destination.table_name = 'foreign';
  expect(mock.snapshot().groups[0].destination.table_name).not.toBe('foreign');
});
it('keeps a custom edit custom and preserves original applied history while making changed schema unready', () => {
  const mock = createGuidedRecordingMock('saved-custom'), saved = mock.snapshot().groups[0];
  const edited = { ...draft(saved), name: 'Renamed custom recording', destination: { connector_id: 'conn-1', connector_revision: 'crev-1',
    table_schema: 'main', table_name: saved.destination.table_name, storage_strategy: 'custom' } };
  const result = data<{ group: WriteGroup }>(call(mock, 'PUT', '/write-groups/G1', { ...cas(mock), group: edited })).group;
  expect(result).not.toHaveProperty('basic_managed_device_id'); expect(result).not.toHaveProperty('recording_profile');
  expect(result.applied_revision).toBe(saved.applied_revision); expect(result.destination.storage_strategy).toBe('custom');
  expect(data(call(mock, 'GET', '/write-groups/G1/readiness'))).toMatchObject({ ready: false, schema_ready: false });
  expect(mock.snapshot().history[0].generations[0].revisions[0].group_revision).toBe('APPLIED1');
});
it('rejects an old preview after a canonical edit even when confirmation carries the new head', () => {
  const mock = createGuidedRecordingMock('saved-long'), saved = mock.snapshot().groups[0];
  const preview = data<SchemaPreviewToken>(call(mock, 'POST', '/write-groups/G1/schema-preview', cas(mock)));
  const edited = draft(saved); edited.destination.table_name = saved.destination.table_name;
  data(call(mock, 'PUT', '/write-groups/G1', { ...cas(mock), group: edited }));
  expect(call(mock, 'POST', '/write-groups/G1/schema-apply', { ...cas(mock), token: preview.token, operation_id: preview.operation_id }).status).toBe(409);
  expect(call(mock, 'POST', '/recording-start', start(mock)).status).toBe(409);
});
it('records an accepted schema confirmation before lost ACK, permits exact lookup but rejects a second apply', () => {
  const mock = createGuidedRecordingMock(); ensure(mock);
  const preview = data<SchemaPreviewToken>(call(mock, 'POST', '/write-groups/G1/schema-preview', cas(mock)));
  const request = { ...cas(mock), token: preview.token, operation_id: preview.operation_id }; mock.loseNext('schema-apply');
  expect(call(mock, 'POST', '/write-groups/G1/schema-apply', request).transport).toBe('abort');
  expect(data(call(mock, 'GET', `/database-operations/${preview.operation_id}`))).toMatchObject({ status: 'succeeded', operation_id: preview.operation_id });
  expect(call(mock, 'POST', '/write-groups/G1/schema-apply', request).status).toBe(409);
  expect(call(mock, 'GET', '/database-operations/foreign').status).toBe(404);
});
it('does not admit a start through blocked connector lifecycle evidence', () => {
  const mock = createGuidedRecordingMock('blocked-wide');
  expect(data(call(mock, 'GET', '/write-groups/G1/readiness'))).toMatchObject({ ready: false, config_ready: false });
  expect(call(mock, 'POST', '/recording-start', start(mock)).status).toBe(409);
  expect(mock.snapshot()).not.toHaveProperty('originalStart');
});
it.each([
  '/write-groups/G-other/lifecycle', '/write-groups/G1?workspace_id=foreign', '/write-groups?unexpected=1',
  '/write-groups?workspace_id=ws-1&workspace_id=ws-1', '/measurements?device_id=foreign',
  '/database-metadata?expected_connector_revision=old', '/database-metadata?expected_connector_revision=crev-1&group_id=G1&expected_group_revision=old',
  '/runtime-context?workspace_id=ws-1&workspace_revision=W1&device_id=foreign',
])('rejects unexpected read scope %s without state changes', path => {
  const mock = createGuidedRecordingMock('saved-long'), before = mock.snapshot();
  expect(call(mock, 'GET', path).status).not.toBe(200); expect(mock.snapshot()).toEqual(before);
  expect(mock.ledger).toHaveLength(1); expect(mock.ledger[0].mutation).toBe(false);
});

function runtime(mock: Mock, deviceId = 'dev-1'): RuleEvidenceView {
  return data<{ rule_evidence: RuleEvidenceView }>(call(mock, 'GET', `/runtime-context?workspace_id=ws-1&workspace_revision=${mock.snapshot().workspaceRevision}&device_id=${deviceId}`)).rule_evidence;
}
it('provides explicit simulated saved R2/30 versus immutable effective R1/15 through the real parser and projection', () => {
  expect(guidedScenarioFrom('saved-different')).toBe('saved-different');
  const mock = createGuidedRecordingMock('saved-different'), saved = mock.snapshot().groups[0];
  expect(mock.snapshot().workspaceRevision).toBe('W2');
  expect(saved).toMatchObject({ revision: 'R2', applied_revision: 'R1', recording_profile: { interval: 30 }, row_policy: { interval_seconds: 30 } });
  const parsed = parseRuntimeRuleEvidence(runtime(mock), { workspace_id: 'ws-1', workspace_revision: 'W2', device_id: 'dev-1' });
  expect(parsed.ok).toBe(true); if (!parsed.ok) throw new Error('Expected valid simulated evidence');
  expect(parsed.value).toMatchObject({ observed_at: '2030-01-01T00:00:00.000Z', freshness: 'fresh', groups: [{
    group_id: 'G1', saved_revision: 'R2', applied_revision: 'R1', effective_revision: 'R1', as_of: '2030-01-01T00:00:00.000Z',
    saved: { binding: { entity_revision: 'R2' }, recording_profile: { interval: 30 }, interval_seconds: 30 },
    effective: { binding: { entity_revision: 'R1' }, recording_profile: { interval: 15 }, interval_seconds: 15 }, delivery: null,
  }] });
  const scope = { workspaceId: 'ws-1', workspaceRevision: 'W2', deviceId: 'dev-1', circuitId: 'cir-saved', groupId: 'G1', groupRevision: 'R2', appliedRevision: 'R1' };
  expect(recordingWorkflowView({ scope, savedGroup: saved, runtimeEvidence: { state: 'available', scope, value: parsed.value } })).toMatchObject({
    state: 'available', runtime: { state: 'available', value: { saved: { interval_seconds: 30 } } },
    effective: { state: 'available', value: { revision: 'R1', facts: { interval_seconds: 15 } } },
  });
  expect(JSON.stringify(parsed.value)).not.toMatch(/mock-private|recordings\.sqlite|password|token|CREATE TABLE/);
  expect(mock.ledger.every(entry => !entry.mutation)).toBe(true);
});
it('keeps every ordinary scenario and other selected devices explicitly unavailable', () => {
  for (const scenario of recordingGuidedScenarios.filter(value => value !== 'saved-different')) {
    expect(runtime(createGuidedRecordingMock(scenario))).toMatchObject({ freshness: 'unavailable', observed_at: null, groups: [] });
  }
  const mock = createGuidedRecordingMock('saved-different');
  for (const device of ['dev-2', 'dev-empty']) expect(runtime(mock, device)).toMatchObject({ device_id: device, freshness: 'unavailable', groups: [] });
});
it.each(['workspace', 'revision', 'device'] as const)('rejects a simulated evidence read or parse with foreign/stale %s scope', invalid => {
  const mock = createGuidedRecordingMock('saved-different'), proof = runtime(mock);
  const scope = { workspace_id: invalid === 'workspace' ? 'foreign' : 'ws-1', workspace_revision: invalid === 'revision' ? 'W1' : 'W2', device_id: invalid === 'device' ? 'foreign' : 'dev-1' };
  expect(call(mock, 'GET', `/runtime-context?${new URLSearchParams(scope)}`).status).not.toBe(200);
  expect(parseRuntimeRuleEvidence(proof, scope)).toMatchObject({ ok: false, reason: invalid === 'revision' ? 'stale-reference' : 'foreign-scope' });
  expect(mock.ledger.every(entry => !entry.mutation)).toBe(true);
});
it.each(['binding', 'as-of', 'saved-binding'] as const)('does not accept mixed simulated %s evidence', invalid => {
  const proof = runtime(createGuidedRecordingMock('saved-different'));
  expect(proof.groups).toHaveLength(1); const group = proof.groups[0];
  if (!group.saved || !group.effective) throw new Error('Expected separate simulated snapshots');
  if (invalid === 'binding') group.effective.binding.entity_revision = 'R2';
  if (invalid === 'saved-binding') group.saved.binding.entity_revision = 'R1';
  if (invalid === 'as-of') group.as_of = '2030-01-01T00:00:01.000Z';
  expect(parseRuntimeRuleEvidence(proof, { workspace_id: 'ws-1', workspace_revision: 'W2', device_id: 'dev-1' })).toMatchObject({ ok: false, reason: 'stale-reference' });
});
it('does not reuse the fixed effective observation after a canonical head changes', () => {
  const mock = createGuidedRecordingMock('saved-different');
  expect(runtime(mock).freshness).toBe('fresh');
  expect(call(mock, 'POST', '/write-groups/G1/disable', cas(mock)).status).toBe(200);
  expect(call(mock, 'GET', '/runtime-context?workspace_id=ws-1&workspace_revision=W2&device_id=dev-1').status).toBe(409);
  expect(runtime(mock)).toMatchObject({ workspace_revision: 'W3', freshness: 'unavailable', observed_at: null, groups: [] });
});
