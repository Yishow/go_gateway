import { StrictMode, useEffect, useRef } from 'react';
import { act, render, renderHook } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { useRecordingRangeBuffer, type RecordingRangeController } from '@/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer';
import type { RangeSourceOutcome, RangeSourcePort, RangeSourcePreparation } from '@/pages/datalink/workbench-v2/rangeSourceSubmission';
import type { StudioV2WorkspaceDeviceRecord } from '@/services/studioV2WorkspaceDevices';

// Prepare lifecycle of the new-range buffer: StrictMode effect replay must not orphan the
// in-flight prepare, while superseded, re-owned and truly unmounted prepares stay discarded.
const record = (id: string) => ({ id, name: id, protocol: 'modbus_tcp' } as StudioV2WorkspaceDeviceRecord);
function deferredPort(observations: RangeSourceOutcome[] = []) {
  const calls: { id: string; resolve: (value: RangeSourcePreparation | undefined) => void; reject: (error: unknown) => void }[] = [];
  const port: RangeSourcePort = {
    prepare: id => new Promise((resolve, reject) => { calls.push({ id, resolve, reject }); }),
    prepareSaved: async (_id, request) => ({ request, status: 'blocked' }), submit: async request => ({ request, status: 'blocked' }),
    inspect: async () => undefined, observations: () => observations,
  };
  const ready = (index: number, activationAllowed = false) => calls[index].resolve({ device: record(calls[index].id), activationAllowed });
  return { port, calls, ready };
}

/** Mirrors the home NewRange owner: start once from a mount effect, guarded by a ref. */
function EffectStart({ port, expose }: { port: RangeSourcePort; expose: (range: RecordingRangeController) => void }) {
  const range = useRecordingRangeBuffer('W', port), started = useRef(false);
  expose(range);
  useEffect(() => { if (!started.current) { started.current = true; void range.start({ id: 'D', name: 'D' }, { address: '40017' }); } }, [range]);
  return null;
}

describe('range buffer prepare lifecycle', () => {
  it('finishes a mount-effect prepare under StrictMode and accepts valid input', async () => {
    const { port, calls, ready } = deferredPort();
    let range!: RecordingRangeController;
    render(<StrictMode><EffectStart port={port} expose={value => { range = value; }} /></StrictMode>);
    expect(calls).toHaveLength(1);
    expect(range.preparing).toBe(true);
    await act(async () => { ready(0); });
    expect(range.preparing).toBe(false);
    expect(range.blocked).toBe(false);
    expect(range.device?.id).toBe('D');
    act(() => { range.change({ ...range.buffer, count: '8', dataType: 'float32', name: 'FLOW_' }); });
    expect(range.parsed?.rule).toMatchObject({ start_address: '40017', count: 8, data_type: 'float32', naming_prefix: 'FLOW_' });
    expect(range.canSubmit).toBe(true);
  });

  it('turns a failed mount-effect prepare into blocked under StrictMode, never endless checking', async () => {
    const { port, calls } = deferredPort();
    let range!: RecordingRangeController;
    render(<StrictMode><EffectStart port={port} expose={value => { range = value; }} /></StrictMode>);
    await act(async () => { calls[0].reject(new Error('readiness 503')); });
    expect(range.preparing).toBe(false);
    expect(range.blocked).toBe(true);
    expect(range.device).toBeUndefined();
    act(() => { range.change({ ...range.buffer, count: '8', dataType: 'float32', name: 'FLOW_' }); });
    expect(range.canSubmit).toBe(false);
  });

  it('never lets an older prepare that returns late overwrite a newer start', async () => {
    const { port, ready } = deferredPort();
    const h = renderHook(() => useRecordingRangeBuffer('W', port));
    let first!: Promise<boolean>, second!: Promise<boolean>;
    act(() => { first = h.result.current.start({ id: 'A', name: 'A' }); });
    act(() => { second = h.result.current.start({ id: 'B', name: 'B' }); });
    await act(async () => { ready(1, true); });
    await expect(second).resolves.toBe(true);
    await act(async () => { ready(0, false); });
    await expect(first).resolves.toBe(false);
    expect(h.result.current.target?.id).toBe('B');
    expect(h.result.current.device?.id).toBe('B');
    expect(h.result.current.activationAllowed).toBe(true);
    expect(h.result.current.preparing).toBe(false);
  });

  it('never lets an older prepare that returns late overwrite a resumed recovery', async () => {
    const saved: RangeSourceOutcome = { status: 'saved-inactive', request: { workspaceId: 'W', sessionId: 'S', requestId: 'REQ', generation: 1, deviceId: 'R',
      buffer: { name: 'OLD_', naming_prefix: 'OLD_', start_address: '40001', count: 2, data_type: 'int16', enabled: false, scale_multiplier: 1, scale_offset: 0, data_format: '',
        skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null } } };
    const { port, ready } = deferredPort([saved]);
    const h = renderHook(() => useRecordingRangeBuffer('W', port));
    let pending!: Promise<boolean>;
    act(() => { pending = h.result.current.start({ id: 'A', name: 'A' }); });
    act(() => { expect(h.result.current.resume('REQ')).toBe(true); });
    await act(async () => { ready(0, true); });
    await expect(pending).resolves.toBe(false);
    expect(h.result.current.target?.id).toBe('R');
    expect(h.result.current.device).toBeUndefined();
    expect(h.result.current.preparing).toBe(false);
    expect(h.result.current.outcome?.request.requestId).toBe('REQ');
    expect(h.result.current.buffer.name).toBe('OLD_');
  });

  it('discards a prepare that returns after the workspace owner changed', async () => {
    const { port, ready } = deferredPort();
    const h = renderHook(({ workspace }) => useRecordingRangeBuffer(workspace, port), { initialProps: { workspace: 'W' } });
    let pending!: Promise<boolean>;
    act(() => { pending = h.result.current.start({ id: 'A', name: 'A' }); });
    h.rerender({ workspace: 'OTHER' });
    await act(async () => { ready(0, true); });
    await expect(pending).resolves.toBe(false);
    expect(h.result.current.device).toBeUndefined();
  });

  it('discards a prepare that returns after a real unmount', async () => {
    const { port, calls, ready } = deferredPort();
    const h = renderHook(() => useRecordingRangeBuffer('W', port));
    let resolved!: Promise<boolean>, rejected!: Promise<boolean>;
    act(() => { resolved = h.result.current.start({ id: 'A', name: 'A' }); });
    const second = renderHook(() => useRecordingRangeBuffer('W', port));
    act(() => { rejected = second.result.current.start({ id: 'B', name: 'B' }); });
    h.unmount(); second.unmount();
    await act(async () => { ready(0, true); calls[1].reject(new Error('late')); });
    await expect(resolved).resolves.toBe(false);
    await expect(rejected).resolves.toBe(false);
    expect(h.result.current.device).toBeUndefined();
    expect(second.result.current.blocked).toBe(false);
  });
});
