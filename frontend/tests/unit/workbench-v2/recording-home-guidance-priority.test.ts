import { describe, expect, it } from 'vitest';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { delivery, fixture, group, groupScope, owners, scope } from './recording-home-guidance-fixtures';

function mixed() {
  const a = group('A'); const b = group('B'); const c = group('C', 'draft');
  c.destination.connector_id = '';
  const input = fixture(['D1'], true, [c, b, a]);
  return { ...input, owners: { ...input.owners, A: { state: 'available' as const, scope: groupScope(a), value: [
    ...owners(a).filter(item => item.owner !== 'start'),
    { owner: 'start' as const, scope: { ...groupScope(a), deviceId: 'D1' }, state: 'unresolved' as const, reason: 'unknown' as const, operationId: 'OP-A' },
  ] } } };
}

describe('mixed named recordings remain reachable through the same resolver', () => {
  it('checks unknown A, retains historical B, and continues named draft C at destination', () => {
    const result = recordingHomeGuidance(mixed());
    expect(result.reason).toBe('operation_recovery');
    expect(result.primaryAction).toMatchObject({ kind: 'navigate', target: 'recording', groupId: 'A', groupRevision: 'A2',
      section: 'operation', intent: { kind: 'query_operation', operationId: 'OP-A', operationOwner: 'start' } });
    expect(result.recordingItems.map(item => item.id)).toEqual(['A', 'B', 'C']);
    expect(result.draftItems.map(item => [item.id, item.name, item.action])).toEqual([
      ['C', 'Recording C', { kind: 'navigate', scope, target: 'recording', groupId: 'C', groupRevision: 'C2', section: 'destination', intent: null }],
    ]);
    expect(result.attentionItems.map(item => [item.id, item.reason])).toEqual([['A', 'operation_recovery']]);
    const b = result.recordingItems.find(item => item.id === 'B');
    expect(b?.action).toMatchObject({ kind: 'navigate', groupId: 'B', section: 'overview' });
    expect(b?.evidence.delivery).toMatchObject({ state: 'available', value: { last_sql_committed_at: '2026-10-30T10:32:00Z' } });
  });
  it('keeps a draft in both named collections even when its own operation is unresolved', () => {
    const input = mixed();
    const c = group('C', 'draft');
    const result = recordingHomeGuidance({ ...input, owners: { ...input.owners,
      C: { state: 'available', scope: groupScope(c), value: [{ owner: 'schema', scope: { ...groupScope(c), deviceId: 'D1' },
        state: 'unresolved', reason: 'lost_ack' }] } } });
    expect(result.draftItems.map(item => item.id)).toEqual(['C']);
    expect(result.attentionItems.map(item => item.id)).toEqual(['A', 'C']);
    expect(result.attentionItems[1].action).toMatchObject({ section: 'operation', intent: { kind: 'recover_operation' } });
  });
  it('orders operation recovery before delivery faults before drafts, then canonical ID for ties', () => {
    const a = group('A'); const z = group('Z'); const d = group('D', 'draft');
    const input = fixture(['D1'], true, [z, d, a]);
    const faulty = { ...delivery(a), intake: { state: 'blocked' as const, reason: 'source_stale' } };
    const result = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, delivery: {
      ...input.evidence.delivery, A: { state: 'available', scope: groupScope(a), value: faulty },
      Z: { state: 'available', scope: groupScope(z), value: { ...faulty, group_id: 'Z' } },
    } } });
    expect(result.reason).toBe('delivery_attention');
    expect(result.primaryAction).toMatchObject({ groupId: 'A' });
    expect(result.attentionItems.map(item => item.id)).toEqual(['A', 'Z']);
    expect(result.draftItems.map(item => item.id)).toEqual(['D']);
    expect(result.recordingItems.map(item => item.id)).toEqual(['A', 'D', 'Z']);
  });
  it('does not call historical success healthy when the current source is stale', () => {
    const b = group('B'); const input = fixture(['D1'], true, [b]);
    const result = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, delivery: { B: {
      state: 'available', scope: groupScope(b), value: { ...delivery(b), source_evidence: {
        state: 'available', group_revision: 'B1', sources: [{ device_id: 'D1', point_id: 'P1', tag_id: 'T1',
          last_sql_committed: null, latest: { available: true, group_revision: 'B1', connector_revision: 'DB2',
            record_id: 'row', effect_key: 'effect', payload_digest: 'digest', profile_kind: 'instantaneous',
            window_start: '2026-10-30T10:33:00Z', window_end: '2026-10-30T10:34:00Z', status: 'stale', raw_present: false, delivery: 'queued' } }],
      } },
    } } } });
    expect(result.reason).toBe('source_attention');
    expect(result.attentionItems.map(item => item.id)).toEqual(['B']);
    expect(result.recordingItems[0].evidence.delivery).toMatchObject({ state: 'available', value: {
      last_sql_committed_at: '2026-10-30T10:32:00Z', intake: { state: 'active' },
    } });
  });
  it('returns the same priority and item order after canonical collection order changes', () => {
    const input = mixed();
    if (input.evidence.groups?.state !== 'available') throw new Error('fixture');
    const shuffled = { ...input, evidence: { ...input.evidence, groups: { ...input.evidence.groups,
      value: { ...input.evidence.groups.value, groups: [...input.evidence.groups.value.groups].reverse() } } } };
    const result = recordingHomeGuidance(shuffled);
    expect(result.primaryAction).toMatchObject({ groupId: 'A', section: 'operation' });
    expect(result.recordingItems.map(item => item.id)).toEqual(['A', 'B', 'C']);
  });
});
