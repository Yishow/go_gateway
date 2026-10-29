import type { EvidenceScope, RecordingGroupFacts, RuleEvidenceView } from '@/types/runtimeRuleEvidence';

// Test-only server response. Never imported by production code.
export const evidenceScope: EvidenceScope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'D' };
export function ruleEvidenceFixture(scope = evidenceScope, effectiveRevision = 'A'): RuleEvidenceView {
  const facts = (revision: string, table: string, interval: number): RecordingGroupFacts => ({
    binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: revision },
    recording_profile: { kind: 'instantaneous', interval, unit: 'second' }, allowed_lateness_seconds: 2,
    ends_at: null, effective_at: revision === 'C' ? null : '2026-10-27T09:00:00Z', enabled: true,
    destination: { connector_id: 'C1', connector_revision: '7', alias: null, kind: 'sqlite', schema: null, table },
    members: [{ device_id: scope.device_id, point_id: 'P', tag_id: 'T', column: 'value', source_value_type: 'float64', circuit_id: null }],
    layout: 'managed_wide_profile_v1', storage_mode: 'plain', interval_seconds: interval,
    reduction: 'mean', write_mode: 'append', incomplete_policy: 'partial',
  });
  const observedAt = effectiveRevision === 'B' ? '2026-10-27T10:00:01Z' : '2026-10-27T10:00:00Z';
  const effective = effectiveRevision === 'B' ? facts('B', 'meter_new', 30) : facts('A', 'meter_old', 15);
  if (effectiveRevision === 'B') effective.effective_at = observedAt;
  return {
    contract_version: 1, ...scope, observed_at: observedAt, freshness: 'fresh', source: null, locator: null,
    groups: [{ group_id: 'G', saved_revision: 'C', applied_revision: 'B', effective_revision: effectiveRevision,
      as_of: observedAt, applied_effective_at: '2026-10-27T10:00:01Z', state: 'effective', proof: 'applied-snapshot-and-layout',
      saved: facts('C', 'meter_draft', 60), effective, delivery: null }],
  };
}
export function evidenceContext(evidence: unknown = ruleEvidenceFixture(), workspaceID = 'W') {
  return { workspace_id: workspaceID, devices: [], default_device_id: null, setup: null, rule_evidence: evidence };
}
