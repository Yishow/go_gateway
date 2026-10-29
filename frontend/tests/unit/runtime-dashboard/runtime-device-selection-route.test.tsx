import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation, useNavigate } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '../../../src/services/datalinkClient';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { RuntimeDashboardRoute } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardRoute';
import { useRuntimeDashboardState } from '../../../src/features/datalink/runtime-dashboard/useRuntimeDashboardState';
import { usePointsQuery } from '../../../src/hooks/datalink/usePoints';

vi.mock('../../../src/services/datalinkClient', async (original) => ({ ...(await original<typeof import('../../../src/services/datalinkClient')>()), api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, fallback?: unknown) => typeof fallback === 'string' ? fallback : key }) }));
const events: string[] = [];
class EventSourceMock {
  onopen = null; onerror = null;
  constructor(url: string) { events.push(url); }
  addEventListener() {} removeEventListener() {} close() {}
}
const context = (ids = ['device-A', 'device-B'], defaultId: string | null = 'device-A') => ({ workspace_id: 'W', default_device_id: defaultId, devices: ids.map(device_id => ({ device_id, name: `Meter ${device_id}`, protocol: 'modbus_tcp', availability_status: 'available', running: false })), setup: null });
const response = (data: unknown) => ({ data: { success: true, data } });
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(done => { resolve = done; }); return { promise, resolve }; }
let currentSearch = '';
let client: QueryClient;
function LocationProbe() { currentSearch = useLocation().search; return null; }
function Controls() {
  const state = useRuntimeDashboardState(); const navigate = useNavigate();
  return <><output data-testid="selection">{state.selectedDeviceId ?? ''}</output><output data-testid="route-state">{state.routeState}</output>
    <button onClick={() => void state.onRetrySnapshot()}>Retry snapshot</button>
    <button onClick={() => state.onSelectDevice('private.example:5432')}>Unsafe select</button>
    <button onClick={() => state.onSelectDevice('foreign-device')}>Foreign select</button>
    <button onClick={() => state.onSelectDevice('device-B')}>Owned select</button>
    <button onClick={() => void navigate('?device_id=device-B&keep=x')}>Go B</button>
    <button onClick={() => void navigate('?device_id=')}>Go invalid</button></>;
}
function mount(search = '', controls = false) {
  client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[`/studio/runtime${search}`]}><Routes><Route path="/studio/runtime" element={<>{controls ? <Controls /> : <RuntimeDashboardRoute />}<LocationProbe /></>} /></Routes></MemoryRouter></QueryClientProvider>);
}
beforeEach(() => {
  vi.clearAllMocks(); events.length = 0; currentSearch = ''; vi.stubGlobal('EventSource', EventSourceMock);
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(context()));
  vi.mocked(api.get).mockImplementation(async (url) => response(url === '/points' ? [] : { collectors: [], running: false, metrics: {}, snapshot_state: { empty: true } }));
});
afterEach(() => { cleanup(); client?.clear(); vi.unstubAllGlobals(); });

describe('runtime device selector at the actual route and transport boundary', () => {
  it('never reflects or requests a private explicit selector while context is deferred', async () => {
    const pending = deferred<ReturnType<typeof response>>(); vi.mocked(studioV2DatalinkApi.get).mockReturnValue(pending.promise);
    const secret = 'postgres://operator:password@private.example/db'; const search = `?keep=a%20b&device_id=${encodeURIComponent(secret)}`;
    const view = mount(search);
    expect(view.container.innerHTML).not.toContain(secret);
    expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]); expect(currentSearch).toBe(search);
    await act(async () => pending.resolve(response(context())));
    expect(await screen.findByTestId('runtime-dashboard-missing-device-context')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-selected-device')).toBeEmptyDOMElement();
    expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]); expect(currentSearch).toBe(search);
  });
  it('waits for matching owned context even for a syntactically valid selector', async () => {
    const pending = deferred<ReturnType<typeof response>>(); vi.mocked(studioV2DatalinkApi.get).mockReturnValue(pending.promise);
    mount('?device_id=device-A'); expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]);
    await act(async () => pending.resolve(response(context())));
    await waitFor(() => expect(api.get).toHaveBeenCalledWith('/runtime/status', { params: { device_id: 'device-A' } }));
    expect(api.get).toHaveBeenCalledWith('/points', { params: { device_id: 'device-A' } });
    expect(events.every(url => new URL(url, 'http://test').searchParams.get('device_id') === 'device-A')).toBe(true);
  });
  it.each(['?device_id=', '?device_id', '?device_id=%ZZ', '?device_id=%E0%A4%A', '?device_id=device-A&device_id=device-B', '?device_id=device-A&%64evice_id=device-B', '?device_id=%2Fprivate%2Fdb', '?device_id=foreign-device'])('fails closed without default/unfiltered reads or URL rewriting: %s', async (search) => {
    mount(search, true); await waitFor(() => expect(screen.getByTestId('route-state')).toHaveTextContent('missing-device-context'));
    expect(screen.getByTestId('selection')).toBeEmptyDOMElement();
    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' }));
    await act(async () => { await Promise.resolve(); });
    expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]); expect(currentSearch).toBe(search);
    expect(api.post).not.toHaveBeenCalled(); expect(api.put).not.toHaveBeenCalled(); expect(api.delete).not.toHaveBeenCalled();
  });
  it.each(['', '?keep=a%20b', '?device_id=device-B&keep=a%20b'])('retains absent default and owned explicit selection without mount rewrite: %s', async search => {
    mount(search); const id = search.includes('device_id') ? 'device-B' : 'device-A';
    await waitFor(() => expect(api.get).toHaveBeenCalledWith('/runtime/status', { params: { device_id: id } }));
    expect(currentSearch).toBe(search); await waitFor(() => expect(screen.getByTestId('runtime-dashboard-selected-device')).toHaveTextContent(`Meter ${id}`));
  });
  it('does not allow unsafe or foreign callbacks to rewrite selection; owned selection remains native', async () => {
    mount('?device_id=device-A&keep=x', true); await waitFor(() => expect(api.get).toHaveBeenCalledWith('/points', { params: { device_id: 'device-A' } }));
    const calls = vi.mocked(api.get).mock.calls.length;
    fireEvent.click(screen.getByRole('button', { name: 'Unsafe select' })); fireEvent.click(screen.getByRole('button', { name: 'Foreign select' }));
    expect(currentSearch).toBe('?device_id=device-A&keep=x'); expect(api.get).toHaveBeenCalledTimes(calls);
    fireEvent.click(screen.getByRole('button', { name: 'Owned select' }));
    await waitFor(() => expect(api.get).toHaveBeenCalledWith('/runtime/status', { params: { device_id: 'device-B' } }));
    expect(currentSearch).toBe('?device_id=device-B&keep=x');
  });
  it('does not fetch global points or status when the observed workspace has no default', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(context([], null))); mount('', true);
    await waitFor(() => expect(screen.getByTestId('route-state')).toHaveTextContent('empty-workspace'));
    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' })); await act(async () => { await Promise.resolve(); });
    expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]);
  });
  it('binds a delayed context to the current route rather than the earlier selector', async () => {
    const pending = deferred<ReturnType<typeof response>>(); vi.mocked(studioV2DatalinkApi.get).mockReturnValue(pending.promise);
    mount('?device_id=device-A', true);
    fireEvent.click(screen.getByRole('button', { name: 'Go B' }));
    expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]);
    await act(async () => pending.resolve(response(context())));
    await waitFor(() => expect(api.get).toHaveBeenCalledWith('/runtime/status', { params: { device_id: 'device-B' } }));
    expect(api.get).not.toHaveBeenCalledWith('/runtime/status', { params: { device_id: 'device-A' } });
    expect(api.get).not.toHaveBeenCalledWith('/points', { params: { device_id: 'device-A' } });
    expect(events.every(url => new URL(url, 'http://test').searchParams.get('device_id') === 'device-B')).toBe(true);
    fireEvent.click(screen.getByRole('button', { name: 'Go invalid' }));
    const requests = vi.mocked(api.get).mock.calls.length;
    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' }));
    await act(async () => { await Promise.resolve(); });
    expect(screen.getByTestId('selection')).toBeEmptyDOMElement(); expect(api.get).toHaveBeenCalledTimes(requests);
    expect(currentSearch).toBe('?device_id=');
  });
  it('does not turn failed context or a manual retry into global runtime transport', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValue(new Error('private endpoint failure'));
    const view = mount('?device_id=device-A', true);
    await waitFor(() => expect(screen.getByTestId('route-state')).toHaveTextContent('error'));
    fireEvent.click(screen.getByRole('button', { name: 'Retry snapshot' }));
    fireEvent.click(screen.getByRole('button', { name: 'Owned select' }));
    await act(async () => { await Promise.resolve(); });
    expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]);
    expect(currentSearch).toBe('?device_id=device-A'); expect(view.container.innerHTML).not.toContain('private endpoint failure');
  });
  it('does not infer first-device ownership from an absent or unsafe observed default', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(context(['device-A'], 'host.example')));
    mount('', true); await waitFor(() => expect(screen.getByTestId('route-state')).toHaveTextContent('empty-workspace'));
    expect(api.get).not.toHaveBeenCalled(); expect(events).toEqual([]); expect(currentSearch).toBe('');
  });
  it('keeps the ordinary unfiltered points hook enabled by default', async () => {
    function Ordinary() { usePointsQuery(); return null; }
    client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(<QueryClientProvider client={client}><Ordinary /></QueryClientProvider>);
    await waitFor(() => expect(api.get).toHaveBeenCalledWith('/points', { params: undefined }));
  });
});
