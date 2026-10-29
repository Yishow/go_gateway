import type {
  CounterReplayQuarantine, CurrentDelivery, DeliveryEvidence, DeliveryStageCounts, EntityLocator, EvidenceAction,
  EvidenceDataType, EvidenceReason, EvidenceResult, EvidenceScope, EvidenceStage,
  HistoricalBacklog, IntakeEvidence, PerSourceEvidence, PollingEvidence, QuotaEvidence, RecordingGroupFacts,
  RecordingGroupView, RecordingMember, RecordingProfile, RevisionBinding, RuleEvidenceView, RuntimeReceipt,
  SafeDestination, SafeTransformStep, SourceEvidenceView, SourceFacts, SourceResult,
  StorageEvidence, StorageEvidenceEvent, TestWriteEvidence,
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
const nullable = <T>(parse: Parser<T>): Parser<T | null> => (v) => v == null ? null : parse(v)
const array = <T>(parse: Parser<T>): Parser<T[]> => (v) => Array.isArray(v) ? v.map(parse) : refuse()
const enumeration = <T extends string>(...allowed: T[]): Parser<T> => (v) =>
  typeof v === 'string' && allowed.some((item) => item === v) ? v as T : refuse()
/** Closed schemas construct new DTOs; no input spreading, arbitrary JSON, or unsafe legacy fallback. */
function object<T extends object>(shape: { [K in keyof T]: Parser<T[K]> }): Parser<T> {
  return (v) => {
    if (v === null || typeof v !== 'object' || Array.isArray(v) || (Object.getPrototypeOf(v) !== Object.prototype && Object.getPrototypeOf(v) !== null)) return refuse()
    const input = v as Record<string, unknown>
    if (Reflect.ownKeys(input).some((key) => typeof key !== 'string' || !Object.hasOwn(shape, key))) return refuse()
    const result = {} as T
    for (const key of Object.keys(shape) as (keyof T)[]) {
      const descriptor = Object.getOwnPropertyDescriptor(input, key)
      if (descriptor && !('value' in descriptor)) return refuse()
      result[key] = shape[key](descriptor?.value)
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
const member = object<RecordingMember>({ source_value_type: nullable(dataType), circuit_id: nullable(id), device_id: id, point_id: id, tag_id: id, column: id })
const profile = object<RecordingProfile>({ kind: enumeration('instantaneous', 'cumulative'), interval: numeric(true, 1, 600), unit: enumeration('second', 'minute') })
const groupFacts = object<RecordingGroupFacts>({
  recording_profile: nullable(profile), allowed_lateness_seconds: nullable(count),
  ends_at: nullable(timestamp), binding, effective_at: nullable(timestamp), enabled: bool, destination, members: array(member),
  layout: enumeration('managed_wide_v1', 'managed_long_v1', 'managed_wide_profile_v1', 'unavailable'),
  storage_mode: enumeration('plain', 'timescale', 'native_monthly', 'unavailable'), interval_seconds: nullable(numeric(true, 1)),
  reduction: enumeration('latest', 'mean', 'cumulative_delta', 'unavailable'), write_mode: enumeration('append', 'latest', 'unavailable'),
  incomplete_policy: enumeration('skip_row', 'partial', 'unavailable'),
})
const stageCounts = object<DeliveryStageCounts>({ collecting: count, queued: count, retrying: count, blocked: count, quarantined: count, unknown: count, sql_committed: count, skipped: count })
const receipt = object<RuntimeReceipt>({
  group_revision: id, connector_revision: id, record_id: id, effect_key: id, payload_digest: id, committed_at: timestamp, provenance: enumeration('runtime'),
})
const sourceResult = object<SourceResult>({
  state: enumeration('available', 'missing', 'unknown', 'partial', 'unavailable'), window_start: nullable(timestamp), window_end: nullable(timestamp),
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
      if ((r.delta_known === false && r.state !== 'unknown') || (r.coverage === 'partial' && r.state !== 'partial')) refuse()
      if (r.receipt && (r.window_start === null || r.state === 'unavailable' || !c.last_commit)) refuse()
    }
    if (s.last_committed && !s.last_committed.receipt) refuse()
  }
}
function checkSummaryProfile(f: RecordingGroupFacts) {
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
  if (g.state === 'effective' && (!g.effective_revision || !g.as_of || !observedAt || g.proof !== 'applied-snapshot-and-layout' || !g.effective || !revision || !g.effective.effective_at || g.effective.layout === 'unavailable' || g.effective.interval_seconds === null)) {
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
