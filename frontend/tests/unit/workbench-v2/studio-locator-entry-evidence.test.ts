import { describe, expect, it } from 'vitest';
import { resolveStudioLocatorEntrySelection } from '@/features/datalink/workbench-v2/state/studioLocatorEntry';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
import type { StudioV2LocatorEntrySnapshot } from '@/types/studioV2LocatorEntry';
import type { StudioV2Workspace } from '@/types/studioV2Workspace';
const workspace: StudioV2Workspace = { id: 'W', kind: 'single', status: 'ready', ordered_device_ids: ['D'], created_at: '', updated_at: '' };
const locator: EntityLocator = { step: 4, workspace_id: 'W', device_id: null, rule_id: null, mapping_id: null, group_id: 'G1', group_revision: 'R1' };
function entry(): StudioV2LocatorEntrySnapshot {
  return { workspace_revision: '', hydration_scope: { complete_workspace: false, device_ids: [], rule_ids: [], mapping_ids: [], group_ids: ['G1'] },
    rule_heads: {}, group_heads: { G1: 'R2' }, revalidated_mapping_ids: [], write_groups: { workspace_id: 'W', workspace_revision: '', groups: [{
      id: 'G1', workspace_id: 'W', revision: 'R2', applied_revision: 'R1', name: 'Same name', status: 'disabled', members: [],
      destination: { connector_id: 'C', connector_revision: 'C1', database: 'private.db', table_schema: '', table_name: 'meter', storage_strategy: 'managed' },
      row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0 }, write_policy: {}, migration: {}, created_at: '', updated_at: '',
    }] } };
}
describe('native entry projects only identity proof to unchanged canonical selection', () => {
  it('uses current R2 for selection and retains historical R1 only as informational hint', () => {
    const result = resolveStudioLocatorEntrySelection(locator, workspace, entry());
    expect(result).toEqual({ state: 'selected', target: { workspaceId: 'W', step: 4, workspaceRevision: '', groupId: 'G1', groupRevision: 'R2' }, history: { groupId: 'G1', revisionHint: 'R1', relation: 'different' } });
    expect(JSON.stringify(result)).not.toContain('private.db'); expect(result).not.toHaveProperty('deviceId');
  });
  it('preserves a deleted G1 unavailable result rather than selecting same-named G2', () => {
    const value = entry(); value.write_groups!.groups[0].status = 'deleted';
    value.write_groups!.groups.push({ ...value.write_groups!.groups[0], id: 'G2', status: 'ready' }); value.group_heads.G2 = 'R2';
    expect(resolveStudioLocatorEntrySelection(locator, workspace, value)).toEqual({ state: 'unavailable', reason: 'target_unavailable' });
  });
  it.each(['head', 'workspace-revision'])('leaves %s conflict to the canonical fail-closed resolver', (kind) => {
    const value = entry(); if (kind === 'head') value.group_heads.G1 = 'R3'; else value.write_groups!.workspace_revision = 'W2';
    expect(resolveStudioLocatorEntrySelection(locator, workspace, value)).toEqual({ state: 'unavailable', reason: 'revision_conflict' });
  });
  it('does not turn missing revision into explicit empty current proof', () => {
    const value = entry(); value.workspace_revision = undefined as unknown as string;
    expect(resolveStudioLocatorEntrySelection(locator, workspace, value)).toEqual({ state: 'unavailable', reason: 'revision_unavailable' });
  });
});
