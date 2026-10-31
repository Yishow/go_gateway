import { queryClient } from '../../../src/lib/queryClient';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { PropsWithChildren } from 'react';
import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceDevicesAPI } from './device-autosave-page.testHarness';
import { useStudioV2AutosaveState } from '../../../src/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import type { DeviceSourceRequest } from '../../../src/pages/datalink/workbench-v2/deviceSourceSubmission';
const buffer = { name: 'New PLC', description: '', protocol: 'modbus_tcp' as const, config: { host: '192.0.2.20', port: 502, slave_id: 1, timeout: 5 } };
const request = (id = 'request-1'): DeviceSourceRequest => ({ workspaceId: 'workspace-1', sessionId: 'modal-1', requestId: id, generation: 1, buffer: structuredClone(buffer) });
function setup() {
  const client = new QueryClient({ defaultOptions: queryClient.getDefaultOptions() });
  const wrapper = ({ children }: PropsWithChildren) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return renderHook(({ scope }) => useStudioV2AutosaveState(true, undefined, false, undefined, scope), { wrapper, initialProps: { scope: 'workspace-1' } });
}
function deferred<T>() { let resolve!: (value: T) => void; let reject!: (reason: unknown) => void; const promise = new Promise<T>((yes, no) => { resolve = yes; reject = no; }); return { promise, resolve, reject }; }

beforeEach(() => { resetDeviceAutosaveMocks(); vi.mocked(studioV2WorkspaceDevicesAPI.create).mockReset(); });
describe('recording editor uses the actual device autosave owner', () => {
  it('does not submit local edits until explicit submit, then confirms only returned identity and readback', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    const input = request(); input.buffer.name = 'Edited local';
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled();
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: String(payload.name), connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    let outcome; await act(async () => { outcome = await current.result.current.deviceSource.submit(input); });
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledTimes(1);
    expect(outcome).toMatchObject({ status: 'confirmed', request: { requestId: 'request-1', buffer: { name: 'Edited local' } }, device: { name: 'Edited local' } });
    expect(current.result.current.state.devices[0]).toMatchObject({ name: 'Edited local', persisted: true });
    expect(studioV2WorkspaceDevicesAPI.remove).not.toHaveBeenCalled();
  });
  it('retains unknown after a lost response; same-name readback and new request IDs cannot authorize replacement create', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockRejectedValue(new Error('lost ack'));
    let outcome; await act(async () => { outcome = await current.result.current.deviceSource.submit(request()); });
    expect(outcome).toMatchObject({ status: 'unknown' });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([deviceFixture({ name: 'New PLC' })]);
    await act(async () => {
      expect(await current.result.current.deviceSource.inspect('request-1')).toMatchObject({ status: 'unknown' });
      expect(await current.result.current.deviceSource.submit(request('replacement'))).toMatchObject({ status: 'blocked' });
    });
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledTimes(1);
  });
  it('keeps an immutable submitted snapshot and de-duplicates the actual request', async () => {
    const ack = deferred<ReturnType<typeof deviceFixture>>();
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(() => ack.promise);
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    const input = request(); let first!: ReturnType<typeof current.result.current.deviceSource.submit>;
    act(() => { first = current.result.current.deviceSource.submit(input); });
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledTimes(1));
    const sent = vi.mocked(studioV2WorkspaceDevicesAPI.create).mock.calls[0][0];
    input.buffer.name = 'Newer unsent name'; input.buffer.config.host = 'invalid raw !!';
    const saved = deviceFixture({ id: sent.id, name: 'New PLC', connection_config: JSON.stringify(buffer.config) });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]);
    await act(async () => { ack.resolve(saved); const outcome = await first; expect(outcome).toMatchObject({ status: 'confirmed', request: { buffer } }); });
    await act(async () => { await current.result.current.deviceSource.submit(request()); });
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledTimes(1);
  });
  it('does not confirm or populate a newer workspace after a late response', async () => {
    const ack = deferred<ReturnType<typeof deviceFixture>>(); vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(() => ack.promise);
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    let pending!: ReturnType<typeof current.result.current.deviceSource.submit>;
    act(() => { pending = current.result.current.deviceSource.submit(request()); });
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce());
    current.rerender({ scope: 'workspace-2' });
    await act(async () => { ack.resolve(deviceFixture({ name: 'New PLC' })); expect(await pending).toMatchObject({ status: 'unknown' }); });
    expect(current.result.current.state.devices.some(device => device.persisted)).toBe(false);
  });
});


describe('device submission readback and existing queue boundaries', () => {
  it('can inspect only an acknowledged exact identity after readback fails, publishing recovered saved metadata', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    let saved!: ReturnType<typeof deviceFixture>;
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config) }); return saved;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'unknown' }); });
    expect(current.result.current.deviceSource.observations()).toHaveLength(1);
    expect(current.result.current.state.devices[0].save_state).toBe('save-error');
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]);
    await act(async () => { expect(await current.result.current.deviceSource.inspect('request-1')).toMatchObject({ status: 'confirmed', device: saved }); });
    expect(current.result.current.state.devices[0]).toMatchObject({ persisted: true, save_state: 'saved' });
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'confirmed' }); });
  });
  it('keeps runtime application failure separate from confirmed persistence', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config), runtime_apply_status: 'apply_failed' });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'persisted-unavailable' }); });
    expect(current.result.current.state.devices[0]).toMatchObject({ persisted: true, save_state: 'save-error' });
  });
  it('uses the same entity queue for a legacy edit made during the submitted save', async () => {
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([deviceFixture()]);
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    const ack = deferred<ReturnType<typeof deviceFixture>>();
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockImplementationOnce(() => ack.promise).mockImplementation(async (id, payload) => {
      const saved = deviceFixture({ id, name: payload.name, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    let first!: Promise<unknown>;
    act(() => { first = current.result.current.deviceSource.submit({ ...request(), deviceId: 'dev-A' }); });
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledTimes(1));
    act(() => { current.result.current.actions.dispatch({ type: 'renameDevice', deviceId: 'dev-A', name: 'Newer legacy edit' }); });
    expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledTimes(1);
    const saved = deviceFixture({ name: buffer.name, connection_config: JSON.stringify(buffer.config) });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]);
    await act(async () => { ack.resolve(saved); expect(await first).toMatchObject({ status: 'confirmed', request: { buffer: { name: 'New PLC' } } }); });
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(current.result.current.state.devices[0]).toMatchObject({ name: 'Newer legacy edit', save_state: 'saved' }));
  });
  it('blocks foreign-scope or malformed input and retains unresolved submissions without deleting them', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    await act(async () => {
      expect(await current.result.current.deviceSource.submit({ ...request(), workspaceId: 'other' })).toMatchObject({ status: 'blocked' });
      expect(await current.result.current.deviceSource.submit({ ...request(), buffer: { ...buffer, name: '' } })).toMatchObject({ status: 'blocked' });
    });
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled();
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockRejectedValue(new Error('response lost'));
    await act(async () => { await current.result.current.deviceSource.submit(request()); });
    const localId = current.result.current.state.devices[0].id;
    act(() => { current.result.current.actions.dispatch({ type: 'removeDevice', deviceId: localId }); });
    expect(current.result.current.state.devices).toHaveLength(1);
    expect(current.result.current.deviceSource.observations()).toEqual([expect.objectContaining({ status: 'unknown' })]);
    expect(studioV2WorkspaceDevicesAPI.remove).not.toHaveBeenCalled();
  });
});

describe('acknowledged recovery preserves queued legacy intent', () => {
  it('records a legacy edit while readback is unknown and resumes it only after exact acknowledgement recovery', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    let saved!: ReturnType<typeof deviceFixture>;
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config) }); return saved;
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockReset().mockImplementation(async (id, payload) => {
      const next = deviceFixture({ id, name: payload.name, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([next]); return next;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'unknown' }); });
    const id = current.result.current.state.devices[0].id;
    act(() => { current.result.current.actions.dispatch({ type: 'renameDevice', deviceId: id, name: 'Legacy edit after unknown' }); });
    expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled();
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]);
    await act(async () => { await current.result.current.deviceSource.inspect('request-1'); });
    await waitFor(() => expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledOnce());
    await waitFor(() => expect(current.result.current.state.devices[0]).toMatchObject({ name: 'Legacy edit after unknown', save_state: 'saved' }));
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
  });
});

describe('runtime application repair remains an update of the acknowledged device', () => {
  it('removes superseded recovery attention without rewriting the old failed-apply observation', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const saved = deviceFixture({ id: payload.id, name: payload.name, connection_config: JSON.stringify(payload.connection_config), runtime_apply_status: 'apply_failed' });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'persisted-unavailable' }); });
    const id = current.result.current.state.devices[0].id;
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockReset().mockImplementation(async (deviceId, payload) => {
      const saved = deviceFixture({ id: deviceId, name: payload.name, connection_config: JSON.stringify(payload.connection_config), runtime_apply_status: 'not_running' });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([saved]); return saved;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit({ ...request('repair-2'), deviceId: id, generation: 2 })).toMatchObject({ status: 'confirmed' }); });
    expect(current.result.current.deviceSource.observations()).toEqual([expect.objectContaining({ status: 'confirmed', request: expect.objectContaining({ requestId: 'repair-2' }) })]);
    expect(await current.result.current.deviceSource.inspect('request-1')).toMatchObject({ status: 'persisted-unavailable' });
    expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledOnce();
    expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledOnce();
  });
});

describe('device identity acknowledgement', () => {
  it('does not mark a requested client identity persisted when the mutation returns a different identity', async () => {
    const current = setup(); await waitFor(() => expect(current.result.current.workspaceHydrated).toBe(true));
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockImplementation(async payload => {
      const different = deviceFixture({ id: 'not-the-requested-device', name: payload.name, connection_config: JSON.stringify(payload.connection_config) });
      vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([different]); return different;
    });
    await act(async () => { expect(await current.result.current.deviceSource.submit(request())).toMatchObject({ status: 'unknown' }); });
    expect(current.result.current.state.devices[0]).toMatchObject({ persisted: false, save_state: 'save-error' });
    expect(current.result.current.deviceSource.observations()[0].device).toBeUndefined();
  });
});
