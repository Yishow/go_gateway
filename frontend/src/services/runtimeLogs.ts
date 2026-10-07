import { z } from 'zod';
import { VITE_API_BASE_URL } from '../env';
import type { RuntimeLogEvent, RuntimeLogFilters, RuntimeLogTransport } from '../types/runtimeLogs';
const id = z.string().regex(/^[A-Za-z0-9_-]{1,128}$/);
const sequence = z.string().regex(/^\d{1,20}$/);
const cursor = z.string().max(160).refine((value) => value === '' || /^[A-Za-z0-9_-]{1,128}:\d{1,20}$/.test(value));
const count = z.number().finite().nonnegative();
const recordSchema = z.object({
  instance_id: id, sequence, timestamp: z.string().max(64).refine((value) => Number.isFinite(Date.parse(value))),
  level: z.enum(['debug', 'info', 'warn', 'error']), source: z.string().max(128), code: z.string().max(128),
  message: z.string().max(8192), fields: z.record(z.string().max(64),
    z.union([z.string().max(8192), z.number().finite(), z.boolean()])).nullish().transform((value) => value ?? {}),
  truncated: z.boolean(),
});
const metadataSchema = z.object({
  instance_id: id, oldest_cursor: cursor, latest_cursor: cursor, capture_dropped: count, evicted: count,
  subscriber_overflow: count.optional(), file_dropped: count.optional(), disk_errors: count.optional(),
  file_recovered_gaps: count.optional(), sink_health: z.enum(['healthy', 'degraded', 'disabled']),
});
const snapshotSchema = metadataSchema.extend({ records: z.array(recordSchema).max(500).nullable().transform((value) => value ?? []) });
const progressSchema = metadataSchema.extend({ progress_cursor: cursor });
const gapSchema = z.object({ reason: z.string().max(64), from: cursor.optional(), to: cursor.optional() });
const encoder = new TextEncoder();
const MAX_RESPONSE_BYTES = 5 * 1024 * 1024;
const MAX_FRAME_BYTES = 16 * 1024;
export class RuntimeLogRequestError extends Error {
  constructor(public readonly status: number) { super('runtime_log_request_failed'); }
}
function safeParse<S extends z.ZodTypeAny>(schema: S, value: unknown): z.output<S> {
  const result = schema.safeParse(value);
  if (!result.success) throw new RuntimeLogRequestError(422);
  return result.data;
}
function json(text: string): unknown {
  try { return JSON.parse(text); } catch { throw new RuntimeLogRequestError(422); }
}
export function runtimeLogsURL(filters: RuntimeLogFilters, stream = false, after = '') {
  const params = new URLSearchParams({ level: filters.level });
  if (filters.source) params.set('source', filters.source);
  if (filters.q) params.set('q', filters.q);
  if (stream) { if (after) params.set('after', after); } else params.set('limit', '200');
  return `${VITE_API_BASE_URL.replace(/\/$/, '')}/system/logs${stream ? '/stream' : ''}?${params}`;
}
export function parseRuntimeLogFrame(frame: string): RuntimeLogEvent | null {
  if (encoder.encode(frame).byteLength > MAX_FRAME_BYTES) throw new RuntimeLogRequestError(422);
  let type = ''; let eventID = ''; const data: string[] = [];
  for (const line of frame.split(/\r?\n/)) {
    if (line.startsWith('event:')) type = line.slice(6).trim();
    if (line.startsWith('id:')) eventID = line.slice(3).trim();
    if (line.startsWith('data:')) data.push(line.slice(5).trimStart());
  }
  if (!data.length) return null;
  const value = json(data.join('\n'));
  if (type === 'log') {
    const record = safeParse(recordSchema, value);
    if (encoder.encode(JSON.stringify(value)).byteLength > 8192 || eventID !== `${record.instance_id}:${record.sequence}`) {
      throw new RuntimeLogRequestError(422);
    }
    return { type, record, cursor: eventID };
  }
  if (type === 'handshake' || type === 'heartbeat') {
    const metadata = safeParse(progressSchema, value);
    return { type, metadata, cursor: metadata.progress_cursor };
  }
  if (type === 'gap' || type === 'reset') {
    const gap = safeParse(gapSchema, value);
    // Loss boundaries describe missing history, never acknowledgement of unsent replay.
    return { type, reason: gap.reason, cursor: '' };
  }
  throw new RuntimeLogRequestError(422);
}
async function request<T>(signal: AbortSignal, work: (signal: AbortSignal, touch: () => void) => Promise<T>) {
  const controller = new AbortController(); let timer: ReturnType<typeof setTimeout>;
  const touch = () => { clearTimeout(timer); timer = setTimeout(() => controller.abort(), 45000); };
  const abort = () => controller.abort();
  signal.addEventListener('abort', abort, { once: true });
  if (signal.aborted) controller.abort();
  touch();
  try { return await work(controller.signal, touch); }
  finally { clearTimeout(timer!); signal.removeEventListener('abort', abort); }
}
async function readerFor(url: string, signal: AbortSignal, stream = false) {
  const response = await fetch(url, { signal, credentials: 'same-origin', redirect: 'error', cache: 'no-store',
    headers: { Accept: stream ? 'text/event-stream' : 'application/json' } });
  if (!response.ok) {
    void response.body?.cancel().catch(() => undefined); throw new RuntimeLogRequestError(response.status);
  }
  if (!response.body || !(response.headers.get('content-type') ?? '').includes(stream ? 'text/event-stream' : 'application/json')) {
    void response.body?.cancel().catch(() => undefined); throw new RuntimeLogRequestError(422);
  }
  return response.body.getReader();
}
export const runtimeLogTransport: RuntimeLogTransport = {
  snapshot: (filters, signal) => request(signal, async (requestSignal) => {
    const reader = await readerFor(runtimeLogsURL(filters), requestSignal);
    const decoder = new TextDecoder(); let text = ''; let bytes = 0;
    try {
      while (true) {
        const { done, value } = await reader.read(); if (done) break;
        bytes += value.byteLength;
        if (bytes > MAX_RESPONSE_BYTES) throw new RuntimeLogRequestError(422);
        text += decoder.decode(value, { stream: true });
      }
      const result = safeParse(snapshotSchema, json(text + decoder.decode()));
      if (result.records.some((record) => encoder.encode(JSON.stringify(record)).byteLength > 8192)) throw new RuntimeLogRequestError(422);
      return result;
    } finally { await reader.cancel().catch(() => undefined); reader.releaseLock(); }
  }),
  stream: (filters, after, signal, onEvent) => request(signal, async (requestSignal, touch) => {
    const reader = await readerFor(runtimeLogsURL(filters, true, after), requestSignal, true);
    const decoder = new TextDecoder(); let pending = '';
    try {
      while (!requestSignal.aborted) {
        const { done, value } = await reader.read(); if (done) break;
        if (value.byteLength > MAX_RESPONSE_BYTES) throw new RuntimeLogRequestError(422);
        pending += decoder.decode(value, { stream: true });
        let boundary: RegExpExecArray | null;
        while ((boundary = /\r?\n\r?\n/.exec(pending))) {
          const event = parseRuntimeLogFrame(pending.slice(0, boundary.index));
          pending = pending.slice(boundary.index + boundary[0].length);
          if (event) { onEvent(event); touch(); }
          if (requestSignal.aborted) return;
        }
        if (encoder.encode(pending).byteLength > MAX_FRAME_BYTES) throw new RuntimeLogRequestError(422);
      }
    } finally { await reader.cancel().catch(() => undefined); reader.releaseLock(); }
  }),
};
