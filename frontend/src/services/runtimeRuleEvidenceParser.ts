import type {
  CounterReplayQuarantine, CurrentDelivery, DeliveryEvidence, DeliveryStageCounts, EntityLocator, EvidenceAction,
  EvidenceDataType, EvidenceReason, EvidenceResult, EvidenceScope, EvidenceStage,
  HistoricalBacklog, IntakeEvidence, PerSourceEvidence, PollingEvidence, QuotaEvidence, RecordingGroupFacts,
  RecordingGroupView, RecordingMember, RecordingProfile, RevisionBinding, RuleEvidenceView, RuntimeReceipt,
  SafeDestination, SafeTransformStep, SourceEvidenceView, SourceFacts, SourceResult,
  StorageEvidence, StorageEvidenceEvent, TestWriteEvidence, MixedRecordingFacts, MemberCalculationFacts, CalculationResult,
} from '../types/runtimeRuleEvidence'

type Parser<T> = (input: unknown) => T
class Rejected extends Error {
  constructor(readonly reason: 'unavailable' | 'stale-reference' | 'foreign-scope' = 'unavailable') { super(reason) }
}
function refuse(reason?: 'stale-reference' | 'foreign-scope'): never { throw new Rejected(reason) }
const id: Parser<string> = (v) => typeof v === 'string' && /^[A-Za-z0-9_-][A-Za-z0-9_.-]{0,127}$/.test(v) ? v : refuse()
const label: Parser<string> = (v) => typeof v === 'string' && /^[\p{L}\p{N}_ -]{1,128}$/u.test(v) ? v : refuse()
const bool: Parser<boolean> = (v) => typeof v === 'boolean' ? v : refuse()
const numeric = (integer: boolean, minimum = 0, maximum = Number.MAX_SAFE_INTEGER): Parser<number> => (v) =>
  typeof v === 'number' && Number.isFinite(v) && (!integer || Number.isInteger(v)) && v >= minimum && v <= maximum ? v : refuse()
const count = numeric(true)
const duration = numeric(false)
const signed = numeric(false, -Number.MAX_SAFE_INTEGER)
const timestamp: Parser<string> = (v) => {
  if (typeof v !== 'string' || !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,9})?Z$/.test(v)) return refuse()
  const at = new Date(v)
  return Number.isFinite(at.getTime()) && at.toISOString().slice(0, 19) === v.slice(0, 19) ? v : refuse()
}
const optional = <T>(parse: Parser<T>): Parser<T | undefined> => (v) => v === undefined ? undefined : parse(v)
const nullable = <T>(parse: Parser<T>): Parser<T | null> => (v) => v == null ? null : parse(v)
const array = <T>(parse: Parser<T>): Parser<T[]> => (v) => Array.isArray(v) ? v.map(parse) : refuse()
const enumeration = <T extends string>(...allowed: T[]): Parser<T> => (v) =>
  typeof v === 'string' && allowed.some((item) => item === v) ? v as T : refuse()
/** Closed schemas construct new DTOs; no input spreading, arbitrary JSON, or unsafe legacy fallback. */
function object<T extends object>(shape: { [K in keyof T]-?: Parser<T[K]> }): Parser<T> {
  return (v) => {
    if (v === null || typeof v !== 'object' || Array.isArray(v) || (Object.getPrototypeOf(v) !== Object.prototype && Object.getPrototypeOf(v) !== null)) return refuse()
    const input = v as Record<string, unknown>
    if (Reflect.ownKeys(input).some((key) => typeof key !== 'string' || !Object.hasOwn(shape, key))) return refuse()
    const result = {} as T
    for (const key of Object.keys(shape) as (keyof T)[]) {
      const descriptor = Object.getOwnPropertyDescriptor(input, key)
      if (descriptor && !('value' in descriptor)) return refuse()
      const parsed = shape[key](descriptor?.value)
      if (parsed !== undefined) result[key] = parsed
    }
    return result
  }
}
const reasons: EvidenceReason[] = ['none', 'unavailable', 'stale-reference', 'foreign-scope', 'deleted', 'unsupported', 'effective-proof-missing', 'counter-replay-proof-missing', 'counter-replay-proof-conflict', 'counter-replay-proof-unsupported', 'recording-ownership-unavailable', 'recording-ownership-conflict', 'recording-handoff-proof-conflict', 'recording-cutover-finalized-overlap']
const stages: EvidenceStage[] = ['unavailable', 'acquisition', 'accepted', 'collecting', 'queued', 'retrying', 'blocked', 'quarantined', 'unknown', 'skipped', 'sql-committed', 'readback-verified', 'cleanup-completed', 'transport-event']
/** Unknown owner codes/stages never reflect arbitrary text into diagnostics. */
export function parseEvidenceReason(v: unknown): EvidenceReason { return reasons.find((r) => r === v) ?? 'unavailable' }
export function parseEvidenceStage(v: unknown): EvidenceStage { return stages.find((s) => s === v) ?? 'unavailable' }
/** Read-only next step; no retries/activation/repair mutations or mutation permission. */
export function actionForEvidenceReason(reason: EvidenceReason): EvidenceAction {
  switch (reason) {
    case 'none': return 'none'
    case 'foreign-scope': case 'deleted': return 'open-workspace-list'
    case 'unsupported': return 'inspect-settings'
    case 'counter-replay-proof-missing': case 'counter-replay-proof-conflict': case 'counter-replay-proof-unsupported':
    case 'recording-ownership-unavailable': case 'recording-ownership-conflict': case 'recording-handoff-proof-conflict':
    case 'recording-cutover-finalized-overlap': return 'inspect-immutable-evidence'
    default: return 'refresh'
  }
}
const dataType = enumeration<EvidenceDataType>('bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string')
const binding = object<RevisionBinding>({ workspace_id: id, workspace_revision: id, entity_revision: id })
const locator = object<EntityLocator>({
  step: (v) => v === 1 || v === 2 || v === 3 || v === 4 ? v : refuse(), workspace_id: id,
  device_id: nullable(id), rule_id: nullable(id), mapping_id: nullable(id), group_id: nullable(id), group_revision: nullable(id),
})
const transform = object<SafeTransformStep>({
  kind: enumeration('byte-order', 'scale', 'cast', 'unsupported'), order: nullable(enumeration('ABCD', 'BADC', 'CDAB', 'DCBA')),
  multiplier: nullable(signed), offset: nullable(signed), target_type: nullable(dataType),
})
const polling = object<PollingEvidence>({ group_id: id, interval_ms: numeric(true, 1), source: enumeration('polling-group') })
const sourceFacts = object<SourceFacts>({
  binding, mapping_revision: nullable(id), enabled: bool, address: nullable(id), protocol_function: nullable(id),
  read_type: nullable(dataType), result_type: nullable(dataType), polling: nullable(polling), transforms: array(transform),
})
const source = object<SourceEvidenceView>({
  rule_id: id, point_id: id, tag_id: nullable(id), mapping_id: nullable(id), saved_revision: nullable(id), effective_revision: nullable(id),
  state: enumeration('saved-only', 'effective', 'unavailable'), proof: enumeration('unavailable', 'immutable-source-facts'),
  saved: nullable(sourceFacts), effective: nullable(sourceFacts),
})
const destination = object<SafeDestination>({
  connector_id: id, connector_revision: id, alias: nullable(label), kind: enumeration('sqlite', 'postgres', 'mysql', 'sqlserver', 'unavailable'),
  schema: nullable(id), table: id,
})
const calculationMethod = enumeration('latest', 'mean', 'delta')
const resultType = enumeration('bool', 'text', 'int64', 'uint64', 'decimal', 'float64')
const mixedRecording = object<MixedRecordingFacts>({
  record_id: v => typeof v === 'string' && /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/.test(v) && v !== '00000000-0000-0000-0000-000000000000' ? v : refuse(),
  result_shape: enumeration('member-calculation-v1'),
})
const memberCalculation = object<MemberCalculationFacts>({
  method: calculationMethod, interval: numeric(true, 1, 9_223_372_036), unit: enumeration('second', 'minute'),
  alignment: enumeration('utc_epoch'), required: bool, source_revision: id, mapping_revision: id,
  measurement_id: nullable(id), result_type: resultType,
})
const mixedTimestamp: Parser<string> = v => {
  const parsed = timestamp(v); const ns = timeNanos(parsed)
  return ns < -9223372036854775808n || ns > 9223372036854775807n ? refuse() : parsed
}
const calculationResult = object<CalculationResult>({
  result_shape: enumeration('member-calculation-v1'), method: calculationMethod,
  interval_seconds: numeric(true, 1, 9_223_372_036), alignment: enumeration('utc_epoch'), due_at: mixedTimestamp,
  next_due: nullable(mixedTimestamp), outcome_digest: v => typeof v === 'string' && /^[0-9a-f]{64}$/.test(v) ? v : refuse(),
  outcome_kind: enumeration('row', 'skipped', 'no_data'), status: enumeration('ok', 'missing', 'bad', 'stale', 'invalid', 'not_due'),
  reason: nullable(id), result_present: bool, row_emitted: bool, raw_present: bool, source_type: dataType, result_type: resultType,
  aligned_start: nullable(mixedTimestamp), aligned_end: nullable(mixedTimestamp), actual_start: nullable(mixedTimestamp), actual_end: nullable(mixedTimestamp), rounded: nullable(bool),
})
const member = object<RecordingMember>({ calculation: optional(memberCalculation), source_value_type: nullable(dataType), circuit_id: nullable(id), device_id: id, point_id: id, tag_id: id, column: id })
const profile = object<RecordingProfile>({ kind: enumeration('instantaneous', 'cumulative'), interval: numeric(true, 1, 600), unit: enumeration('second', 'minute') })
const groupFacts = object<RecordingGroupFacts>({
  mixed_recording: optional(mixedRecording),
  recording_profile: nullable(profile), allowed_lateness_seconds: nullable(count),
  ends_at: nullable(timestamp), binding, effective_at: nullable(timestamp), enabled: bool, destination, members: array(member),
  layout: enumeration('managed_wide_v1', 'managed_long_v1', 'managed_wide_profile_v1', 'unavailable'),
  storage_mode: enumeration('plain', 'timescale', 'native_monthly', 'unavailable'), interval_seconds: nullable(numeric(true, 1)),
  reduction: enumeration('latest', 'mean', 'cumulative_delta', 'member_calculation', 'unavailable'), write_mode: enumeration('append', 'latest', 'unavailable'),
  incomplete_policy: enumeration('skip_row', 'partial', 'unavailable'),
})
const stageCounts = object<DeliveryStageCounts>({ collecting: count, queued: count, retrying: count, blocked: count, quarantined: count, unknown: count, sql_committed: count, skipped: count })
const receipt = object<RuntimeReceipt>({
  group_revision: id, connector_revision: id, record_id: id, effect_key: id, payload_digest: id, committed_at: timestamp, provenance: enumeration('runtime'),
})
const sourceResult = object<SourceResult>({
  calculation: optional(calculationResult),
  state: enumeration('available', 'missing', 'unknown', 'partial', 'not_due', 'unavailable'), window_start: nullable(timestamp), window_end: nullable(timestamp),
  delta_known: nullable(bool), coverage: nullable(enumeration('complete', 'partial', 'unknown')), receipt: nullable(receipt),
})
const perSource = object<PerSourceEvidence>({ device_id: id, point_id: id, tag_id: id, latest: nullable(sourceResult), last_committed: nullable(sourceResult) })
const quota = object<QuotaEvidence>({ scope: enumeration('global', 'group'), state: enumeration('unconfigured', 'ok', 'warning', 'hard_limit', 'unavailable'), used_bytes: count, max_bytes: count, intake_refused: bool })
const storageEvent = object<StorageEvidenceEvent>({ kind: enumeration('confirmed', 'recreated_by_confirmation', 'reconstructed', 'repair_blocked', 'unavailable'), at: nullable(timestamp) })
const storage = object<StorageEvidence>({
  generation: count, history_may_be_lost: bool, reconstructions: count, last_reconstructed_at: nullable(timestamp),
  repair_state: enumeration('unavailable', 'none', 'blocked'), events: array(storageEvent),
})
const intake = object<IntakeEvidence>({ group_revision: id, state: enumeration('active', 'retiring', 'blocked', 'not_running', 'unavailable'), reason: nullable(enumeration(...reasons)) })
const current = object<CurrentDelivery>({
  intake: nullable(intake),
  binding, totals_scope: enumeration('group'), shared: bool, stages: stageCounts, last_commit: nullable(receipt), sources: array(perSource),
  source_state: enumeration('available', 'unavailable'), next_due_at: nullable(timestamp),
  next_retry_at: nullable(timestamp), attempt_duration_ms: nullable(duration), quota: nullable(quota), storage: nullable(storage),
})
const backlog = object<HistoricalBacklog>({ binding, destination, pending: count, reasons: array(enumeration(...reasons)) })
const testWrite = object<TestWriteEvidence>({
  operation_id: id, provenance: enumeration('test-write'), written_at: nullable(timestamp), readback_at: nullable(timestamp), cleanup_at: nullable(timestamp),
})
const quarantine = object<CounterReplayQuarantine>({
  workspace_id: id, group_id: id, group_revision: id, pending_samples: count,
  reason: enumeration('counter-replay-proof-missing', 'counter-replay-proof-conflict', 'counter-replay-proof-unsupported'),
  action: enumeration('preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance'),
})
const delivery = object<DeliveryEvidence>({ counter_replay_quarantines: nullable(array(quarantine)), oldest_pending_scope: enumeration('group-all-revisions'), oldest_pending_seconds: nullable(duration), current: nullable(current), backlog: array(backlog), test_write: nullable(testWrite) })
const group = object<RecordingGroupView>({
  effective_revision: nullable(id), as_of: nullable(timestamp), applied_effective_at: nullable(timestamp),
  group_id: id, saved_revision: nullable(id), applied_revision: nullable(id), state: enumeration('saved-only', 'effective', 'unavailable'),
  proof: enumeration('unavailable', 'applied-snapshot-and-layout'), saved: nullable(groupFacts), effective: nullable(groupFacts), delivery: nullable(delivery),
})
const view = object<RuleEvidenceView>({
  contract_version: (v) => v === 1 ? 1 : refuse(), workspace_id: id, workspace_revision: id, device_id: id,
  observed_at: nullable(timestamp), freshness: enumeration('fresh', 'stale', 'unavailable'), source: nullable(source), groups: array(group), locator: nullable(locator),
})
function checkBinding(b: RevisionBinding, revision: string | null, scope: EvidenceScope) {
  if (b.workspace_id !== scope.workspace_id) refuse('foreign-scope')
  if (b.workspace_revision !== scope.workspace_revision || revision === null || b.entity_revision !== revision) refuse('stale-reference')
}
function checkTransform(s: SafeTransformStep) {
  const valid = s.kind === 'scale' ? s.multiplier !== null && s.offset !== null && s.order === null && s.target_type === null
    : s.kind === 'cast' ? s.target_type !== null && s.multiplier === null && s.offset === null && s.order === null
      : s.kind === 'byte-order' ? s.order !== null && s.multiplier === null && s.offset === null && s.target_type === null
        : s.order === null && s.multiplier === null && s.offset === null && s.target_type === null
  if (!valid) refuse()
}
function checkSource(s: SourceEvidenceView | null, scope: EvidenceScope) {
  if (!s) return
  for (const [facts, revision] of [[s.saved, s.saved_revision], [s.effective, s.effective_revision]] as const) {
    if (!facts) continue
    checkBinding(facts.binding, revision, scope); facts.transforms.forEach(checkTransform)
  }
  if (s.state === 'saved-only' && (!s.saved || s.effective)) refuse()
  if (s.state !== 'effective' && s.effective) refuse()
  if (s.state === 'effective' && (s.effective?.transforms.some((step) => step.kind === 'unsupported') || s.proof !== 'immutable-source-facts' || !s.effective || !s.effective_revision || !s.effective.mapping_revision || !s.effective.read_type || !s.effective.result_type || !s.effective.address || !s.effective.polling)) {
    s.state = 'unavailable'; s.effective = null; s.proof = 'unavailable'
  }
}
function checkLocator(l: EntityLocator, scope: EvidenceScope) {
  if (l.workspace_id !== scope.workspace_id || (l.device_id !== null && l.device_id !== scope.device_id)) refuse('foreign-scope')
  const valid = l.step === 1 ? l.device_id !== null && l.rule_id === null && l.mapping_id === null && l.group_id === null && l.group_revision === null
    : l.step === 2 ? l.device_id !== null && l.rule_id !== null && l.mapping_id === null && l.group_id === null && l.group_revision === null
      : l.step === 3 ? l.device_id !== null && l.mapping_id !== null && l.rule_id === null && l.group_id === null && l.group_revision === null
        : l.group_id !== null && l.rule_id === null && l.mapping_id === null
  if (!valid) refuse()
}
function checkReceipt(r: RuntimeReceipt | null, facts: RecordingGroupFacts) {
  if (r && (r.group_revision !== facts.binding.entity_revision || r.connector_revision !== facts.destination.connector_revision)) refuse('stale-reference')
}
const exactTime = (value: string) => `${value.slice(0, 19)}.${(value.split('.')[1]?.slice(0, -1) ?? '').padEnd(9, '0')}Z`
const sourceKey = (s: RecordingMember | PerSourceEvidence) => `${s.device_id}:${s.point_id}:${s.tag_id}`
function checkCurrent(c: CurrentDelivery, facts: RecordingGroupFacts) {
  checkReceipt(c.last_commit, facts)
  if ((c.stages.sql_committed > 0) !== (c.last_commit !== null)) refuse()
  const devices = new Set(facts.members.map((m) => m.device_id)); const members = new Set(facts.members.map(sourceKey))
  if (c.shared !== (devices.size > 1) || (c.source_state === 'unavailable' && c.sources.length > 0)) refuse()
  const seen = new Set<string>()
  for (const s of c.sources) {
    const key = sourceKey(s)
    if (!members.has(key)) refuse('foreign-scope')
    if (seen.has(key)) refuse()
    seen.add(key)
    for (const r of [s.latest, s.last_committed]) {
      if (!r) continue
      checkReceipt(r.receipt, facts)
      if ((r.window_start === null) !== (r.window_end === null)) refuse()
      if (r.window_start !== null && r.window_end !== null && exactTime(r.window_start) >= exactTime(r.window_end)) refuse()
      checkMixedResult(r, facts, s)
      if (!r.calculation && ((r.delta_known === false && r.state !== 'unknown') || (r.coverage === 'partial' && r.state !== 'partial'))) refuse()
      if (r.receipt && ((r.window_start === null && r.calculation?.status !== 'not_due') || r.state === 'unavailable' || !c.last_commit)) refuse()
    }
    if (s.last_committed && !s.last_committed.receipt) refuse()
  }
}
function timeNanos(value: string): bigint {
  const seconds = BigInt(Date.parse(`${value.slice(0, 19)}Z`)) * 1_000_000n
  return seconds + BigInt((value.split('.')[1]?.slice(0, -1) ?? '').padEnd(9, '0'))
}
function checkMixedFacts(f: RecordingGroupFacts) {
  if (!f.mixed_recording) {
    if (f.reduction === 'member_calculation' || f.members.some(m => m.calculation !== undefined)) refuse()
    return
  }
  if (f.recording_profile !== null || f.interval_seconds !== null || f.reduction !== 'member_calculation') refuse()
  for (const m of f.members) {
    const c = m.calculation
    if (!c || m.source_value_type === null || (c.unit === 'minute' && c.interval > 153_722_867) ||
      (c.method === 'delta' && c.measurement_id === null)) refuse()
  }
}
function checkMixedResult(r: SourceResult, facts: RecordingGroupFacts, source: PerSourceEvidence) {
  const c = r.calculation
  if (!c) {
    if ((facts.mixed_recording && r.state !== 'unavailable') || r.state === 'not_due') refuse()
    return
  }
  if (!facts.mixed_recording || c.result_shape !== facts.mixed_recording.result_shape) refuse()
  const member = facts.members.find(m => sourceKey(m) === sourceKey(source))
  if (!member?.calculation || !member.source_value_type) refuse('foreign-scope')
  const intent = member.calculation
  if (c.method !== intent.method || c.interval_seconds !== intent.interval * (intent.unit === 'minute' ? 60 : 1) ||
    c.alignment !== intent.alignment || c.source_type !== member.source_value_type || c.result_type !== intent.result_type) refuse('stale-reference')
  if (c.row_emitted !== (c.outcome_kind === 'row') || (c.result_present !== (c.status === 'ok')) || (c.result_present && !c.raw_present) || (r.receipt && !c.row_emitted)) refuse()
  if ((r.state === 'available' && !c.result_present) || ((r.state === 'not_due') !== (c.status === 'not_due'))) refuse()
  const width = BigInt(c.interval_seconds) * 1_000_000_000n
  const due = timeNanos(c.due_at)
  if (c.next_due !== null && exactTime(c.next_due) <= exactTime(c.due_at)) refuse()
  if (c.rounded !== null && c.method !== 'mean') refuse()
  if (c.method === 'delta' ? r.delta_known !== c.result_present : r.delta_known !== null) refuse()
  if (c.status === 'not_due') {
    if (facts.layout !== 'managed_wide_profile_v1' || r.state !== 'not_due' || c.result_present || c.raw_present ||
      c.aligned_start !== null || c.aligned_end !== null || c.actual_start !== null || c.actual_end !== null ||
      r.window_start !== null || r.window_end !== null) refuse()
    if (c.next_due !== null) {
      const aligned = due - ((due % width + width) % width)
      if (aligned < -9223372036854775808n || aligned === due || timeNanos(c.next_due) !== aligned + width) refuse()
    }
  } else {
    if (c.actual_start === null || c.actual_end === null || c.aligned_start === null || c.aligned_end === null || r.window_start === null || r.window_end === null) refuse()
    const aligned = timeNanos(c.aligned_start)
    if (aligned % width !== 0n || timeNanos(c.aligned_end) !== aligned + width || timeNanos(c.aligned_end) !== due ||
      (c.next_due !== null && timeNanos(c.next_due) !== due + width)) refuse()
    if (exactTime(c.actual_start) !== exactTime(r.window_start) || exactTime(c.actual_end) !== exactTime(r.window_end) ||
      exactTime(c.actual_start) >= exactTime(c.actual_end) || exactTime(c.aligned_start) > exactTime(c.actual_start) || exactTime(c.actual_end) > exactTime(c.aligned_end)) refuse()
  }
}
function checkSummaryProfile(f: RecordingGroupFacts) {
  checkMixedFacts(f)
  const p = f.recording_profile
  if (!p) return
  if (p.kind === 'instantaneous' ? p.unit !== 'second' || p.interval > 300 || !['unavailable', 'latest', 'mean'].includes(f.reduction)
    : p.unit !== 'minute' || !['unavailable', 'cumulative_delta'].includes(f.reduction)) refuse()
  const seconds = p.kind === 'instantaneous' ? p.interval : p.interval * 60
  if (f.interval_seconds !== null && f.interval_seconds !== seconds) refuse()
}
function checkGroup(g: RecordingGroupView, scope: EvidenceScope, observedAt: string | null) {
  // Legacy fallback checks consistency only, never establishes current proof.
  const revision = g.effective_revision ?? g.effective?.binding.entity_revision ?? null
  if (g.as_of && observedAt && exactTime(g.as_of) !== exactTime(observedAt)) refuse('stale-reference')
  if (g.effective?.effective_at && g.as_of && exactTime(g.effective.effective_at) > exactTime(g.as_of)) refuse('stale-reference')
  let selected = false
  for (const [facts, factRevision] of [[g.saved, g.saved_revision], [g.effective, revision]] as const) {
    if (!facts) continue
    checkSummaryProfile(facts)
    checkBinding(facts.binding, factRevision, scope)
    if (new Set(facts.members.map(sourceKey)).size !== facts.members.length) refuse()
    selected ||= facts.members.some((m) => m.device_id === scope.device_id)
  }
  if (!selected && (g.saved || g.effective)) refuse('foreign-scope')
  if (g.state === 'saved-only' && (!g.saved || g.effective)) refuse()
  if (g.state !== 'effective' && g.effective) refuse()
  if (g.delivery) {
    const seen = new Set<string>()
    for (const q of g.delivery.counter_replay_quarantines ?? []) {
      if (q.workspace_id !== scope.workspace_id || q.group_id !== g.group_id) refuse('foreign-scope')
      if (seen.has(q.group_revision)) refuse()
      seen.add(q.group_revision)
    }
    for (const b of g.delivery.backlog) {
      checkBinding(b.binding, b.binding.entity_revision, scope)
      if (revision !== null && b.binding.entity_revision === revision) refuse('stale-reference')
    }
    if (g.delivery.current) {
      const c = g.delivery.current; checkBinding(c.binding, revision, scope)
      if (c.intake?.state === 'active' && g.effective && (!g.effective.enabled || (g.effective.ends_at && g.as_of && exactTime(g.effective.ends_at) <= exactTime(g.as_of)))) refuse()
      if (c.intake && c.intake.group_revision !== c.binding.entity_revision) refuse('stale-reference')
      const receipts = [c.last_commit, ...c.sources.flatMap((s) => [s.latest?.receipt, s.last_committed?.receipt])]
      const identities = new Map<string, RuntimeReceipt>()
      for (const r of receipts) {
        if (!r) continue
        if (r.group_revision !== revision) refuse('stale-reference')
        for (const key of [`effect:${r.effect_key}`, `record:${r.record_id}`]) {
          const old = identities.get(key)
          if (old && (old.group_revision !== r.group_revision || old.connector_revision !== r.connector_revision || old.record_id !== r.record_id || old.effect_key !== r.effect_key || old.payload_digest !== r.payload_digest || old.committed_at !== r.committed_at || old.provenance !== r.provenance)) refuse('stale-reference')
          identities.set(key, r)
        }
      }
      if (g.effective) checkCurrent(c, g.effective)
    }
  }
  if (g.state === 'effective' && (!g.effective_revision || !g.as_of || !observedAt || g.proof !== 'applied-snapshot-and-layout' || !g.effective || !revision || !g.effective.effective_at || g.effective.layout === 'unavailable' || (!g.effective.mixed_recording && g.effective.interval_seconds === null))) {
    g.state = 'unavailable'; g.effective = null; g.proof = 'unavailable'
  }
  if (!g.delivery?.current) return
  if (!g.effective) { g.delivery.current = null; return }
}
function result<T>(parse: () => T): EvidenceResult<T> {
  try { return { ok: true, value: parse() } }
  catch (error) {
    const reason = error instanceof Rejected ? error.reason : 'unavailable'
    return { ok: false, reason, action: actionForEvidenceReason(reason) }
  }
}
/** Pure unmounted parser; equality checks do not establish ownership or verify immutable proof. */
export function parseRuntimeRuleEvidence(input: unknown, scope: EvidenceScope): EvidenceResult<RuleEvidenceView> {
  return result(() => {
    id(scope.workspace_id); id(scope.workspace_revision); id(scope.device_id)
    const v = view(input)
    if (v.workspace_id !== scope.workspace_id || v.device_id !== scope.device_id) refuse('foreign-scope')
    if (v.workspace_revision !== scope.workspace_revision) refuse('stale-reference')
    if (v.freshness !== 'unavailable' && v.observed_at === null) refuse()
    checkSource(v.source, scope)
    const groups = new Set<string>()
    for (const g of v.groups) { if (groups.has(g.group_id)) refuse(); groups.add(g.group_id); checkGroup(g, scope, v.observed_at) }
    if (v.locator) {
      const l = v.locator; checkLocator(l, scope)
      if ((l.group_id !== null && !groups.has(l.group_id)) || (l.rule_id !== null && l.rule_id !== v.source?.rule_id) || (l.mapping_id !== null && l.mapping_id !== v.source?.mapping_id)) refuse('stale-reference')
    }
    return v
  })
}
/** Pure locator parsing; future resolver still checks ownership/deletion, never substitutes names. */
export function parseEntityLocator(input: unknown, scope: EvidenceScope): EvidenceResult<EntityLocator> {
  return result(() => { const l = locator(input); checkLocator(l, scope); return l })
}
