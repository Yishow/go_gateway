/** TEST ONLY: mock observations, never production/device/database evidence. */
export type RuntimeRecordingScenario = 'normal' | 'many' | 'draft' | 'stale' | 'backlog' | 'uncertain' | 'shared' | 'missing-timing' | 'hostile' | 'deleted' | 'foreign' | 'no-proof' | 'quota-group' | 'quota-global';
export type RuntimeRecordingAction = 'context' | 'status' | 'points' | 'rules' | 'mappings' | 'detail' | 'evidence' | 'locator' | 'reports' | 'native-delivery' | 'native-templates' | 'unknown';
export interface RuntimeRecordingRequest { method: string; url: string; body?: unknown }
export interface RuntimeRecordingReply { status: number; body: unknown }
export interface RuntimeRecordingCall { method: string; path: string; action: RuntimeRecordingAction; status: number; mutation: boolean; params: Record<string, string> }
export const RUNTIME_RECORDING_SCOPE = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'A' } as const;
export const RUNTIME_RECORDING_TIME = '2026-10-29T10:00:00Z';
export const RUNTIME_RECORDING_PATH = '/api/v1/datalink';

import type { EvidenceScope, RuleEvidenceView } from '../../../src/types/runtimeRuleEvidence';
import { ruleEvidenceFixture } from '../../unit/services/runtimeRuleEvidenceFixture';

export class RuntimeRecordingFixtureError extends Error {
  constructor(readonly status: number, readonly code = 'unavailable') { super(code); }
}
/** Mock immutable-shaped observations validated by production parsers; not backend proof. */
export function runtimeRecordingEvidence(scenario: RuntimeRecordingScenario, deviceId: 'A' | 'B'): RuleEvidenceView {
  const scope: EvidenceScope = { workspace_id: 'W', workspace_revision: 'W1', device_id: deviceId };
  const view = ruleEvidenceFixture(scope); view.observed_at = RUNTIME_RECORDING_TIME;
  const shared = ['shared', 'quota-group', 'quota-global'].includes(scenario);
  if (deviceId === 'B' && !shared) { view.groups = []; return view; }
  const group = view.groups[0]; const saved = group.saved; const effective = group.effective;
  if (!saved || !effective) throw new Error('Expected typed test-only group facts');
  group.group_id = 'G1'; group.as_of = RUNTIME_RECORDING_TIME; group.applied_effective_at = null;
  group.saved_revision = scenario === 'draft' ? 'C' : 'B'; saved.binding.entity_revision = group.saved_revision;
  saved.destination.table = scenario === 'draft' ? 'meter_draft' : 'meter_new';
  saved.interval_seconds = scenario === 'draft' ? 45 : 30;
  saved.recording_profile = { kind: 'instantaneous', interval: saved.interval_seconds, unit: 'second' };
  const members = [{ device_id: 'A', point_id: 'PA1', tag_id: 'TA1', column: 'value', source_value_type: 'float64' as const, circuit_id: null }];
  if (scenario === 'uncertain') members.push({ ...members[0], point_id: `P${deviceId}2`, tag_id: `T${deviceId}2`, column: 'second' });
  if (shared) members.push({ ...members[0], device_id: 'B', point_id: 'PB1', tag_id: 'TB1', column: 'other' });
  saved.members = structuredClone(members); effective.members = structuredClone(members);
  effective.destination.alias = '安全資料表'; effective.effective_at = '2026-10-29T09:00:00Z';
  group.delivery = { current: { binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: 'A' }, intake: null, totals_scope: 'group', shared,
    stages: { collecting: 2, queued: 3, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: scenario === 'shared' ? 10 : 0, skipped: 0 },
    last_commit: scenario === 'shared' ? { group_revision: 'A', connector_revision: '7', record_id: 'fixture-group-receipt', effect_key: 'fixture-effect', payload_digest: 'a'.repeat(64), committed_at: '2026-10-29T09:59:50Z', provenance: 'runtime' } : null,
    sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null },
    backlog: [], counter_replay_quarantines: [], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  const current = group.delivery.current;
  if (!current) throw new Error('Expected mock current observation');
  if (scenario === 'uncertain') {
    current.source_state = 'available';
    current.sources = members.map((member, index) => ({ device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id,
      latest: { state: index === 0 ? 'unknown' : 'partial', window_start: '2026-10-29T09:59:30Z', window_end: '2026-10-29T09:59:45Z', delta_known: index === 0 ? false : null, coverage: index === 0 ? 'unknown' : 'partial', receipt: null }, last_committed: null }));
  }
  if (scenario === 'backlog') {
    saved.destination = { ...saved.destination, connector_id: 'C2', connector_revision: '8', table: 'new_table' };
    saved.enabled = false; effective.enabled = false; current.intake = { group_revision: 'A', state: 'retiring', reason: null };
    group.delivery.backlog = [{ binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A0' },
      destination: { connector_id: 'C1', connector_revision: '7', alias: null, kind: 'sqlite', schema: null, table: 'old_table' }, pending: 3, reasons: [] }];
  }
  if (scenario === 'quota-group' || scenario === 'quota-global') {
    current.quota = { scope: scenario === 'quota-group' ? 'group' : 'global', state: scenario === 'quota-group' ? 'warning' : 'hard_limit',
      used_bytes: scenario === 'quota-group' ? 80 : 125, max_bytes: 100, intake_refused: scenario === 'quota-global' };
  }
  if (scenario === 'no-proof') {
    group.state = 'saved-only'; group.proof = 'unavailable'; group.effective = null; group.effective_revision = null; group.delivery.current = null;
  }
  if (scenario === 'stale') view.workspace_revision = 'W2';
  return view;
}
