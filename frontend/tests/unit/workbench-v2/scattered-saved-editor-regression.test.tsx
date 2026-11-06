import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import '../../../src/i18n/config';
import { RecordingHomeSavedRangeEditor } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor';
import type { HomeRuleObserver } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeSourceHost';
import { hydrateStudioV2Device } from '@/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { hydrateStudioV2Rule } from '@/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2RuleAutosave } from '@/pages/datalink/workbench-v2/useStudioV2RuleAutosave';
import { readStudioLocatorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import type { SourceRuleRecord } from '@/types/datalink';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';

vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEntry: vi.fn() }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
beforeEach(() => vi.resetAllMocks());

function mountSavedEditor(record: SourceRuleRecord, protocol: 'fatek_fbs' | 'mc_3e' = 'fatek_fbs') {
  const device = { ...nativeDevice(), protocol,
    connection_config: JSON.stringify({ host: 'fixture.invalid', port: 5000, acquisition: { version: 1,
      profile: protocol === 'mc_3e' ? 'mitsubishi-q-iqr-3e-binary-v1' : 'fatek-fbs-ascii-v1' } }) };
  const mcEvidence = protocol === 'mc_3e' ? { status: 'compatible' as const,
    scope: { workspace_id: 'W', workspace_revision: 'W2', revision_id: 'R2' },
    points: (record.mixed_items ?? []).map((item, index) => ({ point_id: `P${index}`, address: item.address ?? '' })) } : undefined;
  const read = locatorRead(2);
  read.entry.source_rules = [record];
  read.entry.devices = [device];
  vi.mocked(readStudioLocatorEntry).mockResolvedValue({ ok: true, value: read });
  const initial = { ...INITIAL_STATE, devices: [hydrateStudioV2Device(device)], rules: [hydrateStudioV2Rule(record)] };
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  function Host() {
    const [state, dispatch] = React.useReducer(workbenchV2Reducer, initial);
    const stateRef = React.useRef(state);
    const observer = React.useRef<HomeRuleObserver | undefined>(undefined);
    const leave = React.useRef<((proceed: () => void) => void) | undefined>(undefined);
    const autosave = useStudioV2RuleAutosave({ dispatch }, stateRef, undefined, {
      submitted: (rule, generation) => observer.current?.submitted?.(rule, generation),
      acknowledged: () => true,
      saved: (rule, saved) => observer.current?.saved(rule, saved),
      failed: (rule, error) => observer.current?.failed(rule, error),
    }, 'W');
    const send = (action: WorkbenchV2Action) => {
      const previous = stateRef.current;
      stateRef.current = workbenchV2Reducer(previous, action);
      dispatch(action);
      autosave.afterRuleAction(action);
    };
    return <RecordingHomeSavedRangeEditor target={{ workspaceId: 'W', deviceId: 'D', ruleId: 'R', kind: 'points', epoch: 1 }}
      state={state} rule={state.rules[0]} mcEvidence={mcEvidence} dispatch={send} leave={leave} suspended={false} list={null} onLeave={() => undefined}
      onClose={() => undefined} headerSlot={null} onLocked={() => undefined} ruleObserverRef={observer}
      onRetryShare={() => undefined} onOpenShareSettings={() => undefined} />;
  }
  const view = render(<QueryClientProvider client={client}><Host /></QueryClientProvider>);
  return { ...view, client };
}

it('saves an edit to a 70-item FATEK range and accepts its 70-item acknowledgment', async () => {
  const items = Array.from({ length: 70 }, (_, index) => ({ item_id: `item-${index}`, name: `Item ${index}`,
    address: `R${index * 2}`, register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' as const }));
  const record: SourceRuleRecord = { ...nativeRule(), start_address: 'R0', count: 70, layout_mode: 'mixed', mixed_items: items };
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({ ...record, ...request, count: 70,
    revision_id: 'R3', runtime_apply_status: 'not_running' }));
  const view = mountSavedEditor(record);
  expect(screen.getByTestId('scattered-usage-over-limit')).toBeVisible();
  fireEvent.change(screen.getByTestId('scattered-item-name-0'), { target: { value: 'Motor' } });
  fireEvent.click(screen.getByTestId('home-range-save'));
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledOnce());
  expect(vi.mocked(studioV2RulesAPI.update).mock.calls[0][1].mixed_items).toHaveLength(70);
  await screen.findByTestId('home-range-saved');
  expect(screen.queryByTestId('home-range-unknown')).not.toBeInTheDocument();
  expect(screen.getByTestId('scattered-item-name-0')).toHaveValue('Motor');
  view.unmount(); view.client.clear();
});

it('accepts a canonical item name acknowledgment when the submitted name has trailing space', async () => {
  const record: SourceRuleRecord = { ...nativeRule(), start_address: 'R0', count: 1, layout_mode: 'mixed', mixed_items: [
    { item_id: 'pump', name: 'Old', address: 'R0', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
  ] };
  vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({ ...record, ...request, count: 1,
    revision_id: 'R3', runtime_apply_status: 'not_running' }));
  const view = mountSavedEditor(record);
  fireEvent.change(screen.getByTestId('scattered-item-name-0'), { target: { value: 'Pump ' } });
  fireEvent.click(screen.getByTestId('home-range-save'));
  await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledOnce());
  expect(vi.mocked(studioV2RulesAPI.update).mock.calls[0][1].mixed_items?.[0].name).toBe('Pump');
  await screen.findByTestId('home-range-saved');
  expect(screen.getByTestId('scattered-item-name-0')).toHaveValue('Pump');
  expect(screen.queryByTestId('home-range-unknown')).not.toBeInTheDocument();
  view.unmount(); view.client.clear();
});

it.each(['changed-item', 'foreign-workspace', 'foreign-device', 'old-revision'] as const)(
  'keeps a mismatched mixed acknowledgment locked (%s)', async mismatch => {
    const record: SourceRuleRecord = { ...nativeRule(), start_address: 'R0', count: 1, layout_mode: 'mixed', mixed_items: [
      { item_id: 'a', name: 'Old', address: 'R0', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
    ] };
    vi.mocked(studioV2RulesAPI.update).mockImplementation(async (_id, request) => ({ ...record, ...request, count: 1,
      revision_id: mismatch === 'old-revision' ? 'R2' : 'R3',
      workspace_id: mismatch === 'foreign-workspace' ? 'OTHER' : 'W', device_id: mismatch === 'foreign-device' ? 'OTHER' : 'D',
      mixed_items: mismatch === 'changed-item' ? [{ ...request.mixed_items![0], name: 'Unexpected' }] : request.mixed_items,
      runtime_apply_status: 'not_running' }));
    const view = mountSavedEditor(record);
    fireEvent.change(screen.getByTestId('scattered-item-name-0'), { target: { value: 'Pump ' } });
    fireEvent.click(screen.getByTestId('home-range-save'));
    expect(await screen.findByTestId('home-range-unknown')).toBeVisible();
    expect(screen.getByTestId('scattered-item-name-0')).toBeDisabled();
    view.unmount(); view.client.clear();
  },
);

it('blocks saving an overlapping MC mixed draft and keeps it editable', () => {
  const record: SourceRuleRecord = { ...nativeRule(), start_address: 'D3', count: 2, data_type: 'uint32', layout_mode: 'mixed', mixed_items: [
    { item_id: 'wide', name: 'Wide', address: 'D3', register_offset: 0, data_type: 'uint32', semantic_kind: 'gauge' },
    { item_id: 'other', name: 'Other', address: 'D5', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
  ] };
  const view = mountSavedEditor(record, 'mc_3e');
  fireEvent.change(screen.getByTestId('scattered-item-address-1'), { target: { value: 'D4' } });
  expect(screen.getByTestId('scattered-item-error-1')).toBeVisible();
  expect(screen.getByTestId('home-range-save')).toBeDisabled();
  fireEvent.click(screen.getByTestId('home-range-save'));
  expect(readStudioLocatorEntry).not.toHaveBeenCalled();
  expect(studioV2RulesAPI.update).not.toHaveBeenCalled();
  expect(screen.getByTestId('scattered-item-address-1')).toBeEnabled();
  fireEvent.change(screen.getByTestId('scattered-item-address-1'), { target: { value: 'D6' } });
  expect(screen.getByTestId('home-range-save')).toBeEnabled();
  view.unmount(); view.client.clear();
});

it('keeps a known pre-write validation rejection editable with a safe error', async () => {
  const record: SourceRuleRecord = { ...nativeRule(), start_address: 'R0', count: 1, layout_mode: 'mixed', mixed_items: [
    { item_id: 'a', name: 'Old', address: 'R0', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
  ] };
  vi.mocked(studioV2RulesAPI.update)
    .mockRejectedValueOnce({ response: { status: 400, data: { error: { code: 'validation', message: 'private raw error' } } } })
    .mockImplementation(async (_id, request) => ({ ...record, ...request, count: 1, revision_id: 'R3', runtime_apply_status: 'not_running' }));
  const view = mountSavedEditor(record);
  fireEvent.change(screen.getByTestId('scattered-item-name-0'), { target: { value: 'Draft' } });
  fireEvent.click(screen.getByTestId('home-range-save'));
  expect(await screen.findByTestId('home-range-validation-error')).not.toHaveTextContent('private raw error');
  expect(screen.getByTestId('scattered-item-name-0')).toBeEnabled();
  expect(screen.getByTestId('home-range-save')).toBeEnabled();
  expect(screen.queryByTestId('home-range-unknown')).not.toBeInTheDocument();
  fireEvent.change(screen.getByTestId('scattered-item-name-0'), { target: { value: 'Corrected' } });
  fireEvent.click(screen.getByTestId('home-range-save'));
  await screen.findByTestId('home-range-saved');
  expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2);
  view.unmount(); view.client.clear();
});

it('keeps an ambiguous failed save locked', async () => {
  const record: SourceRuleRecord = { ...nativeRule(), start_address: 'R0', count: 1, layout_mode: 'mixed', mixed_items: [
    { item_id: 'a', name: 'Old', address: 'R0', register_offset: 0, data_type: 'int16', semantic_kind: 'gauge' },
  ] };
  vi.mocked(studioV2RulesAPI.update).mockRejectedValue(new Error('connection lost'));
  const view = mountSavedEditor(record);
  fireEvent.change(screen.getByTestId('scattered-item-name-0'), { target: { value: 'Draft' } });
  fireEvent.click(screen.getByTestId('home-range-save'));
  expect(await screen.findByTestId('home-range-unknown')).toBeVisible();
  expect(screen.getByTestId('scattered-item-name-0')).toBeDisabled();
  expect(screen.getByTestId('home-range-save')).toBeDisabled();
  view.unmount(); view.client.clear();
});
