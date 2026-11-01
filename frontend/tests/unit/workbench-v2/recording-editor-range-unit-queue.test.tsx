import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { useRef } from 'react';
import { useRecordingRangeSubmission } from '@/pages/datalink/workbench-v2/useRecordingRangeSubmission';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { RangeSourceRequest } from '@/pages/datalink/workbench-v2/rangeSourceSubmission';
import { toStudioV2RuleCreateRequest } from '@/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { toStudioV2MappingUpdateRequest } from '@/features/datalink/workbench-v2/state/studioV2MappingAutosave';
import { readRuleProof, readRangeCompletion } from '@/pages/datalink/workbench-v2/rangeSourceReadback';
vi.mock('@/pages/datalink/workbench-v2/rangeSourceReadback', async original => ({
  ...await original<typeof import('@/pages/datalink/workbench-v2/rangeSourceReadback')>(), readRuleProof: vi.fn(), readRangeCompletion: vi.fn(),
}));
const device = { id: 'dev-1', name: 'PLC', description: '', protocol: 'modbus_tcp' as const, status: 'draft' as const,
  connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' };
const request: RangeSourceRequest & { mappingDefaults: { unit: string } } = { workspaceId: 'W', sessionId: 'S', requestId: 'REQ', generation: 1,
  deviceId: 'dev-1', mappingDefaults: { unit: 'kWh' }, buffer: { name: 'Meter', naming_prefix: 'METER', start_address: '40001', count: 1,
    data_type: 'float32', enabled: true, scale_multiplier: 1, scale_offset: 0, data_format: 'CDAB', skipped_addresses: [],
    share_enabled: false, share_start_register: null, share_stride: null } };
beforeEach(() => vi.clearAllMocks());
it('writes the frozen unit through the original mapping queue, and rejects mismatched unit ACK', async () => {
  const queueRule = vi.fn<(action: WorkbenchV2Action) => void>();
  const queueMapping = vi.fn();
  const h = renderHook(() => {
    const stateRef = useRef(savedGroupState());
    const submission = useRecordingRangeSubmission({ workspaceId: 'W', enabled: true, stateRef,
      apply: action => { stateRef.current = workbenchV2Reducer(stateRef.current, action); },
      prepareDevice: async () => ({ device, activationAllowed: true }), diagnoseDevice: async () => false, queueRule, queueMapping });
    return { submission, stateRef };
  });
  const sent = structuredClone(request);
  act(() => { void h.result.current.submission.port.submit(sent); });
  await waitFor(() => expect(queueRule).toHaveBeenCalledTimes(1));
  // Mutating a caller's later object must not alter the already submitted intent.
  sent.mappingDefaults.unit = 'MWh';
  const action = queueRule.mock.calls[0][0];
  if (action.type !== 'addRule') throw new Error('Expected the existing addRule owner');
  const rule = action.rule;
  expect(rule.data_format).toBe('CDAB');
  expect(rule).not.toHaveProperty('unit');
  const ack = { ...toStudioV2RuleCreateRequest(rule), id: rule.id, workspace_id: 'W', revision_id: 'R1', locked: false,
    origin: 'manual' as const, skipped_addresses: [], runtime_apply_status: 'not_running' as const, created_at: '', updated_at: '' };
  vi.mocked(readRuleProof).mockResolvedValue({ request, receipt: { submitted: toStudioV2RuleCreateRequest(rule), acknowledged: ack }, record: ack, device,
    snapshot: { source_rule_id: rule.id, workspace_id: 'W', workspace_revision: 'SOURCE', revision_id: 'R1',
      tags: { status: 'ready', candidates: [{ id: 'C', point_id: 'P', address: '40001' }] } } as unknown as NonNullable<Awaited<ReturnType<typeof readRuleProof>>>['snapshot'] });
  act(() => { expect(h.result.current.submission.ruleObservation.acknowledged(rule, ack)).toBe(true); h.result.current.submission.ruleObservation.saved(rule, ack); });
  await waitFor(() => expect(queueMapping).toHaveBeenCalledTimes(2));
  const point = h.result.current.stateRef.current.points.find(row => row.rule_id === rule.id)!;
  const mapping = h.result.current.stateRef.current.mappings[point.id];
  expect(mapping.unit).toBe('kWh');
  expect(queueMapping.mock.calls[1][1]).toMatchObject({ type: 'updateMapping', pointId: point.id, patch: { unit: 'kWh' } });
  const submitted = toStudioV2MappingUpdateRequest(point, mapping);
  act(() => { expect(h.result.current.submission.mappingObservation.acknowledged(point, mapping, { ...submitted, unit: 'MWh', id: 'MAP',
    workspace_id: 'W', device_id: device.id, point_id: 'P', tag_id: 'T', runtime_apply_status: 'not_running', created_at: '', updated_at: '' })).toBe(false); });
  expect(h.result.current.submission.port.observations()[0].status).toBe('unknown');
  expect(readRangeCompletion).not.toHaveBeenCalled();
  expect(queueRule).toHaveBeenCalledTimes(1);
});
it('rejects changed mapping unit during offline preparation before diagnosis or another queue write', async () => {
  const queueRule = vi.fn<(action: WorkbenchV2Action) => void>(), queueMapping = vi.fn(), diagnoseDevice = vi.fn(async () => true);
  const h = renderHook(() => {
    const stateRef = useRef(savedGroupState());
    const submission = useRecordingRangeSubmission({ workspaceId: 'W', enabled: true, stateRef,
      apply: action => { stateRef.current = workbenchV2Reducer(stateRef.current, action); },
      prepareDevice: async () => ({ device, activationAllowed: false }), diagnoseDevice, queueRule, queueMapping });
    return { submission, stateRef };
  });
  const inactive = { ...structuredClone(request), buffer: { ...request.buffer, enabled: false } };
  act(() => { void h.result.current.submission.port.submit(inactive); });
  await waitFor(() => expect(queueRule).toHaveBeenCalledTimes(1));
  const action = queueRule.mock.calls[0][0];
  if (action.type !== 'addRule') throw new Error('Expected the existing addRule owner');
  const rule = action.rule;
  const ack = { ...toStudioV2RuleCreateRequest(rule), id: rule.id, workspace_id: 'W', revision_id: 'R1', locked: false,
    origin: 'manual' as const, skipped_addresses: [], runtime_apply_status: 'not_running' as const, created_at: '', updated_at: '' };
  vi.mocked(readRuleProof).mockResolvedValue({ request: inactive, receipt: { submitted: toStudioV2RuleCreateRequest(rule), acknowledged: ack }, record: ack, device,
    snapshot: { source_rule_id: rule.id, workspace_id: 'W', workspace_revision: 'SOURCE', revision_id: 'R1',
      tags: { status: 'ready', candidates: [{ id: 'C', point_id: 'P', address: '40001' }] } } as unknown as NonNullable<Awaited<ReturnType<typeof readRuleProof>>>['snapshot'] });
  vi.mocked(readRangeCompletion).mockResolvedValue({ candidates: [], excluded: [], inactive: true, savedPoints: [{ device_id: device.id, point_id: 'P', tag_id: 'T' }] });
  act(() => { h.result.current.submission.ruleObservation.acknowledged(rule, ack); h.result.current.submission.ruleObservation.saved(rule, ack); });
  await waitFor(() => expect(queueMapping).toHaveBeenCalledTimes(2));
  const point = h.result.current.stateRef.current.points.find(row => row.rule_id === rule.id)!;
  const mapping = h.result.current.stateRef.current.mappings[point.id];
  const saved = { ...toStudioV2MappingUpdateRequest(point, mapping), id: 'MAP', workspace_id: 'W', device_id: device.id,
    point_id: 'P', tag_id: 'T', runtime_apply_status: 'not_running' as const, created_at: '', updated_at: '' };
  act(() => { expect(h.result.current.submission.mappingObservation.acknowledged(point, mapping, saved)).toBe(true); h.result.current.submission.mappingObservation.saved(point, mapping, saved); });
  await waitFor(() => expect(h.result.current.submission.port.observations()[0].status).toBe('saved-inactive'));
  const changed = { ...inactive, requestId: 'PREPARE', buffer: { ...inactive.buffer, enabled: true }, mappingDefaults: { unit: 'MWh' } };
  await act(async () => { expect(await h.result.current.submission.port.prepareSaved(inactive.requestId, changed)).toMatchObject({ status: 'blocked' }); });
  expect(diagnoseDevice).not.toHaveBeenCalled();
  expect(queueRule).toHaveBeenCalledTimes(1); expect(queueMapping).toHaveBeenCalledTimes(2);
});
