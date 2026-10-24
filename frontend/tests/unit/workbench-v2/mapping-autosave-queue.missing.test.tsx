import { act, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { api, deferred, mountQueue, points, records, response } from './mapping-autosave-queue.harness';

vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() } }));
const missing = { response: { status: 404, data: { error: { code: 'workspace_mapping_not_found', action: 'reload' } } } };
beforeEach(() => vi.resetAllMocks());

it('reloads missing row identities before explicit retry and creates each latest queued draft', async () => {
  api.update.mockRejectedValue(missing);
  api.create.mockImplementation(async request => response(`new-${request.address}`, request));
  const hook = await mountQueue(points.slice(0, 2), records.slice(0, 2));
  hook.send({ type: 'setAllMappingsEnabled', enabled: false });
  await waitFor(() => expect(hook.result.current.state.mappings.a1.save_state).toBe('save-error'));
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { display_name: 'Latest name', scale: 2, offset: 3 } });
  api.list.mockResolvedValue([]);
  hook.send({ type: 'retryMappingSave', pointId: 'a1' });
  await waitFor(() => expect(hook.result.current.state.mappings.a1.save_state).toBe('saved'));
  expect(api.update).toHaveBeenCalledTimes(1);
  expect(api.create).toHaveBeenCalledTimes(2);
  expect(api.create).toHaveBeenCalledWith(expect.objectContaining({ address: '40001', display_name: 'Latest name', scale: 2, offset: 3, enabled: false }));
  expect(hook.result.current.state.mappings.a0).toMatchObject({ mapping_id: 'new-40001', display_name: 'Latest name', scale: 2, offset: 3 });
});

it('uses a replacement persisted identity while preserving edits made during reload', async () => {
  api.update.mockRejectedValueOnce(missing).mockImplementation(async (id, request) => response(id, request));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { display_name: 'Draft' } });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('save-error'));
  const reload = deferred<typeof records>();
  api.list.mockReturnValueOnce(reload.promise);
  hook.send({ type: 'retryMappingSave', pointId: 'a0' });
  await waitFor(() => expect(api.list).toHaveBeenCalledTimes(2));
  hook.send({ type: 'retryMappingSave', pointId: 'a0' });
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { display_name: 'Newest draft', unit: 'kW' } });
  await act(async () => reload.resolve([{ ...records[0], id: 'replacement', display_name: 'Server name' }]));
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
  expect(api.update).toHaveBeenLastCalledWith('replacement', expect.objectContaining({ display_name: 'Newest draft', unit: 'kW' }));
  expect(api.update).toHaveBeenCalledTimes(2);
  expect(api.create).not.toHaveBeenCalled();
});

it('keeps drafts paused if the authoritative reload fails and allows a later explicit retry', async () => {
  api.update.mockRejectedValue(missing);
  api.create.mockImplementation(async request => response('recovered', request));
  const hook = await mountQueue(points.slice(0, 1), records.slice(0, 1));
  hook.send({ type: 'updateMapping', pointId: 'a0', patch: { display_name: 'Keep draft' } });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('save-error'));
  api.list.mockRejectedValueOnce(new Error('reload unavailable'));
  hook.send({ type: 'retryMappingSave', pointId: 'a0' });
  await waitFor(() => expect(hook.result.current.query.isError).toBe(true));
  expect(api.update).toHaveBeenCalledTimes(1);
  expect(api.create).not.toHaveBeenCalled();
  expect(hook.result.current.state.mappings.a0).toMatchObject({ display_name: 'Keep draft', save_state: 'save-error' });
  api.list.mockResolvedValue([]);
  hook.send({ type: 'retryMappingSave', pointId: 'a0' });
  await waitFor(() => expect(hook.result.current.state.mappings.a0.save_state).toBe('saved'));
});
