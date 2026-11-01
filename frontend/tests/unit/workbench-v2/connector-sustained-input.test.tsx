import { useEffect, useRef } from 'react';
import { useStudioV2DatabaseAutosave } from '@/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave';
import { useWorkbenchV2State, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { studioV2WorkspaceDatabaseAPI as databaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
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
  const actions = useWorkbenchV2State(recordingAuthoringHostState(scenario));
  const stateRef = useRef(actions.state);
  useEffect(() => { stateRef.current = actions.state; }, [actions.state]);
  const database = useStudioV2DatabaseAutosave(actions, stateRef, true);
  const dispatch = (action: WorkbenchV2Action) => {
    const next = workbenchV2Reducer(stateRef.current, action);
    stateRef.current = next; actions.dispatch(action); database.afterDatabaseAction(action, next);
  };
  return <Step4Database state={actions.state} dispatch={dispatch} workspaceId="ws-1" />;
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


it.each([
  ['database', 'DB File Path', '/workspace/scratch/b2a0a232be39/linux-gateway-qa/qa-runtime-fixture13/recordings.sqlite'],
  ['database', 'DB File Path', `/tmp/${'long-path-segment/'.repeat(15)}recordings.sqlite`],
  ['name', 'Connection Name', 'Operator connection '.repeat(12)],
] as const)('preserves sustained %s input and the final autosave acknowledgement', async (field, label, expected) => {
  vi.spyOn(databaseAPI, 'getConfig').mockResolvedValue(null);
  vi.spyOn(databaseAPI, 'listTargets').mockResolvedValue([]);
  const saves = vi.spyOn(databaseAPI, 'updateConfig').mockImplementation(async request => ({
    ...request, id: 'conn-1', workspace_id: 'ws-1', identity_revision: 'crev-1', setup_revision: 'setup-1',
    status: 'ready', save_state: 'saved', created_at: '2030-01-01T00:00:00Z', updated_at: '2030-01-01T00:00:00Z',
    runtime_apply_status: 'not_running',
  }));
  await open();
  fireEvent.click(screen.getByTestId('database-connection-card-toggle'));
  const input = screen.getByLabelText(label) as HTMLInputElement;
  const errors: string[] = [];
  const capture = (event: ErrorEvent) => { errors.push(String(event.error ?? event.message)); event.preventDefault(); };
  window.addEventListener('error', capture);
  const setValue = Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, 'value')!.set!;
  const globals = globalThis as typeof globalThis & { IS_REACT_ACT_ENVIRONMENT: boolean };
  const previous = globals.IS_REACT_ACT_ENVIRONMENT; globals.IS_REACT_ACT_ENVIRONMENT = false;
  try {
    setValue.call(input, ''); input.dispatchEvent(new Event('input', { bubbles: true })); await Promise.resolve();
    // Do not wrap each input in act: that drains background lanes between keystrokes
    // and masks the production failure. Append to the actual controlled DOM value.
    for (const character of expected) {
      setValue.call(input, input.value + character);
      input.dispatchEvent(new Event('input', { bubbles: true }));
      await Promise.resolve();
    }
    await new Promise(resolve => setTimeout(resolve, 50));
  } finally { globals.IS_REACT_ACT_ENVIRONMENT = previous; window.removeEventListener('error', capture); }
  await waitFor(() => expect(saves).toHaveBeenCalled());
  expect.soft(errors).toEqual([]);
  expect.soft(input.value).toBe(expected);
  await waitFor(() => expect(saves.mock.calls.at(-1)?.[0][field]).toBe(expected));
});

it.each(['devices', 'rules', 'points', 'mappings'] as const)('still refreshes eligible sources when %s changes', async slice => {
  const mock = createRecordingAuthoringMock('new-sqlite'); connect(mock);
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  clients.push(client);
  const initial = recordingAuthoringHostState('new-sqlite');
  const tree = (state: typeof initial) => <QueryClientProvider client={client}>
    <Step4Database state={state} dispatch={vi.fn()} workspaceId="ws-1" />
  </QueryClientProvider>;
  const view = render(tree(initial));
  await screen.findByTestId('circuit-profile-basic-panel');
  await waitFor(() => expect(client.isFetching()).toBe(0));
  fireEvent.click(screen.getByTestId('circuit-profile-sources-toggle'));
  expect(screen.getByTestId('circuit-profile-member-inst-tag-0')).toBeChecked();
  const changed = { ...initial, ...(slice === 'devices' ? { devices: initial.devices.map(device => ({ ...device, save_state: 'saving' as const })) }
    : slice === 'rules' ? { rules: initial.rules.map(rule => ({ ...rule, enabled: false })) }
      : slice === 'points' ? { points: initial.points.filter(point => point.id !== 'rule-1-p-0') }
        : { mappings: { ...initial.mappings, 'rule-1-p-0': { ...initial.mappings['rule-1-p-0'], enabled: false } } }) };
  view.rerender(tree(changed));
  expect(screen.queryByTestId('circuit-profile-member-inst-tag-0')).not.toBeInTheDocument();
  view.rerender(tree(initial));
  const toggle = await screen.findByTestId('circuit-profile-sources-toggle');
  if (toggle.getAttribute('aria-expanded') !== 'true') fireEvent.click(toggle);
  expect(await screen.findByTestId('circuit-profile-member-inst-tag-0')).toBeChecked();
});
