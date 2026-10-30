import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { AxiosError, type AxiosResponse } from 'axios';
import { BrowserRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { RuntimeLogsHostView } from '../../e2e/fixtures/runtime-logs-host';
import { createRuntimeLogsMock } from '../../e2e/fixtures/runtime-logs-api';
import { runtimeLogsFrame, runtimeLogsSnapshot, domainRecord, type RuntimeLogsScenario } from '../../e2e/fixtures/runtime-logs-state';
import { installRuntimeLogsStream, type RuntimeLogsStreamController } from '../../e2e/fixtures/runtime-logs-stream';

// Real production components and wire parser in jsdom: no browser geometry/keyboard/screenshot or SQL claim.
const originalAdapter = studioV2DatalinkApi.defaults.adapter;
let client: QueryClient; let stream: RuntimeLogsStreamController;
beforeEach(() => { localStorage.clear(); sessionStorage.clear(); window.history.replaceState(null, '', '/studio/logs'); });
afterEach(() => { cleanup(); client?.clear(); stream?.restore(); studioV2DatalinkApi.defaults.adapter = originalAdapter; vi.unstubAllGlobals(); });
async function mount(scenario: RuntimeLogsScenario = 'normal', lang = 'en') {
  await i18n.changeLanguage(lang); const mock = createRuntimeLogsMock(scenario);
  vi.stubGlobal('fetch', vi.fn<typeof globalThis.fetch>(async (input, init) => {
    const reply = mock.handle({ method: init?.method ?? 'GET', url: String(input) });
    return new Response(JSON.stringify(reply.body), { status: reply.status, headers: { 'content-type': 'application/json' } });
  }));
  studioV2DatalinkApi.defaults.adapter = async config => {
    const url = new URL(`/api/v1/datalink${config.url}`, 'http://fixture.invalid');
    for (const [key, value] of Object.entries(config.params ?? {})) url.searchParams.set(key, String(value));
    const reply = mock.handle({ method: (config.method ?? 'GET').toUpperCase(), url: url.href, body: config.data });
    const response: AxiosResponse = { data: reply.body, status: reply.status, statusText: '', headers: {}, config };
    if (reply.status >= 400) throw new AxiosError('Mock rejected', undefined, config, undefined, response);
    return response;
  };
  stream = installRuntimeLogsStream();
  client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><BrowserRouter><RuntimeLogsHostView /></BrowserRouter></QueryClientProvider>);
  await waitFor(() => expect(stream.requests).toHaveLength(1));
  await act(async () => { stream.send(0, runtimeLogsFrame({ type: 'handshake', metadata: runtimeLogsSnapshot(), cursor: 'instance-a:9007199254740993' })); });
  return mock;
}
function noCommands(mock: ReturnType<typeof createRuntimeLogsMock>) {
  expect(mock.unexpected).toEqual([]); expect(mock.ledger.every(call => !call.mutation)).toBe(true);
  expect(mock.ledger.some(call => Object.hasOwn(call.params, 'editor'))).toBe(false);
  for (const id of ['group-editor', 'basic-recording-panel', 'workspace-test-write-recovery']) expect(screen.queryByTestId(id)).not.toBeInTheDocument();
}

describe('formal logs host viability through actual Page, BrowserRouter and native locator receiver', () => {
  it.each(['en', 'zh-TW'])('%s opens immutable event G/R1 as current native G/R2 with no per-row GET or command', async lang => {
    const mock = await mount('normal', lang);
    expect(screen.getByTestId('runtime-logs-fixture-label')).toHaveTextContent('No backend, device or SQL verification');
    expect(mock.ledger).toHaveLength(1); expect(mock.ledger[0]?.path).toBe('/api/v1/system/logs');
    fireEvent.click(screen.getByRole('link', { name: i18n.t('domain.openGroup', { ns: 'runtime-logs' }) }));
    const card = await screen.findByTestId('group-list-item-G'); expect(card).toHaveAttribute('data-locator-selected', 'true');
    expect(window.location.search).toContain('group_revision=R1'); expect(screen.getByText('R2')).toBeVisible();
    await waitFor(() => expect(mock.ledger).toHaveLength(3)); noCommands(mock);
  });
  it.each(['deleted', 'foreign', 'unavailable'] as const)('%s preserves historical identity and fails closed without hydration or mutation', async scenario => {
    const mock = await mount(scenario);
    fireEvent.click(screen.getByRole('link', { name: 'Open group settings' }));
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).not.toHaveTextContent('Checking'));
    expect(screen.queryByTestId('group-list-item-G')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-list-item-G2')).not.toBeInTheDocument();
    expect(mock.ledger.filter(call => call.path.endsWith('/studio-v2/workspace'))).toHaveLength(1); noCommands(mock);
  });
  it('hostile and unknown candidates never render raw messages, arbitrary JSON or secrets', async () => {
    const mock = await mount('hostile');
    for (const marker of ['password_fixture', 'private.example', '<script', 'SELECT secret', 'credential']) expect(document.body.textContent).not.toContain(marker);
    expect(screen.queryByRole('link', { name: 'Open group settings' })).not.toBeInTheDocument(); noCommands(mock);
  });
  it('actual ReadableStream recovery remains G/R1 receipt evidence and pause blocks late bytes', async () => {
    const mock = await mount();
    const record = domainRecord('domain.delivery.recovered', { sequence: '9007199254740994' });
    await act(async () => { stream.send(0, runtimeLogsFrame({ type: 'log', record, cursor: 'instance-a:9007199254740994' })); });
    await screen.findByText(i18n.t('domain.problems.delivery.recovered', { ns: 'runtime-logs' }));
    fireEvent.click(screen.getByRole('button', { name: 'Pause view' }));
    expect(stream.requests[0]).toMatchObject({ aborted: true, closed: true });
    expect(stream.send(0, runtimeLogsFrame({ type: 'log', record, cursor: 'instance-a:9007199254740994' }))).toBe(false);
    noCommands(mock);
  });
});
