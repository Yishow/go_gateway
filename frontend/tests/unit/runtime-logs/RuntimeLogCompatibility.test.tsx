import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import RuntimeLogsPage from '@/features/runtime-logs/RuntimeLogsPage';
import { projectRuntimeLogLegacy } from '@/services/runtimeLogLegacyProjection';
import type { RuntimeLogRecord, RuntimeLogSnapshot } from '@/types/runtimeLogs';
import i18n from '@/i18n/config';
import { snapshot as base } from './fixtures';
import { domainRecord } from './domain-event-fixtures';
import { legacyCases, legacyInstant, legacyRecord, legacyRequestID } from './legacy-event-fixtures';

// Real Page -> session -> HTTP/SSE parser. Only the HTTP boundary is controlled.
const encoder = new TextEncoder();
const secret = 'password_fixture https://private.example/secret <script>alert(1)</script>';
const cursor = 'instance-a:9007199254740993';
type FetchOptions = NonNullable<Parameters<typeof fetch>[1]>;
const requests: { url: URL; init: FetchOptions }[] = [];
const unexpected: string[] = [];
const streams: ReadableStreamDefaultController<Uint8Array>[] = [];
let data: RuntimeLogSnapshot;
let raw: unknown;
let denial: 'snapshot' | 'stream' | undefined;
const frame = (type: string, value: unknown, id = '') => encoder.encode(
  `${id ? `id: ${id}\n` : ''}event: ${type}\ndata: ${JSON.stringify(value)}\n\n`,
);
function row(value = cursor) { return document.querySelector<HTMLTableRowElement>(`[data-log-cursor="${value}"]`); }
async function connected() { await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('Connected')); }
function detailsValue(element: HTMLElement, label: string) {
  const term = within(element).getByText(label, { selector: 'dt' });
  return term.nextElementSibling?.textContent;
}
beforeEach(async () => {
  await i18n.changeLanguage('en'); requests.length = 0; unexpected.length = 0; streams.length = 0;
  raw = undefined; denial = undefined;
  data = { ...base, records: [legacyRecord()], latest_cursor: cursor };
  vi.stubGlobal('fetch', vi.fn(async (input: Parameters<typeof fetch>[0], init: FetchOptions = {}) => {
    const url = new URL(String(input), window.location.origin); requests.push({ url, init });
    if (!/^\/api\/v1\/system\/logs(?:\/stream)?$/.test(url.pathname) || (init.method ?? 'GET') !== 'GET') {
      unexpected.push(`${init.method ?? 'GET'} ${url.pathname}`); throw new Error('Unexpected endpoint or mutation');
    }
    const stream = url.pathname.endsWith('/stream');
    if (denial === (stream ? 'stream' : 'snapshot')) return new Response(secret, { status: 403 });
    if (!stream) return new Response(JSON.stringify(raw ?? data), { headers: { 'Content-Type': 'application/json' } });
    return new Response(new ReadableStream<Uint8Array>({ start(controller) {
      streams.push(controller); controller.enqueue(frame('handshake', { ...data, progress_cursor: data.latest_cursor }));
      init.signal?.addEventListener('abort', () => { try { controller.close(); } catch { /* Already closed by cancellation. */ } }, { once: true });
    } }), { headers: { 'Content-Type': 'text/event-stream' } });
  }));
});
afterEach(() => { cleanup(); vi.useRealTimers(); vi.unstubAllGlobals(); expect(unexpected).toEqual([]); });

describe('closed legacy and domain compatibility on the actual runtime log page', () => {
  it.each(legacyCases)('renders fixed %s meaning without trusting candidate text', async (code) => {
    data.records = [legacyRecord(code, { message: secret })];
    const view = render(<RuntimeLogsPage />); await connected();
    expect(document.body).not.toHaveTextContent(secret);
    expect(view.container.querySelector('script,img')).toBeNull();
    const current = row(); expect(current).not.toBeNull(); if (!current) throw new Error('Missing exact safe row');
    expect(within(current).getByText(code)).toBeInTheDocument();
    if (code === 'runtime.started') expect(current).toHaveTextContent('does not prove successful acquisition, intake or SQL commitment');
    if (code === 'startup.ready') expect(current).toHaveTextContent('acquisition and SQL commitment remain separate');
    if (code === 'http.access' || code === 'http.request') expect(current).toHaveTextContent('HTTP success does not prove acquisition, intake or SQL commitment');
    if (code === 'http.recovered' || code === 'http.recovery') expect(current).toHaveTextContent('not acquisition recovery or SQL commitment');
    expect(current.querySelector('details > summary')).toHaveTextContent('Diagnostic details');
    expect(detailsValue(current, 'Log reference')).toBe(cursor);
    expect(current).toHaveTextContent(legacyInstant);
    expect(current.querySelector('a,button')).toBeNull();
    expect(requests).toHaveLength(2);
  });

  it('keeps safe HTTP scalars and exact zero but never reflects a registered-looking route', async () => {
    data.records = [legacyRecord('http.access', { fields: { route: '/api/v1/private/password_fixture', duration_ms: 0 } })];
    render(<RuntimeLogsPage />); await connected();
    const current = row(); if (!current) throw new Error('Missing HTTP row');
    expect(detailsValue(current, 'HTTP method')).toBe('GET'); expect(detailsValue(current, 'HTTP status')).toBe('200');
    expect(detailsValue(current, 'Duration (ms)')).toBe('0'); expect(detailsValue(current, 'Request ID')).toBe(legacyRequestID);
    expect(current).toHaveTextContent('Only a safe request summary is shown.');
    expect(document.body).not.toHaveTextContent('/api/v1/private/password_fixture');
    expect(document.body).not.toHaveTextContent('"duration_ms"');
  });

  it('shows absent HTTP values as unavailable, not zero or guessed success', async () => {
    const value = legacyRecord('http.request'); value.fields = {}; data.records = [value];
    render(<RuntimeLogsPage />); await connected(); const current = row(); if (!current) throw new Error('Missing row');
    for (const label of ['HTTP method', 'HTTP status', 'Duration (ms)', 'Request ID']) {
      expect(detailsValue(current, label)).toBe('Unavailable');
    }
  });

  type HostileCase = { name: string; overrides: Partial<RuntimeLogRecord> };
  const hostile: HostileCase[] = [
    { name: 'unknown code', overrides: { code: 'private.password_fixture' } },
    { name: 'unsafe outer source', overrides: { source: 'postgres://private.example/password_fixture' } },
    { name: 'unsafe timestamp suffix', overrides: { timestamp: '2026-10-29T10:00:00Z (password_fixture)' } },
    { name: 'arbitrary scalar field', overrides: { fields: { raw_error: secret } } },
    { name: 'invalid request identity', overrides: { fields: { request_id: 'password_fixture' } } },
    { name: 'contradictory level', overrides: { level: 'debug' } },
  ];
  it.each(hostile)('suppresses $name across all outer cells without an action', async ({ overrides }) => {
    data.records = [legacyRecord('http.request', { message: secret, ...overrides })];
    render(<RuntimeLogsPage />); await waitFor(() => expect(requests.length).toBeGreaterThan(0));
    await waitFor(() => expect(screen.getByRole('status').textContent).not.toBe('Loading'));
    expect(document.body).not.toHaveTextContent(/password_fixture|postgres:\/\/|private\.example/);
    const table = screen.getByRole('table'); expect(table.querySelector('a,button')).toBeNull();
    expect(table).not.toHaveTextContent('HTTP request completed');
  });

  it('rejects a nested credential payload at the real transport parser with no fallback feed', async () => {
    raw = { ...data, records: [{ ...legacyRecord('http.request'), fields: { credential: { password: 'password_fixture' } } }] };
    render(<RuntimeLogsPage />); await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('Failed'));
    expect(requests).toHaveLength(1); expect(document.body).not.toHaveTextContent('password_fixture');
    expect(screen.queryByText('No matching retained records')).toBeNull();
  });

  it('never invokes a candidate getter at the presentation boundary (descriptors cannot travel over JSON)', () => {
    const value = legacyRecord(); const getter = vi.fn(() => secret);
    Object.defineProperty(value, 'fields', { enumerable: true, get: getter });
    expect(projectRuntimeLogLegacy(value).kind).toBe('suppressed'); expect(getter).not.toHaveBeenCalled();
  });

  it('mixes legacy and domain rows on snapshot and SSE without confusing log connection with delivery evidence', async () => {
    data.records = [legacyRecord('runtime.started'), domainRecord('domain.delivery.blocked', { sequence: '9007199254740992' })];
    render(<RuntimeLogsPage />); await connected();
    expect(screen.getByRole('status')).toHaveTextContent('Connected · log feed only');
    expect(screen.getByText('Delivery blocked')).toBeInTheDocument();
    expect(screen.getByText(/Connection status describes this log feed only, not acquisition, recording intake, SQL commitment/i)).toBeInTheDocument();
    const recovered = legacyRecord('http.recovered', { sequence: '9007199254740994', message: secret });
    await act(async () => streams[0].enqueue(frame('log', recovered, 'instance-a:9007199254740994')));
    expect(row('instance-a:9007199254740994')).not.toBeNull();
    expect(screen.getByText('Delivery blocked')).toBeInTheDocument();
    expect(document.body).not.toHaveTextContent(secret);
    expect(screen.queryByText(/all healthy|SQL confirmed|acquisition recovered/i)).toBeNull();
    expect(requests).toHaveLength(2);
  });

  it.each(['snapshot', 'stream'] as const)('stops automatic retry after %s 403 and only retries on explicit request', async (where) => {
    denial = where; vi.useFakeTimers(); render(<RuntimeLogsPage />);
    await act(async () => { await vi.advanceTimersByTimeAsync(0); });
    expect(screen.getByRole('status')).toHaveTextContent('Access denied');
    const deniedCount = where === 'snapshot' ? 1 : 2; expect(requests).toHaveLength(deniedCount);
    await act(async () => { await vi.advanceTimersByTimeAsync(120000); });
    expect(requests).toHaveLength(deniedCount); expect(document.body).not.toHaveTextContent(secret);
    expect(screen.getByRole('alert')).toHaveTextContent('Automatic retries have stopped');
    denial = undefined; fireEvent.click(screen.getByRole('button', { name: 'Retry connection' }));
    await act(async () => { await vi.advanceTimersByTimeAsync(0); });
    expect(screen.getByRole('status')).toHaveTextContent('Connected'); expect(requests).toHaveLength(deniedCount + 2);
    expect(requests.every(({ init }) => init.credentials === 'same-origin' && init.redirect === 'error')).toBe(true);
  });

  it('pauses and clears only browser rows, then resumes the exact unrounded progress without snapshot or mutation', async () => {
    render(<RuntimeLogsPage />); await connected(); expect(row()).not.toBeNull();
    fireEvent.click(screen.getByRole('button', { name: 'Pause view' }));
    expect(requests[1].init.signal?.aborted).toBe(true);
    fireEvent.click(screen.getByRole('button', { name: 'Clear view' }));
    expect(row()).toBeNull(); expect(requests).toHaveLength(2);
    expect(screen.getByText(/Pause and clear affect this browser only/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Resume view' })); await connected();
    expect(requests).toHaveLength(3); expect(requests[2].url.pathname).toBe('/api/v1/system/logs/stream');
    expect(requests[2].url.searchParams.get('after')).toBe(cursor); expect(data.records).toHaveLength(1);
    const next = legacyRecord('shutdown.complete', { sequence: '9007199254740994' });
    await act(async () => streams[1].enqueue(frame('log', next, 'instance-a:9007199254740994')));
    expect(row('instance-a:9007199254740994')).not.toBeNull(); expect(row()).toBeNull();
  });

  it('keeps file loss, capture gaps, retention eviction and reset separate from health even with no records', async () => {
    data = { ...data, records: [], evicted: 1, file_dropped: 1, sink_health: 'degraded' };
    render(<RuntimeLogsPage />); await connected();
    expect(screen.getByText(/Current fault or recovery status is unknown/)).toBeInTheDocument();
    expect(screen.getByText(/Diagnostic files have missing history/)).toBeInTheDocument();
    expect(screen.queryByText(/Some capture history is missing/)).toBeNull();
    await act(async () => streams[0].enqueue(frame('gap', { reason: 'retention' })));
    expect(screen.getByText(/Some capture history is missing/)).toBeInTheDocument();
    await act(async () => streams[0].enqueue(frame('reset', { reason: 'reset' })));
    expect(screen.getByText(/The gateway process restarted/)).toBeInTheDocument();
    expect(screen.getByText('No matching retained records')).toBeInTheDocument();
    expect(screen.getByText(/Missing events or recovery records do not establish health/)).toBeInTheDocument();
    expect(screen.queryByText(/all healthy|no faults|resolved successfully/i)).toBeNull(); expect(requests).toHaveLength(2);
  });

  it.each([
    { language: 'en', summary: 'Diagnostic details', pause: 'Pause view' },
    { language: 'zh-TW', summary: '診斷詳細資料', pause: '暫停檢視' },
  ])('provides native focusable localized details and controls in $language', async ({ language, summary, pause }) => {
    await i18n.changeLanguage(language); data.records = [legacyRecord('http.request')];
    render(<RuntimeLogsPage />); await waitFor(() => expect(row()).not.toBeNull());
    const current = row(); if (!current) throw new Error('Missing exact row');
    const disclosure = within(current).getByText(summary, { selector: 'summary' });
    disclosure.focus(); expect(disclosure).toHaveFocus(); fireEvent.click(disclosure, { detail: 0 });
    expect(disclosure.closest('details')).toHaveAttribute('open');
    const button = screen.getByRole('button', { name: pause }); button.focus(); expect(button).toHaveFocus();
    fireEvent.click(button, { detail: 0 }); expect(requests).toHaveLength(2);
    expect(current.querySelector('a')).toBeNull();
    // Native browser keyboard-default behavior and screenshots remain NOT RUN.
  });
});
