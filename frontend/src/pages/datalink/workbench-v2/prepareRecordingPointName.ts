import { studioV2WorkspaceAPI } from '../../../services/studioV2Workspace';
import type { WorkbenchV2State, Mapping, Point } from '../../../features/datalink/workbench-v2/state/types';
import type { SourceIdentity } from './rangeSourceSubmission';
import { studioV2MappingsAPI } from '../../../services/studioV2Mappings';
import { readStudioLocatorEntry } from '../../../services/studioV2LocatorEntry';
import { studioV2LocatorHydration } from './studioV2LocatorHydration';
import { mappingValueFromState, mappingValueFromRecord } from '../../../features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { toStudioV2RuleUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { sourceRuleCandidateAPI } from '../../../services/sourceRuleCandidates';
import { buildPersistedMCRuleState, type PersistedMCRuleState } from '../../../features/datalink/workbench-v2/state/mcPersistedPoints';
import { toStudioV2DeviceUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { sameMappingValue, sameHydratedMapping } from './studioV2MappingQueueHelpers';

// Same canonical value comparison as the original device submission observer.
// Config key order is not an edit; array order and every value remain significant.
function comparable(value: unknown): string {
  if (Array.isArray(value)) return `[${value.map(comparable).join(',')}]`;
  if (value && typeof value === 'object') return `{${Object.entries(value).sort(([a], [b]) => a.localeCompare(b)).map(([key, item]) => `${JSON.stringify(key)}:${comparable(item)}`).join(',')}}`;
  return JSON.stringify(value);
}

export interface PointNamePreparation {
  workspaceId: string; workspaceRevision: string; identity: SourceIdentity; mappingId: string; ruleId: string; ruleRevision: string;
  point: Point; mapping: Mapping; mcEvidence?: PersistedMCRuleState;
}
/** Exact one-point read/merge proposal. Never replaces devices, rules, database or settings.
 * Read freshness is not mutation CAS. The original queue must still accept canMutateRule.
 */
export async function prepareRecordingPointName(workspaceId: string, identity: SourceIdentity,
  state: () => WorkbenchV2State, current: () => boolean,
): Promise<{ expectedState: WorkbenchV2State; preparation: PointNamePreparation; merge?: Pick<WorkbenchV2State, 'devices' | 'rules' | 'points' | 'mappings'> } | undefined> {
  const initial = state();
  const all = await studioV2MappingsAPI.list();
  if (new Set(all.map(row => row.id)).size !== all.length) return;
  const rows = all.filter(row => row.workspace_id === workspaceId && row.device_id === identity.device_id && row.point_id === identity.point_id && row.tag_id === identity.tag_id);
  if (!current() || rows.length !== 1 || !rows[0].id) return;
  const row = rows[0];
  const read = await readStudioLocatorEntry({ step: 3, workspace_id: workspaceId, device_id: identity.device_id, rule_id: null, mapping_id: row.id, group_id: null, group_revision: null });
  if (!current() || !read.ok || read.value.workspace.id !== workspaceId || read.value.entry.mappings?.length !== 1 ||
    !read.value.entry.revalidated_mapping_ids.includes(row.id)) return;
  const record = read.value.entry.mappings[0];
  if (record.id !== row.id || record.workspace_id !== workspaceId || record.device_id !== identity.device_id || record.point_id !== identity.point_id || record.tag_id !== identity.tag_id || record.rule_id !== row.rule_id ||
    !sameMappingValue(mappingValueFromRecord(record), mappingValueFromRecord(row))) return;
  const projected = studioV2LocatorHydration(read.value);
  if (!projected || projected.points.length !== 1 || Object.keys(projected.mappings).length !== 1) return;
  const point = projected.points[0], mapping = projected.mappings[point.id];
  const freshDevice = projected.devices.find(value => value.id === identity.device_id), freshRule = projected.rules.find(value => value.id === record.rule_id);
  if (!freshDevice || !freshRule || !freshRule.revision_id) return;
  let mcEvidence: PersistedMCRuleState | undefined;
  if (freshDevice.protocol === 'mc_3e') {
    const rawRule = read.value.entry.source_rules?.find(value => value.id === record.rule_id);
    const workspace = await studioV2WorkspaceAPI.get();
    const revision = workspace.modbus_share?.workspace_revision;
    if (!current() || !rawRule || workspace.id !== workspaceId || !revision?.trim() || !workspace.ordered_device_ids.includes(identity.device_id)) return;
    const snapshot = await sourceRuleCandidateAPI.get(rawRule.id, { workspace_id: workspaceId, expected_workspace_revision: revision, revision_id: freshRule.revision_id });
    const after = await studioV2WorkspaceAPI.get();
    if (!current() || after.id !== workspaceId || after.modbus_share?.workspace_revision !== revision || !after.ordered_device_ids.includes(identity.device_id)) return;
    mcEvidence = buildPersistedMCRuleState(rawRule, { workspace_id: workspaceId, workspace_revision: revision, revision_id: freshRule.revision_id }, snapshot);
    if (mcEvidence.status !== 'compatible' || !mcEvidence.points.some(value => value.point_id === identity.point_id && value.mapping_id === record.id && value.tag_id === identity.tag_id && value.address === record.address)) return;
  }
  const live = state(), rule = live.rules.find(value => value.id === record.rule_id), device = live.devices.find(value => value.id === identity.device_id);
  if (live.rules.filter(value => value.id === record.rule_id).length > 1 || live.devices.filter(value => value.id === identity.device_id).length > 1 ||
    (rule && (rule.device_id !== identity.device_id || rule.persisted !== true || rule.save_state !== 'saved' ||
      rule.revision_id !== freshRule.revision_id || JSON.stringify(toStudioV2RuleUpdateRequest(rule)) !== JSON.stringify(toStudioV2RuleUpdateRequest(freshRule)) || initial.rules.find(value => value.id === rule.id) !== rule)) ||
    (device && (!device.persisted || device.save_state !== 'saved' || comparable(toStudioV2DeviceUpdateRequest(device)) !== comparable(toStudioV2DeviceUpdateRequest(freshDevice)) || initial.devices.find(value => value.id === device.id) !== device))) return;
  if (Object.values(live.mappings).some(value => value.rule_id === record.rule_id && value.mapping_id !== record.id && value.save_state !== 'saved')) return;
  const ownedRule = rule ?? freshRule;
  const candidates = live.points.filter(value => value.id === identity.point_id || live.mappings[value.id]?.persisted_point_id === identity.point_id);
  if (candidates.length > 1 || Object.values(live.mappings).filter(value => value.mapping_id === record.id).length > 1) return;
  const local = candidates[0], existing = local && live.mappings[local.id];
  if (local && (!existing || local.device_id !== identity.device_id || local.rule_id !== ownedRule.id || local.address !== record.address ||
    existing.workspace_id !== workspaceId || existing.device_id !== identity.device_id || existing.rule_id !== ownedRule.id || existing.address !== record.address ||
    existing.mapping_id !== record.id || existing.tag_id !== identity.tag_id || existing.persisted_point_id !== identity.point_id ||
    !existing.persisted || existing.save_state !== 'saved' || !sameMappingValue(mappingValueFromState(existing), mappingValueFromRecord(record)) ||
    initial.points.find(value => value.id === local.id) !== local || !initial.mappings[local.id] || !sameHydratedMapping(initial.mappings[local.id], existing))) return;
  if (!local && (Object.values(live.mappings).some(value => value.mapping_id === record.id) || live.mappings[point.id] || live.points.some(value => value.id === point.id))) return;
  const preparedPoint = local ?? point, preparedMapping = existing ?? mapping;
  return { expectedState: live, preparation: { workspaceId, workspaceRevision: read.value.entry.workspace_revision, identity: { ...identity }, mappingId: record.id, ruleId: ownedRule.id, ruleRevision: ownedRule.revision_id!, point: preparedPoint, mapping: preparedMapping, mcEvidence },
    merge: local && rule && device && !mcEvidence ? undefined : { devices: device ? live.devices : [...live.devices, freshDevice],
      rules: rule ? live.rules.map(value => value === rule && mcEvidence ? { ...value, mc_persisted_points: mcEvidence } : value) : [...live.rules, { ...freshRule, ...(mcEvidence ? { mc_persisted_points: mcEvidence } : {}) }],
      points: local ? live.points : [...live.points, point], mappings: local ? live.mappings : { ...live.mappings, [point.id]: mapping } } };
}
