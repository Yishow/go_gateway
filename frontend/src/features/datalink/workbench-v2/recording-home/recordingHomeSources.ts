import type { HomeDeviceSummary } from './recordingHomeCollections';
import { studioV2WorkspaceAPI } from '../../../../services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../../services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '../../../../services/studioV2Rules';
import { studioV2MappingsAPI } from '../../../../services/studioV2Mappings';
import { pointAPI } from '../../../../services/datalink';
import { buildGroupCandidates, type GroupCandidates, type GroupCandidateState } from '../state/writeGroup/candidates';
import type { RecordingHomeEvidence, RecordingHomeScope } from './recordingHomeProjection';

export interface RecordingHomeSourceSummary extends GroupCandidates { devices: HomeDeviceSummary[] }
async function readSources() {
  const [workspace, devices, rules, mappings, points] = await Promise.all([
    studioV2WorkspaceAPI.get(), studioV2WorkspaceDevicesAPI.list(), studioV2RulesAPI.list(), studioV2MappingsAPI.list(), pointAPI.list(),
  ]);
  // PointHandler.List currently returns the complete list. Never infer a complete
  // result from a paged endpoint; this projection discards foreign devices/values.
  return { workspace: { id: workspace.id, ordered_device_ids: workspace.ordered_device_ids },
    devices: devices.map(({ id, name, protocol, availability_status }) => ({ id, name, protocol: protocol ?? 'unknown', availability: availability_status ?? 'unknown' })), rules: rules.map(rule => ({ id: rule.id, device_id: rule.device_id,
      workspace_id: rule.workspace_id, revision_id: rule.revision_id, enabled: rule.enabled, skipped_addresses: rule.skipped_addresses })),
    mappings: mappings.map(({ id, workspace_id, point_id, rule_id, device_id, address, tag_id, tag_key, display_name, target_type, enabled }) =>
      ({ id, workspace_id, point_id, rule_id, device_id, address, tag_id, tag_key, display_name, target_type, enabled })),
    points: points.filter(point => workspace.ordered_device_ids.includes(point.device_id)).map(({ id, device_id, address, enabled }) => ({ id, device_id, address, enabled })),
  };
}
type Sources = Awaited<ReturnType<typeof readSources>>;
function fingerprint(value: Sources): string {
  const sort = <T extends { id: string }>(items: T[]) => [...items].sort((a, b) => a.id < b.id ? -1 : a.id > b.id ? 1 : 0);
  return JSON.stringify({ ...value, devices: sort(value.devices), rules: sort(value.rules), mappings: sort(value.mappings), points: sort(value.points) });
}
function projectSources(value: Sources, scope: RecordingHomeScope): RecordingHomeSourceSummary | undefined {
  if (value.workspace.id !== scope.workspaceId || value.workspace.ordered_device_ids.length !== value.devices.length ||
    value.devices.some(device => !value.workspace.ordered_device_ids.includes(device.id))) return;
  for (const items of [value.devices, value.rules, value.mappings, value.points]) if (new Set(items.map(item => item.id)).size !== items.length) return;
  if (value.rules.some(rule => rule.workspace_id !== scope.workspaceId || !value.devices.some(device => device.id === rule.device_id))) return;
  const state: GroupCandidateState = { devices: value.devices.map(device => ({ ...device, persisted: true, save_state: 'saved' })),
    rules: value.rules, points: [], mappings: {} };
  const points: GroupCandidateState['points'][number][] = [];
  const mappings: Record<string, GroupCandidateState['mappings'][string]> = {};
  for (const mapping of value.mappings) {
    const rule = value.rules.find(rule => rule.id === mapping.rule_id);
    const point = value.points.find(point => point.id === mapping.point_id);
    if (mapping.workspace_id !== scope.workspaceId || !rule || !point || point.device_id !== mapping.device_id ||
      rule.device_id !== mapping.device_id || point.address !== mapping.address || mappings[point.id]) return;
    points.push({ ...point, rule_id: rule.id, skipped: rule.skipped_addresses.includes(point.address) });
    mappings[point.id] = { ...mapping, persisted: true, save_state: 'saved', persisted_point_id: point.id };
  }
  return { ...buildGroupCandidates({ ...state, points, mappings }), devices: value.devices };
}
/** Captured-generation readonly observation, not an atomic snapshot or command token.
 * Compare actual source content twice; Share revision cannot prove its freshness. */
export async function readRecordingHomeSources(scope: RecordingHomeScope): Promise<RecordingHomeEvidence<RecordingHomeSourceSummary>> {
  const captured = { ...scope };
  try {
    const before = await readSources();
    const after = await readSources();
    if (fingerprint(before) !== fingerprint(after)) return { state: 'unknown', scope: captured };
    const value = projectSources(after, captured);
    return value ? { state: 'available', scope: captured, value } : { state: 'unknown', scope: captured };
  } catch { return { state: 'error', scope: captured }; }
}
