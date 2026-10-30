import { describe, expect, it } from 'vitest';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { delivery, fixture, group, groupScope, owners } from './recording-home-guidance-fixtures';

const a = group('A');
const base = delivery(a);
const patches: Partial<WriteGroupDelivery>[] = [
  { stages: { ...base.stages, blocked: 1 } }, { stages: { ...base.stages, quarantined: 1 } },
  { stages: { ...base.stages, unknown: 1 } }, { stages: { ...base.stages, retrying: 1 } },
  { quota: { ...base.quota, intake_refused: true } },
  { storage: { generation: 1, history_may_be_lost: false, reconstructions: 0, events: [], repair_state: 'blocked' } },
];

describe('canonical evidence branches never turn faults into healthy claims', () => {
  it.each(patches)('surfaces current delivery fault %j without losing saved name/history', patch => {
    const input = fixture(['D1'], true, [a]);
    const result = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, delivery: { A: {
      state: 'available', scope: groupScope(a), value: { ...base, ...patch },
    } } } });
    expect(result.reason).toBe('delivery_attention');
    expect(result.attentionItems.map(item => [item.id, item.name])).toEqual([['A', 'Recording A']]);
    expect(result.recordingItems[0].evidence.delivery).toMatchObject({ state: 'available', value: { last_sql_committed_at: '2026-10-30T10:32:00Z' } });
  });
  it('keeps acknowledged/collecting without a receipt distinct from SQL commit', () => {
    const input = fixture(['D1'], true, [a]);
    const result = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, delivery: { A: {
      state: 'available', scope: groupScope(a), value: { ...base, last_sql_committed_at: null, stages: { ...base.stages, sql_committed: 0 } },
    } } } });
    expect(result.reason).toBe('recordings');
    expect(result.recordingItems[0].evidence.delivery).toMatchObject({ state: 'available', value: {
      last_sql_committed_at: null, stages: { collecting: 2, sql_committed: 0 },
    } });
    expect(result.primaryAction).toMatchObject({ kind: 'navigate', target: 'recording', section: 'overview', intent: null });
  });
  it.each(['loading', 'error', 'unknown', 'unavailable'] as const)('preserves %s owner degradation as evidence unavailable', state => {
    const input = fixture(['D1'], true, [a]);
    const result = recordingHomeGuidance({ ...input, owners: { A: { state, scope: groupScope(a) } } });
    expect(result.reason).toBe('evidence_unavailable');
    expect(result.attentionItems[0].action).toMatchObject({ groupId: 'A', section: 'overview', intent: null });
  });
  it('inspects the canonical managed-device context even before members are selected', () => {
    const managed = { ...a, members: [], basic_managed_device_id: 'D1' };
    const input = fixture(['D1'], true, [managed]);
    expect(recordingHomeGuidance(input).reason).toBe('recordings');
  });
  it('cannot infer absence of owners for a groupless device context', () => {
    const empty = { ...a, members: [] };
    expect(recordingHomeGuidance(fixture(['D1'], true, [empty])).reason).toBe('evidence_unavailable');
  });
  it('does not hide another member device operation behind a healthy first device', () => {
    const shared = { ...a, members: [...a.members, { ...a.members[0], device_id: 'D2', point_id: 'P2', tag_id: 'T2' }] };
    const input = fixture(['D2', 'D1'], true, [shared]);
    const result = recordingHomeGuidance({ ...input, owners: { A: { state: 'available', scope: groupScope(a), value: [
      ...owners(a), { owner: 'start', scope: { ...groupScope(a), deviceId: 'D2' }, state: 'unresolved', reason: 'query_failed', operationId: 'D2-OP' },
    ] } } });
    expect(result.primaryAction).toMatchObject({ section: 'operation', intent: { operationId: 'D2-OP', scope: { deviceId: 'D2' } } });
  });
  it('does not leak foreign circuit owner evidence into another circuit', () => {
    const circuit = { ...a, members: [{ ...a.members[0], circuit_id: 'C1' }] };
    const input = fixture(['D1'], true, [circuit]);
    const result = recordingHomeGuidance(input);
    expect(result.reason).toBe('evidence_unavailable');
    expect(result.primaryAction).toMatchObject({ section: 'overview', intent: null });
  });
  it('chooses a stable owner when multiple unresolved observations arrive in different orders', () => {
    const input = fixture(['D1'], true, [a]);
    const value = [
      { owner: 'start' as const, scope: { ...groupScope(a), deviceId: 'D1' }, state: 'unresolved' as const, reason: 'unknown' as const, operationId: 'START' },
      { owner: 'schema' as const, scope: { ...groupScope(a), deviceId: 'D1' }, state: 'unresolved' as const, reason: 'unknown' as const, operationId: 'SCHEMA' },
    ];
    const first = recordingHomeGuidance({ ...input, owners: { A: { state: 'available', scope: groupScope(a), value } } });
    const reverse = recordingHomeGuidance({ ...input, owners: { A: { state: 'available', scope: groupScope(a), value: [...value].reverse() } } });
    expect(first.primaryAction).toMatchObject({ intent: { operationOwner: 'schema', operationId: 'SCHEMA' } });
    expect(reverse.primaryAction).toMatchObject({ intent: { operationOwner: 'schema', operationId: 'SCHEMA' } });
  });
});
