import { StrictMode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeAll, beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { deviceAPI, pointAPI } from '@/services/datalink';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import type { SourceRuleRecord } from '@/types/datalink';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';
import i18n from '@/i18n/config';
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { list: vi.fn(), get: vi.fn() }, tagAPI: { get: vi.fn() }, deviceAPI: { checkReadiness: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));

const device = (id: string) => ({ ...nativeDevice(id), connection_config: '{"host":"device.example","port":502,"slave_id":1,"timeout":1500}' });
const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D', 'B'] };
let rules: SourceRuleRecord[] = [];
let gates: Record<string, Promise<void>> = {};
const requested: string[] = [];
const clients: QueryClient[] = [];
/** One scoped locator read per (step, device, rule); a test may hold any one of them open. */
function locatorFor(params: Record<string, string>) {
  const step = params.step === '2' ? 2 : 1, deviceId = params.device_id, ruleId = params.rule_id;
  const read = locatorRead(step, deviceId); read.workspace = workspace; read.entry.devices = [device(deviceId)];
  if (step === 2) {
    const rule = rules.find(item => item.id === ruleId && item.device_id === deviceId)!;
    read.locator = { ...read.locator, rule_id: ruleId };
    read.entry = { ...read.entry, source_rules: [rule], rule_heads: { [ruleId]: rule.revision_id! }, hydration_scope: { ...read.entry.hydration_scope, rule_ids: [ruleId] } };
    read.selection = { ...read.selection, target: { ...(read.selection as { target: object }).target, ruleId, ruleRevision: rule.revision_id } } as typeof read.selection;
  }
  return { ...read.workspace, locator_entry: read.entry };
}
beforeEach(async () => {
  vi.resetAllMocks(); await i18n.changeLanguage('en'); sessionStorage.clear(); requested.length = 0; gates = {};
  rules = [nativeRule('R', 'D'), nativeRule('S', 'B')];
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace);
  vi.mocked(deviceAPI.checkReadiness).mockResolvedValue({ device_id: 'D', status: 'ready', connect_status: 'success', probe_status: 'success', activation_allowed: false, planning_allowed: true, apply_allowed: false, blocking_reasons: [], checks: [] });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'W2', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([device('D'), device('B')]);
  vi.mocked(studioV2RulesAPI.list).mockImplementation(async () => rules);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]); vi.mocked(pointAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error('Unexpected unscoped child read');
    const params = config?.params as Record<string, string>, key = `${params.step}:${params.device_id}${params.rule_id ? `:${params.rule_id}` : ''}`;
    requested.push(key); await gates[key];
    return { data: { success: true, data: locatorFor(params) } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
function mount(strict = false) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  const page = <QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings', hash: '#home' }} navigateTo={vi.fn()} /></QueryClientProvider>;
  return render(strict ? <StrictMode>{page}</StrictMode> : page);
}
async function open(id = 'D') {
  fireEvent.click(await screen.findByTestId(`home-device-points-${id}`));
  return screen.findByRole('dialog');
}
const header = (dialog: HTMLElement) => dialog.querySelector<HTMLElement>('.gw-modal-header')!;
const enabledSwitch = (dialog: HTMLElement) => within(header(dialog)).getByRole('switch', { name: 'Enable rule' });
const devicesOf = (dialog: HTMLElement) => within(dialog).getByRole('group', { name: 'Devices' });
const prefix = (dialog: HTMLElement) => within(dialog).findByTestId('rule-prefix-input');
function hold(key: string) { let release!: () => void; gates[key] = new Promise<void>(resolve => { release = resolve; }); return () => act(async () => release()); }
function noWrites() {
  for (const fn of [studioV2WorkspaceDevicesAPI.create, studioV2WorkspaceDevicesAPI.update, studioV2RulesAPI.create, studioV2RulesAPI.update,
    studioV2RulesAPI.remove, studioV2MappingsAPI.create, studioV2MappingsAPI.update, studioV2WorkspaceWriteGroupsAPI.create,
    studioV2WorkspaceWriteGroupsAPI.update, studioV2WorkspaceWriteGroupsAPI.apply]) expect(fn).not.toHaveBeenCalled();
}

describe('one settings frame switches devices', () => {
  it.each([false, true])('switches to another workspace device through its own scoped locator (strict=%s)', async strict => {
    mount(strict); const dialog = await open();
    expect(await prefix(dialog)).toHaveValue('Rule R');
    const devices = devicesOf(dialog);
    expect(within(devices).getAllByRole('button').map(button => button.dataset.testid)).toEqual(['home-source-device-D', 'home-source-device-B']);
    expect(within(devices).getByTestId('home-source-device-D')).toHaveAttribute('aria-pressed', 'true');
    fireEvent.click(within(devices).getByTestId('home-source-device-B'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule S'));
    expect(requested).toEqual(expect.arrayContaining(['1:B', '2:B:S']));
    expect(within(header(dialog)).getByRole('heading')).toHaveTextContent('Device B');
    expect(within(devicesOf(dialog)).getByTestId('home-source-device-B')).toHaveAttribute('aria-pressed', 'true');
    expect(screen.getAllByRole('dialog')).toHaveLength(1); noWrites();
  });
  it('names the device in the title from the first loading frame', async () => {
    const release = hold('1:D'); mount(); const dialog = await open();
    await waitFor(() => expect(requested).toContain('1:D'));
    expect(within(header(dialog)).getByRole('heading')).toHaveTextContent('Device D');
    await release(); expect(await prefix(dialog)).toHaveValue('Rule R');
  });
  it('keeps or discards unsent edits before another device is read', async () => {
    mount(); const dialog = await open();
    fireEvent.change(await prefix(dialog), { target: { value: 'UNSENT_' } });
    fireEvent.click(within(devicesOf(dialog)).getByTestId('home-source-device-B'));
    fireEvent.click(await within(dialog).findByTestId('home-source-keep'));
    expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('UNSENT_'); expect(requested).not.toContain('1:B');
    fireEvent.click(within(devicesOf(dialog)).getByTestId('home-source-device-B'));
    fireEvent.click(await within(dialog).findByTestId('home-source-discard'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule S')); noWrites();
  });
  it('cannot leave a pending or unconfirmed save by switching device, range or memory view', async () => {
    rules = [nativeRule('R', 'D'), nativeRule('R4', 'D'), nativeRule('S', 'B')];
    let reject!: (error: Error) => void; vi.mocked(studioV2RulesAPI.update).mockImplementation(() => new Promise((_resolve, fail) => { reject = fail; }));
    mount(); const dialog = await open();
    fireEvent.change(await prefix(dialog), { target: { value: 'Pending' } });
    fireEvent.click(within(dialog).getByTestId('home-range-save'));
    await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledOnce());
    const toggle = enabledSwitch(dialog);
    expect(toggle).toBeDisabled(); fireEvent.click(toggle); expect(toggle).toHaveAttribute('aria-checked', 'true');
    fireEvent.click(within(devicesOf(dialog)).getByTestId('home-source-device-B'));
    fireEvent.click(within(dialog).getByTestId('home-range-option-R4'));
    fireEvent.click(within(header(dialog)).getByTestId('home-range-memory'));
    expect(within(dialog).queryByTestId('home-source-discard')).not.toBeInTheDocument();
    await act(async () => reject(new Error('acknowledgement unavailable'))); await within(dialog).findByTestId('home-range-unknown');
    fireEvent.click(within(devicesOf(dialog)).getByTestId('home-source-device-B'));
    fireEvent.click(within(dialog).getByTestId('home-range-option-R4'));
    expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Pending');
    expect(requested.filter(key => key.startsWith('1:B') || key.endsWith(':R4'))).toEqual([]);
    expect(studioV2RulesAPI.update).toHaveBeenCalledOnce();
  });
  it('ignores the late reads of a device the operator already left', async () => {
    const release = hold('1:D'); mount(); const dialog = await open();
    await waitFor(() => expect(requested).toContain('1:D'));
    fireEvent.click(within(devicesOf(dialog)).getByTestId('home-source-device-B'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule S'));
    await release();
    expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule S');
    expect(within(header(dialog)).getByRole('heading')).toHaveTextContent('Device B');
    expect(requested).not.toContain('2:D:R'); noWrites();
  });
});

describe('the device row uses the last validated Home projection', () => {
  it('offers every workspace device even when opened during a background source refresh', async () => {
    mount(); const trigger = await screen.findByTestId('home-device-points-D');
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockImplementation(() => new Promise(() => undefined));
    fireEvent.click(screen.getByRole('button', { name: 'Refresh' }));
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.list).toHaveBeenCalledTimes(3));
    fireEvent.click(trigger);
    const dialog = await screen.findByRole('dialog'); await prefix(dialog);
    expect(within(devicesOf(dialog)).getAllByRole('button').map(button => button.dataset.testid)).toEqual(['home-source-device-D', 'home-source-device-B']);
    // Selecting another device still goes through its own scoped locator read.
    fireEvent.click(within(devicesOf(dialog)).getByTestId('home-source-device-B'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule S'));
    expect(requested).toEqual(expect.arrayContaining(['1:B', '2:B:S'])); noWrites();
  });
});

describe('ranges of one device are switched inline', () => {
  it('opens the first range directly and switches ranges without a separate list', async () => {
    rules = [nativeRule('R', 'D'), { ...nativeRule('R4', 'D'), start_address: '40011' }];
    mount(); const dialog = await open();
    expect(await prefix(dialog)).toHaveValue('Rule R');
    expect(within(dialog).queryByRole('region', { name: 'Read ranges' })).not.toBeInTheDocument();
    const options = within(dialog).getByRole('group', { name: 'Read ranges' });
    expect(within(options).getByTestId('home-range-option-R')).toHaveAttribute('aria-pressed', 'true');
    fireEvent.change(within(dialog).getByTestId('rule-prefix-input'), { target: { value: 'UNSENT_' } });
    fireEvent.click(within(options).getByTestId('home-range-option-R4'));
    fireEvent.click(await within(dialog).findByTestId('home-source-discard'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule R4'));
    expect(within(dialog).getByTestId('home-range-option-R4')).toHaveAttribute('aria-pressed', 'true');
    fireEvent.click(within(dialog).getByTestId('home-range-option-new'));
    expect(await within(dialog).findByLabelText('Point-name prefix', { selector: 'input' })).toHaveValue('');
    noWrites();
  });
  it('fails closed when the scoped range list cannot be read and recovers through its re-read', async () => {
    rules = [nativeRule('R', 'D'), nativeRule('R4', 'D')];
    let failOnce = true;
    vi.mocked(studioV2RulesAPI.list).mockImplementation(async () => {
      if (failOnce && requested.includes('2:D:R')) { failOnce = false; throw new Error('fixture list unavailable'); }
      return rules;
    });
    mount(); const dialog = await open();
    expect(await prefix(dialog)).toHaveValue('Rule R');
    const options = within(dialog).getByRole('group', { name: 'Read ranges' });
    expect(await within(dialog).findByText('The read ranges of this device could not be read. No changes have been sent.')).toBeInTheDocument();
    expect(within(options).getByTestId('home-range-option-new')).toBeDisabled();
    expect(within(options).queryByTestId('home-range-option-R4')).not.toBeInTheDocument();
    fireEvent.click(within(options).getByTestId('home-range-option-new'));
    expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule R');
    fireEvent.click(within(dialog).getByRole('button', { name: 'Read ranges again' }));
    await waitFor(() => expect(within(dialog).getByTestId('home-range-option-new')).toBeEnabled());
    expect(within(dialog).getByTestId('home-range-option-R4')).toBeEnabled(); noWrites();
  });
  it('lets a device without ranges add one and closes from its cancel action', async () => {
    rules = [nativeRule('R', 'D')];
    mount(); const dialog = await open();
    await prefix(dialog);
    fireEvent.click(within(devicesOf(dialog)).getByTestId('home-source-device-B'));
    expect(await within(dialog).findByLabelText('Point-name prefix', { selector: 'input' })).toHaveValue('');
    expect(within(dialog).getByTestId('home-range-option-new')).toHaveAttribute('aria-pressed', 'true');
    expect(within(dialog).getByTestId('recording-range-back')).toHaveTextContent('Close');
    fireEvent.click(within(dialog).getByTestId('recording-range-back'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument()); noWrites();
  });
  it('closes from Cancel and keeps the saved range open after its acknowledged save', async () => {
    vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => { const saved = { ...nativeRule(), ...request, revision_id: 'R3', runtime_apply_status: 'not_running' as const }; rules = [saved, nativeRule('S', 'B')]; return saved; });
    mount(); const dialog = await open();
    fireEvent.change(await prefix(dialog), { target: { value: 'Saved range' } });
    expect(within(dialog).getByTestId('home-range-save')).toHaveTextContent(/^Save$/);
    fireEvent.click(within(dialog).getByTestId('home-range-save'));
    await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledOnce());
    await within(dialog).findByTestId('home-range-saved');
    expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Saved range');
    expect(within(dialog).queryByTestId('home-range-edit-R')).not.toBeInTheDocument();
    expect(within(dialog).getByTestId('home-range-done')).toHaveTextContent('Close');
    fireEvent.click(within(dialog).getByTestId('home-range-done'));
    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    expect(studioV2RulesAPI.update).toHaveBeenCalledOnce();
  });
});

describe('the memory view returns to the setting it was opened from', () => {
  it('returns to a new range (not the first saved one) after the unsent values were explicitly discarded', async () => {
    mount(); const dialog = await open(); await prefix(dialog);
    fireEvent.click(within(dialog).getByTestId('home-range-option-new'));
    const name = await within(dialog).findByLabelText('Point-name prefix', { selector: 'input' });
    fireEvent.change(name, { target: { value: 'UNSENT_' } });
    fireEvent.click(within(header(dialog)).getByTestId('home-range-memory'));
    fireEvent.click(await within(dialog).findByTestId('home-source-keep'));
    expect(within(dialog).getByLabelText('Point-name prefix', { selector: 'input' })).toHaveValue('UNSENT_');
    fireEvent.click(within(header(dialog)).getByTestId('home-range-memory'));
    fireEvent.click(await within(dialog).findByTestId('home-source-discard'));
    fireEvent.click(await within(dialog).findByTestId('home-memory-back'));
    expect(await within(dialog).findByLabelText('Point-name prefix', { selector: 'input' })).toHaveValue('');
    expect(within(dialog).getByTestId('home-range-option-new')).toHaveAttribute('aria-pressed', 'true');
    expect(within(dialog).queryByTestId('rule-prefix-input')).not.toBeInTheDocument(); noWrites();
  });
});

describe('title bar controls follow the one draft', () => {
  it('moves the enabled switch into the title bar without a second copy and without a write', async () => {
    mount(); const dialog = await open(); await prefix(dialog);
    expect(within(dialog).getAllByRole('switch', { name: 'Enable rule' })).toHaveLength(1);
    const toggle = enabledSwitch(dialog);
    expect(within(dialog).getByTestId('home-range-save')).toBeDisabled();
    fireEvent.click(toggle);
    expect(toggle).toHaveAttribute('aria-checked', 'false'); expect(within(dialog).getByTestId('home-range-save')).toBeEnabled();
    fireEvent.click(within(dialog).getByTestId('home-range-done'));
    expect(await within(dialog).findByTestId('home-source-keep')).toBeInTheDocument();
    expect(within(dialog).queryByRole('switch', { name: 'Enable rule' })).not.toBeInTheDocument();
    fireEvent.click(within(dialog).getByTestId('home-source-keep'));
    expect(enabledSwitch(dialog)).toHaveAttribute('aria-checked', 'false'); noWrites();
  });
  it('opens the memory view from the current setting and returns to it', async () => {
    rules = [nativeRule('R', 'D'), { ...nativeRule('R4', 'D'), start_address: '40011' }];
    mount(); const dialog = await open(); await prefix(dialog);
    fireEvent.click(within(dialog).getByTestId('home-range-option-R4'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule R4'));
    fireEvent.change(within(dialog).getByTestId('rule-prefix-input'), { target: { value: 'UNSENT_' } });
    fireEvent.click(within(header(dialog)).getByTestId('home-range-memory'));
    fireEvent.click(await within(dialog).findByTestId('home-source-discard'));
    expect(await within(dialog).findByTestId('home-memory-view')).toBeInTheDocument();
    expect(within(dialog).queryByRole('switch', { name: 'Enable rule' })).not.toBeInTheDocument();
    expect(within(dialog).getByTestId('home-memory-back')).toHaveTextContent('Back to settings');
    fireEvent.click(within(dialog).getByTestId('home-memory-back'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule R4'));
    fireEvent.click(within(header(dialog)).getByTestId('home-range-memory'));
    fireEvent.click(await within(dialog).findByTestId('home-memory-cell-40001'));
    await waitFor(() => expect(within(dialog).getByTestId('rule-prefix-input')).toHaveValue('Rule R')); noWrites();
  });
});

// Cold chunk behavior is covered by advanced-home-source-loading; these cases exercise source owners.
beforeAll(async () => { await import('@/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page'); });
