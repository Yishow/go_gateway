import type { EntityLocator } from '../../../../types/runtimeRuleEvidence';
import type { StudioV2LocatorEntrySnapshot } from '../../../../types/studioV2LocatorEntry';
import type { StudioV2Workspace } from '../../../../types/studioV2Workspace';
import { resolveStudioSelection, type StudioSelectionEvidence, type StudioSelectionResult } from './studioSelection';

/** Project current entry proof into the canonical resolver, never a parallel selection authority. */
export function resolveStudioLocatorEntrySelection(locator: EntityLocator, workspace: StudioV2Workspace, entry: StudioV2LocatorEntrySnapshot): StudioSelectionResult {
  const evidence: StudioSelectionEvidence = { state: 'current', workspace: { id: workspace.id, ordered_device_ids: [...workspace.ordered_device_ids] },
    workspaceRevision: entry.workspace_revision,
    devices: entry.devices?.map(({ id }) => ({ id })),
    rules: entry.source_rules?.map(({ id, workspace_id, device_id, revision_id }) => ({
      value: { id, workspace_id, device_id, revision_id }, currentRevision: entry.rule_heads[id],
    })),
    mappings: entry.mappings?.filter(({ id }) => entry.revalidated_mapping_ids.includes(id)).map(({ id, workspace_id, device_id, rule_id, point_id, tag_id }) =>
      ({ id, workspace_id, device_id, rule_id, point_id, tag_id })),
    groups: entry.write_groups?.groups.map(({ id, workspace_id, revision, status, basic_managed_device_id, members }) => ({
      value: { id, workspace_id, revision, status, basic_managed_device_id, members: members.map(({ device_id, point_id, tag_id }) => ({ device_id, point_id, tag_id })) },
      currentRevision: entry.group_heads[id], workspaceRevision: entry.write_groups!.workspace_revision,
    })),
  };
  return resolveStudioSelection(workspace.id, { workspaceId: locator.workspace_id, step: locator.step,
    ...(locator.device_id !== null ? { deviceId: locator.device_id } : {}), ...(locator.rule_id !== null ? { ruleId: locator.rule_id } : {}),
    ...(locator.mapping_id !== null ? { mappingId: locator.mapping_id } : {}), ...(locator.group_id !== null ? { groupId: locator.group_id } : {}),
    ...(locator.group_revision !== null ? { groupRevisionHint: locator.group_revision } : {}),
  }, { readCurrent: () => evidence });
}
