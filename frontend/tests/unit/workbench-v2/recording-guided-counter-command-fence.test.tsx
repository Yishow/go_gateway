import { useState, type ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, renderHook, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitCounterConfirmation, type CircuitCounterConfirmationProps } from '@/features/datalink/workbench-v2/steps/step4/CircuitCounterConfirmation';
import { CircuitProfileConfigFields } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileConfigFields';
import { defaultCircuitProfileLocalDraft, parseStrictProfileIntervalInput } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { recordingTestWriteRecovery as store } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { recordingPreparationBlocked, type RecordingPreparationFenceInput } from '@/features/datalink/workbench-v2/steps/step4/recordingPreparationFence';
import { useCircuitRecordingCommandFence } from '@/features/datalink/workbench-v2/steps/step4/useCircuitRecordingCommandFence';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { MeasurementDefinition } from '@/types/measurement';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const candidates = buildGroupCandidates(savedGroupState()).candidates;
const group = savedGroup();
const clients: QueryClient[] = [];
function definition(overrides: Partial<MeasurementDefinition> = {}): MeasurementDefinition {
  return { id: 'm-1', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', name: 'Water total',
    equipment_id: 'circuit-water', quantity: 'volume', unit: 'm3', semantic_kind: 'counter', definition_revision: 'rev-9',
    source_binding_revision: 'binding-4', series_epoch: 'epoch-3', numeric_encoding: 'uint64',
    counter_policy: { allow_negative: false, rollover_modulus: '18446744073709551616', max_rate_per_second: '0.125', reset_threshold: 4 }, ...overrides };
}
function fenceInput(): RecordingPreparationFenceInput {
  return { ownerKey: 'ws-1/dev-1/circuit-water', workspaceId: 'ws-1', workspaceRevision: 'w-1', groups: [group], start: { state: 'observed', result: 'none' },
    schema: { [group.id]: { scope: { workspaceId: 'ws-1', workspaceRevision: 'w-1', groupId: group.id, groupRevision: group.revision,
      connectorId: group.destination.connector_id, connectorRevision: group.destination.connector_revision, epoch: 0 },
    observation: { owner: 'schema', state: 'observed', result: 'none' }, phase: 'idle', preview: 'none', commands: {} } } };
}
function wrapper() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  return ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}
function props(overrides: Partial<CircuitCounterConfirmationProps> = {}): CircuitCounterConfirmationProps {
  return { workspaceId: 'ws-1', workspaceRevision: 'w-1', circuitId: 'circuit-water', readonly: false,
    unconfirmedStatuses: [{ candidate: candidates[0], confirmed: false }], onConfirmed: vi.fn(), ...overrides };
}
function enterUnit() { fireEvent.change(screen.getByTestId('circuit-confirm-counter-unit'), { target: { value: ' m3 ' } }); }
const submit = async () => { await act(async () => { fireEvent.click(screen.getByTestId('circuit-confirm-counter-submit')); }); };
function resetStore() {
  localStorage.clear(); sessionStorage.clear();
  for (const ref of store.read().references) store.resolve({ operation_id: ref.operation_id, action: 'test_write', status: 'succeeded',
    write_outcome: 'written_verified', cleanup_status: 'cleaned', created_at: ref.created_at, updated_at: ref.created_at }, ref.workspace_id);
}
beforeEach(() => { vi.restoreAllMocks(); vi.clearAllMocks(); resetStore();
  vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { data: definition() } });
  vi.mocked(studioV2DatalinkApi.put).mockResolvedValue({ data: { data: { measurement: definition(), epoch_transitioned: false } } }); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); });

describe('counter definition commands respect the same current recording owner denial', () => {
  it('blocks only submit while retaining local semantics/unit edits when the current owner reports O1', async () => {
    render(<CircuitCounterConfirmation {...props({ newOperationBlocked: true })} />, { wrapper: wrapper() }); enterUnit();
    fireEvent.change(screen.getByTestId('circuit-confirm-counter-semantics'), { target: { value: 'signed_counter' } });
    expect(screen.getByTestId('circuit-confirm-counter-unit')).toHaveValue(' m3 ');
    expect(screen.getByTestId('circuit-confirm-counter-semantics')).toHaveValue('signed_counter');
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled(); await submit();
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
  it.each(['known O1', 'newly unreadable storage'] as const)('rechecks %s appearing after render before sending POST', async condition => {
    const input = fenceInput(); const isNewOperationBlocked = () => recordingPreparationBlocked(input);
    expect(isNewOperationBlocked()).toBe(false);
    render(<CircuitCounterConfirmation {...props({ isNewOperationBlocked })} />, { wrapper: wrapper() }); enterUnit();
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeEnabled();
    if (condition === 'known O1') store.retain({ workspace_id: 'ws-1', group_id: group.id, group_revision: group.revision,
      operation_id: 'O1', operation_kind: 'test_write', created_at: '2026-10-28T00:00:00Z' });
    else vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private storage path'); });
    await submit(); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
    expect(document.body).not.toHaveTextContent('private storage path');
  });
  it('treats a throwing fresh denial callback as blocked without exposing raw exceptions', async () => {
    render(<CircuitCounterConfirmation {...props({ isNewOperationBlocked: () => { throw new Error('password=secret'); } })} />, { wrapper: wrapper() });
    enterUnit(); await submit(); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(document.body).not.toHaveTextContent('password=secret');
  });
  it('preserves standalone creation payload and only invokes the existing mutation after explicit confirmation', async () => {
    const onConfirmed = vi.fn(); render(<CircuitCounterConfirmation {...props({ onConfirmed })} />, { wrapper: wrapper() });
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); enterUnit();
    fireEvent.change(screen.getByTestId('circuit-confirm-counter-semantics'), { target: { value: 'signed_counter' } }); await submit();
    await waitFor(() => expect(onConfirmed).toHaveBeenCalledOnce());
    expect(studioV2DatalinkApi.post).toHaveBeenCalledExactlyOnceWith('/studio-v2/workspace/measurements', {
      workspace_id: 'ws-1', name: candidates[0].label, equipment_id: 'circuit-water', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
      quantity: '', semantic_kind: 'signed_counter', unit: 'm3', counter_policy: { allow_negative: true } });
  });
  it('preserves the exact saved binding/epoch/policy PUT payload when the fresh delegate permits submission', async () => {
    const saved = definition(); const onConfirmed = vi.fn();
    render(<CircuitCounterConfirmation {...props({ onConfirmed, isNewOperationBlocked: () => false,
      unconfirmedStatuses: [{ candidate: candidates[0], confirmed: false, measurement: saved }] })} />, { wrapper: wrapper() });
    await submit(); await waitFor(() => expect(onConfirmed).toHaveBeenCalledOnce());
    expect(studioV2DatalinkApi.put).toHaveBeenCalledExactlyOnceWith('/studio-v2/workspace/measurements/m-1', saved);
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('does not revive a captured denial delegate after its full owner scope retires A→B→A', async () => {
    const input = { ownerKey: 'owner-A', scopeKey: 'scope-A', workspaceId: 'ws-1', deviceId: 'dev-1', circuitId: 'circuit-water', batchGroups: [],
      updateWriteGroup: vi.fn(), ensureCircuitProfile: vi.fn() };
    const hook = renderHook(current => useCircuitRecordingCommandFence(current), { initialProps: input });
    act(() => hook.result.current.onCommandFence({ ownerKey: 'owner-A', scopeKey: 'scope-A', blocked: false, isBlocked: () => false }));
    const retired = hook.result.current.isBlocked;
    render(<CircuitCounterConfirmation {...props({ isNewOperationBlocked: retired })} />, { wrapper: wrapper() }); enterUnit();
    hook.rerender({ ...input, scopeKey: 'scope-B' }); hook.rerender(input);
    act(() => hook.result.current.onCommandFence({ ownerKey: 'owner-A', scopeKey: 'scope-A', blocked: false, isBlocked: () => false }));
    await submit(); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('keeps readonly fieldset restrictions and does not turn the denial into authority over reload or name-focus NAV', () => {
    const onConfirmed = vi.fn(); const onFocus = vi.fn();
    const view = render(<CircuitCounterConfirmation {...props({ newOperationBlocked: true, onConfirmed,
      unconfirmedStatuses: [{ candidate: candidates[0], confirmed: false, measurement: definition({ source_binding_revision: '' }) }] })} />, { wrapper: wrapper() });
    fireEvent.click(screen.getByTestId('circuit-counter-reload-source')); expect(onConfirmed).toHaveBeenCalledOnce();
    view.rerender(<CircuitCounterConfirmation {...props({ newOperationBlocked: true, namePrerequisite: { descriptionId: 'name', onFocus } })} />);
    fireEvent.click(screen.getByTestId('circuit-confirm-counter-name-focus')); expect(onFocus).toHaveBeenCalledOnce();
    view.rerender(<CircuitCounterConfirmation {...props({ readonly: true })} />);
    expect(screen.getByTestId('circuit-confirm-counter-unit')).toBeDisabled(); expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('passes the same denial through real ConfigFields without disabling role/name/local policy edits', () => {
    function Host() {
      const [draft, setDraft] = useState(() => ({ ...defaultCircuitProfileLocalDraft(candidates, []), circuit_id: 'circuit-water' }));
      return <CircuitProfileConfigFields workspaceId="ws-1" workspaceRevision="w-1" readonly={false}
        newOperationBlocked isNewOperationBlocked={() => true} savedDevices={[{ id: 'dev-1', name: 'Device A' }]} activeDeviceId="dev-1" onSelectDevice={vi.fn()}
        deviceCandidates={candidates} localDraft={draft} updateLocalDraft={setDraft} instParsed={parseStrictProfileIntervalInput('instantaneous', '15')}
        cumParsed={parseStrictProfileIntervalInput('cumulative', '10')} isPostgres={false} storageLocked={false} effectiveStorageMode="plain"
        cumulativeCircuitMissing={false} circuitId="circuit-water" unconfirmedCounterStatuses={props().unconfirmedStatuses} onCounterConfirmed={vi.fn()}
        instantaneousReduction="mean" />;
    }
    render(<Host />, { wrapper: wrapper() }); enterUnit();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Edited locally' } });
    const policy = screen.getByTestId('circuit-profile-recreate-missing') as HTMLInputElement; const wasChecked = policy.checked;
    fireEvent.click(policy);
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Edited locally');
    expect(policy.checked).toBe(!wasChecked);
    fireEvent.click(screen.getByTestId('circuit-profile-sources-toggle'));
    expect(screen.getByTestId('circuit-profile-role-tag-0')).toBeEnabled();
    expect(screen.getByTestId('circuit-confirm-counter-submit')).toBeDisabled();
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
});
