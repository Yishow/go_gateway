import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { describe, expect, it } from 'vitest';
import { projectRuntimeLogDomain } from '@/services/runtimeLogDomainProjection';
import type { RuntimeLogRecord } from '@/types/runtimeLogs';
import { buildStudioEntityLocatorTarget } from '@/services/studioEntityLocator';
import { actionForEvidenceReason } from '@/services/runtimeRuleEvidenceParser';
import { domainCases, domainInstant, domainRecord } from './domain-event-fixtures';

function rejected(record: RuntimeLogRecord) { expect(projectRuntimeLogDomain(record)).toMatchObject({ kind: 'unavailable' }); }

describe('closed runtime domain presentation', () => {
  it('keeps all 24 consumer examples aligned with the actual Go producer matrix', () => {
    const go = readFileSync(resolve(process.cwd(), '../internal/diagnostics/domain_event.go'), 'utf8');
    const actual = [...go.matchAll(/\{"(domain\.[^"]+)", "([^"]+)", "([^"]+)", ([1-4]), (true|false)\}/g)]
      .map(([, code, stage, kind, step, success]) => [code, stage, kind, Number(step), success === 'true']);
    expect(actual).toEqual(domainCases); expect(actual).toHaveLength(24);
  });
  it.each(domainCases)('projects %s using exact code and applicable identity', (code, stage, kind, step, success) => {
    const result = projectRuntimeLogDomain(domainRecord(code));
    expect(result).toMatchObject({ kind: 'domain', code, stage, evidenceKind: kind, level: success ? 'info' : 'warn',
      timestamp: domainInstant, binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'R1' }, locator: { step }, summary: null });
    if (result.kind !== 'domain') throw new Error('Expected validated domain fixture');
    expect(buildStudioEntityLocatorTarget(result.locator)).not.toBeNull();
    if (step === 4) expect(result.locator).toEqual({ step: 4, workspace_id: 'W', device_id: null, rule_id: null, mapping_id: null, group_id: 'G', group_revision: 'R1' });
  });
  it('keeps historical original rule revision separate from current-IDs locator', () => {
    const result = projectRuntimeLogDomain(domainRecord('domain.source.rejected'));
    expect(result).toMatchObject({ binding: { entity_revision: 'R1' }, locator: { rule_id: 'RULE', group_revision: null } });
    if (result.kind !== 'domain') throw new Error('Expected source');
    expect(buildStudioEntityLocatorTarget(result.locator)).toBe('/studio/v2?entry=locator&workspace_id=W&step=2&device_id=D&rule_id=RULE');
  });
  it.each(['9007199254740993', '18446744073709551615'])('preserves uint64 count %s and nanosecond bytes without numeric conversion', occurrences => {
    const first = '2026-10-29T10:00:00.123456788Z';
    expect(projectRuntimeLogDomain(domainRecord('domain.delivery.blocked', { fields: { first_seen_at: first, last_seen_at: domainInstant, occurrences } })))
      .toMatchObject({ summary: { firstSeenAt: first, lastSeenAt: domainInstant, occurrences } });
  });
  it('does not fabricate first/last/count from a transport timestamp', () => {
    expect(projectRuntimeLogDomain(domainRecord())).toMatchObject({ kind: 'domain', summary: null });
  });
  it.each(['0', '01', '-1', '1.0', '1e3', '18446744073709551616', '9'.repeat(200)])('rejects noncanonical occurrence %s', occurrences => {
    rejected(domainRecord('domain.delivery.blocked', { fields: { first_seen_at: domainInstant, last_seen_at: domainInstant, occurrences } }));
  });
  it.each<RuntimeLogRecord['fields']>([
    { first_seen_at: domainInstant }, { last_seen_at: domainInstant }, { occurrences: '1' },
    { first_seen_at: domainInstant, last_seen_at: '2026-10-29T10:00:00.123456788Z', occurrences: '1' },
  ])('rejects incomplete or reversed summary %j', fields => rejected(domainRecord('domain.delivery.blocked', { fields })));
  it.each(['2026-02-30T10:00:00Z', '2026-10-29T24:00:00Z', '2026-10-29T10:00:00+00:00',
    '2026-10-29T10:00:00Z secret', '2026-10-29T10:00:00.1234567890Z', '2026-10-29T10:00:60Z'])('rejects unsafe or invalid UTC instant %s', timestamp => {
    expect(projectRuntimeLogDomain(domainRecord('domain.delivery.blocked', { timestamp }))).toEqual({ kind: 'unavailable', timestamp: null, summary: null });
    rejected(domainRecord('domain.delivery.first-confirmed', { fields: { committed_at: timestamp } }));
    rejected(domainRecord('domain.delivery.blocked', { fields: { first_seen_at: timestamp, last_seen_at: timestamp, occurrences: '1' } }));
  });
  it('ignores candidate message and detaches all nested presentation values', () => {
    const input = domainRecord('domain.delivery.summary', { message: 'password=hunter2 postgresql://private.example SELECT * FROM x', fields: {
      first_seen_at: domainInstant, last_seen_at: domainInstant, occurrences: '2',
    } });
    const result = projectRuntimeLogDomain(input);
    expect(result.kind).toBe('domain'); expect(JSON.stringify(result)).not.toMatch(/hunter2|private\.example|SELECT|Candidate/);
    input.fields.workspace_id = 'OTHER'; input.fields.record_id = 'OTHER'; input.fields.occurrences = '9';
    expect(result).toMatchObject({ binding: { workspace_id: 'W' }, receipt: { record_id: 'REC1' }, summary: { occurrences: '2' } });
  });
  it.each(['host', 'raw_error', 'connection_config', 'password', 'value', 'unknown'])('rejects arbitrary field %s', key => {
    rejected(domainRecord('domain.delivery.blocked', { fields: { [key]: 'secret' } }));
  });
  it('rejects accessor, symbol, inherited and nested fields without invoking getters', () => {
    const accessor = domainRecord(); let reads = 0;
    Object.defineProperty(accessor.fields, 'workspace_id', { get() { reads++; throw new Error('secret'); } });
    rejected(accessor); expect(reads).toBe(0);
    const outer = domainRecord(); Object.defineProperty(outer, 'code', { get() { reads++; throw new Error('secret'); } });
    rejected(outer); expect(reads).toBe(0);
    const symbol = domainRecord(); Object.defineProperty(symbol.fields, Symbol('secret'), { value: 'secret' }); rejected(symbol);
    const inherited = domainRecord(); Object.setPrototypeOf(inherited.fields, { secret: 'secret' }); rejected(inherited);
    const nested = domainRecord(); Object.defineProperty(nested.fields, 'reference_id', { value: { secret: 'secret' } }); rejected(nested);
  });
  it.each(['example.com', '127.0.0.1', '[::1]', '/private/path', 'A B', 'A'.repeat(129), ''])('rejects unsafe identity %s', reference_id => {
    rejected(domainRecord('domain.delivery.blocked', { fields: { reference_id } }));
  });
  it.each<RuntimeLogRecord['fields']>([
    { event_version: 2 }, { locator_step: 2 }, { stage: 'queued' }, { stage: 'made-up' }, { reason: 'none', action: 'none' },
    { reason: 'raw-error' }, { action: 'retry' }, { evidence_kind: 'runtime-receipt' }, { group_revision: 'R2' },
    { point_id: 'P' }, { tag_id: 'T' }, { device_id: 'D' }, { rule_id: 'RULE' }, { mapping_id: 'M' },
    { intake_state: 'active' }, { pending_samples: '0' }, { record_id: 'REC' },
  ])('rejects contradictory identity, stage or proof %j', fields => rejected(domainRecord('domain.delivery.blocked', { fields })));
  it.each(['workspace_id', 'workspace_revision', 'entity_revision', 'reference_id', 'group_id', 'group_revision', 'connector_id', 'connector_revision'])('requires %s', key => {
    const input = domainRecord(); delete input.fields[key]; rejected(input);
  });
  it('requires exact source and fixed code severity', () => {
    rejected(domainRecord('domain.delivery.blocked', { source: 'http' }));
    rejected(domainRecord('domain.delivery.blocked', { level: 'info' }));
    rejected(domainRecord('domain.source.applied', { level: 'warn' }));
  });
  it.each(['counter-replay-proof-missing', 'counter-replay-proof-conflict', 'counter-replay-proof-unsupported'] as const)('keeps counter reason %s on inputs only', reason => {
    const fields = { reason, action: actionForEvidenceReason(reason) };
    expect(projectRuntimeLogDomain(domainRecord('domain.counter-replay.quarantined', { fields }))).toMatchObject({ kind: 'domain', pendingSamples: '9007199254740993' });
    rejected(domainRecord('domain.delivery.blocked', { fields }));
  });
  it.each(['recording-ownership-unavailable', 'recording-ownership-conflict', 'recording-handoff-proof-conflict', 'recording-cutover-finalized-overlap'] as const)('keeps owner reason %s on allowed references', reason => {
    const fields = { reason, action: actionForEvidenceReason(reason) };
    for (const code of ['domain.delivery.blocked', 'domain.recording.refused', 'domain.cutover.rejected'] as const)
      expect(projectRuntimeLogDomain(domainRecord(code, { fields }))).toMatchObject({ kind: 'domain', reason });
    rejected(domainRecord('domain.source.rejected', { fields }));
    rejected(domainRecord('domain.repair.blocked', { fields }));
  });
  it('does not accept ordinary reasons in counter/recording-refusal evidence', () => {
    for (const code of ['domain.counter-replay.quarantined', 'domain.recording.refused'] as const)
      rejected(domainRecord(code, { fields: { reason: 'unavailable', action: 'refresh' } }));
  });
  it('validates receipt effect/reference identity, provenance and complete branch', () => {
    for (const fields of [{ effect_key: 'OTHER' }, { receipt_provenance: 'test-write' }, { intake_state: 'active' }] as RuntimeLogRecord['fields'][])
      rejected(domainRecord('domain.delivery.first-confirmed', { fields }));
    for (const key of ['record_id', 'effect_key', 'payload_digest', 'committed_at', 'receipt_provenance']) {
      const input = domainRecord('domain.delivery.first-confirmed'); delete input.fields[key]; rejected(input);
    }
    expect(projectRuntimeLogDomain(domainRecord('domain.delivery.first-confirmed'))).toMatchObject({ receipt: {
      group_revision: 'R1', connector_revision: 'C1', effect_key: 'REF1', provenance: 'runtime', committed_at: domainInstant,
    } });
  });
  it('checks intake state against its fixed event rather than current aggregate status', () => {
    rejected(domainRecord('domain.group.paused', { fields: { intake_state: 'active' } }));
    rejected(domainRecord('domain.group.retiring', { fields: { intake_state: 'not_running' } }));
    const input = domainRecord('domain.group.applied'); delete input.fields.intake_state; rejected(input);
  });
  it('retains exact global/group quota zero and overshoot without derived threshold or allocation', () => {
    for (const quota_scope of ['global', 'group']) {
      expect(projectRuntimeLogDomain(domainRecord('domain.quota.refused', { fields: { quota_scope, quota_used_bytes: 120 } })))
        .toMatchObject({ quota: { scope: quota_scope, used_bytes: 120, max_bytes: 100, state: 'hard_limit', intake_refused: true } });
      expect(projectRuntimeLogDomain(domainRecord('domain.quota.recovered', { fields: { quota_scope, quota_used_bytes: 0, quota_state: 'ok' } })))
        .toMatchObject({ quota: { used_bytes: 0, intake_refused: false } });
    }
  });
  it.each<RuntimeLogRecord['fields']>([{ quota_scope: 'device' }, { quota_state: 'unavailable' }, { quota_used_bytes: -1 }, { quota_used_bytes: 0.1 },
    { quota_used_bytes: Number.MAX_SAFE_INTEGER + 1 }, { quota_max_bytes: Infinity }, { intake_refused: false }])('rejects malformed refusal quota %j', fields => {
    rejected(domainRecord('domain.quota.refused', { fields }));
  });
  it.each<RuntimeLogRecord['fields']>([{ quota_max_bytes: 0 }, { quota_used_bytes: 100 }, { quota_state: 'hard_limit' }, { intake_refused: true }])('rejects contradictory capacity recovery %j', fields => {
    rejected(domainRecord('domain.quota.recovered', { fields }));
  });
  it.each(['0', '18446744073709551615'])('allows exact pending-input count %s', pending_samples => {
    expect(projectRuntimeLogDomain(domainRecord('domain.recording.refused', { fields: { pending_samples } }))).toMatchObject({ pendingSamples: pending_samples });
  });
  it.each(['-1', '00', '18446744073709551616'])('rejects malformed pending count %s', pending_samples => {
    rejected(domainRecord('domain.recording.refused', { fields: { pending_samples } }));
  });
  it('keeps identity-free unavailable and bounded summary without a locator', () => {
    expect(projectRuntimeLogDomain(domainRecord('domain.unavailable'))).toEqual({ kind: 'unavailable', timestamp: domainInstant, summary: null });
    expect(projectRuntimeLogDomain(domainRecord('domain.unavailable', { fields: { first_seen_at: domainInstant, last_seen_at: domainInstant, occurrences: '5' } })))
      .toEqual({ kind: 'unavailable', timestamp: domainInstant, summary: { firstSeenAt: domainInstant, lastSeenAt: domainInstant, occurrences: '5' } });
    expect(projectRuntimeLogDomain(domainRecord('domain.unavailable', { fields: { workspace_id: 'W' } }))).toEqual({ kind: 'unavailable', timestamp: domainInstant, summary: null });
  });
  it('keeps unknown domain codes unavailable instead of falling into raw legacy rendering', () => {
    rejected(domainRecord('domain.delivery.blocked', { code: 'domain.unknown', message: 'secret' }));
    expect(projectRuntimeLogDomain({ ...domainRecord(), code: 'runtime.started' })).toEqual({ kind: 'legacy' });
  });
});
