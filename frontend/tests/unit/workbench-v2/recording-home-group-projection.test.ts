import { describe, expect, it } from 'vitest';
import { recordingHomeProjection, type RecordingHomeInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';

const scope = { workspaceId: 'W1', workspaceRevision: 'R2', epoch: 2 };
const groupScope = { ...scope, groupId: 'G1', groupRevision: 'G2', appliedRevision: 'G1' };
const group: WriteGroup = { id: 'G1', workspace_id: 'W1', revision: 'G2', applied_revision: 'G1',
  name: 'Temperature', status: 'running', members: [],
  destination: { connector_id: 'DB1', connector_revision: 'D1', database: 'telemetry',
    table_schema: 'public', table_name: 'temperature', storage_strategy: 'managed' },
  row_policy: { interval_seconds: 30, allowed_lateness_seconds: 2 }, write_policy: {}, migration: {},
  created_at: '', updated_at: '' };
const lifecycle: WriteGroupLifecycle = { workspace_id: 'W1', workspace_revision: 'R2', group_id: 'G1',
  group_revision: 'G2', applied_revision: 'G1', status: 'running', connector_id: 'DB1', connector_revision: 'D1',
  actions: [{ action: 'stop', allowed: true, reason_code: '', next_action: 'review', requires_review: true }],
  history: { canonical: true, current_group_id: 'G1', generations: [] } };
const delivery: WriteGroupDelivery = { group_id: 'G1', intake: { state: 'active' },
  stages: { collecting: 2, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
  last_sql_committed_at: null, oldest_pending_seconds: 1, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false } };
function input(): RecordingHomeInput {
  return { scope, groups: { state: 'available', scope, value: { workspace_id: 'W1', workspace_revision: 'R2', groups: [group] } },
    lifecycle: { G1: { state: 'available', scope: groupScope, value: lifecycle } },
    delivery: { G1: { state: 'available', scope: groupScope, value: delivery } } };
}
function item(data: RecordingHomeInput) {
  const result = recordingHomeProjection(data);
  if (result.groups.state !== 'available') throw new Error('Expected available groups');
  return result.groups.value[0];
}

describe('home group evidence binds canonical saved and applied revisions', () => {
  it('forwards server lifecycle and delivery without converting running/collecting into SQL success or a command', () => {
    const result = recordingHomeProjection(input());
    expect(result.primaryAction).toBeNull();
    expect(item(input()).lifecycle).toEqual({ state: 'available', value: lifecycle });
    expect(item(input()).delivery).toEqual({ state: 'available', value: delivery });
    expect(item(input()).group.status).toBe('running');
    expect(Object.isFrozen(item(input()).delivery)).toBe(true);
    expect(Object.isFrozen(delivery)).toBe(false);
  });
  it.each([
    [{ workspaceId: 'old' }, 'scope_mismatch'], [{ workspaceRevision: 'old' }, 'revision_mismatch'],
    [{ epoch: 1 }, 'stale'], [{ groupId: 'old' }, 'scope_mismatch'],
    [{ groupRevision: 'old' }, 'revision_mismatch'], [{ appliedRevision: 'old' }, 'revision_mismatch'],
  ] as const)('rejects late group envelope %j without losing saved recording', (patch, reason) => {
    const result = item({ ...input(),
      lifecycle: { G1: { state: 'available', scope: { ...groupScope, ...patch }, value: lifecycle } },
      delivery: { G1: { state: 'available', scope: { ...groupScope, ...patch }, value: delivery } } });
    expect(result.group.name).toBe('Temperature');
    expect(result.lifecycle).toEqual({ state: 'unavailable', reason });
    expect(result.delivery).toEqual({ state: 'unavailable', reason });
  });
  it.each([
    [{ workspace_id: 'foreign' }, 'scope_mismatch'], [{ group_id: 'foreign' }, 'scope_mismatch'],
    [{ workspace_revision: 'old' }, 'revision_mismatch'], [{ group_revision: 'old' }, 'revision_mismatch'],
    [{ applied_revision: 'old' }, 'revision_mismatch'],
  ] as const)('rejects inconsistent lifecycle payload %j', (patch, reason) => {
    expect(item({ ...input(), lifecycle: { G1: { state: 'available', scope: groupScope,
      value: { ...lifecycle, ...patch } } } }).lifecycle).toEqual({ state: 'unavailable', reason });
  });
  it('rejects a foreign delivery payload', () => {
    expect(item({ ...input(), delivery: { G1: { state: 'available', scope: groupScope,
      value: { ...delivery, group_id: 'foreign' } } } }).delivery)
      .toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
  });
  it.each([
    [{ workspace_id: 'foreign' }, 'scope_mismatch'], [{ id: '' }, 'scope_mismatch'],
    [{ revision: '' }, 'revision_mismatch'],
  ] as const)('rejects invalid saved group %j instead of projecting its permissions', (patch, reason) => {
    expect(recordingHomeProjection({ ...input(), groups: { state: 'available', scope,
      value: { workspace_id: 'W1', workspace_revision: 'R2', groups: [{ ...group, ...patch }] } } }).groups)
      .toEqual({ state: 'unavailable', reason });
  });
  it('rejects duplicate group identity rather than combining ambiguous evidence', () => {
    expect(recordingHomeProjection({ ...input(), groups: { state: 'available', scope,
      value: { workspace_id: 'W1', workspace_revision: 'R2', groups: [group, group] } } }).groups)
      .toEqual({ state: 'unavailable', reason: 'invalid_evidence' });
  });
  it('keeps missing lifecycle/delivery unavailable rather than defaulting healthy', () => {
    const result = item({ ...input(), lifecycle: undefined, delivery: undefined });
    expect(result.lifecycle).toEqual({ state: 'unavailable', reason: 'missing' });
    expect(result.delivery).toEqual({ state: 'unavailable', reason: 'missing' });
  });
});

describe('delivery current-revision proof', () => {
  it.each([
    { revision_stages: { group_revision: 'old', stages: delivery.stages } },
    { last_sql_committed_effect: { group_revision: 'old', connector_revision: 'D1', record_id: 'row', effect_key: 'effect', payload_digest: 'digest', committed_at: '2026-10-30T00:00:00Z' } },
    { source_evidence: { state: 'available' as const, group_revision: 'old', sources: [] } },
  ])('does not accept contradictory current proof %j', patch => {
    expect(item({ ...input(), delivery: { G1: { state: 'available', scope: groupScope,
      value: { ...delivery, ...patch } } } }).delivery).toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
  });
  it('retains historical backlog and unavailable source evidence without retargeting to the saved revision', () => {
    const value: WriteGroupDelivery = { ...delivery,
      revision_stages: { group_revision: 'G1', stages: delivery.stages },
      source_evidence: { state: 'unavailable', unavailable_reason: 'revision_mismatch', group_revision: 'old', sources: [] },
      backlog: [{ group_revision: 'old', connector_id: 'DB0', connector_revision: 'D0', table_schema: 'public', table_name: 'old', pending: 1, error_codes: [] }],
    };
    const result = item({ ...input(), delivery: { G1: { state: 'available', scope: groupScope, value } } });
    expect(result.delivery).toEqual({ state: 'available', value });
  });
});

describe.each(['loading', 'error', 'unknown', 'unavailable'] as const)('group %s evidence', state => {
  it('retains matching degradation without inventing an empty or healthy result', () => {
    const evidence = { state, scope: groupScope };
    const result = item({ ...input(), lifecycle: { G1: evidence }, delivery: { G1: evidence } });
    const expected = { state: 'unavailable', reason: state === 'unavailable' ? 'missing' : state };
    expect(result.lifecycle).toEqual(expected);
    expect(result.delivery).toEqual(expected);
  });
  it.each([
    [{ groupId: 'old' }, 'scope_mismatch'], [{ groupRevision: 'old' }, 'revision_mismatch'],
    [{ appliedRevision: 'old' }, 'revision_mismatch'],
  ] as const)('rejects old group degradation %j', (patch, reason) => {
    const evidence = { state, scope: { ...groupScope, ...patch } };
    const result = item({ ...input(), lifecycle: { G1: evidence }, delivery: { G1: evidence } });
    expect(result.lifecycle).toEqual({ state: 'unavailable', reason });
    expect(result.delivery).toEqual({ state: 'unavailable', reason });
  });
});
