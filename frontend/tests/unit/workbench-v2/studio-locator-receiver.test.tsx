import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import type { WorkbenchV2ShellProps } from '@/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { buildStudioEntityLocatorTarget, parseStudioEntityLocator } from '@/services/studioEntityLocator';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import i18n from '@/i18n/config';
import { locatorRead, nativeMapping } from './studio-locator-fixtures';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
const native = vi.hoisted(() => ({ shell: false }));
vi.mock('@/features/datalink/workbench-v2/shell/WorkbenchV2Shell', async importOriginal => {
  const actual = await importOriginal<typeof import('@/features/datalink/workbench-v2/shell/WorkbenchV2Shell')>();
  return { WorkbenchV2Shell: (props: WorkbenchV2ShellProps) => {
    if (native.shell) return <actual.WorkbenchV2Shell {...props} />;
    const { state, actions } = props;
    return <div data-testid="receiver-shell"><div data-testid="receiver-devices">{state.devices.map(device => device.id).join(',')}</div>
    <div data-testid="receiver-targets">{Object.keys(state.db.targets).join(',')}</div><div data-testid="receiver-mapping">{Object.values(state.mappings).map(mapping => `${mapping.mapping_id}:${mapping.persisted_point_id}:${mapping.scale}`).join(',')}</div>
    <button onClick={() => actions.dispatch({ type: 'SET_STATE', payload: { devices: state.devices.map(device => ({ ...device, name: 'My unsaved draft' })) } })}>Keep draft</button>
    <span>{state.devices[0]?.name}</span>
  </div>;
  } };
});
const clients: QueryClient[] = [];
const response = (read: StudioLocatorEntryRead) => ({ data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } });
function Host() {
  const location = useLocation(), navigate = useNavigate();
  return <><button onClick={() => navigate(buildStudioEntityLocatorTarget(locatorRead(1, 'B').locator)!)}>Device B</button>
    <button onClick={() => navigate(-1)}>Back</button><button onClick={() => navigate(1)}>Forward</button>
    <DatalinkWorkbenchV2Page navigateTo={navigate} locatorIntent={parseStudioEntityLocator(location.search)} /></>;
}
function mount(read = locatorRead(), client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }), target?: string) {
  clients.push(client); const view = render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[target ?? buildStudioEntityLocatorTarget(read.locator)!]}><Host /></MemoryRouter></QueryClientProvider>);
  return { ...view, client };
}
function params(value: unknown): Record<string, unknown> { return value && typeof value === 'object' ? value as Record<string, unknown> : {}; }
beforeEach(async () => { vi.resetAllMocks(); native.shell = false; await i18n.changeLanguage('en'); vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValue(new Error('Legacy bootstrap must not run')); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

describe('locator receiver initialization through native parser and router', () => {
  it('rejects invalid locator before invoking legacy workspace bootstrap or child lists', async () => {
    mount(locatorRead(), undefined, '/studio/v2?entry=locator&workspace_id=W');
    await screen.findByTestId('studio-locator-notice');
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
    expect(screen.queryByTestId('receiver-shell')).not.toBeInTheDocument();
  });
  it('keeps absent intent ordinary-compatible', async () => {
    mount(locatorRead(), undefined, '/studio/v2'); await waitFor(() => expect(studioV2WorkspaceAPI.get).toHaveBeenCalledOnce());
  });
  it('hydrates detached device without writing canonical caches or triggering child GETs', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(locatorRead())); const { client } = mount();
    expect(await screen.findByTestId('receiver-devices')).toHaveTextContent(/^D$/);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledOnce(); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    expect(client.getQueryData(studioV2WorkspaceKeys.devices())).toBeUndefined();
    expect(client.getQueryData(studioV2WorkspaceKeys.bootstrap())).toBeUndefined();
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
  it.each([403, 409, 503])('typed failure %i never mounts a default entity or legacy bootstrap', async status => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValue({ response: { status } }); mount();
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).not.toHaveTextContent('Loading'));
    expect(screen.queryByTestId('receiver-shell')).not.toBeInTheDocument(); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
  });
  it('shows exact tombstone removal and only proven same-workspace device links, never a replacement editor', async () => {
    const read = locatorRead(4); read.entry.write_groups!.groups[0].status = 'deleted';
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(read)); mount(read);
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).toHaveTextContent('removed'));
    expect(screen.queryByTestId('receiver-shell')).not.toBeInTheDocument(); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    const links = screen.getAllByRole('link'); expect(links).toHaveLength(2);
    links.forEach(link => { const intent = parseStudioEntityLocator(new URL(link.getAttribute('href')!, 'https://example.invalid').search);
      expect(intent).toMatchObject({ kind: 'locator', locator: { workspace_id: 'W', step: 1 } }); });
  });
  it('Back and Forward revalidate exact IDs rather than reusing previous entity facts', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => response(locatorRead(1, params(config?.params).device_id === 'B' ? 'B' : 'D')));
    mount(); expect(await screen.findByTestId('receiver-devices')).toHaveTextContent(/^D$/);
    fireEvent.click(screen.getByText('Device B')); await waitFor(() => expect(screen.getByTestId('receiver-devices')).toHaveTextContent(/^B$/));
    fireEvent.click(screen.getByText('Back')); await waitFor(() => expect(screen.getByTestId('receiver-devices')).toHaveTextContent(/^D$/));
    fireEvent.click(screen.getByText('Forward')); await waitFor(() => expect(screen.getByTestId('receiver-devices')).toHaveTextContent(/^B$/));
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
  });
  it('a late A response cannot replace B', async () => {
    let resolveA!: (value: ReturnType<typeof response>) => void;
    vi.mocked(studioV2DatalinkApi.get).mockImplementation((_url, config) => params(config?.params).device_id === 'B'
      ? Promise.resolve(response(locatorRead(1, 'B'))) : new Promise(resolve => { resolveA = resolve; }));
    mount(); fireEvent.click(screen.getByText('Device B'));
    expect(await screen.findByTestId('receiver-devices')).toHaveTextContent(/^B$/);
    await act(async () => { resolveA(response(locatorRead())); }); expect(screen.getByTestId('receiver-devices')).toHaveTextContent(/^B$/);
  });
  it('preserves local draft and recovery storage on a new locator rather than replacing its state', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => response(locatorRead(1, params(config?.params).device_id === 'B' ? 'B' : 'D')));
    sessionStorage.setItem('wbv2.fixture-recovery', 'retained'); const before = JSON.stringify(sessionStorage);
    mount(); await screen.findByTestId('receiver-devices'); fireEvent.click(screen.getByText('Keep draft')); fireEvent.click(screen.getByText('Device B'));
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).toHaveTextContent('edits'));
    expect(JSON.stringify(sessionStorage)).toBe(before); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
  it('does not reconcile cached same-address foreign mapping into the selected native binding', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    client.setQueryData(studioV2WorkspaceKeys.mappings(), [{ ...nativeMapping('FOREIGN'), point_id: 'foreign-point', scale: 999 }]);
    client.setQueryData(studioV2WorkspaceKeys.databaseTargets(), [{ id: 'foreign-target', workspace_id: 'OTHER', point_id: 'P-native', tag_id: 'OTHER', column_name: 'foreign', enabled: true, updated_at: 'old' }]);
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(locatorRead(3))); mount(locatorRead(3), client);
    expect(await screen.findByTestId('receiver-mapping')).toHaveTextContent(/^M:P-native:0.1$/);
    expect(screen.getByTestId('receiver-targets')).toBeEmptyDOMElement();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledOnce(); expect(client.getQueryData(studioV2WorkspaceKeys.mappings())).toMatchObject([{ id: 'FOREIGN' }]);
  });
  it('unopened Step4 card navigates to the exact Step1 locator without a retention denial', async () => {
    native.shell = true;
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
      if (url !== '/studio-v2/workspace') throw { response: { status: 503 } };
      return response(params(config?.params).step === '4' ? locatorRead(4) : locatorRead(1, 'B'));
    });
    mount(locatorRead(4)); await screen.findByTestId('group-open-G'); fireEvent.click(screen.getByText('Device B'));
    expect(await screen.findByTestId('editor-input-name')).toHaveValue('Device B'); expect(screen.queryByTestId('group-editor')).not.toBeInTheDocument();
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
  });
  it('consumes selection once and shows unrequested scope after explicit native step navigation', async () => {
    native.shell = true;
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url) => {
      if (url !== '/studio-v2/workspace') throw { response: { status: 503 } }; return response(locatorRead(4));
    });
    mount(locatorRead(4)); await screen.findByTestId('group-open-G');
    fireEvent.click(screen.getByRole('button', { name: /連接設備/ }));
    expect(await screen.findByTestId('studio-locator-notice')).toHaveTextContent('have not been read');
    expect(screen.queryByTestId('group-open-G')).not.toBeInTheDocument(); expect(screen.queryByTestId('editor-input-name')).not.toBeInTheDocument();
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('retains cold opened editor edits through Step1 navigation until Close revalidates the requested target', async () => {
    native.shell = true;
    const read = locatorRead(4), opened = { ...read, entry: { ...read.entry, editor_hydrated: true as const,
      devices: [], source_rules: [], mappings: [], database_config: null, database_targets: [] } };
    let deviceReads = 0; let fresh!: (value: ReturnType<typeof response>) => void;
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
      if (url !== '/studio-v2/workspace') throw { response: { status: 503 } };
      const query = params(config?.params);
      if (query.device_id === 'B' && ++deviceReads === 2) return new Promise(resolve => { fresh = resolve; });
      return response(query.editor === '1' ? opened : query.step === '4' ? read : locatorRead(1, 'B'));
    });
    mount(read); fireEvent.click(await screen.findByTestId('group-open-G'));
    const input = await screen.findByTestId('group-name');
    expect(screen.queryByText('This group points at a different destination than the saved connection below.')).not.toBeInTheDocument();
    expect(screen.getByText(i18n.t('locatorEntry.destinationUnknown', { ns: 'runtime-dashboard' }))).toBeInTheDocument();
    fireEvent.change(input, { target: { value: 'Keep cold editor draft' } });
    fireEvent.click(screen.getByText('Device B'));
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).toHaveTextContent('edits'));
    expect(screen.getByTestId('group-name')).toBe(input); expect(input).toHaveValue('Keep cold editor draft');
    expect(screen.queryByTestId('editor-input-name')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-editor-close'));
    await waitFor(() => expect(deviceReads).toBe(2));
    expect(screen.queryByTestId('editor-input-name')).not.toBeInTheDocument(); expect(screen.getByTestId('studio-locator-notice')).toBeVisible();
    await act(async () => { fresh(response(locatorRead(1, 'B'))); });
    expect(await screen.findByTestId('editor-input-name')).toHaveValue('Device B');
    const reads = vi.mocked(studioV2DatalinkApi.get).mock.calls.filter(([url, config]) => url === '/studio-v2/workspace' && params(config?.params).device_id === 'B');
    expect(reads).toHaveLength(2); // initial pending intent, then fresh proof after explicit Close
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
  });

  it.each([403, 409])('post-Close refresh failure %i does not consume the earlier accepted target proof', async status => {
    native.shell = true;
    const read = locatorRead(4), opened = { ...read, entry: { ...read.entry, editor_hydrated: true as const,
      devices: [], source_rules: [], mappings: [], database_config: null, database_targets: [] } };
    let deviceReads = 0;
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
      if (url !== '/studio-v2/workspace') throw { response: { status: 503 } };
      const query = params(config?.params);
      if (query.device_id === 'B' && ++deviceReads === 2) throw { response: { status } };
      return response(query.editor === '1' ? opened : query.step === '4' ? read : locatorRead(1, 'B'));
    });
    mount(read); fireEvent.click(await screen.findByTestId('group-open-G')); await screen.findByTestId('group-name');
    fireEvent.click(screen.getByText('Device B')); await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).toHaveTextContent('edits'));
    fireEvent.click(screen.getByTestId('group-editor-close'));
    await waitFor(() => expect(screen.getByTestId('studio-locator-notice')).not.toHaveTextContent('edits'));
    expect(screen.queryByTestId('editor-input-name')).not.toBeInTheDocument(); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });

});
