import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { deviceAPI, pointAPI, tagAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';
import i18n from '@/i18n/config';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn() } }));
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
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
async function editRange() {
  const dialog = await openRange();
  return within(dialog).findByTestId('home-saved-range-form');
}
function noOtherWrites() {
  for (const fn of [studioV2WorkspaceDevicesAPI.create, studioV2WorkspaceDevicesAPI.update, studioV2RulesAPI.create,
    studioV2RulesAPI.remove, studioV2MappingsAPI.create, studioV2MappingsAPI.update,
    studioV2WorkspaceWriteGroupsAPI.create, studioV2WorkspaceWriteGroupsAPI.update, studioV2WorkspaceWriteGroupsAPI.apply]) expect(fn).not.toHaveBeenCalled();
}
describe('Home uses the original scoped source owner inside the demo modal', () => {
  it('opens the captured device and existing ranges without old navigation or writes', async () => {
    const { navigate } = mount(); const dialog = await openRange();
    expect(await within(dialog).findByTestId('rule-prefix-input')).toHaveValue('Rule R'); expect(navigate).not.toHaveBeenCalled();
    expect(screen.queryByTestId('workbench-v2-root')).not.toBeInTheDocument();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace', expect.objectContaining({ params: expect.objectContaining({ step: '1', device_id: 'D', workspace_id: 'W' }) }));
    expect(studioV2RulesAPI.update).not.toHaveBeenCalled(); noOtherWrites();
  });
  it('keeps local edits unsent and discards only after the explicit dirty confirmation', async () => {
    mount(); const form = await editRange(); const name = within(form).getByTestId('rule-prefix-input');
    fireEvent.change(name, { target: { value: 'Unsent range' } }); expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('home-source-close'));
    fireEvent.click(await screen.findByTestId('home-source-keep')); expect(name).toHaveValue('Unsent range');
    fireEvent.click(screen.getByTestId('home-source-close')); fireEvent.click(await screen.findByTestId('home-source-discard'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument()); expect(studioV2RulesAPI.update).not.toHaveBeenCalled(); noOtherWrites();
  });
  it('submits one exact original rule update and keeps the range modal open after ACK', async () => {
    mount(); const form = await editRange(); fireEvent.change(within(form).getByTestId('rule-prefix-input'), { target: { value: 'Saved range' } });
    expect(studioV2RulesAPI.update).not.toHaveBeenCalled(); fireEvent.click(screen.getByTestId('home-range-save'));
    await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledOnce());
    expect(studioV2RulesAPI.update).toHaveBeenCalledWith('R', expect.objectContaining({ naming_prefix: 'Saved range', device_id: 'D' }));
    await screen.findByTestId('home-range-saved'); expect(screen.getByRole('dialog')).toBeInTheDocument();
    expect(screen.getByTestId('rule-prefix-input')).toHaveValue('Saved range');
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([{ ...nativeRule(), naming_prefix: 'Saved range', revision_id: 'R3' }]);
    fireEvent.click(screen.getByTestId('home-range-done'));
    await waitFor(() => expect(screen.queryByTestId('home-saved-range-form')).not.toBeInTheDocument());
    expect(screen.getByRole('dialog')).toBeInTheDocument(); expect((await screen.findByTestId('home-range-edit-R')).closest('li')).toHaveTextContent('Saved range');
    expect(studioV2RulesAPI.update).toHaveBeenCalledOnce(); noOtherWrites();
  });
  it('does not resend or pretend to cancel a pending or unknown original save', async () => {
    let reject!: (error: Error) => void; vi.mocked(studioV2RulesAPI.update).mockImplementation(() => new Promise((_resolve, fail) => { reject = fail; }));
    mount(); const form = await editRange(); fireEvent.change(within(form).getByTestId('rule-prefix-input'), { target: { value: 'Pending range' } });
    fireEvent.click(screen.getByTestId('home-range-save')); fireEvent.click(screen.getByTestId('home-range-save'));
    await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledOnce()); fireEvent.click(screen.getByTestId('home-source-close'));
    expect(screen.getByRole('dialog')).toBeInTheDocument(); expect(screen.queryByTestId('home-source-discard')).not.toBeInTheDocument();
    await act(async () => reject(new Error('acknowledgement unavailable'))); await screen.findByTestId('home-range-unknown');
    expect(screen.getByTestId('home-range-save')).toBeDisabled(); expect(studioV2RulesAPI.update).toHaveBeenCalledOnce(); noOtherWrites();
  });
});

describe('modal exact identity and original completion boundaries', () => {
  it('uses the Home guidance points action without an advanced route', async () => {
    const { navigate } = mount(); await screen.findByTestId('home-device-points-D');
    fireEvent.click(within(screen.getByTestId('recording-home-guidance')).getByRole('button', { name: 'Open setup' }));
    await waitFor(() => expect(screen.getByRole('dialog')).toHaveTextContent('Device D')); expect(navigate).not.toHaveBeenCalled(); noOtherWrites();
  });
  it('never exposes foreign range labels and ignores a closed old list after reopening the same device', async () => {
    mount(); await screen.findByTestId('home-device-points-D');
    let resolve!: (value: ReturnType<typeof nativeRule>[]) => void;
    vi.mocked(studioV2RulesAPI.list).mockImplementationOnce(() => new Promise(done => { resolve = done; }));
    await openRange(); await waitFor(() => expect(resolve).toBeTypeOf('function'));
    fireEvent.click(screen.getByTestId('home-source-close')); await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([nativeRule(), nativeRule('OTHER-DEVICE', 'B'), { ...nativeRule('FOREIGN'), workspace_id: 'OTHER' }]);
    const dialog = await openRange(); await within(dialog).findByTestId('rule-prefix-input');
    await act(async () => resolve([nativeRule('OLD')]));
    expect(within(dialog).queryByText('Rule OLD')).not.toBeInTheDocument(); expect(within(dialog).queryByText('Rule FOREIGN')).not.toBeInTheDocument();
    expect(within(dialog).queryByText('Rule OTHER-DEVICE')).not.toBeInTheDocument(); noOtherWrites();
  });
  it('rejects a known stale range revision before dispatching the original update', async () => {
    mount(); const form = await editRange(); fireEvent.change(within(form).getByTestId('rule-prefix-input'), { target: { value: 'Stale edit' } });
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async () => {
      const read = locatorRead(2); read.workspace = workspace; read.entry.source_rules![0].revision_id = 'R3'; read.entry.rule_heads.R = 'R3';
      return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
    });
    fireEvent.click(screen.getByTestId('home-range-save')); await screen.findByText(/saved range changed or could not be verified/i);
    expect(screen.getByTestId('rule-prefix-input')).toHaveValue('Stale edit'); expect(studioV2RulesAPI.update).not.toHaveBeenCalled(); noOtherWrites();
  });
  it('never lets a late old completion publish into a newly opened modal', async () => {
    let resolve!: (value: Awaited<ReturnType<typeof studioV2RulesAPI.update>>) => void;
    vi.mocked(studioV2RulesAPI.update).mockImplementation(() => new Promise(done => { resolve = done; }));
    const first = mount(); const form = await editRange(); fireEvent.change(within(form).getByTestId('rule-prefix-input'), { target: { value: 'Old submitted prefix' } });
    fireEvent.click(screen.getByTestId('home-range-save')); await waitFor(() => expect(resolve).toBeTypeOf('function')); first.unmount();
    mount(); await editRange(); await act(async () => resolve({ ...nativeRule(), revision_id: 'R3', naming_prefix: 'Old submitted prefix', runtime_apply_status: 'not_running' }));
    expect(screen.getByTestId('rule-prefix-input')).toHaveValue('Rule R'); expect(screen.queryByTestId('home-range-saved')).not.toBeInTheDocument(); noOtherWrites();
  });
  it('revokes a pending pre-save read when the Home workspace changes', async () => {
    const { client } = mount(); const form = await editRange(); fireEvent.change(within(form).getByTestId('rule-prefix-input'), { target: { value: 'Old workspace' } });
    let resolve!: (value: Awaited<ReturnType<typeof studioV2DatalinkApi.get>>) => void;
    vi.mocked(studioV2DatalinkApi.get).mockImplementationOnce(() => new Promise(done => { resolve = done; }));
    fireEvent.click(screen.getByTestId('home-range-save')); await waitFor(() => expect(resolve).toBeTypeOf('function'));
    await act(async () => { client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { ...workspace, id: 'NEXT' }); });
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    const read = locatorRead(2); read.workspace = workspace;
    await act(async () => resolve({ data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } }));
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); expect(studioV2RulesAPI.update).not.toHaveBeenCalled(); noOtherWrites();
  });
  it('starts an existing-device connection buffer without changing or duplicating the saved device', async () => {
    mount(); await screen.findByTestId('home-device-points-D'); fireEvent.click(screen.getByRole('button', { name: 'Manage devices' }));
    fireEvent.click(within(screen.getByTestId('home-management-table')).getByRole('button', { name: 'Connection' }));
    const dialog = await screen.findByRole('dialog'); const name = await within(dialog).findByTestId('recording-source-name');
    expect(name).toHaveValue('Device D'); expect(within(dialog).getByLabelText('Protocol', { selector: 'select' })).toBeDisabled();
    fireEvent.change(name, { target: { value: 'Unsent device name' } }); fireEvent.click(screen.getByTestId('home-source-close'));
    fireEvent.click(await screen.findByTestId('recording-source-discard'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument()); noOtherWrites();
  });
  it('updates the existing connection exactly once and closes only after the original ACK and readback', async () => {
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockImplementation(async (_id, request) => {
      const saved = { ...completeDevice(), ...request, connection_config: JSON.stringify(request.connection_config), runtime_apply_status: 'not_running' as const };
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]);
      return saved;
    });
    mount(); await screen.findByTestId('home-device-points-D'); fireEvent.click(screen.getByRole('button', { name: 'Manage devices' }));
    fireEvent.click(within(screen.getByTestId('home-management-table')).getByRole('button', { name: 'Connection' }));
    const dialog = await screen.findByRole('dialog');
    fireEvent.change(await within(dialog).findByTestId('recording-source-name'), { target: { value: 'Saved device name' } });
    expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled();
    await waitFor(() => expect(within(dialog).getByTestId('recording-source-save')).toBeEnabled());
    fireEvent.click(within(dialog).getByTestId('recording-source-save'));
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledOnce());
    expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledWith('D', expect.objectContaining({ name: 'Saved device name' }));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled(); expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
  });
  it('opens the original new-range buffer and discards unsent values without a source write', async () => {
    mount(); const dialog = await openRange(); fireEvent.click(await within(dialog).findByTestId('home-range-done')); fireEvent.click(await within(dialog).findByTestId('home-range-new'));
    const name = await within(dialog).findByLabelText('Point-name prefix');
    fireEvent.change(name, { target: { value: 'Unsent new prefix' } });
    fireEvent.click(screen.getByTestId('home-source-close'));
    fireEvent.click(await screen.findByTestId('home-source-discard'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument()); noOtherWrites();
  });

  it('creates one new Home range through the original owners and refreshes the acknowledged list and point count', async () => {
    let saved = nativeRule('NEW');
    let mapping: Awaited<ReturnType<typeof studioV2MappingsAPI.create>>;
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready', readiness: false } });
    vi.mocked(studioV2RulesAPI.create).mockImplementation(async request => {
      saved = { ...nativeRule(request.id), ...request, revision_id: 'NEW1' };
      vi.mocked(studioV2RulesAPI.list).mockResolvedValue([nativeRule(), saved]);
      return { ...saved, runtime_apply_status: 'not_running' };
    });
    vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async () => ({ source_rule_id: saved.id, workspace_id: 'W', workspace_revision: 'W2', revision_id: saved.revision_id!,
      tags: { status: 'ready', candidates: [{ id: 'CANDIDATE', identity: { source_rule_id: saved.id, candidate_type: 'tag', candidate_kind: 'point', derived_from_rule_address: '40010' }, proposed_signature: 'SIG',
        address: '40010', point_id: 'NEWPOINT', tag_key: `${saved.naming_prefix}_40010`, display_name: `${saved.naming_prefix}_40010`, data_type: 'int16', status: 'draft' }] },
      database_outputs: { status: 'ready', candidates: [] }, local_modbus_outputs: { status: 'deferred', candidates: [] } }));
    vi.mocked(studioV2MappingsAPI.create).mockImplementation(async request => {
      mapping = { ...request, id: 'NEWMAP', workspace_id: 'W', device_id: 'D', point_id: 'NEWPOINT', tag_id: 'NEWTAG', created_at: '', updated_at: '', runtime_apply_status: 'not_running' };
      vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([mapping]); return mapping;
    });
    vi.mocked(pointAPI.get).mockImplementation(async () => ({ id: 'NEWPOINT', device_id: 'D', address: '40010', data_type: 'int16', enabled: false } as Awaited<ReturnType<typeof pointAPI.get>>));
    vi.mocked(tagAPI.get).mockImplementation(async () => ({ id: 'NEWTAG', key: mapping.tag_key, data_type: mapping.target_type, display_name: mapping.display_name, unit: mapping.unit ?? '', status: 'draft' } as Awaited<ReturnType<typeof tagAPI.get>>));
    mount(); const dialog = await openRange(); fireEvent.click(await within(dialog).findByTestId('home-range-done')); fireEvent.click(await within(dialog).findByTestId('home-range-new'));
    fireEvent.change(await within(dialog).findByLabelText('Point-name prefix'), { target: { value: 'New source' } });
    fireEvent.change(within(dialog).getByLabelText('Starting address', { selector: 'input' }), { target: { value: '40010' } });
    fireEvent.change(within(dialog).getByLabelText('Number of data items', { selector: 'input' }), { target: { value: '1' } });
    const save = within(dialog).getByRole('button', { name: 'Save range' });
    await waitFor(() => expect(save).toBeEnabled()); fireEvent.click(save);
    await within(dialog).findByTestId('recording-range-saved-inactive');
    expect(studioV2RulesAPI.create).toHaveBeenCalledOnce(); expect(studioV2MappingsAPI.create).toHaveBeenCalledOnce();
    fireEvent.click(within(dialog).getByTestId('recording-range-back'));
    const edit = await within(dialog).findByTestId(`home-range-edit-${saved.id}`);
    expect(edit.closest('li')).toHaveTextContent('40010 · 1');
    expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled(); expect(studioV2RulesAPI.remove).not.toHaveBeenCalled();
  });

  it('keeps Save disabled while the original connection prepare is unresolved', async () => {
    mount(); await screen.findByTestId('home-device-points-D');
    let waiting = false;
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockImplementationOnce(() => { waiting = true; return new Promise(() => undefined); });
    fireEvent.click(screen.getByRole('button', { name: 'Manage devices' }));
    fireEvent.click(within(screen.getByTestId('home-management-table')).getByRole('button', { name: 'Connection' }));
    const dialog = await screen.findByRole('dialog');
    fireEvent.change(await within(dialog).findByTestId('recording-source-name'), { target: { value: 'Prepared later' } });
    await waitFor(() => expect(waiting).toBe(true));
    expect(within(dialog).getByTestId('recording-source-save')).toBeDisabled(); noOtherWrites();
  });

});
