import { act, renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useRef, type ReactNode } from 'react';
import { beforeEach, expect, it, vi } from 'vitest';
import { useRecordingRangeSubmission } from '@/pages/datalink/workbench-v2/useRecordingRangeSubmission';
import { useStudioV2MappingAutosave } from '@/pages/datalink/workbench-v2/useStudioV2MappingAutosave';
import { workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { toStudioV2RuleCreateRequest } from '@/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import type { Mapping } from '@/features/datalink/workbench-v2/state/types';
import type { RangeSourceRequest } from '@/pages/datalink/workbench-v2/rangeSourceSubmission';
import type { SourceRuleRecord } from '@/types/datalink';
import type { StudioV2MappingSaveResult, StudioV2WorkspaceMappingRequest } from '@/services/studioV2Mappings';
import { savedGroupState } from '../../fixtures/writeGroupState';

const api = vi.hoisted(() => ({ workspace: vi.fn(), rules: vi.fn(), devices: vi.fn(), candidates: vi.fn(), list: vi.fn(), create: vi.fn(), update: vi.fn(), point: vi.fn(), tag: vi.fn() }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: api.workspace } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: api.rules } }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: api.devices } }));
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: api.candidates } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: api.list, create: api.create, update: api.update, remove: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { get: api.point }, tagAPI: { get: api.tag } }));
const device = { id: 'dev-1', name: 'PLC', description: '', protocol: 'modbus_tcp' as const, status: 'draft' as const,
  connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' };
const request: RangeSourceRequest = { workspaceId: 'W', sessionId: 'S', requestId: 'REQ', generation: 1, deviceId: 'dev-1', mappingDefaults: { unit: 'kWh' }, buffer: {
  name: 'Meter', naming_prefix: 'Meter', start_address: '40001', count: 1, data_type: 'int16', enabled: true, scale_multiplier: 1, scale_offset: 0,
  data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null,
} };
let serverRule: SourceRuleRecord;
let serverMapping: StudioV2MappingSaveResult | undefined;
function saveMapping(submitted: StudioV2WorkspaceMappingRequest): StudioV2MappingSaveResult {
  serverMapping = { ...submitted, id: 'MAP', workspace_id: 'W', device_id: device.id, point_id: 'P', tag_id: 'T', runtime_apply_status: 'not_running', created_at: '', updated_at: '' };
  return structuredClone(serverMapping);
}
beforeEach(() => {
  vi.resetAllMocks(); serverMapping = undefined;
  api.workspace.mockResolvedValue({ id: 'W', ordered_device_ids: [device.id], modbus_share: { workspace_revision: 'SOURCE' } });
  api.rules.mockImplementation(async () => [serverRule]);
  api.devices.mockResolvedValue([device]);
  api.candidates.mockImplementation(async () => ({ source_rule_id: serverRule.id, workspace_id: 'W', workspace_revision: 'SOURCE', revision_id: 'R1',
    tags: { status: 'ready', candidates: [{ id: 'C', point_id: 'P', address: '40001', identity: { source_rule_id: serverRule.id }, status: 'draft' }] } }));
  api.list.mockImplementation(async () => serverMapping ? [structuredClone(serverMapping)] : []);
  api.create.mockImplementation(async (submitted: StudioV2WorkspaceMappingRequest) => saveMapping(submitted));
  api.update.mockImplementation(async (_id: string, submitted: StudioV2WorkspaceMappingRequest) => saveMapping(submitted));
  api.point.mockResolvedValue({ id: 'P', device_id: device.id, address: '40001', data_type: 'int16', enabled: true });
  api.tag.mockImplementation(async () => ({ id: 'T', key: serverMapping?.tag_key, display_name: serverMapping?.display_name,
    unit: serverMapping?.unit, data_type: 'int16', status: 'active' }));
});
async function confirmedRange() {
  const queueRule = vi.fn<(action: WorkbenchV2Action) => void>();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const actions = { dispatch: vi.fn() };
  const hook = renderHook(() => {
    const stateRef = useRef(savedGroupState({ rules: [], points: [], mappings: {} }));
    const queue = useRef<ReturnType<typeof useStudioV2MappingAutosave>['afterMappingAction']>(() => undefined);
    const submission = useRecordingRangeSubmission({ workspaceId: 'W', enabled: true, stateRef,
      apply: action => { stateRef.current = workbenchV2Reducer(stateRef.current, action); },
      prepareDevice: async () => ({ device, activationAllowed: true }), diagnoseDevice: async () => false, queueRule,
      queueMapping: (...args) => queue.current(...args) });
    // Match the production caller: its callback changes when range state changes.
    const mapping = useStudioV2MappingAutosave(actions, stateRef, false, id => submission.canWrite(id), submission.mappingObservation);
    queue.current = mapping.afterMappingAction;
    return { submission, mapping, stateRef };
  }, { wrapper: ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider> });
  act(() => { void hook.result.current.submission.port.submit(request); });
  await waitFor(() => expect(queueRule).toHaveBeenCalledOnce());
  const action = queueRule.mock.calls[0][0];
  if (action.type !== 'addRule') throw new Error('Expected existing addRule owner');
  const rule = action.rule;
  serverRule = { ...toStudioV2RuleCreateRequest(rule), id: rule.id, workspace_id: 'W', revision_id: 'R1', locked: false,
    origin: 'manual', skipped_addresses: [], created_at: '', updated_at: '' };
  const acknowledged = { ...serverRule, runtime_apply_status: 'not_running' as const };
  act(() => {
    expect(hook.result.current.submission.ruleObservation.acknowledged(rule, acknowledged)).toBe(true);
    hook.result.current.submission.ruleObservation.saved(rule, acknowledged);
  });
  await waitFor(() => expect(hook.result.current.submission.port.observations()[0].status).toBe('confirmed'));
  const point = hook.result.current.stateRef.current.points[0];
  expect(api.create).toHaveBeenCalledOnce(); expect(api.update).not.toHaveBeenCalled();
  return { ...hook, point, rule, queueRule,
    edit: (patch: Partial<Mapping>) => act(() => {
      const previous = hook.result.current.stateRef.current;
      const action: WorkbenchV2Action = { type: 'updateMapping', pointId: point.id, patch };
      hook.result.current.stateRef.current = workbenchV2Reducer(previous, action);
      hook.result.current.mapping.afterMappingAction(previous, action, hook.result.current.stateRef.current);
    }),
    inspect: async () => { await act(async () => { await hook.result.current.submission.port.inspect(request.requestId); }); },
  };
}

it.each([{ unit: 'MWh' }, { scale: 2 }])('keeps a later lost mapping ACK unknown and never resends after Inspect: %j', async patch => {
  const h = await confirmedRange();
  api.update.mockRejectedValue(new Error('response lost'));
  h.edit(patch);
  await waitFor(() => expect(h.result.current.submission.port.observations()[0].status).toBe('unknown'));
  expect(api.update).toHaveBeenCalledOnce();
  expect(serverMapping).toMatchObject({ unit: 'kWh', scale: 1 });
  await h.inspect(); await h.inspect();
  expect(api.update).toHaveBeenCalledOnce();
  // An old successful receipt and current old values cannot settle this PUT.
  expect(h.result.current.submission.port.observations()[0]).toMatchObject({ status: 'unknown', inspection: { status: 'current-state' } });
  expect(h.result.current.submission.canWrite(h.rule.id)).toBe(false);
  expect(h.result.current.submission.port.canSelect?.(device.id, 'P')).toBe(false);
  expect(h.result.current.stateRef.current.mappings[h.point.id]).toMatchObject(patch);
  expect(api.update).toHaveBeenCalledOnce(); expect(api.create).toHaveBeenCalledOnce();
  expect(h.queueRule).toHaveBeenCalledOnce();
});

it('cannot confirm a rejected later mapping ACK using the old successful receipt', async () => {
  const h = await confirmedRange();
  api.update.mockImplementation(async () => structuredClone(serverMapping));
  h.edit({ unit: 'MWh' });
  await waitFor(() => expect(h.result.current.submission.port.observations()[0].status).toBe('unknown'));
  await h.inspect();
  expect(h.result.current.submission.port.observations()[0]).toMatchObject({ status: 'unknown', inspection: { status: 'current-state' } });
  expect(h.result.current.submission.canWrite(h.rule.id)).toBe(false);
  expect(h.result.current.submission.port.canSelect?.(device.id, 'P')).toBe(false);
  expect(h.result.current.stateRef.current.mappings[h.point.id].unit).toBe('MWh');
  expect(api.update).toHaveBeenCalledOnce(); expect(api.create).toHaveBeenCalledOnce();
});

it('does not let an older in-flight mapping readback unlock a later lost ACK or resend its PUT', async () => {
  const h = await confirmedRange();
  let finishReadback!: () => void;
  api.tag.mockImplementationOnce(() => new Promise(resolve => {
    // This canonical GET belongs to the first PUT, not the later uncertain one.
    const tag = { id: 'T', key: serverMapping!.tag_key, display_name: serverMapping!.display_name,
      unit: 'MWh', data_type: 'int16', status: 'active' };
    finishReadback = () => resolve(tag);
  }));
  h.edit({ unit: 'MWh' });
  await waitFor(() => expect(api.tag).toHaveBeenCalledTimes(2));
  expect(serverMapping).toMatchObject({ unit: 'MWh' });
  api.update.mockRejectedValue(new Error('later response lost'));
  h.edit({ unit: 'GWh' });
  await waitFor(() => expect(h.result.current.submission.port.observations()[0].status).toBe('unknown'));
  expect(api.update).toHaveBeenCalledTimes(2);
  await act(async () => { finishReadback(); });
  expect(api.update).toHaveBeenCalledTimes(2);
  expect(h.result.current.submission.port.observations()[0].status).toBe('unknown');
  expect(h.result.current.submission.canWrite(h.rule.id)).toBe(false);
  expect(h.result.current.stateRef.current.mappings[h.point.id].unit).toBe('GWh');
  expect(api.create).toHaveBeenCalledOnce(); expect(h.queueRule).toHaveBeenCalledOnce();
});

it('keeps a newer confirmed mapping when an older in-flight readback later fails', async () => {
  const h = await confirmedRange();
  let rejectReadback!: () => void;
  api.tag.mockImplementationOnce(() => new Promise((_resolve, reject) => {
    rejectReadback = () => reject(new Error('older readback failed'));
  }));
  h.edit({ unit: 'MWh' });
  await waitFor(() => expect(api.tag).toHaveBeenCalledTimes(2));
  h.edit({ unit: 'GWh', tag_key: 'meter-final', display_name: 'Final meter' });
  await waitFor(() => expect(h.result.current.submission.port.observations()[0]).toMatchObject({
    status: 'confirmed', candidates: { candidates: [{ tag_key: 'meter-final', label: 'Final meter' }] },
  }));
  await act(async () => { rejectReadback(); });
  expect(h.result.current.submission.port.observations()[0]).toMatchObject({
    status: 'confirmed', candidates: { candidates: [{ tag_key: 'meter-final', label: 'Final meter' }] },
  });
  expect(h.result.current.submission.canWrite(h.rule.id)).toBe(true);
  expect(h.result.current.stateRef.current.mappings[h.point.id]).toMatchObject({ unit: 'GWh', tag_key: 'meter-final', save_state: 'saved' });
  expect(api.update).toHaveBeenCalledTimes(2); expect(api.create).toHaveBeenCalledOnce();
});
