import { act, renderHook, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useRef } from 'react';
import { useRecordingRangeSubmission } from '@/pages/datalink/workbench-v2/useRecordingRangeSubmission';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { workbenchV2Reducer, type WorkbenchV2Action } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { RangeSourceRequest } from '@/pages/datalink/workbench-v2/rangeSourceSubmission';
import type { Rule } from '@/features/datalink/workbench-v2/state/types';
import { toStudioV2RuleCreateRequest } from '@/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { readRuleProof, readRangeCompletion } from '@/pages/datalink/workbench-v2/rangeSourceReadback';
vi.mock('@/pages/datalink/workbench-v2/rangeSourceReadback', () => ({ readRuleProof: vi.fn(), readRangeCompletion: vi.fn() }));
const device = { id: 'dev-1', name: 'PLC', description: '', protocol: 'modbus_tcp' as const, status: 'draft' as const, connection_config: '{"host":"192.0.2.1","port":502}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' };
const request: RangeSourceRequest = { workspaceId: 'W', sessionId: 'S', requestId: 'REQ', generation: 1, deviceId: 'dev-1', buffer: {
  name: 'Range', naming_prefix: 'Range', start_address: '40001', count: 1, data_type: 'int16', enabled: true, scale_multiplier: 1, scale_offset: 0,
  data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null,
} };
function harness() {
 const ruleQueue = vi.fn<(action: WorkbenchV2Action) => void>();
 const mappingQueue = vi.fn();
 const hook = renderHook(({ workspaceId }) => {
  const stateRef = useRef(savedGroupState());
  const submission = useRecordingRangeSubmission({ workspaceId, enabled: true, stateRef,
   apply: action => { stateRef.current = workbenchV2Reducer(stateRef.current, action); }, prepareDevice: async () => device,
   queueRule: ruleQueue, queueMapping: mappingQueue });
  return { submission, stateRef };
 }, { initialProps: { workspaceId: 'W' } });
 return { ...hook, ruleQueue, mappingQueue };
}
describe('range observation uses the existing queue and causal receipts', () => {
 it('queues once, waits beyond rule acknowledgement, and never retransmits a lost create', async () => {
  const h = harness(); let pending!: Promise<unknown>;
  act(() => { pending = h.result.current.submission.port.submit(request); });
  await waitFor(() => expect(h.ruleQueue).toHaveBeenCalledTimes(1));
  const action = h.ruleQueue.mock.calls[0][0]; expect(action.type).toBe('addRule');
  const rule = (action as { type: 'addRule'; rule: Rule }).rule;
  expect(h.result.current.submission.port.observations()[0].status).toBe('pending');
  act(() => h.result.current.submission.ruleObservation.failed(rule));
  await act(async () => { await pending; });
  expect(h.result.current.submission.port.observations()[0].status).toBe('unknown');
  await act(async () => { await h.result.current.submission.port.submit({ ...request, requestId: 'REPLACEMENT' }); });
  expect(h.ruleQueue).toHaveBeenCalledTimes(1); expect(h.mappingQueue).not.toHaveBeenCalled();
 });
 it('does not fabricate mapping completion from a rule POST acknowledgement', async () => {
  vi.mocked(readRuleProof).mockResolvedValue(undefined);
  const h = harness(); let pending!: Promise<unknown>;
  act(() => { pending = h.result.current.submission.port.submit(request); });
  await waitFor(() => expect(h.ruleQueue).toHaveBeenCalledTimes(1));
  const rule = (h.ruleQueue.mock.calls[0][0] as { type: 'addRule'; rule: Rule }).rule;
  const ack = { ...toStudioV2RuleCreateRequest(rule), id: rule.id, workspace_id: 'W', revision_id: 'R1', locked: false, origin: 'manual' as const, skipped_addresses: [], runtime_apply_status: 'not_running' as const, created_at: '', updated_at: '' };
  act(() => { expect(h.result.current.submission.ruleObservation.acknowledged(rule, ack)).toBe(true); h.result.current.submission.ruleObservation.saved(rule, ack); });
  await act(async () => { await pending; });
  expect(h.result.current.submission.port.observations()[0].status).toBe('unknown');
  expect(readRangeCompletion).not.toHaveBeenCalled(); expect(h.mappingQueue).not.toHaveBeenCalled();
 });
 it('rechecks a known rule receipt and starts only never-submitted mappings, not a second rule', async () => {
  vi.mocked(readRuleProof).mockResolvedValue(undefined);
  const h = harness(); let pending!: Promise<unknown>;
  act(() => { pending = h.result.current.submission.port.submit(request); });
  await waitFor(() => expect(h.ruleQueue).toHaveBeenCalledTimes(1));
  const rule = (h.ruleQueue.mock.calls[0][0] as { type: 'addRule'; rule: Rule }).rule;
  const ack = { ...toStudioV2RuleCreateRequest(rule), id: rule.id, workspace_id: 'W', revision_id: 'R1', locked: false, origin: 'manual' as const, skipped_addresses: [], runtime_apply_status: 'not_running' as const, created_at: '', updated_at: '' };
  act(() => { h.result.current.submission.ruleObservation.acknowledged(rule, ack); h.result.current.submission.ruleObservation.saved(rule, ack); });
  await act(async () => { await pending; });
  vi.mocked(readRuleProof).mockResolvedValue({ request, receipt: { submitted: toStudioV2RuleCreateRequest(rule), acknowledged: ack }, record: ack, device,
    snapshot: { source_rule_id: rule.id, workspace_id: 'W', workspace_revision: 'SOURCE', revision_id: 'R1', tags: { status: 'ready', candidates: [{ id: 'C', point_id: 'P', address: '40001' }] } } as unknown as NonNullable<Awaited<ReturnType<typeof readRuleProof>>>['snapshot'] });
  h.mappingQueue.mockImplementation(() => { expect(h.result.current.submission.canWrite(rule.id)).toBe(true); });
  await act(async () => { await h.result.current.submission.port.inspect(request.requestId); });
  expect(h.mappingQueue).toHaveBeenCalledTimes(2);
  expect(h.result.current.submission.port.observations()[0].status).toBe('pending');
  expect(h.ruleQueue).toHaveBeenCalledTimes(1);
  const point = h.result.current.stateRef.current.points.find(p => p.rule_id === rule.id)!;
  act(() => h.result.current.submission.mappingObservation.failed(point));
  await act(async () => { await h.result.current.submission.port.inspect(request.requestId); });
  expect(h.mappingQueue).toHaveBeenCalledTimes(2);
  expect(h.result.current.submission.port.observations()[0].status).toBe('unknown');
 });

 it('does not advance a late proof into mappings after owner workspace changes', async () => {
  let resolve!: (proof: Awaited<ReturnType<typeof readRuleProof>>) => void;
  vi.mocked(readRuleProof).mockImplementationOnce(() => new Promise(yes => { resolve = yes; }));
  const h = harness(); let pending!: Promise<unknown>;
  act(() => { pending = h.result.current.submission.port.submit(request); });
  await waitFor(() => expect(h.ruleQueue).toHaveBeenCalledTimes(1));
  const rule = (h.ruleQueue.mock.calls[0][0] as { type: 'addRule'; rule: Rule }).rule;
  const ack = { ...toStudioV2RuleCreateRequest(rule), id: rule.id, workspace_id: 'W', revision_id: 'R1', locked: false, origin: 'manual' as const, runtime_apply_status: 'not_running' as const, skipped_addresses: [], created_at: '', updated_at: '' };
  act(() => { expect(h.result.current.submission.ruleObservation.acknowledged(rule, ack)).toBe(true); h.result.current.submission.ruleObservation.saved(rule, ack); });
  h.rerender({ workspaceId: 'FOREIGN' });
  await act(async () => { resolve({ request, receipt: { submitted: toStudioV2RuleCreateRequest(rule), acknowledged: ack }, record: ack, device, snapshot: { tags: { candidates: [] } } as unknown as NonNullable<Awaited<ReturnType<typeof readRuleProof>>>['snapshot'] }); await pending; });
  expect(h.mappingQueue).not.toHaveBeenCalled();
  expect(h.result.current.submission.port.observations()).toEqual([]);
  h.rerender({ workspaceId: 'W' });
  await act(async () => { await h.result.current.submission.port.inspect(request.requestId); });
  expect(h.mappingQueue).not.toHaveBeenCalled();
  expect(h.result.current.submission.port.observations()[0].status).toBe('unknown');
 });

 it('does not start the next mutation phase after the original owner unmounts', async () => {
  let resolve!: (proof: Awaited<ReturnType<typeof readRuleProof>>) => void;
  vi.mocked(readRuleProof).mockImplementationOnce(() => new Promise(yes => { resolve = yes; }));
  const h = harness(); let pending!: Promise<unknown>;
  act(() => { pending = h.result.current.submission.port.submit(request); });
  await waitFor(() => expect(h.ruleQueue).toHaveBeenCalledTimes(1));
  const rule = (h.ruleQueue.mock.calls[0][0] as { type: 'addRule'; rule: Rule }).rule;
  const ack = { ...toStudioV2RuleCreateRequest(rule), id: rule.id, workspace_id: 'W', revision_id: 'R1', locked: false, origin: 'manual' as const, runtime_apply_status: 'not_running' as const, skipped_addresses: [], created_at: '', updated_at: '' };
  act(() => { h.result.current.submission.ruleObservation.acknowledged(rule, ack); h.result.current.submission.ruleObservation.saved(rule, ack); });
  h.unmount();
  await act(async () => { resolve({ request, receipt: { submitted: toStudioV2RuleCreateRequest(rule), acknowledged: ack }, record: ack, device,
    snapshot: { tags: { candidates: [{ address: '40001', point_id: 'P' }] } } as unknown as NonNullable<Awaited<ReturnType<typeof readRuleProof>>>['snapshot'] }); await Promise.resolve(); await Promise.resolve(); });
  expect(h.mappingQueue).not.toHaveBeenCalled();
  await pending;
 });

});
