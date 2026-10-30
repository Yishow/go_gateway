import { act, cleanup, renderHook } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { useRuntimeLogs } from '../../../src/features/runtime-logs/useRuntimeLogs';
import type { RuntimeLogFilters, RuntimeLogSnapshot } from '../../../src/types/runtimeLogs';
import { filters, record, snapshot } from './fixtures';

const scoped: RuntimeLogFilters = { ...filters, workspaceId: 'W', deviceId: 'A', ruleId: 'R',
  groupId: 'G', stage: 'blocked', actionable: true };
const precise = '9007199254740993';
const settle = async () => { for (let i = 0; i < 12; i++) await Promise.resolve(); };
function wire() {
  const snapshots: { url: URL; signal: AbortSignal; resolve: (response: Response) => void }[] = [];
  const streams: { url: URL; signal: AbortSignal; push: (frame: string) => void }[] = [];
  const fetch = vi.fn<typeof globalThis.fetch>((input, options) => {
    const url = new URL(String(input), 'http://localhost'); const signal = options?.signal;
    if (!(signal instanceof AbortSignal)) throw new Error('Missing actual cancellation signal');
    if (!url.pathname.endsWith('/stream')) return new Promise<Response>((resolve) => snapshots.push({ url, signal, resolve }));
    let sink: ReadableStreamDefaultController<Uint8Array> | undefined;
    const body = new ReadableStream<Uint8Array>({ start(controller) { sink = controller; } });
    streams.push({ url, signal, push(frame) {
      if (!sink) throw new Error('Stream not initialized'); sink.enqueue(new TextEncoder().encode(frame));
    } });
    return Promise.resolve(new Response(body, { headers: { 'content-type': 'text/event-stream' } }));
  });
  vi.stubGlobal('fetch', fetch);
  const respond = (index: number, sequence = '1', patch: Partial<RuntimeLogSnapshot> = {}) => {
    const request = snapshots[index]; if (!request) throw new Error('Snapshot was not requested');
    request.resolve(new Response(JSON.stringify({ ...snapshot, latest_cursor: `instance-a:${sequence}`,
      records: [record(sequence, request.url.searchParams.get('device_id') ?? 'unscoped')], ...patch }),
    { headers: { 'content-type': 'application/json' } }));
  };
  const push = (index: number, frame: string) => {
    const stream = streams[index]; if (!stream) throw new Error('Stream was not requested'); stream.push(`${frame}\n\n`);
  };
  return { fetch, snapshots, streams, respond, push };
}
function log(sequence: string, message: string) {
  return `event: log\nid: instance-a:${sequence}\ndata: ${JSON.stringify(record(sequence, message))}`;
}
function progress(sequence: string, type = 'heartbeat') {
  return `event: ${type}\ndata: ${JSON.stringify({ ...snapshot, progress_cursor: `instance-a:${sequence}` })}`;
}
afterEach(() => { cleanup(); vi.unstubAllGlobals(); vi.useRealTimers(); });

describe('actual hook sessions bind domain filters, cancellation and precise progress', () => {
  it('spec example A at 9007199254740993 switches to B; late A stream never enters B or lends its cursor', async () => {
    const h = wire(); const { result, rerender } = renderHook(useRuntimeLogs, { initialProps: scoped });
    await act(async () => { h.respond(0, precise); await settle(); });
    expect(h.streams[0]?.url.searchParams.get('after')).toBe(`instance-a:${precise}`);
    rerender({ ...scoped, deviceId: 'B' });
    expect(h.snapshots).toHaveLength(2); expect(h.streams[0]?.signal.aborted).toBe(true);
    expect(result.current.records).toEqual([]);
    await act(async () => { h.respond(1, '7'); await settle(); h.push(0, log('9007199254740994', 'late A')); await settle(); });
    expect(h.streams[1]?.url.searchParams.get('after')).toBe('instance-a:7');
    expect(h.streams[1]?.url.searchParams.get('device_id')).toBe('B');
    expect(result.current.records.map((r) => r.message)).toEqual(['B']);
  });
  it('late A snapshot after B connects is ignored, cannot start an A stream or change B progress', async () => {
    const h = wire(); const { result, rerender } = renderHook(useRuntimeLogs, { initialProps: scoped });
    rerender({ ...scoped, deviceId: 'B' });
    expect(h.snapshots).toHaveLength(2); expect(h.snapshots[0]?.signal.aborted).toBe(true);
    await act(async () => { h.respond(1, '7'); await settle(); h.respond(0, precise); await settle(); });
    expect(h.streams).toHaveLength(1); expect(h.streams[0]?.url.searchParams.get('after')).toBe('instance-a:7');
    expect(result.current.records.map((r) => r.message)).toEqual(['B']);
  });
  it.each<[keyof RuntimeLogFilters, string | boolean]>([
    ['workspaceId', 'W2'], ['deviceId', 'B'], ['ruleId', 'R2'], ['groupId', 'G2'], ['stage', 'unknown'], ['actionable', false],
  ])('changing only %s replaces the session and removes old progress', async (key, value) => {
    const h = wire(); const { result, rerender } = renderHook(useRuntimeLogs, { initialProps: scoped });
    await act(async () => { h.respond(0, precise); await settle(); });
    const next: RuntimeLogFilters = { ...scoped }; Object.assign(next, { [key]: value }); rerender(next);
    expect(h.snapshots).toHaveLength(2); expect(h.streams[0]?.signal.aborted).toBe(true);
    expect(result.current.records).toEqual([]);
    await act(async () => { h.respond(1, '2'); await settle(); });
    expect(h.streams[1]?.url.searchParams.get('after')).toBe('instance-a:2');
  });
  it('same filter object values keep the session; pause/clear/resume preserve exact processed cursor locally', async () => {
    const h = wire(); const { result, rerender } = renderHook(useRuntimeLogs, { initialProps: scoped });
    await act(async () => { h.respond(0, precise); await settle(); h.push(0, progress('9007199254740995', 'handshake')); await settle(); });
    rerender({ ...scoped }); expect(h.snapshots).toHaveLength(1);
    act(() => { result.current.pause(); result.current.clear(); });
    expect(result.current).toMatchObject({ status: 'paused', records: [] });
    await act(async () => { h.push(0, log('9007199254740996', 'late paused')); result.current.resume(); await settle(); });
    expect(h.snapshots).toHaveLength(1); expect(h.streams).toHaveLength(2);
    expect(h.streams[1]?.url.searchParams.get('after')).toBe('instance-a:9007199254740995');
    expect(result.current.records).toEqual([]);
    await act(async () => { h.push(1, log('9007199254740994', 'already processed')); await settle(); });
    expect(result.current.records).toEqual([]);
    expect(h.fetch.mock.calls.every(([, options]) => !options?.method || options.method === 'GET')).toBe(true);
  });
  it('equivalent omitted/empty identity and false/absent actionable filters retain progress', async () => {
    const h = wire(); const initialProps: RuntimeLogFilters = { ...filters };
    const { rerender } = renderHook(useRuntimeLogs, { initialProps });
    await act(async () => { h.respond(0, precise); await settle(); });
    rerender({ ...filters, workspaceId: '', deviceId: '', ruleId: '', groupId: '', actionable: false });
    expect(h.snapshots).toHaveLength(1); expect(h.streams[0]?.signal.aborted).toBe(false);
    rerender({ ...filters }); expect(h.snapshots).toHaveLength(1);
    expect(h.streams[0]?.url.searchParams.get('after')).toBe(`instance-a:${precise}`);
  });
  it('invalid explicit identity cancels the prior scope and fails safely before another request', async () => {
    const h = wire(); const { result, rerender } = renderHook(useRuntimeLogs, { initialProps: scoped });
    await act(async () => { h.respond(0, precise); await settle(); });
    await act(async () => { rerender({ ...scoped, deviceId: 'https://private.example' }); await settle(); });
    expect(h.streams[0]?.signal.aborted).toBe(true); expect(h.fetch).toHaveBeenCalledTimes(2);
    expect(result.current).toMatchObject({ status: 'failed', records: [] });
  });
  it('403 stops automatic retries with scoped filters; explicit retry starts fresh', async () => {
    vi.useFakeTimers(); const h = wire(); const { result } = renderHook(useRuntimeLogs, { initialProps: scoped });
    const first = h.snapshots[0]; if (!first) throw new Error('Missing initial request');
    await act(async () => { first.resolve(new Response('private denied text', { status: 403 })); await settle(); });
    expect(result.current.status).toBe('denied'); expect(vi.getTimerCount()).toBe(0);
    await act(async () => { await vi.advanceTimersByTimeAsync(60000); }); expect(h.fetch).toHaveBeenCalledTimes(1);
    act(() => result.current.retry()); expect(h.snapshots).toHaveLength(2);
    expect(h.snapshots[1]?.url.searchParams.get('device_id')).toBe('A');
  });
  it('retained fault without recovery remains historical; capture/retention/reset uncertainty never becomes healthy', async () => {
    const h = wire(); const { result } = renderHook(useRuntimeLogs, { initialProps: scoped });
    const fault = { ...record('1'), code: 'domain.delivery.blocked', fields: { stage: 'blocked', action: 'inspect-delivery' } };
    await act(async () => { h.respond(0, '1', { records: [fault], capture_dropped: 2 }); await settle();
      h.push(0, progress('1', 'handshake')); h.push(0, 'event: gap\ndata: {"reason":"retention","to":"instance-a:999"}'); await settle(); });
    expect(result.current).toMatchObject({ gap: true, records: [fault], status: 'connected' });
    act(() => result.current.clear()); expect(result.current.gap).toBe(true);
    await act(async () => { h.push(0, 'event: reset\ndata: {"reason":"reset"}'); await settle(); });
    expect(result.current).toMatchObject({ gap: true, reset: true, records: [] });
    expect(result.current).not.toHaveProperty('healthy');
  });
});
