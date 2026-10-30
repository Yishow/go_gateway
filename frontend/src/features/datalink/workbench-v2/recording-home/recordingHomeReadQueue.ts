/** Limit transport concurrency, not query ownership. Only React Query's own signal
 * cancels a read, so another observer of the same query is never aborted by a page. */
export function createRecordingReadQueue(limit = 6) {
  let active = 0;
  const pending: { start: () => void; cancel: () => void }[] = [];
  const pump = () => { while (active < limit && pending.length) pending.shift()!.start(); };
  return { run<T>(read: () => Promise<T>, signal?: AbortSignal): Promise<T> {
    return new Promise<T>((resolve, reject) => {
      const aborted = () => reject(new DOMException('Read cancelled', 'AbortError'));
      if (signal?.aborted) { aborted(); return; }
      let started = false;
      const entry = {
        cancel: () => {
          const index = pending.indexOf(entry); if (index >= 0) pending.splice(index, 1);
          aborted();
        },
        start: () => {
          if (signal?.aborted) { aborted(); return; }
          started = true; active++;
          void Promise.resolve().then(read).then(resolve, reject).finally(() => {
            active--; signal?.removeEventListener('abort', entry.cancel); pump();
          });
        },
      };
      signal?.addEventListener('abort', entry.cancel, { once: true });
      pending.push(entry); pump();
      if (!started && signal?.aborted) entry.cancel();
    });
  } };
}
export const recordingHomeReadQueue = createRecordingReadQueue(6);
