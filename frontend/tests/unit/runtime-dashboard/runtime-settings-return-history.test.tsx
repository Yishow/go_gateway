import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeDashboardRoute } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const control = vi.hoisted(() => ({ revision: 'W1', ready: true, live: false }));
vi.mock('../../../src/features/datalink/runtime-dashboard/HistoryReportsPanel', () => ({ HistoryReportsPanel: () => <div>History reports</div> }));
vi.mock('../../../src/features/datalink/runtime-dashboard/useRuntimeDashboardState', async () => {
  const router = await import('react-router-dom');
  return { useRuntimeDashboardState: () => {
    const [params, setParams] = router.useSearchParams(); const device = params.get('device_id') ?? 'D';
    const devices = ['D', 'E'].map((device_id) => ({ device_id, name: device_id, protocol: 'modbus_tcp', running: false, availability_status: 'available' as const }));
    return { routeState: control.live ? 'live' : 'loading', selectedDeviceId: device, selectedDevice: devices.find((item) => item.device_id === device), devices,
      sourceEvidenceScope: control.ready ? { workspace_id: 'W', workspace_revision: control.revision, device_id: device } : null, setupContext: { source_rules: [], mappings: [], database_targets: [] },
      snapshot: null, snapshotError: null, liveValues: {}, streamState: 'connecting', onSelectDevice: (id: string) => setParams({ device_id: id }), onRetrySnapshot: async () => undefined, onReconnectStream: () => undefined };
  } };
});
function requestParams(value: unknown): Record<string, unknown> {
  if (value === null || typeof value !== 'object' || Array.isArray(value)) throw new Error('Expected scoped request parameters');
  return Object.fromEntries(Object.entries(value));
}
const clients: QueryClient[] = [];
function Navigation() {
  const navigate = useNavigate(); const location = useLocation();
  return <><button onClick={() => navigate(-1)}>Back</button><button onClick={() => navigate(1)}>Forward</button><button onClick={() => navigate('/studio/runtime?device_id=E')}>Device E</button><button onClick={() => { control.ready = true; navigate(location.pathname + location.search, { replace: true, state: location.state }); }}>Proof arrives</button><output data-testid="location">{location.pathname}{location.search}</output><output data-testid="state">{JSON.stringify(location.state)}</output></>;
}
function mount(initialState: unknown = null) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  return render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[{ pathname: '/studio/runtime', search: '?device_id=D', state: initialState }]}><Navigation /><Routes><Route path="/studio/runtime" element={<RuntimeDashboardRoute />} /><Route path="/studio/v2" element={<h1>Settings receiver</h1>} /></Routes></MemoryRouter></QueryClientProvider>);
}
function response(params: Record<string, unknown>) {
  const scope = { workspace_id: 'W', workspace_revision: String(params.workspace_revision), device_id: String(params.device_id) };
  const suffix = scope.device_id;
  const view = ruleEvidenceFixture(scope); view.groups = [];
  if (params.rule_id) view.source = { rule_id: `R${suffix}`, saved_revision: 'R1', point_id: `P${suffix}`, mapping_id: `M${suffix}`, tag_id: null, effective_revision: null, state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: 'W', workspace_revision: scope.workspace_revision, entity_revision: 'R1' }, mapping_revision: 'MR1', enabled: true, address: '40001', protocol_function: '03', read_type: 'int16', result_type: 'int16', polling: null, transforms: [] } };
  return { data: { success: true, data: { workspace_id: 'W', ...(params.source_list ? { source_page: {
    contract_version: 1, ...scope, kind: params.source_list, query: '', limit: 50, list_revision: 'L1', next_cursor: null,
    rules: params.source_list === 'rules' ? [{ id: `R${suffix}`, revision_id: 'R1', device_id: suffix, start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'Same name', enabled: true }] : [],
    mappings: params.source_list === 'mappings' ? [{ id: `M${suffix}`, rule_id: `R${suffix}`, rule_revision: 'R1', device_id: suffix, point_id: `P${suffix}`, address: '40001', enabled: true }] : [],
  } } : { rule_evidence: view }) } } };
}
async function openDetail() {
  const expand = screen.queryByRole('button', { name: '展開配置詳情' }); if (expand) fireEvent.click(expand);
  fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
  return screen.findByRole('region', { name: 'Source rule details' });
}
beforeEach(async () => {
  control.revision = 'W1'; control.ready = true; control.live = false; await i18n.changeLanguage('en'); vi.mocked(studioV2DatalinkApi.get).mockReset();
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => response(requestParams(config?.params)));
});
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });
describe('real runtime router settings history', () => {
  it('Back and Forward preserve exact device and explicit source, storing no facts', async () => {
    mount(); const detail = await openDetail(); await within(detail).findByText('Saved-only configuration');
    expect(screen.getByTestId('state').textContent).toBe('{"runtimeEntityHistory":{"version":1,"workspace_id":"W","workspace_revision":"W1","device_id":"D","selection":{"rule_id":"RD","rule_revision":"R1","point_id":"PD","mapping_id":"MD"}}}');
    fireEvent.click(within(detail).getByRole('link', { name: 'Open mapping settings' }));
    await screen.findByRole('heading', { name: 'Settings receiver' }); expect(screen.getByTestId('location')).toHaveTextContent('/studio/v2?entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=MD');
    fireEvent.click(screen.getByRole('button', { name: 'Back' }));
    const restored = await screen.findByRole('region', { name: 'Source rule details' });
    await within(restored).findByText('Saved-only configuration'); expect(restored).toBeVisible(); expect(within(restored).getByText('MD')).toBeInTheDocument();
    expect(screen.getByTestId('location')).toHaveTextContent('/studio/runtime?device_id=D');
    fireEvent.click(screen.getByRole('button', { name: 'Forward' })); await screen.findByRole('heading', { name: 'Settings receiver' });
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.every(([url]) => url === '/studio-v2/workspace/runtime-context')).toBe(true);
  });
  it('shows a restored explicit selection when current scope arrives after the live page mounts', async () => {
    control.ready = false; control.live = true;
    mount({ runtimeEntityHistory: { version: 1, workspace_id: 'W', workspace_revision: 'W1', device_id: 'D', selection: { rule_id: 'RD', rule_revision: 'R1', point_id: 'PD', mapping_id: 'MD' } } });
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Proof arrives' }));
    const detail = await screen.findByRole('region', { name: 'Source rule details' });
    await within(detail).findByText('Saved-only configuration'); expect(detail.closest('.hidden')).toBeNull(); expect(detail).toBeVisible();
  });
  it('does not restore a revision-stale source after returning from settings', async () => {
    mount(); const detail = await openDetail(); await within(detail).findByText('Saved-only configuration');
    fireEvent.click(within(detail).getByRole('link', { name: 'Open rule settings' })); await screen.findByRole('heading', { name: 'Settings receiver' });
    control.revision = 'W2'; const calls = vi.mocked(studioV2DatalinkApi.get).mock.calls.length;
    fireEvent.click(screen.getByRole('button', { name: 'Back' })); await screen.findByTestId('runtime-dashboard-route');
    await waitFor(() => expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.length).toBeGreaterThan(calls));
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.slice(calls).some(([, config]) => requestParams(config?.params).rule_id)).toBe(false);
  });
  it('revalidates a prior source on Back and keeps deleted/stale results unavailable without selecting a namesake', async () => {
    mount(); const detail = await openDetail(); await within(detail).findByText('Saved-only configuration');
    fireEvent.click(within(detail).getByRole('link', { name: 'Open mapping settings' })); await screen.findByRole('heading', { name: 'Settings receiver' });
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => { if (requestParams(config?.params).rule_id) throw { response: { data: { error: { code: 'stale-reference', message: 'private.example' } } } }; return response(requestParams(config?.params)); });
    fireEvent.click(screen.getByRole('button', { name: 'Back' }));
    const restored = await screen.findByRole('region', { name: 'Source rule details' }); await within(restored).findByRole('alert');
    expect(within(restored).queryByRole('link')).not.toBeInTheDocument(); expect(document.body.textContent).not.toContain('private.example');
  });
  it.each(['stale-reference', 'foreign-scope', 'unavailable'])('keeps an outside-page prior selection %s without namesake fallback or detail fan-out', async (code) => {
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      if (requestParams(config?.params).rule_id) throw { response: { data: { error: { code, message: '/private/password' } } } };
      return response(requestParams(config?.params));
    });
    mount({ runtimeEntityHistory: { version: 1, workspace_id: 'W', workspace_revision: 'W1', device_id: 'D', selection: { rule_id: 'Removed', rule_revision: 'Old', point_id: 'OldPoint', mapping_id: 'OldMapping' } } });
    const detail = await screen.findByRole('region', { name: 'Source rule details' }); await within(detail).findByRole('alert');
    expect(within(detail).queryByRole('link')).not.toBeInTheDocument(); expect(within(detail).queryByText('MD')).not.toBeInTheDocument();
    const reads = vi.mocked(studioV2DatalinkApi.get).mock.calls.filter(([, config]) => requestParams(config?.params).rule_id);
    expect(reads).toHaveLength(1); expect(reads[0][1]?.params).toMatchObject({ rule_id: 'Removed', mapping_id: 'OldMapping' });
    expect(document.body.textContent).not.toContain('/private/password');
  });
  it('clears an ordinary page selection when its current page becomes stale instead of retaining old detail facts', async () => {
    mount(); const detail = await openDetail(); await within(detail).findByText('Saved-only configuration');
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      if (requestParams(config?.params).source_list) throw { response: { data: { error: { code: 'stale-reference' } } } };
      return response(requestParams(config?.params));
    });
    await act(async () => { await clients[0].invalidateQueries({ predicate: (query) => query.queryKey.includes('source-page') }); });
    await waitFor(() => expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument());
    await waitFor(() => expect(screen.getByTestId('state')).toHaveTextContent('"selection":null'));
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.filter(([, config]) => requestParams(config?.params).rule_id)).toHaveLength(1);
  });
  it('does not let a late D detail overwrite selected E or its history', async () => {
    let resolve!: (value: ReturnType<typeof response>) => void; let delayedParams: Record<string, unknown> = {}; let signal: AbortSignal | undefined;
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const params = requestParams(config?.params);
      if (params.rule_id && params.device_id === 'D') {
        if (!(config?.signal instanceof AbortSignal)) throw new Error('Expected native cancellation signal');
        delayedParams = params; signal = config.signal; return new Promise((done) => { resolve = done; });
      }
      return response(requestParams(config?.params));
    });
    mount(); await openDetail(); await waitFor(() => expect(resolve).toBeDefined()); fireEvent.click(screen.getByRole('button', { name: 'Device E' }));
    await waitFor(() => expect(signal?.aborted).toBe(true)); const detail = await openDetail(); await within(detail).findByText('ME');
    await act(async () => { resolve(response(delayedParams)); });
    expect(within(detail).queryByText('MD')).not.toBeInTheDocument(); expect(screen.getByTestId('state')).toHaveTextContent('"device_id":"E"'); expect(screen.getByTestId('state')).not.toHaveTextContent('"mapping_id":"MD"');
  });
});
