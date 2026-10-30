import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';
import { parsePreciseInstant, comparePreciseInstants } from '../utils/preciseInstant';
import { actionForEvidenceReason, parseEvidenceReason, parseEvidenceStage } from './runtimeRuleEvidenceParser';
import { buildStudioEntityLocatorTarget } from './studioEntityLocator';
import type { RuntimeLogLevel, RuntimeLogRecord } from '../types/runtimeLogs';
import type { EntityLocator, EvidenceAction, EvidenceReason, EvidenceStage, IntakeEvidence, QuotaEvidence, RevisionBinding, RuntimeReceipt } from '../types/runtimeRuleEvidence';

export type RuntimeLogDomainCode = 'domain.source.applied' | 'domain.source.rejected' | 'domain.collector.failed' | 'domain.collector.recovered' |
  'domain.mapping.failed' | 'domain.mapping.recovered' | 'domain.group.applied' | 'domain.group.paused' | 'domain.group.retiring' |
  'domain.delivery.retrying' | 'domain.delivery.blocked' | 'domain.delivery.quarantined' | 'domain.delivery.unknown' |
  'domain.delivery.first-confirmed' | 'domain.delivery.recovered' | 'domain.delivery.summary' | 'domain.counter-replay.quarantined' |
  'domain.recording.refused' | 'domain.quota.refused' | 'domain.quota.recovered' | 'domain.repair.completed' | 'domain.repair.blocked' |
  'domain.cutover.applied' | 'domain.cutover.rejected';
export type DomainEvidenceKind = 'source-runtime-transition' | 'acquisition-result' | 'mapping-result' | 'intake-status' |
  'durable-transition' | 'runtime-receipt' | 'counter-replay-inputs' | 'recording-refusal' | 'quota-refusal' | 'quota-status' | 'operation-result';
export interface DomainSummary { firstSeenAt: string; lastSeenAt: string; occurrences: string }
export interface DomainPresentation {
  kind: 'domain'; code: RuntimeLogDomainCode; level: RuntimeLogLevel; timestamp: string;
  binding: RevisionBinding; locator: EntityLocator; stage: EvidenceStage; reason: EvidenceReason; action: EvidenceAction;
  evidenceKind: DomainEvidenceKind; referenceID: string; pointID: string | null; tagID: string | null;
  connectorID: string | null; connectorRevision: string | null; receipt: RuntimeReceipt | null;
  intakeState: IntakeEvidence['state'] | null; quota: QuotaEvidence | null; pendingSamples: string | null; summary: DomainSummary | null;
}
export type RuntimeLogDomainProjection = { kind: 'legacy' } | { kind: 'unavailable'; timestamp: string | null; summary: DomainSummary | null } | DomainPresentation;
// Consumer shape of the fixed Go DomainProducerMatrix; never establishes producer origin or current ownership.
const rules = {
  'domain.source.applied': ['accepted', 'source-runtime-transition', 2, true],
  'domain.source.rejected': ['blocked', 'source-runtime-transition', 2, false],
  'domain.collector.failed': ['acquisition', 'acquisition-result', 1, false],
  'domain.collector.recovered': ['acquisition', 'acquisition-result', 1, true],
  'domain.mapping.failed': ['blocked', 'mapping-result', 3, false],
  'domain.mapping.recovered': ['acquisition', 'mapping-result', 3, true],
  'domain.group.applied': ['collecting', 'intake-status', 4, true],
  'domain.group.paused': ['blocked', 'intake-status', 4, true],
  'domain.group.retiring': ['collecting', 'intake-status', 4, true],
  'domain.delivery.retrying': ['retrying', 'durable-transition', 4, false],
  'domain.delivery.blocked': ['blocked', 'durable-transition', 4, false],
  'domain.delivery.quarantined': ['quarantined', 'durable-transition', 4, false],
  'domain.delivery.unknown': ['unknown', 'durable-transition', 4, false],
  'domain.delivery.first-confirmed': ['sql-committed', 'runtime-receipt', 4, true],
  'domain.delivery.recovered': ['sql-committed', 'runtime-receipt', 4, true],
  'domain.delivery.summary': ['sql-committed', 'runtime-receipt', 4, true],
  'domain.counter-replay.quarantined': ['quarantined', 'counter-replay-inputs', 4, false],
  'domain.recording.refused': ['blocked', 'recording-refusal', 4, false],
  'domain.quota.refused': ['blocked', 'quota-refusal', 4, false],
  'domain.quota.recovered': ['accepted', 'quota-status', 4, true],
  'domain.repair.completed': ['accepted', 'operation-result', 4, true],
  'domain.repair.blocked': ['blocked', 'operation-result', 4, false],
  'domain.cutover.applied': ['accepted', 'operation-result', 4, true],
  'domain.cutover.rejected': ['blocked', 'operation-result', 4, false],
} as const satisfies Record<RuntimeLogDomainCode, readonly [EvidenceStage, DomainEvidenceKind, EntityLocator['step'], boolean]>;
const common = ['event_version', 'workspace_id', 'workspace_revision', 'entity_revision', 'locator_step', 'stage', 'reason', 'action', 'evidence_kind', 'reference_id'];
const summaryKeys = ['first_seen_at', 'last_seen_at', 'occurrences'];
const receiptKeys = ['record_id', 'effect_key', 'payload_digest', 'committed_at', 'receipt_provenance'];
const quotaKeys = ['quota_scope', 'quota_state', 'quota_used_bytes', 'quota_max_bytes', 'intake_refused'];
const envelopeKeys = ['instance_id', 'sequence', 'timestamp', 'level', 'source', 'code', 'message', 'fields', 'truncated'];
const stepKeys = {
  1: ['device_id', 'point_id'], 2: ['device_id', 'rule_id'], 3: ['device_id', 'mapping_id', 'point_id', 'tag_id'],
  4: ['group_id', 'group_revision', 'connector_id', 'connector_revision'],
};
function refuse(): never { throw new Error('domain_presentation_unavailable'); }
/** Read own data descriptors only; hostile getters, prototypes, symbols and extra fields are not copied. */
function object(value: unknown, keys: readonly string[]): Record<string, unknown> {
  if (!value || typeof value !== 'object' || Object.getPrototypeOf(value) !== Object.prototype) return refuse();
  const output: Record<string, unknown> = {};
  for (const key of Reflect.ownKeys(value)) {
    const descriptor = Object.getOwnPropertyDescriptor(value, key);
    if (typeof key !== 'string' || !keys.includes(key) || !descriptor || !('value' in descriptor)) return refuse();
    output[key] = descriptor.value;
  }
  return output;
}
function identity(value: unknown): string { return safeRuntimeIdentity(value) || refuse(); }
function utc(value: unknown): string {
  if (typeof value !== 'string' || !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,9})?Z$/.test(value)) return refuse();
  const at = parsePreciseInstant(value);
  // Date.parse alone normalizes invalid calendar dates; compare whole UTC seconds without rewriting input bytes.
  if (!at || new Date(at.seconds * 1000).toISOString().slice(0, 19) !== value.slice(0, 19)) return refuse();
  return value;
}
function uint64(value: unknown, positive = false): string {
  if (typeof value !== 'string' || !/^(0|[1-9][0-9]{0,19})$/.test(value) || BigInt(value) > 18446744073709551615n || (positive && value === '0')) return refuse();
  return value;
}
function summary(fields: Record<string, unknown>): DomainSummary | null {
  if (!summaryKeys.some(key => Object.hasOwn(fields, key))) return null;
  const firstSeenAt = utc(fields.first_seen_at), lastSeenAt = utc(fields.last_seen_at), occurrences = uint64(fields.occurrences, true);
  const first = parsePreciseInstant(firstSeenAt), last = parsePreciseInstant(lastSeenAt);
  if (!first || !last || comparePreciseInstants(first, last) > 0) return refuse();
  return { firstSeenAt, lastSeenAt, occurrences };
}
function compatibleReason(code: RuntimeLogDomainCode, kind: DomainEvidenceKind, reason: EvidenceReason): boolean {
  switch (reason) {
    case 'counter-replay-proof-missing': case 'counter-replay-proof-conflict': case 'counter-replay-proof-unsupported':
      return code === 'domain.counter-replay.quarantined' && kind === 'counter-replay-inputs';
    case 'recording-ownership-unavailable': case 'recording-ownership-conflict': case 'recording-handoff-proof-conflict': case 'recording-cutover-finalized-overlap':
      return kind === 'durable-transition' || code === 'domain.recording.refused' || code === 'domain.cutover.rejected';
    default: return kind !== 'counter-replay-inputs' && kind !== 'recording-refusal';
  }
}
function bytes(value: unknown): number { return typeof value === 'number' && Number.isSafeInteger(value) && value >= 0 ? value : refuse(); }
function quota(fields: Record<string, unknown>, code: RuntimeLogDomainCode): QuotaEvidence {
  const scope = fields.quota_scope, state = fields.quota_state, used_bytes = bytes(fields.quota_used_bytes), max_bytes = bytes(fields.quota_max_bytes);
  if (scope !== 'global' && scope !== 'group') return refuse();
  if (code === 'domain.quota.recovered') {
    if ((state !== 'ok' && state !== 'warning') || fields.intake_refused !== false || max_bytes === 0 || used_bytes >= max_bytes) return refuse();
    return { scope, state, used_bytes, max_bytes, intake_refused: false };
  }
  if (state !== 'hard_limit' || fields.intake_refused !== true) return refuse();
  return { scope, state, used_bytes, max_bytes, intake_refused: true };
}
function project(record: Record<string, unknown>, code: RuntimeLogDomainCode, timestamp: string): DomainPresentation {
  const [stage, evidenceKind, step, success] = rules[code];
  const branch = evidenceKind === 'runtime-receipt' ? receiptKeys : evidenceKind === 'intake-status' ? ['intake_state'] :
    evidenceKind === 'quota-refusal' || evidenceKind === 'quota-status' ? quotaKeys :
      evidenceKind === 'counter-replay-inputs' || evidenceKind === 'recording-refusal' ? ['pending_samples'] : [];
  const fields = object(record.fields, [...common, ...stepKeys[step], ...branch, ...summaryKeys]);
  const reason = parseEvidenceReason(fields.reason), action = actionForEvidenceReason(reason);
  if (fields.event_version !== 1 || fields.locator_step !== step || fields.stage !== stage || parseEvidenceStage(fields.stage) !== stage ||
      fields.evidence_kind !== evidenceKind || reason !== fields.reason || action !== fields.action || success !== (reason === 'none') ||
      !compatibleReason(code, evidenceKind, reason) || record.level !== (success ? 'info' : 'warn')) return refuse();
  const binding: RevisionBinding = { workspace_id: identity(fields.workspace_id), workspace_revision: identity(fields.workspace_revision), entity_revision: identity(fields.entity_revision) };
  const locator: EntityLocator = { step, workspace_id: binding.workspace_id, device_id: step < 4 ? identity(fields.device_id) : null,
    rule_id: step === 2 ? identity(fields.rule_id) : null, mapping_id: step === 3 ? identity(fields.mapping_id) : null,
    group_id: step === 4 ? identity(fields.group_id) : null, group_revision: step === 4 ? identity(fields.group_revision) : null };
  if (!buildStudioEntityLocatorTarget(locator) || (step === 4 && locator.group_revision !== binding.entity_revision)) return refuse();
  const referenceID = identity(fields.reference_id), connectorID = step === 4 ? identity(fields.connector_id) : null,
    connectorRevision = step === 4 ? identity(fields.connector_revision) : null;
  let receipt: RuntimeReceipt | null = null;
  if (evidenceKind === 'runtime-receipt') {
    if (!connectorRevision || fields.effect_key !== referenceID || fields.receipt_provenance !== 'runtime') return refuse();
    receipt = { group_revision: binding.entity_revision, connector_revision: connectorRevision, record_id: identity(fields.record_id),
      effect_key: identity(fields.effect_key), payload_digest: identity(fields.payload_digest), committed_at: utc(fields.committed_at), provenance: 'runtime' };
  }
  let intakeState: IntakeEvidence['state'] | null = null;
  if (evidenceKind === 'intake-status') {
    intakeState = code === 'domain.group.applied' ? 'active' : code === 'domain.group.paused' ? 'not_running' : 'retiring';
    if (fields.intake_state !== intakeState) return refuse();
  }
  return { kind: 'domain', code, level: success ? 'info' : 'warn', timestamp, binding, locator, stage, reason, action, evidenceKind, referenceID,
    pointID: step === 1 || step === 3 ? identity(fields.point_id) : null, tagID: step === 3 ? identity(fields.tag_id) : null,
    connectorID, connectorRevision, receipt, intakeState, quota: evidenceKind === 'quota-refusal' || evidenceKind === 'quota-status' ? quota(fields, code) : null,
    pendingSamples: evidenceKind === 'counter-replay-inputs' || evidenceKind === 'recording-refusal' ? uint64(fields.pending_samples) : null, summary: summary(fields) };
}
/** Detached presentation only. Current ownership and revision checks remain with the native locator receiver. */
export function projectRuntimeLogDomain(record: RuntimeLogRecord): RuntimeLogDomainProjection {
  let timestamp: string | null = null;
  try {
    if (!record || typeof record !== 'object' || Object.getPrototypeOf(record) !== Object.prototype) return refuse();
    const codeDescriptor = Object.getOwnPropertyDescriptor(record, 'code');
    if (!codeDescriptor || !('value' in codeDescriptor) || typeof codeDescriptor.value !== 'string') return refuse();
    const code = codeDescriptor.value;
    if (!code.startsWith('domain.')) return { kind: 'legacy' };
    const input = object(record, envelopeKeys);
    timestamp = utc(input.timestamp);
    if (input.source !== 'runtime' || typeof input.message !== 'string' || typeof input.truncated !== 'boolean') return refuse();
    if (code === 'domain.unavailable') {
      if (input.level !== 'warn') return refuse();
      return { kind: 'unavailable', timestamp, summary: summary(object(input.fields, summaryKeys)) };
    }
    if (!Object.hasOwn(rules, code)) return refuse();
    return project(input, code as RuntimeLogDomainCode, timestamp);
  } catch { return { kind: 'unavailable', timestamp, summary: null }; }
}
