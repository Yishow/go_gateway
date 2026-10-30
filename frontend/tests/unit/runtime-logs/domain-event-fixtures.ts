import type { RuntimeLogRecord } from '@/types/runtimeLogs';
import type { RuntimeLogDomainCode } from '@/services/runtimeLogDomainProjection';
import { actionForEvidenceReason } from '@/services/runtimeRuleEvidenceParser';
import type { EvidenceReason } from '@/types/runtimeRuleEvidence';

/** Expected flattened producer contracts from diagnostics.DomainProducerMatrix, not runtime authority. */
export const domainCases = [
  ['domain.source.applied', 'accepted', 'source-runtime-transition', 2, true],
  ['domain.source.rejected', 'blocked', 'source-runtime-transition', 2, false],
  ['domain.collector.failed', 'acquisition', 'acquisition-result', 1, false],
  ['domain.collector.recovered', 'acquisition', 'acquisition-result', 1, true],
  ['domain.mapping.failed', 'blocked', 'mapping-result', 3, false],
  ['domain.mapping.recovered', 'acquisition', 'mapping-result', 3, true],
  ['domain.group.applied', 'collecting', 'intake-status', 4, true],
  ['domain.group.paused', 'blocked', 'intake-status', 4, true],
  ['domain.group.retiring', 'collecting', 'intake-status', 4, true],
  ['domain.delivery.retrying', 'retrying', 'durable-transition', 4, false],
  ['domain.delivery.blocked', 'blocked', 'durable-transition', 4, false],
  ['domain.delivery.quarantined', 'quarantined', 'durable-transition', 4, false],
  ['domain.delivery.unknown', 'unknown', 'durable-transition', 4, false],
  ['domain.delivery.first-confirmed', 'sql-committed', 'runtime-receipt', 4, true],
  ['domain.delivery.recovered', 'sql-committed', 'runtime-receipt', 4, true],
  ['domain.delivery.summary', 'sql-committed', 'runtime-receipt', 4, true],
  ['domain.counter-replay.quarantined', 'quarantined', 'counter-replay-inputs', 4, false],
  ['domain.recording.refused', 'blocked', 'recording-refusal', 4, false],
  ['domain.quota.refused', 'blocked', 'quota-refusal', 4, false],
  ['domain.quota.recovered', 'accepted', 'quota-status', 4, true],
  ['domain.repair.completed', 'accepted', 'operation-result', 4, true],
  ['domain.repair.blocked', 'blocked', 'operation-result', 4, false],
  ['domain.cutover.applied', 'accepted', 'operation-result', 4, true],
  ['domain.cutover.rejected', 'blocked', 'operation-result', 4, false],
] as const;
export const domainInstant = '2026-10-29T10:00:00.123456789Z';
/** Independent, scalar-only fixture. Overrides.fields merges known fields for negative cases. */
export function domainRecord(code: RuntimeLogDomainCode | 'domain.unavailable' = 'domain.delivery.blocked', overrides: Partial<RuntimeLogRecord> = {}): RuntimeLogRecord {
  const fields: RuntimeLogRecord['fields'] = {};
  const rule = domainCases.find(row => row[0] === code);
  if (rule) {
    const [, stage, kind, step, success] = rule;
    const reason: EvidenceReason = success ? 'none' : kind === 'counter-replay-inputs' ? 'counter-replay-proof-missing' :
      kind === 'recording-refusal' ? 'recording-ownership-unavailable' : 'unavailable';
    Object.assign(fields, { event_version: 1, workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'R1',
      locator_step: step, stage, reason, action: actionForEvidenceReason(reason), evidence_kind: kind, reference_id: 'REF1' });
    if (step < 4) fields.device_id = 'D';
    if (step === 1 || step === 3) fields.point_id = 'P';
    if (step === 2) fields.rule_id = 'RULE';
    if (step === 3) Object.assign(fields, { mapping_id: 'M', tag_id: 'T' });
    if (step === 4) Object.assign(fields, { group_id: 'G', group_revision: 'R1', connector_id: 'C', connector_revision: 'C1' });
    if (kind === 'runtime-receipt') Object.assign(fields, { record_id: 'REC1', effect_key: 'REF1', payload_digest: 'DIGEST1',
      committed_at: domainInstant, receipt_provenance: 'runtime' });
    if (kind === 'intake-status') fields.intake_state = code === 'domain.group.applied' ? 'active' : code === 'domain.group.paused' ? 'not_running' : 'retiring';
    if (kind === 'quota-refusal' || kind === 'quota-status') Object.assign(fields, { quota_scope: 'global',
      quota_state: success ? 'warning' : 'hard_limit', quota_used_bytes: success ? 80 : 120, quota_max_bytes: 100, intake_refused: !success });
    if (kind === 'counter-replay-inputs' || kind === 'recording-refusal') fields.pending_samples = '9007199254740993';
  }
  return { instance_id: 'instance-a', sequence: '9007199254740993', timestamp: domainInstant,
    source: 'runtime', level: rule?.[4] ? 'info' : 'warn', code, message: 'Candidate text must never drive domain presentation.', truncated: false,
    ...overrides, fields: { ...fields, ...overrides.fields } };
}
