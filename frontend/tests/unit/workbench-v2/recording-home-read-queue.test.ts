import { describe, expect, it } from 'vitest';
import { createRecordingReadQueue } from '@/features/datalink/workbench-v2/recording-home/recordingHomeReadQueue';
const tick = () => new Promise(resolve => setTimeout(resolve, 0));
describe('bounded home reads', () => {
  it('runs at most6 requests while preserving every result', async () => {
    const queue = createRecordingReadQueue(6); let active = 0; let peak = 0; const complete: (() => void)[] = [];
    const results = Array.from({ length: 100 }, (_, id) => queue.run(() => new Promise<number>(resolve => {
      active++; peak = Math.max(peak, active); complete.push(() => { active--; resolve(id); });
    })));
    await tick(); expect(peak).toBe(6);
    while (complete.length || active) { complete.splice(0).forEach(done => done()); await tick(); }
    expect(await Promise.all(results)).toEqual(Array.from({ length: 100 }, (_, id) => id)); expect(peak).toBe(6);
  });
  it('rapid page cancellation never starts queued requests and releases their promises', async () => {
    const queue = createRecordingReadQueue(1); let release!: () => void;
    const first = queue.run(() => new Promise<void>(resolve => { release = resolve; }));
    const controller = new AbortController(); let started = false;
    const second = queue.run(async () => { started = true; }, controller.signal).catch(error => error.name);
    controller.abort(); expect(await second).toBe('AbortError'); expect(started).toBe(false);
    release(); await first; expect(await queue.run(async () => 3)).toBe(3);
  });
  it('an in-flight abort does not falsely free capacity before its transport settles', async () => {
    const queue = createRecordingReadQueue(1); const controller = new AbortController(); let release!: () => void; let next = false;
    const first = queue.run(() => new Promise<void>(resolve => { release = resolve; }), controller.signal).catch(error => error.name);
    const second = queue.run(async () => { next = true; }); controller.abort();
    expect(await first).toBe('AbortError'); expect(next).toBe(false); release(); await second; expect(next).toBe(true);
  });
});
