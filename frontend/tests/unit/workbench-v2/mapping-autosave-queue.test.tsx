import { act, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api, deferred, mountQueue, points, records, response } from './mapping-autosave-queue.harness';
import type { StudioV2WorkspaceMappingRequest } from '../../../src/services/studioV2Mappings';
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
beforeEach(() => { vi.resetAllMocks(); });
describe('mapping rule save queue', () => {
  it('serializes eight same-rule rows, deduplicates repeated batches and lets another rule progress', async () => {
    const first = deferred<ReturnType<typeof response>>();
    api.update.mockImplementation((id: string, request: StudioV2WorkspaceMappingRequest) => id === 'mapping-a0' && api.update.mock.calls.length === 1 ? first.promise : Promise.resolve(response(id, request)));
    const hook = await mountQueue();
    hook.send({ type: 'setAllMappingsEnabled', enabled: false });
    hook.send({ type: 'setAllMappingsEnabled', enabled: true });
    await waitFor(() => { expect(api.update).toHaveBeenCalledWith('mapping-b0', expect.anything()); });
    expect(api.update.mock.calls.filter(([id]) => id !== 'mapping-b0')).toHaveLength(1);
    await act(async () => first.resolve(response('mapping-a0', api.update.mock.calls[0][1])));
    await waitFor(() => { expect(Object.values(hook.result.current.state.mappings).every((m) => m.save_state === 'saved')).toBe(true); });
    for (const point of points) expect(api.update.mock.calls.filter(([id]) => id === `mapping-${point.id}`).at(-1)?.[1]).toMatchObject({ rule_id: point.rule_id, address: point.address, enabled: true });
    for (const point of points.slice(1, 8)) expect(api.update.mock.calls.filter(([id]) => id === `mapping-${point.id}`)).toHaveLength(1);
  });
  it('keeps every newer column draft and continues with its exact request after an old response', async () => {
    const old = deferred<ReturnType<typeof response>>(), latest = deferred<ReturnType<typeof response>>();
    api.update.mockReturnValueOnce(old.promise).mockReturnValueOnce(latest.promise);
    const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'first' } });
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'last', display_name: 'Latest', unit: 'kW', target_type: 'float64', scale: 0.5, offset: 10, enabled: false } });
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    await act(async () => old.resolve(response('mapping-a0', api.update.mock.calls[0][1])));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(2));
    expect(hook.result.current.state.mappings.a0).toMatchObject({ tag_key: 'last', save_state: 'saving', persisted_value: { tag_key: 'first' } });
    expect(api.update.mock.calls[1][1]).toEqual({ rule_id: 'rule-A', address: '40001', tag_key: 'last', display_name: 'Latest', unit: 'kW', target_type: 'float64', scale: 0.5, offset: 10, enabled: false });
    await act(async () => latest.resolve(response('mapping-a0', api.update.mock.calls[1][1])));
    await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
    expect(hook.result.current.state.mappings.a0.persisted_value?.tag_key).toBe('last');
  });
  it('reads the newest transform for queued rows after repeated bulk actions', async () => {
    const first = deferred<ReturnType<typeof response>>();
    api.update.mockImplementation((id: string, request: StudioV2WorkspaceMappingRequest) => api.update.mock.calls.length === 1 ? first.promise : Promise.resolve(response(id, request)));
    const hook = await mountQueue(points.slice(0, 8), records.slice(0, 8));
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { scale: 2, offset: 1 } });
    hook.send({ type: 'bulkApplyTransform', fromPointId: 'a0', fields: ['scale', 'offset', 'target_type', 'unit'] });
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { scale: 0.5, offset: 10, target_type: 'float64', unit: 'C' } });
    hook.send({ type: 'bulkApplyTransform', fromPointId: 'a0', fields: ['scale', 'offset', 'target_type', 'unit'] });
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    await act(async () => first.resolve(response('mapping-a0', api.update.mock.calls[0][1])));
    await waitFor(() => expect(hook.result.current.state.mappings.a7.save_state).toBe('saved'));
    for (const point of points.slice(0, 8)) {
      expect(api.update.mock.calls.filter(([id]) => id === `mapping-${point.id}`).at(-1)?.[1]).toMatchObject({ scale: 0.5, offset: 10, target_type: 'float64', unit: 'C' });
      expect(hook.result.current.state.mappings[point.id]).toMatchObject({ scale: 0.5, offset: 10, target_type: 'float64', unit: 'C', save_state: 'saved' });
    }
  });
});

describe('mapping failure recovery', () => {
  it.each([[409, 'workspace_mapping_conflict', 'review_source'], [404, 'workspace_mapping_not_found', 'reload'], [500, 'workspace_mapping_save_failed', 'retry']])('retains safe metadata and pauses failed rule until explicit retry (%s)', async (status, code, action) => {
    const failing = deferred<ReturnType<typeof response>>();
    api.update.mockReturnValueOnce(failing.promise).mockImplementation((id: string, request: StudioV2WorkspaceMappingRequest) => Promise.resolve(response(id, request)));
    const hook = await mountQueue(points.slice(0, 2), records.slice(0, 2));
    hook.send({ type: 'setAllMappingsEnabled', enabled: false });
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'dirty.latest' } });
    await waitFor(() => expect(api.update).toHaveBeenCalled());
    await act(async () => failing.reject({ response: { status, data: { error: { code, action, request_id: 'req_safe_42', message: 'SQL password https://private' } } } }));
    await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('save-error'));
    expect(api.update).toHaveBeenCalledTimes(1);
    expect(hook.result.current.state.mappings.a0).toMatchObject({ tag_key: 'dirty.latest', save_error_detail: { status, code, action, requestId: 'req_safe_42' } });
    expect(hook.result.current.state.mappings.a1.save_state).toBe('save-error');
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { unit: 'new' } });
    expect(api.update).toHaveBeenCalledTimes(1);
    hook.send({ type: 'retryMappingSave', pointId: 'a0' });
    await waitFor(() => expect(hook.result.current.state.mappings.a1.save_state).toBe('saved'));
    expect(api.update.mock.calls.filter(([id]) => id === 'mapping-a0').at(-1)?.[1]).toMatchObject({ tag_key: 'dirty.latest', unit: 'new' });
  });
});


it('never inherits an automatic mutation retry after a mapping failure', async () => {
  api.update.mockRejectedValueOnce({ response: { status: 409, data: { error: { code: 'workspace_mapping_conflict', action: 'review_source' } } } }).mockImplementation((id, request) => Promise.resolve(response(id, request)));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1), 1);
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'dirty' } });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('save-error'));
  expect(api.update).toHaveBeenCalledTimes(1);
});

it('continues the latest draft entered during the final authoritative refetch', async () => {
  api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  const authoritative = deferred<typeof records>();
  api.list.mockReturnValueOnce(authoritative.promise);
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'first' } });
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(2));
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'during.refetch' } });
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(2));
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  await act(async () => authoritative.resolve([{ ...records[0], tag_key: 'first' }]));
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  expect(hook.result.current.state.mappings.a0.tag_key).toBe('during.refetch');
});

it('accepts a legitimate tag-only authoritative refetch with the same mapping timestamp', async () => {
  api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { display_name: 'saved name' } });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  await waitFor(() => expect(hook.result.current.query.isFetching).toBe(false));
  api.list.mockResolvedValueOnce([{ ...records[0], display_name: 'external tag name', unit: 'external unit' }]);
  await act(async () => { await hook.result.current.query.refetch(); });
  await waitFor(() => expect(hook.result.current.state.mappings.a0).toMatchObject({ display_name: 'external tag name', unit: 'external unit', save_state: 'saved' }));
});

it('keeps already confirmed row values visible while later batch invalidations return older snapshots', async () => {
  const last = deferred<ReturnType<typeof response>>();
  api.update.mockImplementation((id, request) => id === 'mapping-a2' ? last.promise : Promise.resolve(response(id, request)));
  const hook = await mountQueue(points.slice(0, 3), records.slice(0, 3));
  api.list.mockResolvedValue(records.slice(0, 3));
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { scale: 0.5, offset: 10, target_type: 'float64' } });
  hook.send({ type: 'bulkApplyTransform', fromPointId: 'a0', fields: ['scale', 'offset', 'target_type'] });
  await waitFor(() => expect(api.update).toHaveBeenCalledWith('mapping-a2', expect.anything()));
  for (const id of ['a0', 'a1']) expect(hook.result.current.state.mappings[id]).toMatchObject({ scale: 0.5, offset: 10, target_type: 'float64', save_state: 'saved' });
  api.list.mockResolvedValue(records.slice(0, 3).map((record) => ({ ...record, scale: 0.5, offset: 10, target_type: 'float64' as const })));
  const request = api.update.mock.calls.find(([id]) => id === 'mapping-a2')![1];
  await act(async () => last.resolve(response('mapping-a2', request)));
  await waitFor(() => expect(hook.result.current.query.isFetching).toBe(false));
  for (const id of ['a0', 'a1', 'a2']) expect(hook.result.current.state.mappings[id]).toMatchObject({ scale: 0.5, offset: 10, target_type: 'float64', save_state: 'saved' });
});

it('continues another rule while a global mapping GET is unresolved', async () => {
  const b1 = { ...points[8], id: 'b1', address: '40012' };
  const customPoints = [points[0], points[8], b1];
  const customRecords = [records[0], records[8], { ...records[8], id: 'mapping-b1', address: '40012' }];
  api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
  const hook = await mountQueue(customPoints, customRecords);
  const blockedGet = deferred<typeof records>(), finalGet = deferred<typeof records>();
  api.list.mockReturnValueOnce(blockedGet.promise).mockReturnValueOnce(finalGet.promise);
  hook.send({ type: 'updateMapping', pointId: 'b0', patch: { tag_key: 'b.first' } });
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(2));
  hook.send({ type: 'updateMapping', pointId: 'b1', patch: { tag_key: 'b.next' } });
  await waitFor(() => expect(api.update).toHaveBeenCalledWith('mapping-b1', expect.objectContaining({ tag_key: 'b.next' })));
  await waitFor(() => expect(hook.result.current.state.mappings.b1.save_state).toBe('saved'));
  await act(async () => blockedGet.resolve(customRecords));
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(3));
  expect(hook.result.current.state.mappings.b1).toMatchObject({ tag_key: 'b.next', save_state: 'saved' });
  const { studioV2WorkspaceKeys } = await import('../../../src/hooks/datalink/keys');
  expect(hook.client.getQueryData<typeof records>(studioV2WorkspaceKeys.mappings())?.find((record) => record.id === 'mapping-b1')?.tag_key).toBe('b.next');
  await act(async () => finalGet.resolve(customRecords.map((record) => ({ ...record, tag_key: record.id === 'mapping-b1' ? 'b.next' : record.id === 'mapping-b0' ? 'b.first' : record.tag_key }))));
  expect(hook.result.current.state.mappings.b1.tag_key).toBe('b.next');
});

it('retains a POST acknowledged after an older GET began until a later authoritative GET completes', async () => {
  api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
  api.create.mockImplementation((request) => Promise.resolve(response('new-B', request)));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  const oldGet = deferred<typeof records>(), finalGet = deferred<typeof records>();
  api.list.mockReturnValueOnce(oldGet.promise).mockReturnValueOnce(finalGet.promise);
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'a.saved' } });
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(2));
  hook.send({ type: 'initMappingsForPoints', points: [points[0], points[8]] });
  hook.send({ type: 'updateMapping', pointId: 'b0', patch: { tag_key: 'b.created' } });
  await waitFor(() => expect(hook.result.current.state.mappings.b0).toMatchObject({ mapping_id: 'new-B', tag_key: 'b.created', save_state: 'saved' }));
  expect(api.create).toHaveBeenCalledTimes(1);
  await act(async () => oldGet.resolve([{ ...records[0], tag_key: 'a.saved' }]));
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(3));
  const { studioV2WorkspaceKeys } = await import('../../../src/hooks/datalink/keys');
  const cached = () => hook.client.getQueryData<typeof records>(studioV2WorkspaceKeys.mappings());
  expect(hook.result.current.state.mappings.b0).toMatchObject({ mapping_id: 'new-B', tag_key: 'b.created', save_state: 'saved' });
  expect(cached()?.find((record) => record.id === 'new-B')).toMatchObject({ tag_key: 'b.created' });
  // A later authoritative absence must remain authoritative; acknowledgments
  // that predate that GET cannot keep an externally removed row forever.
  await act(async () => finalGet.resolve([{ ...records[0], tag_key: 'a.saved' }]));
  await waitFor(() => expect(hook.result.current.query.isFetching).toBe(false));
  expect(cached()?.find((record) => record.id === 'new-B')).toBeUndefined();
});
