import type * as React from 'react';
import { QueryClientProvider, type QueryClient } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { api, deferred, mountQueue, points, records, response } from './mapping-autosave-queue.harness';
import { workspaceFixture } from './mapping-autosave-page.testHarness';
import { useStudioV2WorkspaceQuery } from '../../../src/hooks/datalink/useStudioV2Workspace';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';

vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
const bootstrapAPI = vi.mocked(studioV2WorkspaceAPI);
const metadata = (ready: boolean) => workspaceFixture({ readiness_summary: { ready, blocking_count: ready ? 0 : 1, warning_count: 0, issues: [] } });
async function observeBootstrap(client: QueryClient) {
  const wrapper = ({ children }: { children: React.ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const observer = renderHook(() => ({ ...useStudioV2WorkspaceQuery() }), { wrapper });
  await waitFor(() => expect(observer.result.current.isSuccess).toBe(true));
  return observer;
}
beforeEach(() => {
  vi.resetAllMocks();
  bootstrapAPI.get.mockResolvedValue(metadata(false));
  api.create.mockImplementation((request) => Promise.resolve(response('created-B', request)));
  api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
  api.remove.mockResolvedValue({ runtime_apply_status: 'not_running' });
});

it.each(['create', 'save', 'delete'])('refreshes active workspace readiness after accepted mapping %s', async (operation) => {
  const queue = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  const bootstrap = await observeBootstrap(queue.client);
  bootstrapAPI.get.mockResolvedValueOnce(metadata(true));
  if (operation === 'create') {
    queue.send({ type: 'initMappingsForPoints', points: [points[0], points[8]] });
    queue.send({ type: 'updateMapping', pointId: 'b0', patch: { tag_key: 'b.created' } });
  } else {
    queue.send(operation === 'save' ? { type: 'updateMapping', pointId: 'a0', patch: { unit: 'updated' } } : { type: 'initMappingsForPoints', points: [] });
  }
  await waitFor(() => expect(bootstrap.result.current.data?.readiness_summary?.ready).toBe(true));
  expect(bootstrapAPI.get).toHaveBeenCalledTimes(2);
  if (operation === 'create') expect(queue.result.current.state.mappings.b0).toMatchObject({ mapping_id: 'created-B', save_state: 'saved' });
  else if (operation === 'save') expect(queue.result.current.state.mappings.a0.save_state).toBe('saved');
  else expect(queue.result.current.state.mappings).toEqual({});
});

it('continues rule B while bootstrap GETs are blocked and ignores older metadata after its latest acknowledgment', async () => {
  const b1 = { ...points[8], id: 'b1', address: '40012' };
  const queue = await mountQueue([points[0], points[8], b1], [records[0], records[8], { ...records[8], id: 'mapping-b1', address: '40012' }]);
  const bootstrap = await observeBootstrap(queue.client);
  const reads: ReturnType<typeof deferred<ReturnType<typeof metadata>>>[] = [];
  bootstrapAPI.get.mockImplementation(() => { const read = deferred<ReturnType<typeof metadata>>(); reads.push(read); return read.promise; });
  queue.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'a.saved' } });
  await waitFor(() => expect(reads).toHaveLength(1));
  queue.send({ type: 'updateMapping', pointId: 'b0', patch: { tag_key: 'b.first' } });
  queue.send({ type: 'updateMapping', pointId: 'b1', patch: { tag_key: 'b.next' } });
  await waitFor(() => expect(queue.result.current.state.mappings.b1).toMatchObject({ tag_key: 'b.next', save_state: 'saved' }));
  expect(api.update).toHaveBeenCalledWith('mapping-b1', expect.objectContaining({ tag_key: 'b.next' }));
  await waitFor(() => expect(reads).toHaveLength(3));
  await act(async () => reads[2].resolve(metadata(true)));
  await waitFor(() => expect(bootstrap.result.current.data?.readiness_summary?.ready).toBe(true));
  await act(async () => { reads[0].resolve(metadata(false)); reads[1].resolve(metadata(false)); });
  expect(bootstrap.result.current.data?.readiness_summary?.ready).toBe(true);
});

it('starts a new authoritative bootstrap read when a pre-save read is already in flight', async () => {
  const queue = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  const bootstrap = await observeBootstrap(queue.client);
  const oldRead = deferred<ReturnType<typeof metadata>>(), authoritative = deferred<ReturnType<typeof metadata>>();
  bootstrapAPI.get.mockReturnValueOnce(oldRead.promise).mockReturnValueOnce(authoritative.promise);
  act(() => { void bootstrap.result.current.refetch({ cancelRefetch: false }); });
  await waitFor(() => expect(bootstrapAPI.get).toHaveBeenCalledTimes(2));
  queue.send({ type: 'updateMapping', pointId: 'a0', patch: { unit: 'latest' } });
  await waitFor(() => expect(bootstrapAPI.get).toHaveBeenCalledTimes(3));
  await waitFor(() => expect(queue.result.current.state.mappings.a0.save_state).toBe('saved'));
  await act(async () => authoritative.resolve(metadata(true)));
  await waitFor(() => expect(bootstrap.result.current.data?.readiness_summary?.ready).toBe(true));
  await act(async () => oldRead.resolve(metadata(false)));
  expect(bootstrap.result.current.data?.readiness_summary?.ready).toBe(true);
});

it('keeps an accepted save successful when bootstrap refresh fails and refreshes on the next acknowledgment', async () => {
  const queue = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  const bootstrap = await observeBootstrap(queue.client);
  bootstrapAPI.get.mockRejectedValueOnce(new Error('owned metadata read failure'));
  queue.send({ type: 'updateMapping', pointId: 'a0', patch: { unit: 'first' } });
  await waitFor(() => expect(bootstrap.result.current.isError).toBe(true));
  expect(queue.result.current.state.mappings.a0).toMatchObject({ unit: 'first', save_state: 'saved', save_error: null });
  expect(api.update).toHaveBeenCalledTimes(1);
  expect(bootstrapAPI.get).toHaveBeenCalledTimes(2);
  bootstrapAPI.get.mockResolvedValueOnce(metadata(true));
  queue.send({ type: 'updateMapping', pointId: 'a0', patch: { unit: 'second' } });
  await waitFor(() => expect(bootstrap.result.current.data?.readiness_summary?.ready).toBe(true));
  expect(queue.result.current.state.mappings.a0).toMatchObject({ unit: 'second', save_state: 'saved', save_error: null });
});
