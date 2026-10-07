import { afterEach, describe, expect, it, vi } from 'vitest';
import { parseRuntimeLogFrame, runtimeLogsURL, runtimeLogTransport, RuntimeLogRequestError } from '../../../src/services/runtimeLogs';
import { filters, record, snapshot } from './fixtures';
const metadata = { ...snapshot, records: undefined, latest_cursor: 'instance-a:100', progress_cursor: 'instance-a:50' };
const signal = () => new AbortController().signal;
afterEach(() => vi.unstubAllGlobals());
describe('bounded runtime log wire contract', () => {
  it('uses configured API base and identical encoded snapshot and stream filters', () => {
    const query = { ...filters, source: 'runtime', q: '告警 <test>' };
    const a = new URL(runtimeLogsURL(query), 'http://localhost');
    const b = new URL(runtimeLogsURL(query, true, 'instance-a:9007199254740993'), 'http://localhost');
    expect(a.pathname).toBe('/api/v1/system/logs'); expect(b.pathname).toBe('/api/v1/system/logs/stream');
    for (const key of ['level', 'source', 'q']) expect(a.searchParams.get(key)).toBe(b.searchParams.get(key));
    expect(a.searchParams.get('limit')).toBe('200'); expect(b.searchParams.has('limit')).toBe(false);
    expect(b.searchParams.get('after')).toBe('instance-a:9007199254740993');
  });
  it('uses processed progress and does not acknowledge replay from a gap boundary', () => {
    expect(parseRuntimeLogFrame(`event: heartbeat\ndata: ${JSON.stringify(metadata)}`)).toMatchObject({ cursor: 'instance-a:50' });
    for (const reason of ['retention', 'file_loss']) {
      expect(parseRuntimeLogFrame(`event: gap\ndata: ${JSON.stringify({ reason, from: 'instance-a:1', to: 'instance-a:100' })}`))
        .toEqual({ type: 'gap', reason, cursor: '' });
    }
  });
  it('accepts Go omitted fields and keeps sequence strings above 2^53', () => {
    const { fields: _fields, ...value } = record('9007199254740993');
    expect(parseRuntimeLogFrame(`event: log\nid: instance-a:9007199254740993\ndata: ${JSON.stringify(value)}`))
      .toMatchObject({ type: 'log', record: { ...value, fields: {} }, cursor: 'instance-a:9007199254740993' });
  });
  it('rejects malformed metadata, mismatched IDs and oversized frames with only a safe error', () => {
    expect(() => parseRuntimeLogFrame(`event: log\nid: instance-a:2\ndata: ${JSON.stringify(record())}`)).toThrow(RuntimeLogRequestError);
    expect(() => parseRuntimeLogFrame('event: log\ndata: {"error":"raw secret"}')).toThrow('runtime_log_request_failed');
    expect(() => parseRuntimeLogFrame(`event: log\ndata: ${'x'.repeat(20000)}`)).toThrow(RuntimeLogRequestError);
    expect(() => parseRuntimeLogFrame('event: handshake\ndata: {}')).toThrow(RuntimeLogRequestError);
  });
  it('handles fragmented UTF-8 and CRLF, then releases its reader', async () => {
    const value = record('2', '告警');
    const bytes = new TextEncoder().encode(`event: handshake\r\ndata: ${JSON.stringify(metadata)}\r\n\r\nevent: log\r\nid: instance-a:2\r\ndata: ${JSON.stringify(value)}\r\n\r\n`);
    const body = new ReadableStream<Uint8Array>({ start(controller) {
      for (let i = 0; i < bytes.length; i += 7) controller.enqueue(bytes.slice(i, i + 7)); controller.close();
    } });
    const fetch = vi.fn().mockResolvedValue(new Response(body, { headers: { 'content-type': 'text/event-stream' } }));
    vi.stubGlobal('fetch', fetch); const events = vi.fn();
    await runtimeLogTransport.stream(filters, 'instance-a:1', signal(), events);
    expect(events).toHaveBeenCalledTimes(2); expect(events.mock.calls[1][0].record.message).toBe('告警');
    expect(body.locked).toBe(false);
    expect(fetch.mock.calls[0][1]).toMatchObject({ credentials: 'same-origin', cache: 'no-store', redirect: 'error' });
  });
  it('does not parse raw denied response bodies', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('raw secret', { status: 403 })));
    await expect(runtimeLogTransport.snapshot(filters, signal())).rejects.toMatchObject({ status: 403, message: 'runtime_log_request_failed' });
  });
  it('rejects oversized snapshots and unexpected content types', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('x'.repeat(5 * 1024 * 1024 + 1), { headers: { 'content-type': 'application/json' } })));
    await expect(runtimeLogTransport.snapshot(filters, signal())).rejects.toMatchObject({ status: 422 });
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('<html>private error</html>', { headers: { 'content-type': 'text/html' } })));
    await expect(runtimeLogTransport.snapshot(filters, signal())).rejects.toMatchObject({ status: 422 });
  });
  it('retains file-loss counters independently from capture counters', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ ...snapshot,
      file_dropped: 7, disk_errors: 2, file_recovered_gaps: 1 }), { headers: { 'content-type': 'application/json' } })));
    await expect(runtimeLogTransport.snapshot(filters, signal())).resolves.toMatchObject({
      capture_dropped: 0, file_dropped: 7, disk_errors: 2, file_recovered_gaps: 1,
    });
  });
});
