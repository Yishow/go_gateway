export interface RuntimeRecordingStreamControls {
  emit(index: number, event: string, data: unknown): void;
  fail(index: number): void;
  list(): { url: string; closed: boolean }[];
}
declare global { interface Window { __runtimeRecordingStreams: RuntimeRecordingStreamControls } }
/** Self-contained for Playwright addInitScript; drives actual EventSource wire listeners. */
export function installRuntimeRecordingStreamFixture() {
  const sources: FixtureEventSource[] = [];
  class FixtureEventSource extends EventTarget {
    static readonly CONNECTING = 0; static readonly OPEN = 1; static readonly CLOSED = 2;
    readonly CONNECTING = 0; readonly OPEN = 1; readonly CLOSED = 2;
    readonly withCredentials = false; readyState = 0; closed = false;
    onopen: ((event: Event) => void) | null = null;
    onerror: ((event: Event) => void) | null = null;
    onmessage: ((event: MessageEvent) => void) | null = null;
    readonly url: string;
    constructor(url: string | URL) { super(); this.url = String(url); sources.push(this); }
    close() { this.closed = true; this.readyState = 2; }
    send(name: string, data: unknown) {
      if (this.closed) return;
      if (name === 'open') { this.readyState = 1; const event = new Event('open'); this.onopen?.(event); this.dispatchEvent(event); return; }
      const event = new MessageEvent(name, { data: typeof data === 'string' ? data : JSON.stringify(data) });
      if (name === 'message') this.onmessage?.(event);
      this.dispatchEvent(event);
    }
    fail() { if (this.closed) return; this.readyState = 0; const event = new Event('error'); this.onerror?.(event); this.dispatchEvent(event); }
  }
  Object.defineProperty(globalThis, 'EventSource', { configurable: true, writable: true, value: FixtureEventSource });
  Object.defineProperty(window, 'EventSource', { configurable: true, writable: true, value: FixtureEventSource });
  window.__runtimeRecordingStreams = { emit(index, event, data) { sources[index]?.send(event, data); },
    fail(index) { sources[index]?.fail(); }, list() { return sources.map(source => ({ url: source.url, closed: source.closed })); } };
}
