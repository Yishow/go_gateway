import { act, renderHook, waitFor, cleanup } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useStudioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { useStudioV2RuleAutosave } from '@/pages/datalink/workbench-v2/useStudioV2RuleAutosave';
import { INITIAL_STATE, workbenchV2Reducer, type WorkbenchV2Action, type useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { SourceRuleRecord } from '@/types/datalink';
import { locatorRead, nativeRule } from './studio-locator-fixtures';

const mutations = vi.hoisted(() => ({ update: vi.fn(), create: vi.fn(), remove: vi.fn() }));
vi.mock('@/hooks/datalink/useStudioV2Rules', () => ({
  useUpdateStudioV2RuleMutation: () => ({ mutateAsync: mutations.update }),
  useCreateStudioV2RuleMutation: () => ({ mutateAsync: mutations.create }),
  useDeleteStudioV2RuleMutation: () => ({ mutateAsync: mutations.remove }),
}));
beforeEach(() => vi.resetAllMocks());
afterEach(cleanup);
function harness() {
  const entry = locatorRead(2), ref = { current: structuredClone(INITIAL_STATE) as WorkbenchV2State };
  entry.entry.source_rules![0].count = 8;
  const actions: ReturnType<typeof useWorkbenchV2State> = { state: INITIAL_STATE, dispatch: vi.fn(), setView: vi.fn(), setCurrent: vi.fn(),
    completeStep: vi.fn(), toggleSidebar: vi.fn(), toggleSummaryRail: vi.fn(), setSidebarCollapsed: vi.fn(),
    setShowSummaryRail: vi.fn(), resetFlow: vi.fn(), selectRule: vi.fn() };
  const pending: Array<{ request: Record<string, unknown>; resolve: (ack: SourceRuleRecord) => void; reject: (error: Error) => void }> = [];
  mutations.update.mockImplementation(({ request }) => new Promise<SourceRuleRecord>((resolve, reject) => pending.push({ request, resolve, reject })));
  const hook = renderHook(({ selected }) => {
    const hydration = useStudioV2LocatorHydration(actions, ref, false, selected);
    const queue = useStudioV2RuleAutosave(actions, ref, undefined, undefined, 'W', hydration.captureRuleSave);
    return { hydration, queue };
  }, { initialProps: { selected: entry } });
  const edit = (count: number) => act(() => {
    const action: WorkbenchV2Action = { type: 'updateRule', ruleId: 'R', patch: { count } };
    ref.current = workbenchV2Reducer(ref.current, action); hook.result.current.queue.afterRuleAction(action);
  });
  const ack = async (index = 0, extra: Partial<SourceRuleRecord> = {}) => {
    await act(async () => pending[index].resolve({ ...nativeRule(), ...pending[index].request, revision_id: `ACK${index}`, runtime_apply_status: 'aligned', ...extra } as SourceRuleRecord));
  };
  const goB = () => hook.rerender({ selected: locatorRead(1, 'B') });
  return { ...hook, ref, entry, pending, edit, ack, goB };
}
describe('locator own-rule queue ACK baseline', () => {
  it('allows a fresh B navigation after count 8→1→8 is genuinely ACKed by the original queue', async () => {
    const h = harness(); h.edit(1); await h.ack(); h.edit(8); await h.ack(1);
    expect(h.ref.current.rules[0]).toMatchObject({ count: 8, revision_id: 'ACK1', save_state: 'saved' });
    expect(mutations.update).toHaveBeenCalledTimes(2); h.goB();
    expect(h.result.current.hydration.ready).toBe(true); expect(h.ref.current.devices[0].id).toBe('B');
  });
  it('does not accept a foreign pre-submit revision as this baseline owner revision', async () => {
    const h = harness(); h.ref.current.rules[0].revision_id = 'FOREIGN-REVISION';
    h.edit(1); await h.ack(); h.goB(); expect(h.result.current.hydration.blocked).toBe(true);
  });
  it('keeps an unsaved count edit refused by a genuinely new B selection', () => {
    const h = harness(); h.edit(1); h.goB(); expect(h.result.current.hydration.blocked).toBe(true);
  });
  it.each(['device', 'other-rule', 'mapping', 'settings', 'database', 'points', 'local-name'] as const)('never clears independent %s dirty state', async reason => {
    const h = harness(); h.edit(1);
    if (reason === 'device') h.ref.current.devices[0].name = 'Unsaved device';
    if (reason === 'other-rule') h.ref.current.rules.push({ ...h.ref.current.rules[0], id: 'OTHER' });
    if (reason === 'mapping') h.ref.current.mappings = { ...h.ref.current.mappings, LOCAL: { point_id: 'LOCAL' } as WorkbenchV2State['mappings'][string] };
    if (reason === 'settings') h.ref.current.settings.modbus_share.enabled = true;
    if (reason === 'database') h.ref.current.db.connector.name = 'Unsaved connector';
    if (reason === 'points') h.ref.current.points = [];
    if (reason === 'local-name') h.ref.current.rules[0].name = 'Local name not submitted';
    await h.ack(); h.goB(); expect(h.result.current.hydration.blocked).toBe(true);
  });
  it.each(['foreign-ack', 'foreign-workspace', 'server-different', 'late-view', 'ABA-view', 'failed'] as const)('does not advance on %s', async reason => {
    const h = harness(); h.edit(1);
    if (reason === 'late-view' || reason === 'ABA-view') h.goB();
    if (reason === 'ABA-view') h.rerender({ selected: structuredClone(h.entry) });
    if (reason === 'failed') await act(async () => h.pending[0].reject(new Error('Rejected')));
    else await h.ack(0, reason === 'foreign-ack' ? { id: 'FOREIGN' } : reason === 'foreign-workspace' ? { workspace_id: 'OTHER' } : reason === 'server-different' ? { count: 3 } : {});
    h.goB(); expect(h.result.current.hydration.blocked).toBe(true);
  });
  it('does not revive a pending rule ACK after returning to the exact retained editor object', async () => {
    const h = harness(); h.edit(1); h.goB(); expect(h.result.current.hydration.blocked).toBe(true);
    h.rerender({ selected: h.entry });
    expect(h.result.current.hydration.ready).toBe(true); expect(h.ref.current.rules[0].count).toBe(1);
    await h.ack();
    expect(h.ref.current.rules[0].save_state).toBe('saved');
    h.goB(); expect(h.result.current.hydration.blocked).toBe(true);
    expect(h.ref.current.rules[0].count).toBe(1);
  });
  it('does not clean a newer A→B→A edit before that newer queued request receives its own ACK', async () => {
    const h = harness(); h.edit(1); h.edit(2); h.edit(1); await h.ack();
    await waitFor(() => expect(h.pending).toHaveLength(2));
    expect(h.ref.current.rules[0].save_state).toBe('saving'); h.goB();
    expect(h.result.current.hydration.blocked).toBe(true);
  });
});
