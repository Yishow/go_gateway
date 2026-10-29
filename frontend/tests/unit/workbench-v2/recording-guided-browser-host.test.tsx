import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { AxiosError, type AxiosAdapter } from 'axios';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { basicRecordingScopeKey, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { groupCircuitId } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { RecordingGuidedHostView } from '../../e2e/fixtures/recording-guided-host-view';
import { createGuidedRecordingMock } from '../../e2e/fixtures/recording-guided-api';
import type { GuidedScenario } from '../../e2e/fixtures/recording-guided-state';

// Same formal browser host and route fixture. jsdom viability is NOT browser/geometry/screenshots or SQL proof.
const originalAdapter = studioV2DatalinkApi.defaults.adapter, originalBase = studioV2DatalinkApi.defaults.baseURL;
const clients: QueryClient[] = [];
type Mock = ReturnType<typeof createGuidedRecordingMock>;
function calls(mock: Mock, action: string) { return mock.ledger.filter(entry => entry.action === action); }
function connect(mock: Mock) {
  studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
  const adapter: AxiosAdapter = async config => {
    const reply = mock.handle({ method: (config.method ?? 'get').toUpperCase(),
      url: new URL(studioV2DatalinkApi.getUri(config), 'http://fixture.test').href,
      body: typeof config.data === 'string' ? JSON.parse(config.data) : config.data });
    if (reply.transport === 'abort') throw new AxiosError('Test-only lost response', 'ERR_NETWORK', config);
    const response = { data: JSON.parse(JSON.stringify(reply.body)), status: reply.status, statusText: 'Mock response', config, headers: {} };
    if (reply.status >= 400) throw new AxiosError('Test-only rejected request', 'ERR_BAD_RESPONSE', config, undefined, response);
    return response;
  };
  studioV2DatalinkApi.defaults.adapter = adapter;
}
async function open(scenario: GuidedScenario = 'new', setup?: (mock: Mock) => void) {
  const mock = createGuidedRecordingMock(scenario); setup?.(mock); connect(mock);
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const tree = () => <QueryClientProvider client={client}><MemoryRouter initialEntries={[`/tests/e2e/fixtures/recording-guided-host.html?scenario=${scenario}`]}>
    <RecordingGuidedHostView scenario={scenario} /></MemoryRouter></QueryClientProvider>;
  const view = render(tree()); await screen.findByTestId('circuit-profile-basic-panel');
  await waitFor(() => expect(client.isFetching()).toBe(0));
  return { mock, client, ...view, rerenderHost: () => view.rerender(tree()) };
}
async function click(testId: string) {
  const button = await screen.findByTestId(testId); await waitFor(() => expect(button).toBeEnabled()); fireEvent.click(button);
}
async function prepare(mock: Mock) {
  await click('circuit-profile-prepare'); await waitFor(() => expect(calls(mock, 'ensure')).toHaveLength(1));
  expect(calls(mock, 'ensure')[0].status).toBe(200); expect(calls(mock, 'schema-preview')).toHaveLength(0);
  await click('group-schema-preview'); await screen.findByTestId('group-schema-preview-result');
  expect(calls(mock, 'schema-apply')).toHaveLength(0); expect(calls(mock, 'start')).toHaveLength(0);
  await click('group-schema-apply'); await waitFor(() => expect(calls(mock, 'schema-apply')).toHaveLength(1));
  await waitFor(() => expect(screen.getByTestId('circuit-profile-start')).toBeEnabled());
}
async function start(mock: Mock) {
  const count = calls(mock, 'start').length; await click('circuit-profile-start');
  expect(screen.getByRole('dialog', { name: 'Confirm this recording scope' })).toBeVisible(); expect(calls(mock, 'start')).toHaveLength(count);
  await click('recording-start-confirm'); await waitFor(() => expect(calls(mock, 'start')).toHaveLength(count + 1));
}
function safeLedger(mock: Mock, allowed: string[] = [], errors: string[] = []) {
  expect(mock.ledger.filter(call => call.status === 404)).toEqual([]);
  expect(mock.ledger.filter(call => call.status >= 400 && !errors.includes(call.action))).toEqual([]);
  expect(mock.ledger.filter(call => call.mutation && !allowed.includes(call.action))).toEqual([]);
}
beforeEach(async () => {
  sessionStorage.clear(); localStorage.clear(); Object.defineProperty(navigator, 'onLine', { configurable: true, value: true });
  await i18n.changeLanguage('en');
});
afterEach(() => {
  cleanup(); clients.splice(0).forEach(client => client.clear());
  studioV2DatalinkApi.defaults.adapter = originalAdapter; studioV2DatalinkApi.defaults.baseURL = originalBase; vi.restoreAllMocks();
});

it('mounts actual Shell and completes new ordinary recording through distinct explicit commands and navigation only', async () => {
  const { mock } = await open(); await prepare(mock); await start(mock);
  expect(await screen.findByTestId('recording-completion-summary')).toHaveTextContent('Confirmed active devices');
  expect(screen.getByTestId('workbench-v2-shell')).toBeVisible(); expect(screen.queryByTestId('fixture-runtime-destination')).not.toBeInTheDocument();
  expect(calls(mock, 'start')[0].body).toMatchObject({ workspace_id: 'ws-1', device_ids: ['dev-1'] });
  const proceed = await screen.findByTestId('basic-recording-continue-runtime'); fireEvent.click(proceed); fireEvent.click(proceed);
  expect(await screen.findByTestId('fixture-runtime-scope')).toHaveTextContent('dev-1'); expect(calls(mock, 'start')).toHaveLength(1);
  fireEvent.click(screen.getByRole('button', { name: 'Back to setup' })); await screen.findByTestId('workbench-v2-shell');
  expect(calls(mock, 'start')).toHaveLength(1); safeLedger(mock, ['ensure', 'schema-preview', 'schema-apply', 'start']);
});
it('shows partial confirmed A and incomplete work before Continue, keeping original query primary', async () => {
  const { mock } = await open('partial'); await start(mock);
  const summary = await screen.findByTestId('recording-completion-summary'); expect(summary).toHaveTextContent('Recording groups not completed');
  const query = await screen.findByTestId('basic-recording-check'); expect(query).toHaveAttribute('data-primary', 'true');
  expect(screen.getByTestId('workbench-v2-shell')).toBeVisible(); await click('basic-recording-continue-runtime');
  expect(screen.getByTestId('fixture-runtime-scope')).toHaveTextContent('dev-1'); safeLedger(mock, ['start']);
});
it('keeps original O1 query after unknown outcome and explicit host remount without resending', async () => {
  const { mock } = await open('unknown'); await start(mock); await click('basic-recording-check');
  expect(calls(mock, 'start-operation').every(call => call.url.endsWith('/O1'))).toBe(true);
  fireEvent.click(screen.getByTestId('fixture-remount-host')); await screen.findByTestId('basic-recording-check');
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument(); expect(calls(mock, 'start')).toHaveLength(1); safeLedger(mock, ['start']);
});
it('lost ACK keeps no-ID manual recovery after remount and never discovers a fabricated operation', async () => {
  const { mock } = await open('saved-long', value => value.loseNext('start')); await start(mock);
  await screen.findByRole('button', { name: 'Inspect operation recovery' }); fireEvent.click(screen.getByTestId('fixture-remount-host'));
  await screen.findByRole('button', { name: 'Inspect operation recovery' });
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument(); expect(calls(mock, 'start-operation')).toHaveLength(0);
  expect(calls(mock, 'start')).toHaveLength(1); safeLedger(mock, ['start']);
});
it('rejects schema failure without DDL/start and preserves a separate explicit preview retry', async () => {
  const { mock } = await open('new', value => value.rejectNext('schema-preview'));
  await click('circuit-profile-prepare'); await click('group-schema-preview'); await screen.findByTestId('group-schema-error');
  expect(calls(mock, 'schema-apply')).toHaveLength(0); expect(calls(mock, 'start')).toHaveLength(0);
  await click('group-schema-preview'); await screen.findByTestId('group-schema-preview-result');
  expect(calls(mock, 'schema-preview')).toHaveLength(2); safeLedger(mock, ['ensure', 'schema-preview'], ['schema-preview']);
});
it.each(['retained_start', 'foreign_cache', 'offline'] as const)('rechecks fresh %s before the visible preview retry can send a command', async blocker => {
  const reads = vi.spyOn(Storage.prototype, 'getItem');
  const { mock, client } = await open('new', value => value.rejectNext('schema-preview'));
  await click('circuit-profile-prepare'); await click('group-schema-preview'); await screen.findByTestId('group-schema-error');
  const retry = screen.getByTestId('group-schema-preview'); await waitFor(() => expect(retry).toBeEnabled());
  const { groups, workspaceRevision } = mock.snapshot(), group = groups[0];
  if (!group) throw new Error('The explicit save must return the canonical group');
  // Change the real owner/cache before React publishes a disabled button. The existing handler must deny dispatch.
  await act(async () => {
    if (blocker === 'retained_start') {
      const scopeKey = basicRecordingScopeKey('ws-1', 'dev-1', encodeURIComponent(JSON.stringify([
        groupCircuitId(group), group.destination.connector_id, group.destination.connector_revision, group.destination.table_schema,
      ])));
      const ownerKeys = reads.mock.calls.map(([key]) => key).filter(key => key.startsWith('wbv2.basic-recording.'));
      expect(scopeKey).toBe(decodeURIComponent(ownerKeys.at(-1)?.slice('wbv2.basic-recording.'.length) ?? ''));
      saveBasicRecordingIntent({ scope_key: scopeKey, operation_id: 'retained-O1', request: {
        request_id: 'original-request', workspace_id: 'ws-1', expected_workspace_revision: workspaceRevision, device_ids: ['dev-1'],
        groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision }],
      } });
    } else if (blocker === 'foreign_cache') {
      client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'foreign', workspace_revision: workspaceRevision, groups });
    } else Object.defineProperty(navigator, 'onLine', { configurable: true, value: false });
    fireEvent.click(retry);
  });
  expect(calls(mock, 'schema-preview')).toHaveLength(1); expect(calls(mock, 'schema-apply')).toHaveLength(0);
  expect(calls(mock, 'start')).toHaveLength(0); safeLedger(mock, ['ensure', 'schema-preview'], ['schema-preview']);
});
it('uses the existing URL-independent selection resolver and never applies a historical/foreign/stale result', async () => {
  const { mock } = await open('saved-long');
  const harness = screen.getByTestId('fixture-selection-contract'); fireEvent.click(within(harness).getByText(/URL-independent selection contract/));
  fireEvent.click(screen.getByTestId('fixture-selection-owned')); expect(screen.getByTestId('fixture-selection-result')).toHaveTextContent('selected:dev-1:G1:R1');
  fireEvent.click(screen.getByTestId('fixture-selection-historical')); expect(screen.getByTestId('fixture-selection-result')).toHaveTextContent('history-different');
  fireEvent.click(screen.getByTestId('fixture-selection-stale')); expect(screen.getByTestId('fixture-selection-result')).toHaveTextContent('unavailable:evidence_stale');
  fireEvent.click(screen.getByTestId('fixture-selection-foreign')); expect(screen.getByTestId('fixture-selection-result')).toHaveTextContent('unavailable:scope_unavailable');
  safeLedger(mock);
});
it('shows independently simulated saved 30 and effective 15 through the real read-only runtime parser and summary', async () => {
  const { mock } = await open('saved-different');
  const summary = screen.getByTestId('recording-scope-summary');
  const saved = within(summary).getByRole('region', { name: 'Saved configuration' });
  const effective = within(summary).getByRole('region', { name: 'Verified effective configuration' });
  expect(saved).toHaveTextContent('30 seconds'); expect(saved).not.toHaveTextContent('15 seconds');
  expect(effective).toHaveTextContent('15 seconds'); expect(effective).not.toHaveTextContent('30 seconds');
  expect(summary).not.toHaveTextContent('/mock-private/'); expect(summary).not.toHaveTextContent('recordings.sqlite');
  const detail = within(summary).getByRole('button', { name: 'View complete rules and delivery evidence' });
  fireEvent.click(detail); const drawer = screen.getByRole('dialog', { name: 'Rules and evidence for this scope' });
  for (const [label, value] of [['Saved revision', 'R2'], ['Latest requested revision', 'R1'], ['Effective revision at observation', 'R1']] as const) {
    expect(within(drawer).getByText(label).nextElementSibling).toHaveTextContent(value);
  }
  expect(calls(mock, 'runtime-context')).toHaveLength(1); safeLedger(mock);
});
it.each(['saved-wide', 'saved-long', 'saved-custom', 'blocked-wide'] as const)('reopens %s without conversion or hidden commands', async scenario => {
  const { mock } = await open(scenario); const before = structuredClone(mock.snapshot().groups);
  const advanced = screen.getByTestId('step4-advanced-recording'), toggle = advanced.querySelector('summary');
  if (!toggle) throw new Error('Missing actual Advanced disclosure'); fireEvent.click(toggle);
  const card = await screen.findByTestId('group-list-item-G1'); await within(card).findByTestId('saved-recording-actions');
  expect(mock.snapshot().groups).toEqual(before); safeLedger(mock);
});
it('preserves a 30-second draft and Cancel focus through device changes with the formal host fixture', async () => {
  const { mock } = await open('saved-long'); await click('circuit-profile-start'); await click('recording-start-cancel');
  await waitFor(() => expect(screen.getByTestId('circuit-profile-start')).toHaveFocus());
  fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
  fireEvent.click(screen.getByRole('button', { name: 'Mock PLC B' })); fireEvent.click(screen.getByRole('button', { name: 'Mock PLC A' }));
  expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30'); safeLedger(mock);
});
it('retains explicitly empty device and uses actual bounded source controls for the long-list fixture', async () => {
  const { mock } = await open('long-sources'); await click('circuit-profile-sources-toggle');
  const list = screen.getByTestId('circuit-profile-member-list'); expect(list).toHaveClass('max-h-72', 'overflow-y-auto');
  expect(within(list).getAllByRole('checkbox', { name: /^Source \d+$/ })).toHaveLength(40);
  fireEvent.click(screen.getByRole('button', { name: 'Empty device' }));
  expect(screen.getByTestId('recording-scope-summary')).toHaveTextContent('Empty device');
  await act(async () => { await Promise.resolve(); }); safeLedger(mock);
});
