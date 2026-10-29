import { StrictMode, useEffect } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { act, cleanup, render, renderHook } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { useCircuitRecordingCommandFence, type CircuitRecordingCommandFenceInput } from '@/features/datalink/workbench-v2/steps/step4/useCircuitRecordingCommandFence';
import { BatchPlanError, executeBatchPlan, type FrozenBatchPlan } from '@/features/datalink/workbench-v2/steps/step4/circuitProfileBatchExecutor';
import { newProfileDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBatchIntent';
import type { RecordingPreparationCommandFence } from '@/features/datalink/workbench-v2/steps/step4/recordingPreparationPresentation';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup } from '../../fixtures/writeGroupState';

const capabilityContext = { connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } };
function group(circuit = 'A', id = `G-${circuit}`): WriteGroup {
  const base = savedGroup();
  return savedGroup({ id, layout_version: 'managed_long_v1', storage_mode: 'plain',
    destination: { ...base.destination, storage_strategy: 'managed' },
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    members: base.members.map(member => ({ ...member, circuit_id: circuit })) });
}
const groups = [group('A'), group('B'), group('C')];
function input(overrides: Partial<CircuitRecordingCommandFenceInput> = {}): CircuitRecordingCommandFenceInput {
  return { ownerKey: 'owner-A', scopeKey: 'scope-A-R1', workspaceId: 'ws-1', deviceId: 'dev-1', circuitId: 'A', batchGroups: groups,
    updateWriteGroup: vi.fn<CircuitRecordingCommandFenceInput['updateWriteGroup']>(async ({ id }) => ({ workspace_revision: 'ws-next', group: group(id.slice(2), id) })),
    ensureCircuitProfile: vi.fn<CircuitRecordingCommandFenceInput['ensureCircuitProfile']>(async () => ({ workspace_revision: 'ws-next', group: group('A', 'NEW') })), ...overrides };
}
function report(isBlocked: RecordingPreparationCommandFence['isBlocked'] = vi.fn(() => false)): RecordingPreparationCommandFence {
  return { ownerKey: 'owner-A', scopeKey: 'scope-A-R1', blocked: false, isBlocked };
}
function update(target = groups[0]) {
  return { id: target.id, request: { workspace_id: 'ws-1', expected_workspace_revision: 'ws-1', expected_group_revision: target.revision,
    expected_connector_revision: 'crev-1', group: target } };
}
function ensure(target = groups[0]) {
  return { request: { workspace_id: 'ws-1', expected_workspace_revision: 'ws-1', expected_connector_revision: 'crev-1',
    group: newProfileDraft('ws-1', { circuit_id: target.members[0].circuit_id!, circuit_name: '', anyGroup: target }, 'cumulative', 10) } };
}
function plan(targets = [groups[1], groups[0], groups[2]]): FrozenBatchPlan {
  return { workspaceId: 'ws-1', workspaceRevision: 'ws-1', connectorRevision: 'crev-1', applyInst: true, instInterval: 30, applyCum: false, cumInterval: 10,
    targets: targets.map(target => ({ circuit_id: target.members[0].circuit_id!, circuit_name: '', instantaneousGroup: target, anyGroup: target })) };
}
afterEach(cleanup);

describe('current circuit batch denial adapter', () => {
  it('defaults blocked and rejects current-circuit writes without a current owner report', async () => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    expect(view.result.current.blocked).toBe(true);
    await expect(view.result.current.updateWriteGroup(update())).rejects.toMatchObject({ outcome: 'failed', code: 'RECORDING_START_BUSY' });
    expect(props.updateWriteGroup).not.toHaveBeenCalled();
  });
  it('uses fresh aggregate denial at dispatch rather than the cached blocked flag or group-only clearance', async () => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props)); let denied = false;
    const isBlocked = vi.fn((groupId?: string) => groupId ? false : denied);
    act(() => view.result.current.onCommandFence(report(isBlocked)));
    expect(view.result.current.blocked).toBe(false); denied = true;
    await expect(view.result.current.updateWriteGroup(update())).rejects.toMatchObject({ outcome: 'failed' });
    expect(isBlocked).toHaveBeenLastCalledWith(); expect(props.updateWriteGroup).not.toHaveBeenCalled();
  });
  it('fails closed when an accepted owner denial callback throws, without leaking its raw exception', async () => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    act(() => view.result.current.onCommandFence(report(() => { throw new Error('private endpoint password'); })));
    await expect(view.result.current.updateWriteGroup(update())).rejects.toEqual({ outcome: 'failed', code: 'RECORDING_START_BUSY' });
    expect(props.updateWriteGroup).not.toHaveBeenCalled();
  });
  it('allows an unrelated canonical circuit without importing selected-circuit recovery authority', async () => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    const request = update(groups[1]); await view.result.current.updateWriteGroup(request);
    expect(props.updateWriteGroup).toHaveBeenCalledExactlyOnceWith(request);
  });
  it('preserves actual executor earlier successes and marks a denied selected target failed, not unknown', async () => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    const outcome = await executeBatchPlan({ plan: plan(), capabilityContext, updateWriteGroup: view.result.current.updateWriteGroup }).catch(error => error);
    expect(outcome).toBeInstanceOf(BatchPlanError);
    expect(outcome.updatedGroups.map((saved: WriteGroup) => saved.id)).toEqual(['G-B']);
    expect(outcome.workspaceRevision).toBe('ws-next');
    expect(outcome.targetResults.B.profiles.instantaneous.status).toBe('saved');
    expect(outcome.targetResults.A.profiles.instantaneous.status).toBe('failed');
    expect(outcome.targetResults.C.profiles.instantaneous.status).toBe('not-attempted');
    expect(props.updateWriteGroup).toHaveBeenCalledTimes(1);
  });
  it('does not block the selected target merely because an earlier batch target refreshed its revision', async () => {
    const props = input(); const view = renderHook(current => useCircuitRecordingCommandFence(current), { initialProps: props });
    act(() => view.result.current.onCommandFence(report())); const guardedUpdate = view.result.current.updateWriteGroup;
    vi.mocked(props.updateWriteGroup).mockImplementationOnce(async () => {
      const updated = { ...groups[1], revision: 'rev-2' };
      act(() => view.rerender({ ...props, batchGroups: [groups[0], updated, groups[2]] }));
      return { workspace_revision: 'ws-next', group: updated };
    });
    const outcome = await executeBatchPlan({ plan: plan([groups[1], groups[0]]), capabilityContext, updateWriteGroup: guardedUpdate });
    expect(outcome.updatedGroups.map(saved => saved.id)).toEqual(['G-B', 'G-A']);
    expect(props.updateWriteGroup).toHaveBeenCalledTimes(2);
    expect(outcome.targetResults.A.profiles.instantaneous?.status).toBe('saved');
  });
  it('keeps same-target earlier profile success when a new denial arrives before legitimate missing-profile ensure', async () => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props)); let denied = false;
    act(() => view.result.current.onCommandFence(report(() => denied)));
    vi.mocked(props.updateWriteGroup).mockImplementationOnce(async () => { denied = true; return { workspace_revision: 'ws-next', group: groups[0] }; });
    const outcome = await executeBatchPlan({ plan: { ...plan([groups[0]]), applyCum: true }, capabilityContext,
      updateWriteGroup: view.result.current.updateWriteGroup, ensureCircuitProfile: view.result.current.ensureCircuitProfile }).catch(error => error);
    expect(outcome).toBeInstanceOf(BatchPlanError); expect(outcome.targetResults.A.status).toBe('partial');
    expect(outcome.targetResults.A.profiles.instantaneous.status).toBe('saved'); expect(outcome.targetResults.A.profiles.cumulative.status).toBe('failed');
    expect(props.ensureCircuitProfile).not.toHaveBeenCalled();
  });
  it('allows legitimate ensure with no synthetic group ID after a healthy selected report', async () => {
    const props = input({ batchGroups: [] }); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    act(() => view.result.current.onCommandFence(report())); const request = ensure();
    await view.result.current.ensureCircuitProfile(request); expect(props.ensureCircuitProfile).toHaveBeenCalledExactlyOnceWith(request);
  });
  it.each(['workspace', 'group', 'device', 'circuit', 'revision'] as const)('rejects an update with foreign %s before native dispatch', async mismatch => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    act(() => view.result.current.onCommandFence(report())); const request = structuredClone(update());
    if (mismatch === 'workspace') request.request.workspace_id = 'foreign';
    if (mismatch === 'group') request.id = 'foreign';
    if (mismatch === 'device') request.request.group.members[0].device_id = 'foreign';
    if (mismatch === 'circuit') request.request.group.members[0].circuit_id = 'foreign';
    if (mismatch === 'revision') request.request.expected_group_revision = 'foreign';
    await expect(view.result.current.updateWriteGroup(request)).rejects.toMatchObject({ outcome: 'failed', code: 'BASIC_INTENT_CHANGED' });
    expect(props.updateWriteGroup).not.toHaveBeenCalled();
  });
  it.each(['foreign', 'deleted', 'duplicate'] as const)('rejects an invalid %s canonical target despite a healthy report', async invalid => {
    const batchGroups = invalid === 'duplicate' ? [groups[0], { ...groups[0] }] :
      [{ ...groups[0], ...(invalid === 'foreign' ? { workspace_id: 'foreign' } : { status: 'deleted' as const }) }];
    const props = input({ batchGroups }); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    act(() => view.result.current.onCommandFence(report()));
    await expect(view.result.current.updateWriteGroup(update())).rejects.toMatchObject({ outcome: 'failed', code: 'BASIC_INTENT_CHANGED' });
    expect(props.updateWriteGroup).not.toHaveBeenCalled();
  });
  it('rejects foreign ensure workspace/device/circuit but permits another canonical circuit unchanged', async () => {
    const props = input(); const view = renderHook(() => useCircuitRecordingCommandFence(props));
    for (const mutation of [(request: ReturnType<typeof ensure>) => { request.request.workspace_id = 'foreign'; },
      (request: ReturnType<typeof ensure>) => { request.request.group.members[0].device_id = 'foreign'; },
      (request: ReturnType<typeof ensure>) => { request.request.group.members.forEach(member => { member.circuit_id = 'foreign'; }); }]) {
      const request = structuredClone(ensure()); mutation(request);
      await expect(view.result.current.ensureCircuitProfile(request)).rejects.toMatchObject({ outcome: 'failed' });
    }
    const request = ensure(groups[1]); await view.result.current.ensureCircuitProfile(request);
    expect(props.ensureCircuitProfile).toHaveBeenCalledExactlyOnceWith(request);
  });
  it('ignores foreign reports, makes A→B→A old callbacks inert, and preserves the new report against old cleanup', async () => {
    const props = input(); const view = renderHook(current => useCircuitRecordingCommandFence(current), { initialProps: props });
    const old = view.result.current; act(() => old.onCommandFence({ ...report(), scopeKey: 'foreign' }));
    expect(view.result.current.blocked).toBe(true);
    view.rerender({ ...props, ownerKey: 'owner-B', scopeKey: 'scope-B', circuitId: 'B' }); view.rerender(props);
    act(() => old.onCommandFence(report())); expect(view.result.current.blocked).toBe(true);
    await expect(old.updateWriteGroup(update())).rejects.toMatchObject({ outcome: 'failed' });
    act(() => view.result.current.onCommandFence(report())); act(() => old.onCommandFence(null));
    expect(view.result.current.blocked).toBe(false); await view.result.current.updateWriteGroup(update());
    expect(props.updateWriteGroup).toHaveBeenCalledTimes(1);
  });
  it('does not reset on batch revision churn but rejects outdated target revisions and unmounted commands', async () => {
    const props = input(); const view = renderHook(current => useCircuitRecordingCommandFence(current), { initialProps: props });
    act(() => view.result.current.onCommandFence(report())); const old = view.result.current;
    view.rerender({ ...props, batchGroups: [...groups] }); expect(view.result.current.blocked).toBe(false);
    view.rerender({ ...props, batchGroups: [{ ...groups[0], revision: 'rev-2' }, ...groups.slice(1)] });
    expect(view.result.current.blocked).toBe(false); await expect(old.updateWriteGroup(update())).rejects.toMatchObject({ outcome: 'failed' });
    const latest = view.result.current; view.unmount(); await expect(latest.ensureCircuitProfile(ensure())).rejects.toMatchObject({ outcome: 'failed' });
    expect(props.updateWriteGroup).not.toHaveBeenCalled(); expect(props.ensureCircuitProfile).not.toHaveBeenCalled();
  });
});

describe('own accepted batch progress waits only for exact current controller handoff', () => {
  function handoffHost(strict = false) {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const cache = (revision: string, workspace_id = 'ws-1') => client.setQueryData(studioV2WorkspaceKeys.writeGroups(),
      { workspace_id, workspace_revision: revision, groups: structuredClone(groups) });
    cache('ws-1');
    const props = { ...input(), consentKey: 'native-controls-A' };
    let ownerBlocked = false;
    vi.mocked(props.updateWriteGroup).mockImplementation(async ({ id }) => {
      cache('ws-next'); return { workspace_revision: 'ws-next', group: group(id.slice(2), id) };
    });
    let current!: ReturnType<typeof useCircuitRecordingCommandFence>;
    let publicationAction: (() => void) | undefined;
    function Publisher() { useEffect(() => { const action = publicationAction; publicationAction = undefined; action?.(); }); return null; }
    function Host({ value }: { value: typeof props }) { current = useCircuitRecordingCommandFence(value); return <Publisher />; }
    const tree = (value: typeof props) => <QueryClientProvider client={client}>{strict ? <StrictMode><Host value={value} /></StrictMode> : <Host value={value} />}</QueryClientProvider>;
    const rendered = render(tree(props));
    const view = { result: { get current() { return current; } }, rerender: (value: typeof props) => rendered.rerender(tree(value)), unmount: rendered.unmount };
    const publish = (revision: string, state: 'settling' | 'current' | 'blocked', deny = revision !== 'ws-next') => {
      view.result.current.onCommandFence({ ...report(() => ownerBlocked || deny),
        handoff: { workspaceRevision: revision, consentKey: 'native-controls-A', state: () => ownerBlocked ? 'blocked' : state } });
    };
    act(() => publish('ws-1', 'current', false));
    let settled = false;
    const run = () => executeBatchPlan({ plan: plan([groups[1], groups[0]]), capabilityContext,
      updateWriteGroup: view.result.current.updateWriteGroup }).then(value => { settled = true; return value; }, (error: unknown) => { settled = true; return error; });
    return { client, props, view, cache, publish, run, onPublication: (action: () => void) => { publicationAction = action; }, settled: () => settled, blockOwner: () => { ownerBlocked = true; } };
  }
  it('fails closed when production consent is supplied without the shared canonical cache provider', async () => {
    const props = { ...input(), consentKey: 'native-controls-A' };
    const view = renderHook(() => useCircuitRecordingCommandFence(props));
    act(() => view.result.current.onCommandFence({ ...report(), handoff: { workspaceRevision: 'ws-1', consentKey: props.consentKey, state: () => 'current' } }));
    await expect(view.result.current.updateWriteGroup(update())).rejects.toMatchObject({ outcome: 'failed' });
    expect(props.updateWriteGroup).not.toHaveBeenCalled();
  });
  it.each(['current', 'settling'] as const)('does not wait for an unchanged idle old canonical head even when the controller says %s', async state => {
    const host = handoffHost(); act(() => host.publish('ws-1', state, false));
    vi.mocked(host.props.updateWriteGroup).mockResolvedValueOnce({ workspace_revision: 'ws-next', group: groups[1] });
    const result = await host.run();
    expect(result).toBeInstanceOf(BatchPlanError); expect(result).toMatchObject({ updatedGroups: [{ id: 'G-B' }], workspaceRevision: 'ws-next' });
    expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1); host.view.unmount(); host.client.clear();
  });
  it('rejects a changed target in the same canonical cache head despite stale React batchGroups', async () => {
    const host = handoffHost(); let result!: Promise<unknown>;
    host.onPublication(() => {
      host.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'ws-1', workspace_revision: 'ws-next',
        groups: [{ ...groups[0], revision: 'foreign-revision' }, ...groups.slice(1)] });
      host.publish('ws-next', 'current', false);
    });
    await act(async () => { result = host.run(); });
    expect(host.settled()).toBe(true); expect(await result).toMatchObject({ updatedGroups: [{ id: 'G-B' }], workspaceRevision: 'ws-next' });
    expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1); host.view.unmount(); host.client.clear();
  });
  it('allows a new explicit reviewed current-head command after a previous receipt without borrowing that receipt', async () => {
    const host = handoffHost(); await host.view.result.current.updateWriteGroup(update(groups[1]));
    act(() => { host.cache('new-explicit-head'); host.publish('new-explicit-head', 'current', false); });
    const request = update(); request.request.expected_workspace_revision = 'new-explicit-head';
    await host.view.result.current.updateWriteGroup(request);
    expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(2);
    expect(vi.mocked(host.props.updateWriteGroup).mock.calls[1][0]).toEqual(request);
    host.view.unmount(); host.client.clear();
  });
  it.each(['old current report', 'exact-head settling report'] as const)('terminates a permanently %s after publication under StrictMode with the accepted prefix intact', async mode => {
    const host = handoffHost(true);
    if (mode === 'exact-head settling report') vi.mocked(host.props.updateWriteGroup).mockImplementationOnce(async () => {
      host.cache('ws-next'); host.publish('ws-next', 'settling');
      return { workspace_revision: 'ws-next', group: groups[1] };
    });
    let result!: Promise<unknown>; await act(async () => { result = host.run(); });
    expect(host.settled()).toBe(true);
    expect(await result).toBeInstanceOf(BatchPlanError);
    expect(await result).toMatchObject({ updatedGroups: [{ id: 'G-B' }], workspaceRevision: 'ws-next',
      targetResults: { B: { status: 'saved' }, A: { profiles: { instantaneous: { status: 'failed' } } } } });
    expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1); host.view.unmount(); host.client.clear();
  });
  it('records the first accepted save immediately, waits for exact returned head report, then rechecks and dispatches the selected target', async () => {
    const host = handoffHost(); let result!: Promise<unknown>;
    host.onPublication(() => {
      expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1); expect(host.settled()).toBe(false);
      host.publish('ws-next', 'settling');
      host.onPublication(() => {
        expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1); host.publish('ws-next', 'current', false);
      });
    });
    await act(async () => { result = host.run(); });
    expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(2);
    expect(await result).toMatchObject({ updatedGroups: [{ id: 'G-B' }, { id: 'G-A' }] });
    expect(vi.mocked(host.props.updateWriteGroup).mock.calls[1][0].request.expected_workspace_revision).toBe('ws-next');
    host.view.unmount(); host.client.clear();
  });
  it.each(['external advance', 'foreign cache', 'cache error', 'cache missing', 'typed owner block', 'null report', 'scope change', 'native consent change', 'unmount'] as const)('aborts %s during handoff without hanging, resending, or forgetting the first accepted save', async event => {
      const host = handoffHost(); let result!: Promise<unknown>;
      host.onPublication(() => {
        expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1); expect(host.settled()).toBe(false);
        if (event === 'external advance') host.cache('external-head');
        if (event === 'foreign cache') host.cache('ws-next', 'foreign');
        if (event === 'cache error') host.client.getQueryCache().find({ queryKey: studioV2WorkspaceKeys.writeGroups(), exact: true })?.setState({ status: 'error', error: new Error('private cache error') });
        if (event === 'cache missing') host.client.removeQueries({ queryKey: studioV2WorkspaceKeys.writeGroups(), exact: true });
        if (event === 'typed owner block') { host.blockOwner(); host.publish('ws-next', 'blocked'); }
        if (event === 'null report') host.view.result.current.onCommandFence(null);
        if (event === 'scope change') host.view.rerender({ ...host.props, scopeKey: 'scope-B' });
        if (event === 'native consent change') host.view.rerender({ ...host.props, consentKey: 'changed-native-controls' });
        if (event === 'unmount') host.view.unmount();
      });
      await act(async () => { result = host.run(); });
      expect(host.settled()).toBe(true);
      const outcome = await result;
      expect(outcome).toBeInstanceOf(BatchPlanError);
      expect(outcome).toMatchObject({ updatedGroups: [{ id: 'G-B' }], workspaceRevision: 'ws-next',
        targetResults: { B: { status: 'saved' }, A: { profiles: { instantaneous: { status: 'failed' } } } } });
      expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1);
      if (event !== 'unmount') { act(() => host.publish('ws-next', 'current', false)); host.view.unmount(); }
      expect(host.props.updateWriteGroup).toHaveBeenCalledTimes(1); host.client.clear();
    });
});
