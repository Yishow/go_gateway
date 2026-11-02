import { studioV2WorkspaceAPI } from '../../../../services/studioV2Workspace';
import { studioV2RulesAPI } from '../../../../services/studioV2Rules';
/** The list supplies labels only. Opening an item still requires its exact scoped
 * locator read and original hydration owner before any mutation is available. */
export async function readHomeDeviceRanges(workspaceId: string, deviceId: string) {
  const before = await studioV2WorkspaceAPI.get();
  if (before.id !== workspaceId || !before.ordered_device_ids.includes(deviceId)) throw new Error('source scope unavailable');
  const rules = await studioV2RulesAPI.list();
  const after = await studioV2WorkspaceAPI.get();
  if (after.id !== workspaceId || JSON.stringify(before.ordered_device_ids) !== JSON.stringify(after.ordered_device_ids)) throw new Error('source scope changed');
  const selected = rules.filter(rule => rule.workspace_id === workspaceId && rule.device_id === deviceId);
  if (new Set(selected.map(rule => rule.id)).size !== selected.length) throw new Error('source identity unavailable');
  return selected;
}
