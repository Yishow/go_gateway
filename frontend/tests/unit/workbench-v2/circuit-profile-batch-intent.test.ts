import { describe, expect, it } from 'vitest';
import { defaultCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { applyUpdatedGroupsToActiveDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBatch';
import { existingProfileUpdateGroup, profileMatchesIntent } from '@/features/datalink/workbench-v2/state/circuitProfileBatchIntent';
import { reconcileBatchResults, type BatchTargetExecutionResult } from '@/features/datalink/workbench-v2/state/circuitProfileBatchOutcome';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { WriteGroup, WriteGroupMember, WriteGroupProfileKind } from '@/types/studioV2WriteGroup';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { savedGroupState } from '../../fixtures/writeGroupState';

const INTERVALS = { instantaneous: 30, cumulative: 20 };
const unknown = (kind: WriteGroupProfileKind): Record<string, BatchTargetExecutionResult> => ({
  'cir-a': { circuit_id: 'cir-a', circuit_name: 'Feeder A', status: 'unknown', profiles: { [kind]: { kind, status: 'unknown' } } },
});
const reviewedInst = () => longProfileGroup({ id: 'inst-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] });
const proofOf = (live: WriteGroup[], reviewed = reviewedInst()) =>
  reconcileBatchResults([{ circuit_id: 'cir-a', circuit_name: 'Feeder A', instantaneousGroup: reviewed, anyGroup: reviewed }], INTERVALS, unknown('instantaneous'), live, 'ws-1');

describe('An existing profile is proven only by the reviewed update', () => {
  const committed = (): WriteGroup => ({ ...existingProfileUpdateGroup(reviewedInst(), 'instantaneous', 30), revision: 'rev-inst-a-2' });

  it('accepts the committed update even when the server restamped its own evidence', () => {
    const live: WriteGroup = {
      ...committed(), status: 'running', applied_revision: 'applied-9', updated_at: '2026-10-24T09:00:00Z', name: 'Renamed by the server',
      destination: { ...committed().destination, database: '/other/path.db', schema_digest: 'digest-new', schema_revision: 'schema-9' },
      row_policy: { ...committed().row_policy, record_key_column: 'record_id', bucket_start_column: 'bucket_start' },
      members: committed().members.map((member) => ({ ...member, source_revision: 'restamped', mapping_revision: 'restamped' })),
    };
    const result = proofOf([live]);
    expect(result.results['cir-a'].profiles.instantaneous?.status).toBe('saved');
    expect(result.proven).toEqual([live]);
  });

  it.each([
    ['another interval', (g: WriteGroup) => ({ ...g, recording_profile: { kind: 'instantaneous' as const, interval: 31, unit: 'second' as const }, row_policy: { ...g.row_policy, interval_seconds: 31 } })],
    ['a storage mode', (g: WriteGroup) => ({ ...g, storage_mode: 'timescale' as const })],
    ['the recreate policy', (g: WriteGroup) => ({ ...g, write_policy: { ...g.write_policy, recreate_missing: false } })],
    ['the lateness policy', (g: WriteGroup) => ({ ...g, row_policy: { ...g.row_policy, allowed_lateness_seconds: 5 } })],
    ['the circuit identity', (g: WriteGroup) => ({ ...g, members: g.members.map((m) => ({ ...m, circuit_id: 'cir-other' })) })],
    ['a source binding', (g: WriteGroup) => ({ ...g, members: g.members.map((m, i) => (i === 0 ? { ...m, tag_id: 'tag-9' } : m)) })],
    ['the measurement binding', (g: WriteGroup) => ({ ...g, members: g.members.map((m) => ({ ...m, measurement_id: 'm-other' })) })],
    ['the destination table', (g: WriteGroup) => ({ ...g, destination: { ...g.destination, table_name: 'someone_elses_table' } })],
  ])('keeps unknown when %s differs from the reviewed update', (_name, change) => {
    const result = proofOf([change(committed())]);
    expect(result.results['cir-a'].profiles.instantaneous?.status).toBe('unknown');
    expect(result.proven).toEqual([]);
  });

  it('keeps unknown while the revision has not moved', () => {
    expect(proofOf([{ ...committed(), revision: reviewedInst().revision }]).proven).toEqual([]);
  });

  it('carries an explicit reduction when only the period of a saved profile changes', () => {
    const saved: WriteGroup = {
      ...reviewedInst(),
      recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' },
      row_policy: { ...reviewedInst().row_policy, interval_seconds: 30 },
    };
    const next = existingProfileUpdateGroup(saved, 'instantaneous', 45);
    expect(next.recording_profile?.interval).toBe(45);
    expect(next.row_policy.interval_seconds).toBe(45);
    expect(next.recording_profile?.sample_reduction).toBe('mean');
  });

  it('does not invent a reduction for a saved profile that never had one', () => {
    const next = existingProfileUpdateGroup(reviewedInst(), 'instantaneous', 45);
    expect(next.recording_profile).not.toHaveProperty('sample_reduction');
  });
});

type MemberMetadata = Pick<WriteGroupMember, 'value_role' | 'source_value_type' | 'result_value_type'>;
// Literal persisted member metadata; the helper keeps it explicit that absence is not a role.
const gaugeMetadata: MemberMetadata = { value_role: 'gauge', source_value_type: 'int16', result_value_type: 'decimal' };

function withMemberMetadata(metadata: Partial<MemberMetadata>): WriteGroup {
  const base = reviewedInst();
  return { ...base, members: base.members.map((item, index) => (index === 0 ? { ...item, ...metadata } : item)) };
}

describe('A reviewed intent distinguishes member roles and value types', () => {
  it('matches an identical role and value-type set', () => {
    expect(profileMatchesIntent(withMemberMetadata(gaugeMetadata), withMemberMetadata(gaugeMetadata))).toBe(true);
  });

  it.each([
    ['the value role', { value_role: 'counter_raw' }],
    ['the source type', { source_value_type: 'int32' }],
    ['the result type', { result_value_type: 'float64' }],
  ] as Array<[string, Partial<MemberMetadata>]>)('keeps unknown when %s of a member differs', (_name, change) => {
    expect(profileMatchesIntent(withMemberMetadata(gaugeMetadata), withMemberMetadata({ ...gaugeMetadata, ...change }))).toBe(false);
  });

  it('separates an explicit role from a legacy member without one', () => {
    expect(profileMatchesIntent(withMemberMetadata(gaugeMetadata), withMemberMetadata({}))).toBe(false);
    expect(profileMatchesIntent(withMemberMetadata({}), withMemberMetadata(gaugeMetadata))).toBe(false);
  });

  it('keeps two legacy members without metadata equivalent', () => {
    expect(profileMatchesIntent(withMemberMetadata({}), withMemberMetadata({}))).toBe(true);
  });
});

describe('A new profile is proven only by a writer that matches what would have been created', () => {
  const anyGroup = reviewedInst();
  const created = (over: Partial<WriteGroup> = {}): WriteGroup => ({
    ...anyGroup, id: 'cum-new', revision: 'rev-cum-1', name: 'Feeder A (Cumulative)',
    recording_profile: { kind: 'cumulative', interval: 20, unit: 'minute' },
    row_policy: { ...anyGroup.row_policy, interval_seconds: 1200 }, ...over,
  });
  const proofOfNew = (live: WriteGroup[], workspaceId: string | null = 'ws-1') =>
    reconcileBatchResults([{ circuit_id: 'cir-a', circuit_name: 'Feeder A', instantaneousGroup: anyGroup, anyGroup }], INTERVALS, unknown('cumulative'), live, workspaceId ?? undefined);

  it('accepts the canonical writer that carries the reviewed sources, destination and period', () => {
    const result = proofOfNew([created()]);
    expect(result.results['cir-a'].profiles.cumulative?.status).toBe('saved');
    expect(result.proven.map((group) => group.id)).toEqual(['cum-new']);
  });

  it('skips a first same-circuit writer that does not match and picks the one that does', () => {
    const wrong = created({ id: 'cum-wrong', members: anyGroup.members.slice(0, 1) });
    expect(proofOfNew([wrong, created()]).proven.map((group) => group.id)).toEqual(['cum-new']);
  });

  it.each([
    ['fewer sources', created({ members: anyGroup.members.slice(0, 1) })],
    ['another destination', created({ destination: { ...anyGroup.destination, connector_id: 'other-connector' } })],
    ['another period', created({ recording_profile: { kind: 'cumulative', interval: 25, unit: 'minute' }, row_policy: { ...anyGroup.row_policy, interval_seconds: 1500 } })],
    ['another circuit', created({ members: anyGroup.members.map((m) => ({ ...m, circuit_id: 'cir-other' })) })],
  ])('keeps unknown for a writer with %s', (_name, live) => {
    const result = proofOfNew([live]);
    expect(result.results['cir-a'].profiles.cumulative?.status).toBe('unknown');
    expect(result.proven).toEqual([]);
  });

  it('cannot prove a new profile without the workspace it would have been created in', () => {
    expect(proofOfNew([created()], null).proven).toEqual([]);
  });
});

describe('Draft sync distinguishes a new profile from a period change', () => {
  const candidates = buildGroupCandidates(savedGroupState()).candidates;
  const inst = reviewedInst();
  const base = { ...defaultCircuitProfileLocalDraft(candidates, [inst], 'cir-a'), instantaneous_keys: ['kept\0only\0one'] };

  it('initialises a created profile from its reviewed sources and enables it', () => {
    const cumulative = longProfileGroup({ id: 'cum-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 20, points: [0, 1], kind: 'cumulative' });
    const next = applyUpdatedGroupsToActiveDraft(base, 'cir-a', [cumulative], [inst]);
    expect(next.cumulative_enabled).toBe(true);
    expect(next.cumulative_interval_input).toBe('20');
    expect(next.cumulative_keys).toEqual(cumulative.members.map((m) => `${m.device_id}\0${m.point_id}\0${m.tag_id}`));
    expect(next.instantaneous_keys).toEqual(base.instantaneous_keys);
  });

  it('changes only the period of an existing profile, including one the operator switched off locally', () => {
    const off = { ...base, instantaneous_enabled: false };
    const next = applyUpdatedGroupsToActiveDraft(off, 'cir-a', [{ ...inst, revision: 'rev-2', recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second' } }], [inst]);
    expect(next.instantaneous_interval_input).toBe('30');
    expect(next.instantaneous_enabled).toBe(false);
    expect(next.instantaneous_keys).toEqual(base.instantaneous_keys);
  });
});
