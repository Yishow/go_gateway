import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { homeSources } from './recording-home-source-fixtures';
import { homeDeviceAction } from './recording-home-device-actions';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeAll, afterEach, beforeEach, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { deviceAPI, modbusShareAPI, pointAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
import type { SourceRuleRecord } from '@/types/datalink';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';
import i18n from '@/i18n/config';
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn(), getTags: vi.fn() } }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { defaults: { baseURL: '/api/v1/datalink' }, get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ modbusShareAPI: { status: vi.fn() }, pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));

// A saved MC range stays open after its ACK. The next save needs current proof for the acknowledged revision.
const clients: QueryClient[] = [];
const device = { ...nativeDevice(), protocol: 'mc_3e' as const, connection_config: '{"host":"fixture.invalid","port":5000,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0}' };
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready' as const, readiness: true } };
let rule: SourceRuleRecord;
let releaseR2!: () => void;
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear();
  rule = { ...nativeRule('I'), revision_id: 'R1', start_address: 'D0', count: 1, data_type: 'int16', naming_prefix: 'I_' };
  const r2Proof = new Promise<void>(resolve => { releaseR2 = resolve; });
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false, port: 502, address: '127.0.0.1:502', bind_state: 'disabled', mapping_count: 0 });
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: 'D', status: 'ready', connect_status: 'success', probe_status: 'success', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device]);
  vi.mocked(recordingHomeSourcesAPI.read).mockResolvedValue(homeSources([device]));
  vi.mocked(studioV2RulesAPI.list).mockImplementation(async () => [rule]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]); vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(sourceRuleCandidateAPI.get).mockRejectedValue(new Error('database output unavailable'));
  vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(async (id, scope) => {
    if (scope.revision_id === 'R2') await r2Proof;
    return { source_rule_id: id, workspace_id: scope.workspace_id, workspace_revision: scope.expected_workspace_revision, revision_id: scope.revision_id, tags: { status: 'ready', candidates: [{
      id: 'C0', identity: { source_rule_id: id, candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: 'D0' }, proposed_signature: 'S0',
      address: 'D0', point_id: 'P0', tag_key: 'T0', display_name: 'T0', data_type: 'int16', status: 'draft' }] } };
  });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error('Unexpected unscoped child read');
    const params = config?.params as Record<string, string>;
    const read = locatorRead(params.step === '2' ? 2 : 1); read.entry.devices = [device];
    if (params.step === '2') {
      read.locator.rule_id = rule.id; read.entry.source_rules = [rule]; read.entry.rule_heads = { [rule.id]: rule.revision_id! };
      read.entry.hydration_scope = { ...read.entry.hydration_scope, rule_ids: [rule.id] };
    }
    const { modbus_share: _share, ...locatorWorkspace } = workspace;
    return { data: { success: true, data: { ...locatorWorkspace, locator_entry: read.entry } } };
  });
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => {
    rule = { ...rule, ...request, revision_id: rule.revision_id === 'R1' ? 'R2' : 'R3' };
    return { ...rule, runtime_apply_status: 'not_running' as const };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

it('saves an MC range twice only with proof scoped to the acknowledged revision', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '#home' }} navigateTo={vi.fn()} /></QueryClientProvider>);
  fireEvent.click(await homeDeviceAction());
  const dialog = await screen.findByRole('dialog');
  const save = () => within(dialog).getByTestId('home-range-save');
  fireEvent.change(await within(dialog).findByTestId('rule-prefix-input'), { target: { value: 'FIRST_' } });
  await waitFor(() => expect(save()).toBeEnabled());
  fireEvent.click(save());
  await within(dialog).findByTestId('home-range-saved');
  expect(studioV2RulesAPI.update).toHaveBeenCalledOnce();
  expect(studioV2RulesAPI.update).toHaveBeenLastCalledWith('I', expect.objectContaining({ naming_prefix: 'FIRST_' }));

  // R1 proof no longer authorizes the R2 range: the next edit waits for R2 proof of the same workspace, device and rule.
  await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('I', { workspace_id: 'W', expected_workspace_revision: 'W2', revision_id: 'R2' }));
  fireEvent.change(await within(dialog).findByTestId('rule-prefix-input'), { target: { value: 'SECOND_' } });
  expect(save()).toBeDisabled(); fireEvent.click(save());
  expect(studioV2RulesAPI.update).toHaveBeenCalledOnce();
  await act(async () => releaseR2());
  await waitFor(() => expect(save()).toBeEnabled());
  fireEvent.click(save());
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
  expect(studioV2RulesAPI.update).toHaveBeenLastCalledWith('I', expect.objectContaining({ naming_prefix: 'SECOND_' }));
  expect(rule.revision_id).toBe('R3');
  await within(dialog).findByTestId('home-range-saved');
});

// Warm only the code: these cases still exercise the real wrapper, scoped reads, and draft owners.
// Cold import timing and failures remain covered by the dedicated advanced-loading tests.
beforeAll(async () => { await import('@/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page'); });
