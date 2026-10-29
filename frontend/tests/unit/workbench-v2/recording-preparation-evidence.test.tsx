import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, render, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { RecordingPreparationEvidence, type RecordingPreparationEvidenceView } from '@/features/datalink/workbench-v2/steps/step4/RecordingPreparationEvidence';
import type { RecordingWorkflowScope } from '@/features/datalink/workbench-v2/state/recordingWorkflowView';
import type { WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { savedGroup } from '../../fixtures/writeGroupState';
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { readiness: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), update: vi.fn() } }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
const clients: QueryClient[] = [];
const scope: RecordingWorkflowScope = { workspaceId: 'ws-1', workspaceRevision: 'W1', deviceId: 'dev-1', groupId: 'group-1', groupRevision: 'rev-1', appliedRevision: '' };
function readiness(selected = scope): WriteGroupReadiness {
  return { workspace_id: selected.workspaceId, workspace_revision: selected.workspaceRevision, group_id: selected.groupId,
    group_revision: selected.groupRevision, applied_revision: selected.appliedRevision, config_ready: true, schema_ready: false, ready: false,
    issues: [{ code: 'schema-unverified', severity: 'blocking', step: 'Step 4', scope: '/private/scope', message: 'password=secret /private/db' }] };
}
function lifecycle(selected = scope): WriteGroupLifecycle {
  return { workspace_id: selected.workspaceId, workspace_revision: selected.workspaceRevision, group_id: selected.groupId,
    group_revision: selected.groupRevision, applied_revision: selected.appliedRevision, status: 'ready', connector_id: 'conn-1', connector_revision: 'crev-1',
    actions: [{ action: 'edit', allowed: true, reason_code: '', next_action: '', requires_review: false },
      { action: 'stop', allowed: false, reason_code: 'already_stopped', next_action: 'inspect', requires_review: false }],
    history: { canonical: false, current_group_id: selected.groupId, generations: [] } };
}
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(done => { resolve = done; }); return { promise, resolve }; }
function host(selected = scope, onEvidence = vi.fn<(value: RecordingPreparationEvidenceView | null) => void>()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const group = (value: RecordingWorkflowScope) => savedGroup({ id: value.groupId, workspace_id: value.workspaceId, revision: value.groupRevision, applied_revision: value.appliedRevision });
  const tree = (value: RecordingWorkflowScope, report = onEvidence) => <QueryClientProvider client={client}>
    <RecordingPreparationEvidence scope={value} group={group(value)} onEvidence={report} />
  </QueryClientProvider>;
  const view = render(tree(selected));
  return { ...view, onEvidence, current: () => onEvidence.mock.calls.at(-1)?.[0], rerenderScope: (value: RecordingWorkflowScope, report = onEvidence) => view.rerender(tree(value, report)) };
}
beforeEach(() => { vi.clearAllMocks(); vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue(readiness()); vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValue(lifecycle()); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

describe('mounted current-scope preparation evidence', () => {
  it('reports only validated current readonly evidence, actual lifecycle actions and safe typed reasons', async () => {
    const view = host();
    await waitFor(() => expect(view.current()?.readiness.state).toBe('available'));
    expect(view.current()?.readiness).toMatchObject({ value: { config_ready: true, schema_ready: false, ready: false } });
    expect(view.current()?.lifecycle).toEqual({ state: 'available', scope, value: {
      edit: { state: 'available', reasonCodes: [] }, stop: { state: 'blocked', reasonCodes: ['already_stopped'] } } });
    expect(JSON.stringify(view.current())).not.toMatch(/secret|private|history|requires_review|schema_digest/);
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledExactlyOnceWith('group-1');
    expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledExactlyOnceWith('group-1');
    expect(studioV2WorkspaceWriteGroupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.schemaApply).not.toHaveBeenCalled();
  });
  it('does not requery or rereport for identical scope object copies and disclosure rerenders', async () => {
    const view = host(); await waitFor(() => expect(view.current()?.lifecycle.state).toBe('available'));
    const reports = view.onEvidence.mock.calls.length;
    view.rerenderScope({ ...scope }); view.rerenderScope({ ...scope });
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(1); expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(1);
    expect(view.onEvidence).toHaveBeenCalledTimes(reports);
  });
  it('shares exactly one readiness and lifecycle fetch for simultaneous same-group consumers', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
    const first = vi.fn(), second = vi.fn();
    render(<QueryClientProvider client={client}>
      <RecordingPreparationEvidence scope={scope} group={savedGroup()} onEvidence={first} />
      <RecordingPreparationEvidence scope={{ ...scope }} group={savedGroup()} onEvidence={second} />
    </QueryClientProvider>);
    await waitFor(() => expect(first.mock.calls.at(-1)?.[0]?.lifecycle.state).toBe('available'));
    expect(second.mock.calls.at(-1)?.[0]?.readiness.state).toBe('available');
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(1); expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(1);
    const nextRead = deferred<WriteGroupReadiness>(), nextLifecycle = deferred<WriteGroupLifecycle>();
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockReturnValue(nextRead.promise);
    vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockReturnValue(nextLifecycle.promise);
    let pending: Promise<void>[] = [];
    act(() => { pending = [first.mock.calls.at(-1)![0].refresh(), second.mock.calls.at(-1)![0].refresh()]; });
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(2); expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(2);
    await act(async () => { nextRead.resolve(readiness()); nextLifecycle.resolve(lifecycle()); await Promise.all(pending); });
  });
  it('changes the viewed device of a shared group without requerying unchanged group evidence', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
    const group = savedGroup(); group.members.push({ ...group.members[0], device_id: 'dev-2', point_id: 'P2', tag_id: 'T2' });
    const onEvidence = vi.fn();
    const tree = (selected: RecordingWorkflowScope) => <QueryClientProvider client={client}>
      <RecordingPreparationEvidence scope={selected} group={group} onEvidence={onEvidence} />
    </QueryClientProvider>;
    const view = render(tree(scope));
    await waitFor(() => expect(onEvidence.mock.calls.at(-1)?.[0]?.readiness.state).toBe('available'));
    view.rerender(tree({ ...scope, deviceId: 'dev-2' }));
    await waitFor(() => expect(onEvidence.mock.calls.at(-1)?.[0]?.scope.deviceId).toBe('dev-2'));
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(1); expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(1);
  });
  it('rejects conflicting managed device ownership without querying another device', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
    const onEvidence = vi.fn();
    render(<QueryClientProvider client={client}><RecordingPreparationEvidence scope={scope}
      group={savedGroup({ basic_managed_device_id: 'foreign-device' })} onEvidence={onEvidence} /></QueryClientProvider>);
    await waitFor(() => expect(onEvidence.mock.calls.at(-1)?.[0]?.readiness.state).toBe('unavailable'));
    await act(async () => { await onEvidence.mock.calls.at(-1)?.[0]?.refresh(); });
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).not.toHaveBeenCalled(); expect(studioV2WriteGroupLifecycleAPI.lifecycle).not.toHaveBeenCalled();
  });
  it('does not infer or overwrite a duplicated lifecycle action', async () => {
    const value = lifecycle(); value.actions.push({ ...value.actions[0], allowed: false, reason_code: 'blocked' });
    vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValue(value);
    const view = host(); await waitFor(() => expect(view.current()?.readiness.state).toBe('available'));
    expect(view.current()?.lifecycle.state).toBe('unavailable');
  });
  it.each(['workspace_id', 'workspace_revision', 'group_id', 'group_revision', 'applied_revision'] as const)('withholds mismatched readiness %s without promoting lifecycle to readiness', async key => {
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue({ ...readiness(), [key]: 'foreign' });
    const view = host(); await waitFor(() => expect(view.current()?.lifecycle.state).toBe('available'));
    expect(view.current()?.readiness.state).toBe('unavailable');
  });
  it.each(['workspace_revision', 'group_revision', 'applied_revision', 'connector_id', 'connector_revision'] as const)('withholds mismatched lifecycle %s without borrowing availability', async key => {
    vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValue({ ...lifecycle(), [key]: 'foreign' });
    const view = host(); await waitFor(() => expect(view.current()?.readiness.state).toBe('available'));
    expect(view.current()?.lifecycle.state).toBe('unavailable');
  });
  it('keeps missing/error responses unavailable and never reports raw exceptions', async () => {
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockRejectedValue(new Error('password=secret /private/db'));
    vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockRejectedValue(new Error('dsn=secret'));
    const view = host(); await waitFor(() => expect(view.current()?.readiness.state).toBe('error'));
    expect(view.current()?.lifecycle.state).toBe('error'); expect(JSON.stringify(view.current())).not.toMatch(/secret|private/);
  });
  it('never queries invalid selected scope and reports unavailable rather than ready', async () => {
    const view = host({ ...scope, deviceId: 'foreign-device' });
    await waitFor(() => expect(view.current()?.readiness.state).toBe('unavailable'));
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).not.toHaveBeenCalled(); expect(studioV2WriteGroupLifecycleAPI.lifecycle).not.toHaveBeenCalled();
  });
  it('fences A→B→A stale refresh closures, revalidates changed revision once, and ignores late A evidence', async () => {
    const oldReadiness = deferred<WriteGroupReadiness>();
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockReturnValueOnce(oldReadiness.promise).mockResolvedValueOnce(readiness({ ...scope, workspaceRevision: 'W2', groupRevision: 'rev-2' })).mockResolvedValue(readiness());
    vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValueOnce(lifecycle()).mockResolvedValueOnce(lifecycle({ ...scope, workspaceRevision: 'W2', groupRevision: 'rev-2' })).mockResolvedValue(lifecycle());
    const view = host(); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(1));
    const firstRefresh = view.current()!.refresh;
    const next = { ...scope, workspaceRevision: 'W2', groupRevision: 'rev-2' };
    view.rerenderScope(next);
    expect(view.current()?.readiness.state).not.toBe('available');
    await act(async () => { oldReadiness.resolve(readiness()); });
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(2));
    expect(view.current()?.scope.groupRevision).toBe('rev-2');
    const current = view.current();
    if (current?.readiness.state === 'available') expect(current.readiness.value.group_revision).toBe('rev-2');
    view.rerenderScope(scope);
    await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(3));
    await act(async () => { await firstRefresh(); });
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(3);
    await waitFor(() => expect(view.current()?.readiness.state).toBe('available'));
    await act(async () => { await view.current()!.refresh(); });
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(4); expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenCalledTimes(4);
  });
  it('retires only the callback that owned a report and disables refresh after unmount', async () => {
    const view = host(); await waitFor(() => expect(view.current()?.readiness.state).toBe('available'));
    const refresh = view.current()!.refresh, old = view.onEvidence;
    const next = vi.fn<(value: RecordingPreparationEvidenceView | null) => void>();
    view.rerenderScope(scope, next);
    expect(old.mock.calls.at(-1)?.[0]).toBeNull(); expect(next.mock.calls.at(-1)?.[0]?.scope).toEqual(scope);
    view.unmount(); expect(next.mock.calls.at(-1)?.[0]).toBeNull();
    await act(async () => { await refresh(); });
    expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalledTimes(1); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  });
});
