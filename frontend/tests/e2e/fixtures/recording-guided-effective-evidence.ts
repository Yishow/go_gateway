import type { EvidenceScope, RecordingGroupFacts, RuleEvidenceView } from '../../../src/types/runtimeRuleEvidence';
import { parseRuntimeRuleEvidence } from '../../../src/services/runtimeRuleEvidenceParser';
import { AUTHORING_TIME } from './recording-authoring-state';
import type { GuidedMockState } from './recording-guided-state';

/** TEST ONLY: explicit simulated snapshots, not actual SQL/runtime evidence or copied editable facts. */
export function guidedRuleEvidence(state: GuidedMockState, deviceId: string): RuleEvidenceView {
  const scope: EvidenceScope = { workspace_id: 'ws-1', workspace_revision: state.workspaceRevision, device_id: deviceId };
  const unavailable: RuleEvidenceView = { ...scope, contract_version: 1, observed_at: null, freshness: 'unavailable', source: null, groups: [], locator: null };
  if (state.scenario !== 'saved-different' || state.workspaceRevision !== 'W2' || deviceId !== 'dev-1') return unavailable;
  const group = state.groups.find(item => item.id === 'G1');
  if (!group || group.revision !== 'R2' || group.applied_revision !== 'R1' || group.recording_profile?.interval !== 30 ||
    group.row_policy.interval_seconds !== 30) return unavailable;
  // Both named snapshots use deliberately fixed safe fixture fields. Never spread a saved WriteGroup here.
  const common: Pick<RecordingGroupFacts, 'enabled' | 'ends_at' | 'allowed_lateness_seconds' | 'destination' | 'members' |
    'layout' | 'storage_mode' | 'reduction' | 'write_mode' | 'incomplete_policy'> = {
    enabled: true, ends_at: null, allowed_lateness_seconds: 0,
    destination: { connector_id: 'conn-1', connector_revision: 'crev-1', alias: 'Mock recording target', kind: 'sqlite', schema: 'main', table: 'mock_recording_G1' },
    members: [{ device_id: 'dev-1', circuit_id: 'cir-saved', point_id: 'pt-0', tag_id: 'tag-0', column: 'value_0', source_value_type: 'int16' },
      { device_id: 'dev-1', circuit_id: 'cir-saved', point_id: 'pt-1', tag_id: 'tag-1', column: 'value_1', source_value_type: 'float64' }],
    layout: 'managed_long_v1', storage_mode: 'plain', reduction: 'mean', write_mode: 'append', incomplete_policy: 'partial',
  };
  const saved: RecordingGroupFacts = { ...common, binding: { workspace_id: 'ws-1', workspace_revision: 'W2', entity_revision: 'R2' },
    effective_at: null, recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second' }, interval_seconds: 30 };
  const effective: RecordingGroupFacts = { ...common, binding: { workspace_id: 'ws-1', workspace_revision: 'W2', entity_revision: 'R1' },
    effective_at: '2029-12-31T23:59:00.000Z', recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, interval_seconds: 15 };
  const response: RuleEvidenceView = { ...scope, contract_version: 1, observed_at: AUTHORING_TIME, freshness: 'fresh', source: null, locator: null,
    groups: [{ group_id: 'G1', saved_revision: 'R2', applied_revision: 'R1', effective_revision: 'R1', as_of: AUTHORING_TIME,
      applied_effective_at: '2029-12-31T23:59:00.000Z', state: 'effective', proof: 'applied-snapshot-and-layout', saved, effective, delivery: null }] };
  const parsed = parseRuntimeRuleEvidence(response, scope);
  if (!parsed.ok) throw new Error('Malformed simulated guided effective-evidence fixture');
  return parsed.value;
}
