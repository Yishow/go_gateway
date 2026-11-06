import { createHash } from 'node:crypto';
import type { SourceRuleRecord } from '../../../src/types/datalink';
import type { WriteGroupLifecycle } from '../../../src/types/studioV2WriteGroupLifecycle';
import { savedGroup } from '../../fixtures/writeGroupState';
import { nativeDevice, nativeMapping, nativeRule, locatorRead } from '../../unit/workbench-v2/studio-locator-fixtures';
import { quietDelivery } from '../../fixtures/cutoverFixtures';

/** Synthetic, offline fixture only. No device or database connection is made. */
export function snapshotWorkspace(id: 'A' | 'B', count = 60) {
  const devices = ['D', 'E'].map(prefix => ({ ...nativeDevice(`${prefix}-${id}`), name: `Snapshot ${id} ${prefix}`,
    connection_config: '{"host":"127.0.0.1","port":15020,"slave_id":1,"timeout":100}' }));
  const rule = { ...nativeRule(`R-${id}`, devices[0].id), workspace_id: id, enabled: false, share_enabled: false };
  const mapping = { ...nativeMapping(`M-${id}`), workspace_id: id, device_id: devices[0].id, rule_id: rule.id,
    point_id: `P-${id}`, tag_id: `T-${id}`, address: '40002' };
  const makeGroup = (index: number) => savedGroup({ id: `${id}-${String(index + 1).padStart(4, '0')}`, workspace_id: id,
    revision: `${id}-G${index + 1}-r1`, name: `Line ${id} ${String(index + 1).padStart(4, '0')}`,
    members: [{ device_id: devices[0].id, point_id: mapping.point_id, tag_id: mapping.tag_id,
      source_revision: 'source-1', mapping_revision: 'mapping-1', target_column: 'temperature', required: true }] });
  const groups = Array.from({ length: count }, (_, index) => makeGroup(index));
  // Four independent controls prove that each non-default filter actually excludes a row.
  groups.push({ ...makeGroup(count), name: `Other ${id}` }, { ...makeGroup(count + 1), status: 'draft' },
    { ...makeGroup(count + 2), members: [{ ...groups[0].members[0], device_id: devices[1].id }] },
    { ...makeGroup(count + 3), destination: { ...groups[0].destination, connector_id: 'conn-2', database: '/other.db' } });
  const workspace = { ...locatorRead().workspace, id, ordered_device_ids: devices.map(device => device.id) };
  return { workspace, revision: `${id}1`, sourceVersion: 1, devices, rules: [rule] as SourceRuleRecord[], mappings: [mapping],
    points: [{ id: mapping.point_id, device_id: devices[0].id, address: mapping.address, enabled: false }], groups };
}
export type SnapshotWorkspace = ReturnType<typeof snapshotWorkspace>;
export function compactSnapshot(state: SnapshotWorkspace) {
  const devices = state.devices.map(({ id, name, protocol }) => ({ id, name, protocol, availability: 'unknown' }));
  // No source is enabled in this fixture. The complete empty candidate set is intentional.
  const excluded: { label: string; reason: string }[] = [];
  return { workspace_id: state.workspace.id, state: 'available', source_digest: createHash('sha256')
    .update(JSON.stringify({ devices, candidates: [], excluded })).digest('hex'), devices, candidates: [], excluded };
}
export function snapshotLifecycle(state: SnapshotWorkspace, group: SnapshotWorkspace['groups'][number]): WriteGroupLifecycle {
  return { workspace_id: state.workspace.id, workspace_revision: state.revision, group_id: group.id, group_revision: group.revision,
    applied_revision: group.applied_revision, status: group.status, connector_id: group.destination.connector_id,
    connector_revision: group.destination.connector_revision,
    actions: (['edit', 'stop', 'remove', 'recreate', 'transition'] as const).map(action => ({ action, allowed: action === 'edit',
      reason_code: action === 'edit' ? '' : 'source_draft_not_active', next_action: action === 'edit' ? '' : 'review_lifecycle',
      requires_review: action === 'recreate' || action === 'transition' })),
    history: { canonical: false, current_group_id: group.id, generations: [{ group_id: group.id, group_revision: group.revision,
      applied_revision: group.applied_revision, status: group.status, current: true, reserved: false, revisions: [],
      saved_destination: { connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision,
        destination_scope: 'a'.repeat(64) } }] } };
}
export function snapshotDelivery(group: SnapshotWorkspace['groups'][number]) {
  return { ...quietDelivery, group_id: group.id, intake: { state: 'not_running' } };
}
export function snapshotLocator(state: SnapshotWorkspace, params: URLSearchParams) {
  const step = Number(params.get('step')) as 1 | 2 | 3 | 4;
  const group = state.groups.find(group => group.id === params.get('group_id'));
  const device = state.devices.find(device => device.id === params.get('device_id'));
  const rule = state.rules.find(rule => rule.id === params.get('rule_id'));
  const mapping = state.mappings.find(mapping => mapping.id === params.get('mapping_id'));
  if (![1, 2, 3, 4].includes(step) || (step === 4 ? !group : !device) || step === 2 && !rule || step === 3 && !mapping) return undefined;
  const hydrated = params.get('editor') === '1';
  const devices = hydrated ? [state.devices[0]] : device ? [device] : [];
  const rules = hydrated ? state.rules : step === 2 && rule ? [rule] : step === 3 ? state.rules : [];
  const mappings = hydrated ? state.mappings : mapping ? [mapping] : [];
  return { ...state.workspace, locator_entry: { workspace_revision: state.revision,
    hydration_scope: { complete_workspace: false, device_ids: devices.map(row => row.id), rule_ids: rules.map(row => row.id),
      mapping_ids: mappings.map(row => row.id), group_ids: group ? [group.id] : [] },
    ...(devices.length || hydrated ? { devices } : {}), ...(rules.length || hydrated ? { source_rules: rules } : {}),
    ...(mappings.length || hydrated ? { mappings } : {}),
    ...(group ? { write_groups: { workspace_id: state.workspace.id, workspace_revision: state.revision, groups: [group] } } : {}),
    rule_heads: Object.fromEntries(rules.map(row => [row.id, row.revision_id])), group_heads: group ? { [group.id]: group.revision } : {},
    revalidated_mapping_ids: mappings.map(row => row.id), ...(hydrated ? { editor_hydrated: true, database_config: null, database_targets: [] } : {}) } };
}
