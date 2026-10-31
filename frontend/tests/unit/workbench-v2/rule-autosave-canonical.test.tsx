import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceDevicesAPI, studioV2RulesAPI } from './device-autosave-page.testHarness';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, expect, it, vi } from 'vitest';
import { queryClient } from '@/lib/queryClient';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { RuleTabRail } from '@/features/datalink/workbench-v2/steps/step2/RuleTabRail';
import { RuleEditor } from '@/features/datalink/workbench-v2/steps/step2/RuleEditor';
import type { SourceRuleRecord } from '@/types/datalink';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const record: SourceRuleRecord = {
  id: 'RULE', device_id: 'DEVICE', workspace_id: 'W', revision_id: 'R1', start_address: '40001', count: 8,
  data_type: 'int16', naming_prefix: 'SOURCE', enabled: true, skipped_addresses: ['40003'],
  scale_multiplier: 1, scale_offset: 0, data_format: '', share_enabled: false,
  share_start_register: null, share_stride: null, locked: false, origin: 'manual', created_at: '', updated_at: '',
};
let serverRules: SourceRuleRecord[];
function Harness({ scope = 'W' }: { scope?: string }) {
  const owner = useStudioV2AutosaveState(true, undefined, false, undefined, scope);
  if (!owner.workspaceHydrated) return <p>Loading</p>;
  const rule = owner.state.rules[0];
  return <><RuleTabRail rules={owner.state.rules} devices={owner.state.devices} selectedRuleId={rule?.id ?? null} dispatch={owner.actions.dispatch} />
    {rule && <><RuleEditor rule={rule} devices={owner.state.devices} globalShareEnabled={false} dispatch={owner.actions.dispatch} />
      <output data-testid="rule-state">{JSON.stringify(rule)}</output></>}</>;
}
function state() { return JSON.parse(screen.getByTestId('rule-state').textContent!) as SourceRuleRecord & { save_state: string; name: string }; }
async function mount() {
  const client = new QueryClient({ defaultOptions: queryClient.getDefaultOptions() });
  const rendered = render(<QueryClientProvider client={client}><Harness /></QueryClientProvider>);
  await screen.findByTestId('rule-add-btn');
  return { setScope: (scope: string) => rendered.rerender(<QueryClientProvider client={client}><Harness scope={scope} /></QueryClientProvider>) };
}
beforeEach(() => {
  resetDeviceAutosaveMocks(); serverRules = [];
  studioV2WorkspaceDevicesAPI.list.mockResolvedValue([deviceFixture({ id: 'DEVICE', last_test_success: null })]);
  studioV2RulesAPI.list.mockImplementation(async () => serverRules);
  studioV2RulesAPI.create.mockImplementation(async request => {
    const saved = { ...record, ...request, enabled: false };
    serverRules = [saved]; return { ...saved, runtime_apply_status: 'not_running' };
  });
  studioV2RulesAPI.update.mockImplementation(async (id, request) => {
    const saved = { ...serverRules[0], ...request, id, revision_id: 'R2' };
    serverRules = [saved]; return { ...saved, runtime_apply_status: 'not_running' };
  });
});
it('legacy add immediately displays the acknowledged inactive rule, without reload or another write', async () => {
  await mount(); fireEvent.click(screen.getByTestId('rule-add-btn'));
  await waitFor(() => expect(state().save_state).toBe('saved'));
  const id = state().id;
  expect(studioV2RulesAPI.create).toHaveBeenCalledWith(expect.objectContaining({ enabled: true, count: 8 }));
  expect(screen.getByTestId(`rule-enabled-count-${id}`)).toHaveTextContent('0/8');
  expect(screen.getByTestId(`rule-toggle-enabled-${id}`)).toHaveAttribute('title', '啟用規則');
  expect(screen.getByRole('switch', { name: '啟用此規則 (停用時將不衍生其所屬點位)' })).toHaveAttribute('aria-checked', 'false');
  expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1); expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
});
it('adopts canonical fields without resetting skipped addresses or the local display name', async () => {
  serverRules = [{ ...record }]; await mount();
  fireEvent.change(screen.getByTestId('rule-name-input'), { target: { value: 'My local label' } });
  studioV2RulesAPI.update.mockImplementationOnce(async () => ({ ...record, revision_id: 'R2', enabled: false, start_address: '40009', naming_prefix: 'SENSOR',
    count: 6, data_type: 'uint32', skipped_addresses: ['40011'], scale_multiplier: 2, scale_offset: 3, data_format: 'CDAB', share_enabled: true, share_start_register: 40021, share_stride: 2 }));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: ' sensor ' } });
  await waitFor(() => expect(state().save_state).toBe('saved'));
  expect(state()).toMatchObject({ enabled: false, start_address: '40009', naming_prefix: 'SENSOR', count: 6, data_type: 'uint32', skipped_addresses: ['40011'],
    scale_multiplier: 2, scale_offset: 3, data_format: 'CDAB', share_enabled: true, share_start_register: 40021, share_stride: 2, name: 'My local label' });
});
it.each([400, 503])('one legacy conflict action causes one PUT under production defaults (HTTP %s)', async status => {
  serverRules = [{ ...record }]; await mount();
  studioV2RulesAPI.update.mockRejectedValue({ response: { status, data: { error: { code: status === 400 ? 'validation' : 'unknown_failure' } } } });
  fireEvent.change(screen.getByTestId('rule-start-input'), { target: { value: '40009' } });
  await waitFor(() => expect(state().save_state).toBe('save-error'), { timeout: 2500 });
  expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
});
it('preserves a pending newer field while untouched fields take canonical normalization into the next request', async () => {
  serverRules = [{ ...record }]; await mount();
  let resolve!: (value: SourceRuleRecord) => void;
  studioV2RulesAPI.update.mockImplementationOnce(() => new Promise<SourceRuleRecord>(yes => { resolve = yes; }));
  studioV2RulesAPI.update.mockImplementationOnce(() => new Promise(() => undefined));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: ' sensor ' } });
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
  fireEvent.change(screen.getByTestId('rule-start-input'), { target: { value: '40017' } });
  await act(async () => { resolve({ ...record, naming_prefix: 'SENSOR', enabled: false, revision_id: 'R2' }); });
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
  expect(state()).toMatchObject({ start_address: '40017', naming_prefix: 'SENSOR', enabled: false, skipped_addresses: [], save_state: 'saving' });
  expect(studioV2RulesAPI.update.mock.calls[1][1]).toMatchObject({ start_address: '40017', naming_prefix: 'SENSOR', enabled: false, skipped_addresses: [] });
});
it('a newer off/on intent returning to the submitted value is not overwritten by an old response', async () => {
  serverRules = [{ ...record }]; await mount();
  let resolve!: (value: SourceRuleRecord) => void;
  studioV2RulesAPI.update.mockImplementationOnce(() => new Promise<SourceRuleRecord>(yes => { resolve = yes; }));
  studioV2RulesAPI.update.mockImplementationOnce(() => new Promise(() => undefined));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: ' sensor ' } });
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
  fireEvent.click(screen.getByTestId('rule-toggle-enabled-RULE')); fireEvent.click(screen.getByTestId('rule-toggle-enabled-RULE'));
  await act(async () => { resolve({ ...record, naming_prefix: 'SENSOR', enabled: false, revision_id: 'R2' }); });
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
  expect(studioV2RulesAPI.update.mock.calls[1][1]).toMatchObject({ enabled: true, naming_prefix: 'SENSOR' });
});

it.each(['id', 'device_id', 'workspace_id'] as const)('rejects a foreign %s acknowledgement without adopting canonical values', async field => {
  serverRules = [{ ...record }]; await mount();
  studioV2RulesAPI.update.mockResolvedValue({ ...record, [field]: 'FOREIGN', enabled: false, revision_id: 'FOREIGN' });
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'LOCAL' } });
  await waitFor(() => expect(state().save_state).toBe('save-error'));
  expect(state()).toMatchObject({ id: 'RULE', device_id: 'DEVICE', workspace_id: 'W', revision_id: 'R1', naming_prefix: 'LOCAL', enabled: true });
  expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
});
it('a late response from a previous workspace epoch cannot normalize the current buffer or replay pending intent', async () => {
  serverRules = [{ ...record }]; const owner = await mount();
  let resolve!: (value: SourceRuleRecord) => void;
  studioV2RulesAPI.update.mockImplementationOnce(() => new Promise<SourceRuleRecord>(yes => { resolve = yes; }));
  fireEvent.change(screen.getByTestId('rule-prefix-input'), { target: { value: 'LOCAL' } });
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
  fireEvent.change(screen.getByTestId('rule-start-input'), { target: { value: '40017' } });
  owner.setScope('OTHER'); owner.setScope('W');
  await act(async () => { resolve({ ...record, naming_prefix: 'CANONICAL-OLD', enabled: false, revision_id: 'R2' }); });
  expect(state()).toMatchObject({ naming_prefix: 'LOCAL', start_address: '40017', enabled: true, revision_id: 'R1' });
  expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
});
