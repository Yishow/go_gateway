import { useReducer } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { AxiosError, type AxiosAdapter } from 'axios';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { workbenchV2Reducer } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import i18n from '@/i18n/config';
import { createRecordingAuthoringMock } from '../../e2e/fixtures/recording-authoring-api';
import { recordingAuthoringHostState, type AuthoringScenario } from '../../e2e/fixtures/recording-authoring-state';

// jsdom fixture viability only: not Playwright execution, screenshots, SQL or production HTTP wiring.
const originalAdapter = studioV2DatalinkApi.defaults.adapter;
const originalBase = studioV2DatalinkApi.defaults.baseURL;
type Mock = ReturnType<typeof createRecordingAuthoringMock>;
const clients: QueryClient[] = [];
function connect(mock: Mock) {
  studioV2DatalinkApi.defaults.baseURL = '/api/v1/datalink';
  const adapter: AxiosAdapter = async config => {
    const reply = mock.handle({ method: (config.method ?? 'get').toUpperCase(),
      url: new URL(studioV2DatalinkApi.getUri(config), 'http://fixture.test').href,
      body: typeof config.data === 'string' ? JSON.parse(config.data) : config.data });
    if (reply.transport === 'abort') throw new AxiosError('Test-only lost response', 'ERR_NETWORK', config);
    // Match route.fulfill(JSON.stringify(...)): optional undefined fields do not exist on the wire.
    const response = { data: JSON.parse(JSON.stringify(reply.body)), status: reply.status, statusText: 'Mock response', config, headers: {} };
    if (reply.status >= 400) throw new AxiosError('Test-only rejected request', 'ERR_BAD_RESPONSE', config, undefined, response);
    return response;
  };
  studioV2DatalinkApi.defaults.adapter = adapter;
}
function Host({ scenario }: { scenario: AuthoringScenario }) {
  const [state, dispatch] = useReducer(workbenchV2Reducer, scenario, recordingAuthoringHostState);
  return <><button data-testid="fixture-receive-revision" onClick={() => dispatch({ type: 'updateDbConnector',
    patch: { identity_revision: 'crev-2', save_state: 'saved', persisted: true } })}>Receive saved connector revision</button>
    <Step4Database state={state} dispatch={dispatch} workspaceId="ws-1" /></>;
}
function mount(scenario: AuthoringScenario) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  clients.push(client);
  return render(<QueryClientProvider client={client}><Host scenario={scenario} /></QueryClientProvider>);
}
async function open(scenario: AuthoringScenario = 'new-sqlite', setup?: (mock: Mock) => void) {
  const mock = createRecordingAuthoringMock(scenario);
  setup?.(mock); connect(mock); const view = mount(scenario);
  await screen.findByTestId('circuit-profile-basic-panel');
  await waitFor(() => expect(mock.ledger.some(call => call.action === 'test-write-list')).toBe(true));
  await waitFor(() => expect(clients.at(-1)?.isFetching()).toBe(0));
  return { mock, view };
}
function sources() {
  const toggle = screen.getByTestId('circuit-profile-sources-toggle');
  if (toggle.getAttribute('aria-expanded') !== 'true') fireEvent.click(toggle);
}
function calls(mock: Mock, action: string) { return mock.ledger.filter(entry => entry.action === action); }
function assertRequests(mock: Mock, mutations: string[] = [], errors: string[] = []) {
  expect(mock.ledger.filter(entry => entry.status >= 400 && !errors.includes(entry.action))).toEqual([]);
  expect(mock.ledger.filter(entry => entry.mutation && !mutations.includes(entry.action))).toEqual([]);
}
async function advancedPanel() {
  const advanced = screen.getByTestId('step4-advanced-recording');
  fireEvent.click(advanced.querySelector('summary')!);
  fireEvent.click(await within(advanced).findByTestId('group-open-G1'));
  return within(advanced).findByTestId('group-test-write');
}
async function submit(mock: Mock) {
  const panel = await advancedPanel();
  const preview = within(panel).getByTestId('group-test-write-preview');
  await waitFor(() => expect(preview).toBeEnabled()); fireEvent.click(preview);
  fireEvent.click(await within(panel).findByTestId('group-test-write-confirm'));
  await waitFor(() => expect(calls(mock, 'test-write-confirm')).toHaveLength(1));
  return panel;
}
beforeEach(async () => {
  sessionStorage.clear(); localStorage.clear();
  for (const ref of recordingTestWriteRecovery.read().references) recordingTestWriteRecovery.resolve({ operation_id: ref.operation_id,
    action: 'test_write', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned', created_at: '', updated_at: '' }, ref.workspace_id);
  await i18n.changeLanguage('en');
});
afterEach(() => {
  cleanup(); clients.splice(0).forEach(client => client.clear());
  studioV2DatalinkApi.defaults.adapter = originalAdapter; studioV2DatalinkApi.defaults.baseURL = originalBase;
  vi.restoreAllMocks();
});

it.each(['sqlite', 'postgres', 'mysql', 'mariadb', 'sqlserver'] as const)('real host can prepare mock %s long without bypassing production parsers', async engine => {
  const { mock } = await open(`new-${engine}`);
  const wide = screen.getByTestId('circuit-profile-purpose-wide');
  if (engine === 'sqlite' || engine === 'postgres') expect(wide).toBeEnabled(); else expect(wide).toBeDisabled();
  fireEvent.click(screen.getByTestId('circuit-profile-purpose-long'));
  await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
  fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
  await screen.findByTestId('group-schema-preview-result');
  expect(calls(mock, 'ensure')).toHaveLength(1); expect(calls(mock, 'ensure')[0].status).toBe(200);
  expect(mock.snapshot().groups[0].layout_version).toBe('managed_long_v1');
  assertRequests(mock, ['ensure', 'schema-preview']);
});
it('real numeric role control authors state with no required counter name and preserves gauge', async () => {
  const { mock } = await open(); sources();
  fireEvent.change(screen.getByTestId('circuit-profile-role-tag-0'), { target: { value: 'state' } });
  expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveAttribute('aria-required', 'false');
  fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
  await screen.findByTestId('group-schema-preview-result');
  expect(mock.snapshot().groups[0].members.map(member => member.value_role)).toEqual(['state', 'gauge', 'state']);
  assertRequests(mock, ['ensure', 'schema-preview']);
});
it('instant-only counter gives a required-name reason and does not save or start from naming', async () => {
  const { mock } = await open(); sources();
  fireEvent.change(screen.getByTestId('circuit-profile-role-tag-0'), { target: { value: 'counter_raw' } });
  expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveAttribute('aria-required', 'true');
  expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
  fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Counter A' } });
  expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled(); assertRequests(mock);
});
it('new cumulative source selection starts empty and persists a deliberate selection on remount', async () => {
  const { mock, view } = await open(); fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
  expect(screen.getByTestId('circuit-profile-cumulative-selected')).toHaveTextContent('Selected cumulative sources: 0');
  fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-0'));
  view.unmount(); mount('new-sqlite'); await screen.findByTestId('circuit-profile-basic-panel'); sources();
  await waitFor(() => expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked());
  expect(screen.getByTestId('circuit-profile-member-cum-tag-1')).not.toBeChecked(); assertRequests(mock);
});
it('counter query failure exposes a real retry without mutating the draft', async () => {
  const { mock } = await open('new-sqlite', value => value.setMeasurementError('dev-1', true)); sources();
  fireEvent.change(screen.getByTestId('circuit-profile-role-tag-0'), { target: { value: 'counter_raw' } });
  fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Counter A' } });
  const error = await screen.findByTestId('circuit-profile-measurements-error');
  mock.setMeasurementError('dev-1', false); fireEvent.click(within(error).getByRole('button', { name: 'Retry' }));
  await waitFor(() => expect(screen.queryByTestId('circuit-profile-measurements-error')).not.toBeInTheDocument());
  expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Counter A');
  expect(screen.getByTestId('circuit-profile-role-tag-0')).toHaveValue('counter_raw');
  assertRequests(mock, [], ['measurements:dev-1']);
});
it('saved role and source meaning survive actual supported lifecycle review and dismissal', async () => {
  const { mock } = await open('saved-gauge'); sources();
  expect(screen.getByTestId('circuit-profile-role-tag-0')).toBeDisabled();
  const advanced = screen.getByTestId('step4-advanced-recording');
  fireEvent.click(advanced.querySelector('summary')!); fireEvent.click(await within(advanced).findByTestId('group-open-G1'));
  fireEvent.click(await within(advanced).findByTestId('saved-recording-transition'));
  const dialog = await screen.findByRole('dialog', { name: 'Change saved recording' });
  fireEvent.click(within(dialog).getByRole('button', { name: 'Review impact' }));
  fireEvent.click(await within(dialog).findByRole('button', { name: 'Back' }));
  fireEvent.click(within(dialog).getByRole('button', { name: 'Cancel' }));
  await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
  expect(mock.snapshot().groups).toEqual(mock.snapshot().frozenGroups);
  assertRequests(mock, ['schema-preview']);
});
it('real host retains lost test-write ACK after remount with only original O1 query access', async () => {
  const { mock, view } = await open('test-write', value => value.loseNext('test-write-confirm'));
  const panel = await submit(mock); await within(panel).findByTestId('group-test-write-unconfirmed');
  expect(within(panel).getByTestId('group-test-write-preview')).toBeDisabled();
  view.unmount(); mount('test-write');
  await waitFor(() => expect(screen.getByTestId('workspace-test-write-recovery')).toHaveTextContent('O1'));
  expect(mock.snapshot().insertCount).toBe(1);
  assertRequests(mock, ['schema-preview', 'test-write-preview', 'test-write-confirm']);
});
it('real SQLServer-capability fixture offers same-O1 cleanup repair without another insert', async () => {
  const { mock } = await open('test-write-sqlserver', value => value.setTestWriteOutcome('succeeded', 'written_verified', 'failed'));
  const panel = await submit(mock);
  const repair = await within(panel).findByTestId('group-test-write-cleanup-repair');
  expect(within(panel).getByTestId('group-test-write-write-phase')).toHaveTextContent('Write committed');
  fireEvent.click(repair);
  await waitFor(() => expect(within(panel).getByTestId('group-test-write-cleanup')).toHaveTextContent('Test row removed'));
  expect(mock.snapshot().insertCount).toBe(1); expect(mock.snapshot().repairCount).toBe(1);
  assertRequests(mock, ['schema-preview', 'test-write-preview', 'test-write-confirm', 'cleanup-repair']);
});
it('unsupported saved wide stays visible and cannot be prepared', async () => {
  const { mock } = await open('saved-unsupported-wide');
  expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeChecked();
  expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeDisabled();
  expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
  expect(screen.getByTestId('circuit-profile-layout-blocked')).toHaveTextContent('wide tables are unavailable');
  expect(mock.snapshot().groups).toEqual(mock.snapshot().frozenGroups); assertRequests(mock, ['schema-preview']);
});
it('the fixture external-revision control makes real capability guards reject old R1 evidence', async () => {
  const { mock } = await open('new-sqlserver');
  fireEvent.click(screen.getByTestId('circuit-profile-purpose-long'));
  expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();
  fireEvent.click(screen.getByTestId('fixture-receive-revision'));
  await waitFor(() => expect(screen.getByTestId('circuit-profile-layout-blocked')).toHaveTextContent('unverified'));
  expect(screen.getByTestId('circuit-profile-prepare')).toBeDisabled();
  expect(calls(mock, 'capability').length).toBeGreaterThan(1); assertRequests(mock);
});
