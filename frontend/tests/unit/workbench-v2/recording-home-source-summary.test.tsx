import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, describe, expect, it, vi } from 'vitest';
import axios from 'axios';
import type { ReactNode } from 'react';
import { useRecordingHomeSources, invalidateRecordingHomeSources } from '@/hooks/datalink/useRecordingHomeSources';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import type { RecordingHomeSourceObservation } from '@/types/recordingHomeSources';
const snapshot = (workspace_id = 'A', name = 'Observed') => ({ workspace_id, state: 'available' as const, source_digest: 'a'.repeat(64),
  devices: [{ id: 'D', name, protocol: 'modbus', availability: 'available' }], candidates: [], excluded: [] });
const deferred = () => { let resolve!: (value: RecordingHomeSourceObservation) => void; let reject!: (error: Error) => void;
  const promise = new Promise<RecordingHomeSourceObservation>((yes, no) => { resolve = yes; reject = no; }); return { promise, resolve, reject }; };
const clients: QueryClient[] = [];
function setup(client = new QueryClient(), backend = axios.create({ baseURL: 'http://one/api/v1/datalink' })) {
  clients.push(client);
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const mount = (workspaceId = 'A', epoch = 1) => renderHook(({ id, visit }) => useRecordingHomeSources({ workspaceId: id, workspaceRevision: '', epoch: visit }, { backend }),
    { wrapper, initialProps: { id: workspaceId, visit: epoch } });
  return { client, backend, mount };
}
afterEach(() => { cleanup(); clients.forEach(client => client.clear()); clients.length = 0; vi.restoreAllMocks(); });
describe('observed source cache and current visit fences', () => {
  it('returns observation immediately but never fresh while refreshing or after failure', async () => {
    const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockResolvedValue(snapshot());
    const { mount } = setup(); const first = mount(); await waitFor(() => expect(first.result.current.fresh).toBe(true)); first.unmount();
    const pending = deferred(); read.mockReturnValue(pending.promise);
    const second = mount('A', 2);
    expect(second.result.current.observed?.devices[0].name).toBe('Observed'); expect(second.result.current.status).toBe('refreshing');
    expect(second.result.current.evidence.state).toBe('loading'); expect(second.result.current.fresh).toBe(false);
    await act(async () => { pending.reject(new Error('private.host password')); });
    await waitFor(() => expect(second.result.current.status).toBe('error'));
    expect(second.result.current.observed?.devices[0].name).toBe('Observed'); expect(second.result.current.evidence.state).toBe('error');
  });
  it('fences ignored-abort A -> B -> A completions and forwards aborted signals', async () => {
    const old = deferred(); const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockReturnValueOnce(old.promise).mockResolvedValueOnce(snapshot('B', 'B')).mockResolvedValueOnce(snapshot('A', 'new A'));
    const { mount } = setup(); const view = mount(); await waitFor(() => expect(read).toHaveBeenCalledTimes(1));
    const signal = read.mock.calls[0][1]!;
    view.rerender({ id: 'B', visit: 2 }); await waitFor(() => expect(view.result.current.fresh).toBe(true)); expect(signal.aborted).toBe(true);
    view.rerender({ id: 'A', visit: 3 }); await waitFor(() => expect(view.result.current.observed?.devices[0].name).toBe('new A'));
    await act(async () => { old.resolve(snapshot('A', 'old A')); });
    expect(view.result.current.observed?.devices[0].name).toBe('new A');
  });
  it('invalidation immediately removes fresh evidence and rejects pending pre-mutation reads', async () => {
    const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockResolvedValue(snapshot());
    const { mount, client, backend } = setup(); const view = mount(); await waitFor(() => expect(view.result.current.fresh).toBe(true));
    const pending = deferred(); read.mockReturnValueOnce(pending.promise);
    act(() => { void view.result.current.refetch(); });
    act(() => { invalidateRecordingHomeSources(client, { workspaceId: 'A', backend }); });
    expect(view.result.current.fresh).toBe(false); expect(view.result.current.observed?.devices[0].name).toBe('Observed');
    await act(async () => { pending.resolve(snapshot('A', 'obsolete')); });
    expect(view.result.current.fresh).toBe(false); expect(view.result.current.observed?.devices[0].name).toBe('Observed');
    await act(async () => { await view.result.current.refetch(); }); await waitFor(() => expect(view.result.current.fresh).toBe(true));
  });
  it('isolates QueryClient, backend object and current API base pathname', async () => {
    const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockResolvedValue(snapshot());
    const first = setup(); const view = first.mount(); await waitFor(() => expect(view.result.current.fresh).toBe(true)); view.unmount();
    read.mockReturnValue(new Promise(() => {}));
    const otherClient = setup(new QueryClient(), first.backend).mount(); expect(otherClient.result.current.observed).toBeUndefined();
    const otherBackend = setup(first.client).mount(); expect(otherBackend.result.current.observed).toBeUndefined();
    first.backend.defaults.baseURL = 'http://one/another/api/v1/datalink';
    const otherBase = first.mount('A', 2); expect(otherBase.result.current.observed).toBeUndefined();
  });
  it('does not reuse pending promises across backend clients with identical base URLs', async () => {
    const pending = deferred(), read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockReturnValue(pending.promise);
    const first = setup(), second = setup(first.client);
    const one = first.mount(), two = second.mount();
    await waitFor(() => expect(read).toHaveBeenCalledTimes(2));
    expect(read.mock.calls[0][2]).not.toBe(read.mock.calls[1][2]);
    one.unmount(); expect(read.mock.calls[0][1]?.aborted).toBe(true); expect(read.mock.calls[1][1]?.aborted).toBe(false);
    await act(async () => { pending.resolve(snapshot()); });
    await waitFor(() => expect(two.result.current.fresh).toBe(true));
  });
  it('keeps failed refresh safely retryable and only successful refresh replaces the digest', async () => {
    const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockResolvedValue(snapshot());
    const { mount } = setup(), view = mount(); await waitFor(() => expect(view.result.current.fresh).toBe(true));
    read.mockRejectedValueOnce(new Error('private.host'));
    await act(async () => { await view.result.current.refetch(); });
    await waitFor(() => expect(view.result.current.status).toBe('error')); expect(view.result.current.observed?.source_digest).toBe('a'.repeat(64));
    read.mockResolvedValueOnce({ ...snapshot('A', 'new'), source_digest: 'b'.repeat(64) });
    await act(async () => { await view.result.current.refetch(); });
    await waitFor(() => expect(view.result.current.fresh).toBe(true));
    expect(view.result.current.observed?.source_digest).toBe('b'.repeat(64));
  });
  it('retains latest only, evicts ninth scope and does not cache unknown as an empty observation', async () => {
    const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockImplementation(async id => snapshot(id));
    const { mount } = setup();
    for (let i = 0; i < 9; i++) { const view = mount(`W${i}`, i + 1); await waitFor(() => expect(view.result.current.fresh).toBe(true)); view.unmount(); }
    read.mockReturnValue(new Promise(() => {})); const evicted = mount('W0', 10); expect(evicted.result.current.observed).toBeUndefined(); evicted.unmount();
    const retained = mount('W8', 11); expect(retained.result.current.observed).toBeDefined(); retained.unmount();
    read.mockResolvedValue(snapshot('W8', 'latest')); const updated = mount('W8', 12); await waitFor(() => expect(updated.result.current.fresh).toBe(true)); updated.unmount();
    read.mockResolvedValue({ workspace_id: 'W8', state: 'unknown', reason_code: 'source-changed' });
    const unknown = mount('W8', 13); await waitFor(() => expect(unknown.result.current.status).toBe('stale'));
    expect(unknown.result.current.observed?.devices[0].name).toBe('latest'); expect(unknown.result.current.fresh).toBe(false);
  });
});
it('retains the captured API destination when a previous render refetches after base changes', async () => {
  const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockResolvedValue(snapshot());
  const { mount, backend } = setup(); const view = mount(); await waitFor(() => expect(view.result.current.fresh).toBe(true));
  const refetch = view.result.current.refetch;
  backend.defaults.baseURL = 'http://two/api/v1/datalink';
  await act(async () => { await refetch(); });
  expect(read.mock.calls[1][3]).toBe('http://one/api/v1/datalink');
});
