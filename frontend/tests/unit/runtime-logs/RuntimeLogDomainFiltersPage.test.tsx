import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import RuntimeLogsPage from '../../../src/features/runtime-logs/RuntimeLogsPage';
import i18n from '../../../src/i18n/config';
import { record, snapshot as base } from './fixtures';
import { legacyRecord } from './legacy-event-fixtures';
import type { RuntimeLogSnapshot } from '../../../src/types/runtimeLogs';

// Actual Page -> hook/session -> transport/parser. Only HTTP responses are controlled.
const encoder = new TextEncoder();
type FetchOptions = NonNullable<Parameters<typeof fetch>[1]>;
const requests: { url: URL; init: FetchOptions }[] = [];
const streams: ReadableStreamDefaultController<Uint8Array>[] = [];
let data: RuntimeLogSnapshot;
let deny = false;
let respond: ((url: URL) => Promise<Response> | undefined) | undefined;
const response = (value: RuntimeLogSnapshot) => new Response(JSON.stringify(value), {
  headers: { 'Content-Type': 'application/json' },
});
const frame = (type: string, value: unknown) => encoder.encode(`event: ${type}\ndata: ${JSON.stringify(value)}\n\n`);
const row = (sequence: string) => document.querySelector(`[data-log-cursor="instance-a:${sequence}"]`);
const query = (index: number) => Object.fromEntries(requests[index].url.searchParams);
async function connected() {
  await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('Connected'));
}
function openScope() {
  const button = screen.getByRole('button', { name: /^Scope filters \(exact IDs\)/ });
  button.focus();
  expect(button).toHaveFocus();
  // Native button keyboard activation dispatches click; browser screenshots remain separate.
  fireEvent.click(button, { detail: 0 });
  expect(button).toHaveAttribute('aria-expanded', 'true');
}
function apply() { fireEvent.click(screen.getByRole('button', { name: 'Apply filters' })); }
function change(label: string, value: string) { fireEvent.change(screen.getByLabelText(label), { target: { value } }); }
beforeEach(async () => {
  await i18n.changeLanguage('en');
  requests.length = 0; streams.length = 0; deny = false; respond = undefined;
  data = { ...base, records: [], latest_cursor: 'instance-a:9007199254740993' };
  vi.stubGlobal('fetch', vi.fn(async (input: Parameters<typeof fetch>[0], init: FetchOptions = {}) => {
    const url = new URL(String(input), window.location.origin); requests.push({ url, init });
    expect(url.pathname).toMatch(/^\/api\/v1\/system\/logs(?:\/stream)?$/);
    expect(init.method ?? 'GET').toBe('GET');
    if (deny) return new Response('private server exception', { status: 403 });
    const override = respond?.(url); if (override) return override;
    if (!url.pathname.endsWith('/stream')) return response(data);
    return new Response(new ReadableStream<Uint8Array>({
      start(controller) {
        streams.push(controller);
        controller.enqueue(frame('handshake', { ...data, progress_cursor: data.latest_cursor }));
        init.signal?.addEventListener('abort', () => { try { controller.close(); } catch { /* already canceled */ } }, { once: true });
      },
    }), { headers: { 'Content-Type': 'text/event-stream' } });
  }));
});
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });

describe('domain filters on the actual runtime log page', () => {
  it.each([
    ['en', 'Evidence stage', 'Suggested action / unknown state', 'Scope filters (exact IDs)', 'Workspace ID', 'All stages'],
    ['zh-TW', '證據階段', '有處理建議／狀態不明', '範圍篩選（完整 ID）', '工作區 ID', '全部階段'],
  ])('provides labeled native keyboard controls and quiet light contrast in %s', async (language, stage, actionable, scope, workspace, all) => {
    await i18n.changeLanguage(language);
    render(<RuntimeLogsPage />);
    const select = screen.getByRole('combobox', { name: stage });
    expect(screen.getByRole('option', { name: all })).toHaveProperty('selected', true);
    expect(select).toHaveClass('bg-white', 'text-slate-900', 'border-slate-300', 'focus-visible:outline-blue-600');
    expect(screen.getByRole('main')).toHaveClass('bg-slate-50', 'text-slate-900');
    const box = screen.getByRole('checkbox', { name: actionable });
    expect(box).toHaveClass('accent-blue-600', 'focus-visible:outline-blue-600');
    const button = screen.getByRole('button', { name: scope });
    expect(button).toHaveAttribute('aria-expanded', 'false');
    expect(screen.queryByLabelText(workspace)).toBeNull();
    button.focus(); expect(button).toHaveFocus(); fireEvent.click(button, { detail: 0 });
    expect(button).toHaveAttribute('aria-expanded', 'true');
    expect(button).toHaveClass('ring-2', 'ring-blue-600', 'focus-visible:outline-blue-600');
    for (const label of language === 'en' ? ['Workspace ID', 'Device ID', 'Source rule ID', 'Recording group ID']
      : ['工作區 ID', '設備 ID', '來源規則 ID', '寫入群組 ID']) {
      expect(screen.getByLabelText(label)).toHaveAttribute('maxlength', '128');
    }
    const section = document.getElementById(button.getAttribute('aria-controls') ?? '');
    expect(section).toContainElement(screen.getByLabelText(workspace));
    screen.getByLabelText(workspace).focus(); expect(screen.getByLabelText(workspace)).toHaveFocus();
    fireEvent.click(button, { detail: 0 }); expect(button).toHaveAttribute('aria-expanded', 'false');
    await waitFor(() => expect(requests).toHaveLength(2));
  });

  it('applies all six typed filters only on submit with identical snapshot/stream scope and no lookup calls', async () => {
    render(<RuntimeLogsPage />); await connected(); openScope();
    change('Workspace ID', 'workspace-A'); change('Device ID', 'device-A');
    change('Source rule ID', 'rule-A'); change('Recording group ID', 'group-A');
    change('Evidence stage', 'retrying'); change('Minimum severity', 'warn');
    change('Source (exact ID)', 'runtime'); change('Search safe text', '溫度.*');
    fireEvent.click(screen.getByRole('checkbox', { name: 'Suggested action / unknown state' }));
    expect(screen.getByRole('checkbox', { name: 'Suggested action / unknown state' })).toBeChecked();
    fireEvent.click(screen.getByRole('button', { name: /^Scope filters \(exact IDs\)/ })); openScope();
    expect(screen.getByLabelText('Device ID')).toHaveValue('device-A');
    expect(requests).toHaveLength(2); apply(); await waitFor(() => expect(requests).toHaveLength(4));
    const shared = { level: 'warn', source: 'runtime', q: '溫度.*', workspace_id: 'workspace-A',
      device_id: 'device-A', rule_id: 'rule-A', group_id: 'group-A', stage: 'retrying', actionable: 'true' };
    expect(query(2)).toEqual({ ...shared, limit: '200' });
    expect(query(3)).toEqual({ ...shared, after: 'instance-a:9007199254740993' });
    expect(requests[1].init.signal?.aborted).toBe(true);
    expect(screen.getByText(/Historical events with a suggested action or unknown state/)).toBeInTheDocument();
    expect(screen.getByText(/not a list of current unresolved faults/)).toBeInTheDocument();
  });

  it('clears exact IDs without forcing a default and keeps same-filter pause/clear/resume progress', async () => {
    render(<RuntimeLogsPage />); await connected(); openScope(); change('Device ID', 'device-A'); apply();
    await waitFor(() => expect(requests).toHaveLength(4));
    fireEvent.click(screen.getByRole('button', { name: 'Pause view' }));
    fireEvent.click(screen.getByRole('button', { name: 'Clear view' }));
    expect(screen.getByRole('status')).toHaveTextContent('Paused'); expect(requests).toHaveLength(4);
    fireEvent.click(screen.getByRole('button', { name: 'Resume view' }));
    await waitFor(() => expect(requests).toHaveLength(5));
    expect(query(4)).toEqual({ level: 'debug', device_id: 'device-A', after: 'instance-a:9007199254740993' });
    change('Device ID', ''); change('Evidence stage', ''); apply();
    await waitFor(() => expect(requests).toHaveLength(7));
    expect(query(5)).toEqual({ level: 'debug', limit: '200' });
    expect(query(6)).toEqual({ level: 'debug', after: 'instance-a:9007199254740993' });
  });

  it('preserves the applied view when an explicit ID is invalid and shows safe inline format guidance', async () => {
    data = { ...data, records: [legacyRecord('http.request', { sequence: '31' })] };
    render(<RuntimeLogsPage />); await connected(); openScope();
    for (const [label, value] of [['Device ID', 'https://private.invalid/secret'], ['Workspace ID', '192.168.1.2'],
      ['Source rule ID', ' rule-A '], ['Recording group ID', '[::1]']]) {
      change(label, value); apply();
      expect(await screen.findByRole('alert')).toHaveTextContent('Invalid ID format. Enter a complete ID.');
      expect(screen.getByRole('alert')).not.toHaveTextContent(value);
      expect(screen.getByRole('status')).toHaveTextContent('Connected');
      expect(row('31')).not.toBeNull(); expect(row('31')).toHaveTextContent('instance-a:31');
      expect(requests).toHaveLength(2);
      change(label, '');
    }
    change('Device ID', 'device-valid'); apply();
    await waitFor(() => expect(requests).toHaveLength(4));
    expect(query(2).device_id).toBe('device-valid'); expect(screen.queryByRole('alert')).toBeNull();
  });

  it('shows the applied exact-ID count while collapsed, never the unsubmitted draft count', async () => {
    render(<RuntimeLogsPage />); await connected(); openScope();
    change('Device ID', 'device-A'); change('Recording group ID', 'group-A');
    expect(screen.queryByText('Applied IDs: 2')).toBeNull(); apply();
    await waitFor(() => expect(requests).toHaveLength(4));
    fireEvent.click(screen.getByRole('button', { name: /^Scope filters \(exact IDs\)/ }));
    const button = screen.getByRole('button', { name: /Scope filters \(exact IDs\).*Applied IDs: 2/ });
    expect(button).toHaveAttribute('aria-expanded', 'false');
    openScope(); change('Device ID', ''); change('Recording group ID', '');
    fireEvent.click(button); expect(button).toHaveTextContent('Applied IDs: 2');
    expect(requests).toHaveLength(4); apply();
    await waitFor(() => expect(requests).toHaveLength(6));
    expect(button).not.toHaveTextContent('Applied IDs:');
  });

  it('keeps old cursor and late device A data out of device B after a pending request', async () => {
    let finish: ((value: Response) => void) | undefined;
    respond = (url) => url.searchParams.get('device_id') === 'device-A' && !url.pathname.endsWith('/stream')
      ? new Promise<Response>((resolve) => { finish = resolve; }) : undefined;
    render(<RuntimeLogsPage />); await connected(); openScope(); change('Device ID', 'device-A'); apply();
    await waitFor(() => expect(requests).toHaveLength(3));
    data = { ...base, latest_cursor: 'instance-a:7', records: [legacyRecord('runtime.started', { sequence: '7' })] };
    change('Device ID', 'device-B'); apply(); await waitFor(() => expect(row('7')).not.toBeNull());
    await waitFor(() => expect(requests).toHaveLength(5));
    expect(query(4)).toEqual({ level: 'debug', device_id: 'device-B', after: 'instance-a:7' });
    expect(requests[2].init.signal?.aborted).toBe(true);
    if (!finish) throw new Error('device A HTTP request was not captured');
    await act(async () => finish?.(response({ ...base, latest_cursor: 'instance-a:9007199254740993', records: [legacyRecord('runtime.started', { sequence: '9007199254740993', message: 'Late A' })] })));
    expect(row('9007199254740993')).toBeNull(); expect(row('7')).not.toBeNull(); expect(screen.queryByText('Late A')).toBeNull(); expect(requests).toHaveLength(5);
  });

  it('keeps gap/reset and missing recovery explicitly unknown, including an empty actionable result', async () => {
    data = { ...data, capture_dropped: 2, records: [{ ...record('1', 'Retained acquisition failure'),
      code: 'domain.collector.failed', level: 'warn', fields: { stage: 'acquisition', action: 'refresh' } }] };
    render(<RuntimeLogsPage />); await connected();
    expect(screen.getByText(/Some capture history is missing/)).toBeInTheDocument();
    expect(screen.getByText(/Current fault or recovery status is unknown/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole('checkbox', { name: 'Suggested action / unknown state' })); apply();
    await waitFor(() => expect(requests).toHaveLength(4));
    await act(async () => streams[1].enqueue(frame('reset', { reason: 'reset' })));
    expect(screen.getByText(/The gateway process restarted/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Clear view' }));
    expect(screen.getByText('No matching retained records')).toBeInTheDocument();
    expect(screen.getByText(/Missing events or recovery records do not establish health/)).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Open runtime monitor' })).toHaveAttribute('href', '/studio/runtime');
    expect(screen.queryByText(/all healthy|no faults|resolved successfully/i)).toBeNull();
  });

  it('keeps the existing local-only 403 stop and explicit retry without disclosure-driven requests', async () => {
    deny = true; render(<RuntimeLogsPage />);
    expect(await screen.findByRole('alert')).toHaveTextContent('Local-only access');
    expect(screen.getByRole('status')).toHaveTextContent('Access denied');
    openScope(); expect(requests).toHaveLength(1);
    expect(screen.getByRole('button', { name: 'Pause view' })).toBeDisabled();
    expect(document.body).not.toHaveTextContent('private server exception');
    deny = false; fireEvent.click(screen.getByRole('button', { name: 'Retry connection' }));
    await connected(); expect(requests).toHaveLength(3);
    expect(requests.every(({ init }) => init.credentials === 'same-origin' && init.redirect === 'error')).toBe(true);
  });
});
