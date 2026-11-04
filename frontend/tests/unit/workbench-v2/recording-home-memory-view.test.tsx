import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
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
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ modbusShareAPI: { status: vi.fn() }, pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
const clients: QueryClient[] = [];
const modbusDevice = () => ({ ...nativeDevice(), connection_config: '{"host":"device.example","port":502,"slave_id":1,"timeout":1500}' });
const mcDevice = () => ({ ...nativeDevice(), protocol: 'mc_3e' as const, connection_config: '{"host":"fixture.invalid","port":5000,"network_no":0,"pc_no":255,"io_no":1023,"station_no":0}' });
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'], modbus_share: { workspace_id: 'W', workspace_revision: 'W2', hydration_state: 'ready' as const, readiness: true } };
const range = (id: string, start_address: string, count: number, data_type: SourceRuleRecord['data_type'], extra: Partial<SourceRuleRecord> = {}): SourceRuleRecord =>
  ({ ...nativeRule(id), revision_id: `${id}-rev`, start_address, count, data_type, naming_prefix: `${id}_`, ...extra });
let device = modbusDevice();
let rules: SourceRuleRecord[] = [];
function arrange(next: SourceRuleRecord[], nextDevice = modbusDevice()) {
  rules = next; device = nextDevice;
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device]);
  vi.mocked(studioV2RulesAPI.list).mockImplementation(async () => rules);
}
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear();
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', configured_enabled: false, enabled: false, port: 502, address: '127.0.0.1:502', bind_state: 'disabled', mapping_count: 0 });
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: 'D', status: 'ready', connect_status: 'success', probe_status: 'success', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]); vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(sourceRuleCandidateAPI.get).mockRejectedValue(new Error('database output unavailable'));
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error('Unexpected unscoped child read');
    const params = config?.params as Record<string, string>;
    const read = locatorRead(params.step === '2' ? 2 : 1); read.entry.devices = [device];
    const selected = rules.find(rule => rule.id === params.rule_id);
    if (params.step === '2' && selected) {
      read.locator.rule_id = selected.id; read.entry.source_rules = [selected]; read.entry.rule_heads = { [selected.id]: selected.revision_id! };
      read.entry.hydration_scope = { ...read.entry.hydration_scope, rule_ids: [selected.id] };
    }
    const { modbus_share: _share, ...locatorWorkspace } = workspace;
    return { data: { success: true, data: { ...locatorWorkspace, locator_entry: read.entry } } };
  });
  arrange([range('A', '40001', 2, 'int16'), range('B', '40010', 2, 'float32')]);
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '#home' }} navigateTo={vi.fn()} /></QueryClientProvider>);
}
async function openList() {
  mount(); fireEvent.click(await screen.findByTestId('home-device-points-D'));
  const dialog = await screen.findByRole('dialog');
  await within(dialog).findByTestId(`home-range-edit-${rules[0].id}`);
  return dialog;
}
async function openMemory() {
  const dialog = await openList();
  fireEvent.click(within(dialog).getByTestId('home-range-memory'));
  await within(dialog).findByTestId('home-memory-view');
  await within(dialog).findAllByRole('tab');
  return dialog;
}
function noWrites() {
  for (const fn of [studioV2WorkspaceDevicesAPI.create, studioV2WorkspaceDevicesAPI.update, studioV2RulesAPI.create, studioV2RulesAPI.update,
    studioV2RulesAPI.remove, studioV2MappingsAPI.create, studioV2MappingsAPI.update,
    studioV2WorkspaceWriteGroupsAPI.create, studioV2WorkspaceWriteGroupsAPI.update, studioV2WorkspaceWriteGroupsAPI.apply]) expect(fn).not.toHaveBeenCalled();
}
describe('Home range list shows saved facts and one top close', () => {
  it('lists several saved ranges with real address span, type, item count and register footprint', async () => {
    const dialog = await openList();
    const first = within(dialog).getByTestId('home-range-edit-A').closest('li')!;
    expect(first).toHaveTextContent('40001 – 40002'); expect(first).toHaveTextContent('int16'); expect(first).toHaveTextContent('2 items'); expect(first).toHaveTextContent('2 registers');
    const second = within(dialog).getByTestId('home-range-edit-B').closest('li')!;
    expect(second).toHaveTextContent('40010 – 40013'); expect(second).toHaveTextContent('float32'); expect(second).toHaveTextContent('4 registers');
    expect(within(dialog).getAllByRole('button', { name: 'Close' })).toHaveLength(1);
    expect(within(dialog).getByTestId('home-source-close')).toBeInTheDocument();
    expect(within(dialog).getByTestId('home-range-memory')).toHaveTextContent('Memory view');
    noWrites();
  });

  it('opens the only saved range directly and its Cancel returns to the range list with memory view', async () => {
    arrange([range('A', '40001', 2, 'int16')]);
    mount(); fireEvent.click(await screen.findByTestId('home-device-points-D'));
    const dialog = await screen.findByRole('dialog');
    expect(await within(dialog).findByTestId('rule-prefix-input')).toHaveValue('A_');
    fireEvent.click(within(dialog).getByTestId('home-range-done'));
    await within(dialog).findByTestId('home-range-edit-A');
    expect(within(dialog).getByTestId('home-range-memory')).toBeInTheDocument();
    noWrites();
  });
});

describe('memory view opens exact cells into the original range flows', () => {
  it('separates areas and shows a float32 value as one item across two registers', async () => {
    const dialog = await openMemory();
    expect(within(dialog).getByRole('tab', { name: /4x/ })).toHaveAttribute('aria-selected', 'true');
    expect(within(dialog).getByRole('tab', { name: /3x/ })).toBeInTheDocument();
    expect(within(dialog).getByTestId('home-memory-cell-40001')).toHaveAttribute('data-state', 'used');
    expect(within(dialog).getByTestId('home-memory-cell-40010')).toHaveAccessibleName('40010, B_, float32, register 1 of 2');
    expect(within(dialog).getByTestId('home-memory-cell-40011')).toHaveAccessibleName('40011, B_, float32, register 2 of 2, starts at 40010');
    expect(within(dialog).getByTestId('home-memory-cell-40001')).toHaveTextContent('int16'); expect(within(dialog).getByTestId('home-memory-cell-40010')).toHaveTextContent(/float32 1\/2.*next row/);
    expect(within(dialog).getByTestId('home-memory-cell-40011')).toHaveTextContent('cont. 40010'); expect(within(dialog).getByTestId('home-memory-cell-40012')).toHaveAttribute('data-join-next', 'true');
    fireEvent.click(within(dialog).getByRole('tab', { name: /3x/ }));
    expect(within(dialog).getByTestId('home-memory-cell-30001')).toHaveAttribute('data-state', 'empty');
    noWrites();
  });

  it('prefills the exact empty address into the original new-range form and waits for an explicit save', async () => {
    const dialog = await openMemory();
    fireEvent.click(within(dialog).getByTestId('home-memory-cell-40003'));
    const address = await within(dialog).findByLabelText('Starting address', { selector: 'input' });
    expect(address).toHaveValue('40003');
    expect(within(dialog).getByRole('button', { name: 'Save range' })).toBeDisabled();
    fireEvent.change(within(dialog).getByLabelText('Point-name prefix', { selector: 'input' }), { target: { value: 'UNSENT_' } });
    fireEvent.click(within(dialog).getByTestId('recording-range-back'));
    fireEvent.click(await within(dialog).findByTestId('home-source-keep'));
    expect(within(dialog).getByLabelText('Starting address', { selector: 'input' })).toHaveValue('40003');
    fireEvent.click(within(dialog).getByTestId('recording-range-back'));
    fireEvent.click(await within(dialog).findByTestId('home-source-discard'));
    await within(dialog).findByTestId('home-range-edit-A');
    expect(within(dialog).queryByTestId('home-memory-view')).not.toBeInTheDocument();
    noWrites();
  });

  it('returns from a used cell editor to the range list and from the memory view with its own back', async () => {
    const dialog = await openMemory();
    fireEvent.click(within(dialog).getByTestId('home-memory-back'));
    await within(dialog).findByTestId('home-range-edit-A');
    fireEvent.click(within(dialog).getByTestId('home-range-memory'));
    fireEvent.click(await within(dialog).findByTestId('home-memory-cell-40001'));
    await within(dialog).findByTestId('rule-prefix-input');
    fireEvent.click(within(dialog).getByTestId('home-range-done'));
    await within(dialog).findByTestId('home-range-edit-A');
    expect(within(dialog).queryByTestId('home-memory-view')).not.toBeInTheDocument();
    noWrites();
  });

  it('opens the owning saved range from a used cell', async () => {
    const dialog = await openMemory();
    fireEvent.click(within(dialog).getByTestId('home-memory-cell-40011'));
    expect(await within(dialog).findByTestId('rule-prefix-input')).toHaveValue('B_');
    noWrites();
  });

  it('asks which range to open when one address belongs to several ranges', async () => {
    arrange([range('A', '40001', 2, 'int16'), range('B', '40002', 1, 'int16')]);
    const dialog = await openMemory();
    const cell = within(dialog).getByTestId('home-memory-cell-40002');
    expect(cell).toHaveAttribute('data-state', 'conflict');
    fireEvent.click(cell);
    const owners = await within(dialog).findByTestId('home-memory-owners');
    expect(within(dialog).queryByTestId('rule-prefix-input')).not.toBeInTheDocument();
    fireEvent.click(within(owners).getByRole('button', { name: /B_/ }));
    expect(await within(dialog).findByTestId('rule-prefix-input')).toHaveValue('B_');
    noWrites();
  });

  it('marks a string area unconfirmed and offers no new-range entry there', async () => {
    arrange([range('A', '40001', 2, 'int16'), range('S', '40010', 2, 'string')]);
    const dialog = await openMemory();
    expect(within(dialog).getByRole('alert')).toHaveTextContent(/cannot be confirmed/i);
    const unknown = within(dialog).getByTestId('home-memory-cell-40003');
    expect(unknown).toHaveAttribute('data-state', 'unconfirmed'); expect(unknown.tagName).not.toBe('BUTTON');
    fireEvent.click(unknown);
    expect(within(dialog).queryByLabelText('Starting address', { selector: 'input' })).not.toBeInTheDocument();
    noWrites();
  });

  it('shows no cells while the range read is slow or failed', async () => {
    const dialog = await openList();
    let fail!: (error: Error) => void;
    vi.mocked(studioV2RulesAPI.list).mockImplementation(() => new Promise((_resolve, reject) => { fail = reject; }));
    fireEvent.click(within(dialog).getByTestId('home-range-memory'));
    const view = await within(dialog).findByTestId('home-memory-view');
    expect(within(view).getByRole('status')).toHaveTextContent('Reading the current workspace');
    expect(within(view).queryByTestId(/home-memory-cell-/)).not.toBeInTheDocument();
    await waitFor(() => expect(fail).toBeTypeOf('function')); fail(new Error('offline'));
    expect(await within(view).findByRole('alert')).toBeInTheDocument();
    expect(within(view).queryByTestId(/home-memory-cell-/)).not.toBeInTheDocument();
    noWrites();
  });

  it('uses current MC persisted proof for D int16 and float32 and prefills D33', async () => {
    arrange([range('I', 'D0', 32, 'int16'), range('F', 'D40', 1, 'float32')], mcDevice());
    vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(async id => {
      const rule = rules.find(item => item.id === id)!; const width = rule.data_type === 'float32' ? 2 : 1, start = Number(rule.start_address.slice(1));
      return { source_rule_id: id, workspace_id: 'W', workspace_revision: 'W2', revision_id: rule.revision_id!, tags: { status: 'ready', candidates: Array.from({ length: rule.count }, (_, index) => ({
        id: `${id}${index}`, identity: { source_rule_id: id, candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: rule.start_address }, proposed_signature: `${id}${index}`,
        address: `D${start + index * width}`, point_id: `${id}P${index}`, tag_key: `${id}${index}`, display_name: `${id}${index}`, data_type: rule.data_type, status: 'draft' as const })) } };
    });
    const dialog = await openMemory();
    await waitFor(() => expect(within(dialog).getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'used'));
    expect(within(dialog).getByTestId('home-memory-cell-D1')).toHaveAccessibleName(/D1.*I_.*int16.*1 of 1/);
    expect(within(dialog).getByTestId('home-memory-cell-D41')).toHaveAccessibleName(/D41.*F_.*float32.*2 of 2.*D40/);
    fireEvent.click(within(dialog).getByTestId('home-memory-cell-D33'));
    expect(await within(dialog).findByLabelText('Starting address', { selector: 'input' })).toHaveValue('D33');
    noWrites();
  });

  it('keeps MC cells unconfirmed when persisted proof cannot be read', async () => {
    arrange([range('I', 'D0', 32, 'int16'), range('F', 'D40', 1, 'float32')], mcDevice());
    vi.mocked(sourceRuleCandidateAPI.getTags).mockRejectedValue(new Error('evidence unavailable'));
    const dialog = await openMemory();
    await waitFor(() => expect(within(dialog).getByTestId('home-memory-cell-D10')).toHaveAttribute('data-state', 'unconfirmed'));
    expect(within(dialog).getByRole('alert')).toHaveTextContent(/saved MC points/);
    expect(within(dialog).getByTestId('home-memory-cell-D10').tagName).not.toBe('BUTTON');
    expect(within(dialog).getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'unconfirmed');
    noWrites();
  });

  it('has complete Traditional Chinese labels', async () => {
    await i18n.changeLanguage('zh-TW');
    mount(); fireEvent.click(await screen.findByTestId('home-device-points-D'));
    const dialog = await screen.findByRole('dialog');
    const entry = await within(dialog).findByTestId('home-range-memory');
    expect(entry).toHaveTextContent('記憶體視圖'); fireEvent.click(entry);
    await within(dialog).findByTestId('home-memory-view');
    expect(within(dialog).getByTestId('home-memory-back')).toHaveTextContent('返回範圍列表');
  });
});

describe('MC proof is bound to the range list workspace revision', () => {
  const proof = (revision: string) => async (id: string) => {
    const rule = rules.find(item => item.id === id)!, start = Number(rule.start_address.slice(1));
    return { source_rule_id: id, workspace_id: 'W', workspace_revision: revision, revision_id: rule.revision_id!, tags: { status: 'ready' as const, candidates: Array.from({ length: rule.count }, (_, index) => ({
      id: `${id}${index}`, identity: { source_rule_id: id, candidate_type: 'tags', candidate_kind: 'tag', derived_from_rule_address: rule.start_address }, proposed_signature: `${id}${index}`,
      address: `D${start + index}`, point_id: `${id}P${index}`, tag_key: `${id}${index}`, display_name: `${id}${index}`, data_type: rule.data_type, status: 'draft' as const })) } };
  };
  beforeEach(() => arrange([range('I', 'D0', 4, 'int16'), range('J', 'D10', 1, 'int16')], mcDevice()));

  it('rejects W1 candidate evidence for a W2 range list', async () => {
    vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(proof('W1'));
    const dialog = await openMemory();
    await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('I', { workspace_id: 'W', expected_workspace_revision: 'W2', revision_id: 'I-rev' }));
    await waitFor(() => expect(within(dialog).queryByText('Checking saved MC points…')).not.toBeInTheDocument());
    expect(within(dialog).getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'unconfirmed');
    expect(within(dialog).getByTestId('home-memory-cell-D5')).toHaveAttribute('data-state', 'unconfirmed');
    expect(within(dialog).getByRole('alert')).toHaveTextContent(/I_.*saved MC points/);
    fireEvent.click(within(dialog).getByRole('tab', { name: /^M/ }));
    expect(within(dialog).getByTestId('home-memory-cell-M0')).toHaveAttribute('data-state', 'empty');
    noWrites();
  });

  it('reads no MC proof when the Share scope belongs to another workspace', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, modbus_share: { ...workspace.modbus_share, workspace_id: 'OTHER' } });
    vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(proof('W2'));
    const dialog = await openMemory();
    expect(within(dialog).getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'unconfirmed');
    fireEvent.click(within(dialog).getByRole('tab', { name: /^M/ }));
    expect(within(dialog).getByTestId('home-memory-cell-M0')).toHaveAttribute('data-state', 'unconfirmed');
    expect(sourceRuleCandidateAPI.getTags).not.toHaveBeenCalledWith('I', expect.objectContaining({ expected_workspace_revision: 'W2' }));
  });

  it('re-reads proof for the new revision after the list is read again', async () => {
    vi.mocked(sourceRuleCandidateAPI.getTags).mockImplementation(proof('W2'));
    const dialog = await openMemory();
    await waitFor(() => expect(within(dialog).getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'used'));
    await waitFor(() => expect(within(dialog).queryByText('Checking saved MC points…')).not.toBeInTheDocument());
    fireEvent.click(within(dialog).getByRole('tab', { name: /^M/ }));
    expect(within(dialog).getByTestId('home-memory-cell-M0')).toHaveAttribute('data-state', 'empty');
    fireEvent.click(within(dialog).getByTestId('home-memory-back'));
    await within(dialog).findByTestId('home-range-edit-I');
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, modbus_share: { ...workspace.modbus_share, workspace_revision: 'W3' } });
    fireEvent.click(within(dialog).getByTestId('home-range-memory'));
    await within(dialog).findAllByRole('tab');
    await waitFor(() => expect(sourceRuleCandidateAPI.getTags).toHaveBeenCalledWith('I', expect.objectContaining({ expected_workspace_revision: 'W3' })));
    await waitFor(() => expect(within(dialog).queryByText('Checking saved MC points…')).not.toBeInTheDocument());
    expect(within(dialog).getByTestId('home-memory-cell-D0')).toHaveAttribute('data-state', 'unconfirmed');
  });
});

describe('memory view blocks every free cell while a recovery is unresolved', () => {
  it('shows the recovery notice and renders empty cells as unconfirmed', async () => {
    const { RecordingHomeMemoryView } = await import('@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView');
    const onAdd = vi.fn();
    render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false} blocked="recovery" onOpen={vi.fn()} onAdd={onAdd} onBack={vi.fn()}
      inputs={[{ id: 'A', label: 'A_', start_address: '40001', count: 1, data_type: 'int16', skipped_addresses: [] }]} />);
    expect(screen.getByRole('alert')).toHaveTextContent(/still unconfirmed/);
    expect(screen.getByTestId('home-memory-cell-40001')).toHaveAttribute('data-state', 'used');
    fireEvent.click(screen.getByTestId('home-memory-cell-40002'));
    expect(screen.getByTestId('home-memory-cell-40002')).toHaveAttribute('data-state', 'unconfirmed');
    expect(onAdd).not.toHaveBeenCalled();
  });
});
