import { act, renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useRecordingPickerDiscovery, type RecordingPickerDiscovery } from '@/features/datalink/workbench-v2/recording-editor/useRecordingPickerDiscovery';
const empty = { candidates: [], devices: [] };
describe('explicit recording picker read projection', () => {
  it('does not read until explicit demand; failure is not a healthy empty result', async () => {
    const load = vi.fn<() => Promise<RecordingPickerDiscovery>>().mockRejectedValue(new Error('private backend details'));
    const hook = renderHook(() => useRecordingPickerDiscovery('S', load));
    expect(load).not.toHaveBeenCalled();
    await act(async () => hook.result.current.refresh());
    expect(hook.result.current.result).toEqual({ scope: 'S', request: 1, status: 'failed' });
    load.mockResolvedValue(empty); await act(async () => hook.result.current.refresh());
    expect(hook.result.current.result).toEqual({ scope: 'S', request: 2, status: 'loaded', value: empty });
  });
  it('late old scope and previous same-scope request cannot replace current observation', async () => {
    const resolves: ((value: RecordingPickerDiscovery) => void)[] = [];
    const load = () => new Promise<RecordingPickerDiscovery>(resolve => resolves.push(resolve));
    const hook = renderHook(({ scope }) => useRecordingPickerDiscovery(scope, load), { initialProps: { scope: 'A' } });
    act(() => { void hook.result.current.refresh(); });
    hook.rerender({ scope: 'B' }); hook.rerender({ scope: 'A' });
    act(() => { void hook.result.current.refresh(); });
    await act(async () => resolves[1](empty));
    await act(async () => resolves[0]({ candidates: [], devices: [{ id: 'OLD', name: 'Old' }] }));
    expect(hook.result.current.result?.value).toEqual(empty);
  });
});
