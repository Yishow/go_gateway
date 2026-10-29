import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { AxiosError, type AxiosAdapter } from 'axios';
import { afterEach, describe, expect, it } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI as groups } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDatabaseAPI as database } from '@/services/studioV2WorkspaceDatabase';
import { studioV2WorkspaceMeasurementsAPI as measurements } from '@/services/studioV2WorkspaceMeasurements';
import { listRecordingTestWriteRecovery } from '@/services/recordingTestWriteRecovery';
import type { WriteGroupCreateRequest } from '@/types/studioV2WriteGroup';
import type { WriteGroupTestWriteConfirmation, WriteGroupTestWriteOperation } from '@/types/studioV2WriteGroupTestWrite';
import { createRecordingAuthoringMock } from '../../e2e/fixtures/recording-authoring-api';
import { AUTHORING_PATH as PATH, authoringKind, authoringScenarios, initialAuthoringGroups,
  type AuthoringScenario } from '../../e2e/fixtures/recording-authoring-state';

// Lower-layer fixture viability only. These tests cannot replace browser, HTTP ownership or SQL acceptance.
const originalAdapter = studioV2DatalinkApi.defaults.adapter;
const originalBase = studioV2DatalinkApi.defaults.baseURL;
function connect(scenario: AuthoringScenario) {
  const mock = createRecordingAuthoringMock(scenario);
  studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
  const adapter: AxiosAdapter = async config => {
    const reply = mock.handle({ method: (config.method ?? 'get').toUpperCase(),
      url: new URL(studioV2DatalinkApi.getUri(config), 'http://fixture.test').href,
      body: typeof config.data === 'string' ? JSON.parse(config.data) : config.data });
    if (reply.transport === 'abort') throw new AxiosError('Test-only lost response', 'ERR_NETWORK', config);
    const response = { data: JSON.parse(JSON.stringify(reply.body)), status: reply.status, statusText: 'Mock response', config, headers: {} };
    if (reply.status >= 400) throw new AxiosError('Test-only rejected request', 'ERR_BAD_RESPONSE', config, undefined, response);
    return response;
  };
  studioV2DatalinkApi.defaults.adapter = adapter;
  return mock;
}
afterEach(() => {
  studioV2DatalinkApi.defaults.adapter = originalAdapter;
  studioV2DatalinkApi.defaults.baseURL = originalBase;
});
function request(): WriteGroupCreateRequest {
  const saved = initialAuthoringGroups('saved-state')[0];
  return { workspace_id: 'ws-1', expected_workspace_revision: 'W1', expected_connector_revision: 'crev-1',
    group: { workspace_id: 'ws-1', name: 'Authored state', members: saved.members, destination: saved.destination,
      layout_version: 'managed_long_v1', storage_mode: 'plain', recording_profile: saved.recording_profile,
      row_policy: saved.row_policy, write_policy: saved.write_policy } };
}
function confirmation(preview: Pick<WriteGroupTestWriteConfirmation, 'token' | 'operation_id'>) {
  return { token: preview.token, operation_id: preview.operation_id };
}
function capability(scenario: AuthoringScenario, revision = 'crev-1') {
  return authoringKind(scenario) === 'sqlserver' ? database.getSQLServerRecordingCapability('conn-1', revision)
    : database.getMySQLRecordingCapability('conn-1', revision);
}
const dynamicScenarios = ['new-mysql', 'new-mariadb', 'new-sqlserver', 'test-write-sqlserver', 'saved-unsupported-wide'] as const;

describe('authoring browser mock through real service parsers', () => {
  it.each(authoringScenarios)('parses %s saved identities without changing meaning', async scenario => {
    const mock = connect(scenario);
    const list = await groups.list();
    expect(list).toEqual({ workspace_id: 'ws-1', workspace_revision: 'W1', groups: initialAuthoringGroups(scenario) });
    for (const saved of list.groups) {
      expect((await groups.get(saved.id)).group).toEqual(saved);
      expect(await groups.readiness(saved.id)).toMatchObject({ group_id: saved.id, schema_ready: true });
    }
    expect(mock.snapshot().groups).toEqual(mock.snapshot().frozenGroups);
    expect(mock.ledger.every(call => call.status === 200 && !call.mutation)).toBe(true);
  });

  it.each(dynamicScenarios)('parses %s matching capabilities and rejects stale/foreign replies', async scenario => {
    const mock = connect(scenario);
    expect(await capability(scenario)).toEqual(mock.snapshot().capability);
    await expect(capability(scenario, 'crev-2')).rejects.toThrow(/saved connector/);
    expect(mock.snapshot().capability.connector_revision).toBe('crev-1');
    mock.setCapability({ connector_revision: 'crev-2' });
    expect(await capability(scenario, 'crev-2')).toMatchObject({ connector_revision: 'crev-2' });
    mock.setCapability({ connector_id: 'foreign' });
    await expect(capability(scenario, 'crev-2')).rejects.toThrow(/saved connector/);
    expect(mock.ledger.every(call => call.status === 200 && !call.mutation)).toBe(true);
  });

  it.each(dynamicScenarios)('retains %s disconnected capability as an explicit negative', async scenario => {
    const mock = connect(scenario);
    mock.setCapability({ connected: false, metadata_status: 'unavailable', supports_managed_long: false,
      supports_native_partition: false, supports_group_test_write: false, supports_missing_repair: false });
    expect(await capability(scenario)).toMatchObject({ connected: false, metadata_status: 'unavailable', supports_managed_long: false });
  });

  it.each(['new-sqlite', 'new-postgres', 'new-mysql', 'new-mariadb', 'new-sqlserver'] as const)(
    '%s ensure succeeds and yields parseable schema preview, not just an attempted request', async scenario => {
      const mock = connect(scenario), input = request();
      const saved = await groups.ensureCircuitProfile(input);
      expect(saved.workspace_revision).toBe('W2');
      expect(saved.group.members).toEqual(input.group.members);
      expect(saved.group).toMatchObject({ id: 'G1', status: 'draft', applied_revision: '', layout_version: 'managed_long_v1' });
      expect(await groups.readiness('G1')).toMatchObject({ schema_ready: false, ready: false });
      const preview = await groups.schemaPreview('G1', { workspace_id: 'ws-1', expected_workspace_revision: 'W2',
        expected_group_revision: 'R1', expected_connector_revision: 'crev-1' });
      expect(preview).toMatchObject({ operation_id: 'SCHEMA-1', plan_id: 'G1', plan_revision: 'R1' });
      expect(mock.ledger.filter(call => call.mutation).map(call => [call.action, call.status])).toEqual([
        ['ensure', 200], ['schema-preview', 200],
      ]);
      expect(mock.snapshot()).toMatchObject({ schemaPreviewCount: 1, insertCount: 0 });
    });

  it('keeps unsupported saved wide readable while refusing conversion', async () => {
    const mock = connect('saved-unsupported-wide'), before = mock.snapshot();
    expect((await groups.list()).groups[0]).toMatchObject({ layout_version: 'managed_wide_profile_v1' });
    await expect(groups.ensureCircuitProfile(request())).rejects.toMatchObject({ response: { status: 409 } });
    expect(mock.snapshot().groups).toEqual(before.groups);
  });

  it.each(['workspace', 'workspace-revision', 'connector-revision', 'source-device', 'source-tag'])(
    'rejects invalid %s on ensure without saving', async field => {
      const mock = connect('new-sqlite'), input = request();
      if (field === 'workspace') input.workspace_id = 'foreign';
      if (field === 'workspace-revision') input.expected_workspace_revision = 'old';
      if (field === 'connector-revision') input.expected_connector_revision = 'old';
      if (field === 'source-device') input.group.members[0].device_id = 'dev-2';
      if (field === 'source-tag') input.group.members[0].tag_id = 'foreign';
      await expect(groups.ensureCircuitProfile(input)).rejects.toBeInstanceOf(AxiosError);
      expect(mock.snapshot().groups).toEqual([]);
      expect(mock.ledger).toHaveLength(1);
      expect(mock.ledger[0]).toMatchObject({ action: 'ensure', mutation: true });
      expect(mock.ledger[0].status).toBeGreaterThanOrEqual(400);
    });

  it('rejects old schema revision and logs unknown mutations as errors', async () => {
    const mock = connect('saved-state');
    await expect(groups.schemaPreview('G1', { workspace_id: 'ws-1', expected_workspace_revision: 'W1',
      expected_group_revision: 'old', expected_connector_revision: 'crev-1' })).rejects.toMatchObject({ response: { status: 409 } });
    await expect(studioV2DatalinkApi.post('/studio-v2/workspace/not-a-route', {})).rejects.toMatchObject({ response: { status: 404 } });
    expect(mock.ledger.map(call => [call.mutation, call.status])).toEqual([[true, 409], [true, 404]]);
    expect(mock.snapshot().groups).toEqual(mock.snapshot().frozenGroups);
  });

  it('scopes measurement listing, fails closed, and permits retry without inventing definitions', async () => {
    const mock = connect('saved-counter');
    expect(await measurements.list('dev-1')).toEqual(mock.snapshot().measurements);
    expect(await measurements.list('dev-2')).toEqual([]);
    await expect(measurements.list('foreign')).rejects.toMatchObject({ response: { status: 400 } });
    mock.setMeasurementError('dev-1', true);
    await expect(measurements.list('dev-1')).rejects.toMatchObject({ response: { status: 503 } });
    mock.setMeasurementError('dev-1', false);
    expect(await measurements.list('dev-1')).toHaveLength(1);
    expect(mock.ledger.every(call => !call.mutation)).toBe(true);
  });
});

type Outcome = Pick<WriteGroupTestWriteOperation, 'status' | 'write_outcome' | 'cleanup_status'>;
const outcomes: Outcome[] = [
  { status: 'pending' }, { status: 'running' },
  { status: 'unknown', write_outcome: 'unknown', cleanup_status: 'unknown' },
  { status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'failed' },
  { status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'unknown' },
  { status: 'partial', write_outcome: 'written_unverified', cleanup_status: 'unknown' },
  { status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned' },
  { status: 'failed', write_outcome: 'failed', cleanup_status: 'not_attempted' },
];
describe('authoring test-write mock recovery contract', () => {
  it.each(outcomes)('parses $status/$cleanup_status and lists only unresolved O1', async outcome => {
    const mock = connect('test-write');
    expect((await listRecordingTestWriteRecovery('ws-1')).operations).toEqual([]);
    mock.setTestWriteOutcome(outcome.status, outcome.write_outcome, outcome.cleanup_status);
    const preview = await groups.testWritePreview('G1');
    expect((await listRecordingTestWriteRecovery('ws-1')).operations).toEqual([]);
    const operation = await groups.testWrite('G1', confirmation(preview));
    expect(operation).toMatchObject({ operation_id: 'O1', ...outcome });
    expect(await groups.testWriteOperation('O1')).toEqual(operation);
    const unresolved = !['cleaned', 'not_attempted'].includes(outcome.cleanup_status ?? '');
    expect((await listRecordingTestWriteRecovery('ws-1')).operations).toHaveLength(unresolved ? 1 : 0);
    expect(mock.snapshot()).toMatchObject({ insertCount: 1, repairCount: 0, testWritePreviewCount: 1 });
    expect(mock.ledger.every(call => call.status === 200)).toBe(true);
  });

  it('recovers lost ACK through the actual services; no replacement token or second insert', async () => {
    const mock = connect('test-write'), preview = await groups.testWritePreview('G1');
    mock.loseNext('test-write-confirm');
    await expect(groups.testWrite('G1', confirmation(preview))).rejects.toMatchObject({ code: 'ERR_NETWORK' });
    expect((await listRecordingTestWriteRecovery('ws-1')).operations[0]).toMatchObject({ operation_id: 'O1', status: 'unknown' });
    expect(await groups.testWriteOperation('O1')).toMatchObject({ operation_id: 'O1', status: 'unknown' });
    await expect(groups.testWritePreview('G1')).rejects.toBeInstanceOf(AxiosError);
    await expect(groups.testWrite('G1', confirmation(preview))).rejects.toBeInstanceOf(AxiosError);
    expect(mock.snapshot()).toMatchObject({ insertCount: 1, testWritePreviewCount: 1 });
  });

  it('rejects foreign list scope and supplied operation scope without inventing query requirements', async () => {
    const mock = connect('test-write');
    // Production GET has session-bound scope and no workspace query; parser binds returned data.
    await expect(listRecordingTestWriteRecovery('foreign')).rejects.toThrow(/unavailable/);
    expect(mock.ledger[0].url).not.toContain('workspace_id');
    for (const path of ['test-writes', 'O1']) {
      expect(mock.handle({ method: 'GET', url: `${PATH}/database-operations/${path}?workspace_id=foreign` }).status).toBe(403);
    }
    await expect(groups.testWriteOperation('foreign')).rejects.toMatchObject({ response: { status: 404 } });
    mock.rejectNext('test-write-list');
    await expect(listRecordingTestWriteRecovery('ws-1')).rejects.toMatchObject({ response: { status: 503 } });
    mock.rejectNext('operation-status');
    await expect(groups.testWriteOperation('O1')).rejects.toMatchObject({ response: { status: 503 } });
    expect(mock.snapshot().insertCount).toBe(0);
  });

  it.each(['failed', 'unknown'] as const)('allows SQL Server %s cleanup once, same O1/token and no insertion', async cleanup => {
    const mock = connect('test-write-sqlserver');
    mock.setTestWriteOutcome('succeeded', 'written_verified', cleanup);
    const preview = await groups.testWritePreview('G1');
    expect(preview.target.dialect).toBe('sqlserver');
    await groups.testWrite('G1', confirmation(preview));
    for (const changed of [{ token: 'foreign' }, { operation_id: 'foreign' }]) {
      await expect(groups.testWrite('G1', { ...confirmation(preview), ...changed, action: 'cleanup_repair' })).rejects.toBeInstanceOf(AxiosError);
    }
    expect(mock.snapshot()).toMatchObject({ insertCount: 1, repairCount: 0 });
    expect(await groups.testWrite('G1', { ...confirmation(preview), action: 'cleanup_repair' })).toMatchObject({
      operation_id: 'O1', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned',
    });
    expect((await listRecordingTestWriteRecovery('ws-1')).operations).toEqual([]);
    await expect(groups.testWrite('G1', { ...confirmation(preview), action: 'cleanup_repair' })).rejects.toBeInstanceOf(AxiosError);
    expect(mock.snapshot()).toMatchObject({ insertCount: 1, repairCount: 1, testWritePreviewCount: 1 });
  });

  it('never permits SQLite cleanup repair even with a verified write and retained token', async () => {
    const mock = connect('test-write');
    mock.setTestWriteOutcome('succeeded', 'written_verified', 'failed');
    const preview = await groups.testWritePreview('G1');
    expect(preview.target.dialect).toBe('sqlite');
    await groups.testWrite('G1', confirmation(preview));
    await expect(groups.testWrite('G1', { ...confirmation(preview), action: 'cleanup_repair' })).rejects.toBeInstanceOf(AxiosError);
    expect(mock.snapshot()).toMatchObject({ insertCount: 1, repairCount: 0 });
  });
});

describe('guided summary read in the authoring fixture', () => {
  it('returns explicit unavailable evidence only for the exact current workspace/revision/device without changing saved configuration', async () => {
    const mock = connect('saved-state'), before = mock.snapshot();
    const scope = { workspace_id: 'ws-1', workspace_revision: 'W1', device_id: 'dev-1' };
    expect(await studioV2RuntimeContextAPI.getRuleEvidence(scope)).toEqual({ ok: true, value: {
      contract_version: 1, ...scope, observed_at: null, freshness: 'unavailable', source: null, groups: [], locator: null,
    } });
    expect(mock.ledger).toHaveLength(1);
    expect(mock.ledger[0]).toMatchObject({ method: 'GET', action: 'runtime-context', status: 200, mutation: false });
    expect(mock.snapshot()).toEqual(before);
  });
  it.each(['workspace', 'revision', 'device', 'missing', 'extra', 'duplicate'] as const)(
    'refuses %s runtime scope without a new proof or mutation', async kind => {
      const mock = connect('saved-state');
      const params = new URLSearchParams({ workspace_id: 'ws-1', workspace_revision: 'W1', device_id: 'dev-1' });
      if (kind === 'workspace') params.set('workspace_id', 'foreign');
      if (kind === 'revision') params.set('workspace_revision', 'old');
      if (kind === 'device') params.set('device_id', 'foreign');
      if (kind === 'missing') params.delete('device_id');
      if (kind === 'extra') params.set('extra', 'unsafe');
      if (kind === 'duplicate') params.append('device_id', 'dev-2');
      const reply = mock.handle({ method: 'GET', url: `${PATH}/runtime-context?${params}` });
      expect(reply.status).toBeGreaterThanOrEqual(400);
      expect(mock.ledger[0]).toMatchObject({ action: 'runtime-context', mutation: false });
      expect(mock.snapshot().groups).toEqual(mock.snapshot().frozenGroups);
    });
});
