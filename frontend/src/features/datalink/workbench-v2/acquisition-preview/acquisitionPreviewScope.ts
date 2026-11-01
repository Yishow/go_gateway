import type { AcquisitionPreviewScope } from '../../../../types/acquisitionPlanPreview';
import type { StudioLocatorEntryRead } from '../../../../types/studioV2LocatorEntry';
import type { DbConnector } from '../state/types-settings';
import type { ModbusShareStatus } from '../../../../types/modbusShare';
export type AcquisitionWorkspaceScope = Pick<AcquisitionPreviewScope, 'workspace_id' | 'expected_workspace_revision'>;
export interface AcquisitionScopeEvidence { workspaceId?: string; databaseConnector?: Pick<DbConnector, 'workspace_id' | 'setup_revision' | 'persisted' | 'save_state'>; locatorEntry?: StudioLocatorEntryRead; locatorSuspended?: boolean; shareStatus?: Pick<ModbusShareStatus, 'workspace_id' | 'workspace_revision'> | null }

/** Existing canonical read evidence only; no lookup, rebase, or mutation authority. */
export function acquisitionPreviewWorkspaceScope({ workspaceId, locatorEntry, locatorSuspended, databaseConnector }: AcquisitionScopeEvidence): AcquisitionWorkspaceScope | undefined {
  if (!workspaceId || locatorSuspended) return undefined;
  if (locatorEntry) {
    const selection = locatorEntry.selection;
    if (locatorEntry.workspace.id !== workspaceId || locatorEntry.locator.workspace_id !== workspaceId || selection.state !== 'selected' ||
      selection.target.workspaceId !== workspaceId || selection.target.workspaceRevision !== locatorEntry.entry.workspace_revision ||
      typeof locatorEntry.entry.workspace_revision !== 'string') return undefined;
    return { workspace_id: workspaceId, expected_workspace_revision: locatorEntry.entry.workspace_revision };
  }
  // This last hydrated/acknowledged setup head is not the unrelated source/Share revision.
  // Editing connection fields does not invalidate this evidence; the API rejects stale heads.
  return databaseConnector?.workspace_id === workspaceId && typeof databaseConnector.setup_revision === 'string'
    ? { workspace_id: workspaceId, expected_workspace_revision: databaseConnector.setup_revision } : undefined;
}
