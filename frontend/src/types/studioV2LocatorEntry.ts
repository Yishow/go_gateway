import type { Device, SourceRuleRecord, StudioV2WorkspaceMappingRecord, StudioV2WorkspaceDatabaseConfigRecord, StudioV2WorkspaceDatabaseTargetRecord } from './datalink';
import type { EntityLocator } from './runtimeRuleEvidence';
import type { StudioV2Workspace } from './studioV2Workspace';
import type { WriteGroupListResponse } from './studioV2WriteGroup';
import type { StudioSelectionResult } from '../features/datalink/workbench-v2/state/studioSelection';

/** Detached, partial native setup hydration. Omission never means a complete empty collection. */
export interface StudioV2LocatorEntrySnapshot {
  workspace_revision: string;
  editor_hydrated?: true;
  database_config?: StudioV2WorkspaceDatabaseConfigRecord | null;
  database_targets?: StudioV2WorkspaceDatabaseTargetRecord[];
  hydration_scope: { complete_workspace: false; device_ids: string[]; rule_ids: string[]; mapping_ids: string[]; group_ids: string[] };
  devices?: Device[];
  source_rules?: SourceRuleRecord[];
  mappings?: StudioV2WorkspaceMappingRecord[];
  write_groups?: WriteGroupListResponse;
  rule_heads: Record<string, string>;
  group_heads: Record<string, string>;
  revalidated_mapping_ids: string[];
}
export type StudioLocatorEntryRead = {
  locator: EntityLocator;
  workspace: StudioV2Workspace;
  entry: StudioV2LocatorEntrySnapshot;
  selection: StudioSelectionResult;
};
export type StudioLocatorEntryResult = { ok: true; value: StudioLocatorEntryRead } |
  { ok: false; reason: 'unavailable' | 'stale-reference' | 'foreign-scope' };

/** Observed native setup read scope, not mutation authorization or a historical revision hint. */
export interface StudioV2LocatorReadScope { readonly workspaceId: string; readonly workspaceRevision: string }
