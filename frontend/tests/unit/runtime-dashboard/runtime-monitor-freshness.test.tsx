import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter, useNavigate } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '../../../src/services/datalinkClient';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { useRuntimeDashboardState, type RuntimeDashboardState } from '../../../src/features/datalink/runtime-dashboard/useRuntimeDashboardState';
vi.mock('../../../src/services/datalinkClient', async original => ({ ...(await original<typeof import('../../../src/services/datalinkClient')>()), api: { get: vi.fn() } }));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const sources: Source[] = [];
class Source {
  onopen: (() => void) | null = null; onerror: (() => void) | null = null;
  listeners = new Map<string, ((event: Event) => void)>();
  constructor(readonly url: string) { sources.push(this); }
  addEventListener(type: string, fn: ((event: Event) => void)) { this.listeners.set(type, fn); }
  removeEventListener(type: string) { this.listeners.delete(type); }
  close() {}
  emit(type: string, data: unknown) { this.listeners.get(type)?.({ data: JSON.stringify(data) } as MessageEvent<string>); }
  ready(id = 'A') { this.emit('stream_state', { device_id: id, timestamp: '2026-10-29T10:00:00Z', stream_state: { state: 'ready', empty: false, stale: false, degraded: false, unavailable: false } }); }
}
const start = Date.parse('2026-10-29T10:00:00Z');
const response = (data: unknown) => ({ data: { success: true, data } });
const snapshot = (id = 'A', state = 'ready') => ({ running: true, uptime_seconds: 7, snapshot_state: { state, empty: state === 'empty', stale: state === 'stale', degraded: state === 'degraded', unavailable: state === 'unavailable' }, collectors: [{ device_id: id, device_name: `Device ${id}`, protocol: 'modbus_tcp', availability_status: 'available', status: 'running', running: true, last_read_at: '2026-10-29T09:59:00Z', points_total: 1, points_healthy: 1, points_stale: 0, points_error: 0, breaker_state: 'closed' }] });
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(done => { resolve = done; }); return { promise, resolve }; }
let client: QueryClient; let current: RuntimeDashboardState; const renders: Array<{ selected: string | null; collector?: string; received?: number | null }> = [];
function Host() {
  current = useRuntimeDashboardState(); const navigate = useNavigate();
  renders.push({ selected: current.selectedDeviceId, collector: current.snapshot?.collectors[0]?.device_id, received: current.snapshotReceivedAt });
  return <><output>{current.routeState}</output><button onClick={() => void current.onRetrySnapshot()}>Refresh</button><button onClick={current.onReconnectStream}>Reconnect</button><button onClick={() => void navigate('?device_id=B')}>B</button></>;
}
async function flush() { await act(async () => { await Promise.resolve(); await vi.advanceTimersByTimeAsync(1); }); }
async function mount(seed?: (client: QueryClient) => void) { client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); seed?.(client); render(<QueryClientProvider client={client}><MemoryRouter initialEntries={['/studio/runtime?device_id=A']}><Host /></MemoryRouter></QueryClientProvider>); await flush(); await flush(); }
beforeEach(() => {
  vi.useFakeTimers(); vi.setSystemTime(start); vi.clearAllMocks(); sources.length = 0; renders.length = 0; vi.stubGlobal('EventSource', Source);
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response({ workspace_id: 'W', default_device_id: 'A', devices: ['A', 'B'].map(device_id => ({ device_id, name: device_id, protocol: 'modbus_tcp', availability_status: 'available' })), setup: null }));
  vi.mocked(api.get).mockImplementation(async (url, config) => { const params = config?.params; const id = params && typeof params === 'object' && 'device_id' in params && typeof params.device_id === 'string' ? params.device_id : 'A'; return response(url === '/points' ? [] : snapshot(id)); });
});
afterEach(() => { cleanup(); client?.clear(); vi.useRealTimers(); vi.unstubAllGlobals(); });
describe('scoped runtime monitor freshness', () => {
  it('retains the 10:00 snapshot after 10:00:05 disconnect, stale with receipt and no LIVE', async () => {
    await mount(); act(() => sources[0].ready()); await flush(); expect(current.routeState).toBe('live');
    const received = current.snapshotReceivedAt; expect(received).toBeGreaterThanOrEqual(start); expect(received).toBeLessThan(start + 10);
    const saved = current.snapshot; vi.setSystemTime(start + 5000); act(() => sources[0].onerror?.()); await flush();
    expect(current.routeState).toBe('degraded'); expect(current.snapshot).toBe(saved); expect(current.snapshotReceivedAt).toBe(received); expect(current.snapshotStale).toBe(true);
    expect(current.snapshotReceivedAt).not.toBe(Date.parse('2026-10-29T09:59:00Z'));
  });
  it('does not freshen retained evidence on reconnect readiness alone; successful scoped fetch restores LIVE', async () => {
    await mount(); act(() => sources[0].ready()); await flush(); const received = current.snapshotReceivedAt;
    vi.setSystemTime(start + 5000); act(() => sources[0].onerror?.()); await flush();
    fireEvent.click(screen.getByRole('button', { name: 'Reconnect' })); act(() => sources[1].ready()); await flush();
    expect(current.routeState).not.toBe('live'); expect(current.snapshotStale).toBe(true); expect(current.snapshotReceivedAt).toBe(received);
    fireEvent.click(screen.getByRole('button', { name: 'Refresh' })); await flush(); await flush();
    expect(current.routeState).toBe('live'); expect(current.snapshotStale).toBe(false); expect(current.snapshotReceivedAt).toBeGreaterThan(received ?? 0);
  });
  it('marks background snapshot failure stale despite a ready stream and recovers by existing polling', async () => {
    await mount(); act(() => sources[0].ready()); await flush(); const received = current.snapshotReceivedAt;
    vi.mocked(api.get).mockRejectedValueOnce(new Error('snapshot unavailable'));
    fireEvent.click(screen.getByRole('button', { name: 'Refresh' })); await flush(); await flush();
    expect(current.routeState).toBe('degraded'); expect(current.snapshotStale).toBe(true); expect(current.snapshotReceivedAt).toBe(received);
    await act(async () => { await vi.advanceTimersByTimeAsync(5000); }); await flush();
    expect(current.routeState).toBe('live'); expect(current.snapshotStale).toBe(false);
  });
  it.each(['stale', 'degraded', 'unavailable', 'unknown'])('never claims LIVE for snapshot truth %s', async state => {
    vi.mocked(api.get).mockImplementation(async url => response(url === '/points' ? [] : snapshot('A', state)));
    await mount(); act(() => sources[0].ready()); await flush(); expect(current.routeState).not.toBe('live');
  });
  it.each(['empty', 'stale', 'degraded', 'unavailable'])('rejects contradictory ready snapshot flag %s', async flag => {
    const data = snapshot(); Object.assign(data.snapshot_state, { [flag]: true });
    vi.mocked(api.get).mockImplementation(async url => response(url === '/points' ? [] : data));
    await mount(); act(() => sources[0].ready()); await flush(); expect(current.routeState).not.toBe('live');
  });
  it('does not replace missing snapshot truth with a synthetic ready state', async () => {
    const { snapshot_state: _truth, ...data } = snapshot();
    vi.mocked(api.get).mockImplementation(async url => response(url === '/points' ? [] : data));
    await mount(); act(() => sources[0].ready()); await flush(); expect(current.routeState).not.toBe('live');
  });
  it('keeps a batched failure followed by ready/value stale until a later scoped success', async () => {
    await mount(); act(() => sources[0].ready()); await flush(); const received = current.snapshotReceivedAt;
    act(() => { sources[0].onerror?.(); sources[0].ready(); sources[0].emit('value', { device_id: 'A', point_id: 'P', address: '40001', raw_value: 1, transformed_value: 1, quality: 'good', stale: false, timestamp: '2026-10-29T10:00:05Z' }); });
    await flush(); expect(current.routeState).not.toBe('live'); expect(current.snapshotStale).toBe(true); expect(current.snapshotReceivedAt).toBe(received);
  });
  it('does not accept a late A HTTP success as B snapshot or receipt', async () => {
    const pendingA = deferred<ReturnType<typeof response>>();
    vi.mocked(api.get).mockImplementation(async (url, config) => {
      if (url === '/points') return response([]);
      const params = config?.params;
      return params && typeof params === 'object' && 'device_id' in params && params.device_id === 'A' ? pendingA.promise : response(snapshot('B'));
    });
    await mount(); fireEvent.click(screen.getByRole('button', { name: 'B' })); await flush(); await flush();
    act(() => sources.at(-1)?.ready('B')); await flush(); const received = current.snapshotReceivedAt;
    await act(async () => pendingA.resolve(response(snapshot('A')))); await flush();
    expect(current.selectedDeviceId).toBe('B'); expect(current.snapshot?.collectors[0]?.device_id).toBe('B'); expect(current.snapshotReceivedAt).toBe(received);
    expect(renders.filter(row => row.selected === 'B').some(row => row.collector === 'A')).toBe(false);
  });
  it('rejects a foreign collector snapshot returned for the current scoped request', async () => {
    vi.mocked(api.get).mockImplementation(async url => response(url === '/points' ? [] : snapshot('B')));
    await mount(); act(() => sources[0].ready()); await flush();
    expect(current.snapshot).toBeNull(); expect(current.snapshotReceivedAt).toBeNull(); expect(current.routeState).toBe('error');
  });
  it('shows cached prior-mount data stale until the current observer obtains a scoped success', async () => {
    const pending = deferred<ReturnType<typeof response>>();
    vi.mocked(api.get).mockImplementation(async url => url === '/points' ? response([]) : pending.promise);
    await mount(cache => cache.setQueryData(['runtime-dashboard-status', 'A'], snapshot(), { updatedAt: start - 60000 }));
    act(() => sources[0].ready()); await flush();
    expect(current.snapshotReceivedAt).toBe(start - 60000); expect(current.routeState).not.toBe('live'); expect(current.snapshotStale).toBe(true);
    await act(async () => pending.resolve(response(snapshot()))); await flush();
    expect(current.routeState).toBe('live'); expect(current.snapshotReceivedAt).toBeGreaterThan(start - 60000);
  });
  it('keeps missing receipt unavailable instead of borrowing acquisition time', async () => {
    const pending = deferred<ReturnType<typeof response>>();
    vi.mocked(api.get).mockImplementation(async url => url === '/points' ? response([]) : pending.promise);
    await mount(cache => cache.setQueryData(['runtime-dashboard-status', 'A'], snapshot(), { updatedAt: 0 }));
    act(() => sources[0].ready()); await flush();
    expect(current.snapshotReceivedAt).toBeNull(); expect(current.snapshotStale).toBe(true); expect(current.routeState).not.toBe('live');
  });
  it('does not treat same-millisecond or structurally shared data as a new receipt identity', async () => {
    await mount(); act(() => sources[0].ready()); await flush(); const received = current.snapshotReceivedAt; const saved = current.snapshot;
    act(() => { sources[0].onerror?.(); sources[0].ready(); }); await flush();
    vi.setSystemTime(received ?? start); fireEvent.click(screen.getByRole('button', { name: 'Refresh' })); await flush(); await flush();
    expect(current.snapshot).toBe(saved); expect(current.snapshotReceivedAt).toBe(received); expect(current.routeState).not.toBe('live');
    vi.setSystemTime((received ?? start) + 1000); fireEvent.click(screen.getByRole('button', { name: 'Refresh' })); await flush(); await flush();
    expect(current.snapshot).toBe(saved); expect(current.snapshotReceivedAt).toBeGreaterThan(received ?? 0); expect(current.routeState).toBe('live');
  });
  it.each(['ready', 'empty'])('keeps zero collectors consistent with producer truth %s', async state => {
    vi.mocked(api.get).mockImplementation(async url => response(url === '/points' ? [] : { ...snapshot('A', state), collectors: [] }));
    await mount(); act(() => sources[0].ready()); await flush();
    expect(current.routeState).not.toBe('live');
    if (state === 'empty') expect(current.routeState).toBe('empty');
  });
  it('requires explicit ready even when a successful snapshot exists', async () => {
    await mount(); act(() => { sources[0].onopen?.(); sources[0].emit('heartbeat', { ts: '2026-10-29T10:00:00Z' }); });
    expect(current.routeState).not.toBe('live'); expect(current.snapshotStale).toBe(true);
  });
  it('binds snapshot and receipt synchronously so no render of B contains retained A', async () => {
    await mount(); act(() => sources[0].ready()); await flush(); const pending = deferred<ReturnType<typeof response>>();
    vi.mocked(api.get).mockImplementation(async url => url === '/points' ? response([]) : pending.promise);
    fireEvent.click(screen.getByRole('button', { name: 'B' })); await flush();
    expect(renders.filter(row => row.selected === 'B').every(row => row.collector === undefined && row.received == null)).toBe(true);
    await act(async () => pending.resolve(response(snapshot('B')))); await flush(); act(() => sources.at(-1)?.ready('B')); await flush();
    expect(current.snapshot?.collectors[0]?.device_id).toBe('B'); expect(current.routeState).toBe('live');
  });
});
