import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { useStudioLocatorEntry } from '@/hooks/datalink/useStudioLocatorEntry';
import { useStudioV2WorkspaceQuery } from '@/hooks/datalink/useStudioV2Workspace';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import type { EntityLocator } from '@/types/runtimeRuleEvidence';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
const locator: EntityLocator = { step: 1, workspace_id: 'W', device_id: 'D', rule_id: null, mapping_id: null, group_id: null, group_revision: null };
const at = '2026-10-29T00:00:00Z';
function response(device = 'D') {
  return { data: { success: true, data: { id: 'W', kind: 'single', status: 'ready', ordered_device_ids: ['D', 'E'], created_at: at, updated_at: at,
    locator_entry: { workspace_revision: 'W1', hydration_scope: { complete_workspace: false, device_ids: [device], rule_ids: [], mapping_ids: [], group_ids: [] },
      devices: [{ id: device, name: device, description: '', protocol: 'modbus_tcp', status: 'active', connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: at, updated_at: at }],
      rule_heads: {}, group_heads: {}, revalidated_mapping_ids: [] } } } };
}
const clients: QueryClient[] = [];
function harness() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return { client, wrapper };
}
function deferred() { let resolve!: (value: ReturnType<typeof response>) => void; const promise = new Promise<ReturnType<typeof response>>(done => { resolve = done; }); return { promise, resolve }; }
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.resetAllMocks(); });
describe('isolated locator query and ordinary bootstrap guard', () => {
  it('aborts A, clears prior evidence, and ignores delayed A after B selection', async () => {
    const a = deferred(); const b = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(a.promise).mockReturnValueOnce(b.promise);
    const { client, wrapper } = harness(); const canonical = { id: 'canonical', sentinel: true }; client.setQueryData(studioV2WorkspaceKeys.bootstrap(), canonical);
    const view = renderHook(({ value }) => useStudioLocatorEntry(value), { wrapper, initialProps: { value: locator as EntityLocator | null } });
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1)); const signal = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal;
    view.rerender({ value: { ...locator, device_id: 'E' } }); expect(view.result.current.data).toBeUndefined();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2)); expect(signal?.aborted).toBe(true);
    await act(async () => b.resolve(response('E'))); await waitFor(() => expect(view.result.current.data).toMatchObject({ ok: true, value: { locator: { device_id: 'E' } } }));
    await act(async () => a.resolve(response())); expect(view.result.current.data).toMatchObject({ ok: true, value: { locator: { device_id: 'E' } } });
    expect(client.getQueryData(studioV2WorkspaceKeys.bootstrap())).toEqual(canonical); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    view.rerender({ value: null }); expect(view.result.current.data).toBeUndefined();
  });
  it.each([null, { ...locator, device_id: 'http://private.local' }])('does not query invalid or missing intent %j, even on refetch', async (value) => {
    const { wrapper } = harness(); const view = renderHook(() => useStudioLocatorEntry(value), { wrapper });
    expect(view.result.current.data).toBeUndefined();
    await act(async () => { const refreshed = await view.result.current.refetch(); expect(refreshed.data).toEqual({ ok: false, reason: 'unavailable' }); }); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
    await waitFor(() => expect(view.result.current.data).toEqual({ ok: false, reason: 'unavailable' }));
  });
  it('keeps ordinary bootstrap enabled by default and disables it only for explicit entry callers', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ id: 'W' } as never); const { wrapper } = harness();
    const blocked = renderHook(() => useStudioV2WorkspaceQuery(false), { wrapper }); expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled(); blocked.unmount();
    renderHook(() => useStudioV2WorkspaceQuery(), { wrapper }); await waitFor(() => expect(studioV2WorkspaceAPI.get).toHaveBeenCalledTimes(1));
  });
});

import { useStudioV2DatabaseMetadataQuery } from '@/hooks/datalink/useStudioV2WorkspaceDatabase';
import { useWriteGroupTestWriteOperationQuery } from '@/hooks/datalink/useStudioV2WriteGroupTestWrite';
import { useRecordingTestWriteRecovery } from '@/hooks/datalink/useRecordingTestWriteRecovery';
const nativeReadScope = { workspaceId: 'W', workspaceRevision: 'W1' };
it('keys native metadata reads by complete verified workspace observation without reusing prior columns', async () => {
  const metadata = { workspace_id: 'W', connector_id: 'C', connector_revision: 'C1', database: 'native', schema: '', table: 'meter', inspection_status: 'missing', columns: [] };
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: metadata } }); const { wrapper } = harness();
  const view = renderHook(({ selected }) => useStudioV2DatabaseMetadataQuery({ connectorId: 'C', connectorRevision: 'C1', database: 'native', schema: '', table: 'meter' }, true, selected), { wrapper, initialProps: { selected: nativeReadScope } });
  await waitFor(() => expect(view.result.current.data?.workspace_id).toBe('W'));
  view.rerender({ selected: { ...nativeReadScope, workspaceRevision: 'W2' } }); expect(view.result.current.data).toBeUndefined();
  await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
  expect(vi.mocked(studioV2DatalinkApi.get).mock.calls[1][1]?.params).toMatchObject({ locator_workspace_id: 'W', locator_workspace_revision: 'W2' });
});
it('forwards native operation and recovery scopes without automatic operation re-execution', async () => {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: { workspace_id: 'W', complete: true, operations: [] } } });
  const { wrapper } = harness(); const list = renderHook(() => useRecordingTestWriteRecovery('W', nativeReadScope), { wrapper });
  await waitFor(() => expect(list.result.current.available).toBe(true));
  expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/test-writes', { params: { locator_workspace_id: 'W', locator_workspace_revision: 'W1' } });
  list.unmount(); vi.mocked(studioV2DatalinkApi.get).mockClear();
  const operation = { operation_id: 'O', workspace_id: 'W', action: 'test_write', status: 'running', created_at: at, updated_at: at };
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: operation } });
  const view = renderHook(() => useWriteGroupTestWriteOperationQuery('O', 'W', nativeReadScope), { wrapper });
  expect(studioV2DatalinkApi.get).not.toHaveBeenCalled(); await act(async () => { await view.result.current.refetch(); });
  expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/O', { params: { locator_workspace_id: 'W', locator_workspace_revision: 'W1' } });
});

import { useStudioV2WorkspaceMeasurementsQuery } from '@/hooks/datalink/useStudioV2WorkspaceMeasurements';
it('does not repurpose the existing cache-only measurement scope as locator read authority', async () => {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: [] } }); const { wrapper } = harness();
  const view = renderHook(({ explicit }) => useStudioV2WorkspaceMeasurementsQuery(true, 'D', nativeReadScope, explicit), { wrapper, initialProps: { explicit: undefined as typeof nativeReadScope | undefined } });
  await waitFor(() => expect(view.result.current.data).toEqual([]));
  expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/measurements', { params: { device_id: 'D' } });
  view.rerender({ explicit: nativeReadScope }); expect(view.result.current.data).toBeUndefined();
  await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
  expect(vi.mocked(studioV2DatalinkApi.get).mock.calls[1][1]?.params).toMatchObject({ locator_workspace_id: 'W', locator_workspace_revision: 'W1' });
});

it('captures native read scope for refetch rather than following a mutated caller under the old cache key', async () => {
  vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: { success: true, data: { workspace_id: 'W', connector_id: 'C', connector_revision: 'C1', database: 'native', schema: '', table: 'meter', inspection_status: 'missing', columns: [] } } });
  const selected = { ...nativeReadScope }; const { wrapper } = harness();
  const view = renderHook(() => useStudioV2DatabaseMetadataQuery({ connectorId: 'C', connectorRevision: 'C1', database: 'native', schema: '', table: 'meter' }, true, selected), { wrapper });
  await waitFor(() => expect(view.result.current.data?.workspace_id).toBe('W')); selected.workspaceRevision = 'MUTATED';
  await act(async () => { await view.result.current.refetch(); });
  expect(vi.mocked(studioV2DatalinkApi.get).mock.calls[1][1]?.params).toMatchObject({ locator_workspace_id: 'W', locator_workspace_revision: 'W1' });
});
