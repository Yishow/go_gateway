import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { PropsWithChildren } from 'react';
import { queryClient } from '../../../src/lib/queryClient';
import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceDevicesAPI } from './device-autosave-page.testHarness';
import { useStudioV2AutosaveState } from '../../../src/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { useRecordingDeviceBuffer } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer';
import type { DeviceSourceRequest } from '../../../src/pages/datalink/workbench-v2/deviceSourceSubmission';
const request = (id = 'one'): DeviceSourceRequest => ({ workspaceId: 'workspace-1', sessionId: 'form-one', requestId: id, generation: 1,
  buffer: { name: 'Fresh PLC', description: '', protocol: 'modbus_tcp', config: { host: '192.0.2.25', port: 502, slave_id: 2, timeout: 5 } } });
function setup() {
  const client = new QueryClient({ defaultOptions: queryClient.getDefaultOptions() });
  const wrapper = ({ children }: PropsWithChildren) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return renderHook(({ scope }) => useStudioV2AutosaveState(true, undefined, false, undefined, scope), { wrapper, initialProps: { scope: 'workspace-1' } });
}
function deferred<T>() { let resolve!: (value: T) => void; const promise = new Promise<T>(yes => { resolve = yes; }); return { promise, resolve }; }
beforeEach(() => { resetDeviceAutosaveMocks(); vi.mocked(studioV2WorkspaceDevicesAPI.create).mockReset(); });
describe('page-session recent connections consume the existing source owner', () => {
  it('offers only a confirmed new-device request snapshot through the page port', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'confirmed' }); });
    expect(current.result.current.deviceSource.recentConnections?.read('workspace-1')).toEqual({ protocol: 'modbus_tcp', config: request().buffer.config });
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
  });
  it('clearing invalidates an earlier pending request without cancelling its source mutation', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    const ack = deferred<ReturnType<typeof deviceFixture>>(); vi.mocked(studioV2WorkspaceDevicesAPI.create).mockReturnValue(ack.promise);
    let pending!: ReturnType<typeof current.result.current.deviceSource.submit>;
    act(() => { pending = current.result.current.deviceSource.submit(request()); });
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce());
    expect(current.result.current.deviceSource.recentConnections).toBeDefined();
    current.result.current.deviceSource.recentConnections?.clear('workspace-1', 'modbus_tcp');
    expect(current.result.current.deviceSource.observations()[0].status).toBe('pending');
    const sent = vi.mocked(studioV2WorkspaceDevicesAPI.create).mock.calls[0][0];
    const saved = deviceFixture({ id: sent.id, name: sent.name, connection_config: JSON.stringify(sent.connection_config) });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]);
    await act(async () => { ack.resolve(saved); expect(await pending).toMatchObject({ status: 'confirmed' }); });
    expect(current.result.current.deviceSource.recentConnections?.read('workspace-1')).toBeUndefined();
  });
});


describe('recent cache lifetime follows the page, not a modal or editor', () => {
  it('survives form unmount, ignores cancelled edits, and is emptied by page unmount', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    await act(async () => { await current.result.current.deviceSource.submit(request()); });
    const port = current.result.current.deviceSource;
    const first = renderHook(() => useRecordingDeviceBuffer('workspace-1', port));
    await act(async () => { await first.result.current.start(); });
    act(() => first.result.current.change({ ...first.result.current.buffer, name: 'Cancelled name', config: { host: 'cancelled endpoint' } }));
    first.unmount();
    const next = renderHook(() => useRecordingDeviceBuffer('workspace-1', port));
    await act(async () => { await next.result.current.start(); });
    expect(next.result.current.buffer.name).toBe('');
    expect(next.result.current.buffer.config.host).toBe('192.0.2.25');
    next.unmount(); current.unmount();
    expect(port.recentConnections?.read('workspace-1')).toBeUndefined();
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
  });
  it('empties on workspace change and never resurrects the previous scope', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    await act(async () => { await current.result.current.deviceSource.submit(request()); });
    expect(current.result.current.deviceSource.recentConnections?.read('workspace-1')).toBeDefined();
    current.rerender({ scope: 'workspace-2' });
    expect(current.result.current.deviceSource.recentConnections?.read('workspace-1')).toBeUndefined();
    current.rerender({ scope: 'workspace-1' });
    expect(current.result.current.deviceSource.recentConnections?.read('workspace-1')).toBeUndefined();
  });
  it('does not populate from failed runtime application even with exact persisted readback', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config), runtime_apply_status: 'apply_failed' });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'persisted-unavailable' }); });
    expect(current.result.current.deviceSource.recentConnections?.read('workspace-1')).toBeUndefined();
  });
});
