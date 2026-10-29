import { describe, expect, it } from 'vitest'
import { parseRuntimeRuleEvidence, parseEntityLocator, parseEvidenceStage, parseEvidenceReason, actionForEvidenceReason } from '../../../src/services/runtimeRuleEvidenceParser'
import type { RuleEvidenceView, EntityLocator, EvidenceReason, EvidenceAction } from '../../../src/types/runtimeRuleEvidence'

// Literal test evidence only, identical to the Go contract test. Not a production data source.
const contractLiteral = `{
 "contract_version":1,"workspace_id":"W","workspace_revision":"W1","device_id":"D",
 "observed_at":"2026-10-27T10:00:00Z","freshness":"fresh",
 "source":{"rule_id":"R","point_id":"P","tag_id":"T","mapping_id":"M","saved_revision":"R2","effective_revision":"R1","state":"saved-only","proof":"unavailable",
 "saved":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"R2"},"mapping_revision":"M2","enabled":true,"address":"40001","protocol_function":"FC03","read_type":"int16","result_type":"float64","polling":{"group_id":"poll","interval_ms":1000,"source":"polling-group"},"transforms":[{"kind":"scale","order":null,"multiplier":0.1,"offset":0,"target_type":null}]},"effective":null},
 "groups":[{"group_id":"G","saved_revision":"B","applied_effective_at":null,"applied_revision":"A","effective_revision":"A","as_of":"2026-10-27T10:00:00Z","state":"effective","proof":"applied-snapshot-and-layout",
 "saved":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"B"},"recording_profile":null,"allowed_lateness_seconds":null,"ends_at":null,"effective_at":null,"enabled":false,"destination":{"connector_id":"C2","connector_revision":"8","alias":"Meter archive","kind":"postgres","schema":"public","table":"meter_new"},"members":[{"device_id":"D","point_id":"P","tag_id":"T","column":"value","source_value_type":null,"circuit_id":null},{"device_id":"E","point_id":"Q","tag_id":"U","column":"other","source_value_type":null,"circuit_id":null}],"layout":"managed_wide_profile_v1","storage_mode":"plain","interval_seconds":30,"reduction":"mean","write_mode":"append","incomplete_policy":"partial"},
 "effective":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"A"},"recording_profile":null,"allowed_lateness_seconds":null,"ends_at":null,"effective_at":"2026-10-27T09:00:00Z","enabled":true,"destination":{"connector_id":"C1","connector_revision":"7","alias":"Meter archive","kind":"postgres","schema":"public","table":"meter_old"},"members":[{"device_id":"D","point_id":"P","tag_id":"T","column":"value","source_value_type":null,"circuit_id":null},{"device_id":"E","point_id":"Q","tag_id":"U","column":"other","source_value_type":null,"circuit_id":null}],"layout":"managed_wide_profile_v1","storage_mode":"plain","interval_seconds":15,"reduction":"mean","write_mode":"append","incomplete_policy":"partial"},
 "delivery":{"oldest_pending_scope":"group-all-revisions","oldest_pending_seconds":null,"current":{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"A"},"intake":null,"totals_scope":"group","shared":true,"stages":{"collecting":1,"queued":3,"retrying":0,"blocked":0,"quarantined":0,"unknown":0,"sql_committed":0,"skipped":0},"last_commit":null,"sources":[],"source_state":"unavailable","next_due_at":null,"next_retry_at":null,"attempt_duration_ms":null,"quota":null,"storage":null},"backlog":[{"binding":{"workspace_id":"W","workspace_revision":"W1","entity_revision":"A0"},"destination":{"connector_id":"C1","connector_revision":"7","alias":null,"kind":"postgres","schema":"public","table":"old_table"},"pending":3,"reasons":["counter-replay-proof-missing"]}],"counter_replay_quarantines":null,"test_write":null}}],
 "locator":{"step":4,"workspace_id":"W","device_id":"D","rule_id":null,"mapping_id":null,"group_id":"G","group_revision":"A0"}
}`
const scope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'D' }
function fixture(): RuleEvidenceView { return JSON.parse(contractLiteral) as RuleEvidenceView }
function parsed(value: unknown = fixture()) {
  const result = parseRuntimeRuleEvidence(value, scope)
  expect(result.ok).toBe(true)
  if (!result.ok) throw new Error(result.reason)
  return result.value
}

describe('named literal contract cases', () => {
  it('saved draft B differs from applied A', () => {
    const v = parsed(); const g = v.groups[0]
    expect([g.effective?.destination.table, g.effective?.interval_seconds]).toEqual(['meter_old', 15])
    expect([g.saved?.destination.table, g.saved?.interval_seconds]).toEqual(['meter_new', 30])
    expect(v.source?.saved?.polling).toEqual({ group_id: 'poll', interval_ms: 1000, source: 'polling-group' })
    expect(v.source?.saved?.transforms[0].multiplier).toBe(0.1)
    expect(v.source?.state).toBe('saved-only'); expect(v.source?.effective).toBeNull()
  })
  it('queued without receipt not committed', () => {
    const d = parsed().groups[0].delivery?.current
    expect(d?.stages.queued).toBe(3); expect(d?.stages.sql_committed).toBe(0)
    expect(d?.last_commit).toBeNull(); expect(d?.shared).toBe(true)
    expect(d?.totals_scope).toBe('group'); expect(d?.source_state).toBe('unavailable')
  })
  it('historical backlog retains old destination revision', () => {
    const b = parsed().groups[0].delivery?.backlog[0]
    expect([b?.binding.entity_revision, b?.destination.connector_id, b?.destination.connector_revision, b?.destination.table, b?.pending]).toEqual(['A0', 'C1', '7', 'old_table', 3])
  })
  it('unknown timing stays null', () => {
    const d = parsed().groups[0].delivery?.current
    expect([d?.next_due_at, d?.next_retry_at, d?.attempt_duration_ms]).toEqual([null, null, null])
    const raw = JSON.stringify(d)
    for (const key of ['next_due_at', 'next_retry_at', 'attempt_duration_ms']) expect(raw).toContain(`"${key}":null`)
  })
  it('missing effective proof unavailable', () => {
    const v = fixture(); v.source!.state = 'effective'
    expect(parsed(v).source).toMatchObject({ state: 'unavailable', effective: null })
    expect(parsed(v).source?.saved?.transforms[0].multiplier).toBe(0.1)
    const g = fixture(); g.groups[0].proof = 'unavailable'
    const result = parsed(g).groups[0]
    expect(result.state).toBe('unavailable'); expect(result.effective).toBeNull()
    expect(result.delivery?.current).toBeNull(); expect(result.delivery?.backlog[0].pending).toBe(3)
  })
  it('foreign locator mixed scope refused', () => {
    const cases: [RuleEvidenceView, 'foreign-scope' | 'stale-reference'][] = []
    const foreign = fixture(); foreign.locator!.workspace_id = 'foreign'; cases.push([foreign, 'foreign-scope'])
    const mixed = fixture(); mixed.groups[0].effective!.binding.workspace_revision = 'W2'; cases.push([mixed, 'stale-reference'])
    const rev = fixture(); rev.groups[0].delivery!.current!.binding.entity_revision = 'B'; cases.push([rev, 'stale-reference'])
    const device = fixture(); device.device_id = 'E'; cases.push([device, 'foreign-scope'])
    for (const [input, reason] of cases) expect(parseRuntimeRuleEvidence(input, scope)).toEqual({ ok: false, reason, action: reason === 'foreign-scope' ? 'open-workspace-list' : 'refresh' })
  })
  it('hostile secret fields rejected', () => {
    for (const injected of [{ host: 'secret-host' }, { password: 'secret-password' }, { dsn: 'secret-dsn' }, { path: '/private/db' }, { raw_error: 'driver-secret' }, { params: { nested: { token: 'secret-token' } } }]) {
      const v = fixture(); Object.assign(v.groups[0].saved!.destination, injected)
      expect(parseRuntimeRuleEvidence(v, scope)).toEqual({ ok: false, reason: 'unavailable', action: 'refresh' })
    }
  })
  it('unknown enums safe', () => {
    for (const [from, to] of [['"freshness":"fresh"', '"freshness":"future"'], ['"state":"effective"', '"state":"future"'], ['"kind":"scale"', '"kind":"formula"'], ['"reduction":"mean"', '"reduction":"magic"'], ['"counter-replay-proof-missing"', '"raw-error-secret"']]) {
      expect(parseRuntimeRuleEvidence(JSON.parse(contractLiteral.replace(from, to)), scope)).toEqual({ ok: false, reason: 'unavailable', action: 'refresh' })
    }
  })
  it('runtime commit separate from test write readback cleanup', () => {
    const v = fixture(); v.groups[0].delivery!.test_write = { operation_id: 'test1', provenance: 'test-write', written_at: '2026-10-27T10:00:00Z', readback_at: '2026-10-27T10:00:01Z', cleanup_at: '2026-10-27T10:00:02Z' }
    const d = parsed(v).groups[0].delivery
    expect(d?.current?.stages.sql_committed).toBe(0); expect(d?.current?.last_commit).toBeNull()
    expect(d?.test_write?.cleanup_at).toBe('2026-10-27T10:00:02Z')
  })
  it('historical locator is selection only', () => {
    expect(parsed().locator).toEqual({ step: 4, workspace_id: 'W', device_id: 'D', rule_id: null, mapping_id: null, group_id: 'G', group_revision: 'A0' })
  })
})

it('unproven effective transform unavailable', () => {
  const v = fixture(); const s = v.source!
  s.effective = structuredClone(s.saved); s.effective!.binding.entity_revision = 'R1'
  s.effective!.transforms = [{ kind: 'unsupported', order: null, multiplier: null, offset: null, target_type: null }]
  s.state = 'effective'; s.proof = 'immutable-source-facts'
  expect(parsed(v).source).toMatchObject({ state: 'unavailable', effective: null })
})

const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' }
it.each([
  ['missing-required', '"contract_version":1,', ''], ['null-primitive', '"enabled":true', '"enabled":null'],
  ['missing-boolean', '"enabled":false,', ''], ['missing-array', '"groups":[', '"unknown_groups":['],
  ['unknown-root', '"contract_version":1', '"contract_version":1,"password":"secret"'],
  ['wrong-case', '"workspace_id":"W"', '"Workspace_id":"W"'], ['negative-count', '"queued":3', '"queued":-1'],
  ['fractional-count', '"queued":3', '"queued":0.5'], ['unsafe-integer', '"queued":3', '"queued":9007199254740992'],
  ['zero-poll', '"interval_ms":1000', '"interval_ms":0'], ['unknown-version', '"contract_version":1', '"contract_version":2'],
  ['null-array', '"sources":[]', '"sources":null'], ['wrong-scalar', '"enabled":true', '"enabled":"true"'],
  ['bad-timestamp', '2026-10-27T10:00:00Z', '2026-02-30T10:00:00Z'],
  ['non-utc', '2026-10-27T10:00:00Z', '2026-10-27T10:00:00+01:00'],
  ['alias-endpoint', 'Meter archive', 'postgres://host/private'], ['private-table', 'meter_new', '/private/db.sqlite'],
  ['sql-table', 'meter_new', 'SELECT * FROM secret'], ['unknown-datatype', '"int16"', '"number"'],
])('rejects %s', (_name, from, to) => {
  expect(parseRuntimeRuleEvidence(JSON.parse(contractLiteral.replace(from, to)), scope)).toEqual(unavailable)
})
it('missing nullable evidence normalizes to null and no input gets mutated', () => {
  const input = JSON.parse(contractLiteral.replace('"next_retry_at":null,', '')) as RuleEvidenceView
  expect(parsed(input).groups[0].delivery?.current?.next_retry_at).toBeNull()
  const v = fixture(); v.groups[0].proof = 'unavailable'
  const before = structuredClone(v); parsed(v); expect(v).toEqual(before)
})
it('rejects non-JSON values accessors prototypes and non-finite numbers safely', () => {
  for (const input of [undefined, null, true, [], new Date(), Object.create({ workspace_id: 'W' })]) expect(parseRuntimeRuleEvidence(input, scope)).toEqual(unavailable)
  for (const number of [NaN, Infinity, -Infinity]) {
    const v = fixture(); v.groups[0].delivery!.current!.attempt_duration_ms = number
    expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
  }
  const v = fixture(); Object.defineProperty(v, 'source', { get() { throw new Error('secret') } })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
  const symbol = fixture(); Object.defineProperty(symbol, Symbol('secret'), { value: 'private' })
  expect(parseRuntimeRuleEvidence(symbol, scope)).toEqual(unavailable)
  const proxy = new Proxy({}, { getPrototypeOf() { throw new Error('secret') } })
  expect(parseRuntimeRuleEvidence(proxy, scope)).toEqual(unavailable)
})
it('proven source preserves independent revisions and structured transform parameters', () => {
  const v = fixture(); const s = v.source!; s.effective = structuredClone(s.saved)!
  s.effective.binding.entity_revision = 'R1'; s.state = 'effective'; s.proof = 'immutable-source-facts'
  s.effective.transforms = [
    { kind: 'byte-order', order: 'CDAB', multiplier: null, offset: null, target_type: null },
    { kind: 'scale', order: null, multiplier: -2, offset: 3, target_type: null },
    { kind: 'cast', order: null, multiplier: null, offset: null, target_type: 'float64' },
  ]
  const got = parsed(v).source
  expect(got?.state).toBe('effective'); expect(got?.effective?.transforms[1].multiplier).toBe(-2)
  expect(got?.saved?.binding.entity_revision).toBe('R2')
  for (const key of ['mapping_revision', 'address', 'read_type', 'result_type', 'polling'] as const) {
    const copy = structuredClone(v); copy.source!.effective![key] = null
    expect(parsed(copy).source?.state).toBe('unavailable')
  }
  for (const step of [
    { kind: 'scale', order: null, multiplier: null, offset: null, target_type: null },
    { kind: 'cast', order: null, multiplier: null, offset: 1, target_type: 'float64' },
    { kind: 'byte-order', order: null, multiplier: null, offset: null, target_type: null },
    { kind: 'unsupported', order: null, multiplier: null, offset: 1, target_type: null },
  ]) {
    const copy = structuredClone(v); Object.assign(copy.source!.effective!.transforms[0], step)
    expect(parseRuntimeRuleEvidence(copy, scope)).toEqual(unavailable)
  }
})
function committedFixture() {
  const v = fixture(); const c = v.groups[0].delivery!.current!
  c.stages.sql_committed = 1
  c.last_commit = { group_revision: 'A', connector_revision: '7', record_id: 'record1', effect_key: 'effect1', payload_digest: 'digest1', committed_at: '2026-10-27T10:00:00Z', provenance: 'runtime' }
  c.source_state = 'available'
  c.sources = [{ device_id: 'D', point_id: 'P', tag_id: 'T',
    latest: { state: 'unknown', window_start: '2026-10-27T10:00:00Z', window_end: '2026-10-27T10:00:15Z', delta_known: false, coverage: null, receipt: null },
    last_committed: { state: 'partial', window_start: '2026-10-27T09:59:45Z', window_end: '2026-10-27T10:00:00Z', delta_known: null, coverage: 'partial', receipt: { ...c.last_commit } },
  }]
  return v
}
it('runtime receipt and latest unknown delta / old partial mean stay independent', () => {
  const s = parsed(committedFixture()).groups[0].delivery!.current!.sources[0]
  expect([s.latest?.state, s.latest?.delta_known, s.latest?.receipt]).toEqual(['unknown', false, null])
  expect([s.last_committed?.state, s.last_committed?.receipt?.record_id]).toEqual(['partial', 'record1'])
})
const receiptFailures: [string, (v: RuleEvidenceView) => void, 'unavailable' | 'stale-reference' | 'foreign-scope'][] = [
  ['receipt-group', (v) => { v.groups[0].delivery!.current!.last_commit!.group_revision = 'B' }, 'stale-reference'],
  ['receipt-connector', (v) => { v.groups[0].delivery!.current!.last_commit!.connector_revision = '8' }, 'stale-reference'],
  ['test-as-runtime', (v) => { Object.assign(v.groups[0].delivery!.current!.last_commit!, { provenance: 'test-write' }) }, 'unavailable'],
  ['missing-receipt', (v) => { v.groups[0].delivery!.current!.last_commit = null }, 'unavailable'],
  ['receipt-with-zero', (v) => { v.groups[0].delivery!.current!.stages.sql_committed = 0 }, 'unavailable'],
  ['shared-hidden', (v) => { v.groups[0].delivery!.current!.shared = false }, 'unavailable'],
  ['foreign-source', (v) => { v.groups[0].delivery!.current!.sources[0].device_id = 'foreign' }, 'foreign-scope'],
  ['duplicate-source', (v) => { const c = v.groups[0].delivery!.current!; c.sources.push(c.sources[0]) }, 'unavailable'],
  ['unavailable-sources', (v) => { v.groups[0].delivery!.current!.source_state = 'unavailable' }, 'unavailable'],
  ['half-window', (v) => { v.groups[0].delivery!.current!.sources[0].latest!.window_end = null }, 'unavailable'],
  ['reverse-window', (v) => { v.groups[0].delivery!.current!.sources[0].latest!.window_end = '2026-10-27T09:00:00Z' }, 'unavailable'],
  ['unknown-as-zero', (v) => { v.groups[0].delivery!.current!.sources[0].latest!.state = 'available' }, 'unavailable'],
  ['partial-as-complete', (v) => { v.groups[0].delivery!.current!.sources[0].last_committed!.state = 'available' }, 'unavailable'],
  ['unmatched-source-receipt', (v) => { v.groups[0].delivery!.current!.sources[0].last_committed!.receipt = null }, 'unavailable'],
  ['unavailable-committed', (v) => { const r = v.groups[0].delivery!.current!.sources[0].last_committed!; r.coverage = null; r.state = 'unavailable' }, 'unavailable'],
]
it.each(receiptFailures)('refuses %s', (_name, mutate, reason) => {
  const v = committedFixture(); mutate(v)
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual({ ok: false, reason, action: reason === 'foreign-scope' ? 'open-workspace-list' : 'refresh' })
})
it('locator steps require applicable IDs only; action mapping is read-only', () => {
  const base = { workspace_id: 'W', device_id: 'D', rule_id: null, mapping_id: null, group_id: null, group_revision: null }
  const locators: EntityLocator[] = [{ ...base, step: 1 }, { ...base, step: 2, rule_id: 'R' }, { ...base, step: 3, mapping_id: 'M' }, { ...base, step: 4, device_id: null, group_id: 'G', group_revision: 'A0' }]
  for (const l of locators) {
    expect(parseEntityLocator(l, scope)).toEqual({ ok: true, value: l })
    const v = fixture(); v.locator = l; expect(parsed(v).locator).toEqual(l)
    expect(parseEntityLocator({ ...l, ...(l.step < 4 ? { device_id: null } : { group_id: null }) }, scope)).toEqual(unavailable)
    expect(parseEntityLocator({ ...l, ...(l.step === 3 ? { rule_id: 'other' } : { mapping_id: 'other' }) }, scope)).toEqual(unavailable)
  }
  const mappings: [EvidenceReason, EvidenceAction][] = [['none', 'none'], ['foreign-scope', 'open-workspace-list'], ['deleted', 'open-workspace-list'], ['unsupported', 'inspect-settings'], ['effective-proof-missing', 'refresh'], ['counter-replay-proof-missing', 'inspect-immutable-evidence'], ['counter-replay-proof-conflict', 'inspect-immutable-evidence'], ['counter-replay-proof-unsupported', 'inspect-immutable-evidence'], ['recording-ownership-unavailable', 'inspect-immutable-evidence'], ['recording-ownership-conflict', 'inspect-immutable-evidence'], ['recording-handoff-proof-conflict', 'inspect-immutable-evidence'], ['recording-cutover-finalized-overlap', 'inspect-immutable-evidence'], ['stale-reference', 'refresh'], ['unavailable', 'refresh']]
  for (const [reason, action] of mappings) { expect(parseEvidenceReason(reason)).toBe(reason); expect(actionForEvidenceReason(reason)).toBe(action) }
  for (const stage of ['unavailable', 'acquisition', 'accepted', 'collecting', 'queued', 'retrying', 'blocked', 'quarantined', 'unknown', 'skipped', 'sql-committed', 'readback-verified', 'cleanup-completed', 'transport-event']) expect(parseEvidenceStage(stage)).toBe(stage)
  expect(parseEvidenceReason('secret-driver-text')).toBe('unavailable'); expect(parseEvidenceStage('secret-driver-text')).toBe('unavailable')
})
it('freshness and real nullable timing/quota/storage evidence remain explicit', () => {
  const v = fixture(); v.freshness = 'stale'; expect(parsed(v).freshness).toBe('stale')
  v.observed_at = null; expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
  v.freshness = 'unavailable'; v.source = null; v.groups = []; v.locator = null; expect(parsed(v).observed_at).toBeNull()
  const real = fixture(); const c = real.groups[0].delivery!.current!
  c.next_due_at = '2026-10-27T10:00:15Z'; c.next_retry_at = '2026-10-27T10:00:16Z'; c.attempt_duration_ms = 12.5
  c.quota = { scope: 'global', state: 'hard_limit', used_bytes: 95, max_bytes: 100, intake_refused: true }
  c.storage = { generation: 2, history_may_be_lost: true, reconstructions: 1, last_reconstructed_at: '2026-10-27T09:00:00Z', repair_state: 'blocked', events: [{ kind: 'repair_blocked', at: '2026-10-27T10:00:00Z' }] }
  expect(parsed(real).groups[0].delivery?.current?.attempt_duration_ms).toBe(12.5)
  expect(parsed(real).groups[0].delivery?.current?.quota?.state).toBe('hard_limit')
})
it('incomplete proof cannot hide receipt revision mismatch', () => {
  const v = committedFixture(); v.groups[0].proof = 'unavailable'; v.groups[0].delivery!.current!.last_commit!.group_revision = 'B'
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' })
})
it('repeated receipt digest conflict is refused but different older receipt is independent', () => {
  const v = committedFixture(); v.groups[0].delivery!.current!.sources[0].last_committed!.receipt!.payload_digest = 'different'
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' })
  const older = committedFixture(); Object.assign(older.groups[0].delivery!.current!.sources[0].last_committed!.receipt!, { record_id: 'older-record', effect_key: 'older-effect', payload_digest: 'older-digest', committed_at: '2026-10-27T09:59:45Z' })
  expect(parsed(older).groups[0].delivery?.current?.sources[0].last_committed?.receipt?.record_id).toBe('older-record')
})
it('draft removal retains applied device without merging memberships', () => {
  const v = fixture(); v.groups[0].saved!.members = v.groups[0].saved!.members.slice(1)
  const g = parsed(v).groups[0]; expect(g.effective?.members).toHaveLength(2); expect(g.saved?.members).toHaveLength(1)
})
it('nanosecond window is preserved consistently with Go', () => {
  const v = committedFixture(); const r = v.groups[0].delivery!.current!.sources[0].latest!
  r.window_start = '2026-10-27T10:00:00.000000001Z'; r.window_end = '2026-10-27T10:00:00.000000002Z'
  expect(parsed(v).groups[0].delivery?.current?.sources[0].latest?.window_end).toBe('2026-10-27T10:00:00.000000002Z')
})
it('cleanup does not imply readback', () => {
  const v = fixture(); v.groups[0].delivery!.test_write = { operation_id: 'test2', provenance: 'test-write', written_at: '2026-10-27T10:00:00Z', readback_at: null, cleanup_at: '2026-10-27T10:00:02Z' }
  expect(parsed(v).groups[0].delivery?.test_write?.readback_at).toBeNull()
})

it('exact literal projection parity and equal/reversed/offset nanosecond windows', () => {
  expect(parsed()).toEqual(fixture())
  const v = committedFixture(); const r = v.groups[0].delivery!.current!.sources[0].latest!
  r.window_start = '2026-10-27T10:00:00.000000001Z'
  for (const end of ['2026-10-27T10:00:00.000000001Z', '2026-10-27T10:00:00Z', '2026-10-27T10:00:00.000000002+00:00']) {
    r.window_end = end; expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
  }
})
it('current intake reason is revision-bound and independent of historical backlog', () => {
  const v = fixture(); Object.assign(v.groups[0].delivery!.current!, { intake: { group_revision: 'A', state: 'blocked', reason: 'counter-replay-proof-missing' } })
  expect(parsed(v).groups[0].delivery?.current).toMatchObject({ intake: { group_revision: 'A', state: 'blocked', reason: 'counter-replay-proof-missing' } })
  Object.assign(v.groups[0].delivery!.current!, { intake: { group_revision: 'A', state: 'active', reason: null } })
  expect(parsed(v).groups[0].delivery?.current).toMatchObject({ intake: { group_revision: 'A', state: 'active', reason: null } })
  Object.assign(v.groups[0].delivery!.current!, { intake: { group_revision: 'A', state: 'active', reason: 'secret-error' } })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
it('safe labels and identifiers reject trailing newlines consistently with Go', () => {
  const v = fixture(); v.groups[0].saved!.destination.alias = 'Meter archive\n'
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
  const id = fixture(); id.groups[0].saved!.destination.table = 'meter_new\n'
  expect(parseRuntimeRuleEvidence(id, scope)).toEqual(unavailable)
})
const stateFailures: [string, (v: RuleEvidenceView) => void, 'unavailable' | 'stale-reference' | 'foreign-scope'][] = [
  ['source-saved-missing', (v) => { v.source!.saved = null }, 'unavailable'],
  ['source-unavailable-with-effective', (v) => { v.source!.state = 'unavailable'; v.source!.effective = v.source!.saved; v.source!.effective_revision = v.source!.saved_revision }, 'unavailable'],
  ['source-saved-with-effective', (v) => { v.source!.effective = v.source!.saved; v.source!.effective_revision = v.source!.saved_revision }, 'unavailable'],
  ['duplicate-group', (v) => { v.groups.push(v.groups[0]) }, 'unavailable'],
  ['duplicate-member', (v) => { v.groups[0].saved!.members.push(v.groups[0].saved!.members[0]) }, 'unavailable'],
  ['foreign-all-members', (v) => { v.groups[0].saved!.members = v.groups[0].saved!.members.slice(1); v.groups[0].effective!.members = v.groups[0].effective!.members.slice(1) }, 'foreign-scope'],
  ['saved-group-with-effective', (v) => { v.groups[0].state = 'saved-only' }, 'unavailable'],
  ['unavailable-group-with-effective', (v) => { v.groups[0].state = 'unavailable' }, 'unavailable'],
  ['backlog-current-revision', (v) => { v.groups[0].delivery!.backlog[0].binding.entity_revision = 'A' }, 'stale-reference'],
  ['backlog-mixed-workspace', (v) => { v.groups[0].delivery!.backlog[0].binding.workspace_id = 'foreign' }, 'foreign-scope'],
  ['locator-stale-group', (v) => { v.locator!.group_id = 'other' }, 'stale-reference'],
  ['locator-stale-rule', (v) => { Object.assign(v.locator!, { step: 2, group_id: null, group_revision: null, rule_id: 'other' }) }, 'stale-reference'],
  ['locator-stale-mapping', (v) => { Object.assign(v.locator!, { step: 3, group_id: null, group_revision: null, mapping_id: 'other' }) }, 'stale-reference'],
  ['locator-foreign-device', (v) => { v.locator!.device_id = 'other' }, 'foreign-scope'],
  ['locator-invalid-step', (v) => { Object.assign(v.locator!, { step: 0 }) }, 'unavailable'],
]
it.each(stateFailures)('state/scope %s', (_name, mutate, reason) => {
  const v = fixture(); mutate(v)
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual({ ok: false, reason, action: reason === 'foreign-scope' ? 'open-workspace-list' : 'refresh' })
})
it('optional saved/effective states degrade rather than borrow facts', () => {
  for (const mutate of [
    (v: RuleEvidenceView) => { v.groups[0].effective!.effective_at = null },
    (v: RuleEvidenceView) => { v.groups[0].effective!.layout = 'unavailable' },
    (v: RuleEvidenceView) => { v.groups[0].effective!.interval_seconds = null },
    (v: RuleEvidenceView) => { v.groups[0].effective = null },
  ]) { const v = fixture(); mutate(v); expect(parsed(v).groups[0].state).toBe('unavailable') }
  const saved = fixture(); const g = saved.groups[0]; g.state = 'saved-only'; g.effective = null; g.proof = 'unavailable'; g.delivery = null
  expect(parsed(saved).groups[0].state).toBe('saved-only')
  g.state = 'unavailable'; g.saved = null; g.saved_revision = null; expect(parsed(saved).groups[0].state).toBe('unavailable')
  const unsupported = fixture(); unsupported.source!.saved!.transforms = [{ kind: 'unsupported', order: null, multiplier: null, offset: null, target_type: null }]
  expect(parsed(unsupported).source?.state).toBe('saved-only')
  const single = fixture(); single.groups[0].saved!.members = single.groups[0].saved!.members.slice(0, 1); single.groups[0].effective!.members = single.groups[0].effective!.members.slice(0, 1); single.groups[0].delivery!.current!.shared = false
  expect(parsed(single).groups[0].delivery?.current?.shared).toBe(false)
  single.groups[0].delivery!.current = null; expect(parsed(single).groups[0].delivery?.current).toBeNull()
  expect(parseRuntimeRuleEvidence(fixture(), { ...scope, workspace_id: '' })).toEqual(unavailable)
  expect(parseRuntimeRuleEvidence(fixture(), { ...scope, workspace_revision: 'W2' })).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' })
})
it('authority labels distinguish revision intake from aggregate quota and pending age', () => {
  const v = fixture(); const d = v.groups[0].delivery!; const c = d.current!
  Object.assign(c, { intake: { group_revision: 'A', state: 'active', reason: null }, quota: { scope: 'global', state: 'hard_limit', used_bytes: 95, max_bytes: 100, intake_refused: true } })
  Reflect.deleteProperty(c, 'oldest_pending_seconds')
  Object.assign(d, { oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 120 })
  const got = parsed(v).groups[0].delivery!
  expect(got).toMatchObject({ oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 120 })
  expect(got.current?.quota).toMatchObject({ scope: 'global', used_bytes: 95 })
  expect(got.current).not.toHaveProperty('oldest_pending_seconds')
  Object.assign(c.quota!, { scope: 'group' }); expect(parsed(v).groups[0].delivery?.current?.quota?.scope).toBe('group')
  Object.assign(c.quota!, { scope: 'future' }); expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
  Object.assign(c.quota!, { scope: 'global' })
  Object.assign(c.intake!, { group_revision: 'B' })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' })
  Reflect.deleteProperty(c.intake!, 'group_revision')
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})

it('legacy requested revision cannot claim current without as-of proof', () => {
  const v = JSON.parse(contractLiteral.replace('"effective_revision":"A","as_of":"2026-10-27T10:00:00Z",', ''))
  expect(parsed(v).groups[0]).toMatchObject({ state: 'unavailable', effective: null, delivery: { current: null } })
})
it('scheduled B keeps effective A and accepted input quarantine separate', () => {
  const v = fixture()
  Object.assign(v.groups[0], { applied_revision: 'B', effective_revision: 'A', as_of: '2026-10-27T10:00:00Z' })
  Object.assign(v.groups[0].delivery!, { counter_replay_quarantines: [{ workspace_id: 'W', group_id: 'G', group_revision: 'A0', pending_samples: 7, reason: 'counter-replay-proof-missing', action: 'preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance' }] })
  const g = parsed(v).groups[0]
  expect(g.effective?.binding.entity_revision).toBe('A')
  expect(g.delivery?.current?.stages.sql_committed).toBe(0)
  expect(g.delivery?.backlog[0].pending).toBe(3)
})

it('missing effective revision or server as-of safely degrades while retaining backlog', () => {
  for (const key of ['effective_revision', 'as_of'] as const) {
    const v = fixture(); v.groups[0][key] = null
    expect(parsed(v).groups[0]).toMatchObject({ state: 'unavailable', effective: null, delivery: { current: null, backlog: [{ pending: 3 }] } })
  }
})
it('mixed as-of, future effective nanosecond and mixed effective revision are stale', () => {
  for (const at of ['2026-10-27T10:00:00.000000001Z', '2026-10-27T09:59:59Z']) {
    const v = fixture(); v.groups[0].as_of = at
    expect(parseRuntimeRuleEvidence(v, scope)).toMatchObject({ ok: false, reason: 'stale-reference' })
  }
  const v = fixture(); v.groups[0].effective!.effective_at = '2026-10-27T10:00:00.000000001Z'
  expect(parseRuntimeRuleEvidence(v, scope)).toMatchObject({ ok: false, reason: 'stale-reference' })
  const mixed = fixture(); mixed.groups[0].effective_revision = 'B'
  expect(parseRuntimeRuleEvidence(mixed, scope)).toMatchObject({ ok: false, reason: 'stale-reference' })
})
it('retained ended or disabled snapshots never infer intake and refuse active contradiction', () => {
  for (const ended of [true, false]) {
    const v = fixture(); const g = v.groups[0]
    if (ended) g.effective!.ends_at = v.observed_at; else g.effective!.enabled = false
    expect(parsed(v).groups[0].delivery?.current?.intake).toBeNull()
    g.delivery!.current!.intake = { group_revision: 'A', state: 'active', reason: null }
    expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
    g.delivery!.current!.intake.state = 'not_running'
    expect(parsed(v).groups[0].delivery?.current?.intake?.state).toBe('not_running')
  }
})
const quarantineFixture = () => ({ workspace_id: 'W', group_id: 'G', group_revision: 'A0', pending_samples: 7, reason: 'counter-replay-proof-missing', action: 'preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance' })
it('quarantine null is unavailable, empty is observed, typed reasons retain exact input counts', () => {
  expect(parsed().groups[0].delivery?.counter_replay_quarantines).toBeNull()
  const v = fixture(); Object.assign(v.groups[0].delivery!, { counter_replay_quarantines: [] })
  expect(parsed(v).groups[0].delivery?.counter_replay_quarantines).toEqual([])
  for (const reason of ['counter-replay-proof-missing', 'counter-replay-proof-conflict', 'counter-replay-proof-unsupported']) {
    Object.assign(v.groups[0].delivery!, { counter_replay_quarantines: [{ ...quarantineFixture(), reason }] })
    expect(parsed(v).groups[0].delivery?.counter_replay_quarantines?.[0].pending_samples).toBe(7)
    expect(parsed(v).groups[0].delivery?.current?.stages.quarantined).toBe(0)
  }
})
it.each([
  ['workspace_id', 'foreign', 'foreign-scope'], ['group_id', 'other', 'foreign-scope'],
  ['group_revision', '', 'unavailable'], ['reason', 'future', 'unavailable'], ['reason', 'recording-cutover-finalized-overlap', 'unavailable'],
  ['action', 'retry', 'unavailable'], ['pending_samples', -1, 'unavailable'], ['pending_samples', 9007199254740992, 'unavailable'],
  ['host', 'secret-host', 'unavailable'], ['database', '/private/db', 'unavailable'],
])('quarantine refuses %s = %s', (key, value, reason) => {
  const v = fixture(); Object.assign(v.groups[0].delivery!, { counter_replay_quarantines: [{ ...quarantineFixture(), [key]: value }] })
  expect(parseRuntimeRuleEvidence(v, scope)).toMatchObject({ ok: false, reason })
})
it('saved C requested B effective A stays separate then boundary represents B', () => {
  const v = fixture(); const g = v.groups[0]
  g.saved_revision = 'C'; g.saved!.binding.entity_revision = 'C'; g.applied_revision = 'B'; g.applied_effective_at = '2026-10-27T10:00:01Z'
  expect(parsed(v).groups[0]).toMatchObject({ saved_revision: 'C', applied_revision: 'B', effective_revision: 'A', effective: { destination: { table: 'meter_old' } } })
  v.observed_at = g.applied_effective_at; g.as_of = v.observed_at; g.effective_revision = 'B'
  g.effective!.binding.entity_revision = 'B'; g.effective!.effective_at = g.applied_effective_at; g.delivery!.current!.binding.entity_revision = 'B'
  expect(parsed(v).groups[0].effective_revision).toBe('B')
})
it.each(['pending_samples', 'reason', 'action', 'group_revision'])('quarantine missing %s never infers proof or zero', (key) => {
  const q: Record<string, unknown> = quarantineFixture(); delete q[key]
  const v = fixture(); Object.assign(v.groups[0].delivery!, { counter_replay_quarantines: [q] })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
it('equivalent nanosecond timestamp remains effective and duplicate quarantine scope is rejected', () => {
  const v = fixture(); v.groups[0].as_of = '2026-10-27T10:00:00.000000000Z'
  expect(parsed(v).groups[0].state).toBe('effective')
  Object.assign(v.groups[0].delivery!, { counter_replay_quarantines: [quarantineFixture(), quarantineFixture()] })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
it('persisted frozen C1 circuit identity round trips without column inference', () => {
  const v = fixture()
  Object.assign(v.groups[0].effective!.members[0], { circuit_id: 'C1' })
  expect(parsed(v).groups[0].effective!.members[0]).toMatchObject({ circuit_id: 'C1', column: 'value' })
})
it('missing persisted circuit stays null; effective A-only C1 never copies newer saved C', () => {
  const legacy = JSON.parse(contractLiteral.replaceAll(',"circuit_id":null', ''))
  expect(parsed(legacy).groups[0].effective!.members[0].circuit_id).toBeNull()
  const v = fixture(); const g = v.groups[0]
  g.saved_revision = 'C'; g.saved!.binding.entity_revision = 'C'; g.saved!.members = g.saved!.members.slice(1)
  g.effective!.members[0].circuit_id = 'C1'
  expect(parsed(v).groups[0].effective!.members[0].circuit_id).toBe('C1')
  expect(parsed(v).groups[0].saved!.members[0].circuit_id).toBeNull()
  g.effective!.binding.workspace_id = 'foreign'
  expect(parseRuntimeRuleEvidence(v, scope)).toMatchObject({ ok: false, reason: 'foreign-scope' })
  g.effective!.binding.workspace_id = 'W'; g.effective!.binding.entity_revision = 'C'
  expect(parseRuntimeRuleEvidence(v, scope)).toMatchObject({ ok: false, reason: 'stale-reference' })
})
it.each(['', '/private/C1', 'postgres://secret', 1, { name: 'C1' }])('refuses unsafe persisted circuit %j', (circuit_id) => {
  const v = fixture(); Object.assign(v.groups[0].effective!.members[0], { circuit_id })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
it('does not accept display circuit names as safe identity evidence', () => {
  const v = fixture(); Object.assign(v.groups[0].effective!.members[0], { circuit_id: 'C1', circuit_name: 'secret' })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
it('safe summary profile lateness and frozen source type round trip', () => {
  const v = fixture(); const f = v.groups[0].effective!
  Object.assign(f, { recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, allowed_lateness_seconds: 2 })
  Object.assign(f.members[0], { source_value_type: 'float64' })
  expect(parsed(v).groups[0].effective).toMatchObject({ recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, allowed_lateness_seconds: 2, members: [{ source_value_type: 'float64' }, {}] })
})
it.each([
  ['instantaneous', 'second', 15, 16, 'mean'], ['instantaneous', 'minute', 1, 60, 'mean'],
  ['instantaneous', 'second', 15, 15, 'cumulative_delta'], ['instantaneous', 'second', 301, 301, 'mean'],
  ['cumulative', 'minute', 1, 60, 'mean'], ['cumulative', 'second', 1, 1, 'cumulative_delta'],
])('refuses contradictory profile %s/%s/%s/%s/%s', (kind, unit, interval, interval_seconds, reduction) => {
  const v = fixture(); Object.assign(v.groups[0].effective!, { recording_profile: { kind, unit, interval }, interval_seconds, reduction })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
it.each([['instantaneous', 'second', 1, 1, 'latest'], ['instantaneous', 'second', 300, 300, 'mean'], ['cumulative', 'minute', 1, 60, 'cumulative_delta'], ['cumulative', 'minute', 600, 36000, 'cumulative_delta']])('preserves valid profile bound %s %s %s', (kind, unit, interval, interval_seconds, reduction) => {
  const v = fixture(); Object.assign(v.groups[0].effective!, { recording_profile: { kind, unit, interval }, interval_seconds, reduction })
  expect(parsed(v).groups[0].effective!.recording_profile?.interval).toBe(interval)
})
it('missing summary fields stay null and profile never synthesizes missing normalized evidence', () => {
  const v = JSON.parse(contractLiteral.replaceAll('"recording_profile":null,', '').replaceAll('"allowed_lateness_seconds":null,', '').replaceAll('"source_value_type":null,', ''))
  expect(parsed(v).groups[0].effective).toMatchObject({ recording_profile: null, allowed_lateness_seconds: null, members: [{ source_value_type: null }, { source_value_type: null }] })
  const f = fixture(); Object.assign(f.groups[0].saved!, { recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, interval_seconds: null, reduction: 'unavailable' })
  expect(parsed(f).groups[0].saved).toMatchObject({ interval_seconds: null, reduction: 'unavailable' })
})
it.each([{ recording_profile: { kind: 'future', interval: 15, unit: 'second' } }, { recording_profile: { kind: 'instantaneous', interval: 0, unit: 'second' } }, { recording_profile: { kind: 'instantaneous', interval: 15, unit: 'hour' } }, { recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', physical_unit: 'secret' } }, { allowed_lateness_seconds: -1 }, { allowed_lateness_seconds: 1.5 }])('refuses unsafe summary %j', (fields) => {
  const v = fixture(); Object.assign(v.groups[0].effective!, fields)
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
it('unknown frozen source type is not coerced into a numeric type', () => {
  const v = fixture(); Object.assign(v.groups[0].effective!.members[0], { source_value_type: 'future' })
  expect(parseRuntimeRuleEvidence(v, scope)).toEqual(unavailable)
})
