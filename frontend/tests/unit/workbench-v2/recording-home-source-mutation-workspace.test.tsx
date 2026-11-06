import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import type { ReactNode } from 'react';
import { afterEach, expect, it, vi } from 'vitest';
import { useRecordingHomeSources, recordingHomeSourceContext, recordingHomeSourcePrefix } from '@/hooks/datalink/useRecordingHomeSources';
import * as devices from '@/hooks/datalink/useStudioV2WorkspaceDevices';
import * as rules from '@/hooks/datalink/useStudioV2Rules';
import * as mappings from '@/hooks/datalink/useStudioV2Mappings';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
const originalBase = studioV2DatalinkApi.defaults.baseURL;
const clients: QueryClient[] = [];
const summary = (workspace_id: string) => ({ workspace_id, state: 'available' as const, source_digest: 'a'.repeat(64), devices: [], candidates: [], excluded: [] });
function setup() {
  const client = new QueryClient(); clients.push(client); client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { id: 'A' });
  const read = vi.spyOn(recordingHomeSourcesAPI, 'read').mockImplementation(async id => summary(id));
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  const source = renderHook(() => useRecordingHomeSources({ workspaceId: 'A', workspaceRevision: '', epoch: 1 }), { wrapper });
  return { client, wrapper, source, read };
}
afterEach(() => { cleanup(); clients.forEach(client => client.clear()); clients.length = 0; studioV2DatalinkApi.defaults.baseURL = originalBase; vi.restoreAllMocks(); });
const operations: { name: string; api: () => unknown; useWrite: () => () => Promise<unknown> }[] = [
  { name: 'device create', api: () => vi.spyOn(studioV2WorkspaceDevicesAPI, 'create').mockResolvedValue({} as never), useWrite: () => { const mutation = devices.useCreateStudioV2WorkspaceDeviceMutation(); return () => mutation.mutateAsync({} as never); } },
  { name: 'device update', api: () => vi.spyOn(studioV2WorkspaceDevicesAPI, 'update').mockResolvedValue({} as never), useWrite: () => { const mutation = devices.useUpdateStudioV2WorkspaceDeviceMutation(); return () => mutation.mutateAsync({ deviceId: 'D', request: {} }); } },
  { name: 'device remove', api: () => vi.spyOn(studioV2WorkspaceDevicesAPI, 'remove').mockResolvedValue(), useWrite: () => { const mutation = devices.useDeleteStudioV2WorkspaceDeviceMutation(); return () => mutation.mutateAsync('D'); } },
  { name: 'device availability', api: () => vi.spyOn(studioV2WorkspaceDevicesAPI, 'updateAvailability').mockResolvedValue({} as never), useWrite: () => { const mutation = devices.useUpdateStudioV2WorkspaceDeviceAvailabilityMutation(); return () => mutation.mutateAsync({ deviceId: 'D', request: {} as never }); } },
  { name: 'device order', api: () => vi.spyOn(studioV2WorkspaceDevicesAPI, 'updateOrder').mockResolvedValue({} as never), useWrite: () => { const mutation = devices.useUpdateStudioV2WorkspaceDeviceOrderMutation(); return () => mutation.mutateAsync(['D']); } },
  { name: 'rule create', api: () => vi.spyOn(studioV2RulesAPI, 'create').mockResolvedValue({} as never), useWrite: () => { const mutation = rules.useCreateStudioV2RuleMutation(); return () => mutation.mutateAsync({} as never); } },
  { name: 'rule update', api: () => vi.spyOn(studioV2RulesAPI, 'update').mockResolvedValue({} as never), useWrite: () => { const mutation = rules.useUpdateStudioV2RuleMutation(); return () => mutation.mutateAsync({ ruleId: 'R', request: {} }); } },
  { name: 'rule remove', api: () => vi.spyOn(studioV2RulesAPI, 'remove').mockResolvedValue(), useWrite: () => { const mutation = rules.useDeleteStudioV2RuleMutation(); return () => mutation.mutateAsync('R'); } },
  { name: 'mapping create', api: () => vi.spyOn(studioV2MappingsAPI, 'create').mockResolvedValue({} as never), useWrite: () => { const mutation = mappings.useCreateStudioV2MappingMutation(); return () => mutation.mutateAsync({} as never); } },
  { name: 'mapping update', api: () => vi.spyOn(studioV2MappingsAPI, 'update').mockResolvedValue({} as never), useWrite: () => { const mutation = mappings.useUpdateStudioV2MappingMutation(); return () => mutation.mutateAsync({ mappingId: 'M', request: {} as never }); } },
  { name: 'mapping remove', api: () => vi.spyOn(studioV2MappingsAPI, 'remove').mockResolvedValue({} as never), useWrite: () => { const mutation = mappings.useDeleteStudioV2MappingMutation(); return () => mutation.mutateAsync('M'); } },
];
it.each(operations)('$name invalidates observed sources on ACK without awaiting independent summary refresh', async ({ api, useWrite }) => {
  const { source, wrapper, read } = setup(); api();
  await waitFor(() => expect(source.result.current.fresh).toBe(true));
  read.mockReturnValue(new Promise(() => {}));
  const write = renderHook(useWrite, { wrapper });
  await act(async () => { await write.result.current(); });
  await waitFor(() => expect(source.result.current.fresh).toBe(false));
  expect(source.result.current.status).toBe('refreshing'); expect(source.result.current.observed).toBeDefined(); expect(read).toHaveBeenCalledTimes(2);
});
it('a rejected config mutation does not invalidate successful observation', async () => {
  const { source, wrapper, read } = setup(); vi.spyOn(studioV2MappingsAPI, 'remove').mockRejectedValue(new Error('failed'));
  await waitFor(() => expect(source.result.current.fresh).toBe(true));
  const write = renderHook(mappings.useDeleteStudioV2MappingMutation, { wrapper });
  await act(async () => { await expect(write.result.current.mutateAsync('M')).rejects.toThrow('failed'); });
  expect(source.result.current.fresh).toBe(true); expect(read).toHaveBeenCalledTimes(1);
});
it('late ACK invalidates captured workspace/backend rather than newly selected context', async () => {
  const { source, wrapper, client, read } = setup(); let finish!: () => void;
  vi.spyOn(studioV2MappingsAPI, 'remove').mockImplementation(() => new Promise(resolve => { finish = () => resolve({} as never); }));
  await waitFor(() => expect(source.result.current.fresh).toBe(true)); const old = recordingHomeSourceContext();
  const write = renderHook(mappings.useDeleteStudioV2MappingMutation, { wrapper }); let pending!: Promise<unknown>;
  act(() => { pending = write.result.current.mutateAsync('M'); }); await waitFor(() => expect(finish).toBeTypeOf('function'));
  source.unmount(); studioV2DatalinkApi.defaults.baseURL = 'http://other/api/v1/datalink'; client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { id: 'B' });
  const current = renderHook(() => useRecordingHomeSources({ workspaceId: 'B', workspaceRevision: '', epoch: 2 }), { wrapper });
  await waitFor(() => expect(current.result.current.fresh).toBe(true));
  const invalidate = vi.spyOn(client, 'invalidateQueries');
  await act(async () => { finish(); await pending; });
  expect(invalidate).toHaveBeenCalledWith({ queryKey: recordingHomeSourcePrefix(old, 'A'), refetchType: 'none' });
  expect(current.result.current.fresh).toBe(true); expect(read).toHaveBeenCalledTimes(2);
});
it('source invalidation does not remove the rule owner existing awaited refresh barrier', async () => {
  const { source, wrapper, read } = setup(); let finish!: () => void;
  vi.spyOn(studioV2RulesAPI, 'update').mockResolvedValue({} as never);
  const list = vi.spyOn(studioV2RulesAPI, 'list').mockResolvedValueOnce([]).mockImplementation(() => new Promise(resolve => { finish = () => resolve([]); }));
  const ownerRead = renderHook(() => rules.useStudioV2RulesQuery(true), { wrapper });
  await waitFor(() => expect(ownerRead.result.current.isSuccess).toBe(true));
  await waitFor(() => expect(source.result.current.fresh).toBe(true));
  read.mockReturnValue(new Promise(() => {}));
  const mutation = renderHook(rules.useUpdateStudioV2RuleMutation, { wrapper }); let done = false; let pending!: Promise<unknown>;
  act(() => { pending = mutation.result.current.mutateAsync({ ruleId: 'R', request: {} }).then(result => { done = true; return result; }); });
  await waitFor(() => expect(list).toHaveBeenCalledTimes(2));
  await waitFor(() => expect(source.result.current.fresh).toBe(false));
  expect(done).toBe(false); expect(read).toHaveBeenCalledTimes(2);
  await act(async () => { finish(); await pending; }); expect(done).toBe(true);
});
it('a mounted home revalidates once after ACK without holding the save for that summary', async () => {
  const { source, wrapper, read } = setup();
  await waitFor(() => expect(source.result.current.fresh).toBe(true));
  let finish!: () => void;
  read.mockImplementation(() => new Promise(resolve => { finish = () => resolve({ ...summary('A'), source_digest: 'b'.repeat(64) }); }));
  vi.spyOn(studioV2MappingsAPI, 'remove').mockResolvedValue({} as never);
  const mutation = renderHook(mappings.useDeleteStudioV2MappingMutation, { wrapper });
  await act(async () => { await mutation.result.current.mutateAsync('M'); });
  await waitFor(() => expect(read).toHaveBeenCalledTimes(2));
  expect(source.result.current.fresh).toBe(false); expect(source.result.current.status).toBe('refreshing');
  await act(async () => { finish(); });
  await waitFor(() => expect(source.result.current.observed?.source_digest).toBe('b'.repeat(64)));
  expect(source.result.current.fresh).toBe(true); expect(read).toHaveBeenCalledTimes(2);
});
it('an explicit locator owner wins over a different cached bootstrap workspace', async () => {
  const { source, wrapper, client, read } = setup();
  await waitFor(() => expect(source.result.current.fresh).toBe(true));
  client.setQueryData(studioV2WorkspaceKeys.bootstrap(), { id: 'B' });
  read.mockReturnValue(new Promise(() => {}));
  vi.spyOn(studioV2RulesAPI, 'update').mockResolvedValue({} as never);
  const mutation = renderHook(() => rules.useUpdateStudioV2RuleMutation('A'), { wrapper });
  await act(async () => { await mutation.result.current.mutateAsync({ ruleId: 'R', request: {} }); });
  await waitFor(() => expect(source.result.current.fresh).toBe(false));
  expect(read).toHaveBeenCalledTimes(2);
});
