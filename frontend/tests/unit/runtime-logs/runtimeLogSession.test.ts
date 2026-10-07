import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { RuntimeLogSession } from '../../../src/features/runtime-logs/runtimeLogSession';
import { RuntimeLogRequestError } from '../../../src/services/runtimeLogs';
import type { RuntimeLogEvent, RuntimeLogSnapshot, RuntimeLogTransport } from '../../../src/types/runtimeLogs';
import { filters, record, snapshot } from './fixtures';
const settle = async () => { await Promise.resolve(); await Promise.resolve(); await Promise.resolve(); };
function harness() {
  let emit: (event: RuntimeLogEvent) => void = () => undefined;
  let fail: (error: unknown) => void = () => undefined;
  const transport: RuntimeLogTransport = {
    snapshot: vi.fn().mockResolvedValue(snapshot),
    stream: vi.fn((_filters, _cursor, _signal, callback) => {
      emit = callback; return new Promise<void>((_resolve, reject) => { fail = reject; });
    }),
  };
  return { transport, session: new RuntimeLogSession(filters, transport, () => 0.5),
    emit: (event: RuntimeLogEvent) => emit(event), drop: (error = new Error('network')) => fail(error) };
}
describe('runtime log connection ownership', () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());
  it('requires handshake, retains rows on disconnect and owns one retry timer', async () => {
    const h = harness(); h.session.start(); await settle();
    expect(h.session.getSnapshot().status).toBe('loading');
    h.emit({ type: 'handshake', metadata: snapshot, cursor: 'instance-a:1' });
    expect(h.session.getSnapshot().status).toBe('connected');
    h.drop(); await settle();
    expect(h.session.getSnapshot()).toMatchObject({ status: 'reconnecting', records: [record()] });
    expect(vi.getTimerCount()).toBe(1);
    await vi.advanceTimersByTimeAsync(1000);
    expect(h.transport.stream).toHaveBeenCalledTimes(2);
    h.session.stop(); expect(vi.getTimerCount()).toBe(0);
  });
  it('does not retain an unbounded hidden snapshot record array in metadata', async () => {
    const h = harness();
    vi.mocked(h.transport.snapshot).mockResolvedValue({ ...snapshot,
      records: Array.from({ length: 500 }, (_, i) => record(String(i + 1), '界'.repeat(2000))),
      latest_cursor: 'instance-a:500',
    });
    h.session.start(); await settle();
    expect(h.session.getSnapshot().metadata).not.toHaveProperty('records');
    expect(h.session.getSnapshot().records.length).toBeLessThan(500);
    h.session.stop();
  });
  it('pauses/clears locally and resumes from processed progress, not latest retained', async () => {
    const h = harness(); h.session.start(); await settle();
    h.emit({ type: 'heartbeat', metadata: snapshot, cursor: 'instance-a:9007199254740993' });
    h.session.pause(); h.session.clear();
    expect(h.session.getSnapshot()).toMatchObject({ records: [], status: 'paused' });
    h.emit({ type: 'log', record: record('3'), cursor: 'instance-a:3' });
    h.session.resume(); await settle();
    expect(h.transport.snapshot).toHaveBeenCalledTimes(1);
    expect(h.transport.stream).toHaveBeenLastCalledWith(filters, 'instance-a:9007199254740993', expect.any(AbortSignal), expect.any(Function));
    h.session.stop();
  });
  it.each(['snapshot', 'stream'] as const)('stops retries on %s 403 and allows explicit fresh retry', async (method) => {
    const h = harness(); vi.mocked(h.transport[method]).mockRejectedValue(new RuntimeLogRequestError(403));
    h.session.start(); await settle();
    expect(h.session.getSnapshot().status).toBe('denied'); expect(vi.getTimerCount()).toBe(0);
    h.session.retry(); await settle(); expect(h.transport.snapshot).toHaveBeenCalledTimes(2); h.session.stop();
  });
  it('ignores canceled filter A snapshot and old stream callbacks', async () => {
    const a = harness(); const b = harness(); let finish!: (value: RuntimeLogSnapshot) => void;
    vi.mocked(a.transport.snapshot).mockReturnValue(new Promise((resolve) => { finish = resolve; }));
    a.session.start(); a.session.stop(); b.session.start(); await settle();
    finish({ ...snapshot, records: [record('99')] }); await settle();
    expect(a.transport.stream).not.toHaveBeenCalled();
    expect(vi.mocked(a.transport.snapshot).mock.calls[0][1].aborted).toBe(true);
    b.session.stop(); b.emit({ type: 'log', record: record('88'), cursor: 'instance-a:88' }); b.drop(); await settle();
    expect(b.session.getSnapshot().records).toEqual(snapshot.records); expect(vi.getTimerCount()).toBe(0);
  });
  it('backs off 1/2/4/8/16/30 seconds with a hard cap', async () => {
    const h = harness(); vi.mocked(h.transport.stream).mockRejectedValue(new Error('network'));
    h.session.start(); await settle();
    for (const [index, delay] of [1000, 2000, 4000, 8000, 16000, 30000, 30000].entries()) {
      expect(vi.getTimerCount()).toBe(1); await vi.advanceTimersByTimeAsync(delay - 1);
      expect(h.transport.stream).toHaveBeenCalledTimes(index + 1); await vi.advanceTimersByTimeAsync(1);
    }
    h.session.stop();
  });
  it('distinguishes file-only loss and process reset without accepting duplicate cleared records', async () => {
    const h = harness(); h.session.start(); await settle(); h.session.clear();
    h.emit({ type: 'log', record: record(), cursor: 'instance-a:1' });
    expect(h.session.getSnapshot().records).toEqual([]);
    h.emit({ type: 'gap', reason: 'file_loss', cursor: '' });
    expect(h.session.getSnapshot()).toMatchObject({ fileGap: true, gap: false });
    h.emit({ type: 'reset', reason: 'reset', cursor: '' });
    h.emit({ type: 'log', record: { ...record(), instance_id: 'instance-b' }, cursor: 'instance-b:1' });
    expect(h.session.getSnapshot()).toMatchObject({ reset: true, gap: true });
    expect(h.session.getSnapshot().records).toHaveLength(1); h.session.stop();
  });
  it('marks a gap only for its own overflow, not another client\'s shared counter', async () => {
    const h = harness(); vi.mocked(h.transport.snapshot).mockResolvedValue({ ...snapshot, subscriber_overflow: 3 });
    h.session.start(); await settle();
    expect(h.session.getSnapshot().gap).toBe(false);
    h.emit({ type: 'handshake', metadata: { ...snapshot, subscriber_overflow: 3 }, cursor: 'instance-a:1' });
    h.emit({ type: 'heartbeat', metadata: { ...snapshot, subscriber_overflow: 4 }, cursor: 'instance-a:1' });
    expect(h.session.getSnapshot().gap).toBe(false);
    h.emit({ type: 'gap', reason: 'subscriber_overflow', cursor: '' });
    expect(h.session.getSnapshot().gap).toBe(true); h.session.stop();
  });
});
