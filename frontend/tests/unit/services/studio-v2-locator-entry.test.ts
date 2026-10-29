import { afterEach, describe, expect, it, vi } from 'vitest';
import { readStudioLocatorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const locator: EntityLocator = { step: 1, workspace_id: 'W', device_id: 'D', rule_id: null, mapping_id: null, group_id: null, group_revision: null };
const at = '2026-10-29T00:00:00Z';
function payload() {
  return { id: 'W', kind: 'single', status: 'ready', readiness_summary: null, ordered_device_ids: ['D'], created_at: at, updated_at: at,
    locator_entry: { workspace_revision: '', hydration_scope: { complete_workspace: false, device_ids: ['D'], rule_ids: [], mapping_ids: [], group_ids: [] },
      devices: [{ id: 'D', name: 'Meter', description: '', protocol: 'modbus_tcp', status: 'active', connection_config: '{"host":"private.local","password":"secret"}',
        last_test_at: null, last_test_success: null, last_test_error: '', created_at: at, updated_at: at }], rule_heads: {}, group_heads: {}, revalidated_mapping_ids: [] } };
}
function respond(value: unknown = payload()) { vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: value } }); }
afterEach(() => vi.resetAllMocks());
describe('detached native locator entry transport', () => {
  it('retains explicit empty observed revision and editable config only inside the detached entry', async () => {
    respond(); const signal = new AbortController().signal; const result = await readStudioLocatorEntry(locator, signal);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace', { params: { entry: 'locator', workspace_id: 'W', step: '1', device_id: 'D' }, signal });
    expect(result).toMatchObject({ ok: true, value: { selection: { state: 'selected', target: { workspaceId: 'W', workspaceRevision: '', deviceId: 'D', step: 1 } } } });
    if (!result.ok) throw new Error('entry unavailable');
    expect(result.value.entry.devices![0].connection_config).toContain('private.local');
    expect(JSON.stringify([result.value.locator, result.value.selection, result.value.workspace])).not.toMatch(/private|secret|password/);
    expect(result.value.entry).not.toHaveProperty('source_rules');
  });
  it.each(['missing-entry', 'missing-revision', 'null-revision', 'complete-workspace', 'missing-heads', 'unknown-entry', 'unknown-workspace', 'duplicate-device', 'unrequested-rules'])('fails closed for malformed %s', async (kind) => {
    const value = payload() as Record<string, unknown>; const entry = value.locator_entry as Record<string, unknown>;
    if (kind === 'missing-entry') delete value.locator_entry;
    if (kind === 'missing-revision') delete entry.workspace_revision;
    if (kind === 'null-revision') entry.workspace_revision = null;
    if (kind === 'complete-workspace') (entry.hydration_scope as Record<string, unknown>).complete_workspace = true;
    if (kind === 'missing-heads') delete entry.rule_heads;
    if (kind === 'unknown-entry') entry.database_config = { password: 'secret' };
    if (kind === 'unknown-workspace') value.modbus_share = { endpoint: 'private.local' };
    if (kind === 'duplicate-device') entry.devices = [...(entry.devices as unknown[]), ...(entry.devices as unknown[])];
    if (kind === 'unrequested-rules') entry.source_rules = [];
    respond(value); expect(await readStudioLocatorEntry(locator)).toEqual({ ok: false, reason: 'unavailable' });
  });
  it.each(['workspace', 'device', 'scope'])('rejects foreign or mixed %s identity', async (kind) => {
    const value = payload(); if (kind === 'workspace') value.id = 'FOREIGN';
    if (kind === 'device') value.locator_entry.devices[0].id = 'OTHER';
    if (kind === 'scope') value.locator_entry.hydration_scope.device_ids = ['OTHER'];
    respond(value); expect((await readStudioLocatorEntry(locator)).ok).toBe(false);
  });
  it('captures caller IDs before await and never binds the response to a mutated caller', async () => {
    let done!: (value: unknown) => void; vi.mocked(studioV2DatalinkApi.get).mockImplementationOnce(() => new Promise(resolve => { done = resolve; }));
    const selected = { ...locator }; const pending = readStudioLocatorEntry(selected); selected.workspace_id = 'OTHER'; selected.device_id = 'OTHER';
    done({ data: { success: true, data: payload() } }); expect(await pending).toMatchObject({ ok: true, value: { locator: { workspace_id: 'W', device_id: 'D' } } });
  });
  it.each([403, 409, 503])('maps HTTP %s to safe typed failure without raw message', async (status) => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { status, data: { error: { message: 'password=secret /private/driver.log' } } } });
    expect(await readStudioLocatorEntry(locator)).toEqual({ ok: false, reason: status === 403 ? 'foreign-scope' : status === 409 ? 'stale-reference' : 'unavailable' });
  });
  it('rejects invalid locator before HTTP and does not fall back to bootstrap', async () => {
    expect(await readStudioLocatorEntry({ ...locator, device_id: 'http://private.local' })).toEqual({ ok: false, reason: 'unavailable' });
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
});

import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
const groupLocator: EntityLocator = { ...locator, step: 4, device_id: null, group_id: 'G', group_revision: 'R1' };
function groupPayload(editor = false) {
  const value = payload();
  return { ...value, locator_entry: { workspace_revision: 'W3', hydration_scope: { complete_workspace: false, device_ids: [], rule_ids: [], mapping_ids: [], group_ids: ['G'] },
    rule_heads: {}, group_heads: { G: 'R3' }, revalidated_mapping_ids: [],
    write_groups: { workspace_id: 'W', workspace_revision: 'W3', groups: [{ id: 'G', workspace_id: 'W', revision: 'R3', applied_revision: 'R1', name: 'Group', status: 'disabled', members: [],
      destination: { connector_id: 'C', connector_revision: 'C1', database: 'native.db', table_schema: '', table_name: 'meter', storage_strategy: 'managed' },
      row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0 }, write_policy: {}, migration: {}, created_at: at, updated_at: at }] },
    ...(editor ? { editor_hydrated: true, devices: [], source_rules: [], mappings: [], database_config: null, database_targets: [] } : {}) } };
}
describe('explicit native editor snapshot stays bound to initial current revisions', () => {
  it('keeps historical R1 informational and sends independently observed current R3/W3 for editor hydration', async () => {
    respond(groupPayload(true)); const result = await readStudioLocatorEditorEntry(groupLocator, { workspaceRevision: 'W3', groupRevision: 'R3' });
    expect(result).toMatchObject({ ok: true, value: { entry: { editor_hydrated: true, database_config: null, database_targets: [] }, selection: { target: { groupRevision: 'R3', workspaceRevision: 'W3' }, history: { revisionHint: 'R1', relation: 'different' } } } });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace', expect.objectContaining({ params: { entry: 'locator', workspace_id: 'W', step: '4', group_id: 'G', group_revision: 'R1', editor: '1', expected_workspace_revision: 'W3', expected_group_revision: 'R3' } }));
  });
  it.each([undefined, { workspaceRevision: 'W3' }, { groupRevision: 'R3' }, { workspaceRevision: 'W3', groupRevision: '' }])('rejects missing current expectation %j before any HTTP', async expected => {
    expect(await readStudioLocatorEditorEntry(groupLocator, expected as { workspaceRevision: string; groupRevision: string })).toEqual({ ok: false, reason: 'unavailable' });
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('preserves explicit empty observed workspace revision rather than fabricating it when missing', async () => {
    const value = groupPayload(true); value.locator_entry.workspace_revision = ''; value.locator_entry.write_groups.workspace_revision = ''; respond(value);
    expect((await readStudioLocatorEditorEntry(groupLocator, { workspaceRevision: '', groupRevision: 'R3' })).ok).toBe(true);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace', expect.objectContaining({ params: expect.objectContaining({ expected_workspace_revision: '' }) }));
  });
  it.each(['unsolicited', 'missing-editor', 'group-changed', 'workspace-changed', 'deleted'])('fails closed for editor mismatch %s', async kind => {
    const value = groupPayload(kind !== 'missing-editor');
    if (kind === 'group-changed') { value.locator_entry.group_heads.G = 'R4'; value.locator_entry.write_groups.groups[0].revision = 'R4'; }
    if (kind === 'workspace-changed') { value.locator_entry.workspace_revision = 'W4'; value.locator_entry.write_groups.workspace_revision = 'W4'; }
    if (kind === 'deleted') value.locator_entry.write_groups.groups[0].status = 'deleted';
    respond(value); const result = kind === 'unsolicited' ? await readStudioLocatorEntry(groupLocator) : await readStudioLocatorEditorEntry(groupLocator, { workspaceRevision: 'W3', groupRevision: 'R3' });
    expect(result.ok).toBe(false);
  });
  it('retains canonical tombstone for cold selection without substituting another group', async () => {
    const value = groupPayload(); value.locator_entry.write_groups.groups[0].status = 'deleted'; respond(value);
    expect(await readStudioLocatorEntry(groupLocator)).toMatchObject({ ok: true, value: { workspace: { id: 'W' }, selection: { state: 'unavailable', reason: 'target_unavailable' } } });
  });
});

const nativeRule = () => ({ id: 'R', workspace_id: 'W', device_id: 'D', revision_id: 'R2', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: '', enabled: true, locked: false, origin: 'manual', skipped_addresses: [], created_at: at, updated_at: at, share_start_register: null, share_stride: null });
const nativeMapping = () => ({ id: 'M', workspace_id: 'W', device_id: 'D', rule_id: 'R', point_id: 'P', tag_id: 'T', address: '40001', tag_key: 'meter.value', display_name: 'Value', unit: '', target_type: 'float64', scale: 1, offset: 0, enabled: true, created_at: at, updated_at: at });
function sourcePayload(step: 2 | 3) {
  const value = payload(); return { ...value, locator_entry: { ...value.locator_entry, source_rules: [nativeRule()], rule_heads: { R: 'R2' },
    hydration_scope: { ...value.locator_entry.hydration_scope, rule_ids: ['R'], mapping_ids: step === 3 ? ['M'] : [] },
    ...(step === 3 ? { mappings: [nativeMapping()] } : {}), revalidated_mapping_ids: step === 3 ? ['M'] : [] } };
}
function editorClosure() {
  const value = groupPayload(true); const member = { device_id: 'D', point_id: 'P', tag_id: 'T', source_revision: 'R2', mapping_revision: 'M2', target_column: 'value', required: true };
  return { ...value, locator_entry: { ...value.locator_entry, devices: payload().locator_entry.devices, source_rules: [nativeRule()], mappings: [nativeMapping()],
    hydration_scope: { complete_workspace: false, device_ids: ['D'], rule_ids: ['R'], mapping_ids: ['M'], group_ids: ['G'] }, rule_heads: { R: 'R2' }, revalidated_mapping_ids: ['M'],
    write_groups: { ...value.locator_entry.write_groups, groups: [{ ...value.locator_entry.write_groups.groups[0], members: [member, { ...member, device_id: 'OLD', point_id: 'OLDP', tag_id: 'OLDT' }] }] },
    database_config: { id: 'C', identity_revision: 'C1', setup_revision: 'W3', workspace_id: 'W', kind: 'sqlite', name: 'Database', host: '', port: 0, database: '/private/meter.sqlite', username: '', schema: '', table: 'meter', write_mode: 'insert', write_interval_seconds: 15, timestamp_column: 'time', status: 'ready', created_at: at, updated_at: at },
    database_targets: [{ id: 'DBT', workspace_id: 'W', point_id: 'P', tag_id: 'T', column_name: 'value', enabled: true, created_at: at, updated_at: at }] } };
}
describe('exact native entity relationships and editor closure', () => {
  it.each([2, 3] as const)('selects exact Step%s ancestry without same-name/address fallback', async step => {
    respond(sourcePayload(step)); const selected = { ...locator, step, rule_id: step === 2 ? 'R' : null, mapping_id: step === 3 ? 'M' : null };
    expect(await readStudioLocatorEntry(selected)).toMatchObject({ ok: true, value: { selection: { state: 'selected', target: { deviceId: 'D', ruleId: 'R', ruleRevision: 'R2', ...(step === 3 ? { mappingId: 'M' } : {}) } } } });
  });
  it.each(['rule-head', 'mapping-ancestor', 'mapping-currentness', 'rule-foreign'])('rejects inconsistent %s source proof', async kind => {
    const value = sourcePayload(3); if (kind === 'rule-head') value.locator_entry.rule_heads.R = 'R3';
    if (kind === 'mapping-ancestor') value.locator_entry.mappings![0].rule_id = 'OTHER';
    if (kind === 'mapping-currentness') value.locator_entry.revalidated_mapping_ids = [];
    if (kind === 'rule-foreign') value.locator_entry.source_rules[0].workspace_id = 'OTHER';
    respond(value); expect((await readStudioLocatorEntry({ ...locator, step: 3, mapping_id: 'M' })).ok).toBe(false);
  });
  it('hydrates only proven native editor closure and retains unhydrated orphan members without fabricating them', async () => {
    respond(editorClosure()); const result = await readStudioLocatorEditorEntry(groupLocator, { workspaceRevision: 'W3', groupRevision: 'R3' });
    expect(result.ok).toBe(true); if (!result.ok) throw new Error('editor unavailable');
    expect(result.value.entry.write_groups!.groups[0].members).toHaveLength(2); expect(result.value.entry.devices!.map(row => row.id)).toEqual(['D']);
    expect(result.value.entry.database_config!.database).toBe('/private/meter.sqlite');
    expect(JSON.stringify([result.value.locator, result.value.selection, result.value.workspace])).not.toContain('/private');
  });
  it.each(['foreign-target', 'unrelated-target', 'unrelated-mapping', 'raw-diagnostic'])('rejects unsafe/mixed native editor %s', async kind => {
    const value = editorClosure(); if (kind === 'foreign-target') value.locator_entry.database_targets[0].workspace_id = 'FOREIGN';
    if (kind === 'unrelated-target') value.locator_entry.database_targets[0].point_id = 'FOREIGN';
    if (kind === 'unrelated-mapping') value.locator_entry.mappings[0].point_id = 'FOREIGN';
    if (kind === 'raw-diagnostic') Object.assign(value.locator_entry.database_config, { last_write_error: 'password=secret /private/error' });
    respond(value); expect((await readStudioLocatorEditorEntry(groupLocator, { workspaceRevision: 'W3', groupRevision: 'R3' })).ok).toBe(false);
  });
  it('never accepts readiness proof or raw device diagnostics in the entry envelope', async () => {
    const value = payload(); value.locator_entry.devices[0].last_test_error = 'password=secret'; respond(value);
    expect(await readStudioLocatorEntry(locator)).toEqual({ ok: false, reason: 'unavailable' });
    respond({ ...payload(), readiness_summary: { ready: true } }); expect(await readStudioLocatorEntry(locator)).toEqual({ ok: false, reason: 'unavailable' });
  });
});

it('captures editor expectations before await and rejects aborted read results', async () => {
  let done!: (value: unknown) => void; vi.mocked(studioV2DatalinkApi.get).mockImplementationOnce(() => new Promise(resolve => { done = resolve; }));
  const expected = { workspaceRevision: 'W3', groupRevision: 'R3' }; const controller = new AbortController();
  const pending = readStudioLocatorEditorEntry(groupLocator, expected, controller.signal); expected.groupRevision = 'R9'; expected.workspaceRevision = 'W9';
  done({ data: { success: true, data: groupPayload(true) } }); expect((await pending).ok).toBe(true);
  respond(payload()); controller.abort(); expect(await readStudioLocatorEntry(locator, controller.signal)).toEqual({ ok: false, reason: 'unavailable' });
});
it('does not permit explicit editor hydration on a source/device locator', async () => {
  expect(await readStudioLocatorEditorEntry(locator, { workspaceRevision: 'W3', groupRevision: 'R3' })).toEqual({ ok: false, reason: 'unavailable' });
  expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
});

import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
const metadata = () => ({ workspace_id: 'W', connector_id: 'C', connector_revision: 'C1', database: 'native.db', schema: '', table: 'meter', inspection_status: 'missing', columns: [] });
describe('locator-origin native database reads remain scope-bound', () => {
  it('forwards observed workspace scope without changing connector/group CAS parameters', async () => {
    respond(metadata()); await studioV2WorkspaceDatabaseAPI.getMetadata('C1', { groupId: 'G', groupRevision: 'R3' }, { workspaceId: 'W', workspaceRevision: '' });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-metadata', { params: { expected_connector_revision: 'C1', group_id: 'G', expected_group_revision: 'R3', locator_workspace_id: 'W', locator_workspace_revision: '' } });
  });
  it('does not accept metadata bound to a different workspace', async () => {
    respond({ ...metadata(), workspace_id: 'FOREIGN' });
    await expect(studioV2WorkspaceDatabaseAPI.getMetadata('C1', undefined, { workspaceId: 'W', workspaceRevision: 'W3' })).rejects.toThrow();
  });
  it('rejects unsafe/missing native scope before HTTP rather than falling back to ordinary metadata', async () => {
    await expect(studioV2WorkspaceDatabaseAPI.getMetadata('C1', undefined, { workspaceId: '/private', workspaceRevision: 'W3' })).rejects.toThrow();
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
});

import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { listRecordingTestWriteRecovery } from '@/services/recordingTestWriteRecovery';
const pendingOperation = () => ({ operation_id: 'O', workspace_id: 'W', workspace_revision: 'OLD', group_id: 'G', group_revision: 'R1', connector_id: 'C', connector_revision: 'C1', target_fingerprint: 'a'.repeat(64), action: 'test_write', status: 'running', created_at: at, updated_at: at });
describe('locator-origin existing operation reads preserve historical provenance', () => {
  it('forwards current read scope but does not relabel an older durable operation revision', async () => {
    respond(pendingOperation()); const result = await studioV2WorkspaceWriteGroupsAPI.testWriteOperation('O', { workspaceId: 'W', workspaceRevision: 'CURRENT' });
    expect(result.operation_id).toBe('O'); expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/O', { params: { locator_workspace_id: 'W', locator_workspace_revision: 'CURRENT' } });
    respond({ workspace_id: 'W', complete: true, operations: [pendingOperation()] });
    const list = await listRecordingTestWriteRecovery('W', { workspaceId: 'W', workspaceRevision: 'CURRENT' }); expect(list.operations[0].workspace_revision).toBe('OLD');
    expect(studioV2DatalinkApi.get).toHaveBeenLastCalledWith('/studio-v2/workspace/database-operations/test-writes', { params: { locator_workspace_id: 'W', locator_workspace_revision: 'CURRENT' } });
  });
  it('rejects foreign operation workspace and malformed listing scope without fallback', async () => {
    respond({ ...pendingOperation(), workspace_id: 'OTHER' });
    await expect(studioV2WorkspaceWriteGroupsAPI.testWriteOperation('O', { workspaceId: 'W', workspaceRevision: 'CURRENT' })).rejects.toThrow();
    vi.mocked(studioV2DatalinkApi.get).mockClear(); await expect(listRecordingTestWriteRecovery('W', { workspaceId: 'OTHER', workspaceRevision: 'CURRENT' })).rejects.toThrow();
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
});

const nativeCapability = (kind: 'mysql' | 'sqlserver') => ({ workspace_id: 'W', connector_id: 'C', connector_revision: 'C1', database: 'native', server_family: 'mysql', server_version: '8', capability_revision: 'CAP', metadata_status: 'verified', connected: true, transactional_engine: true,
  supports_managed_long: false, supports_native_partition: false, supports_group_test_write: false, supports_missing_repair: false, reasons: [],
  ...(kind === 'sqlserver' ? { schema: 'dbo', server_edition: 'dev', engine_edition: 3, compatibility_level: 160, permissions: Object.fromEntries(['verified', 'database_view_definition', 'schema_view_definition', 'create_table', 'alter_schema', 'select_schema', 'insert_schema', 'delete_schema', 'alter_any_dataspace'].map(key => [key, false])) } : {}) });
it.each(['mysql', 'sqlserver'] as const)('forwards scoped %s capability without weakening connector revision matching', async kind => {
  const method = kind === 'mysql' ? studioV2WorkspaceDatabaseAPI.getMySQLRecordingCapability : studioV2WorkspaceDatabaseAPI.getSQLServerRecordingCapability;
  respond(nativeCapability(kind)); await method('C', 'C1', { workspaceId: 'W', workspaceRevision: 'CURRENT' });
  expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/recording-plans/capabilities', { params: { kind, connector_id: 'C', expected_connector_revision: 'C1', locator_workspace_id: 'W', locator_workspace_revision: 'CURRENT' } });
  respond({ ...nativeCapability(kind), workspace_id: 'FOREIGN' }); await expect(method('C', 'C1', { workspaceId: 'W', workspaceRevision: 'CURRENT' })).rejects.toThrow();
  respond({ ...nativeCapability(kind), connector_revision: 'OTHER' }); await expect(method('C', 'C1', { workspaceId: 'W', workspaceRevision: 'CURRENT' })).rejects.toThrow();
});

import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
it('forwards only explicit locator measurement read scope and rejects foreign response ownership', async () => {
  respond([]); await studioV2WorkspaceMeasurementsAPI.list('D', { workspaceId: 'W', workspaceRevision: '' });
  expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/measurements', { params: { device_id: 'D', locator_workspace_id: 'W', locator_workspace_revision: '' } });
  respond([{ id: 'MEASURE', workspace_id: 'OTHER', device_id: 'D' }]); await expect(studioV2WorkspaceMeasurementsAPI.list('D', { workspaceId: 'W', workspaceRevision: 'CURRENT' })).rejects.toThrow();
});

it('does not turn a missing scoped measurement body into a proven empty list or expose native driver errors', async () => {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true } });
  await expect(studioV2WorkspaceMeasurementsAPI.list('D', { workspaceId: 'W', workspaceRevision: 'W1' })).rejects.toThrow('Scoped settings read unavailable.');
  vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce(new Error('password=secret /private/driver'));
  await expect(studioV2WorkspaceDatabaseAPI.getMetadata('C1', undefined, { workspaceId: 'W', workspaceRevision: 'W1' })).rejects.toThrow('Scoped settings read unavailable.');
});
