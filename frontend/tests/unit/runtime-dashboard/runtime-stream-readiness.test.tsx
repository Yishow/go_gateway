import { act, cleanup, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useRuntimeDashboardStream } from '../../../src/features/datalink/runtime-dashboard/useRuntimeStream';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const sources: Source[] = [];
class Source {
  onopen: (() => void) | null = null; onerror: (() => void) | null = null; closed = false;
  listeners = new Map<string, ((event: Event) => void)>();
  constructor(readonly url: string) { sources.push(this); }
  addEventListener(type: string, fn: ((event: Event) => void)) { this.listeners.set(type, fn); }
  removeEventListener(type: string) { this.listeners.delete(type); }
  close() { this.closed = true; }
  emit(type: string, value: unknown) { this.listeners.get(type)?.({ data: JSON.stringify(value) } as MessageEvent<string>); }
}
const truth = (state = 'ready', device_id = 'A') => ({ device_id, timestamp: '2026-10-29T10:00:00Z', stream_state: { state, empty: state === 'empty', degraded: state === 'degraded', unavailable: state === 'unavailable', stale: state === 'stale' } });
const value = (device_id = 'A', stale = false) => ({ device_id, point_id: 'P', address: '40001', raw_value: 10, transformed_value: 10, quality: stale ? 'stale' : 'good', stale, timestamp: '2026-10-29T10:00:00Z' });
const status = (device_id = 'A') => ({ device_id, status: 'running', points_total: 1, points_healthy: 1, points_stale: 0, points_error: 0, breaker_state: 'closed', last_read_at: '2026-10-29T10:00:00Z' });
beforeEach(() => { sources.length = 0; vi.useFakeTimers(); vi.stubGlobal('EventSource', Source); });
afterEach(() => { cleanup(); vi.useRealTimers(); vi.unstubAllGlobals(); });
describe('explicit runtime stream readiness', () => {
  it('does not establish readiness from open, heartbeat, value or status alone', async () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({ deviceId: 'A', pointIds: ['P'] }));
    act(() => { sources[0].onopen?.(); sources[0].emit('heartbeat', { ts: '2026-10-29T10:00:00Z' }); sources[0].emit('value', value()); sources[0].emit('status', status()); });
    await act(async () => { await vi.advanceTimersByTimeAsync(200); });
    expect(result.current.connectionState).not.toBe('connected'); expect(result.current).toMatchObject({ serverReady: false });
    expect(result.current.liveValues.P).toMatchObject({ raw_value: 10 });
    act(() => sources[0].emit('stream_state', truth()));
    expect(result.current).toMatchObject({ serverReady: true, connectionState: 'connected' });
  });
  it('keeps stale point quality local without requiring a nonexistent subsequent ready event', async () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({ deviceId: 'A', pointIds: ['P'] }));
    act(() => { sources[0].emit('stream_state', truth()); sources[0].emit('value', value('A', true)); });
    await act(async () => { await vi.advanceTimersByTimeAsync(200); });
    expect(result.current).toMatchObject({ serverReady: true }); expect(result.current.liveValues.P.stale).toBe(true);
    act(() => sources[0].emit('value', value())); await act(async () => { await vi.advanceTimersByTimeAsync(200); });
    expect(result.current).toMatchObject({ serverReady: true, connectionState: 'connected' }); expect(result.current.liveValues.P.stale).toBe(false);
  });
  it.each(['empty', 'stale', 'degraded', 'unavailable'])('revokes readiness on explicit non-ready truth %s', state => {
    const { result } = renderHook(() => useRuntimeDashboardStream({ deviceId: 'A' }));
    act(() => sources[0].emit('stream_state', truth())); act(() => sources[0].emit('stream_state', truth(state)));
    expect(result.current).toMatchObject({ serverReady: false }); expect(result.current.connectionState).not.toBe('connected');
    act(() => sources[0].emit('value', value())); expect(result.current).toMatchObject({ serverReady: false });
  });
  it('rejects inconsistent ready flags, foreign values/status and malformed events', async () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({ deviceId: 'A', pointIds: ['P'] }));
    act(() => sources[0].emit('stream_state', { ...truth(), stream_state: { ...truth().stream_state, stale: true } }));
    expect(result.current).toMatchObject({ serverReady: false });
    act(() => { sources[0].emit('value', value('B')); sources[0].emit('status', status('B')); sources[0].emit('stream_state', truth('ready', 'B')); sources[0].emit('value', { secret: 'invalid' }); });
    await act(async () => { await vi.advanceTimersByTimeAsync(200); });
    expect(result.current.liveValues).toEqual({}); expect(result.current.latestStatus).toBeNull(); expect(result.current).toMatchObject({ serverReady: false });
  });
  it('does not accept readiness without exact device identity', () => {
    const { result } = renderHook(() => useRuntimeDashboardStream({ deviceId: 'A' }));
    const { device_id: _device, ...unscoped } = truth();
    act(() => sources[0].emit('stream_state', unscoped)); expect(result.current).toMatchObject({ serverReady: false });
  });
  it('never exposes previous-device values, status, connection or event labels during the B render', async () => {
    const renders: Array<{ deviceId: string; state: string; points: string[]; status?: string; messages: string[] }> = [];
    const { rerender } = renderHook(({ deviceId }) => {
      const stream = useRuntimeDashboardStream({ deviceId, pointIds: ['P'] });
      renders.push({ deviceId, state: stream.connectionState, points: Object.keys(stream.liveValues), status: stream.latestStatus?.device_id, messages: stream.logs.map(log => log.message) });
      return stream;
    }, { initialProps: { deviceId: 'A' } });
    act(() => { sources[0].emit('stream_state', truth()); sources[0].emit('value', value('A', true)); sources[0].emit('status', { ...status(), last_error: 'known error' }); });
    await act(async () => { await vi.advanceTimersByTimeAsync(200); });
    expect(renders.at(-1)).toMatchObject({ state: 'connected', points: ['P'], status: 'A' });
    rerender({ deviceId: 'B' });
    for (const render of renders.filter(row => row.deviceId === 'B')) {
      expect(render.state).not.toBe('connected'); expect(render.points).toEqual([]); expect(render.status).toBeUndefined();
      expect(render.messages).not.toContain('logs.point_stale'); expect(render.messages).not.toContain('logs.device_error');
    }
  });
  it('does not let late old-generation events or reconnect itself re-establish readiness', () => {
    const { result, rerender, unmount } = renderHook(({ deviceId }) => useRuntimeDashboardStream({ deviceId }), { initialProps: { deviceId: 'A' } });
    const oldReady = sources[0].listeners.get('stream_state'); act(() => sources[0].emit('stream_state', truth()));
    rerender({ deviceId: 'B' }); expect(sources[0].closed).toBe(true);
    act(() => oldReady?.({ data: JSON.stringify(truth()) } as MessageEvent<string>)); expect(result.current).toMatchObject({ serverReady: false });
    act(() => sources[1].emit('stream_state', truth('ready', 'B'))); expect(result.current).toMatchObject({ serverReady: true });
    act(() => sources[1].onerror?.()); expect(result.current).toMatchObject({ serverReady: false });
    act(() => result.current.reconnect()); expect(result.current).toMatchObject({ serverReady: false });
    unmount(); expect(sources[2].closed).toBe(true);
  });
});
