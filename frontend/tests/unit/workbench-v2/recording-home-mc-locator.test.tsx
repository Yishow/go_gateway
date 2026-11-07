import { configReadContext, configReadKey } from '@/hooks/datalink/recordingConfigRefresh';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { homeSources } from './recording-home-source-fixtures';
import { homeDeviceAction } from './recording-home-device-actions';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeAll, afterEach, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { deviceAPI, pointAPI, modbusShareAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { modbusShareStatusKey } from '@/hooks/datalink/useModbusShareStatus';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';

vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn(), getTags: vi.fn() } }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { defaults: { baseURL: '/api/v1/datalink' }, get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ modbusShareAPI: { status: vi.fn() }, pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));

afterEach(() => { cleanup(); vi.resetAllMocks(); });

async function openMC() {
  const device = { ...nativeDevice(), protocol: 'mc_3e' as const, connection_config: '{"host":"fixture.invalid","port":5000,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0}' };
  const rule = { ...nativeRule(), start_address: 'D0', count: 32, naming_prefix: 'TAG_' };
  const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready' as const, readiness: true } };
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false, port: 502, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0 });
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: 'D', status: 'ready', connect_status: 'success', probe_status: 'success', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device]);
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue(homeSources([device]));
  vi.mocked(studioV2RulesAPI.list).mockResolvedValue([rule]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]); vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error('Unexpected read');
    const params = config?.params as Record<string, string>;
    const read = locatorRead(params.step === '2' ? 2 : 1); read.workspace = workspace; read.entry.devices = [device];
    if (params.step === '2') read.entry.source_rules = [rule];
    const { modbus_share: _share, ...locatorWorkspace } = workspace;
    return { data: { success: true, data: { ...locatorWorkspace, locator_entry: read.entry } } };
  });
  vi.mocked(sourceRuleCandidateAPI.get).mockRejectedValue(new Error('database output unavailable'));
  vi.mocked(sourceRuleCandidateAPI.getTags).mockResolvedValue({ source_rule_id: 'R', workspace_id: 'W', workspace_revision: 'W2', revision_id: 'R2',
    tags: { status: 'ready', candidates: Array.from({ length: 32 }, (_, index) => ({ id: `C${index}`,
      identity: { source_rule_id: 'R', candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: 'D0' }, proposed_signature: `sig${index}`,
      address: `D${index}`, point_id: `P${index}`, tag_key: `tag${index}`, display_name: `TAG_${index}`, data_type: 'int16', status: 'draft' })) } });
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({ ...rule, ...request, revision_id: 'R3', runtime_apply_status: 'not_running' }));
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '#home' }} navigateTo={vi.fn()} /></QueryClientProvider>);
  fireEvent.click(await homeDeviceAction());
  const dialog = await screen.findByRole('dialog');
  await within(dialog).findByTestId('home-saved-range-form');
  await waitFor(() => expect(within(dialog).queryByTestId('rule-mc-persisted-status')).not.toBeInTheDocument());
  expect(within(dialog).getByTestId('range-summary')).toHaveTextContent('D31');
  return { client, dialog, rule, device, workspace };
}

it('opens the saved MC range with matching point evidence and saves an edited prefix', async () => {
  const { dialog } = await openMC();
  fireEvent.change(within(dialog).getByTestId('rule-prefix-input'), { target: { value: 'MC_NEW_' } });
  fireEvent.click(within(dialog).getByTestId('home-range-save'));
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledExactlyOnceWith('R', expect.objectContaining({ naming_prefix: 'MC_NEW_' })));
  expect(sourceRuleCandidateAPI.get).not.toHaveBeenCalled();
});

it('does not confirm a saved MC range when workspace evidence changes and then returns during the write', async () => {
  const { client, dialog, rule, workspace } = await openMC();
  let complete!: (value: Awaited<ReturnType<typeof studioV2RulesAPI.update>>) => void;
  vi.mocked(studioV2RulesAPI.update).mockImplementation(() => new Promise(resolve => { complete = resolve; }));
  fireEvent.change(within(dialog).getByTestId('rule-prefix-input'), { target: { value: 'PENDING_' } });
  fireEvent.click(within(dialog).getByTestId('home-range-save'));
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledOnce());
  await act(async () => {
    client.setQueryData(configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext()), { ...workspace, modbus_share: { ...workspace.modbus_share, workspace_revision: 'W3' } });
    client.setQueryData(modbusShareStatusKey, { workspace_id: 'W', workspace_revision: 'W3', configured_enabled: false, enabled: false,
      port: 502, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0 });
  });
  await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('R', expect.objectContaining({ expected_workspace_revision: 'W3' })));
  await act(async () => {
    client.setQueryData(configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext()), workspace);
    client.setQueryData(modbusShareStatusKey, { workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false,
      port: 502, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0 });
  });
  await waitFor(() => expect(within(dialog).queryByTestId('rule-mc-persisted-status')).not.toBeInTheDocument());
  await act(async () => complete({ ...rule, naming_prefix: 'PENDING_', revision_id: 'R3', runtime_apply_status: 'not_running' }));
  await within(dialog).findByTestId('home-range-unknown');
  expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('PENDING_');
  expect(studioV2RulesAPI.update).toHaveBeenCalledOnce();
});

it('does not queue a range edit if MC evidence changes during the pre-save read', async () => {
  const { client, dialog, rule, device, workspace } = await openMC();
  let completeRead!: (value: Awaited<ReturnType<typeof studioV2DatalinkApi.get>>) => void;
  vi.mocked(studioV2DatalinkApi.get).mockImplementationOnce(() => new Promise(resolve => { completeRead = resolve; }));
  fireEvent.change(within(dialog).getByTestId('rule-prefix-input'), { target: { value: 'KEEP_DRAFT_' } });
  fireEvent.click(within(dialog).getByTestId('home-range-save'));
  await waitFor(() => expect(completeRead).toBeTypeOf('function'));
  await act(async () => {
    client.setQueryData(configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext()), { ...workspace, modbus_share: { ...workspace.modbus_share, workspace_revision: 'W3' } });
    client.setQueryData(modbusShareStatusKey, { workspace_id: 'W', workspace_revision: 'W3', configured_enabled: false, enabled: false,
      port: 502, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0 });
  });
  await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('R', expect.objectContaining({ expected_workspace_revision: 'W3' })));
  const read = locatorRead(2); read.entry.devices = [device]; read.entry.source_rules = [rule];
  const { modbus_share: _share, ...locatorWorkspace } = workspace;
  await act(async () => completeRead({ data: { success: true, data: { ...locatorWorkspace, locator_entry: read.entry } } }));
  await within(dialog).findByText('source.stale');
  await act(async () => {
    client.setQueryData(configReadKey(studioV2WorkspaceKeys.bootstrap(), configReadContext()), workspace);
    client.setQueryData(modbusShareStatusKey, { workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false,
      port: 502, address: '127.0.0.1', bind_state: 'disabled', mapping_count: 0 });
  });
  await waitFor(() => expect(within(dialog).getByTestId('home-range-save')).toBeEnabled());
  expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('KEEP_DRAFT_');
  expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
});

// Warm only the code: these cases still exercise the real wrapper, scoped reads, and draft owners.
// Cold import timing and failures remain covered by the dedicated advanced-loading tests.
beforeAll(async () => { await import('@/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page'); });
