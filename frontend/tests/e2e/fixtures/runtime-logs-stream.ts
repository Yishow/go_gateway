import { RUNTIME_LOGS_PATH } from './runtime-logs-state';
export interface RuntimeLogsStreamRequest { url: string; aborted: boolean; closed: boolean }
export interface RuntimeLogsStreamController {
  readonly requests: RuntimeLogsStreamRequest[];
  send(index: number, rawFrame: string): boolean;
  end(index: number): void;
  fail(index: number): void;
  denyNext(status: number): void;
  restore(): void;
}
declare global { interface Window { __runtimeLogsStream: RuntimeLogsStreamController } }

/** Test-only fetch boundary: real response/readable bytes feed the unchanged production parser/session. */
export function installRuntimeLogsStream(): RuntimeLogsStreamController {
  const original = globalThis.fetch;
  const requests: RuntimeLogsStreamRequest[] = [];
  const pending: { controller: ReadableStreamDefaultController<Uint8Array> | null; removeAbort: () => void }[] = [];
  let nextStatus: number | undefined; let restored = false;
  const close = (index: number, error?: Error) => {
    const request = requests[index], body = pending[index];
    if (!request || !body || request.closed) return;
    request.closed = true;
    if (error) body.controller?.error(error); else body.controller?.close();
  };
  const intercepted: typeof fetch = async (input, init) => {
    const address = input instanceof Request ? input.url : String(input);
    const url = new URL(address, globalThis.location?.href ?? 'http://runtime-logs-fixture.invalid');
    if (url.pathname !== `${RUNTIME_LOGS_PATH}/stream`) return original(input, init);
    const signal = init?.signal ?? (input instanceof Request ? input.signal : undefined);
    const index = requests.length;
    requests.push({ url: url.href, aborted: signal?.aborted ?? false, closed: false });
    const body = { controller: null as ReadableStreamDefaultController<Uint8Array> | null, removeAbort: () => undefined as void }; pending.push(body);
    const abort = () => { requests[index].aborted = true; close(index, new DOMException('Aborted', 'AbortError')); };
    const method = (init?.method ?? (input instanceof Request ? input.method : 'GET')).toUpperCase();
    if (method !== 'GET' || init?.body != null) { requests[index].closed = true; return new Response(null, { status: method !== 'GET' ? 405 : 400 }); }
    if (nextStatus !== undefined) {
      const status = nextStatus; nextStatus = undefined; requests[index].closed = true;
      return new Response(null, { status });
    }
    const bytes = new ReadableStream<Uint8Array>({
      start(controller) {
        body.controller = controller;
        signal?.addEventListener('abort', abort, { once: true });
        body.removeAbort = () => signal?.removeEventListener('abort', abort);
        if (signal?.aborted) abort();
      },
      cancel() { requests[index].closed = true; body.removeAbort(); },
    });
    return new Response(bytes, { headers: { 'content-type': 'text/event-stream' } });
  };
  globalThis.fetch = intercepted;
  return { requests,
    send(index, rawFrame) { const body = pending[index]; if (restored || !body?.controller || requests[index].closed) return false;
      body.controller.enqueue(new TextEncoder().encode(rawFrame)); return true; },
    end(index) { close(index); }, fail(index) { close(index, new TypeError('Fixture stream offline')); },
    denyNext(status) { if (!Number.isInteger(status) || status < 400 || status > 599) throw new Error('Expected HTTP failure status.'); nextStatus = status; },
    restore() { if (restored) return; restored = true; requests.forEach((_request, index) => { close(index); pending[index].removeAbort(); });
      if (globalThis.fetch === intercepted) globalThis.fetch = original; },
  };
}
