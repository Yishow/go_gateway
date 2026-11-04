import { studioV2WorkspaceAPI } from '../../../../services/studioV2Workspace';
import { studioV2RulesAPI } from '../../../../services/studioV2Rules';
/** The list supplies labels only. Opening an item still requires its exact scoped
 * locator read and original hydration owner before any mutation is available.
 * `workspaceRevision` binds derived evidence (MC candidates) to this same read. */
export async function readHomeDeviceRanges(workspaceId: string, deviceId: string) {
  const before = await studioV2WorkspaceAPI.get();
  if (before.id !== workspaceId || !before.ordered_device_ids.includes(deviceId)) throw new Error('source scope unavailable');
  const rules = await studioV2RulesAPI.list();
  const after = await studioV2WorkspaceAPI.get();
  // Same scope rule as the page: a revision counts only when it belongs to this workspace.
  const revisionOf = (read: typeof before) => read.modbus_share?.workspace_id === workspaceId ? read.modbus_share.workspace_revision || undefined : undefined;
  if (after.id !== workspaceId || JSON.stringify(before.ordered_device_ids) !== JSON.stringify(after.ordered_device_ids) ||
    (before.modbus_share?.workspace_revision ?? null) !== (after.modbus_share?.workspace_revision ?? null)) throw new Error('source scope changed');
  const selected = rules.filter(rule => rule.workspace_id === workspaceId && rule.device_id === deviceId);
  if (new Set(selected.map(rule => rule.id)).size !== selected.length) throw new Error('source identity unavailable');
  const revision = revisionOf(before);
  return { rules: selected, workspaceRevision: revision !== undefined && revision === revisionOf(after) ? revision : undefined };
}
