import { studioV2WorkspaceAPI } from '../../../services/studioV2Workspace';
import { studioV2RulesAPI } from '../../../services/studioV2Rules';
import { studioV2WorkspaceDevicesAPI, type StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
import { sourceRuleCandidateAPI } from '../../../services/sourceRuleCandidates';
import { studioV2MappingsAPI } from '../../../services/studioV2Mappings';
import { pointAPI, tagAPI } from '../../../services/datalink';
import type { SourceRuleRecord } from '../../../types/datalink';
import type { SourceRuleCandidateSnapshotView } from '../../../types/sourceRuleCandidates';
import { hydrateStudioV2Rule, toStudioV2RuleCreateRequest } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { buildPersistedMCRuleState, type PersistedMCRuleState } from '../../../features/datalink/workbench-v2/state/mcPersistedPoints';
import { derivePoints } from '../../../features/datalink/workbench-v2/state/sourceRule';
import { buildGroupCandidates, type GroupCandidates, type GroupCandidateState } from '../../../features/datalink/workbench-v2/state/writeGroup/candidates';
import type { RangeMappingReceipt, RangeRuleReceipt, RangeSourceRequest } from './rangeSourceSubmission';

export interface RangeRuleProof {
  request: RangeSourceRequest;
  receipt: RangeRuleReceipt;
  record: SourceRuleRecord;
  snapshot: SourceRuleCandidateSnapshotView;
  device: StudioV2WorkspaceDeviceRecord;
  mcEvidence?: PersistedMCRuleState;
}
const same = (a: unknown, b: unknown): boolean => JSON.stringify(a) === JSON.stringify(b);
const canonical = (rule: SourceRuleRecord) => toStudioV2RuleCreateRequest(hydrateStudioV2Rule(rule));
const matches = (expected: object, actual: object) => Object.entries(expected).every(([key, value]) => same(value, (actual as Record<string, unknown>)[key]));
function ruleMatches(receipt: RangeRuleReceipt, request: RangeSourceRequest) {
  const expected = toStudioV2RuleCreateRequest({ ...request.buffer, id: receipt.submitted.id ?? '', device_id: request.deviceId });
  return Boolean(receipt.submitted.id?.trim() && matches(expected, receipt.submitted) &&
    matches(expected, canonical(receipt.acknowledged)) && receipt.acknowledged.workspace_id === request.workspaceId &&
    receipt.acknowledged.revision_id?.trim() && receipt.acknowledged.runtime_apply_status !== 'apply_failed');
}
/** Read-only evidence after a known acknowledgement. Never locates a lost create by name. */
export async function readRuleProof(receipt: RangeRuleReceipt, request: RangeSourceRequest): Promise<RangeRuleProof | undefined> {
  try {
    if (!ruleMatches(receipt, request)) return;
    const workspace = await studioV2WorkspaceAPI.get();
    const revision = workspace.modbus_share?.workspace_revision;
    if (workspace.id !== request.workspaceId || !revision?.trim() || !workspace.ordered_device_ids.includes(request.deviceId)) return;
    const [rules, devices] = await Promise.all([studioV2RulesAPI.list(), studioV2WorkspaceDevicesAPI.list()]);
    const matching = rules.filter(rule => rule.id === receipt.acknowledged.id);
    const matchingDevices = devices.filter(device => device.id === request.deviceId);
    if (matching.length !== 1 || matchingDevices.length !== 1) return;
    const record = matching[0], device = matchingDevices[0];
    if (record.workspace_id !== request.workspaceId || record.revision_id !== receipt.acknowledged.revision_id ||
      !matches(canonical(receipt.acknowledged), canonical(record)) || device.runtime_apply_status === 'apply_failed') return;
    const snapshot = await sourceRuleCandidateAPI.get(record.id, { workspace_id: workspace.id, expected_workspace_revision: revision, revision_id: record.revision_id! });
    if (snapshot.source_rule_id !== record.id || snapshot.workspace_id !== workspace.id || snapshot.workspace_revision !== revision ||
      snapshot.revision_id !== record.revision_id || snapshot.tags.status !== 'ready') return;
    const candidates = snapshot.tags.candidates;
    if (!candidates.length || new Set(candidates.map(c => c.point_id)).size !== candidates.length ||
      new Set(candidates.map(c => c.address)).size !== candidates.length || new Set(candidates.map(c => c.id)).size !== candidates.length ||
      candidates.some(c => !c.id?.trim() || !c.point_id?.trim() || c.identity.source_rule_id !== record.id)) return;
    const mcEvidence = device.protocol === 'mc_3e' ? buildPersistedMCRuleState(record, { workspace_id: workspace.id, workspace_revision: revision, revision_id: record.revision_id! }, snapshot) : undefined;
    if (mcEvidence && mcEvidence.status !== 'compatible') return;
    if (!mcEvidence) {
      const planned = derivePoints(hydrateStudioV2Rule(record), device.id, new Set(record.skipped_addresses), device.protocol).filter(point => !point.skipped);
      if (planned.length !== candidates.length || planned.some(point => !candidates.some(c => c.address === point.address))) return;
    }
    const after = await studioV2WorkspaceAPI.get();
    if (after.id !== workspace.id || after.modbus_share?.workspace_revision !== revision || !same(after.ordered_device_ids, workspace.ordered_device_ids)) return;
    return { request, receipt, record, snapshot, device, ...(mcEvidence ? { mcEvidence } : {}) };
  } catch { return; }
}
/** Complete canonical identities, not synthetic planned point IDs or an operation ledger. */
export async function readRangeCompletion(proof: RangeRuleProof, receipts: readonly RangeMappingReceipt[]): Promise<GroupCandidates | undefined> {
  try {
    const fresh = await readRuleProof(proof.receipt, proof.request);
    if (!fresh || receipts.length !== fresh.snapshot.tags.candidates.length || new Set(receipts.map(r => r.acknowledged.id)).size !== receipts.length) return;
    const records = await studioV2MappingsAPI.list();
    const points: GroupCandidateState['points'][number][] = [];
    const mappings: Record<string, GroupCandidateState['mappings'][string]> = {};
    for (const { submitted, acknowledged: ack } of receipts) {
      if (!ack.id?.trim() || !ack.point_id?.trim() || !ack.tag_id?.trim() || ack.runtime_apply_status === 'apply_failed' || ack.cleanup_status === 'failed' ||
        ack.rule_id !== fresh.record.id || ack.workspace_id !== fresh.request.workspaceId || ack.device_id !== fresh.device.id ||
        !matches(submitted, ack)) return;
      const rows = records.filter(row => row.id === ack.id);
      if (rows.length !== 1 || !matches(submitted, rows[0]) ||
        !matches({ workspace_id: ack.workspace_id, device_id: ack.device_id, point_id: ack.point_id, tag_id: ack.tag_id }, rows[0])) return;
      // Revisioned candidates are proposals and may predate live mapping links.
      // The exact acknowledged mapping ID above is resolved by workspace GET
      // through current source_rule_links; never locate by display name.
      const candidate = fresh.snapshot.tags.candidates.find(c => c.point_id === ack.point_id && c.address === submitted.address);
      if (!candidate || (candidate.mapping_id && candidate.mapping_id !== ack.id) || (candidate.tag_id && candidate.tag_id !== ack.tag_id) ||
        candidate.status === 'out_of_sync' || candidate.status === 'error' || candidate.blocking_reason || mappings[ack.point_id]) return;
      // Candidate metadata is the rule proposal, not the actual explicitly saved
      // workspace mapping. Verify actual values against its receipt and Tag below.
      const [point, tag] = await Promise.all([pointAPI.get(ack.point_id), tagAPI.get(ack.tag_id)]);
      if (point.id !== ack.point_id || point.device_id !== fresh.device.id || point.address !== ack.address || point.data_type !== fresh.record.data_type || !point.enabled ||
        tag.id !== ack.tag_id || tag.key !== ack.tag_key || tag.data_type !== ack.target_type ||
        tag.display_name !== ack.display_name || (tag.unit === undefined ? '' : tag.unit) !== ack.unit || tag.status === 'retired') return;
      points.push({ ...point, rule_id: fresh.record.id, skipped: false });
      mappings[point.id] = { ...rows[0], persisted: true, save_state: 'saved', persisted_point_id: point.id };
    }
    const after = await studioV2WorkspaceAPI.get();
    if (after.id !== fresh.request.workspaceId || after.modbus_share?.workspace_revision !== fresh.snapshot.workspace_revision) return;
    const result = buildGroupCandidates({ devices: [{ ...fresh.device, persisted: true, save_state: 'saved' }], rules: [fresh.record], points, mappings });
    return result.candidates.length === receipts.length ? result : undefined;
  } catch { return; }
}
