import { StrictMode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeAll, afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
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
vi.mock('@/services/datalink', () => ({ modbusShareAPI: { status: vi.fn() }, pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));

// Full StrictMode create → candidates → mapping ACK for a new Home range, from both entries.
// All services are disposable mocks; no device, database or user session is touched.
const clients: QueryClient[] = [];
const device = () => ({ ...nativeDevice(), connection_config: '{"host":"device.example","port":502,"slave_id":1,"timeout":1500}' });
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready' as const, readiness: true } };
const existing: SourceRuleRecord = { ...nativeRule('I'), revision_id: 'I1', start_address: '40001', count: 16, data_type: 'int16', naming_prefix: 'TEMP_' };
const addresses = Array.from({ length: 8 }, (_, index) => String(40017 + index * 2));
let rules: SourceRuleRecord[] = [];
let created: SourceRuleRecord | undefined;
let mappings: StudioV2MappingSaveResult[] = [];
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear();
  rules = [existing]; created = undefined; mappings = [];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false, port: 502, address: '127.0.0.1:502', bind_state: 'disabled', mapping_count: 0 });
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: 'D', status: 'warning', connect_status: 'unknown', probe_status: 'unknown', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device()]);
  vi.mocked(studioV2RulesAPI.list).mockImplementation(async () => rules);
  vi.mocked(studioV2MappingsAPI.list).mockImplementation(async () => mappings); vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error('Unexpected unscoped child read');
    const params = config?.params as Record<string, string>;
    const read = locatorRead(params.step === '2' ? 2 : 1); read.entry.devices = [device()];
    const selected = rules.find(rule => rule.id === params.rule_id);
    if (params.step === '2' && selected) {
      read.locator.rule_id = selected.id; read.entry.source_rules = [selected]; read.entry.rule_heads = { [selected.id]: selected.revision_id! };
      read.entry.hydration_scope = { ...read.entry.hydration_scope, rule_ids: [selected.id] };
    }
    const { modbus_share: _share, ...plain } = workspace;
    return { data: { success: true, data: { ...plain, locator_entry: read.entry } } };
  });
  vi.mocked(studioV2RulesAPI.create).mockImplementation(async request => {
    created = { ...nativeRule(request.id), ...request, revision_id: 'NEW1' } as SourceRuleRecord;
    rules = [existing, created];
    return { ...created, runtime_apply_status: 'not_running' };
  });
  vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async id => {
    if (!created || id !== created.id) throw new Error('candidates unavailable');
    const rule = created;
    return { source_rule_id: rule.id, workspace_id: 'W', workspace_revision: 'W2', revision_id: rule.revision_id!,
      tags: { status: 'ready', candidates: addresses.map(address => ({ id: `C${address}`, identity: { source_rule_id: rule.id, candidate_type: 'tag', candidate_kind: 'point', derived_from_rule_address: '40017' },
        proposed_signature: `S${address}`, address, point_id: `P${address}`, tag_key: `${rule.naming_prefix}_${address}`, display_name: `${rule.naming_prefix}_${address}`, data_type: 'float32', status: 'draft' })) },
      database_outputs: { status: 'ready', candidates: [] }, local_modbus_outputs: { status: 'deferred', candidates: [] } };
  });
  vi.mocked(studioV2MappingsAPI.create).mockImplementation(async request => {
    const address = request.address;
    const saved = { ...request, id: `M${address}`, workspace_id: 'W', device_id: 'D', point_id: `P${address}`, tag_id: `T${address}`, created_at: '', updated_at: '', runtime_apply_status: 'not_running' } as StudioV2MappingSaveResult;
    mappings = [...mappings, saved]; return saved;
  });
  vi.mocked(pointAPI.get).mockImplementation(async id => ({ id, device_id: 'D', address: id.slice(1), data_type: 'float32', enabled: false } as Awaited<ReturnType<typeof pointAPI.get>>));
  vi.mocked(tagAPI.get).mockImplementation(async id => {
    const mapping = mappings.find(item => item.tag_id === id)!;
    return { id, key: mapping.tag_key, data_type: mapping.target_type, display_name: mapping.display_name, unit: mapping.unit ?? '', status: 'draft' } as Awaited<ReturnType<typeof tagAPI.get>>;
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

async function openNewRange(entry: 'direct' | 'memory') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  render(<StrictMode><QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '#home' }} navigateTo={vi.fn()} /></QueryClientProvider></StrictMode>);
  fireEvent.click(await screen.findByTestId('home-device-points-D'));
  const dialog = await screen.findByRole('dialog');
  await waitFor(() => expect(within(dialog).getByTestId('home-range-option-I')).toHaveAttribute('aria-pressed', 'true'));
  if (entry === 'direct') {
    fireEvent.click(within(dialog).getByTestId('home-range-option-new'));
    fireEvent.change(await within(dialog).findByLabelText('Starting address', { selector: 'input' }), { target: { value: '40017' } });
  } else {
    fireEvent.click(within(dialog).getByTestId('home-range-memory'));
    fireEvent.click(await within(dialog).findByTestId('home-memory-cell-40017'));
    expect(await within(dialog).findByLabelText('Starting address', { selector: 'input' })).toHaveValue('40017');
  }
  fireEvent.change(within(dialog).getByLabelText('Number of data items', { selector: 'input' }), { target: { value: '8' } });
  fireEvent.change(within(dialog).getByLabelText('Data type', { selector: 'select' }), { target: { value: 'float32' } });
  fireEvent.change(within(dialog).getByLabelText('Point-name prefix', { selector: 'input' }), { target: { value: 'FLOW_' } });
  return dialog;
}

describe('StrictMode new Home range reaches an acknowledged saved state', () => {
  for (const entry of ['direct', 'memory'] as const) {
    it(`saves 40017 × 8 float32 FLOW_ from the ${entry} entry with one create`, async () => {
      const dialog = await openNewRange(entry);
      const save = within(dialog).getByRole('button', { name: 'Save range' });
      await waitFor(() => expect(save).toBeEnabled());
      fireEvent.click(save); fireEvent.click(save);
      await within(dialog).findByTestId('recording-range-saved-inactive', undefined, { timeout: 5000 });
      expect(studioV2RulesAPI.create).toHaveBeenCalledOnce();
      expect(studioV2RulesAPI.create).toHaveBeenCalledWith(expect.objectContaining({ device_id: 'D', start_address: '40017', count: 8, data_type: 'float32', naming_prefix: 'FLOW_', enabled: false }));
      expect(studioV2MappingsAPI.create).toHaveBeenCalledTimes(addresses.length);
      expect(vi.mocked(studioV2MappingsAPI.create).mock.calls.map(([request]) => request.address).sort()).toEqual(addresses);
      expect(await within(dialog).findByTestId(`home-range-option-${created!.id}`)).toHaveTextContent('40017');
      expect(studioV2RulesAPI.update).not.toHaveBeenCalled(); expect(studioV2RulesAPI.remove).not.toHaveBeenCalled();
      expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled();
    });
  }
  it('blocks the save when readiness fails and never creates', async () => {
    vi.mocked(deviceAPI.checkReadiness).mockRejectedValue(new Error('503'));
    const dialog = await openNewRange('direct');
    expect(await within(dialog).findByText(i18n.t('recording-editor:device.scopeBlocked'))).toBeInTheDocument();
    expect(within(dialog).queryByText(i18n.t('recording-editor:device.checkingScope'))).not.toBeInTheDocument();
    expect(within(dialog).getByRole('button', { name: 'Save range' })).toBeDisabled();
    expect(studioV2RulesAPI.create).not.toHaveBeenCalled(); expect(studioV2MappingsAPI.create).not.toHaveBeenCalled();
  });
});

// Warm only the code: these cases still exercise the real wrapper, scoped reads, and draft owners.
// Cold import timing and failures remain covered by the dedicated advanced-loading tests.
beforeAll(async () => { await import('@/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page'); });
