import { afterEach, describe, expect, it, vi } from 'vitest';
import { readStudioLocatorEditorEntry, readStudioLocatorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
import response from '../../fixtures/studioLocatorInactiveDraftResponse.json';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));

// Exact HTTP 200 body captured through the native Linux browser after disabling
// the only source, then following the home recording card (request 257).
const locator: EntityLocator = {
  step: 4, workspace_id: 'c69631a9-35da-4f16-bf9a-b2c2e094e19c', device_id: null,
  rule_id: null, mapping_id: null, group_id: '79daf50d-8c46-4bde-8e58-da39feddc83d',
  group_revision: '21c34198-5b7c-441f-84a5-defaec95809b',
};
const expected = {
  workspaceRevision: 'e2ef58b3-45c7-4c0a-87d4-e2c52d9cbd47',
  groupRevision: '21c34198-5b7c-441f-84a5-defaec95809b',
};
function respond(value: unknown = structuredClone(response)) {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: value });
}
afterEach(() => vi.resetAllMocks());

describe('locator readback of canonical calculation capabilities', () => {
  it('accepts the captured existing draft and preserves capabilities without granting readiness', async () => {
    respond();
    const result = await readStudioLocatorEntry(locator);
    expect(result).toMatchObject({ ok: true, value: {
      selection: { state: 'selected', target: { step: 4, workspaceId: locator.workspace_id,
        workspaceRevision: expected.workspaceRevision, groupId: locator.group_id, groupRevision: expected.groupRevision } },
      entry: { write_groups: { groups: [{ name: 'Inactive Draft QA', status: 'draft' }] } },
    } });
    if (!result.ok) throw new Error('captured locator readback was rejected');
    expect(result.value.entry.write_groups?.calculation_capabilities).toEqual(response.data.locator_entry.write_groups.calculation_capabilities);
    expect(result.value.entry.write_groups?.calculation_capabilities?.editor_supported).toBe(false);
    expect(result.value.entry).not.toHaveProperty('editor_hydrated');
    expect(result.value.workspace).not.toHaveProperty('readiness_summary');
  });
});

function editorResponse() {
  const value = structuredClone(response);
  const entry = value.data.locator_entry;
  const group = entry.write_groups.groups[0], member = group.members[0];
  const at = '2026-10-31T08:24:27.303948306Z';
  return { ...value, data: { ...value.data, locator_entry: { ...entry,
    editor_hydrated: true,
    hydration_scope: { ...entry.hydration_scope, device_ids: [member.device_id], rule_ids: ['rule-1'], mapping_ids: ['mapping-1'] },
    devices: [{ id: member.device_id, name: 'Disabled source', description: '', protocol: 'modbus_tcp', status: 'disabled',
      connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: at, updated_at: at }],
    source_rules: [{ id: 'rule-1', workspace_id: group.workspace_id, device_id: member.device_id, revision_id: 'rule-revision-1',
      start_address: '40001', count: 1, data_type: 'int16', naming_prefix: '', enabled: true, locked: false,
      origin: 'manual', skipped_addresses: [], created_at: at, updated_at: at }],
    mappings: [{ id: 'mapping-1', workspace_id: group.workspace_id, device_id: member.device_id, rule_id: 'rule-1',
      point_id: member.point_id, tag_id: member.tag_id, address: '40001', tag_key: 'known', display_name: 'Known', unit: '',
      target_type: 'int16', scale: 1, offset: 0, enabled: true, created_at: at, updated_at: at }],
    rule_heads: { 'rule-1': 'rule-revision-1' }, revalidated_mapping_ids: ['mapping-1'],
    database_config: { id: group.destination.connector_id, identity_revision: group.destination.connector_revision,
      workspace_id: group.workspace_id, setup_revision: entry.workspace_revision, kind: 'sqlite', name: 'Target', host: '', port: 0,
      database: 'target.sqlite', username: '', schema: 'main', table: 'sensor_readings', write_mode: 'insert',
      write_interval_seconds: 60, timestamp_column: 'time', status: 'ready', created_at: at, updated_at: at },
    database_targets: [],
  } } };
}

const read = (editor: boolean) => editor ? readStudioLocatorEditorEntry(locator, expected) : readStudioLocatorEntry(locator);

it('hydrates the existing draft with its only source disabled, without promoting it to active or ready', async () => {
  respond(editorResponse());
  const result = await readStudioLocatorEditorEntry(locator, expected);
  expect(result).toMatchObject({ ok: true, value: { entry: { editor_hydrated: true,
    devices: [{ status: 'disabled' }], write_groups: { groups: [{ status: 'draft', members: [{ target_column: 'known' }] }] },
  } } });
  if (!result.ok) throw new Error('disabled source editor readback was rejected');
  expect(result.value.entry.write_groups?.calculation_capabilities).toEqual(response.data.locator_entry.write_groups.calculation_capabilities);
  expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace', expect.objectContaining({ params: expect.objectContaining({
    editor: '1', expected_workspace_revision: expected.workspaceRevision, expected_group_revision: expected.groupRevision,
  }) }));
});

it.each([false, true])('keeps capabilities optional for an older backend; editor=%s', async editor => {
  const value = editor ? editorResponse() : structuredClone(response);
  Reflect.deleteProperty(value.data.locator_entry.write_groups, 'calculation_capabilities');
  respond(value);
  const result = await read(editor);
  expect(result.ok).toBe(true);
  if (!result.ok) throw new Error('older backend readback was rejected');
  expect(result.value.entry.write_groups).not.toHaveProperty('calculation_capabilities');
});

const capabilities = response.data.locator_entry.write_groups.calculation_capabilities;
const invalidCapabilities: [string, unknown][] = [
  ['null', null], ['empty', {}], ['unknown version', { ...capabilities, version: 2 }],
  ['unknown result shape', { ...capabilities, result_shape: 'member-calculation-v2' }],
  ['unknown key', { ...capabilities, ready: true }], ['invalid boolean', { ...capabilities, editor_supported: 'true' }],
  ['unknown method', { ...capabilities, methods: ['latest', 'other'] }],
  ['missing cadence', { ...capabilities, cadence: undefined }],
  ['unchecked bounds', { ...capabilities, cadence: { ...capabilities.cadence, window_bounds_checked: false } }],
  ['unknown restriction layout', { ...capabilities, unsupported_combinations: [{
    method: 'delta', phase: 'draft', source_transform: 'exact_counter_scale', reason: 'blocked', layout: 'unknown',
  }] }],
];

describe.each([false, true])('strict canonical capability validation; editor=%s', editor => {
  it.each(invalidCapabilities)('fails closed on %s instead of discarding present capability data', async (_name, invalid) => {
    const value = editor ? editorResponse() : structuredClone(response);
    Object.assign(value.data.locator_entry.write_groups, { calculation_capabilities: invalid });
    respond(value);
    expect(await read(editor)).toEqual({ ok: false, reason: 'unavailable' });
  });
});

it.each(['unknown-list-key', 'foreign-list', 'foreign-group', 'workspace-revision', 'group-revision',
  'group-id', 'scope', 'duplicate-group', 'missing-connector-revision', 'malformed-member'])('preserves %s guard with valid capabilities', async kind => {
  const value = structuredClone(response), entry = value.data.locator_entry, list = entry.write_groups, group = list.groups[0];
  if (kind === 'unknown-list-key') Object.assign(list, { authority: true });
  if (kind === 'foreign-list') list.workspace_id = 'foreign';
  if (kind === 'foreign-group') group.workspace_id = 'foreign';
  if (kind === 'workspace-revision') list.workspace_revision = 'other-revision';
  if (kind === 'group-revision') group.revision = 'other-revision';
  if (kind === 'group-id') group.id = 'other-group';
  if (kind === 'scope') entry.hydration_scope.group_ids = ['other-group'];
  if (kind === 'duplicate-group') list.groups.push(structuredClone(group));
  if (kind === 'missing-connector-revision') Reflect.deleteProperty(group.destination, 'connector_revision');
  if (kind === 'malformed-member') Object.assign(group.members[0], { required: 'true' });
  respond(value);
  expect((await readStudioLocatorEntry(locator)).ok).toBe(false);
});
