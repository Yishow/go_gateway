import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { RecordingPreparationTestWriteQuery, type RecordingPreparationTestWriteQueryReport } from '@/features/datalink/workbench-v2/steps/step4/RecordingPreparationTestWriteQuery';
import { recordingTestWriteRecovery as store, recordingTestWriteRecoveryKey, type RecordingTestWriteReference } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { WriteGroupTestWriteOperation } from '@/types/studioV2WriteGroupTestWrite';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const clients: QueryClient[] = [];
const original: RecordingTestWriteReference = { workspace_id: 'W1', group_id: 'G1', group_revision: 'R1', operation_id: 'O1', operation_kind: 'test_write', created_at: '2026-10-28T00:00:00Z' };
function completed(overrides: Partial<WriteGroupTestWriteOperation> = {}): WriteGroupTestWriteOperation {
  return { operation_id: 'O1', action: 'test_write', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned',
    created_at: original.created_at, updated_at: '2026-10-28T00:00:02Z', completed_at: '2026-10-28T00:00:01Z', ...overrides };
}
const response = (operation = completed()) => ({ data: { success: true, data: operation } });
function deferred() { let resolve!: (value: ReturnType<typeof response>) => void; const promise = new Promise<ReturnType<typeof response>>(done => { resolve = done; }); return { promise, resolve }; }
function resetStore() {
  localStorage.clear();
  for (const reference of store.read().references) store.resolve(completed({ operation_id: reference.operation_id }), reference.workspace_id);
}
function host(reference = original, ownerIdentity = 'W1/A/G1/R2') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const onQuery = vi.fn<(report: RecordingPreparationTestWriteQueryReport | null) => void>();
  const tree = (ref: RecordingTestWriteReference, identity: string) => <QueryClientProvider client={client}>
    <RecordingPreparationTestWriteQuery reference={ref} ownerIdentity={identity} onQuery={onQuery} />
  </QueryClientProvider>;
  const view = render(tree(reference, ownerIdentity));
  return { ...view, onQuery, current: () => onQuery.mock.calls.at(-1)?.[0], rerenderOwner: (identity: string, ref = reference) => view.rerender(tree(ref, identity)) };
}
beforeEach(() => { vi.restoreAllMocks(); vi.clearAllMocks(); resetStore(); store.retain(original); vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response()); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); });

describe('explicit readonly preparation test-write query delegate', () => {
  it('does not fetch on mount and explicitly queries original O1/R1 while viewing saved R2', async () => {
    const view = host(); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
    expect(Object.keys(view.current()!)).toEqual(['operationId', 'check']);
    await act(async () => { await view.current()!.check(); });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledExactlyOnceWith('/studio-v2/workspace/database-operations/O1');
    expect(await screen.findByTestId('preparation-test-write-result')).toHaveTextContent('O1');
    expect(store.read().references).toEqual([]);
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
  it.each(['pending', 'running', 'unknown', 'cleanup failed'] as const)('retains the original reference for %s without resending', async kind => {
    const operation: WriteGroupTestWriteOperation = kind === 'pending' || kind === 'running'
      ? { operation_id: 'O1', action: 'test_write', status: kind, created_at: original.created_at, updated_at: original.created_at }
      : kind === 'unknown' ? completed({ status: 'unknown', write_outcome: 'unknown', cleanup_status: 'unknown' }) : completed({ cleanup_status: 'failed' });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(operation));
    const view = host(); await act(async () => { await view.current()!.check(); });
    expect(await screen.findByTestId('preparation-test-write-result')).toBeInTheDocument();
    expect(store.read().references).toEqual([original]); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('allows query of a known unreadable reference but never resolves or overwrites it', async () => {
    const resolve = vi.spyOn(store, 'resolve');
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private storage error'); });
    const setItem = vi.spyOn(Storage.prototype, 'setItem');
    const view = host(); await act(async () => { await view.current()!.check(); });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(await screen.findByTestId('preparation-test-write-result')).toBeInTheDocument();
    expect(resolve).not.toHaveBeenCalled(); expect(setItem).not.toHaveBeenCalled();
    expect(store.read()).toMatchObject({ unreadable: true, references: [original] });
    expect(document.body).toHaveTextContent('step4.group.test_write.recovery_storage_error');
    expect(document.body).not.toHaveTextContent('private storage error');
  });
  it('does not query a supplied locator that is not the exact retained original reference', async () => {
    const view = host({ ...original, group_revision: 'R2' });
    await act(async () => { await view.current()!.check(); });
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled(); expect(store.read().references).toEqual([original]);
  });
  it('preserves a replaced reference and rejects a result if the full original locator changed during the read', async () => {
    const pending = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValue(pending.promise);
    const resolve = vi.spyOn(store, 'resolve'); const view = host();
    let check!: Promise<void>; act(() => { check = view.current()!.check(); });
    const replacement = { ...original, created_at: '2026-10-29T00:00:00Z' };
    localStorage.setItem(recordingTestWriteRecoveryKey, JSON.stringify([replacement]));
    await act(async () => { pending.resolve(response()); await check; });
    expect(resolve).not.toHaveBeenCalled(); expect(store.read().references).toEqual([replacement]);
    expect(screen.queryByTestId('preparation-test-write-result')).not.toBeInTheDocument();
  });
  it('rejects foreign result IDs and malformed operation envelopes without forgetting O1 or exposing raw errors', async () => {
    const view = host(); vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(completed({ operation_id: 'FOREIGN' })))
      .mockRejectedValueOnce(new Error('password=secret /private/endpoint'));
    await act(async () => { await view.current()!.check(); });
    expect(screen.queryByTestId('preparation-test-write-result')).not.toBeInTheDocument();
    await act(async () => { await view.current()!.check(); });
    expect(store.read().references).toEqual([original]); expect(document.body).not.toHaveTextContent(/FOREIGN|password|secret|private/);
    expect(screen.getByRole('alert')).toHaveTextContent('step4.group.test_write.unconfirmed');
  });
  it('uses the existing parser to reject contradictory cleanup proof without a fabricated result', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue(response(completed({ cleanup_status: 'not_attempted' })));
    const view = host(); await act(async () => { await view.current()!.check(); });
    expect(store.read().references).toEqual([original]); expect(screen.queryByTestId('preparation-test-write-result')).not.toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveTextContent('step4.group.test_write.unconfirmed');
  });
  it('keeps a known reference when storage becomes unreadable during an explicit read', async () => {
    const pending = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValue(pending.promise);
    const resolve = vi.spyOn(store, 'resolve'); const view = host();
    let check!: Promise<void>; act(() => { check = view.current()!.check(); });
    localStorage.setItem(recordingTestWriteRecoveryKey, '{broken');
    await act(async () => { pending.resolve(response()); await check; });
    expect(resolve).not.toHaveBeenCalled(); expect(store.read()).toMatchObject({ unreadable: true, references: [original] });
    expect(localStorage.getItem(recordingTestWriteRecoveryKey)).toBe('{broken');
    expect(screen.getByTestId('preparation-test-write-result')).toBeInTheDocument();
  });
  it('copies the original locator and cannot be retargeted by later caller edits or identical rerenders', async () => {
    const ref = { ...original }; const view = host(ref); const check = view.current()!.check;
    const reports = view.onQuery.mock.calls.length; view.rerenderOwner('W1/A/G1/R2', { ...ref });
    expect(view.onQuery).toHaveBeenCalledTimes(reports); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
    ref.operation_id = 'FOREIGN'; ref.group_revision = 'R9';
    await act(async () => { await check(); });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledExactlyOnceWith('/studio-v2/workspace/database-operations/O1');
    expect(screen.getByTestId('preparation-test-write-result')).toHaveTextContent('O1');
  });
  it('fences stale A→B→A reports and late completions, then permits a new explicit original query', async () => {
    const pending = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(pending.promise).mockResolvedValue(response());
    const resolve = vi.spyOn(store, 'resolve'); const view = host(); const oldCheck = view.current()!.check;
    let check!: Promise<void>; act(() => { check = oldCheck(); });
    view.rerenderOwner('W1/B/G1/R2'); view.rerenderOwner('W1/A/G1/R2');
    await act(async () => { pending.resolve(response()); await check; await oldCheck(); });
    expect(resolve).not.toHaveBeenCalled(); expect(store.read().references).toEqual([original]);
    expect(screen.queryByTestId('preparation-test-write-result')).not.toBeInTheDocument(); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    await act(async () => { await view.current()!.check(); });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2); expect(store.read().references).toEqual([]);
  });
  it('deduplicates repeated clicks and multiple same-operation consumers through the existing query cache', async () => {
    const pending = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValue(pending.promise);
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
    const first = vi.fn<(report: RecordingPreparationTestWriteQueryReport | null) => void>(), second = vi.fn<(report: RecordingPreparationTestWriteQueryReport | null) => void>();
    render(<QueryClientProvider client={client}>
      <RecordingPreparationTestWriteQuery reference={original} ownerIdentity="A" onQuery={first} />
      <RecordingPreparationTestWriteQuery reference={{ ...original }} ownerIdentity="A" onQuery={second} />
    </QueryClientProvider>);
    let calls: Promise<void>[] = [];
    act(() => { calls = [first.mock.calls.at(-1)![0]!.check(), first.mock.calls.at(-1)![0]!.check(), second.mock.calls.at(-1)![0]!.check()]; });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    await act(async () => { pending.resolve(response()); await Promise.all(calls); });
    expect(store.read().references).toEqual([]); expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
  it('clears only its callback and does not resolve an operation after unmount', async () => {
    const pending = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValue(pending.promise);
    const resolve = vi.spyOn(store, 'resolve'); const view = host(); const oldCheck = view.current()!.check;
    let check!: Promise<void>; act(() => { check = oldCheck(); }); view.unmount();
    expect(view.onQuery.mock.calls.at(-1)?.[0]).toBeNull();
    await act(async () => { pending.resolve(response()); await check; await oldCheck(); });
    expect(resolve).not.toHaveBeenCalled(); expect(store.read().references).toEqual([original]); expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
});
