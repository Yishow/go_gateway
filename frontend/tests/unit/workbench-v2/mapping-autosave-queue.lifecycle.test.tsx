import { act, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { api, deferred, mountQueue, points, records, response } from './mapping-autosave-queue.harness';
import { buildDefaultMapping } from '../../../src/features/datalink/workbench-v2/state/mappingDefaults';
import { mappingFixture } from './mapping-autosave-page.testHarness';
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
beforeEach(() => vi.resetAllMocks());
describe('mapping queue lifecycle and defaults', () => {
  it('does not resume a removed pending row or overlap its deletion with a save', async () => {
    const save = deferred<ReturnType<typeof response>>();
    api.update.mockReturnValueOnce(save.promise);
    api.remove.mockResolvedValue({ runtime_apply_status: 'not_running' });
    const hook = await mountQueue(points.slice(0, 2), records.slice(0, 2));
    hook.send({ type: 'setAllMappingsEnabled', enabled: false });
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    hook.send({ type: 'initMappingsForPoints', points: [] });
    expect(api.remove).not.toHaveBeenCalled();
    await act(async () => save.resolve(response('mapping-a0', api.update.mock.calls[0][1])));
    await waitFor(() => expect(api.remove).toHaveBeenCalledTimes(2));
    expect(api.update).toHaveBeenCalledTimes(1);
    expect(hook.result.current.state.mappings).toEqual({});
  });
  it('does not apply a removed row response to a new row with the same identity', async () => {
    const old = deferred<ReturnType<typeof response>>();
    api.update.mockReturnValueOnce(old.promise);
    api.remove.mockResolvedValue({ runtime_apply_status: 'not_running' });
    const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'obsolete' } });
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    hook.send({ type: 'initMappingsForPoints', points: [] });
    hook.send({ type: 'initMappingsForPoints', points: points.slice(0, 1) });
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: '', display_name: 'new draft' } });
    await act(async () => old.resolve(response('mapping-a0', api.update.mock.calls[0][1])));
    expect(hook.result.current.state.mappings.a0).toMatchObject({ tag_key: '', display_name: 'new draft', persisted: false });
  });
  it('shows an actionable failed deletion and removes the restored row after explicit retry', async () => {
    api.remove.mockRejectedValueOnce({ response: { status: 500, data: { error: { code: 'workspace_mapping_save_failed', action: 'retry' } } } }).mockResolvedValue({ runtime_apply_status: 'not_running' });
    const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
    hook.send({ type: 'initMappingsForPoints', points: [] });
    await waitFor(() => expect(hook.result.current.state.mappings.a0?.save_state).toBe('save-error'));
    expect(api.remove).toHaveBeenCalledTimes(1);
    hook.send({ type: 'initMappingsForPoints', points: [] });
    expect(hook.result.current.state.mappings.a0?.save_state).toBe('save-error');
    hook.send({ type: 'retryMappingSave', pointId: 'a0' });
    await waitFor(() => expect(api.remove).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(hook.result.current.state.mappings.a0).toBeUndefined());
  });
  it('marks a newly edited row in a paused rule actionable and never automatically retries', async () => {
    api.update.mockRejectedValue({ response: { status: 500, data: { error: { code: 'workspace_mapping_save_failed', action: 'retry' } } } });
    const hook = await mountQueue(points.slice(0, 2), records.slice(0, 2));
    hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'failed' } });
    await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('save-error'));
    hook.send({ type: 'updateMapping', pointId: 'a1', patch: { tag_key: 'also.dirty' } });
    expect(hook.result.current.state.mappings.a1).toMatchObject({ tag_key: 'also.dirty', save_state: 'save-error' });
    expect(api.update).toHaveBeenCalledTimes(1);
  });
  it('stops continuation after unmount', async () => {
    const save = deferred<ReturnType<typeof response>>();
    api.update.mockReturnValueOnce(save.promise);
    const hook = await mountQueue(points.slice(0, 2), records.slice(0, 2));
    hook.send({ type: 'setAllMappingsEnabled', enabled: false });
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    hook.unmount();
    await act(async () => save.resolve(response('mapping-a0', api.update.mock.calls[0][1])));
    expect(api.update).toHaveBeenCalledTimes(1);
  });
  it('keeps default uint64 and scaled int16 request types exact through the queue', async () => {
    const numericPoints = [{ ...points[0], data_type: 'uint64' }, { ...points[1], _rule_scale: 0.5, _rule_offset: 10 }];
    const numericRecords = numericPoints.map((point) => mappingFixture({ id: `mapping-${point.id}`, rule_id: point.rule_id, address: point.address, ...buildDefaultMapping(point, 0).local_value }));
    api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
    const hook = await mountQueue(numericPoints, numericRecords);
    hook.send({ type: 'setAllMappingsEnabled', enabled: true });
    await waitFor(() => expect(hook.result.current.state.mappings.a1.save_state).toBe('saved'));
    expect(api.update.mock.calls[0][1]).toMatchObject({ target_type: 'uint64', scale: 1, offset: 0 });
    expect(api.update.mock.calls[1][1]).toMatchObject({ target_type: 'float64', scale: 0.5, offset: 10 });
  });
});

it('removes a successful deletion and retains a safe warning when unused-tag cleanup fails', async () => {
  const result = { runtime_apply_status: 'not_running' as const, cleanup_status: 'failed' as const, cleanup_message: 'SQL password https://private' };
  api.remove.mockResolvedValueOnce(result).mockResolvedValue({ runtime_apply_status: 'not_running' });
  const hook = await mountQueue(points.slice(0, 2), records.slice(0, 2));
  hook.send({ type: 'initMappingsForPoints', points: [] });
  await waitFor(() => expect(hook.result.current.state.mapping_cleanup_incomplete).toBe(true));
  expect(hook.result.current.state.mappings).toEqual({});
  await waitFor(() => expect(api.remove).toHaveBeenCalledTimes(2));
  expect(hook.result.current.state.mapping_cleanup_incomplete).toBe(true);
  const { Step3Mapping } = await import('../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping');
  const { QueryClient, QueryClientProvider } = await import('@tanstack/react-query');
  render(<QueryClientProvider client={new QueryClient()}><Step3Mapping state={{ ...hook.result.current.state, devices: [], rules: [] }} dispatch={vi.fn()} onContinue={vi.fn()} onBack={vi.fn()} /></QueryClientProvider>);
  expect(screen.getByTestId('mapping-cleanup-warning')).toHaveTextContent('The mapping change completed');
  expect(screen.getByTestId('mapping-cleanup-warning')).not.toHaveTextContent(/SQL|password|https:\/\/private/);
});


it('does not interpret an unknown cleanup status as a failure or claim cleanup completed', async () => {
  const result = { runtime_apply_status: 'not_running' as const, cleanup_status: 'unknown backend details' } as unknown as Awaited<ReturnType<typeof api.remove>>;
  api.remove.mockResolvedValue(result);
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  hook.send({ type: 'initMappingsForPoints', points: [] });
  await waitFor(() => expect(api.remove).toHaveBeenCalledTimes(1));
  expect(hook.result.current.state.mappings).toEqual({});
  expect(hook.result.current.state.mapping_cleanup_incomplete).toBeUndefined();
});

it('does not apply an old rule rejection to a replacement row and still pauses its own pending rows', async () => {
  const old = deferred<ReturnType<typeof response>>();
  api.update.mockReturnValueOnce(old.promise);
  api.create.mockImplementation((request) => Promise.resolve(response('new-B', request)));
  api.remove.mockResolvedValue({ runtime_apply_status: 'not_running' });
  const hook = await mountQueue(points.slice(0, 2), records.slice(0, 2));
  hook.send({ type: 'setAllMappingsEnabled', enabled: false });
  await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
  const replacement = { ...points[8], id: 'a0' };
  hook.send({ type: 'initMappingsForPoints', points: [replacement, points[1]] });
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'new.b' } });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  await act(async () => old.reject({ response: { status: 409, data: { error: { code: 'workspace_mapping_conflict' } } } }));
  expect(hook.result.current.state.mappings.a0).toMatchObject({ tag_key: 'new.b', save_state: 'saved', rule_id: 'rule-B' });
  expect(hook.result.current.state.mappings.a1.save_state).toBe('save-error');
});

it('keeps an old deletion failure actionable without changing its replacement row', async () => {
  const removal = deferred<Awaited<ReturnType<typeof api.remove>>>();
  api.remove.mockReturnValueOnce(removal.promise).mockResolvedValue({ runtime_apply_status: 'not_running' });
  api.create.mockImplementation((request) => Promise.resolve(response('new-B', request)));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  hook.send({ type: 'initMappingsForPoints', points: [] });
  await waitFor(() => expect(api.remove).toHaveBeenCalledTimes(1));
  hook.send({ type: 'initMappingsForPoints', points: [{ ...points[8], id: 'a0' }] });
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'new.b' } });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  await act(async () => removal.reject({ response: { status: 500, data: { error: { code: 'workspace_mapping_save_failed' } } } }));
  expect(hook.result.current.state.mappings.a0).toMatchObject({ tag_key: 'new.b', save_state: 'saved', rule_id: 'rule-B' });
  const failedRemoval = Object.values(hook.result.current.state.mappings).find((mapping) => mapping.save_error_detail?.operation === 'delete');
  expect(failedRemoval?.point_id).not.toBe('a0');
  expect(failedRemoval?.save_state).toBe('save-error');
  hook.send({ type: 'retryMappingSave', pointId: failedRemoval!.point_id });
  await waitFor(() => expect(api.remove).toHaveBeenCalledTimes(2));
  expect(hook.result.current.state.mappings.a0.save_state).toBe('saved');
});


it('acknowledges an accepted save with failed tag cleanup and shows a cleanup warning without retry', async () => {
  api.update.mockImplementation((id, request) => Promise.resolve({ ...response(id, request), cleanup_status: 'failed' as const, cleanup_message: 'SQL password https://private' }));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { tag_key: 'accepted.new' } });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  expect(hook.result.current.state.mapping_cleanup_incomplete).toBe(true);
  expect(hook.result.current.state.mappings.a0.save_error).toBeNull();
  expect(api.update).toHaveBeenCalledTimes(1);
  const { Step3Mapping } = await import('../../../src/features/datalink/workbench-v2/steps/step3/Step3Mapping');
  const { QueryClient, QueryClientProvider } = await import('@tanstack/react-query');
  render(<QueryClientProvider client={new QueryClient()}><Step3Mapping state={{ ...hook.result.current.state, devices: [], rules: [] }} dispatch={vi.fn()} onContinue={vi.fn()} onBack={vi.fn()} /></QueryClientProvider>);
  expect(screen.getByTestId('mapping-cleanup-warning')).toHaveTextContent('The mapping change completed');
  expect(screen.getByTestId('mapping-cleanup-warning')).not.toHaveTextContent(/SQL|password|https:\/\/private/);
  expect(screen.queryByTestId('mapping-save-retry-a0')).not.toBeInTheDocument();
});

it('does not restore a successfully deleted mapping from an older global GET', async () => {
  const customPoints = [points[0], points[8]], customRecords = [records[0], records[8]];
  api.update.mockImplementation((id, request) => Promise.resolve(response(id, request)));
  api.remove.mockResolvedValue({ runtime_apply_status: 'not_running' });
  const hook = await mountQueue(customPoints, customRecords);
  const oldGet = deferred<typeof records>(), finalGet = deferred<typeof records>();
  api.list.mockReturnValueOnce(oldGet.promise).mockReturnValueOnce(finalGet.promise);
  hook.send({ type: 'updateMapping', pointId: 'b0', patch: { tag_key: 'b.saved' } });
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(2));
  hook.send({ type: 'initMappingsForPoints', points: [points[8]] });
  await waitFor(() => expect(api.remove).toHaveBeenCalledWith('mapping-a0'));
  await act(async () => oldGet.resolve(customRecords));
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(3));
  const { studioV2WorkspaceKeys } = await import('../../../src/hooks/datalink/keys');
  expect(hook.client.getQueryData<typeof records>(studioV2WorkspaceKeys.mappings())?.find((record) => record.id === 'mapping-a0')).toBeUndefined();
  expect(hook.result.current.state.mappings.a0).toBeUndefined();
  await act(async () => finalGet.resolve([{ ...records[8], tag_key: 'b.saved' }]));
});
