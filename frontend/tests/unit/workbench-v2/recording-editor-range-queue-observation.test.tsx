import { act, renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useStudioV2RuleAutosave } from '../../../src/pages/datalink/workbench-v2/useStudioV2RuleAutosave';
import { useStudioV2MappingAutosave } from '../../../src/pages/datalink/workbench-v2/useStudioV2MappingAutosave';
import type { MappingQueueObservation, RuleQueueObservation } from '../../../src/pages/datalink/workbench-v2/sourceQueueObservation';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { studioV2WorkspaceKeys } from '../../../src/hooks/datalink/keys';
const api = vi.hoisted(() => ({ rule: vi.fn(), mapping: vi.fn(), list: vi.fn() }));
vi.mock('../../../src/services/studioV2Rules', () => ({ studioV2RulesAPI: { create: api.rule, update: api.rule, remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: api.list, create: api.mapping, update: api.mapping, remove: vi.fn() } }));
function deferred<T>() { let resolve!: (value: T) => void, reject!: (reason: Error) => void; const promise = new Promise<T>((yes, no) => { resolve = yes; reject = no; }); return { promise, resolve, reject }; }
function harness(ruleObservation?: RuleQueueObservation, mappingObservation?: MappingQueueObservation) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const state = savedGroupState(); state.rules[0] = { ...state.rules[0], start_address: '40001', count: 1, save_state: 'saved', persisted: true };
  const point = state.points[0]; state.mappings[point.id] = { ...state.mappings[point.id], persisted: false, mapping_id: undefined };
  const ref = { current: state }, actions = { dispatch: vi.fn() }, allowed = { current: true };
  const hook = renderHook(() => ({
    rule: useStudioV2RuleAutosave(actions, ref, () => allowed.current, ruleObservation),
    mapping: useStudioV2MappingAutosave(actions, ref, false, () => allowed.current, mappingObservation),
  }), { wrapper: ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider> });
  return { ...hook, ref, client, point, allowed,
    setRuleObserver: (next: RuleQueueObservation) => { ruleObservation = next; hook.rerender(); },
    setMappingObserver: (next: MappingQueueObservation) => { mappingObservation = next; hook.rerender(); },
    saveRule: () => { ref.current = { ...ref.current, rules: [{ ...ref.current.rules[0], persisted: false }] }; act(() => hook.result.current.rule.afterRuleAction({ type: 'updateRule', ruleId: 'rule-1', patch: { count: 1 } })); },
    saveMapping: () => act(() => hook.result.current.mapping.afterMappingAction(ref.current, { type: 'updateMapping', pointId: point.id, patch: { scale: 1 } }, ref.current)),
  };
}
const ruleAck = { ...savedGroupState().rules[0], id: 'rule-1', device_id: 'dev-1', start_address: '40001', count: 1, workspace_id: 'w', revision_id: 'r', runtime_apply_status: 'applied' };
const mappingAck = { id: 'mapping-1', workspace_id: 'w', point_id: 'canonical-p', tag_id: 't', rule_id: 'rule-1', device_id: 'dev-1', address: '40001', tag_key: 'line.temperature', display_name: 'line.temperature', unit: '', target_type: 'float64', scale: 1, offset: 0, enabled: true, runtime_apply_status: 'applied' };
beforeEach(() => { vi.resetAllMocks(); api.list.mockResolvedValue([]); });
describe('actual rule and mapping queue observation', () => {
  it('reports captured rule snapshot and saved only after local publication', async () => {
    const held = deferred<typeof ruleAck>(); api.rule.mockReturnValue(held.promise);
    const observer: RuleQueueObservation = { acknowledged: vi.fn(() => true), saved: vi.fn(), failed: vi.fn() };
    const current = harness(observer); current.saveRule();
    current.ref.current = { ...current.ref.current, rules: [{ ...current.ref.current.rules[0], count: 9 }] };
    vi.mocked(observer.saved).mockImplementation(() => expect(current.ref.current.rules[0].persisted).toBe(true));
    await act(async () => { held.resolve(ruleAck); });
    await waitFor(() => expect(observer.saved).toHaveBeenCalledOnce());
    expect(observer.acknowledged).toHaveBeenCalledWith(expect.objectContaining({ count: 1 }), ruleAck);
    expect(current.ref.current.rules[0].count).toBe(9);
  });
  it('rejects rule acknowledgement without adopting foreign persisted identity or resending', async () => {
    api.rule.mockResolvedValue({ ...ruleAck, id: 'foreign' });
    const observer: RuleQueueObservation = { acknowledged: vi.fn(() => false), saved: vi.fn(), failed: vi.fn() };
    const current = harness(observer); current.saveRule();
    await waitFor(() => expect(current.ref.current.rules[0].save_state).toBe('save-error'));
    expect(current.ref.current.rules[0].persisted).toBe(false); expect(observer.saved).not.toHaveBeenCalled(); expect(api.rule).toHaveBeenCalledOnce();
  });
  it('reports mapping snapshot and publishes before saved callback', async () => {
    api.mapping.mockResolvedValue(mappingAck);
    const observer: MappingQueueObservation = { acknowledged: vi.fn(() => true), saved: vi.fn(), failed: vi.fn() };
    const current = harness(undefined, observer);
    vi.mocked(observer.saved).mockImplementation(() => expect(current.ref.current.mappings[current.point.id].persisted).toBe(true));
    current.saveMapping();
    await waitFor(() => expect(observer.saved).toHaveBeenCalledOnce());
    expect(observer.acknowledged).toHaveBeenCalledWith(current.point, expect.objectContaining({ scale: 1 }), mappingAck);
  });
  it('rejects mapping acknowledgement before cache or local persisted claims', async () => {
    api.mapping.mockResolvedValue({ ...mappingAck, id: 'foreign' });
    const observer: MappingQueueObservation = { acknowledged: vi.fn(() => false), saved: vi.fn(), failed: vi.fn() };
    const current = harness(undefined, observer); current.saveMapping();
    await waitFor(() => expect(current.ref.current.mappings[current.point.id].save_state).toBe('save-error'));
    expect(current.ref.current.mappings[current.point.id].persisted).toBe(false);
    expect(current.client.getQueryData(studioV2WorkspaceKeys.mappings())).toBeUndefined();
    expect(observer.saved).not.toHaveBeenCalled(); expect(api.mapping).toHaveBeenCalledOnce();
  });
  it.each(['rule', 'mapping'] as const)('reports lost %s response even if mutation permission changed', async kind => {
    const held = deferred<never>(); api[kind].mockReturnValue(held.promise);
    const observer = { acknowledged: vi.fn(() => true), saved: vi.fn(), failed: vi.fn() };
    const current = harness(kind === 'rule' ? observer : undefined, kind === 'mapping' ? observer : undefined);
    if (kind === 'rule') current.saveRule(); else current.saveMapping();
    current.allowed.current = false;
    await act(async () => { held.reject(new Error('lost response')); });
    await waitFor(() => expect(observer.failed).toHaveBeenCalledOnce());
    expect(observer.saved).not.toHaveBeenCalled();
  });
});

describe('request observer replacement and queued edits', () => {
  it.each(['rule', 'mapping'] as const)('uses latest %s observer and does not resend pending edits after a rejected acknowledgement', async kind => {
    const held = deferred<typeof mappingAck>(); api[kind].mockReturnValue(held.promise);
    const first = { acknowledged: vi.fn(() => true), saved: vi.fn(), failed: vi.fn() };
    const next = { acknowledged: vi.fn(() => false), saved: vi.fn(), failed: vi.fn() };
    const current = harness(kind === 'rule' ? first : undefined, kind === 'mapping' ? first : undefined);
    if (kind === 'rule') { current.saveRule(); current.saveRule(); current.setRuleObserver(next); }
    else { current.saveMapping(); current.saveMapping(); current.setMappingObserver(next); }
    await act(async () => { held.resolve(mappingAck); });
    await waitFor(() => expect(next.acknowledged).toHaveBeenCalledOnce());
    expect(first.acknowledged).not.toHaveBeenCalled(); expect(next.saved).not.toHaveBeenCalled();
    expect(api[kind]).toHaveBeenCalledOnce();
    expect(kind === 'rule' ? current.ref.current.rules[0].persisted : current.ref.current.mappings[current.point.id].persisted).toBe(false);
  });
});

it('observes an actually queued generation before mutation permission allows HTTP, then submits that generation once', async () => {
  api.rule.mockResolvedValue(ruleAck);
  const observer: RuleQueueObservation = { queued: vi.fn(), submitted: vi.fn(), acknowledged: vi.fn(() => true), saved: vi.fn(), failed: vi.fn() };
  const current = harness(observer); current.allowed.current = false; current.saveRule();
  expect(observer.queued).toHaveBeenCalledWith(expect.objectContaining({ id: 'rule-1', count: 1 }), 1);
  expect(observer.submitted).not.toHaveBeenCalled(); expect(api.rule).not.toHaveBeenCalled();
  current.allowed.current = true; current.rerender();
  await waitFor(() => expect(observer.saved).toHaveBeenCalledOnce());
  expect(observer.submitted).toHaveBeenCalledWith(expect.objectContaining({ id: 'rule-1', count: 1 }), 1);
  expect(observer.queued).toHaveBeenCalledOnce(); expect(api.rule).toHaveBeenCalledOnce();
});
