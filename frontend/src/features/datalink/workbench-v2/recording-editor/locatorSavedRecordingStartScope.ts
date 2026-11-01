import type { StudioLocatorEntryRead } from '../../../../types/studioV2LocatorEntry';
import { studioV2LocatorHydration } from '../../../../pages/datalink/workbench-v2/studioV2LocatorHydration';
import { savedRecordingStartDevices, savedRecordingStartSourcesSaved } from './savedRecordingStartScope';
/** A detached exact editor receipt, never the outer URL hint or a complete workspace fallback. */
export function locatorSavedRecordingStartScope(read: StudioLocatorEntryRead, workspaceId: string) {
  if (read.selection.state !== 'selected' || !read.entry.editor_hydrated) return undefined;
  const target = read.selection.target, list = read.entry.write_groups, scope = read.entry.hydration_scope;
  const group = list?.groups.length === 1 ? list.groups[0] : undefined;
  if (target.step !== 4 || target.workspaceId !== workspaceId || read.workspace.id !== workspaceId ||
    read.locator.workspace_id !== workspaceId || list?.workspace_id !== workspaceId || !group || group.workspace_id !== workspaceId ||
    target.groupId !== group.id || read.locator.group_id !== group.id || target.groupRevision !== group.revision ||
    read.entry.group_heads[group.id] !== group.revision || scope.complete_workspace !== false ||
    scope.group_ids.length !== 1 || scope.group_ids[0] !== group.id ||
    !target.workspaceRevision || target.workspaceRevision !== read.entry.workspace_revision || list.workspace_revision !== target.workspaceRevision) return undefined;
  const devices = savedRecordingStartDevices(group), state = studioV2LocatorHydration(read);
  if (!devices.length || !state || !savedRecordingStartSourcesSaved(group, state) || devices.some(id => !scope.device_ids.includes(id)) ||
    group.members.some(member => !(read.entry.mappings ?? []).some(mapping => mapping.device_id === member.device_id &&
      mapping.point_id === member.point_id && mapping.tag_id === member.tag_id && mapping.workspace_id === workspaceId &&
      scope.mapping_ids.includes(mapping.id) && scope.rule_ids.includes(mapping.rule_id) &&
      read.entry.source_rules?.some(rule => rule.id === mapping.rule_id && rule.device_id === member.device_id &&
        rule.workspace_id === workspaceId && read.entry.rule_heads[rule.id] === rule.revision_id)))) return undefined;
  return { group, state, workspaceId, workspaceRevision: list.workspace_revision };
}
