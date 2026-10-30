import { describe, expect, it, vi } from 'vitest';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { cancelPicker, createPicker, filterPicker, pendingAddition, pickerResults, removeMember, restoreMember, sourceKey, togglePending } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingPickerState';
import type { PickerProjection } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingPickerState';

const candidate = (id: string, device = 'd1'): GroupCandidate => ({ key: id, device_id: device, point_id: id, tag_id: `t-${id}`, tag_key: id, label: 'same name', address: '40001', device_name: device, target_type: 'float32' });
const points = Array.from({ length: 60 }, (_, i) => candidate(`p${String(i).padStart(2, '0')}`, i < 30 ? 'd1' : 'd2'));
const projection: PickerProjection = { scope: 'workspace/session', candidates: points, members: points.slice(0, 5), evidence: Object.fromEntries(points.map(p => [sourceKey(p), { availability: 'eligible', semantic: 'gauge' }])) };

describe('recording editor additive picker foundation', () => {
  it('keeps 5 existing identities locked and 3 pending across filters, with 2 hidden', () => {
    let state = createPicker(projection.scope);
    state = togglePending(state, points[0], projection);
    expect(state.pending).toHaveLength(0);
    for (const point of [points[5], points[35], points[36]]) state = togglePending(state, point, projection);
    state = filterPicker(state, { deviceId: 'd1' });
    expect(state.pending).toHaveLength(3);
    expect(pickerResults(state, projection)).toMatchObject({ total: 30, added: 5, hiddenPending: 2, canAdd: true });
    expect(pendingAddition(state, projection)?.map(p => p.point_id)).toEqual(['p05', 'p35', 'p36']);
  });
  it('resets query/type/device/page-size changes and deterministically bounds 25/50 rows', () => {
    const state = { ...createPicker(projection.scope), page: 99 };
    expect(filterPicker(state, { query: 'same' }).page).toBe(0);
    expect(filterPicker(state, { deviceId: 'd2' }).page).toBe(0);
    expect(filterPicker(state, { pointType: 'state' }).page).toBe(0);
    expect(filterPicker(state, { pageSize: 50 }).page).toBe(0);
    const result = pickerResults(createPicker(projection.scope), { ...projection, candidates: [...points].reverse() });
    expect(result.rows).toHaveLength(25);
    expect(result.rows[0].point_id).toBe('p00');
    expect(pickerResults({ ...state, page: 0, pageSize: 50 }, projection).rows).toHaveLength(50);
    expect(pickerResults(state, projection)).toMatchObject({ page: 2, pages: 3 });
  });
  it('keeps page-hidden pending and cancellation discards pending only', () => {
    const original = { ...createPicker(projection.scope), query: 'same', focusKey: 'trigger', scrollTop: 120 };
    const pending = togglePending(original, points[35], projection);
    expect(pickerResults(pending, projection).hiddenPending).toBe(1);
    expect(cancelPicker(pending)).toEqual(original);
    expect(pickerResults(original, projection).canAdd).toBe(false);
  });
  it.each(['invalid', 'stale', 'unknown'] as const)('blocks %s evidence even when the canonical ID remains', availability => {
    const state = togglePending(createPicker(projection.scope), points[5], projection);
    const stale = { ...projection, evidence: { ...projection.evidence, [sourceKey(points[5])]: { availability } } };
    expect(pickerResults(state, stale).invalidPending).toHaveLength(1);
    expect(pendingAddition(state, stale)).toBeUndefined();
    expect(togglePending(createPicker(projection.scope), points[5], stale).pending).toHaveLength(0);
  });
  it('blocks deleted IDs and cross-scope state, never substitutes same name/address', () => {
    const state = togglePending(createPicker(projection.scope), points[5], projection);
    expect(pickerResults(state, { ...projection, candidates: points.filter(p => p !== points[5]) }).invalidPending).toHaveLength(1);
    expect(pendingAddition(state, { ...projection, scope: 'other' })).toBeUndefined();
    expect(pickerResults(state, { ...projection, scope: 'other' }).rows).toEqual([]);
  });
  it('does not infer point meaning from datatype or name', () => {
    const noEvidence = { ...projection, evidence: {} };
    expect(pickerResults(filterPicker(createPicker(projection.scope), { pointType: 'gauge' }), noEvidence).total).toBe(0);
    expect(togglePending(createPicker(projection.scope), points[5], noEvidence).pending).toEqual([]);
  });
  it('preserves full member intent, including invalid text, for undo or same-ID re-add', () => {
    const member = { ...points[5], method: 'mean', intervalText: 'not a number', nested: { original: ' 15x ' } };
    const draft = { name: 'unsaved record', destination: 'unsaved table', members: [member] };
    const removed = removeMember(draft.members, { scope: projection.scope, removed: new Map() }, sourceKey(member));
    expect(removed.members).toEqual([]);
    const initialize = vi.fn(() => ({ ...member, method: 'latest', intervalText: '15' }));
    const restored = restoreMember(removed.members, removed.history, member, projection, initialize);
    expect(restored).toEqual([member]);
    expect(restored[0]).toBe(member);
    expect(initialize).not.toHaveBeenCalled();
    expect(draft.members).toEqual([member]);
    expect(restoreMember([], { ...removed.history, scope: 'other' }, member, projection, initialize)).toEqual([]);
    expect(restoreMember([], removed.history, member, { ...projection, evidence: {} }, initialize)).toEqual([]);
  });
});
