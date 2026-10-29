import axios, { AxiosError, AxiosHeaders, type AxiosAdapter } from 'axios';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { api } from '../../../src/services/datalinkClient';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { studioV2RuntimeContextAPI } from '../../../src/services/studioV2RuntimeContext';
import { readStudioLocatorEntry } from '../../../src/services/studioV2LocatorEntry';
import { buildStudioEntityLocatorTarget } from '../../../src/services/studioEntityLocator';
import type { EntityLocator, EvidenceScope } from '../../../src/types/runtimeRuleEvidence';
import { createRuntimeRecordingMock } from '../../e2e/fixtures/runtime-recording-api';
import { RuntimeRecordingHostView } from '../../e2e/fixtures/runtime-recording-host';
import { installRuntimeRecordingStreamFixture } from '../../e2e/fixtures/runtime-recording-stream';

// Fixture/parser and jsdom viability only: NOT an executed browser, embedded server, SQL or screenshot result.
const scope: EvidenceScope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'A' };
const selection = { rule_id: 'A1', rule_revision: 'R2', point_id: 'PA1', mapping_id: 'MA1' };
const originalAdapters = [api.defaults.adapter, studioV2DatalinkApi.defaults.adapter];
const originalEventSource = globalThis.EventSource;
let mock: ReturnType<typeof createRuntimeRecordingMock>;
const clients: QueryClient[] = [];
function locator(step: 1 | 2 | 3 | 4): EntityLocator {
  return { step, workspace_id: 'W', device_id: 'A', rule_id: step === 2 ? 'A1' : null,
    mapping_id: step === 3 ? 'MA1' : null, group_id: step === 4 ? 'G1' : null, group_revision: step === 4 ? 'A' : null };
}
function installTransport() {
  const adapter: AxiosAdapter = async config => {
    const result = mock.handle({ method: config.method?.toUpperCase() ?? 'GET',
      url: new URL(axios.getUri(config), 'http://fixture.local').href, body: config.data });
    const response = { data: result.body, status: result.status, statusText: String(result.status), headers: new AxiosHeaders(), config };
    if (result.status >= 400) throw new AxiosError('Named fixture transport failure', 'ERR_BAD_RESPONSE', config, undefined, response);
    return response;
  };
  api.defaults.adapter = adapter; studioV2DatalinkApi.defaults.adapter = adapter;
}
function mount(path = '/studio/runtime?device_id=A') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } });
  clients.push(client);
  return render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[path]}><RuntimeRecordingHostView /></MemoryRouter></QueryClientProvider>);
}
function cleanLedger() {
  expect(mock.ledger.filter(entry => entry.mutation || entry.action === 'unknown' || entry.status >= 400)).toEqual([]);
  expect(mock.ledger.some(entry => /logs|proxy/.test(entry.path))).toBe(false);
}
beforeEach(async () => {
  await i18n.changeLanguage('en'); mock = createRuntimeRecordingMock(); installTransport();
  installRuntimeRecordingStreamFixture();
  vi.spyOn(XMLHttpRequest.prototype, 'send').mockImplementation(() => { throw new Error('Unmocked HTTP escaped strict fixture transport'); });
});
afterEach(() => {
  cleanup(); clients.splice(0).forEach(client => client.clear());
  api.defaults.adapter = originalAdapters[0]; studioV2DatalinkApi.defaults.adapter = originalAdapters[1];
  Object.defineProperty(globalThis, 'EventSource', { value: originalEventSource, writable: true, configurable: true });
  vi.restoreAllMocks();
});

describe('formal runtime fixture contracts through production readers', () => {
  it('returns six rules and ten mappings with exact explicit selection and no detail fan-out', async () => {
    const rules = await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', query: '', limit: 50 });
    const mappings = await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'mappings', query: '', limit: 50 });
    expect(rules.ok).toBe(true); expect(mappings.ok).toBe(true);
    if (!rules.ok || !mappings.ok) throw new Error('Expected valid source pages');
    expect(rules.value.rules.map(row => row.id)).toEqual(['A1', 'A2', 'A3', 'A4', 'A5', 'A6']);
    expect(mappings.value.mappings).toHaveLength(10);
    expect(mock.ledger.map(entry => entry.action)).toEqual(['rules', 'mappings']); cleanLedger();
  });
  it('keeps saved R2 scale0.1 and polling1000ms distinct from unproven effective source and mean15s recording', async () => {
    const source = await studioV2RuntimeContextAPI.getSourceEvidence(scope, selection);
    const groups = await studioV2RuntimeContextAPI.getRuleEvidence(scope);
    expect(source.ok).toBe(true); expect(groups.ok).toBe(true);
    if (!source.ok || !groups.ok) throw new Error('Expected valid evidence fixture');
    expect(source.value.source).toMatchObject({ saved_revision: 'R2', state: 'saved-only', proof: 'unavailable', effective: null });
    expect(source.value.source?.saved?.polling).toMatchObject({ interval_ms: 1000, source: 'polling-group' });
    expect(source.value.source?.saved?.transforms).toContainEqual(expect.objectContaining({ kind: 'scale', multiplier: 0.1 }));
    expect(groups.value.groups[0].effective).toMatchObject({ interval_seconds: 15, reduction: 'mean' }); cleanLedger();
  });
  it.each(['normal', 'draft', 'backlog', 'uncertain', 'shared', 'missing-timing', 'no-proof'] as const)('parses %s evidence without bypassing the real proof contract', async scenario => {
    mock.setScenario(scenario); const evidence = await studioV2RuntimeContextAPI.getRuleEvidence(scope);
    expect(evidence.ok).toBe(true);
    if (!evidence.ok) throw new Error(`Invalid ${scenario} evidence`);
    expect(evidence.value.groups[0].group_id).toBe('G1'); cleanLedger();
  });
  it('does not let the missing-proof fixture claim an effective group or display its active facts', async () => {
    mock.setScenario('no-proof'); const evidence = await studioV2RuntimeContextAPI.getRuleEvidence(scope);
    expect(evidence.ok).toBe(true); if (!evidence.ok) throw new Error('Expected readable saved-only group');
    expect(evidence.value.groups[0].proof).toBe('unavailable');
    mount(); const effective = await screen.findByRole('region', { name: 'Effective recording snapshot' });
    expect(effective).toHaveTextContent('Effective recording facts unavailable'); expect(effective).not.toHaveTextContent('meter_old');
    cleanLedger();
  });
  it('bounds120 rows using a revision-bound next cursor and exact search, not hidden nested rows', async () => {
    mock.setScenario('many');
    const first = await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', query: '', limit: 50 });
    expect(first.ok).toBe(true); if (!first.ok) throw new Error('Expected first page');
    expect(first.value.rules).toHaveLength(50); expect(first.value.mappings).toEqual([]); expect(first.value.next_cursor).not.toBeNull();
    if (!first.value.next_cursor) throw new Error('Expected bounded next cursor');
    const next = await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', query: '', limit: 50, cursor: first.value.next_cursor });
    expect(next.ok).toBe(true); if (!next.ok) throw new Error('Expected next page');
    expect(next.value.rules).toHaveLength(50); expect(next.value.rules[0].id).not.toBe(first.value.rules[0].id);
    const found = await studioV2RuntimeContextAPI.getSourcePage(scope, { kind: 'rules', query: 'A120', limit: 50 });
    expect(found.ok).toBe(true); if (!found.ok) throw new Error('Expected search page');
    expect(found.value.rules.map(row => row.id)).toEqual(['A120']); cleanLedger();
  });
  it.each([1, 2, 3, 4] as const)('validates detached native Step%i entry using the actual receiver reader', async step => {
    const target = locator(step); const result = await readStudioLocatorEntry(target);
    expect(result.ok).toBe(true); if (!result.ok) throw new Error('Expected scoped entry');
    expect(result.value.selection.state).toBe('selected'); expect(result.value.entry.hydration_scope.complete_workspace).toBe(false);
    expect(result.value.locator).toEqual(target);
    expect(mock.ledger).toHaveLength(1); expect(mock.ledger[0].action).toBe('locator'); cleanLedger();
  });
  it('preserves typed stale and foreign failures instead of manufacturing successful evidence', async () => {
    mock.setScenario('stale'); expect(await studioV2RuntimeContextAPI.getRuleEvidence(scope)).toMatchObject({ ok: false, reason: 'stale-reference' });
    mock.setScenario('foreign'); expect(await readStudioLocatorEntry(locator(4))).toMatchObject({ ok: false, reason: 'foreign-scope' });
    expect(mock.ledger.every(entry => !entry.mutation)).toBe(true);
  });
  it('rejects unknown reads, mutations and partial/foreign scopes in the mock itself', () => {
    const path = 'http://fixture.local/api/v1/datalink/studio-v2/workspace/runtime-context';
    const requests = [
      { method: 'POST', url: path }, { method: 'GET', url: 'http://fixture.local/api/v1/datalink/unrecognized' },
      { method: 'GET', url: `${path}?workspace_id=W` },
      { method: 'GET', url: `${path}?workspace_id=Foreign&workspace_revision=W1&device_id=A` },
    ];
    for (const request of requests) expect(mock.handle(request).status).toBeGreaterThanOrEqual(400);
    expect(mock.ledger).toHaveLength(4); expect(mock.ledger[0].mutation).toBe(true);
  });
});

describe('real runtime and native receiver fixture viability in jsdom', () => {
  it('mounts the real runtime source/group views with one scoped evidence read and no mutation or log proxy', async () => {
    mount(); const group = await screen.findByRole('article', { name: 'Recording group G1' });
    expect(within(group).getByRole('region', { name: 'Effective recording snapshot' })).toHaveTextContent('meter_old');
    const mappings = await screen.findByRole('region', { name: 'Source mappings' });
    expect(within(mappings).getAllByRole('button', { name: /Inspect source/ })).toHaveLength(10);
    expect(mock.ledger.filter(entry => entry.action === 'evidence')).toHaveLength(1); cleanLedger();
  });
  it('drives actual stream wire readiness and disconnect without injecting component state or freshening the receipt', async () => {
    let browserNow = Date.parse('2026-10-29T10:00:00Z'); vi.spyOn(Date, 'now').mockImplementation(() => browserNow);
    mount(); await screen.findByRole('article', { name: 'Recording group G1' });
    await waitFor(() => expect(window.__runtimeRecordingStreams.list().length).toBeGreaterThan(0));
    const index = window.__runtimeRecordingStreams.list().length - 1;
    act(() => window.__runtimeRecordingStreams.emit(index, 'stream_state', { device_id: 'A', timestamp: '2026-10-29T10:00:00Z',
      stream_state: { state: 'ready', empty: false, stale: false, degraded: false, unavailable: false } }));
    await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live'));
    const receipt = screen.getByTestId('runtime-snapshot-received-at'); expect(receipt).toHaveAttribute('dateTime', '2026-10-29T10:00:00.000Z');
    browserNow += 5000; act(() => window.__runtimeRecordingStreams.fail(index));
    await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('degraded'));
    expect(screen.getByTestId('runtime-snapshot-received-at')).toHaveAttribute('dateTime', '2026-10-29T10:00:00.000Z');
    expect(screen.getByRole('article', { name: 'Recording group G1' })).toHaveTextContent('meter_old'); cleanLedger();
  });
  it('follows a real runtime group link to the exact native cold card with no ordinary bootstrap or mutation', async () => {
    mount(); const group = await screen.findByRole('article', { name: 'Recording group G1' });
    const link = within(group).getByRole('link', { name: 'Open saved group settings' }); fireEvent.click(link);
    await screen.findByTestId('workbench-v2-shell');
    await waitFor(() => expect(mock.ledger.filter(entry => entry.action === 'locator')).toHaveLength(1));
    expect(mock.ledger.filter(entry => entry.path.endsWith('/studio-v2/workspace')).every(entry => entry.params.entry === 'locator')).toBe(true);
    expect(screen.queryByTestId('group-editor')).not.toBeInTheDocument(); cleanLedger();
  });
  it.each([1, 2, 3] as const)('cold-mounts real native Step%i without automatic configuration writes', async step => {
    const target = buildStudioEntityLocatorTarget(locator(step)); if (!target) throw new Error('Expected safe target');
    mount(target); await screen.findByTestId('workbench-v2-shell');
    expect(mock.ledger.filter(entry => entry.action === 'locator')).toHaveLength(1);
    expect(mock.ledger.filter(entry => entry.path.endsWith('/studio-v2/workspace')).every(entry => entry.params.entry === 'locator')).toBe(true);
    cleanLedger();
  });
});
