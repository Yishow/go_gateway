import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { AxiosError, type AxiosResponse } from 'axios';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import i18n from '@/i18n/config';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { createRecordingLifecycleMock, type RecordingLifecycleMock } from '../../e2e/fixtures/recording-lifecycle-api';
import { recordingLifecycleHostState, type LifecycleScenario } from '../../e2e/fixtures/recording-lifecycle-state';

// This is a jsdom fixture/provider reachability check, NOT browser, keyboard, geometry, or screenshot acceptance.
const originalAdapter = studioV2DatalinkApi.defaults.adapter;
let client: QueryClient;
beforeEach(async () => { localStorage.clear(); sessionStorage.clear(); await i18n.changeLanguage('en'); });
afterEach(() => { cleanup(); client?.clear(); studioV2DatalinkApi.defaults.adapter = originalAdapter; });
function mount(scenario: LifecycleScenario, mock: RecordingLifecycleMock = createRecordingLifecycleMock(scenario)) {
  studioV2DatalinkApi.defaults.adapter = async config => {
    const url = new URL(`${config.baseURL}${config.url}`, 'http://fixture.invalid');
    for (const [key, value] of Object.entries(config.params ?? {})) url.searchParams.set(key, String(value));
    const reply = mock.handle({ method: (config.method ?? 'GET').toUpperCase(), url: url.href,
      body: typeof config.data === 'string' ? JSON.parse(config.data) : config.data });
    const response: AxiosResponse = { data: reply.body, status: reply.status, statusText: '', headers: {}, config };
    if (reply.status >= 400) throw new AxiosError('Mock rejected', undefined, config, undefined, response);
    return response;
  };
  client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><Step4Database state={recordingLifecycleHostState()} dispatch={() => undefined} workspaceId="ws-1" /></QueryClientProvider>);
  return mock;
}
describe('real Step4Database host with browser fixture mock contracts (jsdom only)', () => {
  it.each(['applied-wide', 'draft-wide', 'disabled-wide', 'blocked-wide', 'review-blocked-wide', 'stale-wide'] as const)('%s exposes its existing ordinary-preview baseline before lifecycle interaction', async scenario => {
    const mock = mount(scenario);
    expect(await within(screen.getByTestId('step4-basic-recording')).findByTestId('group-schema-preview-result')).toBeInTheDocument();
    expect(mock.ledger.filter(call => call.mutation)).toMatchObject([{ action: 'ordinary-schema-preview', status: 200 }]);
    expect(mock.ledger.filter(call => call.status !== 200)).toEqual([]);
    expect(mock.snapshot().status).toBeUndefined();
  });
  it.each(['applied-wide', 'applied-long', 'deleted-wide', 'deleted-basic'] as const)('%s reaches shared saved action controls and production dialog without a replacement component', async scenario => {
    const mock = mount(scenario);
    const basic = within(screen.getByTestId('step4-basic-recording'));
    if (scenario.startsWith('deleted')) fireEvent.click((await basic.findByTestId('basic-removed-recordings')).querySelector('summary')!);
    const action = await basic.findByTestId(scenario.startsWith('deleted') ? 'saved-recording-recreate' : 'saved-recording-transition');
    await waitFor(() => expect(action).not.toBeDisabled());
    await waitFor(() => expect(client.isMutating()).toBe(0));
    // The existing Basic wide host may automatically read the ORIGINAL group's schema before dialog use.
    expect(mock.ledger.filter(call => call.mutation).every(call => call.action === 'ordinary-schema-preview')).toBe(true);
    const baseline = mock.ledger.length;
    fireEvent.click(action);
    const dialog = await screen.findByRole('dialog');
    expect(await within(dialog).findByTestId('recording-change-review')).not.toBeDisabled();
    expect(mock.ledger.slice(baseline).filter(call => call.mutation)).toEqual([]);
    expect(mock.ledger.filter(call => call.status !== 200)).toEqual([]);
  });
  it('advanced editor uses the same fixture lifecycle and routes managed start without ordinary apply', async () => {
    const mock = mount('draft-wide');
    const advanced = screen.getByTestId('step4-advanced-recording');
    fireEvent.click(advanced.querySelector('summary')!);
    fireEvent.click(await within(advanced).findByTestId('group-open-G1'));
    const apply = await within(advanced).findByTestId('group-apply');
    await waitFor(() => expect(apply).not.toBeDisabled());
    await waitFor(() => expect(client.isMutating()).toBe(0));
    expect(mock.ledger.filter(call => call.mutation).every(call => call.action === 'ordinary-schema-preview')).toBe(true);
    const baseline = mock.ledger.length;
    fireEvent.click(apply);
    expect(mock.ledger.slice(baseline).filter(call => call.mutation)).toEqual([]);
    expect(mock.ledger.filter(call => call.status !== 200)).toEqual([]);
  });
  it('terminal real-host remount restores the recorded root with no unsupported owned-group requests', async () => {
    const mock = mount('applied-wide');
    const basic = within(screen.getByTestId('step4-basic-recording'));
    await basic.findByTestId('group-schema-preview-result');
    const transition = await basic.findByTestId('saved-recording-transition');
    await waitFor(() => expect(transition).not.toBeDisabled()); fireEvent.click(transition);
    const dialog = await screen.findByRole('dialog');
    for (const action of ['review', 'reserve_preview', 'confirm_schema']) {
      const button = await within(dialog).findByTestId(`recording-change-${action}`);
      await waitFor(() => expect(button).not.toBeDisabled()); fireEvent.click(button);
    }
    await waitFor(() => expect(dialog.querySelector('input[type="datetime-local"]')).not.toBeNull());
    fireEvent.change(dialog.querySelector('input[type="datetime-local"]')!, { target: { value: '2099-01-01T01:00:00' } });
    for (const action of ['review_activation', 'confirm_cutover']) {
      const button = await within(dialog).findByTestId(`recording-change-${action}`);
      await waitFor(() => expect(button).not.toBeDisabled()); fireEvent.click(button);
    }
    await waitFor(() => expect(mock.snapshot().status?.cutover.status, JSON.stringify({ recent: mock.ledger.slice(-2), alert: dialog.querySelector('[role="alert"]')?.textContent })).toBe('succeeded'));
    await waitFor(() => expect(client.isMutating()).toBe(0));
    cleanup(); client.clear(); mount('applied-wide', mock);
    await screen.findByTestId('recording-change-recovery');
    await waitFor(() => expect(client.isFetching()).toBe(0));
    await waitFor(() => expect(client.isMutating()).toBe(0));
    expect(mock.ledger.filter(call => call.status !== 200)).toEqual([]);
    expect(mock.snapshot().status?.operation_id).toBeDefined();
  });
});
