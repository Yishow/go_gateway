import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { deviceAPI, pointAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';
import i18n from '@/i18n/config';
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn().mockResolvedValue({ device_id: 'D', activation_allowed: false }) } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
const completeDevice = () => ({ ...nativeDevice(), connection_config: '{"host":"device.example","port":502,"slave_id":1,"timeout":1500}' });
const clients: QueryClient[] = [];
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'] };
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear();
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: 'D', status: 'ready', connect_status: 'success', probe_status: 'success', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([completeDevice()]);
  vi.mocked(studioV2RulesAPI.list).mockResolvedValue([nativeRule()]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]); vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error('Unexpected unscoped child read');
    const params = config?.params as Record<string, string>;
    const read = locatorRead(params.step === '2' ? 2 : 1); read.workspace = workspace; read.entry.devices = [completeDevice()];
    return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
  });
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({ ...nativeRule(), ...request, revision_id: 'R3', runtime_apply_status: 'not_running' }));
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  const navigate = vi.fn();
  const view = render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '#home' }} navigateTo={navigate} /></QueryClientProvider>);
  return { ...view, navigate, client };
}
async function openRange() {
  fireEvent.click(await screen.findByTestId('home-device-points-D'));
  return screen.findByRole('dialog');
}
function noOtherWrites() {
  for (const fn of [studioV2WorkspaceDevicesAPI.create, studioV2WorkspaceDevicesAPI.update, studioV2RulesAPI.create,
    studioV2RulesAPI.remove, studioV2MappingsAPI.create, studioV2MappingsAPI.update,
    studioV2WorkspaceWriteGroupsAPI.create, studioV2WorkspaceWriteGroupsAPI.update, studioV2WorkspaceWriteGroupsAPI.apply]) expect(fn).not.toHaveBeenCalled();
}
describe('direct source entry and one shared form presentation', () => {
  it.each(['points', 'connection'] as const)('can close the %s loading cue before a late locator response', async kind => {
    let release!: () => void, started = false;
    const gate = new Promise<void>(resolve => { release = resolve; });
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
      if (url !== '/studio-v2/workspace') throw new Error('Unexpected unscoped child read');
      started = true; await gate;
      const params = config?.params as Record<string, string>;
      const read = locatorRead(params.step === '2' ? 2 : 1); read.workspace = workspace; read.entry.devices = [completeDevice()];
      return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
    });
    mount();
    const points = await screen.findByTestId('home-device-points-D');
    const trigger = kind === 'points' ? points : within(screen.getByTestId('home-device-summary')).getByRole('button', { name: 'Connection' });
    fireEvent.click(trigger);
    const dialog = await screen.findByRole('dialog');
    await waitFor(() => expect(started).toBe(true));
    const footer = dialog.querySelector('footer');
    expect(footer).not.toBeNull();
    fireEvent.click(within(footer as HTMLElement).getByRole('button', { name: 'Close' }));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    await act(async () => release());
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); noOtherWrites();
  });
  it('keeps one neutral loading state until the automatically selected range is ready', async () => {
    let releaseDevice!: () => void, releaseRule!: () => void, releaseList!: (rules: ReturnType<typeof nativeRule>[]) => void;
    const deviceGate = new Promise<void>(resolve => { releaseDevice = resolve; });
    const ruleGate = new Promise<void>(resolve => { releaseRule = resolve; });
    const listGate = new Promise<ReturnType<typeof nativeRule>[]>(resolve => { releaseList = resolve; });
    const requested: string[] = [];
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
      if (url !== '/studio-v2/workspace') throw new Error('Unexpected unscoped child read');
      const params = config?.params as Record<string, string>, step = params.step;
      requested.push(step);
      await (step === '2' ? ruleGate : deviceGate);
      const read = locatorRead(step === '2' ? 2 : 1); read.workspace = workspace; read.entry.devices = [completeDevice()];
      return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
    });
    mount();
    const trigger = await screen.findByTestId('home-device-points-D');
    let listStarted = false;
    vi.mocked(studioV2RulesAPI.list).mockImplementationOnce(() => { listStarted = true; return listGate; });
    fireEvent.click(trigger);
    const dialog = await screen.findByRole('dialog');
    await waitFor(() => expect(requested).toContain('1'));
    const expectLoading = () => {
      expect(within(dialog).getByRole('status')).toHaveTextContent('Reading the current workspace');
      expect(within(dialog).queryByTestId('studio-locator-notice')).not.toBeInTheDocument();
      expect(within(dialog).queryByRole('region', { name: 'Read ranges' })).not.toBeInTheDocument();
      expect(within(dialog).queryByTestId('rule-prefix-input')).not.toBeInTheDocument();
    };
    expectLoading();
    await act(async () => releaseDevice());
    await waitFor(() => expect(listStarted).toBe(true));
    expectLoading();
    await act(async () => releaseList([nativeRule()]));
    await waitFor(() => expect(requested).toContain('2'));
    expectLoading();
    await act(async () => releaseRule());
    expect(await within(dialog).findByTestId('rule-prefix-input')).toHaveValue('Rule R');
    noOtherWrites();
  });
  it('still shows a real locator failure and allows closing the source dialog', async () => {
    mount();
    const trigger = await screen.findByTestId('home-device-points-D');
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce(new Error('fixture unavailable'));
    fireEvent.click(trigger);
    const dialog = await screen.findByRole('dialog');
    const notice = await within(dialog).findByTestId('studio-locator-notice');
    expect(notice).toHaveAttribute('role', 'alert');
    fireEvent.click(within(dialog).getByTestId('home-source-close'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    noOtherWrites();
  });
  it('shows a real connection locator failure and allows closing', async () => {
    mount();
    await screen.findByTestId('home-device-points-D');
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce(new Error('fixture unavailable'));
    fireEvent.click(within(screen.getByTestId('home-device-summary')).getByRole('button', { name: 'Connection' }));
    const dialog = await screen.findByRole('dialog');
    expect(await within(dialog).findByTestId('studio-locator-notice')).toHaveAttribute('role', 'alert');
    fireEvent.click(within(dialog).getByTestId('home-source-close'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    noOtherWrites();
  });
  it('opens the only saved range directly, then respects an explicit return to the list', async () => {
    const { navigate } = mount(); const dialog = await openRange();
    const form = await within(dialog).findByTestId('home-saved-range-form');
    expect(within(form).getByTestId('rule-prefix-input')).toHaveValue('Rule R');
    expect(within(dialog).queryByTestId('home-range-edit-R')).not.toBeInTheDocument();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace', expect.objectContaining({ params: expect.objectContaining({ step: '2', device_id: 'D', rule_id: 'R', workspace_id: 'W' }) }));
    fireEvent.click(within(dialog).getByTestId('home-range-done'));
    await within(dialog).findByTestId('home-range-edit-R');
    expect(within(dialog).queryByTestId('home-saved-range-form')).not.toBeInTheDocument();
    expect(navigate).not.toHaveBeenCalled(); noOtherWrites(); expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
  });
  it('opens the first new range directly with visible examples and the actual prefix field', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: false } });
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    mount(); const dialog = await openRange();
    const address = await within(dialog).findByLabelText('Starting address', { selector: 'input' });
    await waitFor(() => expect(address).toHaveAttribute('placeholder', '40001'));
    expect(within(dialog).getByLabelText('Point-name prefix', { selector: 'input' })).toHaveValue('');
    expect(within(dialog).queryByLabelText('Range name')).not.toBeInTheDocument();
    expect(within(dialog).getByText('For example, 16 values. Count data items, not register words.')).toBeVisible();
    expect(within(dialog).getByText('Use the data type specified by the device register map.')).toBeVisible();
    expect(within(dialog).queryByTestId('home-range-new')).not.toBeInTheDocument();
    noOtherWrites();
  });
  it('keeps multiple ranges as an explicit choice and uses the existing dirty confirmation on return', async () => {
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([nativeRule(), nativeRule('R2')]);
    mount(); const dialog = await openRange();
    fireEvent.click(await within(dialog).findByTestId('home-range-edit-R'));
    const prefix = await within(dialog).findByTestId('rule-prefix-input');
    fireEvent.change(prefix, { target: { value: 'UNSENT_' } });
    fireEvent.click(within(dialog).getByTestId('home-range-done'));
    fireEvent.click(await within(dialog).findByTestId('home-source-keep'));
    expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('UNSENT_');
    expect(studioV2RulesAPI.update).not.toHaveBeenCalled(); noOtherWrites();
  });
});
