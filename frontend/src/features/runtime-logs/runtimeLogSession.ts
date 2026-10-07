import type { RuntimeLogEvent, RuntimeLogFilters, RuntimeLogMetadata, RuntimeLogRecord, RuntimeLogTransport } from '../../types/runtimeLogs';
import { RuntimeLogRequestError, runtimeLogTransport } from '../../services/runtimeLogs';
import { appendRuntimeLogs } from './runtimeLogBuffer';
export type RuntimeLogStatus = 'loading' | 'connected' | 'reconnecting' | 'paused' | 'denied' | 'failed';
export interface RuntimeLogView {
  records: RuntimeLogRecord[];
  status: RuntimeLogStatus;
  metadata: RuntimeLogMetadata | null;
  gap: boolean;
  fileGap: boolean;
  reset: boolean;
  truncated: boolean;
  lastUpdate: string | null;
}
const backoff = [1000, 2000, 4000, 8000, 16000, 30000];
function alreadyConsumed(candidate: string, consumed: string) {
  const a = candidate.split(':'); const b = consumed.split(':');
  return a[0] === b[0] && /^\d+$/.test(a[1]) && /^\d+$/.test(b[1]) && BigInt(a[1]) <= BigInt(b[1]);
}
/** Single owner for fetch cancellation, processed progress and reconnect timers. */
export class RuntimeLogSession {
  private state: RuntimeLogView = { records: [], status: 'loading', metadata: null,
    gap: false, fileGap: false, reset: false, truncated: false, lastUpdate: null };
  private listeners = new Set<() => void>();
  private controller: AbortController | null = null;
  private timer: ReturnType<typeof setTimeout> | null = null;
  private generation = 0;
  private active = false;
  private loaded = false;
  private attempt = 0;
  private cursor = '';
  private lastRecordCursor = '';
  constructor(private filters: RuntimeLogFilters, private transport: RuntimeLogTransport = runtimeLogTransport,
    private random: () => number = Math.random) {}
  getSnapshot = () => this.state;
  subscribe = (listener: () => void) => { this.listeners.add(listener); return () => { this.listeners.delete(listener); }; };
  private patch(update: Partial<RuntimeLogView>) {
    this.state = { ...this.state, ...update }; this.listeners.forEach((listener) => listener());
  }
  private cancel() {
    this.generation++; this.controller?.abort(); this.controller = null;
    if (this.timer !== null) clearTimeout(this.timer); this.timer = null;
  }
  start = () => { this.active = true; this.connect(); };
  stop = () => { this.active = false; this.cancel(); };
  pause = () => { this.stop(); this.patch({ status: 'paused' }); };
  resume = () => { if (this.state.status === 'paused') this.start(); };
  clear = () => { this.patch({ records: [], truncated: false }); };
  retry = () => { this.attempt = 0; this.loaded = false; this.cursor = ''; this.lastRecordCursor = ''; this.start(); };
  private append(records: RuntimeLogRecord[]) {
    const next = appendRuntimeLogs(this.state.records, records);
    this.patch({ records: next.records, truncated: this.state.truncated || next.truncated || records.some((r) => r.truncated) });
  }
  private consume(event: RuntimeLogEvent) {
    if (event.type === 'log') {
      if (alreadyConsumed(event.cursor, this.cursor) || alreadyConsumed(event.cursor, this.lastRecordCursor)) return;
      this.lastRecordCursor = event.cursor; this.append([event.record]);
    } else if (event.type === 'handshake' || event.type === 'heartbeat') {
      const prior = this.state.metadata;
      this.patch({ metadata: event.metadata,
        gap: this.state.gap || event.metadata.capture_dropped > (prior?.capture_dropped ?? 0) ||
          (event.metadata.subscriber_overflow ?? 0) > (prior?.subscriber_overflow ?? 0),
        fileGap: this.state.fileGap || (event.metadata.file_dropped ?? 0) > 0 || (event.metadata.file_recovered_gaps ?? 0) > 0 });
      if (event.type === 'handshake') { this.attempt = 0; this.patch({ status: 'connected' }); }
    } else if (event.type === 'reset') {
      this.lastRecordCursor = ''; this.patch({ reset: true, gap: true, records: [] });
    } else if (event.type === 'gap' && event.reason === 'file_loss') { this.patch({ fileGap: true }); }
    else { this.patch({ gap: true }); }
    if (event.cursor && !alreadyConsumed(event.cursor, this.cursor)) this.cursor = event.cursor;
    this.patch({ lastUpdate: new Date().toISOString() });
  }
  private connect() {
    this.cancel(); const generation = this.generation;
    const controller = new AbortController(); this.controller = controller;
    const current = () => this.active && generation === this.generation && !controller.signal.aborted;
    this.patch({ status: this.loaded ? 'reconnecting' : 'loading' });
    void (async () => {
      if (!this.loaded) {
        const result = await this.transport.snapshot(this.filters, controller.signal);
        if (!current()) return;
        this.cursor = result.latest_cursor; this.loaded = true;
        const { records, ...metadata } = result;
        this.patch({ metadata, records: [], gap: result.capture_dropped > 0 || (result.subscriber_overflow ?? 0) > 0,
          fileGap: (result.file_dropped ?? 0) > 0 || (result.file_recovered_gaps ?? 0) > 0, lastUpdate: new Date().toISOString() });
        this.append(records);
      }
      if (!current()) return;
      await this.transport.stream(this.filters, this.cursor, controller.signal, (event) => { if (current()) this.consume(event); });
      if (current()) throw new RuntimeLogRequestError(0);
    })().catch((error: unknown) => {
      if (!current()) return;
      controller.abort();
      if (error instanceof RuntimeLogRequestError && error.status === 403) {
        this.active = false; this.patch({ status: 'denied' }); return;
      }
      if (error instanceof RuntimeLogRequestError && error.status >= 400 && error.status < 500 && error.status !== 429) {
        this.active = false; this.patch({ status: 'failed' }); return;
      }
      this.patch({ status: 'reconnecting' });
      const base = backoff[Math.min(this.attempt++, backoff.length - 1)];
      const delay = Math.min(30000, Math.round(base * (0.8 + this.random() * 0.4)));
      this.timer = setTimeout(() => { this.timer = null; if (this.active) this.connect(); }, delay);
    });
  }
}
