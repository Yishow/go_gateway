import { StrictMode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI, type StudioV2MappingSaveResult } from '@/services/studioV2Mappings';
import { deviceAPI, pointAPI, tagAPI, modbusShareAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
import type { SourceRuleRecord } from '@/types/datalink';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';
import i18n from '@/i18n/config';

vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn(), getTags: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ modbusShareAPI: { status: vi.fn() }, pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn(), testConnection: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));

const clients: QueryClient[] = [];
const device = { ...nativeDevice(), name: 'MC', protocol: 'mc_3e' as const,
  connection_config: '{"host":"127.0.0.1","port":15110,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0,"timeout":5}' };
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready' as const, readiness: true } };
let connected = false;
let rule: SourceRuleRecord | undefined;
let mappings: StudioV2MappingSaveResult[] = [];

beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear();
  connected = false; rule = undefined; mappings = [];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false, port: 502, address: '127.0.0.1:502', bind_state: 'disabled', mapping_count: 0 });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device]);
  vi.mocked(studioV2RulesAPI.list).mockImplementation(async () => rule ? [rule] : []);
  vi.mocked(studioV2MappingsAPI.list).mockImplementation(async () => mappings);
  vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(deviceAPI.checkReadiness).mockImplementation(async () => ({ device_id: 'D', status: connected ? 'ready' : 'warning',
    connect_status: connected ? 'success' : 'unknown', probe_status: connected ? 'success' : 'unknown', activation_allowed: connected,
    planning_allowed: true, apply_allowed: connected, blocking_reasons: [], checks: [] }));
  vi.mocked(deviceAPI.testConnection).mockImplementation(async () => { connected = true; return { success: true, error: '', latency_ms: 1, can_activate: true }; });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error('Unexpected locator read');
    const params = config?.params as Record<string, string>;
    const read = locatorRead(params.step === '2' ? 2 : 1); read.entry.devices = [device];
    if (params.step === '2' && rule) {
      read.locator.rule_id = rule.id; read.entry.source_rules = [rule]; read.entry.rule_heads = { [rule.id]: rule.revision_id! };
      read.entry.hydration_scope = { ...read.entry.hydration_scope, rule_ids: [rule.id] };
    }
    const { modbus_share: _share, ...plain } = workspace;
    return { data: { success: true, data: { ...plain, locator_entry: read.entry } } };
  });
  vi.mocked(studioV2RulesAPI.create).mockImplementation(async request => {
    rule = { ...nativeRule(request.id), ...request, revision_id: 'R1' } as SourceRuleRecord;
    return { ...rule, runtime_apply_status: 'not_running' };
  });
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (id, request) => {
    rule = { ...rule!, ...request, id, revision_id: 'R2' }; return { ...rule, runtime_apply_status: 'not_running' };
  });
  vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async () => ({ source_rule_id: rule!.id, workspace_id: 'W', workspace_revision: 'W2', revision_id: rule!.revision_id!,
    tags: { status: 'ready', candidates: [{ id: 'C', identity: { source_rule_id: rule!.id, candidate_type: 'tag', candidate_kind: 'point', derived_from_rule_address: 'D0' },
      proposed_signature: 'S', address: 'D0', point_id: 'P', tag_key: 'mc_d0', display_name: 'MC D0', data_type: 'int16', status: 'draft' }] },
    database_outputs: { status: 'ready', candidates: [] }, local_modbus_outputs: { status: 'deferred', candidates: [] } }));
  vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation((id, scope) => sourceRuleCandidateAPI.get(id, scope));
  vi.mocked(studioV2MappingsAPI.create).mockImplementation(async request => {
    const saved = { ...request, id: 'M', workspace_id: 'W', device_id: 'D', point_id: 'P', tag_id: 'T', created_at: '', updated_at: '', runtime_apply_status: 'not_running' } as StudioV2MappingSaveResult;
    mappings = [saved]; return saved;
  });
  vi.mocked(studioV2MappingsAPI.update).mockImplementation(async (_id, request) => {
    const saved = { ...mappings[0], ...request, runtime_apply_status: 'not_running' as const }; mappings = [saved]; return saved;
  });
  vi.mocked(pointAPI.get).mockImplementation(async () => ({ id: 'P', device_id: 'D', address: 'D0', data_type: 'int16', enabled: rule?.enabled ?? false } as Awaited<ReturnType<typeof pointAPI.get>>));
  vi.mocked(tagAPI.get).mockImplementation(async () => ({ id: 'T', key: mappings[0].tag_key, data_type: mappings[0].target_type,
    display_name: mappings[0].display_name, unit: mappings[0].unit ?? '', status: 'draft' } as Awaited<ReturnType<typeof tagAPI.get>>));
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

async function openInactiveRange() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  render(<StrictMode><QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '' }} navigateTo={vi.fn()} /></QueryClientProvider></StrictMode>);
  fireEvent.click(await screen.findByTestId('home-device-points-D'));
  const dialog = await screen.findByRole('dialog');
  fireEvent.change(await within(dialog).findByLabelText('Starting address', { selector: 'input' }), { target: { value: 'D0' } });
  fireEvent.change(within(dialog).getByLabelText('Number of data items', { selector: 'input' }), { target: { value: '1' } });
  fireEvent.change(within(dialog).getByLabelText('Point-name prefix', { selector: 'input' }), { target: { value: 'MC' } });
  const save = within(dialog).getByRole('button', { name: 'Save range' });
  await waitFor(() => expect(save).toBeEnabled()); fireEvent.click(save);
  await within(dialog).findByTestId('recording-range-saved-inactive', undefined, { timeout: 5000 });
  expect(studioV2RulesAPI.create).toHaveBeenCalledOnce(); expect(studioV2MappingsAPI.create).toHaveBeenCalledOnce();
  expect(rule?.enabled).toBe(false); expect(deviceAPI.testConnection).not.toHaveBeenCalled();
  return dialog;
}

it('prepares a newly saved MC range through the real locator queue after an explicit successful connection test', async () => {
  const dialog = await openInactiveRange();
  fireEvent.click(within(dialog).getByTestId('recording-range-prepare'));
  await waitFor(() => expect(deviceAPI.testConnection).toHaveBeenCalledOnce());
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledWith(rule!.id, expect.objectContaining({ enabled: true })));
  await waitFor(() => expect(studioV2MappingsAPI.update).toHaveBeenCalledOnce());
  await waitFor(() => expect(rule?.enabled).toBe(true));
});

it('does not enable an MC range whose persisted revision changes during diagnostics', async () => {
  const dialog = await openInactiveRange();
  let finish!: (result: Awaited<ReturnType<typeof deviceAPI.testConnection>>) => void;
  vi.mocked(deviceAPI.testConnection).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
  fireEvent.click(within(dialog).getByTestId('recording-range-prepare'));
  await waitFor(() => expect(deviceAPI.testConnection).toHaveBeenCalledOnce());
  rule = { ...rule!, revision_id: 'FOREIGN', naming_prefix: 'FOREIGN' };
  connected = true;
  await act(async () => finish({ success: true, error: '', latency_ms: 1, can_activate: true }));
  await within(dialog).findByText(/The saved source changed/);
  expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
  expect(studioV2MappingsAPI.update).not.toHaveBeenCalled();
  expect(rule.enabled).toBe(false);
});
