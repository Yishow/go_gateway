import { afterEach, describe, expect, it } from 'vitest';
import { AxiosError, type AxiosResponse } from 'axios';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WriteGroupLifecycleAPI as lifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2WriteGroupSuccessorAPI as api } from '@/services/studioV2WriteGroupSuccessor';
import { studioV2WorkspaceWriteGroupsAPI as groupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { changeTargetLayout, reviewRequest } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeHelpers';
import { parseWriteGroupLifecycle } from '@/utils/studioV2WriteGroupLifecycleJson';
import { parseWriteGroupListData } from '@/utils/studioV2WriteGroupJson';
import { createRecordingLifecycleMock, type RecordingLifecycleMock } from '../../e2e/fixtures/recording-lifecycle-api';
import { lifecycleScenarios, fixtureCurrent, fixtureLifecycle, FIXTURE_PATH, type LifecycleScenario } from '../../e2e/fixtures/recording-lifecycle-state';

// Contract and controller-input checks only: not browser execution, screenshots, backend wiring, or SQL evidence.
const originalAdapter = studioV2DatalinkApi.defaults.adapter;
afterEach(() => { studioV2DatalinkApi.defaults.adapter = originalAdapter; });
const path = `${FIXTURE_PATH}/write-groups/G1/successor`;
function connect(mock: RecordingLifecycleMock) {
  studioV2DatalinkApi.defaults.adapter = async config => {
    const url = new URL(`${config.baseURL}${config.url}`, 'http://fixture.invalid');
    for (const [key, value] of Object.entries(config.params ?? {})) url.searchParams.set(key, String(value));
    const reply = mock.handle({ method: (config.method ?? 'GET').toUpperCase(), url: url.href,
      body: typeof config.data === 'string' ? JSON.parse(config.data) : config.data });
    const response: AxiosResponse = { data: reply.body, status: reply.status, statusText: '', headers: {}, config };
    if (reply.transport === 'abort') throw new AxiosError('MOCK lost response', 'ERR_NETWORK', config);
    if (reply.status >= 400) throw new AxiosError('MOCK rejected', undefined, config, undefined, response);
    return response;
  };
}
async function reviewed(scenario: LifecycleScenario = 'applied-wide', layout?: 'managed_long_v1' | 'managed_wide_profile_v1', destination = 'conn-1') {
  const mock = createRecordingLifecycleMock(scenario); connect(mock);
  const lifecycle = await lifecycleAPI.lifecycle('G1');
  let target = lifecycle.successor_seed!.group;
  if (layout) target = changeTargetLayout(target, layout);
  target.destination.connector_id = destination;
  target.destination.connector_revision = destination === 'conn-1' ? 'crev-1' : 'crev-2';
  const review = await api.review('G1', reviewRequest(lifecycle, target, 'ROOT'));
  return { mock, review };
}
async function reserved(scenario: LifecycleScenario = 'applied-wide', layout?: 'managed_long_v1' | 'managed_wide_profile_v1', destination?: string) {
  const { mock, review } = await reviewed(scenario, layout, destination);
  const request = { envelope: review.envelope, review_digest: review.review_digest };
  const preparation = await api.reserve('G1', request);
  const stage = { ...request, current: fixtureCurrent(mock.snapshot()) };
  return { mock, review, preparation, stage };
}
async function prepared(scenario: LifecycleScenario) {
  const result = await reserved(scenario);
  const preview = await api.schemaPreview('G1', result.stage);
  await api.schemaConfirm('G1', { ...result.stage, token: preview.token, digest: preview.digest, schema_operation_id: preview.operation_id });
  return { ...result, preview };
}
async function ready(scenario: LifecycleScenario) {
  const result = await prepared(scenario);
  const review = await api.activationReview('G1', { stage: result.stage, boundary: '2099-01-01T01:00:00.000Z' });
  const readiness = await api.activationReadiness('G1', { envelope: review.envelope, review_digest: review.review_digest });
  return { ...result, request: { envelope: readiness.review.envelope, review_digest: readiness.review.review_digest, readiness_digest: readiness.readiness_digest } };
}
function direct(mock: RecordingLifecycleMock, action: string, body: object) { return mock.handle({ method: 'POST', url: path, body: { ...body, action } }); }

describe('test-only real-host browser fixture: strict API contracts, NOT browser acceptance', () => {
  it.each(lifecycleScenarios)('%s emits list/lifecycle/readiness accepted by real production parsers', async scenario => {
    const mock = createRecordingLifecycleMock(scenario); connect(mock);
    expect(await groupsAPI.list()).toMatchObject({ workspace_id: 'ws-1', groups: [{ id: 'G1' }] });
    expect(await lifecycleAPI.lifecycle('G1')).toMatchObject({ group_id: 'G1', history: { lineage_id: 'L' } });
    expect(await groupsAPI.readiness('G1')).toMatchObject({ group_revision: 'R1' });
    expect(await groupsAPI.delivery('G1')).toMatchObject({ group_id: 'G1' });
    expect(mock.ledger.every(entry => !entry.mutation)).toBe(true);
  });
  it.each([['applied-wide', 'managed_long_v1'], ['applied-long', 'managed_wide_profile_v1']] as const)('reviews %s → %s without mutating or preparing schema', async (scenario, layout) => {
    const { mock, review } = await reviewed(scenario, layout);
    expect(review.envelope.request.group.layout_version).toBe(layout);
    expect(mock.snapshot().target).toBeUndefined();
    expect(mock.ledger.filter(entry => entry.mutation)).toHaveLength(0);
  });
  it('same-mode destination retains source semantics and uses exact reviewed connector', async () => {
    const { mock, stage } = await reserved('applied-long', undefined, 'conn-2');
    const preview = await api.schemaPreview('G1', stage);
    expect(preview.connector_id).toBe('conn-2');
    expect(stage.current.target_connector_revision).toBe('crev-2');
    expect(mock.snapshot().target?.members).toEqual(mock.snapshot().source.members);
  });
  it('reservation is immutable, repeated exact reserve is idempotent, and current lifecycle includes the reserved generation', async () => {
    const { mock, preparation, stage } = await reserved();
    expect(await api.reserve('G1', { envelope: stage.envelope, review_digest: stage.review_digest })).toEqual(preparation);
    for (const group of [mock.snapshot().source, mock.snapshot().target!]) expect(parseWriteGroupLifecycle(fixtureLifecycle(mock.snapshot(), group))).not.toBeNull();
    const preview = await api.schemaPreview('G1', stage);
    await api.schemaConfirm('G1', { ...stage, token: preview.token, digest: preview.digest, schema_operation_id: preview.operation_id });
    expect((await api.operationStatus('G1', { workspace_id: 'ws-1', lineage_id: 'L', operation_id: 'ROOT' })).reserve.result).toEqual(preparation);
    expect(preparation.group.destination.schema_revision).toBeUndefined();
  });
  it.each(['applied-wide', 'deleted-unapplied'] as const)('%s validates separate schema and final activation contracts through real services', async scenario => {
    const { mock, request } = await ready(scenario);
    expect(mock.ledger.filter(entry => entry.action.endsWith('-confirm'))).toHaveLength(1);
    const result = scenario === 'deleted-unapplied' ? await api.confirmFirstStart('G1', request) : await api.confirmCutover('G1', request);
    expect(result.boundary).toBe(request.envelope.request.boundary);
    const status = await api.operationStatus('G1', { workspace_id: 'ws-1', lineage_id: 'L', operation_id: 'ROOT' });
    expect(status[scenario === 'deleted-unapplied' ? 'first_start' : 'cutover'].status).toBe('succeeded');
    expect(status[scenario === 'deleted-unapplied' ? 'cutover' : 'first_start'].status).toBe('not_submitted');
    expect(await groupsAPI.get('G2')).toMatchObject({ group: { applied_revision: 'R2' } });
    expect(parseWriteGroupLifecycle(fixtureLifecycle(mock.snapshot(), mock.snapshot().target!))).not.toBeNull();
  });
  it('unknown/lost preview commits no DDL and retains original operation for explicit same-root repreview', async () => {
    const { mock, stage } = await reserved(); mock.loseNext('schema-preview');
    await expect(api.schemaPreview('G1', stage)).rejects.toMatchObject({ outcome: 'unconfirmed' });
    const status = await api.operationStatus('G1', { workspace_id: 'ws-1', lineage_id: 'L', operation_id: 'ROOT' });
    expect(status.schema_operations).toEqual([{ operation_id: 'SCHEMA-1', status: 'preview_not_submitted' }]);
    expect(mock.ledger.filter(entry => entry.action === 'schema-confirm')).toHaveLength(0);
    expect((await api.schemaPreview('G1', stage)).operation_id).toBe('SCHEMA-2');
    expect(mock.snapshot().envelope?.request.operation_id).toBe('ROOT');
  });
  it('controlled unknown schema never implies preparation, then readback observes the same schema result', async () => {
    const { mock, stage } = await reserved(); const preview = await api.schemaPreview('G1', stage);
    mock.setSchemaOperationStatus('unknown'); mock.loseNext('schema-confirm');
    await expect(api.schemaConfirm('G1', { ...stage, token: preview.token, digest: preview.digest, schema_operation_id: preview.operation_id })).rejects.toMatchObject({ outcome: 'unconfirmed' });
    expect((await api.operationStatus('G1', { workspace_id: 'ws-1', lineage_id: 'L', operation_id: 'ROOT' })).schema_operations[0].status).toBe('unknown');
    expect(direct(mock, 'schema-preview', stage).status).toBe(409);
    mock.setSchemaOperationStatus('succeeded');
    expect((await groupsAPI.get('G2')).group.destination.schema_revision).toBe('SCHEMA-1');
  });
  it('late failed status read does not erase mock terminal result; transport delay/loss are explicit', async () => {
    const { mock, request } = await ready('applied-wide'); await api.confirmCutover('G1', request);
    mock.rejectNext('operation-status');
    await expect(api.operationStatus('G1', { workspace_id: 'ws-1', lineage_id: 'L', operation_id: 'ROOT' })).rejects.toMatchObject({ outcome: 'unconfirmed' });
    expect(mock.snapshot().status?.cutover.status).toBe('succeeded');
    mock.delayNext('operation-status', 75);
    expect(direct(mock, 'operation-status', { workspace_id: 'ws-1', lineage_id: 'L', operation_id: 'ROOT' }).delayMs).toBe(75);
  });
  it('blocked activation readiness passes strict parser but cannot authorize final confirmation', async () => {
    const { mock, stage } = await prepared('applied-wide'); mock.setActivationReady(false);
    const review = await api.activationReview('G1', { stage, boundary: '2099-01-01T00:00:00Z' });
    const readiness = await api.activationReadiness('G1', { envelope: review.envelope, review_digest: review.review_digest });
    expect(readiness.readiness.ready).toBe(false);
    expect(readiness.readiness_digest).toBe('');
    expect(direct(mock, 'activation-confirm', { envelope: review.envelope, review_digest: review.review_digest, readiness_digest: 'a'.repeat(64) }).status).toBe(409);
    expect(mock.snapshot().status?.cutover.status).toBe('not_submitted');
  });
  it('rejects foreign path, workspace, stale current CAS, root/envelope replacement, and malformed requests', async () => {
    const { mock, stage } = await reserved();
    const foreignRoot = structuredClone(stage); foreignRoot.envelope.request.operation_id = 'OTHER';
    const alteredIntent = structuredClone(stage); alteredIntent.envelope.request.group.name = 'replacement';
    for (const body of [foreignRoot, alteredIntent, { ...stage, current: { ...stage.current, workspace_revision: 'OLD' } }, {}]) expect(direct(mock, 'schema-preview', body).status).toBe(409);
    expect(direct(mock, 'operation-status', { workspace_id: 'OTHER', lineage_id: 'L', operation_id: 'ROOT' }).status).not.toBe(200);
    expect(mock.handle({ method: 'POST', url: path.replace('G1', 'FOREIGN'), body: { ...stage, action: 'reserve' } }).status).toBe(404);
    expect(mock.handle({ method: 'GET', url: `${FIXTURE_PATH}/unexpected/write-groups/G1/lifecycle` }).status).toBe(404);
    expect(mock.handle({ method: 'POST', url: `${FIXTURE_PATH}/write-groups/G1/apply`, body: {} }).status).toBe(404);
    expect(mock.ledger.at(-1)).toMatchObject({ action: 'apply', mutation: true, status: 404 });
    expect(mock.snapshot().status?.schema_operations).toEqual([]);
  });
  it('requires original token, schema ID/digest and final activation kind/future boundary', async () => {
    const { mock, stage } = await reserved(); const preview = await api.schemaPreview('G1', stage);
    expect(direct(mock, 'schema-confirm', { ...stage, token: 'wrong', digest: preview.digest, schema_operation_id: preview.operation_id }).status).toBe(409);
    expect(direct(mock, 'activation-review', { stage, boundary: '2000-01-01T00:00:00Z' }).status).toBe(409);
    const next = await ready('deleted-unapplied');
    expect(direct(next.mock, 'activation-confirm', next.request).status).toBe(409);
  });
  it('rejects mixed-action and extra nested fields rather than stripping them into mock success', async () => {
    const { mock, stage } = await reserved();
    expect(direct(mock, 'reserve', stage).status).not.toBe(200);
    expect(direct(mock, 'operation-status', { workspace_id: 'ws-1', lineage_id: 'L', operation_id: 'ROOT', current: stage.current }).status).not.toBe(200);
    const altered = structuredClone(stage) as typeof stage & { current: typeof stage.current & { extra?: string } };
    altered.current.extra = 'not accepted';
    expect(direct(mock, 'schema-preview', altered).status).not.toBe(200);
    expect(mock.snapshot().status?.schema_operations).toEqual([]);
  });
  it('stop/remove preserve accepted history and forbid stale or blocked mutation', async () => {
    const mock = createRecordingLifecycleMock(); connect(mock);
    const history = (await lifecycleAPI.lifecycle('G1')).history.generations[0].revisions;
    const request = () => ({ workspace_id: 'ws-1', expected_workspace_revision: mock.snapshot().workspaceRevision, expected_group_revision: mock.snapshot().source.revision, expected_connector_revision: 'crev-1' });
    const old = request(); await groupsAPI.disable('G1', old);
    await expect(groupsAPI.remove('G1', old)).rejects.toBeDefined();
    expect((await groupsAPI.remove('G1', request())).group.status).toBe('deleted');
    expect((await lifecycleAPI.lifecycle('G1')).history.generations[0].revisions).toEqual(history);
    const blocked = createRecordingLifecycleMock('blocked-wide');
    expect(blocked.handle({ method: 'DELETE', url: `${FIXTURE_PATH}/write-groups/G1`, body: { ...old, expected_workspace_revision: 'W1', expected_group_revision: 'R1' } }).status).toBe(409);
  });
  it('blocked review and stale scenarios return a repairable contract, not unconditional 200', async () => {
    const { mock, review } = await reviewed('review-blocked-wide');
    expect(review.blockers).toEqual([{ reason_code: 'connector_identity_stale', next_action: 'review_connector' }]);
    expect(direct(mock, 'reserve', { envelope: review.envelope, review_digest: review.review_digest }).status).toBe(409);
    await expect(reviewed('stale-wide')).rejects.toMatchObject({ outcome: 'failed' });
    const stale = createRecordingLifecycleMock(); connect(stale); stale.staleNextLifecycle();
    expect((await lifecycleAPI.lifecycle('G1')).group_revision).toBe('R1-STALE');
  });
  it('snapshots and caller-owned response bodies cannot mutate mock server state', () => {
    const mock = createRecordingLifecycleMock(); const snapshot = mock.snapshot(); snapshot.source.name = 'tampered';
    const reply = mock.handle({ method: 'GET', url: `${FIXTURE_PATH}/write-groups` });
    const data = parseWriteGroupListData((reply.body as { data: unknown }).data)!; data.groups[0].name = 'tampered';
    expect(mock.snapshot().source.name).toBe('Mock saved recording');
  });
  it('records the existing host ordinary preview separately and rejects foreign CAS or ordinary DDL', async () => {
    const mock = createRecordingLifecycleMock(); connect(mock);
    const request = { workspace_id: 'ws-1', expected_workspace_revision: 'W1', expected_group_revision: 'R1', expected_connector_revision: 'crev-1' };
    const preview = await groupsAPI.schemaPreview('G1', request);
    expect(preview.operation_id).toBe('ORDINARY-SCHEMA-1');
    expect(mock.ledger.at(-1)).toMatchObject({ action: 'ordinary-schema-preview', mutation: true, status: 200 });
    await expect(groupsAPI.schemaPreview('G1', { ...request, expected_group_revision: 'FOREIGN' })).rejects.toBeDefined();
    await expect(groupsAPI.schemaApply('G1', { ...request, token: preview.token, operation_id: preview.operation_id })).rejects.toBeDefined();
    expect(mock.snapshot().ordinaryPreviewCount).toBe(1);
    expect(mock.snapshot().status).toBeUndefined();
  });
  it('permits observed ordinary G2 preview only after terminal scheduling and with current G2 CAS', async () => {
    const { mock, request } = await ready('applied-wide');
    const scope = { workspace_id: 'ws-1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'crev-1' };
    await expect(groupsAPI.schemaPreview('G2', scope)).rejects.toBeDefined();
    await api.confirmCutover('G1', request);
    await expect(groupsAPI.schemaPreview('G2', scope)).rejects.toBeDefined();
    expect(await groupsAPI.schemaPreview('G2', { ...scope, expected_workspace_revision: 'W3' })).toMatchObject({ plan_id: 'G2', plan_revision: 'R2' });
    expect(mock.snapshot().target?.status).toBe('ready');
    expect(mock.snapshot().status?.reserve.result.group.applied_revision).toBe('');
  });
});
