import { describe, expect, it } from 'vitest';
import { createPicker } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingPickerState';
import { beginSource, chooseSourceDevice, createSourceFlow, editSourceBuffer, leaveSource, settleSource, submitSource } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingSourceFlow';
const devices = [{ id: 'first' }, { id: 'empty-device' }];
const picker = { ...createPicker('scope'), query: 'voltage', page: 2, scrollTop: 130, focusKey: 'source-trigger', pending: [{ device_id: 'first', point_id: 'p1', tag_id: 't1' }] };
const initial = () => createSourceFlow<{ name: string }>(picker, 'modal-session');
function editing() { return editSourceBuffer(chooseSourceDevice(beginSource(initial(), devices), 'empty-device', devices, { name: '' }), { name: 'range draft' }); }
function submitted() { return submitSource(editing(), 'request-1'); }
const match = { scope: 'scope', session: 'modal-session', requestId: 'request-1', generation: 1 };

describe('unmounted source callback state contract, not persistence', () => {
  it('uses a single view without silently picking the first device', () => {
    const choice = beginSource(initial(), devices);
    expect(choice.view).toBe('device-choice');
    expect(choice.targetDeviceId).toBeUndefined();
    expect(chooseSourceDevice(choice, 'missing', devices, { name: '' })).toBe(choice);
    expect(chooseSourceDevice(choice, undefined, devices, { name: '' }).view).toBe('device-buffer');
    expect(chooseSourceDevice(choice, 'empty-device', devices, { name: '' })).toMatchObject({ view: 'range-buffer', targetDeviceId: 'empty-device' });
    const scoped = createSourceFlow({ ...picker, deviceId: 'empty-device' }, 's');
    expect(beginSource(scoped, devices)).toMatchObject({ view: 'range-buffer', targetDeviceId: 'empty-device' });
  });
  it('returns through the unsent dirty guard with picker position/pending untouched', () => {
    const flow = editing();
    expect(leaveSource(flow).outcome).toBe('needs-confirmation');
    expect(leaveSource(flow, true)).toMatchObject({ outcome: 'returned', flow: { view: 'picker', dirty: false, picker } });
  });
  it('returns a typed explicit request once, retaining unresolved work instead of calling it cancelled', () => {
    const { flow, request } = submitted();
    expect(request).toMatchObject({ ...match, kind: 'range', deviceId: 'empty-device', buffer: { name: 'range draft' } });
    expect(submitSource(flow, 'duplicate').request).toBeUndefined();
    expect(leaveSource(flow, true).outcome).toBe('needs-recovery');
    const unknown = settleSource(flow, { ...match, status: 'unknown' });
    expect(unknown.operation?.status).toBe('unknown');
    expect(leaveSource(unknown, true).outcome).toBe('needs-recovery');
    expect(submitSource(unknown, 'duplicate').request).toBeUndefined();
  });
  it('retains failed input and ignores other scope/session/request/generation', () => {
    const { flow } = submitted();
    for (const patch of [{ scope: 'other' }, { session: 'other' }, { requestId: 'other' }, { generation: 8 }]) {
      expect(settleSource(flow, { ...match, ...patch, status: 'confirmed', deviceId: 'empty-device', ruleId: 'r1' })).toBe(flow);
    }
    const failed = settleSource(flow, { ...match, status: 'failed' });
    expect(failed).toMatchObject({ view: 'range-buffer', dirty: true, buffer: { name: 'range draft' }, operation: { status: 'failed' } });
    expect(submitSource(failed, 'retry').request).toBeDefined();
  });
  it('does not close or clear newer dirty input on older submitted generation success', () => {
    const { flow } = submitted();
    const newer = editSourceBuffer(flow, { name: 'new invalid text !!' });
    const result = settleSource(newer, { ...match, status: 'confirmed', deviceId: 'empty-device', ruleId: 'r1' });
    expect(result).toMatchObject({ view: 'range-buffer', dirty: true, buffer: { name: 'new invalid text !!' }, generation: 2 });
    expect(result.confirmed).toEqual([{ kind: 'range', deviceId: 'empty-device', ruleId: 'r1' }]);
  });
  it('requires caller-confirmed matching identity, retaining confirmed source facts on return', () => {
    const { flow } = submitted();
    expect(settleSource(flow, { ...match, status: 'confirmed', deviceId: 'wrong', ruleId: 'r1' })).toBe(flow);
    expect(settleSource(flow, { ...match, status: 'confirmed', deviceId: 'empty-device' })).toBe(flow);
    const confirmed = settleSource(flow, { ...match, status: 'confirmed', deviceId: 'empty-device', ruleId: 'r1' });
    expect(confirmed).toMatchObject({ view: 'picker', dirty: false, picker });
    expect(leaveSource(confirmed, true).flow.confirmed).toHaveLength(1);
  });
  it('retains a separately confirmed device when unsent first-range edits are discarded', () => {
    const device = chooseSourceDevice(beginSource(initial(), devices), undefined, devices, { name: 'new device' });
    const { flow } = submitSource(editSourceBuffer(device, { name: 'new device' }), 'request-1');
    const saved = settleSource(flow, { ...match, status: 'confirmed', deviceId: 'server-device' });
    expect(saved).toMatchObject({ view: 'range-buffer', targetDeviceId: 'server-device', buffer: undefined, dirty: false });
    expect(leaveSource(editSourceBuffer(saved, { name: 'unsent range' }), true).flow.confirmed).toEqual([{ kind: 'device', deviceId: 'server-device' }]);
  });
});
