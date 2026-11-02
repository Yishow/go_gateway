import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { RecordingHomeSourceHost, type HomeSourceControls } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2WorkspaceDatabaseAPI } from '@/services/studioV2WorkspaceDatabase';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { studioV2WorkspaceAuditAPI } from '@/services/studioV2WorkspaceAudit';
import { deviceAPI, pointAPI, modbusShareAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';
import i18n from '@/i18n/config';
vi.mock('@/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: { getConfig: vi.fn(), listTargets: vi.fn(), updateConfig: vi.fn(), upsertTarget: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceAudit', () => ({ studioV2WorkspaceAuditAPI: { list: vi.fn().mockResolvedValue({ entries: [] }) } }));
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ modbusShareAPI: { status: vi.fn().mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: false, enabled: false, port: 502, address: '127.0.0.1:502', bind_state: 'disabled', mapping_count: 0 }) }, pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn().mockResolvedValue({ device_id: 'D', activation_allowed: false }) } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
const completeDevice = () => ({ ...nativeDevice(), connection_config: '{"host":"device.example","port":502,"slave_id":1,"timeout":1500}' });
const clients: QueryClient[] = [];
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'] };
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear();
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: false, enabled: false, port: 502, address: '127.0.0.1:502', bind_state: 'disabled', mapping_count: 0 });
  vi.mocked(studioV2WorkspaceAuditAPI.list).mockResolvedValue({ entries: [] });
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: false } });
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
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
function noOtherWrites() {
  for (const fn of [studioV2WorkspaceDevicesAPI.create, studioV2WorkspaceDevicesAPI.update, studioV2RulesAPI.create,
    studioV2RulesAPI.remove, studioV2MappingsAPI.create, studioV2MappingsAPI.update,
    studioV2WorkspaceWriteGroupsAPI.create, studioV2WorkspaceWriteGroupsAPI.update, studioV2WorkspaceWriteGroupsAPI.apply]) expect(fn).not.toHaveBeenCalled();
}

async function openCreate() {
  fireEvent.click(await screen.findByRole('button', { name: 'Devices' }));
  expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Device management');
  expect(screen.getByRole('searchbox', { name: 'Search device names or protocols' })).toBeVisible();
  expect(screen.queryByRole('button', { name: 'New recording' })).not.toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'Add device' }));
  const dialog = await screen.findByRole('dialog');
  await within(dialog).findByTestId('recording-source-name');
  return dialog;
}
function noConfigurationWrites() {
  noOtherWrites(); expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
  expect(studioV2WorkspaceDatabaseAPI.updateConfig).not.toHaveBeenCalled();
  expect(studioV2WorkspaceDatabaseAPI.upsertTarget).not.toHaveBeenCalled();
  expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
}
describe('Home direct device creation through the original workspace owner', () => {
  it('mounts the actual source form without a recording draft or writes and cancels unsent input', async () => {
    const { navigate } = mount(); const dialog = await openCreate();
    expect(screen.queryByTestId('group-name')).not.toBeInTheDocument();
    expect(screen.queryByTestId('workbench-v2-root')).not.toBeInTheDocument();
    noConfigurationWrites(); expect(deviceAPI.checkReadiness).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
    fireEvent.change(within(dialog).getByTestId('recording-source-name'), { target: { value: 'Unsent device' } });
    fireEvent.click(within(dialog).getByTestId('recording-source-close'));
    fireEvent.click(await within(dialog).findByTestId('recording-source-discard'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    noConfigurationWrites(); expect(navigate).not.toHaveBeenCalled();
  });
  it('uses one original create ACK and real readback before entering the new device points form', async () => {
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async request => {
      if (!request.id) throw new Error('Expected original queue client identity');
      const saved = { ...completeDevice(), id: request.id, name: request.name, protocol: request.protocol, connection_config: JSON.stringify(request.connection_config) };
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([completeDevice(), saved]);
      vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, ordered_device_ids: ['D', saved.id], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: false } });
      vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
        expect(config?.params).toMatchObject({ workspace_id: 'W', device_id: saved.id, step: '1' });
        const read = locatorRead(1, saved.id); read.workspace = { ...workspace, ordered_device_ids: ['D', saved.id] };
        read.entry.devices = [saved];
        return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
      });
      return saved;
    });
    const { navigate } = mount(); const dialog = await openCreate();
    fireEvent.change(within(dialog).getByTestId('recording-source-name'), { target: { value: 'New Home PLC' } });
    fireEvent.change(within(dialog).getByTestId('input-host'), { target: { value: '127.0.0.1' } });
    fireEvent.click(within(dialog).getByTestId('recording-source-save'));
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce());
    await screen.findByLabelText('Starting address', { selector: 'input' });
    expect(screen.queryByTestId('group-name')).not.toBeInTheDocument();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
    expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled();
    expect(studioV2RulesAPI.create).not.toHaveBeenCalled();
    expect(navigate).not.toHaveBeenCalled();
  });
  it('keeps an uncertain create in the original locked inspect flow without a second create', async () => {
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockRejectedValue(new Error('response lost'));
    mount(); const dialog = await openCreate();
    fireEvent.change(within(dialog).getByTestId('recording-source-name'), { target: { value: 'Uncertain PLC' } });
    fireEvent.change(within(dialog).getByTestId('input-host'), { target: { value: '127.0.0.1' } });
    fireEvent.click(within(dialog).getByTestId('recording-source-save'));
    await within(dialog).findByTestId('recording-source-unknown');
    expect(within(dialog).getByTestId('recording-source-save')).toBeDisabled();
    expect(within(dialog).getByTestId('recording-source-inspect')).toBeEnabled();
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
    expect(screen.queryByLabelText('Starting address', { selector: 'input' })).not.toBeInTheDocument();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
});

 it('rejects a late create callback before touching the current modal or its close owner', () => {
  const client = new QueryClient(); clients.push(client);
  let old: HomeSourceControls | undefined;
  const currentClose = vi.fn();
  render(<QueryClientProvider client={client}><RecordingHomeSourceHost renderSource={controls => {
    if (controls.target.kind === 'create') old ??= controls;
    if (controls.target.kind === 'connection') controls.closeRef.current = currentClose;
    return <p data-testid="selected-source">{controls.target.kind}</p>;
  }}>{(open, workspaceChanged, create) => <>
    <button onClick={() => { workspaceChanged('W'); create('W'); }}>Create owned</button>
    <button onClick={() => open('W', 'D', 'connection')}>Open existing</button>
    <button onClick={() => workspaceChanged('OTHER')}>Switch workspace</button>
  </>}</RecordingHomeSourceHost></QueryClientProvider>);
  fireEvent.click(screen.getByText('Create owned'));
  fireEvent.click(screen.getByText('Open existing'));
  expect(old).toBeDefined();
  old!.onCreated('LATE');
  expect(screen.getByTestId('selected-source')).toHaveTextContent('connection');
  fireEvent.click(screen.getByTestId('home-source-close'));
  expect(currentClose).toHaveBeenCalledOnce();
  fireEvent.click(screen.getByText('Switch workspace'));
  old!.onCreated('FOREIGN');
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
});
