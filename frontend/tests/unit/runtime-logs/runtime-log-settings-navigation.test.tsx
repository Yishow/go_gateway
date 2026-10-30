import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { MemoryRouter, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import RuntimeLogsPage from '@/features/runtime-logs/RuntimeLogsPage';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { parseStudioEntityLocator } from '@/services/studioEntityLocator';
import type { RuntimeLogRecord } from '@/types/runtimeLogs';
import i18n from '@/i18n/config';
import { locatorRead } from '../workbench-v2/studio-locator-fixtures';
import { domainRecord } from './domain-event-fixtures';
import { snapshot } from './fixtures';

const clients: QueryClient[] = [];
const originalAdapter = studioV2DatalinkApi.defaults.adapter;
type Read = ReturnType<typeof locatorRead>;
type EntryHandler = (params: URLSearchParams, config: InternalAxiosRequestConfig) => Promise<Read>;
let ledger: { method: string; path: string; params: URLSearchParams; signal: unknown }[];
let unexpected: string[];
function native(group = 'G', revision = 'R2', name = 'Renamed current recorder'): Read {
  const value = locatorRead(4, group); const saved = value.entry.write_groups?.groups[0];
  if (!saved) throw new Error('Missing native group fixture');
  saved.revision = revision; saved.name = name; value.entry.group_heads = { [group]: revision };
  return value;
}
function failure(config: InternalAxiosRequestConfig, status: number): never {
  throw new AxiosError('Safe fixture failure', 'ERR_BAD_RESPONSE', config, undefined,
    { data: { success: false, error: { code: 'unavailable', message: 'private.example/password_fixture' } },
      status, statusText: 'Unavailable', headers: {}, config });
}
function transport(records: RuntimeLogRecord[], entry: EntryHandler = async () => native()) {
  studioV2DatalinkApi.defaults.adapter = async config => {
    const path = config.url ?? ''; const method = (config.method ?? 'get').toUpperCase();
    const params = new URLSearchParams(config.params);
    ledger.push({ method, path, params, signal: config.signal });
    if (method !== 'GET') { unexpected.push(`${method} ${path}`); throw new Error('Unexpected mutation'); }
    if (path === '/studio-v2/workspace' && params.get('entry') === 'locator' && !params.has('editor')) {
      const value = await entry(params, config);
      return { data: { success: true, data: { ...value.workspace, locator_entry: value.entry } }, status: 200, statusText: 'OK', headers: {}, config };
    }
    if (/^\/studio-v2\/workspace\/write-groups\/(G|G1|G2|B)\/delivery$/.test(path)) failure(config, 503);
    unexpected.push(`${method} ${path}`); throw new Error('Unexpected GET or editor hydration');
  };
  const latest = records.at(-1)?.sequence ?? '0';
  const fetch = vi.fn<typeof globalThis.fetch>(async (input, options) => {
    const url = new URL(String(input), 'http://localhost');
    if (options?.method && options.method !== 'GET') { unexpected.push(String(input)); throw new Error('Unexpected log mutation'); }
    if (url.pathname === '/api/v1/system/logs') return new Response(JSON.stringify({ ...snapshot,
      latest_cursor: `instance-a:${latest}`, records }), { headers: { 'content-type': 'application/json' } });
    if (url.pathname === '/api/v1/system/logs/stream') return new Response(new ReadableStream<Uint8Array>({ start(controller) {
      controller.enqueue(new TextEncoder().encode(`event: handshake\ndata: ${JSON.stringify({ ...snapshot, progress_cursor: `instance-a:${latest}` })}\n\n`));
    } }), { headers: { 'content-type': 'text/event-stream' } });
    unexpected.push(url.pathname); throw new Error('Unexpected fetch');
  });
  vi.stubGlobal('fetch', fetch); return fetch;
}
function Receiver() {
  const location = useLocation(), navigate = useNavigate();
  return <DatalinkWorkbenchV2Page navigateTo={navigate} locatorIntent={parseStudioEntityLocator(location.search)} />;
}
function Host() {
  const location = useLocation(), navigate = useNavigate();
  return <><nav aria-label="Test history"><button onClick={() => navigate(-1)}>Back</button><button onClick={() => navigate(1)}>Forward</button></nav>
    <output data-testid="location">{location.pathname}{location.search}</output>
    <Routes><Route path="/studio/logs" element={<RuntimeLogsPage />} /><Route path="/studio/v2" element={<Receiver />} /></Routes></>;
}
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  return render(<QueryClientProvider client={client}><MemoryRouter initialEntries={['/studio/logs']}><Host /></MemoryRouter></QueryClientProvider>);
}
async function openGroup(index = 0) {
  const links = await screen.findAllByRole('link', { name: 'Open group settings' }); const link = links[index];
  if (!link) throw new Error('Missing exact event navigation'); fireEvent.click(link);
}
function expectNoCommands() {
  expect(unexpected).toEqual([]); expect(ledger.every(request => request.method === 'GET')).toBe(true);
  expect(ledger.some(request => request.params.has('editor'))).toBe(false);
  expect(screen.queryByTestId('group-editor')).not.toBeInTheDocument();
  expect(screen.queryByTestId('basic-recording-panel')).not.toBeInTheDocument();
  expect(screen.queryByTestId('workspace-test-write-recovery')).not.toBeInTheDocument();
}
beforeEach(async () => { ledger = []; unexpected = []; await i18n.changeLanguage('en'); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); studioV2DatalinkApi.defaults.adapter = originalAdapter; vi.unstubAllGlobals(); });

describe('real log Page and native settings receiver with strict HTTP-only fixtures', () => {
  it('spec W/G/R1 opens exact current G/R2 without mutation or historical relabeling', async () => {
    const event = domainRecord(); transport([event]); mount();
    await screen.findByText('Delivery blocked'); expect(ledger).toEqual([]);
    expect(screen.queryByText('Renamed current recorder')).not.toBeInTheDocument();
    await openGroup();
    const card = await screen.findByTestId('group-list-item-G'); expect(card).toHaveAttribute('data-locator-selected', 'true');
    expect(card).toHaveFocus(); expect(card).toHaveTextContent('Renamed current recorder');
    expect(screen.getByTestId('location')).toHaveTextContent('workspace_id=W&step=4&group_id=G&group_revision=R1');
    expect(screen.getByText(/Requested revision R1 is informational/)).toBeVisible();
    expect(screen.getByText(/Current settings revision/)).toBeVisible(); expect(screen.getByText('R2')).toBeVisible();
    await waitFor(() => expect(ledger).toHaveLength(2)); expectNoCommands();
    expect(ledger[0]?.params.get('group_revision')).toBe('R1');
    expect(event.fields.group_revision).toBe('R1'); expect(event.fields.entity_revision).toBe('R1');
    expect(document.body.textContent).not.toContain('private.example');
  });
  it('Back re-reads the current receiver, while immutable event R1 remains unchanged after rename', async () => {
    let count = 0; const fetch = transport([domainRecord()], async () => native('G', ++count === 1 ? 'R2' : 'R3', 'Current renamed recorder'));
    mount(); await openGroup(); await screen.findByTestId('group-list-item-G');
    fireEvent.click(screen.getByRole('button', { name: 'Back' })); await screen.findByText('Delivery blocked');
    expect(screen.queryByText('Current renamed recorder')).not.toBeInTheDocument();
    const details = screen.getByText('Event details'); fireEvent.click(details);
    expect(screen.getAllByText('R1').length).toBeGreaterThan(0);
    fireEvent.click(screen.getByRole('button', { name: 'Forward' })); await screen.findByTestId('group-list-item-G');
    await waitFor(() => expect(count).toBe(2));
    expect(ledger.filter(request => request.path === '/studio-v2/workspace').every(request => request.params.get('group_id') === 'G' && request.params.get('group_revision') === 'R1')).toBe(true);
    expect(fetch.mock.calls.filter(([input]) => String(input).includes('/logs?'))).toHaveLength(2); expectNoCommands();
  });
  it('spec deleted G1 never selects active same-name G2 or bootstraps another workspace', async () => {
    const read = native('G1', 'R2', 'Shared recorder'); const group = read.entry.write_groups?.groups[0];
    if (!group) throw new Error('Missing group'); group.status = 'deleted';
    transport([domainRecord('domain.delivery.blocked', { fields: { group_id: 'G1' } })], async params =>
      params.get('group_id') === 'G1' ? read : native('G2', 'R2', 'Shared recorder'));
    mount(); await openGroup();
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).toHaveTextContent('removed'));
    expect(screen.queryByTestId('group-list-item-G1')).not.toBeInTheDocument(); expect(screen.queryByTestId('group-list-item-G2')).not.toBeInTheDocument();
    expect(ledger).toHaveLength(1); expect(ledger[0]?.params.get('group_id')).toBe('G1'); expectNoCommands();
    const links = within(screen.getByTestId('studio-locator-notice')).getAllByRole('link');
    expect(links.every(link => link.getAttribute('href')?.includes('workspace_id=W&step=1'))).toBe(true);
  });
  it('rejects a same-name G2 substitute in the exact G1 response', async () => {
    transport([domainRecord('domain.delivery.blocked', { fields: { group_id: 'G1' } })], async () => native('G2', 'R2', 'Shared recorder'));
    mount(); await openGroup(); await screen.findByTestId('studio-locator-notice');
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).not.toHaveTextContent('Checking'));
    expect(screen.queryByTestId('group-list-item-G2')).not.toBeInTheDocument(); expectNoCommands();
  });
  it.each([403, 409, 503])('safe %i failure preserves historical identity without any fallback or command', async status => {
    transport([domainRecord()], async (_params, config) => failure(config, status)); mount(); await openGroup();
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).not.toHaveTextContent('Checking'));
    expect(ledger).toHaveLength(1); expect(screen.queryByTestId('group-list-item-G')).not.toBeInTheDocument(); expectNoCommands();
    expect(document.body.textContent).not.toContain('private.example');
    fireEvent.click(screen.getByRole('button', { name: 'Back' })); await screen.findByText('Delivery blocked');
    fireEvent.click(screen.getByText('Event details')); expect(screen.getAllByText('R1').length).toBeGreaterThan(0);
  });
  it('foreign workspace payload fails closed without substituting owned G by name', async () => {
    const foreign = native(); foreign.workspace.id = 'OTHER';
    transport([domainRecord()], async () => foreign); mount(); await openGroup();
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).toHaveTextContent('inaccessible'));
    expect(ledger).toHaveLength(1); expectNoCommands();
  });
  it('late A entry cannot replace B after Back and another event navigation', async () => {
    let finishA: ((read: Read) => void) | undefined;
    transport([domainRecord(), domainRecord('domain.delivery.blocked', { sequence: '9007199254740994', fields: { group_id: 'B' } })],
      async params => params.get('group_id') === 'B' ? native('B') : new Promise<Read>(resolve => { finishA = resolve; }));
    mount(); await openGroup(0); await waitFor(() => expect(finishA).toBeDefined());
    fireEvent.click(screen.getByRole('button', { name: 'Back' })); await openGroup(1);
    expect(await screen.findByTestId('group-list-item-B')).toHaveAttribute('data-locator-selected', 'true');
    const a = ledger.find(request => request.params.get('group_id') === 'G'); expect(a?.signal).toBeInstanceOf(AbortSignal);
    if (!(a?.signal instanceof AbortSignal) || !finishA) throw new Error('Missing actual A request'); expect(a.signal.aborted).toBe(true);
    await act(async () => { finishA?.(native()); });
    expect(screen.queryByTestId('group-list-item-G')).not.toBeInTheDocument(); expect(screen.getByTestId('group-list-item-B')).toBeVisible(); expectNoCommands();
  });
});
